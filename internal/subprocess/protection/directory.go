package protection

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/canonicalpath"
)

// ErrRetained marks a release that keeps its resource, because a user is unresolved or the
// observation failed.
var ErrRetained = errors.New("stays retained")

// Directory is the identity adapter of a resource that is one owned directory. Its descriptor
// names the canonical directory, the store that Store derives for it, and the directory's device
// and inode, so a directory that replaced the resource at the same path is a different resource.
// Prefix and Lock are optional rules: a set Prefix starts the directory's base name, and a set
// Lock names a regular file that the directory holds. Noun names the resource in a refusal,
// and Class names the prefixed directory class.
type Directory struct {
	Kind   Kind
	Noun   string
	Prefix string
	Class  string
	Lock   string
}

// Describe names the resource at dir without its generation.
func (a Directory) Describe(dir string) (Descriptor, error) {
	canonical, err := canonicalpath.Resolve(dir)
	if err != nil {
		return Descriptor{}, err
	}
	binding, err := a.identity(canonical)
	if err != nil {
		return Descriptor{}, err
	}
	store, err := a.Store(canonical, binding)
	if err != nil {
		return Descriptor{}, err
	}
	return Descriptor{Store: store, Kind: a.Kind, ResourcePath: canonical, Binding: binding}, nil
}

// Store derives the store of the resource at path with binding. Its name is the kind and a
// digest of the path and the binding, and it resides under the store root, never inside the
// resource. A root that Store cannot certify prevents every launch that names it.
func (a Directory) Store(path, binding string) (string, error) {
	root, err := storeRoot()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(path + "\x00" + binding))
	return filepath.Join(root, a.storePrefix()+hex.EncodeToString(sum[:16])), nil
}

// storePrefix starts the name of each store of the kind. Store names a store with it, and
// ObserveRecorded selects the stores of the kind with it.
func (a Directory) storePrefix() string {
	return string(a.Kind) + "-"
}

// Open opens the generation of the resource at dir under p before any user of it.
func (a Directory) Open(p *Protocol, dir string) (Generation, error) {
	d, err := a.Describe(dir)
	if err == nil {
		d, err = p.OpenResource(d)
	}
	if err != nil {
		return Generation{}, err
	}
	return Generation{Protocol: p, Resource: d}, nil
}

// Validate is the kind's Validator. It accepts only a canonical directory of the kind, the
// store that Store derives for it, and the directory's existing identity.
func (a Directory) Validate(d Descriptor) error {
	if canonical, err := canonicalpath.Resolve(d.ResourcePath); err != nil || canonical != d.ResourcePath {
		return fmt.Errorf("%s path %s is not canonical", a.Noun, d.ResourcePath)
	}
	if a.Prefix != "" && !strings.HasPrefix(filepath.Base(d.ResourcePath), a.Prefix) {
		return fmt.Errorf("%s is not a %s", d.ResourcePath, a.Class)
	}
	if a.Lock != "" {
		if lock, err := os.Lstat(filepath.Join(d.ResourcePath, a.Lock)); err != nil || !lock.Mode().IsRegular() {
			return fmt.Errorf("%s holds no %s lock", d.ResourcePath, a.Noun)
		}
	}
	binding, err := a.identity(d.ResourcePath)
	if err != nil {
		return err
	}
	if binding != d.Binding {
		return fmt.Errorf("binding %s does not name %s %s", d.Binding, a.Noun, d.ResourcePath)
	}
	// The store derives from the checked path and binding, so it comes last.
	return a.derived(d)
}

// ObserveRecorded seals and grades the users of each store of the kind whose generation
// records the resource path path. An owner calls it for a resource that it can no longer
// describe, such as a directory that no longer exists, so no store derives from the
// resource. A store whose generation retired names no user. Each recorded generation must
// reside in the store that Store derives for its recorded path and binding, and
// SealAndObserve grades it. The result is the worst disposition over every such store, with
// each unresolved user. A store root or a store of the kind that ObserveRecorded cannot read
// is uncertain, because its generation can record path. readLimit bounds each record read.
func (a Directory) ObserveRecorded(path string, readLimit int64) (Disposition, error) {
	recorded, err := canonicalpath.Resolve(path)
	if err != nil {
		return uncertain(Descriptor{ResourcePath: path}, err)
	}
	root, entries, err := storeEntries()
	if err != nil {
		return uncertain(Descriptor{ResourcePath: recorded}, err)
	}
	p := &Protocol{Validators: Validators{a.Kind: a.derived}, ReadLimit: readLimit}
	result := Disposition{State: Safe}
	var failed []error
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), a.storePrefix()) {
			continue
		}
		store := filepath.Join(root, entry.Name())
		var r resourceRecord
		switch err := readRecord(filepath.Join(store, resourceFile), readLimit, &r); {
		case errors.Is(err, fs.ErrNotExist):
			continue
		case err != nil:
			failed = append(failed, err)
			continue
		case r.Kind != a.Kind || r.ResourcePath != recorded:
			continue
		}
		disposition, err := p.SealAndObserve(Descriptor{Store: store, ResourceID: r.ResourceID, Kind: r.Kind, ResourcePath: r.ResourcePath, Binding: r.Binding}, nil)
		failed = append(failed, err)
		result.Resources = append(result.Resources, disposition.Resources...)
		result.Users = append(result.Users, disposition.Users...)
		if rank(disposition.State) > rank(result.State) {
			result.State = disposition.State
		}
	}
	if err := errors.Join(failed...); err != nil {
		result.State, result.Err = Uncertain, err
		return result, err
	}
	return result, nil
}

