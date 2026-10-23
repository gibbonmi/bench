package gocache

import (
	"errors"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// Protect opens the protection generation of the cache that h holds shared. The shared lock
// comes first, so an exclusive clean cannot start between the hold and the protection. The
// last release of the hold retires the generation through p.
func (h *Holder) Protect(p *protection.Protocol) (protection.Descriptor, error) {
	if h == nil || h.dir == "" || h.released {
		return protection.Descriptor{}, errors.New("cache protection needs a live shared cache hold")
	}
	open, err := cacheDirectory.Describe(h.dir)
	if err != nil {
		return protection.Descriptor{}, err
	}
	d, err := p.OpenResource(open)
	// A refused clean leaves its generation sealed. Recovery retires it once every user resolves.
	if errors.Is(err, protection.ErrSealed) {
		if _, err = p.Recover(open.Store, nil); err == nil {
			d, err = p.OpenResource(open)
		}
	}
	if err != nil {
		return protection.Descriptor{}, err
	}
	held.Lock()
	defer held.Unlock()
	if entry := held.paths[h.path]; entry != nil {
		entry.protocol, entry.resource = p, d
	}
	return d, nil
}

// Resources is the cache generation that a launch under h registers on. With no hold or no
// installed protocol, a launch registers none.
func (h *Holder) Resources() ([]protection.Descriptor, error) {
	p := protection.Installed()
	if h == nil || p == nil {
		return nil, nil
	}
	d, err := h.Protect(p)
	if err != nil {
		return nil, err
	}
	return []protection.Descriptor{d}, nil
}

// cacheDirectory is the identity adapter of a cache directory, which holds the cache lock.
var cacheDirectory = protection.Directory{Kind: protection.KindCache, Noun: "cache", Lock: LockFile}

// ProtectionStore derives the protection store of the existing cache directory dir.
func ProtectionStore(dir string) (string, error) {
	d, err := cacheDirectory.Describe(dir)
	return d.Store, err
}

// ValidateProtection is the cache identity adapter.
func ValidateProtection(d protection.Descriptor) error { return cacheDirectory.Validate(d) }
