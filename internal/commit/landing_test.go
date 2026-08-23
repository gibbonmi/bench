package commit

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// ticketsSlug carries a space and a glob character: a named path resolves literally,
// so a folder name a shell would expand must land unchanged.
const ticketsSlug = "ft900 tickets [x]*"

// landingRepo builds the minimal repository `bench commit` lands into: a gate with the
// named exit code, one staged tracked change to attribute, and whatever tree shape
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

// FA1: the retired --spec is an unknown flag. The grammar refuses it before any repo
// read, so the command reports the usage line and no ref moves.
func TestSpecFlagIsAGrammarError(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, stdout, stderr := runCommand(t, root, "-m", "m", "--spec", "x", "tracked.txt")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.HasPrefix(stderr, "usage: bench commit") || !strings.Contains(stderr, "--spec") {
		t.Fatalf("stderr = %q, want the usage line naming the unknown argument", stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatalf("HEAD moved from %s to %s", before, after)
	}
}

// FA2: the help text advertises the grammar the parser accepts, so it names no --spec.
func TestHelpAdvertisesNoSpecFlag(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, stdout, stderr := runCommand(t, root, "--help")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%q", code, stderr)
	}
	if !strings.HasPrefix(stdout, "usage: bench commit") {
		t.Fatalf("stdout = %q, want the help text", stdout)
	}
	if strings.Contains(stdout, "--spec") {
		t.Fatalf("help still advertises --spec: %q", stdout)
	}
}

// FA7: the plain path route publishes the one named path on a green gate.
func TestGreenLandingPublishesTheNamedPath(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)
	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after == before {
		t.Fatal("HEAD did not move")
	}
	if !headHasPrefix(t, root, "a.txt") {
		t.Fatalf("published commit does not track a.txt: %v", headPaths(t, root))
	}
}

// FA8: specs/<slug> is a named path like any other. The folder a landing once closed
// now publishes its files and reconciles, so the checkout is clean afterwards.
func TestGreenLandingPublishesATicketsOnlyFolderNamedAsAPath(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	writeTicketsOnly(t, root, ticketsSlug)
	code, stdout, stderr := runCommand(t, root, "-m", "m", filepath.Join("specs", ticketsSlug))
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after == before {
		t.Fatal("HEAD did not move")
	}
	if !headHasPrefix(t, root, "specs/"+ticketsSlug+"/tickets/one.md") {
		t.Fatalf("published commit does not track the folder's files: %v", headPaths(t, root))
	}
	if status := string(runGit(t, root, "status", "--porcelain")); status != "" {
		t.Fatalf("checkout is not clean: %q", status)
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
