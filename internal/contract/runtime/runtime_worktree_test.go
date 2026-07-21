package runtime

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestRuntimeWorktreeContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench worktree create binds ownership contract", testRuntimeWorktreeCreateBindsOwnership)
	contract.RunParallel(t, "bench worktree release verifies and recovers contract", testRuntimeWorktreeReleaseVerifiesAndRecovers)
	contract.RunParallel(t, "bench worktree exact foreign cleanup matrix contract", testRuntimeWorktreeExactForeignCleanup)
	contract.RunParallel(t, "bench worktree recovery exact plan/apply contract", testRuntimeWorktreeRecoveryPlanApply)
	contract.RunParallel(t, "bench worktree clean pool cwd classification contract", testRuntimeWorktreeCleanFromPoolCwd)
	contract.RunParallel(t, "bench worktree usage contract", testRuntimeWorktreeRejectsUnknownArgs)
	contract.RunParallel(t, "bench worktree interactive release contract", testRuntimeWorktreeInteractiveRelease)
	contract.RunParallel(t, "bench worktree lease hardening contract", testRuntimeWorktreeLeaseHardening)
	contract.RunParallel(t, "bench worktree concurrent-acquire contract", testRuntimeWorktreeConcurrentAcquire)
	contract.RunParallel(t, "bench resume-clean and SessionStart preserve foreign worktrees contract", testRuntimeWorktreeForeignPreservation)
	contract.RunParallel(t, "bench resume-clean failure summary contract", testRuntimeResumeFailureSummary)
	contract.RunParallel(t, "bench automatic cleanup eligibility matrix contract", testRuntimeWorktreeAutomaticEligibilityMatrix)
	contract.RunParallel(t, "bench embedded repository preservation contract", testRuntimeWorktreeEmbeddedRepositoryPreservation)
	contract.RunParallel(t, "bench worktree exact dry-run purity contract", testRuntimeWorktreeExactDryRunPurity)
	contract.RunParallel(t, "bench worktree ignored discard contract", testRuntimeWorktreeIgnoredDiscard)
	contract.RunParallel(t, "bench worktree hostile surfaces contract", testRuntimeWorktreeHostileSurfaces)
	contract.RunParallel(t, "bench worktree list rows contract", testRuntimeWorktreeListRows)
	contract.RunParallel(t, "bench worktree list AXI posture contract", testRuntimeWorktreeListAXIPosture)
	contract.RunParallel(t, "bench worktree list shell surface contract", testRuntimeWorktreeListShellSurface)
}

func testRuntimeWorktreeListShellSurface(t *testing.T) {
	f := onMainFixture(t)
	general := f.Bench("--help")
	general.RequireExit(0)
	contract.RequireContains(t, general.Stdout, "bench worktree list")
	worktreeHelp := f.Bench("worktree", "--help")
	worktreeHelp.RequireExit(0)
	contract.RequireContains(t, worktreeHelp.Stdout, "bench worktree list")
	list := f.Bench("worktree", "list")
	list.RequireExit(0)
	if list.Stdout != "worktrees[0]{id,label,state,source,tree,lease,landed,ignored}:\n" || list.Stderr != "" {
		t.Fatalf("wrapper list streams = stdout %q stderr %q", list.Stdout, list.Stderr)
	}
}

func testRuntimeWorktreeListAXIPosture(t *testing.T) {
	f := onMainFixture(t)
	unknown := f.Bench("worktree", "list", "bogus")
	unknown.RequireExit(2)
	if unknown.Stdout != "usage: bench worktree list (unknown argument: bogus)\n" || unknown.Stderr != "" {
		t.Fatalf("unknown argument streams = stdout %q stderr %q", unknown.Stdout, unknown.Stderr)
	}
	for _, flag := range []string{"-h", "--help"} {
		help := f.Bench("worktree", "list", flag)
		help.RequireExit(0)
		if help.Stdout != "usage: bench worktree list\n" || help.Stderr != "" {
			t.Fatalf("%s streams = stdout %q stderr %q", flag, help.Stdout, help.Stderr)
		}
	}
	outside := contract.NewFixture(t, contract.WithNoRepo()).Bench("worktree", "list")
	outside.RequireExit(1)
	if outside.Stdout != "error: not in a git repository — run inside a Bench-linked repo\n" || outside.Stderr != "" {
		t.Fatalf("outside-repository streams = stdout %q stderr %q", outside.Stdout, outside.Stderr)
	}
}

