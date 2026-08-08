package gate

// The needs-edge contract: a phase waits for its needs regardless of table position,
// a need's red or skip skips its dependents with cause while unrelated phases still
// run, and an interrupt is a 130 that launches nothing further, never a graded red —
// in both outer and inner modes.

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

func TestRunnerShellcheckAbsentSkips(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	phase := Phase{Name: "shellcheck", Argv: []string{"definitely-not-installed-shellcheck-for-bench-test"}, Optional: true}
	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, want skip to stay green; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if strings.Contains(out, "gate: red") || strings.Contains(out, "shellcheck reported issues") {
		t.Fatalf("optional missing shellcheck looked red:\n%s", out)
	}
	if !strings.Contains(out, "phase shellcheck: skipped") {
		t.Fatalf("missing shellcheck skip summary:\n%s", out)
	}
}

func r12Contention(id, action string) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := gateTestRepo(t, "#!/usr/bin/env bash\ngitdir=\"$(git rev-parse --absolute-git-dir)\"\nprintf 'run\\n' >> \"$gitdir/owner-runs\"\n[ \"$(wc -l < \"$gitdir/owner-runs\")\" -eq 1 ] || exit 0\ntouch \"$gitdir/owner-started\"\nwhile [ ! -f \"$gitdir/release-owner\" ]; do sleep .01; done\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
		_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".bench-contract-env/\n"), 0o644)
		_, _ = benchgit.Output("-C", root, "add", "-A")
		_, _ = benchgit.Output("-C", root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-m", "base")
		f := contract.NewFixtureAt(t, root, contract.IsolatedEnv(t, root))
		kit := kitRootForTest(t)
		bench := filepath.Join(kit, "bin", "bench.sh")
		adapter, ownerRoot := filepath.Join(root, "adapter"), root
		if action == "commit" {
			_ = os.WriteFile(filepath.Join(root, "work"), []byte("x"), 0o644)
		}
		if action == "gate" || action == "stop" {
			_ = os.WriteFile(filepath.Join(root, "probe-work"), []byte("dirty"), 0o644)
		}
		if action == "shift" {
			_ = os.WriteFile(adapter, []byte("#!/bin/sh\nprintf dirty > shift-work\n"), 0o755)
			_, _ = benchgit.Output("-C", root, "add", "adapter")
			_, _ = benchgit.Output("-C", root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-m", "adapter")
			pool := strings.TrimSpace(contract.RunAt(t, f, root, nil, "bash", bench, "worktree-pool", root).Stdout)
			ownerRoot = filepath.Join(pool, "warm")
			_ = os.MkdirAll(pool, 0o700)
			_, _ = benchgit.Output("-C", root, "worktree", "add", "--detach", ownerRoot, "HEAD")
		}
		if action == "stop" {
			binary := currentBenchBinary(t)
			_ = os.MkdirAll(filepath.Join(root, "bin"), 0o755)
			_ = os.MkdirAll(filepath.Join(root, "dist"), 0o755)
			_ = os.WriteFile(filepath.Join(root, "bin", "bench.sh"), mustRead(t, bench), 0o755)
			_ = os.WriteFile(filepath.Join(root, "dist", "bench"), mustRead(t, binary), 0o755)
		}
		done := make(chan Result, 1)
		go func() {
			done <- executeWithEngineAtKit(context.Background(), ownerRoot, ownerRoot, io.Discard, io.Discard, productionGateEngine{})
		}()
		ownerGitDir := filepath.Dir(cachePath(t, ownerRoot))
		waitFile(t, filepath.Join(ownerGitDir, "owner-started"))
		defer func() {
			_ = os.WriteFile(filepath.Join(ownerGitDir, "release-owner"), nil, 0o644)
			select {
			case got := <-done:
				// FT79: a blocked shift used to leave its adapter's dirty write behind
				// on the shared pool worktree (the old red-gate "retain verbatim"
				// special case), which the owner's later commit phase then picked up
				// unexpectedly — the only reason this wanted action exit 1 for "shift".
				// The uniform snapshot-and-release rule now leaves the worktree clean
				// before the owner's action phase runs, so every action's owner exits 0.
				if got.ActionExit != 0 {
					t.Errorf("owner result = %+v, want action 0", got)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("timed out releasing owner")
			}
		}()
		cacheBefore := mustRead(t, cachePath(t, ownerRoot))
		infoBefore, _ := os.Stat(cachePath(t, ownerRoot))
		headBefore, _ := benchgit.Output("-C", root, "rev-parse", "HEAD")
		indexBefore, _ := benchgit.Output("-C", root, "write-tree")
		statusBefore, _ := benchgit.Output("-C", root, "status", "--porcelain")
		branchesBefore, _ := benchgit.Output("-C", root, "for-each-ref", "--format=%(refname)", "refs/heads")
		var probe contract.Probe
		switch action {
		case "gate":
			probe = contract.RunAt(t, f, root, nil, "bash", bench, "gate")
		case "commit":
			probe = contract.RunAt(t, f, root, nil, "bash", bench, "commit", "-m", "blocked", "work")
		case "shift":
			probe = contract.RunAtWithTimeout(t, f, root, map[string]string{"BENCH_AGENT": adapter, "BENCH_MAX_ITERS": "1"}, 5*time.Second, "bash", bench, "shift", "blocked")
		case "stop":
			probe = contract.RunAtWithInput(t, f, root, map[string]string{"BENCH_SHIFT": "1"}, "{}\n", "bash", filepath.Join(kit, ".bench", "hooks", "stop.sh"))
		}
		if probe.TimedOut || probe.ExitCode == 0 || !strings.Contains(probe.Stdout+probe.Stderr, "gate execution already in progress") {
			t.Fatalf("%s contention = exit %d\n%s%s", action, probe.ExitCode, probe.Stdout, probe.Stderr)
		}
		if runs := strings.Count(string(mustRead(t, filepath.Join(ownerGitDir, "owner-runs"))), "run\n"); runs != 1 {
			t.Fatalf("%s gate runs = %d, want 1", action, runs)
		}
		infoAfter, _ := os.Stat(cachePath(t, ownerRoot))
		if !bytes.Equal(mustRead(t, cachePath(t, ownerRoot)), cacheBefore) || infoAfter.Mode() != infoBefore.Mode() || !infoAfter.ModTime().Equal(infoBefore.ModTime()) {
			t.Fatalf("%s rewrote pending evidence", action)
		}
		headAfter, _ := benchgit.Output("-C", root, "rev-parse", "HEAD")
		indexAfter, _ := benchgit.Output("-C", root, "write-tree")
		statusAfterRoot, _ := benchgit.Output("-C", root, "status", "--porcelain")
		if headAfter != headBefore || indexAfter != indexBefore || statusAfterRoot != statusBefore {
			t.Fatalf("blocked %s changed HEAD/index/worktree", action)
		}
		if action == "commit" && string(mustRead(t, filepath.Join(root, "work"))) != "x" {
			t.Fatalf("blocked commit changed named work bytes")
		}
		if action == "shift" {
			worktreesAfter, _ := benchgit.Output("-C", root, "worktree", "list", "--porcelain")
			branchesAfter, _ := benchgit.Output("-C", root, "for-each-ref", "--format=%(refname)", "refs/heads")
			if !strings.Contains(worktreesAfter, ownerRoot) || branchesAfter == branchesBefore {
				t.Fatalf("blocked shift discarded branch/registration")
			}
			// FT79: gate contention is a red-gate-style post-mutation failure — the
			// uniform rule snapshots the dirty tree to refs/bench/recovery/<branch> and
			// releases the pool worktree, rather than retaining it (and its lease)
			// verbatim as before.
			recoveryRefs, _ := benchgit.Output("-C", root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/")
			if !strings.Contains(recoveryRefs, "refs/bench/recovery/") {
				t.Fatalf("blocked shift did not snapshot its work to a recovery ref")
			}
			ref := strings.TrimSpace(strings.SplitN(recoveryRefs, "\n", 2)[0])
			tree, _ := benchgit.Output("-C", root, "ls-tree", "-r", "--name-only", ref)
			if !strings.Contains(tree, "shift-work") {
				t.Fatalf("blocked shift's recovery snapshot did not preserve shift-work:\n%s", tree)
			}
			lease, _ := benchgit.Output("-C", ownerRoot, "rev-parse", "--git-path", "bench-lease")
			if _, err := os.Stat(lease); err == nil {
				t.Fatalf("blocked shift did not release its pool lease")
			}
		}
	}}
}

func currentBenchBinary(t *testing.T) string {
	t.Helper()
	kit := kitRootForTest(t)
	binary := filepath.Join(t.TempDir(), "bench")
	cmd := exec.Command("bash", filepath.Join(kit, "scripts", "go-build.sh"), kit, binary)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build current bench binary: %v\n%s", err, output)
	}
	return binary
}

