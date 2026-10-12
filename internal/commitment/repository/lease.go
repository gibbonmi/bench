package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// Lease is one landing's hold on its destination ref. It is the ledger record that names
// the landing's assignment, and the record lock that proves that the landing is alive.
type Lease struct {
	root   string
	record intent.LandingLease
	holder *gocache.Holder
}

// LeaseLockPath is the lock file of destination's lease. It lives in the common Git
// directory, so every checkout of the repository locks the same file, and no step removes it.
func LeaseLockPath(root, destination string) (string, error) {
	common, err := git.CommonDir(root)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(destination))
	return filepath.Join(common, "bench-landing-"+hex.EncodeToString(sum[:])[:16]+".lock"), nil
}

// AcquireLease records assignment as the landing that holds destination, in one decision
// under the intent lock. A record whose lock has no holder belongs to a dead landing, so
// the new record replaces it. The compensation releases the lock when the write fails.
func (store Store) AcquireLease(destination, assignment string) (*Lease, error) {
	path, err := LeaseLockPath(store.Root, destination)
	if err != nil {
		return nil, leaseUnavailable(err)
	}
	lease := &Lease{root: store.Root, record: intent.LandingLease{Destination: destination, Assignment: assignment}}
	err = intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		holder, err := gocache.HoldExclusive(path)
		if errors.Is(err, gocache.ErrHeld) {
			return ledger, false, leaseHeld(ledger, destination)
		}
		if err != nil {
			return ledger, false, leaseUnavailable(err)
		}
		lease.holder = holder
		state := storedState(ledger)
		state.Leases = append(slices.DeleteFunc(slices.Clone(state.Leases), func(record intent.LandingLease) bool {
			return record.Destination == destination
		}), lease.record)
		return withRuntime(ledger, state)
	}, lease.unlock)
	if err != nil {
		lease.unlock()
		return nil, err
	}
	return lease, nil
}

// Release removes the lease record in its own transaction and releases the lock. A
// removal that fails leaves a record with no lock holder, which the next acquire takes over.
func (lease *Lease) Release() {
	if lease == nil {
		return
	}
	_ = intent.Transact(lease.root, intent.StrictRead, lease.drop, nil)
	lease.unlock()
}

// drop is the decision that removes the lease record and releases the lock. The lock
// goes inside the decision, so an acquire that waits on the intent lock finds neither.
// A nil lease drops nothing.
func (lease *Lease) drop(ledger intent.Ledger) (intent.Ledger, bool, error) {
	if lease == nil {
		return ledger, false, nil
	}
	lease.unlock()
	state := storedState(ledger)
	state.Leases = slices.DeleteFunc(slices.Clone(state.Leases), func(record intent.LandingLease) bool {
		return record == lease.record
	})
	return withRuntime(ledger, state)
}

func (lease *Lease) unlock() { _ = lease.holder.Release() }

// leaseHeld is the refusal of a landing whose destination another live landing holds. It
// names the assignment of the ledger's record for that destination.
func leaseHeld(ledger intent.Ledger, destination string) error {
	owner := ""
	for _, record := range storedState(ledger).Leases {
		if record.Destination == destination {
			owner = record.Assignment
		}
	}
	return refusalroute.Raised{
		Name:   refusalroute.LandLeaseHeld,
		Err:    fmt.Errorf("landing destination %s is leased by assignment %s", destination, owner),
		Values: map[string]string{refusalroute.FactLeaseOwner: owner},
	}
}

// leaseUnavailable is the refusal of a lock file that cannot be opened or locked.
func leaseUnavailable(err error) error {
	return refusalroute.Raised{Name: refusalroute.LandInfrastructure, Err: fmt.Errorf("landing lease unavailable: %w", err)}
}
