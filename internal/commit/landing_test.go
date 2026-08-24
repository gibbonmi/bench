package commit

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/sanitize"
)

// ticketsSlug carries a space and a glob character: a named path resolves literally,
// so a folder name a shell would expand must land unchanged.
const ticketsSlug = "ft900 tickets [x]*"

// landingRepo builds the minimal linked worktree `bench commit` lands into. It returns
// that worktree and its pre-landing HEAD.
func landingRepo(t *testing.T, gateExit int, write func(t *testing.T, root string)) (root, before string) {
	t.Helper()
	primary := t.TempDir()
	initializeLandingRepo(t, primary, gateExit)
	root = filepath.Join(t.TempDir(), "linked")
	runGit(t, primary, "worktree", "add", "-q", "-b", "topic", root)
	return prepareLandingCheckout(t, root, write)
}

func primaryLandingRepo(t *testing.T, gateExit int, write func(t *testing.T, root string)) (root, before string) {
	t.Helper()
	root = t.TempDir()
	initializeLandingRepo(t, root, gateExit)
	return prepareLandingCheckout(t, root, write)
}

func initializeLandingRepo(t *testing.T, root string, gateExit int) {
	t.Helper()
	git := func(args ...string) { t.Helper(); runGit(t, root, args...) }
	git("init", "-q", "-b", "main")
	git("config", "user.email", "a@b.c")
	git("config", "user.name", "a")
	mustMkdirAll(t, filepath.Join(root, ".bench"))
	mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), "#!/bin/sh\nexit "+strconv.Itoa(gateExit)+"\n", 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-inputs.json"), `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	mustWrite(t, filepath.Join(root, "tracked.txt"), "base\n", 0o644)
	git("add", "-A")
	git("commit", "-qm", "bootstrap")
}

func prepareLandingCheckout(t *testing.T, root string, write func(t *testing.T, root string)) (string, string) {
	t.Helper()
	write(t, root)
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "--allow-empty", "-qm", "base")
	mustWrite(t, filepath.Join(root, "tracked.txt"), "changed\n", 0o644)
	runGit(t, root, "add", "tracked.txt")
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
	if !strings.Contains(stdout, "example: bench commit -m") || !strings.Contains(stdout, " -- ") {
		t.Fatalf("stdout = %q, want the first-line example with the trailing -- <path>... form", stdout)
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

func TestPrimaryCheckoutRefusesBeforePublication(t *testing.T) {
	root, before := primaryLandingRepo(t, 0, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "must stay uncommitted\n", 0o644)

	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "bench worktree create") {
		t.Fatalf("stderr = %q, want the worktree creation action", stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatalf("HEAD moved from %s to %s", before, after)
	}
}

// FA8: specs/<slug> is a named path like any other. A tickets-only folder named as
// a path publishes its files and reconciles, so the checkout is clean afterwards.
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

// unreconcilableGate writes the stub gate the exit-3 fixtures use. It plants a nested
// repository under each named directory: `git clean -f -d` skips a nested repository, so
// the residue survives the reconcile's clean and reaches its untracked-content check.
// The gate runs during authorization, before the publication, so the checkout is already
// unreconcilable when the reconcile reaches it. The absolute root is baked into the
// script because the gate runs from a prospective checkout, not from this repository.
func unreconcilableGate(dirs ...string) func(t *testing.T, root string) {
	return func(t *testing.T, root string) {
		t.Helper()
		script := "#!/bin/sh\n"
		for _, dir := range dirs {
			nested := sanitize.ShellQuote(filepath.Join(root, dir, "nested"))
			script += "mkdir -p " + nested + " && git init -q " + nested + " && printf residue > " + nested + "/r.txt\n"
		}
		mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), script+"exit 0\n", 0o755)
	}
}

// namedDir creates one named path as a directory holding one tracked-to-be file.
func namedDir(t *testing.T, root, dir string) {
	t.Helper()
	mustMkdirAll(t, filepath.Join(root, dir))
	mustWrite(t, filepath.Join(root, dir, "f.txt"), "content\n", 0o644)
}

// recordFields parses one `name{key=value,...}` record into its ordered fields. The
// values carry no comma, because the sanitizer and the shell quoting both leave the
// separator alone, so a plain split is exact here.
func recordFields(t *testing.T, stdout string) (name string, fields map[string]string, order []string) {
	t.Helper()
	line := strings.TrimSuffix(stdout, "\n")
	open := strings.Index(line, "{")
	if open < 0 || !strings.HasSuffix(line, "}") {
		t.Fatalf("stdout is not one record: %q", stdout)
	}
	name = line[:open]
	fields = map[string]string{}
	for _, field := range strings.Split(line[open+1:len(line)-1], ",") {
		key, value, ok := strings.Cut(field, "=")
		if !ok {
			t.Fatalf("record field %q is not key=value", field)
		}
		fields[key] = value
		order = append(order, key)
	}
	return name, fields, order
}

