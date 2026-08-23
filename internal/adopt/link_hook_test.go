package adopt

import (
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
	if got := InspectPrePush(root); got.State != PrePushManaged || got.Branch != "baked" || got.Provenance != PrePushBaked {
		t.Fatalf("missing-git record = %#v", got)
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
