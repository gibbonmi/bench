// The explicit discard of one unrecorded branch. Each fixture drives the clean verb's joins
// form over a real repository under a fixed clock, reads the decoded rows, and reads the refs
// back.
package worktree

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// discardDay is the fixed plan day. It is far from any real run, so a planner that reads the
// wall clock names a different discarded ref.
var discardDay = time.Date(2031, time.February, 3, 15, 4, 5, 0, time.UTC)

// discardCall is call with its clock value fixed at discardDay.
func discardCall(call verbCall) verbCall {
	call.clock = func() time.Time { return discardDay }
	return call
}

// atStep runs act each time the transaction reaches step.
func atStep(step LifecycleStep, act func() error) Fault {
	return func(reached LifecycleStep) error {
		if reached != step {
			return nil
		}
		return act()
	}
}

// uniqueBranch plants one unclaimed assignment branch with a commit main lacks, and returns
// the ref and its tip.
func uniqueBranch(t *testing.T, root, owner string) (string, string) {
	t.Helper()
	ref := addUnclaimedBranch(t, root, owner)
	return ref, commitOnBranch(t, root, ref, owner+".txt", owner+"\n")
}

// uniqueShiftBranch plants one unclaimed shift branch with a commit main lacks.
func uniqueShiftBranch(t *testing.T, root string) (string, string) {
	t.Helper()
	ref := intent.ShiftBranchPrefix() + "20310203-150405"
	gitRun(t, root, "branch", strings.TrimPrefix(ref, "refs/heads/"))
	return ref, commitOnBranch(t, root, ref, "shift.txt", "shift\n")
}

// refTip is the object ref names, or empty when the ref does not exist.
func refTip(root, ref string) string {
	tip, _ := git.Output("-C", root, "rev-parse", "--verify", "--quiet", ref)
	return tip
}

// discardedRefs lists every ref under the discarded namespace.
func discardedRefs(t *testing.T, root string) string {
	t.Helper()
	return gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)", intent.DiscardedRefNamespace)
}

// rowAction is the action cell of the row that names target, or empty when no row does.
func rowAction(rows map[string]textRow, target string) string {
	if row := rows[target]; len(row) == len(cleanupFields) {
		return row["action"]
	}
	return ""
}

// TestDiscardTargetResolvesAnUnrecordedBranch is RI29, RI30, and RI33: each operand form plans
// one unique row with its class and the discarded ref dated by the call's clock value.
func TestDiscardTargetResolvesAnUnrecordedBranch(t *testing.T) {
	t.Parallel()
	short := func(ref string) string { return strings.TrimPrefix(ref, "refs/heads/") }
	assignment := func(t *testing.T, root string) string { ref, _ := uniqueBranch(t, root, "a"); return ref }
	shift := func(t *testing.T, root string) string { ref, _ := uniqueShiftBranch(t, root); return ref }
	for _, tc := range []struct {
		name     string
		plant    func(*testing.T, string) string
		operands func(string) []string
	}{
		{"RI29 the full branch path", assignment, func(ref string) []string { return []string{ref} }},
		{"RI29 the branch path without refs/heads/", assignment, func(ref string) []string { return []string{short(ref)} }},
		{"RI29 a shift branch path", shift, func(ref string) []string { return []string{short(ref)} }},
		{"RI30 the id segment", assignment, func(ref string) []string { return []string{branchSegment(ref)} }},
		{"RI30 the id and the branch path in one call", assignment, func(ref string) []string { return []string{branchSegment(ref), short(ref)} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t)
			ref := tc.plant(t, root)
			args := []string{"--discard-branch"}
			for _, operand := range tc.operands(ref) {
				args = append(args, "--target", operand)
			}
			plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
			rows := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")
			row := rows[ref]
			if plan.exit != 0 || plan.stderr != "" || len(rows) != 1 || len(row) != len(cleanupFields) || row["action"] != string(ActionDiscardRemove) ||
				row["recovery"] != intent.DiscardedRef(discardDay, ref) || !strings.HasPrefix(row["detail"], classField+string(classUnique)) {
				t.Fatalf("plan exit=%d stdout=%q stderr=%q, want one unique row for %s with its discarded ref", plan.exit, plan.stdout, plan.stderr, ref)
			}
			if !strings.Contains(plan.stdout, "--apply "+plan.mustFingerprint(t)) {
				t.Fatalf("plan = %q, want the apply action", plan.stdout)
			}
		})
	}
}

