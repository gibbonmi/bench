package axi

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// statusUnknownCorruption is one way a control record stops being trustworthy while still
// existing. Both reach the status board through the same classifier state, so the unknown
// row's contract is asserted once per state rather than once per file: an implementation
// that handles only the permission case leaves the malformed case rendering a clean board.
type statusUnknownCorruption struct {
	state   string
	corrupt func(f contract.Fixture, path string)
}

var statusUnknownCorruptions = []statusUnknownCorruption{
	{"unreadable", func(f contract.Fixture, path string) { f.WriteUnreadable(path, 0o644) }},
	{"malformed", func(f contract.Fixture, path string) { f.WriteFile(path, "\xff\xfe") }},
}

// statusUnknownFixture builds a repository on its (sole-branch) default with each of story
// 17's three migrated signals unable to read its source: capture/IDEAS.md (the drain signal's
// capture source), ROADMAP.md (the roadmap-reconcile signal's source, which only fires on
// the default branch — hence the commit and single local branch), and decisions/ (the
// decision-map signal's control directory, here a regular file so the directory read is
// wrong-type). corrupt decides how the two files fail. capture/learnings.md and any spec
// ROADMAP.md would name are left absent, so each signal's *other* input reads as the
// ordinary quiet zero: any unknown row that shows up is attributable only to the broken
// source, not to unrelated content.
func statusUnknownFixture(t *testing.T, c statusUnknownCorruption) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "# fixture\n")
	f.CommitAll("init")

	f.WriteFile("capture/IDEAS.md", "- 2026-01-01  parked\n")
	c.corrupt(f, "capture/IDEAS.md")

	f.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. do the thing\n")
	c.corrupt(f, "ROADMAP.md")

	f.WriteFile("decisions", "not a directory\n")

	return f
}

// TestAXIStatusUnknownRow pins story 17: a signal whose underlying read failed renders an
// explicit `unknown` row naming the source and its classified state, for all three migrated
// signals (drain's capture/IDEAS.md, roadmap-reconcile's ROADMAP.md, the decision-map count's
// decisions/), and `bench status` still exits 0 — the dashboard's job is to render whatever
// it can, not to fail closed itself. Asserting the rows are present *and* the exit is 0
// forbids both the fabricated-zero bug (row absent) and a fail-closed over-correction
// (nonzero exit).
func TestAXIStatusUnknownRow(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	for _, c := range statusUnknownCorruptions {
		contract.RunParallel(t, c.state, func(t *testing.T) {
			f := statusUnknownFixture(t, c)

			out := f.Bench("status", "--all")

			out.RequireExit(0)
			requireContainsFold(t, out.Stdout, "unknown (capture/IDEAS.md is "+c.state+")")
			requireContainsFold(t, out.Stdout, "unknown (ROADMAP.md is "+c.state+")")
			requireContainsFold(t, out.Stdout, "unknown (decisions is wrong-type)")
		})
	}
}

// TestAXIStatusUnknownNotSuppressed pins the other half of story 17: every signal on the
// board is gated on `count > 0`, so an unknown state that fell through to that gate
// unaltered would vanish behind a fabricated zero. Here the readable side of each signal
// genuinely is zero — no open learnings, no merged/dangling roadmap rows, no unresolved
// maps — the exact condition that used to suppress the row entirely; the unknown segment
// must still render for every migrated signal at its normal severity.
func TestAXIStatusUnknownNotSuppressed(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	for _, c := range statusUnknownCorruptions {
		contract.RunParallel(t, c.state, func(t *testing.T) {
			f := statusUnknownFixture(t, c)

			out := f.Bench("status", "--all")

			out.RequireExit(0)
			requireAXILine(t, out.Stdout, "  drain      unknown (capture/IDEAS.md is "+c.state+"), 0 open learning(s), 0 pending retro(s) → /bench-what-next")
			requireAXILine(t, out.Stdout, "  roadmap    unknown (ROADMAP.md is "+c.state+") → /bench-what-next")
			requireAXILine(t, out.Stdout, "  decisions  unknown (decisions is wrong-type) → investigate decisions/ (bench maps)")
		})
	}
}

// TestAXIStatusSpecialFileDoesNotBlock drives a FIFO where a spec file belongs. The
// housekeeping counters walk `specs/*/spec.md` and read each entry, and a FIFO with no writer
// never yields EOF, so a counter that opens before checking the type blocks in open(2)
// forever — the SessionStart hook runs `bench status`, which makes a hang worse than any
// wrong count. The failure is an expired deadline rather than a wrong answer, so the
// deadline is what carries the assertion.
func TestAXIStatusSpecialFileDoesNotBlock(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "# fixture\n")
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n- FT1 — specs/hang/spec.md\n\n## Recommended sequence\n\n1. do the thing\n")
	f.CommitAll("init")
	f.WriteFifo("specs/hang/spec.md")

	out := f.BenchDeadlined("status", "--all")

	if out.TimedOut {
		t.Fatal("bench status blocked on a FIFO at specs/hang/spec.md, so a spec read opened the path before checking its type")
	}
	out.RequireExit(0)
}
