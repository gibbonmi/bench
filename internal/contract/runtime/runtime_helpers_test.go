package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/contract/internal/freshnessfixture"
	"github.com/gibbonmi/bench/internal/freshness"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/runbinary"
)

const runtimeLauncherFreshnessChild = "BENCH_FT131_RUNTIME_LAUNCHER_CHILD"

func benchPath(t testing.TB) string {
	t.Helper()
	selection := contract.SelectedBench(t)
	return contract.SelectedBenchPath(t, selection)
}

func benchWrapperPath(t testing.TB) string {
	t.Helper()
	return filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
}

func selectedBenchEnv(t testing.TB, extra map[string]string) map[string]string {
	t.Helper()
	selection := contract.SelectedBench(t)
	env := make(map[string]string, len(extra)+2)
	for key, value := range extra {
		env[key] = value
	}
	env["BENCH_KIT"] = selection.SourceRoot
	env[runbinary.Env] = selection.Path
	return env
}

func selectedSurfaceEnv(t testing.TB, f contract.Fixture, extra map[string]string) []string {
	t.Helper()
	return contract.ProcessEnv(f.Env, selectedBenchEnv(t, extra))
}

func TestRuntimeLauncherRefusesStaleBeforeFalseGreenAndFalseRedAssertions(t *testing.T) {
	for _, assertion := range []string{"false-green", "false-red"} {
		t.Run(assertion, func(t *testing.T) {
			if os.Getenv(runtimeLauncherFreshnessChild) == assertion {
				root := contract.SubjectRoot(t)
				if err := freshness.Verify(root, filepath.Join(root, "dist", "bench")); err != nil {
					t.Fatal(err)
				}
				command := exec.Command("bash", filepath.Join(root, "bin", "bench.sh"), "version")
				output, err := command.CombinedOutput()
				if assertion == "false-green" && err == nil {
					return
				}
				t.Fatalf("correct runtime assertion rejected: %v\n%s", err, output)
			}

			root := freshnessfixture.StaleSubject(t, "old runtime assertion output")
			command := exec.Command(os.Args[0], "-test.run=^TestRuntimeLauncherRefusesStaleBeforeFalseGreenAndFalseRedAssertions$/"+assertion+"$")
			command.Env = append(os.Environ(), runtimeLauncherFreshnessChild+"="+assertion, canary.SubjectRootEnv+"="+root)
			output, err := command.CombinedOutput()
			if err == nil {
				t.Fatalf("stale runtime subject satisfied %s assertion:\n%s", assertion, output)
			}
			for _, want := range []string{"bench binary ", "rebuild with "} {
				if !strings.Contains(string(output), want) {
					t.Fatalf("stale runtime subject did not report %q:\n%s", want, output)
				}
			}
			for _, forbidden := range []string{"old runtime assertion output", "correct runtime assertion rejected"} {
				if strings.Contains(string(output), forbidden) {
					t.Fatalf("stale runtime subject reached assertion output %q:\n%s", forbidden, output)
				}
			}
		})
	}
}

func TestRuntimeBenchPathConsumerEnumeration(t *testing.T) {
	expected := []string{
		"runtime_agent_entry_test.go",
		"runtime_commit_test.go",
		"runtime_gate_action_proof_test.go",
		"runtime_gate_owner_helper_test.go",
		"runtime_gate_owner_test.go",
		"runtime_gate_partial_proof_test.go",
		"runtime_gate_proof_helpers_test.go",
		"runtime_gate_shift_proof_test.go",
		"runtime_gate_test.go",
		"runtime_shift_adapters_test.go",
		"runtime_spec_history_test.go",
		"runtime_spec_test.go",
		"runtime_stale_gate_test.go",
		"runtime_status_test.go",
		"runtime_testreport_test.go",
		"runtime_worktree_test.go",
	}
	files, err := filepath.Glob(filepath.Join(contract.KitRoot(t), "internal", "contract", "runtime", "runtime_*_test.go"))
	if err != nil {
		t.Fatal(err)
	}
	needle := "benchPath" + "(t)"
	var consumers []string
	count := 0
	for _, path := range files {
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatal(readErr)
		}
		occurrences := strings.Count(string(data), needle)
		if occurrences > 0 {
			consumers = append(consumers, filepath.Base(path))
			count += occurrences
		}
	}
	if got, want := strings.Join(consumers, "\n"), strings.Join(expected, "\n"); got != want {
		t.Fatalf("benchPath consumer files:\n%s\nwant:\n%s", got, want)
	}
	if count != 28 {
		t.Fatalf("benchPath occurrences = %d, want 28", count)
	}
}

