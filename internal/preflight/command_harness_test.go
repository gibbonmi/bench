package preflight

import (
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// seedBuildFresh builds the build-mode fresh fixture. A base commit on
// main carries a bootstrap-conformant spec with no tickets/ directory at
// all, then a feature branch adds one authorized change. B1 answers
// not-applicable ticket rows over this tree.
func seedBuildFresh(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug))
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/"+slug+"/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}
