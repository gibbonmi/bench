// When a cleanup-set apply refuses, and when it does not. The fixtures here drive drift at
// each window an apply passes through — before entry, before the first transaction, and
// inside a later member's own lock — and check which members survive. Two sibling files own
// the rest: clean_set_outcomes_test.go owns what a stopped apply reports, and
// clean_set_wiring_test.go owns the command's wiring to the refusal renderer. The fixture
// builders all three share live here.
package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// removableSetFixture creates count landed, clean assignments. Every member of that
// selection plans a removal, so a fault or a drift on one member is always a fault or a
// drift the set reaches with earlier members already removable.
func removableSetFixture(t *testing.T, count int) (string, string, []Creation) {
	t.Helper()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creations := make([]Creation, 0, count)
	for i := 0; i < count; i++ {
		creation := mustCreate(t, root, home, fmt.Sprintf("set-apply-%d", i), fmt.Sprintf("set apply %d", i))
		landAssignment(t, root, creation, fmt.Sprintf("member-%d.txt", i))
		creations = append(creations, creation)
	}
	return root, home, creations
}

// explicitIdentities names every member of creations by its canonical identity.
func explicitIdentities(creations []Creation) []string {
	identities := make([]string, 0, len(creations))
	for _, creation := range creations {
		identities = append(identities, creation.Assignment.ID)
	}
	return identities
}

// explicitTargetArguments is the explicit selection mode's operands for those identities.
func explicitTargetArguments(creations []Creation) []string {
	arguments := make([]string, 0, 2*len(creations))
	for _, identity := range explicitIdentities(creations) {
		arguments = append(arguments, "--target", identity)
	}
	return arguments
}

// landedFileByAssignment maps each member's identity to the landed file in its checkout.
// A rewrite of that file is the tracked drift every fixture below uses.
func landedFileByAssignment(creations []Creation, names ...string) map[string]string {
	files := make(map[string]string, len(creations))
	for i, creation := range creations {
		files[creation.Assignment.ID] = filepath.Join(creation.Path, names[i])
	}
	return files
}

// driftTracked rewrites one member's landed file, which leaves that member's tracked state
// dirty and takes its planned removal away.
func driftTracked(t *testing.T, files map[string]string, assignment string) {
	t.Helper()
	path, known := files[assignment]
	if !known {
		t.Fatalf("assignment %q has no landed file in %#v", assignment, files)
	}
	mustWrite(t, path, []byte("drifted\n"), 0o644)
}

// TestCleanSetPreexistingDrift is CL4. A member that is already dirty when the apply starts
// refuses the whole carried plan, so no earlier member is deleted for a set the repository
// no longer describes. Both selection modes answer the same way.
func TestCleanSetPreexistingDrift(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name      string
		selectors func([]Creation) []string
	}{
		{name: "explicit", selectors: explicitTargetArguments},
		{name: "landed", selectors: func([]Creation) []string { return []string{"--landed"} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home, creations := removableSetFixture(t, 2)
			files := landedFileByAssignment(creations, "member-0.txt", "member-1.txt")
			selectors := tc.selectors(creations)
			plan, planErr, planCode := runCleanup(t, root, home, selectors...)
			if planCode != 0 || planErr != "" {
				t.Fatalf("plan = (%d, %q, %q), want one applicable plan", planCode, plan, planErr)
			}
			later := creations[0]
			if creations[0].Assignment.ID < creations[1].Assignment.ID {
				later = creations[1]
			}
			driftTracked(t, files, later.Assignment.ID)
			before := repositoryState(t, root)

			stdout, stderr, code := runCleanup(t, root, home, append(selectors, "--apply", cleanupRowFingerprint(t, plan))...)
			if code != 1 || stderr != "" || !strings.Contains(stdout, errStaleFingerprint.Error()) {
				t.Fatalf("drifted apply = (%d, %q, %q), want a stale refusal", code, stdout, stderr)
			}
			if strings.Contains(stdout, ",removed,") {
				t.Fatalf("drifted apply = %q, want no completed removal", stdout)
			}
			requireMembersPresent(t, creations)
			if after := repositoryState(t, root); after != before {
				t.Fatalf("drifted apply changed the repository: %q -> %q", before, after)
			}
		})
	}
}

