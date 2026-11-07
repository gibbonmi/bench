package protection

import (
	"errors"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
)

// SealAndObserve closes admission and grades every user of the resource under its protocol
// lock. Only a safe resource runs release and retires its generation, still under that lock.
// The retired store keeps a retired record of the resource in place of the generation. An
// absent store names no user in the current boot, so its resource is safe and SealAndObserve
// runs release without creating the store.
func (p *Protocol) SealAndObserve(d Descriptor, release func() error) (Disposition, error) {
	if err := d.check(); err != nil {
		return uncertain(d, err)
	}
	if err := p.owner(d); err != nil {
		return uncertain(d, err)
	}
	if absent(d.Store) {
		return unused(d, release)
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
	// The retired record goes before the generation. A crash between the two leaves a sealed
	// generation with no user, and the next recovery retires it again.
	retired := retiredRecord{Schema: schema, Kind: d.Kind, ResourcePath: d.ResourcePath, Binding: d.Binding}
	if err := p.publish(filepath.Join(d.Store, retiredFile), retired); err != nil {
		return uncertain(d, err)
	}
	if err := p.remove(d.Store, filepath.Join(d.Store, resourceFile)); err != nil {
		return uncertain(d, err)
	}
	return result, nil
}

// Recover applies the same observation to a resource that its owner found on disk. Its
// published generation names the descriptor, and the kind's validator must accept it. The
// owner derives store, so an absent store is a safe resource with no users.
func (p *Protocol) Recover(store string, release func() error) (Disposition, error) {
	if absent(store) {
		return unused(Descriptor{Store: store}, release)
	}
	if p == nil {
		return uncertain(Descriptor{Store: store}, errors.New("protection: no protocol is installed"))
	}
	var r resourceRecord
	if err := readRecord(filepath.Join(store, resourceFile), p.ReadLimit, &r); err != nil {
		return uncertain(Descriptor{Store: store}, err)
	}
	return p.SealAndObserve(Descriptor{Store: store, ResourceID: r.ResourceID, Kind: r.Kind, ResourcePath: r.ResourcePath, Binding: r.Binding}, release)
}

// NoGeneration reports that the result of Recover found a store that holds no generation,
// such as a store whose generation retired. Such a store names no user, as an absent store
// does. Every other failure of Recover keeps the resource.
func NoGeneration(d Disposition, err error) bool {
	return errors.Is(err, fs.ErrNotExist) && len(d.Resources) == 1 && d.Resources[0].ResourceID == ""
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

// absent reports whether no file exists at path.
func absent(path string) bool {
	_, err := os.Lstat(path)
	return errors.Is(err, fs.ErrNotExist)
}

// unused runs release for a resource whose store is absent and grades it safe with no users.
func unused(d Descriptor, release func() error) (Disposition, error) {
	if release != nil {
		if err := release(); err != nil {
			return uncertain(d, err)
		}
	}
	return Disposition{State: Safe, Resources: []Descriptor{d}}, nil
}

func uncertain(d Descriptor, err error) (Disposition, error) {
	return Disposition{State: Uncertain, Resources: []Descriptor{d}, Err: err}, err
}
