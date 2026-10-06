package transaction_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt/transaction"
)

func TestUndoRefusesEquivalentReplacement(t *testing.T) {
	f := newRepairFixture(t)
	id, err := f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(f.dest)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(f.dest)
	if err != nil {
		t.Fatal(err)
	}
	replacement := filepath.Join(f.root, "later")
	if err := os.WriteFile(replacement, data, info.Mode()); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(replacement, f.dest); err != nil {
		t.Fatal(err)
	}
	later, err := os.Stat(f.dest)
	if err != nil || os.SameFile(info, later) {
		t.Fatalf("replacement fixture did not change identity: %v", err)
	}
	if err := (transaction.Store{Directory: f.store.Directory}).Undo(id); err == nil {
		t.Fatal("undo accepted a different inode with the exact postimage bytes and mode")
	}
	current, err := os.Stat(f.dest)
	if err != nil || !os.SameFile(later, current) {
		t.Fatalf("undo changed the replacement identity: %v", err)
	}
}

func TestUndoReportsEveryPreflightFailure(t *testing.T) {
	f := newRepairFixture(t)
	second := filepath.Join(f.root, "second")
	if err := os.WriteFile(second, []byte("second before"), 0o600); err != nil {
		t.Fatal(err)
	}
	paths := []string{f.dest, second}
	id, err := f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}, {Destination: second, Stage: f.stage}})
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if err := os.Rename(path, path+".aside"); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(path, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	err = f.store.Undo(id)
	if err == nil {
		t.Fatal("undo accepted non-file destinations")
	}
	for _, path := range paths {
		if !strings.Contains(err.Error(), path) {
			t.Errorf("unresolved target %s missing from %v", path, err)
		}
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(path+".aside", path); err != nil {
			t.Fatal(err)
		}
	}
	if err := (transaction.Store{Directory: f.store.Directory}).Undo(id); err != nil {
		t.Fatalf("preflight failure lost recoverable backups: %v", err)
	}
	for path, want := range map[string]string{f.dest: "before", second: "second before"} {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Fatalf("recovery at %s = %q, %v", path, got, err)
		}
	}
}

func TestIdentityAfterReplacementPreparation(t *testing.T) {
	for _, undo := range []bool{false, true} {
		name := "apply"
		if undo {
			name = "undo"
		}
		t.Run(name, func(t *testing.T) {
			f := newRepairFixture(t)
			var id string
			var err error
			if undo {
				id, err = f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}})
				if err != nil {
					t.Fatal(err)
				}
			}
			var replacement os.FileInfo
			f.store.SyncDirectory = func(dir string) error {
				staged, err := filepath.Glob(filepath.Join(f.root, ".bench-adopt-*"))
				if err != nil {
					return err
				}
				if len(staged) > 0 && dir != f.root && replacement == nil {
					later := filepath.Join(f.root, "later")
					if err := os.WriteFile(later, []byte("later user file"), 0o600); err != nil {
						return err
					}
					if err := os.Rename(later, f.dest); err != nil {
						return err
					}
					replacement, err = os.Stat(f.dest)
					if err != nil {
						return err
					}
				}
				return transaction.SyncDirectory(dir)
			}
			if undo {
				err = f.store.Undo(id)
			} else {
				_, err = f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}})
			}
			if replacement == nil {
				t.Fatal("fixture did not replace the destination after staging")
			}
			if err == nil {
				t.Fatal("publication accepted a destination replaced after staging")
			}
			current, err := os.Stat(f.dest)
			if err != nil || !os.SameFile(replacement, current) {
				t.Fatalf("publication overwrote the later identity: %v", err)
			}
		})
	}
}
