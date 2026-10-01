package worktree

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
)

func landAssignment(t *testing.T, root string, creation Creation, name string) {
	t.Helper()
	commitInWorktree(t, creation.Path, name, "landed\n", "landed")
	gitRun(t, root, "cherry-pick", strings.TrimPrefix(creation.Assignment.Branch, "refs/heads/"))
}

func TestResumeSummaryCountsLandedAssignments(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	landed := mustCreate(t, root, home, "count-landed", "landed")
	active := mustCreate(t, root, home, "count-active", "active")
	landAssignment(t, root, landed, "landed.txt")
	commitInWorktree(t, active.Path, "active.txt", "active\n", "active")

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	if !strings.Contains(resumed.stdout, "retained active=1 landed=1") {
		t.Fatalf("summary=%q, want active and landed counts", resumed.stdout)
	}
	if _, err := os.Stat(landed.Path); err != nil {
		t.Fatalf("landed worktree was removed: %v", err)
	}
	if _, err := os.Stat(active.Path); err != nil {
		t.Fatalf("active worktree was removed: %v", err)
	}
}

func TestResumeSummaryCountsLandedProofWithoutBranchAdvance(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	mustCreate(t, root, home, "landed-proof", "landed proof")

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	if !strings.Contains(resumed.stdout, "retained landed=1") {
		t.Fatalf("summary=%q, want a landed proof to count without a branch advance", resumed.stdout)
	}
}

func TestResumeSummaryPartitionsLandedLeases(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	live := mustCreate(t, root, home, "lease-live", "live")
	dead := mustCreate(t, root, home, "lease-dead", "dead")
	unknown := mustCreate(t, root, home, "lease-unknown", "unknown")
	landAssignment(t, root, live, "live.txt")
	landAssignment(t, root, dead, "dead.txt")
	landAssignment(t, root, unknown, "unknown.txt")
	liveLease, err := LeaseFile(live.Path)
	mustNoError(t, err)
	mustWrite(t, liveLease, []byte(strconv.Itoa(os.Getpid())+" 2026-07-15T00:00:00Z\n"), 0o600)
	deadLease, err := LeaseFile(dead.Path)
	mustNoError(t, err)
	mustWrite(t, deadLease, []byte(deadPidLine(t)), 0o600)
	unknownLease, err := LeaseFile(unknown.Path)
	mustNoError(t, err)
	mustWrite(t, unknownLease, []byte("not-a-lease\n"), 0o600)

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	summary := resumed.stdout
	if !strings.Contains(summary, "landed=2") || !strings.Contains(summary, "live-lease=1") || strings.Contains(summary, "landed=3") {
		t.Fatalf("summary=%q, want dead/unknown landed and live lease separate", summary)
	}
}

func TestResumeSummaryLiveLeaseWinsOverResidue(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	creation := mustCreate(t, root, home, "live-residue", "live residue")
	landAssignment(t, root, creation, "live-residue.txt")
	mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
	lease, err := LeaseFile(creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(strconv.Itoa(os.Getpid())+" 2026-07-15T00:00:00Z\n"), 0o600)

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	summary := resumed.stdout
	if !strings.Contains(summary, "live-lease=1") || strings.Contains(summary, "landed=") || strings.Contains(summary, "ignored=") {
		t.Fatalf("summary=%q, want live lease precedence", summary)
	}
}

func TestResumeSummaryKeepsLandedClassificationAboveResidue(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	ignored := mustCreate(t, root, home, "landed-ignored", "ignored")
	dirty := mustCreate(t, root, home, "landed-dirty", "dirty")
	landAssignment(t, root, ignored, "ignored-landed.txt")
	landAssignment(t, root, dirty, "dirty-landed.txt")
	mustWrite(t, filepath.Join(ignored.Path, "ignored.txt"), []byte("residue\n"), 0o644)
	mustWrite(t, filepath.Join(dirty.Path, "dirty-landed.txt"), []byte("changed\n"), 0o644)

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	summary := resumed.stdout
	if !strings.Contains(summary, "retained landed=2") || strings.Contains(summary, "ignored=") || strings.Contains(summary, "dirty=") || strings.Contains(summary, "stale-active") {
		t.Fatalf("summary=%q, want residue-independent landed count", summary)
	}
	if _, err := os.Stat(ignored.Path); err != nil {
		t.Fatalf("ignored landed worktree was removed: %v", err)
	}
	if _, err := os.Stat(dirty.Path); err != nil {
		t.Fatalf("dirty landed worktree was removed: %v", err)
	}
}

