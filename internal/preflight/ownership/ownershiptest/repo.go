// Package ownershiptest owns the committed proposal repository that the proposal command
// tests of package preflight and of package ownership share. Only tests import it. It does
// not import package preflight, because that package's own internal tests import it.
package ownershiptest

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// ProposalRepo is one committed proposal repository over the slug example. The selected
// ticket one.md writes the conformant fence plus Writes, and each Others ticket writes its
// own list. A nil Fence fences every write, so only an explicit Fence separates the spec
// from the ticket union. Tail closes the spec, so its bytes are part of what the fixture
// admits. Files are committed in the base; a charge source that Files omits takes a
// one-line body.
type ProposalRepo struct {
	Writes []string
	Others map[string][]string
	Files  map[string]string
	Fence  []string
	Tail   string
}

// Seed commits the repository as the base on main, then adds one authorized change on a
// feature branch, and returns the root, which is also the working directory.
func (r ProposalRepo) Seed(t *testing.T) string {
	t.Helper()
	const slug = "example"
	root := preflighttest.StartRepo(t)
	others := make([]string, 0, len(r.Others))
	for name := range r.Others {
		others = append(others, name)
	}
	slices.Sort(others)
	fence := r.Fence
	if fence == nil {
		for _, write := range r.Writes {
			fence = append(fence, "- `"+write+"`")
		}
		for _, name := range others {
			for _, write := range r.Others[name] {
				fence = append(fence, "- `"+write+"`")
			}
		}
	}
	spec := "specs/" + slug + "/spec.md"
	preflighttest.MustWriteFile(t, spec, preflighttest.SpecBody(slug, fence...)+r.Tail)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.WritesTicketDoc("One", append(preflighttest.FenceWrites(preflighttest.ConformantFence), r.Writes...), "PF1", "PF2"))
	for _, name := range others {
		preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/"+name, preflighttest.WritesTicketDoc(strings.TrimSuffix(name, ".md"), r.Others[name], "PF1", "PF2"))
	}
	for path, body := range r.Files {
		preflighttest.MustWriteFile(t, path, body)
	}
	for source, body := range map[string]string{chargesource.DelegateSkill: "# Delegation skill\n", chargesource.DelegateProcedure: "# Delegation procedure\n", chargesource.BuildPhase: "# Build phase\n"} {
		if _, given := r.Files[source]; !given {
			preflighttest.MustWriteFile(t, source, body)
		}
	}
	commitmenttest.SeedAdmission(t, root, spec)
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/example/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	return root
}
