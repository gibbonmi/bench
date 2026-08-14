package worktree

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var errStaleFingerprint = errors.New("cleanup fingerprint is stale")

const cleanupOperation = "worktree-clean"

func cleanupIdentity(root, path string) (string, string, error) {
	target, err := canonicalPath(path)
	if err != nil {
		return "", "", err
	}
	address, err := intent.Address(root)
	if err != nil {
		return "", "", err
	}
	return filepath.Dir(address), target, nil
}
func receiptFromPlan(repo string, plan CleanupPlan, state string) intent.CleanupReceipt {
	receipt := intent.CleanupReceipt{
		Schema: intent.CleanupReceiptSchema, Repo: repo, Operation: cleanupOperation,
		Target: plan.Target, Fingerprint: plan.Fingerprint, State: state,
		Phase:  intent.ReceiptPhasePlanned,
		Action: string(plan.Action), Tracked: plan.Tracked, Ignored: plan.Ignored.Summary(),
		Recovery: plan.Recovery, Detail: plan.Reason, Owned: plan.owned,
	}
	if plan.owned && plan.assignment != nil && plan.deleteBranch {
		receipt.Branch, receipt.BranchOID = plan.assignment.Branch, plan.branchOID
	}
	if plan.owned && plan.assignment != nil {
		receipt.Owner, receipt.Assignment, receipt.Request = plan.assignment.OwnerID, plan.assignment.ID, plan.assignment.Request
	}
	if state == intent.ReceiptComplete {
		receipt.Phase = intent.ReceiptPhaseTerminal
	}
	return receipt
}
func planFromReceipt(receipt intent.CleanupReceipt) CleanupPlan {
	return CleanupPlan{
		Target: receipt.Target, Action: CleanupAction(receipt.Action), Tracked: receipt.Tracked,
		Recovery: receipt.Recovery, Fingerprint: receipt.Fingerprint, Reason: receipt.Detail,
		ignoredSummary: receipt.Ignored,
	}
}
func assignmentByID(root, id string) (intent.Assignment, error) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return intent.Assignment{}, err
	}
	for _, assignment := range assignments {
		if assignment.ID == id {
			return assignment, nil
		}
	}
	return intent.Assignment{}, errors.New("assignment record disappeared")
}
func renderReleaseReceipt(stdout io.Writer, receipt intent.CleanupReceipt) int {
	return renderRelease(stdout, intent.Assignment{ID: receipt.Tracked, Worktree: receipt.Target, State: intent.AssignmentState(receipt.Detail)}, receipt.Action)
}

func ApplyExplicit(root, path, fingerprint string) (CleanupPlan, error) {
	return ApplyExplicitWithOptions(root, path, fingerprint, CleanupOptions{})
}
func ApplyExplicitWithOptions(root, path, fingerprint string, options CleanupOptions) (CleanupPlan, error) {
	planner := func(target string) (CleanupPlan, error) {
		return PlanExplicitWithOptions(root, target, options)
	}
	return applyCleanupTransaction(root, path, fingerprint, planner, nil, nil)
}

type cleanupPlanner func(string) (CleanupPlan, error)
type cleanupTerminal func(CleanupPlan) error

