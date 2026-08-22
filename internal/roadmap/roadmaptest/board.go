// Package roadmaptest owns the split-board on-disk test fixture writer that every package
// sharing a split ROADMAP.md against a real directory uses. Only tests import this
// package. It does not import package roadmap. Package roadmap's own internal tests,
// package roadmap not roadmap_test, need this writer too, and a package roadmap import
// here would create a cycle. The index filename and row directory name are therefore
// repeated as literals rather than read from roadmap.RoadmapFile and roadmap.RoadmapDir.
package roadmaptest

import (
	"os"
	"path/filepath"
	"testing"
)

// indexFile and rowDir mirror roadmap.RoadmapFile and roadmap.RoadmapDir.
const (
	indexFile = "ROADMAP.md"
	rowDir    = "roadmap"
)

// WriteSplitBoard writes a split board into root: the index text as ROADMAP.md and one
// roadmap/<name> per entry in files. A test with no row files to write passes a nil or
// empty files map. That writes the index alone and skips the directory.
func WriteSplitBoard(t testing.TB, root, index string, files map[string]string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, indexFile), []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		return
	}
	if err := os.MkdirAll(filepath.Join(root, rowDir), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(root, rowDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
