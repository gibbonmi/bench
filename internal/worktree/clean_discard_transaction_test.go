// The discard transaction's refusals: a handle at the planned discarded path that the write
// must not replace or follow, a candidate the class function faulted, and a member whose
// holder an earlier member of the same set removed.
package worktree

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestDiscardTargetRefusesASymrefAtTheDiscardedPath is RI96: a symref at the planned path
// that points to the target branch reads as that branch's tip through a dereferencing read.
// The apply must refuse it, keep both refs, and write nothing.
func TestDiscardTargetRefusesASymrefAtTheDiscardedPath(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	f := repoHome{root, home}
	ref, tip := uniqueBranch(t, root, "a")
	discarded := intent.DiscardedRef(discardDay, ref)
	gitRun(t, root, "symbolic-ref", discarded, ref)
	args := []string{"--discard-branch", "--target", strings.TrimPrefix(ref, "refs/heads/")}
	plan := runVerb(t, verbClean, f.callWith(discardJoins(), args...))
	applied := runVerb(t, verbClean, f.callWith(discardJoins(), append(args, "--apply", plan.mustFingerprint(t))...))
	rows := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")
	if plan.stderr != "" || applied.stderr != "" || applied.exit != 1 || rowAction(rows, ref) != string(ActionError) || refTip(root, ref) != tip || gitOutput(t, root, "symbolic-ref", "--quiet", discarded) != ref {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q, want an error row, the branch at %s, and the symref unchanged", applied.exit, applied.stdout, plan.stderr+applied.stderr, tip)
	}
}

// TestDiscardTargetRefusesARefPlantedAfterTheRead is RI101: a direct ref that appears at the
// planned path after the read found it absent makes the write fail on its zero old value.
func TestDiscardTargetRefusesARefPlantedAfterTheRead(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	f := repoHome{root, home}
	ref, tip := uniqueBranch(t, root, "a")
	discarded := intent.DiscardedRef(discardDay, ref)
	planted := gitOutput(t, root, "rev-parse", "main")
	j := discardJoins()
	j.cleanupBoundary = atStep(StepDiscardedRefAbsent, func() error {
		gitRun(t, root, "update-ref", discarded, planted)
		return nil
	})
	args := []string{"--discard-branch", "--target", strings.TrimPrefix(ref, "refs/heads/")}
	plan := runVerb(t, verbClean, f.callWith(discardJoins(), args...))
	applied := runVerb(t, verbClean, f.callWith(j, append(args, "--apply", plan.mustFingerprint(t))...))
	rows := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")
	if plan.stderr != "" || applied.stderr != "" || applied.exit != 1 || rowAction(rows, ref) != string(ActionError) || refTip(root, ref) != tip || refTip(root, discarded) != planted {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q, want an error row, the branch at %s, and the planted ref at %s", applied.exit, applied.stdout, plan.stderr+applied.stderr, tip, planted)
	}
}

// TestDiscardTargetNeverFollowsASymrefPlantedAfterTheRead is RI102: a symref that appears at
// the planned path after the read found it absent is never followed. A resolving symref fails
// the write and keeps the branch; a dangling one is replaced by the discarded ref at the row's
// tip, and no ref appears at its old target.
func TestDiscardTargetNeverFollowsASymrefPlantedAfterTheRead(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, target string
		resolves     bool
	}{
		{"resolving", "refs/heads/main", true},
		{"dangling", "refs/heads/bench/assign/zz/gone", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t)
			f := repoHome{root, home}
			ref, tip := uniqueBranch(t, root, "a")
			discarded := intent.DiscardedRef(discardDay, ref)
			targetTip := refTip(root, tc.target)
			j := discardJoins()
			j.cleanupBoundary = atStep(StepDiscardedRefAbsent, func() error {
				gitRun(t, root, "symbolic-ref", discarded, tc.target)
				return nil
			})
			args := []string{"--discard-branch", "--target", strings.TrimPrefix(ref, "refs/heads/")}
			plan := runVerb(t, verbClean, f.callWith(discardJoins(), args...))
			applied := runVerb(t, verbClean, f.callWith(j, append(args, "--apply", plan.mustFingerprint(t))...))
			rows := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")
			symref, _ := git.Output("-C", root, "symbolic-ref", "--quiet", discarded)
			if plan.stderr != "" || applied.stderr != "" {
				t.Fatalf("plan stderr=%q apply stderr=%q, want no stderr", plan.stderr, applied.stderr)
			}
			if tc.resolves && (applied.exit != 1 || rowAction(rows, ref) != string(ActionError) || refTip(root, ref) != tip || refTip(root, tc.target) != targetTip || symref != tc.target) {
				t.Fatalf("apply exit=%d stdout=%q, want an error row, the branch at %s, %s unmoved, and the symref kept", applied.exit, applied.stdout, tip, tc.target)
			}
			if !tc.resolves && (applied.exit != 0 || refTip(root, ref) != "" || symref != "" || refTip(root, discarded) != tip || refTip(root, tc.target) != "") {
				t.Fatalf("apply exit=%d stdout=%q, want the branch removed, %s a direct ref at %s, and no ref at %s", applied.exit, applied.stdout, discarded, tip, tc.target)
			}
		})
	}
}