func testRuntimeWorktreeListRows(t *testing.T) {
	empty := onMainFixture(t)
	emptyList := empty.Bench("worktree", "list")
	emptyList.RequireExit(0)
	if emptyList.Stdout != "worktrees[0]{id,label,state,source,tree,lease,landed,ignored}:\n" || emptyList.Stderr != "" {
		t.Fatalf("empty list streams = stdout %q stderr %q", emptyList.Stdout, emptyList.Stderr)
	}
	unresolved := contract.NewFixture(t)
	commitAllowEmpty(t, unresolved, "unresolved default fixture")
	unresolved.Git("branch", "other-default-candidate")
	unresolvedEnv := map[string]string{"BENCH_HOME": filepath.Join(unresolved.Root, ".bench-home")}
	createRuntimeAssignment(t, unresolved, unresolvedEnv, "list-unresolved")
	unresolvedList := unresolved.BenchEnv(unresolvedEnv, "worktree", "list")
	unresolvedList.RequireExit(0)
	contract.RequireContains(t, unresolvedList.Stdout, ",list-unresolved,active,assignment,present,none,unknown,0")

	f := onMainFixture(t)
	f.WriteFile(".gitignore", "list-ignored\n")
	f.CommitAll("list fixture")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
	present := createRuntimeAssignment(t, f, env, "list-present")
	contract.WriteFileAbs(t, filepath.Join(present.Worktree, "list-ignored"), "residue\n")
	dead := exec.Command("sh", "-c", "exit 0")
	if err := dead.Start(); err != nil {
		t.Fatal(err)
	}
	deadPID := dead.Process.Pid
	if err := dead.Wait(); err != nil {
		t.Fatal(err)
	}
	lease := strings.TrimSpace(contract.RunAt(t, f, present.Worktree, nil, "git", "rev-parse", "--path-format=absolute", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, lease, fmt.Sprintf("%d 2026-01-01T00:00:00Z\n", deadPID))
	missing := createRuntimeAssignment(t, f, env, "list-missing")
	f.Git("worktree", "unlock", missing.Worktree)
	f.Git("worktree", "remove", "--force", missing.Worktree)
	live := createRuntimeAssignment(t, f, env, "list-live")
	liveLease := strings.TrimSpace(contract.RunAt(t, f, live.Worktree, nil, "git", "rev-parse", "--path-format=absolute", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, liveLease, fmt.Sprintf("%d 2026-01-01T00:00:00Z\n", os.Getpid()))
	unknown := createRuntimeAssignment(t, f, env, "list-unknown")
	unknownLease := strings.TrimSpace(contract.RunAt(t, f, unknown.Worktree, nil, "git", "rev-parse", "--path-format=absolute", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, unknownLease, "not a lease\n")
	foreignPath := filepath.Join(t.TempDir(), "foreign-list")
	f.Git("worktree", "add", "-q", "-b", "foreign-list", foreignPath, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(foreignPath, "foreign.txt"), "unique\n")
	contract.RunAt(t, f, foreignPath, nil, "git", "add", "foreign.txt").RequireExit(0)
	contract.RunAt(t, f, foreignPath, nil, "git", "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "foreign unique").RequireExit(0)
	missingForeignPath := filepath.Join(t.TempDir(), "foreign-list-missing")
	f.Git("worktree", "add", "-q", "-b", "foreign-list-missing", missingForeignPath, "HEAD")
	if err := os.Rename(missingForeignPath, missingForeignPath+"-moved"); err != nil {
		t.Fatal(err)
	}
	deep := filepath.Join(f.Root, "deep", "cwd")
	contract.Mkdir(t, deep)

	out := contract.RunAt(t, f, deep, env, "bash", benchPath(t), "worktree", "list")
	out.RequireExit(0)
	if out.Stderr != "" {
		t.Fatalf("list wrote stderr: %q", out.Stderr)
	}
	lines := contract.NonEmptyLines(out.Stdout)
	if len(lines) != 7 || lines[0] != "worktrees[6]{id,label,state,source,tree,lease,landed,ignored}:" {
		t.Fatalf("list block = %q", out.Stdout)
	}
	for _, want := range []string{
		",list-present,active,assignment,present,dead,true,1",
		",list-missing,active,assignment,missing,none,true,unknown",
		",list-live,active,assignment,present,live,true,0",
		",list-unknown,active,assignment,present,unknown,true,0",
		"foreign,foreign-list,foreign,foreign,present,none,false,0",
		"foreign,foreign-list-missing,foreign,foreign,missing,none,true,unknown",
	} {
		contract.RequireContains(t, out.Stdout, want)
	}

	hostile := createRuntimeAssignment(t, f, env, "hostile-label")
	hostile.Label = "unsafe\x1blabel"
	if err := intent.PutAssignment(f.Root, hostile); err != nil {
		t.Fatal(err)
	}
	refused := f.BenchEnv(env, "worktree", "list")
	refused.RequireExit(1)
	contract.RequireContains(t, refused.Stdout, "error: unrepresentable TOON cell")
	if refused.Stderr != "" || strings.Contains(refused.Stdout, "\x1b") {
		t.Fatalf("hostile list streams = stdout %q stderr %q", refused.Stdout, refused.Stderr)
	}
}

// onMainFixture returns a fixture whose default branch resolves to a real `main` commit — the
// sweep's happy-path precondition. git.DefaultBranch falls back to "main", but a bare `git init`
// fixture is born on "master", so an explicit HEAD symref lands the first commit on main and makes
// the default branch resolve.
func onMainFixture(t *testing.T, opts ...contract.FixtureOption) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t, opts...)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	commitAllowEmpty(t, f, "init")
	return f
}

// headExists reports whether refs/heads/<name> resolves — an exit-code probe so a slashed or
// unicode branch name is checked verbatim, never through a substring or glob match.
func headExists(f contract.Fixture, name string) bool {
	return f.GitAllow("show-ref", "--verify", "--quiet", "refs/heads/"+name).ExitCode == 0
}

// FT77 story 2: the two accused unattended surfaces must classify ordinary and
// detached-unique registrations as foreign, report that reason, and leave every
// filesystem/Git-visible layer intact.
func testRuntimeWorktreeForeignPreservation(t *testing.T) {
	f := onMainFixture(t)
	f.Bench("link").RequireExit(0)
	f.CommitAll("link fixture")
	ordinary := filepath.Join(t.TempDir(), "ordinary foreign")
	detached := filepath.Join(t.TempDir(), "detached foreign")
	f.Git("worktree", "add", "-q", "-b", "foreign-ordinary", ordinary, "HEAD")
	f.Git("worktree", "add", "-q", "--detach", detached, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(detached, "unique.txt"), "unique\n")
	contract.RunAt(t, f, detached, nil, "git", "add", "unique.txt").RequireExit(0)
	contract.RunAt(t, f, detached, nil, "git", "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "unique detached").RequireExit(0)
	unique := strings.TrimSpace(contract.RunAt(t, f, detached, nil, "git", "rev-parse", "HEAD").Stdout)

	out := f.Bench("resume-clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "retained foreign=2")
	for _, path := range []string{ordinary, detached} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("resume-clean removed foreign worktree %q: %v", path, err)
		}
	}
	if !headExists(f, "foreign-ordinary") {
		t.Fatal("resume-clean deleted foreign branch")
	}
	if got := strings.TrimSpace(contract.RunAt(t, f, detached, nil, "git", "rev-parse", "HEAD").Stdout); got != unique {
		t.Fatalf("resume-clean changed detached unique HEAD: got %s want %s", got, unique)
	}

	hook := filepath.Join(f.Root, ".bench", "hooks", "session-start.sh")
	session := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook)
	session.RequireExit(0)
	contract.RequireContains(t, session.Stdout+session.Stderr, "retained foreign=2")
	for _, path := range []string{ordinary, detached} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("SessionStart removed foreign worktree %q: %v", path, err)
		}
	}
}

func testRuntimeResumeFailureSummary(t *testing.T) {
	f := onMainFixture(t)
	f.Bench("link").RequireExit(0)
	f.CommitAll("link fixture")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
	assignment := createRuntimeAssignment(t, f, env, "resume-failure")
	assignment.State = intent.StateCleanupPending
	if err := intent.PutAssignment(f.Root, assignment); err != nil {
		t.Fatal(err)
	}
	contract.WriteFileAbs(t, filepath.Join(assignment.Worktree, "dirty.txt"), "preserve after failure\n")
	refs := filepath.Join(gitDir(t, f), "refs")
	if err := os.Chmod(refs, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(refs, 0o700) })

	want := "bench resume: removed 0, recovered 0; reconciled 0; failed 1; open assignments 1\n"
	resume := f.BenchEnv(env, "resume-clean")
	resume.RequireExit(1)
	if resume.Stdout != want || !strings.Contains(resume.Stderr, "bench resume-clean:") {
		t.Fatalf("failed resume streams = stdout %q stderr %q", resume.Stdout, resume.Stderr)
	}
	if _, err := os.Stat(assignment.Worktree); err != nil {
		t.Fatalf("failed resume removed assignment: %v", err)
	}
	contract.RequireContains(t, f.Git("worktree", "list", "--porcelain").Stdout, "locked bench owner=")

	hook := filepath.Join(f.Root, ".bench", "hooks", "session-start.sh")
	session := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin", "BENCH_HOME": env["BENCH_HOME"]}, "bash", hook)
	session.RequireExit(0)
	contract.RequireContains(t, session.Stdout, want)
	contract.RequireContains(t, session.Stderr, "warning: bench session-start: resume-clean failed; inspect retained worktree state")
}

