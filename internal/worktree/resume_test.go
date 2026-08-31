package worktree

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestResumeCleanSurfacesMalformedWorktreeAdmin(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	journeyFIFOWorktreeAdmin(t, root, "resume")
	var stdout, stderr bytes.Buffer
	done := make(chan int, 1)
	go func() { done <- ResumeCleanCommand(root, Home(), nil, &stdout, &stderr) }()
	select {
	case code := <-done:
		out := stdout.String() + stderr.String()
		requireTest(t, code == 1, "ResumeCleanCommand code=%d output=%q", code, out)
		requireTest(t, strings.Contains(out, "worktrees/resume/gitdir") && strings.Contains(out, "fifo") && strings.Contains(out, "inspect and remove it") && !strings.Contains(out, "git worktree list failed"), "resume output = %q", out)
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("resume cleanup blocked on malformed worktree admin entry")
	}
}

func TestResumeCleanRemovesOnlyVerifiedOwnedAssignment(t *testing.T) {
	home := t.TempDir()
	bindEnv(t, "BENCH_HOME", home)
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	clean := filepath.Join(filepath.Dir(root), "auto clean")
	dirty := filepath.Join(filepath.Dir(root), "auto dirty")
	locked := filepath.Join(filepath.Dir(root), "auto locked")
	pool := filepath.Join(Pool(root), "leased")
	mustMkdirAll(t, filepath.Dir(pool), 0o700)
	for _, path := range []string{clean, dirty, locked, pool} {
		gitRun(t, root, "worktree", "add", "-q", "--detach", path, "HEAD")
	}
	mustWrite(t, filepath.Join(dirty, "dirty.txt"), []byte("recover me\n"), 0o644)
	gitRun(t, root, "worktree", "lock", locked)
	lease, err := LeaseFile(pool)
	mustNoError(t, err)
	mustWrite(t, lease, []byte("123 2026-07-11T00:00:00Z\n"), 0o600)
	// Three unclaimed orphans cover what the prune path may and may not delete.
	// One is landed by ancestry, one has a distinct commit but a tree the
	// default branch holds, and one carries content that is nowhere else.
	gitRun(t, root, "branch", "worktree-agent-landed")
	empty := gitOutput(t, root, "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", "empty scratch")
	gitRun(t, root, "update-ref", "refs/heads/worktree-agent-empty", empty)
	gitRun(t, root, "switch", "-q", "-c", "worktree-agent-unique")
	mustWrite(t, filepath.Join(root, "unique.txt"), []byte("unique\n"), 0o644)
	gitRun(t, root, "add", "unique.txt")
	gitRun(t, root, "commit", "-qm", "unique work")
	gitRun(t, root, "switch", "-q", "main")
	created := time.Unix(1, 0).UTC()
	mustNoError(t, intent.Upsert(root, intent.Entry{Key: "auto-cleaned", Kind: intent.KindWorktree, CreatedAt: created, Worktree: clean}))
	mustNoError(t, intent.Upsert(root, intent.Entry{Key: "unrelated", Kind: intent.KindShift, CreatedAt: created}))
	owned := mustCreate(t, root, home, "resume-owned", "owned cleanup")
	markPending(t, root, owned.Assignment)
	before, _ := os.ReadFile(filepath.Join(dirty, "dirty.txt"))
	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(root, home, nil, &stdout, &stderr)
	requireTest(t, code == 0, "ResumeCleanCommand exit=%d\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	requireTest(t, stdout.String() == expectedResumeSummary(1, 0, "; retained foreign=2 live-lease=1 unexpected-lock=1", 2, 0, 0, 0), "resume report = %q", stdout.String())
	_, err = os.Stat(owned.Path)
	requireTest(t, os.IsNotExist(err), "verified owned worktree remains: %v", err)
	for _, path := range []string{clean, dirty, locked, pool} {
		_, err := os.Stat(path)
		requireTest(t, err == nil, "kept worktree %q: %v", path, err)
	}
	after, _ := os.ReadFile(filepath.Join(dirty, "dirty.txt"))
	requireTest(t, bytes.Equal(before, after), "resume cleanup changed dirty bytes")
	requireTest(t, !git.OK("-C", root, "show-ref", "--verify", "--quiet", "refs/heads/worktree-agent-landed"), "resume cleanup retained a name-only ancestry-landed orphan")
	requireTest(t, !git.OK("-C", root, "show-ref", "--verify", "--quiet", "refs/heads/worktree-agent-empty"), "resume cleanup retained an orphan holding the default branch's own tree")
	requireTest(t, git.OK("-C", root, "show-ref", "--verify", "--quiet", "refs/heads/worktree-agent-unique"), "resume cleanup deleted unique orphan")
}

func TestResumeCleanKeepsIgnoredOnlyOutOfPoolWorktree(t *testing.T) {
	bindEnv(t, "BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ignore")
	owned := mustCreate(t, root, Home(), "ignored-only", "ignored residual")
	candidate := owned.Path
	markPending(t, root, owned.Assignment)
	ignored := filepath.Join(candidate, "ignored.txt")
	mustWrite(t, ignored, []byte("retain me\n"), 0o644)
	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(root, Home(), nil, &stdout, &stderr)
	requireTest(t, code == 0, "ResumeCleanCommand exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	_, err := os.Stat(ignored)
	requireTest(t, err == nil, "ignored-only WIP was not retained: %v", err)
	requireTest(t, strings.Contains(stdout.String(), "retained ignored=1"), "ignored-only WIP not classified as retained dirty state: %q", stdout.String())
}

func TestConcurrentCleanupRecordsOneTransaction(t *testing.T) {
	t.Parallel()
	for _, automatic := range []bool{false, true} {
		t.Run(fmt.Sprintf("automatic=%t", automatic), func(t *testing.T) {
			t.Parallel()
			root, creation, _ := newOwnedAssignment(t, fmt.Sprintf("concurrent-%t", automatic))
			// The two planners differ on dirt: only the explicit one preserves it. Each
			// side of this race is driven with the dirtiest tree its planner still
			// removes. That is what makes the recovery-ref count below meaningful for
			// one and zero for the other.
			wantAction := ActionRecoverRemove
			wantRefs := 1
			if automatic {
				markPending(t, root, creation.Assignment)
				wantAction, wantRefs = ActionRemove, 0
			} else {
				mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("recover once\n"), 0o644)
			}
			plan, err := PlanExplicit(root, creation.Path)
			if automatic {
				plan, err = PlanAutomatic(root, creation.Path)
			}
			requireTest(t, err == nil && plan.Action == wantAction, "plan = %#v, %v; want %q", plan, err, wantAction)
			attempted, locked, proceed := make(chan string, 8), make(chan struct{}), make(chan struct{})
			j := defaultJoins()
			j.cleanupLockAttempt = func(target string) { attempted <- target }
			var once sync.Once
			j.cleanupBoundary = func(step LifecycleStep) error {
				if step == StepApplyLocked {
					once.Do(func() { close(locked); <-proceed })
				}
				return nil
			}
			type outcome struct {
				plan CleanupPlan
				err  error
			}
			results := make(chan outcome, 2)
			apply := func() {
				got, applyErr := applyExplicitWith(j, root, creation.Path, plan.Fingerprint, CleanupOptions{})
				if automatic {
					got, applyErr = applyAutomaticWithTerminal(j, root, creation.Path, nil, nil)
				}
				results <- outcome{got, applyErr}
			}
			go apply()
			<-attempted
			<-locked
			go apply()
			<-attempted
			close(proceed)
			for range 2 {
				got := <-results
				requireTest(t, got.err == nil && got.plan.Action == ActionRemoved && got.plan.Fingerprint == plan.Fingerprint,
					"concurrent apply = %#v, %v", got.plan, got.err)
			}
			refs := strings.Fields(gitOutput(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/"))
			ledger, err := intent.Read(root)
			requireTest(t, len(refs) == wantRefs && err == nil && len(ledger.CleanupReceipts) == 1 && ledger.CleanupReceipts[0].State == intent.ReceiptComplete,
				"transaction refs=%#v receipts=%#v error=%v", refs, ledger.CleanupReceipts, err)
			if !automatic {
				mustNoError(t, intent.DeleteAssignment(root, creation.Assignment.ID))
				replay, err := applyExplicitWith(j, root, creation.Path, plan.Fingerprint, CleanupOptions{})
				requireTest(t, err == nil && replay.Action == ActionRemoved, "compacted replay = %#v, %v", replay, err)
			}
		})
	}
	t.Run("release", func(t *testing.T) {
		orderRoot, ordered, home := newOwnedAssignment(t, "release-receipt-order")
		tip := gitOutput(t, orderRoot, "rev-parse", ordered.Assignment.Branch)
		stop := errors.New("stop after terminal release receipt")
		faulted := defaultJoins()
		faulted.cleanupBoundary = func(step LifecycleStep) error {
			if step == StepTerminalReceipt {
				return stop
			}
			return nil
		}
		code := releaseCommandWith(faulted, orderRoot, home, []string{"--request", "landed-release-receipt-order", ordered.Path}, io.Discard, io.Discard)
		requireTest(t, code != 0, "terminal receipt fault unexpectedly succeeded")
		repo, _, _ := cleanupIdentity(orderRoot, ordered.Path)
		receipt, found, err := intent.CleanupReceiptFor(orderRoot, repo, releaseOperation, ordered.Path, intent.RequestDigest("landed-release-receipt-order"))
		requireTest(t, err == nil && found && receipt.Branch == ordered.Assignment.Branch && receipt.BranchOID == tip, "terminal release receipt = %#v, found=%t error=%v", receipt, found, err)
		_, err = assignmentByID(orderRoot, ordered.Assignment.ID)
		requireTest(t, err == nil, "assignment compacted before terminal receipt checkpoint: %v", err)
		code = ReleaseCommand(orderRoot, home, []string{"--request", "landed-release-receipt-order", ordered.Path}, io.Discard, io.Discard)
		requireTest(t, code == 0, "terminal receipt replay exit=%d", code)
		receipt, found, err = intent.CleanupReceiptFor(orderRoot, repo, releaseOperation, ordered.Path, intent.RequestDigest("landed-release-receipt-order"))
		requireTest(t, err == nil && found && receipt.Branch == ordered.Assignment.Branch && receipt.BranchOID == tip, "replayed terminal release receipt = %#v, found=%t error=%v", receipt, found, err)
		root, creation, home := newOwnedAssignment(t, "concurrent-release")
		attempted, locked, proceed := make(chan string, 8), make(chan struct{}), make(chan struct{})
		j := defaultJoins()
		j.cleanupLockAttempt = func(target string) { attempted <- target }
		var once sync.Once
		j.cleanupBoundary = func(step LifecycleStep) error {
			if step == StepApplyLocked {
				once.Do(func() { close(locked); <-proceed })
			}
			return nil
		}
		type outcome struct {
			code           int
			stdout, stderr string
		}
		results := make(chan outcome, 2)
		apply := func() {
			var stdout, stderr bytes.Buffer
			code := releaseCommandWith(j, root, home, []string{"--request", "landed-concurrent-release", creation.Path}, &stdout, &stderr)
			results <- outcome{code, stdout.String(), stderr.String()}
		}
		go apply()
		<-attempted
		<-locked
		go apply()
		<-attempted
		close(proceed)
		first, second := <-results, <-results
		requireTest(t, first.code == 0 && second.code == 0 && first.stderr == "" && second.stderr == "" && first.stdout == second.stdout,
			"concurrent release = %#v / %#v", first, second)
		var replay, replayErr bytes.Buffer
		code = ReleaseCommand(root, home, []string{"--request", "landed-concurrent-release", creation.Path}, &replay, &replayErr)
		requireTest(t, code == 0 && replay.String() == first.stdout, "compacted release replay code=%d stdout=%q stderr=%q", code, replay.String(), replayErr.String())
		code = ReleaseCommand(root, home, []string{"--request", "changed", creation.Path}, io.Discard, io.Discard)
		requireTest(t, code != 0, "changed request replay was authorized")
		code = ReleaseCommand(root, home, []string{"--request", "landed-concurrent-release", root}, io.Discard, io.Discard)
		requireTest(t, code != 0, "changed path replay was authorized")
		_, err = assignmentByID(root, creation.Assignment.ID)
		requireTest(t, err != nil, "terminal release did not compact assignment")
	})
}

// TestApplyAutomaticHonorsLockScopedReplanOverPreLockCheck proves the lock-scoped
// replan is authoritative over applyAutomaticWithTerminal's pre-lock fast path.
// The pre-lock plan says remove, but state changes once the transaction lock is
// held (the deterministic StepApplyLocked seam). So the fresh replan under lock
// says retain. Execution must honor the later, authoritative verdict rather than
// the earlier one that let it through the fast path.
func TestApplyAutomaticHonorsLockScopedReplanOverPreLockCheck(t *testing.T) {
	t.Parallel()
	root, creation, _ := newPendingAssignment(t, "lock-scoped-replan")
	requirePlanAction(t, root, creation.Path, ActionRemove)
	j := defaultJoins()
	var raced bool
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step == StepApplyLocked && !raced {
			raced = true
			mustWrite(t, filepath.Join(creation.Path, "raced.txt"), []byte("raced\n"), 0o644)
		}
		return nil
	}
	plan, err := applyAutomaticWithTerminal(j, root, creation.Path, nil, nil)
	requireTest(t, err == nil && plan.Action == ActionRetain,
		"raced apply = %#v, %v; want the lock-scoped replan's retain honored", plan, err)
	requireTest(t, raced, "boundary seam never fired; race was not exercised")
	_, statErr := os.Stat(creation.Path)
	requireTest(t, statErr == nil,
		"apply removed a worktree the lock-scoped replan retained: %v", statErr)
}

func TestPlanAutomaticUsesLandedInDefaultMatrix(t *testing.T) {
	t.Parallel()
	t.Run("ancestry eligible", func(t *testing.T) {
		root, creation, _ := newPendingAssignment(t, "ancestry")
		requirePlanAction(t, root, creation.Path, ActionRemove)
	})
	t.Run("complete patch equivalence eligible", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "patch")
		commitInWorktree(t, creation.Path, "patch.txt", "landed\n", "patch")
		gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "--allow-empty", "-qm", "diverge")
		gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "cherry-pick", strings.TrimPrefix(creation.Assignment.Branch, "refs/heads/"))
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRemove)
	})
	t.Run("unique patch retained", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "unique")
		commitInWorktree(t, creation.Path, "unique.txt", "unique\n", "unique")
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRetain)
	})
	t.Run("evil merge retained", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "evil")
		commitInWorktree(t, creation.Path, "feature.txt", "feature\n", "feature")
		mustWrite(t, filepath.Join(root, "main.txt"), []byte("mainline\n"), 0o644)
		gitRun(t, root, "add", "main.txt")
		gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "mainline")
		gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "merge", "--no-commit", "--no-ff", "main")
		mustWrite(t, filepath.Join(creation.Path, "merge-only.txt"), []byte("merge only\n"), 0o644)
		gitRun(t, creation.Path, "add", "merge-only.txt")
		gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "evil merge")
		gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "cherry-pick", strings.TrimPrefix(creation.Assignment.Branch, "refs/heads/")+"~1")
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRetain)
	})
	t.Run("squash landing eligible", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "squash")
		commitInWorktree(t, creation.Path, "one.txt", "one\n", "one")
		commitInWorktree(t, creation.Path, "two.txt", "two\n", "two")
		short := strings.TrimPrefix(creation.Assignment.Branch, "refs/heads/")
		gitRun(t, root, "cherry-pick", "--no-commit", creation.Assignment.Start+".."+short)
		gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "squashed")
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRemove)
	})
	t.Run("missing default retained", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "missing-default")
		gitRun(t, root, "branch", "-m", "main", "trunk")
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRetain)
	})
	t.Run("landedness query failure retained", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "query-failure")
		commitInWorktree(t, creation.Path, "query.txt", "query\n", "query")
		branchOID := gitOutput(t, creation.Path, "rev-parse", "HEAD")
		replace := filepath.Join(root, ".git", "refs", "replace", branchOID)
		mustMkdirAll(t, filepath.Dir(replace), 0o700)
		mustWrite(t, replace, []byte(strings.Repeat("f", 40)+"\n"), 0o600)
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRetain)
	})
}

