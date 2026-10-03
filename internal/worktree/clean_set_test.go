package worktree

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestCleanExplicitSetPlan(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--target", f.first.Assignment.ID, "--target", f.second.Assignment.ID))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("set plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	if rows := plan.mustRows(t, cleanupTable); len(rows) != 2 {
		t.Fatalf("set plan rows = %#v, want one row per selected target", rows)
	}
	shared := plan.mustFingerprint(t)
	if !strings.Contains(plan.stdout, "bench worktree clean --target ") || !strings.Contains(plan.stdout, "--apply "+shared) {
		t.Fatalf("set plan = %q, want the whole-set apply action", plan.stdout)
	}
	for _, creation := range []Creation{f.first, f.second} {
		if _, err := os.Stat(creation.Path); err != nil {
			t.Fatalf("bare set plan removed %s: %v", creation.Path, err)
		}
	}

	narrow := runVerb(t, verbClean, f.call("--target", f.first.Assignment.ID))
	if narrow.exit != 0 || narrow.stderr != "" {
		t.Fatalf("narrow plan exit=%d stdout=%q stderr=%q", narrow.exit, narrow.stdout, narrow.stderr)
	}
	if narrow.mustFingerprint(t) == shared {
		t.Fatalf("set fingerprint did not bind complete membership: pair=%q one=%q", plan.stdout, narrow.stdout)
	}

	applied := runVerb(t, verbClean, f.call("--target", f.first.Assignment.ID, "--target", f.second.Assignment.ID, "--apply", shared))
	if applied.exit != 0 || applied.stderr != "" || strings.Count(applied.stdout, ",removed,") != 2 {
		t.Fatalf("set apply = (%d, %q, %q), want two removals", applied.exit, applied.stdout, applied.stderr)
	}
	for _, creation := range []Creation{f.first, f.second} {
		if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) {
			t.Fatalf("set apply left %s: %v", creation.Path, err)
		}
	}
}

func TestCleanExplicitSetAliases(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	aliases := []string{
		"--target", f.first.Assignment.ID, "--target", f.first.Path, "--target", f.first.Assignment.Label,
		"--target", f.second.Assignment.Label, "--target", f.second.Assignment.ID,
	}
	plan := runVerb(t, verbClean, f.call(aliases...))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("alias plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	rows := textRows(t, plan.mustRows(t, cleanupTable))
	if len(rows) != 2 {
		t.Fatalf("alias plan rows = %#v, want one row per cleanup identity", rows)
	}
	lower, higher := f.first, f.second
	if higher.Assignment.ID < lower.Assignment.ID {
		lower, higher = f.second, f.first
	}
	if rows[0]["target"] != lower.Path || rows[1]["target"] != higher.Path {
		t.Fatalf("alias plan targets = %q, %q, want canonical identity order", rows[0]["target"], rows[1]["target"])
	}
	canonical := runVerb(t, verbClean, f.call("--target", f.first.Assignment.ID, "--target", f.second.Assignment.ID))
	if canonical.exit != 0 || plan.mustFingerprint(t) != canonical.mustFingerprint(t) {
		t.Fatalf("alias plan = %q and canonical plan = %q, want one identity-ordered fingerprint", plan.stdout, canonical.stdout)
	}

	applied := runVerb(t, verbClean, f.call(append(aliases, "--apply", plan.mustFingerprint(t))...))
	if applied.exit != 0 || applied.stderr != "" || strings.Count(applied.stdout, ",removed,") != 2 || len(applied.mustRows(t, cleanupTable)) != 2 {
		t.Fatalf("alias apply = (%d, %q, %q), want one removal per identity", applied.exit, applied.stdout, applied.stderr)
	}
	for _, creation := range []Creation{f.first, f.second} {
		if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) {
			t.Fatalf("alias apply left %s: %v", creation.Path, err)
		}
	}
}

