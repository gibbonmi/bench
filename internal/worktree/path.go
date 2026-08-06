package worktree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/usage"
)

// PathCommand resolves one active Bench-owned assignment and prints its portable path.
func PathCommand(root string, args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		fmt.Fprintln(stderr, "usage: "+usage.WorktreePath)
		return 2
	}
	path, err := resolvePath(root, args[0])
	if err != nil {
		fmt.Fprintln(stderr, "bench worktree path: target is not one active Bench-owned worktree")
		return 1
	}
	if !lineSafe(path) {
		fmt.Fprintln(stderr, "bench worktree path: resolved path is not safe for line output")
		return 1
	}
	fmt.Fprintln(stdout, path)
	return 0
}

func resolvePath(root, target string) (string, error) {
	path, err := resolveWorktree(root, target)
	if err != nil {
		return "", err
	}
	return compactHomePath(path)
}

func resolveWorktree(root, target string) (string, error) {
	if !lineSafe(target) {
		return "", errors.New("target contains control characters")
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return "", err
	}
	selected, err := selectAssignment(assignments, target)
	if err != nil {
		return "", err
	}
	if selected.State != intent.StateActive {
		return "", errors.New("assignment is not active")
	}
	if err := validateCreationBundle(root, selected); err != nil {
		return "", err
	}
	return selected.Worktree, nil
}

var errAbandonMismatch = errors.New("abandon request, assignment, or path mismatch; checkout retained")

// errSpecialMetadata refuses an abandon over a special .git entry — a FIFO, device,
// socket, or symlink — before anything invokes git against the path. The shape answers
// for no checkout and is not leftover residue either: routing it to either would either
// hang on an unwritten FIFO or dispose of bytes nothing has proven safe to touch.
var errSpecialMetadata = errors.New("assignment .git entry is a special file; refusing before invoking git")

func planAbandon(root, request, path string) (CleanupPlan, error) {
	plan, _, err := planAbandonWithPlanner(root, request, path)
	return plan, err
}

// planAbandonWithPlanner returns the abandon plan together with the planner that
// reproduces it, or a nil planner when the ordinary explicit planner does. The apply
// transaction re-plans the target at each checkpoint, and every probe PlanExplicit runs
// is rooted in the checkout itself, so a plan decided without that checkout has to carry
// its own planner forward or the apply would fail where the plan succeeded.
func planAbandonWithPlanner(root, request, path string) (CleanupPlan, cleanupPlanner, error) {
	repo, target, err := cleanupIdentity(root, path)
	if err != nil {
		return CleanupPlan{}, nil, err
	}
	digest := requestDigest(request)
	if plan, found, err := abandonReceipt(root, repo, target, digest); err != nil || found {
		return plan, nil, err
	}
	shape, err := ClassifyPathShape(target)
	if err != nil {
		return CleanupPlan{}, nil, err
	}
	if shape == ShapeAbsent {
		planner := func(target string) (CleanupPlan, error) { return planRemovedCheckout(root, target, digest) }
		plan, err := planner(target)
		return plan, planner, err
	}
	if shape.decayed() {
		planner := func(target string) (CleanupPlan, error) { return planLeftoverEntry(root, target, digest) }
		plan, err := planner(target)
		return plan, planner, err
	}
	if shape == ShapeSpecialMetadata {
		return CleanupPlan{}, nil, errSpecialMetadata
	}
	plan, err := PlanExplicit(root, target)
	if err != nil {
		return CleanupPlan{}, nil, err
	}
	if !plan.owned || plan.assignment == nil || plan.assignment.Request != digest {
		return CleanupPlan{}, nil, errAbandonMismatch
	}
	return plan, nil, nil
}

