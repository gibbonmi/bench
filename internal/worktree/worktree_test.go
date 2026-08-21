package worktree

import (
	"bytes"
	"errors"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// cksumGolden pins POSIX cksum(root+newline), including a spaced path.
var cksumGolden = []struct {
	root string
	sum  uint32
}{
	{"/home/mgibs/workspace/bench", 2826441890},
	{"/tmp/a b/c", 889650394},
}

func TestCksumMatchesGolden(t *testing.T) {
	for _, g := range cksumGolden {
		got := cksum([]byte(g.root + "\n"))
		if got != g.sum {
			t.Errorf("cksum(%q+NL) = %d, want %d", g.root, got, g.sum)
		}
	}
}

// TestCksumMatchesSystemTool cross-checks against the live `cksum` when it is on
// PATH, so the pinned goldens can never silently drift from the real tool. It is
// skipped where `cksum` is unavailable, keeping the suite hermetic there.

func TestCksumMatchesSystemTool(t *testing.T) {
	if _, err := exec.LookPath("cksum"); err != nil {
		capability.Capability(t, capability.Tool, "cksum not available")
	}
	for _, g := range cksumGolden {
		// printf '%s\n' "<root>" | cksum
		printf := exec.Command("printf", "%s\n", g.root)
		ck := exec.Command("cksum")
		pipe, err := printf.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		ck.Stdin = pipe
		out, err := ck.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := ck.Start(); err != nil {
			t.Fatal(err)
		}
		if err := printf.Run(); err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 64)
		n, _ := out.Read(buf)
		if err := ck.Wait(); err != nil {
			t.Fatal(err)
		}
		field := strings.Fields(string(buf[:n]))
		if len(field) == 0 {
			t.Fatalf("empty cksum output for %q", g.root)
		}
		want, err := strconv.ParseUint(field[0], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		if got := cksum([]byte(g.root + "\n")); uint64(got) != want {
			t.Errorf("cksum(%q) = %d, system tool = %d", g.root, got, want)
		}
	}
}

func TestPool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := "/home/mgibs/workspace/bench"
	want := filepath.Join(home, "worktrees", "bench-2826441890")
	if got := Pool(root); got != want {
		t.Errorf("Pool(%q) = %q, want %q", root, got, want)
	}
}

func TestPoolDefaultBenchHome(t *testing.T) {
	// With BENCH_HOME unset, Pool falls back to <home>/.bench.
	t.Setenv("BENCH_HOME", "")
	root := "/tmp/a b/c"
	got := Pool(root)
	suffix := filepath.Join(".bench", "worktrees", "c-889650394")
	if !strings.HasSuffix(got, suffix) {
		t.Errorf("Pool(%q) = %q, want suffix %q", root, got, suffix)
	}
}

func TestClassifyRegisteredWorktrees(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	pool := Pool(root)
	warm := filepath.Join(pool, "warm")
	leased := filepath.Join(pool, "leased")
	outOfPool := filepath.Join(filepath.Dir(root), "outside pool")
	if err := os.MkdirAll(pool, 0o755); err != nil {
		t.Fatalf("mkdir pool: %v", err)
	}
	gitRun(t, root, "worktree", "add", "-q", "--detach", warm, "HEAD")
	gitRun(t, root, "worktree", "add", "-q", "--detach", leased, "HEAD")
	gitRun(t, root, "worktree", "add", "-q", "--detach", outOfPool, "HEAD")
	lease, err := LeaseFile(leased)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	if err := os.WriteFile(lease, []byte("123 2026-07-06T00:00:00Z\n"), 0o644); err != nil {
		t.Fatalf("write lease: %v", err)
	}
	entries, err := ClassifyRegisteredWorktrees(root)
	if err != nil {
		t.Fatalf("ClassifyRegisteredWorktrees: %v", err)
	}
	got := map[string]Class{}
	for _, entry := range entries {
		got[entry.Path] = entry.Class
	}
	want := map[string]Class{
		root:      ClassRoot,
		warm:      ClassPoolWarm,
		leased:    ClassPoolLease,
		outOfPool: ClassOutOfPool,
	}
	for path, class := range want {
		if got[path] != class {
			t.Errorf("class %q = %q, want %q", path, got[path], class)
		}
	}
	linkedEntries, err := ClassifyRegisteredWorktrees(leased)
	if err != nil {
		t.Fatalf("ClassifyRegisteredWorktrees from linked worktree: %v", err)
	}
	got = map[string]Class{}
	for _, entry := range linkedEntries {
		got[entry.Path] = entry.Class
	}
	for path, class := range want {
		if got[path] != class {
			t.Errorf("class from linked cwd %q = %q, want %q", path, got[path], class)
		}
	}
}

