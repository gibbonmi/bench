package worktree

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
	"os"
	"path/filepath"
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
		f := newOwnedAssignment(t, "terminal-clean")
		sibling := mustCreate(t, f.root, f.home, "terminal-clean-sibling", "sibling")
		siblingRef := "refs/bench/recovery/" + sibling.Assignment.OwnerID + "/" + sibling.Assignment.ID + "/1"
		gitRun(t, f.root, "update-ref", siblingRef, f.creation.Assignment.Start)
		markPending(t, f.root, f.creation.Assignment)
		if _, err := ApplyAutomatic(f.root, f.creation.Path, nil); err != nil {
			t.Fatal(err)
		}
		if descendant(t, "git", "-C", f.root, "show-ref", "--verify", "--quiet", f.creation.Assignment.Branch).Run() == nil {
			t.Fatal("exact cleanup left target branch")
		}
		gitRun(t, f.root, "show-ref", "--verify", "--quiet", sibling.Assignment.Branch)
		gitRun(t, f.root, "show-ref", "--verify", "--quiet", siblingRef)
		assignments, err := intent.Assignments(f.root)
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
	created := runVerb(t, verbCreate, repoHome{root, home}.call("--request", "next-hint", "--label", "next hint label"))
	if created.exit != 0 {
		t.Fatalf("CreateCommand exit = %d, stderr = %q", created.exit, created.stderr)
	}
	got := created.stdout
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
		r := runVerb(t, verbCreate, verbCall{home: Home(), args: args})
		if r.exit != 0 || r.stdout != want || r.stderr != "" {
			t.Fatalf("create %q = (%d, %q, %q), want (0, %q, empty)", args, r.exit, r.stdout, r.stderr, want)
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
			r := runVerb(t, verbRelease, verbCall{home: Home(), args: tc.args})
			if r.exit != tc.wantCode || r.stdout != tc.wantStdout || r.stderr != tc.wantStderr {
				t.Fatalf("release %q = (%d, %q, %q), want (%d, %q, %q)", tc.args, r.exit, r.stdout, r.stderr, tc.wantCode, tc.wantStdout, tc.wantStderr)
			}
		})
	}
}

