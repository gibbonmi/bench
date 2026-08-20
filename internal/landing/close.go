// close.go owns the tickets-only predicate and the light-path close it drives. A
// light-path change writes specs/<slug>/ carrying tickets and no spec.md; its green
// landing commit consumes that folder instead of flipping a status line. The predicate
// is exported because `bench status` counts the same shape this step closes — one
// definition of "tickets-only", read by both.
package landing

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// specsDir is the one parent every spec folder is a direct child of.
const specsDir = "specs"

// TicketsOnlyFolder reports whether name is a tickets-only spec folder under root: a
// direct child of specs/ that exists as a directory and contains no spec.md. A name
// that is not a single path element names no direct child and is never tickets-only,
// so a slug spelled as a path or carrying `..` cannot reach the close step. A spec.md
// whose presence cannot be determined answers false, keeping the deletion branch
// closed when the tree cannot be read.
func TicketsOnlyFolder(root, name string) bool {
	if name == "" || name == "." || name == ".." || strings.ContainsRune(name, '/') || strings.ContainsRune(name, filepath.Separator) {
		return false
	}
	folder := TicketsOnlyFolderPath(root, name)
	if info, err := os.Stat(folder); err != nil || !info.IsDir() {
		return false
	}
	_, err := os.Stat(filepath.Join(folder, "spec.md"))
	return errors.Is(err, fs.ErrNotExist)
}

// TicketsOnlyFolderPath is the exact filesystem path of the spec folder name denotes.
// The name is taken verbatim, so spaces and glob characters in a folder name resolve
// to that folder rather than to a pattern.
func TicketsOnlyFolderPath(root, name string) string {
	return filepath.Join(root, specsDir, name)
}

// TicketsOnlyFolders returns every tickets-only slug under root's specs/, sorted. An
// absent specs/ directory is an empty result, not a failure.
func TicketsOnlyFolders(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, specsDir))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var slugs []string
	for _, entry := range entries {
		if TicketsOnlyFolder(root, entry.Name()) {
			slugs = append(slugs, entry.Name())
		}
	}
	sort.Strings(slugs)
	return slugs, nil
}
