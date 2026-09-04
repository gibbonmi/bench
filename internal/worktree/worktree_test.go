package worktree

import (
	"bytes"
	"errors"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestPool(t *testing.T) {
	home := t.TempDir()
	bindEnv(t, "BENCH_HOME", home)
	root := "/home/mgibs/workspace/bench"
	want := filepath.Join(home, "worktrees", "bench-2826441890")
	if got := Pool(root); got != want {
		t.Errorf("Pool(%q) = %q, want %q", root, got, want)
	}
}

func TestPoolDefaultBenchHome(t *testing.T) {
	// With BENCH_HOME unset, Pool falls back to <home>/.bench.
	bindEnv(t, "BENCH_HOME", "")
	root := "/tmp/a b/c"
	got := Pool(root)
	suffix := filepath.Join(".bench", "worktrees", "c-889650394")
	if !strings.HasSuffix(got, suffix) {
		t.Errorf("Pool(%q) = %q, want suffix %q", root, got, suffix)
	}
}

func TestClassifyRegisteredWorktrees(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	root := newWorktreeRepo(t)
	pool := poolAt(home, root)
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
	entries, err := classifyRegisteredWorktreesAt(root, home)
	if err != nil {
		t.Fatalf("classifyRegisteredWorktreesAt: %v", err)
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
	linkedEntries, err := classifyRegisteredWorktreesAt(leased, home)
	if err != nil {
		t.Fatalf("classifyRegisteredWorktreesAt from linked worktree: %v", err)
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
	markProof(t, "lifecycle/journey/registration")
}

func TestCleanupDeletesOnlyExactBranchAndSparesSiblingRefs(t *testing.T) {
	t.Parallel()
	t.Run("clean assignment compacts and spares sibling", func(t *testing.T) {
		root, target, home := newOwnedAssignment(t, "terminal-clean")
		sibling := mustCreate(t, root, home, "terminal-clean-sibling", "sibling")
		siblingRef := "refs/bench/recovery/" + sibling.Assignment.OwnerID + "/" + sibling.Assignment.ID + "/1"
		gitRun(t, root, "update-ref", siblingRef, target.Assignment.Start)
		markPending(t, root, target.Assignment)
		if _, err := ApplyAutomatic(root, target.Path, nil); err != nil {
			t.Fatal(err)
		}
		if descendant(t, "git", "-C", root, "show-ref", "--verify", "--quiet", target.Assignment.Branch).Run() == nil {
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

// TestCreateCommandPrintsNextHint pins the next-step hint CreateCommand appends
// after its worktree_create table: two literal lines addressing the freshly
// created worktree by the actual --label value. A caller never has to invent
// the exec/path syntax.
func TestCreateCommandPrintsNextHint(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	var stdout, stderr bytes.Buffer
	code := CreateCommand(root, home, []string{"--request", "next-hint", "--label", "next hint label"}, &stdout, &stderr)
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

// TestCreateCommandAnswersHelpSpellings pins the create grammar move onto
// usage.Parse: every help spelling prints the declared grammar on stdout and
// exits 0, whether or not required flags are present.
func TestCreateCommandAnswersHelpSpellings(t *testing.T) {
	t.Parallel()
	// WF31: the grammar the command prints is exact, so the flag's spelling and its
	// optional brackets are pinned against the literal rather than against the const.
	if usage.WorktreeCreate != strings.TrimPrefix(createFromGrammar, "usage: ") {
		t.Fatalf("create grammar = %q, want %q", usage.WorktreeCreate, createFromGrammar)
	}
	want := createFromGrammar + "\n"
	for _, args := range [][]string{
		{"--help"},
		{"-h"},
		{"help"},
		{"--request", "x", "--help"},
	} {
		var stdout, stderr bytes.Buffer
		code := CreateCommand("", Home(), args, &stdout, &stderr)
		if code != 0 || stdout.String() != want || stderr.Len() != 0 {
			t.Fatalf("CreateCommand(%q) = (%d, %q, %q), want (0, %q, empty)", args, code, stdout.String(), stderr.String(), want)
		}
	}
}

func TestReleaseCommandHelpAndInvalidArguments(t *testing.T) {
	t.Parallel()
	want := "usage: " + usage.WorktreeRelease + "\n"
	for _, tc := range []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{name: "help", args: []string{"--help"}, wantCode: 0, wantStdout: want},
		{name: "invalid", args: []string{"invalid"}, wantCode: 2, wantStderr: want},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := ReleaseCommand("", Home(), tc.args, &stdout, &stderr)
			if code != tc.wantCode || stdout.String() != tc.wantStdout || stderr.String() != tc.wantStderr {
				t.Fatalf("ReleaseCommand(%q) = (%d, %q, %q), want (%d, %q, %q)", tc.args, code, stdout.String(), stderr.String(), tc.wantCode, tc.wantStdout, tc.wantStderr)
			}
		})
	}
}

// TestCreateCommandHelpPerformsNoRefresh pins that a help request never fetches:
// usage.Parse answers --help before refreshop.Consume ever sees the args, so a
// --refresh alongside --help prints only the help line, not a worktree_refresh table.
func TestCreateCommandHelpPerformsNoRefresh(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	code := CreateCommand("", Home(), []string{"--request", "x", "--refresh", "--help"}, &stdout, &stderr)
	want := "usage: " + usage.WorktreeCreate + "\n"
	if code != 0 || stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("CreateCommand with --refresh --help = (%d, %q, %q), want (0, %q, empty)", code, stdout.String(), stderr.String(), want)
	}
}

// TestCreateCommandRequiredFlagsKeepDeclaredHelp pins that a missing required
// flag exits 2 with the declared grammar, matching the reauthorize sibling.
func TestCreateCommandRequiredFlagsKeepDeclaredHelp(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"--request", "r"},
		{"--label", "l"},
	} {
		var stdout, stderr bytes.Buffer
		if code := CreateCommand("", Home(), args, &stdout, &stderr); code != 2 || stdout.Len() != 0 || stderr.String() != createGrammar.Help+"\n" {
			t.Fatalf("CreateCommand(%q) = (%d, %q, %q), want (2, empty, %q)", args, code, stdout.String(), stderr.String(), createGrammar.Help+"\n")
		}
	}
}

// TestCreateCommandRejectsEmptyFlagValues pins the shared empty-value rule on
// --request and --label: an empty string names nothing and exits 2 naming it.
func TestCreateCommandRejectsEmptyFlagValues(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"--request", "", "--label", "l"},
		{"--request", "r", "--label", ""},
	} {
		var stdout, stderr bytes.Buffer
		if code := CreateCommand("", Home(), args, &stdout, &stderr); code != 2 || stdout.Len() != 0 {
			t.Fatalf("CreateCommand(%q) = (%d, %q, %q), want exit 2 with empty stdout", args, code, stdout.String(), stderr.String())
		}
	}
}

