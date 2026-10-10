package adopt

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// teamProse is the project AGENTS.md prose that a fixture writes before Link.
const teamProse = "# Team rules\n\nKeep this rule.\n"

// goldenParent holds the golden repositories of this package run. A golden repository
// outlives the test that builds it, so TestMain removes the directory after the run.
var goldenParent string

// golden is a consumer repository that one real Link builds at most once for each package
// run. When arrange is set, it prepares the repository before that Link. A test that starts
// from a linked repository takes a linkedCopy of a golden repository and does not run its
// own Link, because a real Link costs seconds of synchronized writes. The Link result does
// not depend on the repository path.
type golden struct {
	name    string
	arrange func(t *testing.T, root string)
	once    sync.Once
	root    string
}

var (
	goldenPlain = &golden{name: "plain"}
	goldenProse = &golden{name: "prose", arrange: func(t *testing.T, root string) {
		writeFixtureFile(t, filepath.Join(root, "AGENTS.md"), teamProse, 0o644)
	}}
	goldenDiverted = &golden{name: "diverted", arrange: func(t *testing.T, root string) { hookRepo(t, root, true) }}
)

func TestMain(m *testing.M) {
	parent, err := os.MkdirTemp("", "adopt-golden-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	goldenParent = parent
	code := m.Run()
	if err := os.RemoveAll(parent); err != nil {
		fmt.Fprintln(os.Stderr, err)
		code = 1
	}
	os.Exit(code)
}

// build links g in the first test that asks for it and returns its root. The test returns
// to its working directory, because consumerRepo finds the kit from there.
func (g *golden) build(t *testing.T) string {
	t.Helper()
	g.once.Do(func() {
		wd := mustGetwd(t)
		root := filepath.Join(goldenParent, g.name)
		if err := os.Mkdir(root, 0o755); err != nil {
			t.Fatal(err)
		}
		consumerRepoAt(t, root)
		if g.arrange != nil {
			g.arrange(t, root)
		}
		linkConsumer(t, root, 0)
		t.Chdir(wd)
		g.root = root
	})
	if g.root == "" {
		t.Fatalf("the %s golden repository failed to build in an earlier test", g.name)
	}
	return g.root
}

// linkedCopy returns a fresh consumer repository, made the working directory, that holds a
// copy of the golden repository g: each file with its bytes, each symlink with its target,
// and each entry with its mode.
func linkedCopy(t *testing.T, g *golden) string {
	t.Helper()
	base := g.build(t)
	root := consumerRepo(t)
	if err := os.RemoveAll(filepath.Join(root, ".git")); err != nil {
		t.Fatal(err)
	}
	if err := os.CopyFS(root, os.DirFS(base)); err != nil {
		t.Fatal(err)
	}
	err := filepath.WalkDir(base, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.Type()&fs.ModeSymlink != 0 {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		return os.Chmod(filepath.Join(root, rel), info.Mode().Perm())
	})
	if err != nil {
		t.Fatal(err)
	}
	return root
}
