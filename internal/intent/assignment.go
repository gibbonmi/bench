package intent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/intent/admissionpolicy"
)

// CleanupReceiptFor is Repo.CleanupReceiptFor with a nil run.
func CleanupReceiptFor(root, repo, operation, target, fingerprint string) (CleanupReceipt, bool, error) {
	return In(nil, root).CleanupReceiptFor(repo, operation, target, fingerprint)
}

func (r Repo) CleanupReceiptFor(repo, operation, target, fingerprint string) (CleanupReceipt, bool, error) {
	ledger, err := r.Read()
	if err != nil {
		return CleanupReceipt{}, false, err
	}
	for _, receipt := range ledger.CleanupReceipts {
		if receipt.Repo == repo && receipt.Operation == operation && receipt.Target == target && receipt.Fingerprint == fingerprint {
			return receipt, true, nil
		}
	}
	return CleanupReceipt{}, false, nil
}

// CleanupReceiptForRequest is Repo.CleanupReceiptForRequest with a nil run.
func CleanupReceiptForRequest(root, repo, operation, target, request string) (CleanupReceipt, bool, error) {
	return In(nil, root).CleanupReceiptForRequest(repo, operation, target, request)
}

func (r Repo) CleanupReceiptForRequest(repo, operation, target, request string) (CleanupReceipt, bool, error) {
	ledger, err := r.Read()
	if err != nil {
		return CleanupReceipt{}, false, err
	}
	var found *CleanupReceipt
	for i := range ledger.CleanupReceipts {
		receipt := &ledger.CleanupReceipts[i]
		if receipt.Repo != repo || receipt.Operation != operation || receipt.Target != target || receipt.Request != request || !receipt.Owned || receipt.Assignment == "" {
			continue
		}
		if found != nil {
			return CleanupReceipt{}, false, errors.New("cleanup receipt request is ambiguous")
		}
		found = receipt
	}
	if found == nil {
		return CleanupReceipt{}, false, nil
	}
	return *found, true, nil
}

// PutCleanupReceipt is Repo.PutCleanupReceipt with a nil run.
func PutCleanupReceipt(root string, receipt CleanupReceipt) error {
	return In(nil, root).PutCleanupReceipt(receipt)
}

// PutCleanupReceipt records one cleanup receipt and holds the completed receipts to
// their retention window. The rule validates the receipt, so this adapter does not.
func (r Repo) PutCleanupReceipt(receipt CleanupReceipt) error {
	return r.Transact(StrictRead, func(current Ledger) (Ledger, bool, error) {
		return admissionpolicy.PutCleanupReceipt(current, receipt)
	}, nil)
}

// LifecycleEvidence is Repo.LifecycleEvidence with a nil run.
func LifecycleEvidence(root string) ([]byte, error) { return In(nil, root).LifecycleEvidence() }

