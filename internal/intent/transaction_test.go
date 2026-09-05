package intent

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/intent/admissionpolicy"
)

// mutatorFixture writes one deterministic ledger at the address root resolves, so every
// mutator case starts from the same bytes. It writes the file directly, because a seed
// built through the mutators would grade them against themselves.
func mutatorFixture(t *testing.T) (string, string) {
	t.Helper()
	root := newRepo(t)
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	seed := Ledger{
		Schema:      Schema,
		Entries:     []Entry{{Key: "seed", Kind: KindShift, CreatedAt: time.Unix(1, 0).UTC()}},
		Assignments: []Assignment{activeAssignment()},
	}
	if err := writePath(path, seed); err != nil {
		t.Fatal(err)
	}
	return root, path
}

func fixtureReceipt() CleanupReceipt {
	return CleanupReceipt{
		Schema: CleanupReceiptSchema, Repo: "/repo", Operation: "worktree-clean",
		Target: "/repo/target", Fingerprint: strings.Repeat("1", 64),
		State: ReceiptComplete, Phase: ReceiptPhaseTerminal, Action: "removed",
		Tracked: "clean", Ignored: "count=0 bytes=0 shown=0 truncated=false", Recovery: "none",
	}
}

func secondAssignment() Assignment {
	owner, id := strings.Repeat("a", 32), strings.Repeat("e", 32)
	second := activeAssignment()
	second.ID, second.OwnerID = id, owner
	second.Request = strings.Repeat("f", 64)
	second.Branch = AssignmentBranchRef(owner, id)
	second.Worktree = "/pool/second"
	return second
}

// The fixture's own JSON spellings. They are the operands this test writes, not values
// any mutator derives, so the expected ledgers below compose from them.
const (
	seedEntryJSON = `{"key":"seed","kind":"shift","created_at":"1970-01-01T00:00:01Z"}`
	addedJSON     = `{"key":"added","kind":"worktree","created_at":"1970-01-01T00:00:02Z"}`
	firstJSON     = `{"schema":"bench-assignment/v1","id":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","owner_id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","request":"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc","label":"delegate","start":"dddddddddddddddddddddddddddddddddddddddd","branch":"refs/heads/bench/assign/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","worktree":"/pool/delegate","state":"active","recovery":null}`
	swappedJSON   = `{"schema":"bench-assignment/v1","id":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","owner_id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","request":"80f3b5ddef868a2b4422c3203ef39bec9a3a4db048441b8c968ad92acfd966ab","label":"delegate","start":"dddddddddddddddddddddddddddddddddddddddd","branch":"refs/heads/bench/assign/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","worktree":"/pool/delegate","state":"active","recovery":null}`
	secondJSON    = `{"schema":"bench-assignment/v1","id":"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","owner_id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","request":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff","label":"delegate","start":"dddddddddddddddddddddddddddddddddddddddd","branch":"refs/heads/bench/assign/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","worktree":"/pool/second","state":"active","recovery":null}`
	receiptJSON   = `{"schema":"bench-cleanup-receipt/v1","repo":"/repo","operation":"worktree-clean","target":"/repo/target","fingerprint":"1111111111111111111111111111111111111111111111111111111111111111","state":"complete","phase":"terminal","action":"removed","tracked":"clean","ignored":"count=0 bytes=0 shown=0 truncated=false","recovery":"none","detail":""}`
	entriesOnly   = `{"schema":2,"entries":[` + seedEntryJSON + `]}`
	fixtureLedger = `{"schema":2,"entries":[` + seedEntryJSON + `],"assignments":[` + firstJSON + `]}`
)

