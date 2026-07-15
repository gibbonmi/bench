package gate

// The serial-phase contract: a Serial phase completes before any concurrent phase
// starts regardless of table position, its red fails the run fast, and an interrupt
// during it is a 130, never a graded red — in both outer and inner modes.

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

func r12Contention(id, action string) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := gateTestRepo(t, "#!/usr/bin/env bash\ngitdir=\"$(git rev-parse --absolute-git-dir)\"\nprintf 'run\\n' >> \"$gitdir/owner-runs\"\n[ \"$(wc -l < \"$gitdir/owner-runs\")\" -eq 1 ] || exit 0\ntouch \"$gitdir/owner-started\"\nwhile [ ! -f \"$gitdir/release-owner\" ]; do sleep .01; done\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
		_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".bench-contract-env/\n"), 0o644)
		_, _ = benchgit.Output("-C", root, "add", "-A")
		_, _ = benchgit.Output("-C", root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-m", "base")
		f := contract.NewFixtureAt(t, root, contract.IsolatedEnv(t, root))
		bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
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
			_ = os.MkdirAll(filepath.Join(root, "bin"), 0o755)
			_ = os.MkdirAll(filepath.Join(root, "dist"), 0o755)
			_ = os.WriteFile(filepath.Join(root, "bin", "bench.sh"), mustRead(t, bench), 0o755)
			_ = os.WriteFile(filepath.Join(root, "dist", "bench"), mustRead(t, filepath.Join(contract.SubjectRoot(t), "dist", "bench")), 0o755)
		}
		done := make(chan Result, 1)
		go func() { done <- Execute(context.Background(), ownerRoot, io.Discard, io.Discard) }()
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
			probe = contract.RunAtWithInput(t, f, root, map[string]string{"BENCH_SHIFT": "1"}, "{}\n", "bash", filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "stop.sh"))
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

func TestRunnerSerialPhaseCompletesBeforeConcurrentPhasesStart(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "built")
	phases := []Phase{
		{
			Name:   "build",
			Argv:   []string{"bash", "-c", `sleep 0.1; touch "$1"`, "bash", marker},
			Serial: true,
		},
		// Each reader phase goes red unless the serial phase already finished:
		// a runner that launches everything concurrently starts these ~100ms
		// before the marker exists.
		{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
		{Name: "bravo", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; a reader phase started before the serial phase finished\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase build: green") {
		t.Fatalf("serial phase missing from summaries:\n%s", stdout.String())
	}
}

func TestRunnerSerialRedFailsFast(t *testing.T) {
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want 1", rc)
	}
	out := stdout.String() + stderr.String()
	if !strings.Contains(out, "phase build: red (exit 3)") || !strings.Contains(out, "gate: red") {
		t.Fatalf("serial red not attributed:\n%s", out)
	}
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("phase after failed serial phase still ran (it would grade a stale or absent binary)")
	}
}

func TestRunnerSerialRedFailsFastInnerMode(t *testing.T) {
	root := t.TempDir()
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
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
		t.Fatalf("inner-mode phase after failed serial phase still ran")
	}
}

func TestRunnerCancelDuringSerialPhaseReturns130(t *testing.T) {
	root := t.TempDir()
	pidfile := filepath.Join(root, "sleep.pid")
	leak := filepath.Join(root, "leak")
	phases := []Phase{
		{
			Name:   "build",
			Argv:   []string{"bash", "-c", `sleep 30 & echo $! > "$1"; wait`, "bash", pidfile},
			Serial: true,
		},
		{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
	}
	ctx, cancel := context.WithCancel(context.Background())
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
		t.Fatal("runPhases did not return after cancellation during the serial phase")
	}
	waitForProcessExit(t, pid)
	// An interrupt is not a verdict: no phase behind the interrupted build may run,
	// and the run must not present itself as a graded red.
	if _, err := os.Stat(leak); err == nil {
		t.Fatalf("concurrent phase ran after the serial phase was interrupted")
	}
	if strings.Contains(stdout.String()+stderr.String(), "gate: red") {
		t.Fatalf("interrupted serial phase was graded as red:\nstdout=%q\nstderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunnerSerialPhaseNotFirstStillRunsFirst(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "built")
	phases := []Phase{
		// Readers listed ahead of the serial phase: only the runner's own
		// reordering keeps them from starting before the marker exists.
		{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
		{Name: "bravo", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
		{
			Name:   "build",
			Argv:   []string{"bash", "-c", `sleep 0.1; touch "$1"`, "bash", marker},
			Serial: true,
		},
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, phases, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d; a reader phase ran before the non-first serial phase\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase build: green") {
		t.Fatalf("serial phase missing from summaries:\n%s", stdout.String())
	}
}

func TestRunnerSerialPhaseNotFirstInnerMode(t *testing.T) {
	t.Run("runs first", func(t *testing.T) {
		root := t.TempDir()
		marker := filepath.Join(root, "built")
		phases := []Phase{
			// A reader listed ahead of the serial phase: sequential execution in
			// table order would run it before the marker exists.
			{Name: "alpha", Argv: []string{"bash", "-c", `test -f "$1"`, "bash", marker}},
			{
				Name:   "build",
				Argv:   []string{"bash", "-c", `touch "$1"`, "bash", marker},
				Serial: true,
			},
		}

		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
		if rc != 0 {
			t.Fatalf("runPhases rc = %d; the reader ran before the non-first serial phase\nstdout=%q\nstderr=%q", rc, stdout.String(), stderr.String())
		}
	})

	t.Run("red fails fast", func(t *testing.T) {
		root := t.TempDir()
		leak := filepath.Join(root, "leak")
		phases := []Phase{
			{Name: "alpha", Argv: []string{"bash", "-c", `touch "$1"`, "bash", leak}},
			{Name: "build", Argv: []string{"bash", "-c", "exit 3"}, Serial: true},
		}

		var stdout, stderr bytes.Buffer
		rc := runPhases(context.Background(), root, phases, innerMode, &stdout, &stderr)
		if rc != 1 {
			t.Fatalf("runPhases rc = %d, want 1", rc)
		}
		if _, err := os.Stat(leak); err == nil {
			t.Fatalf("phase listed ahead of the failed serial phase still ran")
		}
	})
}
