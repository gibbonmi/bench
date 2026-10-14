package responsebound_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	rb "github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

// retainCall answers the call of a retained response with no read of its worktree.
func retainCall() (rb.Retained, error) { return rb.Retained{Verb: "test"}, nil }

// RT7, RT65: a response with no complete copy clears the earlier retained response. A spill
// file that failed midway holds only part of the output, and a spill that could not open
// passed the output through. A call that fails to read its worktree clears it too.
func TestRetainClearsAnIncompleteResponse(t *testing.T) {
	midway := func(path string) (io.WriteCloser, error) {
		file, err := rb.ExclusiveCreate(path)
		if err != nil {
			return nil, err
		}
		return rb.FaultAfter(file, 100), nil
	}
	refused := func(string) (io.WriteCloser, error) { return nil, rb.ErrInjected }
	failed := func() (rb.Retained, error) { return rb.Retained{}, rb.ErrInjected }
	for _, row := range []struct {
		name   string
		create func(string) (io.WriteCloser, error)
		call   func() (rb.Retained, error)
	}{
		{"spill-failed", midway, retainCall},
		{"create failure", refused, retainCall},
		{"failed call", rb.ExclusiveCreate, failed},
	} {
		t.Run(row.name, func(t *testing.T) {
			home := rb.PrivateHome(t)
			_, checkout, _ := responseboundtest.AssignmentCheckout(t, home)
			root := func() string { return checkout }
			var sink bytes.Buffer
			earlier := rb.New(home, &sink, &sink, root, false)
			rb.RespondWith(t, earlier, &sink, rb.StdoutLines("kept\n"))
			if err := earlier.Retain(home, checkout, retainCall); err != nil {
				t.Fatal(err)
			}
			if _, err := rb.ReadRetained(home, checkout); err != nil {
				t.Fatalf("earlier retained response: %v", err)
			}
			owner := rb.New(home, &sink, &sink, root, false)
			rb.SetCreate(owner, row.create)
			rb.RespondWith(t, owner, &sink, rb.StdoutLines(rb.Numbered(1, 25)...))
			_ = owner.Retain(home, checkout, row.call)
			if got, err := rb.ReadRetained(home, checkout); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("retained after the %s response = (%+v, %v), want none", row.name, got, err)
			}
		})
	}
}

// RT9: the release of a worktree removes its retained response and its staged excerpts.
func TestDropRemovesTheRetainedResponseAndStages(t *testing.T) {
	home := rb.PrivateHome(t)
	repo, checkout, id := responseboundtest.AssignmentCheckout(t, home)
	var sink bytes.Buffer
	owner := rb.New(home, &sink, &sink, func() string { return checkout }, false)
	rb.RespondWith(t, owner, &sink, rb.StdoutLines("kept\n"))
	if err := owner.Retain(home, checkout, retainCall); err != nil {
		t.Fatal(err)
	}
	if err := rb.Stage(home, checkout, "example", "1", "tests", rb.Staged{Mutant: 1}); err != nil {
		t.Fatal(err)
	}
	if files := responseboundtest.StoreFiles(t, home); len(files) != 2 {
		t.Fatalf("store files before the drop = %q, want the retained response and the stage", files)
	}
	if err := rb.Drop(home, repo, id); err != nil {
		t.Fatal(err)
	}
	if files := responseboundtest.StoreFiles(t, home); len(files) != 0 {
		t.Fatalf("store files after the drop = %q, want none", files)
	}
}

// A scope directory that is a symlink refuses the stage, and the link target receives no file.
func TestStageRefusesASymlinkedScope(t *testing.T) {
	home := rb.PrivateHome(t)
	_, checkout, _ := responseboundtest.AssignmentCheckout(t, home)
	if err := rb.Stage(home, checkout, "example", "1", "tests", rb.Staged{}); err != nil {
		t.Fatal(err)
	}
	files := responseboundtest.StoreFiles(t, home)
	if len(files) != 1 {
		t.Fatalf("store files = %q, want one stage", files)
	}
	stage := files[0]
	scope := filepath.Dir(stage)
	target := t.TempDir()
	if err := os.RemoveAll(scope); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, scope); err != nil {
		t.Fatal(err)
	}
	if err := rb.Stage(home, checkout, "example", "1", "tests", rb.Staged{}); err == nil {
		t.Fatal("stage through a symlinked scope = nil, want a refusal")
	}
	if entries, err := os.ReadDir(target); err != nil || len(entries) != 0 {
		t.Fatalf("link target entries = %v (%v), want none", entries, err)
	}
}

// A stage or a retained response that is not a regular file refuses its read, even when each
// store level is a real directory. A symlink to a valid copy refuses, and a FIFO whose writer
// sends a valid copy refuses too, so the read never follows a link or reads a pipe.
func TestReadRefusesAScopedFileThatIsNotRegular(t *testing.T) {
	for _, store := range scopedStores() {
		for _, kind := range []string{"symlink", "fifo"} {
			t.Run(store.name+" "+kind, func(t *testing.T) {
				home, checkout, file, _ := writtenStore(t, store)
				valid, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(file); err != nil {
					t.Fatal(err)
				}
				if kind == "symlink" {
					copied := filepath.Join(t.TempDir(), "copy")
					if err := os.WriteFile(copied, valid, 0o600); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(copied, file); err != nil {
						t.Fatal(err)
					}
				} else {
					defer feedFIFO(t, file, valid)()
				}
				if err := store.read(home, checkout); err == nil {
					t.Fatalf("read of a %s %s = nil, want a refusal", kind, store.name)
				}
			})
		}
	}
}

