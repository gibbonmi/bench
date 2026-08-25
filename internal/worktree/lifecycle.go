package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ResumeResult struct {
	Removed        int
	Retained       map[CleanupReason]int
	Failed, Open   int
	SweptRefs      int
	Reconciled     int
	PrunedBranches int
	Orphans        []OrphanCandidate
	// ReclaimableKeys counts the pool keys `bench worktree reclaim` would target. It is
	// the same predicate that command plans with, so the ambient number and the verb's
	// target count cannot disagree. Resume reports it and removes nothing.
	ReclaimableKeys int
	// PoolUnreadable records why the count above is absent rather than zero. A resume
	// that cannot read the pool still succeeds at its own work, but it must not report a
	// zero it has no basis for.
	PoolUnreadable error
}

// OrphanCandidate names an assignment the sweep judged abandoned while its worktree is
// still on disk. Path is the argument `bench worktree clean` needs to start retiring it.
// The sweep only reports these — removal stays behind that explicit path-addressed
// command, which preserves dirty work before it removes anything.
type OrphanCandidate struct{ ID, Path string }

var ErrCleanupInterrupted = errors.New("cleanup interrupted")

const leaseTimeLayout = lifecyclepolicy.LeaseTimeLayout

const unknownLeaseReason = lifecyclepolicy.UnknownLeaseReason

var chmodPool = os.Chmod

// LeaseState is the policy package's lease liveness verdict;
// internal/worktree/lifecyclepolicy owns its values and semantics.
type LeaseState = lifecyclepolicy.LeaseState

const (
	LeaseLive    = lifecyclepolicy.LeaseLive
	LeaseDead    = lifecyclepolicy.LeaseDead
	LeaseUnknown = lifecyclepolicy.LeaseUnknown
)

