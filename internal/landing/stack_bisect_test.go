package landing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/gittest"
)

const (
	green = authorization.Green
	red   = authorization.Candidate
)

// bisectRun is one stack whose prefix p grades verdicts[p-1]: the result, the refusal, each
// graded prefix in order, and the stack members of each grade.
type bisectRun struct {
	f         stackLanding
	got       StackResult
	err       error
	prefixes  []int
	members   [][]gate.StackMember
	record    *stackRecord
	admission *stackAdmissionRecord
}

// landBisectStack lands a stack of len(verdicts) members. The stub reads the prefix of each
// graded tree from the member files that the tree holds, so the table does not depend on
// the trees the owner writes. A non-nil during runs inside the grade of each prefix p.
func landBisectStack(t *testing.T, during func(f stackLanding, p int), verdicts ...authorization.Kind) bisectRun {
	t.Helper()
	run := bisectRun{f: newStackLanding(t, len(verdicts), nil), admission: &stackAdmissionRecord{}}
	o, record := stackOwner(green)
	recorded := o.authorize
	o.authorize = func(ctx context.Context, root, tree string, stdout, stderr io.Writer) authorization.Result {
		names := treeNames(t, run.f.root, tree)
		p := 0
		for _, name := range names {
			if strings.HasPrefix(name, "member-") {
				p++
			}
		}
		for i := 1; i <= p; i++ {
			if !slices.Contains(names, fmt.Sprintf("member-%d", i)) {
				t.Fatalf("graded tree %s is not a prefix: it lacks member-%d", tree, i)
			}
		}
		run.prefixes = append(run.prefixes, p)
		run.members = append(run.members, gate.StackMembers(ctx))
		if during != nil {
			record.duringGrade = func() { during(run.f, p) }
		}
		recorded(ctx, root, tree, stdout, stderr)
		return authorization.Result{Kind: verdicts[p-1]}
	}
	run.record = record
	run.got, run.err = o.LandStack(context.Background(), run.f.request(t), run.admission)
	return run
}

// assertPublishedPrefix requires the red refusal, the attributed member, and exactly the
// first k merges published by one ref update to merge k.
func (run bisectRun) assertPublishedPrefix(t *testing.T, attributed, k int) {
	t.Helper()
	var refusal AuthorizationRefusal
	if !errors.As(run.err, &refusal) || refusal.Result.Kind != red || run.got.Attributed != attributed {
		t.Fatalf("bisected stack = %+v, %v, want a red refusal attributed to member %d", run.got, run.err, attributed)
	}
	if run.got.Published() != k || run.admission.published != k || !slices.Equal(run.admission.trees, run.admission.checked[:k]) {
		t.Fatalf("published %d members, admission decided %d over %q, want %d", run.got.Published(), run.admission.published, run.admission.trees, k)
	}
	merge := run.got.Merges[k-1]
	if want := []string{"refs/heads/main " + run.f.base + " " + merge}; !slices.Equal(run.record.updates, want) {
		t.Fatalf("ref updates = %q, want %q", run.record.updates, want)
	}
	if parents := gittest.Output(t, run.f.root, "rev-list", "--parents", "-n", "1", merge); !strings.HasSuffix(parents, " "+run.f.members[k-1].tip) {
		t.Fatalf("published merge %s has parents %q, want the merge of member %d", merge, parents, k)
	}
}

// assertStopped requires a stop with kind after exactly the graded prefixes, with no
// attributed member and no publication.
func (run bisectRun) assertStopped(t *testing.T, kind authorization.Kind, prefixes []int) {
	t.Helper()
	var refusal AuthorizationRefusal
	if !errors.As(run.err, &refusal) || refusal.Result.Kind != kind || run.got.Kind != kind || run.got.Attributed != 0 {
		t.Fatalf("stopped stack = %+v, %v, want a %s refusal with no attributed member", run.got, run.err, kind)
	}
	if !slices.Equal(run.prefixes, prefixes) || run.got.Published() != 0 || run.admission.published != 0 {
		t.Fatalf("stopped stack graded prefixes %v and published %d, want %v and none", run.prefixes, run.got.Published(), prefixes)
	}
	assertStackUnpublished(t, run.f, run.record)
}

// TestLandStackGradesEachPrefixWithItsMembers is LS29. Member 3 of four is red, so the
// grades are prefixes 4, 2, and 3, and each grade names the members of its prefix alone.
func TestLandStackGradesEachPrefixWithItsMembers(t *testing.T) {
	run := landBisectStack(t, nil, green, green, red, red)
	if len(run.prefixes) < 2 {
		t.Fatalf("graded prefixes %v, want a bisect", run.prefixes)
	}
	for g, p := range run.prefixes {
		var want []gate.StackMember
		for _, m := range run.f.members[:p] {
			want = append(want, gate.StackMember{Source: m.tip, Folder: m.folder})
		}
		if !slices.Equal(run.members[g], want) {
			t.Fatalf("grade %d of prefix %d carries members %+v, want %+v", g+1, p, run.members[g], want)
		}
	}
}

