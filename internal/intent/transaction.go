package intent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gibbonmi/bench/internal/durablefile"
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
// transaction runs it exactly once, and only when the write fails before its rename
// completes. After a completed rename the new ledger may be visible, so the effect that
// matches it stays.
type Compensation func()

// Transact applies one ledger-to-ledger decision under the ledger's lock. It resolves
// the address, acquires the lock, reads in the caller's mode, calls the decision, and
// replaces the file durably on a reported change. It releases the lock on every
// exit, including every error exit. A nil compensation runs nothing, so only a caller
// with an external effect pays for one.
func Transact(root string, mode ReadMode, decide Decision, compensate Compensation) error {
	return TransactWithReplacer(durablefile.Native(), root, mode, decide, compensate)
}

// TransactWithReplacer is Transact with the ledger replacement run through replacer, so
// that a caller can fail one stage of the write under the real lock and decision.
func TransactWithReplacer(replacer durablefile.Replacer, root string, mode ReadMode, decide Decision, compensate Compensation) error {
	return transact(replacer, root, mode, func(ledger Ledger, _ int) (Ledger, bool, error) { return decide(ledger) }, compensate)
}

// countedDecision also receives how many assignment records the file held. Only the
// tolerant mode reports a count a decision cannot recover from the ledger it gets,
// because that mode drops the records it cannot read. The purge is the one decision
// that needs it, so the count stays off the exported decision.
type countedDecision func(Ledger, int) (Ledger, bool, error)

func transact(replacer durablefile.Replacer, root string, mode ReadMode, decide countedDecision, compensate Compensation) error {
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(lockOf(path))
	if err != nil {
		return err
	}
	defer release()
	ledger, records, err := readMode(path, mode)
	if err != nil {
		return err
	}
	next, changed, err := decide(ledger, records)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	if err := writePathWith(replacer, path, next); err != nil {
		if compensate != nil && !published(err) {
			compensate()
		}
		return err
	}
	return nil
}

// published reports that err is a replacement failure after its rename completed. Any
// other failure, including one the replacement did not classify, left the ledger as it was.
func published(err error) bool {
	var replaceErr *durablefile.Error
	return errors.As(err, &replaceErr) && replaceErr.Published
}

func readMode(path string, mode ReadMode) (Ledger, int, error) {
	if mode == TolerantRead {
		return readPathTolerant(path)
	}
	ledger, err := readPath(path)
	return ledger, len(ledger.Assignments), err
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
		Entries            []Entry             `json:"entries"`
		CleanupReceipts    []CleanupReceipt    `json:"cleanup_receipts"`
		CommitmentReceipts []CommitmentReceipt `json:"commitment_receipts"`
		MilestoneReceipts  []MilestoneReceipt  `json:"milestone_receipts"`
		Commitment         *CommitmentState    `json:"commitment,omitempty"`
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
		Schema:             stored.Schema,
		Entries:            rest.Entries,
		Assignments:        kept,
		CleanupReceipts:    rest.CleanupReceipts,
		CommitmentReceipts: rest.CommitmentReceipts,
		MilestoneReceipts:  rest.MilestoneReceipts,
		Commitment:         rest.Commitment,
	}
	return ledger, len(stored.Assignments), nil
}