// TestCreateCommandHelpPerformsNoRefresh pins that a help request never fetches:
// usage.Parse answers --help before refreshop.Consume ever sees the args, so a
// --refresh alongside --help prints only the help line, not a worktree_refresh table.
func TestCreateCommandHelpPerformsNoRefresh(t *testing.T) {
	t.Parallel()
	r := runVerb(t, verbCreate, verbCall{home: Home(), args: []string{"--request", "x", "--refresh", "--help"}})
	want := "usage: " + usage.WorktreeCreate + "\n"
	if r.exit != 0 || r.stdout != want || r.stderr != "" {
		t.Fatalf("CreateCommand with --refresh --help = (%d, %q, %q), want (0, %q, empty)", r.exit, r.stdout, r.stderr, want)
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
		if r := runVerb(t, verbCreate, verbCall{home: Home(), args: args}); r.exit != 2 || r.stdout != "" || r.stderr != createGrammar.Help+"\n" {
			t.Fatalf("create %q = (%d, %q, %q), want (2, empty, %q)", args, r.exit, r.stdout, r.stderr, createGrammar.Help+"\n")
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
		if r := runVerb(t, verbCreate, verbCall{home: Home(), args: args}); r.exit != 2 || r.stdout != "" {
			t.Fatalf("create %q = (%d, %q, %q), want exit 2 with empty stdout", args, r.exit, r.stdout, r.stderr)
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

	call := repoHome{root, home}.call("--request", "retain-verdict", creation.Path)
	retained := runVerb(t, verbRelease, call)
	msg := retained.stderr
	requireTest(t, retained.exit != 0, "retained release exit = %d, want non-zero", retained.exit)
	requireTest(t, !strings.Contains(msg, "terminal receipt missing"), "masking error still present: %q", msg)
	requireTest(t, strings.Contains(msg, "retained") && strings.Contains(msg, "bench worktree release"),
		"retained verdict is not actionable: %q", msg)

	mustNoError(t, os.Remove(residual))
	recovered := runVerb(t, verbRelease, call)
	requireTest(t, recovered.exit == 0, "recovery release exit = %d, want 0; out=%q", recovered.exit, recovered.stdout)
}

// TestReleaseUnknownRequestNamesReauthorizeRecovery is LR19: release names the request
// component, its own retained clause, and the same recovery command the landing names.
func TestReleaseUnknownRequestNamesReauthorizeRecovery(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "release-reauthorize-recovery")
	r := runVerb(t, verbRelease, f.call("--request", "unknown-request", f.creation.Path))
	wantNext := "bench worktree reauthorize --assignment " + f.creation.Assignment.ID + " --request <new-request> --base <full-base-commit> --source-tip <full-source-tip-commit> '" + f.creation.Path + "'"
	want := "bench worktree release: request token matches no assignment; checkout retained; observed=assignment:" + f.creation.Assignment.ID + ",next=" + wantNext + "\n"
	if r.exit != 1 || r.stdout != "" || r.stderr != want {
		t.Fatalf("unknown-request release = (%d, %q, %q), want exit 1 and stderr %q", r.exit, r.stdout, r.stderr, want)
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
	f := newOwnedAssignment(t, "oob-residue")
	a, err := assignmentByID(f.root, f.creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, len(a.Recovery) == 0, "fixture already holds recovery metadata")
	removeOutOfBand(t, f.root, a, ActionRemoved)

	call := f.call("--request", "landed-oob-residue", f.creation.Path)
	first := runVerb(t, verbRelease, call)
	requireTest(t, first.exit == 0, "residue release exit=%d stderr=%q", first.exit, first.stderr)
	if _, err := assignmentByID(f.root, a.ID); err == nil {
		t.Fatal("residue record survived reconcile")
	}
	replay := runVerb(t, verbRelease, call)
	requireTest(t, replay.exit == 0 && replay.stdout == first.stdout, "replay exit=%d out=%q", replay.exit, replay.stdout)
}

// TestReleaseNamesRecoveryForPreservedOrphan pins FT93(b), preserved path. Here
// a release's tree was removed out of band but still holds preserved work. It
// returns a verdict handing over the ref itself and leaves the record and its
// recovery pointer intact: release never silently discards preserved work.
func TestReleaseNamesRecoveryForPreservedOrphan(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "oob-preserved")
	a, err := assignmentByID(f.root, f.creation.Assignment.ID)
	mustNoError(t, err)
	ref := intent.RecoveryRefPrefix(a.OwnerID, a.ID) + "1"
	a.State, a.Recovery = intent.StateRecovered, []intent.Recovery{{Ref: ref, Root: strings.Repeat("a", 40), Payloads: []string{strings.Repeat("b", 40)}}}
	mustNoError(t, intent.PutAssignment(f.root, a))
	removeOutOfBand(t, f.root, a, ActionRemoved)

	r := runVerb(t, verbRelease, f.call("--request", "landed-oob-preserved", f.creation.Path))
	requireTest(t, r.exit != 0, "preserved release exit=%d, want non-zero", r.exit)
	requireTest(t, strings.Contains(r.stderr, "git show "+ref),
		"preserved verdict does not hand over the ref: %q", r.stderr)
	got, err := assignmentByID(f.root, a.ID)
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

	result := mustSweep(t, root, home)
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
	f := newPendingAssignment(t, "release-in-flight")
	stop := errors.New("crash after removal")
	_, err := ApplyAutomatic(f.root, f.creation.Path, func(step LifecycleStep) error {
		if step == StepRemoval {
			return stop
		}
		return nil
	})
	requireTest(t, errors.Is(err, stop), "automatic interruption = %v", err)
	call := f.call("--request", "landed-release-in-flight", f.creation.Path)
	first := runVerb(t, verbRelease, call)
	requireTest(t, first.exit == 0 && first.stderr == "", "in-flight release code=%d stderr=%q", first.exit, first.stderr)
	replay := runVerb(t, verbRelease, call)
	requireTest(t, replay.exit == 0 && replay.stdout == first.stdout, "in-flight replay code=%d stdout=%q", replay.exit, replay.stdout)
	requireTest(t, runVerb(t, verbRelease, f.call("--request", "changed", f.creation.Path)).exit != 0, "changed request authorized")
	requireTest(t, runVerb(t, verbRelease, f.call("--request", call.args[1], f.root)).exit != 0, "changed path authorized")
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

// TestIgnoredInventoryStatRaceRetains denies search on an ignored file's directory, so Git lists
// the name and os.Lstat fails. The nested-state read fails too; only the reason names the stat fault.
func TestIgnoredInventoryStatRaceRetains(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot deny search access to fail the stat")
	}
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "main")
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("locked/\n"), 0o644)
	target := filepath.Join(filepath.Dir(root), "ignored stat race")
	gitRun(t, root, "worktree", "add", "-q", "-b", "ignored-stat-race", target, "HEAD")
	locked := filepath.Join(target, "locked")
	ignored := filepath.Join(locked, "ignored.txt")
	mustMkdirAll(t, locked, 0o755)
	mustWrite(t, ignored, []byte("secret\n"), 0o644)
	t.Cleanup(func() { _ = os.Chmod(locked, 0o700) })
	mustNoError(t, os.Chmod(locked, 0o600))
	plan, err := PlanExplicitWithOptions(root, target, CleanupOptions{DiscardIgnored: true})
	statFault := lifecyclepolicy.DecideExplicit(lifecyclepolicy.ExplicitFacts{IgnoredErr: os.ErrPermission}).Reason
	requireTest(t, err == nil && plan.Action == ActionRetain && plan.ReasonCode == ReasonUncertain && plan.Reason == statFault,
		"stat-race plan = %#v, %v; want the reason %q", plan, err, statFault)
	mustNoError(t, os.Chmod(locked, 0o700))
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
	r := runVerb(t, verbLeaseFile, verbCall{})
	if r.exit != 2 {
		t.Errorf("exit = %d, want 2", r.exit)
	}
	if !strings.HasPrefix(r.stdout, "usage:") {
		t.Errorf("out = %q, want usage line", r.stdout)
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
	r := runVerb(t, verbCreate, repoHome{root, explicit}.call("--request", "wf15-explicit-home", "--label", "explicit home"))
	requireTest(t, r.exit == 0, "create exit = %d, stderr %q", r.exit, r.stderr)
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
	r := runVerb(t, verbPool, verbCall{home: home, args: []string{"/home/mgibs/workspace/bench"}})
	if r.exit != 0 {
		t.Errorf("exit = %d, want 0", r.exit)
	}
	want := filepath.Join(home, "worktrees", "bench-2826441890") + "\n"
	if r.stdout != want {
		t.Errorf("out = %q, want %q", r.stdout, want)
	}
}

// TestReleaseNamesTheOwnerMarkerAndRetainsTheCheckout pins release's retained clause on a
// bundle component other than the request token: a rewritten owner marker names the marker
// and keeps the checkout.
func TestReleaseNamesTheOwnerMarkerAndRetainsTheCheckout(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "release-owner-marker")
	rewriteMarkerOwner(t, f.creation.Path, strings.Repeat("a", 32))
	r := runVerb(t, verbRelease, f.call("--request", "landed-release-owner-marker", f.creation.Path))
	want := "bench worktree release: owner marker does not match assignment " + f.creation.Assignment.ID + "; checkout retained\n"
	if r.exit != 1 || r.stdout != "" || r.stderr != want {
		t.Fatalf("owner-marker release = (%d, %q, %q), want exit 1 and stderr %q", r.exit, r.stdout, r.stderr, want)
	}
}

