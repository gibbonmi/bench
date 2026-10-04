package repository_test

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

const publicationDeliverable = "specs/x/spec.md"

// publicationRepo approves one delivery and registers two assignments. Only the bound
// assignment starts that delivery. The returned tree adds a production file to main.
func publicationRepo(t *testing.T) (root string, bound, unbound commitrepo.Publication, tree string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, root, publicationDeliverable, "# x\n\nStatus: staged\n")
	commitmenttest.SeedAdmission(t, root, publicationDeliverable)
	commitmenttest.Commit(t, root, "approve fixture delivery")
	boundPath := commitmenttest.Assignment(t, root, "bound")
	commitmenttest.Admit(t, boundPath, "bound", publicationDeliverable)
	unboundPath := commitmenttest.Assignment(t, root, "unbound")
	commitmenttest.Write(t, boundPath, "tool.go", "package tool\n")
	commitmenttest.Commit(t, boundPath, "production file")
	tree = gittest.Output(t, boundPath, "rev-parse", "HEAD^{tree}")
	return root, publication(t, root, boundPath), publication(t, root, unboundPath), tree
}

// publication is the frozen identity of the one assignment that owns worktree.
func publication(t *testing.T, root, worktree string) commitrepo.Publication {
	t.Helper()
	ledger, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	owners := intent.AssignmentsOwning(ledger.Assignments, worktree)
	if len(owners) != 1 {
		t.Fatalf("assignments owning %s = %+v, want one", worktree, owners)
	}
	return commitrepo.Publication{Assignment: owners[0].ID, Request: owners[0].Request, Worktree: owners[0].Worktree}
}

// The decision grades the frozen source assignment, not any assignment in the ledger.
// The ledger holds two active records, so a decision that ignores the assignment ID
// grades the wrong record for either the bound row or the borrowed row, whichever record
// comes first.
func TestAdmitPublicationFrozenIdentity(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	const mismatch = "is not active with its presented request and worktree"
	for _, row := range []struct {
		name   string
		source commitrepo.Publication
		state  intent.AssignmentState
		want   string
	}{
		{name: "bound", source: bound},
		{name: "unbound", source: unbound, want: "assignment has no current delivery binding"},
		{name: "borrowed-request-and-worktree", source: commitrepo.Publication{Assignment: unbound.Assignment, Request: bound.Request, Worktree: bound.Worktree}, want: mismatch},
		{name: "other-request", source: commitrepo.Publication{Assignment: bound.Assignment, Request: unbound.Request, Worktree: bound.Worktree}, want: mismatch},
		{name: "other-worktree", source: commitrepo.Publication{Assignment: bound.Assignment, Request: bound.Request, Worktree: unbound.Worktree}, want: mismatch},
		{name: "bound-not-active", source: bound, state: intent.StateComplete, want: mismatch},
	} {
		t.Run(row.name, func(t *testing.T) {
			if row.state != "" {
				moveAssignment(t, root, row.source.Request, row.state)
			}
			err := (commitrepo.Store{Root: root}).AdmitPublication(row.source, tree)
			if row.want == "" {
				if err != nil {
					t.Fatalf("AdmitPublication = %v, want admission", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("AdmitPublication = %v, want a refusal naming %q", err, row.want)
			}
		})
	}
}

// The verified closure of the reviewed deliverable needs the owner's authority to close
// it. The candidate is the closure alone and changes no production path, so no later
// readiness check can decide it. A bound owner or a listed scope that names the
// deliverable admits it. An unbound owner receives the start guidance, and a listed scope
// that omits the deliverable receives the scope refusal.
func TestAdmitPublicationClosureAuthority(t *testing.T) {
	for _, row := range []struct {
		name  string
		bind  bool
		scope []string
		want  string
	}{
		{name: "bound", bind: true},
		{name: "unbound", want: "assignment has no current delivery binding; run bench commitment start"},
		{name: "listed", scope: []string{closureSpec}},
		{name: "listed-without-deliverable", scope: []string{"owned.txt"}, want: "legacy continuation scope excludes \"" + closureSpec + "\""},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, head := closureRoot(t, func(t testing.TB, root string) { commitmenttest.SeedAdmission(t, root, closureSpec) })
			store := commitrepo.Store{Root: root}
			edits, err := store.Closure(head, commitrepo.Delivery{Spec: closureSpec, Source: head})
			if err != nil || len(edits) != 1 {
				t.Fatalf("Closure = %+v, %v; want the policy edit alone", edits, err)
			}
			worktree := commitmenttest.Assignment(t, root, "closer")
			if row.bind {
				commitmenttest.Admit(t, worktree, "closer", closureSpec)
			}
			commitmenttest.Write(t, worktree, commitment.PolicyPath, string(edits[0].Data))
			commitmenttest.Commit(t, worktree, "closure")
			published := publication(t, root, worktree)
			published.Source, published.Deliverable = head, closureSpec
			if row.scope != nil {
				err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
					ledger.Commitment = &intent.CommitmentState{Continuations: []intent.LegacyContinuation{{Assignment: published.Assignment, Request: published.Request, Scope: row.scope}}}
					return ledger, true, nil
				}, nil)
				if err != nil {
					t.Fatal(err)
				}
			}
			err = store.AdmitPublication(published, gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}"))
			if row.want == "" {
				if err != nil {
					t.Fatalf("AdmitPublication = %v, want the authorized closure admitted", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("AdmitPublication = %v, want a refusal naming %q", err, row.want)
			}
		})
	}
}

