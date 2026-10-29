package roadmap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestRevisionIndexReadsOnce proves that two board index reads through one run at one full
// tree id start ls-tree once and cat-file once. (Coverage row GQ33.)
func TestRevisionIndexReadsOnce(t *testing.T) {
	const board = "# Roadmap\n"
	root := gittest.RepoOnBranch(t, "main")
	if err := os.WriteFile(filepath.Join(root, RoadmapFile), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	gittest.Output(t, root, "add", RoadmapFile)
	gittest.Output(t, root, "commit", "-qm", "board")
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	logPath := filepath.Join(t.TempDir(), "argv")
	gittest.StubGit(t, root, "log", logPath)
	run := new(git.Run)
	for call := range 2 {
		data, present, err := RevisionIndexIn(run, root, tree)
		if err != nil || !present || string(data) != board {
			t.Fatalf("read %d = %q, %t, %v, want %q, true, nil", call+1, data, present, err, board)
		}
	}
	for needle, want := range map[string]int{"ls-tree": 1, "cat-file": 1} {
		if got := gittest.StubGitStarts(t, logPath, needle); got != want {
			t.Errorf("starts of %q = %d, want %d", needle, got, want)
		}
	}
}
