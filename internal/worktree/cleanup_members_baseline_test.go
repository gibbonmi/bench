// The frozen pre-change baseline of every cleanup apply mode. Each case of the manifest below
// builds its repository through the shared cleanup fixtures, drives the current mode adapter,
// and renders one record: the fingerprint inputs, the commit relations, the refs before the
// apply, the response rows and error, and the durable state after it. The records under
// testdata/cleanup-member-traversal came from one run of the original adapters, and
// TestCleanupTraversalBaseline compares every later run against them.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// cleanupBaselineCaptureEnv is the explicit opt-in that rewrites the frozen records from the
// current adapters. No ordinary test run sets it.
const cleanupBaselineCaptureEnv = "BENCH_CAPTURE_CLEANUP_BASELINE"

// cleanupBaselineDir holds one frozen record per manifest case.
const cleanupBaselineDir = "testdata/cleanup-member-traversal"

// baselineCase is one manifest entry: the mode adapter it drives and the case it builds.
type baselineCase struct {
	mode, name string
	build      func(t *testing.T) baselineRun
}

// baselineRun is one built case: its repository, the inputs its plan fingerprint binds, the
// planned selection order, the digests whose role the case knows, and the apply it runs.
type baselineRun struct {
	repoHome
	// byRole names each member by its assignment label instead of its identity rank. A case
	// whose members have distinct roles sets it, because their random identities fix the rank.
	byRole bool
	inputs []string
	// planned is the target of each planned row, in the order the plan reports them.
	planned []string
	// planDigest is the set fingerprint, and targetDigests maps each per-target digest to the
	// assignment it binds.
	planDigest    string
	targetDigests map[string]string
	apply         func() ([]CleanupPlan, error)
}

// cleanupBaselineManifest is the complete case list. Each mode carries only the cases that
// apply to its adapter, so the list is enumerated rather than a product of modes and cases.
// The empty explicit selection serves the unrecorded stage too: one empty set has no recorded
// or unrecorded partition.
var cleanupBaselineManifest = []baselineCase{
	{"explicit", "empty", explicitEmptyCase},
	{"explicit", "retained-only", explicitRetainedOnlyCase},
	{"explicit", "mixed", explicitMixedCase},
	{"explicit", "all-success", explicitAllSuccessCase},
	{"explicit", "first-failure", explicitFailureCase(func(j *joins, stop error) { j.cleanupBoundary = atStep(StepApplyLocked, func() error { return stop }) })},
	{"explicit", "later-failure", explicitFailureCase(faultAtSecondLock)},
	{"explicit", "outer-preflight-drift", explicitOuterDriftCase},
	{"explicit", "locked-drift", explicitLockedDriftCase},
	{"unrecorded", "retained-only", unrecordedRetainedOnlyCase},
	{"unrecorded", "mixed", unrecordedMixedCase},
	{"unrecorded", "all-success", unrecordedAllSuccessCase},
	{"unrecorded", "first-failure", unrecordedFailureCase(2, func(stop error) Fault { return atStep(StepDiscardedRefWrite, func() error { return stop }) })},
	{"unrecorded", "later-failure", unrecordedFailureCase(3, func(stop error) Fault { return atSecondHit(StepDiscardedRefWrite, func() error { return stop }) })},
	{"unrecorded", "holder-removal", unrecordedHolderRemovalCase},
	{"unrecorded", "recovery-ref-conflict", unrecordedRecoveryConflictCase},
	{"landed", "empty", landedEmptyCase},
	{"landed", "retained-only", landedRetainedOnlyCase},
	{"landed", "mixed", landedMixedCase},
	{"landed", "all-success", landedAllSuccessCase},
	{"landed", "first-failure", landedFailureCase(func(j *joins, stop error) { j.cleanupBoundary = atStep(StepApplyLocked, func() error { return stop }) })},
	{"landed", "later-failure", landedFailureCase(faultAtSecondLock)},
	{"landed", "outer-preflight-drift", landedOuterDriftCase},
	{"landed", "locked-drift", landedLockedDriftCase},
	{"landed", "missing-tree-release", landedMissingTreeCase},
	{"unclaimed", "empty", unclaimedCase(nil, nil)},
	{"unclaimed", "retained-only", unclaimedCase([]string{"e"}, nil)},
	{"unclaimed", "mixed", unclaimedCase([]string{"a", "e"}, nil)},
	{"unclaimed", "all-success", unclaimedCase([]string{"a", "b"}, nil)},
	{"unclaimed", "first-failure", unclaimedCase([]string{"a", "b"}, lockRefOf(0))},
	{"unclaimed", "later-failure", unclaimedCase([]string{"a", "b", "c"}, lockRefOf(1))},
	{"unclaimed", "selection-drift", unclaimedCase([]string{"a", "b"}, func(t *testing.T, root string, _ unclaimedAssignmentSet) { addUnclaimedBranch(t, root, "c") })},
}