func TestRunnerOptionalBrokenSymlinkSkips(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	interpreters := filepath.Join(root, "interpreters")
	if err := os.Mkdir(interpreters, 0o755); err != nil {
		t.Fatalf("mkdir interpreters: %v", err)
	}
	brokenInterpreter := filepath.Join(interpreters, "sh")
	if err := os.Symlink(filepath.Join(root, "missing-sh"), brokenInterpreter); err != nil {
		t.Fatalf("symlink broken interpreter: %v", err)
	}
	shellcheck := filepath.Join(bin, "shellcheck")
	if err := os.WriteFile(shellcheck, []byte("#!"+brokenInterpreter+"\n"), 0o755); err != nil {
		t.Fatalf("write shellcheck shim: %v", err)
	}
	t.Setenv("PATH", bin)
	phase := Phase{Name: "shellcheck", Argv: []string{"shellcheck"}, Optional: true}
	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, want broken optional executable skipped; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if out := stdout.String() + stderr.String(); !strings.Contains(out, "phase shellcheck: skipped") || strings.Contains(out, "gate: red") {
		t.Fatalf("broken optional executable did not produce a green skip:\n%s", out)
	}
}

func TestRunnerShellcheckAbsentSkipVerdictNamesNotInstalled(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	phase := Phase{Name: "shellcheck", Argv: []string{"definitely-absent-shellcheck-binary-for-bench-test"}, Optional: true}
	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("absent optional binary went red rc=%d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if out := stdout.String() + stderr.String(); !strings.Contains(out, "phase shellcheck: skipped (not installed)") {
		t.Fatalf("absent shellcheck skip verdict is silent about why (want 'skipped (not installed)'):\n%s", out)
	}
}