func TestPlanAutomaticRetainsDirtyNestedState(t *testing.T) {
	t.Parallel()
	t.Run("dirty nested repository", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "dirty-nested")
		nested := filepath.Join(creation.Path, "nested")
		mustMkdirAll(t, nested, 0o755)
		gitRun(t, nested, "init", "-q", "-b", "main")
		mustWrite(t, filepath.Join(nested, "nested.txt"), []byte("base\n"), 0o644)
		gitRun(t, nested, "add", "nested.txt")
		gitRun(t, nested, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "nested")
		mustWrite(t, filepath.Join(nested, "nested.txt"), []byte("dirty\n"), 0o644)
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRetain)
	})
	t.Run("dirty submodule", func(t *testing.T) {
		root, creation, _ := newOwnedSubmoduleAssignment(t, "dirty-submodule")
		mustWrite(t, filepath.Join(creation.Path, "sub", "sub.txt"), []byte("dirty\n"), 0o644)
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRetain)
	})
	t.Run("clean gitlink remains classifiable", func(t *testing.T) {
		root, creation, _ := newOwnedSubmoduleAssignment(t, "clean-submodule")
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRemove)
	})
	// Ordinary dirt is classified, not undecided: the automatic planner names it
	// as the reason it retains. This separates a checkout holding uncommitted
	// work from one whose state it could not read.
	t.Run("ordinary parent dirt remains classifiable", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "ordinary-dirt")
		mustWrite(t, filepath.Join(creation.Path, "ordinary.txt"), []byte("ordinary\n"), 0o644)
		markPending(t, root, creation.Assignment)
		requirePlanAction(t, root, creation.Path, ActionRetain)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.ReasonCode == ReasonDirty, "dirty plan reason = %q, want %q", plan.ReasonCode, ReasonDirty)
	})
}

