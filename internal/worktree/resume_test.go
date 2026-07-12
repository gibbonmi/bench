package worktree

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestResumeCleanRemovesOnlyCleanUnlockedOutOfPool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	clean := filepath.Join(filepath.Dir(root), "auto clean")
	dirty := filepath.Join(filepath.Dir(root), "auto dirty")
	locked := filepath.Join(filepath.Dir(root), "auto locked")
	pool := filepath.Join(Pool(root), "leased")
	if err := os.MkdirAll(filepath.Dir(pool), 0o700); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{clean, dirty, locked, pool} {
		gitRun(t, root, "worktree", "add", "-q", "--detach", path, "HEAD")
	}
	if err := os.WriteFile(filepath.Join(dirty, "dirty.txt"), []byte("recover me\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "worktree", "lock", locked)
	lease, err := LeaseFile(pool)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lease, []byte("123 2026-07-11T00:00:00Z\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "branch", "worktree-agent-landed")
	unique := gitOutput(t, root, "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", "unique scratch")
	gitRun(t, root, "update-ref", "refs/heads/worktree-agent-unique", unique)
	created := time.Unix(1, 0).UTC()
	if err := intent.Upsert(root, intent.Entry{Key: "auto-cleaned", Kind: intent.KindWorktree, Objective: "clean me", CreatedAt: created, Worktree: clean}); err != nil {
		t.Fatal(err)
	}
	if err := intent.Upsert(root, intent.Entry{Key: "unrelated", Kind: intent.KindShift, Objective: "keep me", CreatedAt: created}); err != nil {
		t.Fatal(err)
	}
	chdir(t, root)
	before, _ := os.ReadFile(filepath.Join(dirty, "dirty.txt"))
	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
	if got := stdout.String(); got != "bench resume: cleaned 1 worktree(s), 1 landed branch(es); kept 1 dirty, 1 locked, 1 leased; 1 open intent(s)\n" {
		t.Fatalf("resume report = %q", got)
	}
	if _, err := os.Stat(clean); !os.IsNotExist(err) {
		t.Fatalf("clean worktree remains: %v", err)
	}
	for _, path := range []string{dirty, locked, pool} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("kept worktree %q: %v", path, err)
		}
	}
	after, _ := os.ReadFile(filepath.Join(dirty, "dirty.txt"))
	if !bytes.Equal(before, after) {
		t.Fatal("resume cleanup changed dirty bytes")
	}
	if git.OK("-C", root, "show-ref", "--verify", "--quiet", "refs/heads/worktree-agent-landed") {
		t.Fatal("resume cleanup kept ancestry-landed orphan")
	}
	if !git.OK("-C", root, "show-ref", "--verify", "--quiet", "refs/heads/worktree-agent-unique") {
		t.Fatal("resume cleanup deleted unique orphan")
	}
}

func TestResumeCleanKeepsIgnoredOnlyOutOfPoolWorktree(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ignore")
	candidate := filepath.Join(filepath.Dir(root), "ignored-only")
	gitRun(t, root, "worktree", "add", "-q", "--detach", candidate, "HEAD")
	ignored := filepath.Join(candidate, "ignored.txt")
	if err := os.WriteFile(ignored, []byte("retain me\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(ignored); err != nil {
		t.Fatalf("ignored-only WIP was not retained: %v", err)
	}
	if !strings.Contains(stdout.String(), "kept 1 dirty") {
		t.Fatalf("ignored-only WIP not classified as retained dirty state: %q", stdout.String())
	}
}

func TestResumeCleanRemovesNewlinePathWorktree(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	candidate := filepath.Join(filepath.Dir(root), "newline\nworktree")
	gitRun(t, root, "worktree", "add", "-q", "--detach", candidate, "HEAD")
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(candidate); !os.IsNotExist(err) {
		t.Fatalf("newline-bearing clean worktree remains: %v", err)
	}
	if !strings.Contains(stdout.String(), "cleaned 1 worktree") {
		t.Fatalf("newline-bearing cleanup report = %q", stdout.String())
	}
}

func TestResumeCleanFailsWhenLandednessQueryFails(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	gitRun(t, root, "checkout", "-q", "-b", "worktree-agent-query-failure")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "--allow-empty", "-qm", "unique")
	gitRun(t, root, "checkout", "-q", "main")
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	wrapper := filepath.Join(bin, "git")
	script := "#!/bin/sh\nfor arg in \"$@\"; do\n  [ \"$arg\" = rev-list ] && exit 17\ndone\nexec \"" + realGit + "\" \"$@\"\n"
	if err := os.WriteFile(wrapper, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 1 {
		t.Fatalf("ResumeCleanCommand exit=%d, want 1\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	}
	if !git.OK("-C", root, "show-ref", "--verify", "--quiet", "refs/heads/worktree-agent-query-failure") {
		t.Fatal("resume cleanup deleted branch after landedness query failure")
	}
}
