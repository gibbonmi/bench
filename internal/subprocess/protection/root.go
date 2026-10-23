package protection

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// rootMode is the only mode of each store root directory.
const rootMode fs.FileMode = 0o700

// storeBase names the system directory that holds each user's store root. The package tests
// replace it, because a refusal test must change a root, and the real root serves every
// concurrent Bench process of the user.
var storeBase = systemBase

// storeRoot is the root of every store of the effective user in the current boot. No
// environment variable selects it, so a nested CLI with a rewritten environment derives the
// same root. Each root directory must be a real directory with mode 0700 that the effective
// user owns.
func storeRoot() (string, error) {
	boot, err := bootID()
	if err != nil {
		return "", err
	}
	if boot == "" || boot == "." || boot == ".." || strings.ContainsRune(boot, filepath.Separator) {
		return "", fmt.Errorf("boot identifier %q is not one path name", boot)
	}
	user := filepath.Join(storeBase(), fmt.Sprintf("bench-process-%d", os.Geteuid()))
	root := filepath.Join(user, boot)
	for _, dir := range []string{user, root} {
		if err := privateRoot(dir); err != nil {
			return "", err
		}
	}
	return root, nil
}

// privateRoot creates dir when it is absent and certifies it. It never repairs an existing
// directory, so a link, a foreign owner, or another mode stays a refusal.
func privateRoot(dir string) error {
	switch err := os.Mkdir(dir, rootMode); {
	case err == nil:
		// The process umask can narrow a fresh directory, and only this process created it.
		if err := os.Chmod(dir, rootMode); err != nil {
			return err
		}
	case !errors.Is(err, fs.ErrExist):
		return err
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.IsDir() || !ok || info.Mode()&(fs.ModePerm|fs.ModeSetuid|fs.ModeSetgid|fs.ModeSticky) != rootMode || int(stat.Uid) != os.Geteuid() {
		return fmt.Errorf("store root %s is not a private directory of user %d", dir, os.Geteuid())
	}
	return nil
}
