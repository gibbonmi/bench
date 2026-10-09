package responsebound_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	rb "github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

// sizeOf sends writes through owner, finishes it, and answers its size.
func sizeOf(t *testing.T, owner *rb.Owner, sink *bytes.Buffer, writes []rb.Tagged) rb.Size {
	t.Helper()
	rb.RespondWith(t, owner, sink, writes)
	return owner.Size()
}

// The size counts the complete output, both within the bound and past it, and reports
// whether a spill file holds that output.
func TestOwnerSizeCountsTheCompleteOutput(t *testing.T) {
	for _, row := range []struct {
		name    string
		writes  []rb.Tagged
		spilled bool
	}{
		{"inline", rb.StdoutLines(rb.Numbered(1, 3)...), false},
		{"spilled", append(rb.StdoutLines(rb.Numbered(1, 24)...), rb.StderrWrite("open")), true},
	} {
		t.Run(row.name, func(t *testing.T) {
			rb.PrivateHome(t)
			var sink bytes.Buffer
			got := sizeOf(t, rb.New(benchhome.Dir(), &sink, &sink, rb.OutsideRepository, false), &sink, row.writes)
			want := rb.Size{Lines: len(row.writes), Bytes: int64(len(rb.Joined(row.writes))), Spilled: row.spilled}
			if got != want {
				t.Fatalf("Size = %+v, want %+v", got, want)
			}
		})
	}
}

// A failed spill create prints every byte, so the size reads inline. A spill file that
// failed midway still holds part of the output, so the size reads spilled.
func TestOwnerSizeFollowsTheSpillFile(t *testing.T) {
	for _, row := range []struct {
		name    string
		create  func(string) (io.WriteCloser, error)
		spilled bool
	}{
		{"create failure", func(string) (io.WriteCloser, error) { return nil, rb.ErrInjected }, false},
		{"midway failure", func(path string) (io.WriteCloser, error) {
			file, err := rb.ExclusiveCreate(path)
			if err != nil {
				return nil, err
			}
			return rb.FaultAfter(file, 100), nil
		}, true},
	} {
		t.Run(row.name, func(t *testing.T) {
			rb.PrivateHome(t)
			var sink bytes.Buffer
			owner := rb.New(benchhome.Dir(), &sink, &sink, rb.OutsideRepository, false)
			rb.SetCreate(owner, row.create)
			if got := sizeOf(t, owner, &sink, rb.StdoutLines(rb.Numbered(1, 25)...)); got.Spilled != row.spilled || got.Lines != 25 {
				t.Fatalf("Size = %+v, want 25 lines and spilled=%v", got, row.spilled)
			}
		})
	}
}

// A process outside any repository and a repository root outside the pool both take the
// primary scope, so neither names an assignment.
func TestAssignmentScopeRefusesThePrimaryScope(t *testing.T) {
	home := rb.PrivateHome(t)
	if id, ok := rb.AssignmentScope(home, ""); ok {
		t.Fatalf("AssignmentScope outside a repository = %q, want none", id)
	}
	if id, ok := rb.AssignmentScope(home, t.TempDir()); ok {
		t.Fatalf("AssignmentScope of a root outside the pool = %q, want none", id)
	}
}

// The root of a linked worktree at an assignment segment of the pool names that
// assignment, with no retirement input.
func TestAssignmentScopeAnswersTheWorktreeAssignment(t *testing.T) {
	home := rb.PrivateHome(t)
	_, checkout, id := responseboundtest.AssignmentCheckout(t, home)
	if got, ok := rb.AssignmentScope(home, checkout); !ok || got != id {
		t.Fatalf("AssignmentScope of the worktree root = (%q, %v), want (%q, true)", got, ok, id)
	}
}
