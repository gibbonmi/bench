package gittest

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// bounded is how long a test waits for a helper that runs on its own goroutine.
const bounded = 30 * time.Second

// ended is how a helper call on a recorder goroutine ended.
type ended struct {
	failed, skipped, returned bool
	message                   string
}

// recorder keeps a Fatal, FailNow, or Skip call of a helper instead of ending the test.
// Each such call then ends the helper goroutine, as the testing package does.
type recorder struct {
	testing.TB
	ended ended
}

func (r *recorder) Helper() {}

func (r *recorder) Fatal(args ...any) { r.fail(fmt.Sprint(args...)) }

func (r *recorder) Fatalf(format string, args ...any) { r.fail(fmt.Sprintf(format, args...)) }

func (r *recorder) FailNow() { r.fail("") }

func (r *recorder) Skip(args ...any) { r.skip(fmt.Sprint(args...)) }

func (r *recorder) Skipf(format string, args ...any) { r.skip(fmt.Sprintf(format, args...)) }

func (r *recorder) SkipNow() { r.skip("") }

func (r *recorder) fail(message string) {
	r.ended.failed, r.ended.message = true, message
	runtime.Goexit()
}

func (r *recorder) skip(message string) {
	r.ended.skipped, r.ended.message = true, message
	runtime.Goexit()
}

// observe calls helper with a recorder on its own goroutine and reports how it ended.
func observe(t *testing.T, helper func(testing.TB)) ended {
	t.Helper()
	r := &recorder{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		helper(r)
		r.ended.returned = true
	}()
	select {
	case <-done:
	case <-time.After(bounded):
		t.Fatal("helper did not end")
	}
	return r.ended
}

// requireFailed holds that a fatal wrapper failed its caller and did not skip or return.
func requireFailed(t *testing.T, name string, got ended) {
	t.Helper()
	if !got.failed || got.skipped || got.returned {
		t.Fatalf("%s ended %+v; want a test failure", name, got)
	}
}

// withoutGit leaves no git executable on PATH for the rest of the test.
func withoutGit(t *testing.T) {
	t.Helper()
	t.Setenv("PATH", t.TempDir())
	if _, err := exec.LookPath("git"); err == nil {
		t.Fatal("git is still on PATH")
	}
}

// both is a git invocation that writes " out" to stdout and "err" to stderr, then succeeds.
var both = []string{"-c", "alias.both=!printf ' out\\n'; printf 'err\\n' >&2", "both"}

// absent is a git invocation that fails in a repository with no commits.
var absent = []string{"rev-parse", "--verify", "refs/heads/absent"}

func TestRunOutputContract(t *testing.T) {
	var _ func(string, ...string) (string, error) = Run
	root := RepoOnBranch(t, "main")
	if got, err := Run(root, both...); err != nil || got != "out\nerr" {
		t.Fatalf("Run(both) = %q, %v; want the trimmed combined output", got, err)
	}
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(nested)
	for operand, want := range map[string]string{root: "", "": "nested/"} {
		if got, err := Run(operand, "rev-parse", "--show-prefix"); err != nil || got != want {
			t.Fatalf("Run(%q, rev-parse --show-prefix) = %q, %v; want %q", operand, got, err, want)
		}
	}
}

func TestRunFailureReturn(t *testing.T) {
	root := RepoOnBranch(t, "main")
	type result struct {
		out string
		err error
	}
	results := make(chan result, 1)
	go func() {
		out, err := Run(root, absent...)
		results <- result{out, err}
	}()
	var got result
	select {
	case got = <-results:
	case <-time.After(bounded):
		t.Fatal("Run did not return to its goroutine")
	}
	var exit *exec.ExitError
	if got.out != "" || !errors.As(got.err, &exit) {
		t.Fatalf("Run(absent) = %q, %v; want empty output and the wrapped process error", got.out, got.err)
	}
	if message := got.err.Error(); !strings.Contains(message, strings.Join(absent, " ")) || !strings.Contains(message, "fatal") {
		t.Fatalf("Run(absent) error = %q; want the git arguments and git's output", message)
	}
}

