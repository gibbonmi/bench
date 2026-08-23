package worktree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestExplicitEligibilityOutcomeMatrix is a characterization oracle for
// PlanExplicitWithOptions: it pins the nine reachable explicit `(Action,
// ReasonCode)` tuples. Each expected tuple below is hand-authored from the
// block order decideExplicit documents, never derived by calling production
// decision code a second time. A rule's current position matters only when
// a later block can overwrite an earlier one. There the fixture carries
// both pieces of evidence, so the assertion pins the actual winner rather
// than an idealized precedence.
func TestExplicitEligibilityOutcomeMatrix(t *testing.T) {
	t.Run("retain-uncertain", func(t *testing.T) {
		// EX1a: PlanExplicitWithOptions refuses the primary checkout before any
		// other evidence is gathered, an early return rather than a
		// last-write-wins overwrite.
		t.Run("primary", func(t *testing.T) {
			root := newWorktreeRepo(t)
			marker, err := markerPath(root)
			mustNoError(t, err)
			before := lifecycleSnapshot(t, root, root, marker)
			plan, err := PlanExplicitWithOptions(root, root, CleanupOptions{})
			mustNoError(t, err)
			requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonUncertain && plan.Reason == "primary checkout is never removable",
				"primary checkout plan = %#v, want retain/uncertain \"primary checkout is never removable\"", plan)
			after := lifecycleSnapshot(t, root, root, marker)
			requireTest(t, before == after, "primary-checkout planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
		})
		// EX1b: an otherwise-clean, otherwise-removable, correctly owned and
		// locked registration would end `remove` with no reason code at every
		// earlier block. But unsafe control bytes in the target path win
		// decideExplicit's final override, beating even a clean removal.
		t.Run("unsafe-target-override", func(t *testing.T) {
			root, target := newUnsafeTargetOwnedWorktree(t, "ex1-unsafe-override")
			marker, err := markerPath(target)
			mustNoError(t, err)
			before := lifecycleSnapshot(t, root, target, marker)
			plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{})
			mustNoError(t, err)
			requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonUncertain && plan.Reason == "target contains unsafe control bytes" && plan.Recovery == "none",
				"unsafe-target plan = %#v, want retain/uncertain \"target contains unsafe control bytes\" with recovery reset to none", plan)
			after := lifecycleSnapshot(t, root, target, marker)
			requireTest(t, before == after, "unsafe-target planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
		})
	})

	// EX2: an unregistered target never reaches ownership evidence at all —
	// path convention alone proves nothing.
	t.Run("retain-foreign", func(t *testing.T) {
		root := newWorktreeRepo(t)
		target := filepath.Join(t.TempDir(), "unregistered")
		mustMkdirAll(t, target, 0o700)
		marker := filepath.Join(target, OwnerMarkerFile)
		before := lifecycleSnapshot(t, root, target, marker)
		plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonForeign && plan.Reason == "target is not registered",
			"unregistered target plan = %#v, want retain/foreign \"target is not registered\"", plan)
		after := lifecycleSnapshot(t, root, target, marker)
		requireTest(t, before == after, "unregistered-target planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// EX3: two different malformed-evidence sources compete: a malformed
	// owner marker, decided in decideExplicit's early marker block, and a
	// malformed build-output declaration, decided in its later ignored-residue
	// block. Both land the same ReasonCode. Only the later block's Reason
	// text survives, pinning the current later-rule winner rather than the
	// marker block's own detail.
	t.Run("retain-malformed", func(t *testing.T) {
		root := newWorktreeRepo(t)
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		creation := mustCreate(t, root, "ex3-malformed", "malformed marker and declaration")
		markerFile, err := markerPath(creation.Path)
		mustNoError(t, err)
		mustWrite(t, markerFile, []byte("{"), 0o600)
		mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
		mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{"), 0o644)
		before := lifecycleSnapshot(t, root, creation.Path, markerFile)
		plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonMalformed && plan.Reason == "build-output declaration is malformed",
			"competing-malformed plan = %#v, want the later build-output-declaration reason to win over the earlier marker reason", plan)
		after := lifecycleSnapshot(t, root, creation.Path, markerFile)
		requireTest(t, before == after, "competing-malformed planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// EX4: a foreign lock on an unmarked registration is isolated evidence.
	// No lease, ignored residue, or nested-state condition is present, so it
	// survives to the end exactly as decideExplicit's ownership block decided it.
	t.Run("retain-unexpected-lock", func(t *testing.T) {
		root := newWorktreeRepo(t)
		target := filepath.Join(t.TempDir(), "foreign locked")
		gitRun(t, root, "worktree", "add", "-q", "-b", "ex4-foreign-locked", target, "HEAD")
		gitRun(t, root, "worktree", "lock", "--reason", "foreign", target)
		marker, err := markerPath(target)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, target, marker)
		plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonUnexpectedLock && plan.Reason == "foreign or unexpected lock is retained",
			"foreign-lock plan = %#v, want retain/unexpected-lock \"foreign or unexpected lock is retained\"", plan)
		after := lifecycleSnapshot(t, root, target, marker)
		requireTest(t, before == after, "foreign-lock planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// EX5: decideExplicit's marker block first refuses on a lock-reason
	// mismatch (an owned registration whose Bench lock no longer matches its
	// assignment), an earlier applicable refusal. A live lease on the same
	// fixture is decided afterward in the lease block and overwrites it. This
	// pins the lease block's effective position after the marker block
	// rather than an idealized order.
	t.Run("retain-live-lease", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "ex5-live-lease")
		gitRun(t, root, "worktree", "unlock", creation.Path)
		gitRun(t, root, "worktree", "lock", "--reason", "manually altered lock reason", creation.Path)
		lease, err := LeaseFile(creation.Path)
		mustNoError(t, err)
		mustWrite(t, lease, []byte(fmt.Sprintf("%d 2026-07-15T00:00:00Z\n", os.Getpid())), 0o600)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonLiveLease && plan.Reason == "assignment has a live lease",
			"competing lock/lease plan = %#v, want the later live-lease reason to beat the earlier lock-mismatch refusal", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "competing lock/lease planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// EX6: decideExplicit's lease block is an earlier applicable refusal.
	// Undeclared ignored residue without --discard-ignored is decided
	// afterward in the ignored-residue block and overwrites it, pinning that
	// block's position after lease/nested evaluation.
	t.Run("retain-ignored", func(t *testing.T) {
		root := newWorktreeRepo(t)
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
		gitRun(t, root, "add", ".gitignore")
		gitRun(t, root, "commit", "-qm", "ignore")
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		creation := mustCreate(t, root, "ex6-ignored", "ignored residue")
		lease, err := LeaseFile(creation.Path)
		mustNoError(t, err)
		mustWrite(t, lease, []byte(fmt.Sprintf("%d 2026-07-15T00:00:00Z\n", os.Getpid())), 0o600)
		mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonIgnored && plan.Reason == "ignored residuals require --discard-ignored" && plan.Ignored.Count == 1,
			"competing lease/ignored plan = %#v, want the later ignored-residue reason to beat the earlier live-lease refusal", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "competing lease/ignored planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// EX7: a clean, owned, correctly locked, registered checkout carries no
	// refusal evidence anywhere in the chain and ends `remove` with an empty
	// reason code.
	t.Run("remove", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "ex7-remove")
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRemove && plan.ReasonCode == "",
			"clean removable plan = %#v, want remove with empty reason code", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "clean-removable planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// EX8: an otherwise-removable checkout carrying an untracked file is
	// dirty rather than clean. So decideExplicit's recovery promotion chooses
	// `recover-remove` and computes a fresh recovery ref rather than removing
	// outright.
	t.Run("recover-remove", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "ex8-recover-remove")
		mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o644)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRecoverRemove && plan.ReasonCode == "" && plan.Recovery != "" && plan.Recovery != "none",
			"dirty removable plan = %#v, want recover-remove with a computed recovery ref", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "recover-remove planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// EX9: bounded ignored residue that is declared as build output in
	// .bench/build-outputs.json is authorized rather than undeclared. This is
	// the opposite authorization path from EX6's --discard-ignored flag, so
	// decideExplicit's discard-remove promotion applies instead of the ignored
	// refusal.
	t.Run("discard-remove", func(t *testing.T) {
		root := newWorktreeRepo(t)
		gitRun(t, root, "branch", "-M", "main")
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644)
		gitRun(t, root, "add", ".gitignore")
		gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore build output")
		mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
		mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte(`{"schema":1,"paths":["dist/"]}`+"\n"), 0o644)
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		creation := mustCreate(t, root, "ex9-discard-remove", "declared build output")
		mustMkdirAll(t, filepath.Join(creation.Path, "dist"), 0o755)
		mustWrite(t, filepath.Join(creation.Path, "dist", "bench"), []byte("binary\n"), 0o755)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionDiscardRemove && plan.ReasonCode == "" && plan.Ignored.Count > 0,
			"declared-output plan = %#v, want discard-remove with an empty reason code", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "discard-remove planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})
}

// newUnsafeTargetOwnedWorktree hand-builds one owned, correctly locked, clean
// registration whose target path carries a control byte the toon output
// encoder refuses. It mirrors the exact bundle Create writes: matching
// marker, matching assignment join, matching branch, matching lock reason.
// So every block before the final unsafe-target override would resolve
// this fixture to `remove` with an empty reason code. Only the override
// itself explains the retained result.
func newUnsafeTargetOwnedWorktree(t *testing.T, request string) (root, target string) {
	t.Helper()
	root = newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	ownerID, err := randomID()
	mustNoError(t, err)
	assignmentID, err := randomID()
	mustNoError(t, err)
	rawTarget := filepath.Join(t.TempDir(), "unsafe\x07byte-"+request)
	target, err = canonicalPath(rawTarget)
	mustNoError(t, err)
	branch := intent.AssignmentBranchRef(ownerID, assignmentID)
	shortBranch := strings.TrimPrefix(branch, "refs/heads/")
	start := gitOutput(t, root, "rev-parse", "HEAD")
	assignment := intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: assignmentID, OwnerID: ownerID,
		Request: intent.RequestDigest(request), Label: "unsafe target fixture", Start: start,
		Branch: branch, Worktree: target, State: intent.StateActive, Recovery: []intent.Recovery{},
	}
	gitRun(t, root, "worktree", "add", "-q", "--lock", "--reason", lockReason(assignment), "-b", shortBranch, target, start)
	admin := gitOutput(t, target, "rev-parse", "--path-format=absolute", "--git-dir")
	marker := Marker{Schema: OwnerMarkerSchema, OwnerID: ownerID, Path: target}
	body, err := json.Marshal(marker)
	mustNoError(t, err)
	mustWrite(t, filepath.Join(admin, OwnerMarkerFile), append(body, '\n'), 0o600)
	mustNoError(t, intent.PutAssignment(root, assignment))
	return root, target
}

func TestExplicitEligibilityAllowsRuntimeIgnoredResidue(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte(".logs/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore runtime records")
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, "runtime-eligibility", "runtime eligibility")
	mustMkdirAll(t, filepath.Join(creation.Path, ".logs"), 0o755)
	mustWrite(t, filepath.Join(creation.Path, ".logs", "gate.jsonl"), []byte("record\n"), 0o644)
	plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
	if err != nil || plan.Action != ActionDiscardRemove || plan.ReasonCode != "" {
		t.Fatalf("runtime residue plan = (%#v, %v), want discard-remove without refusal", plan, err)
	}
}

// TestAutomaticEligibilityOutcomeMatrix is a characterization oracle for PlanAutomatic:
// it pins the thirteen reachable automatic `(Action, ReasonCode)` tuples. PlanAutomatic
// always calls PlanExplicit with an empty CleanupOptions{} first and then applies its own
// ordered overrides on top. This structurally means DiscardBranch, an explicit-only
// operator assertion, can never reach automatic evaluation. Each expected tuple below
// is hand-authored from that ordering, never derived by calling production decision
// code a second time.
//
// A case can exist only because automatic's own check overrides,
// or fails to override, what explicit alone already decided. There the fixture carries
// both pieces of evidence, so the assertion pins the actual winner.
func TestAutomaticEligibilityOutcomeMatrix(t *testing.T) {
	// AU1: an active assignment registered under its declared branch, whose branch
	// object cannot actually be resolved. Git object replacement corrupts the
	// landedness query the same way TestPlanAutomaticUsesLandedInDefaultMatrix's
	// "landedness query failure retained" case does. The corruption reaches deep
	// enough that PlanExplicit itself returns a genuine error: a raw git failure,
	// not a formatted "unknown:" landed string. It does not succeed with an
	// uncertain plan.
	//
	// decideAutomatic's explicit-planning-error fallback catches
	// exactly this: an explicit-planning error whose text names neither
	// "assignment" nor "intent ledger" retains uncertain, rather than guessing
	// unmerged or removable from a fact it never obtained. The evidence pinned
	// here is deliberately the tuple only. The exact git error string is not a
	// stable sentence to pin, and the acceptance row does not require one.
	t.Run("retain-uncertain", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "au1-unknown-landed")
		commitInWorktree(t, creation.Path, "au1.txt", "au1\n", "au1 unmerged commit")
		branchOID := gitOutput(t, creation.Path, "rev-parse", "HEAD")
		replace := filepath.Join(root, ".git", "refs", "replace", branchOID)
		mustMkdirAll(t, filepath.Dir(replace), 0o700)
		mustWrite(t, replace, []byte(strings.Repeat("f", 40)+"\n"), 0o600)
		markPending(t, root, creation.Assignment)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonUncertain && plan.Reason != "",
			"unresolvable-branch-object plan = %#v, want retain/uncertain with a non-empty detail", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "unresolvable-branch-object planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU2: a worktree git itself registers, carrying no owner marker at all, is
	// clean and unlocked. Explicit's own registration-only reading (the same
	// foreign-clean shape EX7 pins) would allow `remove` with no refusal. But
	// automatic never reaches that: since explicit's Action here is not Retain,
	// decideAutomatic's stricter ownership join fires and refuses foreign
	// regardless of what explicit alone permitted.
	t.Run("retain-foreign", func(t *testing.T) {
		root := newWorktreeRepo(t)
		target := filepath.Join(t.TempDir(), "au2-registered-unowned")
		gitRun(t, root, "worktree", "add", "-q", "-b", "au2-registered-unowned", target, "HEAD")
		marker, err := markerPath(target)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, target, marker)
		plan, err := PlanAutomatic(root, target)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonForeign && plan.Reason == "registration is not a verified owned assignment",
			"unowned-registered plan = %#v, want retain/foreign \"registration is not a verified owned assignment\" even though explicit alone would remove it", plan)
		after := lifecycleSnapshot(t, root, target, marker)
		requireTest(t, before == after, "unowned-registered planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU3: a fully owned, cleanup-pending assignment whose ledger record carries a
	// recovery entry that names no actual recovery manifest. Explicit itself never
	// checks recovery agreement, so its own action here is `recover-remove`: the
	// declared Recovery entry promotes it through decideExplicit's recovery
	// promotion. But decideAutomatic's recovery conjunct (ADR 0005) refuses
	// malformed instead of trusting a Recovery entry nothing on disk backs.
	t.Run("retain-malformed", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "au3-recovery-mismatch")
		assignment := creation.Assignment
		assignment.State = intent.StateCleanupPending
		assignment.Recovery = []intent.Recovery{{
			Ref:      intent.RecoveryRefPrefix(assignment.OwnerID, assignment.ID) + "1",
			Root:     assignment.Start,
			Payloads: []string{assignment.Start},
		}}
		mustNoError(t, intent.PutAssignment(root, assignment))
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonMalformed && plan.Reason == "assignment recovery metadata does not match refs",
			"recovery-mismatch plan = %#v, want retain/malformed \"assignment recovery metadata does not match refs\"", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "recovery-mismatch planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU4: a foreign lock on an unmarked, unregistered-by-marker registration,
	// EX4's exact fixture, is isolated evidence with nothing later to override it
	// in explicit planning. So explicit's own decision is already
	// retain/unexpected-lock. decideAutomatic's retain-passthrough passes any
	// already-Retain explicit result straight through, so the reason survives
	// rather than being replaced by the generic ownership refusal below it.
	t.Run("retain-unexpected-lock", func(t *testing.T) {
		root := newWorktreeRepo(t)
		target := filepath.Join(t.TempDir(), "au4 foreign locked")
		gitRun(t, root, "worktree", "add", "-q", "-b", "au4-foreign-locked", target, "HEAD")
		gitRun(t, root, "worktree", "lock", "--reason", "foreign", target)
		marker, err := markerPath(target)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, target, marker)
		plan, err := PlanAutomatic(root, target)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonUnexpectedLock && plan.Reason == "foreign or unexpected lock is retained",
			"foreign-lock plan = %#v, want retain/unexpected-lock \"foreign or unexpected lock is retained\"", plan)
		after := lifecycleSnapshot(t, root, target, marker)
		requireTest(t, before == after, "foreign-lock planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU5: EX6's exact live-lease-plus-undeclared-ignored-residue fixture. Inside
	// decideExplicit the ignored-residue block runs after the lease block and
	// overwrites it, so explicit's own final answer here is retain/ignored, not
	// retain/live-lease. PlanAutomatic reprobes the same lease file independently.
	// decideAutomatic's live-lease override reads that reprobe before any
	// retain-passthrough or ownership check, so it overrides explicit's ignored
	// answer with live-lease regardless.
	t.Run("retain-live-lease", func(t *testing.T) {
		root := newWorktreeRepo(t)
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
		gitRun(t, root, "add", ".gitignore")
		gitRun(t, root, "commit", "-qm", "ignore")
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		creation := mustCreate(t, root, "au5-live-lease", "competing ignored/live-lease")
		lease, err := LeaseFile(creation.Path)
		mustNoError(t, err)
		mustWrite(t, lease, []byte(fmt.Sprintf("%d 2026-07-15T00:00:00Z\n", os.Getpid())), 0o600)
		mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonLiveLease && plan.Reason == "assignment has a live lease" && plan.Ignored.Count == 1,
			"competing ignored/live-lease plan = %#v, want automatic's own lease recheck to beat explicit's ignored answer", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "competing ignored/live-lease planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU6: TestPlanAutomaticKeepsEarlierRetainReason's exact fixture: an aged,
	// still-Active, unlanded assignment carrying undeclared ignored residue.
	// decideExplicit's ignored-residue refusal already makes the plan Retain, so
	// decideAutomatic's retain-passthrough passes it straight through. The
	// branch has not landed, so assignmentLanded is false and no override fires.
	// The not-cleanup-pending block below, which would relabel an aged active
	// assignment orphaned, is never reached because the function already
	// returned.
	t.Run("retain-ignored", func(t *testing.T) {
		root := newWorktreeRepo(t)
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
		gitRun(t, root, "add", ".gitignore")
		gitRun(t, root, "commit", "-qm", "ignore")
		creation := mustCreate(t, root, "au6-aged-ignored", "aged ignored residue")
		commitInWorktree(t, creation.Path, "au6.txt", "au6\n", "au6 unmerged commit")
		mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
		backdate(t, root, creation.Assignment, 8*24*time.Hour)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonIgnored && plan.Reason == "ignored residuals require --discard-ignored",
			"aged ignored plan = %#v, want the earlier explicit ignored refusal kept rather than relabeled orphaned", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "aged ignored planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU7: TestPlanAutomaticLabelsOrphaned's shape without the backdate. Explicit
	// has no refusal here (a clean owned checkout, still Active rather than
	// cleanup-pending), so decideAutomatic reaches its not-cleanup-pending block.
	// The branch is unlanded, so reason starts Active. The assignment is young,
	// so the orphaned age override never fires.
	t.Run("retain-active", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "au7-young-active")
		commitInWorktree(t, creation.Path, "au7.txt", "au7\n", "au7 unmerged commit")
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonActive && plan.Reason == "assignment is not cleanup-pending",
			"young active plan = %#v, want retain/active \"assignment is not cleanup-pending\"", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "young active planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU8 mirrors TestPlanAutomaticKeepsEarlierRetainReason's fixture shape:
	// owned, Active, explicit already refusing on undeclared ignored residue.
	// It carries no divergent commit, so the branch is trivially landed by ancestry.
	// decideAutomatic's retain-passthrough still passes the plan through as
	// Retain, but this time assignmentLanded is true, so the reason is overwritten
	// to landed. This is the opposite of AU6: the override applies only when
	// landed, and it discards whatever explicit refusal preceded it.
	t.Run("retain-landed", func(t *testing.T) {
		root := newWorktreeRepo(t)
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
		gitRun(t, root, "add", ".gitignore")
		gitRun(t, root, "commit", "-qm", "ignore")
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		creation := mustCreate(t, root, "au8-landed-override", "landed active with competing ignored refusal")
		mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonLanded && plan.Reason == "assignment branch has landed",
			"landed-active plan = %#v, want the landed override to replace the earlier ignored refusal with retain/landed", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "landed-active planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU9: TestPlanAutomaticLabelsOrphaned's exact fixture. Explicit again has no
	// refusal, so automatic reaches the not-cleanup-pending block. The branch is
	// unlanded, so reason starts Active. But the assignment was backdated past
	// bounds.AssignmentStale, so decideAutomatic's orphaned-age override relabels
	// it before returning.
	t.Run("retain-orphaned", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "au9-aged-active")
		commitInWorktree(t, creation.Path, "au9.txt", "au9\n", "au9 unmerged commit")
		backdate(t, root, creation.Assignment, 8*24*time.Hour)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonOrphaned && plan.Reason == "assignment is not cleanup-pending",
			"aged active plan = %#v, want retain/orphaned \"assignment is not cleanup-pending\"", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "aged active planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU10: a cleanup-pending, owned assignment whose branch carries a commit not
	// reachable from the default branch. PlanAutomatic's signature takes only
	// (root, path string). It cannot accept CleanupOptions at all, so
	// DiscardBranch structurally cannot reach it. The assertion below pins that
	// as behavior rather than as a signature reading.
	//
	// Explicit planning over the
	// identical fixture WITH DiscardBranch authorizes exact branch deletion
	// (plan.deleteBranch). PlanAutomatic over the same fixture still retains
	// unmerged and carries no branch-deletion authority.
	t.Run("retain-unmerged", func(t *testing.T) {
		root, creation := newPendingAssignment(t, "au10-unmerged-branch")
		commitInWorktree(t, creation.Path, "au10.txt", "au10\n", "au10 unmerged commit")
		discardBranchPlan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{DiscardBranch: true})
		mustNoError(t, err)
		requireTest(t, discardBranchPlan.deleteBranch,
			"explicit-with-DiscardBranch plan = %#v, want branch-deletion authority so the automatic contrast is meaningful", discardBranchPlan)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonUnmerged && plan.Reason == "assignment branch has not landed" && !plan.deleteBranch,
			"unmerged-branch plan = %#v, want retain/unmerged with no branch-deletion authority regardless of the explicit DiscardBranch assertion", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "unmerged-branch planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU11: a cleanup-pending assignment whose branch is trivially landed by
	// ancestry but whose checkout carries an untracked file. decideExplicit's
	// recovery promotion already makes this recover-remove, so plan.preserves()
	// is true. decideAutomatic's final preservation check retains rather than
	// authoring a recovery ref unattended. The before/after lifecycle snapshot
	// (which includes refs/bench/) proves that no recovery evidence gets
	// authored.
	t.Run("retain-dirty", func(t *testing.T) {
		root, creation := newPendingAssignment(t, "au11-dirty-landed")
		mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o644)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRetain && plan.ReasonCode == ReasonDirty && plan.Reason == "automatic cleanup does not preserve uncommitted work",
			"dirty landed plan = %#v, want retain/dirty \"automatic cleanup does not preserve uncommitted work\"", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "dirty landed planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU12: a verified, cleanup-pending, trivially-landed, clean assignment
	// is the one case automatic's stricter reading still admits. It falls
	// through every retain check in decideAutomatic to its final passthrough,
	// with an empty reason code.
	t.Run("remove", func(t *testing.T) {
		root, creation := newPendingAssignment(t, "au12-remove")
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionRemove && plan.ReasonCode == "" && plan.Assignment == creation.Assignment.ID,
			"clean pending plan = %#v, want remove with empty reason code", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "clean pending planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})

	// AU13: AU12's exact fixture with declared bounded ignored build output added
	// (EX9's declaration shape). decideExplicit's discard-remove promotion already
	// applies since the residue is declared. Automatic's preservation check does
	// not fire, because plan.preserves() requires DiscardRemove with a non-clean
	// Tracked state, and Tracked stays "clean" here. So the plan falls through
	// unchanged.
	t.Run("discard-remove", func(t *testing.T) {
		root := newWorktreeRepo(t)
		gitRun(t, root, "branch", "-M", "main")
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644)
		gitRun(t, root, "add", ".gitignore")
		gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore build output")
		mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
		mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte(`{"schema":1,"paths":["dist/"]}`+"\n"), 0o644)
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		creation := mustCreate(t, root, "au13-discard-remove", "declared build output pending")
		mustMkdirAll(t, filepath.Join(creation.Path, "dist"), 0o755)
		mustWrite(t, filepath.Join(creation.Path, "dist", "bench"), []byte("binary\n"), 0o755)
		markPending(t, root, creation.Assignment)
		marker, err := markerPath(creation.Path)
		mustNoError(t, err)
		before := lifecycleSnapshot(t, root, creation.Path, marker)
		plan, err := PlanAutomatic(root, creation.Path)
		mustNoError(t, err)
		requireTest(t, plan.Action == ActionDiscardRemove && plan.ReasonCode == "" && plan.Ignored.Count > 0,
			"declared-output pending plan = %#v, want discard-remove with an empty reason code", plan)
		after := lifecycleSnapshot(t, root, creation.Path, marker)
		requireTest(t, before == after, "declared-output pending planning mutated durable state\nbefore:\n%s\nafter:\n%s", before, after)
	})
}

// TestEligibilityVerdictProjectsWithoutSecondDecision proves decideExplicit is the whole
// decision. It independently gathers the same typed facts PlanExplicitWithOptions
// gathers, through the same package-private evidence functions: validateOwnerMarker,
// ProbeLease, classifyNestedState, inventoryIgnored, git.LandedInDefault. It never
// copies production's own answer. Calling decideExplicit directly reproduces the exact
// action, reason, typed landedness, and recovery/branch-deletion authority
// PlanExplicitWithOptions itself returns. A second decision made in subshell.go, or a
// plan mutated after decideExplicit returned, would make this independent comparison
// diverge from PlanExplicitWithOptions's own projection.
func TestEligibilityVerdictProjectsWithoutSecondDecision(t *testing.T) {
	t.Run("clean-remove", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "ev1-clean-remove")
		assertVerdictMatchesPlan(t, root, creation.Path, CleanupOptions{})
	})
	t.Run("dirty-recover-remove", func(t *testing.T) {
		root, creation := newOwnedAssignment(t, "ev1-dirty-recover-remove")
		mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o644)
		assertVerdictMatchesPlan(t, root, creation.Path, CleanupOptions{})
	})
}

// assertVerdictMatchesPlan gathers explicitFacts for path independently of
// PlanExplicitWithOptions and decides a verdict from them directly. It asserts every
// field the projection is responsible for carrying onto CleanupPlan: action, reason,
// typed landedness, and branch/recovery authority. Every field must agree with what
// PlanExplicitWithOptions actually returns for the identical fixture.
func assertVerdictMatchesPlan(t *testing.T, root, path string, options CleanupOptions) {
	t.Helper()
	target, err := canonicalPath(path)
	mustNoError(t, err)
	facts := gatherExplicitFactsForTest(t, root, target, options)
	verdict := decideExplicit(facts)

	plan, err := PlanExplicitWithOptions(root, path, options)
	mustNoError(t, err)

	requireTest(t, verdict.Action == plan.Action && verdict.ReasonCode == plan.ReasonCode && verdict.Reason == plan.Reason,
		"independently decided verdict = %#v, want it to match PlanExplicitWithOptions's projection %#v", verdict, plan)
	requireTest(t, verdict.Landed.String() == plan.landed,
		"independently decided typed landedness %q, want it to match the projected landed string %q", verdict.Landed.String(), plan.landed)
	requireTest(t, verdict.DeleteBranch == plan.deleteBranch && verdict.BranchRef == plan.branchRef && verdict.BranchOID == plan.branchOID,
		"independently decided branch-deletion authority = %v/%s/%s, want it to match the projected plan %v/%s/%s",
		verdict.DeleteBranch, verdict.BranchRef, verdict.BranchOID, plan.deleteBranch, plan.branchRef, plan.branchOID)

	wantRecovery := verdict.Recovery
	switch verdict.RecoveryLookup {
	case recoveryLookupOwned:
		wantRecovery, err = nextRecoveryRef(root, *verdict.Assignment)
		mustNoError(t, err)
	case recoveryLookupForeign:
		admin, adminErr := git.Output("-C", target, "rev-parse", "--path-format=absolute", "--git-dir")
		mustNoError(t, adminErr)
		wantRecovery, err = predictedForeignRef(root, target, admin)
		mustNoError(t, err)
	}
	if !cleanupOutputSafe(target) {
		wantRecovery = "none"
	}
	requireTest(t, wantRecovery == plan.Recovery,
		"independently resolved recovery ref %q, want it to match the projected plan recovery %q", wantRecovery, plan.Recovery)
}

// gatherExplicitFactsForTest independently gathers the same explicitFacts
// PlanExplicitWithOptions gathers for target, in the same order. This lets
// TestEligibilityVerdictProjectsWithoutSecondDecision call decideExplicit on them
// without going through subshell.go's own projection at all.
func gatherExplicitFactsForTest(t *testing.T, root, target string, options CleanupOptions) explicitFacts {
	t.Helper()
	root = canonicalRoot(root)
	worktrees, err := git.Worktrees(root)
	mustNoError(t, err)
	var registration *git.Worktree
	for i := range worktrees {
		candidate, pathErr := canonicalPath(worktrees[i].Path)
		if pathErr == nil && candidate == target {
			registration = &worktrees[i]
			break
		}
	}
	requireTest(t, registration != nil, "fixture target %s is not a registered worktree", target)
	admin, err := git.Output("-C", target, "rev-parse", "--path-format=absolute", "--git-dir")
	mustNoError(t, err)

	facts := explicitFacts{
		registrationBranchRef:  registration.BranchRef,
		registrationLockReason: registration.LockReason,
		registrationLocked:     registration.Locked,
		registrationDetached:   registration.Detached,
		discardIgnored:         options.DiscardIgnored,
	}
	_, markerStatErr := os.Lstat(filepath.Join(admin, OwnerMarkerFile))
	if markerStatErr == nil {
		facts.markerPresent = true
		evidence, markerErr := validateOwnerMarker(root, target)
		if markerErr != nil {
			facts.markerErr = markerErr
		} else {
			assignments, assignmentErr := intent.Assignments(root)
			if assignmentErr != nil {
				facts.assignmentLedgerErr = assignmentErr
			} else {
				var matched *intent.Assignment
				for i := range assignments {
					if assignments[i].Worktree == target && assignments[i].OwnerID == evidence.marker.OwnerID {
						if matched != nil {
							facts.assignmentAmbiguous = true
							break
						}
						candidate := assignments[i]
						matched = &candidate
					}
				}
				facts.matchedAssignment = matched
			}
		}
	} else if !errors.Is(markerStatErr, os.ErrNotExist) {
		t.Fatalf("stat owner marker: %v", markerStatErr)
	} else if !registration.Locked {
		facts.foreignAssignment = foreignRecoveryAssignment(root, target)
	}

	head, err := git.Output("-C", target, "rev-parse", "HEAD")
	mustNoError(t, err)
	headRef, _ := git.Output("-C", target, "symbolic-ref", "--quiet", "HEAD")
	if headRef == "" {
		headRef = "detached"
	}
	status, err := git.Raw("--no-optional-locks", "-C", target, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules=none")
	mustNoError(t, err)
	tracked := "clean"
	if len(status) > 0 {
		tracked = "dirty"
		for record := range bytes.SplitSeq(status, []byte{0}) {
			if len(record) >= 2 && (bytes.Contains(record[:2], []byte("U")) || bytes.Equal(record[:2], []byte("AA")) || bytes.Equal(record[:2], []byte("DD"))) {
				tracked = "conflicted"
			}
		}
	}
	facts.initialTracked = tracked

	leasePath, err := LeaseFile(target)
	mustNoError(t, err)
	if _, statErr := os.Lstat(leasePath); statErr == nil {
		facts.leasePresent = true
		facts.leaseState = ProbeLease(leasePath)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		facts.leaseStatErr = statErr
	}

	nested, nestedErr := classifyNestedState(target)
	facts.nestedState, facts.nestedErr = nested, nestedErr

	buildOutputs, _, buildOutputErr := loadBuildOutputs(root)
	ignored, _, ignoredErr := inventoryIgnored(target, options.Full)
	facts.buildOutputErr = buildOutputErr
	facts.ignoredErr = ignoredErr
	facts.ignoredOverLimit = ignored.OverLimit
	facts.ignoredCount = ignored.Count
	facts.declaredIgnored = buildOutputErr == nil && ignoredWithinLandingAllowance(ignored, buildOutputs)

	defaultRef, defaultOID := "none", "none"
	if def, ok := git.ResolvedDefault(root); ok {
		defaultRef = def
		if oid, oidErr := git.Output("-C", root, "rev-parse", "--verify", def+"^{commit}"); oidErr == nil {
			defaultOID = oid
		}
	}
	facts.headDetached = headRef == "detached"
	facts.defaultKnown = defaultOID != "none"
	facts.headRef, facts.head = headRef, head
	if !facts.headDetached && facts.defaultKnown {
		facts.landedOK, facts.landedByContent, facts.landedErr = git.LandedInDefault(root, strings.TrimPrefix(headRef, "refs/heads/"), defaultRef)
	}
	facts.unsafeTarget = !cleanupOutputSafe(target)
	return facts
}