// FT77 story 2: one real CLI run proves the two eligible cases are not hidden by a
// safe no-op while each fail-closed class survives with its own stable reason.
func testRuntimeWorktreeAutomaticEligibilityMatrix(t *testing.T) {
	f := onMainFixture(t)
	f.Bench("link").RequireExit(0)
	f.WriteFile(".gitignore", "ignored.txt\n")
	f.CommitAll("matrix fixture")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}

	clean := createRuntimeAssignment(t, f, env, "matrix-clean")
	dirty := createRuntimeAssignment(t, f, env, "matrix-dirty")
	active := createRuntimeAssignment(t, f, env, "matrix-active")
	live := createRuntimeAssignment(t, f, env, "matrix-live")
	unmerged := createRuntimeAssignment(t, f, env, "matrix-unmerged")
	ignored := createRuntimeAssignment(t, f, env, "matrix-ignored")
	malformed := createRuntimeAssignment(t, f, env, "matrix-malformed")
	uncertain := createRuntimeAssignment(t, f, env, "matrix-uncertain")
	locked := createRuntimeAssignment(t, f, env, "matrix-locked")
	for _, assignment := range []intent.Assignment{clean, dirty, live, unmerged, ignored, malformed, uncertain, locked} {
		assignment.State = intent.StateCleanupPending
		if err := intent.PutAssignment(f.Root, assignment); err != nil {
			t.Fatal(err)
		}
	}
	contract.WriteFileAbs(t, filepath.Join(dirty.Worktree, "dirty.txt"), "recover me\n")
	lease := strings.TrimSpace(contract.RunAt(t, f, live.Worktree, nil, "git", "rev-parse", "--path-format=absolute", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, lease, fmt.Sprintf("%d 2026-01-01T00:00:00Z\n", os.Getpid()))
	contract.WriteFileAbs(t, filepath.Join(unmerged.Worktree, "unique.txt"), "unique\n")
	contract.RunAt(t, f, unmerged.Worktree, nil, "git", "add", "unique.txt").RequireExit(0)
	contract.RunAt(t, f, unmerged.Worktree, nil, "git", "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "unique").RequireExit(0)
	contract.WriteFileAbs(t, filepath.Join(ignored.Worktree, "ignored.txt"), "secret\n")
	malformedAdmin := strings.TrimSpace(contract.RunAt(t, f, malformed.Worktree, nil, "git", "rev-parse", "--path-format=absolute", "--git-dir").Stdout)
	contract.WriteFileAbs(t, filepath.Join(malformedAdmin, "bench-owner"), "not json\n")
	contract.RunAt(t, f, uncertain.Worktree, nil, "git", "switch", "-q", "--detach").RequireExit(0)
	f.Git("worktree", "unlock", locked.Worktree)
	f.Git("worktree", "lock", "--reason", "foreign lock", locked.Worktree)

	out := f.BenchEnv(env, "resume-clean")
	out.RequireExit(0)
	want := "bench resume: removed 1, recovered 1; retained active=1 live-lease=1 unmerged=1 ignored=1 malformed=1 uncertain=1 unexpected-lock=1; reconciled 0; failed 0; open assignments 8\n"
	if out.Stdout != want {
		t.Fatalf("automatic matrix summary = %q, want %q", out.Stdout, want)
	}
	for _, path := range []string{active.Worktree, live.Worktree, unmerged.Worktree, ignored.Worktree, malformed.Worktree, uncertain.Worktree, locked.Worktree} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("automatic cleanup removed retained matrix target %q: %v", path, err)
		}
	}
	for _, path := range []string{clean.Worktree, dirty.Worktree} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("automatic cleanup kept eligible target %q: %v", path, err)
		}
	}
	contract.RequireContains(t, f.Git("for-each-ref", "--format=%(refname)", "refs/bench/recovery/").Stdout, "refs/bench/recovery/")

	hook := filepath.Join(f.Root, ".bench", "hooks", "session-start.sh")
	session := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin", "BENCH_HOME": env["BENCH_HOME"]}, "bash", hook)
	session.RequireExit(0)
	for _, reason := range []string{"active=1", "live-lease=1", "unmerged=1", "ignored=1", "malformed=1", "uncertain=1", "unexpected-lock=1"} {
		contract.RequireContains(t, session.Stdout+session.Stderr, reason)
	}
}

