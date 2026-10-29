package recordcmd_test

import (
	"bytes"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
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

// stageOwed stages a biting answer at the current tip for each probe requirement of chunk that
// ticket owes, with dirty as the run's dirty paths.
func stageOwed(f *recordtest.Fixture, home, chunk, ticket string, dirty ...string) {
	f.T.Helper()
	for _, requirement := range f.Plan.Chunk(chunk).Verification {
		if requirement.Ticket != ticket {
			continue
		}
		run := f.Run(0, biting)
		run.Dirty = append(run.Dirty, dirty...)
		f.Stage(home, chunk, requirement.ID, run)
	}
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

// stagedBeforeTheTip commits 2.md and then 3.md in chunk 2 of a twoTickets fixture, stages the
// runs of each ticket at its own commit, and records 3.md. The 2.md runs carry dirty. It answers
// the 2.md commit.
func stagedBeforeTheTip(t *testing.T, f *recordtest.Fixture, home string, dirty ...string) string {
	t.Helper()
	freezeChunk1(t, f, home)
	_, own := advance(f, "ticket 2.md")
	stageOwed(f, home, "2", "2.md", dirty...)
	advance(f, "ticket 3.md")
	stageOwed(f, home, "2", "3.md")
	succeed(t, f, ticketArgs(t, "2", "3.md", "last"))
	return own
}

// RT21, RT22, RT23, RT24, RT28, RT29, RT30, RT31, RT32, RT36, RT60, RT61, RC24: each refusal names
// its cause and leaves the record bytes and every staged excerpt as they were.
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
		{"version 1 plan with no requirement", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := recordtest.Pooled(t, 1, recordtest.Owing(nil))
			return f, home, ticketArgs(t, "1", "1.md", "p"), "the ticket form reads a version 2 plan"
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
		{"stage at the chunk base", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			stageOwed(f, home, "1", "1.md")
			base, _ := advance(f, "chunk 1")
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "requirement tests: the staged run at " + base + " is not on the chunk line " + base + ".."
		}},
		{"later stage off the chunk tip", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 1)
			base, tip := advance(f, "chunk 1")
			stageOwed(f, home, "1", "1.md")
			run := f.Run(0, "ok\n")
			run.Head = f.Git(append(gittest.IdentityArgs(), "commit-tree", tip+"^{tree}", "-p", base, "-m", "off the chunk line")...)
			f.Stage(home, "1", "additional", run)
			return f, home, ticketArgs(t, "1", "1.md", "p", "--base", base), "requirement additional: the staged run at " + run.Head + " is not on the chunk line " + base + ".." + tip
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
		{"dirty stage at an earlier ticket commit", func(t *testing.T) (*recordtest.Fixture, string, []string, string) {
			f, home := pooled(t, 3, twoTickets)
			stagedBeforeTheTip(t, f, home, "main.go")
			return f, home, ticketArgs(t, "2", "2.md", "early"), "requirement tests: the staged run had the dirty path main.go"
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
	if ids, want := verificationIDs(chunk), []string{"p-tests", "p-additional"}; !slices.Equal(ids, want) {
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

// RC16: each written entry carries the commit of its staged run.
func TestTicketWritesTheSourceCommit(t *testing.T) {
	f, home := pooled(t, 1)
	base, tip := advance(f, "chunk 1")
	stageOwed(f, home, "1", "1.md")
	succeed(t, f, ticketArgs(t, "1", "1.md", "p", "--base", base))
	commits := []string{}
	for _, entry := range chunkEntry(t, f, "1").Verification {
		commits = append(commits, entry.SourceCommit)
	}
	if want := []string{tip, tip}; !slices.Equal(commits, want) {
		t.Errorf("entry source commits = %q, want the staged run head %q", commits, want)
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

// The first record call of a chunk creates its entry, whichever ticket of the chunk makes it.
func TestTicketEarlierTicketCreatesTheChunkEntry(t *testing.T) {
	f, home := pooled(t, 3, twoTickets)
	freezeChunk1(t, f, home)
	_, tip := advance(f, "chunk 2")
	stageOwed(f, home, "2", "2.md")
	succeed(t, f, ticketArgs(t, "2", "2.md", "p"))
	if chunk := chunkEntry(t, f, "2"); chunk.Tip != tip || !slices.Equal(verificationIDs(chunk), []string{"p-tests"}) {
		t.Errorf("chunk 2 after the 2.md call = tip %s with %q, want tip %s with [p-tests]", chunk.Tip, verificationIDs(chunk), tip)
	}
}

// RC9: after a later ticket moved the chunk tip past an earlier ticket's commit, the earlier
// ticket's record call writes its entries at the commit where it staged them.
func TestTicketRecordsAtItsOwnCommitOnTheChunkLine(t *testing.T) {
	f, home := pooled(t, 3, twoTickets)
	own := stagedBeforeTheTip(t, f, home)
	succeed(t, f, ticketArgs(t, "2", "2.md", "early"))
	var want rr.Verification
	f.Bind(&want, own)
	results := chunkEntry(t, f, "2").Verification
	if got := results[len(results)-1]; got.ID != "early-tests" || got.SourceCommit != own || got.SourceDigest != want.SourceDigest {
		t.Errorf("earlier ticket entry = %s at %s with digest %s, want early-tests at %s with digest %s", got.ID, got.SourceCommit, got.SourceDigest, own, want.SourceDigest)
	}
}

// repairedEarlierTicket records both tickets of chunk 2 of a twoTickets fixture at the chunk tip,
// then commits a repair of 2.md and records 2.md with its runs staged at that repair commit. It
// answers the chunk 2 entry before the repair and the repair commit.
func repairedEarlierTicket(t *testing.T) (*recordtest.Fixture, rr.Chunk, string) {
	t.Helper()
	f, home := pooled(t, 3, twoTickets)
	freezeChunk1(t, f, home)
	advance(f, "chunk 2")
	stageOwed(f, home, "2", "3.md")
	succeed(t, f, ticketArgs(t, "2", "3.md", "last"))
	stageOwed(f, home, "2", "2.md")
	succeed(t, f, ticketArgs(t, "2", "2.md", "early"))
	frozen := chunkEntry(t, f, "2")
	_, repair := advance(f, "repair 2.md")
	stageOwed(f, home, "2", "2.md")
	succeed(t, f, ticketArgs(t, "2", "2.md", "repair"))
	return f, frozen, repair
}

// RC11: a non-last ticket's record call at a commit past the chunk tip refreezes the chunk at
// that commit and keeps the chunk base.
func TestNonLastTicketRefreezesAtItsRepairCommit(t *testing.T) {
	f, frozen, repair := repairedEarlierTicket(t)
	if chunk := chunkEntry(t, f, "2"); chunk.Base != frozen.Base || chunk.Tip != repair {
		t.Errorf("chunk 2 after the 2.md repair = %s..%s, want %s..%s", chunk.Base, chunk.Tip, frozen.Base, repair)
	}
}

// RC10: a repair session's record call appends the entries of its own ticket and leaves every
// other entry as it was.
func TestRepairRecordLeavesOtherTicketsEntries(t *testing.T) {
	f, frozen, _ := repairedEarlierTicket(t)
	chunk, kept := chunkEntry(t, f, "2"), len(frozen.Verification)
	ids := verificationIDs(chunk)
	if len(ids) != kept+1 || ids[kept] != "repair-tests" || !reflect.DeepEqual(chunk.Verification[:kept], frozen.Verification) {
		t.Errorf("chunk 2 results after the 2.md repair = %q, want %q unchanged and then repair-tests", ids, verificationIDs(frozen))
	}
}

// RC12: a run staged at a record-only descendant of the chunk tip binds the tip, and its record
// call leaves the tip unchanged.
func TestTicketStagedAfterARecordCommitBindsTheTip(t *testing.T) {
	f, home := pooled(t, 3, twoTickets)
	freezeChunk1(t, f, home)
	_, tip := advance(f, "chunk 2")
	stageOwed(f, home, "2", "3.md")
	succeed(t, f, ticketArgs(t, "2", "3.md", "last"))
	f.Commit("record 3.md")
	stageOwed(f, home, "2", "2.md")
	succeed(t, f, ticketArgs(t, "2", "2.md", "early"))
	chunk := chunkEntry(t, f, "2")
	if got := chunk.Verification[len(chunk.Verification)-1]; chunk.Tip != tip || got.SourceCommit != tip {
		t.Errorf("after a record commit, chunk 2 tip = %s and early-tests binds %s, want both %s", chunk.Tip, got.SourceCommit, tip)
	}
}