func TestCleanupDeletesOnlyExactBranchAndSparesSiblingRefs(t *testing.T) {
	t.Run("clean assignment compacts and spares sibling", func(t *testing.T) {
		root, target := newOwnedAssignment(t, "terminal-clean")
		sibling, err := Create(root, "terminal-clean-sibling", "sibling", nil)
		if err != nil {
			t.Fatal(err)
		}
		siblingRef := "refs/bench/recovery/" + sibling.Assignment.OwnerID + "/" + sibling.Assignment.ID + "/1"
		gitRun(t, root, "update-ref", siblingRef, target.Assignment.Start)
		markPending(t, root, target.Assignment)
		if _, err := ApplyAutomatic(root, target.Path, nil); err != nil {
			t.Fatal(err)
		}
		if exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", target.Assignment.Branch).Run() == nil {
			t.Fatal("exact cleanup left target branch")
		}
		gitRun(t, root, "show-ref", "--verify", "--quiet", sibling.Assignment.Branch)
		gitRun(t, root, "show-ref", "--verify", "--quiet", siblingRef)
		assignments, err := intent.Assignments(root)
		if err != nil || len(assignments) != 1 || assignments[0].ID != sibling.Assignment.ID {
			t.Fatalf("clean compaction assignments = %#v, %v", assignments, err)
		}
	})
}

// TestCreateCommandPrintsNextHint pins the next-step hint CreateCommand appends after
// its worktree_create table: two literal lines addressing the freshly created worktree
// by the actual --label value, so a caller never has to invent the exec/path syntax.
func TestCreateCommandPrintsNextHint(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	var stdout, stderr bytes.Buffer
	code := CreateCommand(root, []string{"--request", "next-hint", "--label", "next hint label"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("CreateCommand exit = %d, stderr = %q", code, stderr.String())
	}
	got := stdout.String()
	if !strings.Contains(got, "next[2]:\n") {
		t.Fatalf("CreateCommand output missing next[2] header: %q", got)
	}
	if !strings.Contains(got, `  bench worktree exec "next hint label" -- <command>`+"\n") {
		t.Fatalf("CreateCommand output missing exec hint line: %q", got)
	}
	if !strings.Contains(got, `  bench worktree path "next hint label"`+"\n") {
		t.Fatalf("CreateCommand output missing path hint line: %q", got)
	}
}

// TestReleaseSurfacesRetainedVerdict pins FT93(a): when the automatic plan retains
// (here, an ignored residual the safe planner will not discard), release must report
// the retained verdict and the exact next command, not the internal-bookkeeping
// "terminal receipt missing". The pre-fix path discards the retain plan, finds no
// terminal receipt, and returns the masking error; this goes red on that message.
func TestReleaseSurfacesRetainedVerdict(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("residual.txt\n"), 0o644)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, "retain-verdict", "retain verdict")
	residual := filepath.Join(creation.Path, "residual.txt")
	mustWrite(t, residual, []byte("build output\n"), 0o600)
	requirePlanAction(t, root, creation.Path, ActionRetain)

	args := []string{"--request", "retain-verdict", creation.Path}
	var out, errb strings.Builder
	code := ReleaseCommand(root, args, &out, &errb)
	msg := errb.String()
	requireTest(t, code != 0, "retained release exit = %d, want non-zero", code)
	requireTest(t, !strings.Contains(msg, "terminal receipt missing"), "masking error still present: %q", msg)
	requireTest(t, strings.Contains(msg, "retained") && strings.Contains(msg, "bench worktree release"),
		"retained verdict is not actionable: %q", msg)

	mustNoError(t, os.Remove(residual))
	var out2 strings.Builder
	code = ReleaseCommand(root, args, &out2, io.Discard)
	requireTest(t, code == 0, "recovery release exit = %d, want 0; out=%q", code, out2.String())
}