func createRuntimeAssignment(t *testing.T, f contract.Fixture, env map[string]string, request string) intent.Assignment {
	t.Helper()
	f.BenchEnv(env, "worktree", "create", "--request", request, "--label", request).RequireExit(0)
	assignments, err := intent.Assignments(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	for _, assignment := range assignments {
		if assignment.Label == request {
			return assignment
		}
	}
	t.Fatalf("created assignment %q not found", request)
	return intent.Assignment{}
}

func testRuntimeWorktreeEmbeddedRepositoryPreservation(t *testing.T) {
	for _, gitFile := range []bool{false, true} {
		t.Run(fmt.Sprintf("git-file=%t", gitFile), func(t *testing.T) {
			f := onMainFixture(t)
			env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
			assignment := createRuntimeAssignment(t, f, env, fmt.Sprintf("embedded-%t", gitFile))
			nested := filepath.Join(assignment.Worktree, "embedded")
			contract.Mkdir(t, nested)
			initArgs := []string{"init", "-q", "-b", "main"}
			if gitFile {
				initArgs = append(initArgs, "--separate-git-dir", filepath.Join(t.TempDir(), "embedded.git"))
			}
			contract.RunAt(t, f, nested, nil, "git", initArgs...).RequireExit(0)
			contract.WriteFileAbs(t, filepath.Join(nested, "unique.txt"), "unique nested checkout\n")
			contract.RunAt(t, f, nested, nil, "git", "add", "unique.txt").RequireExit(0)
			contract.RunAt(t, f, nested, nil, "git", "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "unique nested").RequireExit(0)
			unique := strings.TrimSpace(contract.RunAt(t, f, nested, nil, "git", "rev-parse", "HEAD").Stdout)

			plan := f.Bench("worktree", "clean", assignment.Worktree)
			plan.RequireExit(0)
			contract.RequireContains(t, plan.Stdout, ",retain,")
			f.Bench("worktree", "clean", assignment.Worktree, "--apply", cleanupFingerprint(t, plan.Stdout)).RequireExit(0)
			assignment.State = intent.StateCleanupPending
			if err := intent.PutAssignment(f.Root, assignment); err != nil {
				t.Fatal(err)
			}
			resume := f.BenchEnv(env, "resume-clean")
			resume.RequireExit(0)
			contract.RequireContains(t, resume.Stdout, "retained uncertain=1")
			if got := strings.TrimSpace(contract.RunAt(t, f, nested, nil, "git", "rev-parse", "HEAD").Stdout); got != unique {
				t.Fatalf("embedded checkout changed: got %s want %s", got, unique)
			}
		})
	}
}

func testRuntimeWorktreeExactDryRunPurity(t *testing.T) {
	f := onMainFixture(t)
	target := filepath.Join(t.TempDir(), "pure exact target")
	f.Git("worktree", "add", "-q", "-b", "pure-exact", target, "HEAD")
	before := snapshotRuntimePaths(t, gitDir(t, f), target)

	first := f.Bench("worktree", "clean", target)
	first.RequireExit(0)
	if first.Stderr != "" {
		t.Fatalf("dry-run wrote stderr: %q", first.Stderr)
	}
	if !strings.HasSuffix(first.Stdout, "\n") {
		t.Fatalf("dry-run lacks trailing newline: %q", first.Stdout)
	}
	lines := contract.NonEmptyLines(first.Stdout)
	if len(lines) != 2 || lines[0] != "worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:" {
		t.Fatalf("dry-run literal schema = %q", first.Stdout)
	}
	fields := strings.Split(strings.TrimSpace(lines[1]), ",")
	if len(fields) != 7 {
		t.Fatalf("dry-run row fields = %#v", fields)
	}
	if got := strings.Trim(fields[3], `"`); got != "count=0 bytes=0 shown=0 truncated=false" {
		t.Fatalf("ignored summary = %q", got)
	}
	fingerprint := strings.Trim(fields[5], `"`)
	if len(fingerprint) != 64 || strings.ToLower(fingerprint) != fingerprint {
		t.Fatalf("fingerprint = %q, want lowercase SHA-256", fingerprint)
	}
	if got := strings.Trim(fields[6], `"`); got != "apply with exact fingerprint" {
		t.Fatalf("detail = %q", got)
	}

	second := f.Bench("worktree", "clean", target)
	second.RequireExit(0)
	if second.Stdout != first.Stdout || second.Stderr != "" {
		t.Fatalf("dry-run is nondeterministic\nfirst=%q\nsecond=%q", first.Stdout, second.Stdout)
	}
	if after := snapshotRuntimePaths(t, gitDir(t, f), target); after != before {
		t.Fatalf("dry-run mutated Git/private/filesystem state\nbefore:\n%s\nafter:\n%s", before, after)
	}

	contract.WriteFileAbs(t, filepath.Join(target, "README.md"), "drifted after plan\n")
	beforeStale := snapshotRuntimePaths(t, gitDir(t, f), target)
	stale := f.Bench("worktree", "clean", target, "--apply", fingerprint)
	stale.RequireExit(1)
	contract.RequireContains(t, stale.Stdout, "worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:")
	contract.RequireNotContains(t, stale.Stdout, fingerprint)
	if stale.Stderr != "" {
		t.Fatalf("stale apply wrote an unrelated stderr shape: %q", stale.Stderr)
	}
	if after := snapshotRuntimePaths(t, gitDir(t, f), target); after != beforeStale {
		t.Fatalf("stale built apply mutated target\nbefore:\n%s\nafter:\n%s", beforeStale, after)
	}
}

func snapshotRuntimePaths(t *testing.T, roots ...string) string {
	t.Helper()
	var snapshot strings.Builder
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			fmt.Fprintf(&snapshot, "%s %s %o %d %d", root, rel, info.Mode(), info.Size(), info.ModTime().UnixNano())
			if info.Mode().IsRegular() {
				body, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				digest := sha256.Sum256(body)
				fmt.Fprintf(&snapshot, " %x", digest)
			}
			snapshot.WriteByte('\n')
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return snapshot.String()
}

func testRuntimeWorktreeIgnoredDiscard(t *testing.T) {
	f := onMainFixture(t)
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "info", "exclude"), "ignored-*\n")
	target := filepath.Join(t.TempDir(), "ignored exact")
	f.Git("worktree", "add", "-q", "-b", "ignored-exact", target, "HEAD")
	for i := 0; i < 21; i++ {
		contract.WriteFileAbs(t, filepath.Join(target, fmt.Sprintf("ignored-%02d", i)), "secret\n")
	}
	external := filepath.Join(t.TempDir(), "external-secret")
	contract.WriteFileAbs(t, external, "outside\n")
	if err := os.Symlink(external, filepath.Join(target, "ignored-link")); err != nil {
		t.Fatal(err)
	}

	plan := f.Bench("worktree", "clean", "--discard-ignored", target)
	plan.RequireExit(0)
	contract.RequireContains(t, plan.Stdout, "discard-remove")
	contract.RequireContains(t, plan.Stdout, "ignored_paths[20]{path}:")
	fingerprint := cleanupFingerprint(t, plan.Stdout)
	full := f.Bench("worktree", "clean", "--discard-ignored", "--full", target)
	full.RequireExit(0)
	contract.RequireContains(t, full.Stdout, "ignored_paths[22]{path}:")
	if cleanupFingerprint(t, full.Stdout) != fingerprint {
		t.Fatal("--full changed cleanup fingerprint")
	}

	contract.WriteFileAbs(t, filepath.Join(target, "ignored-new"), "new\n")
	stale := f.Bench("worktree", "clean", "--discard-ignored", target, "--apply", fingerprint)
	stale.RequireExit(1)
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("inventory drift removed target: %v", err)
	}
	fresh := f.Bench("worktree", "clean", "--discard-ignored", target)
	fresh.RequireExit(0)
	apply := f.Bench("worktree", "clean", "--discard-ignored", target, "--apply", cleanupFingerprint(t, fresh.Stdout))
	apply.RequireExit(0)
	if body := contract.ReadFileAbs(t, external); body != "outside\n" {
		t.Fatalf("discard followed symlink: %q", body)
	}
}

func cleanupFingerprint(t *testing.T, output string) string {
	t.Helper()
	lines := contract.NonEmptyLines(output)
	if len(lines) < 2 || !strings.HasPrefix(lines[0], "worktree_cleanup[") {
		t.Fatalf("cleanup output = %q", output)
	}
	fields := strings.Split(strings.TrimSpace(lines[1]), ",")
	if len(fields) != 7 {
		t.Fatalf("cleanup fields = %#v", fields)
	}
	return strings.Trim(fields[5], `"`)
}

func testRuntimeWorktreeCreateBindsOwnership(t *testing.T) {
	f := onMainFixture(t)
	benchHome := filepath.Join(f.Root, ".bench-home")
	env := map[string]string{"BENCH_HOME": benchHome}

	out := f.BenchEnv(env, "worktree", "create", "--request", "FT77-probe", "--label", "FT77 probe")
	out.RequireExit(0)
	path := worktreeCreatePath(t, out.Stdout)
	if !filepath.IsAbs(path) {
		t.Fatalf("worktree path = %q, want absolute", path)
	}
	branch := strings.TrimSpace(contract.RunAt(t, f, path, nil, "git", "symbolic-ref", "--quiet", "--short", "HEAD").Stdout)
	if !strings.HasPrefix(branch, "bench/assign/") {
		t.Fatalf("branch = %q, want dedicated bench/assign branch", branch)
	}
	admin := strings.TrimSpace(contract.RunAt(t, f, path, nil, "git", "rev-parse", "--path-format=absolute", "--git-dir").Stdout)
	marker := filepath.Join(admin, "bench-owner")
	info, err := os.Lstat(marker)
	if err != nil {
		t.Fatalf("ownership marker missing: %v", err)
	}
	if info.Mode().Perm() != 0o600 || !info.Mode().IsRegular() {
		t.Fatalf("ownership marker mode/type = %v, want regular 0600", info.Mode())
	}
	markerBody := contract.ReadFileAbs(t, marker)
	contract.RequireContains(t, markerBody, `"schema":"bench-owner/v1"`)
	contract.RequireContains(t, markerBody, `"path":"`)
	ledger := contract.ReadFileAbs(t, filepath.Join(f.Root, ".git", "bench-intent.json"))
	contract.RequireContains(t, ledger, path)
	contract.RequireContains(t, ledger, "refs/heads/"+branch)
	registration := f.Git("worktree", "list", "--porcelain").Stdout
	contract.RequireContains(t, registration, "locked bench owner=")
	contract.RequireContains(t, registration, " assignment=")

	retry := f.BenchEnv(env, "worktree", "create", "--request", "FT77-probe", "--label", "FT77 probe")
	retry.RequireExit(0)
	if got := worktreeCreatePath(t, retry.Stdout); got != path {
		t.Fatalf("idempotent retry path = %q, want %q", got, path)
	}
}