// TestDiscardTargetRefusesAnAmbiguousOrForeignOperand is RI31 and RI32: a shared id segment
// names both refs and plans no fingerprint, a foreign branch keeps the path refusal, and an
// id that matches nothing keeps the unassigned refusal.
func TestDiscardTargetRefusesAnAmbiguousOrForeignOperand(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t, "a", "b")
	gitRun(t, root, "branch", "archive/x")
	for _, tc := range []struct {
		name, operand string
		detail        func(string) bool
	}{
		{"RI31 a shared id segment", branchSegment(unclaimedBranchRef("a")), func(detail string) bool {
			return strings.Contains(detail, unclaimedBranchRef("a")) && strings.Contains(detail, unclaimedBranchRef("b"))
		}},
		{"RI32 a foreign branch", "archive/x", func(detail string) bool { return detail == "relative path targets are unsupported" }},
		{"an id that matches nothing", strings.Repeat("0", 32), func(detail string) bool { return detail == errTargetUnassigned.Error() }},
	} {
		plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call("--discard-branch", "--target", tc.operand)))
		rows := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")
		row := rows[tc.operand]
		if plan.exit != 1 || plan.stderr != "" || len(rows) != 1 || len(row) != len(cleanupFields) || row["action"] != string(ActionError) || !tc.detail(row["detail"]) {
			t.Fatalf("%s: plan exit=%d stdout=%q stderr=%q, want one error row and no fingerprint", tc.name, plan.exit, plan.stdout, plan.stderr)
		}
		plan.mustNoFingerprint(t)
	}
}

// TestDiscardTargetPrefixReadsTheCandidateSet is RI74, RI75, RI76, and RI41: a prefix operand
// for a checkout or the default branch refuses, and a recorded branch, by prefix, by id, or by
// a label that matches an unclaimed id, plans through its record.
func TestDiscardTargetPrefixReadsTheCandidateSet(t *testing.T) {
	t.Parallel()
	t.Run("RI74 a checked-out shift branch", func(t *testing.T) {
		t.Parallel()
		root, home := unclaimedBranchFixture(t)
		short := strings.TrimPrefix(intent.ShiftBranchPrefix(), "refs/heads/") + "20310203-150405"
		checkout := filepath.Join(t.TempDir(), "checked-out-shift")
		gitRun(t, root, "branch", short)
		gitRun(t, root, "worktree", "add", "-q", checkout, short)
		plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call("--discard-branch", "--target", short)))
		if row := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")[short]; plan.exit != 1 || plan.stderr != "" || len(row) != len(cleanupFields) || row["action"] != string(ActionError) || !strings.Contains(row["detail"], "checked-out-shift") {
			t.Fatalf("plan exit=%d stdout=%q stderr=%q, want an error row naming the checkout and no fingerprint", plan.exit, plan.stdout, plan.stderr)
		}
		plan.mustNoFingerprint(t)
	})
	t.Run("RI75 the default branch", func(t *testing.T) {
		t.Parallel()
		root, home := unclaimedBranchFixture(t)
		short := strings.TrimPrefix(intent.AssignmentBranchRef(strings.Repeat("d", 32), strings.Repeat("e", 32)), "refs/heads/")
		gitRun(t, root, "branch", "-m", short)
		gitRun(t, root, "checkout", "--detach", "-q")
		plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call("--discard-branch", "--target", short)))
		if row := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")[short]; plan.exit != 1 || plan.stderr != "" || len(row) != len(cleanupFields) || row["action"] != string(ActionError) || !strings.Contains(row["detail"], "refs/heads/"+short) {
			t.Fatalf("plan exit=%d stdout=%q stderr=%q, want an error row naming the branch and no fingerprint", plan.exit, plan.stdout, plan.stderr)
		}
		plan.mustNoFingerprint(t)
	})
	for _, tc := range []struct {
		name     string
		label    string
		operands func(Creation) []string
	}{
		{"RI76 a recorded branch by prefix and by id", "discard recorded", func(c Creation) []string {
			return []string{"--target", strings.TrimPrefix(c.Assignment.Branch, "refs/heads/"), "--target", c.Assignment.ID}
		}},
		{"RI41 a recorded label that matches an unclaimed id", branchSegment(unclaimedBranchRef("a")), func(c Creation) []string {
			return []string{"--target", c.Assignment.Label}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t, "a")
			recorded := mustCreate(t, root, home, "discard-recorded", tc.label)
			plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(tc.operands(recorded)...)))
			rows := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")
			if row := rows[recorded.Path]; plan.stderr != "" || len(rows) != 1 || len(row) != len(cleanupFields) || strings.Contains(row["detail"], classField) {
				t.Fatalf("plan = %q stderr=%q, want one row through the record of %s and no class", plan.stdout, plan.stderr, recorded.Path)
			}
		})
	}
}