func applyCleanupTransaction(root, path, fingerprint string, planner cleanupPlanner, localFault Fault, terminal cleanupTerminal) (CleanupPlan, error) {
	repo, target, err := cleanupIdentity(root, path)
	if err != nil {
		return CleanupPlan{}, err
	}
	receipt, found, err := intent.CleanupReceiptFor(root, repo, cleanupOperation, target, fingerprint)
	if err != nil {
		return CleanupPlan{}, err
	}
	if found && receipt.State == intent.ReceiptComplete {
		return planFromReceipt(receipt), nil
	}
	release, err := lockCleanupRegistration(repo, target)
	if err != nil {
		return CleanupPlan{}, err
	}
	defer release()
	fault := func(step LifecycleStep) error {
		if err := hit(cleanupTransactionBoundary, step); err != nil {
			return err
		}
		return hit(localFault, step)
	}
	if err := hit(fault, StepApplyLocked); err != nil {
		return CleanupPlan{}, err
	}
	receipt, found, err = intent.CleanupReceiptFor(root, repo, cleanupOperation, target, fingerprint)
	if err != nil {
		return CleanupPlan{}, err
	}
	if found && receipt.State == intent.ReceiptComplete {
		return planFromReceipt(receipt), nil
	}
	if found && receipt.State == intent.ReceiptInFlight {
		resumable, resumeErr := interruptedCleanupIsPastReplanning(root, receipt, target)
		if resumeErr != nil {
			return planFromReceipt(receipt), resumeErr
		}
		if resumable {
			return finishInterruptedExplicit(root, receipt, terminal, fault)
		}
	}
	plan, err := planner(target)
	if err != nil {
		return plan, err
	}
	if plan.Fingerprint != fingerprint {
		if !found || receipt.State != intent.ReceiptInFlight || receipt.Checkpoint != plan.Fingerprint {
			return plan, errStaleFingerprint
		}
		plan.Fingerprint = fingerprint
	}
	if plan.Action == ActionRetain {
		if err := intent.PutCleanupReceipt(root, receiptFromPlan(repo, plan, intent.ReceiptComplete)); err != nil {
			return plan, err
		}
		return plan, nil
	}
	releasePersistence, err := lockCleanupPersistence(repo, target)
	if err != nil {
		return plan, err
	}
	defer releasePersistence()
	if !found {
		receipt = receiptFromPlan(repo, plan, intent.ReceiptInFlight)
		if err := intent.PutCleanupReceipt(root, receipt); err != nil {
			return plan, err
		}
		if err := hit(fault, StepReceipt); err != nil {
			return plan, err
		}
	}
	checkpoint := func(phase string) error {
		receipt.Phase = phase
		if phase == intent.ReceiptPhasePlanned || phase == intent.ReceiptPhasePreserved || phase == intent.ReceiptPhaseRemoving {
			current, planErr := planner(target)
			if planErr != nil {
				return planErr
			}
			receipt.Checkpoint = current.Fingerprint
		}
		return intent.PutCleanupReceipt(root, receipt)
	}
	plan, err = executeCleanup(root, plan, checkpoint, fault)
	if err != nil {
		return plan, err
	}
	return completeCleanupTransaction(root, plan, receipt, terminal, fault)
}
func completeCleanupTransaction(root string, plan CleanupPlan, receipt intent.CleanupReceipt, terminal cleanupTerminal, fault Fault) (CleanupPlan, error) {
	receipt.State, receipt.Phase, receipt.Action, receipt.Detail = intent.ReceiptComplete, intent.ReceiptPhaseTerminal, string(ActionRemoved), ""
	plan.Action = ActionRemoved
	if err := intent.PutCleanupReceipt(root, receipt); err != nil {
		return plan, err
	}
	if terminal != nil {
		if err := terminal(plan); err != nil {
			return plan, err
		}
		if err := hit(fault, StepTerminalReceipt); err != nil {
			return plan, err
		}
	}
	if assignment, readErr := assignmentByID(root, receipt.Assignment); readErr == nil && assignment.State == intent.StateComplete {
		if err := intent.DeleteAssignment(root, assignment.ID); err != nil {
			return plan, err
		}
	}
	return planFromReceipt(receipt), nil
}

// interruptedCleanupIsPastReplanning reports whether an in-flight cleanup already spent
// the thing its plan was decided from, leaving the receipt as the only thing that can
// finish it. Re-planning such a target answers a different question and refuses the retry
// on a stale fingerprint, wedging an abandon that has nothing left to do but land.
//
// A removal proves it by the tree being gone. A release-leftover never removes the bytes,
// so presence at the target says nothing about its progress; the registration is what it
// spends, and a target this repository no longer registers carries the same proof.
func interruptedCleanupIsPastReplanning(root string, receipt intent.CleanupReceipt, target string) (bool, error) {
	shape, err := ClassifyPathShape(target)
	if err != nil {
		return false, err
	}
	if shape == ShapeAbsent {
		return true, nil
	}
	if CleanupAction(receipt.Action) != actionReleaseLeftover {
		return false, nil
	}
	registered, err := registeredAt(root, target)
	return !registered, err
}