// TestCleanupTraversalBaseline is R26. Every manifest case must reproduce its frozen record,
// and the record directory must hold exactly one record per case.
func TestCleanupTraversalBaseline(t *testing.T) {
	t.Parallel()
	capture := os.Getenv(cleanupBaselineCaptureEnv) == "1"
	want := make(map[string]bool, len(cleanupBaselineManifest))
	for _, c := range cleanupBaselineManifest {
		want[c.mode+"-"+c.name+".txt"] = true
	}
	if !capture {
		entries, err := os.ReadDir(cleanupBaselineDir)
		mustNoError(t, err)
		for _, entry := range entries {
			if !want[entry.Name()] {
				t.Errorf("record %s names no manifest case", entry.Name())
			}
		}
	}
	for _, c := range cleanupBaselineManifest {
		t.Run(c.mode+"/"+c.name, func(t *testing.T) {
			t.Parallel()
			got := renderBaselineRecord(t, c)
			path := filepath.Join(cleanupBaselineDir, c.mode+"-"+c.name+".txt")
			if capture {
				mustNoError(t, os.MkdirAll(cleanupBaselineDir, 0o755))
				mustNoError(t, os.WriteFile(path, []byte(got), 0o644))
				return
			}
			frozen, err := os.ReadFile(path)
			mustNoError(t, err)
			if got != string(frozen) {
				t.Fatalf("%s differs from its frozen record at %s\n--- got\n%s--- frozen\n%s", c.mode+"/"+c.name, firstRecordDifference(got, string(frozen)), got, frozen)
			}
		})
	}
}

// baselineAmbient is the fixture's ambient at the fixed plan day, so a discarded ref name
// carries no wall-clock date.
func baselineAmbient(f repoHome) ambient {
	a := f.ambient()
	a.now = discardDay
	return a
}

// renderFingerprintInputs renders a production input list, so an input the fingerprint omits
// leaves the record, and the plan digest must be the production hash of exactly these inputs.
// Only the empty landed selection plans no digest; the frozen record pins that absence.
func renderFingerprintInputs(t *testing.T, digest string, inputs []explicitFingerprintInput, version string) (lines []string) {
	t.Helper()
	hash := fingerprintInputs(version, inputs)
	requireTest(t, digest == "" || hash == digest, "plan digest %s is not the fingerprint %s of its rendered inputs", digest, hash)
	for _, input := range inputs {
		fields := append([]string{input.kind}, input.values...)
		for i, key := range input.keys {
			fields[i+1] = key + "=" + fields[i+1]
		}
		lines = append(lines, strings.Join(fields, " "))
	}
	return lines
}

// --- explicit and unrecorded cases: applyExplicitSet and its unrecorded stage ---