func testRuntimeWorktreeReleaseVerifiesAndRecovers(t *testing.T) {
	f := onMainFixture(t)
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
	created := f.BenchEnv(env, "worktree", "create", "--request", "release-1", "--label", "release probe")
	created.RequireExit(0)
	lines := contract.NonEmptyLines(created.Stdout)
	if len(lines) != 2 {
		t.Fatalf("create output = %q", created.Stdout)
	}
	path := strings.Trim(strings.TrimSpace(strings.Split(lines[1], ",")[0]), `"`)
	ledgerPath := filepath.Join(f.Root, ".git", "bench-intent.json")
	before := contract.ReadFileAbs(t, ledgerPath)
	mismatch := f.BenchEnv(env, "worktree", "release", "--request", "wrong", path)
	mismatch.RequireExit(1)
	if after := contract.ReadFileAbs(t, ledgerPath); after != before {
		t.Fatal("release mismatch mutated assignment ledger")
	}
	contract.RequireContains(t, f.Git("worktree", "list", "--porcelain").Stdout, "locked bench owner=")
	contract.WriteFileAbs(t, filepath.Join(path, "dirty.txt"), "dirty\n")
	released := f.BenchEnv(env, "worktree", "release", "--request", "release-1", path)
	released.RequireExit(0)
	contract.RequireContains(t, released.Stdout, "recovered")
	contract.RequireNotContains(t, f.Git("worktree", "list", "--porcelain").Stdout, path)
	contract.RequireContains(t, f.Git("for-each-ref", "--format=%(refname)", "refs/bench/recovery/").Stdout, "refs/bench/recovery/")
}