func gitDir(t testing.TB, f contract.Fixture) string {
	t.Helper()
	return strings.TrimSpace(f.Git("rev-parse", "--absolute-git-dir").Stdout)
}

func commitAllowEmpty(t testing.TB, f contract.Fixture, message string) {
	t.Helper()
	f.Run("git", "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", message).RequireExit(0)
}

func cksum(t testing.TB, value string) string {
	t.Helper()
	cmd := exec.Command("cksum")
	cmd.Stdin = strings.NewReader(value + "\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("cksum %q: %v", value, err)
	}
	return strings.Fields(string(out))[0]
}

type runtimePoolWorktrees struct {
	Pool, Warm, Leased, LeaseFile string
}

func addRuntimePoolWorktrees(t testing.TB, f contract.Fixture, benchHome string) runtimePoolWorktrees {
	t.Helper()
	repoRoot := strings.TrimSpace(f.Git("rev-parse", "--show-toplevel").Stdout)
	pool := filepath.Join(benchHome, "worktrees", filepath.Base(repoRoot)+"-"+cksum(t, repoRoot))
	warm := filepath.Join(pool, "warm")
	leased := filepath.Join(pool, "leased")
	contract.Mkdir(t, pool)
	f.Git("worktree", "add", "-q", "--detach", warm, "HEAD")
	f.Git("worktree", "add", "-q", "--detach", leased, "HEAD")
	lease := strings.TrimSpace(contract.RunAt(t, f, leased, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, lease, fmt.Sprintf("%d 2026-01-01T00:00:00Z\n", os.Getpid()))
	return runtimePoolWorktrees{Pool: pool, Warm: warm, Leased: leased, LeaseFile: lease}
}

func repeatLines(n int, line string) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(line)
	}
	return b.String()
}

func copyRuntimeFile(t testing.TB, src, dst string, mode os.FileMode) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, data, mode); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
}

type normalizedLifecycleState struct {
	MarkerSchema, MarkerPath, AssignmentPath, OwnerJoin, BranchJoin, LockJoin bool
}

type normalizedCleanupState struct {
	PathGone, RegistrationGone, BranchGone, AssignmentGone bool
}

func surfaceEnv(f contract.Fixture, extra map[string]string) []string {
	return contract.ProcessEnv(f.Env, extra)
}

func waitForSurfacePath(t *testing.T, ready string, cmd *exec.Cmd) string {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(ready); err == nil && strings.TrimSpace(string(data)) != "" {
			return strings.TrimSpace(string(data))
		}
		time.Sleep(10 * time.Millisecond)
	}
	_ = cmd.Process.Kill()
	t.Fatal("interactive shell never exposed its lifecycle state")
	return ""
}

func worktreeCreatePath(t *testing.T, output string) string {
	t.Helper()
	lines := contract.NonEmptyLines(output)
	if len(lines) != 2 {
		t.Fatalf("worktree create output = %q", output)
	}
	return strings.Trim(strings.TrimSpace(strings.Split(lines[1], ",")[0]), `"`)
}

