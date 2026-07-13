package gate

// The runner engine's core behavior: concurrency, prefixed output shape, exit
// codes, aggregation, and cancellation — plus the helpers shared by the serial
// (runner_serial_test.go) and optional/start (runner_optional_test.go) families.
// The PhasesCommand surface (signal handling, the phase table, shellcheck
// expansion) is tested in phases_command_test.go.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/toon"
)

func TestRunnerRunsPhasesConcurrently(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("alpha", "printf 'alpha\\n'; sleep 0.25"),
		fakePhase("bravo", "printf 'bravo\\n'; sleep 0.25"),
		fakePhase("charlie", "printf 'charlie\\n'; sleep 0.25"),
		fakePhase("delta", "printf 'delta\\n'; sleep 0.25"),
	}

	var stdout, stderr bytes.Buffer
	start := time.Now()
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	elapsed := time.Since(start)

	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
	}
	if elapsed >= 750*time.Millisecond {
		t.Fatalf("runPhases took %v, want concurrent execution under serial sum", elapsed)
	}
	out := stdout.String() + stderr.String()
	for _, marker := range []string{"[alpha] alpha", "[bravo] bravo", "[charlie] charlie", "[delta] delta"} {
		if !strings.Contains(out, marker) {
			t.Fatalf("output missing %q:\n%s", marker, out)
		}
	}
}

func TestRunnerPrefixesAndKeepsLinesIntact(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("alpha", "printf 'alpha'; sleep 0.05; printf ' one\\nalpha two\\n'"),
		fakePhase("bravo", "printf 'bravo'; sleep 0.05; printf ' one\\nbravo two\\n'"),
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	got := map[string]bool{}
	for _, line := range lines {
		if strings.HasPrefix(line, "phase ") || strings.HasPrefix(line, "gate: ") {
			continue
		}
		got[line] = true
		if !strings.HasPrefix(line, "[alpha] ") && !strings.HasPrefix(line, "[bravo] ") {
			t.Fatalf("unprefixed phase line %q in output:\n%s", line, stdout.String())
		}
	}
	for _, want := range []string{
		"[alpha] alpha one",
		"[alpha] alpha two",
		"[bravo] bravo one",
		"[bravo] bravo two",
	} {
		if !got[want] {
			t.Fatalf("missing intact line %q; got %#v\nfull output:\n%s", want, got, stdout.String())
		}
	}
}

func TestRunnerFinalLineAndExitCodes(t *testing.T) {
	t.Run("green", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), t.TempDir(), []Phase{fakePhase("ok", "true")}, outerMode, &stdout, &stderr)
		if rc != 0 {
			t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
		}
		if !strings.HasSuffix(stdout.String(), "gate: green\n") {
			t.Fatalf("stdout final line = %q, want gate: green", stdout.String())
		}
	})

	t.Run("red", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), t.TempDir(), []Phase{fakePhase("bad", "exit 7")}, outerMode, &stdout, &stderr)
		if rc != 1 {
			t.Fatalf("runPhases rc = %d, want 1", rc)
		}
		if !strings.HasSuffix(stderr.String(), "gate: red\n") {
			t.Fatalf("stderr final line = %q, want gate: red", stderr.String())
		}
	})

	t.Run("not in repo", func(t *testing.T) {
		cwd := mustGetwd(t)
		t.Cleanup(func() { mustChdir(t, cwd) })
		mustChdir(t, t.TempDir())

		var stdout, stderr bytes.Buffer
		rc := PhasesCommand(nil, &stdout, &stderr)
		if rc != 3 {
			t.Fatalf("PhasesCommand rc = %d, want 3; stderr:\n%s", rc, stderr.String())
		}
		if strings.TrimSpace(stderr.String()) != toon.NotInRepo() {
			t.Fatalf("stderr = %q, want gate not-in-repo line", stderr.String())
		}
	})
}

func TestRunnerAggregatesAllPhases(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("first", "printf 'first\\n'"),
		fakePhase("bad", "printf 'bad\\n'; exit 1"),
		fakePhase("third", "printf 'third\\n'"),
		fakePhase("fourth", "printf 'fourth\\n'"),
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	out := stdout.String() + stderr.String()
	for _, marker := range []string{"[first] first", "[bad] bad", "[third] third", "[fourth] fourth"} {
		if !strings.Contains(out, marker) {
			t.Fatalf("output missing %q:\n%s", marker, out)
		}
	}
}