func TestRunnerOptionalUnexecutableStubGoesRed(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	stub := filepath.Join(bin, "shellcheck")
	if err := os.WriteFile(stub, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o000); err != nil {
		t.Fatalf("write unexecutable stub: %v", err)
	}
	t.Setenv("PATH", bin)
	phase := Phase{Name: "shellcheck", Argv: []string{"shellcheck"}, Optional: true}
	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("present-but-unexecutable shellcheck did not go red rc=%d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if strings.Contains(out, "skipped") {
		t.Fatalf("present-but-unexecutable shellcheck was masked as a skip:\n%s", out)
	}
	if !strings.Contains(out, "phase shellcheck: red") || !strings.Contains(out, "permission denied") {
		t.Fatalf("red verdict does not name the exec failure:\n%s", out)
	}
}

func TestRunnerRequiredStartFailureRed(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	phase := Phase{Name: "required", Argv: []string{filepath.Join(root, "missing-required")}}
	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want required start failure red; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if !strings.Contains(out, "phase required: red") || !strings.Contains(out, "gate: red") {
		t.Fatalf("required start failure did not stay red:\n%s", out)
	}
}

func TestRunnerNeededPhaseRedSkipsDependentsInnerMode(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "build", Argv: []string{"bash", "-c", "exit 3"}},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}, Needs: []string{"build"}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	if !strings.HasSuffix(stderr.String(), "gate: red\n") {
		t.Fatalf("stderr final line = %q, want gate: red", stderr.String())
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("inner-mode phase behind a red need still ran")
	}
}

