package preflight

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// newSystemTest reports whether a Go test path the tree does not hold yet joins the
// system suite. The file has no bytes to read, so its directory answers: a directory
// that holds a system-tagged test file is where the system tag lives. The tag parse is
// systemTagged's own, so a new file and a landed one are graded by one rule. Only a
// regular file is read, because a FIFO named like a test file would block the read.
func newSystemTest(root, testPath string) bool {
	if !strings.HasSuffix(testPath, "_test.go") {
		return false
	}
	dir := path.Dir(testPath)
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() && systemTagged(root, path.Join(dir, entry.Name())) {
			return true
		}
	}
	return false
}
