package repository_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

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
	lock, err := intent.LockPath(root)
	if err != nil {
		t.Fatal(err)
	}

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
	// A reader that opens without blocking releases the writer that waits for one.
	releaseWriter := func() {
		if reader, err := os.OpenFile(lock, os.O_RDONLY|syscall.O_NONBLOCK, 0); err == nil {
			_ = reader.Close()
		}
	}
	window := bounds.TestDeadline(bounds.IntentLockTimeout)
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
		releaseWriter()
		t.Fatalf("PublishAdmitted = %v (published=%t) before it waited for the held intent lock", err, published)
	case <-time.After(window):
		releaseWriter()
		t.Fatal(bounds.TestTimeoutVerdict("PublishAdmitted to read the held intent lock", window))
	}

	if err := finish(); err != nil {
		t.Fatal(err)
	}
	err = <-done
	if published || err == nil || !strings.Contains(err.Error(), "blocked") {
		t.Fatalf("PublishAdmitted = %v (published=%t), want a blocked refusal that never publishes", err, published)
	}
}