// TestReleaseDropsTheCensusRecords is EC23. The release retires the assignment, so
// its records leave with it; a kept file shows a stale row on every later board.
func TestReleaseDropsTheCensusRecords(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "census-release")
	recordRawCalls(t, f.home, f.root, f.creation.Path, 2)
	survivor := seedHandoffSections(t, f.root, f.creation.Assignment)
	if r := runVerb(t, verbRelease, f.call("--request", "landed-census-release", f.creation.Path)); r.exit != 0 {
		t.Fatalf("release = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(censusRecordPath(f.home, f.root, f.creation.Assignment.ID)); !os.IsNotExist(err) {
		t.Fatalf("the release kept the census record: %v", err)
	}
	requireHandoffSections(t, f.root, handoffdoc.MainKey, survivor)
}

// TestRetirementLeavesMainInTheDocument is HS20. The last assignment section leaves with
// its retirement, and the document still carries main without a later `bench handoff`.
func TestRetirementLeavesMainInTheDocument(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "handoff-last-section")
	seedOneHandoffSection(t, f.root, f.creation.Assignment.Request)
	if r := runVerb(t, verbRelease, f.call("--request", "landed-handoff-last-section", f.creation.Path)); r.exit != 0 {
		t.Fatalf("release = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	requireHandoffSections(t, f.root, handoffdoc.MainKey)
}

// TestRetirementPrintsTheSectionRemovalError is HS30. A document the parser refuses
// keeps its section, so the retirement says which document and which line refused it.
// The verdict is the retirement's own, because the removal is advisory.
func TestRetirementPrintsTheSectionRemovalError(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "handoff-unparseable")
	seedOneHandoffSection(t, f.root, f.creation.Assignment.Request)
	path := handoffDocumentPath(f.root)
	seeded, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read seeded handoff document: %v", err)
	}
	// The heading is one the grammar has no key for, so Parse refuses the document and
	// the removal never reaches the section.
	document := string(seeded) + "\n## Foo\n\n"
	mustWrite(t, path, []byte(document), 0o644)
	refused := 0
	for i, line := range strings.Split(document, "\n") {
		if line == "## Foo" {
			refused = i + 1
		}
	}
	r := runVerb(t, verbRelease, f.call("--request", "landed-handoff-unparseable", f.creation.Path))
	if r.exit != 0 {
		t.Fatalf("release = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if !strings.Contains(r.stdout, string(ActionRemoved)) {
		t.Fatalf("release verdict = %q, want %s", r.stdout, ActionRemoved)
	}
	if want := fmt.Sprintf("%s:%d:", path, refused); !strings.Contains(r.stderr, want) {
		t.Fatalf("release stderr = %q, want the file and line %q", r.stderr, want)
	}
}

