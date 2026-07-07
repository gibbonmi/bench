package stophook

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestActive(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"json true is active", `{"stop_hook_active":true}`, true},
		{"json false is not active", `{"stop_hook_active":false}`, false},
		{"top-level false is not active", `false`, false},
		{"absent field is not active", `{}`, false},
		{"string True is active", `{"stop_hook_active":"True"}`, true},
		{"number is not active", `{"stop_hook_active":1}`, false},
		{"other string is not active", `{"stop_hook_active":"yes"}`, false},
		{"invalid json is not active", `not json`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Active([]byte(c.in)); got != c.want {
				t.Errorf("Active(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

// linesOf builds an n-line string joined by "\n" (no trailing newline), each line
// numbered so the boundary of a tail is checkable.
func linesOf(n int) string {
	xs := make([]string, n)
	for i := range xs {
		xs[i] = fmt.Sprintf("line %d", i+1)
	}
	return strings.Join(xs, "\n")
}

func TestTail(t *testing.T) {
	cases := []struct {
		name string
		in   string
		n    int
		want string
	}{
		{"31 lines yields last 30", linesOf(31), 30, linesOf(31)[len("line 1\n"):]},
		{"exactly 30 lines is unchanged", linesOf(30), 30, linesOf(30)},
		{"5 lines is unchanged", linesOf(5), 30, linesOf(5)},
		{"empty is empty", "", 30, ""},
		{"trailing newline is not a blank line", "a\nb\n", 30, "a\nb"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Tail(c.in, c.n); got != c.want {
				t.Errorf("Tail(%q, %d) = %q, want %q", c.in, c.n, got, c.want)
			}
		})
	}
}

// TestTailBoundaryCount is the sharp boundary check: 31 lines under a 30-line tail
// drops exactly the first line and keeps 30.
func TestTailBoundaryCount(t *testing.T) {
	got := Tail(linesOf(31), 30)
	lines := strings.Split(got, "\n")
	if len(lines) != 30 {
		t.Fatalf("Tail(31 lines, 30) produced %d lines, want 30", len(lines))
	}
	if lines[0] != "line 2" {
		t.Errorf("first kept line = %q, want %q", lines[0], "line 2")
	}
	if lines[29] != "line 31" {
		t.Errorf("last kept line = %q, want %q", lines[29], "line 31")
	}
}

func TestBlockMessage(t *testing.T) {
	msg := BlockMessage(linesOf(40))

	if !strings.HasPrefix(msg, "BLOCKED: the gate is red, so this shift is not done.\n") {
		t.Errorf("BlockMessage missing header, got:\n%s", msg)
	}
	if !strings.Contains(msg, "do not weaken or skip a check") {
		t.Errorf("BlockMessage missing middle header line, got:\n%s", msg)
	}
	if !strings.HasSuffix(msg, "Gate output:\n"+Tail(linesOf(40), 30)) {
		t.Errorf("BlockMessage tail mismatch, got:\n%s", msg)
	}

	// Only the last 30 of the 40 gate lines survive into the message.
	if strings.Contains(msg, "line 10\n") {
		t.Errorf("BlockMessage kept a line beyond the 30-line tail (line 10), got:\n%s", msg)
	}
	if !strings.Contains(msg, "line 11") {
		t.Errorf("BlockMessage dropped the first kept line (line 11), got:\n%s", msg)
	}
	if !strings.Contains(msg, "line 40") {
		t.Errorf("BlockMessage dropped the last line (line 40), got:\n%s", msg)
	}
}

// newGitRepo makes a fresh git-init temp repo and chdirs into it, so Run's internal
// git.Root() resolves there and gate.Record can write <git-dir>/bench-last-gate. It
// returns the repo dir. No seed commit is needed: git.TreeHash falls back to the
// empty tree in a commit-less repo.
func newGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v\n%s", err, out)
	}
	t.Chdir(dir)
	return dir
}

// writeScript writes an executable-mode-controlled stub at dir/name with the given
// body and returns its path. mode 0o755 makes a runnable gate; 0o644 makes the
// non-executable gate whose exec fails at start (a non-*exec.ExitError).
func writeScript(t *testing.T, dir, name, body string, mode os.FileMode) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
	return path
}

// exitGate is a runnable stub that ignores its args and exits with code.
func exitGate(t *testing.T, dir string, code int) string {
	t.Helper()
	return writeScript(t, dir, "gate", fmt.Sprintf("#!/usr/bin/env bash\nexit %d\n", code), 0o755)
}

// sentinelGate is a runnable stub that creates sentinel when invoked, then exits 0.
// A test asserts the sentinel is absent to prove Run never reached the gate.
func sentinelGate(t *testing.T, dir, sentinel string) string {
	t.Helper()
	body := fmt.Sprintf("#!/usr/bin/env bash\ntouch %q\nexit 0\n", sentinel)
	return writeScript(t, dir, "gate", body, 0o755)
}

// readCache reads the verdict cache the same way gate.Record wrote it: at
// <absolute-git-dir>/bench-last-gate. It returns (line, present).
func readCache(t *testing.T, dir string) (string, bool) {
	t.Helper()
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--absolute-git-dir")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("rev-parse --absolute-git-dir failed: %v", err)
	}
	gitdir := strings.TrimSpace(string(out))
	data, err := os.ReadFile(filepath.Join(gitdir, "bench-last-gate"))
	if os.IsNotExist(err) {
		return "", false
	}
	if err != nil {
		t.Fatalf("reading verdict cache: %v", err)
	}
	return string(data), true
}

