package prospectiveartifact

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// ProtectionStore names the protection store inside a bundle root.
const ProtectionStore = "process-users"

// protect opens the generation of the bundle root before the record names it and before any
// checkout use.
func protect(p *protection.Protocol, root string) (protection.Descriptor, error) {
	d, err := describe(root)
	if err != nil {
		return protection.Descriptor{}, err
	}
	return p.OpenResource(d)
}

// describe names the bundle root's resource without its generation.
func describe(root string) (protection.Descriptor, error) {
	canonical, err := canonicalpath.Resolve(root)
	if err != nil {
		return protection.Descriptor{}, err
	}
	binding, err := bundleIdentity(canonical)
	if err != nil {
		return protection.Descriptor{}, err
	}
	return protection.Descriptor{Store: filepath.Join(canonical, ProtectionStore), Kind: protection.KindProspective, ResourcePath: canonical, Binding: binding}, nil
}

// ValidateProtection is the bundle identity adapter. A prospective descriptor names a canonical
// bundle root, that root's fixed store, and its identity.
func ValidateProtection(d protection.Descriptor) error {
	if canonical, err := canonicalpath.Resolve(d.ResourcePath); err != nil || canonical != d.ResourcePath {
		return fmt.Errorf("bundle path %s is not canonical", d.ResourcePath)
	}
	if d.Store != filepath.Join(d.ResourcePath, ProtectionStore) {
		return fmt.Errorf("store %s is not the protection store of bundle %s", d.Store, d.ResourcePath)
	}
	if !strings.HasPrefix(filepath.Base(d.ResourcePath), BundlePrefix) {
		return fmt.Errorf("%s is not a prospective artifact bundle", d.ResourcePath)
	}
	binding, err := bundleIdentity(d.ResourcePath)
	if err != nil {
		return err
	}
	if binding != d.Binding {
		return fmt.Errorf("binding %s does not name bundle %s", d.Binding, d.ResourcePath)
	}
	return nil
}

// bundleIdentity is the bundle root's device and inode, so a directory that replaced the bundle
// at the same path is a different bundle.
func bundleIdentity(root string) (string, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.IsDir() || !ok {
		return "", fmt.Errorf("%s is not a bundle directory", root)
	}
	return fmt.Sprintf("%d:%d", stat.Dev, stat.Ino), nil
}