func testRuntimeWorktreeExactForeignCleanup(t *testing.T) {
	t.Run("one sibling only and branch survives", func(t *testing.T) {
		f := onMainFixture(t)
		one := filepath.Join(t.TempDir(), "foreign one")
		two := filepath.Join(t.TempDir(), "foreign two")
		f.Git("worktree", "add", "-q", "-b", "foreign-one", one, "HEAD")
		f.Git("worktree", "add", "-q", "-b", "foreign-two", two, "HEAD")
		plan := f.Bench("worktree", "clean", one)
		plan.RequireExit(0)
		fingerprint := cleanupFingerprint(t, plan.Stdout)
		apply := f.Bench("worktree", "clean", one, "--apply", fingerprint)
		apply.RequireExit(0)
		replay := f.Bench("worktree", "clean", one, "--apply", fingerprint)
		replay.RequireExit(0)
		if replay.Stdout != apply.Stdout || replay.Stderr != "" {
			t.Fatalf("idempotent replay changed result\napply=%q\nreplay=%q\nstderr=%q", apply.Stdout, replay.Stdout, replay.Stderr)
		}
		if _, err := os.Stat(one); !os.IsNotExist(err) {
			t.Fatalf("exact target remains: %v", err)
		}
		if _, err := os.Stat(two); err != nil {
			t.Fatalf("sibling was removed: %v", err)
		}
		for _, branch := range []string{"foreign-one", "foreign-two"} {
			if !headExists(f, branch) {
				t.Fatalf("exact foreign cleanup deleted branch %q", branch)
			}
		}
	})

	t.Run("dirty attached recovers before removal", func(t *testing.T) {
		f := onMainFixture(t)
		target := filepath.Join(t.TempDir(), "dirty attached")
		f.Git("worktree", "add", "-q", "-b", "foreign-dirty", target, "HEAD")
		contract.WriteFileAbs(t, filepath.Join(target, "dirty.txt"), "recover me\n")
		apply := applyCleanupPlan(t, f, target)
		apply.RequireExit(0)
		contract.RequireContains(t, apply.Stdout, "removed")
		if _, err := os.Stat(target); !os.IsNotExist(err) {
			t.Fatalf("dirty exact target remains: %v", err)
		}
		if !headExists(f, "foreign-dirty") {
			t.Fatal("dirty foreign branch was deleted")
		}
		refs := f.Git("for-each-ref", "--format=%(refname)", "refs/bench/recovery/").Stdout
		contract.RequireContains(t, refs, "refs/bench/recovery/")
	})

	t.Run("detached unique anchored", func(t *testing.T) {
		f := onMainFixture(t)
		target := filepath.Join(t.TempDir(), "detached exact")
		f.Git("worktree", "add", "-q", "--detach", target, "HEAD")
		contract.WriteFileAbs(t, filepath.Join(target, "unique.txt"), "unique\n")
		contract.RunAt(t, f, target, nil, "git", "add", "unique.txt").RequireExit(0)
		contract.RunAt(t, f, target, nil, "git", "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "detached unique").RequireExit(0)
		unique := strings.TrimSpace(contract.RunAt(t, f, target, nil, "git", "rev-parse", "HEAD").Stdout)
		applyCleanupPlan(t, f, target).RequireExit(0)
		refs := strings.Fields(f.Git("for-each-ref", "--format=%(objectname)", "refs/bench/recovery/").Stdout)
		if len(refs) != 1 || refs[0] != unique {
			t.Fatalf("detached recovery refs = %#v, want %s", refs, unique)
		}
	})

	t.Run("refusals", func(t *testing.T) {
		f := onMainFixture(t)
		registered := filepath.Join(t.TempDir(), "registered parent")
		f.Git("worktree", "add", "-q", "-b", "foreign-parent", registered, "HEAD")
		inside := filepath.Join(registered, "child")
		contract.Mkdir(t, inside)
		unregistered := t.TempDir()
		other := onMainFixture(t)
		crossRepo := filepath.Join(t.TempDir(), "cross repo")
		other.Git("worktree", "add", "-q", "-b", "cross-repo", crossRepo, "HEAD")
		for name, target := range map[string]string{"primary": f.Root, "inside": inside, "unregistered": unregistered, "cross-repository": crossRepo} {
			plan := f.Bench("worktree", "clean", target)
			plan.RequireExit(0)
			contract.RequireContains(t, plan.Stdout, "retain")
			applied := f.Bench("worktree", "clean", target, "--apply", cleanupFingerprint(t, plan.Stdout))
			applied.RequireExit(0)
			contract.RequireContains(t, applied.Stdout, "retain")
			if _, err := os.Stat(target); err != nil {
				t.Fatalf("%s refusal lost target: %v", name, err)
			}
		}
	})

	t.Run("classification failure stays structured", func(t *testing.T) {
		f := onMainFixture(t)
		target := filepath.Join(t.TempDir(), "failed classification")
		f.Git("worktree", "add", "-q", "-b", "failed-classification", target, "HEAD")
		contract.WriteFileAbs(t, filepath.Join(f.Root, ".git", "bench-intent.json"), "{broken\n")
		out := f.Bench("worktree", "clean", target)
		out.RequireExit(1)
		if out.Stderr != "" || !strings.HasSuffix(out.Stdout, "\n") {
			t.Fatalf("failed intent streams = stdout %q stderr %q", out.Stdout, out.Stderr)
		}
		contract.RequireContains(t, out.Stdout, "worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\n")
		contract.RequireContains(t, out.Stdout, ",error,unknown,")
	})
}

func applyCleanupPlan(t *testing.T, f contract.Fixture, target string) contract.Probe {
	t.Helper()
	plan := f.Bench("worktree", "clean", target)
	plan.RequireExit(0)
	lines := contract.NonEmptyLines(plan.Stdout)
	if len(lines) != 2 {
		t.Fatalf("cleanup plan = %q", plan.Stdout)
	}
	fields := strings.Split(strings.TrimSpace(lines[1]), ",")
	if len(fields) != 7 {
		t.Fatalf("cleanup fields = %#v", fields)
	}
	return f.Bench("worktree", "clean", target, "--apply", strings.Trim(fields[5], `"`))
}

func testRuntimeWorktreeRecoveryPlanApply(t *testing.T) {
	f := onMainFixture(t)
	target := filepath.Join(t.TempDir(), "recovery retirement")
	f.Git("worktree", "add", "-q", "-b", "recovery-retirement", target, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(target, "recovered.txt"), "land me\n")
	applyCleanupPlan(t, f, target).RequireExit(0)
	assignments, err := intent.Assignments(f.Root)
	if err != nil || len(assignments) != 1 || len(assignments[0].Recovery) != 1 {
		t.Fatalf("foreign recovery assignment = %#v, %v", assignments, err)
	}
	assignment := assignments[0]
	first := assignment.Recovery[0]
	base := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	for _, payload := range first.Payloads {
		f.Git("-c", "user.name=bench", "-c", "user.email=bench@local", "cherry-pick", payload)
	}
	first.Payloads = append(first.Payloads, base)
	assignment.Recovery[0] = first

	f.Git("checkout", "-q", "-b", "recovery-unique")
	f.WriteFile("unique-recovery.txt", "unique\n")
	f.CommitAll("unique recovery")
	unique := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	uniqueTree := strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout)
	f.Git("checkout", "-q", "main")
	uniqueRoot := strings.TrimSpace(f.Git("-c", "user.name=bench", "-c", "user.email=bench@local", "commit-tree", uniqueTree, "-p", unique, "-m", "unique root").Stdout)
	second := first
	second.Ref = strings.TrimSuffix(first.Ref, "/1") + "/2"
	second.Root, second.Payloads = uniqueRoot, []string{unique}
	f.Git("update-ref", second.Ref, uniqueRoot)
	assignment.Recovery = append(assignment.Recovery, second)
	if err := intent.PutAssignment(f.Root, assignment); err != nil {
		t.Fatal(err)
	}
	uniquePlan := f.Bench("worktree", "recovery", second.Ref)
	uniquePlan.RequireExit(0)
	contract.RequireContains(t, uniquePlan.Stdout, "unlanded")

	mergePayload := strings.TrimSpace(f.Git("-c", "user.name=bench", "-c", "user.email=bench@local", "commit-tree", uniqueTree, "-p", unique, "-p", base, "-m", "ambiguous merge payload").Stdout)
	mergeRoot := strings.TrimSpace(f.Git("-c", "user.name=bench", "-c", "user.email=bench@local", "commit-tree", uniqueTree, "-p", mergePayload, "-m", "merge root").Stdout)
	second.Root, second.Payloads = mergeRoot, []string{mergePayload}
	f.Git("update-ref", second.Ref, mergeRoot)
	assignment.Recovery[1] = second
	if err := intent.PutAssignment(f.Root, assignment); err != nil {
		t.Fatal(err)
	}
	mergePlan := f.Bench("worktree", "recovery", second.Ref)
	mergePlan.RequireExit(0)
	contract.RequireContains(t, mergePlan.Stdout, "unlanded")

	second.Payloads = append(second.Payloads, strings.Repeat("f", 40))
	assignment.Recovery[1] = second
	if err := intent.PutAssignment(f.Root, assignment); err != nil {
		t.Fatal(err)
	}
	queryPlan := f.Bench("worktree", "recovery", second.Ref)
	queryPlan.RequireExit(0)
	contract.RequireContains(t, queryPlan.Stdout, "retain")

	plan := f.BenchEnv(map[string]string{"SHELL": "/bin/true"}, "worktree", "recovery", first.Ref)
	plan.RequireExit(0)
	contract.RequireContains(t, plan.Stdout, "recovery_cleanup[1]{ref,root,payloads,landed,action,fingerprint,detail}:")
	contract.RequireContains(t, plan.Stdout, first.Root)
	for _, payload := range first.Payloads {
		contract.RequireContains(t, plan.Stdout, payload)
	}
	fields := strings.Split(strings.TrimSpace(contract.NonEmptyLines(plan.Stdout)[1]), ",")
	if len(fields) != 7 {
		t.Fatalf("recovery plan fields = %#v", fields)
	}
	fingerprint := strings.Trim(fields[5], `"`)
	f.Git("update-ref", second.Ref, base)
	stale := f.Bench("worktree", "recovery", first.Ref, "--apply", fingerprint)
	stale.RequireExit(1)
	contract.RequireContains(t, stale.Stdout, "recovery_cleanup[1]{ref,root,payloads,landed,action,fingerprint,detail}:")
	if f.GitAllow("show-ref", "--verify", "--quiet", first.Ref).ExitCode != 0 {
		t.Fatal("stale recovery apply deleted named ref")
	}
	f.Git("update-ref", second.Ref, second.Root)
	plan = f.Bench("worktree", "recovery", first.Ref)
	plan.RequireExit(0)
	apply := f.Bench("worktree", "recovery", first.Ref, "--apply", strings.Trim(strings.Split(strings.TrimSpace(contract.NonEmptyLines(plan.Stdout)[1]), ",")[5], `"`))
	apply.RequireExit(0)
	contract.RequireContains(t, apply.Stdout, "retired")
	if f.GitAllow("show-ref", "--verify", "--quiet", first.Ref).ExitCode == 0 {
		t.Fatal("exact retired ref survived")
	}
	if f.GitAllow("show-ref", "--verify", "--quiet", second.Ref).ExitCode != 0 {
		t.Fatal("sibling recovery ref was deleted")
	}
	current, err := intent.Assignments(f.Root)
	if err != nil || len(current) != 1 || len(current[0].Recovery) != 1 || current[0].Recovery[0].Ref != second.Ref {
		t.Fatalf("sibling recovery context = %#v, %v", current, err)
	}
}