func TestReleaseUnknownRequestNamesReauthorizeRecovery(t *testing.T) {
	root, creation := newOwnedAssignment(t, "release-reauthorize-recovery")
	var stdout, stderr strings.Builder
	code := ReleaseCommand(root, []string{"--request", "unknown-request", creation.Path}, &stdout, &stderr)
	wantNext := "bench worktree reauthorize --assignment " + creation.Assignment.ID + " --request <new-request> --base <full-base-commit> --source-tip <full-source-tip-commit> '" + creation.Path + "'"
	want := "bench worktree release: request, assignment, or path mismatch; checkout retained; observed=assignment:" + creation.Assignment.ID + ",next=" + wantNext + "\n"
	if code != 1 || stdout.String() != "" || stderr.String() != want {
		t.Fatalf("unknown-request release = (%d, %q, %q), want exit 1 and stderr %q", code, stdout.String(), stderr.String(), want)
	}
}

// removeOutOfBand simulates a request-less `bench worktree clean --discard-ignored
// --apply` that removed an owned tree: it drops the git registration and directory and
// writes the completed explicit-clean cleanup receipt (owned, request-bound, no
// automatic-registration fingerprint), leaving the assignment record stranded.
func removeOutOfBand(t *testing.T, root string, a intent.Assignment, action CleanupAction) {
	t.Helper()
	gitRun(t, root, "worktree", "remove", "-f", "-f", a.Worktree)
	repo, target, err := cleanupIdentity(root, a.Worktree)
	mustNoError(t, err)
	mustNoError(t, intent.PutCleanupReceipt(root, intent.CleanupReceipt{
		Schema: intent.CleanupReceiptSchema, Repo: repo, Operation: cleanupOperation,
		Target: target, Fingerprint: strings.Repeat("c", 64), State: intent.ReceiptComplete,
		Phase: intent.ReceiptPhaseTerminal, Action: string(action), Tracked: a.ID,
		Ignored: "count=0 bytes=0 shown=0 truncated=false", Recovery: "none",
		Owned: true, Owner: a.OwnerID, Assignment: a.ID, Request: a.Request,
	}))
}

// TestReleaseReconcilesOutOfBandResidue pins FT93(b), residue path: a release whose
// tree was removed out of band, holding no preserved work, reconciles and compacts the
// record instead of dead-ending on "cleanup receipt does not authorize release
// reconciliation". Replay is idempotent.
func TestReleaseReconcilesOutOfBandResidue(t *testing.T) {
	root, creation := newOwnedAssignment(t, "oob-residue")
	a, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, len(a.Recovery) == 0, "fixture already holds recovery metadata")
	removeOutOfBand(t, root, a, ActionRemoved)

	args := []string{"--request", "landed-oob-residue", creation.Path}
	var out, errb strings.Builder
	code := ReleaseCommand(root, args, &out, &errb)
	requireTest(t, code == 0, "residue release exit=%d stderr=%q", code, errb.String())
	if _, err := assignmentByID(root, a.ID); err == nil {
		t.Fatal("residue record survived reconcile")
	}
	var replay strings.Builder
	code = ReleaseCommand(root, args, &replay, io.Discard)
	requireTest(t, code == 0 && replay.String() == out.String(), "replay exit=%d out=%q", code, replay.String())
}

// TestReleaseNamesRecoveryForPreservedOrphan pins FT93(b), preserved path: a release
// whose tree was removed out of band but still holds preserved work returns a verdict
// handing over the ref itself and leaves the record and its recovery pointer intact —
// release never silently discards preserved work.
func TestReleaseNamesRecoveryForPreservedOrphan(t *testing.T) {
	root, creation := newOwnedAssignment(t, "oob-preserved")
	a, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	ref := intent.RecoveryRefPrefix(a.OwnerID, a.ID) + "1"
	a.State, a.Recovery = intent.StateRecovered, []intent.Recovery{{Ref: ref, Root: strings.Repeat("a", 40), Payloads: []string{strings.Repeat("b", 40)}}}
	mustNoError(t, intent.PutAssignment(root, a))
	removeOutOfBand(t, root, a, ActionRemoved)

	var out, errb strings.Builder
	code := ReleaseCommand(root, []string{"--request", "landed-oob-preserved", creation.Path}, &out, &errb)
	requireTest(t, code != 0, "preserved release exit=%d, want non-zero", code)
	requireTest(t, strings.Contains(errb.String(), "git show "+ref),
		"preserved verdict does not hand over the ref: %q", errb.String())
	got, err := assignmentByID(root, a.ID)
	requireTest(t, err == nil && len(got.Recovery) == 1, "preserved record was mutated or deleted: %v", err)
}

