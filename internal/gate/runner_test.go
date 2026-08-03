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
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
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
		if strings.HasPrefix(line, "phase ") || strings.HasPrefix(line, "gate: ") || strings.HasPrefix(line, skipRowPrefix) {
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

// TestRunnerSummaryLineByteShape pins the outer summary vocabulary byte for byte. The
// canary sweep matches EXPECT substrings against gate output and downstream parsers key
// on these lines, so a reworded or dropped summary is a contract break that no other
// test in this package would notice.
func TestRunnerSummaryLineByteShape(t *testing.T) {
	root := t.TempDir()
	phases := []Phase{
		fakePhase("green-phase", "true"),
		fakePhase("red-phase", "exit 4"),
		{Name: "absent-phase", Argv: []string{"definitely-absent-binary-for-bench-summary-test"}, Optional: true},
		{Name: "blocked-phase", Argv: []string{"bash", "-c", "true"}, Needs: []string{"red-phase"}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"phase green-phase: green\n",
		"phase red-phase: red (exit 4)\n",
		"phase absent-phase: skipped (not installed)\n",
		"phase blocked-phase: skipped (needs red-phase)\n",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("summary line %q missing from stdout:\n%s", want, stdout.String())
		}
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

// TestRunnerPhaseDirIsAbsoluteOrRoot pins Phase.Dir's one rule where the process is
// launched: an absolute directory is used as it stands, and an empty one means the
// runner's root. Anchoring a declared path belongs to whoever produced the phase.
func TestRunnerPhaseDirIsAbsoluteOrRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	printCwd := []string{"bash", "-c", `printf 'cwd=%s\n' "$PWD"`}
	phases := []Phase{
		{Name: "rooted", Argv: printCwd},
		{Name: "nested", Argv: printCwd, Dir: filepath.Join(root, "sub")},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"[rooted] cwd=" + root + "\n",
		"[nested] cwd=" + filepath.Join(root, "sub") + "\n",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("phase cwd line %q missing:\n%s", want, stdout.String())
		}
	}
}

// TestRunnerPhaseDirUnusableIsRed grades a working directory the phase cannot enter.
// chdir fails ENOENT exactly as a missing binary does, so without its own check an
// optional phase would answer a directory typo with "not installed" and quietly take
// its check off the gate.
func TestRunnerPhaseDirUnusableIsRed(t *testing.T) {
	root := t.TempDir()
	regular := filepath.Join(root, "regular-file")
	writeFile(t, regular, "not a directory\n")
	for _, tc := range []struct {
		name string
		dir  string
	}{
		{name: "missing", dir: filepath.Join(root, "absent-subdir")},
		{name: "not-a-directory", dir: regular},
	} {
		t.Run(tc.name, func(t *testing.T) {
			phases := []Phase{{
				Name:     "opt",
				Argv:     []string{"definitely-absent-binary-for-bench-dir-test"},
				Optional: true,
				Dir:      tc.dir,
			}}

			var stdout, stderr bytes.Buffer
			rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
			if rc != 1 {
				t.Fatalf("runPhases rc = %d, want 1; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
			}
			out := stdout.String() + stderr.String()
			if strings.Contains(out, "skipped") {
				t.Fatalf("an unusable working directory was masked as a skip:\n%s", out)
			}
			if !strings.Contains(out, "phase opt: red (unusable working directory "+tc.dir) || !strings.Contains(out, "gate: red") {
				t.Fatalf("red verdict does not name the unusable working directory:\n%s", out)
			}
		})
	}
}

// TestMergeEnvStripsThenSets grades the merge rule where it is decidable. os/exec
// collapses a repeated key on its way to the child in favor of the later value, so no
// subprocess can distinguish strip-then-set from plain appending — the slice this
// function returns is the last place the difference is visible.
func TestMergeEnvStripsThenSets(t *testing.T) {
	base := []string{"KEEP=base", "REPLACED=base", "REPLACED=base-again", "PREFIXY_KEY=base"}
	got := mergeEnv(base, []string{"REPLACED=phase", "ADDED=phase"})
	want := []string{"KEEP=base", "PREFIXY_KEY=base", "REPLACED=phase", "ADDED=phase"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mergeEnv = %#v, want %#v", got, want)
	}
	if merged := mergeEnv(base, nil); !reflect.DeepEqual(merged, base) {
		t.Fatalf("mergeEnv with no overrides = %#v, want the base unchanged", merged)
	}
}

// TestRunnerPhaseEnvStripsThenSets pins what the child actually receives for an
// overridden key. It execs env(1) directly rather than through a shell: a shell folds a
// repeated key into one variable as it imports the environment, so only the raw environ
// block can show what the phase was handed.
func TestRunnerPhaseEnvStripsThenSets(t *testing.T) {
	const key = "BENCH_PHASE_ENV_PROBE"
	t.Setenv(key, "from-gate-env")
	root := t.TempDir()
	phase := Phase{Name: "env", Argv: []string{"env"}, Env: []string{key + "=from-phase"}}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	var got []string
	for _, line := range strings.Split(stdout.String(), "\n") {
		if strings.HasPrefix(line, "[env] "+key+"=") {
			got = append(got, strings.TrimPrefix(line, "[env] "))
		}
	}
	if len(got) != 1 || got[0] != key+"=from-phase" {
		t.Fatalf("child environ carries %#v for %s, want exactly one entry with the phase's value", got, key)
	}
}

func TestRunnerTransportsGateOwnedCanarySelection(t *testing.T) {
	const selection = "line-routing,package-core-guard"
	phase := Phase{
		Name:           "canary",
		Argv:           []string{"bash", "-c", `proof=$(cat <&"$BENCH_CANARY_FAMILIES_FD"); test "$proof" = "$BENCH_CANARY_FAMILIES" && printf '%s\n' "$proof"`},
		Env:            []string{canary.FamilySelectionEnv + "=" + selection, canary.FamilySelectionOwnerEnv + "=gate"},
		canaryFamilies: []string{"line-routing", "package-core-guard"},
	}

	var stdout, stderr bytes.Buffer
	result := runPhase(context.Background(), t.TempDir(), phase, &stdout, &stderr)
	if result.Code != 0 || result.StartErr != nil {
		t.Fatalf("runPhase = %+v; stdout=%q stderr=%q", result, stdout.String(), stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != selection {
		t.Fatalf("authority payload = %q, want %q", got, selection)
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

func TestExecuteCancellationKillsDescendantAfterLeaderExits(t *testing.T) {
	root := gateTestRepo(t, `#!/usr/bin/env bash
(
  trap '' INT TERM
  exec >/dev/null 2>&1
  while :; do sleep 1; done
) &
child=$!
echo "$child" > .git/stubborn-child
trap 'exit 130' INT
wait
`, `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan Result, 1)
	go func() { done <- Execute(ctx, root, io.Discard, io.Discard) }()
	child := waitForPIDFile(t, filepath.Join(root, ".git", "stubborn-child"))
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	cancel()
	select {
	case got := <-done:
		if got.ActionExit != 130 {
			t.Fatalf("cancelled execution = %+v, want action exit 130", got)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("cancelled execution did not return")
	}
	waitForProcessExit(t, child)
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
		capability.Environment(t, "private FT78 bridge")
	}
	if !valid[op] {
		t.Fatalf("unknown R17 fault %q", op)
	}
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	// The seed is recorded a full freshness window in the past. Faults after lock
	// acquisition are reached only past the reuse short-circuit, which a fresh green
	// answers without ever running them, so they are driven at the present clock where
	// the seed has expired. The pre-acquire pair is driven at the seed's own clock,
	// where the green is still reusable and its interrupted-pending demotion fires.
	seedNow := time.Now().UTC().Truncate(time.Second).Add(-freshness - time.Minute)
	seed := executeWithEngine(context.Background(), root, io.Discard, io.Discard, &faultEngine{now: seedNow})
	if seed.ActionExit != 0 || !seed.Inspection.ReusableGreen {
		t.Fatalf("seed = %+v, want reusable green", seed)
	}
	durable := "ready-green"
	wantAttempts := 2
	faultNow := time.Now().UTC().Truncate(time.Second)
	if op == "lock-open" || op == "lock-acquisition" {
		wantAttempts, faultNow = 1, seedNow
	} else if op == "post-run-subject-rebuild" {
		wantAttempts = 1
	}
	for call := 1; call <= 2; call++ {
		engine := &faultEngine{now: faultNow, failOp: op}
		got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
		if got.GateExit != 0 || got.ActionExit != 1 || got.Inspection.ReusableGreen || !engine.failed || engine.opCounts[op] != wantAttempts {
			t.Fatalf("call %d tuple = gate:%d action:%d reusable:%v failed:%v hits:%d trace:%v", call, got.GateExit, got.ActionExit, got.Inspection.ReusableGreen, engine.failed, engine.opCounts[op], engine.trace)
		}
	}
	if got := Inspect(root); got.State == Pending && !got.ReusableGreen {
		durable = "interrupted-pending"
	} else {
		t.Fatalf("durable tuple = %+v for %s", got, op)
	}
	temps, err := filepath.Glob(filepath.Join(filepath.Dir(cachePath(t, root)), ".bench-last-gate-*"))
	if err != nil || len(temps) != 0 {
		t.Fatalf("temporary evidence = %v/%v, want none", temps, err)
	}
	fmt.Printf("R17-TUPLE op=%s calls=2 gate=0 action=1 returned_reusable=false durable=%s attempts=%d,%d temps=0\n", op, durable, wantAttempts, wantAttempts)
}

func prepareSymlinkMutationProof(t *testing.T, f contract.Fixture) {
	t.Helper()
	f.WriteFile("inputs/target", "green\n")
	if err := os.Symlink("target", filepath.Join(f.Root, "inputs", "link-b")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("link-b", filepath.Join(f.Root, "inputs", "link-a")); err != nil {
		t.Fatal(err)
	}
	f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\necho run >> .git/ft78-runs\ntest \"$(cat inputs/link-a)\" = green\n")
	f.WriteFile("work.txt", "changed\n")
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
}

func finishSymlinkMutationProof(t *testing.T, f contract.Fixture) {
	t.Helper()
	before := story3ReadVerdict(t, f)
	head := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	f.WriteFile("inputs/target", "red\n")
	f.Bench("commit", "-m", "must not commit stale symlink evidence", "work.txt").RequireExit(1)
	after := story3ReadVerdict(t, f)
	if before.Oracle == after.Oracle {
		t.Fatal("resolved symlink target mutation did not change oracle identity")
	}
	if got := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout); got != head {
		t.Fatalf("HEAD = %s, want unchanged %s after red gate", got, head)
	}
	lines := contract.NonEmptyLines(contract.ReadFileAbs(t, filepath.Join(story3GitDir(f), "ft78-runs")))
	if len(lines) != 2 {
		t.Fatalf("gate runs = %d, want 2 after symlink target mutation", len(lines))
	}
	if got := Inspect(f.Root); got.State != Ready || got.Status != "red" || got.ReusableGreen {
		t.Fatalf("mutated symlink target inspection = %+v, want non-reusable ready red", got)
	}
}
