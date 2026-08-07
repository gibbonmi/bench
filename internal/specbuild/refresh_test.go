package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// refreshFixture is one stranded-consumer scenario: a consumer assignment with an
// in-fence partial edit whose repro is blocked by a defect in a disjoint fence,
// and a repair ticket staged mid-run, recomposed, assigned, checkpointed, and
// integrated through the public lifecycle.
type refreshFixture struct {
	root                    string
	service                 *Service
	consumer, repair        Assignment
	dirtyPath, dirtyContent string
	refusedBeforeRecompose  error
}

const consumerRequest = "consumer request"

func newRefreshFixture(t *testing.T) refreshFixture {
	t.Helper()
	fixture := newCleanRefreshFixture(t)
	fixture.dirtyPath, fixture.dirtyContent = "internal/specbuild/consumer-partial.go", "package specbuild\n\nvar consumerPartial = true\n"
	write(t, filepath.Join(fixture.consumer.Path, fixture.dirtyPath), fixture.dirtyContent)
	return fixture
}

// newCleanRefreshFixture is the same scenario without the partial edit: the
// consumer worktree sits byte-clean at its base, so the preservation patch the
// refresh carries has zero bytes.
func newCleanRefreshFixture(t *testing.T) refreshFixture {
	t.Helper()
	return newCleanRefreshFixtureWithRepair(t, "# Repair landing\n\nOwnership fence: internal/landing\n\n- [ ] [R90] repair the prerequisite defect\n")
}

// newCleanRefreshFixtureWithRepair runs the same scenario with a caller-shaped
// repair ticket, for tests probing what the repair's own metadata declares.
func newCleanRefreshFixtureWithRepair(t *testing.T, repairBody string) refreshFixture {
	t.Helper()
	root := repo(t)
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	consumer, _, err := service.Assign(t.Context(), "build demo", "one.md", consumerRequest)
	if err != nil {
		t.Fatalf("Assign consumer: %v", err)
	}
	fixture := refreshFixture{root: root, service: service, consumer: consumer}
	// The repair ticket lands mid-run, exactly the insertion that used to strand
	// the consumer: the ticket commit moves the working tip, so assign refuses
	// until promote recomposes the run onto it.
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "repair.md"), repairBody)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "stage repair ticket")
	_, _, fixture.refusedBeforeRecompose = service.Assign(t.Context(), "build demo", "repair.md", "repair request")
	if _, err := service.Promote(t.Context(), "build demo"); err != nil {
		t.Fatalf("recomposing Promote: %v", err)
	}
	fixture.repair, _, err = service.Assign(t.Context(), "build demo", "repair.md", "repair request")
	if err != nil {
		t.Fatalf("Assign repair: %v", err)
	}
	write(t, filepath.Join(fixture.repair.Path, "internal", "landing", "fix.go"), "package landing\n\nvar repaired = true\n")
	checkpointAssignment(t, root, service, fixture.repair, []string{"internal/landing/fix.go"})
	if _, err := service.Integrate(t.Context(), "build demo", fixture.repair.ID); err != nil {
		t.Fatalf("Integrate repair: %v", err)
	}
	return fixture
}

// debugReceiptFor builds the receipt the stranded consumer's repro earns; tests
// mutate the result to forge each hostile variant rather than restating the shape.
func debugReceiptFor(t *testing.T, fixture refreshFixture) debugReceipt {
	t.Helper()
	receipt := debugReceipt{
		Version: debugReceiptVersion, Run: loadRun(t, fixture.service).Run, Assignment: fixture.consumer.ID, Base: fixture.consumer.Base,
		Repro:         debugRepro{Command: "go test ./internal/contract/runtime -run TestRuntimeCommitContracts", Exit: 1, OutputDigest: digest("deterministic red"), Produced: time.Now().UTC().Format(time.RFC3339Nano)},
		Cause:         "prerequisite defect in the landing fence discards gate output",
		RequiredFence: []string{"internal/landing"},
		Resumable:     true,
	}
	if fixture.dirtyPath != "" {
		receipt.DirtyPaths = []string{fixture.dirtyPath}
	}
	return receipt
}

