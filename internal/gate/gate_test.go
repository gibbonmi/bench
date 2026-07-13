package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// TestResolvePrecedence pins the ordered chain as a pure function: `.bench/gate.sh`
// beats `$BENCH_GATE` beats the auto-detect lockfiles in their fixed order, and an
// absent gate resolves to None. A reordered chain would silently run the wrong oracle;
// no black-box assertion pins this cheaply, so the table is the genuinely-red signal.
func TestResolvePrecedence(t *testing.T) {
	// present names the set of paths the injected probes report as existing/executable.
	fs := func(present ...string) FS {
		set := map[string]bool{}
		for _, p := range present {
			set[p] = true
		}
		has := func(p string) bool { return set[p] }
		return FS{Executable: has, Exists: has}
	}
	const gateSh = "/r/.bench/gate.sh"

	cases := []struct {
		name      string
		benchGate string
		fs        FS
		want      Kind
	}{
		{"gate.sh beats BENCH_GATE and lockfiles", "echo hi", fs(gateSh, "/r/package.json"), GateSh},
		{"BENCH_GATE beats lockfiles", "echo hi", fs("/r/package.json"), BenchGate},
		{"pnpm beats npm", "", fs("/r/pnpm-lock.yaml", "/r/package.json"), Pnpm},
		{"package.json picks npm", "", fs("/r/package.json"), Npm},
		{"pyproject after npm", "", fs("/r/pyproject.toml"), Pyproject},
		{"cargo last", "", fs("/r/Cargo.toml"), Cargo},
		{"nothing resolves to None", "", fs(), None},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Resolve("/r", tc.benchGate, tc.fs)
			if got.Kind != tc.want {
				t.Errorf("Resolve = %v, want %v", got.Kind, tc.want)
			}
			if tc.want == BenchGate && got.Command != tc.benchGate {
				t.Errorf("BenchGate command = %q, want %q", got.Command, tc.benchGate)
			}
		})
	}
}

func TestGateEnvStripsWrapperRoutingInternals(t *testing.T) {
	t.Setenv("BENCH_KIT", "/wrong/kit")
	t.Setenv("BENCH_WRAPPER", "/wrong/wrapper")
	t.Setenv("BENCH_GATE", "echo ok")

	env := gateEnv()
	sawBenchGate := false
	for _, kv := range env {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") {
			t.Fatalf("gateEnv leaked wrapper-routing internal %q", kv)
		}
		if kv == "BENCH_GATE=echo ok" {
			sawBenchGate = true
		}
	}
	if !sawBenchGate {
		t.Fatal("gateEnv stripped BENCH_GATE, which is part of the project gate contract")
	}
}

func writeGateTestFile(t *testing.T, root, name, body string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
}

func gateTestRepo(t *testing.T, script, manifest string) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
	writeGateTestFile(t, root, ".bench/gate.sh", script, 0o755)
	if manifest != "" {
		writeGateTestFile(t, root, ".bench/gate-inputs.json", manifest, 0o644)
	}
	return root
}

func TestClosedSubjectReuseAndMutations(t *testing.T) {
	const manifest = `{"schema":1,"closure":"local","environment":["ORACLE_INPUT"],"paths":["ignored-input"],"tools":[]}`
	root := gateTestRepo(t, "#!/usr/bin/env bash\necho run >> .git/runs\n[ \"$ORACLE_INPUT\" = green ] && [ \"$(cat ignored-input)\" = green ]\n", manifest)
	writeGateTestFile(t, root, ".gitignore", "ignored-input\n", 0o644)
	writeGateTestFile(t, root, "ignored-input", "green\n", 0o644)
	t.Setenv("ORACLE_INPUT", "green")
	var stderr bytes.Buffer
	first := Execute(context.Background(), root, io.Discard, &stderr)
	if first.ActionExit != 0 || !Inspect(root).ReusableGreen {
		t.Fatalf("closed green did not become reusable: result=%+v inspect=%+v stderr=%s", first, Inspect(root), stderr.String())
	}

	t.Setenv("ORACLE_INPUT", "red")
	if got := Inspect(root); got.ReusableGreen || got.Reason != "oracle changed" {
		t.Fatalf("environment mutation reused verdict: %+v", got)
	}
	second := Execute(context.Background(), root, io.Discard, io.Discard)
	if second.GateExit == 0 || second.ActionExit == 0 {
		t.Fatalf("red mutated oracle normalized to success: %+v", second)
	}

	t.Setenv("ORACLE_INPUT", "green")
	writeGateTestFile(t, root, "ignored-input", "red\n", 0o644)
	if got := Inspect(root); got.ReusableGreen {
		t.Fatalf("declared ignored-input mutation reused verdict: %+v", got)
	}
}