func testRuntimeWorktreeCleanFromPoolCwd(t *testing.T) {
	f := onMainFixture(t)
	benchHome := filepath.Join(f.Root, ".bh")
	pool := addRuntimePoolWorktrees(t, f, benchHome)

	out := contract.RunAt(t, f, pool.Leased, map[string]string{"BENCH_HOME": benchHome}, "bash", benchPath(t), "resume-clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "retained foreign=1 live-lease=1")
	contract.RequireNotContains(t, out.Stdout, f.Root)
	contract.RequireNotContains(t, out.Stdout, pool.Warm)
	contract.RequireNotContains(t, out.Stdout, pool.Leased)
}

func testRuntimeWorktreeRejectsUnknownArgs(t *testing.T) {
	wantRecovery := "recovery_cleanup[1]{ref,root,payloads,landed,action,fingerprint,detail}:\n  unknown,unknown,none,unknown,error,none,\"invalid invocation; run bench worktree recovery <ref> [--apply <fingerprint>]\"\n"
	noRepo := contract.NewFixture(t, contract.WithNoRepo())
	outside := noRepo.Bench("worktree", "recovery")
	outside.RequireExit(2)
	if outside.Stdout != wantRecovery || outside.Stderr != "" {
		t.Fatalf("outside-repository recovery usage streams = stdout %q stderr %q", outside.Stdout, outside.Stderr)
	}

	f := onMainFixture(t)
	canonicalUsage := []string{
		"bench worktree create [--refresh] --request <opaque-id> --label <work-item>",
		"bench worktree release --request <opaque-id> <path>",
		"bench worktree clean [--discard-ignored] [--full] <path> [--apply <fingerprint>]",
		"bench worktree recovery <ref> [--apply <fingerprint>]",
	}
	help := f.Bench("--help")
	help.RequireExit(0)
	contract.RequireContains(t, help.Stdout, "bench worktree --help")
	worktreeHelp := f.Bench("worktree", "--help")
	worktreeHelp.RequireExit(0)
	for _, usage := range canonicalUsage {
		contract.RequireContains(t, worktreeHelp.Stdout, usage)
	}
	branch := "worktree-same-prefix-sibling"
	registeredBranch := "worktree-registered-sibling"
	target := filepath.Join(t.TempDir(), "same prefix target")
	f.Git("branch", branch)
	f.Git("worktree", "add", "-q", "-b", registeredBranch, target, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(target, "dirty.txt"), "must survive\n")
	registrations := f.Git("worktree", "list", "--porcelain").Stdout
	wantClean := "worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\n  unknown,error,unknown,unknown,none,none,\"invalid invocation; run bench worktree clean [--discard-ignored] [--full] <path> [--apply <fingerprint>]\"\n"
	for _, args := range [][]string{{"worktree", "clean"}, {"worktree", "clean", "one", "two"}, {"worktree", "clean", "one", "--apply"}, {"worktree", "clean", "one", "--apply", "bad"}} {
		out := f.Bench(args...)
		out.RequireExit(2)
		if out.Stdout != wantClean || out.Stderr != "" {
			t.Fatalf("usage streams = stdout %q stderr %q", out.Stdout, out.Stderr)
		}
	}
	for _, args := range [][]string{{"worktree", "recovery"}, {"worktree", "recovery", "one", "two"}, {"worktree", "recovery", "ref", "--apply"}, {"worktree", "recovery", "ref", "--apply", "bad"}} {
		out := f.Bench(args...)
		out.RequireExit(2)
		if out.Stdout != wantRecovery || out.Stderr != "" {
			t.Fatalf("recovery usage streams = stdout %q stderr %q", out.Stdout, out.Stderr)
		}
	}
	if got := f.Git("worktree", "list", "--porcelain").Stdout; got != registrations {
		t.Fatalf("usage error discovered or mutated worktrees:\nbefore:\n%s\nafter:\n%s", registrations, got)
	}
	for _, name := range []string{branch, registeredBranch} {
		if !headExists(f, name) {
			t.Fatalf("usage error swept same-prefix branch %q", name)
		}
	}
	if got := contract.ReadFileAbs(t, filepath.Join(target, "dirty.txt")); got != "must survive\n" {
		t.Fatalf("usage error changed sibling bytes: %q", got)
	}
}

func testRuntimeWorktreeHostileSurfaces(t *testing.T) {
	f := onMainFixture(t)
	f.Bench("link").RequireExit(0)
	shim := t.TempDir()
	for _, tool := range []string{"bash", "basename", "dirname", "git", "readlink", "tr", "uname"} {
		source, err := exec.LookPath(tool)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(source, filepath.Join(shim, tool)); err != nil {
			t.Fatal(err)
		}
	}
	env := map[string]string{"PATH": shim, "BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
	realCLI := filepath.Join(contract.SubjectRoot(t), "dist", "bench")
	byPath := filepath.Join(f.Root, ".bench", "bin", "bench.sh")
	launcher := filepath.Join(t.TempDir(), "bench linked launcher")
	if err := os.Symlink(byPath, launcher); err != nil {
		t.Fatal(err)
	}
	deep := filepath.Join(f.Root, "deep", "cwd")
	contract.Mkdir(t, deep)

	hostile := filepath.Join(t.TempDir(), "space [glob]* $(pwn) Ünicode\nline")
	f.Git("worktree", "add", "-q", "-b", "hostile-combined", hostile, "HEAD")
	relative, err := filepath.Rel(deep, hostile)
	if err != nil {
		t.Fatal(err)
	}
	plan := contract.RunAt(t, f, deep, env, realCLI, "worktree", "clean", relative)
	plan.RequireExit(0)
	contract.RequireContains(t, plan.Stdout, strings.ReplaceAll(hostile, "\n", `\n`))
	apply := contract.RunAt(t, f, deep, env, realCLI, "worktree", "clean", relative, "--apply", cleanupFingerprint(t, plan.Stdout))
	apply.RequireExit(0)
	if _, err := os.Stat(hostile); !os.IsNotExist(err) {
		t.Fatalf("hostile exact target remains: %v", err)
	}

	leading := filepath.Join(f.Root, "-leading")
	f.Git("worktree", "add", "-q", "-b", "hostile-leading", leading, "HEAD")
	leadingPlan := contract.RunAt(t, f, f.Root, env, launcher, "worktree", "clean", "--", "-leading")
	leadingPlan.RequireExit(0)
	contract.RunAt(t, f, f.Root, env, launcher, "worktree", "clean", "--", "-leading", "--apply", cleanupFingerprint(t, leadingPlan.Stdout)).RequireExit(0)

	created := contract.RunAt(t, f, deep, env, byPath, "worktree", "create", "--request", "hostile-label", "--label", "quoted multiword label")
	created.RequireExit(0)
	contract.RequireContains(t, contract.ReadFileAbs(t, filepath.Join(f.Root, ".git", "bench-intent.json")), `"label":"quoted multiword label"`)

	unsafeTarget := filepath.Join(t.TempDir(), "unsafe\x1bpath")
	f.Git("worktree", "add", "-q", "-b", "hostile-control", unsafeTarget, "HEAD")
	unsafePlan := contract.RunAt(t, f, deep, env, realCLI, "worktree", "clean", unsafeTarget)
	unsafePlan.RequireExit(0)
	digest := sha256.Sum256([]byte(unsafeTarget))
	contract.RequireContains(t, unsafePlan.Stdout, fmt.Sprintf("sha256:%x", digest))
	contract.RequireContains(t, unsafePlan.Stdout, "retain")
	if strings.Contains(unsafePlan.Stdout+unsafePlan.Stderr, "\x1b") {
		t.Fatal("unsafe control byte reached cleanup output")
	}
	hook := filepath.Join(f.Root, ".bench", "hooks", "session-start.sh")
	session := contract.RunAt(t, f, deep, env, "/bin/bash", hook)
	session.RequireExit(0)
	if strings.Contains(session.Stdout+session.Stderr, "\x1b") {
		t.Fatal("unsafe control byte reached SessionStart output")
	}
	contract.RequireContains(t, session.Stdout, "bench not on PATH; invoke by path")
	if _, err := os.Stat(unsafeTarget); err != nil {
		t.Fatalf("SessionStart removed unsafe target: %v", err)
	}
}

func testRuntimeWorktreeInteractiveRelease(t *testing.T) {
	f := onMainFixture(t)
	f.WriteExecutable("wt-shell", `#!/usr/bin/env bash
: "${BENCH_WT_RECORD:?}"
pwd >> "$BENCH_WT_RECORD"
`)
	record := filepath.Join(f.Root, "paths")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "SHELL": filepath.Join(f.Root, "wt-shell")}
	f.BenchEnv(env, "worktree", "same objective").RequireExit(0)
	f.BenchEnv(env, "worktree", "same objective").RequireExit(0)
	paths := contract.NonEmptyLines(contract.ReadFileAbs(t, record))
	contract.RequireIntEqual(t, len(paths), 2, "worktree shell did not run twice")
	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("interactive exit did not release %q: %v", path, err)
		}
	}
	assignments, err := intent.Assignments(f.Root)
	if err != nil || len(assignments) != 0 {
		t.Fatalf("interactive releases left assignments = %#v, %v", assignments, err)
	}
}

