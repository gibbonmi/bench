package worktree

// This file is the serial journey harness: the one place a worktree test may
// create a disposable repository, bind environment, change directory, or start
// a descendant process. It records every repository and descendant start, and
// its source census fails the package when a test file outside the harness
// starts either effect itself. The selected Bench executable lives with the
// harness in test_run_test.go; the package environment root lives with it in
// main_test.go. (Coverage rows EI1, RJ1, FA1.)

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// journeyEffects is the harness effect log: one entry per repository or
// descendant start, carrying the effect kind and the starting test's name.
var journeyEffects struct {
	mu      sync.Mutex
	records []string
}

func recordJourneyEffect(t testing.TB, kind, detail string) {
	t.Helper()
	journeyEffects.mu.Lock()
	defer journeyEffects.mu.Unlock()
	journeyEffects.records = append(journeyEffects.records, kind+" "+t.Name()+" "+detail)
}

// journeyEffectLog returns a copy of the recorded starts.
func journeyEffectLog() []string {
	journeyEffects.mu.Lock()
	defer journeyEffects.mu.Unlock()
	return append([]string(nil), journeyEffects.records...)
}

// --- disposable repositories ---

func journeyRepo(t testing.TB) string {
	t.Helper()
	recordJourneyEffect(t, "repository", "bare-default-branch")
	return gittest.Repo(t)
}

func journeyRepoOnBranch(t testing.TB, branch string) string {
	t.Helper()
	recordJourneyEffect(t, "repository", "branch "+branch)
	return gittest.RepoOnBranch(t, branch)
}

// journeyStubGit installs the file-backed git stub on the process PATH. It
// binds the PATH through bindEnv, not through gittest.StubGit, so the effect
// log records the environment bind and the parallel census sees the edge: every
// caller is serial by construction. A caller that only needs the stub for one
// child process wants gittest.StubGitDir through its own child environment.
func journeyStubGit(t testing.TB, root, mode, logPath string) string {
	t.Helper()
	recordJourneyEffect(t, "descendant-stub", "git mode "+mode)
	dir, commonDir := gittest.StubGitDir(t, root, mode, logPath)
	bindEnv(t, "PATH", dir)
	return commonDir
}

// journeyStubGitStarts counts the starts that a journeyStubGit stub logged at logPath with
// flag among their arguments.
func journeyStubGitStarts(t testing.TB, logPath, flag string) int {
	t.Helper()
	return gittest.StubGitStarts(t, logPath, flag)
}

func journeyFIFOWorktreeAdmin(t testing.TB, root, id string) string {
	t.Helper()
	recordJourneyEffect(t, "repository", "fifo worktree admin "+id)
	return gittest.FIFOWorktreeAdmin(t, root, id)
}

