package recordcmd_test

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// pooled is a delegated fixture in an assignment worktree with count planned chunks. Each option
// reshapes the plan before Delegate maps the obligations onto the tickets.
func pooled(t *testing.T, count int, options ...recordtest.Option) (*recordtest.Fixture, string) {
	t.Helper()
	return recordtest.Pooled(t, count, append(options, recordtest.Delegate)...)
}

// twoTickets moves the ticket of chunk 3 into chunk 2, so chunk 2 lists 2.md and then 3.md.
// Delegate then gives requirement tests to 2.md and requirement additional to 3.md.
func twoTickets(plan *rr.Plan) {
	plan.Chunks[1].Tickets = append(plan.Chunks[1].Tickets, plan.Chunks[2].Tickets...)
	plan.Chunks = plan.Chunks[:2]
}

// ticketArgs is the ticket argv for ticket of chunk with the evidence prefix, and the extra flags
// after it.
func ticketArgs(t *testing.T, chunk, ticket, prefix string, extra ...string) []string {
	t.Helper()
	return append([]string{"ticket", slug(t), "--chunk", chunk, "--ticket", ticket, "--id-prefix", prefix, "--performer", recordtest.Author(ticket),
		"--model", "unknown", "--effort", "unknown", "--ref", "fixture:ticket"}, extra...)
}

// biting is a biting single-form probe answer whose command exited 0.
var biting = answer("bit", "yes", 0, "1", report)

// stageOwed stages a passing run at the current tip for each requirement of chunk that ticket
// owes, with dirty as the run's dirty paths: a biting answer for a probe requirement.
func stageOwed(f *recordtest.Fixture, home, chunk, ticket string, dirty ...string) {
	f.T.Helper()
	for _, requirement := range f.Plan.Chunk(chunk).Verification {
		if requirement.Ticket != ticket {
			continue
		}
		run := f.Run(0, "ok\n")
		if requirement.Probe != "" {
			run = f.Run(0, biting)
		}
		run.Dirty = append(run.Dirty, dirty...)
		f.Stage(home, chunk, requirement.ID, run)
	}
}

