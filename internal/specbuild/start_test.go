package specbuild

import (
	"context"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/worktree"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type abandonOwner struct{ plans, applies int }

type rejectGate struct{}

func (rejectGate) Bootstrap(context.Context, string, string, string, string) error {
	return fmt.Errorf("missing evidence")
}
func stringsSplitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

type reuseGreenGate struct{}

func (reuseGreenGate) Bootstrap(_ context.Context, root, branch, tip, expected string) error {
	return updateRef(root, "refs/bench/green/"+branch, tip, expected)
}

func TestStartCreatesRunFromExactGreenEvidence(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	status, err := service.Start(context.Background(), "build demo")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if status.State != "active" || status.Subject != git(t, root, "rev-parse", "HEAD") || status.Next != "bench spec build assign build demo" {
		t.Fatalf("Start status = %#v", status)
	}
	branch := git(t, root, "symbolic-ref", "--short", "HEAD")
	if git(t, root, "rev-parse", "refs/bench/green/"+branch) != git(t, root, "rev-parse", "HEAD") {
		t.Fatal("green marker does not name the initial tip")
	}
	refs := git(t, root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/specbuild/candidate/")
	if refs == "" || refs[len(refs)-40:] != git(t, root, "rev-parse", "HEAD") {
		t.Fatalf("candidate ref = %q", refs)
	}
}
func TestLifecycleMutatorsRefuseSharedPreconditionDriftWithoutMutation(t *testing.T) {
	operations := []struct {
		name, prerequisite string
		invoke             func(*testing.T, preconditionFixture) error
	}{
		{"start", "gate evidence", func(t *testing.T, f preconditionFixture) error {
			_, err := f.service.Start(t.Context(), "build demo")
			return err
		}},
		{"assign", "worktree owner", func(t *testing.T, f preconditionFixture) error {
			_, _, err := f.service.Assign(t.Context(), "build demo", "one.md", "next request")
			return err
		}},
		{"checkpoint", "assignment", func(t *testing.T, f preconditionFixture) error {
			_, err := f.service.Checkpoint(t.Context(), "build demo", "missing", "")
			return err
		}},
		{"integrate", "checkpoint", func(t *testing.T, f preconditionFixture) error {
			_, err := f.service.Integrate(t.Context(), "build demo", "missing")
			return err
		}},
		{"review", "review receipt", func(t *testing.T, f preconditionFixture) error {
			_, err := f.service.Review(t.Context(), "build demo", "")
			return err
		}},
	}
	conditions := []struct {
		name     string
		needsRun bool
		apply    func(*testing.T, *preconditionFixture)
	}{
		{"tracked dirt", false, func(t *testing.T, f *preconditionFixture) {
			write(t, filepath.Join(f.root, "tracked.txt"), "changed\n")
		}},
		{"untracked dirt", false, func(t *testing.T, f *preconditionFixture) {
			write(t, filepath.Join(f.root, "untracked.txt"), "changed\n")
		}},
		{"detached checkout", false, func(t *testing.T, f *preconditionFixture) { git(t, f.root, "checkout", "--detach", "-q") }},
		{"wrong branch", false, func(t *testing.T, f *preconditionFixture) { git(t, f.root, "checkout", "-qb", "other") }},
		{"working advance", false, func(t *testing.T, f *preconditionFixture) { advanceWorking(t, f.root) }},
		{"unrecognized head move", true, func(t *testing.T, f *preconditionFixture) { rewriteWorkingHead(t, f.root) }},
		{"candidate drift", true, func(t *testing.T, f *preconditionFixture) { moveCandidate(t, f) }},
		{"spec identity drift", true, func(t *testing.T, f *preconditionFixture) {
			updateRun(t, f, func(run *record) { run.SpecTip = "changed" })
		}},
		{"ownership mismatch", true, func(t *testing.T, f *preconditionFixture) {
			updateRun(t, f, func(run *record) {
				assigned := run.Assignments[f.assignmentKey]
				assigned.OwnerRequest = "changed"
				run.Assignments[f.assignmentKey] = assigned
			})
		}},
	}
	for _, operation := range operations {
		for _, condition := range conditions {
			if operation.name == "start" && (condition.needsRun || condition.name == "wrong branch" || condition.name == "working advance") {
				continue
			}
			t.Run(operation.name+"/"+condition.name, func(t *testing.T) {
				fixture := newPreconditionFixture(t, operation.name != "start")
				condition.apply(t, &fixture)
				before := snapshotPrecondition(t, fixture)
				if err := operation.invoke(t, fixture); err == nil {
					t.Fatalf("%s accepted %s", operation.name, condition.name)
				}
				after := snapshotPrecondition(t, fixture)
				if after != before {
					t.Fatalf("%s/%s mutated:\n before=%#v\n after=%#v", operation.name, condition.name, before, after)
				}
			})
		}
	}
	for _, operation := range operations {
		t.Run(operation.name+"/missing prerequisite", func(t *testing.T) {
			fixture := newPreconditionFixture(t, operation.name != "start")
			if operation.name == "start" {
				fixture.service.gate = rejectGate{}
			}
			if operation.name == "assign" {
				fixture.service.worktrees = nil
			}
			before := snapshotPrecondition(t, fixture)
			if err := operation.invoke(t, fixture); err == nil {
				t.Fatalf("%s accepted missing %s", operation.name, operation.prerequisite)
			}
			if after := snapshotPrecondition(t, fixture); after != before {
				t.Fatalf("%s missing %s mutated", operation.name, operation.prerequisite)
			}
		})
	}
}
func TestSharedPreconditionsAllowExpectedAssignmentDirt(t *testing.T) {
	fixture := newPreconditionFixture(t, true)
	write(t, filepath.Join(fixture.assignment.Path, "in-progress.txt"), "expected dirt\n")
	if _, err := fixture.service.preconditions(mutationCheckpoint, "build demo", fixture.run.Spec, &fixture.run, fixture.assignment.ID, "expected receipt"); err != nil {
		t.Fatalf("preconditions rejected expected assignment dirt: %v", err)
	}
}

type preconditionFixture struct {
	root, assignmentKey string
	service             *Service
	owner               *preconditionOwner
	run                 record
	assignment          Assignment
}

func newPreconditionFixture(t *testing.T, started bool) preconditionFixture {
	t.Helper()
	root := repo(t)
	owner := &preconditionOwner{}
	service := New(root, &countingGate{}, owner)
	fixture := preconditionFixture{root: root, service: service, owner: owner}
	if !started {
		return fixture
	}
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	assignment, _, err := service.Assign(t.Context(), "build demo", "one.md", "precondition fixture")
	if err != nil {
		t.Fatalf("Assign: %v", err)
	}
	run := loadRun(t, service)
	key, _, ok := assignmentFor(run, assignment.ID)
	if !ok {
		t.Fatal("assigned fixture has no durable row")
	}
	fixture.run, fixture.assignment, fixture.assignmentKey = run, assignment, key
	return fixture
}

type preconditionOwner struct{ calls int }

func (o *preconditionOwner) Create(_ context.Context, root, request, label, start string) (OwnedWorktree, error) {
	o.calls++
	created, err := worktree.Create(root, request, label, nil, start)
	if err != nil {
		return OwnedWorktree{}, err
	}
	return OwnedWorktree{ID: created.Assignment.ID, Path: created.Path}, nil
}

type preconditionSnapshot struct {
	state, refs, head, tree, status, worktrees string
	calls                                      int
}

func snapshotPrecondition(t *testing.T, fixture preconditionFixture) preconditionSnapshot {
	t.Helper()
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	state, err := os.ReadFile(statePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	head := git(t, fixture.root, "rev-parse", "HEAD")
	return preconditionSnapshot{state: string(state), refs: git(t, fixture.root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/"), head: head, tree: git(t, fixture.root, "rev-parse", "HEAD^{tree}"), status: git(t, fixture.root, "status", "--porcelain", "--untracked-files=all"), worktrees: git(t, fixture.root, "worktree", "list", "--porcelain"), calls: fixture.owner.calls}
}
func updateRun(t *testing.T, fixture *preconditionFixture, change func(*record)) {
	t.Helper()
	change(&fixture.run)
	saveRun(t, fixture.service, fixture.run)
}
func advanceWorking(t *testing.T, root string) {
	t.Helper()
	write(t, filepath.Join(root, "advanced.txt"), "advance\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "advance")
}
func moveCandidate(t *testing.T, fixture *preconditionFixture) {
	t.Helper()
	tree := git(t, fixture.root, "rev-parse", fixture.run.Candidate+"^{tree}")
	commit := git(t, fixture.root, "commit-tree", tree, "-p", fixture.run.Candidate, "-m", "candidate drift")
	git(t, fixture.root, "update-ref", fixture.run.Candidate, commit, fixture.run.CandidateTip)
}
func rewriteWorkingHead(t *testing.T, root string) {
	t.Helper()
	tree := git(t, root, "rev-parse", "HEAD^{tree}")
	commit := git(t, root, "commit-tree", tree, "-m", "rewritten head")
	branch := git(t, root, "symbolic-ref", "--short", "HEAD")
	git(t, root, "update-ref", "refs/heads/"+branch, commit, "HEAD")
}
func TestRecompositionErrorIsStable(t *testing.T) {
	fixture := newPreconditionFixture(t, true)
	advanceWorking(t, fixture.root)
	_, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "later")
	if err == nil || !strings.Contains(err.Error(), "bench spec build promote build demo") {
		t.Fatalf("recomposition error = %v", err)
	}
}
func TestStartResumeAndConflictsDoNotDuplicateRun(t *testing.T) {
	tests := []struct {
		name   string
		change func(*testing.T, string)
	}{
		{"branch", func(t *testing.T, root string) { git(t, root, "checkout", "-qb", "other") }},
		{"tip", func(t *testing.T, root string) {
			write(t, filepath.Join(root, "later"), "x")
			git(t, root, "add", ".")
			git(t, root, "commit", "-qm", "later")
		}},
		{"dirt", func(t *testing.T, root string) { write(t, filepath.Join(root, "dirty"), "x") }},
		{"detached", func(t *testing.T, root string) { git(t, root, "checkout", "--detach", "-q") }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := repo(t)
			service := New(root, greenGate{}, nil)
			if _, err := service.Start(context.Background(), "build demo"); err != nil {
				t.Fatal(err)
			}
			tc.change(t, root)
			if _, err := service.Start(context.Background(), "build demo"); err == nil {
				t.Fatal("Start unexpectedly accepted a changed working subject")
			}
			if got := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/candidate/"); len(stringsSplitLines(got)) != 1 {
				t.Fatalf("candidate identities = %q", got)
			}
		})
	}
	for _, point := range []string{"start/bootstrap", "start/state", "start/candidate-ref"} {
		t.Run(point, func(t *testing.T) {
			root, gate := repo(t), &countingGate{}
			service := New(root, gate, nil)
			service.fault = func(got string) error {
				if got == point {
					return errors.New("injected")
				}
				return nil
			}
			if _, err := service.Start(context.Background(), "build demo"); err == nil {
				t.Fatal("start fault did not interrupt")
			}
			service.fault = nil
			first, err := service.Start(context.Background(), "build demo")
			second, replayErr := service.Start(context.Background(), "build demo")
			if err != nil || replayErr != nil || second != first || gate.calls != 1 || len(stringsSplitLines(git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/candidate/"))) != 1 {
				t.Fatal("start replay duplicated or changed result")
			}
		})
	}
}
func TestStartWithoutEvidenceNamesOneRecoveryAndDoesNotMutate(t *testing.T) {
	root := repo(t)
	service := New(root, rejectGate{}, nil)
	if _, err := service.Start(context.Background(), "build demo"); err == nil || !strings.Contains(err.Error(), "bench gate, then retry start") {
		t.Fatalf("Start error = %v", err)
	}
	if got := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/"); got != "" {
		t.Fatalf("failed start created refs: %q", got)
	}
	status, err := service.Status("build demo")
	if err != nil || status.State != "empty" {
		t.Fatalf("status after failed start = %#v, %v", status, err)
	}
}
func TestLiteralSlugAndTicketUseOpaqueIdentities(t *testing.T) {
	root := repo(t)
	slug := "build [special]*"
	write(t, filepath.Join(root, "specs", slug, "spec.md"), "# Special\n\nStatus: staged\n")
	write(t, filepath.Join(root, "specs", slug, "tickets", "ticket [*].md"), "# Special ticket\n\nOwnership fence: internal/specbuild\n\n- [ ] [R52] literal input\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "literal spec")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), slug); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if _, _, err := service.Assign(context.Background(), slug, "ticket [*].md", "literal request"); err != nil {
		t.Fatalf("Assign: %v", err)
	}
	refs := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/candidate/")
	if strings.Contains(refs, " ") || strings.Contains(refs, "[") || strings.Contains(refs, "*") {
		t.Fatalf("candidate ref is not opaque: %q", refs)
	}
	if status, err := service.Status("build demo"); err != nil || status.State != "empty" {
		t.Fatalf("neighbor spec status = %#v, %v", status, err)
	}
}
func TestStartRefusesConflictingCandidateAndInvalidPriorState(t *testing.T) {
	t.Run("candidate compare-and-swap", func(t *testing.T) {
		root := repo(t)
		gate := &countingGate{}
		service := New(root, gate, nil)
		_, resolved, _, ok, err := spec.Resolve(root, "build demo")
		if err != nil || !ok {
			t.Fatalf("Resolve: %v, %v", ok, err)
		}
		candidate := "refs/bench/specbuild/candidate/" + digest(resolved)
		before := git(t, root, "rev-parse", "HEAD")
		write(t, filepath.Join(root, "later"), "x")
		git(t, root, "add", ".")
		git(t, root, "commit", "-qm", "later")
		git(t, root, "update-ref", candidate, before)
		if _, err := service.Start(context.Background(), "build demo"); err == nil {
			t.Fatal("Start overwrote a pre-existing candidate")
		}
		if got := git(t, root, "rev-parse", candidate); got != before {
			t.Fatalf("candidate = %s, want %s", got, before)
		}
		if gate.calls != 0 {
			t.Fatalf("gate calls = %d", gate.calls)
		}
		if got := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/green/"); got != "" {
			t.Fatalf("candidate conflict created a green marker: %q", got)
		}
		if _, found, err := service.load("build demo"); err != nil || found {
			t.Fatalf("run state after candidate conflict = found:%v err:%v", found, err)
		}
	})
	for _, tc := range []struct {
		name   string
		mutate func(*record)
	}{
		{"conflicting run identity", func(run *record) { run.Run = "another-run" }},
		{"incomplete assignments", func(run *record) { run.Assignments = nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := repo(t)
			service := New(root, greenGate{}, nil)
			if _, err := service.Start(context.Background(), "build demo"); err != nil {
				t.Fatal(err)
			}
			run := loadRun(t, service)
			candidate := run.Candidate
			before := git(t, root, "rev-parse", candidate)
			tc.mutate(&run)
			saveRun(t, service, run)
			if _, err := service.Start(context.Background(), "build demo"); err == nil {
				t.Fatal("Start accepted invalid durable state")
			}
			if got := git(t, root, "rev-parse", candidate); got != before {
				t.Fatalf("candidate changed from %s to %s", before, got)
			}
		})
	}
}

func saveRun(t *testing.T, service *Service, run record) {
	t.Helper()
	if err := service.save(run); err != nil {
		t.Fatal(err)
	}
}
