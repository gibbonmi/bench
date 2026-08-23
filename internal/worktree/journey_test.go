package worktree

// This file is the serial journey harness: the one place a worktree test may
// create a disposable repository, bind environment, change directory, or start
// a descendant process. It records every repository and descendant start, and
// its source census fails the package when a test file outside the harness
// starts either effect itself. The selected Bench executable lives with the
// harness in test_run_test.go; the package environment root lives with it in
// main_test.go. (Coverage rows EI1, RJ1, FA1.)

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

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

func journeyStubGit(t testing.TB, root, mode, logPath string) string {
	t.Helper()
	recordJourneyEffect(t, "descendant-stub", "git mode "+mode)
	return gittest.StubGit(t, root, mode, logPath)
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

// --- explicit environment and directories ---

func bindEnv(t testing.TB, key, value string) {
	t.Helper()
	recordJourneyEffect(t, "environment", key)
	t.Setenv(key, value)
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
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, "_test.go") || journeyHarnessFiles[name] {
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