// TestReleaseSurfacesRetainedVerdict pins FT93(a). Here the automatic plan
// retains an ignored residual the safe planner will not discard. Release must
// report the retained verdict and the exact next command, not the
// internal-bookkeeping "terminal receipt missing". A path that discards the
// retain plan finds no terminal receipt and returns the masking error; this
// test goes red on that message.
func TestReleaseSurfacesRetainedVerdict(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("residual.txt\n"), 0o644)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "retain-verdict", "retain verdict")
	residual := filepath.Join(creation.Path, "residual.txt")
	mustWrite(t, residual, []byte("build output\n"), 0o600)
	requirePlanAction(t, root, creation.Path, ActionRetain)

	args := []string{"--request", "retain-verdict", creation.Path}
	var out, errb strings.Builder
	code := ReleaseCommand(root, home, args, &out, &errb)
	msg := errb.String()
	requireTest(t, code != 0, "retained release exit = %d, want non-zero", code)
	requireTest(t, !strings.Contains(msg, "terminal receipt missing"), "masking error still present: %q", msg)
	requireTest(t, strings.Contains(msg, "retained") && strings.Contains(msg, "bench worktree release"),
		"retained verdict is not actionable: %q", msg)

	mustNoError(t, os.Remove(residual))
	var out2 strings.Builder
	code = ReleaseCommand(root, home, args, &out2, io.Discard)
	requireTest(t, code == 0, "recovery release exit = %d, want 0; out=%q", code, out2.String())
}

// TestReleaseUnknownRequestNamesReauthorizeRecovery is LR19: release names the request
// component, its own retained clause, and the same recovery command the landing names.
func TestReleaseUnknownRequestNamesReauthorizeRecovery(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "release-reauthorize-recovery")
	var stdout, stderr strings.Builder
	code := ReleaseCommand(root, home, []string{"--request", "unknown-request", creation.Path}, &stdout, &stderr)
	wantNext := "bench worktree reauthorize --assignment " + creation.Assignment.ID + " --request <new-request> --base <full-base-commit> --source-tip <full-source-tip-commit> '" + creation.Path + "'"
	want := "bench worktree release: request token matches no assignment; checkout retained; observed=assignment:" + creation.Assignment.ID + ",next=" + wantNext + "\n"
	if code != 1 || stdout.String() != "" || stderr.String() != want {
		t.Fatalf("unknown-request release = (%d, %q, %q), want exit 1 and stderr %q", code, stdout.String(), stderr.String(), want)
	}
}