func TestExplicitApplyBindsRecoveryActionsAndDiscardFlag(t *testing.T) {
	t.Parallel()
	t.Run("recovery action", func(t *testing.T) {
		root := newWorktreeRepo(t)
		gitRun(t, root, "branch", "-M", "main")
		target := filepath.Join(filepath.Dir(root), "recovery action drift")
		gitRun(t, root, "worktree", "add", "-q", "--detach", target, "HEAD")
		plan, err := PlanExplicit(root, target)
		requireTest(t, err == nil && plan.Recovery != "none", "detached plan = %#v, %v", plan, err)
		head := gitOutput(t, target, "rev-parse", "HEAD")
		gitRun(t, root, "update-ref", plan.Recovery, head)
		before := gitOutput(t, root, "worktree", "list", "--porcelain")
		current, err := ApplyExplicit(root, target, plan.Fingerprint)
		requireTest(t, errors.Is(err, errStaleFingerprint) && current.Recovery != plan.Recovery, "recovery drift apply = %#v, %v", current, err)
		requireTest(t, gitOutput(t, root, "worktree", "list", "--porcelain") == before, "recovery drift removed the target")
	})
	t.Run("discard flag", func(t *testing.T) {
		root := newWorktreeRepo(t)
		gitRun(t, root, "branch", "-M", "main")
		mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("ignored.txt\n"), 0o644)
		target := filepath.Join(filepath.Dir(root), "discard flag drift")
		gitRun(t, root, "worktree", "add", "-q", "-b", "discard-flag", target, "HEAD")
		ignored := filepath.Join(target, "ignored.txt")
		mustWrite(t, ignored, []byte("secret\n"), 0o600)
		plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{DiscardIgnored: true})
		requireTest(t, err == nil && plan.Action == ActionDiscardRemove, "discard plan = %#v, %v", plan, err)
		current, err := ApplyExplicit(root, target, plan.Fingerprint)
		requireTest(t, errors.Is(err, errStaleFingerprint) && current.Fingerprint != plan.Fingerprint, "discard flag drift apply = %#v, %v", current, err)
		body, err := os.ReadFile(ignored)
		requireTest(t, err == nil && string(body) == "secret\n", "discard flag drift changed ignored file: %q, %v", body, err)
	})
}