func TestCleanExplicitSetSelectionFailure(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		target func(*testing.T, string, string) string
		detail string
	}{
		{name: "unassigned", target: func(*testing.T, string, string) string { return "absent-selection-target" }, detail: "target is unassigned"},
		{name: "ambiguous", target: func(t *testing.T, root, home string) string {
			mustCreate(t, root, home, "ambiguous-alias-one", "shared alias")
			mustCreate(t, root, home, "ambiguous-alias-two", "shared alias")
			return "shared alias"
		}, detail: "target is ambiguous"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := landedSetFixture(t)
			failing := tc.target(t, f.root, f.home)
			plan := runVerb(t, verbClean, f.call("--target", f.first.Assignment.ID, "--target", failing))
			if plan.exit != 1 || plan.stderr != "" {
				t.Fatalf("failed selection exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
			}
			if rows := plan.mustRows(t, cleanupTable); len(rows) != 2 {
				t.Fatalf("failed selection rows = %#v, want every selection outcome", rows)
			}
			if !strings.Contains(plan.stdout, tc.detail) || !strings.Contains(plan.stdout, ",error,") {
				t.Fatalf("failed selection = %q, want the reported failure", plan.stdout)
			}
			plan.mustNoFingerprint(t)
			if strings.Contains(plan.stdout, "bench worktree clean --target") {
				t.Fatalf("failed selection = %q, want no applicable fingerprint and no apply action", plan.stdout)
			}
			applied := runVerb(t, verbClean, f.call("--target", f.first.Assignment.ID, "--target", failing, "--apply", strings.Repeat("a", 64)))
			if applied.exit != 1 || applied.stderr != "" {
				t.Fatalf("failed selection apply exit=%d stderr=%q, want a refusal", applied.exit, applied.stderr)
			}
			for _, creation := range []Creation{f.first, f.second} {
				if _, err := os.Stat(creation.Path); err != nil {
					t.Fatalf("failed selection removed %s: %v", creation.Path, err)
				}
			}
		})
	}
}

func TestCleanSetRetainsAuthority(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	active := mustCreate(t, f.root, f.home, "set-authority-active", "active label")
	commitInWorktree(t, active.Path, "active.txt", "active\n", "active work")
	lease, err := LeaseFile(active.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(strconv.Itoa(os.Getpid())+" 2026-07-15T00:00:00Z\n"), 0o600)
	members := []Creation{f.first, f.dirty, active}

	// Every explicit target keeps the verdict the single-target form reaches for it, so a
	// set selection cannot widen what the command is allowed to remove.
	verdicts := make([]string, 0, len(members))
	for _, member := range members {
		single := runVerb(t, verbClean, f.call(member.Path))
		if single.exit != 0 || single.stderr != "" {
			t.Fatalf("single-target plan for %s exit=%d stdout=%q stderr=%q", member.Path, single.exit, single.stdout, single.stderr)
		}
		set := runVerb(t, verbClean, f.call("--target", member.Assignment.ID))
		if set.exit != 0 || set.stderr != "" {
			t.Fatalf("set plan for %s exit=%d stdout=%q stderr=%q", member.Path, set.exit, set.stdout, set.stderr)
		}
		singleRow, setRow := textRows(t, single.mustRows(t, cleanupTable))[0], textRows(t, set.mustRows(t, cleanupTable))[0]
		singleRow["fingerprint"], setRow["fingerprint"] = "fingerprint", "fingerprint"
		if !maps.Equal(singleRow, setRow) {
			t.Fatalf("set row = %q, want the single-target verdict %q", setRow, singleRow)
		}
		verdicts = append(verdicts, setRow["action"])
	}
	if verdicts[2] != string(ActionRetain) {
		t.Fatalf("a live-leased target planned %q, want the existing refusal", verdicts[2])
	}

	args := []string{"--target", f.first.Assignment.ID, "--target", f.dirty.Assignment.ID, "--target", active.Assignment.ID}
	plan := runVerb(t, verbClean, f.call(args...))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("mixed set plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
	}
	sorted := append([]string(nil), verdicts...)
	sort.Strings(sorted)
	var planned []string
	for _, row := range textRows(t, plan.mustRows(t, cleanupTable)) {
		planned = append(planned, row["action"])
	}
	sort.Strings(planned)
	if strings.Join(planned, ",") != strings.Join(sorted, ",") {
		t.Fatalf("mixed set verdicts = %#v, want each target's own verdict %#v", planned, sorted)
	}
	applied := runVerb(t, verbClean, f.call(append(args, "--apply", plan.mustFingerprint(t))...))
	if applied.exit != 0 || applied.stderr != "" {
		t.Fatalf("mixed set apply = (%d, %q), want success", applied.exit, applied.stderr)
	}
	if _, statErr := os.Stat(active.Path); statErr != nil {
		t.Fatalf("set apply removed the live-leased target %s: %v", active.Path, statErr)
	}
}