// abandonRegistration reconciles the live worktree registration at target against the
// intent ledger for one request. It is the whole identity proof available to an abandon
// decided without a checkout: the registration must name the target, and its branch and
// Bench lock must match the assignment this request owns, which is the same identity pair
// the explicit plan asserts.
func abandonRegistration(root, target, request string) (git.Worktree, intent.Assignment, error) {
	registrations, err := git.Worktrees(canonicalRoot(root))
	if err != nil {
		return git.Worktree{}, intent.Assignment{}, err
	}
	var registration *git.Worktree
	for i := range registrations {
		registered, pathErr := canonicalPath(registrations[i].Path)
		if pathErr != nil || registered != target {
			continue
		}
		if registration != nil {
			return git.Worktree{}, intent.Assignment{}, errAbandonMismatch
		}
		registration = &registrations[i]
	}
	assignment, found, err := intent.FindAssignmentByRequest(root, request)
	if err != nil {
		return git.Worktree{}, intent.Assignment{}, err
	}
	assignmentPath, pathErr := canonicalPath(assignment.Worktree)
	if registration == nil || !found || pathErr != nil || assignmentPath != target || !cleanupOutputSafe(target) ||
		registration.Detached || registration.BranchRef != assignment.Branch || registration.LockReason != lockReason(assignment) {
		return git.Worktree{}, intent.Assignment{}, errAbandonMismatch
	}
	return *registration, assignment, nil
}

// planLeftoverEntry plans the abandon of an assignment whose path still holds bytes that
// are not a checkout. The registration and the ledger entry are the whole of what this
// abandon spends: the bytes are left exactly as they are, travelling in the plan as the
// leftover, because nothing here can tell what they hold and no recovery ref stands behind
// them. The identity proof is the removed-checkout path's, so releasing a registration
// that belongs to different work refuses the same way.
func planLeftoverEntry(root, target, request string) (CleanupPlan, error) {
	registration, assignment, err := abandonRegistration(root, target, request)
	if err != nil {
		return CleanupPlan{}, err
	}
	plan := CleanupPlan{
		Target: target, Action: actionReleaseLeftover, Tracked: "unknown", Recovery: "none",
		registration: registration, owned: true, assignment: &assignment, leftover: target,
	}
	// A recorded recovery ref names work that outlived the checkout, so the release
	// re-asserts it rather than completing the assignment out from under it.
	if len(assignment.Recovery) > 0 {
		plan.Recovery = assignment.Recovery[0].Ref
	}
	plan.Fingerprint = leftoverFingerprint(plan.leftover, registration, assignment, plan.Action, plan.Recovery)
	return plan, nil
}

func leftoverFingerprint(leftover string, registration git.Worktree, assignment intent.Assignment, action CleanupAction, recovery string) string {
	return fingerprintParts(
		[]byte("bench-abandon-leftover/v1"), []byte(leftover),
		[]byte(registration.BranchRef), []byte(registration.LockReason),
		[]byte(assignment.Schema), []byte(assignment.OwnerID), []byte(assignment.ID), []byte(assignment.Request),
		[]byte(assignment.Start), []byte(assignment.Branch), []byte(assignment.State),
		[]byte(action), []byte(recovery),
	)
}

// planRemovedCheckout plans the abandon of an assignment whose checkout is gone from disk
// while this repository still registers it. Only a path Lstat reports absent reaches here;
// a path holding anything at all is planned from its checkout or released as leftover, so
// a stranger's repository sitting at the target is never mistaken for an absent one.
// Absence alone still licenses nothing — abandonRegistration owns that proof.
//
// The branch is never deleted here: landedness is derived from the checkout, and a
// removed one cannot supply the proof.
func planRemovedCheckout(root, target, request string) (CleanupPlan, error) {
	registration, assignment, err := abandonRegistration(root, target, request)
	if err != nil {
		return CleanupPlan{}, err
	}
	plan := CleanupPlan{
		Target: target, Action: ActionRemove, Tracked: "clean", Recovery: "none",
		registration: registration, owned: true, assignment: &assignment,
	}
	// Recovery refs already recorded name work that outlived the checkout, so the removal
	// must re-assert them rather than complete the assignment out from under them.
	if len(assignment.Recovery) > 0 {
		plan.Action, plan.Recovery = ActionRecoverRemove, assignment.Recovery[0].Ref
	}
	plan.Fingerprint = fingerprintParts(
		[]byte("bench-abandon-removed/v1"), []byte(target),
		[]byte(registration.BranchRef), []byte(registration.LockReason),
		[]byte(assignment.Schema), []byte(assignment.OwnerID), []byte(assignment.ID), []byte(assignment.Request),
		[]byte(assignment.Start), []byte(assignment.Branch), []byte(assignment.State),
		[]byte(plan.Action), []byte(plan.Recovery),
	)
	return plan, nil
}