// removeOutOfBand simulates a request-less `bench worktree clean
// --discard-ignored --apply` that removed an owned tree. It drops the git
// registration and directory, and writes the completed explicit-clean cleanup
// receipt (owned, request-bound, no automatic-registration fingerprint),
// leaving the assignment record stranded.
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

// TestReleaseReconcilesOutOfBandResidue pins FT93(b), residue path. Here a
// release's tree was removed out of band and holds no preserved work. It
// reconciles and compacts the record instead of dead-ending on "cleanup
// receipt does not authorize release reconciliation". Replay is idempotent.
func TestReleaseReconcilesOutOfBandResidue(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "oob-residue")
	a, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, len(a.Recovery) == 0, "fixture already holds recovery metadata")
	removeOutOfBand(t, root, a, ActionRemoved)

	args := []string{"--request", "landed-oob-residue", creation.Path}
	var out, errb strings.Builder
	code := ReleaseCommand(root, home, args, &out, &errb)
	requireTest(t, code == 0, "residue release exit=%d stderr=%q", code, errb.String())
	if _, err := assignmentByID(root, a.ID); err == nil {
		t.Fatal("residue record survived reconcile")
	}
	var replay strings.Builder
	code = ReleaseCommand(root, home, args, &replay, io.Discard)
	requireTest(t, code == 0 && replay.String() == out.String(), "replay exit=%d out=%q", code, replay.String())
}

// TestReleaseNamesRecoveryForPreservedOrphan pins FT93(b), preserved path. Here
// a release's tree was removed out of band but still holds preserved work. It
// returns a verdict handing over the ref itself and leaves the record and its
// recovery pointer intact: release never silently discards preserved work.
func TestReleaseNamesRecoveryForPreservedOrphan(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "oob-preserved")
	a, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	ref := intent.RecoveryRefPrefix(a.OwnerID, a.ID) + "1"
	a.State, a.Recovery = intent.StateRecovered, []intent.Recovery{{Ref: ref, Root: strings.Repeat("a", 40), Payloads: []string{strings.Repeat("b", 40)}}}
	mustNoError(t, intent.PutAssignment(root, a))
	removeOutOfBand(t, root, a, ActionRemoved)

	var out, errb strings.Builder
	code := ReleaseCommand(root, home, []string{"--request", "landed-oob-preserved", creation.Path}, &out, &errb)
	requireTest(t, code != 0, "preserved release exit=%d, want non-zero", code)
	requireTest(t, strings.Contains(errb.String(), "git show "+ref),
		"preserved verdict does not hand over the ref: %q", errb.String())
	got, err := assignmentByID(root, a.ID)
	requireTest(t, err == nil && len(got.Recovery) == 1, "preserved record was mutated or deleted: %v", err)
}

// TestResumeReconcilesTreeGoneRecordsAndSparesYoungActive pins the standing
// cleaner's blast radius over the ledger. A tree-gone record is dropped
// whether it was mid-cleanup or holding preserved work the removed lifecycle
// wrote. A record whose tree still exists survives untouched.
//
// What holds the active, tree-gone record here is its age, not its state. The
// reconcile drops an orphaned active record, but this one survives only
// because it was stamped moments ago and so is not aged. That is the race
// this fixture guards: a reconcile that dropped on tree-absence alone would
// catch a session between `worktree add` and its first write.
func TestResumeReconcilesTreeGoneRecordsAndSparesYoungActive(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	residue := mustCreate(t, root, home, "landed-sweep-residue", "residue")
	preserved := mustCreate(t, root, home, "landed-sweep-preserved", "preserved")
	activeGone := mustCreate(t, root, home, "landed-sweep-active", "active gone")
	live := mustCreate(t, root, home, "landed-sweep-live", "live present")

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

	result, err := conservativeCleanupAt(defaultJoins(), root, home, currentTime())
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
	t.Parallel()
	root, creation, home := newPendingAssignment(t, "release-in-flight")
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
	code := ReleaseCommand(root, home, args, &first, &firstErr)
	requireTest(t, code == 0 && firstErr.String() == "", "in-flight release code=%d stderr=%q", code, firstErr.String())
	var replay strings.Builder
	code = ReleaseCommand(root, home, args, &replay, io.Discard)
	requireTest(t, code == 0 && replay.String() == first.String(), "in-flight replay code=%d stdout=%q", code, replay.String())
	requireTest(t, ReleaseCommand(root, home, []string{"--request", "changed", creation.Path}, io.Discard, io.Discard) != 0, "changed request authorized")
	requireTest(t, ReleaseCommand(root, home, []string{"--request", args[1], root}, io.Discard, io.Discard) != 0, "changed path authorized")
}
func TestExplicitApplyRejectsContentDriftWithoutMutation(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	j := defaultJoins()
	j.ignoredLstat = func(path string) (os.FileInfo, error) {
		if path == ignored {
			return nil, os.ErrNotExist
		}
		return os.Lstat(path)
	}
	plan, err := planExplicitWith(j, root, target, CleanupOptions{DiscardIgnored: true})
	if err != nil || plan.Action != ActionRetain || plan.ReasonCode != ReasonUncertain {
		t.Fatalf("stat-race plan = %#v, %v", plan, err)
	}
	if _, err := os.Lstat(ignored); err != nil {
		t.Fatalf("stat-race plan mutated ignored file: %v", err)
	}
}

