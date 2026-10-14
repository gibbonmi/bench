package responsebound_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
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