// TestDiscardTargetNeverFollowsASymrefPlantedBeforeTheDelete is RI103: a symref planted at the
// branch path after the write, pointing at the discarded ref, passes an exact-tip check
// through its referent. The delete must not follow it, so the discarded ref survives at the
// row's tip.
func TestDiscardTargetNeverFollowsASymrefPlantedBeforeTheDelete(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	f := repoHome{root, home}
	ref, tip := uniqueBranch(t, root, "a")
	discarded := intent.DiscardedRef(discardDay, ref)
	j := discardJoins()
	j.cleanupBoundary = atStep(StepDiscardedBranchDelete, func() error {
		gitRun(t, root, "symbolic-ref", ref, discarded)
		return nil
	})
	args := []string{"--discard-branch", "--target", strings.TrimPrefix(ref, "refs/heads/")}
	plan := runVerb(t, verbClean, f.callWith(discardJoins(), args...))
	applied := runVerb(t, verbClean, f.callWith(j, append(args, "--apply", plan.mustFingerprint(t))...))
	rows := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")
	if plan.stderr != "" || applied.stderr != "" || refTip(root, discarded) != tip {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q, want %s kept at %s", applied.exit, applied.stdout, plan.stderr+applied.stderr, discarded, tip)
	}
	// Git checks the old value against the referent and deletes only the symref.
	if applied.exit != 0 || rowAction(rows, ref) != string(ActionRemoved) || refTip(root, ref) != "" {
		t.Fatalf("apply exit=%d stdout=%q, want a removed row with the branch path gone", applied.exit, applied.stdout)
	}
	if holders := gitOutput(t, root, "for-each-ref", "--format=%(refname)", "--points-at", tip); holders != discarded {
		t.Fatalf("refs at %s = %q, want only %s", tip, holders, discarded)
	}
}

// TestDiscardTargetRefusesAFaultedCandidate is RI98: a Bench assignment branch that is a
// resolving symref prints an error row that names the symref and no fingerprint, and neither
// apply form changes a ref.
func TestDiscardTargetRefusesAFaultedCandidate(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	f := repoHome{root, home}
	symref := unclaimedBranchRef("a")
	gitRun(t, root, "symbolic-ref", symref, "refs/heads/main")
	short := strings.TrimPrefix(symref, "refs/heads/")
	identities := func() string {
		return gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname) %(symref)")
	}
	refs := identities()
	plan := runVerb(t, verbClean, f.callWith(discardJoins(), "--discard-branch", "--target", short))
	if row := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")[short]; plan.exit != 1 || plan.stderr != "" || len(row) != len(cleanupFields) || row["action"] != string(ActionError) || !strings.Contains(row["detail"], "symref") {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q, want an error row naming the symref and no fingerprint", plan.exit, plan.stdout, plan.stderr)
	}
	plan.mustNoFingerprint(t)
	// No plan offers a fingerprint, so the fingerprint form carries a well-formed guess, and the
	// current form is outside the explicit grammar.
	for _, tc := range []struct {
		apply []string
		code  int
	}{{[]string{"--apply", strings.Repeat("a", 64)}, 1}, {[]string{"--apply-current"}, 2}} {
		applied := runVerb(t, verbClean, f.callWith(discardJoins(), append([]string{"--discard-branch", "--target", short}, tc.apply...)...))
		if applied.exit != tc.code || strings.Contains(applied.stdout, ",removed,") || identities() != refs {
			t.Fatalf("%s exit=%d stdout=%q stderr=%q, want exit %d and no ref change", tc.apply[0], applied.exit, applied.stdout, applied.stderr, tc.code)
		}
	}
}

// TestDiscardTargetRequalifiesAfterARecordedRemoval is RI99: an unrecorded member subsumed by a
// recorded member of the same set becomes unique once that member's removal takes its holder
// away, so its own requalify refuses it as stale and it survives at its tip.
func TestDiscardTargetRequalifiesAfterARecordedRemoval(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	f := repoHome{root, home}
	ref, holder := recordedBranch(t, root, home, intent.StateActive, false)
	tip := refTip(root, ref)
	args := []string{"--discard-branch", "--target", branchSegment(holder), "--target", strings.TrimPrefix(ref, "refs/heads/")}
	plan := runVerb(t, verbClean, f.callWith(discardJoins(), args...))
	if row := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")[ref]; plan.stderr != "" || len(row) != len(cleanupFields) || !strings.HasPrefix(row["detail"], classField+string(classSubsumed)+holderField+holder) {
		t.Fatalf("plan = %q stderr=%q, want %s subsumed by %s", plan.stdout, plan.stderr, ref, holder)
	}
	applied := runVerb(t, verbClean, f.callWith(discardJoins(), append(args, "--apply", plan.mustFingerprint(t))...))
	rows := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")
	if applied.stderr != "" || refTip(root, holder) != "" {
		t.Fatalf("apply stdout=%q stderr=%q, want the recorded holder %s removed first", applied.stdout, applied.stderr, holder)
	}
	if applied.exit != 1 || !strings.Contains(applied.stdout, errStaleFingerprint.Error()) || rowAction(rows, ref) == string(ActionRemoved) || refTip(root, ref) != tip || discardedRefs(t, root) != "" {
		t.Fatalf("apply exit=%d stdout=%q, want a stale refusal for %s, the branch at %s, and no discarded ref", applied.exit, applied.stdout, ref, tip)
	}
}