func inspectLifecycleState(t *testing.T, f contract.Fixture, path string) (normalizedLifecycleState, intent.Assignment) {
	t.Helper()
	admin := strings.TrimSpace(contract.RunAt(t, f, path, nil, "git", "rev-parse", "--path-format=absolute", "--git-dir").Stdout)
	var marker struct {
		Schema  string `json:"schema"`
		OwnerID string `json:"owner_id"`
		Path    string `json:"path"`
	}
	if err := json.Unmarshal([]byte(contract.ReadFileAbs(t, filepath.Join(admin, "bench-owner"))), &marker); err != nil {
		t.Fatal(err)
	}
	assignments, err := intent.Assignments(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	var assignment intent.Assignment
	for _, candidate := range assignments {
		if candidate.Worktree == path {
			assignment = candidate
			break
		}
	}
	branch := strings.TrimSpace(contract.RunAt(t, f, path, nil, "git", "symbolic-ref", "HEAD").Stdout)
	registrations, err := benchgit.Worktrees(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	lock := ""
	for _, registration := range registrations {
		if registration.Path == path {
			lock = registration.LockReason
		}
	}
	return normalizedLifecycleState{
		MarkerSchema:   marker.Schema == "bench-owner/v1",
		MarkerPath:     marker.Path == path,
		AssignmentPath: assignment.Worktree == path,
		OwnerJoin:      marker.OwnerID != "" && marker.OwnerID == assignment.OwnerID,
		BranchJoin:     strings.HasPrefix(branch, "refs/heads/bench/assign/") && branch == assignment.Branch,
		LockJoin:       strings.Contains(lock, "owner="+assignment.OwnerID) && strings.Contains(lock, "assignment="+assignment.ID),
	}, assignment
}

func inspectCleanupState(t *testing.T, f contract.Fixture, assignment intent.Assignment) normalizedCleanupState {
	t.Helper()
	_, pathErr := os.Stat(assignment.Worktree)
	registrations, err := benchgit.Worktrees(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	registered := false
	for _, registration := range registrations {
		registered = registered || registration.Path == assignment.Worktree
	}
	assignments, err := intent.Assignments(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	persisted := false
	for _, candidate := range assignments {
		persisted = persisted || candidate.ID == assignment.ID
	}
	return normalizedCleanupState{
		PathGone:         os.IsNotExist(pathErr),
		RegistrationGone: !registered,
		BranchGone:       !benchgit.OK("-C", f.Root, "show-ref", "--verify", "--quiet", assignment.Branch),
		AssignmentGone:   !persisted,
	}
}

type configuredClaudeWorktreeHandler struct {
	Command string
	Args    []string
}

func configuredClaudeWorktreeCommand(t *testing.T, f contract.Fixture, event string) configuredClaudeWorktreeHandler {
	t.Helper()
	var cfg struct {
		Hooks map[string][]map[string]json.RawMessage `json:"hooks"`
	}
	if err := json.Unmarshal([]byte(f.ReadFile(".claude/settings.json")), &cfg); err != nil {
		t.Fatal(err)
	}
	groups := cfg.Hooks[event]
	if len(groups) != 1 {
		t.Fatalf("claude %s configured command missing: groups = %d, want 1", event, len(groups))
	}
	if _, present := groups[0]["matcher"]; present {
		t.Fatalf("claude %s group carries ignored matcher; lifecycle events must be matcher-free", event)
	}
	var hooks []struct {
		Type, Command string
		Args          []string
	}
	if err := json.Unmarshal(groups[0]["hooks"], &hooks); err != nil {
		t.Fatal(err)
	}
	if len(hooks) != 1 || hooks[0].Type != "command" {
		t.Fatalf("claude %s command hooks = %#v, want one command", event, hooks)
	}
	action := map[string]string{"WorktreeCreate": "create", "WorktreeRemove": "remove"}[event]
	wantCommand := "${CLAUDE_PROJECT_DIR}/.bench/hooks/worktree-lifecycle.sh"
	if hooks[0].Command != wantCommand || len(hooks[0].Args) != 1 || hooks[0].Args[0] != action {
		t.Fatalf("claude %s handler = command %q args %#v, want exec form %q [%q]", event, hooks[0].Command, hooks[0].Args, wantCommand, action)
	}
	return configuredClaudeWorktreeHandler{
		Command: strings.Replace(wantCommand, "${CLAUDE_PROJECT_DIR}", f.Root, 1),
		Args:    hooks[0].Args,
	}
}

func mustJSON(t *testing.T, value string) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func registrationLockReason(t *testing.T, f contract.Fixture, path string) string {
	t.Helper()
	registrations, err := benchgit.Worktrees(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	for _, registration := range registrations {
		if registration.Path == path {
			return registration.LockReason
		}
	}
	return ""
}

func countStatusRows(out string) int {
	n := 0
	for _, line := range strings.Split(out, "\n") {
		if regexp.MustCompile(`^  [a-z]`).MatchString(line) {
			n++
		}
	}
	return n
}
