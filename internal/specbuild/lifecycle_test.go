package specbuild

import (
	"context"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/worktree"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestAssignBindsTicketAtCandidateInOwnedWorktree(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\nAssumptions: Go is installed\n\n- [ ] [R06-R09] create a real assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "ticket details")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	assignment, status, err := service.Assign(context.Background(), "build demo", "one.md", "request 1")
	if err != nil {
		t.Fatalf("Assign: %v", err)
	}
	if assignment.ID == "" || assignment.Path == "" || assignment.Base != status.Subject {
		t.Fatalf("assignment = %#v, status = %#v", assignment, status)
	}
	if got := git(t, assignment.Path, "rev-parse", "HEAD"); got != status.Subject {
		t.Fatalf("assignment base = %s, want %s", got, status.Subject)
	}
}
func TestStatusRefusesMalformedAndUnknownDurableState(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func([]byte) []byte
	}{
		{"malformed", func([]byte) []byte { return []byte("not JSON\n") }},
		{"unknown field", func(data []byte) []byte {
			return []byte(strings.TrimSuffix(string(data), "}\n") + `,"unknown":true}` + "\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := repo(t)
			service := New(root, greenGate{}, nil)
			if _, err := service.Start(context.Background(), "build demo"); err != nil {
				t.Fatal(err)
			}
			path, err := service.statePath("build demo")
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, tc.mutate(data), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := service.Status("build demo"); err == nil {
				t.Fatal("Status accepted invalid durable state")
			}
		})
	}
}
func TestAssignAllowsSiblingsAndScopesRequestIdentity(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R07] sibling assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatal(err)
	}
	for _, point := range []string{"assign/state", "assign/worktree"} {
		request := "fault " + point
		var prepared OwnedWorktree
		service.fault = func(got string) error {
			if got == point {
				return errors.New("injected")
			}
			return nil
		}
		if _, _, err := service.Assign(context.Background(), "build demo", "one.md", request); err == nil {
			t.Fatal("assign fault did not interrupt")
		}
		if point == "assign/worktree" {
			run := loadRun(t, service)
			requestID := digest(run.Run + "\x00" + request)
			op, found := service.operation(run, "assign", requestID)
			parts := strings.Split(op.Result, "\x00")
			_, stored := run.Assignments[requestID]
			owned, ownedFound, ownerErr := intent.FindAssignmentByRequest(root, digest(requestID))
			if !found || op.State != "prepared" || len(parts) != 2 || stored || ownerErr != nil || !ownedFound || owned.ID != parts[0] || owned.Worktree != parts[1] {
				t.Fatalf("worktree fault state = %#v, %#v, %v", op, owned, ownerErr)
			}
			prepared = OwnedWorktree{ID: parts[0], Path: parts[1]}
		}
		service.fault = nil
		assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", request)
		if err != nil {
			t.Fatal(err)
		}
		if point == "assign/worktree" && (assigned.ID != prepared.ID || assigned.Path != prepared.Path) {
			t.Fatalf("retry changed owner: %#v != %#v", assigned, prepared)
		}
		if _, _, err := service.Assign(context.Background(), "build demo", "two.md", request); err == nil {
			t.Fatal("assign changed input was accepted")
		}
	}
	first, _, err := service.Assign(context.Background(), "build demo", "one.md", "same")
	if err != nil {
		t.Fatal(err)
	}
	if replay, _, err := service.Assign(context.Background(), "build demo", "one.md", "same"); err != nil || replay.ID != first.ID || replay.Path != first.Path || replay.Base != first.Base {
		t.Fatalf("replay = %#v, %v", replay, err)
	}
	if _, _, err := service.Assign(context.Background(), "build demo", "two.md", "same"); err == nil {
		t.Fatal("request reuse for another ticket was accepted")
	}
	second, _, err := service.Assign(context.Background(), "build demo", "two.md", "other")
	if err != nil || first.Base != second.Base || first.Path == second.Path {
		t.Fatalf("siblings = %#v, %#v, %v", first, second, err)
	}
}
func TestAssignRejectsHostileTicketsBeforeWorktreeCreation(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "malformed.md"), "# Malformed\n")
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "no-fence.md"), "# No fence\n\n- [ ] [R09] missing authority\n")
	if err := os.Symlink("one.md", filepath.Join(root, "specs", "build demo", "tickets", "one-link")); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "hostile ticket fixtures")
	owner := &countingOwner{}
	service := New(root, greenGate{}, owner)
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatal(err)
	}
	fifo := filepath.Join(root, "specs", "build demo", "tickets", "one-fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "info", "exclude"), []byte("specs/build demo/tickets/one-fifo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, ticket := range []string{"missing.md", "../spec.md", "malformed.md", "no-fence.md", "one-dir", "one-fifo", "one-link"} {
		if ticket == "one-dir" {
			if err := os.Mkdir(filepath.Join(root, "specs", "build demo", "tickets", ticket), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := service.Assign(context.Background(), "build demo", ticket, ticket); err == nil {
			t.Fatalf("Assign(%q) unexpectedly succeeded", ticket)
		}
	}
	if owner.calls != 0 {
		t.Fatalf("worktree owner calls = %d", owner.calls)
	}
}
func TestStatusHasDefinitiveEmptyAndActiveProjections(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	empty, err := service.Status("build demo")
	if err != nil || empty != (Status{Slug: "build demo", State: "empty", Next: "bench spec build start build demo"}) {
		t.Fatalf("empty status = %#v, %v", empty, err)
	}
	active, err := service.Start(context.Background(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	got, err := service.Status("build demo")
	if err != nil || got != active {
		t.Fatalf("active status = %#v, %v", got, err)
	}
}
func TestStartProjectsPromotedTerminalStateAfterDescendantAdvance(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatal(err)
	}
	run := loadRun(t, service)
	run.Terminal = true
	saveRun(t, service, run)
	advanceWorking(t, root)
	status, err := service.Start(t.Context(), "build demo")
	if err != nil || status != (Status{Slug: "build demo", State: "terminal", Subject: run.CandidateTip}) {
		t.Fatalf("terminal status = %#v, %v", status, err)
	}
}
func TestStatusUsesTheGitCommonDirectory(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	active, err := service.Start(context.Background(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(t.TempDir(), "linked")
	git(t, root, "worktree", "add", "--detach", "-q", linked, "HEAD")
	fromLinked, err := New(linked, greenGate{}, nil).Status("build demo")
	if err != nil || fromLinked != active {
		t.Fatalf("linked status = %#v, %v", fromLinked, err)
	}
}
func TestReviewRecordsExactCandidateAndRetainsOnlyProjection(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatal(err)
	}
	status, err := fixture.service.Status("build demo")
	if err != nil || status.Next != "bench spec build review build demo" {
		t.Fatalf("status before review = %#v, %v", status, err)
	}
	secret := "review-body-control-token"
	receipt := reviewReceipt{Version: 1, Run: fixture.run.Run, Candidate: status.Subject, Body: secret, Axes: []reviewAxis{{Axis: "Standards", Findings: []reviewFinding{{ID: "finding-1", Disposition: "accepted"}}}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatalf("Review: %v", err)
	}
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	durable, err := os.ReadFile(statePath)
	if err != nil || strings.Contains(string(durable), secret) {
		t.Fatalf("durable review state retained receipt body: %v", err)
	}
	full, err := fixture.service.FullStatus("build demo")
	if err != nil || full.Review == nil || full.Review.Candidate != status.Subject || len(full.Review.Axes) != 3 || len(full.Review.Axes[0].Findings) != 1 || full.Review.Axes[0].Findings[0].Disposition != "accepted" || len(full.Assignments) != 1 || full.Assignments[0].Checkpoint == "" || full.Assignments[0].Integrated == "" || full.Assignments[0].Cleanup != "released" || strings.Contains(fmt.Sprintf("%#v", full), secret) {
		t.Fatalf("full status = %#v, %v", full, err)
	}
	repair, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "review repair")
	if err != nil {
		t.Fatalf("Assign repair: %v", err)
	}
	write(t, filepath.Join(repair.Path, "internal", "specbuild", "repair.go"), "package specbuild\n")
	git(t, repair.Path, "add", ".")
	git(t, repair.Path, "commit", "-qm", "review repair")
	checkpointAssignment(t, fixture.root, fixture.service, repair, []string{"internal/specbuild/repair.go"})
	if _, err := fixture.service.Integrate(t.Context(), "build demo", repair.ID); err != nil {
		t.Fatalf("Integrate repair: %v", err)
	}
	full, err = fixture.service.FullStatus("build demo")
	if err != nil || full.Review != nil || full.Next != "bench spec build review build demo" || len(full.Assignments) != 2 {
		t.Fatalf("full status after repair = %#v, %v", full, err)
	}
	receipt.Candidate, receipt.Body = full.Subject, ""
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatalf("fresh Review: %v", err)
	}
}
func TestReviewRoutesCleanAndAcceptedFindingsDifferently(t *testing.T) {
	for _, test := range []struct {
		name, disposition, want string
	}{{"clean", "", "bench spec build promote build demo"}, {"accepted", "accepted", "bench spec build assign build demo"}} {
		t.Run(test.name, func(t *testing.T) {
			fixture := checkpointedReleaseFixture(t)
			if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
				t.Fatal(err)
			}
			run := loadRun(t, fixture.service)
			candidate, tree := git(t, fixture.root, "rev-parse", run.Candidate), git(t, fixture.root, "rev-parse", "HEAD^{tree}")
			receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
			if test.disposition != "" {
				receipt.Axes[0].Findings = []reviewFinding{{ID: "repair", Disposition: test.disposition}}
			}
			if test.name == "clean" {
				fixture.service.fault = func(point string) error {
					if point == "review/state" {
						return errors.New("injected")
					}
					return nil
				}
				if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err == nil {
					t.Fatal("review fault did not interrupt")
				}
				fixture.service.fault = nil
			}
			status, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt))
			if err != nil || status.Next != test.want {
				t.Fatalf("Review status = %#v, %v", status, err)
			}
			if test.name == "clean" {
				receipt.Body = "different"
				if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err == nil {
					t.Fatal("review changed input was accepted")
				}
			}
			if git(t, fixture.root, "rev-parse", run.Candidate) != candidate || git(t, fixture.root, "rev-parse", "HEAD^{tree}") != tree {
				t.Fatal("Review mutated the project tree or candidate ref")
			}
		})
	}
}
func TestTerminalStatusIgnoresRetainedReview(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatal(err)
	}
	run := loadRun(t, service)
	run.Terminal, run.Review = true, &reviewEvidence{Candidate: run.CandidateTip, Digest: "retained", Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	saveRun(t, service, run)
	status, err := service.Status("build demo")
	if err != nil || status != (Status{Slug: "build demo", State: "terminal", Subject: run.CandidateTip}) {
		t.Fatalf("terminal status = %#v, %v", status, err)
	}
}

