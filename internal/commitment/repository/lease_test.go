package repository_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
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

// TestPublishReleasesTheLeaseInItsDecision is BW7. A competing acquire that the publish
// decision makes wait on the intent lock runs before any later transaction could, so it
// takes the lease only when the publish decision itself released it.
func TestPublishReleasesTheLeaseInItsDecision(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	store := commitrepo.Store{Root: root}
	lease, err := store.AcquireLease(leasedDestination, bound.Assignment)
	if err != nil {
		t.Fatal(err)
	}
	bound.Lease = lease
	type acquired struct {
		lease *commitrepo.Lease
		err   error
	}
	competing := make(chan acquired, 1)
	if err := store.PublishAdmitted(bound, tree, func() error {
		go func() {
			next, err := store.AcquireLease(leasedDestination, unbound.Assignment)
			competing <- acquired{next, err}
		}()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	got := <-competing
	if got.err != nil {
		t.Fatalf("competing acquire after the publish decision = %v, want the lease", got.err)
	}
	t.Cleanup(got.lease.Release)
	want := []intent.LandingLease{{Destination: leasedDestination, Assignment: unbound.Assignment}}
	if leases := ledgerLeases(t, root); !reflect.DeepEqual(leases, want) {
		t.Fatalf("leases after the competing acquire = %+v, want %+v", leases, want)
	}
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
	lease, err := commitrepo.Store{Root: root}.AcquireLease(leasedDestination, bound.Assignment)
	if err != nil {
		t.Fatalf("acquire over a dead owner's record = %v, want the takeover", err)
	}
	t.Cleanup(lease.Release)
	want := []intent.LandingLease{{Destination: leasedDestination, Assignment: bound.Assignment}}
	if leases := ledgerLeases(t, root); !reflect.DeepEqual(leases, want) {
		t.Fatalf("leases after the takeover = %+v, want %+v", leases, want)
	}
}