// registeredAt reports whether this repository still registers a worktree at target.
func registeredAt(root, target string) (bool, error) {
	worktrees, err := git.Worktrees(root)
	if err != nil {
		return false, err
	}
	for _, worktree := range worktrees {
		if samePath(worktree.Path, target) {
			return true, nil
		}
	}
	return false, nil
}

func finishInterruptedExplicit(root string, receipt intent.CleanupReceipt, terminal cleanupTerminal, fault Fault) (CleanupPlan, error) {
	plan := planFromReceipt(receipt)
	switch receipt.Phase {
	case intent.ReceiptPhaseRemoving, intent.ReceiptPhaseRemoved, intent.ReceiptPhaseBranch, intent.ReceiptPhaseTerminal:
	default:
		return plan, errStaleFingerprint
	}
	registered, err := registeredAt(root, receipt.Target)
	if err != nil {
		return plan, err
	}
	if registered {
		return plan, errStaleFingerprint
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return plan, err
	}
	var assignment *intent.Assignment
	for i := range assignments {
		if assignments[i].Worktree == receipt.Target {
			if assignment != nil {
				return plan, errors.New("interrupted cleanup has ambiguous assignments")
			}
			candidate := assignments[i]
			assignment = &candidate
		}
	}
	if receipt.Recovery != "" && receipt.Recovery != "none" {
		if assignment == nil {
			return plan, errors.New("interrupted cleanup lost recovery context")
		}
		verified := false
		for _, recovery := range assignment.Recovery {
			if recovery.Ref == receipt.Recovery && verifyRecovery(root, *assignment, recovery) == nil {
				verified = true
			}
		}
		if !verified {
			return plan, errors.New("interrupted cleanup recovery does not verify")
		}
	}
	if receipt.Branch != "" {
		current, oidErr := git.Output("-C", root, "rev-parse", "--verify", receipt.Branch)
		if oidErr == nil {
			if current != receipt.BranchOID {
				return plan, errors.New("interrupted cleanup branch changed")
			}
			if out, err := exec.Command("git", "-C", root, "update-ref", "-d", receipt.Branch, receipt.BranchOID).CombinedOutput(); err != nil {
				return plan, fmt.Errorf("delete interrupted assignment branch: %s", strings.TrimSpace(string(out)))
			}
		}
	}
	if assignment != nil {
		if receipt.Recovery != "" && receipt.Recovery != "none" {
			assignment.State = intent.StateRecovered
			if err := intent.PutAssignment(root, *assignment); err != nil {
				return plan, err
			}
		} else {
			assignment.State = intent.StateComplete
			if err := intent.PutAssignment(root, *assignment); err != nil {
				return plan, err
			}
		}
	}
	return completeCleanupTransaction(root, plan, receipt, terminal, fault)
}
func recoveryAssignmentForPlan(root string, plan CleanupPlan) (intent.Assignment, error) {
	if plan.assignment != nil {
		return *plan.assignment, nil
	}
	parts := strings.Split(plan.Recovery, "/")
	if len(parts) != 6 || parts[0] != "refs" || parts[1] != "bench" || parts[2] != "recovery" {
		return intent.Assignment{}, errors.New("foreign recovery identity is malformed")
	}
	head, err := git.Output("-C", plan.Target, "rev-parse", "HEAD")
	if err != nil {
		return intent.Assignment{}, err
	}
	return intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: parts[4], OwnerID: parts[3],
		Request: intent.RequestDigest("foreign:" + plan.Target), Label: "foreign exact cleanup", Start: head,
		Branch: intent.AssignmentBranchRef(parts[3], parts[4]), Worktree: plan.Target,
		State: intent.StateCleanupPending, Recovery: []intent.Recovery{},
	}, nil
}
func ApplyAutomatic(root, path string, fault Fault) (CleanupPlan, error) {
	return applyAutomaticWithTerminal(root, path, fault, nil)
}
func applyAutomaticWithTerminal(root, path string, fault Fault, terminal cleanupTerminal) (CleanupPlan, error) {
	plan, err := PlanAutomatic(root, path)
	if err != nil || plan.Action == ActionRetain {
		return plan, err
	}
	planner := func(target string) (CleanupPlan, error) { return PlanAutomatic(root, target) }
	return applyCleanupTransaction(root, path, plan.Fingerprint, planner, fault, terminal)
}