func TestResumeSummarySeparatesAgedLandedAndActiveAssignments(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	landed := mustCreate(t, root, home, "aged-landed", "aged landed")
	active := mustCreate(t, root, home, "aged-active", "aged active")
	landAssignment(t, root, landed, "aged-landed.txt")
	commitInWorktree(t, active.Path, "aged-active.txt", "active\n", "active")
	backdate(t, root, landed.Assignment, 8*24*time.Hour)
	backdate(t, root, active.Assignment, 8*24*time.Hour)

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	summary := resumed.stdout
	if !strings.Contains(summary, "retained landed=1 stale-active=1;") || !strings.Contains(summary, "\nstale-active "+active.Assignment.ID+": ") || strings.Contains(summary, landed.Assignment.ID) || strings.Contains(summary, "orphan") {
		t.Fatalf("summary=%q, want only the non-landed stale-active line", summary)
	}
}

func TestResumeSummaryAdvertisesLandedSweep(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "advertise-landed", "advertised")
	landAssignment(t, root, creation, "advertised.txt")

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	summary := resumed.stdout
	if strings.Count(summary, "bench worktree clean --landed") != 1 || strings.Contains(summary, "--discard-ignored") {
		t.Fatalf("summary=%q, want one safe landed sweep advertisement", summary)
	}
}

func TestLandedClassifierUnknownDefaultStaysActive(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "unknown-default", "unknown default")
	landAssignment(t, root, creation, "unknown-default.txt")
	gitRun(t, root, "branch", "-m", "main", "trunk")

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 1 || !strings.Contains(resumed.stderr, "no resolvable default branch") {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q, want the existing no-default refusal", resumed.exit, resumed.stdout, resumed.stderr)
	}
	if !strings.Contains(resumed.stdout, "retained active=1") || strings.Contains(resumed.stdout, "landed=") {
		t.Fatalf("summary=%q, want unknown landedness under active", resumed.stdout)
	}
	list := runVerb(t, verbList, repoHome{root, home}.call())
	if list.exit != 0 || strings.Contains(list.stdout, "clean --landed") {
		t.Fatalf("ListCommand exit=%d output=%q, want no landed action", list.exit, list.stdout)
	}
}

func TestLandedClassifierErroredProofStaysActive(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "errored-proof", "errored proof")
	landAssignment(t, root, creation, "errored-proof.txt")
	gitRun(t, root, "update-ref", "-d", creation.Assignment.Branch)

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	if !strings.Contains(resumed.stdout, "retained active=1") || strings.Contains(resumed.stdout, "landed=") {
		t.Fatalf("summary=%q, want errored proof under active", resumed.stdout)
	}
	list := runVerb(t, verbList, repoHome{root, home}.call())
	if list.exit != 0 || strings.Contains(list.stdout, "clean --landed") {
		t.Fatalf("ListCommand exit=%d output=%q, want no landed action", list.exit, list.stdout)
	}
}

