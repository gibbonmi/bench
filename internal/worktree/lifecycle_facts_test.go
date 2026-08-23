package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
)

// The tests in this file are the FA1 focused fact-adapter proofs for the
// lifecycle policy: one real-Git test per named lifecycle fact group —
// ownership, lease, eligibility, age, ignored output, preservation, and
// action — each asserting the exact typed facts the parent boundary hands
// internal/worktree/lifecyclepolicy.

// TestLifecycleOwnershipFactAdapterTranslatesRealBundle is the ownership fact
// group. A request-created bundle must translate into an owner-marker join
// against the real registration: the matched assignment, its branch, and the
// prescribed Bench lock reason the policy compares against the recorded one.
func TestLifecycleOwnershipFactAdapterTranslatesRealBundle(t *testing.T) {
	root, creation := newOwnedAssignment(t, "fa-ownership")
	facts := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, facts.MarkerPresent && facts.MarkerErr == nil, "ownership facts = %+v, want a validated owner marker", facts)
	requireTest(t, facts.MatchedAssignment != nil && facts.MatchedAssignment.ID == creation.Assignment.ID,
		"matched assignment = %+v, want the created assignment %s", facts.MatchedAssignment, creation.Assignment.ID)
	requireTest(t, facts.RegistrationBranchRef == creation.Assignment.Branch,
		"registration branch fact %q, want the assignment branch %q", facts.RegistrationBranchRef, creation.Assignment.Branch)
	requireTest(t, facts.AssignmentLockReason == lockReason(creation.Assignment) && facts.RegistrationLockReason == facts.AssignmentLockReason,
		"lock-reason facts = %q vs %q, want both to carry the assignment's prescribed reason", facts.RegistrationLockReason, facts.AssignmentLockReason)
	requireTest(t, facts.RegistrationLocked && !facts.RegistrationDetached, "registration facts = %+v, want locked and attached", facts)

	plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
	mustNoError(t, err)
	requireTest(t, plan.owned && plan.Action == ActionRemove,
		"owned clean plan = %#v, want the production adapter to reach the same owned removal over these facts", plan)
}

// TestLifecycleLeaseFactAdapterTranslatesRealLeases is the lease fact group.
// Real lease files — live owner, dead owner, malformed content, absent — must
// translate into the exact typed lease facts the policy consumes.
func TestLifecycleLeaseFactAdapterTranslatesRealLeases(t *testing.T) {
	root, creation := newOwnedAssignment(t, "fa-lease")
	lease, err := LeaseFile(creation.Path)
	mustNoError(t, err)

	absent := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, !absent.LeasePresent && absent.LeaseStatErr == nil, "absent-lease facts = %+v, want no lease evidence", absent)

	mustWrite(t, lease, []byte(fmt.Sprintf("%d 2026-07-15T00:00:00Z\n", os.Getpid())), 0o600)
	live := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, live.LeasePresent && live.LeaseState == LeaseLive, "live-lease facts = %+v, want a present live lease", live)

	mustWrite(t, lease, []byte(deadPidLine(t)), 0o600)
	dead := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, dead.LeasePresent && dead.LeaseState == LeaseDead, "dead-lease facts = %+v, want a present dead lease", dead)

	mustWrite(t, lease, []byte("junk\n"), 0o600)
	malformed := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, malformed.LeasePresent && malformed.LeaseState == LeaseUnknown, "malformed-lease facts = %+v, want an unknown lease", malformed)
}

// TestLifecycleEligibilityFactAdapterTranslatesTrackedState is the eligibility
// fact group. A real dirty checkout must translate into the typed tracked and
// nested facts, and the production planner must reach the policy's verdict
// over exactly those facts.
func TestLifecycleEligibilityFactAdapterTranslatesTrackedState(t *testing.T) {
	root, creation := newOwnedAssignment(t, "fa-eligibility")
	clean := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, clean.InitialTracked == "clean" && clean.NestedState == nestedClean && clean.NestedErr == nil,
		"clean facts = %+v, want clean tracked and nested state", clean)

	mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o644)
	dirty := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, dirty.InitialTracked == "dirty", "dirty facts = %+v, want tracked=dirty", dirty)

	plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
	mustNoError(t, err)
	requireTest(t, plan.Action == ActionRecoverRemove && plan.Tracked == "dirty",
		"dirty plan = %#v, want the production adapter to promote these facts to recover-remove", plan)
}

