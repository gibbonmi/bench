package env

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// ProtectionStore names the protection store inside a private test run.
const ProtectionStore = "process-users"

// protect opens the private generation of the run directory dir before any child uses it.
func protect(p *protection.Protocol, dir string) (protection.Descriptor, error) {
	canonical, err := canonicalpath.Resolve(dir)
	if err != nil {
		return protection.Descriptor{}, err
	}
	binding, err := runIdentity(canonical)
	if err != nil {
		return protection.Descriptor{}, err
	}
	return p.OpenResource(protection.Descriptor{Store: filepath.Join(canonical, ProtectionStore), Kind: protection.KindKitTest, ResourcePath: canonical, Binding: binding})
}

// ValidateProtection is the private test run's identity adapter. A kit-test descriptor names a
// canonical run directory, that directory's fixed store, and its identity.
func ValidateProtection(d protection.Descriptor) error {
	if canonical, err := canonicalpath.Resolve(d.ResourcePath); err != nil || canonical != d.ResourcePath {
		return fmt.Errorf("kit test run path %s is not canonical", d.ResourcePath)
	}
	if d.Store != filepath.Join(d.ResourcePath, ProtectionStore) {
		return fmt.Errorf("store %s is not the protection store of kit test run %s", d.Store, d.ResourcePath)
	}
	binding, err := runIdentity(d.ResourcePath)
	if err != nil {
		return err
	}
	if binding != d.Binding {
		return fmt.Errorf("binding %s does not name kit test run %s", d.Binding, d.ResourcePath)
	}
	return nil
}

// runIdentity is the run directory's device and inode, so a directory that replaced the run
// at the same path is a different run.
func runIdentity(dir string) (string, error) {
	info, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.IsDir() || !ok {
		return "", fmt.Errorf("%s is not a kit test run directory", dir)
	}
	return fmt.Sprintf("%d:%d", stat.Dev, stat.Ino), nil
}
