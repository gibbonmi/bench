package gittest

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// TreeEntry is one index entry of a fixture tree. Mode is the six-digit Git mode. Body is
// the blob content, except for a gitlink, where Body is the commit object ID.
type TreeEntry struct{ Mode, Body string }

// Trees builds Git tree objects in one fixture repository, so a test can grade a delta
// between two trees with any mode, path, or object that a worktree cannot hold.
type Trees struct {
	t    testing.TB
	Root string
}

// NewTrees returns a tree builder over a fresh repository on branch main.
func NewTrees(t testing.TB) Trees {
	t.Helper()
	return Trees{t, RepoOnBranch(t, "main")}
}

// Tree writes entries as one tree object and returns its object ID. A nil map writes the
// empty tree.
func (f Trees) Tree(entries map[string]TreeEntry) string {
	f.t.Helper()
	Output(f.t, f.Root, "read-tree", "--empty")
	paths := make([]string, 0, len(entries))
	for path := range entries {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		entry := entries[path]
		object := entry.Body
		if entry.Mode != "160000" {
			object = f.object("blob", []byte(entry.Body))
		}
		Output(f.t, f.Root, "update-index", "--add", "--cacheinfo", entry.Mode+","+object+","+path)
	}
	return Output(f.t, f.Root, "write-tree")
}

// Commits writes two commits of the empty tree, the second a child of the first, and
// returns their object IDs for use as gitlink bodies.
func (f Trees) Commits() (string, string) {
	f.t.Helper()
	tree := f.Tree(nil)
	first := Output(f.t, f.Root, "commit-tree", tree, "-m", "first")
	second := Output(f.t, f.Root, "commit-tree", tree, "-p", first, "-m", "second")
	return first, second
}

// EmptyDirectoryTree returns a tree that holds one empty subtree. Git reports no change
// between it and the empty tree, although the two object IDs differ.
func (f Trees) EmptyDirectoryTree() string {
	f.t.Helper()
	object, err := hex.DecodeString(f.Tree(nil))
	if err != nil {
		f.t.Fatal(err)
	}
	return f.object("tree", append([]byte("40000 empty\x00"), object...))
}

func (f Trees) object(kind string, body []byte) string {
	f.t.Helper()
	if err := os.WriteFile(filepath.Join(f.Root, "payload"), body, 0600); err != nil {
		f.t.Fatal(err)
	}
	return Output(f.t, f.Root, "hash-object", "-w", "-t", kind, "payload")
}
