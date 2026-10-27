// Stack landing tests: a green stack lands every member with one grade, and the stack
// refuses before its lease or under a held lease.
package worktree

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
)

// writeStack writes the stack file of members at the base of f, with tips as each
// member's source_tip, and answers its path.
func (f stackFixture) writeStack(t *testing.T, members []stackFixtureMember, tips ...string) string {
	t.Helper()
	type entry struct {
		Path      string `json:"path"`
		Request   string `json:"request"`
		Base      string `json:"base"`
		SourceTip string `json:"source_tip"`
		Spec      string `json:"spec"`
		Message   string `json:"message"`
	}
	var doc struct {
		Members []entry `json:"members"`
	}
	for i, m := range members {
		doc.Members = append(doc.Members, entry{Path: m.creation.Path, Request: m.request, Base: f.base, SourceTip: tips[i], Spec: m.slug, Message: "land " + m.request})
	}
	data, err := json.Marshal(doc)
	mustNoError(t, err)
	path := filepath.Join(t.TempDir(), "stack.json")
	mustWrite(t, path, data, 0o644)
	return path
}

// stack writes the stack file of every member at its own tip and answers the stack
// call's arguments, followed by overrides as `--source-tip` values.
func (f stackFixture) stack(t *testing.T, overrides ...string) (file string, args []string) {
	t.Helper()
	tips := make([]string, len(f.members))
	for i, m := range f.members {
		tips[i] = m.tip
	}
	file = f.writeStack(t, f.members, tips...)
	args = []string{"--stack", file}
	for _, o := range overrides {
		args = append(args, "--source-tip", o)
	}
	return file, args
}

// requireNoLease fails t unless the ledger at root holds no landing lease.
func requireNoLease(t *testing.T, root string) {
	t.Helper()
	leases, err := commitmenttest.Leases(root)
	mustNoError(t, err)
	if len(leases) != 0 {
		t.Fatalf("leases = %+v, want none", leases)
	}
}

// TestLandStackJourney is LS40, LS41, LS42, LS43, LS44, LS45, and LS101. Two tickets-only
// members land with one command and one grade.
func TestLandStackJourney(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-journey", 2, nil)
	_, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 0 {
		t.Fatalf("stack = (%d, %q, %q), want exit 0", r.exit, r.stdout, r.stderr)
	}
	top := gitOutput(t, f.root, "rev-parse", "main")
	first := gitOutput(t, f.root, "rev-parse", "main^1")
	if tally, err := os.ReadFile(f.tally); err != nil || string(tally) != "g" {
		t.Fatalf("gate tally = %q, %v; want one run", tally, err)
	}
	var wantSources []string
	for _, m := range f.members {
		// The landed-sibling cleanup also retires a landed member, so the release proof
		// reads that the cleanup found nothing of the member's left to retire.
		if assignmentActive(t, f.root, m.creation.Assignment.ID) || strings.Contains(r.stderr, "landing cleanup{assignment="+m.creation.Assignment.ID+",") {
			t.Fatalf("member %s was not released by the stack: stderr %q", m.request, r.stderr)
		}
		if descendant(t, "git", "-C", f.root, "cat-file", "-e", top+":specs/"+m.slug).Run() == nil {
			t.Fatalf("destination still holds specs/%s", m.slug)
		}
		wantSources = append(wantSources, "landing source{review_base="+f.base+",assignment_start="+m.creation.Assignment.Start+"}")
	}
	if marker := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); marker != top {
		t.Fatalf("green marker = %s, want the stack tip %s", marker, top)
	}
	var sources []string
	for _, line := range strings.Split(r.stderr, "\n") {
		if strings.HasPrefix(line, "landing source{") {
			sources = append(sources, line)
		}
	}
	if strings.Join(sources, "\n") != strings.Join(wantSources, "\n") {
		t.Fatalf("landing source lines = %q, want %q", sources, wantSources)
	}
	table, err := stackMembersTable([]stackMemberRow{
		{member: 1, label: f.members[0].creation.Assignment.Label, sourceTip: f.members[0].tip, published: first, outcome: "landed", census: 0},
		{member: 2, label: f.members[1].creation.Assignment.Label, sourceTip: f.members[1].tip, published: top, outcome: "landed", census: 0},
	})
	mustNoError(t, err)
	if !strings.Contains(r.stdout, table) {
		t.Fatalf("stack stdout = %q, want the members table %q", r.stdout, table)
	}
	if want := "stack{destination_base=" + f.base + ",published_commit=" + top + ",landed=2,worktree=released}\n"; !strings.HasSuffix(r.stdout, want) {
		t.Fatalf("stack stdout = %q, want the last line %q", r.stdout, want)
	}
}

