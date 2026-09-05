package worktree

// This file covers the migration of the six worktree-package sites onto the
// git adapter's named readers: the decided fallbacks the migration keeps, and
// the glossary term the refusal texts adopt. (Coverage rows GR11, GR12,
// GR35.)

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLockCleanupRegistrationFallsBackForNonWorktree proves the cleanup
// registration lock keeps its repository-level lock file fallback when the
// reader fails: over a target that is no repository, the lock opens the
// pre-existing lock file at cleanupLockPath instead of propagating the
// reader's typed failure. (Coverage row GR11.)
func TestLockCleanupRegistrationFallsBackForNonWorktree(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	target := t.TempDir()
	lockPath := cleanupLockPath(repo, target)
	if err := os.WriteFile(lockPath, nil, 0o600); err != nil {
		t.Fatalf("plant the repository-level lock file: %v", err)
	}
	release, err := lockCleanupRegistration(defaultJoins(), repo, target)
	if err != nil {
		t.Fatalf("lockCleanupRegistration over a non-repository target = %v, want a release function over the repository-level lock file", err)
	}
	if release == nil {
		t.Fatal("lockCleanupRegistration returned a nil release function")
	}
	release()
}

// TestSourceMergePendingIsFalseWhenUndecided proves the merge-pending probe
// answers false when the reader fails, so the commit-and-review route stays
// correct under the undecided state. (Coverage row GR12.)
func TestSourceMergePendingIsFalseWhenUndecided(t *testing.T) {
	t.Parallel()
	source := t.TempDir()
	if sourceMergePending(source) {
		t.Fatal("sourceMergePending over a non-repository source = true, want false")
	}
}

// TestOwnershipRefusalTextsUseTheGlossaryTerm proves the classifier's
// owner-evidence refusal and the ownership marker refusal each hold the
// glossary term checkout administration directory. The fixture creates the
// owned registration before the fail-git-dir stub goes on PATH, so creation
// itself still reaches the real git. (Coverage row GR35.)
func TestOwnershipRefusalTextsUseTheGlossaryTerm(t *testing.T) {
	root, creation, _ := newOwnedAssignment(t, "admin-glossary")
	journeyStubGit(t, root, "fail-git-dir", filepath.Join(t.TempDir(), "argv"))

	want := "checkout administration directory"
	if _, err := validateOwnerMarker(root, creation.Path); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("validateOwnerMarker error = %v, want it to hold %q", err, want)
	}
	if _, err := markerPath(creation.Path); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("markerPath error = %v, want it to hold %q", err, want)
	}
}