// TestResumeReconcilesTreeGoneRecordsAndSparesYoungActive pins the standing cleaner's
// blast radius over the ledger: a tree-gone record is dropped whether it was mid-cleanup
// or holding preserved work the removed lifecycle wrote, while a record whose tree still
// exists survives untouched.
//
// What holds the active, tree-gone record here is its age, not its state: the reconcile
// drops an orphaned active record, and this one survives only because it was stamped
// moments ago and so is not aged. That is the race this fixture guards — a reconcile that
// dropped on tree-absence alone would catch a session between `worktree add` and its
// first write.
func TestResumeReconcilesTreeGoneRecordsAndSparesYoungActive(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	residue := mustCreate(t, root, "landed-sweep-residue", "residue")
	preserved := mustCreate(t, root, "landed-sweep-preserved", "preserved")
	activeGone := mustCreate(t, root, "landed-sweep-active", "active gone")
	live := mustCreate(t, root, "landed-sweep-live", "live present")

	ra, err := assignmentByID(root, residue.Assignment.ID)
	mustNoError(t, err)
	gitRun(t, root, "worktree", "remove", "-f", "-f", ra.Worktree)
	ra.State = intent.StateCleanupPending
	mustNoError(t, intent.PutAssignment(root, ra))

	pa, err := assignmentByID(root, preserved.Assignment.ID)
	mustNoError(t, err)
	gitRun(t, root, "worktree", "remove", "-f", "-f", pa.Worktree)
	ref := intent.RecoveryRefPrefix(pa.OwnerID, pa.ID) + "1"
	pa.State, pa.Recovery = intent.StateRecovered, []intent.Recovery{{Ref: ref, Root: strings.Repeat("a", 40), Payloads: []string{strings.Repeat("b", 40)}}}
	mustNoError(t, intent.PutAssignment(root, pa))

	ag, err := assignmentByID(root, activeGone.Assignment.ID)
	mustNoError(t, err)
	gitRun(t, root, "worktree", "remove", "-f", "-f", ag.Worktree) // active, tree gone, unregistered

	result, err := ConservativeCleanup(root)
	mustNoError(t, err)
	requireTest(t, result.Reconciled == 2, "Reconciled=%d, want 2", result.Reconciled)
	for _, dropped := range []string{ra.ID, pa.ID} {
		if _, err := assignmentByID(root, dropped); err == nil {
			t.Fatalf("tree-gone record %s survived the reconcile", dropped)
		}
	}
	for _, keep := range []string{ag.ID, live.Assignment.ID} {
		if _, err := assignmentByID(root, keep); err != nil {
			t.Fatalf("record %s was dropped but must survive: %v", keep, err)
		}
	}
}

