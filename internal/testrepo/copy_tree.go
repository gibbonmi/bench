package testrepo

import (
	"io/fs"
	"os"
	"path/filepath"
)

// CopyTree copies the tree at src into dst: each file with its bytes, each symbolic link
// with its target, and each entry with its permission bits. dst can exist, but it must not
// hold a path that src holds. os.CopyFS creates each file with mode 0o666 plus the source
// execute bits and each directory with mode 0o777, before the umask, so CopyTree sets the
// source permission bits on each copied entry after the copy. It skips a symbolic link,
// because os.Chmod follows the link to its target.
func CopyTree(dst, src string) error {
	if err := os.CopyFS(dst, os.DirFS(src)); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.Type()&fs.ModeSymlink != 0 {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		return os.Chmod(filepath.Join(dst, rel), info.Mode().Perm())
	})
}
