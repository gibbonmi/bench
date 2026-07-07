package gate

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/toon"
)

// TestPinCommandNotInRepo pins the one-phrase not-in-repo contract for `bench gate
// pin` at the Go seam rather than through the built binary: the command refuses
// non-TTY stdin before it ever checks git.Root(), so an outside-a-repo cwd is
// unreachable through a non-interactive contract probe. Injecting isTerminal=true
// bypasses that unrelated precondition and exercises the same not-in-repo branch
// every other operational command shares.
func TestPinCommandNotInRepo(t *testing.T) {
	chdir(t, t.TempDir())

	var stdout, stderr bytes.Buffer
	code := pinCommand(nil, strings.NewReader(""), &stdout, &stderr, func(io.Reader) bool { return true })
	if code != 1 {
		t.Fatalf("pinCommand rc = %d, want 1; stderr:\n%s", code, stderr.String())
	}
	if strings.TrimSpace(stderr.String()) != toon.NotInRepo() {
		t.Fatalf("stderr = %q, want %q", stderr.String(), toon.NotInRepo())
	}
}

func TestPinCommandWritesCommittedBenchTree(t *testing.T) {
	root := newPinRepo(t)
	chdir(t, root)
	committedTree := gitOutput(t, root, "rev-parse", "HEAD:.bench")
	os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte("#!/usr/bin/env bash\nexit 1\n"), 0o755)

	var stdout, stderr bytes.Buffer
	code := pinCommand(nil, strings.NewReader("pin .bench\n"), &stdout, &stderr, func(io.Reader) bool { return true })
	if code != 0 {
		t.Fatalf("pinCommand exit = %d\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	lines := readPinLines(t, root)
	if lines[0] != committedTree {
		t.Fatalf("pin tree = %q, want committed HEAD:.bench %q", lines[0], committedTree)
	}
	if lines[1] != gitOutput(t, root, "rev-parse", "HEAD") {
		t.Fatalf("pin commit = %q, want HEAD", lines[1])
	}
	if !strings.Contains(stderr.String(), "uncommitted changes") {
		t.Fatalf("dirty .bench warning missing from stderr:\n%s", stderr.String())
	}
}

func TestPinCommandDeclineAndSecondPinOverwrite(t *testing.T) {
	root := newPinRepo(t)
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := pinCommand(nil, strings.NewReader("no\n"), &stdout, &stderr, func(io.Reader) bool { return true })
	if code == 0 {
		t.Fatal("pinCommand accepted a declined confirmation")
	}
	if _, err := os.Stat(pinPath(root)); err == nil {
		t.Fatal("pinCommand wrote a pin file after decline")
	}

	code = pinCommand(nil, strings.NewReader("pin .bench\n"), &stdout, &stderr, func(io.Reader) bool { return true })
	if code != 0 {
		t.Fatalf("first pin exit = %d", code)
	}
	firstCommit := readPinLines(t, root)[1]
	os.WriteFile(filepath.Join(root, ".bench", "extra"), []byte("ok\n"), 0o644)
	gitRun(t, root, "add", ".bench/extra")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "gate change")
	secondCommit := gitOutput(t, root, "rev-parse", "HEAD")

	code = pinCommand(nil, strings.NewReader("pin .bench\n"), &stdout, &stderr, func(io.Reader) bool { return true })
	if code != 0 {
		t.Fatalf("second pin exit = %d", code)
	}
	lines := readPinLines(t, root)
	if lines[1] != secondCommit {
		t.Fatalf("second pin commit = %q, want %q", lines[1], secondCommit)
	}
	if lines[1] == firstCommit {
		t.Fatal("second pin did not overwrite the previous commit")
	}
}

func newPinRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatalf("mkdir .bench: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte("#!/usr/bin/env bash\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write gate: %v", err)
	}
	gitRun(t, root, "add", ".bench/gate.sh")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "init")
	return root
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
}

func gitRun(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func gitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

func readPinLines(t *testing.T, root string) []string {
	t.Helper()
	data, err := os.ReadFile(pinPath(root))
	if err != nil {
		t.Fatalf("read pin: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("pin file has %d lines, want 3: %q", len(lines), data)
	}
	return lines
}