// TestCleanSetPreflightAllRows is CL5. It hands an approved set straight to the applier
// after the member the applier reaches second has drifted. The command's own entry re-plan
// cannot answer that window, so only a preflight over every removable row before the first
// transaction keeps the earlier member on disk.
func TestCleanSetPreflightAllRows(t *testing.T) {
	t.Parallel()
	t.Run("explicit", func(t *testing.T) {
		t.Parallel()
		root, _, creations := removableSetFixture(t, 2)
		files := landedFileByAssignment(creations, "member-0.txt", "member-1.txt")
		j := defaultJoins()
		set := planExplicitSet(j, root, explicitIdentities(creations), CleanupOptions{})
		if set.fingerprint == "" || len(set.rows) != 2 {
			t.Fatalf("explicit set = %#v, want two applicable members", set)
		}
		driftTracked(t, files, set.rows[1].assignment.ID)

		plans, err := applyExplicitSet(j, root, set, CleanupOptions{})
		if !errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want a stale refusal before the first transaction", err)
		}
		requireMembersPresent(t, creations)
		requireAllNotAttempted(t, plans)
	})
	t.Run("landed", func(t *testing.T) {
		t.Parallel()
		root, _, creations := removableSetFixture(t, 2)
		files := landedFileByAssignment(creations, "member-0.txt", "member-1.txt")
		j := defaultJoins()
		set, planErr := planLandedSet(j, root, CleanupOptions{}, "")
		if planErr != nil || len(set.rows) != 2 {
			t.Fatalf("landed set = %#v, %v; want two applicable members", set, planErr)
		}
		driftTracked(t, files, set.rows[1].assignment.ID)

		plans, err := applyLandedSet(j, root, set, CleanupOptions{}, "")
		if !errors.Is(err, errStaleFingerprint) {
			t.Fatalf("apply error = %v, want a stale refusal before the first transaction", err)
		}
		requireMembersPresent(t, creations)
		requireAllNotAttempted(t, plans)
	})
}

// requireAllNotAttempted holds a refused set to one report: the whole selection, every row
// of it unstarted. A refusal that reports fewer rows hides part of the selected intent.
func requireAllNotAttempted(t *testing.T, plans []CleanupPlan) {
	t.Helper()
	if len(plans) == 0 {
		t.Fatal("refused apply reported no outcome for the selected set")
	}
	for _, plan := range plans {
		if plan.Action != ActionNotAttempted {
			t.Fatalf("refused apply row = %q/%q, want every selected row not attempted", plan.Target, plan.Action)
		}
	}
}

// TestCleanSetLateDrift is CL6. The drift here lands inside the second member's own locked
// transaction, past the set preflight and past that member's pre-transaction re-plan. Only
// the recheck the transaction runs under its own lock can still refuse it, and the earlier
// member's removal stays completed.
func TestCleanSetLateDrift(t *testing.T) {
	t.Parallel()
	root, _, creations := removableSetFixture(t, 2)
	files := landedFileByAssignment(creations, "member-0.txt", "member-1.txt")
	j := defaultJoins()
	set := planExplicitSet(j, root, explicitIdentities(creations), CleanupOptions{})
	if set.fingerprint == "" || len(set.rows) != 2 {
		t.Fatalf("explicit set = %#v, want two applicable members", set)
	}
	settled, drifted := memberByID(t, creations, set.rows[0].assignment.ID), memberByID(t, creations, set.rows[1].assignment.ID)
	mutated, locks := false, 0
	j.cleanupBoundary = func(step LifecycleStep) error {
		if step == StepApplyLocked {
			locks++
		}
		if locks == 2 && !mutated {
			mutated = true
			driftTracked(t, files, drifted.Assignment.ID)
		}
		return nil
	}

	plans, err := applyExplicitSet(j, root, set, CleanupOptions{})
	if !errors.Is(err, errStaleFingerprint) || !mutated {
		t.Fatalf("apply = (%v, mutated=%t), want the under-lock recheck to refuse the raced member", err, mutated)
	}
	if _, statErr := os.Lstat(settled.Path); !os.IsNotExist(statErr) {
		t.Fatalf("the completed member %s was not settled: %v", settled.Path, statErr)
	}
	if _, statErr := os.Stat(drifted.Path); statErr != nil {
		t.Fatalf("the raced member %s was removed: %v", drifted.Path, statErr)
	}
	if len(plans) != 2 || plans[0].Action != ActionRemoved || plans[1].Action != ActionError {
		t.Fatalf("late-drift rows = %#v, want one completed and one failed outcome", plans)
	}
	if plans[1].Target != drifted.Path {
		t.Fatalf("failed row target = %q, want the raced member %q", plans[1].Target, drifted.Path)
	}
}

