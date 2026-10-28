package recordcmd_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// BC56: the answer of a set whose one row bit is short: it ends at the execution row with
// trailing_bytes=0. `bench record excerpt --mutant 1` stages it, and `bench record ticket` then
// writes a probe entry that bit and restored, with the exact answer bytes.
func TestExcerptStagesAShortMutantSet(t *testing.T) {
	set := setAnswer(t)
	if strings.Contains(set, "mutant_report") || !strings.HasSuffix(set, ",0,0\n") {
		t.Fatalf("set answer = %q, want the short answer with trailing_bytes=0 and no report", set)
	}
	f, home := recordtest.Staging(t)
	base, tip := advance(f, "chunk 1")
	stageOwed(f, home, "1", "1.md")
	f.Retain(home, "probe", []string{"clamp.go", "--mutants", "m.json"}, 0, set)
	succeed(t, f, stageArgs(t, "tests", "--mutant", "1"))
	succeed(t, f, ticketArgs(t, "1", "1.md", "s", "--base", base))
	results := chunkEntry(t, f, "1").Verification
	i := slices.IndexFunc(results, func(entry rr.Verification) bool { return entry.ID == "s-tests" })
	if i < 0 {
		t.Fatalf("entries = %q, want s-tests", verificationIDs(rr.Chunk{Verification: results}))
	}
	got := results[i]
	want := rr.Probe{Outcome: "bit", ExitCode: 1, Restore: "pass", Mutation: f.Plan.Chunk("1").Verification[0].Probe}
	if got.Probe != nil {
		want.NativeRef = got.Probe.NativeRef
	}
	if got.Probe == nil || !reflect.DeepEqual(*got.Probe, want) || *got.ExitCode != 0 || got.SourceCommit != tip || got.NativeRef.Excerpt != set {
		t.Fatalf("entry s-tests = probe %+v, exit %d, at %s; want %+v, exit 0, at %s, and the exact answer bytes", got.Probe, *got.ExitCode, got.SourceCommit, want, tip)
	}
}
