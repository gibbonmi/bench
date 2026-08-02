package gate

// The one capture of the tree every identity in this package is computed over. The stripped
// whole-tree identity and the per-component identities both read it, so they answer for the
// same `git add -A` snapshot rather than for two listings taken moments apart, and the
// listing is parsed in exactly one place — a second parse beside it would let one identity
// see a path the other did not.

import (
	"errors"
	"fmt"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// treeEntry is one entry of the snapshot: the raw metadata git recorded — mode, type, and
// object id — and the repository-relative path it recorded them under.
type treeEntry struct {
	Path     string
	Metadata string
}

// treeSnapshot is a parsed listing. Entries stay in listing order, because an identity that
// hashes them in sequence needs an order that is the same on every read; they are indexed by
// path as well, because a declaration that names one file resolves it by lookup rather than
// by scanning the whole tree for it.
type treeSnapshot struct {
	entries []treeEntry
	byPath  map[string]string
}

// readTreeSnapshot lists root's `git add -A` tree. It reads back the whole-tree object
// rather than materializing a second tree, so every identity taken from it answers for one
// snapshot. Any failure to run or parse the listing is an error rather than a shorter set of
// entries: a missing entry silently drops content from whatever hashes it.
func readTreeSnapshot(root string) (treeSnapshot, error) {
	tree := benchgit.TreeHash(root)
	if !treeHashRE.MatchString(tree) {
		return treeSnapshot{}, errors.New("tree unavailable")
	}
	// -z keeps the paths raw: without it git quotes and escapes unusual names, and a quoted
	// path is not the one a declaration names or Scope.Member is defined over.
	listing, err := benchgit.Output("-C", root, "ls-tree", "-r", "-z", "--full-tree", tree)
	if err != nil {
		return treeSnapshot{}, fmt.Errorf("list the tree snapshot: %w", err)
	}
	snapshot := treeSnapshot{byPath: map[string]string{}}
	for _, entry := range strings.Split(listing, "\x00") {
		if entry == "" {
			continue
		}
		metadata, path, separated := strings.Cut(entry, "\t")
		if !separated {
			return treeSnapshot{}, errors.New("unparsable tree snapshot entry")
		}
		snapshot.entries = append(snapshot.entries, treeEntry{Path: path, Metadata: metadata})
		snapshot.byPath[path] = metadata
	}
	return snapshot, nil
}

// entry returns what git recorded for path, and whether the snapshot tracks it at all.
func (s treeSnapshot) entry(path string) (treeEntry, bool) {
	metadata, tracked := s.byPath[path]
	return treeEntry{Path: path, Metadata: metadata}, tracked
}
