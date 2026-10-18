// Package ownershiptest owns the committed build repository and the premise collection that
// the tests of package ownership and of package staleness share. Only tests import it.
package ownershiptest

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// Cut removes the one occurrence of part from s. It fails the test unless s holds part
// exactly once, so a drifted fixture cannot pass unchanged.
func Cut(t *testing.T, s, part string) string {
	t.Helper()
	if n := strings.Count(s, part); n != 1 {
		t.Fatalf("fixture text holds %q %d times, want once", part, n)
	}
	return strings.Replace(s, part, "", 1)
}

// Repo is one committed build repository over the slug example. The selected ticket one.md
// writes the conformant fence plus Writes, and each Others ticket writes its own list. A nil
// Fence fences every write, so only an explicit Fence separates the spec from the ticket
// union. The spec declares the complete-empty inventory unless OwnInventory is set; Tail then
// closes the spec, so its bytes are part of what the fixture admits. Files are committed in
// the base; a charge source that Files omits takes a one-line body.
type Repo struct {
	Writes       []string
	Others       map[string][]string
	Files        map[string]string
	Fence        []string
	OwnInventory bool
	Tail         string
}

// Seed commits the repository as the base on main, then adds one authorized change on a
// feature branch, and returns the root, which is also the working directory.
func (r Repo) Seed(t *testing.T) string {
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
	body := preflighttest.SpecBody(slug, fence...)
	if r.OwnInventory {
		body = preflighttest.SpecBodyNoInventory(slug, fence...)
	}
	preflighttest.MustWriteFile(t, spec, body+r.Tail)
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

// Premises collects the planning inventory of the spec at source.SpecPath, as the source tip
// holds it, and answers each premise row, then each diagnostic as an unknown row at the tip,
// so an incomplete collection never reads as complete proof. A collection refusal fails the
// test.
func Premises(t *testing.T, root string, source ownership.Source) []ownership.Premise {
	t.Helper()
	content, err := git.ReadTreeFile(root, source.Tip, source.SpecPath)
	if err != nil {
		t.Fatalf("spec at the source tip: %v", err)
	}
	snapshot, err := ownership.Collect(root, source, ownership.ParseManifest(content), nil, nil)
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	rows := slices.Clone(snapshot.Premises)
	for _, diagnostic := range snapshot.Diagnostics {
		rows = append(rows, ownership.Premise{ID: diagnostic.ID, State: ownership.StateUnknown, Source: source.Tip, Reason: diagnostic.Reason})
	}
	return rows
}
