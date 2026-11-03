package preflight

import (
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// seedDefaultBranchRepo starts a repository whose main branch holds one commit, and returns
// its root.
func seedDefaultBranchRepo(t *testing.T) string {
	t.Helper()
	root := preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "internal/x/y.go", "package x\n")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	return root
}

// seedUnresolvedDefaultRepo leaves the default branch unresolved: two local branches, neither
// named main, and no origin.
func seedUnresolvedDefaultRepo(t *testing.T) string {
	t.Helper()
	root := seedDefaultBranchRepo(t)
	preflighttest.RunGit(t, "branch", "-m", "main", "trunk")
	preflighttest.RunGit(t, "branch", "other")
	return root
}

// baseCurrentFacts reads the default-branch tip through git.DefaultTip, which answers the
// commit identity. An annotated tag named like the branch therefore does not make a current
// base stale. An unresolved default branch answers no name and both predicates false.
func TestBaseCurrentFactsReadsTheDefaultTip(t *testing.T) {
	t.Run("an annotated tag named like the branch", func(t *testing.T) {
		root := seedDefaultBranchRepo(t)
		preflighttest.RunGit(t, "tag", "-a", "main", "-m", "a tag named like the branch")
		preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")

		branch, resolved, current := baseCurrentFacts(root, headRef)
		if branch != "main" || !resolved || !current {
			t.Fatalf("baseCurrentFacts = (%q, %t, %t), want (main, true, true)", branch, resolved, current)
		}
	})
	t.Run("an unresolved default branch", func(t *testing.T) {
		root := seedUnresolvedDefaultRepo(t)

		branch, resolved, current := baseCurrentFacts(root, headRef)
		if branch != "" || resolved || current {
			t.Fatalf("baseCurrentFacts = (%q, %t, %t), want (\"\", false, false)", branch, resolved, current)
		}
	})
}

// An unresolved default branch answers no default-equal path and no bootstrap failure, so
// the fence rule grades every changed path.
func TestDefaultTipEqualPathsUnresolvedAnswersNone(t *testing.T) {
	root := seedUnresolvedDefaultRepo(t)
	sourceTip := preflighttest.RunGit(t, "rev-parse", "HEAD")

	equal, failure := defaultTipEqualPaths(root, headRef, sourceTip, []string{"internal/x/y.go"})
	if failure != nil {
		t.Fatalf("failure = %+v, want none for an unresolved default branch", failure)
	}
	if equal != nil {
		t.Fatalf("equal = %q, want none for an unresolved default branch", equal)
	}
}

// The HEAD grade keeps the porcelain dirty test: a staged change that the worktree reverts
// still makes the source dirty against HEAD, so no changed path keeps the default-tip
// exemption.
func TestDefaultTipEqualPathsHeadCountsTheIndex(t *testing.T) {
	root := seedDefaultBranchRepo(t)
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	sourceTip := preflighttest.RunGit(t, "rev-parse", "HEAD")
	preflighttest.MustWriteFile(t, "internal/x/y.go", "package x\n\nvar Staged = 1\n")
	preflighttest.RunGit(t, "add", "internal/x/y.go")
	preflighttest.MustWriteFile(t, "internal/x/y.go", "package x\n")

	equal, failure := defaultTipEqualPaths(root, headRef, sourceTip, []string{"internal/x/y.go"})
	if failure != nil || equal != nil {
		t.Fatalf("defaultTipEqualPaths = (%q, %+v), want none for a staged change", equal, failure)
	}
}
