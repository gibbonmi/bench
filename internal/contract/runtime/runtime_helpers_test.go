package runtime

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func benchPath(t testing.TB) string {
	t.Helper()
	return filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
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
	contract.WriteFileAbs(t, lease, "")
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
	merged := make(map[string]string, len(os.Environ())+len(f.Env)+len(extra))
	for _, entry := range os.Environ() {
		if key, value, ok := strings.Cut(entry, "="); ok {
			merged[key] = value
		}
	}
	for key, value := range f.Env {
		merged[key] = value
	}
	for key, value := range extra {
		merged[key] = value
	}
	env := make([]string, 0, len(merged))
	for key, value := range merged {
		env = append(env, key+"="+value)
	}
	return env
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
