package commentgap

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

type treeEntry struct{ mode, body string }
type trees struct {
	t    *testing.T
	root string
}

func newTrees(t *testing.T) trees {
	t.Helper()
	return trees{t, gittest.RepoOnBranch(t, "main")}
}

func (f trees) tree(entries map[string]treeEntry) string {
	f.t.Helper()
	gittest.Output(f.t, f.root, "read-tree", "--empty")
	paths := make([]string, 0, len(entries))
	for path := range entries {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		entry := entries[path]
		object := entry.body
		if entry.mode != "160000" {
			payload := filepath.Join(f.root, "payload")
			if err := os.WriteFile(payload, []byte(entry.body), 0600); err != nil {
				f.t.Fatal(err)
			}
			object = gittest.Output(f.t, f.root, "hash-object", "-w", "payload")
		}
		gittest.Output(f.t, f.root, "update-index", "--add", "--cacheinfo", entry.mode+","+object+","+path)
	}
	return gittest.Output(f.t, f.root, "write-tree")
}

func (f trees) commits() (string, string) {
	f.t.Helper()
	tree := f.tree(nil)
	first := gittest.Output(f.t, f.root, "commit-tree", tree, "-m", "first")
	second := gittest.Output(f.t, f.root, "commit-tree", tree, "-p", first, "-m", "second")
	return first, second
}
