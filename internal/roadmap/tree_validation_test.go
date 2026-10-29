package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestRevisionIndexReadsOnce proves that two board index reads, and two board document
// reads, through one run at one full tree id start each listing once and cat-file once.
// (Coverage row GQ33.)
func TestRevisionIndexReadsOnce(t *testing.T) {
	const board = "# Roadmap\n\n## Recommended sequence\n"
	root := gittest.RepoOnBranch(t, "main")
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	gittest.Output(t, root, "add", RoadmapFile)
	gittest.Output(t, root, "commit", "-qm", "board")
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	index := func(run *git.Run) error {
		data, present, err := RevisionIndexIn(run, root, tree)
		if err == nil && (!present || string(data) != board) {
			err = fmt.Errorf("index = %q, %t, want %q, true", data, present, board)
		}
		return err
	}
	document := func(run *git.Run) error {
		_, err := RevisionDocumentIn(run, root, tree)
		return err
	}
	// A document read lists the index and the board directory, so it starts two listings.
	cases := []struct {
		name   string
		read   func(*git.Run) error
		starts map[string]int
	}{
		{"index", index, map[string]int{"ls-tree": 1, "cat-file": 1}},
		{"document", document, map[string]int{"ls-tree": 2, "cat-file": 1}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			logPath := filepath.Join(t.TempDir(), "argv")
			gittest.StubGit(t, root, "log", logPath)
			run := new(git.Run)
			for call := range 2 {
				if err := c.read(run); err != nil {
					t.Fatalf("read %d: %v", call+1, err)
				}
			}
			for needle, want := range c.starts {
				if got := gittest.StubGitStarts(t, logPath, needle); got != want {
					t.Errorf("starts of %q = %d, want %d", needle, got, want)
				}
			}
		})
	}
}