func testRuntimeWorktreeLeaseHardening(t *testing.T) {
	f := onMainFixture(t)
	f.WriteExecutable("rec-shell", "#!/usr/bin/env bash\n: \"${BENCH_WT_RECORD:?}\"\npwd >> \"$BENCH_WT_RECORD\"\n")
	record := filepath.Join(f.Root, "paths")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "SHELL": filepath.Join(f.Root, "rec-shell")}
	created := f.BenchEnv(env, "worktree", "create", "--request", "lease-first", "--label", "first assignment")
	created.RequireExit(0)
	p := worktreeCreatePath(t, created.Stdout)
	lease := strings.TrimSpace(contract.RunAt(t, f, p, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, lease, fmt.Sprintf("%d 2026-01-01T00:00:00Z\n", os.Getpid()))
	f.BenchEnv(env, "worktree", "second assignment").RequireExit(0)
	paths := contract.NonEmptyLines(contract.ReadFileAbs(t, record))
	contract.RequireIntEqual(t, len(paths), 1, "expected one interactive worktree run")
	if p == paths[0] {
		t.Fatal("different requests shared an assignment registration")
	}
	if _, err := os.Stat(lease); err != nil {
		t.Fatal("foreign legacy lease was mutated by ownership creation")
	}
}

func testRuntimeWorktreeConcurrentAcquire(t *testing.T) {
	// Release is barriered, not detection: each shell records its worktree and then holds
	// until the test drops the go-file, which the test creates only after seeing both
	// records — so once both are recorded, neither run can exit before the other has
	// genuinely acquired. Overlap *detection* itself is not barriered; it is the
	// event-keyed 60s deadline below (overlapDeadline), armed only once the first record
	// lands so spawn latency under load never counts against it. That window is coupled to
	// the shell's own ~60s self-timeout (the `seq 600` / 0.1s poll loop, a leak backstop
	// that exits loud) — the two must not converge, or a slow second spawn could hit the
	// shell's timeout before the deadline gives it credit for overlapping, misreading a
	// schedule race as the by-design reuse the FT37 flake exercises.
	f := onMainFixture(t)
	f.WriteExecutable("rv-shell", `#!/usr/bin/env bash
: "${BENCH_WT_RECORD:?}" "${BENCH_WT_GO:?}"
pwd >> "$BENCH_WT_RECORD"
for _ in $(seq 600); do
  { [ -e "$BENCH_WT_GO.$(basename "$PWD")" ] || [ -e "$BENCH_WT_GO" ]; } && exit 0
  sleep 0.1
done
exit 1
`)
	record := filepath.Join(f.Root, "paths")
	goFile := filepath.Join(f.Root, "go-file")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "BENCH_WT_GO": goFile, "SHELL": filepath.Join(f.Root, "rv-shell")}
	done := make(chan contract.Probe, 2)
	for i := 0; i < 2; i++ {
		go func() { done <- f.BenchEnv(env, "worktree", "same concurrent objective") }()
	}
	backstop := time.Now().Add(2 * time.Minute)
	var overlapDeadline time.Time
	finished := make([]contract.Probe, 0, 2)
	for {
		raw, _ := os.ReadFile(record)
		lines := contract.NonEmptyLines(string(raw))
		if len(lines) >= 2 {
			break
		}
		if overlapDeadline.IsZero() && len(lines) >= 1 {
			// First record confirms one run has acquired and is holding — arm the overlap
			// window from this event rather than from spawn, so spawn latency under load
			// never counts against the assertion this window actually carries.
			overlapDeadline = time.Now().Add(time.Minute)
		}
	drainDone:
		for len(finished) < 2 {
			select {
			case p := <-done:
				finished = append(finished, p)
			default:
				break drainDone
			}
		}
		if len(finished) == 2 {
			// An exited run can never record (the shell records before it holds), but the
			// second record may land between the read above and the exits — re-read once.
			raw, _ := os.ReadFile(record)
			if len(contract.NonEmptyLines(string(raw))) < 2 {
				t.Fatal("both worktree runs exited with fewer than two acquires recorded — the runs never overlapped")
			}
			break
		}
		if overlapDeadline.IsZero() {
			if time.Now().After(backstop) {
				contract.WriteFileAbs(t, goFile, "") // release the straggler before failing
				t.Fatal("no bench worktree run recorded within 2m — acquire appears wedged")
			}
		} else if time.Now().After(overlapDeadline) {
			contract.WriteFileAbs(t, goFile, "") // release the straggler before failing
			t.Fatal("second acquire did not record within 60s of the first — the runs never overlapped")
		}
		time.Sleep(50 * time.Millisecond)
	}
	paths := contract.NonEmptyLines(contract.ReadFileAbs(t, record))
	sort.Strings(paths)
	contract.RequireIntEqual(t, len(paths), 2, "concurrent worktree runs did not both acquire")
	if paths[0] == paths[1] {
		contract.WriteFileAbs(t, goFile, "")
		<-done
		<-done
		t.Fatal("concurrent acquires shared a worktree")
	}
	if len(finished) != 0 {
		contract.WriteFileAbs(t, goFile, "")
		t.Fatal("a worktree run exited before both acquires overlapped")
	}
	contract.WriteFileAbs(t, goFile+"."+filepath.Base(paths[0]), "")
	first := <-done
	if first.ExitCode != 0 {
		ledger, readErr := intent.Read(f.Root)
		t.Fatalf("first exact release failed: stdout=%q stderr=%q ledger=%#v read=%v", first.Stdout, first.Stderr, ledger, readErr)
	}
	if _, err := os.Stat(paths[0]); !os.IsNotExist(err) {
		t.Fatalf("first exact release left %q: %v", paths[0], err)
	}
	if _, err := os.Stat(paths[1]); err != nil {
		t.Fatalf("first exact release disturbed live sibling %q: %v", paths[1], err)
	}
	assignments, err := intent.Assignments(f.Root)
	if err != nil || len(assignments) != 1 {
		t.Fatalf("first exact release assignments = %#v, %v", assignments, err)
	}
	contract.WriteFileAbs(t, goFile+"."+filepath.Base(paths[1]), "")
	(<-done).RequireExit(0)
	if _, err := os.Stat(paths[1]); !os.IsNotExist(err) {
		t.Fatalf("second exact release left %q: %v", paths[1], err)
	}
	assignments, err = intent.Assignments(f.Root)
	if err != nil || len(assignments) != 0 {
		t.Fatalf("interactive releases left assignments = %#v, %v", assignments, err)
	}
}
