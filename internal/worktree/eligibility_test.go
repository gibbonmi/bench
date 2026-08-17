package worktree

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
