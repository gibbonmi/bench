package recordcmd_test

import (
	"bytes"
	"reflect"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// pairedSecond moves the third ticket into the second chunk, so the chunks are 1, 2, and 4,
// and chunk 2 holds two tickets. It precedes Delegate in the option list.
func pairedSecond(plan *rr.Plan) {
	plan.Chunks[1].Tickets = append(plan.Chunks[1].Tickets, plan.Chunks[2].Tickets...)
	plan.Chunks = slices.Delete(plan.Chunks, 2, 3)
}

// boundBuild is an undispatched delegated plan of chunks 1, 2, and 4 in the checkout of the
// assignment that holds the open build of its approved spec. It also answers the primary
// checkout, which holds main.
func boundBuild(t *testing.T) (*recordtest.Fixture, string) {
	t.Helper()
	t.Setenv("BENCH_HOME", t.TempDir())
	primary := gittest.RepoOnBranch(t, "main")
	f := recordtest.Prepare(t, primary, 4, recordtest.Spec, recordtest.Body, pairedSecond, recordtest.Delegate, recordtest.Undispatched)
	commitmenttest.SeedAdmission(t, primary, recordtest.Spec)
	f.Commit("approve the delivery")
	build := *f
	build.Root = commitmenttest.BindBuild(t, primary, "build", recordtest.Spec)
	return &build, primary
}

// advanceMain commits one change to main in the primary checkout, so the build's source lacks
// the default-branch tip.
func advanceMain(t *testing.T, primary string) {
	t.Helper()
	commitmenttest.Write(t, primary, "advance.txt", "main advances\n")
	commitmenttest.Commit(t, primary, "advance main")
}

// ledgerBytes is the exact intent ledger of the repository at root.
func ledgerBytes(t *testing.T, root string) []byte {
	t.Helper()
	path, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	return bytesOf(t, path)
}

// requireGap fails unless the ledger at the root of f holds only the gap record want, in the
// milestone and outcome of the build's claim.
func requireGap(t *testing.T, f *recordtest.Fixture, want intent.ChunkGap) {
	t.Helper()
	got, claim := commitmenttest.ClaimedGaps(t, f.Root)
	want.Milestone, want.Outcome = claim.Milestone, claim.Outcome
	if !reflect.DeepEqual(got, []intent.ChunkGap{want}) {
		t.Fatalf("gap records = %+v, want only %+v", got, want)
	}
}

// requireLedgerKept runs dispatch and fails unless the ledger at the root of f keeps its bytes.
func requireLedgerKept(t *testing.T, f *recordtest.Fixture, dispatch func()) {
	t.Helper()
	before := ledgerBytes(t, f.Root)
	dispatch()
	if after := ledgerBytes(t, f.Root); !bytes.Equal(before, after) {
		t.Fatalf("ledger after the dispatch = %s, want the bytes %s", after, before)
	}
}

// TestFirstDispatchClosesTheGap is BW65 and BW69. The first dispatch of the first chunk, with
// no gap record, writes a closed gap for that chunk. The first dispatch of chunk 2 closes the
// gap that the chunk 1 checkpoint opened, and the last chunk's first dispatch records Last.
func TestFirstDispatchClosesTheGap(t *testing.T) {
	f, _ := boundBuild(t)
	dispatch(t, f, "1.md")
	requireGap(t, f, intent.ChunkGap{Chunk: "1"})
	f.Commit("dispatch 1.md")
	commitmenttest.PutGap(t, f.Root, intent.ChunkGap{Outcome: commitmenttest.DeliveryOutcome, Chunk: "1", Open: true})
	dispatch(t, f, "2.md")
	requireGap(t, f, intent.ChunkGap{Chunk: "2"})
	f.Commit("dispatch 2.md")
	dispatch(t, f, "4.md")
	requireGap(t, f, intent.ChunkGap{Chunk: "4", Last: true})
}

// TestDispatchRefusesAnUnfoldedMain is BW66 and BW67. A first dispatch while main holds a
// commit that the source lacks refuses with the fold route. The plan fence and the ledger
// keep their bytes.
func TestDispatchRefusesAnUnfoldedMain(t *testing.T) {
	f, primary := boundBuild(t)
	advanceMain(t, primary)
	owner, owned := intent.AssignmentForWorktree(f.Root)
	slug, err := rr.Slug(recordtest.Spec)
	if !owned || err != nil {
		t.Fatalf("build owner = %+v (%t), slug %v", owner, owned, err)
	}
	want := "dispatch refused: the default branch tip " + gittest.Output(t, primary, "rev-parse", "main") +
		" is not an ancestor of the source; fold it with bench worktree merge --from main " + owner.ID +
		" and record it with bench record fold " + slug + " --source <fold-commit> before the dispatch"
	requireLedgerKept(t, f, func() { refuseAssignment(t, f, 1, want, assignArgs(t, *f.Plan.Execution, "1.md")) })
}

// TestDispatchInsideAChunkKeepsTheGap is BW68 and BW95. A second dispatch inside chunk 2, and
// a repair dispatch to a chunk 1 ticket while the record holds the last chunk, each keep the
// ledger bytes. Main moves first, so a dispatch that read as a chunk's first would refuse.
func TestDispatchInsideAChunkKeepsTheGap(t *testing.T) {
	f, primary := boundBuild(t)
	for _, ticket := range []string{"1.md", "2.md"} {
		dispatch(t, f, ticket)
		f.Commit("dispatch " + ticket)
	}
	advanceMain(t, primary)
	requireLedgerKept(t, f, func() { dispatch(t, f, "3.md") })
	f.Commit("dispatch 3.md")
	commitmenttest.PutGap(t, f.Root, intent.ChunkGap{Outcome: commitmenttest.DeliveryOutcome, Chunk: "4", Last: true})
	repair := map[string]string{"--performer": "repair-author", "--assignment": "repair-1.md", "--ref": "native:repair-1.md",
		"--predecessor": recordtest.Author("1.md"), "--trigger": "terminal-failure", "--stopped": "native:stopped-1", "--preserved": "refs/bench/reset/1"}
	requireLedgerKept(t, f, func() { dispatch(t, f, "1.md", repair) })
	if history := committedPlan(t, f, "repair 1.md").Execution.Assignments["1.md"]; len(history) != 2 {
		t.Fatalf("1.md history = %+v, want the first dispatch and the repair", history)
	}
}