// TestDiscardTargetWithoutTheFlagRetains is RI58: the row retains, its detail names the flag,
// and its apply leaves the branch and writes no discarded ref.
func TestDiscardTargetWithoutTheFlagRetains(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	ref, tip := uniqueBranch(t, root, "a")
	plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call("--target", branchSegment(ref))))
	if row := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")[ref]; plan.exit != 0 || plan.stderr != "" || len(row) != len(cleanupFields) || row["action"] != string(ActionRetain) || row["recovery"] != "none" || !strings.Contains(row["detail"], "--discard-branch") {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q, want a retained row that names --discard-branch", plan.exit, plan.stdout, plan.stderr)
	}
	applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.call("--target", branchSegment(ref), "--apply", plan.mustFingerprint(t))))
	if applied.exit != 0 || applied.stderr != "" || strings.Contains(applied.stdout, ",removed,") || refTip(root, ref) != tip || discardedRefs(t, root) != "" {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q, want the branch kept and no discarded ref", applied.exit, applied.stdout, applied.stderr)
	}
}

// TestDiscardTargetWritesTheDiscardedRefFirst is RI35 and RI39: the discarded ref holds the tip
// before the delete runs, and the outcome row names it.
func TestDiscardTargetWritesTheDiscardedRefFirst(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	ref, tip := uniqueBranch(t, root, "a")
	discarded := intent.DiscardedRef(discardDay, ref)
	j := defaultJoins()
	var beforeDelete string
	j.cleanupBoundary = atStep(StepDiscardedBranchDelete, func() error {
		beforeDelete = refTip(root, discarded) + " " + refTip(root, ref)
		return nil
	})
	args := []string{"--discard-branch", "--target", branchSegment(ref)}
	plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
	applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.callWith(j, append(args, "--apply", plan.mustFingerprint(t))...)))
	if beforeDelete != tip+" "+tip {
		t.Fatalf("before the delete the discarded ref and the branch read %q, want both at %s", beforeDelete, tip)
	}
	if row := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")[ref]; plan.stderr != "" || applied.exit != 0 || applied.stderr != "" || len(row) != len(cleanupFields) || row["action"] != string(ActionRemoved) || row["recovery"] != discarded {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q, want a removed row whose recovery is %s", applied.exit, applied.stdout, plan.stderr+applied.stderr, discarded)
	}
	if refTip(root, discarded) != tip || refTip(root, ref) != "" {
		t.Fatalf("after the apply %s reads %q and %s reads %q", discarded, refTip(root, discarded), ref, refTip(root, ref))
	}
}