// memberByID resolves one planned row back to the fixture checkout it names.
func memberByID(t *testing.T, creations []Creation, assignment string) Creation {
	t.Helper()
	for _, creation := range creations {
		if creation.Assignment.ID == assignment {
			return creation
		}
	}
	t.Fatalf("assignment %q names no fixture member", assignment)
	return Creation{}
}

// retainedMemberFixture creates one landed clean member the plan marks removable and one
// landed member whose ignored residue the plan retains without `--discard-ignored`. It is
// the selection shape that proves what an apply does to a member it will not touch, and it
// leaves one member resolvable after the removable one is gone.
func retainedMemberFixture(t *testing.T) (string, string, Creation, Creation) {
	t.Helper()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored-*.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore retained residue")
	removable := mustCreate(t, root, home, "set-retained-removable", "removable member")
	retained := mustCreate(t, root, home, "set-retained-residue", "retained member")
	landAssignment(t, root, removable, "removable.txt")
	landAssignment(t, root, retained, "retained.txt")
	mustWrite(t, filepath.Join(retained.Path, "ignored-one.txt"), []byte("residue\n"), 0o644)
	return root, home, removable, retained
}

// TestCleanSetSpentPlan is CL12. A plan whose removals already completed cannot be applied
// a second time, so a retry can neither repeat a completed side effect nor report one. The
// retained member outlives the apply, which leaves a selection that still resolves under the
// spent digest and therefore a refusal with a live source.
func TestCleanSetSpentPlan(t *testing.T) {
	t.Parallel()
	root, home, removable, retained := retainedMemberFixture(t)
	both := []string{"--target", removable.Assignment.ID, "--target", retained.Assignment.ID}
	plan, planErr, planCode := runCleanup(t, root, home, both...)
	if planCode != 0 || planErr != "" {
		t.Fatalf("plan = (%d, %q, %q), want one applicable plan", planCode, plan, planErr)
	}
	digest := cleanupRowFingerprint(t, plan)
	applied, applyErr, applyCode := runCleanup(t, root, home, append(both, "--apply", digest)...)
	if applyCode != 0 || applyErr != "" || strings.Count(applied, ",removed,") != 1 {
		t.Fatalf("apply = (%d, %q, %q), want exactly one removal", applyCode, applied, applyErr)
	}
	spent := repositoryState(t, root)

	replay, replayErr, replayCode := runCleanup(t, root, home, append(both, "--apply", digest)...)
	if replayCode != 1 || replayErr != "" {
		t.Fatalf("replay = (%d, %q, %q), want a refusal", replayCode, replay, replayErr)
	}
	if strings.Contains(replay, ",removed,") {
		t.Fatalf("replay = %q, want no repeated completion claim", replay)
	}

	// The surviving member still resolves, so this replay reaches the digest comparison with
	// an applicable plan of its own. The spent digest names a set of two, and this selection
	// is a set of one, so the refusal here rests on the comparison rather than on a target
	// that has ceased to exist.
	narrowed := []string{"--target", retained.Assignment.ID, "--apply", digest}
	stale, staleErr, staleCode := runCleanup(t, root, home, narrowed...)
	if staleCode != 1 || staleErr != "" || !strings.Contains(stale, errStaleFingerprint.Error()) {
		t.Fatalf("narrowed replay = (%d, %q, %q), want a stale refusal", staleCode, stale, staleErr)
	}
	if _, statErr := os.Stat(retained.Path); statErr != nil {
		t.Fatalf("narrowed replay removed the retained member %s: %v", retained.Path, statErr)
	}
	if after := repositoryState(t, root); after != spent {
		t.Fatalf("replays changed the repository: %q -> %q", spent, after)
	}
}

// repositoryState is the durable state a refused or replayed apply must not change: the
// assignment ledger and the destination tip.
func repositoryState(t *testing.T, root string) string {
	t.Helper()
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	ids := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		ids = append(ids, assignment.ID)
	}
	sort.Strings(ids)
	head, err := git.Output("-C", root, "rev-parse", "HEAD")
	mustNoError(t, err)
	return strings.Join(ids, ",") + "@" + head
}

// requireMembersPresent fails when any selected checkout is gone. A refusal that removes
// nothing is the whole claim of the preflight rows, so the assertion names every member.
func requireMembersPresent(t *testing.T, creations []Creation) {
	t.Helper()
	for _, creation := range creations {
		if _, err := os.Stat(creation.Path); err != nil {
			t.Fatalf("refused apply removed %s: %v", creation.Path, err)
		}
	}
}
