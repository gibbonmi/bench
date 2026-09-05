package worktree

// This file covers the migration of the six worktree-package sites onto the
// git adapter's named readers: the decided fallbacks the migration keeps, and
// the glossary term the refusal texts adopt. (Coverage rows GR11, GR12,
// GR35.) It also covers LeaseFile's migration onto the file reader, over a
// linked worktree and under the failing file-query stub. (Coverage rows GR18,
// GR24, GR25.)

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
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

// mustAdminPath reads the git administration path reader for name and fails the
// test on its typed error. The lifecycle and reauthorize fixtures share this call
// so the reader's error handling has one source. (Coverage row GR18.)
func mustAdminPath(t testing.TB, path, name string) string {
	t.Helper()
	admin, err := git.AdminPath(path, name)
	if err != nil {
		t.Fatal(err)
	}
	return admin
}

// TestLeaseFileMatchesIndependentRevParse proves LeaseFile answers the same
// path an independent `rev-parse --git-path bench-lease` reaches, over a
// linked worktree. (Coverage row GR18.)
func TestLeaseFileMatchesIndependentRevParse(t *testing.T) {
	t.Parallel()
	_, creation, _ := newOwnedAssignment(t, "lease-file-matches")

	want := gitOutput(t, creation.Path, "rev-parse", "--path-format=absolute", "--git-path", git.BenchLeaseFilename)

	got, err := LeaseFile(creation.Path)
	if err != nil {
		t.Fatalf("LeaseFile(%q) = %v, want a resolved path", creation.Path, err)
	}
	if got != want {
		t.Fatalf("LeaseFile(%q) = %q, want the independently resolved %q", creation.Path, got, want)
	}
}

// TestLeaseFileRefusesUnresolvedAnswer proves LeaseFile propagates the git
// administration path reader's typed failure instead of guessing a path under
// the worktree root. The fail-git-path stub refuses every file query.
// (Coverage rows GR24, GR25.)
func TestLeaseFileRefusesUnresolvedAnswer(t *testing.T) {
	root, creation, _ := newOwnedAssignment(t, "lease-file-unresolved")
	journeyStubGit(t, root, "fail-git-path", filepath.Join(t.TempDir(), "argv"))

	_, err := LeaseFile(creation.Path)
	var resolution *git.ResolutionError
	if !errors.As(err, &resolution) {
		t.Fatalf("LeaseFile under the fail-git-path stub = %v, want a *git.ResolutionError", err)
	}
}