// newWorktreeRepo is the standard disposable repository: one commit on main
// carrying tracked.txt and README.md.
func newWorktreeRepo(t testing.TB) string {
	t.Helper()
	root := journeyRepoOnBranch(t, "main")
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

// nestedRepositoryDir and nestedRepositoryFile name the repository that
// plantNestedRepository builds and the one file that its commit tracks.
const nestedRepositoryDir, nestedRepositoryFile = "nested", "nested.txt"

// plantNestedRepository builds a separate repository with one commit under parent and
// returns its path. The parent repository does not track it, so it is an embedded
// repository, not a submodule. A caller that writes nestedRepositoryFile makes it dirty.
func plantNestedRepository(t *testing.T, parent string) string {
	t.Helper()
	nested := filepath.Join(parent, nestedRepositoryDir)
	mustMkdirAll(t, nested, 0o755)
	gitRun(t, nested, "init", "-q", "-b", "main")
	commitInWorktree(t, nested, nestedRepositoryFile, "base\n", "nested")
	return nested
}

// --- explicit environment and directories ---

func bindEnv(t testing.TB, key, value string) {
	t.Helper()
	recordJourneyEffect(t, "environment", key)
	t.Setenv(key, value)
}

// journeyChildEnv is the environment one journey child runs under: the private
// home the caller owns, and the values git needs from the process environment.
// A child environment leaves the process environment alone, so the journey that
// builds one stays parallel-eligible.
func journeyChildEnv(t testing.TB, home string) []string {
	t.Helper()
	recordJourneyEffect(t, "child-environment", home)
	env := []string{"BENCH_HOME=" + home, "PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}
	if temp := os.Getenv("TMPDIR"); temp != "" {
		env = append(env, "TMPDIR="+temp)
	}
	return env
}

func chdir(t testing.TB, dir string) {
	t.Helper()
	recordJourneyEffect(t, "directory", dir)
	t.Chdir(dir)
}

// --- descendant processes ---

// descendant prepares one child process start. The caller sets Dir, Env, and
// streams, then runs it; the harness records the start and kills a child the
// test left running.
func descendant(t testing.TB, name string, args ...string) *exec.Cmd {
	t.Helper()
	recordJourneyEffect(t, "descendant", name)
	cmd := exec.Command(name, args...)
	t.Cleanup(func() {
		if cmd.Process != nil && (cmd.ProcessState == nil || !cmd.ProcessState.Exited()) {
			_ = cmd.Process.Kill()
		}
	})
	return cmd
}

func gitOutput(t testing.TB, dir string, args ...string) string {
	t.Helper()
	cmd := descendant(t, "git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

func gitRun(t testing.TB, dir string, args ...string) {
	t.Helper()
	cmd := descendant(t, "git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// withIdentity prefixes one git argument list with the shared default fixture identity,
// so a gitRun or gitOutput child that commits or merges records that identity. It is the
// one route by which a worktree test outside this harness reaches the gittest identity;
// the child still starts through descendant.
func withIdentity(args ...string) []string {
	return append(gittest.IdentityArgs(), args...)
}

// --- harness census (coverage row EI1) ---

// journeyHarnessFiles are the files that may start repositories or descendants
// and mutate environment or directory: this harness, the package environment
// root, and the selected-executable owner.
var journeyHarnessFiles = map[string]bool{
	"journey_test.go":  true,
	"main_test.go":     true,
	"test_run_test.go": true,
}

// outsideHarnessEffect matches a repository start, a descendant start, an
// environment mutation, or a directory mutation at a test call site.
var outsideHarnessEffect = regexp.MustCompile(`\b(exec\.Command|exec\.CommandContext|os\.StartProcess|gittest\.|os\.Setenv\(|os\.Clearenv\(|os\.Chdir\(|[a-zA-Z_]\.Setenv\(|[a-zA-Z_]\.Chdir\()`)

// TestSerialJourneyHarnessCensus fails when any worktree test file outside the
// harness starts a repository or a descendant, or mutates environment or the
// current directory. Those effects route through the harness, which records
// them. (Coverage row EI1.)
func TestSerialJourneyHarnessCensus(t *testing.T) {
	t.Parallel()
	// The walk is parseTestFiles, which the parallel census also uses, so the
	// two censuses see the same file set and the same special-file skip.
	_, _, names, err := parseTestFiles(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if journeyHarnessFiles[name] {
			continue
		}
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			code := line
			if idx := strings.Index(code, "//"); idx >= 0 {
				code = code[:idx]
			}
			if match := outsideHarnessEffect.FindString(code); match != "" {
				t.Errorf("%s:%d: effect %s outside the serial journey harness; route it through journey_test.go", name, i+1, match)
			}
		}
	}
}

// TestJourneyHarnessRecordsStarts proves the harness records a repository start
// and a descendant start under the starting test's name. (Coverage row EI1.)
func TestJourneyHarnessRecordsStarts(t *testing.T) {
	t.Parallel()
	before := len(journeyEffectLog())
	newWorktreeRepo(t)
	log := journeyEffectLog()[before:]
	var repo, child bool
	for _, record := range log {
		if strings.HasPrefix(record, "repository "+t.Name()+" ") {
			repo = true
		}
		if strings.HasPrefix(record, "descendant "+t.Name()+" ") {
			child = true
		}
	}
	if !repo || !child {
		t.Fatalf("harness log %q lacks a repository and a descendant start for %s", log, t.Name())
	}
}

// --- location reads through the run value (coverage rows LGP16, LGP17) ---

// TestOwnerMarkerReadsTheCommonDirOnce proves one owner-marker validation under a
// canonical primary root starts the common-directory query once. The fail-git-path stub
// passes that query through to the real git and logs each start. It also proves that two
// verb runs build two run values, so no answer that a run caches crosses verb runs.
// (Coverage rows LGP17, GQ37.)
func TestOwnerMarkerReadsTheCommonDirOnce(t *testing.T) {
	f := newOwnedAssignment(t, "owner-marker-common-dir")
	root, logPath := canonicalRoot(f.root), filepath.Join(t.TempDir(), "argv")
	journeyStubGit(t, root, "fail-git-path", logPath)
	if _, err := validateOwnerMarkerIn(&git.Run{}, root, f.creation.Path); err != nil {
		t.Fatalf("validateOwnerMarkerIn(%q, %q) = %v, want the owner evidence", root, f.creation.Path, err)
	}
	if starts := journeyStubGitStarts(t, logPath, "--git-common-dir"); starts != 1 {
		t.Fatalf("one owner-marker validation started the common-directory query %d times, want 1", starts)
	}
	first, second := callAmbient(verbCall{}, io.Discard), callAmbient(verbCall{}, io.Discard)
	if first.run == second.run {
		t.Fatalf("two verb runs share the run value %p, want a fresh run value for each verb run", first.run)
	}
}

// TestExplicitPlanReadsEachLocationOnce proves one explicit plan for an owned target under
// a canonical primary root starts the common-directory query once and the target's
// administration-directory query once, because its location reads and its two ledger
// reads all go through its one run. The fail-git-path stub refuses the index path query,
// so the plan stops after its owner-marker and ledger reads. (Coverage rows LGP16, GQ35,
// GQ36.)
func TestExplicitPlanReadsEachLocationOnce(t *testing.T) {
	f := newOwnedAssignment(t, "explicit-plan-common-dir")
	root, logPath := canonicalRoot(f.root), filepath.Join(t.TempDir(), "argv")
	journeyStubGit(t, root, "fail-git-path", logPath)
	if _, err := planExplicitWith(defaultJoins(), &git.Run{}, root, f.creation.Path, CleanupOptions{}); err == nil {
		t.Fatalf("planExplicitWith(%q, %q) = nil error, want the refused index path query", root, f.creation.Path)
	}
	if starts := journeyStubGitStarts(t, logPath, "--git-common-dir"); starts != 1 {
		t.Fatalf("one explicit plan started the common-directory query %d times, want 1", starts)
	}
	adminQuery := "-C " + f.creation.Path + " rev-parse --path-format=absolute --git-dir"
	if starts := journeyStubGitStarts(t, logPath, adminQuery); starts != 1 {
		t.Fatalf("one explicit plan started the target administration-directory query %d times, want 1", starts)
	}
}

// --- proof inventory (coverage rows RJ1, FA1) ---

// proofRequirement names the policy owner and journey class one required proof
// identifier belongs to.
type proofRequirement struct{ owner, class string }

// requiredProofs is the retained-coverage inventory: the representative serial
// journeys per policy owner, and one fact-adapter proof per owner. TestMain
// fails an unfiltered package run when any identifier was never marked.
var requiredProofs = map[string]proofRequirement{
	"landing/journey/publish-release":    {"landing", "journey"},
	"landing/journey/conflict-refusal":   {"landing", "journey"},
	"landing/journey/interrupted-resume": {"landing", "journey"},
	"landing/journey/hostile-residue":    {"landing", "journey"},
	"lifecycle/journey/create-remove":    {"lifecycle", "journey"},
	"lifecycle/journey/registration":     {"lifecycle", "journey"},
	"lifecycle/journey/lock":             {"lifecycle", "journey"},
	"lifecycle/journey/recovery":         {"lifecycle", "journey"},
	"reclaim/journey/lease":              {"reclaim", "journey"},
	"reclaim/journey/registration":       {"reclaim", "journey"},
	"reclaim/journey/process-liveness":   {"reclaim", "journey"},
	"reclaim/journey/deletion":           {"reclaim", "journey"},
	"landing/adapter/facts":              {"landing", "adapter"},
	"lifecycle/adapter/facts":            {"lifecycle", "adapter"},
	"reclaim/adapter/facts":              {"reclaim", "adapter"},
}

var proofMarks struct {
	mu   sync.Mutex
	done map[string]bool
}

// markProof records one required proof as observed. A test calls it only after
// its named observation held; a test that already failed marks nothing.
func markProof(t testing.TB, id string) {
	t.Helper()
	if _, ok := requiredProofs[id]; !ok {
		t.Fatalf("markProof: %q is not a required proof identifier", id)
	}
	if t.Failed() {
		return
	}
	proofMarks.mu.Lock()
	defer proofMarks.mu.Unlock()
	if proofMarks.done == nil {
		proofMarks.done = map[string]bool{}
	}
	proofMarks.done[id] = true
}

// missingProofs lists every required identifier that was never marked, each
// line naming the identifier, its owner, and its class.
func missingProofs() []string {
	proofMarks.mu.Lock()
	defer proofMarks.mu.Unlock()
	var missing []string
	for id, req := range requiredProofs {
		if !proofMarks.done[id] {
			missing = append(missing, "missing required proof "+id+" owner="+req.owner+" class="+req.class)
		}
	}
	sort.Strings(missing)
	return missing
}