// scopedStore writes and reads one file in the assignment scope of a checkout.
type scopedStore struct {
	name  string
	write func(t *testing.T, home, checkout string) error
	read  func(home, checkout string) error
}

// scopedStores answers the two scoped files: a staged excerpt and a retained response.
func scopedStores() []scopedStore {
	return []scopedStore{
		{"stage", func(_ *testing.T, home, checkout string) error {
			return rb.Stage(home, checkout, "example", "1", "tests", rb.Staged{Mutant: 1})
		}, func(home, checkout string) error {
			_, err := rb.ReadStaged(home, checkout, "example", "1", "tests")
			return err
		}},
		{"retained", func(t *testing.T, home, checkout string) error {
			var sink bytes.Buffer
			owner := rb.New(home, &sink, &sink, func() string { return checkout }, false)
			rb.RespondWith(t, owner, &sink, rb.StdoutLines("kept\n"))
			return owner.Retain(home, checkout, retainCall)
		}, func(home, checkout string) error {
			_, err := rb.ReadRetained(home, checkout)
			return err
		}},
	}
}

// A symlink at any store level refuses the read of a stage and of a retained response. The
// levels come from the written store, so a new level gets its own subtest. Each link points to
// the moved level, so the valid file stays reachable through the link, and only the no-follow
// check of each level refuses the read.
func TestReadRefusesASymlinkedStoreLevel(t *testing.T) {
	for _, store := range scopedStores() {
		_, _, _, layout := writtenStore(t, store)
		for depth := range layout {
			level := fmt.Sprintf("level %d of %d above the file", depth+1, len(layout))
			t.Run(store.name+" "+level, func(t *testing.T) {
				home, checkout, _, levels := writtenStore(t, store)
				if len(levels) != len(layout) {
					t.Fatalf("store levels = %q, want %d", levels, len(layout))
				}
				if err := store.read(home, checkout); err != nil {
					t.Fatalf("read before the link: %v", err)
				}
				dir := levels[depth]
				moved := filepath.Join(t.TempDir(), "moved")
				if err := os.Rename(dir, moved); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(moved, dir); err != nil {
					t.Fatal(err)
				}
				if err := store.read(home, checkout); err == nil {
					t.Fatalf("read through a symlinked %s = nil, want a refusal", level)
				}
			})
		}
	}
}

// writtenStore writes the one file of store in a new home and answers the home, the checkout,
// the file, and each directory between the file and the home, from the file's own directory
// upward. The home is not a level. A store with no level below the home fails the test, so a level list
// is never empty.
func writtenStore(t *testing.T, store scopedStore) (home, checkout, file string, levels []string) {
	t.Helper()
	home = rb.PrivateHome(t)
	_, checkout, _ = responseboundtest.AssignmentCheckout(t, home)
	if err := store.write(t, home, checkout); err != nil {
		t.Fatal(err)
	}
	files := responseboundtest.StoreFiles(t, home)
	if len(files) != 1 {
		t.Fatalf("store files = %q, want one", files)
	}
	file = files[0]
	for dir := filepath.Dir(file); dir != home; dir = filepath.Dir(dir) {
		if dir == filepath.Dir(dir) {
			t.Fatalf("store file %s is not below the home %s", file, home)
		}
		levels = append(levels, dir)
	}
	if len(levels) == 0 {
		t.Fatalf("store file %s has no store level below the home", file)
	}
	return home, checkout, file, levels
}

// feedFIFO makes path a FIFO whose writer sends data once a reader opens it. The answered
// function opens and closes a reader without blocking, so the writer ends when no read came.
func feedFIFO(t *testing.T, path string, data []byte) func() {
	t.Helper()
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		writer, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			return
		}
		_, _ = writer.Write(data)
		_ = writer.Close()
	}()
	return func() {
		if reader, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0); err == nil {
			_ = reader.Close()
		}
		<-done
	}
}

// The census record and retention answer one assignment scope for the same inputs. A
// relative Bench home names the assignment of its absolute form, and a response retained
// through the relative home reads back through the absolute home.
func TestAssignmentScopeAndRetentionAgree(t *testing.T) {
	home := rb.PrivateHome(t)
	_, checkout, id := responseboundtest.AssignmentCheckout(t, home)
	t.Chdir(filepath.Dir(home))
	relative := filepath.Base(home)
	if got, ok := rb.AssignmentScope(relative, checkout); !ok || got != id {
		t.Fatalf("AssignmentScope with a relative home = (%q, %v), want (%q, true)", got, ok, id)
	}
	var sink bytes.Buffer
	owner := rb.New(relative, &sink, &sink, func() string { return checkout }, false)
	rb.RespondWith(t, owner, &sink, rb.StdoutLines("kept\n"))
	if err := owner.Retain(relative, checkout, retainCall); err != nil {
		t.Fatal(err)
	}
	if _, err := rb.ReadRetained(home, checkout); err != nil {
		t.Fatalf("retained response through the absolute home: %v", err)
	}
}
