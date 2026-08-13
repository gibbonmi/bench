//go:build system

package systemtest

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
)

type executableIdentity struct {
	path   string
	inode  uint64
	digest [sha256.Size]byte
}

type processResult struct {
	stdout, stderr string
	code           int
}

type systemOwner struct {
	selected executableIdentity
	kit      string
	root     string
	repos    []string
	mu       sync.Mutex
	starts   int
	seen     []executableIdentity
	terminal map[string]bool
}

var owner *systemOwner

func TestMain(m *testing.M) {
	if pidFile := os.Getenv("BENCH_SYSTEM_INTERRUPT_CHILD"); pidFile != "" {
		os.Exit(runInterruptChild(pidFile))
	}
	created, err := newSystemOwner()
	if err != nil {
		fmt.Fprintln(os.Stderr, "system owner:", err)
		os.Exit(1)
	}
	owner = created
	code := m.Run()
	if err := owner.verify(); err != nil {
		fmt.Fprintln(os.Stderr, "system owner verification:", err)
		code = 1
	}
	if err := owner.cleanup(); err != nil {
		fmt.Fprintln(os.Stderr, "system owner cleanup:", err)
		code = 1
	}
	os.Exit(code)
}

func newSystemOwner() (*systemOwner, error) {
	selected, err := identifyExecutable(os.Getenv("BENCH_RUN_BINARY"))
	if err != nil {
		return nil, err
	}
	kit, err := filepath.Abs(os.Getenv("BENCH_KIT"))
	if err != nil || kit == "" {
		return nil, errors.New("BENCH_KIT must name an absolute kit root")
	}
	root, err := os.MkdirTemp("", "bench system [owner]-")
	if err != nil {
		return nil, err
	}
	o := &systemOwner{selected: selected, kit: kit, root: root, terminal: map[string]bool{}}
	for range 3 {
		repo, err := os.MkdirTemp(root, "repository [journey]-")
		if err != nil {
			_ = o.cleanup()
			return nil, err
		}
		if result := o.runAt(repo, nil, "git", "init", "-q"); result.code != 0 {
			_ = o.cleanup()
			return nil, fmt.Errorf("git init: %s", result.stderr)
		}
		o.repos = append(o.repos, repo)
	}
	return o, nil
}

func identifyExecutable(path string) (executableIdentity, error) {
	cleaned, err := filepath.Abs(path)
	if err != nil || path == "" {
		return executableIdentity{}, errors.New("BENCH_RUN_BINARY must name an absolute selected executable")
	}
	info, err := os.Stat(cleaned)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		return executableIdentity{}, errors.New("BENCH_RUN_BINARY is not an executable regular file")
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return executableIdentity{}, errors.New("selected executable inode is unavailable")
	}
	data, err := os.ReadFile(cleaned)
	if err != nil {
		return executableIdentity{}, err
	}
	return executableIdentity{path: cleaned, inode: stat.Ino, digest: sha256.Sum256(data)}, nil
}

func (o *systemOwner) runSelected(repo string, args ...string) processResult {
	if err := o.observeSelected(); err != nil {
		return processResult{code: 1, stderr: "selected executable identity changed"}
	}
	env := []string{"BENCH_RUN_BINARY=" + o.selected.path, "BENCH_KIT=" + o.kit, "BENCH_SYSTEM_ROOT=" + repo, "BENCH_COMMAND_OBSERVE=1"}
	return o.runAt(repo, env, o.selected.path, args...)
}

func (o *systemOwner) runWrapper(repo, wrapper string, args ...string) processResult {
	if err := o.observeSelected(); err != nil {
		return processResult{code: 1, stderr: "selected executable identity changed"}
	}
	env := []string{"BENCH_RUN_BINARY=" + o.selected.path, "BENCH_KIT=" + filepath.Join(repo, "missing-kit"), "BENCH_COMMAND_OBSERVE=1"}
	return o.runAt(repo, env, "bash", append([]string{wrapper}, args...)...)
}

func (o *systemOwner) observeSelected() error {
	identity, err := identifyExecutable(o.selected.path)
	if err != nil || identity != o.selected {
		return errors.New("selected executable identity changed")
	}
	o.mu.Lock()
	o.seen = append(o.seen, identity)
	o.mu.Unlock()
	return nil
}