func TestLandedClassifierOnlyActiveStateQualifies(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	cleanupPending := mustCreate(t, root, home, "state-cleanup-pending", "cleanup pending")
	recovered := mustCreate(t, root, home, "state-recovered", "recovered")
	complete := mustCreate(t, root, home, "state-complete", "complete")
	landAssignment(t, root, cleanupPending, "cleanup-pending.txt")
	landAssignment(t, root, recovered, "recovered.txt")
	landAssignment(t, root, complete, "complete.txt")
	cleanupPending.Assignment.State = intent.StateCleanupPending
	mustNoError(t, intent.PutAssignment(root, cleanupPending.Assignment))
	recovered.Assignment.State = intent.StateRecovered
	recovered.Assignment.Recovery = []intent.Recovery{{
		Ref:  intent.RecoveryRefPrefix(recovered.Assignment.OwnerID, recovered.Assignment.ID) + "1",
		Root: strings.Repeat("a", 40), Payloads: []string{strings.Repeat("b", 40)},
	}}
	mustNoError(t, intent.PutAssignment(root, recovered.Assignment))
	complete.Assignment.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(root, complete.Assignment))

	before := runVerb(t, verbList, repoHome{root, home}.call())
	if before.exit != 0 || !strings.Contains(before.stdout, cleanupPending.Path) || !strings.Contains(before.stdout, "resume the cleanup-pending assignment") {
		t.Fatalf("ListCommand = (%d, %q), want the cleanup-pending release action", before.exit, before.stdout)
	}
	argv, err := axitest.RecoverHelpCommandArgv(before.stdout)
	if err != nil {
		t.Fatal(err)
	}
	release := []string{"bench", "worktree", "release", "--request", cleanupPending.Assignment.RequestToken, cleanupPending.Path}
	if !slices.Equal(argv, release) {
		t.Fatalf("release action argv = %q, want %q", argv, release)
	}

	resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call())
	if resumed.exit != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", resumed.exit, resumed.stdout, resumed.stderr)
	}
	if strings.Contains(resumed.stdout, "landed=") {
		t.Fatalf("summary=%q, want no landed count for settled states", resumed.stdout)
	}
	list := runVerb(t, verbList, repoHome{root, home}.call())
	if list.exit != 0 || strings.Contains(list.stdout, "clean --landed") {
		t.Fatalf("ListCommand exit=%d output=%q, want no landed action", list.exit, list.stdout)
	}
}

// TestMissingBranchScanKeepsTheActiveStateFilter pins the state filter the any-state
// path match cannot supply. A retired record names a tree whose phase is over, so its
// gone branch is the expected end state, not the recoverable loss the scan reports.
func TestMissingBranchScanKeepsTheActiveStateFilter(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	retired := mustCreate(t, root, home, "missing-branch-retired", "retired")
	gitRun(t, root, "update-ref", "-d", retired.Assignment.Branch)
	retired.Assignment.State = intent.StateCleanupPending
	mustNoError(t, intent.PutAssignment(root, retired.Assignment))

	if got, missing := activeAssignmentWithMissingBranch(root, retired.Path); missing {
		t.Fatalf("activeAssignmentWithMissingBranch = (%q, true), want no assignment for a retired row", got.ID)
	}

	active := mustCreate(t, root, home, "missing-branch-active", "active")
	gitRun(t, root, "update-ref", "-d", active.Assignment.Branch)
	retired.Assignment.Worktree = active.Path
	mustNoError(t, intent.PutAssignment(root, retired.Assignment))

	got, missing := activeAssignmentWithMissingBranch(root, active.Path)
	if !missing || got.ID != active.Assignment.ID {
		t.Fatalf("activeAssignmentWithMissingBranch = (%q, %t), want (%q, true)", got.ID, missing, active.Assignment.ID)
	}
}

func TestListCommandAdvertisesOneLandedSweep(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	landed := mustCreate(t, root, home, "list-landed", "landed")
	active := mustCreate(t, root, home, "list-active", "active")
	landAssignment(t, root, landed, "list-landed.txt")

	list := runVerb(t, verbList, repoHome{root, home}.call())
	if list.exit != 0 {
		t.Fatalf("ListCommand exit=%d output=%q", list.exit, list.stdout)
	}
	out := list.stdout
	if strings.Count(out, "bench worktree clean --landed") != 1 {
		t.Fatalf("ListCommand output=%q, want one landed action", out)
	}
	if strings.Count(out, activePathHelpRow) != 1 || strings.Count(out, activeExecHelpRow) != 1 {
		t.Fatalf("ListCommand output=%q, want one target-slot path action and one target-slot exec action", out)
	}
	for _, id := range []string{landed.Assignment.ID, active.Assignment.ID} {
		if strings.Contains(out, "bench worktree path "+id) || strings.Contains(out, "bench worktree exec "+id) {
			t.Fatalf("ListCommand output=%q, want no per-row action for %s", out, id)
		}
	}
}
