package recordcmd_test

import (
	"bytes"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// ticketTable is the TOON table of a ticket response: the bytes before its next= line.
func ticketTable(out string) string {
	table, _, _ := strings.Cut(out, refusalroute.NextField+"=")
	return table
}

// BC29: a ticket that owes no requirement writes the chunk step alone. The call exits 0, freezes
// the chunk at the tip, and prints an empty verifications cell.
func TestTicketFreezesWithoutRequirements(t *testing.T) {
	f, _ := pooled(t, 1, recordtest.Owing(nil))
	base, tip := advance(f, "chunk 1")
	out := succeed(t, f, ticketArgs(t, "1", "1.md", "p", "--base", base))
	if chunk := chunkEntry(t, f, "1"); chunk.Base != base || chunk.Tip != tip || len(chunk.Verification) != 0 {
		t.Errorf("chunk 1 = %s..%s with %d results, want %s..%s with none", chunk.Base, chunk.Tip, len(chunk.Verification), base, tip)
	}
	want, err := toon.Table("ticket", []string{"chunk", "ticket", "amendment", "freeze", "base", "tip", "verifications", "carried"},
		[][]string{{"1", "1.md", "none", "created", base, tip, "", "none"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, want) {
		t.Errorf("ticket output = %q, want the table %q first", out, want)
	}
}

// requireBitEntry requires that chunk's verification entry id is a probe that bit and restored,
// at tip, with exit 0 and the excerpt bytes. Its mutation is the plan probe of the requirement
// that the ID names after its prefix.
func requireBitEntry(t *testing.T, f *recordtest.Fixture, chunk, id, tip, excerpt string) {
	t.Helper()
	results := chunkEntry(t, f, chunk).Verification
	i := slices.IndexFunc(results, func(entry rr.Verification) bool { return entry.ID == id })
	if i < 0 {
		t.Fatalf("entries = %q, want %s", verificationIDs(rr.Chunk{Verification: results}), id)
	}
	got := results[i]
	requirement := strings.SplitN(id, "-", 2)[1]
	planned := f.Plan.Chunk(chunk).Verification
	j := slices.IndexFunc(planned, func(item rr.Requirement) bool { return item.ID == requirement })
	if j < 0 {
		t.Fatalf("chunk %s plans no requirement %s", chunk, requirement)
	}
	want := rr.Probe{Outcome: "bit", ExitCode: 1, Restore: "pass", Mutation: planned[j].Probe}
	if got.Probe != nil {
		want.NativeRef = got.Probe.NativeRef
	}
	if got.SourceCommit != tip || got.NativeRef.Excerpt != excerpt || *got.ExitCode != 0 || got.Probe == nil || !reflect.DeepEqual(*got.Probe, want) {
		t.Errorf("entry %s = at %s, exit %d, probe %+v; want at %s, exit 0, %+v, and the excerpt bytes", id, got.SourceCommit, *got.ExitCode, got.Probe, tip, want)
	}
}

// ticketCase prepares one call of TestTicketStagesTheRetainedProbe. It answers the fixture, the
// home, the chunk base, and the ticket argv before the extra flags.
type ticketCase func(t *testing.T) (*recordtest.Fixture, string, string, []string)

// firstChunk commits chunk 1 of a pooled fixture, and stages the tests requirement when staged.
func firstChunk(staged bool) ticketCase {
	return func(t *testing.T) (*recordtest.Fixture, string, string, []string) {
		f, home := pooled(t, 1)
		base, _ := advance(f, "chunk 1")
		if staged {
			f.Stage(home, "1", "tests", f.Run(0, biting))
		}
		return f, home, base, ticketArgs(t, "1", "1.md", "p", "--base", base)
	}
}

// laterChunk freezes chunk 1, commits chunk 2, and stages its tests requirement. Its first ticket
// names the recorded chunk 1 tip with --base, and the case answers that tip as the chunk base.
func laterChunk(t *testing.T) (*recordtest.Fixture, string, string, []string) {
	f, home := pooled(t, 2)
	base := freezeChunk1(t, f, home)
	advance(f, "chunk 2")
	f.Stage(home, "2", "tests", f.Run(0, biting))
	return f, home, base, ticketArgs(t, "2", "2.md", "p", "--base", base)
}

// carriedRepairCase is the repair of carriedRepair with no stage: system and other can carry,
// and pkg cannot.
func carriedRepairCase(t *testing.T) (*recordtest.Fixture, string, string, []string) {
	f, home, frozen, _ := carriedRepair(t)
	return f, home, frozen.Base, repairArgs(t, f)
}

// BC41, BC42, BC43, BC44, BC45, BC46: with the tests requirement staged, the retained `bench
// probe` response is the entry of the additional requirement, and --mutant names its row of a
// mutant set. The call removes no other stage, prints at most 10 lines, and ends with the
// record commit. A retained `bench test` response and two unstaged requirements refuse, name
// the excerpt form, and leave the record and the store as they were. A retained run off the
// chunk line or with a dirty path refuses the same way and names the staged run. A requirement
// that can carry carries, so the retained response is the entry of the one requirement left.
// The first ticket of a later chunk prints the recorded predecessor tip as the next= base.
func TestTicketStagesTheRetainedProbe(t *testing.T) {
	produce, _ := setProducer(t, linked(t, 1))
	set, setExit := silentThenBiting(produce)
	retained := answer("bit", "yes", 0, "1", "packages[1]{package,outcome}:\n  retained,failed\n")
	for _, tc := range []struct {
		name, verb   string
		chunk        string
		output       string
		exit         int
		prepare      ticketCase
		extra        []string
		offLine      bool
		dirty        []string
		entry        string
		used         int
		carried      string
		refusalNames string
	}{
		{name: "single-form probe", verb: "probe", output: retained, prepare: firstChunk(true), entry: "p-additional", used: 1, carried: "none"},
		{name: "mutant row 2", verb: "probe", output: set, exit: setExit, prepare: firstChunk(true), extra: []string{"--mutant", "2"}, entry: "p-additional", used: 1, carried: "none"},
		{name: "first ticket of a later chunk", verb: "probe", chunk: "2", output: retained, prepare: laterChunk, entry: "p-additional", used: 1, carried: "none"},
		{name: "carried requirements", verb: "probe", output: retained, prepare: carriedRepairCase, entry: "repair-pkg", carried: "system,other"},
		{name: "bench test response", verb: "test", output: retained, prepare: firstChunk(true), refusalNames: "these owed requirements have no staged excerpt: additional"},
		{name: "two unstaged requirements", verb: "probe", output: retained, prepare: firstChunk(false), refusalNames: "these owed requirements have no staged excerpt: tests, additional"},
		{name: "head off the chunk line", verb: "probe", output: retained, prepare: firstChunk(true), offLine: true, refusalNames: "requirement additional: the staged run at "},
		{name: "dirty path", verb: "probe", output: retained, prepare: firstChunk(true), dirty: []string{"x.go"}, refusalNames: "requirement additional: the staged run had the dirty path x.go"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, home, base, args := tc.prepare(t)
			if tc.chunk == "" {
				tc.chunk = "1"
			}
			tip := f.Tip()
			f.Stage(home, "9", "other", f.Run(0, "ok\n"))
			head := tip
			if tc.offLine {
				head = f.Git("rev-parse", base+"^")
				tc.refusalNames += head
			}
			f.RetainAt(home, head, append([]string{}, tc.dirty...), tc.verb, []string{"clamp.go"}, tc.exit, tc.output)
			before, stages := recordBytes(t, f), responseboundtest.StoreFiles(t, home)
			out, code := recordcmd.Command(f.Root, append(args, tc.extra...))
			if tc.refusalNames != "" {
				if want := "; write it with bench record excerpt"; code != 1 || !strings.Contains(out, tc.refusalNames) || !strings.Contains(out, want) {
					t.Fatalf("ticket = exit %d, output %q; want exit 1 naming %q and %q", code, out, tc.refusalNames, want)
				}
				if after := recordBytes(t, f); !bytes.Equal(before, after) {
					t.Errorf("the refusal changed the record bytes:\n%s", after)
				}
				if after := responseboundtest.StoreFiles(t, home); !slices.Equal(stages, after) {
					t.Errorf("store files after the refusal = %q, want %q", after, stages)
				}
				return
			}
			if code != 0 {
				t.Fatalf("ticket = exit %d, output %q; want exit 0", code, out)
			}
			if cell := carriedCell(t, out); cell != tc.carried {
				t.Errorf("carried cell = %v, want %s", cell, tc.carried)
			}
			requireBitEntry(t, f, tc.chunk, tc.entry, tip, tc.output)
			if after := responseboundtest.StoreFiles(t, home); len(after) != len(stages)-tc.used {
				t.Errorf("store files after the write = %q, want %d used stages gone from %q", after, tc.used, stages)
			}
			lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
			next := "next=bench commit -m <msg> --preflight-build " + sanitize.ShellQuote(slug(t)) + " --base " + sanitize.ShellQuote(base) + " " + sanitize.ShellQuote(recordPath(t))
			if len(lines) > 10 || lines[len(lines)-1] != next || !strings.HasSuffix(out, "\n") {
				t.Errorf("ticket output = %q, want at most 10 lines that end with %q", out, next)
			}
		})
	}
}
