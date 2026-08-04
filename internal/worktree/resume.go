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

// errRecoveryUnauthorized marks a refusal the plan itself decided, as opposed to a failure
// while acting, so the command can answer with the receipt the operator planned from.
var errRecoveryUnauthorized = errors.New("recovery plan does not authorize this verb")

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

// recoveryVerb is the flag an operator typed, carried through to the one place that acts
// on a plan. Its value is the flag itself so the grammar and the authority stay one fact.
type recoveryVerb string

const (
	recoveryRetire  recoveryVerb = "--apply"
	recoveryDiscard recoveryVerb = "--discard"
)

// authorizes reports whether a plan carrying this action licenses the verb, and the
// refusal detail when it does not. The two verbs partition the vocabulary: only the
// landedness proof's own verdict retires, and only the verdicts that still hold work the
// proof judged and refused — discard-eligible, and orphaned, which has no row left to
// judge through — are the operator's to discard. Every other verdict refuses both:
// retain means the plan could not classify the ref, and a verdict that proved nothing
// must carry no destructive authority, however the fingerprint arrived.
func (verb recoveryVerb) authorizes(action RecoveryAction) (bool, string) {
	switch action {
	case RecoveryRetire:
		if verb == recoveryRetire {
			return true, ""
		}
		return false, "a proven-landed payload retires with --apply, not --discard"
	case RecoveryDiscard, RecoveryOrphaned:
		if verb == recoveryDiscard {
			return true, ""
		}
		return false, "only a proven-landed payload retires; drop unproven work with --discard"
	case RecoveryForeign:
		return false, "only a ref inside the recovery namespace is the operator's to discard"
	case RecoveryRetain:
		return false, "the plan could not classify this ref; neither verb is authorized"
	default:
		return false, "plan action holds no discardable work"
	}
}

// terminal names the claim a completed verb records in the receipt.
func (verb recoveryVerb) terminal() RecoveryAction {
	if verb == recoveryRetire {
		return RecoveryRetired
	}
	return RecoveryDiscarded
}

// ApplyRecovery retires a recovery ref whose payloads the landedness proof accepts.
func ApplyRecovery(root, ref, fingerprint string) (RecoveryPlan, error) {
	return applyRecoveryVerb(root, ref, fingerprint, recoveryRetire)
}

// applyRecoveryVerb is the one actor on a recovery plan. Both verbs spend the same
// fingerprint over the same plan and share the ref deletion and the row compaction; only
// which verdicts they accept and which claim they record differ.
func applyRecoveryVerb(root, ref, fingerprint string, verb recoveryVerb) (RecoveryPlan, error) {
	plan, err := PlanRecovery(root, ref)
	if err != nil {
		return plan, err
	}
	if plan.Fingerprint != fingerprint {
		return plan, errStaleFingerprint
	}
	// A ref that is already gone is what a completed discard leaves behind, so a re-run
	// converges on success rather than refusing work nobody can still do.
	if verb == recoveryDiscard && plan.Action == RecoveryAbsent {
		return plan, nil
	}
	// A recovered row naming a ref nothing resolves is what an interruption between the
	// two halves of either verb leaves behind: the ref delete landed and the row close did
	// not. Closing the row is all that remains, and it happens before the authorization
	// check because the vanished ref is exactly what makes the plan unclassifiable — asking
	// the landedness proof about it would refuse the only command that can finish the work.
	// The claim recorded is the discard for both verbs: retired asserts the proof accepted
	// the payload, and no proof can run over a ref that no longer resolves, so this receipt
	// can only honestly say the work is gone without the proof's backing.
	if plan.assignment != nil && plan.assignment.State == intent.StateRecovered && !refExists(root, plan.Ref) {
		if err := compactRecoveredAssignment(root, *plan.assignment, ref); err != nil {
			return plan, err
		}
		plan.Action, plan.Detail = RecoveryDiscarded, ""
		return plan, nil
	}
	if authorized, detail := verb.authorizes(plan.Action); !authorized {
		plan.Detail = detail
		return plan, fmt.Errorf("%w: %s", errRecoveryUnauthorized, detail)
	}
	// Retire reaches an assignment only in the recovered state, because the plan refuses
	// every other one before it can prove landedness. Discard holds itself to the same
	// bar: closing a row mid-release would spend a transaction another command owns.
	if plan.assignment != nil && plan.assignment.State != intent.StateRecovered {
		plan.Detail = "recovery ref has no recovered assignment"
		return plan, fmt.Errorf("%w: %s", errRecoveryUnauthorized, plan.Detail)
	}
	if err := deleteRecoveryRef(root, plan); err != nil {
		return plan, err
	}
	if err := hit(cleanupTransactionBoundary, StepRecoveryRowClose); err != nil {
		return plan, err
	}
	// An orphaned ref has no row to close, and inventing one would record an intent that
	// no longer exists.
	if plan.assignment != nil {
		if err := compactRecoveredAssignment(root, *plan.assignment, ref); err != nil {
			return plan, err
		}
	}
	plan.Action, plan.Detail = verb.terminal(), ""
	return plan, nil
}

