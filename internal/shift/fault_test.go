package shift

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// shiftBranchName extracts the branch a Loop run started, from its "▶ shift on <branch>"
// line — the same pattern runtime_gate_shift_proof_test.go's shiftBranchFromStart uses.
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

// requireRecoveryRefResolves asserts refs/bench/recovery/<branch> exists and resolves.
func requireRecoveryRefResolves(t *testing.T, root, branch string) {
	t.Helper()
	ref := "refs/bench/recovery/" + branch
	if out, err := exec.Command("git", "-C", root, "show-ref", "--verify", ref).CombinedOutput(); err != nil {
		t.Fatalf("recovery ref %s did not resolve: %v\n%s", ref, err, out)
	}
}

// faultFixtureCore builds a throwaway repo with a gate and chdirs the test into it — the
// shared repo-setup every in-process Seam B fault test in this file drives Loop against.
// extra runs after the gate is written but before the init commit, so a caller can add
// its own tracked files (an agent script) to that same commit; nil skips it. Every
// caller still owns its own BENCH_AGENT: this only sets the env both fixtures share.
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
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte(gateScript), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(`{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if extra != nil {
		extra(root)
	}
	runGit("-c", "user.email=bench@local", "-c", "user.name=bench", "add", "-A")
	runGit("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "init")

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
// set to it — the shared setup every in-process Seam B fault test in this file that
// wants a mutating adapter drives Loop against.
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
// it on cleanup — the same discipline internal/worktree's Fault tests use.
func armFault(t *testing.T, f fault) {
	t.Helper()
	old := shiftFault
	shiftFault = f
	t.Cleanup(func() { shiftFault = old })
}

// TestLoopStagingFaultPreservesAndSplitsEvidence covers row 14 (TDD): a fault injected
// at the staging step must propagate — snapshot the dirty tree and split by evidence —
// rather than being swallowed the way stageTouched ignored a real `git add` failure
// before FT79. No partial commit lands on the branch.
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
	requireRecoveryRefResolves(t, root, branch)
}

// TestLoopTeardownFaultReportsFailed covers row 15 (TDD): a fault at teardown must exit
// failed/1 with a detail naming the branch (and any recovery ref) rather than reporting
// the shift's own outcome while silently swallowing a real cleanup failure.
func TestLoopTeardownFaultReportsFailed(t *testing.T) {
	// A no-op adapter (`true`) leaves nothing dirty, so the only failure this run can hit
	// is the teardown fault itself — isolating the row's assertion from the evidence split.
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
// wants its own no-op adapter (BENCH_AGENT set by the caller).
func faultFixtureNoAgentOverride(t *testing.T, gateScript string) (root string) {
	t.Helper()
	return faultFixtureCore(t, gateScript, nil)
}

// TestFinishReportsUpsertFailure covers C2: a finish-time intent.Upsert failure (forced
// here via the stepIntentUpsert fault, since a real ledger write only fails on a broken
// BENCH_HOME, a much harder repro) must not be swallowed with `_ =`. The outcome and
// exit code still resolve from the gate's real verdict — the ledger record is
// enrichment, not the oracle — but a warning naming the failure lands on stderr.
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

// TestLoopRetainsAndLocksOnSnapshotFailure covers row 4's Seam B half: when the
// snapshot itself fails — here forced by pre-creating the recovery ref the shift will
// try to create, so SnapshotDirty's own fail-closed update-ref rejects it — the shift
// retains and locks the worktree instead of releasing it, drops the lease, and names
// the worktree path as the recovery pointer.
func TestLoopRetainsAndLocksOnSnapshotFailure(t *testing.T) {
	root := faultFixture(t, "#!/usr/bin/env bash\nexit 1\n") // red gate: forces the preserving path
	fixed := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	oldNow := timeNow
	timeNow = func() time.Time { return fixed }
	t.Cleanup(func() { timeNow = oldNow })
	branch := "bench/shift-" + fixed.Format("20060102-150405")
	ref := "refs/bench/recovery/" + branch

	// Pre-create the exact ref this shift will try to create, at some arbitrary commit —
	// SnapshotDirty's fail-closed update-ref (zero old-oid) then rejects the create.
	head := strings.TrimSpace(runGitOutput(t, root, "rev-parse", "HEAD"))
	runGitCmd(t, root, "update-ref", ref, head)

	var stdout, stderr bytes.Buffer
	code := Loop("snapshot conflict", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Loop returned %d, want 1 (failed, zero commits): stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	wt := shiftWorktreePath(t, stdout.String())
	if !contains(stdout.String(), "worktree:"+wt) {
		t.Fatalf("stdout did not name the retained worktree as the recovery pointer:\n%s", stdout.String())
	}
	porcelain := runGitOutput(t, root, "worktree", "list", "--porcelain")
	if !contains(porcelain, "worktree "+wt) || !contains(porcelain, "locked") {
		t.Fatalf("worktree was not retained and locked after a snapshot failure:\n%s", porcelain)
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
