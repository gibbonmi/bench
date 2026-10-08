// Refusal route tests for the commit: each commit refusal prints the route of its face in
// the shared refusal-route registry, and that route clears the refusal.
package commit

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// TestPublishedUnreconciledRouteIsTheResetPlan is RR32 through RR34, the collision 8 restore
// repro. A commit exit 3 names the reset plan at the published commit for the commit's own
// root, and never the restore that the destructive-git guard denies. A root that is not
// line-safe prints the checkout placeholder, so no control byte reaches the record.
func TestPublishedUnreconciledRouteIsTheResetPlan(t *testing.T) {
	// requireResetPlan requires the reset plan at the published commit for checkout, the
	// root as the route prints it, in the next= field of the exit 3 record.
	requireResetPlan := func(t *testing.T, stdout, published, checkout string) {
		t.Helper()
		_, fields, _ := recordFields(t, stdout)
		next := fields[refusalroute.NextField]
		if want := "bench worktree reset --to " + sanitize.ShellQuote(published) + " " + checkout; next != want {
			t.Errorf("next = %q, want %q", next, want)
		}
		if strings.Contains(next, "git restore") {
			t.Errorf("next = %q, want no restore step", next)
		}
	}
	t.Run("line-safe root", func(t *testing.T) {
		root, stdout, published := publishedUnreconciled(t, "named")
		requireResetPlan(t, stdout, published, sanitize.ShellQuote(root))
	})
	// A test's first TempDir call makes the one parent in TMPDIR that all its temporary
	// directories share. So this test names the unsafe directory, and the subtest sets
	// TMPDIR to it before its first call: the commit's root then sits under the ESC byte.
	unsafe := filepath.Join(t.TempDir(), "esc\x1bdir")
	t.Run("root that is not line-safe", func(t *testing.T) {
		mustMkdirAll(t, unsafe)
		t.Setenv("TMPDIR", unsafe)
		root, stdout, published := publishedUnreconciled(t, "named")
		if sanitize.LineSafe(root) {
			t.Fatalf("root = %q, want a root that is not line-safe", root)
		}
		requireResetPlan(t, stdout, published, "<checkout>")
	})
}
