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
	return filepath.Join(root, string(a.Kind)+"-"+hex.EncodeToString(sum[:16])), nil
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
		return fmt.Errorf("%s %w: protection is %s for users %v", name, ErrRetained, disposition.State, disposition.Users)
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
// store without a readable generation stays.
func (p *Protocol) discardNested(dir string) error {
	root, err := storeRoot()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	var failed []error
	for _, entry := range entries {
		store := filepath.Join(root, entry.Name())
		var r resourceRecord
		if !entry.IsDir() || readRecord(filepath.Join(store, resourceFile), p.ReadLimit, &r) != nil || !strings.HasPrefix(r.ResourcePath, dir+string(filepath.Separator)) {
			continue
		}
		if _, err := os.Lstat(r.ResourcePath); errors.Is(err, fs.ErrNotExist) {
			failed = append(failed, os.RemoveAll(store))
		}
	}
	return errors.Join(failed...)
}
