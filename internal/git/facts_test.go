package git

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestTreeHashCleanTreeMatchesHEAD is the same-tree property the gate relies on. A
// clean working tree hashes to HEAD's tree object.
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

func TestAllFilesFactsExpandsWithoutChangingFacts(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init", "-q", "-b", "main")
	runGit(t, root, "config", "user.email", "t@example.com")
	runGit(t, root, "config", "user.name", "t")
	if err := os.MkdirAll(filepath.Join(root, "nested", "new"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "nested", "new", "file.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	legacy, err := Facts(root)
	if err != nil {
		t.Fatal(err)
	}
	all, err := AllFilesFacts(root)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := legacy.Changes[0].Path, "nested/"; got != want {
		t.Fatalf("Facts changed its legacy collapsed-directory result: got %q, want %q", got, want)
	}
	if got, want := all.Changes[0].Path, "nested/new/file.txt"; got != want {
		t.Fatalf("AllFilesFacts path = %q, want %q", got, want)
	}
}

func TestTreeHashExcludesIgnoredContent(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", ".gitignore")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "ignore fixture")
	clean := TreeHash(root)
	if err := os.MkdirAll(filepath.Join(root, "ignored"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored", "smuggled.txt"), []byte("ignored\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := TreeHash(root); got != clean {
		t.Fatalf("ignored content changed tree from %q to %q", clean, got)
	}
}

// TestTreeHashNonRepoReturnsNone checks the failure posture for a non-repo dir.
func TestTreeHashNonRepoReturnsNone(t *testing.T) {
	dir := t.TempDir()
	if got := TreeHash(dir); got != "none" {
		t.Fatalf("non-repo dir: got %q, want \"none\"", got)
	}
}

// TestTreeHashLeavesNoStrayIndex confirms the throwaway index is external. The call
// must not drop an index file into the working tree.
func TestTreeHashLeavesNoStrayIndex(t *testing.T) {
	root := newRepo(t)
	before := listDir(t, root)
	index := runGit(t, root, "write-tree")
	_ = TreeHash(root)
	after := listDir(t, root)

	if _, err := os.Stat(filepath.Join(root, "index")); err == nil {
		t.Fatal("TreeHash created root/index")
	}
	if len(after) != len(before) {
		t.Fatalf("TreeHash changed the working tree: before %v, after %v", before, after)
	}
	if got := runGit(t, root, "write-tree"); got != index {
		t.Fatalf("TreeHash changed the real index from %q to %q", index, got)
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

func TestLandedStateCountsDirtyPathsOnlyInNamedCheckout(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("root dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(linked, "linked.txt"), []byte("linked dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name      string
		root      string
		wantDirty int
	}{
		{name: "primary", root: root, wantDirty: 1},
		{name: "linked", root: linked, wantDirty: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state, err := LandedState(tc.root)
			if err != nil {
				t.Fatal(err)
			}
			if state.DirtyPaths != tc.wantDirty || state.UnpushedCommits != 1 || state.UniqueBranches != 2 || !reflect.DeepEqual(state.UniqueBranchNames, []string{"ahead", "ahead-copy"}) {
				t.Fatalf("LandedState = %#v, want dirty=%d ahead=1 unique=2", state, tc.wantDirty)
			}
		})
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

func TestLandedStateIgnoresDirtySiblingWithNewlinePath(t *testing.T) {
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
	if state.DirtyPaths != 0 {
		t.Fatalf("newline sibling worktree dirty count = %d, want 0", state.DirtyPaths)
	}
}

func TestLandedStateIgnoresUnreadableSibling(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	linked := filepath.Join(t.TempDir(), "missing-worktree")
	runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	if err := os.RemoveAll(linked); err != nil {
		t.Fatal(err)
	}

	state, err := LandedState(root)
	if err != nil {
		t.Fatalf("LandedState with unreadable sibling: %v", err)
	}
	if state.DirtyPaths != 0 {
		t.Fatalf("unreadable sibling dirty count = %d, want 0", state.DirtyPaths)
	}
}

func TestLandedStateGitFailureIsUnknown(t *testing.T) {
	root := newRepo(t)
	t.Setenv("PATH", "")
	if _, err := LandedState(root); err == nil {
		t.Fatal("LandedState swallowed missing git")
	}
}

// TestFactsUnresolvableDefault pins what Facts reports when no default branch resolves:
// the state, not a guess. Ahead/Behind stay zero because there is nothing to measure
// against, and the rest of the snapshot still has to arrive.
func TestFactsUnresolvableDefault(t *testing.T) {
	root := newTwoBranchRepo(t)
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	facts, err := Facts(root)

	if err != nil {
		t.Fatalf("Facts on an unresolvable default: %v", err)
	}
	if facts.DefaultResolved {
		t.Error("DefaultResolved = true, want false")
	}
	if facts.DefaultBranch != "" {
		t.Errorf("DefaultBranch = %q, want empty", facts.DefaultBranch)
	}
	if facts.Ahead != 0 || facts.Behind != 0 {
		t.Errorf("Ahead/Behind = %d/%d, want 0/0", facts.Ahead, facts.Behind)
	}
	if facts.Branch != "master" || !facts.Dirty {
		t.Errorf("rest of the snapshot lost: branch %q, dirty %v", facts.Branch, facts.Dirty)
	}
}
