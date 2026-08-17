package worktree

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

// TestExplicitEligibilityOutcomeMatrix is a characterization oracle for
// PlanExplicitWithOptions: it pins the nine reachable explicit `(Action,
// ReasonCode)` tuples against the pre-refactor tree. Each expected tuple below
// is hand-authored from the block order documented in subshell.go, never
// derived by calling production decision code a second time. Where a rule's
// current position only matters because a later block can overwrite an
// earlier one, the fixture carries both pieces of evidence so the assertion
// pins the actual winner rather than an idealized precedence.
func TestExplicitEligibilityOutcomeMatrix(t *testing.T) {
	t.Run("retain-uncertain", func(t *testing.T) {
		// EX1a: the primary checkout is refused before any other evidence is
		// gathered (subshell.go L93-97), an early return rather than a
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
		// earlier block, but unsafe control bytes in the target path win the
		// final override at L290-293, beating even a clean removal.
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
	// path convention alone proves nothing (L98-102).
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

	// EX3: two different malformed-evidence sources compete — a malformed
	// owner marker (L120, decided in the early marker block) and a malformed
	// build-output declaration (L232, decided in the later ignored-residue
	// block). Both land the same ReasonCode, but only the later block's
	// Reason text survives, pinning the current later-rule winner rather than
	// the marker block's own detail.
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

	// EX4: a foreign lock on an unmarked registration is isolated evidence —
	// no lease, ignored residue, or nested-state condition present — so it
	// survives to the end exactly as the marker block (L150-151) decided it.
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

	// EX5: the marker block first refuses on a lock-reason mismatch
	// (L141-144, an owned registration whose Bench lock no longer matches its
	// assignment) — an earlier applicable refusal. A live lease on the same
	// fixture is decided afterward at L196-216 and overwrites it, pinning the
	// lease block's effective position after the marker block rather than an
	// idealized order.
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

	// EX6: a live lease is decided at L196-216, an earlier applicable
	// refusal. Undeclared ignored residue without --discard-ignored is
	// decided afterward at L228-239 and overwrites it, pinning the
	// ignored-residue block's position after lease/nested evaluation.
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
	// dirty rather than clean, so the recovery promotion at L270-286 chooses
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
	// .bench/build-outputs.json is authorized rather than undeclared — the
	// opposite authorization path from EX6's --discard-ignored flag — so the
	// discard-remove promotion at L287-289 applies instead of the ignored
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
// encoder refuses. It mirrors the exact bundle Create writes — matching
// marker, matching assignment join, matching branch, matching lock reason —
// so that every block before the final unsafe-target override would resolve
// this fixture to `remove` with an empty reason code, and only the override
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

// TestAutomaticEligibilityOutcomeMatrix is a characterization oracle for PlanAutomatic
// (classifier.go:284-347): it pins the thirteen reachable automatic `(Action,
// ReasonCode)` tuples against the pre-refactor tree. PlanAutomatic always calls
// PlanExplicit with an empty CleanupOptions{} first and then applies its own ordered
// overrides on top — this structurally means DiscardBranch, an explicit-only operator
// assertion, can never reach automatic evaluation. Each expected tuple below is
// hand-authored from that ordering, never derived by calling production decision code a
// second time. Where a case exists only because automatic's own check overrides — or
// fails to override — what explicit alone already decided, the fixture carries both
// pieces of evidence so the assertion pins the actual winner.
func TestAutomaticEligibilityOutcomeMatrix(t *testing.T) {
	// AU1: an active assignment registered under its declared branch, whose branch
	// object cannot actually be resolved (git object replacement corrupts the
	// landedness query the same way TestPlanAutomaticUsesLandedInDefaultMatrix's
	// "landedness query failure retained" case does). The corruption reaches deep
	// enough that PlanExplicit itself returns a genuine error (a raw git failure,
	// not a formatted "unknown:" landed string) rather than succeeding with an
	// uncertain plan. PlanAutomatic's top-level fallback (classifier.go L296-300)
	// catches exactly this: an explicit-planning error whose text names neither
	// "assignment" nor "intent ledger" retains uncertain rather than guessing
	// unmerged or removable from a fact it never obtained. The evidence pinned
	// here is deliberately the tuple only — the exact git error string is not a
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
	// clean and unlocked. Explicit's own registration-only reading (subshell.go
	// L150-154, the same foreign-clean shape EX7 pins) would allow `remove` with
	// no refusal — but automatic never reaches that: since explicit's Action here
	// is not Retain, automatic's stricter ownership join (classifier.go L312-314)
	// fires and refuses foreign regardless of what explicit alone permitted.
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
	// checks recovery agreement, so its own action here is `recover-remove` (the
	// declared Recovery entry promotes it, subshell.go L270-286) — but automatic's
	// recovery conjunct (classifier.go L330-332, ADR 0005) refuses malformed
	// instead of trusting a Recovery entry nothing on disk backs.
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

	// AU4: a foreign lock on an unmarked, unregistered-by-marker registration —
	// EX4's exact fixture — is isolated evidence with nothing later to override it
	// in explicit planning, so explicit's own decision is already
	// retain/unexpected-lock. Automatic's first check (classifier.go L306-311)
	// passes any already-Retain explicit result straight through, so the reason
	// survives rather than being replaced by the generic ownership refusal at
	// L312-314.
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
	// explicit planning the ignored-residue block (subshell.go L228-239) runs
	// after the lease block (L196-216) and overwrites it, so explicit's own final
	// answer here is retain/ignored, not retain/live-lease. Automatic's dedicated
	// lease recheck (classifier.go L302-305) runs before any retain-passthrough
	// or ownership check and reprobes the same lease file independently, so it
	// overrides explicit's ignored answer with live-lease regardless.
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

	// AU6: TestPlanAutomaticKeepsEarlierRetainReason's exact fixture — an aged,
	// still-Active, unlanded assignment carrying undeclared ignored residue.
	// Explicit's own ignored-residue refusal (subshell.go L228-239) already makes
	// the plan Retain, so automatic's first check (classifier.go L306-311) passes
	// it straight through: since the branch has not landed, assignmentLanded is
	// false and no override fires, and the block below that would relabel an aged
	// active assignment orphaned (L316-328) is never reached at all because the
	// function already returned.
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
	// cleanup-pending), so automatic reaches the not-cleanup-pending block
	// (classifier.go L316-328); the branch is unlanded so reason starts Active,
	// and the assignment is young so the orphaned age override never fires.
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

	// AU8: mirrors TestPlanAutomaticKeepsEarlierRetainReason's fixture shape
	// (owned, Active, explicit already refusing on undeclared ignored residue),
	// but with no divergent commit so the branch is trivially landed by ancestry.
	// Automatic's first check (classifier.go L306-311) still passes the plan
	// through as Retain, but this time assignmentLanded is true, so the reason is
	// overwritten to landed — the opposite of AU6, proving the override applies
	// only when landed and discards whatever explicit refusal preceded it.
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
	// refusal, so automatic reaches the not-cleanup-pending block; the branch is
	// unlanded so reason starts Active, but the assignment was backdated past
	// bounds.AssignmentStale, so the orphaned override (classifier.go L324-326)
	// relabels it before returning.
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
	// reachable from the default branch. PlanAutomatic's own signature (classifier.go
	// L284, confirmed by reading it, not guessed) takes only (root, path string) —
	// it cannot accept CleanupOptions at all, so DiscardBranch structurally cannot
	// reach it. The dedicated assertion below proves this rather than asserting it
	// by inspection alone: explicit planning over the identical fixture WITH
	// DiscardBranch authorizes exact branch deletion (plan.deleteBranch), yet
	// PlanAutomatic over the same fixture still retains unmerged and carries no
	// branch-deletion authority.
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
	// ancestry but whose checkout carries an untracked file. Explicit's own
	// dirty-preservation promotion (subshell.go L270-286) already makes this
	// recover-remove, so plan.preserves() is true; automatic's final check
	// (classifier.go L343-345) retains rather than authoring a recovery ref
	// unattended. The before/after lifecycle snapshot (which includes
	// refs/bench/) is the proof that no recovery evidence gets authored.
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

	// AU12: a verified, cleanup-pending, trivially-landed, clean assignment —
	// the one case automatic's stricter reading still admits — falls through
	// every retain check to classifier.go L346 with an empty reason code.
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
	// (EX9's declaration shape). Explicit's own discard-remove promotion
	// (subshell.go L287-289) already applies since the residue is declared;
	// automatic's preservation check does not fire because plan.preserves()
	// requires DiscardRemove with a non-clean Tracked state, and Tracked stays
	// "clean" here, so the plan falls through unchanged.
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
