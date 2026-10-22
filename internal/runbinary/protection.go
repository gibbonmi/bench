package runbinary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// ProtectionStore names the protection store inside a private selection directory.
const ProtectionStore = "process-users"

// runDirPrefix starts the name of every private selection directory.
const runDirPrefix = "bench-run-"

// protect opens the private generation of the selection directory dir before any build or use.
func protect(p *protection.Protocol, dir string) (protection.Descriptor, error) {
	canonical, err := canonicalpath.Resolve(dir)
	if err != nil {
		return protection.Descriptor{}, err
	}
	binding, err := runIdentity(canonical)
	if err != nil {
		return protection.Descriptor{}, err
	}
	return p.OpenResource(protection.Descriptor{Store: filepath.Join(canonical, ProtectionStore), Kind: protection.KindBinary, ResourcePath: canonical, Binding: binding})
}

// ValidateProtection is the selection identity adapter. A binary descriptor names a canonical
// private selection directory, that directory's fixed store, and its identity.
func ValidateProtection(d protection.Descriptor) error {
	if canonical, err := canonicalpath.Resolve(d.ResourcePath); err != nil || canonical != d.ResourcePath {
		return fmt.Errorf("selection path %s is not canonical", d.ResourcePath)
	}
	if d.Store != filepath.Join(d.ResourcePath, ProtectionStore) {
		return fmt.Errorf("store %s is not the protection store of selection %s", d.Store, d.ResourcePath)
	}
	if !strings.HasPrefix(filepath.Base(d.ResourcePath), runDirPrefix) {
		return fmt.Errorf("%s is not a private selection directory", d.ResourcePath)
	}
	binding, err := runIdentity(d.ResourcePath)
	if err != nil {
		return err
	}
	if binding != d.Binding {
		return fmt.Errorf("binding %s does not name selection %s", d.Binding, d.ResourcePath)
	}
	return nil
}

// runIdentity is the selection directory's device and inode, so a directory that replaced the
// selection at the same path is a different selection.
func runIdentity(dir string) (string, error) {
	info, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.IsDir() || !ok {
		return "", fmt.Errorf("%s is not a selection directory", dir)
	}
	return fmt.Sprintf("%d:%d", stat.Dev, stat.Ino), nil
}