// TestResumeSupersedesPlannedReceiptOverAbsentTarget proves a cleanup that stopped in
// the window between its receipt write and its first checkpoint does not wedge the
// retry once the target is gone. Phase planned means no side effect ran, so the retry
// re-plans the current tree instead of finishing a cleanup that never started. A
// preserved phase and a present target keep the refusal they already had.
func TestResumeSupersedesPlannedReceiptOverAbsentTarget(t *testing.T) {
	t.Parallel()
	t.Run("absent target", func(t *testing.T) {
		t.Parallel()
		root, creation, _ := newPendingAssignment(t, "planned-absent")
		fingerprint := crashedInCleanupReceiptWindow(t, root, creation.Path, false)
		mustNoError(t, os.RemoveAll(creation.Path))
		plan, err := resumeCleanupTransaction(root, creation.Path, fingerprint)
		requireTest(t, err == nil, "resume over absent target = %#v, %v", plan, err)
		receipt := cleanupReceiptAt(t, root, creation.Path, fingerprint)
		requireTest(t, receipt.State == intent.ReceiptComplete && receipt.Phase == intent.ReceiptPhaseTerminal,
			"resumed receipt = state %q phase %q", receipt.State, receipt.Phase)
		replay, err := resumeCleanupTransaction(root, creation.Path, fingerprint)
		requireTest(t, err == nil && replay.Fingerprint == fingerprint,
			"replayed resume = %#v, %v", replay, err)
	})
	t.Run("absent target of a preserving plan", func(t *testing.T) {
		t.Parallel()
		root, creation, _ := newOwnedAssignment(t, "planned-absent-preserving")
		mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("preserve me\n"), 0o644)
		fingerprint := crashedInCleanupReceiptWindow(t, root, creation.Path, true)
		mustNoError(t, os.RemoveAll(creation.Path))
		plan, err := resumeCleanupTransaction(root, creation.Path, fingerprint)
		requireTest(t, err == nil, "resume over absent preserving target = %#v, %v", plan, err)
		receipt := cleanupReceiptAt(t, root, creation.Path, fingerprint)
		requireTest(t, receipt.State == intent.ReceiptComplete && receipt.Phase == intent.ReceiptPhaseTerminal,
			"resumed preserving receipt = state %q phase %q", receipt.State, receipt.Phase)
		ledger, err := intent.Read(root)
		mustNoError(t, err)
		for _, current := range ledger.CleanupReceipts {
			requireTest(t, current.State != intent.ReceiptInFlight, "superseded receipt survived: %#v", current)
		}
	})
	t.Run("absent target at the preserved phase", func(t *testing.T) {
		t.Parallel()
		root, creation, _ := newOwnedAssignment(t, "preserved-absent")
		mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("preserve me\n"), 0o644)
		plan, err := PlanExplicit(root, creation.Path)
		requireTest(t, err == nil && plan.preserves(), "preserving plan = %#v, %v", plan, err)
		stop := errors.New("stop after the preserved checkpoint")
		faulted := defaultJoins()
		faulted.cleanupBoundary = func(step LifecycleStep) error {
			if step == StepRecoveryRef {
				return stop
			}
			return nil
		}
		_, err = applyExplicitWith(faulted, root, creation.Path, plan.Fingerprint, CleanupOptions{})
		requireTest(t, errors.Is(err, stop), "preserved window fault = %v", err)
		receipt := cleanupReceiptAt(t, root, creation.Path, plan.Fingerprint)
		requireTest(t, receipt.State == intent.ReceiptInFlight && receipt.Phase == intent.ReceiptPhasePreserved,
			"preserved receipt = state %q phase %q", receipt.State, receipt.Phase)
		mustNoError(t, os.RemoveAll(creation.Path))
		got, err := resumeCleanupTransaction(root, creation.Path, plan.Fingerprint)
		requireTest(t, errors.Is(err, errStaleFingerprint) && got.Fingerprint == receipt.Fingerprint &&
			got.Action == CleanupAction(receipt.Action) && got.Recovery == receipt.Recovery,
			"preserved resume = %#v, %v; want the receipt's own plan refused", got, err)
	})
	// The explicit fingerprint is content-derived, so drift in the surviving tree is
	// what makes the retry ask a different question than the receipt answers.
	t.Run("present target", func(t *testing.T) {
		t.Parallel()
		root, creation, _ := newOwnedAssignment(t, "planned-present")
		mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("preserve me\n"), 0o644)
		fingerprint := crashedInCleanupReceiptWindow(t, root, creation.Path, true)
		mustWrite(t, filepath.Join(creation.Path, "drift.txt"), []byte("drift\n"), 0o644)
		got, err := ApplyExplicit(root, creation.Path, fingerprint)
		requireTest(t, errors.Is(err, errStaleFingerprint), "present-target resume = %#v, %v", got, err)
	})
}