// explicitRun plans targets through the explicit planner. prepare runs after the plan and
// before the apply, so it can drift a member or install a fault for the planned set.
func explicitRun(t *testing.T, f repoHome, byRole bool, targets []string, options CleanupOptions, prepare func(*joins, explicitCleanupSet)) baselineRun {
	t.Helper()
	j := defaultJoins()
	a := baselineAmbient(f)
	set := planExplicitSet(j, a, f.root, targets, options)
	if len(set.failures) != 0 {
		t.Fatalf("explicit set failures = %#v, want a resolved selection", set.failures)
	}
	run := baselineRun{repoHome: f, byRole: byRole, planDigest: set.fingerprint, targetDigests: map[string]string{}}
	run.inputs = renderFingerprintInputs(t, set.fingerprint, explicitFingerprintInputs(set, options), explicitSetFingerprintVersion)
	for _, row := range set.rows {
		run.targetDigests[row.targetFingerprint] = row.assignment.ID
	}
	for _, plan := range set.plans() {
		run.planned = append(run.planned, plan.Target)
	}
	if prepare != nil {
		prepare(&j, set)
	}
	run.apply = func() ([]CleanupPlan, error) { return applyExplicitSet(j, a, f.root, set, options) }
	return run
}

func explicitEmptyCase(t *testing.T) baselineRun {
	root := newWorktreeRepo(t)
	return explicitRun(t, repoHome{root, filepath.Join(root, ".bench-home")}, false, nil, CleanupOptions{}, nil)
}

func explicitRetainedOnlyCase(t *testing.T) baselineRun {
	f := retainedMemberFixture(t)
	return explicitRun(t, f.repoHome, true, []string{f.retained.Assignment.ID}, CleanupOptions{}, nil)
}

func explicitMixedCase(t *testing.T) baselineRun {
	f := retainedMemberFixture(t)
	return explicitRun(t, f.repoHome, true, []string{f.removable.Assignment.ID, f.retained.Assignment.ID}, CleanupOptions{}, nil)
}

func explicitAllSuccessCase(t *testing.T) baselineRun {
	f := removableSetFixture(t, 2)
	return explicitRun(t, f.repoHome, false, explicitIdentities(f.creations), CleanupOptions{}, nil)
}

// explicitFailureCase stops a set of three at the transaction fault inject installs.
func explicitFailureCase(inject func(*joins, error)) func(*testing.T) baselineRun {
	return func(t *testing.T) baselineRun {
		f := removableSetFixture(t, 3)
		return explicitRun(t, f.repoHome, false, explicitIdentities(f.creations), CleanupOptions{}, func(j *joins, _ explicitCleanupSet) {
			inject(j, errors.New("baseline transaction fault"))
		})
	}
}

func explicitOuterDriftCase(t *testing.T) baselineRun {
	f := removableSetFixture(t, 2)
	return explicitRun(t, f.repoHome, false, explicitIdentities(f.creations), CleanupOptions{}, func(_ *joins, set explicitCleanupSet) {
		driftTracked(t, f.files, set.rows[1].assignment.ID)
	})
}

func explicitLockedDriftCase(t *testing.T) baselineRun {
	f := removableSetFixture(t, 2)
	return explicitRun(t, f.repoHome, false, explicitIdentities(f.creations), CleanupOptions{}, func(j *joins, set explicitCleanupSet) {
		j.cleanupBoundary = atSecondHit(StepApplyLocked, func() error {
			driftTracked(t, f.files, set.rows[1].assignment.ID)
			return nil
		})
	})
}

// branchTargets names each unclaimed ref by the selector its discard route prints.
func branchTargets(refs ...string) []string {
	targets := make([]string, 0, len(refs))
	for _, ref := range refs {
		targets = append(targets, discardSelector(ref))
	}
	return targets
}

func unrecordedRetainedOnlyCase(t *testing.T) baselineRun {
	root, home := unclaimedBranchFixture(t)
	ref, _ := uniqueBranch(t, root, "e")
	return explicitRun(t, repoHome{root, home}, false, branchTargets(ref), CleanupOptions{}, nil)
}