func writeDebugReceipt(t *testing.T, receipt debugReceipt) string {
	t.Helper()
	data, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "debug.json")
	write(t, path, string(data)+"\n")
	return path
}

func refreshConsumer(t *testing.T, fixture refreshFixture) Assignment {
	t.Helper()
	refreshed, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, debugReceiptFor(t, fixture)))
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	return refreshed
}

func TestRefreshRepairTraceThroughPublicLifecycle(t *testing.T) {
	fixture := newRefreshFixture(t)
	if !errors.Is(fixture.refusedBeforeRecompose, errRecompose) {
		t.Fatalf("mid-run ticket insertion refusal = %v, want recomposition route", fixture.refusedBeforeRecompose)
	}
	run := loadRun(t, fixture.service)
	refreshed := refreshConsumer(t, fixture)
	if refreshed.Base != run.CandidateTip {
		t.Fatalf("refreshed base = %s, want repaired candidate %s", refreshed.Base, run.CandidateTip)
	}
	if got := git(t, fixture.consumer.Path, "rev-parse", "HEAD"); got != run.CandidateTip {
		t.Fatalf("worktree HEAD = %s, want repaired candidate %s", got, run.CandidateTip)
	}
	surviving, err := os.ReadFile(filepath.Join(fixture.consumer.Path, fixture.dirtyPath))
	if err != nil || string(surviving) != fixture.dirtyContent {
		t.Fatalf("in-fence bytes did not survive: %q, %v", surviving, err)
	}
	if got := git(t, fixture.root, "cat-file", "-p", refreshIdentity(run.Run, fixture.consumer.ID)+":"+fixture.dirtyPath); got != strings.TrimSpace(fixture.dirtyContent) {
		t.Fatalf("preservation ref payload = %q", got)
	}
	// The resumed delegate completes the ticket against the repaired base; the
	// ordinary checkpoint and integrate path composes the exact candidate.
	write(t, filepath.Join(fixture.consumer.Path, "internal", "specbuild", "consumer-final.go"), "package specbuild\n\nvar consumerFinal = true\n")
	checkpointAssignment(t, fixture.root, fixture.service, refreshed, []string{fixture.dirtyPath, "internal/specbuild/consumer-final.go"})
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.consumer.ID); err != nil {
		t.Fatalf("Integrate consumer: %v", err)
	}
	run = loadRun(t, fixture.service)
	for _, path := range []string{"internal/landing/fix.go", fixture.dirtyPath, "internal/specbuild/consumer-final.go"} {
		git(t, fixture.root, "cat-file", "-e", run.CandidateTip+":"+path)
	}
	fixture.service.gate = &promotionGate{accept: true}
	review := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, review)); err != nil {
		t.Fatalf("Review: %v", err)
	}
	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.State != "terminal" {
		t.Fatalf("Promote = %#v, %v", status, err)
	}
	promoted := git(t, fixture.root, "rev-parse", "HEAD")
	for _, path := range []string{"internal/landing/fix.go", fixture.dirtyPath, "internal/specbuild/consumer-final.go"} {
		git(t, fixture.root, "cat-file", "-e", promoted+":"+path)
	}
}

func TestRefreshAdvancesACleanAssignmentOntoTheRepairedCandidate(t *testing.T) {
	fixture := newCleanRefreshFixture(t)
	run := loadRun(t, fixture.service)
	refreshed := refreshConsumer(t, fixture)
	if refreshed.Base != run.CandidateTip {
		t.Fatalf("refreshed base = %s, want repaired candidate %s", refreshed.Base, run.CandidateTip)
	}
	if got := git(t, fixture.consumer.Path, "rev-parse", "HEAD"); got != run.CandidateTip {
		t.Fatalf("worktree HEAD = %s, want repaired candidate %s", got, run.CandidateTip)
	}
	if status := git(t, fixture.consumer.Path, "status", "--porcelain", "--untracked-files=all"); status != "" {
		t.Fatalf("refreshed worktree is not byte-clean: %q", status)
	}
}