// FB1, FB2, FB9: the production reconcile, with no injected function, fails on a named
// path after the ref update succeeded. The command exits 3, reports the published commit
// that is now HEAD, and names the path whose reconcile failed.
func TestPublishedButUnreconciledCommitExitsThreeAndNamesTheFailedPath(t *testing.T) {
	root, before := landingRepo(t, 0, unreconcilableGate("named"))
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	namedDir(t, root, "named")

	code, stdout, stderr := runCommand(t, root, "-m", "m", "named")
	if code != 3 {
		t.Fatalf("exit = %d, want 3; stdout=%q stderr=%q", code, stdout, stderr)
	}
	name, fields, order := recordFields(t, stdout)
	if name != "committed" || !reflect.DeepEqual(order, []string{"published_commit", "path", "next"}) {
		t.Fatalf("record = %q, want committed{published_commit=…,path=…,next=…}", stdout)
	}
	head := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
	if head == before {
		t.Fatal("HEAD did not move: the commit was not published")
	}
	if fields["published_commit"] != head {
		t.Fatalf("published_commit = %q, want the new HEAD %q", fields["published_commit"], head)
	}
	if fields["path"] != "named" {
		t.Fatalf("path = %q, want the failed named path", fields["path"])
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want the record on stdout alone", stderr)
	}
}

// FB4: next= is one restore over every named path, shell-quoted, in the owner's sorted
// order rather than the argv order. One path carries a space, so an unquoted render
// would break the paste.
func TestPublicationRemainderNextRestoresEveryNamedPathShellQuoted(t *testing.T) {
	root, _ := landingRepo(t, 0, unreconcilableGate("zzz"))
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	namedDir(t, root, "zzz")
	namedDir(t, root, "one dir")

	code, stdout, stderr := runCommand(t, root, "-m", "m", "zzz", "one dir")
	if code != 3 {
		t.Fatalf("exit = %d, want 3; stdout=%q stderr=%q", code, stdout, stderr)
	}
	_, fields, _ := recordFields(t, stdout)
	head := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
	want := "git restore --source=" + head + " --staged --worktree -- 'one dir' 'zzz'"
	if fields["next"] != want {
		t.Fatalf("next = %q, want %q", fields["next"], want)
	}
}

// FB7: a failed path carrying an ESC byte renders in path= as the sanitizer spells it,
// and next= takes the placeholder rather than a command that would emit the raw byte.
func TestPublicationRemainderSanitizesTheFailedPathAndPointsAtNamedPaths(t *testing.T) {
	const escaped = "esc\x1bdir"
	root, _ := landingRepo(t, 0, unreconcilableGate(escaped))
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	namedDir(t, root, escaped)

	code, stdout, stderr := runCommand(t, root, "-m", "m", escaped)
	if code != 3 {
		t.Fatalf("exit = %d, want 3; stdout=%q stderr=%q", code, stdout, stderr)
	}
	_, fields, _ := recordFields(t, stdout)
	if strings.ContainsRune(stdout, '\x1b') {
		t.Fatalf("record carries a raw control byte: %q", stdout)
	}
	if fields["path"] != sanitize.Controls(escaped) {
		t.Fatalf("path = %q, want the sanitizer's spelling %q", fields["path"], sanitize.Controls(escaped))
	}
	head := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
	want := "git restore --source=" + head + " --staged --worktree -- <named-paths>"
	if fields["next"] != want {
		t.Fatalf("next = %q, want %q", fields["next"], want)
	}
}

// FB5: a refusal before publication keeps exit 1 with the error: prefix and moves no ref.
func TestRedGateRefusesWithExitOneAndMovesNoRef(t *testing.T) {
	root, before := landingRepo(t, 1, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "landed\n", 0o644)
	code, stdout, stderr := runCommand(t, root, "-m", "m", "a.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "error: ") {
		t.Fatalf("stderr = %q, want the error: prefix", stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatalf("HEAD moved from %s to %s", before, after)
	}
}

// FB6: a grammar error keeps exit 2, so the new exit code did not displace it.
func TestMissingMessageIsAGrammarError(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, stdout, stderr := runCommand(t, root, "tracked.txt")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "-m <msg> is required") {
		t.Fatalf("stderr = %q, want the missing -m usage error", stderr)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatalf("HEAD moved from %s to %s", before, after)
	}
}

// FB8: the help text names the exit-3 meaning, so the code is discoverable without the
// source.
func TestHelpNamesTheExitThreeMeaning(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, stdout, stderr := runCommand(t, root, "--help")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "exit 3:") || !strings.Contains(stdout, "reconcile") {
		t.Fatalf("help = %q, want a line naming exit 3 and the reconcile", stdout)
	}
}
