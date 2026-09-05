// close.go owns the tickets-only predicate and the close the reviewed landing drives.
// A light-path change writes specs/<slug>/ with tickets and no spec.md. The reviewed
// landing commit consumes that folder and flips no status line. This file exports the
// predicate: `bench status` counts the same shape that close consumes. One definition
// of "tickets-only" serves both readers.
package landing

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// specsDir is the one parent every spec folder is a direct child of.
const specsDir = "specs"

// TreeReader answers the two questions the tickets-only rule asks of one spec folder,
// over whichever tree the caller reads: the working tree on the first run, or the source
// commit's objects on a resume, when the source checkout is already released.
type TreeReader interface {
	// FolderIsDirectory reports whether specs/<name> exists in the tree as a directory.
	FolderIsDirectory(name string) bool
	// SpecFileAbsent reports whether specs/<name>/spec.md is absent from the tree. A read
	// that fails for any other reason answers false, which keeps the deletion branch closed.
	SpecFileAbsent(name string) bool
}

// TicketsOnly states the rule once, for every reader: name is a single path element,
// specs/<name> exists as a directory, and it carries no spec.md. A name that is not a
// single path element matches no direct child of specs/, so a slug spelled as a path, or
// carrying `..`, is never tickets-only and cannot reach the close step.
func TicketsOnly(reader TreeReader, name string) bool {
	if name == "" || name == "." || name == ".." || strings.ContainsRune(name, '/') || strings.ContainsRune(name, filepath.Separator) {
		return false
	}
	return reader.FolderIsDirectory(name) && reader.SpecFileAbsent(name)
}

// WorkTree reads the spec folder from the checkout at root.
func WorkTree(root string) TreeReader { return workTreeReader{root: root} }

type workTreeReader struct{ root string }

func (r workTreeReader) FolderIsDirectory(name string) bool {
	info, err := os.Stat(TicketsOnlyFolderPath(r.root, name))
	return err == nil && info.IsDir()
}

// A stat error other than not-exist leaves spec.md's presence undetermined, so the
// reader answers false rather than guessing the folder is a close.
func (r workTreeReader) SpecFileAbsent(name string) bool {
	_, err := os.Stat(filepath.Join(TicketsOnlyFolderPath(r.root, name), "spec.md"))
	return errors.Is(err, fs.ErrNotExist)
}

// CommitTree reads the spec folder from one commit's objects in the repository at root.
func CommitTree(root, commit string) TreeReader {
	return commitTreeReader{root: root, commit: commit}
}

type commitTreeReader struct{ root, commit string }

func (r commitTreeReader) FolderIsDirectory(name string) bool {
	return benchgit.OK("-C", r.root, "cat-file", "-e", r.commit+":"+ClosedFolderPath(name))
}

// A `show` that fails for any reason answers absent: over a commit's objects the only
// reachable failure is the missing path, and a repository this landing cannot read
// refuses earlier.
func (r commitTreeReader) SpecFileAbsent(name string) bool {
	_, err := benchgit.Raw("-C", r.root, "show", r.commit+":"+ClosedFolderPath(name)+"/spec.md")
	return err != nil
}

// TicketsOnlyFolder reports whether name is a tickets-only spec folder in the checkout
// at root. It is the working-tree spelling of TicketsOnly, and `bench status` and the
// first run's close both read it.
func TicketsOnlyFolder(root, name string) bool {
	return TicketsOnly(WorkTree(root), name)
}

// TicketsOnlyFolderPath is the exact filesystem path of the spec folder name denotes.
// The function takes the name verbatim, so spaces and glob characters in a folder
// name resolve to that folder, not to a pattern.
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

// ClosedFolderPath is the repository-relative slash path of the tickets-only folder the
// reviewed landing's close consumes. Every caller names the folder through this one
// spelling.
func ClosedFolderPath(name string) string {
	return specsDir + "/" + name
}