// unrecordedMixedCase selects a retained recorded member and a removable landed branch.
func unrecordedMixedCase(t *testing.T) baselineRun {
	f := retainedMemberFixture(t)
	ref := addUnclaimedBranch(t, f.root, "a")
	return explicitRun(t, f.repoHome, true, append([]string{f.retained.Assignment.ID}, branchTargets(ref)...), CleanupOptions{DiscardBranch: true}, nil)
}

// unrecordedAllSuccessCase removes a landed branch with no recovery ref and a unique branch
// through its discarded ref.
func unrecordedAllSuccessCase(t *testing.T) baselineRun {
	root, home := unclaimedBranchFixture(t, "a")
	unique, _ := uniqueBranch(t, root, "e")
	return explicitRun(t, repoHome{root, home}, false, branchTargets(unclaimedBranchRef("a"), unique), CleanupOptions{DiscardBranch: true}, nil)
}

// unrecordedFailureCase stops count unique branches at the discard fault fault installs.
func unrecordedFailureCase(count int, fault func(error) Fault) func(*testing.T) baselineRun {
	return func(t *testing.T) baselineRun {
		root, home := unclaimedBranchFixture(t)
		refs := make([]string, 0, count)
		for _, owner := range []string{"a", "b", "c"}[:count] {
			ref, _ := uniqueBranch(t, root, owner)
			refs = append(refs, ref)
		}
		return explicitRun(t, repoHome{root, home}, false, branchTargets(refs...), CleanupOptions{DiscardBranch: true}, func(j *joins, _ explicitCleanupSet) {
			j.cleanupBoundary = fault(errors.New("baseline discard fault"))
		})
	}
}

// unrecordedHolderRemovalCase selects a recorded holder and the branch it subsumes. The
// holder's removal takes the proof away, so the branch must survive.
func unrecordedHolderRemovalCase(t *testing.T) baselineRun {
	root, home := unclaimedBranchFixture(t)
	ref, holder := recordedBranch(t, root, home, intent.StateActive, false)
	return explicitRun(t, repoHome{root, home}, false, append([]string{branchSegment(holder)}, branchTargets(ref)...), CleanupOptions{DiscardBranch: true}, nil)
}

// unrecordedRecoveryConflictCase plants a ref at the planned discarded path before the plan.
func unrecordedRecoveryConflictCase(t *testing.T) baselineRun {
	root, home := unclaimedBranchFixture(t)
	ref, _ := uniqueBranch(t, root, "e")
	gitRun(t, root, "update-ref", intent.DiscardedRef(discardDay, ref), gitOutput(t, root, "rev-parse", "main"))
	return explicitRun(t, repoHome{root, home}, false, branchTargets(ref), CleanupOptions{DiscardBranch: true}, nil)
}

// --- landed cases: applyLandedSet ---

// landedRun plans the whole-repository landed set. prepare runs after the plan, before the apply.
func landedRun(t *testing.T, f repoHome, byRole bool, prepare func(*joins, landedCleanupSet)) baselineRun {
	t.Helper()
	j := defaultJoins()
	options := CleanupOptions{}
	set, err := planLandedSet(j, nil, f.root, options, "")
	mustNoError(t, err)
	run := baselineRun{repoHome: f, byRole: byRole, planDigest: set.fingerprint}
	run.inputs = renderFingerprintInputs(t, set.fingerprint, landedFingerprintInputs(set.rows, options), landedSetFingerprintVersion)
	for _, row := range set.rows {
		run.planned = append(run.planned, row.plan.Target)
	}
	if prepare != nil {
		prepare(&j, set)
	}
	a := baselineAmbient(f)
	run.apply = func() ([]CleanupPlan, error) { return applyLandedSet(j, a, f.root, set, options, "") }
	return run
}

func landedEmptyCase(t *testing.T) baselineRun {
	root := newWorktreeRepo(t)
	return landedRun(t, repoHome{root, filepath.Join(root, ".bench-home")}, false, nil)
}

