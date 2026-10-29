package repository_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestStoreReadsAPolicyOnce proves that a store with a run starts each pure policy read and
// its ledger address query once, and that a store with no run starts both on each read.
// A tree without the policy reads as absent through the run's listing, and a malformed
// listing record refuses. (Coverage rows GQ29, GQ30, GQ31, and GQ32.)
func TestStoreReadsAPolicyOnce(t *testing.T) {
	present := commitmenttest.Repo(t, policyOf("A"))
	absent := gittest.RepoOnBranch(t, "main")
	gittest.Output(t, absent, "commit", "-q", "--allow-empty", "-m", "no policy")
	policy := func(want bool) func(commitrepo.Store) error {
		return func(store commitrepo.Store) error {
			_, exists, err := store.Policy()
			if err == nil && exists != want {
				err = fmt.Errorf("policy exists = %t, want %t", exists, want)
			}
			return err
		}
	}
	ledger := func(store commitrepo.Store) error {
		_, err := store.HoldsOpenBuild()
		return err
	}
	// A present policy is listed again by its file read, so a store with no run lists it
	// twice for each read. The no-run policy row reads the tree without the policy, which lists
	// once.
	cases := []struct {
		name   string
		root   string
		run    *git.Run
		read   func(commitrepo.Store) error
		starts map[string]int
	}{
		{"policy", present, new(git.Run), policy(true), map[string]int{"ls-tree": 1, "cat-file": 1}},
		{"ledger", present, new(git.Run), ledger, map[string]int{"--git-common-dir": 1}},
		{"absent", absent, new(git.Run), policy(false), map[string]int{"ls-tree": 1, "cat-file": 0}},
		{"no run", absent, nil, policy(false), map[string]int{"ls-tree": 2}},
		{"ledger no run", present, nil, ledger, map[string]int{"--git-common-dir": 2}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			logPath := filepath.Join(t.TempDir(), "argv")
			gittest.StubGit(t, c.root, "log", logPath)
			store := commitrepo.Store{Root: c.root, Run: c.run}
			for call := range 2 {
				if err := c.read(store); err != nil {
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
	t.Run("malformed listing", func(t *testing.T) {
		gittest.StubGit(t, present, "malformed-ls-tree", filepath.Join(t.TempDir(), "argv"))
		_, _, err := commitrepo.Store{Root: present, Run: new(git.Run)}.Policy()
		want := fmt.Sprintf("read commitment policy: git: unreadable tree entry %q", gittest.MalformedTreeRecord)
		if err == nil || err.Error() != want {
			t.Fatalf("policy read error = %v, want %s", err, want)
		}
	})
}