// crashedInCleanupReceiptWindow drives one cleanup to the deterministic StepReceipt seam,
// which is the window between the in-flight receipt write and the first checkpoint. It
// returns the fingerprint that receipt carries. The explicit planner is the preserving
// one; the automatic planner removes a pending assignment without preserving.
func crashedInCleanupReceiptWindow(t *testing.T, root, target string, preserving bool) string {
	t.Helper()
	plan, err := PlanAutomatic(root, target)
	if preserving {
		plan, err = PlanExplicit(root, target)
	}
	requireTest(t, err == nil && plan.preserves() == preserving, "crash-window plan = %#v, %v", plan, err)
	stop := errors.New("stop in the cleanup receipt window")
	faulted := defaultJoins()
	faulted.cleanupBoundary = func(step LifecycleStep) error {
		if step == StepReceipt {
			return stop
		}
		return nil
	}
	if preserving {
		_, err = applyExplicitWith(faulted, root, target, plan.Fingerprint, CleanupOptions{})
	} else {
		_, err = applyAutomaticWithTerminal(faulted, root, target, nil, nil)
	}
	requireTest(t, errors.Is(err, stop), "cleanup receipt window fault = %v", err)
	receipt := cleanupReceiptAt(t, root, target, plan.Fingerprint)
	requireTest(t, receipt.State == intent.ReceiptInFlight && receipt.Phase == intent.ReceiptPhasePlanned,
		"crash-window receipt = state %q phase %q", receipt.State, receipt.Phase)
	return plan.Fingerprint
}