// recordPath is the record path of the fixture spec, as a run's dirty paths name it.
func recordPath(t *testing.T) string {
	t.Helper()
	path, err := rr.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// freezeChunk1 commits a chunk 1 change and records ticket 1.md with the prefix c1 and an
// explicit base. It answers the chunk 1 tip.
func freezeChunk1(t *testing.T, f *recordtest.Fixture, home string) string {
	t.Helper()
	base, tip := advance(f, "chunk 1")
	stageOwed(f, home, "1", "1.md")
	succeed(t, f, ticketArgs(t, "1", "1.md", "c1", "--base", base))
	return tip
}

// recordBytes answers the record bytes of f, or nil when the record is absent.
func recordBytes(t *testing.T, f *recordtest.Fixture) []byte {
	t.Helper()
	data, err := os.ReadFile(recordFile(t, f))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// chunkEntry is the entry of chunk id in the record of f.
func chunkEntry(t *testing.T, f *recordtest.Fixture, id string) rr.Chunk {
	t.Helper()
	for _, chunk := range read(t, f).Chunks {
		if chunk.ID == id {
			return chunk
		}
	}
	t.Fatalf("the record holds no chunk %s", id)
	return rr.Chunk{}
}

// RT21, RT22, RT23, RT24, RT28, RT29, RT30, RT31, RT32, RT36, RT60, RT61: each refusal names its
// cause and leaves the record bytes and every staged excerpt as they were.
func TestTicketRefusals(t *testing.T) {
	for _, tc := range []struct {
		name    string
		prepare func(t *testing.T) (f *recordtest.Fixture, home string, args []string, want string)
	}{
		{"ticket outside the chunk", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 2)
			return f, home, ticketArgs(t, "2", "1.md", "p"), "ticket 1.md is not a ticket of chunk 2"
		}},
		{"requirements that name no ticket", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := recordtest.Pooled(t, 1)
			return f, home, ticketArgs(t, "1", "1.md", "p"), "the requirements of chunk 1 in the plan at "
		}},
		{"planned seam cell", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			withMap(planned)(f)
			base, _ := advance(f, "chunk 1")
			stageOwed(f, home, "1", "1.md")
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "rows E1 hold a planned seam cell"
		}},
		{"two requirements with no stage", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			base, _ := advance(f, "chunk 1")
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "these owed requirements have no staged excerpt: tests, additional"
		}},
		{"stage from another commit", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			stageOwed(f, home, "1", "1.md")
			base, _ := advance(f, "chunk 1")
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "requirement tests: the staged run is at " + f.Git("rev-parse", "HEAD^")
		}},
		{"later stage from another commit", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			base, _ := advance(f, "chunk 1")
			stageOwed(f, home, "1", "1.md")
			run := f.Run(0, "ok\n")
			run.Head = base
			f.Stage(home, "1", "additional", run)
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "requirement additional: the staged run is at " + base
		}},
		{"stage with a dirty Go file", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			base, _ := advance(f, "chunk 1")
			stageOwed(f, home, "1", "1.md", "main.go")
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "requirement tests: the staged run had the dirty path main.go"
		}},
		{"stage with the record path and a dirty Go file", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			base, _ := advance(f, "chunk 1")
			stageOwed(f, home, "1", "1.md", recordPath(t), "main.go")
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "requirement tests: the staged run had the dirty path main.go"
		}},
		{"earlier ticket with no chunk entry", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 3, twoTickets)
			freezeChunk1(t, f, home)
			advance(f, "chunk 2")
			stageOwed(f, home, "2", "2.md")
			return f, home, ticketArgs(t, "2", "2.md", "p"), "chunk 2 has no entry: the result needs a current chunk entry; write it with bench record ticket '" + slug(t) + "' --chunk '2' --ticket '3.md' "
		}},
		{"first chunk with no base", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			advance(f, "chunk 1")
			stageOwed(f, home, "1", "1.md")
			return f, home, ticketArgs(t, "1", "1.md", "p"), "name its base with --base"
		}},
		{"plan change that removes a recorded chunk", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 2)
			freezeChunk1(t, f, home)
			rename(f, "1x", "2")
			_, tip := advance(f, "chunk 2")
			stageOwed(f, home, "2", "2.md")
			return f, home, ticketArgs(t, "2", "2.md", "p"),
				"recorded chunk 1 is not a chunk of the plan at " + tip + "; map it to its new chunks; write it with bench record amendment '" + slug(t) + "' --source '" + tip + "'"
		}},
		{"last owed verification", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			// Chunk 2 owes first and then additional, and chunk 1 already holds c1-additional.
			f, home := pooled(t, 2, func(plan *rr.Plan) { plan.Chunks[1].Verification[0].ID = "first" })
			freezeChunk1(t, f, home)
			replan(f, "go test -count=1 ./...")
			advance(f, "chunk 2")
			stageOwed(f, home, "2", "2.md")
			return f, home, ticketArgs(t, "2", "2.md", "c1"), "evidence ID c1-additional is already in the record"
		}},
		{"chunk step after a moved plan digest", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			// Chunk 2 has no entry, so the chunk 3 freeze has no base after the amendment.
			f, home := pooled(t, 3)
			freezeChunk1(t, f, home)
			replan(f, "go test -count=1 ./...")
			advance(f, "chunk 3")
			stageOwed(f, home, "3", "3.md")
			return f, home, ticketArgs(t, "3", "3.md", "p"), "name its base with --base"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, home, args, want := tc.prepare(t)
			before, stages := recordBytes(t, f), responseboundtest.StoreFiles(t, home)
			out, code := recordcmd.Command(f.Root, args)
			if code != 1 || !strings.Contains(out, want) {
				t.Fatalf("ticket = exit %d, output %q; want exit 1 naming %q", code, out, want)
			}
			if after := recordBytes(t, f); !bytes.Equal(before, after) {
				t.Errorf("the refusal changed the record bytes:\n%s", after)
			}
			if after := responseboundtest.StoreFiles(t, home); !slices.Equal(stages, after) {
				t.Errorf("store files after the refusal = %q, want the staged excerpts %q", after, stages)
			}
		})
	}
}