// pidAlive treats kill-0 success and EPERM as alive. Only ESRCH means gone.
func pidAlive(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// leaseOwnerPID is the policy lease-content parse.
func leaseOwnerPID(content []byte) (int, bool) { return lifecyclepolicy.LeaseOwnerPID(content) }

// ProbeLease reports whether a well-formed lease's recorded owner is live.
// Every unreadable or malformed lease is unknown so lifecycle consumers fail closed.
func ProbeLease(leasePath string) LeaseState {
	info, err := os.Lstat(leasePath)
	if err != nil || !info.Mode().IsRegular() {
		return LeaseUnknown
	}
	content, err := os.ReadFile(leasePath)
	if err != nil {
		return LeaseUnknown
	}
	pid, ok := leaseOwnerPID(content)
	if !ok {
		return LeaseUnknown
	}
	if pidAlive(pid) {
		return LeaseLive
	}
	return LeaseDead
}

// reclaimable is the policy staleness decision over a lease's translated
// content, mtime, and the caller's liveness probe, judged against the
// bounds.LeaseStale window this boundary supplies.
func reclaimable(content []byte, mtime, now time.Time, alive func(int) bool) bool {
	return lifecyclepolicy.Reclaimable(content, mtime, now, alive, bounds.LeaseStale)
}

// candidateName keeps each unique mint attempt inside the pool.
func candidateName(pool string, unixSecs int64, pid, try int) string {
	return filepath.Join(pool, fmt.Sprintf("%d-%d-%d", unixSecs, pid, try))
}

// leaseLine is the bytes an owner writes into its lease: "<pid> <utc-time>\n".
// The instant is the caller's explicit boundary resolution, never an ambient read.
func leaseLine(now time.Time) []byte {
	return []byte(fmt.Sprintf("%d %s\n", os.Getpid(), now.UTC().Format(leaseTimeLayout)))
}

// tryCreate wins a lease only through an atomic O_EXCL create.
func tryCreate(leasePath string, now time.Time) bool {
	f, err := os.OpenFile(leasePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return false
	}
	_, werr := f.Write(leaseLine(now))
	cerr := f.Close()
	return werr == nil && cerr == nil
}

// claimAt atomically creates a lease or identity-checks a provably stale takeover,
// judging staleness against the caller's explicitly resolved instant.
func claimAt(leasePath string, now time.Time) bool {
	if tryCreate(leasePath, now) {
		return true
	}
	info, err := os.Stat(leasePath)
	if err != nil {
		return false // lease vanished under us (a racing reclaim); respect and rescan
	}
	content, _ := os.ReadFile(leasePath)
	if !reclaimable(content, info.ModTime(), now, pidAlive) {
		return false
	}
	claimTakeoverGap(leasePath)
	stale := leasePath + ".stale." + strconv.Itoa(os.Getpid())
	if os.Rename(leasePath, stale) != nil {
		return false // another reclaimer moved it first
	}
	claimStealGap(leasePath)
	// A competing reclaimer may have replaced the judged lease. Concede unless the
	// renamed bytes still match, and never clobber a first-writer in the vacated slot.
	moved, err := os.ReadFile(stale)
	if err != nil || !bytes.Equal(moved, content) {
		_ = os.Link(stale, leasePath) // best-effort; EEXIST means a first-writer won the slot
		os.Remove(stale)
		return false
	}
	os.Remove(stale)
	return tryCreate(leasePath, now)
}

// claimTakeoverGap drives the post-judgment, pre-rename reclaimer interleave.
var claimTakeoverGap = func(leasePath string) {}

// claimStealGap drives the post-rename, pre-identity-check first-writer interleave.
var claimStealGap = func(leasePath string) {}

// isWorktree accepts primary (.git dir) and linked (.git file) checkouts.
func isWorktree(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// isClean reports whether the worktree at dir has no dirty or untracked paths.
func isClean(dir string) bool {
	out, err := git.Output("-C", dir, "status", "--porcelain")
	return err == nil && out == ""
}

// Acquire is the boundary form of acquireAt for a caller in another package: it
// resolves the Bench home and the instant at the effect boundary.
func Acquire(root, resetRef, resetMode string) (string, error) {
	return acquireAt(root, resetRef, resetMode, Home(), currentTime())
}

// acquireAt claims a clean pool entry or mints one in three bounded attempts. It
// resets to resetRef (HEAD when empty). Soft mode tolerates an unresolved ref.
// The home and the instant are the caller's explicit boundary resolutions.
func acquireAt(root, resetRef, resetMode, home string, now time.Time) (string, error) {
	pool := poolAt(home, root)
	if err := os.MkdirAll(pool, 0o700); err != nil {
		return "", err
	}
	_ = chmodPool(pool, 0o700)
	var wt string
	entries, _ := os.ReadDir(pool) // sorted by name, matching the shell glob order
	for _, e := range entries {
		d := filepath.Join(pool, e.Name())
		if !isWorktree(d) || !isClean(d) {
			continue
		}
		lease, err := LeaseFile(d)
		if err != nil || !claimAt(lease, now) {
			continue
		}
		wt = d
		break
	}
	for try := 1; wt == "" && try <= 3; try++ {
		cand := candidateName(pool, now.Unix(), os.Getpid(), try)
		// An unresolved default makes the first attempt the HEAD one already. The HEAD
		// fallback is a genuine second attempt only when the remote ref was non-empty.
		remote := git.RemoteDefaultRef(root)
		if !worktreeAdd(root, cand, remote) && (remote == "" || !worktreeAdd(root, cand, "")) {
			break
		}
		lease, err := LeaseFile(cand)
		if err != nil {
			continue
		}
		if claimAt(lease, now) {
			wt = cand
		}
	}
	if wt == "" {
		return "", errors.New("could not lease a pool worktree")
	}
	_ = exec.Command("git", "-C", wt, "switch", "-q", "--detach").Run()
	if resetRef != "" {
		if exec.Command("git", "-C", wt, "reset", "-q", "--hard", resetRef).Run() != nil {
			if resetMode != "soft" {
				return "", fmt.Errorf("could not reset pool worktree to %s", resetRef)
			}
			_ = exec.Command("git", "-C", wt, "reset", "-q", "--hard").Run()
		}
	} else {
		_ = exec.Command("git", "-C", wt, "reset", "-q", "--hard").Run()
	}
	if err := exec.Command("git", "-C", wt, "clean", "-qfdx").Run(); err != nil {
		return "", err
	}
	return strings.TrimRight(wt, "/"), nil
}

// worktreeAdd runs `git -C root worktree add -q --detach cand [ref]`, reporting success.
// An empty ref adds at the current HEAD (the fallback when origin/<branch> is absent).
func worktreeAdd(root, cand, ref string) bool {
	args := []string{"-C", root, "worktree", "add", "-q", "--detach", cand}
	if ref != "" {
		args = append(args, ref)
	}
	return exec.Command("git", args...).Run() == nil
}

// Release restores cleanliness before unleasing, and leaves a lease owned by another
// live process untouched. Once unleased, a concurrent Acquire owns the checkout.
func Release(wt string) {
	if wt == "" {
		return
	}
	lease, err := LeaseFile(wt)
	if err != nil {
		return
	}
	content, _ := os.ReadFile(lease)
	if field := strings.Fields(string(content)); len(field) > 0 {
		if pid, err := strconv.Atoi(field[0]); err == nil && pid != os.Getpid() && pidAlive(pid) {
			return
		}
	}
	restoreClean(wt)
	os.Remove(lease)
}

// restoreClean is the release-ordering test seam between cleanup and unlease.
var restoreClean = func(wt string) {
	_ = exec.Command("git", "-C", wt, "switch", "-q", "--detach").Run()
	_ = exec.Command("git", "-C", wt, "reset", "-q", "--hard").Run()
	_ = exec.Command("git", "-C", wt, "clean", "-qfdx").Run()
}

// cleanupTransactionBoundary is the deterministic transaction fault seam.
var cleanupTransactionBoundary Fault
var cleanupLockAttempt = func(string) {}

func receiptFromRelease(repo, request string, assignment intent.Assignment, action, branch, branchOID string) intent.CleanupReceipt {
	receipt := intent.CleanupReceipt{Schema: intent.CleanupReceiptSchema, Repo: repo, Operation: releaseOperation, Target: assignment.Worktree, Fingerprint: request, State: intent.ReceiptComplete, Phase: intent.ReceiptPhaseTerminal, Action: action, Tracked: assignment.ID, Recovery: "none", Detail: string(assignment.State), Owned: true}
	if branch != "" && branchOID != "" {
		receipt.Branch, receipt.BranchOID = branch, branchOID
	}
	return receipt
}
func ensureRecoveryRef(root string, assignment intent.Assignment, recovery intent.Recovery) error {
	if current, err := git.Output("-C", root, "show-ref", "--verify", "--hash", recovery.Ref); err == nil {
		if current != recovery.Root {
			return errors.New("existing recovery ref conflicts with recorded metadata")
		}
		return verifyRecovery(root, assignment, recovery)
	}
	if !recoveryEnvelopeValid(root, recovery) {
		return errors.New("recorded recovery envelope is invalid")
	}
	zero := strings.Repeat("0", len(recovery.Root))
	if out, err := exec.Command("git", "-C", root, "update-ref", recovery.Ref, recovery.Root, zero).CombinedOutput(); err != nil {
		return fmt.Errorf("create exact recovery ref: %s", strings.TrimSpace(string(out)))
	}
	return verifyRecovery(root, assignment, recovery)
}
func cleanupLockPath(repo, target string) string {
	return filepath.Join(repo, "bench-cleanup-"+fingerprintParts([]byte(target))+".lock")
}
func lockCleanupFile(file *os.File, target string) (func(), error) {
	cleanupLockAttempt(target)
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		file.Close()
		return nil, fmt.Errorf("lock cleanup registration: %w", err)
	}
	return func() { _ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN); _ = file.Close() }, nil
}
func lockCleanupRegistration(repo, target string) (func(), error) {
	admin, err := git.Output("-C", target, "rev-parse", "--path-format=absolute", "--git-dir")
	var file *os.File
	if err == nil {
		file, err = os.Open(admin)
	} else {
		file, err = os.OpenFile(cleanupLockPath(repo, target), os.O_RDWR, 0o600)
		if errors.Is(err, os.ErrNotExist) {
			file, err = os.Open(repo)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("open cleanup transaction lock: %w", err)
	}
	return lockCleanupFile(file, target)
}
func lockCleanupPersistence(repo, target string) (func(), error) {
	file, err := os.OpenFile(cleanupLockPath(repo, target), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	return lockCleanupFile(file, target)
}

// releaseLeftover completes a release-leftover plan: the registration and the ledger entry
// go, the bytes at the leftover path stay. It never reaches the removal steps below.
// `git worktree remove` deletes the tree it is pointed at, which is the one thing this
// plan exists to avoid.
func releaseLeftover(root string, plan CleanupPlan, checkpoint func(string) error, fault Fault) (CleanupPlan, error) {
	assignment := *plan.assignment
	if len(assignment.Recovery) > 0 {
		if len(assignment.Recovery) != 1 {
			return plan, errors.New("existing recovery metadata is ambiguous")
		}
		if err := ensureRecoveryRef(root, assignment, assignment.Recovery[0]); err != nil {
			return plan, err
		}
		if err := checkpoint(intent.ReceiptPhasePreserved); err != nil {
			return plan, err
		}
		if err := hit(fault, StepRecoveryRef); err != nil {
			return plan, err
		}
	}
	if err := checkpoint(intent.ReceiptPhaseRemoving); err != nil {
		return plan, err
	}
	if err := releaseRegistration(root, plan.leftover); err != nil {
		return plan, err
	}
	if err := checkpoint(intent.ReceiptPhaseRemoved); err != nil {
		return plan, err
	}
	if err := hit(fault, StepRemoval); err != nil {
		return plan, err
	}
	assignment.State = intent.StateComplete
	if len(assignment.Recovery) > 0 {
		assignment.State = intent.StateRecovered
	}
	if err := intent.PutAssignment(root, assignment); err != nil {
		return plan, err
	}
	if err := checkpoint(intent.ReceiptPhaseTerminal); err != nil {
		return plan, err
	}
	plan.Action, plan.Reason, plan.ReasonCode = ActionRemoved, "", ""
	return plan, nil
}

// releaseRegistration deletes the private administration directory git keeps for exactly
// one registered worktree, which is what stops git registering it. This is the scoped form
// of `git worktree prune`: prune decides for every prunable registration at once. An
// abandon answers for one target, so a stranger's stale entry is never swept along with
// it. Only a target with no git metadata entry reaches here, so the administration
// directory it names is already dangling.
func releaseRegistration(root, target string) error {
	common, err := git.CommonDir(root)
	if err != nil {
		return fmt.Errorf("resolve common Git directory: %w", err)
	}
	pool := filepath.Join(filepath.Clean(common), "worktrees")
	entries, err := os.ReadDir(pool)
	if err != nil {
		return fmt.Errorf("read private worktree administration pool: %w", err)
	}
	admin := ""
	for _, entry := range entries {
		candidate := filepath.Join(pool, entry.Name())
		record := filepath.Join(candidate, "gitdir")
		// The whole pool is swept to find one entry, so an unrelated entry's control record
		// gets read too. Anything but a regular file is skipped rather than opened: a FIFO
		// planted here has no writer, and reading it would block the abandon forever.
		if info, statErr := os.Lstat(record); statErr != nil || !info.Mode().IsRegular() {
			continue
		}
		gitdir, readErr := os.ReadFile(record)
		if readErr != nil {
			continue
		}
		registered, pathErr := canonicalPath(filepath.Dir(strings.TrimSpace(string(gitdir))))
		if pathErr != nil || registered != target {
			continue
		}
		if admin != "" {
			return errors.New("target has ambiguous private administration directories")
		}
		admin = candidate
	}
	if admin == "" {
		return errors.New("target has no private administration directory to release")
	}
	return os.RemoveAll(admin)
}

func executeCleanup(root string, plan CleanupPlan, checkpoint func(string) error, fault Fault) (CleanupPlan, error) {
	if plan.Action == actionReleaseLeftover {
		return releaseLeftover(root, plan, checkpoint, fault)
	}
	var recovered *intent.Assignment
	if plan.preserves() {
		assignment, err := recoveryAssignmentForPlan(root, plan)
		if err != nil {
			return plan, err
		}
		if plan.owned && assignment.State == intent.StateActive {
			assignment.State = intent.StateCleanupPending
			if err := intent.PutAssignment(root, assignment); err != nil {
				return plan, err
			}
			if err := checkpoint(intent.ReceiptPhasePlanned); err != nil {
				return plan, err
			}
		}
		if plan.Tracked == "clean" && plan.registration.Detached {
			if err := anchorDetached(root, plan); err != nil {
				return plan, err
			}
			head, err := git.Output("-C", plan.Target, "rev-parse", "HEAD")
			if err != nil {
				return plan, err
			}
			assignment.Recovery = []intent.Recovery{{Ref: plan.Recovery, Root: head, Payloads: []string{head}}}
			if err := intent.PutAssignment(root, assignment); err != nil {
				return plan, err
			}
		} else if assignment, err = recoverAssignmentWithFault(root, assignment, fault); err != nil {
			return plan, err
		}
		recovered = &assignment
		if err := checkpoint(intent.ReceiptPhasePreserved); err != nil {
			return plan, err
		}
		if err := hit(fault, StepRecoveryRef); err != nil {
			return plan, err
		}
	}
	if plan.Action == ActionDiscardRemove || plan.Action == actionReleaseRemove {
		if err := discardIgnored(plan); err != nil {
			return plan, err
		}
		if err := checkpoint(intent.ReceiptPhasePreserved); err != nil {
			return plan, err
		}
	}
	if err := checkpoint(intent.ReceiptPhaseRemoving); err != nil {
		return plan, err
	}
	interruptContext, stopInterrupts := subprocess.NotifyCancel(context.Background())
	defer stopInterrupts()
	if plan.owned {
		if out, err := exec.Command("git", "-C", root, "worktree", "unlock", plan.Target).CombinedOutput(); err != nil {
			return plan, fmt.Errorf("unlock exact assignment: %s", strings.TrimSpace(string(out)))
		}
		if err := hit(fault, StepUnlock); err != nil {
			return plan, errors.Join(err, relock(root, *plan.assignment, fault))
		}
	}
	force := plan.Tracked != "clean" || plan.Action == ActionDiscardRemove || plan.registration.Detached
	removeArgs := []string{"-C", root, "worktree", "remove"}
	if force {
		removeArgs = append(removeArgs, "--force")
	}
	removeArgs = append(removeArgs, plan.Target)
	if err := hit(fault, StepRemovalAttempt); err != nil {
		if plan.owned {
			err = errors.Join(err, relock(root, *plan.assignment, fault))
		}
		return plan, err
	}
	if interruptContext.Err() != nil {
		interrupted := error(ErrCleanupInterrupted)
		if plan.owned {
			interrupted = errors.Join(interrupted, relock(root, *plan.assignment, fault))
		}
		return plan, interrupted
	}
	if out, err := exec.CommandContext(interruptContext, "git", removeArgs...).CombinedOutput(); err != nil {
		if interruptContext.Err() != nil {
			interrupted := error(ErrCleanupInterrupted)
			if plan.owned {
				interrupted = errors.Join(interrupted, relock(root, *plan.assignment, fault))
			}
			return plan, interrupted
		}
		removeErr := fmt.Errorf("remove exact worktree: %s", strings.TrimSpace(string(out)))
		if plan.owned {
			removeErr = errors.Join(removeErr, relock(root, *plan.assignment, fault))
		}
		return plan, removeErr
	}
	stopInterrupts()
	if err := checkpoint(intent.ReceiptPhaseRemoved); err != nil {
		return plan, err
	}
	if err := hit(fault, StepRemoval); err != nil {
		return plan, err
	}
	if plan.deleteBranch {
		if err := git.DeleteBranchExact(root, plan.branchRef, plan.branchOID); err != nil {
			return plan, fmt.Errorf("delete exact landed branch: %w", err)
		}
		if err := checkpoint(intent.ReceiptPhaseBranch); err != nil {
			return plan, err
		}
		if err := hit(fault, StepBranch); err != nil {
			return plan, err
		}
	}
	if recovered != nil {
		recovered.State = intent.StateRecovered
		if err := intent.PutAssignment(root, *recovered); err != nil {
			return plan, err
		}
	} else if plan.owned && plan.assignment != nil {
		complete := *plan.assignment
		complete.State = intent.StateComplete
		if err := intent.PutAssignment(root, complete); err != nil {
			return plan, err
		}
	}
	if err := checkpoint(intent.ReceiptPhaseTerminal); err != nil {
		return plan, err
	}
	plan.Action, plan.Reason, plan.ReasonCode = ActionRemoved, "", ""
	return plan, nil
}
