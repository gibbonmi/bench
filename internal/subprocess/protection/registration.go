package protection

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// State is a disposition. Safe permits release; protected and uncertain retain the resource.
type State string

const (
	Safe      State = "safe"
	Protected State = "protected"
	Uncertain State = "uncertain"
)

// Disposition is a state with its resources, the users that hold them, and its error.
type Disposition struct {
	State     State
	Resources []Descriptor
	Users     []string
	Err       error
}

// Unresolved names the state and the users of a disposition that retains its resource. A
// safe disposition returns nil.
func (d Disposition) Unresolved() error {
	if d.State == Safe {
		return nil
	}
	return fmt.Errorf("protection is %s for users %v", d.State, d.Users)
}

// Terminal holds the live owner's facts about one launch. NoStart is a proved failed or
// canceled start.
type Terminal struct {
	NoStart     bool
	Process     Status
	Streams     Status
	Cleanup     Status
	Obligations []string
}

// Registration is one launch's user across its resource union. Closed means that no record
// of it may change again.
type Registration struct {
	protocol  *Protocol
	resources []Descriptor
	record    userRecord
	state     State
	closed    bool
	err       error
}

// Prepare registers one launch across the union before Start. It checks every descriptor
// and record before its first mutation, and it locks stores in canonical path order.
func (p *Protocol) Prepare(union []Descriptor) (*Registration, error) {
	r := &Registration{protocol: p, state: Safe, closed: true}
	resources, err := p.union(union)
	if err != nil || len(resources) == 0 {
		return r.fail(Safe, err)
	}
	r.resources = resources
	r.record = userRecord{
		Schema: schema, UserID: newID(), ActorPID: os.Getpid(), Identity: identityPending,
		Process: Pending, Streams: Pending, Cleanup: Pending, Obligations: []string{ObligationIdentity},
	}
	err = r.locked(func() error {
		for _, d := range resources {
			v, err := p.read(d)
			if err == nil && v.resource.Admission != admissionOpen {
				err = fmt.Errorf("%s: %w", d.Store, ErrSealed)
			}
			if err != nil {
				return err
			}
		}
		for i, d := range resources {
			if published, err := r.add(d); err != nil {
				return r.rollback(i, published, err)
			}
		}
		return nil
	})
	if err != nil {
		return r.fail(r.state, err)
	}
	r.state, r.closed = Protected, false
	return r, nil
}

// union checks every descriptor and returns the distinct generations in canonical store
// order. Two generations at one store refuse, because one of them is stale.
func (p *Protocol) union(union []Descriptor) ([]Descriptor, error) {
	var out []Descriptor
	for _, d := range union {
		if err := d.check(); err != nil {
			return nil, err
		}
		if err := p.owner(d); err != nil {
			return nil, err
		}
		switch i := slices.IndexFunc(out, func(o Descriptor) bool { return o.Store == d.Store }); {
		case i < 0:
			out = append(out, d)
		case out[i] != d:
			return nil, fmt.Errorf("two generations claim the store %s", d.Store)
		}
	}
	slices.SortFunc(out, func(a, b Descriptor) int { return strings.Compare(a.Store, b.Store) })
	return out, nil
}

// Started publishes the actual identity across the union. A partial update leaves the other
// records pending and the launch uncertain.
func (r *Registration) Started(pid, pgid int) error {
	if r.closed || r.record.Identity != identityPending {
		return fmt.Errorf("protection: Started needs an open pending registration, not %s", r.state)
	}
	next := r.record
	next.Identity, next.PID, next.PGID, next.Obligations = identityKnown, pid, pgid, []string{ObligationProcess}
	if err := r.locked(func() error { return r.writeAll(next) }); err != nil {
		_, err = r.fail(Uncertain, err)
		return err
	}
	r.record = next
	return nil
}