// TestLandStackRefusesAnUnreadableStack is LS104. A stack file that the reader refuses
// exits 2 before the ledger read.
func TestLandStackRefusesAnUnreadableStack(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-unreadable", 2, nil)
	r := runVerb(t, verbLand, f.call("--stack", filepath.Join(t.TempDir(), "absent.json")))
	if r.exit != 2 || r.stdout != "" || !strings.Contains(r.stderr, "absent or not a regular file") {
		t.Fatalf("unreadable stack = (%d, %q, %q), want exit 2 with the reader refusal", r.exit, r.stdout, r.stderr)
	}
	requireNoLease(t, f.root)
}

// TestLandStackRefusesARepeatedMember is LS63. Two members that name one worktree exit 2
// and name both indexes, whether the second spells the path or the label. The label row
// repeats the request too, so a repeated request reaches the same refusal.
func TestLandStackRefusesARepeatedMember(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-repeated", 2, nil)
	for _, row := range []struct {
		name, path string
	}{
		{"path", f.members[0].creation.Path},
		{"label", f.members[0].creation.Assignment.Label},
	} {
		t.Run(row.name, func(t *testing.T) {
			again := f.members[0]
			again.creation.Path = row.path
			file := f.writeStack(t, []stackFixtureMember{f.members[0], again}, f.members[0].tip, f.members[0].tip)
			r := runVerb(t, verbLand, f.call("--stack", file))
			if r.exit != 2 || r.stdout != "" || !strings.Contains(r.stderr, "stack members 1 and 2 name one worktree") {
				t.Fatalf("repeated member = (%d, %q, %q), want exit 2 naming members 1 and 2", r.exit, r.stdout, r.stderr)
			}
			requireNoLease(t, f.root)
		})
	}
}

// TestLandStackPrintsEveryMemberRefusal is LS50 and LS51. Two members that fail a source
// proof print one refusal each, with each member index, and leave no lease.
func TestLandStackPrintsEveryMemberRefusal(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-every-refusal", 2, nil)
	for _, m := range f.members {
		mustWrite(t, filepath.Join(m.creation.Path, m.file), []byte("dirty\n"), 0o644)
	}
	_, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 1 {
		t.Fatalf("stack = (%d, %q, %q), want exit 1", r.exit, r.stdout, r.stderr)
	}
	for i, m := range f.members {
		if !strings.Contains(r.stdout, "refused{detail=stack member "+strconv.Itoa(i+1)+" ("+m.creation.Assignment.Label+"): ") {
			t.Fatalf("stack stdout = %q, want the refusal of member %d", r.stdout, i+1)
		}
	}
	if got := strings.Count(r.stdout, "refused{"); got != 2 {
		t.Fatalf("stack stdout = %q, want two refusals, got %d", r.stdout, got)
	}
	requireNoLease(t, f.root)
}

// TestLandStackRefusesAHeldLease is LS52. Under another landing's lease the stack refuses
// before it composes, so the gate never runs.
func TestLandStackRefusesAHeldLease(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-held-lease", 2, nil)
	holdLandingLease(t, f.root, otherLandingID)
	_, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=landing destination refs/heads/main is leased by assignment "+otherLandingID+",") {
		t.Fatalf("stack under a held lease = (%d, %q, %q), want the lease refusal", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("stack under a held lease ran the gate: %v", err)
	}
}

// TestLandStackHoldsTheLease is LS53. A single landing that starts during the stack's grade
// refuses with the lease of member 1's assignment.
func TestLandStackHoldsTheLease(t *testing.T) {
	t.Parallel()
	started, release, step := gateBarrier(t)
	f := stackLandingFixture(t, "stack-holds-lease", 2, step)
	single := "stack-holds-lease-single"
	other := mustCreate(t, f.root, f.home, single, "single landing")
	commitmenttest.Admit(t, other.Path, single, closureSpec)
	commitInWorktree(t, other.Path, "single.txt", "single bytes\n", "single source")
	otherTip := gitOutput(t, other.Path, "rev-parse", "HEAD")
	var r2 verbResult
	done := barrierLanding{landingFixture: f.landingFixture, started: started, release: release}.during(func() {
		r2 = runVerb(t, verbLand, f.call(specLessLandArgs(single, f.base, otherTip, other.Path)...))
	})
	_, args := f.stack(t)
	r1 := runVerb(t, verbLand, f.call(args...))
	select {
	case err := <-done:
		mustNoError(t, err)
	default:
		t.Fatalf("stack finished without reaching its gate barrier: (%d, %q, %q)", r1.exit, r1.stdout, r1.stderr)
	}
	owner := f.members[0].creation.Assignment.ID
	if r2.exit != 1 || !strings.Contains(r2.stdout, "is leased by assignment "+owner+",") {
		t.Fatalf("single landing during the stack = (%d, %q, %q), want the lease refusal naming %s", r2.exit, r2.stdout, r2.stderr, owner)
	}
	if r1.exit != 0 {
		t.Fatalf("stack = (%d, %q, %q), want a publication", r1.exit, r1.stdout, r1.stderr)
	}
}