func cleanupReceiptAt(t *testing.T, root, target, fingerprint string) intent.CleanupReceipt {
	t.Helper()
	repo, canonical, err := cleanupIdentity(root, target)
	mustNoError(t, err)
	receipt, found, err := intent.CleanupReceiptFor(root, repo, cleanupOperation, canonical, fingerprint)
	requireTest(t, err == nil && found, "cleanup receipt for %q found=%t error=%v", canonical, found, err)
	return receipt
}

// resumeCleanupTransaction is the retry an abandon runs when it finds an in-flight
// cleanup receipt over an absent target: the receipt's own fingerprint carried back into
// the transaction behind the automatic planner.
func resumeCleanupTransaction(root, target, fingerprint string) (CleanupPlan, error) {
	j := defaultJoins()
	planner := func(path string) (CleanupPlan, error) { return planAutomaticAt(j, root, path, currentTime()) }
	return applyCleanupTransaction(j, root, target, fingerprint, planner, nil, nil)
}

func newOwnedSubmoduleAssignment(t *testing.T, request string) (string, Creation, string) {
	t.Helper()
	root := newWorktreeRepo(t)
	source := journeyRepoOnBranch(t, "main")
	mustWrite(t, filepath.Join(source, "sub.txt"), []byte("clean\n"), 0o644)
	gitRun(t, source, "add", "sub.txt")
	gitRun(t, source, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "submodule")
	gitRun(t, root, "-c", "protocol.file.allow=always", "submodule", "add", "-q", source, "sub")
	gitRun(t, root, "add", ".gitmodules", "sub")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "add submodule")
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "nested-"+request, "nested state")
	gitRun(t, creation.Path, "-c", "protocol.file.allow=always", "submodule", "update", "--init", "-q")
	return root, creation, home
}