// derived accepts only the store that Store derives for the descriptor's path and binding.
// It completes Validate, and it alone validates a recorded generation whose resource
// ObserveRecorded cannot describe.
func (a Directory) derived(d Descriptor) error {
	store, err := a.Store(d.ResourcePath, d.Binding)
	if err != nil {
		return err
	}
	if d.Store != store {
		return fmt.Errorf("store %s is not the protection store of %s %s", d.Store, a.Noun, d.ResourcePath)
	}
	return nil
}

// identity is the directory's device and inode.
func (a Directory) identity(dir string) (string, error) {
	info, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.IsDir() || !ok {
		return "", fmt.Errorf("%s is not a %s directory", dir, a.Noun)
	}
	return fmt.Sprintf("%d:%d", stat.Dev, stat.Ino), nil
}

// Generation is the generation that one owner holds on its resource. Its zero value names
// no generation, which a process with no installed protocol holds.
type Generation struct {
	Protocol *Protocol
	Resource Descriptor
}

// Discard deletes the store of a generation. A zero generation has no store.
func (g Generation) Discard() error {
	if g.Resource == (Descriptor{}) {
		return nil
	}
	return os.RemoveAll(g.Resource.Store)
}

// Resources names the generation for a launch that uses the resource. A zero generation
// names none.
func (g Generation) Resources() []Descriptor {
	if g.Resource == (Descriptor{}) {
		return nil
	}
	return []Descriptor{g.Resource}
}

// Release seals the generation and observes every user. A nil result lets the owner remove
// the resource. Every other result retains it: an unresolved user or a failed observation
// returns ErrRetained, and a failed release returns itself. Release runs only when the
// generation is safe, under its lock, and a zero generation runs it alone. Name names the
// resource in the retention.
func (g Generation) Release(name string, release func() error) error {
	if g.Resource == (Descriptor{}) {
		if release != nil {
			return release()
		}
		return nil
	}
	var failed error
	observed := release
	if release != nil {
		observed = func() error {
			failed = release()
			return failed
		}
	}
	disposition, err := g.Protocol.SealAndObserve(g.Resource, observed)
	switch {
	case failed != nil:
		return failed
	case err != nil:
		return fmt.Errorf("%s %w: %w", name, ErrRetained, err)
	case disposition.State != Safe:
		return fmt.Errorf("%s %w: %w", name, ErrRetained, disposition.Unresolved())
	}
	return nil
}

// Remove is the one safe removal of an owner's resource. It releases the generation as Release
// does, and only a nil release runs remove, which deletes the resource. Then it deletes the
// generation's store, and the store of each resource that the removal deleted inside the
// resource. A store goes only after its resource, so no new generation can open on the resource
// between the two deletions. A zero generation has no store.
func (g Generation) Remove(name string, release, remove func() error) error {
	if err := g.Release(name, release); err != nil {
		return err
	}
	if err := remove(); err != nil {
		return err
	}
	if g.Resource == (Descriptor{}) {
		return nil
	}
	if err := errors.Join(g.Discard(), g.Protocol.discardNested(g.Resource.ResourcePath)); err != nil {
		return fmt.Errorf("%s: remove protection store: %w", name, err)
	}
	return nil
}

// discardNested deletes each store under the store root whose generation names a resource
// inside dir that no longer exists. Such a resource has no user that can still reach it. A
// store without a readable generation goes only through discardRetired.
func (p *Protocol) discardNested(dir string) error {
	root, entries, err := storeEntries()
	if err != nil {
		return err
	}
	inside := func(path string) bool { return strings.HasPrefix(path, dir+string(filepath.Separator)) }
	var failed []error
	for _, entry := range entries {
		store := filepath.Join(root, entry.Name())
		if !entry.IsDir() {
			continue
		}
		var r resourceRecord
		if readRecord(filepath.Join(store, resourceFile), p.ReadLimit, &r) != nil {
			failed = append(failed, p.discardRetired(store, func(retired Descriptor) bool { return inside(retired.ResourcePath) }))
			continue
		}
		if inside(r.ResourcePath) && absent(r.ResourcePath) {
			failed = append(failed, os.RemoveAll(store))
		}
	}
	return errors.Join(failed...)
}

// discardRetired deletes store when its generation retired, its retired record names a
// resource that selected accepts and that no longer exists, and its users directory holds no
// entry. Each condition holds under a store lock that it never waits for, so a store that
// another process holds, or that it cannot lock or read, stays with no error. Only a failed
// deletion is an error.
func (p *Protocol) discardRetired(store string, selected func(Descriptor) bool) error {
	// Only a retired store is locked, so a live store never contends with a sweep.
	generation := filepath.Join(store, resourceFile)
	if !absent(generation) || absent(filepath.Join(store, retiredFile)) {
		return nil
	}
	unlock, err := tryLock(store)
	if err != nil {
		return nil
	}
	defer unlock()
	if !absent(generation) {
		return nil
	}
	r, err := readRetired(store, p.ReadLimit)
	if err != nil || !selected(r.resource()) || !absent(r.ResourcePath) {
		return nil
	}
	users := filepath.Join(store, usersDir)
	if err := checkDir(users); err != nil {
		return nil
	}
	if entries, err := os.ReadDir(users); err != nil || len(entries) != 0 {
		return nil
	}
	return os.RemoveAll(store)
}