func TestRunnerCancelDuringNeededPhaseReturns130(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	pidfile := filepath.Join(root, "sleep.pid")
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{
			Name: "build",
			Argv: []string{"bash", "-c", `sleep 30 & echo $! > "$1"; wait`, "bash", pidfile},
		},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}, Needs: []string{"build"}},
	}
	ctx, cancel := context.WithCancel(withProcessGroupCancelGrace(context.Background(), fastProcessGroupCancelGrace))
	var stdout, stderr bytes.Buffer
	done := make(chan int, 1)
	go func() {
		done <- runPhases(ctx, root, phases, outerMode, &stdout, &stderr)
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
		t.Fatal("runPhases did not return after cancellation during the needed phase")
	}
	waitForProcessExit(t, pid)
	// An interrupt is not a verdict: no phase behind the interrupted build may run,
	// and the run must not present itself as a graded red.
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("dependent phase ran after its need was interrupted")
	}
	if strings.Contains(stdout.String()+stderr.String(), "gate: red") {
		t.Fatalf("interrupted needed phase was graded as red:\nstdout=%q\nstderr=%q", stdout.String(), stderr.String())
	}
}

// TestRunnerUnsatisfiableGraphIsRed pins the fail-closed rule for a table whose edges no
// run can satisfy. The loader refuses a cycle, but the built-in table and any injected
// one reach the scheduler unchecked, and a run that executed nothing at all must not be
// able to report green.
func TestRunnerUnsatisfiableGraphIsRed(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}, Needs: []string{"bravo"}},
		{Name: "bravo", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}, Needs: []string{"alpha"}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1 for a graph that can execute nothing; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if strings.Contains(stdout.String(), ": green") {
		t.Fatalf("a phase that never launched reported green:\n%s", stdout.String())
	}
	out := stdout.String() + stderr.String()
	for _, want := range []string{
		"phase alpha: red (stuck behind unsatisfied need bravo)\n",
		"phase bravo: red (stuck behind unsatisfied need alpha)\n",
		"gate: red\n",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("unsatisfiable graph output missing %q:\n%s", want, out)
		}
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("a phase on an unsatisfiable edge still ran")
	}

	var innerOut, innerErr bytes.Buffer
	if rc := runPhases(context.Background(), root, phases, innerMode, &innerOut, &innerErr); rc != 1 {
		t.Fatalf("inner runPhases rc = %d, want 1; stdout=%q stderr=%q", rc, innerOut.String(), innerErr.String())
	}
}

// TestRunnerPhaseExit130IsRedNotCancellation separates the two things exit 130 can mean.
// Only the context decides that a run was cancelled; a phase that chooses 130 is an
// ordinary red, and reading it as an interrupt suppresses the summaries and the gate
// line and leaves the verdict recorded as pending rather than red.
func TestRunnerPhaseExit130IsRedNotCancellation(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	phases := []Phase{
		fakePhase("self130", "exit 130"),
		fakePhase("other", "true"),
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if strings.Contains(out, "gate: cancelled") {
		t.Fatalf("a phase's own exit 130 was read as an interrupt:\n%s", out)
	}
	for _, want := range []string{
		"phase self130: red (exit 130)\n",
		"phase other: green\n",
		"gate: red\n",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("self-inflicted 130 output missing %q:\n%s", want, out)
		}
	}
}

func TestRunnerNeededPhaseNotFirstStillRunsFirst(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	marker := filepath.Join(root, "built")
	phases := []Phase{
		// Readers listed ahead of their need: only the edges keep them from
		// starting before the marker exists.
		{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}, Needs: []string{"build"}},
		{Name: "bravo", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}, Needs: []string{"build"}},
		{
			Name: "build",
			Argv: []string{"bash", "-c", `sleep 0.1; touch "$1"`, "bash", marker},
		},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; a reader phase ran before its later-declared need\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase build: green") {
		t.Fatalf("needed phase missing from summaries:\n%s", stdout.String())
	}
}