// continuationInput is the plan input: one adoption policy and its listed runs.
type continuationInput struct {
	commitment.Policy
	Continuations []intent.LegacyContinuation `json:"continuations,omitempty"`
}

// continuationRepo has no policy and two runs. Only the legacy run's branch holds
// owned.txt, its existing scope.
func continuationRepo(t *testing.T) (root string, legacy, other intent.Assignment, input continuationInput) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, root, "ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. old\n")
	commitmenttest.Commit(t, root, "board")
	legacyPath := commitmenttest.Assignment(t, root, "legacy")
	commitmenttest.Write(t, legacyPath, "owned.txt", "owned\n")
	commitmenttest.Commit(t, legacyPath, "legacy scope")
	otherPath := commitmenttest.Assignment(t, root, "other")
	legacyRun := publication(t, root, legacyPath)
	otherRun := publication(t, root, otherPath)
	legacy = intent.Assignment{ID: legacyRun.Assignment, Request: legacyRun.Request}
	other = intent.Assignment{ID: otherRun.Assignment, Request: otherRun.Request, Worktree: otherRun.Worktree}
	input.Policy = commitment.Policy{Version: 1, ActiveMilestone: "M1", Milestones: []commitment.Milestone{{ID: "M1", Outcomes: []commitment.Outcome{{ID: "A", Criteria: []commitment.Criterion{{ID: "A.done", Text: "The outcome is delivered."}}}}}}}
	input.Continuations = []intent.LegacyContinuation{{Assignment: legacy.ID, Request: legacy.Request, Scope: []string{"owned.txt"}}}
	return root, legacy, other, input
}

