package intent

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// D22: a directory failure after the rename fails the ledger write, and the ledger holds
// the sorted encoding at mode 0600.
func TestIntentReplacement(t *testing.T) {
	ledger := Ledger{Entries: []Entry{
		{Key: "seed", Kind: KindShift, CreatedAt: time.Unix(1, 0).UTC()},
		{Key: "added", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()},
	}}
	for _, stage := range durablefiletest.PostRenameStages {
		t.Run(stage, func(t *testing.T) {
			_, path := mutatorFixture(t)
			cause := errors.New("injected ledger " + stage + " failure")
			rec := &durablefiletest.Recorder{Fail: map[string]error{stage: cause}}

			err := writePathWith(rec.Replacer(), path, ledger)

			durablefiletest.RequireFailure(t, err, stage, true, cause)
			durablefiletest.RequireFile(t, path, 0o600, `{"schema":2,"entries":[`+addedJSON+`,`+seedEntryJSON+`]}`+"\n")
		})
	}
}

// D36: a failure before the rename refuses the ledger write and leaves the prior ledger.
func TestIntentPreRenameFailure(t *testing.T) {
	_, path := mutatorFixture(t)
	prior, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cause := errors.New("injected ledger rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	err = writePathWith(rec.Replacer(), path, Ledger{Entries: []Entry{{Key: "added", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()}}})

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	durablefiletest.RequireBytes(t, path, string(prior))
}

// D44 and D45: the new ledger may be visible after a failure past the rename, so the
// transaction keeps the decision's effects. A rename failure leaves the prior ledger, so
// the compensation runs once while the lock still holds. Either way the transaction
// releases the lock and returns the classified failure.
func TestTransactionPublishedFailure(t *testing.T) {
	type row struct {
		stage     string
		published bool
	}
	cases := []row{{stage: "rename"}}
	for _, stage := range durablefiletest.PostRenameStages {
		cases = append(cases, row{stage: stage, published: true})
	}
	for _, testCase := range cases {
		t.Run(testCase.stage, func(t *testing.T) {
			root, path := mutatorFixture(t)
			prior, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			cause := errors.New("injected ledger " + testCase.stage + " failure")
			rec := &durablefiletest.Recorder{Fail: map[string]error{testCase.stage: cause}}
			compensated := 0

			err = TransactWithReplacer(rec.Replacer(), root, StrictRead, func(ledger Ledger) (Ledger, bool, error) {
				ledger.Entries = append(ledger.Entries, Entry{Key: "added", Kind: KindWorktree, CreatedAt: time.Unix(2, 0).UTC()})
				return ledger, true, nil
			}, func() {
				compensated++
				if _, err := os.Stat(lockOf(path)); err != nil {
					t.Errorf("lock file during the compensation: %v", err)
				}
			})

			durablefiletest.RequireFailure(t, err, testCase.stage, testCase.published, cause)
			if testCase.published {
				if compensated != 0 {
					t.Fatalf("compensation ran %d times after a completed rename, want 0", compensated)
				}
				durablefiletest.RequireBytes(t, path, `{"schema":2,"entries":[`+addedJSON+`,`+seedEntryJSON+`],"assignments":[`+firstJSON+`]}`+"\n")
			} else {
				if compensated != 1 {
					t.Fatalf("compensation ran %d times before a completed rename, want 1", compensated)
				}
				durablefiletest.RequireBytes(t, path, string(prior))
			}
			if _, err := os.Stat(lockOf(path)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("lock file after a replacement failure: %v", err)
			}
		})
	}
}

// D30: a ledger the caller cannot encode refuses before the replacement starts and
// leaves the ledger as it was, so the compensation runs once. A creation instant past
// year 9999 has no JSON encoding.
func TestTransactionUnclassifiedFailureCompensates(t *testing.T) {
	root, path := mutatorFixture(t)
	prior, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	rec := &durablefiletest.Recorder{}
	compensated := 0

	err = TransactWithReplacer(rec.Replacer(), root, StrictRead, func(ledger Ledger) (Ledger, bool, error) {
		ledger.Entries = append(ledger.Entries, Entry{Key: "unencodable", Kind: KindWorktree, CreatedAt: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)})
		return ledger, true, nil
	}, func() { compensated++ })

	durablefiletest.RequireRefusal(t, err, rec.Trace)
	if compensated != 1 {
		t.Fatalf("compensation ran %d times, want 1", compensated)
	}
	durablefiletest.RequireBytes(t, path, string(prior))
}

// D52: the swapped request may be visible after a failed directory sync, so the
// reauthorization keeps the external effect its transition applied, and that effect
// matches the ledger.
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
	durablefiletest.RequireBytes(t, path, `{"schema":2,"entries":[`+seedEntryJSON+`],"assignments":[`+swappedJSON+`]}`+"\n")
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
}
