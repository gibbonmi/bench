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
		if result := o.runAt(repo, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
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
	return o.runWithInput(dir, overrides, "", program, args...)
}

// runWithInput is the package's one owned process start: runAt is the no-stdin form, so
// a launch that feeds a hook its envelope on stdin joins the same starts ledger rather
// than opening a second launch path beside it.
func (o *systemOwner) runWithInput(dir string, overrides []string, input string, program string, args ...string) processResult {
	o.mu.Lock()
	o.starts++
	o.mu.Unlock()
	cmd := exec.Command(program, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnvironment(os.Environ(), overrides)
	cmd.Stdin = strings.NewReader(input)
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

func TestWorktreeReauthorizeJourney(t *testing.T) {
	repo, err := os.MkdirTemp(owner.root, "reauthorize [journey]-")
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "base.txt"}, {"-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "base"}} {
		if len(args) == 2 && args[0] == "add" {
			if err := os.WriteFile(filepath.Join(repo, "base.txt"), []byte("base\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		if result := owner.runAt(repo, nil, "git", args...); result.code != 0 {
			t.Fatalf("git %q = (%d, %q)", args, result.code, result.stderr)
		}
	}
	base := systemGitOutput(t, repo, "rev-parse", "HEAD")
	if err := os.WriteFile(filepath.Join(repo, "reviewed.txt"), []byte("reviewed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "reviewed.txt"}, {"-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "reviewed"}} {
		if result := owner.runAt(repo, nil, "git", args...); result.code != 0 {
			t.Fatalf("git %q = (%d, %q)", args, result.code, result.stderr)
		}
	}
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	home, err := os.MkdirTemp("", "bench-system-reauthorize-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(home) })
	shared := []string{"BENCH_HOME=" + home, "BENCH_COMMAND_OBSERVE=1"}
	created := owner.runAt(repo, shared, owner.selected.path, "worktree", "create", "--request", "lost-token", "--label", "owned")
	if created.code != 0 || !strings.Contains(created.stderr, "command-registry:worktree") {
		t.Fatalf("worktree create = (%d, %q, %q)", created.code, created.stdout, created.stderr)
	}
	lines := strings.Split(strings.TrimSpace(created.stdout), "\n")
	if len(lines) != 5 {
		t.Fatalf("worktree create output = %q", created.stdout)
	}
	fields := strings.Split(strings.TrimSpace(lines[1]), ",")
	if len(fields) != 3 {
		t.Fatalf("worktree create row = %q", lines[1])
	}
	path, err := systemTOONCell(fields[0])
	if err != nil {
		t.Fatalf("worktree create path = %q: %v", fields[0], err)
	}
	assignment, err := systemTOONCell(fields[1])
	if err != nil {
		t.Fatalf("worktree create assignment = %q: %v", fields[1], err)
	}
	tip := systemGitOutput(t, path, "rev-parse", "HEAD")
	before := systemReauthorizeEvidence(t, repo, path)
	result := owner.runAt(repo, shared, owner.selected.path, "worktree", "reauthorize", "--assignment", assignment, "--request", "replacement-token", "--base", base, "--source-tip", tip, path)
	if result.code != 0 || !strings.Contains(result.stderr, "command-registry:worktree") {
		t.Fatalf("worktree reauthorize = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	if want := "reauthorized{assignment=" + assignment + ",recorded_start=" + tip + ",approved_base=" + base + ",source_tip=" + tip + ",state=active}\n"; result.stdout != want {
		t.Fatalf("worktree reauthorize stdout = %q, want %q", result.stdout, want)
	}
	after := systemReauthorizeEvidence(t, repo, path)
	if after != before {
		t.Fatalf("worktree reauthorize changed retained contents or state: before=%#v after=%#v", before, after)
	}
	released := owner.runAt(repo, shared, owner.selected.path, "worktree", "release", "--request", "replacement-token", path)
	if released.code != 0 {
		t.Fatalf("replacement token did not authenticate release: (%d, %q, %q)", released.code, released.stdout, released.stderr)
	}
}

func TestWorktreeLandPublicRaceAndRerun(t *testing.T) {
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	loser := systemCreateLandingWorktree(t, root, home, "race-loser", "race loser")
	base := systemGitOutput(t, root, "rev-parse", "main")
	systemCommit(t, loser.path, "loser.txt", "loser\n", "reviewed loser")
	loser.tip = systemGitOutput(t, loser.path, "rev-parse", "HEAD")
	review := systemSelected(t, loser.path, systemLandEnv(root, home, tally, trees, ready, release), "preflight", "review", "x", "--base", base)
	if review.code != 0 {
		t.Fatalf("reviewed source %s = (%d, %q, %q)", loser.request, review.code, review.stdout, review.stderr)
	}

	loserCommand, loserOut, loserErr := systemStartLand(t, root, home, tally, trees, ready, release, loser, base)
	loserDone := make(chan error, 1)
	go func() { loserDone <- loserCommand.Wait() }()
	t.Cleanup(func() {
		if loserCommand.ProcessState == nil {
			_ = loserCommand.Process.Kill()
			<-loserDone
		}
	})
	type gateReady struct {
		count int
		byte  []byte
		err   error
	}
	readyResult := make(chan gateReady, 1)
	go func() {
		reader, openErr := os.Open(ready)
		if openErr != nil {
			readyResult <- gateReady{err: openErr}
			return
		}
		readyByte := make([]byte, 1)
		count, readErr := reader.Read(readyByte)
		closeErr := reader.Close()
		if readErr == nil {
			readErr = closeErr
		}
		readyResult <- gateReady{count: count, byte: readyByte, err: readErr}
	}()
	select {
	case readyState := <-readyResult:
		if readyState.err != nil || readyState.count != 1 || string(readyState.byte) != "r" {
			t.Fatalf("loser gate barrier = (%d, %q, %v)", readyState.count, readyState.byte, readyState.err)
		}
	case loserWaitErr := <-loserDone:
		t.Fatalf("loser exited before gate barrier = (%d, %q, %v)", systemExitCode(loserWaitErr), loserOut.String(), loserWaitErr)
	}

	systemCommit(t, root, "winner.txt", "winner\n", "winner moves destination")
	winnerCommit := systemGitOutput(t, root, "rev-parse", "main")

	priorRefs := systemGitOutput(t, root, "show-ref", "--head")
	priorMarker := systemGitOutput(t, root, "rev-parse", "refs/bench/green/main")
	releaseWriter, err := os.OpenFile(release, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := releaseWriter.Write([]byte("go\n")); err != nil {
		t.Fatal(err)
	}
	if err := releaseWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if code := systemExitCode(<-loserDone); code != 1 || !strings.HasPrefix(loserOut.String(), "refused{detail=") || !strings.Contains(loserOut.String(), "landing destination checkout changed") || strings.Contains(loserOut.String(), "infrastructure") || !strings.Contains(loserErr.String(), "command-registry:worktree") {
		t.Fatalf("loser land = (%d, %q, %q)", code, loserOut.String(), loserErr.String())
	}
	if got := systemGitOutput(t, root, "show-ref", "--head"); got != priorRefs {
		t.Fatalf("loser moved winner refs: got %q want %q", got, priorRefs)
	}
	if got := systemGitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != priorMarker || priorMarker != base {
		t.Fatalf("loser moved project-green: got %s want original base %s", got, base)
	}
	if got := systemGitOutput(t, loser.path, "rev-parse", "HEAD"); got != loser.tip {
		t.Fatalf("loser source tip moved: got %s want %s", got, loser.tip)
	}
	if got := systemGitOutput(t, loser.path, "symbolic-ref", "--quiet", "HEAD"); got != loser.branch {
		t.Fatalf("loser branch changed: got %s want %s", got, loser.branch)
	}
	if got := systemGitOutput(t, loser.path, "status", "--porcelain=v1", "--untracked-files=all"); got != "" {
		t.Fatalf("loser source status = %q", got)
	}
	active := systemSelected(t, root, systemLandEnv(root, home, tally, trees, ready, release), "worktree", "list")
	if active.code != 0 || !strings.Contains(active.stdout, loser.assignment) || !strings.Contains(active.stdout, ",active,assignment,present,") {
		t.Fatalf("loser active assignment proof = (%d, %q, %q)", active.code, active.stdout, active.stderr)
	}
	path := systemSelected(t, root, systemLandEnv(root, home, tally, trees, ready, release), "worktree", "path", "race-loser")
	if path.code != 0 || strings.TrimSpace(path.stdout) != loser.path {
		t.Fatalf("loser owner path proof = (%d, %q, %q)", path.code, path.stdout, path.stderr)
	}
	if got, readErr := os.ReadFile(tally); readErr != nil || string(got) != "l" {
		t.Fatalf("loser gate tally = %q, %v", got, readErr)
	}

	rerun := systemLand(t, root, home, tally, trees, ready, release, loser, base)
	if rerun.code != 0 || !strings.Contains(rerun.stderr, "command-registry:worktree") {
		t.Fatalf("rerun land = (%d, %q, %q)", rerun.code, rerun.stdout, rerun.stderr)
	}
	published := systemGitOutput(t, root, "rev-parse", "main")
	publishedTree := systemGitOutput(t, root, "rev-parse", published+"^{tree}")
	rerunEnvelope := "landed{source_base=" + base + ",source_tip=" + loser.tip + ",destination_base=" + winnerCommit + ",published_commit=" + published + ",tree=" + publishedTree + ",worktree=released}\n"
	if rerun.stdout != rerunEnvelope {
		t.Fatalf("rerun terminal envelope = %q, want %q", rerun.stdout, rerunEnvelope)
	}
	if parents := strings.Fields(systemGitOutput(t, root, "rev-list", "--parents", "-n", "1", published)); len(parents) != 3 || parents[1] != winnerCommit || parents[2] != loser.tip {
		t.Fatalf("recomposed parents = %q, want winner %s and source %s", parents, winnerCommit, loser.tip)
	}
	if got := systemGitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("rerun project-green = %s, want %s", got, published)
	}
	if got := systemGitOutput(t, root, "status", "--porcelain=v1", "--untracked-files=all"); got != "" {
		t.Fatalf("rerun destination status = %q", got)
	}
	if _, statErr := os.Stat(loser.path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("rerun retained assignment worktree: %v", statErr)
	}
	if got, readErr := os.ReadFile(tally); readErr != nil || string(got) != "ll" {
		t.Fatalf("rerun gate tally = %q, %v", got, readErr)
	}
	gateTrees, readErr := os.ReadFile(trees)
	if readErr != nil {
		t.Fatal(readErr)
	}
	observedTrees := strings.Fields(string(gateTrees))
	if len(observedTrees) != 2 || observedTrees[1] != publishedTree || observedTrees[0] == observedTrees[1] {
		t.Fatalf("prospective gate trees = %q, want distinct loser rerun tree %s", observedTrees, publishedTree)
	}
}

type systemLandingWorktree struct {
	path, assignment, branch, request, tip string
}

func systemLandingRaceFixture(t *testing.T) (root, home, tally, trees, ready, release string) {
	t.Helper()
	var err error
	root, err = os.MkdirTemp(owner.root, "landing-race [journey]-")
	if err != nil {
		t.Fatal(err)
	}
	if result := owner.runAt(root, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init landing race = (%d, %q)", result.code, result.stderr)
	}
	home, err = os.MkdirTemp(owner.root, "landing-race [home]-")
	if err != nil {
		t.Fatal(err)
	}
	tally = filepath.Join(home, "gate-tally")
	trees = filepath.Join(home, "gate-trees")
	ready = filepath.Join(home, "loser-ready")
	release = filepath.Join(home, "loser-release")
	if err := syscall.Mkfifo(ready, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(release, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	gate := "#!/bin/sh\nset -eu\nruntime=$1\nrg -q '^Status: implemented$' \"$runtime/specs/x/spec.md\"\ntree=$(git -C \"$runtime\" write-tree)\nprintf '%s\\n' \"$tree\" >> \"$LAND_GATE_TREES\"\nif [ -f \"$runtime/loser.txt\" ]; then\n  printf l >> \"$LAND_GATE_TALLY\"\n  if [ ! -f \"$runtime/winner.txt\" ]; then\n    printf r > \"$LAND_RACE_READY\"\n    IFS= read -r _ < \"$LAND_RACE_RELEASE\"\n  fi\nelse\n  printf w >> \"$LAND_GATE_TALLY\"\nfi\n"
	for _, file := range []string{"gate.sh", "gate-prospective.sh"} {
		if err := os.WriteFile(filepath.Join(root, ".bench", file), []byte(gate), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	inputs := "{\"schema\":1,\"closure\":\"local\",\"environment\":[\"LAND_GATE_TALLY\",\"LAND_GATE_TREES\",\"LAND_RACE_READY\",\"LAND_RACE_RELEASE\"],\"paths\":[],\"tools\":[]}\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(inputs), 0o644); err != nil {
		t.Fatal(err)
	}
	specBody := "# x\n\nStatus: staged\n\n## User stories\n1. Land source.\n\n### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n| LX1 | 1 | lands | command | catches failure |\n\n## Ownership fences\n\n- `loser.txt`\n- `winner.txt`\n"
	if err := os.MkdirAll(filepath.Join(root, "specs", "x", "tickets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "x", "spec.md"), []byte(specBody), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "x", "tickets", "one.md"), []byte("Ticket covers LX1.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "add", ".")
	systemGit(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "landing race base")
	base := systemGitOutput(t, root, "rev-parse", "HEAD")
	systemGit(t, root, "update-ref", "refs/bench/green/main", base)
	return root, home, tally, trees, ready, release
}

func systemCreateLandingWorktree(t *testing.T, root, home, label, request string) systemLandingWorktree {
	t.Helper()
	result := systemSelected(t, root, []string{"BENCH_HOME=" + home, "BENCH_SYSTEM_ROOT=" + root, "BENCH_COMMAND_OBSERVE=1"}, "worktree", "create", "--request", request, "--label", label)
	if result.code != 0 || !strings.Contains(result.stderr, "command-registry:worktree") {
		t.Fatalf("worktree create %s = (%d, %q, %q)", label, result.code, result.stdout, result.stderr)
	}
	lines := strings.Split(strings.TrimSpace(result.stdout), "\n")
	if len(lines) != 5 {
		t.Fatalf("worktree create %s output = %q", label, result.stdout)
	}
	fields := strings.Split(lines[1], ",")
	if len(fields) != 3 || fields[2] != "active" {
		t.Fatalf("worktree create %s row = %q", label, lines[1])
	}
	path, err := systemTOONCell(fields[0])
	if err != nil {
		t.Fatal(err)
	}
	assignment, err := systemTOONCell(fields[1])
	if err != nil {
		t.Fatal(err)
	}
	return systemLandingWorktree{path: path, assignment: assignment, branch: systemGitOutput(t, path, "symbolic-ref", "--quiet", "HEAD"), request: request}
}

func systemCommit(t *testing.T, root, name, body, message string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "add", name)
	systemGit(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", message)
}

func systemGit(t *testing.T, root string, args ...string) {
	t.Helper()
	result := owner.runAt(root, nil, "git", args...)
	if result.code != 0 {
		t.Fatalf("git %q in %q = (%d, %q)", args, root, result.code, result.stderr)
	}
}

func systemLandEnv(root, home, tally, trees, ready, release string) []string {
	return []string{
		"BENCH_RUN_BINARY=" + owner.selected.path,
		"BENCH_KIT=" + owner.kit,
		"BENCH_COMMAND_OBSERVE=1",
		"BENCH_SYSTEM_ROOT=" + root,
		"BENCH_HOME=" + home,
		"LAND_GATE_TALLY=" + tally,
		"LAND_GATE_TREES=" + trees,
		"LAND_RACE_READY=" + ready,
		"LAND_RACE_RELEASE=" + release,
	}
}

func systemSelected(t *testing.T, root string, env []string, args ...string) processResult {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runAt(root, env, owner.selected.path, args...)
}

func systemLand(t *testing.T, root, home, tally, trees, ready, release string, source systemLandingWorktree, base string) processResult {
	t.Helper()
	return systemSelected(t, root, systemLandEnv(root, home, tally, trees, ready, release), systemLandArgs(source, base)...)
}

func systemStartLand(t *testing.T, root, home, tally, trees, ready, release string, source systemLandingWorktree, base string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(owner.selected.path, systemLandArgs(source, base)...)
	cmd.Dir = root
	cmd.Env = mergeEnvironment(os.Environ(), systemLandEnv(root, home, tally, trees, ready, release))
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	owner.mu.Lock()
	owner.starts++
	owner.mu.Unlock()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	return cmd, stdout, stderr
}

// systemLandArgs is the one landing invocation both race harnesses drive, so a
// blocking and a foreground land can never diverge in what they actually ran.
func systemLandArgs(source systemLandingWorktree, base string) []string {
	return []string{"worktree", "land", "--request", source.request, "--base", base, "--source-tip", source.tip, "--spec", "x", "-m", "land race source", source.path}
}

func systemExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode()
	}
	return -1
}

func systemGitOutput(t *testing.T, repo string, args ...string) string {
	t.Helper()
	result := owner.runAt(repo, nil, "git", args...)
	if result.code != 0 {
		t.Fatalf("git %q in %q = (%d, %q)", args, repo, result.code, result.stderr)
	}
	return strings.TrimSpace(result.stdout)
}

func systemTOONCell(value string) (string, error) {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, `"`) {
		return strconv.Unquote(value)
	}
	return value, nil
}

type systemReauthorizeState struct {
	Tree, Index, Status, Refs string
}

func systemReauthorizeEvidence(t *testing.T, root, path string) systemReauthorizeState {
	t.Helper()
	tree, err := os.ReadFile(filepath.Join(path, "reviewed.txt"))
	if err != nil {
		t.Fatal(err)
	}
	indexPath := systemGitOutput(t, path, "rev-parse", "--git-path", "index")
	index, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	return systemReauthorizeState{
		Tree:   string(tree),
		Index:  string(index),
		Status: systemGitOutput(t, path, "status", "--porcelain=v1"),
		Refs:   systemGitOutput(t, root, "show-ref", "--head"),
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
	for _, path := range []string{"capture", "decisions", "specs", "roadmap"} {
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
	for _, path := range []string{"capture", "decisions", "specs", "roadmap", "ROADMAP.md", ".bench-notes.md"} {
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