func TestRunMissingGit(t *testing.T) {
	root := RepoOnBranch(t, "main")
	withoutGit(t)
	if got, err := Run(root, "status"); got != "" || !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("Run without git = %q, %v; want empty output and the lookup error", got, err)
	}
}

func TestOutputContract(t *testing.T) {
	root := RepoOnBranch(t, "main")
	if got := Output(t, root, both...); got != "out\nerr" {
		t.Fatalf("Output(both) = %q; want the trimmed combined output", got)
	}
}

func TestOutputFailure(t *testing.T) {
	root := RepoOnBranch(t, "main")
	_, want := Run(root, absent...)
	got := observe(t, func(tb testing.TB) { Output(tb, root, absent...) })
	requireFailed(t, "Output(absent)", got)
	if want == nil || got.message != want.Error() {
		t.Fatalf("Output(absent) failure = %q; want the Run error %v", got.message, want)
	}
	withoutGit(t)
	requireFailed(t, "Output without git", observe(t, func(tb testing.TB) { Output(tb, root, "status") }))
}

func TestMustRunFailure(t *testing.T) {
	root := RepoOnBranch(t, "main")
	got := observe(t, func(tb testing.TB) { MustRun(tb, root, absent...) })
	requireFailed(t, "MustRun(absent)", got)
	if !strings.Contains(got.message, strings.Join(absent, " ")) {
		t.Fatalf("MustRun(absent) failure = %q; want the git arguments", got.message)
	}
}

func TestMustRunMissingGit(t *testing.T) {
	root := RepoOnBranch(t, "main")
	withoutGit(t)
	requireFailed(t, "MustRun without git", observe(t, func(tb testing.TB) { MustRun(tb, root, "status") }))
}

func TestOutputBytesContract(t *testing.T) {
	root := RepoOnBranch(t, "main")
	payload := []byte("first\x00second \n\n  ")
	path := filepath.Join(root, "payload")
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	blob := Output(t, root, "hash-object", "-w", "payload")
	if got := OutputBytes(t, root, "cat-file", "blob", blob); !bytes.Equal(got, payload) {
		t.Fatalf("OutputBytes(cat-file) = %q; want %q", got, payload)
	}
	if got := OutputBytes(t, root, both...); string(got) != " out\n" {
		t.Fatalf("OutputBytes(both) = %q; want the raw stdout alone", got)
	}
}

func TestOutputBytesFailure(t *testing.T) {
	root := RepoOnBranch(t, "main")
	requireFailed(t, "OutputBytes(absent)", observe(t, func(tb testing.TB) { OutputBytes(tb, root, absent...) }))
	withoutGit(t)
	requireFailed(t, "OutputBytes without git", observe(t, func(tb testing.TB) { OutputBytes(tb, root, "status") }))
}

func TestCommandArgv(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, `repo with spaces; touch sentinel $(touch sentinel) 'q" *`)
	MustRun(t, "", "init", "-q", root)
	current := filepath.Join(RepoOnBranch(t, "main"), "nested")
	if err := os.Mkdir(current, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(current)
	MustRun(t, root, append(IdentityArgs(), "commit", "--quiet", "--allow-empty", "-m", "explicit root")...)
	if got := string(OutputBytes(t, root, "log", "-1", "--format=%s")); got != "explicit root\n" {
		t.Fatalf("explicit root subject = %q; want the commit in the named repository", got)
	}
	if _, err := Run("", "rev-parse", "--verify", "HEAD"); err == nil {
		t.Fatal("the current directory's repository gained a commit from an explicit root")
	}
	MustRun(t, "", append(IdentityArgs(), "commit", "--quiet", "--allow-empty", "-m", "current directory")...)
	if got := Output(t, "", "log", "-1", "--format=%s"); got != "current directory" {
		t.Fatalf("empty root subject = %q; want the commit in the current directory", got)
	}
	for _, dir := range []string{parent, root, current} {
		if _, err := os.Lstat(filepath.Join(dir, "sentinel")); !os.IsNotExist(err) {
			t.Fatalf("a shell ran the root operand: sentinel in %s: %v", dir, err)
		}
	}
}

