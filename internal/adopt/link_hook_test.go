package adopt

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestClassifyPrePushIsAbsentFromProductionAPI(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "link_hook.go", nil, 0)
	if err != nil {
		t.Fatalf("parse link_hook.go: %v", err)
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == "ClassifyPrePush" {
			t.Fatal("ClassifyPrePush must remain absent; InspectPrePush is the hook-health entry point")
		}
	}
}

func TestInspectPrePushRecord(t *testing.T) {
	root := hookTestRepo(t)
	path := filepath.Join(root, ".git", "hooks", "pre-push")
	writeHook(t, path, "baked")
	runHookGit(t, root, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/live")

	live := InspectPrePush(root)
	wantLive := PrePushHealth{State: PrePushManaged, Path: path, Branch: "live", Provenance: PrePushLive, Currency: PrePushCurrent}
	if live != wantLive {
		t.Fatalf("live record = %#v, want %#v", live, wantLive)
	}

	runHookGit(t, root, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
	baked := InspectPrePush(root)
	wantBaked := PrePushHealth{State: PrePushManaged, Path: path, Branch: "baked", Provenance: PrePushBaked, Currency: PrePushCurrent}
	if baked != wantBaked {
		t.Fatalf("baked record = %#v, want %#v", baked, wantBaked)
	}

	writeHook(t, path, fallbackProtectedBranch)
	fallback := InspectPrePush(root)
	if !fallback.Fallback || fallback.Provenance != PrePushBaked || fallback.Branch != fallbackProtectedBranch {
		t.Fatalf("fallback record = %#v", fallback)
	}
}

func TestInspectPrePushCurrencyAndStates(t *testing.T) {
	root := hookTestRepo(t)
	path := filepath.Join(root, ".git", "hooks", "pre-push")
	writeHook(t, path, "baked")

	truncated := strings.TrimSuffix(strings.ReplaceAll(prePushTemplate, prePushBranchToken, "baked"), "\n")
	if err := os.WriteFile(path, []byte(truncated), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := InspectPrePush(root); got.State != PrePushManaged || got.Currency != PrePushStale {
		t.Fatalf("truncated record = %#v, want managed stale", got)
	}

	malformed := strings.ReplaceAll(prePushTemplate, "#!/usr/bin/env bash", "#!/bin/false")
	malformed = strings.ReplaceAll(malformed, prePushBranchToken, "baked")
	if err := os.WriteFile(path, []byte(malformed), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := InspectPrePush(root); got.State != PrePushManaged || got.Currency != PrePushStale {
		t.Fatalf("malformed record = %#v, want managed stale", got)
	}

	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if got := InspectPrePush(root); got.State != PrePushAbsent || got.Currency != PrePushNotApplicable {
		t.Fatalf("absent record = %#v", got)
	}
	if err := os.WriteFile(path, nil, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := InspectPrePush(root); got.State != PrePushForeign || got.Currency != PrePushNotApplicable {
		t.Fatalf("empty record = %#v", got)
	}

	divert := filepath.Join(root, "hooks [one]")
	if err := os.Mkdir(divert, 0o755); err != nil {
		t.Fatal(err)
	}
	runHookGit(t, root, "config", "core.hooksPath", divert)
	if got := InspectPrePush(root); got.State != PrePushDiverted || got.Currency != PrePushNotApplicable || got.Path != filepath.Join(divert, "pre-push") {
		t.Fatalf("diverted record = %#v", got)
	}
	writeHook(t, filepath.Join(divert, "pre-push"), "baked")
	if got := InspectPrePush(root); got.State != PrePushManaged || got.Currency != PrePushCurrent {
		t.Fatalf("literal-path record = %#v", got)
	}
}

func TestInspectPrePushRefusesSpecialFilesAndMissingGit(t *testing.T) {
	root := hookTestRepo(t)
	path := filepath.Join(root, ".git", "hooks", "pre-push")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}
	if got := InspectPrePush(root); got.State != PrePushForeign || got.Currency != PrePushNotApplicable {
		t.Fatalf("FIFO record = %#v", got)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	writeHook(t, path, "baked")
	t.Setenv("PATH", "")
	// With no git executable reachable at all, hooksDir's reader cannot resolve the hooks
	// directory, and InspectPrePush reports the absent state with an empty path rather than
	// falling back to a guessed .git/hooks. This is a direct consequence of GR23, not a
	// missing-git-specific case: any unresolved hooks directory reads the same way.
	if got := InspectPrePush(root); got.State != PrePushAbsent || got.Path != "" {
		t.Fatalf("missing-git record = %#v, want absent with empty path", got)
	}
}

// TestInspectPrePushRefusesSymlinkToSpecialFile plants a symlink to a writerless FIFO
// where the hook belongs. This is the hostile shape a hooks directory can carry that a
// direct FIFO does not. The inspection runs off the test goroutine behind a deadline
// because opening a writerless FIFO never returns. Without the deadline the failure is a
// wedged suite rather than a named red, and every ambient surface that reads hook health
// hangs with it.
func TestInspectPrePushRefusesSymlinkToSpecialFile(t *testing.T) {
	root := hookTestRepo(t)
	path := filepath.Join(root, ".git", "hooks", "pre-push")
	fifo := filepath.Join(t.TempDir(), "fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(fifo, path); err != nil {
		t.Fatal(err)
	}

	done := make(chan PrePushHealth, 1)
	go func() { done <- InspectPrePush(root) }()
	select {
	case got := <-done:
		if got.State != PrePushForeign || got.Currency != PrePushNotApplicable {
			t.Fatalf("symlinked-FIFO record = %#v, want foreign", got)
		}
	case <-time.After(bounds.TestDeadline(bounds.TestDeadlineFloor)):
		t.Fatal("InspectPrePush blocked on a symlinked FIFO: the link is read before its target is typed")
	}

	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "no-such-hook"), path); err != nil {
		t.Fatal(err)
	}
	if got := InspectPrePush(root); got.State != PrePushForeign || got.Currency != PrePushNotApplicable {
		t.Fatalf("dangling-symlink record = %#v, want foreign", got)
	}
}

func TestInspectPrePushAvoidsRemoteGit(t *testing.T) {
	root := hookTestRepo(t)
	path := filepath.Join(root, ".git", "hooks", "pre-push")
	writeHook(t, path, "baked")
	bin := t.TempDir()
	log := filepath.Join(t.TempDir(), "git.log")
	recorder := filepath.Join(bin, "git")
	if err := os.WriteFile(recorder, []byte("#!/bin/sh\nprintf '%s\\n' \"$*\" >> "+log+"\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	_ = InspectPrePush(root)
	contents, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(contents), "ls-remote") || strings.Contains(string(contents), "fetch") {
		t.Fatalf("inspection invoked remote git: %q", contents)
	}
}

func hookTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runHookGit(t, root, "init")
	return root
}

func writeHook(t *testing.T, path, branch string) {
	t.Helper()
	content := strings.ReplaceAll(prePushTemplate, prePushBranchToken, branch)
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func runHookGit(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// runPrePushHook executes the rendered hook the way git does: cwd at the repo root, one
// "<local ref> <local oid> <remote ref> <remote oid>" line on stdin. It returns the exit
// error, if any, and the hook's stderr.
func runPrePushHook(t *testing.T, root, path, remoteRef string) (error, string) {
	t.Helper()
	command := exec.Command("bash", path)
	command.Dir = root
	oid := strings.Repeat("a", 40)
	command.Stdin = strings.NewReader("refs/heads/topic " + oid + " " + remoteRef + " " + oid + "\n")
	var stderr strings.Builder
	command.Stderr = &stderr
	return command.Run(), stderr.String()
}

// TestLinkRefusesUnresolvedHooksDirectory pins GR20: a failed hooks-directory
// query makes transactionalLink refuse with exit 1 before it stages any file, and
// the refusal names the root path and the reader's action. It is serial because
// gittest.StubGit binds the process PATH.
func TestLinkRefusesUnresolvedHooksDirectory(t *testing.T) {
	root := hookTestRepo(t)
	gittest.StubGit(t, root, "fail-git-path", filepath.Join(t.TempDir(), "argv"))

	plan := []planEntry{{rel: "accepted.txt", kind: "inline", content: "accepted\n"}}
	var stdout, stderr strings.Builder
	code, changed := transactionalLink(root, t.TempDir(), "copy", "test", plan, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("transactionalLink under unresolved hooks dir = %d, want 1\nstderr:\n%s", code, stderr.String())
	}
	if changed {
		t.Fatal("transactionalLink reported a change under the refusal")
	}
	if _, err := os.Lstat(filepath.Join(root, "accepted.txt")); !os.IsNotExist(err) {
		t.Fatalf("accepted.txt staged before the refusal: %v", err)
	}
	got := stderr.String()
	const investigateGitFailure = "investigate the git failure"
	if !strings.Contains(got, root) || !strings.Contains(got, investigateGitFailure) {
		t.Fatalf("stderr = %q, want it to name %s and %q", got, root, investigateGitFailure)
	}
}

// TestInspectPrePushReportsAbsentWhenHooksDirectoryIsUnresolved pins GR23: a
// failed hooks-directory query reads as the absent state with an empty path,
// never a guessed .git/hooks. It is serial because gittest.StubGit binds PATH.
func TestInspectPrePushReportsAbsentWhenHooksDirectoryIsUnresolved(t *testing.T) {
	root := hookTestRepo(t)
	gittest.StubGit(t, root, "fail-git-path", filepath.Join(t.TempDir(), "argv"))

	got := InspectPrePush(root)
	if got.State != PrePushAbsent || got.Path != "" {
		t.Fatalf("InspectPrePush under unresolved hooks dir = %#v, want absent with empty path", got)
	}
}

// TestDoctorFixRefusesUnresolvedHooksDirectory pins GR21: an unresolved hooks
// directory makes `bench doctor --fix` exit 1 in a consumer repository, naming the
// root path and the reader's action, before the repair touches the hook state. It
// is serial because gittest.StubGit binds the process PATH.
func TestDoctorFixRefusesUnresolvedHooksDirectory(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	t.Setenv("BENCH_KIT", t.TempDir())
	bin := t.TempDir()
	target := filepath.Join(bin, "bench-wrapper")
	writeFixtureFile(t, target, "#!/bin/sh\n", 0o755)
	t.Setenv("BENCH_WRAPPER", target)
	t.Chdir(root)

	gittest.StubGit(t, root, "fail-git-path", filepath.Join(t.TempDir(), "argv"))

	var stdout, stderr bytes.Buffer
	if code := Doctor([]string{"--fix"}, &stdout, &stderr, "1.0.0"); code != 1 {
		t.Fatalf("Doctor --fix under unresolved hooks dir = %d, want 1\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	got := stderr.String()
	if !strings.Contains(got, root) || !strings.Contains(got, "investigate the git failure") {
		t.Fatalf("stderr = %q, want it to name %s and %q", got, root, "investigate the git failure")
	}
}

// TestUnlinkRefusesUnresolvedHooksDirectory pins GR22: an unresolved hooks
// directory makes `bench unlink` exit 1 over a present managed hook, naming the
// root path and the reader's action, and it leaves the hook in place. It is
// serial because gittest.StubGit binds the process PATH.
func TestUnlinkRefusesUnresolvedHooksDirectory(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, filepath.Join(root, ".bench", "link-manifest.tsv"), "#kit\t1.0.0\n", 0o644)
	hook := filepath.Join(root, ".git", "hooks", "pre-push")
	writeHook(t, hook, fallbackProtectedBranch)
	t.Chdir(root)

	gittest.StubGit(t, root, "fail-git-path", filepath.Join(t.TempDir(), "argv"))

	var stdout, stderr bytes.Buffer
	if code := Unlink(nil, &stdout, &stderr); code != 1 {
		t.Fatalf("Unlink under unresolved hooks dir = %d, want 1\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	got := stderr.String()
	if !strings.Contains(got, root) || !strings.Contains(got, "investigate the git failure") {
		t.Fatalf("stderr = %q, want it to name %s and %q", got, root, "investigate the git failure")
	}
	if _, err := os.Lstat(hook); err != nil {
		t.Fatalf("managed hook removed despite the refusal: %v", err)
	}
}

func TestPrePushHookAllowProtectedPushConfig(t *testing.T) {
	root := hookTestRepo(t)
	path := filepath.Join(root, ".git", "hooks", "pre-push")
	writeHook(t, path, "main")

	if err, out := runPrePushHook(t, root, path, "refs/heads/main"); err == nil || !strings.Contains(out, "blocked: direct push to main") {
		t.Fatalf("default hook: err=%v stderr=%q, want a block", err, out)
	}
	if err, out := runPrePushHook(t, root, path, "refs/heads/topic"); err != nil {
		t.Fatalf("default hook on a topic branch: %v\n%s", err, out)
	}

	runHookGit(t, root, "config", "bench.allowProtectedPush", "true")
	if err, out := runPrePushHook(t, root, path, "refs/heads/main"); err != nil {
		t.Fatalf("allowed hook: %v\n%s", err, out)
	}

	runHookGit(t, root, "config", "bench.allowProtectedPush", "false")
	if err, _ := runPrePushHook(t, root, path, "refs/heads/main"); err == nil {
		t.Fatal("bench.allowProtectedPush=false must keep the block")
	}
}
