package worktree

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestExplicitEligibilityAllowsRuntimeIgnoredResidue(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte(".logs/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore runtime records")
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, Home(), "runtime-eligibility", "runtime eligibility")
	mustMkdirAll(t, filepath.Join(creation.Path, ".logs"), 0o755)
	mustWrite(t, filepath.Join(creation.Path, ".logs", "gate.jsonl"), []byte("record\n"), 0o644)
	plan, err := PlanExplicitWithOptions(root, creation.Path, CleanupOptions{})
	if err != nil || plan.Action != ActionDiscardRemove || plan.ReasonCode != "" {
		t.Fatalf("runtime residue plan = (%#v, %v), want discard-remove without refusal", plan, err)
	}
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
	t.Parallel()
	t.Run("clean-remove", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "ev1-clean-remove")
		assertVerdictMatchesPlan(t, root, creation.Path, CleanupOptions{})
	})
	t.Run("dirty-recover-remove", func(t *testing.T) {
		root, creation, _ := newOwnedAssignment(t, "ev1-dirty-recover-remove")
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
		RegistrationBranchRef:  registration.BranchRef,
		RegistrationLockReason: registration.LockReason,
		RegistrationLocked:     registration.Locked,
		RegistrationDetached:   registration.Detached,
		DiscardIgnored:         options.DiscardIgnored,
	}
	_, markerStatErr := os.Lstat(filepath.Join(admin, OwnerMarkerFile))
	if markerStatErr == nil {
		facts.MarkerPresent = true
		evidence, markerErr := validateOwnerMarker(root, target)
		if markerErr != nil {
			facts.MarkerErr = markerErr
		} else {
			assignments, assignmentErr := intent.Assignments(root)
			if assignmentErr != nil {
				facts.AssignmentLedgerErr = assignmentErr
			} else {
				var matched *intent.Assignment
				for i := range assignments {
					if assignments[i].Worktree == target && assignments[i].OwnerID == evidence.marker.OwnerID {
						if matched != nil {
							facts.AssignmentAmbiguous = true
							break
						}
						candidate := assignments[i]
						matched = &candidate
					}
				}
				facts.MatchedAssignment = matched
				if matched != nil {
					facts.AssignmentLockReason = lockReason(*matched)
				}
			}
		}
	} else if !errors.Is(markerStatErr, os.ErrNotExist) {
		t.Fatalf("stat owner marker: %v", markerStatErr)
	} else if !registration.Locked {
		facts.ForeignAssignment = foreignRecoveryAssignment(root, target)
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
	facts.InitialTracked = tracked

	leasePath, err := LeaseFile(target)
	mustNoError(t, err)
	if _, statErr := os.Lstat(leasePath); statErr == nil {
		facts.LeasePresent = true
		facts.LeaseState = ProbeLease(leasePath)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		facts.LeaseStatErr = statErr
	}

	nested, nestedErr := classifyNestedState(target)
	facts.NestedState, facts.NestedErr = nested, nestedErr

	buildOutputs, _, buildOutputErr := loadBuildOutputs(root)
	ignored, _, ignoredErr := inventoryIgnored(target, options.Full)
	facts.BuildOutputErr = buildOutputErr
	facts.IgnoredErr = ignoredErr
	facts.IgnoredOverLimit = ignored.OverLimit
	facts.IgnoredCount = ignored.Count
	facts.DeclaredIgnored = buildOutputErr == nil && ignoredWithinLandingAllowance(ignored, buildOutputs)

	defaultRef, defaultOID := "none", "none"
	if def, ok := git.ResolvedDefault(root); ok {
		defaultRef = def
		if oid, oidErr := git.Output("-C", root, "rev-parse", "--verify", def+"^{commit}"); oidErr == nil {
			defaultOID = oid
		}
	}
	facts.HeadDetached = headRef == "detached"
	facts.DefaultKnown = defaultOID != "none"
	facts.HeadRef, facts.Head = headRef, head
	if !facts.HeadDetached && facts.DefaultKnown {
		facts.LandedOK, facts.LandedByContent, facts.LandedErr = git.LandedInDefault(root, strings.TrimPrefix(headRef, "refs/heads/"), defaultRef)
	}
	facts.UnsafeTarget = !cleanupOutputSafe(target)
	return facts
}
