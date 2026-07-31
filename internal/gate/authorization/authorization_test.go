package authorization

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

func TestAuthorizeClassifiesAndValidatesOwnerEvidence(t *testing.T) {
	root := authorizationRepo(t)
	base := gitOutput(t, root, "write-tree")
	if got := Authorize(t.Context(), root, base); got.Kind != Inherited || Validate(root, base, got.Evidence) {
		t.Fatalf("unproved red = %+v, validate=%v; want inherited non-reusable", got, Validate(root, base, got.Evidence))
	}

	os.Remove(filepath.Join(root, "fail"))
	gitRun(t, root, "add", "-u")
	gitRun(t, root, "commit", "-q", "-m", "green base")
	branch := gitOutput(t, root, "branch", "--show-current")
	tip := gitOutput(t, root, "rev-parse", "HEAD")
	if err := Bootstrap(root, branch, tip, ""); err == nil {
		t.Fatal("bootstrap accepted a subject with no green evidence")
	}
	if out, err := exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", "refs/bench/green/"+branch).CombinedOutput(); err == nil {
		t.Fatalf("missing-evidence bootstrap created a marker: %s", out)
	}
	if got := gate.Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed inherited green = %+v", got)
	}
	if err := Bootstrap(root, branch, tip, ""); err != nil {
		t.Fatalf("bootstrap exact green: %v", err)
	}
	if err := Bootstrap(root, branch, tip, ""); err != nil {
		t.Fatalf("idempotent bootstrap: %v", err)
	}
	if got := gate.ValidateProjectGreen(root, branch); !got.ReusableGreen {
		t.Fatalf("bootstrapped project green = %+v", got)
	}
	gitRun(t, root, "checkout", "-q", "--detach")
	if err := Bootstrap(root, branch, tip, ""); err == nil {
		t.Fatal("bootstrap accepted detached HEAD")
	}
	gitRun(t, root, "checkout", "-q", branch)
	previous := gitOutput(t, root, "rev-parse", "HEAD^")
	gitRun(t, root, "update-ref", "refs/bench/green/"+branch, previous, tip)
	if err := Bootstrap(root, branch, tip, ""); err == nil {
		t.Fatal("bootstrap overwrote a conflicting marker")
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/"+branch); got != previous {
		t.Fatalf("conflicting marker = %s, want %s", got, previous)
	}
	if err := Bootstrap(root, branch, tip, previous); err != nil {
		t.Fatalf("advance expected descendant marker: %v", err)
	}
	if err := Bootstrap(root, branch, tip, previous); err != nil {
		t.Fatalf("replay advanced marker: %v", err)
	}
	os.WriteFile(filepath.Join(root, "later"), []byte("stale\n"), 0o644)
	gitRun(t, root, "add", "later")
	gitRun(t, root, "commit", "-q", "-m", "stale descendant")
	stale := gitOutput(t, root, "rev-parse", "HEAD")
	if err := Bootstrap(root, branch, stale, tip); err == nil {
		t.Fatal("bootstrap accepted a descendant without exact green evidence")
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/"+branch); got != tip {
		t.Fatalf("stale-evidence marker = %s, want %s", got, tip)
	}
	gitRun(t, root, "reset", "--hard", tip)
	greenTree := gitOutput(t, root, "write-tree")
	green := Authorize(t.Context(), root, greenTree)
	if green.Kind != Green || !Validate(root, greenTree, green.Evidence) {
		t.Fatalf("green authorization = %+v, validate=%v", green, Validate(root, greenTree, green.Evidence))
	}

	os.WriteFile(filepath.Join(root, "fail"), []byte("red\n"), 0o644)
	gitRun(t, root, "add", "fail")
	redTree := gitOutput(t, root, "write-tree")
	gitRun(t, root, "reset", "--hard", "HEAD")
	red := Authorize(t.Context(), root, redTree)
	if red.Kind != Candidate || Validate(root, redTree, red.Evidence) || Validate(root, greenTree, red.Evidence) {
		t.Fatalf("candidate red = %+v, same=%v other=%v", red, Validate(root, redTree, red.Evidence), Validate(root, greenTree, red.Evidence))
	}
}

func TestAuthorizeClassifiesMissingGateAsInfrastructure(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	os.WriteFile(filepath.Join(root, "tracked"), []byte("x\n"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "base")
	tree := gitOutput(t, root, "write-tree")
	got := Authorize(t.Context(), root, tree)
	if got.Kind != Infrastructure || Validate(root, tree, got.Evidence) {
		t.Fatalf("missing gate = %+v, validate=%v", got, Validate(root, tree, got.Evidence))
	}
}

func authorizationRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	os.MkdirAll(filepath.Join(root, ".bench"), 0o755)
	os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte("#!/bin/sh\ntest ! -f fail\n"), 0o755)
	os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(`{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`), 0o644)
	os.WriteFile(filepath.Join(root, "fail"), []byte("red\n"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "red")
	return root
}

func gitRun(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

func gitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return string(bytesTrimSpace(out))
}

func bytesTrimSpace(value []byte) []byte {
	start, end := 0, len(value)
	for start < end && (value[start] == ' ' || value[start] == '\n' || value[start] == '\r' || value[start] == '\t') {
		start++
	}
	for end > start && (value[end-1] == ' ' || value[end-1] == '\n' || value[end-1] == '\r' || value[end-1] == '\t') {
		end--
	}
	return value[start:end]
}
