package axi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

// TestAXILearningsFailClosed pins story 8: an unreadable journal exits 1 with a
// structured `error:` line naming the state. The current implementation swallows the
// read error and prints an empty table at exit 0, which the exit assertion rejects.
func TestAXILearningsFailClosed(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — first learning  [open]\n")
	path := filepath.Join(f.Root, ".bench", "learnings.md")
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, "cannot strip permissions: "+err.Error())
	}
	if fh, err := os.Open(path); err == nil {
		fh.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}

	out := f.Bench("learnings")

	out.RequireExit(1)
	if out.Stdout == "" {
		t.Fatalf("bench learnings printed nothing on an unreadable journal")
	}
	requireContainsFold(t, out.Stdout, "error:")
	requireContainsFold(t, out.Stdout, "unreadable")
}

// TestAXILearningsAbsentIsEmpty pairs with TestAXILearningsFailClosed: a repository with
// no journal at all still exits 0 and renders the definitive empty table. This forbids
// the over-correction that fails closed on absence too — absence is the one
// authoritative empty state.
func TestAXILearningsAbsentIsEmpty(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)

	out := f.Bench("learnings")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "learnings[0]{date,title}:")
}

// TestAXILearningsMalformedRows pins story 9: a journal mixing one well-formed entry
// with one malformed heading renders both — the good row plus an explicit row naming
// the malformed heading's line and reason — and exits 1. This forbids a hard fail that
// drops the good entries along with the bad one. The fixture is written without a
// trailing newline, the hand-edited-file class the spec's edge inventory assigns to
// this parser.
func TestAXILearningsMalformedRows(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — good entry  [open]\n## bad heading with no date")

	out := f.Bench("learnings")

	out.RequireExit(1)
	requireAXILine(t, out.Stdout, "  2026-01-01,good entry")
	requireContainsFold(t, out.Stdout, "line 2")
	requireContainsFold(t, out.Stdout, "malformed learning heading")
}
