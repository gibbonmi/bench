package worktree

import (
	"errors"
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
		t.Skip("cksum not available")
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

func TestCleanupDeletesOnlyExactBranchAndRetiresLastRecoveryRef(t *testing.T) {
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
	t.Run("recovered context leaves after last exact ref", func(t *testing.T) {
		root, target := newOwnedAssignment(t, "terminal-recovered")
		sibling, err := Create(root, "terminal-recovered-sibling", "sibling", nil)
		if err != nil {
			t.Fatal(err)
		}
		siblingRef := "refs/bench/recovery/" + sibling.Assignment.OwnerID + "/" + sibling.Assignment.ID + "/1"
		gitRun(t, root, "update-ref", siblingRef, target.Assignment.Start)
		mustWrite(t, filepath.Join(target.Path, "recovered.txt"), []byte("recovered\n"), 0o644)
		markPending(t, root, target.Assignment)
		if _, err := ApplyAutomatic(root, target.Path, nil); err != nil {
			t.Fatal(err)
		}
		recovered, err := assignmentByID(root, target.Assignment.ID)
		if err != nil || recovered.State != intent.StateRecovered || len(recovered.Recovery) != 1 {
			t.Fatalf("recovered assignment = %#v, %v", recovered, err)
		}
		first := recovered.Recovery[0]
		second := first
		second.Ref = strings.TrimSuffix(first.Ref, "/1") + "/2"
		gitRun(t, root, "update-ref", second.Ref, second.Root)
		recovered.Recovery = append(recovered.Recovery, second)
		if err := intent.PutAssignment(root, recovered); err != nil {
			t.Fatal(err)
		}
		for _, payload := range first.Payloads {
			gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "cherry-pick", payload)
		}
		if err := RetireRecovery(root, first.Ref); err != nil {
			t.Fatal(err)
		}
		if exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", first.Ref).Run() == nil {
			t.Fatal("first exact recovery ref survived retirement")
		}
		gitRun(t, root, "show-ref", "--verify", "--quiet", second.Ref)
		if current, err := assignmentByID(root, target.Assignment.ID); err != nil || current.State != intent.StateRecovered || len(current.Recovery) != 1 {
			t.Fatalf("intermediate recovered state = %#v, %v", current, err)
		}
		if err := RetireRecovery(root, second.Ref); err != nil {
			t.Fatal(err)
		}
		if _, err := assignmentByID(root, target.Assignment.ID); err == nil {
			t.Fatal("last-ref retirement did not compact recovered assignment")
		}
		gitRun(t, root, "show-ref", "--verify", "--quiet", sibling.Assignment.Branch)
		gitRun(t, root, "show-ref", "--verify", "--quiet", siblingRef)
	})
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
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
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
	root := t.TempDir()
	gitRun(t, root, "init", "-q", "-b", "main")
	gitRun(t, root, "config", "user.email", "bench@local")
	gitRun(t, root, "config", "user.name", "bench")
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
