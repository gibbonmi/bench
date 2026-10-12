package intent

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// requireUnlocked fails t unless the lock of the ledger at path is released.
func requireUnlocked(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(lockOf(path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lock file after a failed write: %v", err)
	}
}

// D22: the ledger writer returns the injected directory-sync failure of its replacement
// after the rename, and the ledger holds the sorted encoding at mode 0600. A rename
// failure returns an unpublished failure and keeps the previous ledger bytes.
func TestIntentReplacement(t *testing.T) {
	ledger := Ledger{Entries: []Entry{
		{Key: "seed", Kind: KindShift, CreatedAt: time.Unix(1, 0).UTC()},
		{Key: "added", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()},
	}}

	t.Run("directory-sync after the rename", func(t *testing.T) {
		_, path := mutatorFixture(t)
		cause := errors.New("injected ledger directory-sync failure")
		rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

		err := writePathWith(rec.Replacer(), path, ledger)

		durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
		durablefiletest.RequireFile(t, path, 0o600, `{"schema":2,"entries":[`+addedJSON+`,`+seedEntryJSON+`]}`+"\n")
	})

	t.Run("rename", func(t *testing.T) {
		_, path := mutatorFixture(t)
		cause := errors.New("injected ledger rename failure")
		rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

		err := writePathWith(rec.Replacer(), path, ledger)

		durablefiletest.RequireFailure(t, err, "rename", false, cause)
		durablefiletest.RequireFile(t, path, 0o600, fixtureLedger+"\n")
	})
}

// D44: a ledger write that fails after its rename runs no compensation, keeps the new
// ledger, releases the lock, and returns the classified failure.
func TestTransactionPublishedFailure(t *testing.T) {
	root, path := mutatorFixture(t)
	cause := errors.New("injected ledger directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}
	compensated := 0

	err := TransactWithReplacer(rec.Replacer(), root, StrictRead, func(ledger Ledger) (Ledger, bool, error) {
		ledger.Entries = append(ledger.Entries, Entry{Key: "added", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()})
		return ledger, true, nil
	}, func() { compensated++ })

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	if compensated != 0 {
		t.Fatalf("compensation ran %d times after a completed rename, want 0", compensated)
	}
	durablefiletest.RequireFile(t, path, 0o600, `{"schema":2,"entries":[`+addedJSON+`,`+seedEntryJSON+`],"assignments":[`+firstJSON+`]}`+"\n")
	requireUnlocked(t, path)
}

// A failure the replacement did not classify follows the pre-rename contract: the
// compensation runs once and the ledger keeps its previous bytes. A creation instant
// past year 9999 has no JSON encoding, so the write fails before any replacement.
func TestTransactionUnclassifiedFailureCompensates(t *testing.T) {
	root, path := mutatorFixture(t)
	rec := &durablefiletest.Recorder{}
	compensated := 0

	err := TransactWithReplacer(rec.Replacer(), root, StrictRead, func(ledger Ledger) (Ledger, bool, error) {
		ledger.Entries = append(ledger.Entries, Entry{Key: "unencodable", Kind: KindWorktree, CreatedAt: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)})
		return ledger, true, nil
	}, func() { compensated++ })

	if err == nil || len(rec.Trace) != 0 {
		t.Fatalf("err = %v, trace = %v, want an encoding failure before any replacement stage", err, rec.Trace)
	}
	if compensated != 1 {
		t.Fatalf("compensation ran %d times, want 1", compensated)
	}
	durablefiletest.RequireFile(t, path, 0o600, fixtureLedger+"\n")
	requireUnlocked(t, path)
}

// D52: a reauthorization whose ledger write fails after its rename keeps the next request
// and the external effect its transition applied, and the effect matches the ledger.
func TestReauthorizePublishedFailure(t *testing.T) {
	root, path := mutatorFixture(t)
	witness := filepath.Join(t.TempDir(), "witness")
	if err := os.WriteFile(witness, []byte(activeAssignment().Request), 0o600); err != nil {
		t.Fatal(err)
	}
	rolledBack := 0
	transition := func(old, next Assignment) (func(), error) {
		if err := os.WriteFile(witness, []byte(next.Request), 0o600); err != nil {
			return nil, err
		}
		return func() {
			rolledBack++
			_ = os.WriteFile(witness, []byte(old.Request), 0o600)
		}, nil
	}
	cause := errors.New("injected ledger directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	_, err := reauthorizeAssignment(rec.Replacer(), root, activeAssignment().ID, "next-request",
		func(Assignment) error { return nil }, transition, nil)

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	if rolledBack != 0 {
		t.Fatalf("the rollback ran %d times after a completed rename, want 0", rolledBack)
	}
	durablefiletest.RequireFile(t, path, 0o600, `{"schema":2,"entries":[`+seedEntryJSON+`],"assignments":[`+swappedJSON+`]}`+"\n")
	stored, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	effect, err := os.ReadFile(witness)
	if err != nil {
		t.Fatal(err)
	}
	if string(effect) != stored.Assignments[0].Request {
		t.Fatalf("external effect = %q, want the ledger request %q", effect, stored.Assignments[0].Request)
	}
	requireUnlocked(t, path)
}
