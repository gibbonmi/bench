// How the clean command wires each selection mode to the refusal renderer. Every route that
// can return a stale refusal from inside an apply has to reach the refusal form through the
// command, so a call site that reverted to the plain outcome rows turns red.
//
// The clean_set_*_test.go files split by the question each answers; this one answers how the
// command reaches the renderer, and how the unclaimed adapter it reaches visits its re-planned
// selection through the shared cleanup traversal.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
)

// TestCleanSetApplyTimeStaleWiring is COV-8. The renderer is proved elsewhere; this holds the
// command's wiring to it. Each case reaches the apply-time stale arm through the command and
// reads the refusal form back, so a wiring that reverted to the plain rows turns red.
func TestCleanSetApplyTimeStaleWiring(t *testing.T) {
	t.Parallel()
	t.Run("landed apply-time drift", func(t *testing.T) {
		t.Parallel()
		f := landedSetFixture(t)
		plan := runVerb(t, verbClean, f.call("--landed"))
		if plan.exit != 0 || plan.stderr != "" {
			t.Fatalf("plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
		}
		// The set removes in canonical identity order and the fixture's third member retains,
		// so the later of these two identities is the member the apply reaches second.
		drifted := f.first
		if f.second.Assignment.ID > f.first.Assignment.ID {
			drifted = f.second
		}
		// The drift lands after the first member's terminal receipt, so the command's entry
		// check passes and the refusal comes from inside the apply.
		j := defaultJoins()
		mutated := false
		j.cleanupBoundary = func(step LifecycleStep) error {
			if step == StepTerminalReceipt && !mutated {
				mutated = true
				mustWrite(t, filepath.Join(drifted.Path, "drifted.txt"), []byte("drifted\n"), 0o644)
			}
			return nil
		}

		digest := plan.mustFingerprint(t)
		stale := runVerb(t, verbClean, f.callWith(j, "--landed", "--apply", digest))
		if stale.exit != 1 || stale.stderr != "" || !mutated {
			t.Fatalf("landed apply = (%d, %q, %q) mutated=%t, want an apply-time refusal", stale.exit, stale.stdout, stale.stderr, mutated)
		}
		requireStaleRefusalRow(t, stale, "unknown", "unknown", digest)
		if !strings.Contains(stale.stdout, "bench 'worktree' 'clean' '--landed'") {
			t.Fatalf("landed apply = %q, want the selector-preserving re-plan command", stale.stdout)
		}
	})
	// Both unclaimed routes read the branch namespace twice, once for the caller's plan and
	// once inside the apply. A branch appearing between those reads is the race, and the seam
	// stands in that window.
	for _, tc := range []struct {
		name  string
		apply func(string) []string
	}{
		{name: "unclaimed apply", apply: func(d string) []string {
			return []string{"--discard-branch", "--unclaimed", "--apply", d}
		}},
		{name: "unclaimed apply-current", apply: func(string) []string {
			return []string{"--discard-branch", "--unclaimed", "--apply-current"}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home := unclaimedBranchFixture(t, "a", "b")
			set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
			mustNoError(t, err)
			if len(set.rows) != 2 {
				t.Fatalf("unclaimed set = %#v, want two selected refs", set.rows)
			}
			j, raced := defaultJoins(), false
			j.cleanupBoundary = func(step LifecycleStep) error {
				if step == StepUnlockedReplan && !raced {
					raced = true
					addUnclaimedBranch(t, root, "d")
				}
				return nil
			}

			stale := runVerb(t, verbClean, repoHome{root, home}.callWith(j, tc.apply(set.fingerprint)...))
			if stale.exit != 1 || stale.stderr != "" || !raced {
				t.Fatalf("%s = (%d, %q, %q) raced=%t, want an apply-time refusal", tc.name, stale.exit, stale.stdout, stale.stderr, raced)
			}
			requireStaleRefusalRow(t, stale, "unclaimed", "none", set.fingerprint)
			if !strings.Contains(stale.stdout, "bench 'worktree' 'clean' '--discard-branch' '--unclaimed'") {
				t.Fatalf("%s = %q, want the selector-preserving re-plan command", tc.name, stale.stdout)
			}
			for _, row := range set.rows {
				if !git.OK("-C", root, "show-ref", "--verify", "--quiet", row.ref) {
					t.Fatalf("%s deleted %q", tc.name, row.ref)
				}
			}
		})
	}
}

// unclaimedApply is an unclaimed selection planned over a unique branch 1, the landed branches
// a, b, and c, and a unique branch e, so the set reads retained, three removable, retained. A
// branch outside the namespace holds one commit no selected ref names.
type unclaimedApply struct {
	root, spare string
	set         unclaimedAssignmentSet
	tips        map[string]string
}

func unclaimedApplyFixture(t *testing.T) unclaimedApply {
	t.Helper()
	root, _ := unclaimedBranchFixture(t, "a", "b", "c")
	uniqueBranch(t, root, "1")
	uniqueBranch(t, root, "e")
	gitRun(t, root, "branch", "spare")
	f := unclaimedApply{root: root, spare: commitOnBranch(t, root, "refs/heads/spare", "spare.txt", "spare\n"), tips: map[string]string{}}
	set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
	mustNoError(t, err)
	f.set = set
	for i, row := range set.rows {
		f.tips[row.ref] = refTip(root, row.ref)
		if row.action().Removes() != (i >= 1 && i <= 3) || set.fingerprint == "" || len(set.rows) != 5 {
			t.Fatalf("unclaimed set = %#v, want retained, three removable, retained", set)
		}
	}
	return f
}

func (f unclaimedApply) apply(j joins) ([]CleanupPlan, error) {
	return applyUnclaimedAssignmentSet(j, f.root, f.set, unclaimedOptions())
}

// requireRefs fails unless each planned ref reads the tip changed names for it, where an empty
// tip is a deleted ref, and every other planned ref still reads its planned tip.
func (f unclaimedApply) requireRefs(t *testing.T, changed map[string]string) {
	t.Helper()
	for _, row := range f.set.rows {
		want, ok := changed[row.ref]
		if !ok {
			want = f.tips[row.ref]
		}
		if got := refTip(f.root, row.ref); got != want {
			t.Fatalf("ref %s reads %q, want %q", row.ref, got, want)
		}
	}
}

// requirePrefix fails unless plans opens with one row per planned row before n, in plan
// order: a retained row exactly as planned and a removable row removed.
func (f unclaimedApply) requirePrefix(t *testing.T, plans []CleanupPlan, n int) {
	t.Helper()
	for i, row := range f.set.rows[:n] {
		got := plans[i]
		if row.action().Removes() && (got.Target != row.ref || got.Action != ActionRemoved || got.Reason != "") ||
			!row.action().Removes() && !reflect.DeepEqual(got, row.plan(f.set.fingerprint)) {
			t.Fatalf("row %d = %#v, want the outcome of planned %s", i, got, row.ref)
		}
	}
}

// TestCleanupUnclaimedModeWitnesses is R01, R02, R04, R05, R27, and R28 for the unclaimed
// adapter. Its only lifecycle boundary is the one before its re-plan. An empty set changes
// nothing. A mixed set reports every row in plan order, passes each unique row through as
// planned and keeps its tip, removes each landed row, and leaves the checked-out ref alone.
func TestCleanupUnclaimedModeWitnesses(t *testing.T) {
	t.Parallel()
	t.Run("R01 empty", func(t *testing.T) {
		t.Parallel()
		root, _ := unclaimedBranchFixture(t)
		set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
		mustNoError(t, err)
		j := defaultJoins()
		steps, trace := stepTrace()
		j.cleanupBoundary = trace
		before := strings.Join(baselineRefs(t, root), "\n")

		plans, err := applyUnclaimedAssignmentSet(j, root, set, unclaimedOptions())
		mustNoError(t, err)
		if len(plans) != 0 || !reflect.DeepEqual(*steps, []LifecycleStep{StepUnlockedReplan}) || strings.Join(baselineRefs(t, root), "\n") != before {
			t.Fatalf("rows = %#v, boundaries = %v, want no row, the re-plan boundary alone, and unchanged refs", plans, *steps)
		}
	})
	t.Run("R02 R04 R05 R27 R28 mixed", func(t *testing.T) {
		t.Parallel()
		f := unclaimedApplyFixture(t)
		checkedOut := addUnclaimedBranchAt(t, f.root, "d", "main")
		gitRun(t, f.root, "worktree", "add", "-q", filepath.Join(t.TempDir(), "checked-out"), strings.TrimPrefix(checkedOut, "refs/heads/"))
		set, err := planUnclaimedAssignmentSet(f.root, unclaimedOptions())
		mustNoError(t, err)
		f.set = set
		j := defaultJoins()
		steps, trace := stepTrace()
		j.cleanupBoundary = trace

		plans, err := f.apply(j)
		mustNoError(t, err)
		if len(plans) != len(f.set.rows) || !reflect.DeepEqual(*steps, []LifecycleStep{StepUnlockedReplan}) {
			t.Fatalf("rows = %#v, boundaries = %v, want one row per planned row and the re-plan boundary alone", plans, *steps)
		}
		f.requirePrefix(t, plans, len(f.set.rows))
		f.requireRefs(t, map[string]string{f.set.rows[1].ref: "", f.set.rows[2].ref: "", f.set.rows[3].ref: ""})
		if refTip(f.root, checkedOut) == "" {
			t.Fatalf("the checked-out ref %s was deleted", checkedOut)
		}
	})
}

// TestCleanupUnclaimedSelectionDrift is R20. A ref that enters the namespace, leaves it, or
// moves after the plan refuses the whole set before any deletion, and the refusal carries no
// row.
func TestCleanupUnclaimedSelectionDrift(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		// drift changes the selection and returns the planned refs it changed.
		drift func(*testing.T, unclaimedApply) map[string]string
	}{
		{"an addition", func(t *testing.T, f unclaimedApply) map[string]string {
			addUnclaimedBranch(t, f.root, "f")
			return nil
		}},
		{"a removal", func(t *testing.T, f unclaimedApply) map[string]string {
			ref := f.set.rows[3].ref
			gitRun(t, f.root, "branch", "-D", strings.TrimPrefix(ref, "refs/heads/"))
			return map[string]string{ref: ""}
		}},
		{"a tuple change", func(t *testing.T, f unclaimedApply) map[string]string {
			ref := f.set.rows[2].ref
			return map[string]string{ref: commitOnBranch(t, f.root, ref, "moved.txt", "moved\n")}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := unclaimedApplyFixture(t)
			changed := tc.drift(t, f)

			plans, err := f.apply(defaultJoins())
			if !errors.Is(err, errStaleFingerprint) || plans != nil {
				t.Fatalf("apply = (%#v, %v), want the stale refusal with no row", plans, err)
			}
			f.requireRefs(t, changed)
		})
	}
}