func encodeInput(t *testing.T, input continuationInput) []byte {
	t.Helper()
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// Approval records a continuation for exactly the listed run, with its assignment,
// request, and scope, in the transaction that approves the receipt. The unlisted run
// receives none, and the plan identity binds the list. Each compared element of a listed
// run refuses on its own and writes nothing.
func TestCommitmentContinuationApproval(t *testing.T) {
	root, legacy, other, input := continuationRepo(t)
	store := commitrepo.Store{Root: root}
	plan, err := store.Plan(encodeInput(t, input))
	if err != nil {
		t.Fatalf("Plan = %v, want the listed run accepted", err)
	}
	bare := input
	bare.Continuations = nil
	unlisted, err := store.Plan(encodeInput(t, bare))
	if err != nil || unlisted.ID == plan.ID {
		t.Fatalf("plan without the list = %s, %v; want an identity other than %s", unlisted.ID, err, plan.ID)
	}
	if _, err := store.Approve(plan.ID, "reviewer adoption", nil, nil); err != nil {
		t.Fatalf("Approve = %v", err)
	}
	ledger, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []intent.LegacyContinuation{{Assignment: legacy.ID, Request: legacy.Request, Scope: []string{"owned.txt"}}}
	if ledger.Commitment == nil || !reflect.DeepEqual(ledger.Commitment.Continuations, want) {
		t.Fatalf("stored continuations = %+v, want %+v", ledger.Commitment, want)
	}
	for _, continuation := range ledger.Commitment.Continuations {
		if continuation.Assignment == other.ID {
			t.Fatalf("unlisted run %s received a continuation", other.ID)
		}
	}
	// The approval that bound the list also authorizes its staged policy at commit.
	for _, path := range []string{commitment.PolicyPath, "ROADMAP.md"} {
		data, err := os.ReadFile(root + "/" + path)
		if err != nil {
			t.Fatal(err)
		}
		commitmenttest.Write(t, other.Worktree, path, string(data))
	}
	commitmenttest.Commit(t, other.Worktree, "approved adoption")
	if err := (commitrepo.Store{Root: other.Worktree}).AuthorizeCandidate(gittest.Output(t, other.Worktree, "rev-parse", "HEAD^{tree}")); err != nil {
		t.Fatalf("AuthorizeCandidate(approved policy) = %v", err)
	}

	for _, row := range []struct {
		name string
		edit func(*intent.LegacyContinuation, *continuationInput)
		want string
	}{
		{name: "unknown-run", edit: func(c *intent.LegacyContinuation, _ *continuationInput) { c.Assignment = strings.Repeat("d", 32) }, want: "is unknown"},
		{name: "other-request", edit: func(c *intent.LegacyContinuation, _ *continuationInput) { c.Request = other.Request }, want: "request does not match run"},
		{name: "empty-scope", edit: func(c *intent.LegacyContinuation, _ *continuationInput) { c.Scope = nil }, want: "invalid legacy continuation for assignment"},
		{name: "scope-outside-run", edit: func(c *intent.LegacyContinuation, _ *continuationInput) { c.Scope = []string{"absent.txt"} }, want: `scope "absent.txt" is outside run`},
		{name: "duplicate-run", edit: func(c *intent.LegacyContinuation, input *continuationInput) {
			input.Continuations = append(input.Continuations, *c)
		}, want: "invalid legacy continuation for assignment"},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, _, _, input := continuationRepo(t)
			row.edit(&input.Continuations[0], &input)
			before := commitmenttest.MilestoneState(t, root, root)
			_, err := (commitrepo.Store{Root: root}).Plan(encodeInput(t, input))
			if err == nil || !strings.Contains(err.Error(), row.want) || commitmenttest.MilestoneState(t, root, root) != before {
				t.Fatalf("Plan = %v, want a refusal naming %q that writes nothing", err, row.want)
			}
		})
	}

	// Approval checks the listed run again under the intent lock. A run whose request
	// changed after the plan refuses and records neither the approval nor a continuation.
	t.Run("run-changed-before-approval", func(t *testing.T) {
		root, _, _, input := continuationRepo(t)
		store := commitrepo.Store{Root: root}
		plan, err := store.Plan(encodeInput(t, input))
		if err != nil {
			t.Fatal(err)
		}
		rewriteAssignment(t, root, intent.RequestDigest("legacy"), func(run *intent.Assignment) { run.Request = intent.RequestDigest("changed") })
		before := commitmenttest.MilestoneState(t, root, root)
		_, err = store.Approve(plan.ID, "reviewer adoption", nil, nil)
		if err == nil || !strings.Contains(err.Error(), "request does not match run") || commitmenttest.MilestoneState(t, root, root) != before {
			t.Fatalf("Approve = %v, want a refusal that writes nothing", err)
		}
	})
}

// rewriteAssignment applies edit to the stored assignment for the request digest and
// returns the record before the edit.
func rewriteAssignment(t *testing.T, root, request string, edit func(*intent.Assignment)) intent.Assignment {
	t.Helper()
	stored, ok, err := intent.FindAssignmentByRequest(root, request)
	if err != nil || !ok {
		t.Fatalf("assignment for request %s = %t, %v", request, ok, err)
	}
	moved := stored
	edit(&moved)
	if err := intent.PutAssignment(root, moved); err != nil {
		t.Fatal(err)
	}
	return stored
}

// moveAssignment moves the assignment for request to state until the test ends, so later
// rows read the ledger the fixture wrote.
func moveAssignment(t *testing.T, root, request string, state intent.AssignmentState) {
	t.Helper()
	stored := rewriteAssignment(t, root, request, func(moved *intent.Assignment) { moved.State = state })
	t.Cleanup(func() {
		if err := intent.PutAssignment(root, stored); err != nil {
			t.Error(err)
		}
	})
}

// An admitted publication runs publish once and returns its error unchanged. A refused
// publication never runs publish.
func TestPublishAdmitted(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	failed := errors.New("ref update failed")
	for _, row := range []struct {
		name      string
		source    commitrepo.Publication
		publish   error
		published int
	}{
		{name: "admitted", source: bound, published: 1},
		{name: "publish-fails", source: bound, publish: failed, published: 1},
		{name: "refused", source: unbound},
	} {
		t.Run(row.name, func(t *testing.T) {
			published := 0
			err := (commitrepo.Store{Root: root}).PublishAdmitted(row.source, tree, func() error {
				published++
				return row.publish
			})
			if published != row.published {
				t.Fatalf("publish ran %d times, want %d (err=%v)", published, row.published, err)
			}
			switch {
			case row.published == 0 && err == nil:
				t.Fatal("refused publication returned no error")
			case row.published == 1 && err != row.publish:
				t.Fatalf("PublishAdmitted = %v, want the publish result %v", err, row.publish)
			}
		})
	}
}
