package axi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

// TestAXIMapsUnsupportedSchema pins the maps half of story 7: a decisions file that
// attempts neither a `## #<n>:` ticket heading nor a marker is unsupported-schema —
// distinct from a decisions file whose bytes are invalid UTF-8, which is malformed —
// asserted together so collapsing either onto the other or onto a generic error fails
// one of the two rows. The no-heading fixture is written without a trailing newline,
// the hand-edited-file class this parser owns.
func TestAXIMapsUnsupportedSchema(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/noheading.md", "# Notes\nJust prose here, no ticket heading and no marker.")
	f.WriteFile("decisions/badutf8.md", "## #1: q?\nType: Grill\n### Answer\n\xff\xfe— (open)\n")

	out := f.Bench("maps")

	out.RequireExit(1)
	requireAXILine(t, out.Stdout, "  noheading,error,unsupported-schema,no ticket heading or marker found")
	requireContainsFold(t, out.Stdout, "badutf8,error,malformed")
	requireContainsFold(t, out.Stdout, "invalid utf-8")
}

// TestAXIMapsUnreadableFileRow pins story 10: a decision file this process cannot read
// gets a row naming it and its state, alongside every other readable file's rows, and
// the command still exits 1. The prior implementation `continue`s past the read error
// and silently drops the file from both the listing and the count.
func TestAXIMapsUnreadableFileRow(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/good.md", "## #1: q?\nType: Grill\n### Answer\n— (open)\n")
	f.WriteFile("decisions/locked.md", "## #1: q?\nType: Grill\n### Answer\n— (open)\n")
	path := filepath.Join(f.Root, "decisions", "locked.md")
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, "cannot strip permissions: "+err.Error())
	}
	if fh, err := os.Open(path); err == nil {
		fh.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}

	out := f.Bench("maps")

	out.RequireExit(1)
	requireAXILine(t, out.Stdout, "  good,1,Grill,open")
	requireContainsFold(t, out.Stdout, "locked,error,unreadable")
}

// TestAXIMapsAbsentIsEmpty pairs with the fail-closed rows above: a repository with no
// decisions/ directory at all still exits 0 and renders the definitive empty table —
// absence is the one authoritative empty state, and an unconditional exit-1 stub would
// satisfy the rows above while destroying this contract.
func TestAXIMapsAbsentIsEmpty(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "maps[0]{map,ticket,type,state}:")
}

// TestAXIMapsCountMatchesListing pins story 11: the status dashboard's decision-map
// count and the `bench maps` listing derive from the same scan, driven against one
// fixture with an unreadable map file. A count that fabricated zero for the failed
// scan would disagree with a listing that shows the row.
func TestAXIMapsCountMatchesListing(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/good.md", "## #1: q?\nType: Grill\n### Answer\n— (open)\n")
	f.WriteFile("decisions/locked.md", "## #1: q?\nType: Grill\n### Answer\n— (open)\n")
	path := filepath.Join(f.Root, "decisions", "locked.md")
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, "cannot strip permissions: "+err.Error())
	}
	if fh, err := os.Open(path); err == nil {
		fh.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}

	count := f.Bench("maps", "--count")
	count.RequireExit(0)
	if got := strings.TrimSpace(count.Stdout); got != "2" {
		t.Fatalf("maps --count = %q, want 2\nstdout:\n%s", got, count.Stdout)
	}

	rows := f.Bench("maps")
	rows.RequireExit(1)
	requireAXILine(t, rows.Stdout, "  good,1,Grill,open")
	requireContainsFold(t, rows.Stdout, "locked,error,unreadable")

	status := f.Bench("status", "--all")
	status.RequireExit(0)
	requireContainsFold(t, status.Stdout, "decisions")
	requireContainsFold(t, status.Stdout, "2 unresolved map")
}