// TestRun drives the Stop-hook orchestration against a stub gate in a temp git repo
// across the four verdict cases plus the two short-circuit guards, asserting the
// return code, the BLOCKED-on-stderr posture, and the verdict-cache state. It pins
// the rc→verdict mapping and the no-forged-verdict guarantee that the shell gate
// contracts otherwise exercise only end-to-end.
func TestRun(t *testing.T) {
	t.Run("green gate allows and caches green", func(t *testing.T) {
		dir := newGitRepo(t)
		wrapper := exitGate(t, dir, 0)
		var stderr bytes.Buffer
		if rc := Run(nil, wrapper, true, &stderr); rc != 0 {
			t.Fatalf("Run(green) = %d, want 0", rc)
		}
		if strings.Contains(stderr.String(), "BLOCKED") {
			t.Errorf("green run wrote BLOCKED to stderr:\n%s", stderr.String())
		}
		line, ok := readCache(t, dir)
		if !ok {
			t.Fatal("green run recorded no verdict cache, want one")
		}
		if !strings.HasPrefix(line, "green ") {
			t.Errorf("cache line = %q, want a green-prefixed verdict", line)
		}
	})

	t.Run("red gate blocks and caches red", func(t *testing.T) {
		dir := newGitRepo(t)
		wrapper := exitGate(t, dir, 2)
		var stderr bytes.Buffer
		if rc := Run(nil, wrapper, true, &stderr); rc != 2 {
			t.Fatalf("Run(red) = %d, want 2", rc)
		}
		if !strings.Contains(stderr.String(), "BLOCKED: the gate is red") {
			t.Errorf("red run missing BLOCKED header on stderr:\n%s", stderr.String())
		}
		line, ok := readCache(t, dir)
		if !ok {
			t.Fatal("red run recorded no verdict cache, want one")
		}
		if !strings.HasPrefix(line, "red ") {
			t.Errorf("cache line = %q, want a red-prefixed verdict", line)
		}
	})

	t.Run("rc==3 no-gate blocks without forging a verdict", func(t *testing.T) {
		dir := newGitRepo(t)
		wrapper := exitGate(t, dir, 3)
		var stderr bytes.Buffer
		if rc := Run(nil, wrapper, true, &stderr); rc != 2 {
			t.Fatalf("Run(no-gate) = %d, want 2", rc)
		}
		if !strings.Contains(stderr.String(), "BLOCKED: the gate is red") {
			t.Errorf("no-gate run missing BLOCKED header on stderr:\n%s", stderr.String())
		}
		if line, ok := readCache(t, dir); ok {
			t.Errorf("no-gate run forged a verdict cache = %q, want none (no-forged-verdict guarantee)", line)
		}
	})

	t.Run("non-executable gate is treated as red", func(t *testing.T) {
		dir := newGitRepo(t)
		// 0o644: the exec fails at start with a non-*exec.ExitError, driving the
		// rc = 1 treat-as-red branch. The body/exit code never runs.
		wrapper := writeScript(t, dir, "gate", "#!/usr/bin/env bash\nexit 0\n", 0o644)
		var stderr bytes.Buffer
		if rc := Run(nil, wrapper, true, &stderr); rc != 2 {
			t.Fatalf("Run(non-executable) = %d, want 2", rc)
		}
		if !strings.Contains(stderr.String(), "BLOCKED: the gate is red") {
			t.Errorf("non-executable run missing BLOCKED header on stderr:\n%s", stderr.String())
		}
		line, ok := readCache(t, dir)
		if !ok {
			t.Fatal("non-executable run recorded no verdict cache, want a red one")
		}
		if !strings.HasPrefix(line, "red ") {
			t.Errorf("cache line = %q, want a red-prefixed verdict", line)
		}
	})

	t.Run("not armed allows without running the gate", func(t *testing.T) {
		dir := newGitRepo(t)
		sentinel := filepath.Join(dir, "invoked")
		wrapper := sentinelGate(t, dir, sentinel)
		var stderr bytes.Buffer
		if rc := Run(nil, wrapper, false, &stderr); rc != 0 {
			t.Fatalf("Run(armed=false) = %d, want 0", rc)
		}
		if _, err := os.Stat(sentinel); err == nil {
			t.Error("unarmed run invoked the gate, want it skipped")
		}
		if line, ok := readCache(t, dir); ok {
			t.Errorf("unarmed run touched the verdict cache = %q, want it untouched", line)
		}
	})

	t.Run("stop_hook_active allows without running the gate", func(t *testing.T) {
		dir := newGitRepo(t)
		sentinel := filepath.Join(dir, "invoked")
		wrapper := sentinelGate(t, dir, sentinel)
		var stderr bytes.Buffer
		if rc := Run([]byte(`{"stop_hook_active":true}`), wrapper, true, &stderr); rc != 0 {
			t.Fatalf("Run(stop_hook_active) = %d, want 0", rc)
		}
		if _, err := os.Stat(sentinel); err == nil {
			t.Error("already-active run invoked the gate, want it skipped")
		}
		if line, ok := readCache(t, dir); ok {
			t.Errorf("already-active run touched the verdict cache = %q, want it untouched", line)
		}
	})
}