// TestLedgerTransactionPreservesEveryMutatorOutcome runs each of the seven mutators over
// one fixture ledger and compares the resulting file bytes with the bytes the
// pre-transaction implementation wrote. The expected values were recorded by running this
// same test against the pre-migration tree. (Coverage row LS1.)
func TestLedgerTransactionPreservesEveryMutatorOutcome(t *testing.T) {
	cases := []struct {
		name string
		run  func(*testing.T, string) error
		want string
	}{
		{
			name: "upsert",
			run: func(_ *testing.T, root string) error {
				return Upsert(root, Entry{Key: "added", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()})
			},
			want: `{"schema":2,"entries":[` + addedJSON + `,` + seedEntryJSON + `],"assignments":[` + firstJSON + `]}`,
		},
		{
			name: "compact",
			run:  func(_ *testing.T, root string) error { return Compact(root) },
			want: fixtureLedger,
		},
		{
			name: "cleanup receipt",
			run:  func(_ *testing.T, root string) error { return PutCleanupReceipt(root, fixtureReceipt()) },
			want: `{"schema":2,"entries":[` + seedEntryJSON + `],"assignments":[` + firstJSON + `],"cleanup_receipts":[` + receiptJSON + `]}`,
		},
		{
			name: "assignment",
			run:  func(_ *testing.T, root string) error { return PutAssignment(root, secondAssignment()) },
			want: `{"schema":2,"entries":[` + seedEntryJSON + `],"assignments":[` + firstJSON + `,` + secondJSON + `]}`,
		},
		{
			name: "reauthorization",
			run: func(_ *testing.T, root string) error {
				_, err := ReauthorizeAssignment(root, activeAssignment().ID, "next-request",
					func(Assignment) error { return nil }, noReauthorizeTransition, nil)
				return err
			},
			want: `{"schema":2,"entries":[` + seedEntryJSON + `],"assignments":[` + swappedJSON + `]}`,
		},
		{
			name: "purge",
			run: func(t *testing.T, root string) error {
				dropped, err := PurgeAssignments(root, func(Assignment) bool { return false })
				if err == nil && dropped != 1 {
					t.Fatalf("PurgeAssignments dropped %d records, want 1", dropped)
				}
				return err
			},
			want: entriesOnly,
		},
		{
			name: "delete",
			run:  func(_ *testing.T, root string) error { return DeleteAssignment(root, activeAssignment().ID) },
			want: entriesOnly,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			root, path := mutatorFixture(t)
			if err := testCase.run(t, root); err != nil {
				t.Fatalf("the mutator answered an unexpected error: %v", err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != testCase.want+"\n" {
				t.Fatalf("the mutator wrote\n%s\nwant\n%s", got, testCase.want)
			}
		})
	}
}

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
		capability.Capability(t, capability.Privilege, "root ignores the read-only directory mode; cannot make the ledger directory unwritable")
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

// TestReauthorizeCompensatesOnTheTerminalWriteFailure pins the two compensation arms.
// The reauthorization supplies the rollback its transition returns, and every other
// mutator supplies none. The directory turns read-only under the lock, because a
// directory made read-only before the call fails at the lock rather than at the write.
// (Coverage rows LS7 and LS8.)
func TestReauthorizeCompensatesOnTheTerminalWriteFailure(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root ignores the read-only directory mode; cannot make the ledger directory unwritable")
	}

	t.Run("the reauthorization rollback runs exactly once", func(t *testing.T) {
		root, path := mutatorFixture(t)
		dir := filepath.Dir(path)
		t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		rolledBack := 0
		transition := func(Assignment, Assignment) (func(), error) {
			return func() {
				rolledBack++
				_ = os.Chmod(dir, 0o700)
			}, nil
		}
		// beforeCAS is the one hook that runs under the lock and before the write.
		lockDirectory := func(*Assignment) {
			if err := os.Chmod(dir, 0o500); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := ReauthorizeAssignment(root, activeAssignment().ID, "next-request",
			func(Assignment) error { return nil }, transition, lockDirectory); err == nil {
			t.Fatal("ReauthorizeAssignment accepted a write into a read-only directory")
		}
		if rolledBack != 1 {
			t.Fatalf("the rollback ran %d times, want 1", rolledBack)
		}
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(before, after) {
			t.Fatalf("the failed reauthorization changed the bytes\nbefore=%s\nafter=%s", before, after)
		}
	})

	t.Run("a mutator without a step answers the failure", func(t *testing.T) {
		root, path := mutatorFixture(t)
		dir := filepath.Dir(path)
		t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
		// PutAssignment supplies no step and offers no hook that runs under the lock, so
		// the nil step is observed through its own rule inside a transaction over the same
		// ledger. A transaction that required a step panics here rather than answer the
		// error. PutAssignment then answers the same read-only directory itself.
		err := Transact(root, StrictRead, func(current Ledger) (Ledger, bool, error) {
			if err := os.Chmod(dir, 0o500); err != nil {
				return current, false, err
			}
			return admissionpolicy.PutAssignment(current, secondAssignment())
		}, nil)
		if err == nil {
			t.Fatal("a nil-step transaction accepted a write into a read-only directory")
		}
		if err := PutAssignment(root, secondAssignment()); err == nil {
			t.Fatal("PutAssignment accepted a read-only ledger directory")
		}
		if err := os.Chmod(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		stored, err := Read(root)
		if err != nil || len(stored.Assignments) != 1 {
			t.Fatalf("the failed writes changed the ledger: %#v, %v", stored.Assignments, err)
		}
	})
}
