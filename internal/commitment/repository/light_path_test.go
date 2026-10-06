package repository_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

// lightFolder is the light-path folder that most rows write, and lightTicket is its ticket.
const (
	lightFolder = "specs/lp"
	lightTicket = lightFolder + "/tickets/one.md"
)

// approvedFolder is the one-ticket folder that the fixture's active milestone approves.
const approvedFolder = "specs/approved"

// unpinnedRow is the board index line of FT9, a row that no outcome pins.
const unpinnedRow = "**FT9 — C**\n\n"

// lightPathRepo commits on main the SeedClosure milestone over MilestoneSpec, its review
// record, the unpinned row FT9, and a one-ticket folder that outcome B approves as a
// deliverable. A second assignment then holds the active binding and claim of the delivery
// outcome. It returns root and that bound worktree.
func lightPathRepo(t *testing.T) (root, bound string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, root, commitmenttest.MilestoneSpec, "# x\n\nStatus: staged\n")
	commitmenttest.Write(t, root, commitmenttest.MilestoneRecord, "record\n")
	commitmenttest.SeedClosure(t, root, commitmenttest.MilestoneSpec)
	commitmenttest.Write(t, root, "ROADMAP.md", strings.Replace(commitmenttest.ClosureIndex(), "## Recommended", unpinnedRow+"## Recommended", 1))
	commitmenttest.Write(t, root, "roadmap/FT9.md", unpinnedRow+"Settle the C finding.\n")
	commitmenttest.WriteLightTicket(t, root, approvedFolder+"/tickets/one.md", "approved.go")
	identity := commitmenttest.FolderIdentity(t, root, approvedFolder)
	commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
		other := &policy.Milestones[0].Outcomes[1]
		other.Deliverables = append(other.Deliverables, commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: "approved", Path: approvedFolder, Identity: identity}})
	})
	commitmenttest.Commit(t, root, "approve the light-path fixture")
	bound = commitmenttest.Assignment(t, root, "bound")
	commitmenttest.Admit(t, bound, "bound", commitmenttest.MilestoneSpec)
	return root, bound
}

