package preflight

import (
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// An unresolved default branch answers no default-equal path and no bootstrap failure, so
// the fence rule grades every changed path. Two local branches, neither named main, and no
// origin leave the default branch unresolved.
func TestDefaultTipEqualPathsUnresolvedAnswersNone(t *testing.T) {
	root := preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "internal/x/y.go", "package x\n")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "branch", "-m", "main", "trunk")
	preflighttest.RunGit(t, "branch", "other")
	sourceTip := preflighttest.RunGit(t, "rev-parse", "HEAD")

	equal, failure := defaultTipEqualPaths(root, sourceTip, []string{"internal/x/y.go"})
	if failure != nil {
		t.Fatalf("failure = %+v, want none for an unresolved default branch", failure)
	}
	if equal != nil {
		t.Fatalf("equal = %q, want none for an unresolved default branch", equal)
	}
}