// TestLifecycleAgeFactAdapterTranslatesLedgerStamp is the age fact group. The
// real ledger's creation stamp must feed the policy's age decision through the
// boundary, which supplies the bounds.AssignmentStale window: a fresh record
// is not orphaned, and the same record backdated past the window is.
func TestLifecycleAgeFactAdapterTranslatesLedgerStamp(t *testing.T) {
	root, creation := newOwnedAssignment(t, "fa-age")
	young, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, young.CreatedAt != nil, "created assignment carries no creation stamp")
	requireTest(t, !orphaned(young, currentTime()), "fresh assignment reads orphaned")

	backdate(t, root, creation.Assignment, 8*24*time.Hour)
	aged, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, orphaned(aged, currentTime()), "backdated assignment does not read orphaned")
}

// TestLifecycleIgnoredFactAdapterTranslatesDeclaration is the ignored-output
// fact group. A real ignored inventory must translate into the typed count and
// declared-allowance facts: undeclared residue reads undeclared, and the same
// residue under a build-output declaration reads declared.
func TestLifecycleIgnoredFactAdapterTranslatesDeclaration(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore build output")
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, "fa-ignored", "ignored declaration")
	mustMkdirAll(t, filepath.Join(creation.Path, "dist"), 0o755)
	mustWrite(t, filepath.Join(creation.Path, "dist", "bench"), []byte("binary\n"), 0o755)

	undeclared := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, undeclared.IgnoredCount == 1 && !undeclared.DeclaredIgnored && undeclared.IgnoredErr == nil && !undeclared.IgnoredOverLimit,
		"undeclared facts = %+v, want one undeclared ignored path", undeclared)

	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte(`{"schema":1,"paths":["dist/"]}`+"\n"), 0o644)
	declared := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, declared.IgnoredCount == 1 && declared.DeclaredIgnored && declared.BuildOutputErr == nil,
		"declared facts = %+v, want the same path inside the declared allowance", declared)
}

// TestLifecyclePreservationFactAdapterTranslatesPlanShape is the preservation
// fact group. The production plan's preserves reading must be the policy
// Preserves predicate over the plan's own translated action, tracked state,
// and registration shape.
func TestLifecyclePreservationFactAdapterTranslatesPlanShape(t *testing.T) {
	root, creation := newOwnedAssignment(t, "fa-preservation")
	clean, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
	mustNoError(t, err)
	requireTest(t, !clean.preserves() && !lifecyclepolicy.Preserves(clean.Action, clean.Tracked, clean.registration.Detached),
		"clean plan = %#v, want no preservation ahead", clean)

	mustWrite(t, filepath.Join(creation.Path, "dirty.txt"), []byte("uncommitted\n"), 0o644)
	dirty, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
	mustNoError(t, err)
	requireTest(t, dirty.preserves() && lifecyclepolicy.Preserves(dirty.Action, dirty.Tracked, dirty.registration.Detached),
		"dirty plan = %#v, want a preserving removal ahead", dirty)
}

// TestLifecycleActionFactAdapterTranslatesLandedness is the action fact group.
// Real branch ancestry must translate into the typed landedness facts that
// authorize branch deletion: a branch at the default tip is proven landed with
// its exact ref and OID, and a diverged branch is proven unmerged.
func TestLifecycleActionFactAdapterTranslatesLandedness(t *testing.T) {
	root, creation := newOwnedAssignment(t, "fa-action")
	landed := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	head := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	requireTest(t, !landed.HeadDetached && landed.DefaultKnown && landed.LandedOK && landed.LandedErr == nil,
		"landed facts = %+v, want a proven landed branch", landed)
	requireTest(t, landed.HeadRef == creation.Assignment.Branch && landed.Head == head,
		"landed identity facts = %q/%q, want %q/%q", landed.HeadRef, landed.Head, creation.Assignment.Branch, head)

	commitInWorktree(t, creation.Path, "fa-action.txt", "diverge\n", "fa-action diverges")
	diverged := gatherExplicitFactsForTest(t, root, creation.Path, CleanupOptions{})
	requireTest(t, diverged.DefaultKnown && !diverged.LandedOK && diverged.LandedErr == nil,
		"diverged facts = %+v, want proven unmerged", diverged)
}