func TestLeaseFile(t *testing.T) {
	t.Parallel()
	dir := journeyRepo(t)
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
	t.Parallel()
	out, code := LeaseFileCommand(nil)
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if !strings.HasPrefix(out, "usage:") {
		t.Errorf("out = %q, want usage line", out)
	}
}

// TestCreateCommandWritesBelowTheExplicitHome covers WF15. The verb takes its home from
// the caller, so the bound environment naming a different directory changes nothing. A
// verb that still read the environment would put the pool under the bound home, and the
// second assertion names that directory empty. The bind keeps this test serial.
func TestCreateCommandWritesBelowTheExplicitHome(t *testing.T) {
	root := newWorktreeRepo(t)
	bound, explicit := t.TempDir(), t.TempDir()
	bindEnv(t, homeEnv, bound)
	var stdout, stderr bytes.Buffer
	code := CreateCommand(root, explicit, []string{"--request", "wf15-explicit-home", "--label", "explicit home"}, &stdout, &stderr)
	requireTest(t, code == 0, "create exit = %d, stderr %q", code, stderr.String())
	canonical, err := canonicalPath(root)
	requireTest(t, err == nil, "canonical root: %v", err)
	entries, err := os.ReadDir(poolAt(explicit, canonical))
	requireTest(t, err == nil && len(entries) == 1, "explicit home pool holds %d entries (%v), want one worktree", len(entries), err)
	_, err = os.Stat(poolKeysDirAt(bound))
	requireTest(t, os.IsNotExist(err), "the verb wrote below the bound home %s: %v", bound, err)
}

func TestPoolCommandExplicitRoot(t *testing.T) {
	home := t.TempDir()
	bindEnv(t, "BENCH_HOME", home)
	out, code := PoolCommand(home, []string{"/home/mgibs/workspace/bench"})
	if code != 0 {
		t.Errorf("exit = %d, want 0", code)
	}
	want := filepath.Join(home, "worktrees", "bench-2826441890") + "\n"
	if out != want {
		t.Errorf("out = %q, want %q", out, want)
	}
}

// TestReleaseNamesTheOwnerMarkerAndRetainsTheCheckout pins release's retained clause on a
// bundle component other than the request token: a rewritten owner marker names the marker
// and keeps the checkout.
func TestReleaseNamesTheOwnerMarkerAndRetainsTheCheckout(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "release-owner-marker")
	rewriteMarkerOwner(t, creation.Path, strings.Repeat("a", 32))
	var stdout, stderr strings.Builder
	code := ReleaseCommand(root, home, []string{"--request", "landed-release-owner-marker", creation.Path}, &stdout, &stderr)
	want := "bench worktree release: owner marker does not match assignment " + creation.Assignment.ID + "; checkout retained\n"
	if code != 1 || stdout.String() != "" || stderr.String() != want {
		t.Fatalf("owner-marker release = (%d, %q, %q), want exit 1 and stderr %q", code, stdout.String(), stderr.String(), want)
	}
}