// TestLandStackBisectsToTheRedMember is LS30, LS31, and LS32. A red member above a green
// prefix must carry the red, and that green prefix alone publishes.
func TestLandStackBisectsToTheRedMember(t *testing.T) {
	run := landBisectStack(t, nil, green, green, red, red)
	if want := []int{4, 2, 3}; !slices.Equal(run.prefixes, want) {
		t.Fatalf("graded prefixes %v, want %v", run.prefixes, want)
	}
	run.assertPublishedPrefix(t, 3, 2)
	t.Run("an odd member count rounds the midpoint down", func(t *testing.T) {
		run := landBisectStack(t, nil, green, green, red)
		if want := []int{3, 1, 2}; !slices.Equal(run.prefixes, want) {
			t.Fatalf("graded prefixes %v, want %v", run.prefixes, want)
		}
		run.assertPublishedPrefix(t, 3, 2)
	})
	moveDuringPrefix3 := func(member int) func(stackLanding, int) {
		return func(f stackLanding, p int) {
			if p == 3 {
				gittest.MustRun(t, f.members[member-1].worktree, "commit", "--allow-empty", "-qm", "moved during the grade")
			}
		}
	}
	t.Run("a moved top published member refuses the publication", func(t *testing.T) {
		run := landBisectStack(t, moveDuringPrefix3(2), green, green, red, red)
		var member StackMemberError
		if !errors.As(run.err, &member) || member.Member != 2 || run.got.Published() != 0 || run.admission.published != 0 {
			t.Fatalf("moved member 2 = %+v, %v, want a refusal of member 2 and nothing published", run.got, run.err)
		}
		assertStackUnpublished(t, run.f, run.record)
	})
	t.Run("a moved unpublished member keeps the green prefix", func(t *testing.T) {
		run := landBisectStack(t, moveDuringPrefix3(3), green, green, red, red)
		run.assertPublishedPrefix(t, 3, 2)
	})
}

// TestLandStackPublishesNothingForAFirstMemberRed is LS33. With no green prefix above the
// destination, no member can publish.
func TestLandStackPublishesNothingForAFirstMemberRed(t *testing.T) {
	run := landBisectStack(t, nil, red, red, red, red)
	if want := []int{4, 2, 1}; !slices.Equal(run.prefixes, want) {
		t.Fatalf("graded prefixes %v, want %v", run.prefixes, want)
	}
	if run.got.Attributed != 1 || run.got.Published() != 0 || run.admission.published != 0 {
		t.Fatalf("first-member red = %+v, want member 1 attributed and nothing published", run.got)
	}
	assertStackUnpublished(t, run.f, run.record)
}

// TestLandStackAttributesARedTransition is LS34. A green prefix above a red one does not
// publish, because the bisect never grades it.
func TestLandStackAttributesARedTransition(t *testing.T) {
	run := landBisectStack(t, nil, green, red, green, red)
	if want := []int{4, 2, 1}; !slices.Equal(run.prefixes, want) {
		t.Fatalf("graded prefixes %v, want %v", run.prefixes, want)
	}
	run.assertPublishedPrefix(t, 2, 1)
}

// TestLandStackBoundsTheBisect is LS35. Each grade is a whole gate, so a long stack must
// not grade each prefix.
func TestLandStackBoundsTheBisect(t *testing.T) {
	verdicts := slices.Repeat([]authorization.Kind{green}, 16)
	verdicts[15] = red
	run := landBisectStack(t, nil, verdicts...)
	if want := []int{16, 8, 12, 14, 15}; !slices.Equal(run.prefixes, want) {
		t.Fatalf("graded prefixes %v, want %v", run.prefixes, want)
	}
	run.assertPublishedPrefix(t, 16, 15)
}

// TestLandStackStopsOnInfrastructureInTheBisect is LS37. An infrastructure grade proves no
// prefix, so the stack stops at it and publishes nothing, not even a green prefix.
func TestLandStackStopsOnInfrastructureInTheBisect(t *testing.T) {
	run := landBisectStack(t, nil, green, authorization.Infrastructure, red, red)
	run.assertStopped(t, authorization.Infrastructure, []int{4, 2})
	t.Run("a fault above a green prefix publishes nothing", func(t *testing.T) {
		run := landBisectStack(t, nil, green, green, authorization.Infrastructure, red)
		run.assertStopped(t, authorization.Infrastructure, []int{4, 2, 3})
	})
}

// TestLandStackStopsOnAnInheritedRed is LS99. An inherited top tree means the destination
// held no green, so no prefix counts as green and the stack stops after one grade.
func TestLandStackStopsOnAnInheritedRed(t *testing.T) {
	run := landBisectStack(t, nil, green, green, green, authorization.Inherited)
	run.assertStopped(t, authorization.Inherited, []int{4})
}
