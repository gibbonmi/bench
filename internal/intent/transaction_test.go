package intent

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The refusal arm. A closure that returns an error must persist nothing, so the closure
// also reports a change: a transaction that writes before it reads the closure's error
// changes the bytes this test compares.
func TestLedgerTransactionRefusalPersistsNothingAndReleasesTheLock(t *testing.T) {
	root := newRepo(t)
	if err := Upsert(root, Entry{Key: "seed", Kind: KindShift, CreatedAt: time.Unix(1, 0).UTC()}); err != nil {
		t.Fatal(err)
	}
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	refusal := errors.New("the decision refuses")
	err = Transact(root, StrictRead, func(ledger Ledger) (Ledger, bool, error) {
		ledger.Entries = append(ledger.Entries, Entry{Key: "unwanted", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()})
		return ledger, true, refusal
	}, nil)
	if !errors.Is(err, refusal) {
		t.Fatalf("Transact error = %v, want %v", err, refusal)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("a refused transaction changed the bytes\nbefore=%s\nafter=%s", before, after)
	}
	if _, err := os.Stat(path + ".lock"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lock file after a refusal: %v", err)
	}
}

// The terminal write arm. The directory turns read-only inside the closure so that the
// failure lands on the write rather than on the lock the transaction already holds, and
// the compensation restores the directory so the release can still remove the lock.
func TestLedgerTransactionTerminalWriteFailureKeepsThePreviousBytes(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the read-only directory mode")
	}
	root := newRepo(t)
	if err := Upsert(root, Entry{Key: "seed", Kind: KindShift, CreatedAt: time.Unix(1, 0).UTC()}); err != nil {
		t.Fatal(err)
	}
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Dir(path)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	compensated := 0
	err = Transact(root, StrictRead, func(ledger Ledger) (Ledger, bool, error) {
		if err := os.Chmod(dir, 0o500); err != nil {
			return Ledger{}, false, err
		}
		ledger.Entries = append(ledger.Entries, Entry{Key: "unwritable", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()})
		return ledger, true, nil
	}, func() {
		compensated++
		_ = os.Chmod(dir, 0o700)
	})
	if err == nil {
		t.Fatal("Transact accepted a write into a read-only directory")
	}
	if compensated != 1 {
		t.Fatalf("compensation ran %d times, want 1", compensated)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("a failed write changed the bytes\nbefore=%s\nafter=%s", before, after)
	}
	if _, err := os.Stat(path + ".lock"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lock file after a failed write: %v", err)
	}
}
