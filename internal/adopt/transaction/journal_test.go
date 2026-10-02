package transaction_test

import (
	"bytes"
	"github.com/gibbonmi/bench/internal/adopt/transaction"
	"os"
	"path/filepath"
	"testing"
)

func TestUndoRefusesInvalidRecord(t *testing.T) {
	for _, kind := range []string{"unknown state", "public record", "public backup", "public directory"} {
		t.Run(kind, func(t *testing.T) {
			f := newRepairFixture(t)
			id, err := f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}})
			if err != nil {
				t.Fatal(err)
			}
			dir := filepath.Join(f.store.Directory, id)
			path := filepath.Join(dir, "record.json")
			switch kind {
			case "unknown state":
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				changed := bytes.Replace(data, []byte(`"state":"applied"`), []byte(`"state":"unrecognized"`), 1)
				if bytes.Equal(data, changed) {
					t.Fatal("fixture did not corrupt the saved state")
				}
				if err := os.WriteFile(path, changed, 0o600); err != nil {
					t.Fatal(err)
				}
			case "public record":
				if err := os.Chmod(path, 0o644); err != nil {
					t.Fatal(err)
				}
			case "public backup":
				if err := os.Chmod(filepath.Join(dir, "0.before"), 0o644); err != nil {
					t.Fatal(err)
				}
			case "public directory":
				if err := os.Chmod(dir, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			before, err := os.ReadFile(f.dest)
			if err != nil {
				t.Fatal(err)
			}
			if err := f.store.Undo(id); err == nil {
				t.Fatal("undo accepted an invalid recovery record")
			}
			after, err := os.ReadFile(f.dest)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("invalid recovery record changed destination: %q, %v", after, err)
			}
		})
	}
}

func TestUndoRestoresCompleteMode(t *testing.T) {
	f := newRepairFixture(t)
	mode := os.ModeSetuid | os.ModeSetgid | 0o640
	if err := os.Chmod(f.dest, mode); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(f.dest)
	if err != nil {
		t.Fatal(err)
	}
	if before.Mode() != mode {
		t.Fatalf("fixture mode = %v, want %v", before.Mode(), mode)
	}
	id, err := f.store.Apply([]transaction.Change{{Destination: f.dest, Stage: f.stage}})
	if err != nil {
		t.Fatal(err)
	}
	if err := f.store.Undo(id); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(f.dest)
	if err != nil {
		t.Fatal(err)
	}
	if after.Mode() != before.Mode() {
		t.Fatalf("undo lost original file mode: got %v, want %v; %v", after.Mode(), before.Mode(), err)
	}
}