// Each row commits one candidate in its own unbound assignment, or in the bound one, and
// grades that tree. An admitted row wants no refusal. Every other row names a fragment of
// its refusal.
func TestLightPathCandidate(t *testing.T) {
	root, bound := lightPathRepo(t)
	for _, row := range []struct {
		name  string
		bound bool
		steps []func(*testing.T, string)
		want  string
	}{
		// The bound assignment's active claim does not block the light path.
		{name: "beside-active-claim", steps: steps(ticket(lightTicket, "change.go"), write("change.go"))},
		{name: "new-marker", steps: steps(ticket(lightTicket, "new.go (new)"), write("new.go"))},
		{name: "directory-entry", steps: steps(ticket(lightTicket, "pkg"), write("pkg/a.go"))},
		{name: "directory-sibling", steps: steps(ticket(lightTicket, "pkg"), write("pkgx/a.go")), want: `production path "pkgx/a.go" is outside the Writes line of light-path ticket "` + lightTicket + `"`},
		// A retirement deletes its review record and adds its ADR.
		{name: "retirement", steps: steps(ticket(lightTicket, commitmenttest.MilestoneRecord, "docs/adr/9999-x.md"), remove(commitmenttest.MilestoneRecord), write("docs/adr/9999-x.md"))},
		{name: "second-ticket", steps: steps(ticket(lightTicket, "change.go"), ticket(lightFolder+"/tickets/two.md", "change.go"), write("change.go")), want: unbound},
		{name: "nested-second-ticket", steps: steps(ticket(lightTicket, "change.go"), ticket(lightFolder+"/tickets/sub/two.md", "change.go"), write("change.go")), want: unbound},
		{name: "span", steps: steps(ticket("specs/lp1/tickets/one.md", "a.go"), ticket("specs/lp2/tickets/one.md", "b.go"), write("a.go", "b.go")), want: "production paths span more than one light-path ticket; a light-path change carries one ticket"},
		// The uncovered path wins over the span.
		{name: "uncovered-beside-span", steps: steps(ticket("specs/lp1/tickets/one.md", "a.go"), ticket("specs/lp2/tickets/one.md", "b.go"), write("a.go", "b.go", "c.go")), want: `production path "c.go"`},
		{name: "approved-folder", steps: steps(write("approved.go")), want: unbound},
		{name: "no-writes-line", steps: steps(ticket(lightTicket), write("change.go")), want: unbound},
		{name: "linked-ticket", steps: steps(write("change.go"), func(t *testing.T, worktree string) {
			full := filepath.Join(worktree, filepath.FromSlash(lightTicket))
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("../../../outside.md", full); err != nil {
				t.Fatal(err)
			}
		}), want: unbound},
		{name: "policy-edit", steps: steps(ticket(lightTicket, "change.go"), write("change.go"), func(t *testing.T, worktree string) {
			commitmenttest.EditPolicy(t, worktree, func(policy *commitment.Policy) {
				policy.Milestones[0].Outcomes[0].Criteria[0].Text = "A light-path criterion."
			})
		}), want: "candidate policy has no exact approval"},
		{name: "pinned-row-deletion", steps: steps(ticket(lightTicket, "change.go"), write("change.go"), remove("roadmap/FT1.md")), want: "candidate changes protected commitment"},
		{name: "continuation-scope", steps: steps(ticket(lightTicket, "change.go", "outside.go"), write("change.go", "outside.go"), func(t *testing.T, worktree string) {
			listContinuation(t, root, worktree, "change.go")
		}), want: `legacy continuation scope excludes "outside.go"`},
		{name: "unpinned-row-removal", steps: steps(ticket(lightTicket, "change.go"), write("change.go"), remove("roadmap/FT9.md"), func(t *testing.T, worktree string) {
			commitmenttest.Write(t, worktree, "ROADMAP.md", commitmenttest.ClosureIndex())
		})},
		{name: "asset-beside-ticket", steps: steps(ticket(lightTicket, "change.go", lightFolder+"/tickets/asset.txt"), write("change.go", lightFolder+"/tickets/asset.txt"))},
		{name: "asset-alone", steps: steps(write(lightFolder + "/tickets/asset.txt")), want: unbound},
		// The path holds a space and an ESC byte.
		{name: "hostile-path", steps: steps(ticket(lightTicket, "change.go"), write("change.go", "bad \x1b.go")), want: `production path "bad \x1b.go"`},
		{name: "bound-beside-uncovering-folder", bound: true, steps: steps(ticket(lightTicket, "other.go"), write("bound.go"))},
	} {
		t.Run(row.name, func(t *testing.T) {
			worktree := bound
			if !row.bound {
				worktree = commitmenttest.Assignment(t, root, row.name)
			}
			for _, step := range row.steps {
				step(t, worktree)
			}
			commitmenttest.Commit(t, worktree, row.name)
			err := (commitrepo.Store{Root: worktree}).AuthorizeCandidate(gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}"))
			if row.want == "" {
				if err != nil {
					t.Fatalf("AuthorizeCandidate = %v, want admission", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("AuthorizeCandidate = %v, want a refusal naming %q", err, row.want)
			}
		})
	}
}

// Each row commits a light-path change in its own unbound assignment, or in the bound one,
// and grades a publication of that commit. A row that closes its deliverable grades the
// commit's tree less the folder, as the landing composes it. Every other row grades the
// commit's tree whole. A row with no deliverable is a spec-less landing.
func TestLightPathPublication(t *testing.T) {
	root, bound := lightPathRepo(t)
	for _, row := range []struct {
		name, deliverable string
		closes, bound     bool
		steps             []func(*testing.T, string)
		want, absent      string
	}{
		{name: "delivery-names-folder", deliverable: lightFolder, closes: true, steps: steps(ticket(lightTicket, "change.go"), write("change.go"))},
		{name: "delivery-uncovered", deliverable: lightFolder, closes: true, steps: steps(ticket(lightTicket, "change.go"), write("change.go", "other.go")), want: `production path "other.go" is outside the Writes line of light-path ticket "` + lightTicket + `"`},
		{name: "delivery-two-tickets", deliverable: lightFolder, closes: true, steps: steps(ticket(lightTicket, "change.go"), ticket(lightFolder+"/tickets/two.md", "change.go"), write("change.go")), want: unbound},
		{name: "spec-less", steps: steps(ticket(lightTicket, "change.go"), write("change.go")), want: `--deliverable <path>; land the light-path change with --spec "lp"`},
		{name: "spec-less-spaced-slug", steps: steps(ticket("specs/a b/tickets/one.md", "change.go"), write("change.go")), want: `--spec "a b"`},
		{name: "spec-less-uncovered", steps: steps(ticket(lightTicket, "change.go"), write("change.go", "other.go")), want: unbound, absent: "--spec"},
		{name: "bound-spec-delivery", deliverable: commitmenttest.MilestoneSpec, bound: true, steps: steps(ticket(lightTicket, "other.go"), write("bound.go"))},
	} {
		t.Run(row.name, func(t *testing.T) {
			worktree := bound
			if !row.bound {
				worktree = commitmenttest.Assignment(t, root, row.name)
			}
			for _, step := range row.steps {
				step(t, worktree)
			}
			commitmenttest.Commit(t, worktree, row.name)
			published := publication(t, root, worktree)
			published.Source, published.Deliverable = gittest.Output(t, worktree, "rev-parse", "HEAD"), row.deliverable
			if row.closes {
				if err := os.RemoveAll(filepath.Join(worktree, filepath.FromSlash(row.deliverable))); err != nil {
					t.Fatal(err)
				}
				commitmenttest.Commit(t, worktree, "close "+row.deliverable)
			}
			err := (commitrepo.Store{Root: root}).AdmitPublication(published, gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}"))
			if row.want == "" {
				if err != nil {
					t.Fatalf("AdmitPublication = %v, want admission", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.want) || (row.absent != "" && strings.Contains(err.Error(), row.absent)) {
				t.Fatalf("AdmitPublication = %v, want a refusal naming %q and not %q", err, row.want, row.absent)
			}
		})
	}
}

// unbound is the binding refusal that a candidate keeps when it is not light-path work.
const unbound = "assignment has no current delivery binding"

func steps(each ...func(*testing.T, string)) []func(*testing.T, string) { return each }

// write writes each path as a production file.
func write(paths ...string) func(*testing.T, string) {
	return func(t *testing.T, worktree string) {
		for _, path := range paths {
			commitmenttest.Write(t, worktree, path, "package fixture\n")
		}
	}
}

// ticket writes one light-path ticket at path whose Writes line lists writes.
func ticket(path string, writes ...string) func(*testing.T, string) {
	return func(t *testing.T, worktree string) { commitmenttest.WriteLightTicket(t, worktree, path, writes...) }
}

func remove(paths ...string) func(*testing.T, string) {
	return func(t *testing.T, worktree string) {
		for _, path := range paths {
			if err := os.Remove(filepath.Join(worktree, filepath.FromSlash(path))); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// listContinuation lists the legacy continuation of the assignment that owns worktree with
// scope. The fixture's other records stay, so the bound claim holds.
func listContinuation(t *testing.T, root, worktree string, scope ...string) {
	t.Helper()
	run := publication(t, root, worktree)
	err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		state := *ledger.Commitment
		state.Continuations = append(state.Continuations, intent.LegacyContinuation{Assignment: run.Assignment, Request: run.Request, Scope: scope})
		ledger.Commitment = &state
		return ledger, true, nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