// TestReleaseDropsTheCensusRecords is EC23. The release retires the assignment, so
// its records leave with it; a kept file shows a stale row on every later board.
func TestReleaseDropsTheCensusRecords(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "census-release")
	recordRawCalls(t, home, root, creation.Path, 2)
	survivor := seedHandoffSections(t, root, creation.Assignment)
	var stdout, stderr strings.Builder
	code := ReleaseCommand(root, home, []string{"--request", "landed-census-release", creation.Path}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("release = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(censusRecordPath(home, root, creation.Assignment.ID)); !os.IsNotExist(err) {
		t.Fatalf("the release kept the census record: %v", err)
	}
	requireHandoffSections(t, root, handoffdoc.MainKey, survivor)
}

// TestRetirementLeavesMainInTheDocument is HS20. The last assignment section leaves with
// its retirement, and the document still carries main without a later `bench handoff`.
func TestRetirementLeavesMainInTheDocument(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "handoff-last-section")
	seedOneHandoffSection(t, root, creation.Assignment.Request)
	var stdout, stderr strings.Builder
	code := ReleaseCommand(root, home, []string{"--request", "landed-handoff-last-section", creation.Path}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("release = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	requireHandoffSections(t, root, handoffdoc.MainKey)
}

// TestCleanDropsTheCensusRecords is the clean half of EC24. Release and clean reach
// the one retirement path, so neither leaves a stale record file behind.
func TestCleanDropsTheCensusRecords(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "census-clean")
	recordRawCalls(t, home, root, creation.Path, 2)
	survivor := seedHandoffSections(t, root, creation.Assignment)
	var planned, stderr bytes.Buffer
	if code := CleanCommand(root, home, []string{creation.Path}, &planned, &stderr); code != 0 {
		t.Fatalf("clean plan = (%d, %q, %q)", code, planned.String(), stderr.String())
	}
	fingerprint := regexp.MustCompile(`[0-9a-f]{64}`).FindString(planned.String())
	if fingerprint == "" {
		t.Fatalf("clean plan carried no fingerprint: %s", planned.String())
	}
	var applied bytes.Buffer
	if code := CleanCommand(root, home, []string{creation.Path, "--apply", fingerprint}, &applied, &stderr); code != 0 || !strings.Contains(applied.String(), ",removed,") {
		t.Fatalf("clean apply = (%d, %q, %q)", code, applied.String(), stderr.String())
	}
	if _, err := os.Stat(censusRecordPath(home, root, creation.Assignment.ID)); !os.IsNotExist(err) {
		t.Fatalf("the clean kept the census record: %v", err)
	}
	requireHandoffSections(t, root, handoffdoc.MainKey, survivor)
}

// seedOneHandoffSection writes one section under key into the document the retirement
// path resolves for root. The document is excluded first, because the repositories that
// carry one ignore it, and a tracked copy would refuse the landing as a dirty destination.
func seedOneHandoffSection(t *testing.T, root, key string) {
	t.Helper()
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte(handoffFile+"\n"+handoffFile+".lock\n"), 0o644)
	section := handoffdoc.Section{Key: key, Next: "bench status", Fields: []handoffdoc.Field{{Label: handoffdoc.LabelLabel, Value: key}}}
	path := handoffDocumentPath(root)
	if err := handoffdoc.WriteSection(path, section); err != nil {
		t.Fatalf("seed handoff section %s: %v", key, err)
	}
	// The write leaves its lock file behind. The seed drops it, so the landing's residue
	// scan grades the checkout the way a repository that has never written one does.
	if err := os.Remove(handoffdoc.LockPath(path)); err != nil {
		t.Fatalf("drop the seed lock: %v", err)
	}
}

// seedHandoffSections gives the document the retired assignment's section and one other
// assignment's, and returns the other's key. A one-section document cannot tell a removal
// that drops the right section from one that empties the file.
func seedHandoffSections(t *testing.T, root string, assignment intent.Assignment) string {
	t.Helper()
	survivor := intent.RequestDigest("handoff-survivor-" + assignment.ID)
	seedOneHandoffSection(t, root, assignment.Request)
	seedOneHandoffSection(t, root, survivor)
	return survivor
}

// requireHandoffSections fails unless the document holds exactly these section keys, in
// this order.
func requireHandoffSections(t *testing.T, root string, want ...string) {
	t.Helper()
	doc, err := handoffdoc.Read(handoffDocumentPath(root))
	if err != nil {
		t.Fatalf("read handoff document: %v", err)
	}
	var got []string
	for _, section := range doc.Sections {
		got = append(got, section.Key)
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("handoff sections = %v, want %v", got, want)
	}
}

// TestCensusDropHasOneCallSiteInThisPackage pins the one drop owner. A second call
// site is a second retirement rule, and the two drift.
func TestCensusDropHasOneCallSiteInThisPackage(t *testing.T) {
	t.Parallel()
	requireOneCallSite(t, "census.Drop(", "lifecycle.go")
}

