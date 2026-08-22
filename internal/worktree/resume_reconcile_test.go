package worktree

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// seedLifecycleDebris writes the state a repository carries out of the removed
// provisional lifecycle: one ref in each swept namespace, plus the gate's
// verdict ref. It returns the green ref's object name for the survival
// assertion.
func seedLifecycleDebris(t *testing.T, root string) string {
	t.Helper()
	gitRun(t, root, "update-ref", "refs/bench/specbuild/candidate/deadbeef", "HEAD")
	gitRun(t, root, "update-ref", "refs/bench/recovery/"+strings.Repeat("a", 32)+"/"+strings.Repeat("b", 32)+"/1", "HEAD")
	gitRun(t, root, "update-ref", "refs/bench/green/main", "HEAD")
	return gitOutput(t, root, "rev-parse", "refs/bench/green/main")
}

// seedLegacyAssignments splices records the removed lifecycle left in the ledger
// into an existing one. It adds a recovered row naming preserved work, and a
// row this build's decoder cannot read. They are written as raw JSON because
// every ledger writer refuses them; only an older binary could have produced
// them.
func seedLegacyAssignments(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, ".git", intent.Filename)
	body, err := os.ReadFile(path)
	mustNoError(t, err)
	var ledger map[string]any
	mustNoError(t, json.Unmarshal(body, &ledger))
	owner, id := strings.Repeat("a", 32), strings.Repeat("b", 32)
	legacy := []any{
		map[string]any{
			"schema": "bench-assignment/v1", "id": id, "owner_id": owner,
			"request": strings.Repeat("c", 64), "label": "spec build ticket", "start": strings.Repeat("d", 40),
			"branch": "refs/heads/bench/assign/" + owner + "/" + id, "worktree": filepath.Join(root, "gone"),
			"state": "recovered",
			"recovery": []any{map[string]any{
				"ref": "refs/bench/recovery/" + owner + "/" + id + "/1", "root": strings.Repeat("e", 40),
				"payloads": []any{strings.Repeat("f", 40)},
			}},
		},
		map[string]any{
			"schema": "bench-assignment/v1", "id": strings.Repeat("9", 32), "owner_id": owner,
			"request": strings.Repeat("8", 64), "label": "provisional run", "start": strings.Repeat("7", 40),
			"branch":   "refs/heads/bench/assign/" + owner + "/" + strings.Repeat("9", 32),
			"worktree": filepath.Join(root, "gone-too"), "state": "provisional",
			"ticket": "specs/removed/tickets/one.md",
		},
	}
	assignments, _ := ledger["assignments"].([]any)
	ledger["assignments"] = append(assignments, legacy...)
	encoded, err := json.Marshal(ledger)
	mustNoError(t, err)
	mustWrite(t, path, append(encoded, '\n'), 0o600)
}

func refsUnder(t *testing.T, root string, namespaces ...string) string {
	t.Helper()
	args := append([]string{"for-each-ref", "--format=%(refname) %(objectname)"}, namespaces...)
	return gitOutput(t, root, args...)
}

func runResume(t *testing.T, root string) string {
	t.Helper()
	chdir(t, root)
	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	requireTest(t, code == 0, "ResumeCleanCommand exit=%d\nstdout=%s\nstderr=%s", code, stdout.String(), stderr.String())
	return stdout.String()
}

// TestResumeReconcileSparesGreenVerdictRefs is the RM10 guard: the sweep's blast
// radius is the two lifecycle namespaces and nothing else. The gate's verdict
// store shares the refs/bench/ prefix, so an over-broad delete would destroy
// green evidence at every session start.
func TestResumeReconcileSparesGreenVerdictRefs(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	green := seedLifecycleDebris(t, root)
	diagnostic := "refs/bench/diagnostic/probe"
	gitRun(t, root, "update-ref", diagnostic, "HEAD")

	runResume(t, root)

	remaining := refsUnder(t, root, "refs/bench/specbuild/", "refs/bench/recovery/")
	requireTest(t, remaining == "", "lifecycle refs survived the reconcile: %q", remaining)
	requireTest(t, gitOutput(t, root, "rev-parse", "refs/bench/green/main") == green, "reconcile moved the gate's green verdict ref")
	requireTest(t, gitOutput(t, root, "rev-parse", diagnostic) == green, "reconcile deleted a diagnostic ref outside the lifecycle namespaces")
}

