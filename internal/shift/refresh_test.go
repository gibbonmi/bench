package shift

import (
	"bytes"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/testrepo"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoopRefreshResolvesStartAfterFetch(t *testing.T) {
	root, _ := shiftCollisionFixture(t)
	remote := filepath.Join(t.TempDir(), "origin.git")
	run := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	run(root, "init", "-q", "--bare", "-b", "main", remote)
	branch := run(root, "branch", "--show-current")
	run(root, "remote", "add", "origin", remote)
	run(root, "push", "-q", "-u", "origin", branch)
	run(root, "remote", "set-head", "origin", "-a")
	publisher := filepath.Join(t.TempDir(), "publisher")
	run(root, "clone", "-q", remote, publisher)
	if err := os.WriteFile(filepath.Join(publisher, "remote-only"), []byte("refreshed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(publisher, "-c", "user.name=bench", "-c", "user.email=bench@local", "add", "remote-only")
	run(publisher, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "remote advance")
	run(publisher, "push", "-q", "origin", branch)
	remoteHead := run(publisher, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	if code := loop("delivery", "refresh start ref", true, &stdout, &stderr); code != 4 {
		t.Fatalf("loop = %d, want no-op/4; stdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "log --oneline "+remoteHead+"..") {
		t.Fatalf("refresh did not select fetched start %s:\n%s", remoteHead, stdout.String())
	}
}

func TestCommitmentShiftBeforeEffects(t *testing.T) {
	root, _ := shiftCollisionFixture(t)
	marker := filepath.Join(root, ".git", "adapter-ran")
	agent := filepath.Join(root, "agent")
	if err := os.WriteFile(agent, []byte("#!/bin/sh\nprintf ran > '"+marker+"'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "-C", root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-am", "mark adapter execution")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("commit: %v %s", err, out)
	}
	var stdout, stderr bytes.Buffer
	code := Command([]string{"--outcome", "uncommitted", "--refresh", "attempt delivery"}, &stdout, &stderr)
	ledger, err := intent.Read(root)
	_, markerErr := os.Stat(marker)
	if code != 2 || !os.IsNotExist(markerErr) || err != nil || len(ledger.Entries) != 0 {
		t.Fatalf("refused shift had effects: code=%d marker=%v entries=%d err=%v stdout=%s stderr=%s", code, markerErr, len(ledger.Entries), err, &stdout, &stderr)
	}
}

// shiftCollisionFixture builds a bare repo plus a passing gate and agent, and points
// timeNow at a fixed instant so the derived branch name is deterministic. preExisting
// names additional branches, relative to the base bench/shift-<ts> name (e.g. "-2"), for
// the fixture to pre-create. This exercises the loop's collision retry.
func shiftCollisionFixture(t *testing.T, preExisting ...string) (tmp, baseBranch string) {
	t.Helper()
	tmp = t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = tmp
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}
	runGit("init", "-q", "-b", "main")
	f := testrepo.NewGateFixture(t.TempDir())
	if err := f.Write(tmp, f.Command("bash")+" -c 'exit 0'\n", ""); err != nil {
		t.Fatal(err)
	}
	agentPath := filepath.Join(tmp, "agent")
	if err := os.WriteFile(agentPath, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGit("-c", "user.email=bench@local", "-c", "user.name=bench", "add", "-A")
	runGit("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "init")
	admitShiftFixture(t, tmp)

	fixed := time.Date(2026, 7, 4, 9, 30, 0, 0, time.UTC)
	baseBranch = "bench/shift-" + fixed.Format("20060102-150405")
	runGit("branch", baseBranch)
	for _, suffix := range preExisting {
		runGit("branch", baseBranch+suffix)
	}
	oldNow := timeNow
	timeNow = func() time.Time { return fixed }
	t.Cleanup(func() { timeNow = oldNow })

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("BENCH_AGENT", agentPath)
	t.Setenv("BENCH_HOME", filepath.Join(tmp, "bench-home"))
	t.Setenv("BENCH_MAX_ITERS", "1")
	return tmp, baseBranch
}

func admitShiftFixture(t *testing.T, root string) {
	t.Helper()
	const deliverable = "specs/shift/spec.md"
	commitmenttest.Write(t, root, deliverable, "# Shift fixture\n\nStatus: staged\n")
	commitmenttest.SeedAdmission(t, root, deliverable)
	for _, args := range [][]string{{"add", "specs/shift/spec.md", ".bench/commitment.json"}, {"-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "approve shift fixture"}} {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	commitmenttest.Register(t, root, "shift-fixture")
	commitmenttest.Admit(t, root, "shift-fixture", deliverable)
}
