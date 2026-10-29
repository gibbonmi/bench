// Stack red tests: a candidate red bisects to one member, publishes the green prefix, and
// refuses that member with stack-red. An inherited red hands back after one grade.
package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// redMemberStep is the gate line of a red stack fixture: a tree in which the file of a
// member reads red grades red.
func redMemberStep(*testrepo.GateFixture, string) string {
	return "for file in owned.txt member-*.txt; do if [ -e \"$file\" ] && IFS= read -r bytes < \"$file\" && [ \"$bytes\" = red ]; then exit 1; fi; done\n"
}

// redStackFixture is the stack fixture of n members whose destination holds a project
// green, with member red (0-based) turned red. The green makes a red top tree a candidate,
// so the stack bisects. The tally holds no grade before the stack runs.
func redStackFixture(t *testing.T, prefix string, n, red int) stackFixture {
	t.Helper()
	f := stackLandingFixture(t, prefix, n, redMemberStep)
	greenDestination(t, &f)
	turnMemberRed(t, &f, red)
	return f
}

// greenDestination records a project green of the destination of f. The ordinary gate of
// the fixture asks for the file that member 1 adds, so the destination first drops that
// line, and each member takes the new destination as its base.
func greenDestination(t *testing.T, f *stackFixture) {
	t.Helper()
	script := filepath.Join(f.root, ".bench", "gate.sh")
	body, err := os.ReadFile(script)
	mustNoError(t, err)
	const memberLine = "[ -f owned.txt ]\n"
	if !strings.Contains(string(body), memberLine) {
		t.Fatalf("destination gate %q holds no line %q", body, memberLine)
	}
	mustWrite(t, script, []byte(strings.Replace(string(body), memberLine, "", 1)), 0o755)
	gitRun(t, f.root, "add", ".bench/gate.sh")
	gitRun(t, f.root, withIdentity("commit", "-qm", "the destination gate grades the destination alone")...)
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	for i := range f.members {
		gitRun(t, f.members[i].creation.Path, "rebase", "-q", "main")
		f.members[i].tip = gitOutput(t, f.members[i].creation.Path, "rev-parse", "HEAD")
	}
	f.tip = f.members[0].tip
	// The fixture's own marker names the old base, and the advance creates the marker anew.
	gitRun(t, f.root, "update-ref", "-d", "refs/bench/green/main")
	markTargetGreen(t, Creation{Path: f.root, Assignment: intent.Assignment{Branch: "refs/heads/main"}})
	mustRemove(t, f.tally)
}

// turnMemberRed commits the red bytes to the file of member i and moves its tip there.
func turnMemberRed(t *testing.T, f *stackFixture, i int) {
	t.Helper()
	m := &f.members[i]
	commitInWorktree(t, m.creation.Path, m.file, "red\n", "member goes red")
	m.tip = gitOutput(t, m.creation.Path, "rev-parse", "HEAD")
}

// TestLandStackRedJourney is LS46, LS47, LS48, and LS49. With three members and member 2
// red, the bisect grades 1 + ⌈log2 3⌉ = 3 trees, publishes member 1, leaves member 3
// active, and refuses member 2 with stack-red at its own label.
func TestLandStackRedJourney(t *testing.T) {
	t.Parallel()
	f := redStackFixture(t, "stack-red", 3, 1)
	file, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 1 {
		t.Fatalf("red stack = (%d, %q, %q), want exit 1", r.exit, r.stdout, r.stderr)
	}
	top := gitOutput(t, f.root, "rev-parse", "main")
	if first, parent := gitOutput(t, f.root, "rev-parse", "main^1"), gitOutput(t, f.root, "rev-parse", "main^2"); first != f.base || parent != f.members[0].tip {
		t.Fatalf("destination %s has parents %s and %s, want the merge of member 1 onto %s", top, first, parent, f.base)
	}
	requireTally(t, f.tally, "ggg", "after the red stack")
	if assignmentActive(t, f.root, f.members[0].creation.Assignment.ID) || !assignmentActive(t, f.root, f.members[2].creation.Assignment.ID) {
		t.Fatalf("stack = %q, want member 1 released and member 3 active", r.stdout)
	}
	var rows []stackMemberRow
	for i, m := range f.members {
		rows = append(rows, stackMemberRow{member: i + 1, label: m.creation.Assignment.Label, sourceTip: m.tip, published: "none", outcome: []string{"landed", "red", "unlanded"}[i]})
	}
	rows[0].published = top
	table, err := stackMembersTable(rows)
	mustNoError(t, err)
	if !strings.Contains(r.stdout, table) {
		t.Fatalf("stack stdout = %q, want the members table %q", r.stdout, table)
	}
	m := f.members[1]
	rerun := "bench worktree land --stack '" + file + "' --source-tip '" + m.request + "'=<repaired-tip>"
	if want, next := landingRoute(faceStackRed, rerun, labelOf(m.creation)), stackRefusalNext(t, r.stdout, "2", m); next != want {
		t.Fatalf("stack-red next = %q, want %q", next, want)
	}
	if want := "stack{destination_base=" + f.base + ",published_commit=" + top + ",landed=1,"; !strings.Contains(r.stdout, want) {
		t.Fatalf("stack stdout = %q, want the record %q", r.stdout, want)
	}
}

// TestLandStackInheritedRedHandsBack is LS106. A destination with no green makes a red top
// tree inherited, so no member carries the red: the stack hands back after one grade and
// publishes nothing.
func TestLandStackInheritedRedHandsBack(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-inherited", 2, redMemberStep)
	gitRun(t, f.root, "update-ref", "-d", "refs/bench/green/main")
	turnMemberRed(t, &f, 1)
	file, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	want := landingRoute(faceLandHandback, "bench worktree land --stack '"+file+"'", nil)
	if next, printed := landingFaceNext(r.stdout, ""); r.exit != 1 || !printed || next != want {
		t.Fatalf("inherited stack = (%d, %q, %q), want exit 1 and next %q", r.exit, r.stdout, r.stderr, want)
	}
	requireOneRun(t, f.tally, "after the inherited stack")
	if top := gitOutput(t, f.root, "rev-parse", "main"); top != f.base {
		t.Fatalf("destination = %s, want the unpublished base %s", top, f.base)
	}
	requireNoLease(t, f.root)
}

// TestStackAttributedReadsOnlyACandidateRed: a refusal after the bisect leaves the
// attribution set, so only a candidate red names the attributed member.
func TestStackAttributedReadsOnlyACandidateRed(t *testing.T) {
	t.Parallel()
	result := landing.StackResult{Attributed: 2}
	for _, c := range []struct {
		name string
		err  error
		want int
	}{
		{"candidate red", fmt.Errorf("grade: %w", landing.AuthorizationRefusal{Result: authorization.Result{Kind: authorization.Candidate}}), 1},
		{"inherited red", landing.AuthorizationRefusal{Result: authorization.Result{Kind: authorization.Inherited}}, -1},
		{"admission refusal", errors.New("the broker refused the admission"), -1},
	} {
		if got := stackAttributed(c.err, result); got != c.want {
			t.Errorf("%s: stackAttributed = %d, want %d", c.name, got, c.want)
		}
	}
}