// TestDiscardTargetFaultWindows is RI36, RI37, and RI67: a fault before the write leaves no
// discarded ref, a fault after it leaves both refs for a retry, and a branch that moves after
// the write survives the exact delete.
func TestDiscardTargetFaultWindows(t *testing.T) {
	t.Parallel()
	stop := errors.New("stop in the discard window")
	fault := func(step LifecycleStep) joins {
		j := defaultJoins()
		j.cleanupBoundary = atStep(step, func() error { return stop })
		return j
	}
	t.Run("RI36 a fault before the write", func(t *testing.T) {
		t.Parallel()
		root, home := unclaimedBranchFixture(t)
		ref, tip := uniqueBranch(t, root, "a")
		args := []string{"--discard-branch", "--target", branchSegment(ref)}
		plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
		applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.callWith(fault(StepDiscardedRefWrite), append(args, "--apply", plan.mustFingerprint(t))...)))
		if plan.stderr != "" || applied.stderr != "" || applied.exit != 1 || refTip(root, ref) != tip || discardedRefs(t, root) != "" {
			t.Fatalf("apply exit=%d stdout=%q stderr=%q discarded=%q, want the branch kept and no discarded ref", applied.exit, applied.stdout, plan.stderr+applied.stderr, discardedRefs(t, root))
		}
	})
	t.Run("RI37 a fault after the write", func(t *testing.T) {
		t.Parallel()
		root, home := unclaimedBranchFixture(t)
		ref, tip := uniqueBranch(t, root, "a")
		discarded := intent.DiscardedRef(discardDay, ref)
		args := []string{"--discard-branch", "--target", branchSegment(ref)}
		plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
		applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.callWith(fault(StepDiscardedBranchDelete), append(args, "--apply", plan.mustFingerprint(t))...)))
		if plan.stderr != "" || applied.stderr != "" || applied.exit != 1 || refTip(root, ref) != tip || refTip(root, discarded) != tip {
			t.Fatalf("apply exit=%d stdout=%q stderr=%q, want the branch and the discarded ref both at %s", applied.exit, applied.stdout, plan.stderr+applied.stderr, tip)
		}
		replan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
		retried := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(append(args, "--apply", replan.mustFingerprint(t))...)))
		rows := textRowsBy(t, retried.mustRows(t, cleanupTable), "target")
		if replan.stderr != "" || retried.stderr != "" || retried.exit != 0 || rowAction(rows, ref) != string(ActionRemoved) || refTip(root, ref) != "" || refTip(root, discarded) != tip {
			t.Fatalf("retry exit=%d stdout=%q stderr=%q, want the branch removed and the discarded ref at %s", retried.exit, retried.stdout, replan.stderr+retried.stderr, tip)
		}
	})
	t.Run("RI67 a branch moved after the write", func(t *testing.T) {
		t.Parallel()
		root, home := unclaimedBranchFixture(t)
		ref, tip := uniqueBranch(t, root, "a")
		var moved string
		j := defaultJoins()
		j.cleanupBoundary = atStep(StepDiscardedBranchDelete, func() error {
			moved = commitOnBranch(t, root, ref, "moved.txt", "moved\n")
			return nil
		})
		args := []string{"--discard-branch", "--target", branchSegment(ref)}
		plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
		applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.callWith(j, append(args, "--apply", plan.mustFingerprint(t))...)))
		rows := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")
		if plan.stderr != "" || applied.stderr != "" || applied.exit != 1 || rowAction(rows, ref) != string(ActionError) || moved == "" || refTip(root, ref) != moved || refTip(root, intent.DiscardedRef(discardDay, ref)) != tip {
			t.Fatalf("apply exit=%d stdout=%q stderr=%q, want an error row, the moved branch kept, and the discarded ref at %s", applied.exit, applied.stdout, plan.stderr+applied.stderr, tip)
		}
	})
}

// TestDiscardTargetRefusesAConflictingDiscardedRef is RI38: a ref already at the planned path
// at another commit makes the apply refuse and keep both refs.
func TestDiscardTargetRefusesAConflictingDiscardedRef(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	ref, tip := uniqueBranch(t, root, "a")
	discarded := intent.DiscardedRef(discardDay, ref)
	planted := gitOutput(t, root, "rev-parse", "main")
	gitRun(t, root, "update-ref", discarded, planted)
	args := []string{"--discard-branch", "--target", branchSegment(ref)}
	plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
	applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(append(args, "--apply", plan.mustFingerprint(t))...)))
	rows := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")
	if plan.stderr != "" || applied.stderr != "" || applied.exit != 1 || rowAction(rows, ref) != string(ActionError) || refTip(root, ref) != tip || refTip(root, discarded) != planted {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q, want a refusal that keeps the branch and the planted ref", applied.exit, applied.stdout, plan.stderr+applied.stderr)
	}
}