// LifecycleEvidence returns the exact persisted schema and assignment JSON values.
// Cleanup receipts deliberately live in the same ledger but outside this evidence:
// recording an apply transaction must not make its already-approved plan stale.
func (r Repo) LifecycleEvidence() ([]byte, error) {
	path, err := r.Address()
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		schema, _ := json.Marshal(Schema)
		return lifecycleEvidenceValues(schema, json.RawMessage("[]")), nil
	} else if err != nil {
		return nil, err
	}
	if _, err := readPath(path); err != nil {
		return nil, err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Schema      json.RawMessage `json:"schema"`
		Assignments json.RawMessage `json:"assignments"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if len(raw.Assignments) == 0 {
		raw.Assignments = json.RawMessage("[]")
	}
	return lifecycleEvidenceValues(raw.Schema, raw.Assignments), nil
}

func lifecycleEvidenceValues(schema, assignments []byte) []byte {
	return []byte(fmt.Sprintf("%d:%s%d:%s", len(schema), schema, len(assignments), assignments))
}

// FindAssignmentByRequest is Repo.FindAssignmentByRequest with a nil run.
func FindAssignmentByRequest(root, requestDigest string) (Assignment, bool, error) {
	return In(nil, root).FindAssignmentByRequest(requestDigest)
}

func (r Repo) FindAssignmentByRequest(requestDigest string) (Assignment, bool, error) {
	ledger, err := r.Read()
	if err != nil {
		return Assignment{}, false, err
	}
	for _, assignment := range ledger.Assignments {
		if assignment.Request == requestDigest {
			return assignment, true, nil
		}
	}
	return Assignment{}, false, nil
}

// FindAssignmentForRequest is Repo.FindAssignmentForRequest with a nil run.
func FindAssignmentForRequest(root, request string) (Assignment, bool, error) {
	return In(nil, root).FindAssignmentForRequest(request)
}

// FindAssignmentForRequest resolves an opaque caller token through the ledger's
// request-digest owner. Lifecycle commands never compare caller tokens directly.
func (r Repo) FindAssignmentForRequest(request string) (Assignment, bool, error) {
	return r.FindAssignmentByRequest(RequestDigest(request))
}

// AssignmentForWorktree is Repo.AssignmentForWorktree with a nil run, for the tree at path.
func AssignmentForWorktree(path string) (Assignment, bool) {
	return In(nil, path).AssignmentForWorktree()
}

// AssignmentForWorktree returns the active assignment that owns the tree at the root of r,
// and reports whether one does. The tree is matched by canonical path, because the
// ledger records a resolved path while a caller's root can arrive through a symlink.
// An unreadable ledger and an unresolvable path both answer no owner, which is the
// answer a path outside the pool gets: an unowned tree is a state to report, not a
// failure to raise.
//
// This is the one worktree-to-assignment lookup. Preflight names the owning
// assignment in its stale-base remedy, and `bench handoff` resolves the section it
// owns from it. A second lookup would let the two disagree about which assignment
// owns a tree.
func (r Repo) AssignmentForWorktree() (Assignment, bool) {
	assignments, err := r.Assignments()
	if err != nil {
		return Assignment{}, false
	}
	for _, a := range AssignmentsOwning(assignments, r.Root) {
		if a.State == StateActive {
			return a, true
		}
	}
	return Assignment{}, false
}

// Assignments is Repo.Assignments with a nil run.
func Assignments(root string) ([]Assignment, error) { return In(nil, root).Assignments() }

func (r Repo) Assignments() ([]Assignment, error) {
	ledger, err := r.Read()
	if err != nil {
		return nil, err
	}
	return append([]Assignment(nil), ledger.Assignments...), nil
}

// PutAssignment is Repo.PutAssignment with a nil run.
func PutAssignment(root string, assignment Assignment) error {
	return In(nil, root).PutAssignment(assignment)
}

// PutAssignment atomically inserts or updates one assignment ID. A request digest
// can name only one assignment, and an identical write preserves the ledger bytes.
func (r Repo) PutAssignment(assignment Assignment) error {
	if err := ValidateAssignment(assignment); err != nil {
		return err
	}
	return r.Transact(StrictRead, func(current Ledger) (Ledger, bool, error) {
		return admissionpolicy.PutAssignment(current, assignment)
	}, nil)
}

// ReauthorizeAssignment is Repo.ReauthorizeAssignment with a nil run.
func ReauthorizeAssignment(root, id, request string, verify func(Assignment) error, transition func(Assignment, Assignment) (func(), error), beforeCAS func(*Assignment)) (Assignment, error) {
	return In(nil, root).ReauthorizeAssignment(id, request, verify, transition, beforeCAS)
}

// ReauthorizeAssignment swaps one request digest only after its caller's reversible
// external transition succeeds. The rollback keeps that derived state coherent when
// the expected-old write cannot commit.
func (r Repo) ReauthorizeAssignment(id, request string, verify func(Assignment) error, transition func(Assignment, Assignment) (func(), error), beforeCAS func(*Assignment)) (Assignment, error) {
	return r.reauthorizeAssignment(durablefile.Native(), id, request, verify, transition, beforeCAS)
}

// reauthorizeAssignment takes the replacer so that a test can fail the ledger write after
// the real request transition.
func (r Repo) reauthorizeAssignment(replacer durablefile.Replacer, id, request string, verify func(Assignment) error, transition func(Assignment, Assignment) (func(), error), beforeCAS func(*Assignment)) (Assignment, error) {
	var previous Assignment
	var rollback func()
	err := r.TransactWithReplacer(replacer, StrictRead, func(current Ledger) (Ledger, bool, error) {
		index := -1
		for i := range current.Assignments {
			if current.Assignments[i].ID == id {
				index = i
				break
			}
		}
		if index < 0 {
			return current, false, errors.New("assignment not found")
		}
		stored := current.Assignments[index]
		if err := verify(stored); err != nil {
			return current, false, err
		}
		if !ValidIdentity(id) || request == "" {
			return current, false, errors.New("invalid reauthorization identity")
		}
		newDigest := RequestDigest(request)
		if err := admissionpolicy.OtherAssignmentOwnsRequest(current, id, newDigest); err != nil {
			return current, false, err
		}
		next := stored
		next.Request = newDigest
		next.RequestToken = request
		step, err := transition(stored, next)
		if err != nil {
			return current, false, err
		}
		rollback = step
		if beforeCAS != nil {
			beforeCAS(&current.Assignments[index])
		}
		swapped, changed, err := admissionpolicy.ReauthorizeAssignment(current, id, stored.Request, newDigest)
		if err != nil {
			// The compensation covers the write alone, so the compare-and-swap arm
			// runs the rollback itself.
			rollback()
			return current, false, err
		}
		previous = stored
		return swapped, changed, nil
	}, func() { rollback() })
	if err != nil {
		return Assignment{}, err
	}
	return previous, nil
}

// PurgeAssignments is Repo.PurgeAssignments with a nil run.
func PurgeAssignments(root string, keep func(Assignment) bool) (int, error) {
	return In(nil, root).PurgeAssignments(keep)
}

// PurgeAssignments drops every assignment record keep rejects, plus every record this
// build can no longer read at all, and reports how many it dropped. Read is strict
// because a record it cannot account for must not authorize anything. The purge is
// deliberately not. A ledger written by an older binary is unreadable exactly when the
// standing cleaner is the only thing able to clear it. Refusing there would leave every
// later command reading the same unreadable file. A purge that drops nothing leaves the
// file's bytes untouched, so a re-run over converged state is a no-op.
func (r Repo) PurgeAssignments(keep func(Assignment) bool) (int, error) {
	path, err := r.Address()
	if err != nil {
		return 0, err
	}
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	dropped := 0
	err = r.transact(durablefile.Native(), TolerantRead, func(current Ledger, records int) (Ledger, bool, error) {
		next, count, changed, err := admissionpolicy.PurgeAssignments(current, records, keep)
		dropped = count
		return next, changed, err
	}, nil)
	if err != nil {
		return 0, err
	}
	return dropped, nil
}

// DeleteAssignment is Repo.DeleteAssignment with a nil run.
func DeleteAssignment(root, id string) error { return In(nil, root).DeleteAssignment(id) }

func (r Repo) DeleteAssignment(id string) error {
	return r.Transact(StrictRead, func(current Ledger) (Ledger, bool, error) {
		return admissionpolicy.DeleteAssignment(current, id)
	}, nil)
}
