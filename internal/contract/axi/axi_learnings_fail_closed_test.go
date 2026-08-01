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
	f.WriteFile("capture/learnings.md", "## 2026-01-01 — first learning  [open]\n")
	path := filepath.Join(f.Root, "capture", "learnings.md")
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

func TestAXILearningsPresentNonJournalIsNotEmpty(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	for _, tc := range []struct {
		name, content, state string
	}{
		{"zero byte", "", "empty"},
		{"arbitrary prose", "a note without journal headings\n", "unsupported-schema"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := contract.NewFixture(t)
			f.WriteFile("capture/learnings.md", tc.content)

			out := f.Bench("learnings")

			out.RequireExit(1)
			requireContainsFold(t, out.Stdout, "error: capture/learnings.md is "+tc.state)
		})
	}
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
	f.WriteFile("capture/learnings.md", "## 2026-01-01 — good entry  [open]\n## bad heading with no date")

	out := f.Bench("learnings")

	out.RequireExit(1)
	requireAXILine(t, out.Stdout, "  2026-01-01,good entry")
	requireContainsFold(t, out.Stdout, "line 2")
	requireContainsFold(t, out.Stdout, "malformed learning heading")
}

func TestAXILearningsDatedStateRequired(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	for _, tc := range []struct {
		name    string
		heading string
	}{
		{"missing state", "## 2026-01-01 — missing state"},
		{"non-open state", "## 2026-01-01 — resolved  [resolved]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := contract.NewFixture(t)
			f.WriteFile("capture/learnings.md", "## 2026-01-02 — good entry  [open]\n"+tc.heading)

			out := f.Bench("learnings")

			out.RequireExit(1)
			requireAXILine(t, out.Stdout, "  2026-01-02,good entry")
			requireContainsFold(t, out.Stdout, "line 2")
			requireContainsFold(t, out.Stdout, "must end with [open]")
		})
	}
}

func TestAXILearningsMalformedStateKeepsDrainVisible(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. first\n2. second\n")
	f.WriteFile("capture/learnings.md", "## 2026-01-01 — missing state\n")

	roadmap := f.Bench("roadmap")
	roadmap.RequireExit(0)
	requireContainsFold(t, roadmap.Stdout, "unknown (capture/learnings.md is malformed)")

	status := f.Bench("status", "--all")
	status.RequireExit(0)
	requireContainsFold(t, status.Stdout, "unknown (capture/learnings.md is malformed)")
}

// TestAXILearningsWrongType drives the wrong-type state through two real commands. It
// needs nothing of the host — a directory parked at the journal's path is wrong-type on
// every filesystem — which is what makes it the surface-level assertion the unreadable
// and socket fixtures cannot supply on a box without them.
//
// The two surfaces are asserted against one fixture on purpose: `bench learnings` fails
// closed with the `<path> is <state>` error line while `bench status` degrades to the
// `unknown (<path> is <state>)` cell built from the same clause, so a change that moved
// one grammar without the other breaks one of these two assertions.
func TestAXILearningsWrongType(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. first\n2. second\n")
	if err := os.MkdirAll(filepath.Join(f.Root, "capture", "learnings.md"), 0o755); err != nil {
		t.Fatalf("park a directory at the journal path: %v", err)
	}

	out := f.Bench("learnings")
	out.RequireExit(1)
	requireContainsFold(t, out.Stdout, "error: capture/learnings.md is wrong-type")

	board := f.Bench("status", "--all")
	board.RequireExit(0)
	requireContainsFold(t, board.Stdout, "unknown (capture/learnings.md is wrong-type)")

	roadmap := f.Bench("roadmap")
	roadmap.RequireExit(0)
	requireContainsFold(t, roadmap.Stdout, "unknown (capture/learnings.md is wrong-type)")
}