func abandonReceipt(root, repo, target, request string) (CleanupPlan, bool, error) {
	receipt, found, err := intent.CleanupReceiptForRequest(root, repo, cleanupOperation, target, request)
	if err != nil || !found {
		return CleanupPlan{}, found, err
	}
	if receipt.State != intent.ReceiptInFlight && receipt.State != intent.ReceiptComplete {
		return CleanupPlan{}, true, errAbandonMismatch
	}
	assigned, active, err := intent.FindAssignmentByRequest(root, request)
	if err != nil {
		return CleanupPlan{}, true, err
	}
	if !active {
		if receipt.State == intent.ReceiptComplete {
			return planFromReceipt(receipt), true, nil
		}
		return CleanupPlan{}, true, errAbandonMismatch
	}
	assignmentPath, err := canonicalPath(assigned.Worktree)
	if err != nil || assigned.ID != receipt.Assignment || assigned.OwnerID != receipt.Owner || assigned.Request != request || assignmentPath != target {
		return CleanupPlan{}, true, errAbandonMismatch
	}
	return planFromReceipt(receipt), true, nil
}

// PlanAbandon returns the exact recovery-aware cleanup fingerprint for one owned assignment.
func PlanAbandon(root, request, path string) (string, error) {
	plan, err := planAbandon(root, request, path)
	return plan.Fingerprint, err
}

// ApplyAbandon preserves Git-visible work before removing one exact owned assignment.
func ApplyAbandon(root, request, path, fingerprint string) (CleanupPlan, error) {
	plan, planner, err := planAbandonWithPlanner(root, request, path)
	if err != nil {
		return plan, err
	}
	if plan.Fingerprint != fingerprint {
		return plan, errStaleFingerprint
	}
	if planner == nil {
		planner = func(target string) (CleanupPlan, error) { return PlanExplicit(root, target) }
	}
	applied, err := applyCleanupTransaction(root, path, plan.Fingerprint, planner, nil, nil)
	if err != nil || applied.Action != ActionRetain {
		return applied, err
	}
	return applied, retainedReleaseError(applied, request, path)
}

// ProvisionalEvidence identifies the durable commits that already preserve an
// assignment payload outside its disposable checkout.
type ProvisionalEvidence struct {
	Base, CheckpointRef, Checkpoint, IntegratedRef, Integrated string
}

