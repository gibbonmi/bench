package gate

// The one capture of the tree every identity in this package is computed over. The stripped
// whole-tree identity and the per-component identities both read it, so they answer for the
// same `git add -A` snapshot rather than for two listings taken moments apart, and the
// listing is parsed in exactly one place — a second parse beside it would let one identity
// see a path the other did not.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// treeEntryMetadataRE is the complete set of shapes a recursive listing records: a file or
// symlink blob, or a gitlink commit, each with a full hex object id. An entry outside it is
// data git never wrote, and a snapshot that kept it would hand identities a fact with no
// source — refusal at parse is the fail-closed direction the blob-read check alone cannot
// give, because entries are hashed whether or not their blobs are ever requested.
var treeEntryMetadataRE = regexp.MustCompile(`^(100644 blob|100755 blob|120000 blob|160000 commit) ([0-9a-f]{40}|[0-9a-f]{64})$`)

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

// treeSource supplies one immutable Git tree and its blob objects to a generation.
type treeSource interface {
	tree() (string, error)
	list(tree string) (string, error)
	blob(object string) ([]byte, error)
}

type treeGeneration struct {
	tree     string
	snapshot treeSnapshot
	source   treeSource
	blobs    map[string]blobResult
}

type blobResult struct {
	data []byte
	err  error
}

type workingTreeSource struct{ root string }

func (s workingTreeSource) tree() (string, error) {
	tree := benchgit.TreeHash(s.root)
	if !treeHashRE.MatchString(tree) {
		return "", errors.New("tree unavailable")
	}
	return tree, nil
}

func (s workingTreeSource) list(tree string) (string, error) {
	return benchgit.Output("-C", s.root, "ls-tree", "-r", "-z", "--full-tree", tree)
}

func (s workingTreeSource) blob(object string) ([]byte, error) {
	return benchgit.Raw("-C", s.root, "cat-file", "blob", object)
}

type prospectiveTreeSource struct {
	root   string
	treeID string
}

func (s prospectiveTreeSource) tree() (string, error) {
	if !treeHashRE.MatchString(s.treeID) {
		return "", errors.New("tree unavailable")
	}
	return s.treeID, nil
}

func (s prospectiveTreeSource) list(tree string) (string, error) {
	return benchgit.Output("-C", s.root, "ls-tree", "-r", "-z", "--full-tree", tree)
}

func (s prospectiveTreeSource) blob(object string) ([]byte, error) {
	return benchgit.Raw("-C", s.root, "cat-file", "blob", object)
}

func captureWorkingTree(root string) (*treeGeneration, error) {
	return captureTreeGeneration(workingTreeSource{root: root})
}

func captureProspectiveTree(root, tree string) (*treeGeneration, error) {
	return captureTreeGeneration(prospectiveTreeSource{root: root, treeID: tree})
}

func (g *treeGeneration) entry(path string) (treeEntry, bool) { return g.snapshot.entry(path) }

func (g *treeGeneration) blob(entry treeEntry) ([]byte, error) {
	fields := strings.Fields(entry.Metadata)
	if len(fields) != 3 || fields[1] != "blob" {
		return nil, errors.New("snapshot entry is not a blob")
	}
	object := fields[2]
	if result, found := g.blobs[object]; found {
		return append([]byte(nil), result.data...), result.err
	}
	data, err := g.source.blob(object)
	result := blobResult{data: append([]byte(nil), data...), err: err}
	g.blobs[object] = result
	return append([]byte(nil), result.data...), result.err
}

func captureTreeGeneration(source treeSource) (*treeGeneration, error) {
	tree, err := source.tree()
	if err != nil {
		return nil, err
	}
	listing, err := source.list(tree)
	if err != nil {
		return nil, fmt.Errorf("list the tree snapshot: %w", err)
	}
	snapshot, err := parseTreeSnapshot(listing)
	if err != nil {
		return nil, err
	}
	return &treeGeneration{tree: tree, snapshot: snapshot, source: source, blobs: map[string]blobResult{}}, nil
}

// parseTreeSnapshot consumes Git's NUL-delimited listing once, so every consumer sees the
// same raw path spelling and malformed listings refuse before an identity can omit an entry.
func parseTreeSnapshot(listing string) (treeSnapshot, error) {
	snapshot := treeSnapshot{byPath: map[string]string{}}
	for _, entry := range strings.Split(listing, "\x00") {
		if entry == "" {
			continue
		}
		metadata, path, separated := strings.Cut(entry, "\t")
		if !separated {
			return treeSnapshot{}, errors.New("unparsable tree snapshot entry")
		}
		if !treeEntryMetadataRE.MatchString(metadata) {
			return treeSnapshot{}, errors.New("malformed tree snapshot metadata")
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