// TestResumeReconcileIsIdempotent covers RM3's re-run half. The first run is the
// settling one; it reports what it swept and so legitimately differs. The two
// after it are the compared pair: once the tree has settled, a reconcile writes
// nothing and reports exactly what its predecessor did.
func TestResumeReconcileIsIdempotent(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	seedLifecycleDebris(t, root)

	requireTest(t, strings.Contains(runResume(t, root), "swept refs 2"), "settling reconcile did not report the sweep")
	before := refsUnder(t, root, "refs/bench/")
	ledgerPath := filepath.Join(root, ".git", intent.Filename)
	ledgerBefore, _ := os.ReadFile(ledgerPath)

	second := runResume(t, root)
	third := runResume(t, root)
	ledgerAfter, _ := os.ReadFile(ledgerPath)
	requireTest(t, second == third, "settled reconciles disagree:\nsecond=%q\nthird=%q", second, third)
	requireTest(t, refsUnder(t, root, "refs/bench/") == before, "a settled reconcile churned refs")
	requireTest(t, bytes.Equal(ledgerBefore, ledgerAfter), "a settled reconcile rewrote the ledger")
}

// TestResumeReconcilePurgesLegacyAssignments is RM5 and RM3's ledger half. A
// reconcile over records the removed lifecycle wrote, including one no current
// decoder can read, exits zero and drops them. It leaves the pool's own record
// alone and authors nothing.
func TestResumeReconcilePurgesLegacyAssignments(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	pool := mustCreate(t, root, "reconcile-pool", "pool work")
	green := seedLifecycleDebris(t, root)
	seedLegacyAssignments(t, root)

	first := runResume(t, root)

	surviving, err := intent.Assignments(root)
	requireTest(t, err == nil, "ledger unreadable after purge: %v", err)
	requireTest(t, len(surviving) == 1 && surviving[0].ID == pool.Assignment.ID, "surviving assignments = %#v", surviving)
	_, statErr := os.Stat(surviving[0].Worktree)
	requireTest(t, statErr == nil, "surviving assignment names a missing worktree: %v", statErr)
	requireTest(t, refsUnder(t, root, "refs/bench/specbuild/", "refs/bench/recovery/") == "", "reconcile left lifecycle refs")
	requireTest(t, gitOutput(t, root, "rev-parse", "refs/bench/green/main") == green, "reconcile moved the gate's green verdict ref")

	requireTest(t, strings.Contains(first, "swept refs 2") && strings.Contains(first, "reconciled 2"),
		"settling reconcile did not report what it removed: %q", first)
	ledgerPath := filepath.Join(root, ".git", intent.Filename)
	ledgerBefore, _ := os.ReadFile(ledgerPath)
	second := runResume(t, root)
	third := runResume(t, root)
	ledgerAfter, _ := os.ReadFile(ledgerPath)
	requireTest(t, second == third, "settled reconciles disagree:\nsecond=%q\nthird=%q", second, third)
	requireTest(t, bytes.Equal(ledgerBefore, ledgerAfter), "a settled reconcile rewrote the ledger")
}

// TestResumeReconcileConvergesAfterInterruption is RM9: a reconcile killed between two ref
// deletions leaves a partial namespace and an unpurged ledger, and re-entry finishes both.
func TestResumeReconcileConvergesAfterInterruption(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	root := newWorktreeRepo(t)
	pool := mustCreate(t, root, "interrupted-pool", "pool work")
	green := seedLifecycleDebris(t, root)
	seedLegacyAssignments(t, root)
	for _, name := range []string{"refs/bench/specbuild/checkpoint/one", "refs/bench/specbuild/checkpoint/two"} {
		gitRun(t, root, "update-ref", name, "HEAD")
	}
	interrupted := errors.New("killed mid-sweep")
	restore, sweeps := cleanupTransactionBoundary, 0
	cleanupTransactionBoundary = func(step LifecycleStep) error {
		if step != StepLifecycleSweep {
			return nil
		}
		if sweeps++; sweeps > 1 {
			return interrupted
		}
		return nil
	}
	chdir(t, root)
	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	cleanupTransactionBoundary = restore
	requireTest(t, code != 0, "interrupted reconcile exit=%d stdout=%s", code, stdout.String())
	partial := refsUnder(t, root, "refs/bench/specbuild/", "refs/bench/recovery/")
	requireTest(t, partial != "" && len(strings.Split(partial, "\n")) < 4, "interruption left %q, want a partial namespace", partial)

	runResume(t, root)

	requireTest(t, refsUnder(t, root, "refs/bench/specbuild/", "refs/bench/recovery/") == "", "re-entry did not finish the sweep")
	surviving, err := intent.Assignments(root)
	requireTest(t, err == nil && len(surviving) == 1 && surviving[0].ID == pool.Assignment.ID, "surviving assignments = %#v, %v", surviving, err)
	requireTest(t, gitOutput(t, root, "rev-parse", "refs/bench/green/main") == green, "re-entry moved the gate's green verdict ref")
}
