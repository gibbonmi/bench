package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/usage"
)

func landedSetFixture(t *testing.T) landedSet {
	t.Helper()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	first := mustCreate(t, root, home, "landed-set-first", "first label")
	second := mustCreate(t, root, home, "landed-set-second", "second label")
	dirty := mustCreate(t, root, home, "landed-set-dirty", "dirty label")
	landAssignment(t, root, first, "first.txt")
	landAssignment(t, root, second, "second.txt")
	landAssignment(t, root, dirty, "dirty.txt")
	mustWrite(t, filepath.Join(dirty.Path, "dirty.txt"), []byte("changed\n"), 0o644)
	return landedSet{repoHome: repoHome{root, home}, first: first, second: second, dirty: dirty}
}

func TestCleanLandedPlansRepositoryWideSet(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 {
		t.Fatalf("CleanCommand exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	wantByID := map[string]string{
		f.first.Assignment.ID: f.first.Path, f.second.Assignment.ID: f.second.Path, f.dirty.Assignment.ID: f.dirty.Path,
	}
	ids := []string{f.first.Assignment.ID, f.second.Assignment.ID, f.dirty.Assignment.ID}
	sort.Strings(ids)
	positions := make([]int, len(ids))
	for i, id := range ids {
		positions[i] = strings.Index(plan.stdout, wantByID[id])
		if positions[i] < 0 {
			t.Fatalf("output=%q, want assignment %s", plan.stdout, id)
		}
		if i > 0 && positions[i] <= positions[i-1] {
			t.Fatalf("output=%q, want assignment-id order %v", plan.stdout, ids)
		}
	}
	if strings.Count(plan.stdout, ",remove,") != 2 || !strings.Contains(plan.stdout, ",retain,dirty,") || strings.Contains(plan.stdout, "refs/bench/recovery/") {
		t.Fatalf("output=%q, want two removes, one dirty retain, and no recovery ref", plan.stdout)
	}
}

func TestCleanLandedPlanSharesOneFingerprint(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 {
		t.Fatalf("CleanCommand exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	if rows := plan.mustRows(t, cleanupTable); len(rows) != 3 {
		t.Fatalf("output=%q, want one shared row fingerprint", plan.stdout)
	}
	plan.mustFingerprint(t)
}

func TestCleanLandedFingerprintBindsSetMembership(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	first := mustCreate(t, f.root, f.home, "fingerprint-member-first", "first")
	landAssignment(t, f.root, first, "first.txt")
	before := runVerb(t, verbClean, f.call("--landed"))
	if before.exit != 0 || before.stderr != "" {
		t.Fatalf("first plan exit=%d stdout=%q stderr=%q", before.exit, before.stdout, before.stderr)
	}

	second := mustCreate(t, f.root, f.home, "fingerprint-member-second", "second")
	landAssignment(t, f.root, second, "second.txt")
	after := runVerb(t, verbClean, f.call("--landed"))
	if after.exit != 0 || after.stderr != "" {
		t.Fatalf("second plan exit=%d stdout=%q stderr=%q", after.exit, after.stdout, after.stderr)
	}
	if before.mustFingerprint(t) == after.mustFingerprint(t) {
		t.Fatalf("set fingerprint did not change when membership changed: before=%q after=%q", before.stdout, after.stdout)
	}
}

func TestCleanLandedPlanAdvertisesApplyAndRemedies(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	if plan.exit != 0 {
		t.Fatalf("CleanCommand exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	fingerprint := plan.mustFingerprint(t)
	if !strings.Contains(plan.stdout, "bench worktree clean --landed --apply "+fingerprint) || !strings.Contains(plan.stdout, "bench worktree clean "+f.dirty.Path) {
		t.Fatalf("output=%q, want apply and dirty-row remedy", plan.stdout)
	}
	for _, creation := range []Creation{f.first, f.second, f.dirty} {
		if _, err := os.Stat(creation.Path); err != nil {
			t.Fatalf("bare plan removed %s: %v", creation.Path, err)
		}
	}
}

func TestCleanLandedEmptySetExitsClean(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	for attempt := 0; attempt < 2; attempt++ {
		plan := runVerb(t, verbClean, f.call("--landed"))
		if plan.exit != 0 || plan.stdout != cleanupTable+"[0]{target,action,tracked,ignored,recovery,fingerprint,detail}:\n" || plan.stderr != "" {
			t.Fatalf("attempt %d exit=%d stdout=%q stderr=%q", attempt, plan.exit, plan.stdout, plan.stderr)
		}
	}
}

func TestCleanLandedApplyOnEmptySetRefused(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	result := runVerb(t, verbClean, repoHome{root, filepath.Join(root, ".bench-home")}.call("--landed", "--apply", strings.Repeat("a", 64)))
	if result.exit != 2 || !strings.Contains(result.stdout, "invalid invocation; run "+usage.WorktreeClean) || result.stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", result.exit, result.stdout, result.stderr)
	}
}

func TestCleanLandedRefusesPathOperand(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	wantUsage := usage.WorktreeClean
	for _, args := range [][]string{{"--landed", root}, {root, "--landed"}} {
		result := runVerb(t, verbClean, f.call(args...))
		if result.exit != 2 || !strings.Contains(result.stdout, wantUsage) || result.stderr != "" {
			t.Fatalf("args=%q exit=%d stdout=%q stderr=%q", args, result.exit, result.stdout, result.stderr)
		}
	}
}

func TestCleanLandedRefusesMalformedFingerprint(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	f := repoHome{root, filepath.Join(root, ".bench-home")}
	for _, fingerprint := range []string{strings.Repeat("a", 63), strings.Repeat("a", 65), strings.Repeat("g", 64), strings.Repeat("A", 64)} {
		result := runVerb(t, verbClean, f.call("--landed", "--apply", fingerprint))
		if result.exit != 2 || !strings.Contains(result.stdout, usage.WorktreeClean) || result.stderr != "" {
			t.Fatalf("fingerprint=%q exit=%d stdout=%q stderr=%q", fingerprint, result.exit, result.stdout, result.stderr)
		}
	}
}

func TestCleanLandedSelectorPartition(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	unknown := mustCreate(t, root, home, "selector-unknown", "unknown lease")
	live := mustCreate(t, root, home, "selector-live", "live lease")
	unlanded := mustCreate(t, root, home, "selector-unlanded", "unlanded")
	pending := mustCreate(t, root, home, "selector-pending", "pending")
	recovered := mustCreate(t, root, home, "selector-recovered", "recovered")
	complete := mustCreate(t, root, home, "selector-complete", "complete")
	for i, creation := range []Creation{unknown, live, pending, recovered, complete} {
		landAssignment(t, root, creation, fmt.Sprintf("landed-%d.txt", i))
	}
	commitInWorktree(t, unlanded.Path, "unlanded.txt", "unlanded\n", "unlanded")
	unknownLease, err := LeaseFile(unknown.Path)
	mustNoError(t, err)
	mustWrite(t, unknownLease, []byte("not-a-lease\n"), 0o600)
	liveLease, err := LeaseFile(live.Path)
	mustNoError(t, err)
	mustWrite(t, liveLease, []byte(strconv.Itoa(os.Getpid())+" 2026-07-15T00:00:00Z\n"), 0o600)
	pending.Assignment.State = intent.StateCleanupPending
	mustNoError(t, intent.PutAssignment(root, pending.Assignment))
	recovered.Assignment.State = intent.StateRecovered
	recovered.Assignment.Recovery = []intent.Recovery{{Ref: intent.RecoveryRefPrefix(recovered.Assignment.OwnerID, recovered.Assignment.ID) + "1", Root: strings.Repeat("a", 40), Payloads: []string{strings.Repeat("b", 40)}}}
	mustNoError(t, intent.PutAssignment(root, recovered.Assignment))
	complete.Assignment.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(root, complete.Assignment))

	plan := runVerb(t, verbClean, repoHome{root, home}.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" || !strings.Contains(plan.stdout, unknown.Path+",retain,") || !strings.Contains(plan.stdout, "assignment lease state is unknown") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	for _, excluded := range []Creation{live, unlanded, pending, recovered, complete} {
		if strings.Contains(plan.stdout, excluded.Path) {
			t.Fatalf("output=%q, unexpectedly selected %s", plan.stdout, excluded.Path)
		}
	}
}

// TestCleanLandedPlanApplyCarriesModifiers holds the advertised apply command to the
// modifiers the plan answered under. The re-plan a refusal renders reads them from the same
// source; without this case only the modifier-free spelling of that head is pinned.
func TestCleanLandedPlanApplyCarriesModifiers(t *testing.T) {
	t.Parallel()
	root, home := ignoringRepo(t)
	creation := landedMember(t, root, home, "landed-apply-modifiers", "landed.txt")
	mustWrite(t, filepath.Join(creation.Path, "ignored-one.txt"), []byte("residue\n"), 0o644)

	plan := runVerb(t, verbClean, repoHome{root, home}.call("--discard-ignored", "--full", "--landed"))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
	}
	want := "bench worktree clean --discard-ignored --full --landed --apply " + plan.mustFingerprint(t)
	if !strings.Contains(plan.stdout, want) {
		t.Fatalf("plan = %q, want the advertised apply command %q", plan.stdout, want)
	}
}