func TestSchedulerRespectsNeeds(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	marker := filepath.Join(root, "made")
	phases := []Phase{
		{Name: "maker", Argv: []string{"bash", "-c", `sleep 0.1; touch "$1"`, "bash", marker}},
		// Red unless the edge is honored: ignoring it starts the reader ~100ms
		// before the marker exists.
		{Name: "reader", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}, Needs: []string{"maker"}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; the reader started before its need completed\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase reader: green\n") {
		t.Fatalf("dependent phase missing from summaries:\n%s", stdout.String())
	}
}

// TestSchedulerOverlapsIndependents holds four edge-free phases at a barrier that only
// opens once all four subprocesses exist at the same time, so a serializing or
// fixed-width scheduler starves rather than merely running slower. The waiters are
// syscall-blocked subprocesses, not goroutines, so the barrier is GOMAXPROCS-safe.
func TestSchedulerOverlapsIndependents(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	started := filepath.Join(root, "started")
	if err := os.Mkdir(started, 0o755); err != nil {
		t.Fatalf("mkdir started: %v", err)
	}
	const barrier = `d="$1"; touch "$d/$2"; for _ in $(seq 500); do [ "$(ls "$d" | wc -l)" -ge 4 ] && exit 0; sleep 0.01; done; exit 1`
	var phases []Phase
	for _, name := range []string{"alpha", "bravo", "charlie", "delta"} {
		phases = append(phases, Phase{Name: name, Argv: []string{"bash", "-c", barrier, "bash", started, name}})
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; four independent phases never overlapped\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
}

func TestSchedulerSkipsDependentsOfRed(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	leak, ran := filepath.Join(root, "leak"), filepath.Join(root, "ran")
	phases := []Phase{
		fakePhase("build", "exit 3"),
		{Name: "dependent", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}, Needs: []string{"build"}},
		{Name: "independent", Argv: []string{"bash", "-c", `touch "$1"`, "bash", ran}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase dependent: skipped (needs build)\n") {
		t.Fatalf("dependent of a red phase is not reported as skipped-with-cause:\n%s", stdout.String())
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("dependent of a red phase still ran (it would grade an artifact its need failed to produce)")
	}
	if !strings.Contains(stdout.String(), "phase independent: green\n") {
		t.Fatalf("phase with no path from the red one lost its summary:\n%s", stdout.String())
	}
	if _, err := os.Stat(ran); err != nil {
		t.Fatalf("phase with no path from the red one did not run: %v", err)
	}
}

func TestSchedulerPropagatesOptionalSkip(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "linter", Argv: []string{"definitely-absent-linter-for-bench-test"}, Optional: true},
		{Name: "dependent", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}, Needs: []string{"linter"}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, want a propagated skip to stay green; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase dependent: skipped (needs linter)\n") {
		t.Fatalf("dependent of a skipped optional phase is not reported as skipped-with-cause:\n%s", stdout.String())
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("dependent ran against an artifact its skipped need never produced")
	}
}

func TestRunnerNeededPhaseNotFirstInnerMode(t *testing.T) {
	t.Parallel()
	t.Run("runs first", func(t *testing.T) {
		root := t.TempDir()
		marker := filepath.Join(root, "built")
		phases := []Phase{
			// A reader listed ahead of its need: sequential execution in table
			// order would run it before the marker exists.
			{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}, Needs: []string{"build"}},
			{
				Name: "build",
				Argv: []string{"bash", "-c", `touch "$1"`, "bash", marker},
			},
		}

		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
		if rc != 0 {
			t.Fatalf("runPhases rc = %d; the reader ran before its later-declared need\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
		}
	})

	t.Run("red need skips its dependent", func(t *testing.T) {
		root := t.TempDir()
		leak := filepath.Join(root, "leak")
		phases := []Phase{
			{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}, Needs: []string{"build"}},
			{Name: "build", Argv: []string{"bash", "-c", "exit 3"}},
		}

		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
		if rc != 1 {
			t.Fatalf("runPhases rc = %d, want 1", rc)
		}
		if _, err := os.Stat(leak); err == nil {
			t.Fatalf("phase listed ahead of its failed need still ran")
		}
	})
}
