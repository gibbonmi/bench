package intent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ReadMode names how a transaction reads the ledger before it hands it to the decision.
type ReadMode int

const (
	// StrictRead refuses a ledger it cannot account for, so no unaccounted record
	// authorizes anything.
	StrictRead ReadMode = iota
	// TolerantRead drops each assignment record it cannot decode or validate and keeps
	// the rest, so a ledger an older binary wrote still clears.
	TolerantRead
)

// Decision answers the next ledger and whether the transaction writes it. A decision
// reads no file and takes no lock: the transaction supplies the ledger it read under
// the lock, and it persists the answer.
type Decision func(Ledger) (Ledger, bool, error)

// Compensation rolls back an effect a decision applied outside the ledger. The
// transaction runs it exactly once, and only when the write fails.
type Compensation func()

// Transact applies one ledger-to-ledger decision under the ledger's lock. It resolves
// the address, acquires the lock, reads in the caller's mode, calls the decision, and
// replaces the file atomically on a reported change. It releases the lock on every
// exit, including every error exit. A nil compensation runs nothing, so only a caller
// with an external effect pays for one.
func Transact(root string, mode ReadMode, decide Decision, compensate Compensation) error {
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readMode(path, mode)
	if err != nil {
		return err
	}
	next, changed, err := decide(ledger)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	if err := writePath(path, next); err != nil {
		if compensate != nil {
			compensate()
		}
		return err
	}
	return nil
}

func readMode(path string, mode ReadMode) (Ledger, error) {
	if mode == TolerantRead {
		ledger, _, err := readPathTolerant(path)
		return ledger, err
	}
	return readPath(path)
}

// readPathTolerant is the tolerant mode's read. It decodes each assignment as a raw
// value and drops each record it cannot decode or validate, so one unreadable record
// does not take the whole ledger with it. It reads the entries and the cleanup receipts
// separately, and it reports how many assignment records the file held, which is the
// count a caller needs to say how many it dropped.
func readPathTolerant(path string) (Ledger, int, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Ledger{Schema: Schema, Entries: []Entry{}, Assignments: []Assignment{}}, 0, nil
	}
	if err != nil {
		return Ledger{}, 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	var stored struct {
		Schema      int               `json:"schema"`
		Assignments []json.RawMessage `json:"assignments"`
	}
	if err := json.Unmarshal(data, &stored); err != nil {
		return Ledger{}, 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	var rest struct {
		Entries         []Entry          `json:"entries"`
		CleanupReceipts []CleanupReceipt `json:"cleanup_receipts"`
	}
	if err := json.Unmarshal(data, &rest); err != nil {
		return Ledger{}, 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	kept := make([]Assignment, 0, len(stored.Assignments))
	// A record the legacy schema carries was never authorized to be there: the strict read
	// refuses the whole file over one. Under that schema every record is debris, whatever
	// it says about itself, so the loop that would judge them individually is skipped.
	if stored.Schema != LegacySchema {
		for _, record := range stored.Assignments {
			var assignment Assignment
			if json.Unmarshal(record, &assignment) != nil || ValidateAssignment(assignment) != nil {
				continue
			}
			kept = append(kept, assignment)
		}
	}
	ledger := Ledger{
		Schema:          stored.Schema,
		Entries:         rest.Entries,
		Assignments:     kept,
		CleanupReceipts: rest.CleanupReceipts,
	}
	return ledger, len(stored.Assignments), nil
}
