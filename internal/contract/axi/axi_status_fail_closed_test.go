package axi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

// chmodUnreadable strips all permissions from path, matching the mode-0o000 fixture
// pattern axi_maps_fail_closed_test.go establishes: skip via capability.Privilege
// instead of asserting on a host (root, some CI containers) that does not honor 0o000,
// rather than a bare t.Skip.
func chmodUnreadable(t *testing.T, path string) {
	t.Helper()
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, "cannot strip permissions: "+err.Error())
	}
	if fh, err := os.Open(path); err == nil {
		fh.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}
}

// statusUnknownFixture builds a repository on its (sole-branch) default with both of
// story 17's migrated sources unreadable: IDEAS.md (the drain signal's capture source)
// and ROADMAP.md (the roadmap-reconcile signal's source, which only fires on the
// default branch — hence the commit and single local branch). .bench/learnings.md and
// specs/ROADMAP.md's own row content are left absent, so each signal's *other* input
// reads as the ordinary quiet zero: any unknown row that shows up is attributable only
// to the unreadable file, not to unrelated content.
func statusUnknownFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.WriteFile("README.md", "# fixture\n")
	f.CommitAll("init")

	ideas := filepath.Join(f.Root, "IDEAS.md")
	f.WriteFile("IDEAS.md", "- 2026-01-01  parked\n")
	chmodUnreadable(t, ideas)

	roadmapPath := filepath.Join(f.Root, "ROADMAP.md")
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. do the thing\n")
	chmodUnreadable(t, roadmapPath)

	return f
}

// TestAXIStatusUnknownRow pins story 17: a signal whose underlying read failed renders
// an explicit `unknown` row naming the source and its classified state, for both
// migrated signals (drain's IDEAS.md, roadmap-reconcile's ROADMAP.md), and `bench
// status` still exits 0 — the dashboard's job is to render whatever it can, not to fail
// closed itself. Asserting the rows are present *and* the exit is 0 forbids both the
// fabricated-zero bug (row absent) and a fail-closed over-correction (nonzero exit).
func TestAXIStatusUnknownRow(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := statusUnknownFixture(t)

	out := f.Bench("status", "--all")

	out.RequireExit(0)
	requireContainsFold(t, out.Stdout, "unknown (IDEAS.md is unreadable)")
	requireContainsFold(t, out.Stdout, "unknown (ROADMAP.md is unreadable)")
}

// TestAXIStatusUnknownNotSuppressed pins the other half of story 17: every signal on
// the board is gated on `count > 0`, so an unknown state that fell through to that gate
// unaltered would vanish behind a fabricated zero. Here the readable side of each
// signal genuinely is zero — no open learnings, no merged/dangling roadmap rows — the
// exact condition that used to suppress the row entirely; the unknown segment must
// still render for both migrated signals at their normal severity.
func TestAXIStatusUnknownNotSuppressed(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := statusUnknownFixture(t)

	out := f.Bench("status", "--all")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  drain      unknown (IDEAS.md is unreadable), 0 open learning(s) → /bench-what-next")
	requireAXILine(t, out.Stdout, "  roadmap    unknown (ROADMAP.md is unreadable) → /bench-what-next")
}
