package roadmap

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// indexTree is the index-only tree a test drives when the row files are not the
// subject: parsed index bytes and no roadmap/ listing at all.
func indexTree(index string) Tree {
	state := bounds.StateParsed
	if index == "" {
		state = bounds.StateEmpty
	}
	return Tree{Index: bounds.Classified{State: state, Data: []byte(index)}}
}

// splitTree is the split board a test drives: parsed index bytes plus one parsed
// row file per named basename, in the directory order os.ReadDir would report.
func splitTree(index string, files map[string]string) Tree {
	tree := indexTree(index)
	if files == nil {
		return tree
	}
	tree.DirState = bounds.StateParsed
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		state := bounds.StateParsed
		if files[name] == "" {
			state = bounds.StateEmpty
		}
		tree.Files = append(tree.Files, RowFile{Name: name, State: state, Data: []byte(files[name])})
		tree.DirBytes += len(files[name])
	}
	return tree
}

// writeSplitBoard writes a split board into root: the index text as ROADMAP.md and one
// roadmap/<name> per entry. It is the on-disk twin of splitTree, for the tests that
// drive a command rather than the parse.
func writeSplitBoard(t *testing.T, root, index string, files map[string]string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		return
	}
	if err := os.MkdirAll(filepath.Join(root, RoadmapDir), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(root, RoadmapDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// writeBoard writes the board built from one heading-and-body pair per row.
func writeBoard(t *testing.T, root string, rows ...[2]string) {
	t.Helper()
	index, files := board(rows...)
	writeSplitBoard(t, root, index, files)
}

// board renders an index and its row files from one row heading and body, the shape
// most fixtures need: `**FT1 — one.**` in ROADMAP.md and heading plus body in
// roadmap/FT1.md.
func board(rows ...[2]string) (string, map[string]string) {
	var index strings.Builder
	files := map[string]string{}
	for _, row := range rows {
		heading, body := row[0], row[1]
		index.WriteString(heading + "\n\n")
		id := heading[2:strings.Index(heading, " ")]
		content := heading + "\n"
		if body != "" {
			content += body
		}
		files[id+".md"] = content
	}
	return index.String(), files
}