// TestCleanDropsTheCensusRecords is the clean half of EC24. Release and clean reach
// the one retirement path, so neither leaves a stale record file behind.
func TestCleanDropsTheCensusRecords(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "census-clean")
	recordRawCalls(t, f.home, f.root, f.creation.Path, 2)
	survivor := seedHandoffSections(t, f.root, f.creation.Assignment)
	planned := runVerb(t, verbClean, f.call(f.creation.Path))
	if planned.exit != 0 {
		t.Fatalf("clean plan = (%d, %q, %q)", planned.exit, planned.stdout, planned.stderr)
	}
	if applied := runVerb(t, verbClean, f.call(f.creation.Path, "--apply", planned.mustFingerprint(t))); applied.exit != 0 || !strings.Contains(applied.stdout, ",removed,") {
		t.Fatalf("clean apply = (%d, %q, %q)", applied.exit, applied.stdout, applied.stderr)
	}
	if _, err := os.Stat(censusRecordPath(f.home, f.root, f.creation.Assignment.ID)); !os.IsNotExist(err) {
		t.Fatalf("the clean kept the census record: %v", err)
	}
	requireHandoffSections(t, f.root, handoffdoc.MainKey, survivor)
}

// seedOneHandoffSection writes one section under key into the document the retirement
// path resolves for root. The document is excluded first, because the repositories that
// carry one ignore it, and a tracked copy would refuse the landing as a dirty destination.
func seedOneHandoffSection(t *testing.T, root, key string) {
	t.Helper()
	// Only the document is excluded. Update unlinks its lock on release, so no lock file
	// is left for the landing to trip over.
	mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte(handoffdoc.DocumentPath+"\n"), 0o644)
	section := handoffdoc.Section{Key: key, Next: "bench status", Fields: []handoffdoc.Field{{Label: handoffdoc.LabelLabel, Value: key}}}
	path := handoffDocumentPath(root)
	if err := handoffdoc.WriteSection(path, section); err != nil {
		t.Fatalf("seed handoff section %s: %v", key, err)
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

// requireCreateFromRefusal pins the surface every `--from` refusal shares: exit 1, an empty
// stdout, and the shared target printer's two stderr lines.
func requireCreateFromRefusal(t *testing.T, r verbResult, fragments ...string) {
	t.Helper()
	if r.exit != 1 || r.stdout != "" {
		t.Fatalf("create --from = (%d, %q, %q), want (1, empty stdout, a refusal)", r.exit, r.stdout, r.stderr)
	}
	for _, fragment := range fragments {
		if !strings.Contains(r.stderr, fragment) {
			t.Fatalf("create --from stderr = %q, want it to hold %q", r.stderr, fragment)
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
	f := mergeFixture(t, "delegate")
	sibling := f.created[0]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	tip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")

	const request = "create-from-tip"
	r := runVerb(t, verbCreate, f.call(
		"--request", request, "--label", "dependent", "--from", sibling.Assignment.Label))
	if r.exit != 0 {
		t.Fatalf("create --from = (%d, %q, %q), want 0", r.exit, r.stdout, r.stderr)
	}
	record, ok, err := intent.FindAssignmentForRequest(f.root, request)
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
	f := mergeFixture(t, "delegate")
	sibling := f.created[0]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")

	const request = "create-from-replay"
	call := f.call("--request", request, "--label", "dependent", "--from", sibling.Assignment.Label)
	first := runVerb(t, verbCreate, call)
	if first.exit != 0 {
		t.Fatalf("create --from = (%d, %q, %q), want 0", first.exit, first.stdout, first.stderr)
	}
	mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)

	second := runVerb(t, verbCreate, call)
	if second.exit != 0 {
		t.Fatalf("create --from replay = (%d, %q, %q), want 0", second.exit, second.stdout, second.stderr)
	}
	if second.stdout != first.stdout {
		t.Errorf("replay stdout = %q, want the first run's %q", second.stdout, first.stdout)
	}
	assignments, err := intent.Assignments(f.root)
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
	f := mergeFixture(t, "delegate")
	before := assignmentCount(t, f.root)

	requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call(
		"--request", "create-from-unknown", "--label", "dependent", "--from", "no-such-label")),
		"bench worktree create: --from names no active assignment\n", "next=bench worktree list\n")
	if after := assignmentCount(t, f.root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// WF25 and WF26: a sibling contributes its committed branch tip alone, so a dirty sibling
// names `bench commit` at the sibling and a detached sibling names its assignment branch.
func TestCreateFromRefusesADirtyOrDetachedSibling(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "delegate")
	sibling := f.created[0]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	before := assignmentCount(t, f.root)
	mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)

	requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call(
		"--request", "create-from-dirty", "--label", "dependent", "--from", sibling.Assignment.Label)),
		"bench worktree create: sibling checkout is not clean\n",
		"next=bench worktree exec "+sibling.Assignment.ID+" -- bench commit\n")

	gitRun(t, sibling.Path, "checkout", "-q", "--", "sibling.txt")
	gitRun(t, sibling.Path, "checkout", "-q", "--detach", "HEAD")
	requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call(
		"--request", "create-from-detached", "--label", "dependent", "--from", sibling.Assignment.Label)),
		"bench worktree create: sibling is not on its assignment branch\n", "next=bench worktree list\n")
	if after := assignmentCount(t, f.root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusals", after, before)
	}
}