// TestCleanupUnclaimedDeletionFailureAfterReplan is R14, R18, and R19. A Git ref lock on the
// first or a later removable ref leaves the re-plan unchanged and fails that ref's deletion.
// The apply returns the completed prefix and the failed row alone, keeps each earlier deletion,
// and leaves the failed ref and every later ref at its planned tip.
func TestCleanupUnclaimedDeletionFailureAfterReplan(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		failed int
	}{{"first", 1}, {"later", 2}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := unclaimedApplyFixture(t)
			lockRefOf(tc.failed)(t, f.root, f.set)
			current, err := planUnclaimedAssignmentSet(f.root, unclaimedOptions())
			mustNoError(t, err)
			if current.fingerprint != f.set.fingerprint || len(current.rows) != len(f.set.rows) {
				t.Fatalf("the lock changed the re-plan: %#v", current)
			}

			plans, err := f.apply(defaultJoins())
			if err == nil || errors.Is(err, errStaleFingerprint) {
				t.Fatalf("apply error = %v, want a deletion fault", err)
			}
			failed := f.set.rows[tc.failed]
			if len(plans) != tc.failed+1 || plans[tc.failed].Target != failed.ref || plans[tc.failed].Action != ActionError || plans[tc.failed].Reason != err.Error() {
				t.Fatalf("rows = %#v, want the completed prefix and the failed row of %s alone", plans, failed.ref)
			}
			f.requirePrefix(t, plans, tc.failed)
			deleted := map[string]string{}
			for _, row := range f.set.rows[1:tc.failed] {
				deleted[row.ref] = ""
			}
			f.requireRefs(t, deleted)
		})
	}
}