// Finish publishes the launch's terminal facts. A complete user or a proved no-start removes
// its records durably; any other user keeps them with each unresolved obligation.
func (r *Registration) Finish(t Terminal) error {
	if r.closed {
		return r.err
	}
	if t.NoStart && r.record.Identity != identityPending {
		return fmt.Errorf("protection: a %s user cannot record no-start", r.record.Identity)
	}
	next := r.record.terminal(t)
	if next.Identity == identityKnown && observeGroup(r.protocol.observer(), next.PGID).State != GroupAbsent && !slices.Contains(next.Obligations, ObligationGroup) {
		next.Obligations = append(next.Obligations, ObligationGroup)
	}
	state := r.protocol.grade(next)
	err := r.locked(func() error {
		if state != Safe {
			return r.writeAll(next)
		}
		for _, d := range r.resources {
			if err := r.protocol.removeUser(d.Store, next.UserID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		_, err = r.fail(Uncertain, err)
		return err
	}
	r.record, r.state, r.closed = next, state, state == Safe
	return nil
}

// Resources returns the distinct generations in canonical store order.
func (r *Registration) Resources() []Descriptor { return slices.Clone(r.resources) }

// Disposition returns the registration's state with its resources and its user.
func (r *Registration) Disposition() Disposition {
	d := Disposition{State: r.state, Resources: r.Resources(), Err: r.err}
	if r.record.UserID != "" {
		d.Users = []string{r.record.UserID}
	}
	return d
}

func (r *Registration) fail(state State, err error) (*Registration, error) {
	r.state, r.err, r.closed = state, err, true
	return r, err
}

// locked runs fn while it holds the protocol lock of every store in the union.
func (r *Registration) locked(fn func() error) error {
	release, err := lock(stores(r.resources))
	if err != nil {
		return err
	}
	defer release()
	return fn()
}

// add durably creates the user's directory in one store and publishes its pending record.
// published reports that the record may be visible.
func (r *Registration) add(d Descriptor) (bool, error) {
	if err := r.protocol.makeDir(r.dir(d)); err != nil {
		return false, err
	}
	return r.write(d, r.record)
}

// rollback publishes no-start in each store that may show the pending user, and it removes
// the failed store's empty directory. It deletes no written record; a failure is uncertain.
func (r *Registration) rollback(failed int, published bool, cause error) error {
	noStart := r.record.terminal(Terminal{NoStart: true})
	var undone []error
	for i, d := range r.resources[:failed+1] {
		if i < failed || published {
			_, err := r.write(d, noStart)
			undone = append(undone, err)
			continue
		}
		undone = append(undone, r.protocol.remove(filepath.Dir(r.dir(d)), r.dir(d)))
	}
	if err := errors.Join(undone...); err != nil {
		r.state = Uncertain
		return errors.Join(cause, err)
	}
	return cause
}

// writeAll publishes next in every store and stops at the first failure.
func (r *Registration) writeAll(next userRecord) error {
	for _, d := range r.resources {
		if _, err := r.write(d, next); err != nil {
			return err
		}
	}
	return nil
}

// write publishes one store's copy of rec. published reports that new bytes may be visible.
func (r *Registration) write(d Descriptor, rec userRecord) (bool, error) {
	rec.ResourceID = d.ResourceID
	if err := rec.valid(rec.UserID); err != nil {
		return false, err
	}
	err := r.protocol.publish(filepath.Join(r.dir(d), userFile), rec)
	var replaceErr *durablefile.Error
	return err == nil || errors.As(err, &replaceErr) && replaceErr.Published, err
}

func (r *Registration) dir(d Descriptor) string {
	return filepath.Join(d.Store, usersDir, r.record.UserID)
}

// terminal applies the owner's facts to a copy of the record.
func (u userRecord) terminal(t Terminal) userRecord {
	if t.NoStart {
		u.Identity, u.PID, u.PGID = identityNoStart, 0, 0
		t = Terminal{Process: Complete, Streams: Complete, Cleanup: Complete}
	}
	u.Process, u.Streams, u.Cleanup = t.Process, t.Streams, t.Cleanup
	u.Obligations = append([]string{}, t.Obligations...)
	return u
}

func stores(resources []Descriptor) []string {
	paths := make([]string, len(resources))
	for i, d := range resources {
		paths[i] = d.Store
	}
	return paths
}
