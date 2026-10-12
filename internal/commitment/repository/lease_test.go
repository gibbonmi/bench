package repository_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/intent"
)

// leasedDestination is the destination ref that each lease test leases.
const leasedDestination = "refs/heads/main"

// ledgerLeases is the lease list of the ledger at root.
func ledgerLeases(t *testing.T, root string) []intent.LandingLease {
	t.Helper()
	leases, err := commitmenttest.Leases(root)
	if err != nil {
		t.Fatal(err)
	}
	return leases
}

// atRelease is what a competing acquire finds when it takes the intent lock first after the
// publish decision. It holds the ledger's leases and the result of a lock attempt on the lease.
type atRelease struct {
	leases       []intent.LandingLease
	lockErr, err error
}

// takeIntentLockFirst polls for the intent lock with no pause. A transaction resolves its
// address through Git before its first poll and pauses between polls, so this poll takes the
// lock first. Under that lock it reads what a competing acquire reads, and then it releases
// the lock. It writes no owner into the lock file, because it holds the lock for a short time.
func takeIntentLockFirst(root, intentLock, leaseLock string) atRelease {
	window := bounds.TestDeadline(bounds.IntentLockTimeout)
	deadline := time.Now().Add(window)
	for {
		file, err := os.OpenFile(intentLock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			defer os.Remove(intentLock)
			if err := file.Close(); err != nil {
				return atRelease{err: err}
			}
			break
		}
		if !errors.Is(err, os.ErrExist) {
			return atRelease{err: err}
		}
		if time.Now().After(deadline) {
			return atRelease{err: errors.New(bounds.TestTimeoutVerdict("the end of the publish decision", window))}
		}
	}
	leases, err := commitmenttest.Leases(root)
	holder, lockErr := gocache.HoldExclusive(leaseLock)
	_ = holder.Release()
	return atRelease{leases: leases, lockErr: lockErr, err: err}
}

// leased gives source the destination's lease, as a landing holds it before it publishes.
func leased(t *testing.T, root string, source commitrepo.Publication) commitrepo.Publication {
	t.Helper()
	lease, err := commitrepo.Store{Root: root}.AcquireLease(leasedDestination, source.Assignment)
	if err != nil {
		t.Fatal(err)
	}
	source.Lease = lease
	return source
}

// requireAcquire takes the destination's lease for assignment and requires that the ledger
// then holds only the record of assignment.
func requireAcquire(t *testing.T, root, assignment, after string) {
	t.Helper()
	lease, err := commitrepo.Store{Root: root}.AcquireLease(leasedDestination, assignment)
	if err != nil {
		t.Fatalf("acquire %s = %v, want the lease", after, err)
	}
	t.Cleanup(lease.Release)
	want := []intent.LandingLease{{Destination: leasedDestination, Assignment: assignment}}
	if leases := ledgerLeases(t, root); !reflect.DeepEqual(leases, want) {
		t.Fatalf("leases after the acquire %s = %+v, want %+v", after, leases, want)
	}
}

// TestPublishReleasesTheLeaseInItsDecision is BW7. The test takes the intent lock first when
// the publish decision ends. So a release in a later transaction cannot occur before the test reads.
func TestPublishReleasesTheLeaseInItsDecision(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	intentLock, err := intent.LockPath(root)
	if err != nil {
		t.Fatal(err)
	}
	leaseLock, err := commitrepo.LeaseLockPath(root, leasedDestination)
	if err != nil {
		t.Fatal(err)
	}
	competing := make(chan atRelease, 1)
	if err := (commitrepo.Store{Root: root}).PublishAdmitted(leased(t, root, bound), tree, func() error {
		go func() { competing <- takeIntentLockFirst(root, intentLock, leaseLock) }()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	got := <-competing
	if got.err != nil {
		t.Fatal(got.err)
	}
	if len(got.leases) != 0 || got.lockErr != nil {
		t.Fatalf("at the end of the publish decision: leases %+v, lease lock %v; want no lease and a free lock", got.leases, got.lockErr)
	}
	requireAcquire(t, root, unbound.Assignment, "after the publish decision")
}

// TestPublishKeepsTheLeaseWhenTheLedgerWriteFails covers a ledger write that fails after
// the ref update. The publication still succeeds. The lease record stays with no lock
// holder, and the next acquire takes it over.
func TestPublishKeepsTheLeaseWhenTheLedgerWriteFails(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	ledgerPath, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	saved := ledgerPath + ".saved"
	// A non-empty directory at the ledger path makes the rename of the new ledger fail.
	if err := (commitrepo.Store{Root: root}).PublishAdmitted(leased(t, root, bound), tree, func() error {
		if err := os.Rename(ledgerPath, saved); err != nil {
			return err
		}
		return os.MkdirAll(filepath.Join(ledgerPath, "block"), 0o700)
	}); err != nil {
		t.Fatalf("publication whose ledger write fails = %v, want nil after the ref update", err)
	}
	if err := os.RemoveAll(ledgerPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(saved, ledgerPath); err != nil {
		t.Fatal(err)
	}
	stale := []intent.LandingLease{{Destination: leasedDestination, Assignment: bound.Assignment}}
	if leases := ledgerLeases(t, root); !reflect.DeepEqual(leases, stale) {
		t.Fatalf("leases after the failed ledger write = %+v, want %+v", leases, stale)
	}
	requireAcquire(t, root, unbound.Assignment, "over the publication's stale record")
}

// TestLeaseTakesOverADeadOwner is BW9. A killed landing leaves its record, and the kernel
// drops its record lock, so the next acquire finds no holder and replaces the record.
func TestLeaseTakesOverADeadOwner(t *testing.T) {
	root, bound, _, _ := publicationRepo(t)
	dead := intent.LandingLease{Destination: leasedDestination, Assignment: strings.Repeat("d", 32)}
	if err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		state := intent.CommitmentState{}
		if ledger.Commitment != nil {
			state = *ledger.Commitment
		}
		state.Leases = append(state.Leases, dead)
		ledger.Commitment = &state
		return ledger, true, nil
	}, nil); err != nil {
		t.Fatal(err)
	}
	requireAcquire(t, root, bound.Assignment, "over a dead owner's record")
}
