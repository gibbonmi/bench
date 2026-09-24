package shift

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// shiftBranchName extracts the branch a Loop run started, from its "▶ shift on <branch>"
// line. This is the same pattern runtime_gate_shift_proof_test.go's shiftBranchFromStart uses.
func shiftBranchName(t *testing.T, output string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^▶ shift on (\S+) `)
	m := re.FindStringSubmatch(output)
	if len(m) != 2 {
		t.Fatalf("shift start did not name branch:\n%s", output)
	}
	return m[1]
}

// shiftWorktreePath extracts the worktree path a Loop run reports, from its
// "  worktree: <path>" line.
func shiftWorktreePath(t *testing.T, output string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^  worktree: (.+)$`)
	m := re.FindStringSubmatch(output)
	if len(m) != 2 {
		t.Fatalf("shift did not report its worktree:\n%s", output)
	}
	return m[1]
}

func runGitCmd(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func runGitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// requireRetainedWorktree asserts the shift named its charged worktree as the recovery
// pointer and left the tree locked in place with its dirty work on disk.
func requireRetainedWorktree(t *testing.T, root, stdout string) {
	t.Helper()
	wt := shiftWorktreePath(t, stdout)
	if !contains(stdout, "worktree:"+wt) {
		t.Fatalf("stdout did not name the retained worktree as the recovery pointer:\n%s", stdout)
	}
	porcelain := runGitOutput(t, root, "worktree", "list", "--porcelain")
	if !contains(porcelain, "worktree "+wt) || !contains(porcelain, "locked") {
		t.Fatalf("worktree was not retained and locked:\n%s", porcelain)
	}
	if _, err := os.Stat(filepath.Join(wt, "work.txt")); err != nil {
		t.Fatalf("retained worktree lost the work it was preserving: %v", err)
	}
}

// faultFixtureCore builds a throwaway repo with a gate and chdirs the test into it. This
// is the shared repo-setup every in-process Seam B fault test in this file drives Loop
// against. extra runs after the gate is written but before the init commit. It lets a
// caller add its own tracked files, such as an agent script, to that same commit; nil
// skips it. Every caller still owns its own BENCH_AGENT: this only sets the env both
// fixtures share.
func faultFixtureCore(t *testing.T, gateScript string, extra func(root string)) (root string) {
	t.Helper()
	root = t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}
	runGit("init", "-q", "-b", "main")
	// The shift loop runs its own git commit here, so the identity has to sit in the
	// repository config. A per-command -c leaves the product's commit without an author.
	runGit("config", "user.email", "bench@local")
	runGit("config", "user.name", "bench")
	f := testrepo.NewGateFixture(t.TempDir())
	if err := f.Write(root, f.Command("bash")+" -c "+sanitize.ShellQuote(gateScript), ""); err != nil {
		t.Fatal(err)
	}
	if extra != nil {
		extra(root)
	}
	runGit("add", "-A")
	runGit("commit", "-q", "-m", "init")

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	t.Setenv("BENCH_MAX_ITERS", "1")
	return root
}

