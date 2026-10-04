package repository_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

const publicationDeliverable = "specs/x/spec.md"

// publicationRepo approves one delivery and registers two assignments. Only the bound
// assignment starts that delivery. The returned tree adds a production file to main.
func publicationRepo(t *testing.T) (root string, bound, unbound commitrepo.Publication, tree string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, root, publicationDeliverable, "# x\n\nStatus: staged\n")
	commitmenttest.SeedAdmission(t, root, publicationDeliverable)
	commitmenttest.Commit(t, root, "approve fixture delivery")
	boundPath := commitmenttest.Assignment(t, root, "bound")
	commitmenttest.Admit(t, boundPath, "bound", publicationDeliverable)
	unboundPath := commitmenttest.Assignment(t, root, "unbound")
	commitmenttest.Write(t, boundPath, "tool.go", "package tool\n")
	commitmenttest.Commit(t, boundPath, "production file")
	tree = gittest.Output(t, boundPath, "rev-parse", "HEAD^{tree}")
	return root, publication(t, root, boundPath), publication(t, root, unboundPath), tree
}

// publication is the frozen identity of the one assignment that owns worktree.
func publication(t *testing.T, root, worktree string) commitrepo.Publication {
	t.Helper()
	ledger, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	owners := intent.AssignmentsOwning(ledger.Assignments, worktree)
	if len(owners) != 1 {
		t.Fatalf("assignments owning %s = %+v, want one", worktree, owners)
	}
	return commitrepo.Publication{Assignment: owners[0].ID, Request: owners[0].Request, Worktree: owners[0].Worktree}
}

// The decision grades the frozen source assignment, not any assignment in the ledger.
// The ledger holds two active records, so a decision that ignores the assignment ID
// grades the wrong record for either the bound row or the borrowed row, whichever record
// comes first.
func TestAdmitPublicationFrozenIdentity(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	const mismatch = "is not active with its presented request and worktree"
	for _, row := range []struct {
		name   string
		source commitrepo.Publication
		want   string
	}{
		{name: "bound", source: bound},
		{name: "unbound", source: unbound, want: "assignment has no current delivery binding"},
		{name: "borrowed-request-and-worktree", source: commitrepo.Publication{Assignment: unbound.Assignment, Request: bound.Request, Worktree: bound.Worktree}, want: mismatch},
		{name: "other-request", source: commitrepo.Publication{Assignment: bound.Assignment, Request: unbound.Request, Worktree: bound.Worktree}, want: mismatch},
		{name: "other-worktree", source: commitrepo.Publication{Assignment: bound.Assignment, Request: bound.Request, Worktree: unbound.Worktree}, want: mismatch},
	} {
		t.Run(row.name, func(t *testing.T) {
			err := (commitrepo.Store{Root: root}).AdmitPublication(row.source, tree)
			if row.want == "" {
				if err != nil {
					t.Fatalf("AdmitPublication = %v, want admission", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("AdmitPublication = %v, want a refusal naming %q", err, row.want)
			}
		})
	}
}

// An admitted publication runs publish once and returns its error unchanged. A refused
// publication never runs publish.
func TestPublishAdmitted(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	failed := errors.New("ref update failed")
	for _, row := range []struct {
		name      string
		source    commitrepo.Publication
		publish   error
		published int
	}{
		{name: "admitted", source: bound, published: 1},
		{name: "publish-fails", source: bound, publish: failed, published: 1},
		{name: "refused", source: unbound},
	} {
		t.Run(row.name, func(t *testing.T) {
			published := 0
			err := (commitrepo.Store{Root: root}).PublishAdmitted(row.source, tree, func() error {
				published++
				return row.publish
			})
			if published != row.published {
				t.Fatalf("publish ran %d times, want %d (err=%v)", published, row.published, err)
			}
			switch {
			case row.published == 0 && err == nil:
				t.Fatal("refused publication returned no error")
			case row.published == 1 && err != row.publish:
				t.Fatalf("PublishAdmitted = %v, want the publish result %v", err, row.publish)
			}
		})
	}
}

// A blocker that a concurrent transaction adds while PublishAdmitted waits for the intent
// lock decides the publication. The test swaps the held lock file for a FIFO, so opening
// that FIFO to write returns only when the waiter reads the lock to judge its staleness.
func TestPublishAdmittedDecidesUnderTheLock(t *testing.T) {
	root, bound, _, tree := publicationRepo(t)
	store := commitrepo.Store{Root: root}
	policy, _, err := store.Policy()
	if err != nil {
		t.Fatal(err)
	}
	address, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	lock := address + ".lock"

	held, proceed := make(chan struct{}), make(chan struct{})
	holder := make(chan error, 1)
	go func() {
		holder <- intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
			close(held)
			<-proceed
			if ledger.Commitment == nil {
				return ledger, false, errors.New("ledger has no commitment state")
			}
			next, err := commitment.SetBlocker(policy, *ledger.Commitment, commitmenttest.DeliveryOutcome, "competing blocker", true)
			if err != nil {
				return ledger, false, err
			}
			ledger.Commitment = &next
			return ledger, true, nil
		}, nil)
	}()
	finished := false
	finish := func() error {
		if finished {
			return nil
		}
		finished = true
		close(proceed)
		return <-holder
	}
	defer func() { _ = finish() }()
	select {
	case <-held:
	case err := <-holder:
		t.Fatalf("holding transaction ended before its decision: %v", err)
	}

	owner, err := os.ReadFile(lock)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(lock+".fifo", 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	if err := os.Rename(lock+".fifo", lock); err != nil {
		t.Fatal(err)
	}

	published := false
	done := make(chan error, 1)
	go func() {
		done <- store.PublishAdmitted(bound, tree, func() error {
			published = true
			return nil
		})
	}()
	type opened struct {
		file *os.File
		err  error
	}
	writer := make(chan opened, 1)
	go func() {
		file, err := os.OpenFile(lock, os.O_WRONLY, 0)
		writer <- opened{file, err}
	}()
	select {
	case got := <-writer:
		if got.err != nil {
			t.Fatal(got.err)
		}
		// The regular lock file replaces the FIFO before the waiter's read ends, so every
		// later poll reads a live owner and keeps waiting.
		if err := os.WriteFile(lock+".held", owner, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(lock+".held", lock); err != nil {
			t.Fatal(err)
		}
		if _, err := got.file.Write(owner); err != nil {
			t.Fatal(err)
		}
		if err := got.file.Close(); err != nil {
			t.Fatal(err)
		}
	case err := <-done:
		// A reader that opens without blocking releases the writer that waits for one.
		if reader, openErr := os.OpenFile(lock, os.O_RDONLY|syscall.O_NONBLOCK, 0); openErr == nil {
			_ = reader.Close()
		}
		t.Fatalf("PublishAdmitted = %v (published=%t) before it waited for the held intent lock", err, published)
	}

	if err := finish(); err != nil {
		t.Fatal(err)
	}
	err = <-done
	if published || err == nil || !strings.Contains(err.Error(), "blocked") {
		t.Fatalf("PublishAdmitted = %v (published=%t), want a blocked refusal that never publishes", err, published)
	}
}
