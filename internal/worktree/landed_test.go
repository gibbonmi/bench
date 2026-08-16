package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

func landAssignment(t *testing.T, root string, creation Creation, name string) {
	t.Helper()
	commitInWorktree(t, creation.Path, name, "landed\n", "landed")
	gitRun(t, root, "cherry-pick", strings.TrimPrefix(creation.Assignment.Branch, "refs/heads/"))
}

func TestResumeSummaryCountsLandedAssignments(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	landed := mustCreate(t, root, "count-landed", "landed")
	active := mustCreate(t, root, "count-active", "active")
	landAssignment(t, root, landed, "landed.txt")
	commitInWorktree(t, active.Path, "active.txt", "active\n", "active")
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "retained active=1 landed=1") {
		t.Fatalf("summary=%q, want active and landed counts", stdout.String())
	}
	if _, err := os.Stat(landed.Path); err != nil {
		t.Fatalf("landed worktree was removed: %v", err)
	}
	if _, err := os.Stat(active.Path); err != nil {
		t.Fatalf("active worktree was removed: %v", err)
	}
}

func TestResumeSummaryCountsLandedProofWithoutBranchAdvance(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	mustCreate(t, root, "landed-proof", "landed proof")
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "retained landed=1") {
		t.Fatalf("summary=%q, want a landed proof to count without a branch advance", stdout.String())
	}
}

func TestResumeSummaryPartitionsLandedLeases(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	live := mustCreate(t, root, "lease-live", "live")
	dead := mustCreate(t, root, "lease-dead", "dead")
	unknown := mustCreate(t, root, "lease-unknown", "unknown")
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
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	summary := stdout.String()
	if !strings.Contains(summary, "landed=2") || !strings.Contains(summary, "live-lease=1") || strings.Contains(summary, "landed=3") {
		t.Fatalf("summary=%q, want dead/unknown landed and live lease separate", summary)
	}
}

func TestResumeSummaryLiveLeaseWinsOverResidue(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	creation := mustCreate(t, root, "live-residue", "live residue")
	landAssignment(t, root, creation, "live-residue.txt")
	mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
	lease, err := LeaseFile(creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(strconv.Itoa(os.Getpid())+" 2026-07-15T00:00:00Z\n"), 0o600)
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	summary := stdout.String()
	if !strings.Contains(summary, "live-lease=1") || strings.Contains(summary, "landed=") || strings.Contains(summary, "ignored=") {
		t.Fatalf("summary=%q, want live lease precedence", summary)
	}
}

func TestResumeSummaryKeepsLandedClassificationAboveResidue(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "commit", "-qm", "ignore")
	ignored := mustCreate(t, root, "landed-ignored", "ignored")
	dirty := mustCreate(t, root, "landed-dirty", "dirty")
	landAssignment(t, root, ignored, "ignored-landed.txt")
	landAssignment(t, root, dirty, "dirty-landed.txt")
	mustWrite(t, filepath.Join(ignored.Path, "ignored.txt"), []byte("residue\n"), 0o644)
	mustWrite(t, filepath.Join(dirty.Path, "dirty-landed.txt"), []byte("changed\n"), 0o644)
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	summary := stdout.String()
	if !strings.Contains(summary, "retained landed=2") || strings.Contains(summary, "ignored=") || strings.Contains(summary, "dirty=") || strings.Contains(summary, "orphan ") {
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
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	landed := mustCreate(t, root, "aged-landed", "aged landed")
	active := mustCreate(t, root, "aged-active", "aged active")
	landAssignment(t, root, landed, "aged-landed.txt")
	commitInWorktree(t, active.Path, "aged-active.txt", "active\n", "active")
	backdate(t, root, landed.Assignment, 8*24*time.Hour)
	backdate(t, root, active.Assignment, 8*24*time.Hour)
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	summary := stdout.String()
	if !strings.Contains(summary, "retained landed=1") || !strings.Contains(summary, "orphan "+active.Assignment.ID) || strings.Contains(summary, "orphan "+landed.Assignment.ID) {
		t.Fatalf("summary=%q, want only the non-landed orphan line", summary)
	}
}

func TestResumeSummaryAdvertisesLandedSweep(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, "advertise-landed", "advertised")
	landAssignment(t, root, creation, "advertised.txt")
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	summary := stdout.String()
	if strings.Count(summary, "bench worktree clean --landed") != 1 || strings.Contains(summary, "--discard-ignored") {
		t.Fatalf("summary=%q, want one safe landed sweep advertisement", summary)
	}
}

func TestLandedClassifierUnknownDefaultStaysActive(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, "unknown-default", "unknown default")
	landAssignment(t, root, creation, "unknown-default.txt")
	gitRun(t, root, "branch", "-m", "main", "trunk")
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 1 || !strings.Contains(stderr.String(), "no resolvable default branch") {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q, want the existing no-default refusal", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "retained active=1") || strings.Contains(stdout.String(), "landed=") {
		t.Fatalf("summary=%q, want unknown landedness under active", stdout.String())
	}
	list, code := ListCommand(nil)
	if code != 0 || strings.Contains(list, "clean --landed") {
		t.Fatalf("ListCommand exit=%d output=%q, want no landed action", code, list)
	}
}

func TestLandedClassifierErroredProofStaysActive(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, "errored-proof", "errored proof")
	landAssignment(t, root, creation, "errored-proof.txt")
	gitRun(t, root, "update-ref", "-d", creation.Assignment.Branch)
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "retained active=1") || strings.Contains(stdout.String(), "landed=") {
		t.Fatalf("summary=%q, want errored proof under active", stdout.String())
	}
	list, code := ListCommand(nil)
	if code != 0 || strings.Contains(list, "clean --landed") {
		t.Fatalf("ListCommand exit=%d output=%q, want no landed action", code, list)
	}
}

func TestLandedClassifierOnlyActiveStateQualifies(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	cleanupPending := mustCreate(t, root, "state-cleanup-pending", "cleanup pending")
	recovered := mustCreate(t, root, "state-recovered", "recovered")
	complete := mustCreate(t, root, "state-complete", "complete")
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
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("ResumeCleanCommand exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if strings.Contains(stdout.String(), "landed=") {
		t.Fatalf("summary=%q, want no landed count for settled states", stdout.String())
	}
	list, code := ListCommand(nil)
	if code != 0 || strings.Contains(list, "clean --landed") {
		t.Fatalf("ListCommand exit=%d output=%q, want no landed action", code, list)
	}
}

func TestListCommandAdvertisesOneLandedSweep(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	landed := mustCreate(t, root, "list-landed", "landed")
	active := mustCreate(t, root, "list-active", "active")
	landAssignment(t, root, landed, "list-landed.txt")
	chdir(t, root)

	out, code := ListCommand(nil)
	if code != 0 {
		t.Fatalf("ListCommand exit=%d output=%q", code, out)
	}
	if strings.Count(out, "bench worktree clean --landed") != 1 {
		t.Fatalf("ListCommand output=%q, want one landed action", out)
	}
	for _, id := range []string{landed.Assignment.ID, active.Assignment.ID} {
		if !strings.Contains(out, "bench worktree path "+id) || !strings.Contains(out, "bench worktree exec "+id+" -- <command>") {
			t.Fatalf("ListCommand output=%q, want path/exec actions for %s", out, id)
		}
	}
}