func TestRefreshRefusesForgedAndMissingDebugReceipts(t *testing.T) {
	fixture := newRefreshFixture(t)
	sound := debugReceiptFor(t, fixture)
	for _, test := range []struct {
		name string
		path func(*testing.T) string
	}{
		{"missing file", func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }},
		{"relative path", func(*testing.T) string { return "debug.json" }},
		{"foreign run", func(t *testing.T) string { r := sound; r.Run = digest("other"); return writeDebugReceipt(t, r) }},
		{"foreign assignment", func(t *testing.T) string {
			r := sound
			r.Assignment = fixture.repair.ID
			return writeDebugReceipt(t, r)
		}},
		{"stale base", func(t *testing.T) string {
			r := sound
			r.Base = loadRun(t, fixture.service).CandidateTip
			return writeDebugReceipt(t, r)
		}},
		{"green repro", func(t *testing.T) string { r := sound; r.Repro.Exit = 0; return writeDebugReceipt(t, r) }},
		{"receipt inside the worktree", func(t *testing.T) string {
			data, err := json.Marshal(sound)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(fixture.consumer.Path, "debug.json")
			write(t, path, string(data)+"\n")
			return path
		}},
		{"fence entirely inside the assignment", func(t *testing.T) string {
			r := sound
			r.RequiredFence = []string{"internal/specbuild/deeper"}
			return writeDebugReceipt(t, r)
		}},
		{"dirty-path claim mismatch", func(t *testing.T) string {
			r := sound
			r.DirtyPaths = []string{"internal/specbuild/other.go"}
			return writeDebugReceipt(t, r)
		}},
		{"produced before assignment", func(t *testing.T) string {
			r := sound
			r.Repro.Produced = time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339Nano)
			return writeDebugReceipt(t, r)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := test.path(t)
			before := refreshSnapshotFor(t, fixture)
			if _, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, path); err == nil {
				t.Fatal("Refresh accepted a forged or missing receipt")
			}
			if after := refreshSnapshotFor(t, fixture); after != before {
				t.Fatalf("refusal mutated state: before=%#v after=%#v", before, after)
			}
		})
	}
}

// refreshSnapshot pins everything a refresh refusal must leave intact; the
// receipt file layout differs per case, so unlike checkpointSnapshot it also
// covers the assignment's recorded base.
type refreshSnapshot struct {
	state, refs, worktreeHead, worktreeStatus, recordedBase string
}