// WF27: `--from` and `--refresh` name two starts, so the pair is invalid usage. The
// refusal runs before refreshop.Consume, and an empty stdout is the evidence: a refresh
// that ran would have written its own table there.
func TestCreateFromWithRefreshRefusesBeforeTheRefresh(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "delegate")
	sibling := f.created[0]
	before := assignmentCount(t, f.root)

	r := runVerb(t, verbCreate, f.call("--request", "create-from-refresh",
		"--label", "dependent", "--refresh", "--from", sibling.Assignment.Label))
	want := toon.Usage(createGrammar.Cmd, "--from with --refresh") + "\n"
	if r.exit != 2 || r.stdout != "" || r.stderr != want {
		t.Fatalf("create --refresh --from = (%d, %q, %q), want (2, empty, %q)", r.exit, r.stdout, r.stderr, want)
	}
	if after := assignmentCount(t, f.root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// WF28: two siblings whose labels share the prefix make it ambiguous, and the refusal names
// both ids. A first-match lookup would start the new worktree at the wrong sibling.
func TestCreateFromRefusesAnAmbiguousPrefix(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "delegate-alpha", "delegate-beta")
	before := assignmentCount(t, f.root)

	requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call(
		"--request", "create-from-ambiguous", "--label", "dependent", "--from", "delegate-")),
		"bench worktree create: target is ambiguous: ",
		f.created[0].Assignment.ID, f.created[1].Assignment.ID, "next=bench worktree list\n")
	if after := assignmentCount(t, f.root); after != before {
		t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
	}
}

