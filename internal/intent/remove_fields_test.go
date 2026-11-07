package intent

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestRemoveFieldsRefusesAResultThatDoesNotDecode proves the raw-field removal runs the
// strict decode on its result before the replace: a ledger whose field this schema does not
// know keeps its bytes when a removal elsewhere would write it, and the preserved copy holds
// the bytes that the removal read.
func TestRemoveFieldsRefusesAResultThatDoesNotDecode(t *testing.T) {
	root := newRepo(t)
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	stored := `{"schema":2,"entries":[` + seedEntryJSON + `],"added":1}` + "\n"
	if err := os.WriteFile(path, []byte(stored), 0o600); err != nil {
		t.Fatal(err)
	}
	preserve := filepath.Join(t.TempDir(), "restore", "copy.json")
	if _, err := RemoveFields(root, preserve, [][]string{{"entries"}}); err == nil || !strings.Contains(err.Error(), `"added"`) {
		t.Fatalf("removal = %v, want the strict decode refusal of the unknown field", err)
	}
	for _, file := range []string{path, preserve} {
		if data, err := os.ReadFile(file); err != nil || string(data) != stored {
			t.Fatalf("%s = %q, %v; want the stored bytes", file, data, err)
		}
	}
}

// TestRemoveFieldsTakesTheLedgerLock proves the raw-field removal takes the ledger lock: while
// another holder has the lock, the removal refuses when its lock window ends, and it writes
// neither the ledger nor the copy.
func TestRemoveFieldsTakesTheLedgerLock(t *testing.T) {
	t.Cleanup(setLockWindowForTest(50 * time.Millisecond))
	root, path := mutatorFixture(t)
	release, err := acquire(lockOf(path))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(release)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	preserve := filepath.Join(t.TempDir(), "copy.json")
	if copied, err := RemoveFields(root, preserve, [][]string{{"assignments"}}); err == nil || copied {
		t.Fatalf("removal under a held lock = (%t, %v), want the lock refusal and no copy", copied, err)
	}
	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(before, after) {
		t.Fatalf("ledger after the refused removal = %q, %v; want the stored bytes", after, err)
	}
	if _, err := os.Lstat(preserve); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("copy stat = %v, want no copy", err)
	}
}

// TestRemoveFieldsWritesNoCopyOfAnAbsentLedger proves that a removal over a ledger file that
// does not exist reports no copy and writes neither the copy nor the ledger.
func TestRemoveFieldsWritesNoCopyOfAnAbsentLedger(t *testing.T) {
	root := newRepo(t)
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	preserve := filepath.Join(t.TempDir(), "copy.json")
	if copied, err := RemoveFields(root, preserve, [][]string{{"entries"}}); err != nil || copied {
		t.Fatalf("removal of an absent ledger = (%t, %v), want no copy and no error", copied, err)
	}
	for _, file := range []string{path, preserve} {
		if _, err := os.Lstat(file); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("%s stat = %v, want no file", file, err)
		}
	}
}

// TestRemoveFieldsHoldsTheLedgerLockThroughItsWrite proves that a ledger writer that runs
// while the removal reads the ledger loses no entry. The ledger is a FIFO, so the removal
// stays in its read until the test feeds it. While it waits, a regular file with the same
// bytes replaces the FIFO path, so the writer reads a complete ledger.
func TestRemoveFieldsHoldsTheLedgerLockThroughItsWrite(t *testing.T) {
	t.Cleanup(setLockWindowForTest(50 * time.Millisecond))
	root, path := mutatorFixture(t)
	stored, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}
	preserve := filepath.Join(t.TempDir(), "copy.json")
	done, finished := make(chan error, 1), make(chan struct{})
	go func() {
		defer close(finished)
		_, err := RemoveFields(root, preserve, [][]string{{"assignments"}})
		done <- err
	}()
	// A failed assertion still joins the removal before the lock window is restored; the
	// feed's close, registered later, runs first and ends the removal's read.
	t.Cleanup(func() { <-finished })
	type opening struct {
		feed *os.File
		err  error
	}
	opened := make(chan opening, 1)
	go func() {
		// The open of the write end returns when the removal opens the read end.
		feed, err := os.OpenFile(path, os.O_WRONLY, 0)
		opened <- opening{feed, err}
	}()
	var feed *os.File
	select {
	case open := <-opened:
		if open.err != nil {
			t.Fatal(open.err)
		}
		feed = open.feed
		t.Cleanup(func() { _ = feed.Close() })
	case err := <-done:
		t.Fatalf("removal returned before it read the ledger: %v", err)
	}
	regular := path + ".regular"
	if err := os.WriteFile(regular, stored, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(regular, path); err != nil {
		t.Fatal(err)
	}
	// A writer that the lock refuses wrote no entry, so only a reported write must survive.
	written := []string{"seed"}
	during := Entry{Key: "during", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()}
	if err := Upsert(root, during); err == nil {
		written = append(written, during.Key)
	}
	if _, err := feed.Write(stored); err != nil {
		t.Fatal(err)
	}
	if err := feed.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatalf("removal = %v, want success", err)
	}
	after := Entry{Key: "after", Kind: KindWorktree, CreatedAt: time.Unix(3, 0).UTC()}
	if err := Upsert(root, after); err != nil {
		t.Fatal(err)
	}
	written = append(written, after.Key)
	ledger, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	var kept []string
	for _, entry := range ledger.Entries {
		kept = append(kept, entry.Key)
	}
	slices.Sort(kept)
	slices.Sort(written)
	if !slices.Equal(kept, written) || len(ledger.Assignments) != 0 {
		t.Fatalf("ledger entries %v, assignments %d; want every written entry %v and no assignment", kept, len(ledger.Assignments), written)
	}
}