func refreshSnapshotFor(t *testing.T, fixture refreshFixture) refreshSnapshot {
	t.Helper()
	path, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	state, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, stored, ok := assignmentFor(loadRun(t, fixture.service), fixture.consumer.ID)
	if !ok {
		t.Fatal("missing consumer assignment")
	}
	return refreshSnapshot{
		state: string(state), refs: git(t, fixture.root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/"),
		worktreeHead: git(t, fixture.consumer.Path, "rev-parse", "HEAD"), worktreeStatus: git(t, fixture.consumer.Path, "status", "--porcelain", "--untracked-files=all"),
		recordedBase: stored.Base,
	}
}

func TestRefreshRefusesOutOfFencePayloadAndCheckpointedAssignment(t *testing.T) {
	t.Run("payload outside the fence", func(t *testing.T) {
		fixture := newRefreshFixture(t)
		write(t, filepath.Join(fixture.consumer.Path, "internal", "landing", "smuggled.go"), "package landing\n")
		receipt := debugReceiptFor(t, fixture)
		receipt.DirtyPaths = append(receipt.DirtyPaths, "internal/landing/smuggled.go")
		if _, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, receipt)); err == nil || !strings.Contains(err.Error(), "leaves the ownership fence") {
			t.Fatalf("out-of-fence payload = %v", err)
		}
	})
	t.Run("checkpointed assignment", func(t *testing.T) {
		fixture := newRefreshFixture(t)
		receiptPath := writeDebugReceipt(t, debugReceiptFor(t, fixture))
		refreshed := refreshConsumer(t, fixture)
		checkpointAssignment(t, fixture.root, fixture.service, refreshed, []string{fixture.dirtyPath})
		if _, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, receiptPath); err == nil || !strings.Contains(err.Error(), "uncheckpointed") {
			t.Fatalf("checkpointed refresh = %v", err)
		}
	})
	t.Run("no candidate advance", func(t *testing.T) {
		root := repo(t)
		service := New(root, greenGate{}, realOwner{})
		if _, err := service.Start(t.Context(), "build demo"); err != nil {
			t.Fatal(err)
		}
		consumer, _, err := service.Assign(t.Context(), "build demo", "one.md", consumerRequest)
		if err != nil {
			t.Fatal(err)
		}
		fixture := refreshFixture{root: root, service: service, consumer: consumer, dirtyPath: "internal/specbuild/consumer-partial.go", dirtyContent: "package specbuild\n"}
		write(t, filepath.Join(consumer.Path, fixture.dirtyPath), fixture.dirtyContent)
		if _, _, err := service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, debugReceiptFor(t, fixture))); err == nil || !strings.Contains(err.Error(), "no candidate advance") {
			t.Fatalf("no-advance refresh = %v", err)
		}
	})
}

func TestRefreshRefusesConflictWithoutTouchingTheWorktree(t *testing.T) {
	fixture := newRefreshFixture(t)
	// A sibling ticket sharing the consumer's fence lands a different byte
	// content at the consumer's dirty path, so the preserved patch cannot apply.
	write(t, filepath.Join(fixture.root, "specs", "build demo", "tickets", "sibling.md"), "# Sibling\n\nOwnership fence: internal/specbuild\n\n- [ ] [R91] conflicting sibling\n")
	git(t, fixture.root, "add", ".")
	git(t, fixture.root, "commit", "-qm", "stage sibling ticket")
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err != nil {
		t.Fatalf("recomposing Promote: %v", err)
	}
	sibling, _, err := fixture.service.Assign(t.Context(), "build demo", "sibling.md", "sibling request")
	if err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(sibling.Path, fixture.dirtyPath), "package specbuild\n\nvar conflicting = true\n")
	checkpointAssignment(t, fixture.root, fixture.service, sibling, []string{fixture.dirtyPath})
	if _, err := fixture.service.Integrate(t.Context(), "build demo", sibling.ID); err != nil {
		t.Fatal(err)
	}
	before := refreshSnapshotFor(t, fixture)
	if _, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, debugReceiptFor(t, fixture))); err == nil || !strings.Contains(err.Error(), "conflicts with the repaired candidate") {
		t.Fatalf("conflicting refresh = %v", err)
	}
	after := refreshSnapshotFor(t, fixture)
	// The preservation ref is durable evidence and may exist; everything the
	// delegate owns — checkout, dirty bytes, recorded base — must be untouched.
	if after.worktreeHead != before.worktreeHead || after.worktreeStatus != before.worktreeStatus || after.recordedBase != before.recordedBase {
		t.Fatalf("conflict refusal mutated the assignment: before=%#v after=%#v", before, after)
	}
}

// refMovingRunner advances the candidate ref the moment the refresh replay
// writes its prospective tree, driving the interior compare-and-swap window.
type refMovingRunner struct {
	inner     Runner
	root, ref string
	moved     bool
}

func (r *refMovingRunner) Output(ctx context.Context, program string, args ...string) (string, error) {
	return r.Run(ctx, Command{Program: program, Args: args})
}