// WF29: a `--from` with a control byte refuses before the ledger read. The malformed
// ledger file is the proof: a lookup that ran first would report the unreadable file
// instead, and the file's bytes stay as written.
func TestCreateFromRefusesControlBytes(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "delegate")
	address, err := intent.Address(f.root)
	mustNoError(t, err)
	malformed := []byte("{ this is not a ledger\n")
	mustWrite(t, address, malformed, 0o600)

	requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call(
		"--request", "create-from-control", "--label", "dependent", "--from", "a\x01b")),
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
	f := mergeFixture(t, "delegate")
	sibling := f.created[0]
	retired := sibling.Assignment
	retired.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(f.root, retired))
	before := assignmentCount(t, f.root)

	requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call(
		"--request", "create-from-retired", "--label", "dependent", "--from", sibling.Assignment.Label)),
		"bench worktree create: --from names no active assignment\n", "next=bench worktree list\n")
	if after := assignmentCount(t, f.root); after != before {
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
			f := mergeFixture(t, "delegate")
			sibling := f.created[0]
			before := assignmentCount(t, f.root)
			fixture := identityComponentFixtureFor(t, component)
			fixture.mutate(t, f.root, sibling)

			requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call("--request", "create-from-"+component,
				"--label", "dependent", "--from", sibling.Assignment.Label)),
				"bench worktree create: "+identityComponentDetail(t, fixture, sibling)+"\n",
				"next=bench worktree list\n")
			if after := assignmentCount(t, f.root); after != before {
				t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
			}
		})
	}
	t.Run(componentAssignmentState, func(t *testing.T) {
		t.Parallel()
		f := mergeFixture(t, "delegate")
		sibling := f.created[0]
		before := assignmentCount(t, f.root)
		identityComponentFixtureFor(t, componentAssignmentState).mutate(t, f.root, sibling)

		requireCreateFromRefusal(t, runVerb(t, verbCreate, f.call("--request", "create-from-state",
			"--label", "dependent", "--from", sibling.Assignment.Label)),
			"bench worktree create: --from names no active assignment\n", "next=bench worktree list\n")
		if after := assignmentCount(t, f.root); after != before {
			t.Fatalf("ledger holds %d records, want the %d it held before the refusal", after, before)
		}
	})
}
