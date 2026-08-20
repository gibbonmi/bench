package commit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// closeSlug carries a space and a glob character: the close step resolves an exact
// path, so a folder name a shell would expand must land unchanged.
const closeSlug = "ft900 close [x]*"

// landingRepo builds the minimal repository `bench commit` lands into: a gate with the
// named exit code, one staged tracked change to attribute, and whatever specs/ shape
// write leaves behind. It returns the root and the pre-landing HEAD.
func landingRepo(t *testing.T, gateExit int, write func(t *testing.T, root string)) (root, before string) {
	t.Helper()
	root = t.TempDir()
	git := func(args ...string) { t.Helper(); runGit(t, root, args...) }
	git("init", "-q", "-b", "main")
	git("config", "user.email", "a@b.c")
	git("config", "user.name", "a")
	mustMkdirAll(t, filepath.Join(root, ".bench"))
	mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), "#!/bin/sh\nexit "+strconv.Itoa(gateExit)+"\n", 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-inputs.json"), `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	mustWrite(t, filepath.Join(root, "tracked.txt"), "base\n", 0o644)
	write(t, root)
	git("add", "-A")
	git("commit", "-qm", "base")
	mustWrite(t, filepath.Join(root, "tracked.txt"), "changed\n", 0o644)
	git("add", "tracked.txt")
	return root, strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
}

// runCommand invokes the real command from root, the way the CLI does.
func runCommand(t *testing.T, root string, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	var out, errOut bytes.Buffer
	code = Command(args, &out, &errOut)
	return code, out.String(), errOut.String()
}

// headPaths lists every path the commit at HEAD tracks.
func headPaths(t *testing.T, root string) []string {
	t.Helper()
	raw := runGit(t, root, "ls-tree", "-r", "--name-only", "-z", "HEAD")
	var paths []string
	for _, p := range strings.Split(string(raw), "\x00") {
		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}

func headHasPrefix(t *testing.T, root, prefix string) bool {
	t.Helper()
	for _, p := range headPaths(t, root) {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}

// writeTicketsOnly writes a tickets-only folder: a direct child of specs/ carrying
// tickets and no spec.md.
func writeTicketsOnly(t *testing.T, root, slug string) {
	t.Helper()
	mustMkdirAll(t, filepath.Join(root, "specs", slug, "tickets"))
	mustWrite(t, filepath.Join(root, "specs", slug, "tickets", "one.md"), "# One\n", 0o644)
}

// H01: a green landing on a tickets-only slug removes the folder in the same commit.
// Asserted in the tracked configuration — an untracked folder is absent from the
// published tree either way, so it could not distinguish a close step from no close step.
func TestGreenLandingOnTicketsOnlySlugDeletesFolderInItsCommit(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) { writeTicketsOnly(t, root, closeSlug) })
	if !headHasPrefix(t, root, "specs/"+closeSlug+"/") {
		t.Fatal("fixture is not tracked; the assertion below could not distinguish a close step")
	}
	code, stdout, stderr := runCommand(t, root, "-m", "land", "--spec", closeSlug, "tracked.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after == before {
		t.Fatal("HEAD did not move")
	}
	if headHasPrefix(t, root, "specs/"+closeSlug+"/") {
		t.Fatalf("published commit still tracks the tickets-only folder: %v", headPaths(t, root))
	}
	if _, err := os.Stat(filepath.Join(root, "specs", closeSlug)); !os.IsNotExist(err) {
		t.Fatalf("tickets-only folder still present in the worktree: %v", err)
	}
}

// H02: a slug whose folder holds spec.md takes the status flip and keeps its folder.
// --spec consults no run state: a green landing on the named spec's tracked path flips
// its one `Status: staged` line, keeps the folder, and moves HEAD.
func TestGreenLandingOnSpecBackedSlugFlipsStatusAndKeepsFolder(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {
		mustMkdirAll(t, filepath.Join(root, "specs", "active"))
		mustWrite(t, filepath.Join(root, "specs", "active", "spec.md"), "# Active\n\nStatus: staged\n", 0o644)
	})
	code, stdout, stderr := runCommand(t, root, "-m", "land", "--spec", "active", "tracked.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after == before {
		t.Fatal("HEAD did not move")
	}
	got := string(mustReadFile(t, filepath.Join(root, "specs", "active", "spec.md")))
	if strings.Contains(got, "Status: staged") || !strings.Contains(got, "Status: implemented") {
		t.Fatalf("spec not flipped: %q", got)
	}
	if !headHasPrefix(t, root, "specs/active/") {
		t.Fatalf("published commit dropped the spec-backed folder: %v", headPaths(t, root))
	}
	if _, err := os.Stat(filepath.Join(root, "specs", "active", "spec.md")); err != nil {
		t.Fatalf("spec-backed folder removed from the worktree: %v", err)
	}
}

// H03: a slug naming no folder returns a structured error and lands nothing.
func TestLandingOnAbsentSlugReportsAndLandsNothing(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) { writeTicketsOnly(t, root, closeSlug) })
	code, stdout, stderr := runCommand(t, root, "-m", "land", "--spec", "no-such-slug", "tracked.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.HasPrefix(stderr, "error: ") || !strings.Contains(stderr, "no-such-slug") {
		t.Fatalf("stderr = %q, want a structured error naming the slug", stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatalf("HEAD moved from %s to %s", before, after)
	}
	if !headHasPrefix(t, root, "specs/"+closeSlug+"/") {
		t.Fatal("an unrelated tickets-only folder was closed by a slug that names no folder")
	}
}

// H04: a red gate leaves the tickets-only folder present.
func TestRedGateLeavesTicketsOnlyFolderPresent(t *testing.T) {
	root, before := landingRepo(t, 1, func(t *testing.T, root string) { writeTicketsOnly(t, root, closeSlug) })
	code, stdout, stderr := runCommand(t, root, "-m", "land", "--spec", closeSlug, "tracked.txt")
	if code == 0 {
		t.Fatalf("exit = 0 on a red gate; stdout=%q stderr=%q", stdout, stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatalf("HEAD moved on a red gate: %s -> %s", before, after)
	}
	if _, err := os.Stat(filepath.Join(root, "specs", closeSlug, "tickets", "one.md")); err != nil {
		t.Fatalf("red gate discarded the tickets-only folder: %v", err)
	}
	if !headHasPrefix(t, root, "specs/"+closeSlug+"/") {
		t.Fatal("red gate removed the tickets-only folder from the tree")
	}
}

func runGit(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return out
}
func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(fmt.Errorf("read %s: %w", path, err))
	}
	return out
}
func mustMkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}
func mustWrite(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}