// deleteRecoveryRef removes the ref only while it still holds the object the plan
// classified, so a ref something else moved is refused rather than dropped blind.
func deleteRecoveryRef(root string, plan RecoveryPlan) error {
	expected := plan.Root
	if plan.assignment == nil {
		// No row records an orphan's root, so the ref itself is the only thing naming the
		// object the plan just read.
		resolved, err := git.Output("-C", root, "rev-parse", "--verify", plan.Ref+"^{commit}")
		if err != nil {
			return fmt.Errorf("resolve orphaned recovery ref: %w", err)
		}
		expected = resolved
	}
	if out, err := exec.Command("git", "-C", root, "update-ref", "-d", plan.Ref, expected).CombinedOutput(); err != nil {
		return fmt.Errorf("delete exact recovery ref: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// compactRecoveredAssignment drops one retired or discarded ref from its row, closing the
// row entirely once nothing preserved is left to point at.
func compactRecoveredAssignment(root string, assignment intent.Assignment, ref string) error {
	next := assignment.Recovery[:0]
	for _, candidate := range assignment.Recovery {
		if candidate.Ref != ref {
			next = append(next, candidate)
		}
	}
	assignment.Recovery = next
	if len(next) > 0 {
		assignment.State = intent.StateRecovered
		return intent.PutAssignment(root, assignment)
	}
	assignment.State = intent.StateComplete
	if err := intent.PutAssignment(root, assignment); err != nil {
		return err
	}
	return intent.DeleteAssignment(root, assignment.ID)
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

// sweepOrphanAssignments turns the ledger's own view of each record into a verdict. It
// asks orphaned directly rather than reading a cleanup plan's reason code: PlanAutomatic
// returns at its first retain reason, and ignored build output is the normal state of a
// worktree a shift ran in, so a plan-derived sweep would report nothing for exactly the
// population this exists for.
//
// An active record is swept only once it is orphaned; a younger one is left alone
// because a live session may still own it. An orphan whose tree survives is reported as
// an OrphanCandidate and never touched; that type owns why removal stays behind an
// explicit command. The tree-gone verdicts are the FT93(c) contract, in this order: one
// git still registers belongs to the prune path, one preserving no work is compacted and
// counted, and one holding recovery metadata is reported for a deliberate
// recover-or-retire and left intact. residualAssignment is the single guard between that
// compaction and a pointer to preserved work.
func sweepOrphanAssignments(root string, registered []Registered, result *ResumeResult) error {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return err
	}
	// One instant for the whole pass, so two records of the same age cannot straddle the
	// window and disagree.
	now := time.Now()
	for _, a := range assignments {
		abandoned := orphaned(a, now)
		if a.State == intent.StateActive && !abandoned {
			continue
		}
		_, statErr := os.Stat(a.Worktree)
		if statErr == nil {
			if abandoned {
				result.Orphans = append(result.Orphans, OrphanCandidate{ID: a.ID, Path: a.Worktree})
			}
			continue // the tree still exists
		}
		// Only absence licenses a tree-gone verdict. Any other stat error — an unreadable
		// pool, an I/O failure — leaves the tree's existence unknown, and compacting on
		// unknown deletes the row while the worktree and its uncommitted work are still
		// there. Unknown is left out of the listing too: the line would name a retirement
		// command for a path this host cannot reach, and the record stays visible in the
		// open-assignment count either way.
		if !os.IsNotExist(statErr) {
			continue
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