// faultFixture is faultFixtureCore plus an agent that writes work.txt, with BENCH_AGENT
// set to it. This is the shared setup every in-process Seam B fault test in this file
// that wants a mutating adapter drives Loop against.
func faultFixture(t *testing.T, gateScript string) (root string) {
	t.Helper()
	var agentPath string
	root = faultFixtureCore(t, gateScript, func(root string) {
		agentPath = filepath.Join(root, "agent")
		if err := os.WriteFile(agentPath, []byte("#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	})
	t.Setenv("BENCH_AGENT", agentPath)
	return root
}

// armFault sets the package-level fault seam for the duration of the test and restores
// it on cleanup, the same discipline internal/worktree's Fault tests use.
func armFault(t *testing.T, f fault) {
	t.Helper()
	old := shiftFault
	shiftFault = f
	t.Cleanup(func() { shiftFault = old })
}

// TestLoopStagingFaultPreservesAndSplitsEvidence covers row 14 (TDD) and pins FT79's
// rule. A fault injected at the staging step must propagate, retaining the dirty tree
// and splitting by evidence, rather than being swallowed. No partial commit lands on
// the branch.
func TestLoopStagingFaultPreservesAndSplitsEvidence(t *testing.T) {
	root := faultFixture(t, "#!/usr/bin/env bash\nexit 0\n") // gate is green; only staging fails
	armFault(t, func(step shiftStep) error {
		if step == stepStage {
			return exec.Command("false").Run()
		}
		return nil
	})

	var stdout, stderr bytes.Buffer
	code := Loop("staging fault", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Loop returned %d, want 1 (failed, zero commits): stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	branch := shiftBranchName(t, stdout.String())
	if out, err := exec.Command("git", "-C", root, "rev-list", "--count", "HEAD.."+branch).CombinedOutput(); err != nil || string(bytes.TrimSpace(out)) != "0" {
		t.Fatalf("staging-fault shift committed a partial tree: count=%q err=%v", out, err)
	}
	requireRetainedWorktree(t, root, stdout.String())
}

// TestLoopTeardownFaultReportsFailed covers row 15 (TDD): a fault at teardown must exit
// failed/1 with a detail naming the branch, and any recovery ref. It must not report
// the shift's own outcome while silently swallowing a real cleanup failure.
func TestLoopTeardownFaultReportsFailed(t *testing.T) {
	// A no-op adapter (`true`) leaves nothing dirty, so the only failure this run can hit
	// is the teardown fault itself. This isolates the row's assertion from the evidence split.
	t.Setenv("BENCH_AGENT", "true")
	faultFixtureNoAgentOverride(t, "#!/usr/bin/env bash\nexit 0\n")
	armFault(t, func(step shiftStep) error {
		if step == stepTeardown {
			return exec.Command("false").Run()
		}
		return nil
	})

	var stdout, stderr bytes.Buffer
	code := Loop("teardown fault", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Loop returned %d, want 1 (teardown failure is always failed/1): stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	branch := shiftBranchName(t, stdout.String())
	if !contains(stdout.String(), "teardown failed") || !contains(stdout.String(), branch) {
		t.Fatalf("teardown-fault result did not name the branch as safe:\n%s", stdout.String())
	}
}

// faultFixtureNoAgentOverride is faultFixtureCore without an agent file, for a test that
// wants its own no-op adapter, with BENCH_AGENT set by the caller.
func faultFixtureNoAgentOverride(t *testing.T, gateScript string) (root string) {
	t.Helper()
	return faultFixtureCore(t, gateScript, nil)
}

// TestFinishReportsUpsertFailure covers C2: a finish-time intent.Upsert failure must not
// be swallowed with `_ =`. It is forced here via the stepIntentUpsert fault, since a
// real ledger write only fails on a broken BENCH_HOME, a much harder repro. The outcome
// and exit code still resolve from the gate's real verdict; the ledger record is
// enrichment, not the oracle. A warning naming the failure lands on stderr.
func TestFinishReportsUpsertFailure(t *testing.T) {
	faultFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	armFault(t, func(step shiftStep) error {
		if step == stepIntentUpsert {
			return fmt.Errorf("injected upsert failure")
		}
		return nil
	})

	var stdout, stderr bytes.Buffer
	code := Loop("upsert fault", &stdout, &stderr)

	if code != 3 {
		t.Fatalf("Loop returned %d, want 3 (incomplete — the fault must not change the outcome): stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !contains(stderr.String(), "could not record shift outcome") || !contains(stderr.String(), "injected upsert failure") {
		t.Fatalf("stderr did not warn about the discarded Upsert failure:\n%s", stderr.String())
	}
}

// TestLoopRetainsAndLocksDirtyWorktree covers row 4's Seam B half: a preserving failure
// retains and locks the charged worktree instead of releasing it. It drops the lease,
// and names the worktree path as the recovery pointer. Nothing is written to a ref; the
// dirty tree stays where an operator can read it.
func TestLoopRetainsAndLocksDirtyWorktree(t *testing.T) {
	root := faultFixture(t, "#!/usr/bin/env bash\nexit 1\n") // red gate: forces the preserving path

	var stdout, stderr bytes.Buffer
	code := Loop("preserving failure", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Loop returned %d, want 1 (failed, zero commits): stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	requireRetainedWorktree(t, root, stdout.String())
	if refs := runGitOutput(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/"); strings.TrimSpace(refs) != "" {
		t.Fatalf("a preserving failure authored recovery refs:\n%s", refs)
	}
	leaseFound := false
	_ = filepath.WalkDir(os.Getenv("BENCH_HOME"), func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && d.Name() == "bench-lease" {
			leaseFound = true
		}
		return nil
	})
	if leaseFound {
		t.Fatal("retain-and-lock fallback left the pool lease in place")
	}
}

// helperRoleEnv selects the role of a re-executed test binary.
const helperRoleEnv = "BENCH_SHIFT_HELPER_ROLE"

// helperWindow bounds every wait on a helper process and every adapter start.
var helperWindow = bounds.TestDeadline(0)

// TestShiftHelperProcess is the re-exec target of the helper-process tests. The interrupt
// role runs one shift and exits with its code, so the checkpoint's os.Exit ends a real
// process. The recover-act role runs one recovery pass that stops between its entry
// write and its act, and exits. Without a role it does nothing.
func TestShiftHelperProcess(t *testing.T) {
	switch os.Getenv(helperRoleEnv) {
	case "interrupt":
		os.Exit(Loop("interrupted shift", io.Discard, io.Discard))
	case "recover-act":
		shiftFault = func(step shiftStep) error {
			if step == stepRecoveryAct {
				return fmt.Errorf("stop before the act")
			}
			return nil
		}
		root, _ := git.Root()
		Recover(root, io.Discard)
		os.Exit(0)
	}
}

// helper is one re-executed test binary and the channel its exit arrives on.
type helper struct {
	cmd  *exec.Cmd
	done chan error
}

// startHelper re-executes the test binary into role in root. The test's cleanup kills
// the helper's process group.
func startHelper(t *testing.T, root, role string) *helper {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestShiftHelperProcess$")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), helperRoleEnv+"="+role)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) })
	h := &helper{cmd: cmd, done: make(chan error, 1)}
	go func() { h.done <- cmd.Wait() }()
	return h
}

// wait waits for the helper's exit and returns its exit code.
func (h *helper) wait(t *testing.T) int {
	t.Helper()
	select {
	case <-h.done:
	case <-time.After(helperWindow):
		t.Fatal("the helper process never exited")
	}
	return h.cmd.ProcessState.ExitCode()
}

// sleepingAgent installs an adapter that runs script, records its own pid, marks its
// start, and then sleeps past the helper window. It returns the file the start marks,
// and the test's cleanup kills the adapter.
func sleepingAgent(t *testing.T, script string) string {
	t.Helper()
	dir := t.TempDir()
	started, pidFile := filepath.Join(dir, "started"), filepath.Join(dir, "adapter")
	// The adapter outlasts the window, so a shift that waits for it instead of stopping it reds.
	withAgent(t, fmt.Sprintf("%secho $$ > '%s'\ntouch '%s'\nexec sleep %d\n", script, pidFile, started, 2*int(helperWindow/time.Second)))
	t.Cleanup(func() {
		if raw, err := os.ReadFile(pidFile); err == nil {
			if pid, err := strconv.Atoi(strings.TrimSpace(string(raw))); err == nil {
				_ = syscall.Kill(pid, syscall.SIGKILL)
			}
		}
	})
	return started
}

// waitForFile waits until path exists.
func waitForFile(t *testing.T, path string) {
	t.Helper()
	for deadline := time.Now().Add(helperWindow); ; time.Sleep(20 * time.Millisecond) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("%s never appeared", path)
		}
	}
}

// reapedPID is the process id of a child that exited and was waited for.
func reapedPID(t *testing.T) int {
	t.Helper()
	child := exec.Command("true")
	if err := child.Run(); err != nil {
		t.Fatal(err)
	}
	return child.Process.Pid
}

// crashedShift runs a shift in a helper process whose adapter appends MEMMARK to the
// notes, writes work.txt when dirty, and sleeps. It kills the helper with SIGKILL when
// kill is set and waits for its exit, so the lease names a reaped owner. It returns the
// repository root and the open shift entry.
func crashedShift(t *testing.T, dirty, kill bool) (string, intent.Entry) {
	t.Helper()
	root := faultFixtureCore(t, greenGate, nil)
	script := "echo MEMMARK >> " + notesFile + "\n"
	if dirty {
		script += "echo work > work.txt\n"
	}
	started := sleepingAgent(t, script)
	h := startHelper(t, root, "interrupt")
	waitForFile(t, started)
	if kill {
		_ = h.cmd.Process.Kill()
		h.wait(t)
	}
	current, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range current.Entries {
		if entry.Kind == intent.KindShift && entry.Lease != "" {
			return root, entry
		}
	}
	t.Fatal("the crashed shift left no leased entry")
	return "", intent.Entry{}
}
