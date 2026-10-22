package gocache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// ProtectionStore names the protection store inside the Bench build cache. Go's clean
// removes the two-hex subdirectories alone, so the store survives a clean.
const ProtectionStore = "bench-process-users"

// Protect opens the protection generation of the cache that h holds shared. The shared lock
// comes first, so an exclusive clean cannot start between the hold and the protection.
func (h *Holder) Protect(p *protection.Protocol) (protection.Descriptor, error) {
	if h == nil || h.dir == "" || h.released {
		return protection.Descriptor{}, errors.New("cache protection needs a live shared cache hold")
	}
	dir, err := canonicalpath.Resolve(h.dir)
	if err != nil {
		return protection.Descriptor{}, err
	}
	binding, err := cacheIdentity(dir)
	if err != nil {
		return protection.Descriptor{}, err
	}
	return p.OpenResource(protection.Descriptor{Store: filepath.Join(dir, ProtectionStore), Kind: protection.KindCache, ResourcePath: dir, Binding: binding})
}

// ValidateProtection is the cache identity adapter. A cache descriptor names a canonical
// directory that holds the cache lock, that directory's fixed store, and its identity.
func ValidateProtection(d protection.Descriptor) error {
	if canonical, err := canonicalpath.Resolve(d.ResourcePath); err != nil || canonical != d.ResourcePath {
		return fmt.Errorf("cache path %s is not canonical", d.ResourcePath)
	}
	if d.Store != filepath.Join(d.ResourcePath, ProtectionStore) {
		return fmt.Errorf("store %s is not the protection store of cache %s", d.Store, d.ResourcePath)
	}
	if lock, err := os.Lstat(filepath.Join(d.ResourcePath, LockFile)); err != nil || !lock.Mode().IsRegular() {
		return fmt.Errorf("%s holds no cache lock", d.ResourcePath)
	}
	binding, err := cacheIdentity(d.ResourcePath)
	if err != nil {
		return err
	}
	if binding != d.Binding {
		return fmt.Errorf("binding %s does not name cache %s", d.Binding, d.ResourcePath)
	}
	return nil
}

// cacheIdentity is the cache directory's device and inode, so a directory that replaced the
// cache at the same path is a different cache.
func cacheIdentity(dir string) (string, error) {
	info, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.IsDir() || !ok {
		return "", fmt.Errorf("%s is not a cache directory", dir)
	}
	return fmt.Sprintf("%d:%d", stat.Dev, stat.Ino), nil
}