// ReleaseProvisional removes one owned assignment whose exact payload remains
// retained by durable checkpoint and integration commits.
func ReleaseProvisional(root, requestArg, path string, evidence ProvisionalEvidence) error {
	target, err := canonicalPath(path)
	if err != nil {
		return err
	}
	repo, _, err := cleanupIdentity(root, target)
	if err != nil {
		return err
	}
	request := requestDigest(requestArg)
	if receipt, found, readErr := intent.CleanupReceiptFor(root, repo, releaseOperation, target, request); readErr != nil {
		return readErr
	} else if found && receipt.State == intent.ReceiptComplete {
		return compactProvisionalAssignment(root, receipt.Tracked)
	}
	assignment, ok, err := intent.FindAssignmentByRequest(root, request)
	if err != nil || !ok || assignment.Worktree != target {
		return errors.New("provisional release request, assignment, or path mismatch; checkout retained")
	}
	if assignment.State != intent.StateActive && assignment.State != intent.StateCleanupPending {
		return errors.New("provisional release assignment state does not accept cleanup")
	}
	if assignment.State == intent.StateActive {
		if err := validateCreationBundle(root, assignment); err != nil {
			return fmt.Errorf("validate provisional release owner: %w", err)
		}
		if _, err := planProvisionalRelease(root, target, assignment.ID, evidence); err != nil {
			return err
		}
		assignment.State = intent.StateCleanupPending
		if err := intent.PutAssignment(root, assignment); err != nil {
			return err
		}
	}
	planner := func(path string) (CleanupPlan, error) {
		return planProvisionalRelease(root, path, assignment.ID, evidence)
	}
	plan, err := planner(target)
	if err != nil {
		return err
	}
	terminal := func(plan CleanupPlan) error {
		current, readErr := assignmentByID(root, assignment.ID)
		if readErr != nil {
			return readErr
		}
		return intent.PutCleanupReceipt(root, receiptFromRelease(repo, request, current, string(plan.Action)))
	}
	if _, err := applyCleanupTransaction(root, target, plan.Fingerprint, planner, nil, terminal); err != nil {
		return err
	}
	return compactProvisionalAssignment(root, assignment.ID)
}

func planProvisionalRelease(root, target, assignmentID string, evidence ProvisionalEvidence) (CleanupPlan, error) {
	plan, err := PlanExplicit(root, target)
	if err != nil {
		return plan, err
	}
	if !plan.owned || plan.assignment == nil || plan.assignment.ID != assignmentID || plan.Action == ActionRetain {
		return plan, errors.New("provisional release checkout is not exact and removable; checkout retained")
	}
	checkpointCommit, err := validateProvisionalEvidence(root, target, evidence)
	if err != nil {
		return plan, err
	}
	if checkpointCommit {
		if plan.Action != ActionRemove || plan.Tracked != "clean" {
			return plan, errors.New("provisional release checkout is not exact and removable; checkout retained")
		}
		return plan, nil
	}
	cleanNoOp := plan.Tracked == "clean" && plan.Action == ActionRemove
	dirtyPayload := plan.Tracked == "dirty" && (plan.Action == ActionRecoverRemove || plan.Action == ActionDiscardRemove)
	if !cleanNoOp && !dirtyPayload {
		return plan, errors.New("provisional release checkout is not exact and removable; checkout retained")
	}
	original := plan.Fingerprint
	plan.Action = actionReleaseRemove
	plan.Recovery = "none"
	plan.Fingerprint = fingerprintParts(
		[]byte("bench-provisional-release/v1"), []byte(original), []byte(evidence.Base),
		[]byte(evidence.CheckpointRef), []byte(evidence.Checkpoint), []byte(evidence.IntegratedRef), []byte(evidence.Integrated),
		[]byte(plan.Action),
	)
	return plan, nil
}

