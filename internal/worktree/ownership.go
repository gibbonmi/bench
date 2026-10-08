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
	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
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

type LifecycleStep string

const (
	StepRegistration     LifecycleStep = "registration"
	StepMarker           LifecycleStep = "marker"
	StepRecord           LifecycleStep = "record"
	StepRollbackRemove   LifecycleStep = "rollback-remove"
	StepRelock           LifecycleStep = "relock"
	StepRecoveryMetadata LifecycleStep = "recovery-metadata"
	StepRecoveryRef      LifecycleStep = "recovery-ref"
	// StepRecoveryRowClose names the boundary a retire or a discard crosses once the
	// recovery ref is deleted and only the assignment row is left to close.
	StepRecoveryRowClose LifecycleStep = "recovery-row-close"
	// StepLifecycleSweep names the boundary the standing cleaner crosses once per debris
	// ref. An interruption can then be placed between two deletions rather than only
	// around the whole sweep.
	StepLifecycleSweep  LifecycleStep = "lifecycle-sweep"
	StepUnlock          LifecycleStep = "unlock"
	StepRemovalAttempt  LifecycleStep = "removal-attempt"
	StepRemoval         LifecycleStep = "removal"
	StepBranch          LifecycleStep = "branch-removal"
	StepApplyLocked     LifecycleStep = "apply-locked"
	StepReceipt         LifecycleStep = "cleanup-receipt"
	StepTerminalReceipt LifecycleStep = "terminal-receipt"
)

type Fault func(LifecycleStep) error

// CleanupOptions carries the operator's invocation choices into every plan and the apply
// that must reproduce it. DiscardBranch is an assertion, not a force: it supplies the
// landedness proof git.LandedInDefault refuses to derive from an ambiguous shape. It
// authorizes nothing else — every ownership, identity, and path-safety refusal is decided
// without reading it.
type CleanupOptions struct {
	DiscardIgnored bool
	DiscardBranch  bool
	Full           bool
	Unclaimed      bool
}

func hit(fault Fault, step LifecycleStep) error {
	if fault == nil {
		return nil
	}
	return fault(step)
}

