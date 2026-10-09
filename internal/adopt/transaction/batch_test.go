package transaction_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt/transaction"
)

// batchFixture holds several file destinations in one root, each with its own staged replacement.
type batchFixture struct {
	root    string
	dests   []string
	before  []os.FileInfo
	changes []transaction.Change
	store   transaction.Store
}

func newBatchFixture(t *testing.T, n int) batchFixture {
	t.Helper()
	root := t.TempDir()
	f := batchFixture{root: root, store: transaction.Store{Directory: filepath.Join(root, "records")}}
	for i := range n {
		dest, stage := filepath.Join(root, fmt.Sprintf("managed-%d", i)), filepath.Join(root, fmt.Sprintf("replacement-%d", i))
		if err := os.WriteFile(dest, preimage(i), 0o640); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(stage, replacement(i), 0o750); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(dest)
		if err != nil {
			t.Fatal(err)
		}
		f.dests, f.before = append(f.dests, dest), append(f.before, info)
		f.changes = append(f.changes, transaction.Change{Destination: dest, Stage: stage})
	}
	return f
}

func preimage(i int) []byte { return []byte(fmt.Sprintf("before %d", i)) }

func replacement(i int) []byte { return []byte(fmt.Sprintf("after %d", i)) }

// requirePreimages fails unless every destination holds its preimage content and mode.
// Each destination from index rewritten on must also keep its original file identity.
func (f batchFixture) requirePreimages(t *testing.T, rewritten int) {
	t.Helper()
	for i, dest := range f.dests {
		data, err := os.ReadFile(dest)
		info, statErr := os.Stat(dest)
		if err != nil || statErr != nil || string(data) != string(preimage(i)) || info.Mode().Perm() != 0o640 {
			t.Fatalf("%s holds %q, %v, %v; want its preimage", dest, data, err, statErr)
		}
		if i >= rewritten && !os.SameFile(f.before[i], info) {
			t.Fatalf("%s was replaced", dest)
		}
	}
}

func TestRecordSavesIndependentOfEntryCount(t *testing.T) {
	saves := func(n int) int {
		f := newBatchFixture(t, n)
		count := 0
		var last os.FileInfo
		f.store.SyncDirectory = func(dir string) error {
			if filepath.Dir(dir) == f.store.Directory {
				// A save publishes a new record file; a backup publication leaves the record as it was.
				info, err := os.Lstat(filepath.Join(dir, "record.json"))
				if err == nil && (last == nil || !os.SameFile(last, info)) {
					count, last = count+1, info
				}
			}
			return transaction.SyncDirectory(dir)
		}
		if _, err := f.store.Apply(f.changes); err != nil {
			t.Fatal(err)
		}
		return count
	}
	if one, many := saves(1), saves(5); one != many {
		t.Fatalf("record saves grow with entry count: 1 entry = %d, 5 entries = %d", one, many)
	}
}

func TestPublicationFailureAfterChangedDestinations(t *testing.T) {
	f := newBatchFixture(t, 5)
	const published = 2
	renames := 0
	// Fail every destination rename after the first published ones, so Apply's own undo cannot restore them either.
	f.store.Rename = func(from, to string) error {
		if filepath.Dir(to) == f.root {
			renames++
			if renames > published {
				return os.ErrPermission
			}
		}
		return os.Rename(from, to)
	}
	id, err := f.store.Apply(f.changes)
	if err == nil || id == "" || renames <= published {
		t.Fatalf("fixture did not fail publication after changed destinations: id=%q, renames=%d, err=%v", id, renames, err)
	}
	for i, dest := range f.dests[:published] {
		if data, err := os.ReadFile(dest); err != nil || string(data) != string(replacement(i)) {
			t.Fatalf("%s holds %q, %v; want its replacement before Undo", dest, data, err)
		}
	}
	if err := (transaction.Store{Directory: f.store.Directory}).Undo(id); err != nil {
		t.Fatalf("undo after failed publication: %v", err)
	}
	f.requirePreimages(t, published)
}

func TestStagingFailureLeavesDestinations(t *testing.T) {
	f := newBatchFixture(t, 4)
	late, stage := filepath.Join(f.root, "late", "managed"), filepath.Join(f.root, "late-replacement")
	if err := os.WriteFile(stage, []byte("late after"), 0o600); err != nil {
		t.Fatal(err)
	}
	f.changes = append(f.changes, transaction.Change{Destination: late, Stage: stage})
	f.store.SyncDirectory = func(dir string) error {
		// Staging the late entry creates its parent; refuse to persist that creation.
		if _, err := os.Lstat(filepath.Dir(late)); err == nil && dir == f.root {
			return os.ErrPermission
		}
		return transaction.SyncDirectory(dir)
	}
	if _, err := f.store.Apply(f.changes); err == nil {
		t.Fatal("apply succeeded after a staging failure")
	}
	f.requirePreimages(t, 0)
	if _, err := os.Lstat(late); !os.IsNotExist(err) {
		t.Fatalf("late destination changed: %v", err)
	}
	for _, dir := range []string{f.root, filepath.Dir(late)} {
		leaked, err := filepath.Glob(filepath.Join(dir, ".bench-adopt-*"))
		if err != nil || len(leaked) != 0 {
			t.Fatalf("staged temporaries left in %s: %v, %v", dir, leaked, err)
		}
	}
}
