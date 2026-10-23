package protection

import (
	"errors"
	"maps"
	"path/filepath"
	"slices"
)

// SealAndObserve closes admission and grades every user of the resource under its protocol
// lock. Only a safe resource runs release and retires its generation, still under that lock.
func (p *Protocol) SealAndObserve(d Descriptor, release func() error) (Disposition, error) {
	if err := d.check(); err != nil {
		return uncertain(d, err)
	}
	if err := p.owner(d); err != nil {
		return uncertain(d, err)
	}
	unlock, err := lock([]string{d.Store})
	if err != nil {
		return uncertain(d, err)
	}
	defer unlock()
	v, err := p.read(d)
	if err != nil {
		return uncertain(d, err)
	}
	if v.resource.Admission == admissionOpen {
		v.resource.Admission = admissionSealed
		if err := p.publish(filepath.Join(d.Store, resourceFile), v.resource); err != nil {
			return uncertain(d, err)
		}
	}
	result := Disposition{State: Safe, Resources: []Descriptor{d}}
	for _, id := range slices.Sorted(maps.Keys(v.users)) {
		state := p.grade(v.users[id])
		if state != Safe {
			result.Users = append(result.Users, id)
		}
		if rank(state) > rank(result.State) {
			result.State = state
		}
	}
	if result.State != Safe {
		return result, nil
	}
	if release != nil {
		if err := release(); err != nil {
			return uncertain(d, err)
		}
	}
	for id := range v.users {
		if err := p.removeUser(d.Store, id); err != nil {
			return uncertain(d, err)
		}
	}
	if err := p.remove(d.Store, filepath.Join(d.Store, resourceFile)); err != nil {
		return uncertain(d, err)
	}
	return result, nil
}

// Recover applies the same observation to a resource that its owner found on disk. Its
// published generation names the descriptor, and the kind's validator must accept it.
func (p *Protocol) Recover(store string, release func() error) (Disposition, error) {
	if p == nil {
		return uncertain(Descriptor{Store: store}, errors.New("protection: no protocol is installed"))
	}
	var r resourceRecord
	if err := readRecord(filepath.Join(store, resourceFile), p.ReadLimit, &r); err != nil {
		return uncertain(Descriptor{Store: store}, err)
	}
	return p.SealAndObserve(Descriptor{Store: store, ResourceID: r.ResourceID, Kind: r.Kind, ResourcePath: r.ResourcePath, Binding: r.Binding}, release)
}

// grade is the one release predicate for a user. A no-start is complete. A pending user needs
// a currently absent actor: its start gate opens only after the actor publishes the known
// identity, so the program of a pending user whose actor is gone never ran. A known user
// needs current group absence, then finished local workers or a currently absent actor. It
// sends only signal zero, and record age grants nothing.
func (p *Protocol) grade(u userRecord) State {
	switch u.Identity {
	case identityNoStart:
		return Safe
	case identityPending:
	default:
		if state := stateOf(observeGroup(p.observer(), u.PGID)); state != Safe {
			return state
		}
		if u.Process == Complete && u.Streams != Pending && u.Cleanup != Pending {
			return Safe
		}
	}
	return stateOf(Classify(p.observer()(u.ActorPID, 0)))
}

func stateOf(o Observation) State {
	switch o.State {
	case GroupAbsent:
		return Safe
	case GroupPresent:
		return Protected
	}
	return Uncertain
}

func rank(s State) int { return slices.Index([]State{Safe, Protected, Uncertain}, s) }

func uncertain(d Descriptor, err error) (Disposition, error) {
	return Disposition{State: Uncertain, Resources: []Descriptor{d}, Err: err}, err
}