// RT18, RT33, RT34, RT35, RT56: the last ticket of the first chunk creates the absent record with
// the chunk entry at the tip and one verification for each owed requirement under the prefix. The
// probe requirement records the staged answer, and the used stages are gone.
func TestTicketFreezesTheFirstChunk(t *testing.T) {
	f, home := pooled(t, 1)
	base, tip := advance(f, "chunk 1")
	stageOwed(f, home, "1", "1.md")
	succeed(t, f, ticketArgs(t, "1", "1.md", "p", "--base", base))
	chunk := chunkEntry(t, f, "1")
	if chunk.Base != base || chunk.Tip != tip {
		t.Errorf("chunk 1 pair = %s..%s, want %s..%s", chunk.Base, chunk.Tip, base, tip)
	}
	ids := []string{}
	for _, entry := range chunk.Verification {
		ids = append(ids, entry.ID)
	}
	if want := []string{"p-tests", "p-additional"}; !slices.Equal(ids, want) {
		t.Fatalf("verification IDs = %q, want %q", ids, want)
	}
	want := rr.Probe{Outcome: "bit", ExitCode: 1, Restore: "pass", Mutation: f.Plan.Chunks[0].Verification[0].Probe, NativeRef: chunk.Verification[0].NativeRef}
	if got := chunk.Verification[0]; got.Probe == nil || !reflect.DeepEqual(*got.Probe, want) || *got.ExitCode != 0 || got.NativeRef.Excerpt != biting {
		t.Errorf("probe entry = %+v with probe %+v, want exit 0, %+v, and the staged bytes", got, got.Probe, want)
	}
	if store := responseboundtest.StoreFiles(t, home); len(store) != 0 {
		t.Errorf("store files after the write = %q, want no staged excerpt", store)
	}
}

// RT19, RT20, RT27, RT55, RT62: a later chunk takes the covered commit after its predecessor as
// base, and a moved plan digest takes one amendment. A run with only the record path dirty stages.
// A --tip before HEAD freezes the chunk at that tip.
func TestTicketFreezesALaterChunk(t *testing.T) {
	for _, tc := range []struct {
		name    string
		between func(f *recordtest.Fixture, tip string) string
		amended bool
		dirty   bool
		pastTip bool
	}{
		{"unchanged plan", func(_ *recordtest.Fixture, tip string) string { return tip }, false, false, false},
		{"moved plan", func(f *recordtest.Fixture, tip string) string { replan(f, "go test -count=1 ./..."); return tip }, true, false, false},
		{"recorded fold", func(f *recordtest.Fixture, _ string) string {
			f.Commit("retain chunk 1")
			fold := f.HandFold(f.Mainline(map[string]string{"main.txt": "main\n"}), nil)
			f.RecordFold(rr.FoldDelta{})
			return fold
		}, false, false, false},
		{"dirty record path", func(_ *recordtest.Fixture, tip string) string { return tip }, false, true, false},
		{"tip before HEAD", func(_ *recordtest.Fixture, tip string) string { return tip }, false, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, home := pooled(t, 2)
			wantBase := tc.between(f, freezeChunk1(t, f, home))
			from := read(t, f).PlanDigest
			_, tip := advance(f, "chunk 2")
			var dirty []string
			if tc.dirty {
				dirty = []string{recordPath(t)}
			}
			stageOwed(f, home, "2", "2.md", dirty...)
			var extra []string
			if tc.pastTip {
				advance(f, "past the chunk 2 tip")
				extra = []string{"--tip", tip}
			}
			succeed(t, f, ticketArgs(t, "2", "2.md", "c2", extra...))
			record, chunk := read(t, f), chunkEntry(t, f, "2")
			if chunk.Base != wantBase || chunk.Tip != tip || len(chunk.Verification) != 2 {
				t.Errorf("chunk 2 = %s..%s with %d results, want %s..%s with 2", chunk.Base, chunk.Tip, len(chunk.Verification), wantBase, tip)
			}
			amendments := []rr.Amendment{}
			if tc.amended {
				amendments = append(amendments, rr.Amendment{From: from, To: chunk.PlanDigest, ChunkIDs: map[string][]string{"1": {"1"}}})
			}
			if !reflect.DeepEqual(append([]rr.Amendment{}, record.Amendments...), amendments) {
				t.Errorf("amendments = %+v, want %+v", record.Amendments, amendments)
			}
		})
	}
}