// deletePlannedBranch deletes the landed branch a cleanup plan names, at the commit the plan
// read, once the plan's checkout is retired. A plan that proves no landed branch keeps it.
func deletePlannedBranch(root string, plan CleanupPlan, checkpoint func(string) error, fault Fault) error {
	if !plan.deleteBranch {
		return nil
	}
	if err := git.DeleteBranchExact(root, plan.branchRef, plan.branchOID); err != nil {
		return fmt.Errorf("delete exact landed branch: %w", err)
	}
	if err := checkpoint(intent.ReceiptPhaseBranch); err != nil {
		return err
	}
	return hit(fault, StepBranch)
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
func lockCreationRequest(j joins, root, digest string) (func(), error) {
	address, err := intent.Address(root)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filepath.Join(filepath.Dir(address), "bench-create.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open creation request lock: %w", err)
	}
	j.creationLockAttempt(digest)
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("lock creation request: %w", err)
	}
	return func() { _ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN); _ = file.Close() }, nil
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
	admin, err := git.AdminDir(path)
	if err != nil {
		return "", errors.New("resolve checkout administration directory")
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

// validateCreationBundle names the one identity component the assignment's own records
// fail. The checked-out branch is the bundle's own predicate, not a registry component,
// so it keeps its own sentence between the marker read and the recorded evidence.
func validateCreationBundle(root string, assignment intent.Assignment) error {
	// The marker step runs before the branch check, because the branch check decides
	// nothing on evidence the marker step has not accepted. The caller already holds the
	// assignment, so the state component stays the caller's own check.
	evidence, err := ownerMarkerRefusal(root, assignment.Worktree, assignment)
	if err != nil {
		return err
	}
	if _, ok := assignmentBranchCheckedOut(assignment); !ok {
		return errors.New("assignment branch is not checked out")
	}
	return registrationRefusal(evidence, assignment)
}

const releaseOperation = "worktree-release"

func releaseAssignment(j joins, a ambient, root, requestArg, targetArg string) (intent.CleanupReceipt, error) {
	target, err := canonicalPath(targetArg)
	if err != nil {
		return intent.CleanupReceipt{}, err
	}
	repo, _, err := cleanupIdentity(root, target)
	if err != nil {
		return intent.CleanupReceipt{}, err
	}
	request := intent.RequestDigest(requestArg)
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
			receipt := receiptFromRelease(repo, request, completed, string(ActionRemoved), cleanup.Branch, cleanup.BranchOID)
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
		foundAssignment, findErr := assignmentForRequest(root, requestArg, assignmentRecoveryContext{
			target: target,
			suffix: retainedSuffix,
		})
		if findErr != nil {
			return intent.CleanupReceipt{}, findErr
		}
		if foundAssignment.Worktree != target {
			retained := componentRefusal(componentAssignmentPath, foundAssignment.ID, foundAssignment.Worktree, target)
			retained.detail += retainedSuffix
			return intent.CleanupReceipt{}, retained
		}
		assignment = &foundAssignment
	}
	if assignment.State == intent.StateRecovered {
		receipt := receiptFromRelease(repo, request, *assignment, "retained", "", "")
		return receipt, intent.PutCleanupReceipt(root, receipt)
	}
	if assignment.State != intent.StateActive && assignment.State != intent.StateCleanupPending {
		return intent.CleanupReceipt{}, errors.New("assignment state does not accept release")
	}
	// A tree that is gone answers for no lease and no bundle, so the release reads its lease
	// through the shared administration directory and releases what the tree left behind.
	absent := false
	if shape, shapeErr := ClassifyPathShape(target); resumeFingerprint == "" && shapeErr == nil && shape == ShapeAbsent {
		absent = true
		if err := missingTreeLeaseRefusal(root, targetArg, *assignment); err != nil {
			return intent.CleanupReceipt{}, err
		}
	} else if resumeFingerprint == "" {
		lease, leaseErr := LeaseFile(target)
		if leaseErr != nil {
			return intent.CleanupReceipt{}, leaseErr
		}
		if _, statErr := os.Lstat(lease); statErr == nil && ProbeLease(lease) == LeaseUnknown {
			return intent.CleanupReceipt{}, retainedReleaseError(retainedPlan(target, ReasonUncertain, unknownLeaseReason), targetArg, assignment.ID)
		}
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
		return intent.PutCleanupReceipt(root, receiptFromRelease(repo, request, current, string(plan.Action), plan.branchRef, plan.branchOID))
	}
	var plan CleanupPlan
	if absent {
		planner := func(string) (CleanupPlan, error) { return missingTreeRelease(*assignment), nil }
		plan, err = applyCleanupTransaction(j, a, root, target, missingTreeRelease(*assignment).Fingerprint, planner, nil, terminal)
	} else if resumeFingerprint == "" {
		plan, err = applyAutomaticWithTerminal(j, a, root, target, nil, terminal)
	} else {
		planner := func(path string) (CleanupPlan, error) { return planAutomaticAt(j, root, path, a.now) }
		plan, err = applyCleanupTransaction(j, a, root, target, resumeFingerprint, planner, nil, terminal)
	}
	if err != nil {
		return intent.CleanupReceipt{}, err
	}
	if plan.Action == ActionRetain {
		return intent.CleanupReceipt{}, retainedReleaseError(plan, targetArg, assignment.ID)
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

// missingTreeLeaseRefusal keeps the lease verdicts of a present tree's release for a tree
// that is gone: an unknown lease and a live lease each retain the assignment. The policy
// decides the live-lease verdict, so its reason has one source.
func missingTreeLeaseRefusal(root, targetArg string, assignment intent.Assignment) error {
	leases, err := assignmentLeaseStates(root)
	if err != nil {
		return err
	}
	switch LeaseState(leases[assignment.OwnerID]) {
	case LeaseUnknown:
		return retainedReleaseError(retainedPlan(assignment.Worktree, ReasonUncertain, unknownLeaseReason), targetArg, assignment.ID)
	case LeaseLive:
		live := lifecyclepolicy.DecideAutomatic(lifecyclepolicy.AutomaticFacts{Explicit: lifecyclepolicy.ExplicitOutcome{HasAssignment: true}, LiveLease: true})
		return retainedReleaseError(retainedPlan(assignment.Worktree, live.ReasonCode, live.Reason), targetArg, assignment.ID)
	}
	return nil
}

// retainedReleaseError turns a retain plan — the safe planner declining to remove the
// tree — into the verdict a session can act on. That verdict names what blocked, and a
// route that can succeed. It never points at
// `bench worktree clean --discard-ignored`, whose request-less form orphans the
// assignment. The tree stays for the caller to resolve; releaseNext owns the route.
func retainedReleaseError(plan CleanupPlan, target, assignment string) error {
	reason := plan.Reason
	if reason == "" {
		reason = string(plan.ReasonCode)
	}
	return refusalError{refusal{detail: fmt.Sprintf("worktree retained (%s): %s", plan.ReasonCode, reason), next: releaseNext(plan.ReasonCode, target, assignment), paths: plan.Ignored.Paths}}
}

// residualAssignment is the policy preservation-residue decision;
// internal/worktree/lifecyclepolicy owns its judgment.
func residualAssignment(a intent.Assignment) bool { return lifecyclepolicy.Residual(a) }

// reconcileOutOfBand resolves a release whose tree was removed out of band. Its input is
// a completed, owned cleanup receipt that does not match the automatic-registration
// reconcile shape (e.g. a request-less `bench worktree clean --discard-ignored --apply`).
// A residue record is compacted to a terminal release receipt. A record still holding
// preserved work is left intact, and its recovery command is named.
func reconcileOutOfBand(root, repo, request, target string, cleanup intent.CleanupReceipt) (intent.CleanupReceipt, error) {
	unauthorized := errors.New("cleanup receipt does not authorize release reconciliation")
	if cleanup.State != intent.ReceiptComplete || !cleanup.Owned || cleanup.Assignment == "" {
		return intent.CleanupReceipt{}, unauthorized
	}
	assignment, err := assignmentByID(root, cleanup.Assignment)
	if err != nil {
		// The record was already compacted out of band. Synthesize the terminal receipt
		// so a replay of this release is idempotent.
		completed := intent.Assignment{ID: cleanup.Assignment, Worktree: target, State: intent.StateComplete}
		receipt := receiptFromRelease(repo, request, completed, string(ActionRemoved), cleanup.Branch, cleanup.BranchOID)
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
	receipt := receiptFromRelease(repo, request, assignment, string(ActionRemoved), cleanup.Branch, cleanup.BranchOID)
	return receipt, intent.PutCleanupReceipt(root, receipt)
}

// recoveryPendingError names what a release cannot compact: its tree was removed out of
// band while the record still points at preserved work. No command retires that ref any
// more. The standing cleaner sweeps the namespace at the next session start. The line
// hands over the ref itself instead — the only handle left for reading the work back.
func recoveryPendingError(a intent.Assignment) error {
	ref := "(none)"
	if len(a.Recovery) > 0 {
		ref = a.Recovery[0].Ref
	}
	return fmt.Errorf("worktree removed out of band; its work is preserved until the next session start sweeps it: git show %s", ref)
}
func renderRelease(stdout io.Writer, assignment intent.Assignment, action string) int {
	out, err := toon.Table(releaseTable, []string{"path", "assignment", "state", "action"}, [][]string{{assignment.Worktree, assignment.ID, string(assignment.State), action}})
	if err != nil {
		return 1
	}
	fmt.Fprint(stdout, out)
	return 0
}
