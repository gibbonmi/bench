package admissionpolicy

import (
	"errors"
	"reflect"

	"github.com/gibbonmi/bench/internal/intent/ledger"
)

// PutCleanupReceipt replaces the receipt that carries the same repo, operation,
// target, and fingerprint, or appends a new one. It then holds the completed
// receipts to the retention window by dropping the oldest completed records over
// it. Every accepted receipt is a change, because a receipt records a fresh
// observation even when its fields repeat.
func PutCleanupReceipt(current ledger.Ledger, receipt ledger.CleanupReceipt) (ledger.Ledger, bool, error) {
	if err := ledger.ValidateCleanupReceipts([]ledger.CleanupReceipt{receipt}); err != nil {
		return current, false, err
	}
	replaced := false
	for i, stored := range current.CleanupReceipts {
		if stored.Repo == receipt.Repo && stored.Operation == receipt.Operation && stored.Target == receipt.Target && stored.Fingerprint == receipt.Fingerprint {
			current.CleanupReceipts[i], replaced = receipt, true
			break
		}
	}
	if !replaced {
		current.CleanupReceipts = append(current.CleanupReceipts, receipt)
	}
	completed := 0
	for _, stored := range current.CleanupReceipts {
		if stored.State == ledger.ReceiptComplete {
			completed++
		}
	}
	drop := completed - ledger.MaxCleanupReceipts
	if drop > 0 {
		next := current.CleanupReceipts[:0]
		for _, stored := range current.CleanupReceipts {
			if drop > 0 && stored.State == ledger.ReceiptComplete {
				drop--
				continue
			}
			next = append(next, stored)
		}
		current.CleanupReceipts = next
	}
	return current, true, nil
}

// PutAssignment inserts or updates one assignment ID. A request digest can name
// only one assignment, and an assignment identical to the stored one reports no
// change, so the caller preserves the ledger's exact bytes.
func PutAssignment(current ledger.Ledger, assignment ledger.Assignment) (ledger.Ledger, bool, error) {
	for i, stored := range current.Assignments {
		if stored.Request == assignment.Request && stored.ID != assignment.ID {
			return current, false, errors.New("assignment request already belongs to another assignment")
		}
		if stored.ID != assignment.ID {
			continue
		}
		if reflect.DeepEqual(stored, assignment) {
			return current, false, nil
		}
		current.Assignments[i] = assignment
		return current, true, nil
	}
	current.Assignments = append(current.Assignments, assignment)
	return current, true, nil
}

// CompareAndSwapRequestDigest replaces one assignment's request digest only while
// the stored digest is still the expected one. It is the pure half of the
// reauthorization: the caller's verification and its reversible external
// transition stay in the adapter that owns those effects.
func CompareAndSwapRequestDigest(assignment *ledger.Assignment, expectedOld, replacement string) error {
	if assignment.Request != expectedOld {
		return errors.New("assignment request changed during reauthorization")
	}
	assignment.Request = replacement
	return nil
}

// ReauthorizeAssignment applies the compare-and-swap to the assignment named id.
// The adapter runs its verification and its external transition before this rule,
// and it applies any pre-swap enrichment to the ledger it supplies here.
func ReauthorizeAssignment(current ledger.Ledger, id, expectedOld, replacement string) (ledger.Ledger, bool, error) {
	for i := range current.Assignments {
		if current.Assignments[i].ID != id {
			continue
		}
		if err := CompareAndSwapRequestDigest(&current.Assignments[i], expectedOld, replacement); err != nil {
			return current, false, err
		}
		return current, true, nil
	}
	return current, false, errors.New("assignment not found")
}

// OtherAssignmentOwnsRequest refuses a digest that already belongs to an
// assignment other than id. The reauthorization adapter asks this before it runs
// its external transition, so a refusal costs no rollback.
func OtherAssignmentOwnsRequest(current ledger.Ledger, id, digest string) error {
	for _, stored := range current.Assignments {
		if stored.ID != id && stored.Request == digest {
			return errors.New("request digest already belongs to another assignment")
		}
	}
	return nil
}

// PurgeAssignments keeps every record keep accepts under a first-wins identity,
// and it reports how many of the records read it dropped. The records count is
// what the file held, so a record this build could not read at all counts as
// dropped. A purge that drops nothing reports no change, so a re-run over
// converged state leaves the ledger's bytes untouched.
func PurgeAssignments(current ledger.Ledger, records int, keep func(ledger.Assignment) bool) (ledger.Ledger, int, bool, error) {
	kept := make([]ledger.Assignment, 0, len(current.Assignments))
	ids, requests := map[string]bool{}, map[string]bool{}
	for _, assignment := range current.Assignments {
		if ids[assignment.ID] || requests[assignment.Request] || !keep(assignment) {
			continue
		}
		ids[assignment.ID], requests[assignment.Request] = true, true
		kept = append(kept, assignment)
	}
	dropped := records - len(kept)
	if dropped == 0 {
		return current, 0, false, nil
	}
	current.Assignments = kept
	return current, dropped, true, nil
}

// DeleteAssignment removes the assignment named id. A ledger that holds no such
// assignment reports no change.
func DeleteAssignment(current ledger.Ledger, id string) (ledger.Ledger, bool, error) {
	next := make([]ledger.Assignment, 0, len(current.Assignments))
	for _, assignment := range current.Assignments {
		if assignment.ID != id {
			next = append(next, assignment)
		}
	}
	if len(next) == len(current.Assignments) {
		return current, false, nil
	}
	current.Assignments = next
	return current, true, nil
}