func validateProvisionalEvidence(root, target string, evidence ProvisionalEvidence) (bool, error) {
	if evidence.Base == "" || evidence.CheckpointRef == "" || evidence.Checkpoint == "" || evidence.IntegratedRef == "" || evidence.Integrated == "" {
		return false, errors.New("provisional release evidence is incomplete")
	}
	paths, pathErr := git.Raw("--no-optional-locks", "-C", target, "ls-files", "--cached", "--others", "--exclude-standard", "-z", "--")
	for record := range strings.SplitSeq(string(paths), "\x00") {
		if record != "" && !lineSafe(record) {
			return false, errors.New("provisional release checkout contains an unsafe path; checkout retained")
		}
	}
	retained, refErr := git.Output("-C", root, "rev-parse", "--verify", evidence.CheckpointRef+"^{commit}")
	head, headErr := git.Output("-C", target, "rev-parse", "HEAD")
	headTree, headTreeErr := git.Output("-C", target, "rev-parse", "HEAD^{tree}")
	baseTree, baseTreeErr := git.Output("-C", root, "rev-parse", evidence.Base+"^{tree}")
	indexTree, indexErr := git.Output("-C", target, "write-tree")
	checkpointTree, checkpointTreeErr := git.Output("-C", root, "rev-parse", evidence.Checkpoint+"^{tree}")
	checkpointParent, checkpointParentErr := git.Output("-C", root, "rev-parse", evidence.Checkpoint+"^")
	integratedRetained, integratedRefErr := git.Output("-C", root, "rev-parse", "--verify", evidence.IntegratedRef+"^{commit}")
	_, integratedTreeErr := git.Output("-C", root, "rev-parse", evidence.Integrated+"^{tree}")
	liveTree := git.TreeHash(target)
	// Releasing the checkout loses nothing exactly when its live content is the checkpoint
	// tree and the checkpoint commit is still retained at its ref. The integrated tree is
	// deliberately not compared against the checkpoint tree: a sibling replayed onto a
	// candidate an earlier sibling already advanced integrates a different tree by
	// construction, and that tree says nothing about what the checkout holds.
	if pathErr != nil || refErr != nil || headErr != nil || headTreeErr != nil || baseTreeErr != nil || indexErr != nil || checkpointTreeErr != nil || checkpointParentErr != nil || integratedRefErr != nil || integratedTreeErr != nil ||
		retained != evidence.Checkpoint || checkpointParent != evidence.Base || integratedRetained != evidence.Integrated || liveTree == "none" || liveTree != checkpointTree {
		return false, errors.New("provisional release evidence drifted; checkout retained")
	}
	if head == evidence.Base {
		if indexTree != baseTree {
			return false, errors.New("provisional release evidence drifted; checkout retained")
		}
		return false, nil
	}
	if indexTree != headTree || headTree != checkpointTree {
		return false, errors.New("provisional release evidence drifted; checkout retained")
	}
	return true, nil
}

func compactProvisionalAssignment(root, assignmentID string) error {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return err
	}
	var assignment *intent.Assignment
	for i := range assignments {
		if assignments[i].ID == assignmentID {
			assignment = &assignments[i]
			break
		}
	}
	if assignment == nil {
		return nil
	}
	if assignment.State != intent.StateComplete {
		return errors.New("provisional release did not reach terminal assignment state")
	}
	return intent.DeleteAssignment(root, assignment.ID)
}

func selectAssignment(assignments []intent.Assignment, target string) (intent.Assignment, error) {
	path, isPath, err := targetPath(target)
	if err != nil {
		return intent.Assignment{}, err
	}
	var selected *intent.Assignment
	for i := range assignments {
		matches := assignments[i].Label == target
		if isPath {
			worktree, worktreeErr := canonicalPath(assignments[i].Worktree)
			matches = worktreeErr == nil && worktree == path
		}
		if !matches {
			continue
		}
		if selected != nil {
			return intent.Assignment{}, errors.New("target is ambiguous")
		}
		selected = &assignments[i]
	}
	if selected == nil {
		return intent.Assignment{}, errors.New("target is unassigned")
	}
	return *selected, nil
}

func targetPath(target string) (string, bool, error) {
	if target == "~" || strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", false, err
		}
		if target == "~" {
			path, err := canonicalPath(home)
			return path, true, err
		}
		path, err := canonicalPath(filepath.Join(home, target[2:]))
		return path, true, err
	}
	if strings.HasPrefix(target, "~") {
		return "", false, errors.New("unsupported home target")
	}
	if filepath.IsAbs(target) {
		path, err := canonicalPath(target)
		return path, true, err
	}
	if target == "." || strings.HasPrefix(target, "."+string(filepath.Separator)) || strings.HasPrefix(target, ".."+string(filepath.Separator)) || strings.Contains(target, string(filepath.Separator)) {
		return "", false, errors.New("relative path targets are unsupported")
	}
	return "", false, nil
}

func compactHomePath(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	home, err = canonicalPath(home)
	if err != nil {
		return "", err
	}
	path, err = canonicalPath(path)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(home, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return path, nil
	}
	if rel == "." {
		return "~", nil
	}
	return filepath.Join("~", rel), nil
}