func TestCommitContract(t *testing.T) {
	root := RepoOnBranch(t, "main")
	writeKitCopyFile(t, root, "tracked", "tracked bytes", 0o644)
	MustRun(t, root, "add", "tracked")
	message := "-a subject that looks like an option\n\nbody line"
	Commit(t, root, message)
	if got := Output(t, root, "log", "-1", "--format=%B"); got != message {
		t.Fatalf("commit message = %q; want %q", got, message)
	}
	before := Output(t, root, "rev-parse", "HEAD")
	Commit(t, root, "explicit empty", "--allow-empty")
	if got := Output(t, root, "rev-parse", "HEAD^"); got != before {
		t.Fatalf("empty commit parent = %s; want %s", got, before)
	}
	if got := Output(t, root, "diff-tree", "--no-commit-id", "-r", "HEAD"); got != "" {
		t.Fatalf("empty commit changes = %q; want none", got)
	}
}

func TestCommitDoesNotStage(t *testing.T) {
	root := RepoOnBranch(t, "main")
	writeKitCopyFile(t, root, "tracked", "first", 0o644)
	MustRun(t, root, "add", "tracked")
	Commit(t, root, "first")
	writeKitCopyFile(t, root, "tracked", "unstaged change", 0o644)
	writeKitCopyFile(t, root, "staged", "staged bytes", 0o644)
	writeKitCopyFile(t, root, "sentinel", "unstaged sentinel", 0o644)
	MustRun(t, root, "add", "staged")
	Commit(t, root, "index only")
	if got := Output(t, root, "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD"); got != "staged" {
		t.Fatalf("commit paths = %q; want only the staged path", got)
	}
	if got := Output(t, root, "status", "--porcelain"); got != "M tracked\n?? sentinel" {
		t.Fatalf("status after commit = %q; want the unstaged change and sentinel left out", got)
	}
}

func TestCommitNoChanges(t *testing.T) {
	root := RepoOnBranch(t, "main")
	Commit(t, root, "initial", "--allow-empty")
	before := Output(t, root, "rev-parse", "HEAD")
	requireFailed(t, "Commit(unchanged index)", observe(t, func(tb testing.TB) { Commit(tb, root, "no change") }))
	if got := Output(t, root, "rev-parse", "HEAD"); got != before {
		t.Fatalf("HEAD after a refused commit = %s; want %s", got, before)
	}
}

// personalIdentity makes a personal global identity visible to every git child of the
// test. A fixture that falls back to it commits as "personal".
func personalIdentity(t *testing.T) {
	t.Helper()
	global := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(global, []byte("[user]\n\tname = personal\n\temail = personal@example.invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", global)
	for _, name := range []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL"} {
		t.Setenv(name, "")
		if err := os.Unsetenv(name); err != nil {
			t.Fatal(err)
		}
	}
}

// identities is the author and committer of the HEAD commit of root.
func identities(t *testing.T, root string) string {
	t.Helper()
	return Output(t, root, "log", "-1", "--format=%an <%ae>|%cn <%ce>")
}

func TestFixtureIdentity(t *testing.T) {
	personalIdentity(t)
	root := RepoOnBranch(t, "main")
	writeKitCopyFile(t, root, "tracked", "tracked bytes", 0o644)
	MustRun(t, root, "add", "tracked")
	Commit(t, root, "fixture")
	unconfigured := Repo(t)
	MustRun(t, unconfigured, append(IdentityArgs(), "commit", "--quiet", "--allow-empty", "-m", "identity arguments")...)
	canonical := IdentityName + " <" + IdentityEmail + ">"
	want := canonical + "|" + canonical
	for name, repository := range map[string]string{"RepoOnBranch": root, "snapshot": KitCopy(t, root), "IdentityArgs": unconfigured} {
		if got := identities(t, repository); got != want {
			t.Fatalf("%s identity = %q; want %q", name, got, want)
		}
	}
}

func TestFixtureIdentityOverride(t *testing.T) {
	personalIdentity(t)
	root := RepoOnBranch(t, "main")
	MustRun(t, root, "config", "user.name", "subject")
	MustRun(t, root, "config", "user.email", "subject@example.invalid")
	Commit(t, root, "subject", "--allow-empty")
	if got, want := identities(t, root), "subject <subject@example.invalid>|subject <subject@example.invalid>"; got != want {
		t.Fatalf("override identity = %q; want %q", got, want)
	}
}
