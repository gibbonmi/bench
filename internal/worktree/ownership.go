package worktree

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/toon"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	OwnerMarkerSchema = "bench-owner/v1"
	OwnerMarkerFile   = "bench-owner"
	ignoredEntryLimit = 1000
	ignoredByteLimit  = int64(1 << 30)
)

type Marker struct {
	Schema  string `json:"schema"`
	OwnerID string `json:"owner_id"`
	Path    string `json:"path"`
}
type Creation struct {
	Path       string
	Assignment intent.Assignment
}

func requestDigest(request string) string { return textDigest(request) }

type LifecycleStep string

const (
	StepRegistration     LifecycleStep = "registration"
	StepMarker           LifecycleStep = "marker"
	StepRecord           LifecycleStep = "record"
	StepRollbackRemove   LifecycleStep = "rollback-remove"
	StepRelock           LifecycleStep = "relock"
	StepRecoveryMetadata LifecycleStep = "recovery-metadata"
	StepRecoveryRef      LifecycleStep = "recovery-ref"
	StepUnlock           LifecycleStep = "unlock"
	StepRemovalAttempt   LifecycleStep = "removal-attempt"
	StepRemoval          LifecycleStep = "removal"
	StepBranch           LifecycleStep = "branch-removal"
	StepApplyLocked      LifecycleStep = "apply-locked"
	StepReceipt          LifecycleStep = "cleanup-receipt"
	StepTerminalReceipt  LifecycleStep = "terminal-receipt"
)

type Fault func(LifecycleStep) error

var creationLockAttempt = func(string) {}

type CleanupOptions struct {
	DiscardIgnored bool
	Full           bool
}

func hit(fault Fault, step LifecycleStep) error {
	if fault == nil {
		return nil
	}
	return fault(step)
}
func randomID() (string, error) {
	raw := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", raw), nil
}
func lockReason(assignment intent.Assignment) string {
	return "bench owner=" + assignment.OwnerID +
		" assignment=" + assignment.ID +
		" request=" + assignment.Request +
		" label=" + textDigest(assignment.Label) +
		" start=" + assignment.Start
}
func relock(root string, assignment intent.Assignment, fault Fault) error {
	if out, err := exec.Command("git", "-C", root, "worktree", "lock", "--reason", lockReason(assignment), assignment.Worktree).CombinedOutput(); err != nil {
		return fmt.Errorf("HIGH-SEVERITY residual: restore exact Bench lock: %s", strings.TrimSpace(string(out)))
	}
	if err := validateCreationBundle(root, assignment); err != nil {
		return fmt.Errorf("HIGH-SEVERITY residual: verify restored exact Bench lock: %w", err)
	}
	if err := hit(fault, StepRelock); err != nil {
		return fmt.Errorf("HIGH-SEVERITY residual: verify restored exact Bench lock: %w", err)
	}
	return nil
}
func lockCreationRequest(root, digest string) (func(), error) {
	address, err := intent.Address(root)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filepath.Join(filepath.Dir(address), "bench-create.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open creation request lock: %w", err)
	}
	creationLockAttempt(digest)
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("lock creation request: %w", err)
	}
	return func() { _ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN); _ = file.Close() }, nil
}

