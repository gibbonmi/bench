package protection

import (
	"errors"
	"fmt"
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
// names the canonical directory, the fixed store inside it, and the directory's device and
// inode, so a directory that replaced the resource at the same path is a different resource.
// Prefix and Lock are optional rules: a set Prefix starts the directory's base name, and a set
// Lock names a regular file that the directory holds. Noun names the resource in a refusal,
// and Class names the prefixed directory class.
type Directory struct {
	Kind   Kind
	Store  string
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
	return Descriptor{Store: filepath.Join(canonical, a.Store), Kind: a.Kind, ResourcePath: canonical, Binding: binding}, nil
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

// Validate is the kind's Validator. It accepts only a canonical directory of the kind, its
// fixed store, and the directory's existing identity.
func (a Directory) Validate(d Descriptor) error {
	if canonical, err := canonicalpath.Resolve(d.ResourcePath); err != nil || canonical != d.ResourcePath {
		return fmt.Errorf("%s path %s is not canonical", a.Noun, d.ResourcePath)
	}
	if d.Store != filepath.Join(d.ResourcePath, a.Store) {
		return fmt.Errorf("store %s is not the protection store of %s %s", d.Store, a.Noun, d.ResourcePath)
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
