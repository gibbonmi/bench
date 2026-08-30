package commit

import (
	"path/filepath"
	"strings"
	"testing"
)

// A dry run answers "would this composed set land green" without paying a junk
// commit: the same compose-and-authorize half a landing runs, then a full stop.
func TestDryRunReportsGreenAndMovesNothing(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "prospective\n", 0o644)
	code, stdout, stderr := runCommand(t, root, "--dry-run", "-m", "m", "a.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "nothing committed") {
		t.Fatalf("stdout = %q, want the dry-run verdict line", stdout)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatal("dry run moved HEAD")
	}
	if headHasPrefix(t, root, "a.txt") {
		t.Fatalf("dry run published the path: %v", headPaths(t, root))
	}
}

// A red composed set reports the refusal as the diagnosis and still moves nothing.
func TestDryRunRedReportsRefusalAndMovesNothing(t *testing.T) {
	root, before := landingRepo(t, 1, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "prospective\n", 0o644)
	code, stdout, stderr := runCommand(t, root, "--dry-run", "-m", "m", "a.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	want := "prospective authorization refused: inherited (the gate ran red on the composed tree and no green baseline attributes the red to this diff); run bench gate --fresh"
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want the operator-facing inherited refusal %q", stderr, want)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatal("red dry run moved HEAD")
	}
}

// A grammar error prints the one-line usage alone; the example stays a help-only cost.
func TestGrammarErrorPrintsNoExample(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, _, stderr := runCommand(t, root, "a.txt")
	if code != 2 || strings.Contains(stderr, "example:") {
		t.Fatalf("grammar error = (%d, %q), want exit 2 with no example line", code, stderr)
	}
}

// The help text advertises the flag the grammar accepts.
func TestHelpAdvertisesDryRun(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, stdout, _ := runCommand(t, root, "--help")
	if code != 0 || !strings.Contains(stdout, "--dry-run") {
		t.Fatalf("help = (%d, %q), want --dry-run advertised", code, stdout)
	}
}
