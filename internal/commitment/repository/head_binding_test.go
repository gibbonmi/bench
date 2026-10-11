package repository_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
)

// stageSpec is the spec that the spec-stage worktree commits and main lacks. approvedBody
// is the body that the approval binds, and editedBody is a later body that no approval names.
const (
	stageSpec    = "specs/stage/spec.md"
	approvedBody = commitmenttest.StagedBody
	editedBody   = approvedBody + "An edit after the approval.\n"
)

// specStage commits the protected milestone of outcomes A then B on main. A new assignment
// worktree then commits approvedBody at stageSpec, and plans and approves there the binding
// of that spec to A as the complete delivery of A's first row. It returns both checkouts.
func specStage(t *testing.T) (root, worktree string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	commitmenttest.SeedProtected(t, root)
	commitmenttest.Commit(t, root, "protect the rows")
	worktree = commitmenttest.Assignment(t, root, "spec-stage")
	commitmenttest.Write(t, worktree, stageSpec, approvedBody)
	commitmenttest.Commit(t, worktree, "stage the spec")
	policy, _, err := (commitrepo.Store{Root: root}).Policy()
	if err != nil {
		t.Fatal(err)
	}
	binding := commitment.DeliveryBinding{Source: stageSource(approvedBody), Obligations: []string{policy.Milestones[0].Outcomes[0].Sources[0].ID}}
	store := commitrepo.Store{Root: worktree}
	plan, err := store.Plan(commitmenttest.Bound(t, root, binding))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Approve(plan.ID, "reviewer sign-off", plan.Effects.Delayed, plan.Effects.Removed); err != nil {
		t.Fatal(err)
	}
	return root, worktree
}

// stageSource binds the spec bytes body at stageSpec.
func stageSource(body string) commitment.SourceBinding {
	return commitment.SourceBinding{ID: "spec", Path: stageSpec, Identity: commitment.Identity([]byte(body))}
}

// changedStage is the refusal of the approved spec source whose bytes became editedBody.
func changedStage() string {
	approved := stageSource(approvedBody)
	return "commitment source " + strconv.Quote(approved.ID) + " changed from " + approved.Identity + " to " + stageSource(editedBody).Identity
}

// The worktree commits its approved binding beside the staged spec that its HEAD holds and
// main lacks, and the approval check admits that candidate tree.
func TestAuthorizeCandidateAdmitsAHeadHeldBinding(t *testing.T) {
	t.Parallel()
	_, worktree := specStage(t)
	commitmenttest.Commit(t, worktree, "commit the approved binding")
	if err := (commitrepo.Store{Root: worktree}).AuthorizeCandidate(gittest.Output(t, worktree, "rev-parse", "HEAD^{tree}")); err != nil {
		t.Fatalf("AuthorizeCandidate(the approved binding beside its spec) = %v, want admission", err)
	}
}

// A candidate tree whose spec bytes changed after the approval refuses, although HEAD keeps
// the approved bytes.
func TestAuthorizeCandidateRefusesAChangedBoundSpec(t *testing.T) {
	t.Parallel()
	_, worktree := specStage(t)
	commitmenttest.Write(t, worktree, stageSpec, editedBody)
	gittest.Output(t, worktree, "add", "-A")
	err := (commitrepo.Store{Root: worktree}).AuthorizeCandidate(gittest.Output(t, worktree, "write-tree"))
	if want := changedStage(); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("AuthorizeCandidate(a changed bound spec) = %v, want a refusal naming %q", err, want)
	}
}

// A spec-less publication grades the tree that composes the worktree tip onto a moved main.
// It admits the approved binding beside its spec, and it refuses a spec whose bytes changed
// after the approval.
func TestAdmitPublicationGradesANewBindingAtTheComposedTree(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ name, body, want string }{
		{name: "approved", body: approvedBody},
		{name: "changed-after-approval", body: editedBody, want: changedStage()},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			root, worktree := specStage(t)
			commitmenttest.Write(t, worktree, stageSpec, row.body)
			commitmenttest.Commit(t, worktree, "commit the approved binding")
			landProduction(t, root)
			composed := gittest.Output(t, root, "merge-tree", "--write-tree", "main", gittest.Output(t, worktree, "rev-parse", "HEAD"))
			err := (commitrepo.Store{Root: root}).AdmitPublication(publication(t, root, worktree), composed)
			if row.want == "" {
				if err != nil {
					t.Fatalf("AdmitPublication(the composed approved binding) = %v, want admission", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("AdmitPublication(a composed changed spec) = %v, want a refusal naming %q", err, row.want)
			}
		})
	}
}