// cleanupEffects is one clean form's normalized lifecycle result: the exit code, every
// rendered verdict, and the durable state each fixture checkout, branch, and assignment
// record is left in. The verdicts sort, because a set orders its rows by an assignment
// identity the fixture draws at random.
func cleanupEffects(t *testing.T, root string, creations []Creation, result verbResult) string {
	t.Helper()
	verdicts := make([]string, 0, len(creations))
	for _, row := range textRows(t, result.mustRows(t, cleanupTable)) {
		verdicts = append(verdicts, row["action"]+"/"+row["tracked"])
	}
	sort.Strings(verdicts)
	trees := make([]string, 0, len(creations))
	for _, creation := range creations {
		tree := "absent"
		if _, err := os.Stat(creation.Path); err == nil {
			tree = "present"
		}
		trees = append(trees, fmt.Sprintf("%s/%t", tree, git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch)))
	}
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	return fmt.Sprintf("exit=%d stderr=%t; rows=%s; trees=%s; assignments=%d",
		result.exit, result.stderr == "", strings.Join(verdicts, ","), strings.Join(trees, ","), len(assignments))
}

// TestCleanSetCompatibility is the differential that holds every existing cleanup form to
// the lifecycle effects it produces without explicit sets. Each baseline below is that
// pre-set effect, so a changed branch, receipt, or retention verdict on any existing form
// turns this row red.
func TestCleanSetCompatibility(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		run  func(*testing.T, landedSet) verbResult
		want string
	}{
		{
			name: "single target plan",
			run: func(t *testing.T, f landedSet) verbResult {
				return runVerb(t, verbClean, f.call(f.first.Path))
			},
			want: "exit=0 stderr=true; rows=remove/clean; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "single target apply",
			run: func(t *testing.T, f landedSet) verbResult {
				plan := runVerb(t, verbClean, f.call(f.first.Path))
				return runVerb(t, verbClean, f.call(f.first.Path, "--apply", plan.mustFingerprint(t)))
			},
			want: "exit=0 stderr=true; rows=removed/clean; trees=absent/false,present/true,present/true; assignments=2",
		},
		{
			name: "landed plan",
			run: func(t *testing.T, f landedSet) verbResult {
				return runVerb(t, verbClean, f.call("--landed"))
			},
			want: "exit=0 stderr=true; rows=remove/clean,remove/clean,retain/dirty; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "landed apply",
			run: func(t *testing.T, f landedSet) verbResult {
				plan := runVerb(t, verbClean, f.call("--landed"))
				return runVerb(t, verbClean, f.call("--landed", "--apply", plan.mustFingerprint(t)))
			},
			want: "exit=0 stderr=true; rows=removed/clean,removed/clean,retain/dirty; trees=absent/false,absent/false,present/true; assignments=1",
		},
		{
			name: "unclaimed plan",
			run: func(t *testing.T, f landedSet) verbResult {
				return runVerb(t, verbClean, f.call("--discard-branch", "--unclaimed"))
			},
			want: "exit=0 stderr=true; rows=discard-remove/unclaimed; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "unclaimed apply",
			run: func(t *testing.T, f landedSet) verbResult {
				plan := runVerb(t, verbClean, f.call("--discard-branch", "--unclaimed"))
				return runVerb(t, verbClean, f.call("--discard-branch", "--unclaimed", "--apply", plan.mustFingerprint(t)))
			},
			want: "exit=0 stderr=true; rows=removed/unclaimed; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "dirty single target apply",
			run: func(t *testing.T, f landedSet) verbResult {
				plan := runVerb(t, verbClean, f.call(f.dirty.Path))
				return runVerb(t, verbClean, f.call(f.dirty.Path, "--apply", plan.mustFingerprint(t)))
			},
			want: "exit=0 stderr=true; rows=removed/dirty; trees=present/true,present/true,absent/false; assignments=3",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := landedSetFixture(t)
			orphan := intent.AssignmentBranchRef(strings.Repeat("c", 32), strings.Repeat("d", 32))
			gitRun(t, f.root, "branch", strings.TrimPrefix(orphan, "refs/heads/"))
			creations := []Creation{f.first, f.second, f.dirty}
			if got := cleanupEffects(t, f.root, creations, tc.run(t, f)); got != tc.want {
				t.Fatalf("%s effects = %q, want the captured baseline %q", tc.name, got, tc.want)
			}
		})
	}
}

