package repository_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
)

const closureSpec = "specs/x/spec.md"

// closureRoot commits the staged fixture spec, its review record, and the policy and board
// that seed writes, and returns the repository and its head commit.
func closureRoot(t *testing.T, seed func(testing.TB, string)) (root, head string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, root, closureSpec, "# x\n\nStatus: staged\n")
	commitmenttest.Write(t, root, "reviews/x.md", "record\n")
	seed(t, root)
	commitmenttest.Commit(t, root, "approve fixture delivery")
	return root, gittest.Output(t, root, "rev-parse", "HEAD")
}

// closedEdit is one expected closure edit. A written board carries its exact bytes. The
// policy edit carries its one recorded fact instead.
type closedEdit struct {
	path, board string
	delete      bool
}

// Each delivery that closes no full spec derives its closure from the one owner. A rowless
// delivery records its fact and closes no row: with no board it edits the policy alone,
// and a board keeps its rows and projects its sequence. A tickets-only delivery closes its
// approved row. Its evidence is the folder's tree in the reviewed source.
func TestCommitmentLightClosure(t *testing.T) {
	policyEdit := closedEdit{path: commitment.PolicyPath}
	for _, row := range []struct {
		name, deliverable, binding, evidence string
		seed                                 func(testing.TB, string)
		want                                 []closedEdit
	}{
		{name: "rowless-without-board", deliverable: closureSpec, binding: "spec", evidence: ":reviews/x.md",
			seed: func(t testing.TB, root string) { commitmenttest.SeedAdmission(t, root, closureSpec) }, want: []closedEdit{policyEdit}},
		{name: "rowless-with-board", deliverable: closureSpec, binding: "spec", evidence: ":reviews/x.md",
			seed: func(t testing.TB, root string) { commitmenttest.SeedRowless(t, root, closureSpec) },
			want: []closedEdit{policyEdit, {path: "ROADMAP.md", board: commitmenttest.DeliveredIndex}}},
		{name: "tickets-only", deliverable: commitmenttest.TicketsFolder, binding: "tickets", evidence: ":" + commitmenttest.TicketsFolder,
			seed: func(t testing.TB, root string) { commitmenttest.SeedTicketsOnly(t, root, closureSpec) },
			want: []closedEdit{policyEdit, {path: "roadmap/FT1.md", delete: true}, {path: "ROADMAP.md", board: commitmenttest.DeliveredIndex}}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, head := closureRoot(t, row.seed)
			edits, err := (commitrepo.Store{Root: root}).Closure(head, commitrepo.Delivery{Spec: row.deliverable, Source: head})
			if err != nil || len(edits) != len(row.want) {
				t.Fatalf("Closure = %+v, %v; want %d edits", edits, err, len(row.want))
			}
			for i, want := range row.want {
				got := edits[i]
				if got.Path != want.path || got.Delete != want.delete || (want.board != "" && string(got.Data) != want.board) {
					t.Fatalf("edit %d = %s delete=%t %q, want %+v", i, got.Path, got.Delete, got.Data, want)
				}
			}
			policy, err := commitment.Parse(edits[0].Data)
			if err != nil || len(policy.Deliveries) != 1 {
				t.Fatalf("closed policy = %+v, %v; want one fact", policy.Deliveries, err)
			}
			evidence := gittest.Output(t, root, "rev-parse", head+row.evidence)
			if row.deliverable == commitmenttest.TicketsFolder {
				evidence = commitrepo.TreeIdentity(evidence)
			}
			if fact := policy.Deliveries[0]; fact.Outcome != commitmenttest.DeliveryOutcome || fact.Binding != row.binding || fact.Source != head || fact.Evidence != evidence {
				t.Fatalf("fact = %+v, want binding %s, source %s, evidence %s", fact, row.binding, head, evidence)
			}
		})
	}
}

// The closure refuses a tickets-only binding whose reviewed folder carries a spec file,
// because that folder is no light-path evidence. It refuses a row to close on a tree with
// no board.
func TestCommitmentLightClosureRefusal(t *testing.T) {
	for _, row := range []struct {
		name, want string
		change     func(t *testing.T, root string)
	}{
		{name: "folder-with-spec", want: "invalid checkpoint spec path", change: func(t *testing.T, root string) {
			commitmenttest.Write(t, root, commitmenttest.TicketsFolder+"/spec.md", "# t\n\nStatus: staged\n")
		}},
		{name: "row-without-board", want: "closed row FT1 appears 0 times", change: func(t *testing.T, root string) {
			gittest.Output(t, root, "rm", "-q", "ROADMAP.md")
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, _ := closureRoot(t, func(t testing.TB, root string) { commitmenttest.SeedTicketsOnly(t, root, closureSpec) })
			row.change(t, root)
			commitmenttest.Commit(t, root, row.name)
			source := gittest.Output(t, root, "rev-parse", "HEAD")
			edits, err := (commitrepo.Store{Root: root}).Closure(source, commitrepo.Delivery{Spec: commitmenttest.TicketsFolder, Source: source})
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("Closure = %+v, %v; want a refusal naming %q", edits, err, row.want)
			}
		})
	}
}

// A project with no board admits the bound rowless closure, and its candidate keeps no
// sequence: a candidate that adds a board with a sequence refuses.
func TestCommitmentRowlessAdmissionWithoutBoard(t *testing.T) {
	root, _ := closureRoot(t, func(t testing.TB, root string) { commitmenttest.SeedAdmission(t, root, closureSpec) })
	worktree := commitmenttest.Assignment(t, root, "bound")
	commitmenttest.Admit(t, worktree, "bound", closureSpec)
	commitmenttest.Write(t, worktree, "tool.go", "package tool\n")
	commitmenttest.Commit(t, worktree, "production file")
	source := gittest.Output(t, worktree, "rev-parse", "HEAD")
	edits, err := (commitrepo.Store{Root: root}).Closure(source, commitrepo.Delivery{Spec: closureSpec, Source: source})
	if err != nil || len(edits) != 1 {
		t.Fatalf("Closure = %+v, %v; want the policy edit alone", edits, err)
	}
	commitmenttest.Write(t, worktree, commitment.PolicyPath, string(edits[0].Data))
	commitmenttest.Commit(t, worktree, "closure")
	published := publication(t, root, worktree)
	published.Source, published.Deliverable = source, closureSpec
	store := commitrepo.Store{Root: root}
	if err := store.AdmitPublication(published, gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}")); err != nil {
		t.Fatalf("AdmitPublication = %v, want the rowless closure admitted", err)
	}
	commitmenttest.Write(t, worktree, "ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. other\n")
	commitmenttest.Commit(t, worktree, "add a board")
	err = store.AdmitPublication(published, gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}"))
	if err == nil || !strings.Contains(err.Error(), "protected recommended sequence changed") {
		t.Fatalf("AdmitPublication = %v, want the protected-sequence refusal", err)
	}
}
