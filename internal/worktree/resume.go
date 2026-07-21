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
func ApplyRecovery(root, ref, fingerprint string) (RecoveryPlan, error) {
	plan, err := PlanRecovery(root, ref)
	if err != nil {
		return plan, err
	}
	if plan.Fingerprint != fingerprint {
		return plan, errStaleFingerprint
	}
	if plan.Action == RecoveryRetain {
		return plan, nil
	}
	if out, err := exec.Command("git", "-C", root, "update-ref", "-d", ref, plan.Root).CombinedOutput(); err != nil {
		return plan, fmt.Errorf("delete exact recovery ref: %s", strings.TrimSpace(string(out)))
	}
	assignment := *plan.assignment
	next := assignment.Recovery[:0]
	for _, candidate := range assignment.Recovery {
		if candidate.Ref != ref {
			next = append(next, candidate)
		}
	}
	assignment.Recovery = next
	if len(next) > 0 {
		assignment.State = intent.StateRecovered
		if err := intent.PutAssignment(root, assignment); err != nil {
			return plan, err
		}
	} else {
		assignment.State = intent.StateComplete
		if err := intent.PutAssignment(root, assignment); err != nil {
			return plan, err
		}
		if err := intent.DeleteAssignment(root, assignment.ID); err != nil {
			return plan, err
		}
	}
	plan.Action, plan.Detail = RecoveryRetired, ""
	return plan, nil
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
		if _, statErr := os.Lstat(target); errors.Is(statErr, os.ErrNotExist) {
			return finishInterruptedExplicit(root, receipt, terminal, fault)
		} else if statErr != nil {
			return planFromReceipt(receipt), statErr
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
func finishInterruptedExplicit(root string, receipt intent.CleanupReceipt, terminal cleanupTerminal, fault Fault) (CleanupPlan, error) {
	plan := planFromReceipt(receipt)
	switch receipt.Phase {
	case intent.ReceiptPhaseRemoving, intent.ReceiptPhaseRemoved, intent.ReceiptPhaseBranch, intent.ReceiptPhaseTerminal:
	default:
		return plan, errStaleFingerprint
	}
	worktrees, err := git.Worktrees(root)
	if err != nil {
		return plan, err
	}
	for _, worktree := range worktrees {
		if samePath(worktree.Path, receipt.Target) {
			return plan, errStaleFingerprint
		}
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
		Request: requestDigest("foreign:" + plan.Target), Label: "foreign exact cleanup", Start: head,
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

// ConservativeCleanup cleans owned worktrees and unclaimed landed branch residue.
func ConservativeCleanup(root string) (ResumeResult, error) {
	registered, err := ClassifyRegisteredWorktrees(root)
	if err != nil {
		return ResumeResult{}, fmt.Errorf("git worktree list failed: %w", err)
	}
	result := ResumeResult{Retained: map[CleanupReason]int{}}
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
		if plan.Action == ActionRecoverRemove {
			result.Recovered++
		} else {
			result.Removed++
		}
	}
	if err := sweepOrphanAssignments(root, registered, &result); err != nil {
		return result, err
	}
	result.PrunedBranches, err = intent.PruneUnclaimedLandedBranches(root)
	return result, err
}

// sweepOrphanAssignments reconciles tree-gone records absent from registered worktrees.
// Residue records are compacted and counted; records that still hold preserved work are
// reported for a deliberate recover-or-retire and left intact. Active records are never
// swept — a live session owns them. See specs/worktree-orphan-reconcile.md (c).
func sweepOrphanAssignments(root string, registered []Registered, result *ResumeResult) error {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return err
	}
	for _, a := range assignments {
		if a.State == intent.StateActive {
			continue
		}
		if _, statErr := os.Stat(a.Worktree); statErr == nil {
			continue // the tree still exists
		}
		if isRegisteredWorktree(registered, a.Worktree) {
			continue // registered (prunable) — the git-worktree path owns it
		}
		if residualAssignment(a) {
			if err := intent.DeleteAssignment(root, a.ID); err != nil {
				return err
			}
			result.Reconciled++
			continue
		}
		result.Preserved = append(result.Preserved, PreservedOrphan{ID: a.ID, Ref: a.Recovery[0].Ref})
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
