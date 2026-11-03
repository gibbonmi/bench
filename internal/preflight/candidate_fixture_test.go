package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// candidateFixture is one candidate commit over an admitted conformant repository. HEAD still
// names the candidate's first parent, and the checkout and the index hold the candidate tree,
// as the landing leaves them before it moves the ref.
type candidateFixture struct {
	root, slug, candidate string
}

// writeCandidate stages every checkout change and writes the staged tree as a commit over
// parents. No ref moves.
func writeCandidate(t *testing.T, root, slug string, parents ...string) candidateFixture {
	t.Helper()
	preflighttest.RunGit(t, "add", "-A")
	args := []string{"commit-tree", preflighttest.RunGit(t, "write-tree"), "-m", "candidate"}
	for _, parent := range parents {
		args = append(args, "-p", preflighttest.RunGit(t, "rev-parse", parent))
	}
	return candidateFixture{root: root, slug: slug, candidate: preflighttest.RunGit(t, args...)}
}

// editFenced changes the fenced example file of slug in the checkout.
func editFenced(t *testing.T, slug string) {
	t.Helper()
	preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n\nfunc Candidate() {}\n")
}

// seedGreenCandidate is a candidate over HEAD with one fenced change and a green verdict.
func seedGreenCandidate(t *testing.T) candidateFixture {
	t.Helper()
	root, slug := preflighttest.SeedAdmitted(t)
	editFenced(t, slug)
	return writeCandidate(t, root, slug, "HEAD")
}

// seedPlanCandidate is a candidate whose spec carries an unreadable completion-plan fence
// over a HEAD whose fence is valid.
func seedPlanCandidate(t *testing.T) candidateFixture {
	t.Helper()
	root, slug := preflighttest.SeedAdmitted(t)
	spec := strings.Replace(preflighttest.SpecBody(slug), "bench-completion-plan\n{", "bench-completion-plan\n[", 1)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", spec)
	return writeCandidate(t, root, slug, "HEAD")
}

// seedMapCandidate is a candidate that deletes docs/kept.md, which the decision map at HEAD
// cites as a source.
func seedMapCandidate(t *testing.T) candidateFixture {
	t.Helper()
	root, slug := preflighttest.SeedConformant(t)
	preflighttest.RestartFeature(t, slug, func() {
		preflighttest.MustWriteFile(t, "docs/kept.md", "kept\n")
		preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/kept.md"))
	})
	preflighttest.ActiveAssignment(t, root, root)
	preflighttest.RunGit(t, "rm", "-q", "docs/kept.md")
	return writeCandidate(t, root, slug, "HEAD")
}

// seedTwoPredicateCandidate is a candidate whose coverage map holds a row with two predicates.
func seedTwoPredicateCandidate(t *testing.T) candidateFixture {
	t.Helper()
	root, slug := preflighttest.SeedAdmitted(t)
	spec := strings.Replace(preflighttest.SpecBody(slug), "| does x |", "| does x; does v |", 1)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", spec)
	return writeCandidate(t, root, slug, "HEAD")
}

// advanceMain commits path on main and returns to the feature branch, so HEAD falls behind
// the default-branch tip.
func advanceMain(t *testing.T, path string) {
	t.Helper()
	preflighttest.RunGit(t, "checkout", "-q", "main")
	preflighttest.MustWriteFile(t, path, "main\n")
	preflighttest.CommitAll(t, "main moves")
	preflighttest.RunGit(t, "checkout", "-q", "feature")
}

// seedMergeCandidate is a merge candidate whose second parent is the default-branch tip, over
// a HEAD behind that tip. Its tree keeps the feature side, so other/p.txt from main is a
// changed path only when the review base is the default-branch tip.
func seedMergeCandidate(t *testing.T) candidateFixture {
	t.Helper()
	root, slug := preflighttest.SeedAdmitted(t)
	advanceMain(t, "other/p.txt")
	editFenced(t, slug)
	return writeCandidate(t, root, slug, "HEAD", "main")
}

// seedDefaultEqualCandidate is a candidate over HEAD whose staged changes include
// other/u.txt, an unfenced path that equals the default-branch tip.
func seedDefaultEqualCandidate(t *testing.T) candidateFixture {
	t.Helper()
	root, slug := preflighttest.SeedAdmitted(t)
	advanceMain(t, "other/u.txt")
	preflighttest.RunGit(t, "checkout", "main", "--", "other/u.txt")
	editFenced(t, slug)
	return writeCandidate(t, root, slug, "HEAD")
}

// publish moves HEAD to the candidate, as the landing's compare-and-swap does.
func (f candidateFixture) publish(t *testing.T) {
	t.Helper()
	preflighttest.RunGit(t, "update-ref", "HEAD", f.candidate)
}