func (r *refMovingRunner) Run(ctx context.Context, command Command) (string, error) {
	for _, arg := range command.Args {
		if arg == "write-tree" && !r.moved {
			r.moved = true
			current, err := refValue(r.root, r.ref)
			if err != nil {
				return "", err
			}
			tree, err := exec.Command("git", "-C", r.root, "rev-parse", current+"^{tree}").Output()
			if err != nil {
				return "", err
			}
			drift, err := exec.Command("git", "-C", r.root, "commit-tree", strings.TrimSpace(string(tree)), "-p", current, "-m", "candidate drift").Output()
			if err != nil {
				return "", err
			}
			if err := updateRef(r.root, r.ref, strings.TrimSpace(string(drift)), current); err != nil {
				return "", err
			}
		}
	}
	return r.inner.Run(ctx, command)
}

func TestRefreshRefusesCandidateMovementMidFlight(t *testing.T) {
	fixture := newRefreshFixture(t)
	run := loadRun(t, fixture.service)
	fixture.service.runner = &refMovingRunner{inner: processRunner{}, root: fixture.root, ref: run.Candidate}
	before := refreshSnapshotFor(t, fixture)
	if _, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, debugReceiptFor(t, fixture))); err == nil || !strings.Contains(err.Error(), "candidate moved during refresh") {
		t.Fatalf("mid-flight movement = %v", err)
	}
	after := refreshSnapshotFor(t, fixture)
	if after.worktreeHead != before.worktreeHead || after.worktreeStatus != before.worktreeStatus || after.recordedBase != before.recordedBase {
		t.Fatalf("movement refusal mutated the assignment: before=%#v after=%#v", before, after)
	}
}

func TestRefreshInterruptionConvergesOnReEntry(t *testing.T) {
	for _, point := range []string{"refresh/preserve", "refresh/worktree", "refresh/state"} {
		t.Run(point, func(t *testing.T) {
			fixture := newRefreshFixture(t)
			receiptPath := writeDebugReceipt(t, debugReceiptFor(t, fixture))
			fixture.service.fault = func(got string) error {
				if got == point {
					return errors.New("injected")
				}
				return nil
			}
			if _, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, receiptPath); err == nil {
				t.Fatal("fault did not interrupt refresh")
			}
			fixture.service.fault = nil
			refreshed, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, receiptPath)
			if err != nil {
				t.Fatalf("re-entry: %v", err)
			}
			run := loadRun(t, fixture.service)
			if refreshed.Base != run.CandidateTip {
				t.Fatalf("converged base = %s, want %s", refreshed.Base, run.CandidateTip)
			}
			surviving, err := os.ReadFile(filepath.Join(fixture.consumer.Path, fixture.dirtyPath))
			if err != nil || string(surviving) != fixture.dirtyContent {
				t.Fatalf("in-fence bytes did not survive interruption: %q, %v", surviving, err)
			}
			replayed, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, receiptPath)
			if err != nil || replayed.Base != run.CandidateTip {
				t.Fatalf("idempotent replay = %#v, %v", replayed, err)
			}
		})
	}
}

func TestAbandonRemainsThePositiveControlForAStrandedRun(t *testing.T) {
	fixture := newRefreshFixture(t)
	fixture.service.worktrees = &abandonOwner{}
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon plan: %v", err)
	}
	if len(plan.Worktrees) != 1 || plan.Worktrees[0].ID != fixture.consumer.ID {
		t.Fatalf("abandon plan worktrees = %#v", plan.Worktrees)
	}
	status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil || status.State != "terminal" {
		t.Fatalf("ApplyAbandon = %#v, %v", status, err)
	}
	if recovery := git(t, fixture.root, "for-each-ref", "refs/bench/recovery/"); recovery == "" {
		t.Fatal("abandon preserved no recovery ref for the dirty consumer")
	}
}
