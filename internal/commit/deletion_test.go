package commit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The help text names the deletion route, so an author deletes a tracked path through
// the verb rather than through a raw git rm.
func TestHelpNamesTheDeletionRoute(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, stdout, stderr := runCommand(t, root, "--help")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%q", code, stderr)
	}
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, "deleted") && strings.Contains(line, "deletion") && strings.Contains(line, "git rm") {
			return
		}
	}
	t.Fatalf("help = %q, want one line naming a deleted path, the deletion it commits, and git rm", stdout)
}

// A named path absent from the worktree commits as a deletion, for a file and for a
// folder, and the checkout reconciles clean.
func TestNamedDeletedPathsPublishAsDeletions(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {
		mustWrite(t, filepath.Join(root, "gone.txt"), "gone\n", 0o644)
		namedDir(t, root, "gone dir")
	})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	if err := os.Remove(filepath.Join(root, "gone.txt")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, "gone dir")); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runCommand(t, root, "-m", "m", "gone.txt", "gone dir")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after == before {
		t.Fatal("HEAD did not move")
	}
	if headHasPrefix(t, root, "gone.txt") || headHasPrefix(t, root, "gone dir/") {
		t.Fatalf("published commit still tracks a deleted path: %v", headPaths(t, root))
	}
	if !headHasPrefix(t, root, "tracked.txt") {
		t.Fatalf("published commit lost an unnamed path: %v", headPaths(t, root))
	}
	if status := string(runGit(t, root, "status", "--porcelain")); status != "" {
		t.Fatalf("checkout is not clean: %q", status)
	}
}
