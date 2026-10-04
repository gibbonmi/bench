package worktree

import (
	"errors"
	"fmt"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/poolkey"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// poolRootMode is the one mode a Bench-selected pool root may carry: the checkouts,
// leases, and markers under it are the operator's private working state.
const poolRootMode os.FileMode = 0o700

// The pool-root refusal details, one per refusing branch, so a test can tell which
// branch refused.
const (
	poolRootSymlink      = "worktree pool root is a symlink"
	poolRootNotDirectory = "worktree pool root is not a directory"
	poolRootModeUnset    = "worktree pool root mode cannot be set"
)

// preparePoolRoot creates the pool root when it is absent and sets it to poolRootMode on
// every call, so a root a crash left permissive is tightened on the next entry and an
// already tight root passes unchanged. A root that is a symlink or not a directory is
// refused without being followed. A mode Bench cannot set is a refusal, so the operation
// that needs the root stops before it puts a checkout in it.
func preparePoolRoot(j joins, pool string) error {
	info, err := os.Lstat(pool)
	if errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(pool, poolRootMode); err == nil {
			info, err = os.Lstat(pool)
		}
	}
	if err != nil {
		return err
	}
	wanted := fmt.Sprintf("directory at mode %04o", poolRootMode)
	if info.Mode()&os.ModeSymlink != 0 {
		return refusalError{refusal{detail: poolRootSymlink, observed: pool, wanted: wanted}}
	}
	if !info.IsDir() {
		return refusalError{refusal{detail: poolRootNotDirectory, observed: pool, wanted: wanted}}
	}
	if err := j.chmodPool(pool, poolRootMode); err != nil {
		return refusalError{refusal{detail: poolRootModeUnset + ": " + err.Error(), observed: pool, wanted: wanted}}
	}
	return nil
}

// Create makes one locked, request-idempotent owned worktree and persists its ownership
// bundle. A repeated request returns the existing active assignment only when its label
// and bundle still match; requestedStart selects the branch start when it is present.
func Create(root, request, label string, fault Fault, requestedStart ...string) (Creation, error) {
	start := ""
	if len(requestedStart) > 0 {
		start = requestedStart[0]
	}
	return createAt(defaultJoins(), root, Home(), request, label, fault, currentTime(), func() (creationStart, error) { return creationStart{ref: start}, nil })
}

type creationStart struct {
	ref     string
	binding *intent.DeliveryBinding
}

// startResolver answers the start a creation branches from. It is a function, not a value,
// because a caller whose start costs a lookup pays for it only when the request is new.
type startResolver func() (creationStart, error)

// createAt is Create with the creation instant and the Bench home resolved explicitly at
// the caller's effect boundary. The start resolver runs after the request replay lookup
// misses, so a replay returns its existing record without resolving a start at all.
func createAt(j joins, root, home, request, label string, fault Fault, now time.Time, resolveStart ...startResolver) (Creation, error) {
	if request == "" || label == "" {
		return Creation{}, errors.New("worktree create requires request and label")
	}
	root, err := canonicalPath(root)
	if err != nil {
		return Creation{}, err
	}
	digest := intent.RequestDigest(request)
	release, err := lockCreationRequest(j, root, digest)
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
	origin := creationStart{}
	startRef := ""
	if len(resolveStart) > 0 && resolveStart[0] != nil {
		requested, err := resolveStart[0]()
		if err != nil {
			return Creation{}, err
		}
		origin, startRef = requested, requested.ref
	}
	if startRef == "" {
		startRef = "HEAD"
		if def, ok := git.ResolvedDefault(root); ok {
			startRef = def
		}
	}
	start, err := git.ResolveCommit(root, startRef)
	if err != nil {
		return Creation{}, fmt.Errorf("resolve assignment start: %w", err)
	}
	pool := poolAt(home, root)
	if err := preparePoolRoot(j, pool); err != nil {
		return Creation{}, fmt.Errorf("create worktree pool: %w", err)
	}
	path := filepath.Join(pool, poolkey.AssignmentSegment(ownerID, assignmentID))
	path, err = canonicalPath(path)
	if err != nil {
		return Creation{}, err
	}
	branch := intent.AssignmentBranchRef(ownerID, assignmentID)
	shortBranch := strings.TrimPrefix(branch, "refs/heads/")
	createdAt := now.UTC().Format(time.RFC3339)
	assignment := intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: assignmentID, OwnerID: ownerID,
		Request: digest, RequestToken: request, Label: label, Start: start, Branch: branch, Worktree: path,
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
		removeRecord := func() error { return intent.DeleteAssignment(root, assignmentID) }
		if origin.binding != nil {
			removeRecord = func() error { return (commitrepo.Store{Root: root}).RollbackSibling(assignment) }
		}
		if err := removeRecord(); err != nil {
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
	assignment.Worktree = path
	marker := Marker{Schema: OwnerMarkerSchema, OwnerID: ownerID, Path: path}
	if err := writeMarker(path, marker); err != nil {
		return rollback(fmt.Errorf("write owner marker: %w", err))
	}
	if err := hit(fault, StepMarker); err != nil {
		return rollback(err)
	}
	record := func() error { return intent.PutAssignment(root, assignment) }
	if origin.binding != nil {
		record = func() error { return (commitrepo.Store{Root: root}).RegisterSibling(assignment, *origin.binding) }
	}
	if err := record(); err != nil {
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