// Create makes one request-idempotent owned registration and persists its bundle.
func Create(root, request, label string, fault Fault, requestedStart ...string) (Creation, error) {
	if request == "" || label == "" {
		return Creation{}, errors.New("worktree create requires request and label")
	}
	root, err := canonicalPath(root)
	if err != nil {
		return Creation{}, err
	}
	digest := requestDigest(request)
	release, err := lockCreationRequest(root, digest)
	if err != nil {
		return Creation{}, err
	}
	defer release()
	if existing, ok, err := intent.FindAssignmentByRequest(root, digest); err != nil {
		return Creation{}, err
	} else if ok {
		if existing.Label != label || existing.State != intent.StateActive {
			return Creation{}, errors.New("worktree create request conflicts with its existing assignment")
		}
		if err := validateCreationBundle(root, existing); err != nil {
			return Creation{}, fmt.Errorf("worktree create request has incomplete prior state: %w", err)
		}
		return Creation{Path: existing.Worktree, Assignment: existing}, nil
	}
	ownerID, err := randomID()
	if err != nil {
		return Creation{}, fmt.Errorf("generate owner ID: %w", err)
	}
	assignmentID, err := randomID()
	if err != nil {
		return Creation{}, fmt.Errorf("generate assignment ID: %w", err)
	}
	startRef := "HEAD"
	if len(requestedStart) > 0 && requestedStart[0] != "" {
		startRef = requestedStart[0]
	} else if def, ok := git.ResolvedDefault(root); ok {
		startRef = def
	}
	start, err := git.Output("-C", root, "rev-parse", "--verify", startRef+"^{commit}")
	if err != nil {
		return Creation{}, fmt.Errorf("resolve assignment start: %w", err)
	}
	pool := Pool(root)
	if err := os.MkdirAll(pool, 0o700); err != nil {
		return Creation{}, fmt.Errorf("create worktree pool: %w", err)
	}
	_ = os.Chmod(pool, 0o700)
	path := filepath.Join(pool, ownerID+"-"+assignmentID)
	path, err = canonicalPath(path)
	if err != nil {
		return Creation{}, err
	}
	branch := intent.AssignmentBranchRef(ownerID, assignmentID)
	shortBranch := strings.TrimPrefix(branch, "refs/heads/")
	createdAt := time.Now().UTC().Format(time.RFC3339)
	assignment := intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: assignmentID, OwnerID: ownerID,
		Request: digest, Label: label, Start: start, Branch: branch, Worktree: path,
		State: intent.StateActive, Recovery: []intent.Recovery{}, CreatedAt: &createdAt,
	}
	args := []string{"-C", root, "worktree", "add", "-q", "--lock", "--reason", lockReason(assignment), "-b", shortBranch, path, start}
	if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		return Creation{}, fmt.Errorf("create locked worktree: %s", strings.TrimSpace(string(out)))
	}
	rollback := func(cause error) (Creation, error) {
		attributionErr := ensureRollbackAttribution(root, assignment)
		if out, unlockErr := exec.Command("git", "-C", root, "worktree", "unlock", path).CombinedOutput(); unlockErr != nil {
			return Creation{}, errors.Join(cause, attributionErr, fmt.Errorf("rollback unlock request-created registration: %s", strings.TrimSpace(string(out))))
		}
		removeErr := hit(fault, StepRollbackRemove)
		if removeErr == nil {
			if out, err := exec.Command("git", "-C", root, "worktree", "remove", "--force", path).CombinedOutput(); err != nil {
				removeErr = fmt.Errorf("rollback remove request-created registration: %s", strings.TrimSpace(string(out)))
			}
		}
		if removeErr != nil {
			return Creation{}, errors.Join(cause, attributionErr, removeErr, relock(root, assignment, fault))
		}
		if out, err := exec.Command("git", "-C", root, "update-ref", "-d", branch, start).CombinedOutput(); err != nil {
			return Creation{}, errors.Join(cause, attributionErr, fmt.Errorf("rollback delete exact request branch: %s", strings.TrimSpace(string(out))))
		}
		if err := intent.DeleteAssignment(root, assignmentID); err != nil {
			return Creation{}, errors.Join(cause, attributionErr, fmt.Errorf("rollback delete assignment attribution: %w", err))
		}
		return Creation{}, errors.Join(cause, attributionErr)
	}
	if err := hit(fault, StepRegistration); err != nil {
		return rollback(err)
	}
	path, err = canonicalPath(path)
	if err != nil {
		return rollback(err)
	}
	marker := Marker{Schema: OwnerMarkerSchema, OwnerID: ownerID, Path: path}
	if err := writeMarker(path, marker); err != nil {
		return rollback(fmt.Errorf("write owner marker: %w", err))
	}
	if err := hit(fault, StepMarker); err != nil {
		return rollback(err)
	}
	if err := intent.PutAssignment(root, assignment); err != nil {
		return rollback(fmt.Errorf("write assignment record: %w", err))
	}
	if err := hit(fault, StepRecord); err != nil {
		return rollback(err)
	}
	if err := validateCreationBundle(root, assignment); err != nil {
		return rollback(err)
	}
	return Creation{Path: path, Assignment: assignment}, nil
}
func ensureRollbackAttribution(root string, assignment intent.Assignment) error {
	marker := Marker{Schema: OwnerMarkerSchema, OwnerID: assignment.OwnerID, Path: assignment.Worktree}
	markerFile, markerErr := markerPath(assignment.Worktree)
	if markerErr == nil {
		if data, err := os.ReadFile(markerFile); errors.Is(err, os.ErrNotExist) {
			markerErr = writeMarker(assignment.Worktree, marker)
		} else if err != nil {
			markerErr = err
		} else if recorded, err := decodeMarker(data); err != nil || recorded != marker {
			markerErr = errors.New("rollback owner marker conflicts with request-created registration")
		}
	}
	recordErr := intent.PutAssignment(root, assignment)
	return errors.Join(markerErr, recordErr)
}
func markerPath(path string) (string, error) {
	admin, err := git.Output("-C", path, "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil || !filepath.IsAbs(admin) {
		return "", errors.New("resolve private worktree administration directory")
	}
	return filepath.Join(admin, OwnerMarkerFile), nil
}
func writeMarker(path string, marker Marker) error {
	markerFile, err := markerPath(path)
	if err != nil {
		return err
	}
	data, err := json.Marshal(marker)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	file, err := os.OpenFile(markerFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
func decodeMarker(data []byte) (Marker, error) {
	var marker Marker
	if err := jsonfile.Decode(data, &marker); err != nil {
		return Marker{}, err
	}
	return marker, nil
}
func validateCreationBundle(root string, assignment intent.Assignment) error {
	evidence, err := validateOwnerMarker(root, assignment.Worktree)
	if err != nil || evidence.marker.OwnerID != assignment.OwnerID {
		return errors.New("owner marker does not match assignment")
	}
	branch, err := git.Output("-C", assignment.Worktree, "symbolic-ref", "--quiet", "HEAD")
	if err != nil || branch != assignment.Branch {
		return errors.New("assignment branch is not checked out")
	}
	if evidence.registration.BranchRef == assignment.Branch && evidence.registration.LockReason == lockReason(assignment) {
		return nil
	}
	return errors.New("registration or Bench lock does not match assignment")
}

const releaseOperation = "worktree-release"

func releaseAssignment(root, requestArg, targetArg string) (intent.CleanupReceipt, error) {
	target, err := canonicalPath(targetArg)
	if err != nil {
		return intent.CleanupReceipt{}, err
	}
	repo, _, err := cleanupIdentity(root, target)
	if err != nil {
		return intent.CleanupReceipt{}, err
	}
	request := requestDigest(requestArg)
	if receipt, found, readErr := intent.CleanupReceiptFor(root, repo, releaseOperation, target, request); readErr != nil {
		return intent.CleanupReceipt{}, readErr
	} else if found && receipt.State == intent.ReceiptComplete {
		return receipt, nil
	}
	var assignment *intent.Assignment
	resumeFingerprint := ""
	if _, statErr := os.Lstat(target); errors.Is(statErr, os.ErrNotExist) {
		cleanup, found, readErr := intent.CleanupReceiptForRequest(root, repo, cleanupOperation, target, request)
		if readErr != nil {
			return intent.CleanupReceipt{}, readErr
		}
		expected := fingerprintParts([]byte("bench-automatic-registration/v1"), []byte(target), []byte(cleanup.Owner), []byte(cleanup.Assignment), []byte(request))
		matches := found && cleanup.Branch != "" && cleanup.Fingerprint == expected
		if matches && cleanup.State == intent.ReceiptComplete && cleanup.Action == string(ActionRemoved) && !git.OK("-C", root, "show-ref", "--verify", "--quiet", cleanup.Branch) {
			completed := intent.Assignment{ID: cleanup.Assignment, Worktree: target, State: intent.StateComplete}
			receipt := receiptFromRelease(repo, request, completed, string(ActionRemoved))
			return receipt, intent.PutCleanupReceipt(root, receipt)
		}
		if matches && cleanup.State == intent.ReceiptInFlight {
			current, currentErr := assignmentByID(root, cleanup.Assignment)
			if currentErr != nil || current.OwnerID != cleanup.Owner || current.Request != request || current.Worktree != target || current.Branch != cleanup.Branch || current.State != intent.StateCleanupPending {
				return intent.CleanupReceipt{}, errors.New("in-flight cleanup assignment does not match receipt")
			}
			assignment, resumeFingerprint = &current, cleanup.Fingerprint
		} else if found {
			return reconcileOutOfBand(root, repo, request, target, cleanup)
		}
	} else if statErr != nil {
		return intent.CleanupReceipt{}, statErr
	}
	if assignment == nil {
		foundAssignment, ok, findErr := intent.FindAssignmentByRequest(root, request)
		if findErr != nil {
			return intent.CleanupReceipt{}, findErr
		}
		if !ok || foundAssignment.Worktree != target {
			return intent.CleanupReceipt{}, errors.New("request, assignment, or path mismatch; checkout retained")
		}
		assignment = &foundAssignment
	}
	if assignment.State == intent.StateRecovered {
		receipt := receiptFromRelease(repo, request, *assignment, "retained")
		return receipt, intent.PutCleanupReceipt(root, receipt)
	}
	if assignment.State != intent.StateActive && assignment.State != intent.StateCleanupPending {
		return intent.CleanupReceipt{}, errors.New("assignment state does not accept release")
	}
	if resumeFingerprint == "" {
		if err := validateCreationBundle(root, *assignment); err != nil {
			return intent.CleanupReceipt{}, fmt.Errorf("%w; checkout retained", err)
		}
	}
	if assignment.State == intent.StateActive {
		assignment.State = intent.StateCleanupPending
		if err := intent.PutAssignment(root, *assignment); err != nil {
			return intent.CleanupReceipt{}, err
		}
	}
	terminal := func(plan CleanupPlan) error {
		current, readErr := assignmentByID(root, assignment.ID)
		if readErr != nil {
			return readErr
		}
		return intent.PutCleanupReceipt(root, receiptFromRelease(repo, request, current, string(plan.Action)))
	}
	var plan CleanupPlan
	if resumeFingerprint == "" {
		plan, err = applyAutomaticWithTerminal(root, target, nil, terminal)
	} else {
		planner := func(path string) (CleanupPlan, error) { return PlanAutomatic(root, path) }
		plan, err = applyCleanupTransaction(root, target, resumeFingerprint, planner, nil, terminal)
	}
	if err != nil {
		return intent.CleanupReceipt{}, err
	}
	if plan.Action == ActionRetain {
		return intent.CleanupReceipt{}, retainedReleaseError(plan, requestArg, targetArg)
	}
	receipt, found, readErr := intent.CleanupReceiptFor(root, repo, releaseOperation, target, request)
	if readErr != nil {
		return intent.CleanupReceipt{}, readErr
	}
	if !found {
		return intent.CleanupReceipt{}, errors.New("terminal receipt missing")
	}
	return receipt, nil
}

// retainedReleaseError turns a retain plan — the safe planner declining to remove the
// tree — into the verdict a session can act on: what blocked, and the exact command to
// re-run once it is cleared. It never points at `bench worktree clean --discard-ignored`,
// whose request-less form orphans the assignment (FT93b); the tree stays for the caller
// to resolve, then release again.
func retainedReleaseError(plan CleanupPlan, request, target string) error {
	reason := plan.Reason
	if reason == "" {
		reason = string(plan.ReasonCode)
	}
	return fmt.Errorf("worktree retained (%s): %s; resolve the retained state in %s, then re-run: bench worktree release --request %s %s",
		plan.ReasonCode, reason, target, request, target)
}

// residualAssignment reports whether a record preserves no work and is therefore safe
// to compact. Its recovery set is the single source of that judgment — an empty set
// means residue; a non-empty set means preserved work that must never be silently
// discarded. Both the release reconcile and the resume sweep consult this one predicate.
func residualAssignment(a intent.Assignment) bool { return len(a.Recovery) == 0 }

// reconcileOutOfBand resolves a release whose tree was removed out of band — a completed,
// owned cleanup receipt that does not match the automatic-registration reconcile shape
// (e.g. a request-less `bench worktree clean --discard-ignored --apply`). A residue
// record is compacted to a terminal release receipt; a record still holding preserved
// work is left intact and its recovery command named — the FT93(c) contract.
func reconcileOutOfBand(root, repo, request, target string, cleanup intent.CleanupReceipt) (intent.CleanupReceipt, error) {
	unauthorized := errors.New("cleanup receipt does not authorize release reconciliation")
	if cleanup.State != intent.ReceiptComplete || !cleanup.Owned || cleanup.Assignment == "" {
		return intent.CleanupReceipt{}, unauthorized
	}
	assignment, err := assignmentByID(root, cleanup.Assignment)
	if err != nil {
		// The record was already compacted out of band; synthesize the terminal receipt
		// so a replay of this release is idempotent.
		completed := intent.Assignment{ID: cleanup.Assignment, Worktree: target, State: intent.StateComplete}
		receipt := receiptFromRelease(repo, request, completed, string(ActionRemoved))
		return receipt, intent.PutCleanupReceipt(root, receipt)
	}
	if assignment.Request != request || assignment.Worktree != target {
		return intent.CleanupReceipt{}, unauthorized
	}
	if !residualAssignment(assignment) {
		return intent.CleanupReceipt{}, recoveryPendingError(assignment)
	}
	assignment.State = intent.StateComplete
	if err := intent.PutAssignment(root, assignment); err != nil {
		return intent.CleanupReceipt{}, err
	}
	receipt := receiptFromRelease(repo, request, assignment, string(ActionRemoved))
	return receipt, intent.PutCleanupReceipt(root, receipt)
}

// recoveryPendingError names the deliberate recover-or-retire path for a release that
// cannot compact because its tree was removed out of band while preserved work remains.
func recoveryPendingError(a intent.Assignment) error {
	ref := "(none)"
	if len(a.Recovery) > 0 {
		ref = a.Recovery[0].Ref
	}
	return fmt.Errorf("worktree removed out of band; its work is preserved — recover or retire it: bench worktree recovery %s", ref)
}
func renderRelease(stdout io.Writer, assignment intent.Assignment, action string) int {
	out, err := toon.Table("worktree_release", []string{"path", "assignment", "state", "action"}, [][]string{{assignment.Worktree, assignment.ID, string(assignment.State), action}})
	if err != nil {
		return 1
	}
	fmt.Fprint(stdout, out)
	return 0
}