// TestHandoffSectionRemovalHasOneCallSiteInThisPackage is HS19. The section removal rides
// the retirement path the census drop rides, so it is pinned the same way: a second site
// is a second retirement rule.
func TestHandoffSectionRemovalHasOneCallSiteInThisPackage(t *testing.T) {
	t.Parallel()
	requireOneCallSite(t, "handoffdoc.RemoveSection(", "lifecycle.go")
}

// requireOneCallSite fails unless needle appears once in the package's production files,
// in the named file.
func requireOneCallSite(t *testing.T, needle, owner string) {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	sites := map[string]int{}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		body, readErr := os.ReadFile(name)
		if readErr != nil {
			t.Fatal(readErr)
		}
		if count := strings.Count(string(body), needle); count > 0 {
			sites[name] = count
		}
	}
	if len(sites) != 1 || sites[owner] != 1 {
		t.Fatalf("%s call sites = %v, want one in %s", strings.TrimSuffix(needle, "("), sites, owner)
	}
}

// --- create --from: the sibling start ---

// createFromGrammar is the exact create usage line. The help rows below read this literal
// rather than the const the command prints, so a grammar that loses `[--from <target>]`
// turns them red instead of following the change.
const createFromGrammar = "usage: bench worktree create [--refresh] --request <opaque-id> --label <work-item> [--from <target>]"