// RT26: a refreeze of the last ticket keeps the base of the chunk entry, and a --base on the
// refreeze takes precedence over that base.
func TestTicketRefreezeKeepsTheBase(t *testing.T) {
	for _, tc := range []struct {
		name  string
		given bool
	}{
		{"no --base", false},
		{"--base given", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, home := pooled(t, 2)
			chunk1 := freezeChunk1(t, f, home)
			_, base := advance(f, "between")
			advance(f, "chunk 2")
			stageOwed(f, home, "2", "2.md")
			succeed(t, f, ticketArgs(t, "2", "2.md", "a", "--base", base))
			_, tip := advance(f, "chunk 2 repair")
			stageOwed(f, home, "2", "2.md")
			var extra []string
			if tc.given {
				base, extra = chunk1, []string{"--base", chunk1}
			}
			succeed(t, f, ticketArgs(t, "2", "2.md", "b", extra...))
			if chunk := chunkEntry(t, f, "2"); chunk.Base != base || chunk.Tip != tip {
				t.Errorf("refrozen chunk 2 = %s..%s, want %s..%s", chunk.Base, chunk.Tip, base, tip)
			}
		})
	}
}

// RT25: an earlier ticket's call records its result at the chunk tip and leaves the chunk entry,
// even when HEAD moved past it.
func TestTicketEarlierTicketKeepsTheEntry(t *testing.T) {
	f, home := pooled(t, 3, twoTickets)
	freezeChunk1(t, f, home)
	advance(f, "chunk 2")
	stageOwed(f, home, "2", "3.md")
	succeed(t, f, ticketArgs(t, "2", "3.md", "last"))
	frozen := chunkEntry(t, f, "2")
	advance(f, "after the freeze")
	run := f.Run(0, biting)
	run.Head = frozen.Tip
	f.Stage(home, "2", "tests", run)
	succeed(t, f, ticketArgs(t, "2", "2.md", "early"))
	chunk := chunkEntry(t, f, "2")
	if chunk.Base != frozen.Base || chunk.Tip != frozen.Tip || len(chunk.Verification) != 2 || chunk.Verification[1].ID != "early-tests" {
		t.Errorf("chunk 2 after the earlier ticket = %+v, want the pair %s..%s and the result early-tests", chunk, frozen.Base, frozen.Tip)
	}
}

// A mutant-set stage records its row: the row's exit, not the set's command exit, and the row's
// probe values.
func TestTicketImportsAMutantRow(t *testing.T) {
	produce, _ := setProducer(t, linked(t, 1))
	data, command := silentThenBiting(produce)
	f, home := pooled(t, 1)
	base, _ := advance(f, "chunk 1")
	stageOwed(f, home, "1", "1.md")
	run := f.Run(command, data)
	run.Mutant = 2
	f.Stage(home, "1", "tests", run)
	succeed(t, f, ticketArgs(t, "1", "1.md", "p", "--base", base))
	got := chunkEntry(t, f, "1").Verification[0]
	if got.Probe == nil || *got.ExitCode != 0 || got.Probe.Outcome != "bit" || got.Probe.ExitCode != 1 || got.Probe.Restore != "pass" {
		t.Errorf("mutant row entry = exit %d, probe %+v; want exit 0 and the biting row", *got.ExitCode, got.Probe)
	}
}

// RT63: the ticket help lists the ticket steps of the completion order, in order.
func TestTicketHelpFollowsCompletionOrder(t *testing.T) {
	want := []string{}
	for _, step := range rr.TicketSteps() {
		want = append(want, step.Form)
	}
	var description string
	for _, row := range recordcmd.HelpRows() {
		if strings.HasPrefix(row.Suffix, " ticket ") {
			description = row.Description
		}
	}
	if _, order, ok := strings.Cut(description, "write in the order "); !ok || !slices.Equal(strings.Split(order, ", "), want) {
		t.Fatalf("ticket help %q, want the order %q", description, want)
	}
}
