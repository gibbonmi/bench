package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

func TestPruneLandedBranchesUsesNeutralDiscoveryFailure(t *testing.T) {
	root := newRepo(t)
	common, err := CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	id := filepath.Join(common, "worktrees", "fifo")
	if err := os.MkdirAll(id, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(id, "gitdir"), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}

	_, err = PruneLandedBranches(root, nil)
	if err == nil || !strings.Contains(err.Error(), "worktree discovery failed") || strings.Contains(err.Error(), "git worktree list") {
		t.Fatalf("worktree discovery refusal = %v", err)
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

// TestResolvedDefaultSoleMaster is the sole-local-branch fallback. A master-only
// repository has no origin/HEAD and no "main" to verify. The lone local branch is
// the only evidence of its default.
func TestResolvedDefaultSoleMaster(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "master")

	def, ok := ResolvedDefault(root)

	if !ok || def != "master" {
		t.Fatalf("ResolvedDefault = (%q, %v), want (\"master\", true)", def, ok)
	}
}

// TestResolvedDefaultNoLocalBranches covers the empty end of the sole-local-branch
// fallback. A repository with no commits has no branch to fall back to. If the code
// indexes the list before counting it, it panics instead of reporting the unresolved
// state.
func TestResolvedDefaultNoLocalBranches(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")

	def, ok := ResolvedDefault(root)

	if ok || def != "" {
		t.Fatalf("ResolvedDefault = (%q, %v), want (\"\", false)", def, ok)
	}
}

// TestResolvedDefaultUnresolvableNamesNothing pins the ok=false return as an empty name.
// No caller can put a branch this repository does not have into a message or a ref.
func TestResolvedDefaultUnresolvableNamesNothing(t *testing.T) {
	def, ok := ResolvedDefault(newTwoBranchRepo(t))

	if ok || def != "" {
		t.Fatalf("ResolvedDefault = (%q, %v), want (\"\", false)", def, ok)
	}
}