func TestCleanSetPresentEmptyInventory(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "present-empty-inventory", "present empty")
	mustNoError(t, intent.DeleteAssignment(root, creation.Assignment.ID))
	ledger, err := intent.Address(root)
	mustNoError(t, err)
	if _, statErr := os.Stat(ledger); statErr != nil {
		t.Fatalf("present-empty fixture has no ledger at %q: %v", ledger, statErr)
	}
	requireEmptyInventoryOutcome(t, repoHome{root, home}, creation)
}

func TestCleanSetAbsentInventory(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	ledger, err := intent.Address(root)
	mustNoError(t, err)
	if _, statErr := os.Lstat(ledger); !os.IsNotExist(statErr) {
		t.Fatalf("absent fixture already holds a ledger at %q: %v", ledger, statErr)
	}
	requireEmptyInventoryOutcome(t, repoHome{root, home}, Creation{Assignment: intent.Assignment{ID: strings.Repeat("e", 32)}})
}

// requireEmptyInventoryOutcome holds both empty landed inventories to the same answer: the
// selector reports an empty plan, an explicit target reports its own failed selection, and
// neither mutates the repository.
func requireEmptyInventoryOutcome(t *testing.T, f repoHome, absent Creation) {
	t.Helper()
	before, err := git.Output("-C", f.root, "rev-parse", "HEAD")
	mustNoError(t, err)
	landed := runVerb(t, verbClean, f.call("--landed"))
	if landed.exit != 0 || landed.stderr != "" || !strings.HasPrefix(landed.stdout, cleanupTable+"[0]") || len(landed.mustRows(t, cleanupTable)) != 0 {
		t.Fatalf("landed plan = (%d, %q, %q), want an empty plan", landed.exit, landed.stdout, landed.stderr)
	}
	targeted := runVerb(t, verbClean, f.call("--target", absent.Assignment.ID))
	if targeted.exit != 1 || targeted.stderr != "" || len(targeted.mustRows(t, cleanupTable)) != 1 {
		t.Fatalf("explicit plan = (%d, %q, %q), want one unapplicable selection outcome", targeted.exit, targeted.stdout, targeted.stderr)
	}
	targeted.mustNoFingerprint(t)
	assignments, err := intent.Assignments(f.root)
	mustNoError(t, err)
	after, err := git.Output("-C", f.root, "rev-parse", "HEAD")
	mustNoError(t, err)
	if len(assignments) != 0 || after != before {
		t.Fatalf("empty inventory mutated the repository: assignments=%#v head %q -> %q", assignments, before, after)
	}
}