func TestRunnerInnerModeByteShape(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("conformance", "printf 'conformance\\n'"),
		fakePhase("contract", "printf 'contract\\n'"),
		fakePhase("shellcheck", "printf 'shellcheck\\n'"),
		fakePhase("canary", "printf 'canary\\n'"),
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phasesForMode(phases, innerMode), innerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stderr:\n%s", rc, stderr.String())
	}
	if got, want := stdout.String(), "conformance\ncontract\nshellcheck\ngate: green\n"; got != want {
		t.Fatalf("inner output = %q, want %q", got, want)
	}
	if strings.Contains(stdout.String()+stderr.String(), "[") {
		t.Fatalf("inner output gained outer prefix bytes:\nstdout=%q\nstderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunnerCancelKillsGroup(t *testing.T) {
	root := t.TempDir()
	pidfile := filepath.Join(root, "sleep.pid")
	phase := Phase{
		Name: "slow",
		Argv: []string{"bash", "-c", `sleep 30 & echo $! > "$1"; wait`, "bash", pidfile},
	}
	ctx, cancel := context.WithCancel(context.Background())
	var stdout, stderr bytes.Buffer
	done := make(chan int, 1)
	go func() {
		done <- runPhases(ctx, root, []Phase{phase}, outerMode, &stdout, &stderr)
	}()

	pid := waitForPIDFile(t, pidfile)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	cancel()

	select {
	case rc := <-done:
		if rc != 130 {
			t.Fatalf("runPhases rc = %d, want 130; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runPhases did not return after cancellation")
	}
	waitForProcessExit(t, pid)
}

func TestRunnerRootWithSpace(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "root with space")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatalf("mkdir root with space: %v", err)
	}
	phase := Phase{
		Name: "pwd",
		Argv: []string{"bash", "-c", `test "$PWD" = "$1"; printf 'root=%s\n' "$1"`, "bash", root},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "[pwd] root="+root) {
		t.Fatalf("stdout does not show literal root argv:\n%s", stdout.String())
	}
}

func fakePhase(name, script string) Phase {
	return Phase{Name: name, Argv: []string{"bash", "-c", script}}
}

func mustGetwd(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	return cwd
}

func mustChdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
}

func waitForPIDFile(t *testing.T, path string) int {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(path)
		if err == nil {
			var pid int
			if _, scanErr := fmt.Sscanf(strings.TrimSpace(string(data)), "%d", &pid); scanErr == nil && pid > 0 {
				return pid
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("pid file %s was not written", path)
	return 0
}

func waitForProcessExit(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("process %d still exists after cancellation", pid)
}

func TestR17PrivateFaultBridge(t *testing.T) {
	op := os.Getenv("FT78_R17_FAULT")
	valid := map[string]bool{
		"lock-open": true, "lock-acquisition": true, "temporary-create": true,
		"mode-establishment": true, "write": true, "file-sync": true,
		"file-close": true, "atomic-rename": true, "directory-open": true,
		"directory-sync": true, "directory-close": true, "post-run-subject-rebuild": true,
	}
	if op == "" {
		t.Skip("private FT78 bridge")
	}
	if !valid[op] {
		t.Fatalf("unknown R17 fault %q", op)
	}
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	seed := Execute(context.Background(), root, io.Discard, io.Discard)
	if seed.ActionExit != 0 || !seed.Inspection.ReusableGreen {
		t.Fatalf("seed = %+v, want reusable green", seed)
	}
	durable := "ready-green"
	wantAttempts := 2
	if op == "lock-open" || op == "lock-acquisition" || op == "post-run-subject-rebuild" {
		wantAttempts = 1
	}
	for call := 1; call <= 2; call++ {
		engine := &faultEngine{now: time.Now().UTC().Truncate(time.Second), failOp: op}
		got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
		if got.GateExit != 0 || got.ActionExit != 1 || got.Inspection.ReusableGreen || !engine.failed || engine.opCounts[op] != wantAttempts {
			t.Fatalf("call %d tuple = gate:%d action:%d reusable:%v failed:%v hits:%d trace:%v", call, got.GateExit, got.ActionExit, got.Inspection.ReusableGreen, engine.failed, engine.opCounts[op], engine.trace)
		}
	}
	if got := Inspect(root); got.State == Pending && !got.ReusableGreen {
		durable = "interrupted-pending"
	} else if got.State != Ready || got.Status != "green" || !got.ReusableGreen || (op != "lock-open" && op != "lock-acquisition") {
		t.Fatalf("durable tuple = %+v for %s", got, op)
	}
	temps, err := filepath.Glob(filepath.Join(filepath.Dir(cachePath(t, root)), ".bench-last-gate-*"))
	if err != nil || len(temps) != 0 {
		t.Fatalf("temporary evidence = %v/%v, want none", temps, err)
	}
	fmt.Printf("R17-TUPLE op=%s calls=2 gate=0 action=1 returned_reusable=false durable=%s attempts=%d,%d temps=0\n", op, durable, wantAttempts, wantAttempts)
}
