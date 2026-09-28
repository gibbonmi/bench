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

func TestResolveCommitMatchesOriginalQueries(t *testing.T) {
	root := newRepo(t)
	runGit(t, root, "branch", "topic")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "tag", "-a", "annotated", "-m", "fixture")
	commit := strings.TrimSpace(runGit(t, root, "rev-parse", "HEAD"))
	for _, flags := range [][]string{nil, {"--quiet"}, {"--quiet", "--end-of-options"}} {
		for _, revision := range []string{"HEAD", "topic", commit, "annotated", "HEAD^{tree}", "HEAD:a.txt", "missing", "--help"} {
			t.Run(fmt.Sprint(flags)+"/"+revision, func(t *testing.T) {
				args := append([]string{"-C", root, "rev-parse", "--verify"}, flags...)
				want, wantErr := Output(append(args, revision+"^{commit}")...)
				got, err := ResolveCommit(root, revision, flags...)
				if got != want || fmt.Sprint(err) != fmt.Sprint(wantErr) {
					t.Fatalf("ResolveCommit = %q, %v; original query = %q, %v", got, err, want, wantErr)
				}
			})
		}
	}
	for _, root := range []string{t.TempDir(), filepath.Join(root, "missing")} {
		want, wantErr := Output("-C", root, "rev-parse", "--verify", "HEAD^{commit}")
		got, err := ResolveCommit(root, "HEAD")
		if got != want || fmt.Sprint(err) != fmt.Sprint(wantErr) {
			t.Fatalf("invalid root = %q, %v; original query = %q, %v", got, err, want, wantErr)
		}
	}
}

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

func TestDeleteBranchExactRefusesMovedRef(t *testing.T) {
	t.Parallel()
	root := newRepo(t)
	ref := "refs/heads/topic"
	runGit(t, root, "branch", "topic")
	oldOID := strings.TrimSpace(runGit(t, root, "rev-parse", ref))
	if err := os.WriteFile(filepath.Join(root, "advanced.txt"), []byte("advanced\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "advanced.txt")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "advance")
	newOID := strings.TrimSpace(runGit(t, root, "rev-parse", "HEAD"))
	runGit(t, root, "update-ref", ref, newOID, oldOID)

	if err := DeleteBranchExact(root, ref, oldOID); err == nil {
		t.Fatal("DeleteBranchExact accepted a stale OID")
	}
	if got := strings.TrimSpace(runGit(t, root, "rev-parse", ref)); got != newOID {
		t.Fatalf("moved ref = %q, want retained %q", got, newOID)
	}
}

func TestBranchDeletionPreservesSymbolicRefTarget(t *testing.T) {
	for _, prune := range []bool{false, true} {
		t.Run(fmt.Sprintf("prune=%t", prune), func(t *testing.T) {
			root := newRepo(t)
			runGit(t, root, "branch", "-M", "main")
			ref := "refs/heads/landed-alias"
			target := "refs/heads/main"
			oid := strings.TrimSpace(runGit(t, root, "rev-parse", target))
			runGit(t, root, "symbolic-ref", ref, target)
			if prune {
				count, err := PruneLandedBranches(root, nil)
				if err != nil || count != 1 {
					t.Fatalf("prune = %d, %v; want one removed alias", count, err)
				}
			} else if err := DeleteBranchExact(root, ref, oid); err != nil {
				t.Fatal(err)
			}
			if got, err := Output("-C", root, "rev-parse", "--verify", target); err != nil || got != oid {
				t.Fatalf("symbolic target = %q, %v; want retained %s", got, err, oid)
			}
			if _, err := Output("-C", root, "symbolic-ref", "-q", ref); err == nil {
				t.Fatal("deleted alias still exists")
			}
		})
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