// ConservativeCleanup reconciles the lifecycle debris, then cleans owned worktrees and
// unclaimed landed branch residue. The reconcile runs first because it is the only thing
// that can make a ledger an older binary wrote readable again, and every step below reads
// that ledger.
func ConservativeCleanup(root string) (ResumeResult, error) {
	registered, err := ClassifyRegisteredWorktrees(root)
	if err != nil {
		return ResumeResult{}, fmt.Errorf("git worktree list failed: %w", err)
	}
	result := ResumeResult{Retained: map[CleanupReason]int{}}
	result.SweptRefs, result.Reconciled, err = reconcileLifecycleDebris(root, registered)
	if err != nil {
		return result, err
	}
	for _, wt := range registered {
		if wt.Class == ClassRoot {
			continue
		}
		if _, err := os.Stat(wt.Path); os.IsNotExist(err) {
			continue
		}
		plan, err := PlanAutomatic(root, wt.Path)
		if err != nil {
			return result, err
		}
		if plan.Action == ActionRetain {
			reason := plan.ReasonCode
			if reason == "" {
				reason = ReasonUncertain
			}
			result.Retained[reason]++
			continue
		}
		if _, err := ApplyAutomatic(root, wt.Path, nil); err != nil {
			result.Failed++
			return result, err
		}
		result.Removed++
	}
	if err := sweepOrphanAssignments(root, &result); err != nil {
		return result, err
	}
	result.PrunedBranches, err = intent.PruneUnclaimedLandedBranches(root)
	return result, err
}

// sweepOrphanAssignments turns the ledger's own view of each record into a verdict. It
// asks orphaned directly rather than reading a cleanup plan's reason code: PlanAutomatic
// returns at its first retain reason, and ignored build output is the normal state of a
// worktree a shift ran in, so a plan-derived sweep would report nothing for exactly the
// population this exists for.
//
// An active record is swept only once it is orphaned; a younger one is left alone
// because a live session may still own it. An orphan whose tree survives is reported as
// an OrphanCandidate and never touched; that type owns why removal stays behind an
// explicit command. A record whose tree is gone is not this sweep's to judge — the
// reconcile that already ran is the one place a record is dropped, so the two can never
// disagree about which ones the pool still answers for.
func sweepOrphanAssignments(root string, result *ResumeResult) error {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return err
	}
	// One instant for the whole pass, so two records of the same age cannot straddle the
	// window and disagree.
	now := time.Now()
	for _, a := range assignments {
		if !orphaned(a, now) {
			continue
		}
		if _, statErr := os.Stat(a.Worktree); statErr == nil {
			result.Orphans = append(result.Orphans, OrphanCandidate{ID: a.ID, Path: a.Worktree})
		}
	}
	return nil
}

func isRegisteredWorktree(registered []Registered, path string) bool {
	for _, wt := range registered {
		if samePath(wt.Path, path) {
			return true
		}
	}
	return false
}
func ResumeCleanCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "usage: bench resume-clean")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	result, cleanupErr := ConservativeCleanup(root)
	assignments, snapshotErr := intent.Assignments(root)
	if snapshotErr == nil {
		result.Open = len(assignments)
	}
	fmt.Fprint(stdout, renderResumeSummary(result))
	if err := errors.Join(cleanupErr, snapshotErr); err != nil {
		fmt.Fprintf(stderr, "bench resume-clean: %v\n", err)
		return 1
	}
	return 0
}
