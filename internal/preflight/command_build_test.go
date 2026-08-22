package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCommandBuildFresh is B1 (the fresh-build contract test). With no
// tickets/ directory at all, `build` prints a not-applicable row for each
// of rows-owned, rows-membership, and diff-nonempty, individually
// asserted. It runs the rest for real, and exits 0 when they are green.
func TestCommandBuildFresh(t *testing.T) {
	_, slug := seedBuildFresh(t)

	out, code := Command([]string{"build", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	for _, name := range []string{"base-current,green", "paths-authorized,green"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row:\n%s", name, out)
		}
	}
	for _, name := range []string{"rows-owned,not-applicable", "rows-membership,not-applicable", "diff-nonempty,not-applicable"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row:\n%s", name, out)
		}
	}
}

// TestCommandBuildResumedTicketsRunForReal is B2's present-tickets half (the
// resumed-build contract test). With a present tickets/ directory citing
// every declared row, `build` runs rows-owned and rows-membership for real
// — green, not not-applicable — while diff-nonempty stays not-applicable.
func TestCommandBuildResumedTicketsRunForReal(t *testing.T) {
	_, slug := seedConformant(t) // seedConformant's tickets/one.md cites PF1 and PF2

	out, code := Command([]string{"build", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	for _, name := range []string{"rows-owned,green", "rows-membership,green"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row (present tickets/ must run the check for real):\n%s", name, out)
		}
	}
	if !strings.Contains(out, "diff-nonempty,not-applicable") {
		t.Errorf("output missing diff-nonempty,not-applicable row:\n%s", out)
	}
}

// TestCommandBuildEmptyTicketsRed is B2's empty-tickets half (the empty-tickets
// contract test): a present-but-empty tickets/ directory is red — declared rows
// unowned — rather than not-applicable.
func TestCommandBuildEmptyTicketsRed(t *testing.T) {
	_, slug := seedBuildFresh(t)
	if err := os.MkdirAll(filepath.Join("specs", slug, "tickets"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "rows-owned,red") {
		t.Errorf("output missing red rows-owned for a present-but-empty tickets/ dir:\n%s", out)
	}
}

// TestCommandBuildStaleBaseRedDespiteNA is B3 (the stale-base build contract test):
// base-current red in build mode exits 1 even while the ticket checks and
// diff-nonempty are not-applicable.
func TestCommandBuildStaleBaseRedDespiteNA(t *testing.T) {
	_, slug := seedBuildFresh(t)
	runGit(t, "checkout", "-q", "main")
	mustWriteFile(t, "unrelated.txt", "advance main\n")
	runGit(t, "add", "unrelated.txt")
	runGit(t, "commit", "-q", "-m", "advance main")
	runGit(t, "checkout", "-q", "feature")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "base-current,red") {
		t.Errorf("output missing red base-current row:\n%s", out)
	}
	for _, name := range []string{"rows-owned,not-applicable", "rows-membership,not-applicable", "diff-nonempty,not-applicable"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing %s row:\n%s", name, out)
		}
	}
}

// TestCommandBuildOutOfFenceRed is B4 (the build out-of-fence contract test): a
// tracked change outside every fence entry makes paths-authorized red in build mode,
// exit 1.
func TestCommandBuildOutOfFenceRed(t *testing.T) {
	_, slug := seedBuildFresh(t)
	mustWriteFile(t, "unfenced/other.go", "package other\n")
	runGit(t, "add", "unfenced/other.go")
	runGit(t, "commit", "-q", "-m", "out of fence")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,red") || !strings.Contains(out, "unfenced/other.go") {
		t.Errorf("output missing red paths-authorized naming unfenced/other.go:\n%s", out)
	}
}