func TestReleaseReconcilesInFlightAutomaticCleanup(t *testing.T) {
	root, creation := newPendingAssignment(t, "release-in-flight")
	stop := errors.New("crash after removal")
	_, err := ApplyAutomatic(root, creation.Path, func(step LifecycleStep) error {
		if step == StepRemoval {
			return stop
		}
		return nil
	})
	requireTest(t, errors.Is(err, stop), "automatic interruption = %v", err)
	args := []string{"--request", "landed-release-in-flight", creation.Path}
	var first, firstErr strings.Builder
	code := ReleaseCommand(root, args, &first, &firstErr)
	requireTest(t, code == 0 && firstErr.String() == "", "in-flight release code=%d stderr=%q", code, firstErr.String())
	var replay strings.Builder
	code = ReleaseCommand(root, args, &replay, io.Discard)
	requireTest(t, code == 0 && replay.String() == first.String(), "in-flight replay code=%d stdout=%q", code, replay.String())
	requireTest(t, ReleaseCommand(root, []string{"--request", "changed", creation.Path}, io.Discard, io.Discard) != 0, "changed request authorized")
	requireTest(t, ReleaseCommand(root, []string{"--request", args[1], root}, io.Discard, io.Discard) != 0, "changed path authorized")
}
func TestExplicitApplyRejectsContentDriftWithoutMutation(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	target := filepath.Join(filepath.Dir(root), "content drift target")
	gitRun(t, root, "worktree", "add", "-q", "--detach", target, "HEAD")
	file := filepath.Join(target, "untracked.txt")
	if err := os.WriteFile(file, []byte("planned bytes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	plan, err := PlanExplicit(root, target)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Fingerprint == "" {
		t.Fatal("explicit detached plan has no fingerprint")
	}
	if err := os.WriteFile(file, []byte("drifted bytes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	beforeWorktrees := gitOutput(t, root, "worktree", "list", "--porcelain")
	beforeRefs := gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/recovery/")
	current, err := ApplyExplicit(root, target, plan.Fingerprint)
	if !errors.Is(err, errStaleFingerprint) {
		t.Fatalf("ApplyExplicit content drift error = %v, want stale fingerprint (current=%#v)", err, current)
	}
	if current.Fingerprint == plan.Fingerprint {
		t.Fatal("content drift did not change the current plan fingerprint")
	}
	if got := gitOutput(t, root, "worktree", "list", "--porcelain"); got != beforeWorktrees {
		t.Fatalf("stale apply mutated registration\nbefore=%s\nafter=%s", beforeWorktrees, got)
	}
	if got := gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/recovery/"); got != beforeRefs {
		t.Fatalf("stale apply created recovery ref\nbefore=%s\nafter=%s", beforeRefs, got)
	}
	if body, err := os.ReadFile(file); err != nil || string(body) != "drifted bytes\n" {
		t.Fatalf("stale apply changed target content: %q, %v", body, err)
	}
}

func TestIgnoredInventoryStatRaceRetains(t *testing.T) {
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	if err := os.WriteFile(filepath.Join(root, ".git", "info", "exclude"), []byte("ignored.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(filepath.Dir(root), "ignored stat race")
	gitRun(t, root, "worktree", "add", "-q", "-b", "ignored-stat-race", target, "HEAD")
	ignored := filepath.Join(target, "ignored.txt")
	if err := os.WriteFile(ignored, []byte("secret\n"), 0o000); err != nil {
		t.Fatal(err)
	}
	original := ignoredLstat
	ignoredLstat = func(path string) (os.FileInfo, error) {
		if path == ignored {
			return nil, os.ErrNotExist
		}
		return os.Lstat(path)
	}
	t.Cleanup(func() { ignoredLstat = original })
	plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{DiscardIgnored: true})
	if err != nil || plan.Action != ActionRetain || plan.ReasonCode != ReasonUncertain {
		t.Fatalf("stat-race plan = %#v, %v", plan, err)
	}
	if _, err := os.Lstat(ignored); err != nil {
		t.Fatalf("stat-race plan mutated ignored file: %v", err)
	}
}

func TestLeaseFile(t *testing.T) {
	dir := gittest.Repo(t)
	lease, err := LeaseFile(dir)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	if !strings.HasSuffix(lease, "bench-lease") {
		t.Errorf("LeaseFile = %q, want suffix bench-lease", lease)
	}
	if !filepath.IsAbs(lease) {
		t.Errorf("LeaseFile = %q, want absolute — a relative path resolves against the caller's CWD, not the worktree", lease)
	}
}

func TestLeaseFileCommandMissingArg(t *testing.T) {
	out, code := LeaseFileCommand(nil)
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if !strings.HasPrefix(out, "usage:") {
		t.Errorf("out = %q, want usage line", out)
	}
}

func TestPoolCommandExplicitRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	out, code := PoolCommand([]string{"/home/mgibs/workspace/bench"})
	if code != 0 {
		t.Errorf("exit = %d, want 0", code)
	}
	want := filepath.Join(home, "worktrees", "bench-2826441890") + "\n"
	if out != want {
		t.Errorf("out = %q, want %q", out, want)
	}
}
func newWorktreeRepo(t testing.TB) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write tracked.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("initial\n"), 0o644); err != nil {
		t.Fatalf("write README.md: %v", err)
	}
	gitRun(t, root, "add", "tracked.txt", "README.md")
	gitRun(t, root, "commit", "-q", "-m", "base")
	return root
}
func chdir(t testing.TB, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
}
func gitOutput(t testing.TB, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}
func gitRun(t testing.TB, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}