type countingGate struct{ calls int }

func (g *countingGate) Bootstrap(ctx context.Context, root, branch, tip, expected string) error {
	g.calls++
	return greenGate{}.Bootstrap(ctx, root, branch, tip, expected)
}

type realOwner struct{}

func (realOwner) Create(_ context.Context, root, request, label, start string) (OwnedWorktree, error) {
	created, err := worktree.Create(root, request, label, nil, start)
	if err != nil {
		return OwnedWorktree{}, err
	}
	return OwnedWorktree{ID: created.Assignment.ID, Path: created.Path}, nil
}
func (realOwner) Release(context.Context, string, string, string, ReleaseEvidence) error { return nil }

type greenGate struct{}

func (greenGate) Bootstrap(_ context.Context, root, branch, tip, expected string) error {
	cmd := exec.Command("git", "-C", root, "update-ref", "refs/bench/green/"+branch, tip, expected)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("establish green: %s", out)
	}
	return nil
}
func repo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.email", "test@example.invalid"}, {"config", "user.name", "Test"}} {
		git(t, root, args...)
	}
	write(t, filepath.Join(root, "specs", "build demo", "spec.md"), "# Build demo\n\nStatus: staged\n")
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\n\n- [ ] [R01] exact start\n")
	write(t, filepath.Join(root, "tracked.txt"), "base\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "initial")
	return root
}
func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(bytesTrimSpace(out))
}
func bytesTrimSpace(v []byte) []byte {
	for len(v) > 0 && (v[0] == ' ' || v[0] == '\n' || v[0] == '\r' || v[0] == '\t') {
		v = v[1:]
	}
	for len(v) > 0 && (v[len(v)-1] == ' ' || v[len(v)-1] == '\n' || v[len(v)-1] == '\r' || v[len(v)-1] == '\t') {
		v = v[:len(v)-1]
	}
	return v
}
func TestPromoteConflictPreservesWorkingCandidateAndState(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	write(t, filepath.Join(fixture.root, "internal", "specbuild", "checkpoint-change.go"), "package specbuild\n\nvar working = true\n")
	git(t, fixture.root, "add", ".")
	git(t, fixture.root, "commit", "-qm", "conflicting working advance")
	statePath, _ := fixture.service.statePath("build demo")
	beforeState, _ := os.ReadFile(statePath)
	working, candidate := git(t, fixture.root, "rev-parse", "HEAD"), git(t, fixture.root, "rev-parse", fixture.run.Candidate)
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("conflicting recomposition promoted")
	}
	afterState, _ := os.ReadFile(statePath)
	if git(t, fixture.root, "rev-parse", "HEAD") != working || git(t, fixture.root, "rev-parse", fixture.run.Candidate) != candidate || string(afterState) != string(beforeState) {
		t.Fatal("conflicting recomposition mutated protected state")
	}
}
func (g *promotionGate) Validate(_ context.Context, _ string, tree, evidence string) (bool, error) {
	g.validations++
	return g.accept && tree == g.tree && evidence == "owner-proof", nil
}
