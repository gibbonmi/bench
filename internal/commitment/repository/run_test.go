package repository_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestStoreReadsAPolicyOnce proves that a store with a run starts each pure policy read and
// its ledger address query once, and that a store with no run lists the policy on each read.
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
	// twice for each read. The no-run row reads the tree without the policy, which lists once.
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
	// The git stub has no mode that answers a listing with a record that git never writes,
	// so this case binds its own shim: each ls-tree answers one record with no tab.
	t.Run("malformed listing", func(t *testing.T) {
		realGit, err := exec.LookPath("git")
		if err != nil {
			t.Fatal(err)
		}
		dir := t.TempDir()
		shim := fmt.Sprintf("#!/bin/sh\ncase \"$*\" in *ls-tree*) printf 'no tab\\000'; exit 0;; esac\nexec %q \"$@\"\n", realGit)
		if err := os.WriteFile(filepath.Join(dir, "git"), []byte(shim), 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("PATH", dir)
		_, _, err = commitrepo.Store{Root: present, Run: new(git.Run)}.Policy()
		if want := `read commitment policy: git: unreadable tree entry "no tab"`; err == nil || err.Error() != want {
			t.Fatalf("policy read error = %v, want %s", err, want)
		}
	})
}
