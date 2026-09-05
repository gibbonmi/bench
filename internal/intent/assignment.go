package intent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

func CleanupReceiptFor(root, repo, operation, target, fingerprint string) (CleanupReceipt, bool, error) {
	ledger, err := Read(root)
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

func CleanupReceiptForRequest(root, repo, operation, target, request string) (CleanupReceipt, bool, error) {
	ledger, err := Read(root)
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

func PutCleanupReceipt(root string, receipt CleanupReceipt) error {
	if err := validateCleanupReceipts([]CleanupReceipt{receipt}); err != nil {
		return err
	}
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return err
	}
	replaced := false
	for i, current := range ledger.CleanupReceipts {
		if current.Repo == receipt.Repo && current.Operation == receipt.Operation && current.Target == receipt.Target && current.Fingerprint == receipt.Fingerprint {
			ledger.CleanupReceipts[i], replaced = receipt, true
			break
		}
	}
	if !replaced {
		ledger.CleanupReceipts = append(ledger.CleanupReceipts, receipt)
	}
	completed := 0
	for _, current := range ledger.CleanupReceipts {
		if current.State == ReceiptComplete {
			completed++
		}
	}
	drop := completed - MaxCleanupReceipts
	if drop > 0 {
		next := ledger.CleanupReceipts[:0]
		for _, current := range ledger.CleanupReceipts {
			if drop > 0 && current.State == ReceiptComplete {
				drop--
				continue
			}
			next = append(next, current)
		}
		ledger.CleanupReceipts = next
	}
	return writePath(path, ledger)
}

// LifecycleEvidence returns the exact persisted schema and assignment JSON values.
// Cleanup receipts deliberately live in the same ledger but outside this evidence:
// recording an apply transaction must not make its already-approved plan stale.
func LifecycleEvidence(root string) ([]byte, error) {
	path, err := Address(root)
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

func FindAssignmentByRequest(root, requestDigest string) (Assignment, bool, error) {
	ledger, err := Read(root)
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

// FindAssignmentForRequest resolves an opaque caller token through the ledger's
// request-digest owner. Lifecycle commands never compare caller tokens directly.
func FindAssignmentForRequest(root, request string) (Assignment, bool, error) {
	return FindAssignmentByRequest(root, RequestDigest(request))
}

// AssignmentForWorktree returns the active assignment that owns the tree at path,
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
func AssignmentForWorktree(path string) (Assignment, bool) {
	assignments, err := Assignments(path)
	if err != nil {
		return Assignment{}, false
	}
	for _, a := range AssignmentsOwning(assignments, path) {
		if a.State == StateActive {
			return a, true
		}
	}
	return Assignment{}, false
}

func Assignments(root string) ([]Assignment, error) {
	ledger, err := Read(root)
	if err != nil {
		return nil, err
	}
	return append([]Assignment(nil), ledger.Assignments...), nil
}

// PutAssignment atomically inserts or updates one assignment ID. A request digest
// can name only one assignment, and an identical write preserves the ledger bytes.
func PutAssignment(root string, assignment Assignment) error {
	if err := ValidateAssignment(assignment); err != nil {
		return err
	}
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return err
	}
	for i, current := range ledger.Assignments {
		if current.Request == assignment.Request && current.ID != assignment.ID {
			return errors.New("assignment request already belongs to another assignment")
		}
		if current.ID != assignment.ID {
			continue
		}
		if reflect.DeepEqual(current, assignment) {
			return nil
		}
		ledger.Assignments[i] = assignment
		return writePath(path, ledger)
	}
	ledger.Assignments = append(ledger.Assignments, assignment)
	return writePath(path, ledger)
}

// ReauthorizeAssignment swaps one request digest only after its caller's reversible
// external transition succeeds. The rollback keeps that derived state coherent when
// the expected-old write cannot commit.
func ReauthorizeAssignment(root, id, request string, verify func(Assignment) error, transition func(Assignment, Assignment) (func(), error), beforeCAS func(*Assignment)) (Assignment, error) {
	path, err := Address(root)
	if err != nil {
		return Assignment{}, err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return Assignment{}, err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return Assignment{}, err
	}
	for i := range ledger.Assignments {
		current := ledger.Assignments[i]
		if current.ID != id {
			continue
		}
		expectedOld := current.Request
		if err := verify(current); err != nil {
			return Assignment{}, err
		}
		if !ValidIdentity(id) || request == "" {
			return Assignment{}, errors.New("invalid reauthorization identity")
		}
		newDigest := RequestDigest(request)
		for j, other := range ledger.Assignments {
			if j != i && other.Request == newDigest {
				return Assignment{}, errors.New("request digest already belongs to another assignment")
			}
		}
		next := current
		next.Request = newDigest
		next.RequestToken = request
		rollback, err := transition(current, next)
		if err != nil {
			return Assignment{}, err
		}
		if beforeCAS != nil {
			beforeCAS(&ledger.Assignments[i])
		}
		if err := compareAndSwapRequestDigest(&ledger.Assignments[i], expectedOld, newDigest); err != nil {
			rollback()
			return Assignment{}, err
		}
		if err := writePath(path, ledger); err != nil {
			rollback()
			return Assignment{}, err
		}
		return current, nil
	}
	return Assignment{}, errors.New("assignment not found")
}

func compareAndSwapRequestDigest(assignment *Assignment, expectedOld, replacement string) error {
	if assignment.Request != expectedOld {
		return errors.New("assignment request changed during reauthorization")
	}
	assignment.Request = replacement
	return nil
}

// PurgeAssignments drops every assignment record keep rejects, plus every record this
// build can no longer read at all, and reports how many it dropped. Read is strict
// because a record it cannot account for must not authorize anything. The purge is
// deliberately not. A ledger written by an older binary is unreadable exactly when the
// standing cleaner is the only thing able to clear it. Refusing there would leave every
// later command reading the same unreadable file. A purge that drops nothing leaves the
// file's bytes untouched, so a re-run over converged state is a no-op.
func PurgeAssignments(root string, keep func(Assignment) bool) (int, error) {
	path, err := Address(root)
	if err != nil {
		return 0, err
	}
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return 0, err
	}
	defer release()
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	// One tolerant pass for the records the purge decides about, and one for everything
	// else in the file. Assignments are read as raw values so a single unreadable record
	// is dropped on its own rather than taking the whole ledger with it.
	var stored struct {
		Schema      int               `json:"schema"`
		Assignments []json.RawMessage `json:"assignments"`
	}
	if err := json.Unmarshal(data, &stored); err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	kept := make([]Assignment, 0, len(stored.Assignments))
	ids, requests := map[string]bool{}, map[string]bool{}
	// A record the legacy schema carries was never authorized to be there: Read refuses the
	// whole file over one. Under that schema every record is debris, whatever it says about
	// itself, so the loop that would judge them individually is skipped.
	if stored.Schema != LegacySchema {
		for _, record := range stored.Assignments {
			var assignment Assignment
			if json.Unmarshal(record, &assignment) != nil || ValidateAssignment(assignment) != nil {
				continue
			}
			if ids[assignment.ID] || requests[assignment.Request] || !keep(assignment) {
				continue
			}
			ids[assignment.ID], requests[assignment.Request] = true, true
			kept = append(kept, assignment)
		}
	}
	dropped := len(stored.Assignments) - len(kept)
	if dropped == 0 {
		return 0, nil
	}
	var rest struct {
		Entries         []Entry          `json:"entries"`
		CleanupReceipts []CleanupReceipt `json:"cleanup_receipts"`
	}
	if err := json.Unmarshal(data, &rest); err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	ledger := Ledger{Entries: rest.Entries, Assignments: kept, CleanupReceipts: rest.CleanupReceipts}
	return dropped, writePath(path, ledger)
}

func DeleteAssignment(root, id string) error {
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return err
	}
	next := ledger.Assignments[:0]
	for _, assignment := range ledger.Assignments {
		if assignment.ID != id {
			next = append(next, assignment)
		}
	}
	if len(next) == len(ledger.Assignments) {
		return nil
	}
	ledger.Assignments = next
	return writePath(path, ledger)
}