// runCreate drives CreateCommand against a fixture repository and returns the exit code
// with both streams.
func runCreate(t *testing.T, root, home string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := CreateCommand(root, home, args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// requireCreateFromRefusal pins the surface every `--from` refusal shares: exit 1, an empty
// stdout, and the shared target printer's two stderr lines.
func requireCreateFromRefusal(t *testing.T, code int, stdout, stderr string, fragments ...string) {
	t.Helper()
	if code != 1 || stdout != "" {
		t.Fatalf("create --from = (%d, %q, %q), want (1, empty stdout, a refusal)", code, stdout, stderr)
	}
	for _, fragment := range fragments {
		if !strings.Contains(stderr, fragment) {
			t.Fatalf("create --from stderr = %q, want it to hold %q", stderr, fragment)
		}
	}
}

// assignmentCount reads how many records the ledger holds, so a refusal row can prove the
// verb registered nothing.
func assignmentCount(t *testing.T, root string) int {
	t.Helper()
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	return len(assignments)
}

// identityComponentDetail is the detail sentence inside a fixture's whole refused record.
// The shared target printer prints that sentence alone, so a row reads the registry's text
// through the same fixture the merge rows read.
func identityComponentDetail(t *testing.T, fixture identityComponentFixture, creation Creation) string {
	t.Helper()
	field, _, _ := strings.Cut(fixture.want(creation, "", ""), ",")
	detail, ok := strings.CutPrefix(field, "detail=")
	if !ok {
		t.Fatalf("fixture record %q names no detail field", field)
	}
	return detail
}

// WF23: `--from <sibling>` starts the new assignment at the sibling's committed tip. The
// ledger start, the checkout HEAD, and the new assignment's own branch all read that tip,
// so the sibling's branch stays untouched.
func TestCreateFromStartsAtTheSiblingTip(t *testing.T) {
	t.Parallel()
	_, root, home, _, created := mergeFixture(t, "delegate")
	sibling := created[0]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	tip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")

	const request = "create-from-tip"
	code, stdout, stderr := runCreate(t, root, home,
		"--request", request, "--label", "dependent", "--from", sibling.Assignment.Label)
	if code != 0 {
		t.Fatalf("create --from = (%d, %q, %q), want 0", code, stdout, stderr)
	}
	record, ok, err := intent.FindAssignmentForRequest(root, request)
	mustNoError(t, err)
	if !ok {
		t.Fatalf("request %q registered no assignment", request)
	}
	if record.Start != tip {
		t.Errorf("ledger start = %s, want the sibling tip %s", record.Start, tip)
	}
	if head := gitOutput(t, record.Worktree, "rev-parse", "HEAD"); head != tip {
		t.Errorf("new worktree HEAD = %s, want the sibling tip %s", head, tip)
	}
	want := intent.AssignmentBranchRef(record.OwnerID, record.ID)
	if branch := gitOutput(t, record.Worktree, "rev-parse", "--symbolic-full-name", "HEAD"); branch != want {
		t.Errorf("new worktree branch = %s, want its own assignment branch %s", branch, want)
	}
	if head := gitOutput(t, sibling.Path, "rev-parse", "HEAD"); head != tip {
		t.Errorf("sibling HEAD = %s, want it unmoved at %s", head, tip)
	}
}

// WF45: a replay of `create --from <sibling>` with the same request returns the existing
// record whatever the sibling's checkout now holds. The request lookup is the first fact
// the creation reads, so an edit the sibling took after the first run refuses nothing and
// the ledger gains no second record.
func TestCreateFromReplayReturnsTheRecord(t *testing.T) {
	t.Parallel()
	_, root, home, _, created := mergeFixture(t, "delegate")
	sibling := created[0]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")

	const request = "create-from-replay"
	code, first, stderr := runCreate(t, root, home,
		"--request", request, "--label", "dependent", "--from", sibling.Assignment.Label)
	if code != 0 {
		t.Fatalf("create --from = (%d, %q, %q), want 0", code, first, stderr)
	}
	mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)

	code, second, stderr := runCreate(t, root, home,
		"--request", request, "--label", "dependent", "--from", sibling.Assignment.Label)
	if code != 0 {
		t.Fatalf("create --from replay = (%d, %q, %q), want 0", code, second, stderr)
	}
	if second != first {
		t.Errorf("replay stdout = %q, want the first run's %q", second, first)
	}
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	held := 0
	for _, a := range assignments {
		if a.RequestToken == request {
			held++
		}
	}
	if held != 1 {
		t.Errorf("ledger holds %d records for request %q, want 1", held, request)
	}
}

// WF24: a `--from` that names no active assignment refuses through the shared printer and
// registers nothing. The flag composes no commit lookup, so a typo never falls through to
// the default tip.
func TestCreateFromRefusesAnUnknownSibling(t *testing.T) {
	t.Parallel()
	_, root, home, _, _ := mergeFixture(t, "delegate")
	before := assignmentCount(t, root)

	code, stdout, stderr := runCreate(t, root, home,
		"--request", "create-from-unknown", "--label", "dependent", "--from", "no-such-label")
	requireCreateFromRefusal(t, code, stdout, stderr,
		"bench worktree create: --from names no active assignment\n", "next=bench worktree list\n")
	if after := assignmentCount(t, root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// WF25 and WF26: a sibling contributes its committed branch tip alone, so a dirty sibling
// names `bench commit` at the sibling and a detached sibling names its assignment branch.
func TestCreateFromRefusesADirtyOrDetachedSibling(t *testing.T) {
	t.Parallel()
	_, root, home, _, created := mergeFixture(t, "delegate")
	sibling := created[0]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	before := assignmentCount(t, root)
	mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)

	code, stdout, stderr := runCreate(t, root, home,
		"--request", "create-from-dirty", "--label", "dependent", "--from", sibling.Assignment.Label)
	requireCreateFromRefusal(t, code, stdout, stderr,
		"bench worktree create: sibling checkout is not clean\n",
		"next=bench worktree exec "+sibling.Assignment.ID+" -- bench commit\n")

	gitRun(t, sibling.Path, "checkout", "-q", "--", "sibling.txt")
	gitRun(t, sibling.Path, "checkout", "-q", "--detach", "HEAD")
	code, stdout, stderr = runCreate(t, root, home,
		"--request", "create-from-detached", "--label", "dependent", "--from", sibling.Assignment.Label)
	requireCreateFromRefusal(t, code, stdout, stderr,
		"bench worktree create: sibling is not on its assignment branch\n", "next=bench worktree list\n")
	if after := assignmentCount(t, root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusals", after, before)
	}
}

// WF27: `--from` and `--refresh` name two starts, so the pair is invalid usage. The
// refusal runs before refreshop.Consume, and an empty stdout is the evidence: a refresh
// that ran would have written its own table there.
func TestCreateFromWithRefreshRefusesBeforeTheRefresh(t *testing.T) {
	t.Parallel()
	_, root, home, _, created := mergeFixture(t, "delegate")
	sibling := created[0]
	before := assignmentCount(t, root)

	code, stdout, stderr := runCreate(t, root, home, "--request", "create-from-refresh",
		"--label", "dependent", "--refresh", "--from", sibling.Assignment.Label)
	want := toon.Usage(createGrammar.Cmd, "--from with --refresh") + "\n"
	if code != 2 || stdout != "" || stderr != want {
		t.Fatalf("create --refresh --from = (%d, %q, %q), want (2, empty, %q)", code, stdout, stderr, want)
	}
	if after := assignmentCount(t, root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// WF28: two siblings whose labels share the prefix make it ambiguous, and the refusal names
// both ids. A first-match lookup would start the new worktree at the wrong sibling.
func TestCreateFromRefusesAnAmbiguousPrefix(t *testing.T) {
	t.Parallel()
	_, root, home, _, created := mergeFixture(t, "delegate-alpha", "delegate-beta")
	before := assignmentCount(t, root)

	code, stdout, stderr := runCreate(t, root, home,
		"--request", "create-from-ambiguous", "--label", "dependent", "--from", "delegate-")
	requireCreateFromRefusal(t, code, stdout, stderr,
		"bench worktree create: target is ambiguous: ",
		created[0].Assignment.ID, created[1].Assignment.ID, "next=bench worktree list\n")
	if after := assignmentCount(t, root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// WF29: a `--from` with a control byte refuses before the ledger read. The malformed
// ledger file is the proof: a lookup that ran first would report the unreadable file
// instead, and the file's bytes stay as written.
func TestCreateFromRefusesControlBytes(t *testing.T) {
	t.Parallel()
	_, root, home, _, _ := mergeFixture(t, "delegate")
	address, err := intent.Address(root)
	mustNoError(t, err)
	malformed := []byte("{ this is not a ledger\n")
	mustWrite(t, address, malformed, 0o600)

	code, stdout, stderr := runCreate(t, root, home,
		"--request", "create-from-control", "--label", "dependent", "--from", "a\x01b")
	requireCreateFromRefusal(t, code, stdout, stderr,
		"bench worktree create: --from contains control characters\n", "next=bench worktree list\n")
	after, err := os.ReadFile(address)
	mustNoError(t, err)
	if string(after) != string(malformed) {
		t.Fatalf("ledger bytes = %q, want them unread and unwritten at %q", after, malformed)
	}
}

// WF30: a sibling whose state is no longer active authenticates nothing, so its label names
// no sibling. A lookup over every state would start the new worktree at a retired tip.
func TestCreateFromRefusesARetiredSibling(t *testing.T) {
	t.Parallel()
	_, root, home, _, created := mergeFixture(t, "delegate")
	sibling := created[0]
	retired := sibling.Assignment
	retired.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(root, retired))
	before := assignmentCount(t, root)

	code, stdout, stderr := runCreate(t, root, home,
		"--request", "create-from-retired", "--label", "dependent", "--from", sibling.Assignment.Label)
	requireCreateFromRefusal(t, code, stdout, stderr,
		"bench worktree create: --from names no active assignment\n", "next=bench worktree list\n")
	if after := assignmentCount(t, root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// WF43: the sibling's creation bundle is the flag's whole authority, so a broken component
// names itself and the verb registers nothing. The assignment-state component is the
// exception this package can reach: the lookup narrows to the active assignments before
// the bundle runs, so a non-active sibling meets WF30's sentence instead.
func TestCreateFromRefusesAFailedSiblingIdentityComponent(t *testing.T) {
	t.Parallel()
	for _, component := range []string{componentOwnerMarker, componentLock} {
		t.Run(component, func(t *testing.T) {
			t.Parallel()
			_, root, home, _, created := mergeFixture(t, "delegate")
			sibling := created[0]
			before := assignmentCount(t, root)
			fixture := identityComponentFixtureFor(t, component)
			fixture.mutate(t, root, sibling)

			code, stdout, stderr := runCreate(t, root, home, "--request", "create-from-"+component,
				"--label", "dependent", "--from", sibling.Assignment.Label)
			requireCreateFromRefusal(t, code, stdout, stderr,
				"bench worktree create: "+identityComponentDetail(t, fixture, sibling)+"\n",
				"next=bench worktree list\n")
			if after := assignmentCount(t, root); after != before {
				t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
			}
		})
	}
	t.Run(componentAssignmentState, func(t *testing.T) {
		t.Parallel()
		_, root, home, _, created := mergeFixture(t, "delegate")
		sibling := created[0]
		before := assignmentCount(t, root)
		identityComponentFixtureFor(t, componentAssignmentState).mutate(t, root, sibling)

		code, stdout, stderr := runCreate(t, root, home, "--request", "create-from-state",
			"--label", "dependent", "--from", sibling.Assignment.Label)
		requireCreateFromRefusal(t, code, stdout, stderr,
			"bench worktree create: --from names no active assignment\n", "next=bench worktree list\n")
		if after := assignmentCount(t, root); after != before {
			t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
		}
	})
}