// newOwnedAssignment returns the repository, its one owned registration, and the
// private home that registration lives under. The caller passes that home to every
// verb, so the fixture binds no process environment.
func newOwnedAssignment(t *testing.T, request string) (string, Creation, string) {
	t.Helper()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "landed-"+request, "landedness")
	return root, creation, home
}
func newPendingAssignment(t *testing.T, request string) (string, Creation, string) {
	t.Helper()
	root, creation, home := newOwnedAssignment(t, request)
	markPending(t, root, creation.Assignment)
	return root, creation, home
}
func markPending(t *testing.T, root string, assignment intent.Assignment) {
	t.Helper()
	assignment.State = intent.StateCleanupPending
	mustNoError(t, intent.PutAssignment(root, assignment))
}
func commitInWorktree(t *testing.T, path, name, body, message string) {
	t.Helper()
	mustWrite(t, filepath.Join(path, name), []byte(body), 0o644)
	gitRun(t, path, "add", name)
	gitRun(t, path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", message)
}
func requirePlanAction(t *testing.T, root, path string, want CleanupAction) {
	t.Helper()
	plan, err := PlanAutomatic(root, path)
	mustNoError(t, err)
	requireTest(t, plan.Action == want, "PlanAutomatic = %#v, want action %q", plan, want)
}
func assignmentString(a intent.Assignment) string { return fmt.Sprintf("%s/%s", a.OwnerID, a.ID) }

// mustCreate creates one owned registration under the home the caller names. The
// home is explicit, so a parallel test keeps its pool private to itself.
func mustCreate(t *testing.T, root, home, request, label string) Creation {
	t.Helper()
	creation, err := createAt(defaultJoins(), root, home, request, label, nil, currentTime())
	mustNoError(t, err)
	return creation
}

func mustNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func requireTest(t *testing.T, ok bool, format string, args ...any) {
	t.Helper()
	if !ok {
		t.Fatalf(format, args...)
	}
}
