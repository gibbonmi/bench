package worktree

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/usage"
)

func landedSetFixture(t *testing.T) (string, Creation, Creation, Creation) {
	t.Helper()
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	first := mustCreate(t, root, "landed-set-first", "first label")
	second := mustCreate(t, root, "landed-set-second", "second label")
	dirty := mustCreate(t, root, "landed-set-dirty", "dirty label")
	landAssignment(t, root, first, "first.txt")
	landAssignment(t, root, second, "second.txt")
	landAssignment(t, root, dirty, "dirty.txt")
	mustWrite(t, filepath.Join(dirty.Path, "dirty.txt"), []byte("changed\n"), 0o644)
	return root, first, second, dirty
}

func runCleanLanded(t *testing.T, root string, args ...string) (string, string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := CleanCommand(root, args, &stdout, &stderr)
	return stdout.String(), stderr.String(), code
}

func landedRowFingerprint(t *testing.T, output string) string {
	t.Helper()
	match := regexp.MustCompile(`,\"?([0-9a-f]{64})\"?,`).FindStringSubmatch(output)
	if len(match) != 2 {
		t.Fatalf("output has no row fingerprint: %q", output)
	}
	return match[1]
}

func TestCleanLandedPlansRepositoryWideSet(t *testing.T) {
	root, first, second, dirty := landedSetFixture(t)
	stdout, stderr, code := runCleanLanded(t, root, "--landed")
	if code != 0 {
		t.Fatalf("CleanCommand exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	wantByID := map[string]string{
		first.Assignment.ID: first.Path, second.Assignment.ID: second.Path, dirty.Assignment.ID: dirty.Path,
	}
	ids := []string{first.Assignment.ID, second.Assignment.ID, dirty.Assignment.ID}
	sort.Strings(ids)
	positions := make([]int, len(ids))
	for i, id := range ids {
		positions[i] = strings.Index(stdout, wantByID[id])
		if positions[i] < 0 {
			t.Fatalf("output=%q, want assignment %s", stdout, id)
		}
		if i > 0 && positions[i] <= positions[i-1] {
			t.Fatalf("output=%q, want assignment-id order %v", stdout, ids)
		}
	}
	if strings.Count(stdout, ",remove,") != 2 || !strings.Contains(stdout, ",retain,dirty,") || strings.Contains(stdout, "refs/bench/recovery/") {
		t.Fatalf("output=%q, want two removes, one dirty retain, and no recovery ref", stdout)
	}
}

func TestCleanLandedPlanSharesOneFingerprint(t *testing.T) {
	root, _, _, _ := landedSetFixture(t)
	stdout, stderr, code := runCleanLanded(t, root, "--landed")
	if code != 0 {
		t.Fatalf("CleanCommand exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	matches := regexp.MustCompile(`[0-9a-f]{64}`).FindAllString(stdout, -1)
	if len(matches) < 3 || matches[0] != matches[1] || matches[1] != matches[2] {
		t.Fatalf("output=%q, want one shared row fingerprint", stdout)
	}
}

func TestCleanLandedFingerprintBindsSetMembership(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	first := mustCreate(t, root, "fingerprint-member-first", "first")
	landAssignment(t, root, first, "first.txt")
	before, beforeErr, beforeCode := runCleanLanded(t, root, "--landed")
	if beforeCode != 0 || beforeErr != "" {
		t.Fatalf("first plan exit=%d stdout=%q stderr=%q", beforeCode, before, beforeErr)
	}

	second := mustCreate(t, root, "fingerprint-member-second", "second")
	landAssignment(t, root, second, "second.txt")
	after, afterErr, afterCode := runCleanLanded(t, root, "--landed")
	if afterCode != 0 || afterErr != "" {
		t.Fatalf("second plan exit=%d stdout=%q stderr=%q", afterCode, after, afterErr)
	}
	if landedRowFingerprint(t, before) == landedRowFingerprint(t, after) {
		t.Fatalf("set fingerprint did not change when membership changed: before=%q after=%q", before, after)
	}
}

func TestCleanLandedPlanAdvertisesApplyAndRemedies(t *testing.T) {
	root, first, second, dirty := landedSetFixture(t)
	stdout, stderr, code := runCleanLanded(t, root, "--landed")
	if code != 0 {
		t.Fatalf("CleanCommand exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	fingerprint := regexp.MustCompile(`[0-9a-f]{64}`).FindString(stdout)
	if !strings.Contains(stdout, "bench worktree clean --landed --apply "+fingerprint) || !strings.Contains(stdout, "bench worktree clean "+dirty.Path) {
		t.Fatalf("output=%q, want apply and dirty-row remedy", stdout)
	}
	for _, creation := range []Creation{first, second, dirty} {
		if _, err := os.Stat(creation.Path); err != nil {
			t.Fatalf("bare plan removed %s: %v", creation.Path, err)
		}
	}
}

func TestCleanLandedEmptySetExitsClean(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	for attempt := 0; attempt < 2; attempt++ {
		stdout, stderr, code := runCleanLanded(t, root, "--landed")
		if code != 0 || stdout != "worktree_cleanup[0]{target,action,tracked,ignored,recovery,fingerprint,detail}:\n" || stderr != "" {
			t.Fatalf("attempt %d exit=%d stdout=%q stderr=%q", attempt, code, stdout, stderr)
		}
	}
}

func TestCleanLandedApplyOnEmptySetRefused(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	stdout, stderr, code := runCleanLanded(t, root, "--landed", "--apply", strings.Repeat("a", 64))
	if code != 2 || !strings.Contains(stdout, "invalid invocation; run "+usage.WorktreeClean) || stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestCleanLandedRefusesPathOperand(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	wantUsage := "bench worktree clean [--discard-ignored] [--discard-branch] [--full] (<path> | --landed) [--apply <fingerprint>]"
	for _, args := range [][]string{{"--landed", root}, {root, "--landed"}} {
		stdout, stderr, code := runCleanLanded(t, root, args...)
		if code != 2 || !strings.Contains(stdout, wantUsage) || stderr != "" {
			t.Fatalf("args=%q exit=%d stdout=%q stderr=%q", args, code, stdout, stderr)
		}
	}
}

func TestCleanLandedRefusesMalformedFingerprint(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	for _, fingerprint := range []string{strings.Repeat("a", 63), strings.Repeat("a", 65), strings.Repeat("g", 64), strings.Repeat("A", 64)} {
		stdout, stderr, code := runCleanLanded(t, root, "--landed", "--apply", fingerprint)
		if code != 2 || !strings.Contains(stdout, usage.WorktreeClean) || stderr != "" {
			t.Fatalf("fingerprint=%q exit=%d stdout=%q stderr=%q", fingerprint, code, stdout, stderr)
		}
	}
}

func TestCleanLandedSelectorPartition(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	unknown := mustCreate(t, root, "selector-unknown", "unknown lease")
	live := mustCreate(t, root, "selector-live", "live lease")
	unlanded := mustCreate(t, root, "selector-unlanded", "unlanded")
	pending := mustCreate(t, root, "selector-pending", "pending")
	recovered := mustCreate(t, root, "selector-recovered", "recovered")
	complete := mustCreate(t, root, "selector-complete", "complete")
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

	stdout, stderr, code := runCleanLanded(t, root, "--landed")
	if code != 0 || stderr != "" || !strings.Contains(stdout, unknown.Path+",retain,") || !strings.Contains(stdout, "assignment lease state is unknown") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	for _, excluded := range []Creation{live, unlanded, pending, recovered, complete} {
		if strings.Contains(stdout, excluded.Path) {
			t.Fatalf("output=%q, unexpectedly selected %s", stdout, excluded.Path)
		}
	}
}