// TestCleanupUnclaimedMovedRefAfterReplan is R21. A real reference-transaction hook moves the
// second removable ref once the deletion of the first commits, after the re-plan passed. The
// exact deletion of the moved ref then fails, so it keeps its new tip, the first stays deleted,
// and no row follows the failed one.
func TestCleanupUnclaimedMovedRefAfterReplan(t *testing.T) {
	t.Parallel()
	f := unclaimedApplyFixture(t)
	first, later := f.set.rows[1].ref, f.set.rows[2].ref
	hooks := gitOutput(t, f.root, "rev-parse", "--git-path", "hooks")
	if !filepath.IsAbs(hooks) {
		hooks = filepath.Join(f.root, hooks)
	}
	ran := filepath.Join(t.TempDir(), "ran")
	// The marker both records the run and ignores the hook's own nested update.
	script := "#!/bin/sh\ninput=$(cat)\ntest \"$1\" = committed || exit 0\ntest -e '" + ran + "' && exit 0\n" +
		"printf '%s\\n' \"$input\" | while read -r old new ref; do\n" +
		"\tif test \"$ref\" = '" + first + "'; then : > '" + ran + "'; git update-ref '" + later + "' " + f.spare + "; fi\ndone\n"
	mustMkdirAll(t, hooks, 0o755)
	mustWrite(t, filepath.Join(hooks, "reference-transaction"), []byte(script), 0o755)

	plans, err := f.apply(defaultJoins())
	if err == nil || errors.Is(err, errStaleFingerprint) {
		t.Fatalf("apply error = %v, want a deletion fault", err)
	}
	if _, statErr := os.Stat(ran); statErr != nil {
		t.Fatalf("the reference-transaction hook did not run: %v", statErr)
	}
	if len(plans) != 3 || plans[2].Target != later || plans[2].Action != ActionError {
		t.Fatalf("rows = %#v, want the completed prefix and the failed row of %s alone", plans, later)
	}
	f.requirePrefix(t, plans, 2)
	f.requireRefs(t, map[string]string{first: "", later: f.spare})
}