func (o *systemOwner) runAt(dir string, overrides []string, program string, args ...string) processResult {
	o.mu.Lock()
	o.starts++
	o.mu.Unlock()
	cmd := exec.Command(program, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnvironment(os.Environ(), overrides)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		code = 1
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		}
	}
	return processResult{stdout: stdout.String(), stderr: stderr.String(), code: code}
}

func (o *systemOwner) timeoutProcessGroup(timeout time.Duration) (int, error) {
	o.mu.Lock()
	o.starts++
	o.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.Command("sh", "-c", "sleep 30 & echo $!; wait")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return 0, err
	}
	<-ctx.Done()
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	_ = cmd.Wait()
	pid, err := strconv.Atoi(strings.TrimSpace(stdout.String()))
	if err != nil {
		return 0, fmt.Errorf("descendant pid: %w", err)
	}
	return pid, nil
}

func (o *systemOwner) interruptProcessGroup() (int, error) {
	o.mu.Lock()
	o.starts++
	o.mu.Unlock()
	pidFile := filepath.Join(o.root, "interrupt-descendant.pid")
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = mergeEnvironment(os.Environ(), []string{"BENCH_SYSTEM_INTERRUPT_CHILD=" + pidFile})
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return 0, err
	}
	var data []byte
	deadline := time.Now().Add(2 * time.Second)
	for len(data) == 0 && time.Now().Before(deadline) {
		data, _ = os.ReadFile(pidFile)
		if len(data) == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
		_ = cmd.Wait()
		return 0, fmt.Errorf("interrupt descendant pid: %w", err)
	}
	if err := syscall.Kill(-pgid, syscall.SIGINT); err != nil {
		return 0, err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
		<-done
		return 0, errors.New("process group ignored interrupt")
	}
	return pid, nil
}

func runInterruptChild(pidFile string) int {
	descendant := exec.Command("sleep", "30")
	if err := descendant.Start(); err != nil {
		return 1
	}
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(descendant.Process.Pid)), 0o600); err != nil {
		_ = descendant.Process.Kill()
		_ = descendant.Wait()
		return 1
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)
	<-interrupt
	_ = descendant.Process.Signal(os.Interrupt)
	_ = descendant.Wait()
	return 130
}

func (o *systemOwner) markTerminal(outcome string) {
	o.mu.Lock()
	o.terminal[outcome] = true
	o.mu.Unlock()
}

func (o *systemOwner) verify() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.repos) != 3 {
		return fmt.Errorf("repository count = %d, want 3", len(o.repos))
	}
	if o.starts == 0 {
		return errors.New("no owned process starts recorded")
	}
	if len(o.seen) == 0 {
		return errors.New("no selected executable observations recorded")
	}
	for _, identity := range o.seen {
		if identity != o.selected {
			return errors.New("selected executable identity ledger diverged")
		}
	}
	for _, outcome := range []string{"green", "red", "interrupt", "timeout"} {
		if !o.terminal[outcome] {
			return fmt.Errorf("terminal outcome %q was not observed", outcome)
		}
	}
	return nil
}

func (o *systemOwner) cleanup() error {
	if o.root == "" {
		return nil
	}
	root := o.root
	o.root = ""
	if err := os.RemoveAll(root); err != nil {
		return err
	}
	if _, err := os.Stat(root); !errors.Is(err, os.ErrNotExist) {
		return errors.New("system repositories remain after cleanup")
	}
	return nil
}

func mergeEnvironment(base, overrides []string) []string {
	want := map[string]string{}
	for _, entry := range overrides {
		key, _, _ := strings.Cut(entry, "=")
		want[key] = entry
	}
	out := make([]string, 0, len(base)+len(overrides))
	for _, entry := range base {
		key, _, _ := strings.Cut(entry, "=")
		if _, replaced := want[key]; replaced {
			continue
		}
		out = append(out, entry)
	}
	for _, entry := range overrides {
		out = append(out, entry)
	}
	return out
}

func TestSelectedExecutableComposition(t *testing.T) {
	var first string
	for _, repo := range owner.repos {
		result := owner.runSelected(repo, "version")
		if result.code != 0 {
			t.Fatalf("version exit = %d: %s", result.code, result.stderr)
		}
		if !strings.Contains(result.stderr, "command-registry:version") {
			t.Fatalf("selected command bypassed the production registry: %q", result.stderr)
		}
		if first == "" {
			first = result.stdout
		} else if result.stdout != first {
			t.Fatalf("selected executable changed behavior: first=%q current=%q", first, result.stdout)
		}
	}
	owner.markTerminal("green")
}