func landedRetainedOnlyCase(t *testing.T) baselineRun {
	root, home := ignoringRepo(t)
	retained := landedMember(t, root, home, "set-retained-residue", "retained.txt")
	mustWrite(t, filepath.Join(retained.Path, "ignored-one.txt"), []byte("residue\n"), 0o644)
	return landedRun(t, repoHome{root, home}, false, nil)
}

func landedMixedCase(t *testing.T) baselineRun {
	return landedRun(t, retainedMemberFixture(t).repoHome, true, nil)
}

func landedAllSuccessCase(t *testing.T) baselineRun {
	return landedRun(t, removableSetFixture(t, 2).repoHome, false, nil)
}

// landedFailureCase stops a set of three at the transaction fault inject installs.
func landedFailureCase(inject func(*joins, error)) func(*testing.T) baselineRun {
	return func(t *testing.T) baselineRun {
		return landedRun(t, removableSetFixture(t, 3).repoHome, false, func(j *joins, _ landedCleanupSet) {
			inject(j, errors.New("baseline transaction fault"))
		})
	}
}

func landedOuterDriftCase(t *testing.T) baselineRun {
	f := removableSetFixture(t, 2)
	return landedRun(t, f.repoHome, false, func(_ *joins, set landedCleanupSet) {
		driftTracked(t, f.files, set.rows[1].assignment.ID)
	})
}

func landedLockedDriftCase(t *testing.T) baselineRun {
	f := removableSetFixture(t, 2)
	return landedRun(t, f.repoHome, false, func(j *joins, set landedCleanupSet) {
		j.cleanupBoundary = atSecondHit(StepApplyLocked, func() error {
			driftTracked(t, f.files, set.rows[1].assignment.ID)
			return nil
		})
	})
}

// landedMissingTreeCase removes the tree of the lowest identity before the plan, so that row
// plans a release, which the landed preflight skips and the landed visits apply.
func landedMissingTreeCase(t *testing.T) baselineRun {
	f := removableSetFixture(t, 2)
	ids := explicitIdentities(f.creations)
	sort.Strings(ids)
	mustNoError(t, os.RemoveAll(memberByID(t, f.creations, ids[0]).Path))
	return landedRun(t, f.repoHome, false, nil)
}

// --- unclaimed cases: applyUnclaimedAssignmentSet ---

// unclaimedCase plants one branch per owner letter: "e" plants a unique branch and every other
// letter a landed one. after runs between the plan and the apply.
func unclaimedCase(owners []string, after func(*testing.T, string, unclaimedAssignmentSet)) func(*testing.T) baselineRun {
	return func(t *testing.T) baselineRun {
		root, home := unclaimedBranchFixture(t)
		for _, owner := range owners {
			if owner == "e" {
				uniqueBranch(t, root, owner)
				continue
			}
			addUnclaimedBranch(t, root, owner)
		}
		options := unclaimedOptions()
		set, err := planUnclaimedAssignmentSet(root, options)
		mustNoError(t, err)
		run := baselineRun{repoHome: repoHome{root, home}, planDigest: set.fingerprint}
		run.inputs = renderFingerprintInputs(t, set.fingerprint, unclaimedFingerprintInputs(set.rows, options), unclaimedAssignmentFingerprintVersion)
		for i, row := range set.rows { // a faulted set has no digest, so the record pins each fault
			run.inputs[i+1] += " " + row.fault
			run.planned = append(run.planned, row.ref)
		}
		if after != nil {
			after(t, root, set)
		}
		run.apply = func() ([]CleanupPlan, error) { return applyUnclaimedAssignmentSet(defaultJoins(), root, set, options) }
		return run
	}
}

// lockRefOf holds a Git ref lock on the planned removable ref at index. Selection planning
// reads past the lock, so only the exact deletion of that ref fails.
func lockRefOf(index int) func(*testing.T, string, unclaimedAssignmentSet) {
	return func(t *testing.T, root string, set unclaimedAssignmentSet) {
		mustWrite(t, filepath.Join(root, ".git", filepath.FromSlash(set.rows[index].ref)+".lock"), nil, 0o644)
	}
}
