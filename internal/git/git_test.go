package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// runGit runs `git -C root <args>` and fails the test on a nonzero exit.
func runGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return string(out)
}

// newRepo initialises a repo with one commit and returns its root.
func newRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init")
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "-A")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-m", "init")
	return root
}

// TestTreeHashCleanTreeMatchesHEAD is the same-tree property the gate relies on:
// a clean working tree hashes to HEAD's tree object.
func TestTreeHashCleanTreeMatchesHEAD(t *testing.T) {
	root := newRepo(t)
	want := runGit(t, root, "rev-parse", "HEAD^{tree}")
	want = want[:len(want)-1] // strip trailing newline
	got := TreeHash(root)
	if got != want {
		t.Fatalf("clean tree: got %q, want %q", got, want)
	}
	if got == "none" {
		t.Fatal("clean repo yielded none")
	}
}

// TestTreeHashDirtyTreeDiffers checks that changing on-disk content changes the hash.
func TestTreeHashDirtyTreeDiffers(t *testing.T) {
	root := newRepo(t)
	clean := TreeHash(root)

	// Modify an existing tracked file.
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if modified := TreeHash(root); modified == clean {
		t.Fatalf("modifying a file did not change the hash (%q)", modified)
	}

	// Add a new untracked file to a fresh clean repo.
	root2 := newRepo(t)
	base := TreeHash(root2)
	if err := os.WriteFile(filepath.Join(root2, "b.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if added := TreeHash(root2); added == base {
		t.Fatalf("adding a file did not change the hash (%q)", added)
	}
}

// TestTreeHashNonRepoReturnsNone checks the failure posture for a non-repo dir.
func TestTreeHashNonRepoReturnsNone(t *testing.T) {
	dir := t.TempDir()
	if got := TreeHash(dir); got != "none" {
		t.Fatalf("non-repo dir: got %q, want \"none\"", got)
	}
}

// TestTreeHashLeavesNoStrayIndex confirms the throwaway index is external: the call
// must not drop an index file into the working tree.
func TestTreeHashLeavesNoStrayIndex(t *testing.T) {
	root := newRepo(t)
	before := listDir(t, root)
	_ = TreeHash(root)
	after := listDir(t, root)

	if _, err := os.Stat(filepath.Join(root, "index")); err == nil {
		t.Fatal("TreeHash created root/index")
	}
	if len(after) != len(before) {
		t.Fatalf("TreeHash changed the working tree: before %v, after %v", before, after)
	}
}

func TestLandedStateAggregatesAndDeduplicates(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	runGit(t, root, "remote", "add", "origin", root)
	runGit(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	linked := filepath.Join(t.TempDir(), "feature worktree")
	runGit(t, root, "worktree", "add", "-q", "-b", "feature", linked, "HEAD")
	if err := os.WriteFile(filepath.Join(linked, "a.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "checkout", "-q", "-b", "ahead")
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("ahead\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "b.txt")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ahead")
	runGit(t, root, "branch", "ahead-copy")
	for _, branch := range []string{"ahead", "ahead-copy"} {
		runGit(t, root, "config", "branch."+branch+".remote", "origin")
		runGit(t, root, "config", "branch."+branch+".merge", "refs/heads/main")
	}
	state, err := LandedState(root)
	if err != nil {
		t.Fatal(err)
	}
	if state.DirtyPaths != 1 || state.UnpushedCommits != 1 || state.UniqueBranches != 2 {
		t.Fatalf("LandedState = %#v, want dirty=1 ahead=1 unique=2", state)
	}
}

func TestLandedStateCountsDefaultBranchUpstreamAhead(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	runGit(t, root, "remote", "add", "origin", root)
	runGit(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	runGit(t, root, "config", "branch.main.remote", "origin")
	runGit(t, root, "config", "branch.main.merge", "refs/heads/main")
	if err := os.WriteFile(filepath.Join(root, "ahead.txt"), []byte("ahead\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "ahead.txt")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ahead")

	state, err := LandedState(root)
	if err != nil {
		t.Fatal(err)
	}
	if state.UnpushedCommits != 1 {
		t.Fatalf("default branch ahead count = %d, want 1", state.UnpushedCommits)
	}
}

func TestLandedStatePreservesNewlineWorktreePath(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	linked := filepath.Join(t.TempDir(), "linked\nworktree")
	runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	if err := os.WriteFile(filepath.Join(linked, "a.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	state, err := LandedState(root)
	if err != nil {
		t.Fatalf("LandedState with newline path: %v", err)
	}
	if state.DirtyPaths != 1 {
		t.Fatalf("newline worktree dirty count = %d, want 1", state.DirtyPaths)
	}
}

func TestLandedStateGitFailureIsUnknown(t *testing.T) {
	root := newRepo(t)
	t.Setenv("PATH", "")
	if _, err := LandedState(root); err == nil {
		t.Fatal("LandedState swallowed missing git")
	}
}

// TestRefResolvesAndBranchExists exercises the two guard probes and their fail-safe
// posture. They run in the process cwd (the agent's working dir), so the test chdirs
// into a fixture repo and restores.
func TestRefResolvesAndBranchExists(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "known-branch")
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	if !RefResolves("HEAD") {
		t.Error("RefResolves(HEAD) = false, want true")
	}
	if RefResolves("definitely-not-a-ref-xyz") {
		t.Error("RefResolves(bogus) = true, want false")
	}
	if !BranchExists("known-branch") {
		t.Error("BranchExists(known-branch) = false, want true")
	}
	if BranchExists("no-such-branch-xyz") {
		t.Error("BranchExists(absent) = true, want false")
	}
}

func listDir(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}