// TestDiscardTargetStaleClassWritesNothing is RI64 and RI65: a class change between the plan
// and the apply refuses the old fingerprint before any discarded ref is written. The second
// case changes the class alone, because without the flag no row plans a discarded ref.
func TestDiscardTargetStaleClassWritesNothing(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		flags  []string
		mutate func(t *testing.T, root, home, ref string)
	}{
		{"RI64 a recorded active branch at a descendant", []string{"--discard-branch"}, func(t *testing.T, root, home, ref string) {
			active := mustCreate(t, root, home, "discard-stale-active", "active")
			gitRun(t, active.Path, "reset", "-q", "--hard", gitOutput(t, root, "rev-parse", ref))
			commitInWorktree(t, active.Path, "descendant.txt", "descendant\n", "descendant")
		}},
		{"RI64 main fast-forwarded onto the ref without the flag", nil, func(t *testing.T, root, _, ref string) {
			gitRun(t, root, "merge", "-q", "--ff-only", ref)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t)
			ref, tip := uniqueBranch(t, root, "a")
			args := append(tc.flags, "--target", branchSegment(ref))
			plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
			tc.mutate(t, root, home, ref)
			applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(append(args, "--apply", plan.mustFingerprint(t))...)))
			if plan.stderr != "" || applied.stderr != "" || applied.exit != 1 || !strings.Contains(applied.stdout, errStaleFingerprint.Error()) || refTip(root, ref) != tip || discardedRefs(t, root) != "" {
				t.Fatalf("apply exit=%d stdout=%q stderr=%q discarded=%q, want a stale refusal that writes nothing", applied.exit, applied.stdout, plan.stderr+applied.stderr, discardedRefs(t, root))
			}
		})
	}
}

// TestDiscardTargetLandedAndSubsumedWriteNoRef is RI40 and RI68: a landed and a subsumed
// target each apply with recovery none and write no discarded ref.
func TestDiscardTargetLandedAndSubsumedWriteNoRef(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		plant func(t *testing.T, root, home string) string
	}{
		{"RI40 a landed target", func(t *testing.T, root, _ string) string { return addUnclaimedBranch(t, root, "a") }},
		{"RI68 a subsumed target", func(t *testing.T, root, home string) string {
			ref, _ := recordedBranch(t, root, home, intent.StateActive, false)
			return ref
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t)
			ref := tc.plant(t, root, home)
			args := []string{"--discard-branch", "--target", branchSegment(ref)}
			plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(args...)))
			applied := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(append(args, "--apply", plan.mustFingerprint(t))...)))
			row := textRowsBy(t, applied.mustRows(t, cleanupTable), "target")[ref]
			if plan.stderr != "" || applied.stderr != "" || applied.exit != 0 || len(row) != len(cleanupFields) || row["action"] != string(ActionRemoved) || row["recovery"] != "none" || refTip(root, ref) != "" || discardedRefs(t, root) != "" {
				t.Fatalf("apply exit=%d stdout=%q stderr=%q, want the branch removed with recovery none and no discarded ref", applied.exit, applied.stdout, plan.stderr+applied.stderr)
			}
		})
	}
}

// TestDiscardTargetRoutesAUniqueShiftRow is RI78: the unclaimed plan ends a unique shift row
// with the branch path form of the discard command, and that command plans one row for it.
func TestDiscardTargetRoutesAUniqueShiftRow(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t)
	ref, _ := uniqueShiftBranch(t, root)
	plan := runVerb(t, verbClean, discardCall(repoHome{root, home}.call("--discard-branch", "--unclaimed")))
	route := "bench worktree clean --discard-branch --target " + strings.TrimPrefix(ref, "refs/heads/")
	if row := textRowsBy(t, plan.mustRows(t, cleanupTable), "target")[ref]; plan.stderr != "" || len(row) != len(cleanupFields) || !strings.HasSuffix(row["detail"], "; "+route) {
		t.Fatalf("unclaimed plan = %q stderr=%q, want the unique shift row to end with %q", plan.stdout, plan.stderr, route)
	}
	routed := runVerb(t, verbClean, discardCall(repoHome{root, home}.call(strings.Fields(strings.TrimPrefix(route, "bench worktree clean "))...)))
	routedRows := textRowsBy(t, routed.mustRows(t, cleanupTable), "target")
	if routed.exit != 0 || routed.stderr != "" || len(routedRows) != 1 || routedRows[ref] == nil {
		t.Fatalf("routed plan exit=%d stdout=%q stderr=%q, want one row for %s", routed.exit, routed.stdout, routed.stderr, ref)
	}
}
