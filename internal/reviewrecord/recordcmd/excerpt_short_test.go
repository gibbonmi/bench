package recordcmd_test

import (
	"strings"
	"testing"

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
	requireBitEntry(t, f, "1", "s-tests", tip, set)
}