func TestGateEnvironmentIsPasslisted(t *testing.T) {
	const manifest = `{"schema":1,"closure":"local","environment":["DECLARED"],"paths":[],"tools":[]}`
	root := gateTestRepo(t, "#!/usr/bin/env bash\nprintf '%s:%s:%s:%s\\n' \"${DECLARED-unset}\" \"${AMBIENT-unset}\" \"${HOME-unset}\" \"${BENCH_GATE-unset}\" > .git/environment\n", manifest)
	t.Setenv("DECLARED", "kept")
	t.Setenv("AMBIENT", "secret")
	t.Setenv("HOME", "/secret/home")
	t.Setenv("BENCH_GATE", "should-not-win-because-gate-script-does")
	result := Execute(context.Background(), root, io.Discard, io.Discard)
	if result.ActionExit != 0 {
		t.Fatalf("passlist gate failed: %+v", result)
	}
	data, err := os.ReadFile(filepath.Join(root, ".git", "environment"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "kept:unset:unset:unset\n" {
		t.Fatalf("gate environment leaked ambient values: %q", data)
	}
}

func TestStrictVerdictInspection(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	s, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	gitdir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		t.Fatal(err)
	}
	cache := filepath.Join(gitdir, benchgit.GateCacheFile)
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	ready := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":"2026-07-13T11:50:01Z"}`, s.Tree, s.Oracle)
	write := func(body string, mode os.FileMode) {
		t.Helper()
		if err := os.RemoveAll(cache); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(cache, []byte(body), mode); err != nil {
			t.Fatal(err)
		}
	}
	write(ready, 0o600)
	if got := inspectAt(root, now); !got.ReusableGreen {
		t.Fatalf("fresh exact ready not reusable: %+v", got)
	}
	if got := inspectAt(root, now.Add(time.Second)); got.ReusableGreen || got.Reason != "verdict expired" {
		t.Fatalf("exact ten-minute boundary reused: %+v", got)
	}
	for _, invalid := range []string{
		"green " + s.Tree + " 2026-07-13T11:50:01Z\n",
		strings.Replace(ready, `"schema":1`, `"schema":1,"schema":1`, 1),
		strings.Replace(ready, `"status":"green",`, "", 1),
		strings.TrimSuffix(ready, "}") + `,"unknown":1}`,
		ready + `{}`,
	} {
		write(invalid, 0o600)
		if got := inspectAt(root, now); got.State != Invalid {
			t.Fatalf("invalid record classified %q: %+v", invalid, got)
		}
	}
	write(ready, 0o644)
	if got := inspectAt(root, now); got.State != Invalid {
		t.Fatalf("broad cache mode classified %+v", got)
	}
	pending := fmt.Sprintf(`{"schema":1,"state":"pending","tree":%q,"oracle":%q,"started_at":"2026-07-13T11:59:00Z","owner_pid":123}`, s.Tree, s.Oracle)
	write(pending, 0o600)
	if got := inspectAt(root, now); got.State != Pending || got.PendingStatus != "interrupted-pending" {
		t.Fatalf("released pending classified %+v", got)
	}
}

func TestExecutionLockAndDriftFailClosed(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\ntouch .git/started\nwhile [ ! -f .git/release ]; do sleep .01; done\nprintf drift > drift.txt\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	done := make(chan Result, 1)
	go func() { done <- Execute(context.Background(), root, io.Discard, io.Discard) }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(root, ".git", "started")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("gate did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if got := Inspect(root); got.State != Pending || got.PendingStatus != "locked-pending" {
		t.Fatalf("live owner not projected: %+v", got)
	}
	if second := Execute(context.Background(), root, io.Discard, io.Discard); second.ActionExit != 1 {
		t.Fatalf("concurrent execution did not fail immediately: %+v", second)
	}
	writeGateTestFile(t, root, ".git/release", "\n", 0o600)
	first := <-done
	if first.ActionExit != 1 || first.GateExit != 0 || first.Inspection.State != Pending {
		t.Fatalf("mid-run drift did not leave pending operational failure: %+v", first)
	}
}

func TestUnlockErrorClearsSameProcessOwnership(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bench-gate.lock")
	engine := productionGateEngine{}
	first, err := engine.OpenLock(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Acquire(first); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		executionLockOwners.Lock()
		delete(executionLockOwners.paths, path)
		executionLockOwners.Unlock()
	})
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	if err := engine.Unlock(first); err == nil {
		t.Fatal("unlock on a closed descriptor succeeded, want EBADF")
	}
	second, err := engine.OpenLock(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	if err := engine.Acquire(second); err != nil {
		t.Fatalf("same-process reacquire after unlock error: %v", err)
	}
	defer engine.Unlock(second)
}

func TestNoGateReturnsThreeWithoutRecord(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", "")
	if err := os.Remove(filepath.Join(root, ".bench", "gate.sh")); err != nil {
		t.Fatal(err)
	}
	result := Execute(context.Background(), root, io.Discard, io.Discard)
	if result.GateExit != 3 || result.ActionExit != 3 {
		t.Fatalf("no gate result = %+v", result)
	}
	if _, err := os.Stat(filepath.Join(root, ".git", benchgit.GateCacheFile)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("no gate wrote verdict: %v", err)
	}
}

// TestPinCommandNotInRepo injects the terminal precondition so the shared
// not-in-repo branch remains reachable through this in-process seam.
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
