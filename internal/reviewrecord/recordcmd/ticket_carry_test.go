package recordcmd_test

import (
	"bytes"
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// carryPlan is the chunk plan of carriedRepair. The repair changes a test file of
// ./internal/a, so only pkg is reached. The carried system requirement comes before pkg, so an
// owed-order write puts a carried entry before a staged one.
var carryPlan = []rr.Requirement{
	{ID: "system", Command: "bench test --check system", Probe: "omit the check"},
	{ID: "pkg", Command: "bench test --package ./internal/a"},
	{ID: "other", Command: "bench test --package ./internal/b"},
}

// carriedRepair records ticket 1.md of a one-chunk delegated fixture that owes carryPlan, with
// each requirement staged. It then dispatches the repair of 1.md through Reassign and commits a
// test-only repair in ./internal/a. It answers the fixture, the home, the chunk entry before the
// repair, and the repair commit.
func carriedRepair(t *testing.T) (*recordtest.Fixture, string, rr.Chunk, string) {
	t.Helper()
	f, home := pooled(t, 1, recordtest.Owing(carryPlan))
	base, _ := advance(f, "chunk 1")
	stageOwed(f, home, "1", "1.md")
	succeed(t, f, ticketArgs(t, "1", "1.md", "first", "--base", base))
	frozen := chunkEntry(t, f, "1")
	f.Reassign("1.md")
	f.Write("internal/a/repair_test.go", "package a\n")
	f.Commit("repair 1.md")
	return f, home, frozen, f.Tip()
}

// repairArgs is the ticket argv of the repair session, the author of 1.md that the plan names.
func repairArgs(t *testing.T, f *recordtest.Fixture) []string {
	t.Helper()
	args := ticketArgs(t, "1", "1.md", "repair")
	args[slices.Index(args, "--performer")+1] = repairSession(t, f)
	return args
}

func repairSession(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	session, ok := f.Plan.Author("1.md")
	if !ok {
		t.Fatal("the plan names no author of 1.md")
	}
	return session
}

// carriedCell is the carried cell of the one ticket output row.
func carriedCell(t *testing.T, out string) any {
	t.Helper()
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("ticket")
	if err != nil || len(rows) != 1 {
		t.Fatalf("ticket output %q: %v; want one row", out, err)
	}
	return rows[0].(map[string]any)["carried"]
}

// CF33, CF34, CF35, CF36, CF37, CF40: a call that stages only the reached requirement carries
// the others in owed order. Each carried entry keeps the run of its candidate, names the repair
// session, and references the candidate.
func TestTicketCarriesTheUnreachedRequirement(t *testing.T) {
	f, home, frozen, repair := carriedRepair(t)
	f.Stage(home, "1", "pkg", f.Run(0, "ok\n"))
	out := succeed(t, f, repairArgs(t, f))
	if cell := carriedCell(t, out); cell != "system,other" {
		t.Errorf("carried cell = %v, want system,other", cell)
	}
	added := chunkEntry(t, f, "1").Verification[len(frozen.Verification):]
	if ids := verificationIDs(rr.Chunk{Verification: added}); !slices.Equal(ids, []string{"repair-system", "repair-pkg", "repair-other"}) {
		t.Fatalf("repair entries = %q, want the owed order", ids)
	}
	if added[1].SourceCommit != repair {
		t.Errorf("staged entry source = %s, want the repair commit %s", added[1].SourceCommit, repair)
	}
	for _, i := range []int{0, 2} {
		prior := frozen.Verification[i]
		want := prior
		want.ID, want.Performer = added[i].ID, repairSession(t, f)
		want.NativeRef = rr.NativeRef{Ref: "carried:" + prior.ID, Digest: prior.NativeRef.Digest, Excerpt: prior.NativeRef.Excerpt}
		if prior.Probe != nil {
			probe := *prior.Probe
			probe.NativeRef = want.NativeRef
			want.Probe = &probe
		}
		if !reflect.DeepEqual(added[i], want) {
			t.Errorf("carried entry = %+v, want %+v", added[i], want)
		}
	}
}

// CF38, CF39: a stage wins over a carry, so a staged requirement that the delta does not reach
// records from its stage at the repair commit, and the carried cell is none.
func TestTicketCarriesNothingWhenEachRequirementIsStaged(t *testing.T) {
	f, home, _, repair := carriedRepair(t)
	stageOwed(f, home, "1", "1.md")
	if cell := carriedCell(t, succeed(t, f, repairArgs(t, f))); cell != "none" {
		t.Errorf("carried cell = %v, want none", cell)
	}
	results := chunkEntry(t, f, "1").Verification
	if got := results[len(results)-len(carryPlan)]; got.ID != "repair-system" || got.SourceCommit != repair || got.NativeRef.Ref != "fixture:ticket" {
		t.Errorf("system entry = %s at %s with reference %s, want repair-system at %s with the call reference", got.ID, got.SourceCommit, got.NativeRef.Ref, repair)
	}
}

// CF41, CF42, CF43: the refusal names only the unstaged requirement that cannot carry, with its
// reason, and leaves the record bytes and every stage unchanged.
func TestTicketNamesTheChecksToRun(t *testing.T) {
	f, home, _, _ := carriedRepair(t)
	f.Stage(home, "1", "other", f.Run(0, "ok\n"))
	before, stages := recordBytes(t, f), responseboundtest.StoreFiles(t, home)
	out, code := recordcmd.Command(f.Root, repairArgs(t, f))
	// The excerpt-form repair route follows the cause.
	if want := "these owed requirements have no staged excerpt: pkg; carry: pkg (the delta reaches the command); write it with "; code != 1 || !strings.Contains(out, want) {
		t.Fatalf("ticket = exit %d, output %q; want exit 1 naming %q", code, out, want)
	}
	if after := recordBytes(t, f); !bytes.Equal(before, after) {
		t.Errorf("the refusal changed the record bytes:\n%s", after)
	}
	if after := responseboundtest.StoreFiles(t, home); !slices.Equal(stages, after) {
		t.Errorf("store files after the refusal = %q, want %q", after, stages)
	}
}

// CF48: after the carry call and a record commit, the chunk grade passes the verification step
// of the chunk and refuses at its review step, because no review is recorded.
func TestCarriedEntryPassesTheChunkGrade(t *testing.T) {
	f, home, _, _ := carriedRepair(t)
	f.Stage(home, "1", "pkg", f.Run(0, "ok\n"))
	succeed(t, f, repairArgs(t, f))
	f.Commit("record the repair of 1.md")
	err := rr.Check(f.Root, f.Tree(), f.Tip(), recordtest.Spec, "1", false)
	if repair := (rr.Repair{}); !errors.As(err, &repair) || repair.Form != rr.FormReview || repair.Chunk != "1" {
		t.Fatalf("chunk grade = %v, want a review-form repair of chunk 1", err)
	}
}