func TestWrapperInstallFreshnessAndReloadJourneys(t *testing.T) {
	linked := owner.repos[0]
	if result := owner.runSelected(linked, "link", "copy"); result.code != 0 {
		t.Fatalf("link exit = %d: %s", result.code, result.stderr)
	}
	wrapper := filepath.Join(linked, ".bench", "bin", "bench.sh")
	wrapped := owner.runWrapper(linked, wrapper, "version")
	if direct := owner.runSelected(linked, "version"); wrapped.code != 0 || wrapped.stdout != direct.stdout {
		t.Fatalf("wrapper route = (%d, %q, %q), direct = (%d, %q, %q)", wrapped.code, wrapped.stdout, wrapped.stderr, direct.code, direct.stdout, direct.stderr)
	}
	if stale := owner.runSelected(linked, "freshness-check", linked); stale.code == 0 {
		t.Fatal("freshness-check accepted a repository that does not match the selected executable")
	}
	if reload := owner.runSelected(linked, "doctor"); !strings.Contains(reload.stdout, "ok: repo-local bench resolvable at .bench/bin/bench.sh") {
		t.Fatalf("fresh process did not reload installed state: exit=%d stdout=%q stderr=%q", reload.code, reload.stdout, reload.stderr)
	}
}

func TestCanaryInventoryAndSelectedExecutable(t *testing.T) {
	fixtures, err := canary.Fixtures(filepath.Join(owner.kit, "tests", "canary"))
	if err != nil {
		t.Fatal(err)
	}
	result := owner.runSelected(owner.repos[1], "canary", owner.kit)
	want := fmt.Sprintf("canary inventory ok (%d fixture bindings)\n", len(fixtures))
	if result.code != 0 || result.stdout != want {
		t.Fatalf("canary inventory = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	if !strings.Contains(result.stderr, "command-registry:canary") {
		t.Fatalf("canary bypassed selected command-registry inventory route: %q", result.stderr)
	}
}

func TestRedTimeoutAndDescendantTeardown(t *testing.T) {
	if red := owner.runSelected(owner.repos[1], "not-a-command"); red.code == 0 {
		t.Fatal("unknown command unexpectedly green")
	}
	owner.markTerminal("red")
	interrupted, err := owner.interruptProcessGroup()
	if err != nil {
		t.Fatal(err)
	}
	requireProcessGone(t, interrupted, "interrupt")
	owner.markTerminal("interrupt")
	descendant, err := owner.timeoutProcessGroup(100 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	requireProcessGone(t, descendant, "timeout")
	owner.markTerminal("timeout")
}

func requireProcessGone(t *testing.T, pid int, outcome string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("descendant %d remains after process-group %s", pid, outcome)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

const strippedJourneyMarker = "branch-native-stripped-distribution"

func TestStrippedDistributionJourney(t *testing.T) {
	repo := owner.repos[2]
	if result := owner.runSelected(repo, "link", "copy"); result.code != 0 {
		t.Fatalf("stripped install exit = %d: %s", result.code, result.stderr)
	}
	for _, path := range []string{"capture", "decisions", "specs"} {
		if err := os.RemoveAll(filepath.Join(repo, path)); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"ROADMAP.md", ".bench-notes.md"} {
		if err := os.Remove(filepath.Join(repo, path)); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	}
	for _, path := range []string{".bench/bin/bench.sh", ".agents/commands/bench-implement-spec.md", "AGENTS.md"} {
		if info, err := os.Stat(filepath.Join(repo, filepath.FromSlash(path))); err != nil || !info.Mode().IsRegular() {
			t.Fatalf("installed package path %s is unavailable: %v", path, err)
		}
	}
	for _, path := range []string{"capture", "decisions", "specs", "ROADMAP.md", ".bench-notes.md"} {
		if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(path))); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("excluded path %s remains in stripped subject", path)
		}
	}
	wrapper := filepath.Join(repo, ".bench", "bin", "bench.sh")
	result := owner.runWrapper(repo, wrapper, "version")
	if result.code != 0 || !strings.Contains(result.stderr, "command-registry:version") {
		t.Fatalf("stripped wrapper route = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	if _, err := os.Stat(filepath.Join(repo, "missing-kit")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("wrapper reached the deliberately unavailable source checkout")
	}
	if strippedJourneyMarker == "" {
		t.Fatal("stripped journey marker is empty")
	}
}
