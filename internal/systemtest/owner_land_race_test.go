//go:build system

package systemtest

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

func TestWorktreeLandPublicRaceAndRerun(t *testing.T) {
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	loser := systemCreateLandingWorktree(t, root, home, "race-loser", "race loser")
	base := systemGitOutput(t, root, "rev-parse", "main")
	systemCommit(t, loser.path, "loser.txt", "loser\n", "reviewed loser")
	loser.tip = systemRetainLandingEvidence(t, loser.path, base)
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

	unreviewed := systemLand(t, root, home, tally, trees, ready, release, loser, base)
	if unreviewed.code != 1 || !strings.Contains(unreviewed.stdout, "completion composition changes winner.txt") || systemGitOutput(t, root, "rev-parse", "main") != winnerCommit {
		t.Fatalf("unreviewed destination retry = (%d, %q, %q)", unreviewed.code, unreviewed.stdout, unreviewed.stderr)
	}
	systemGit(t, loser.path, "merge", "--no-edit", winnerCommit)
	loser.tip = systemRetainLandingEvidence(t, loser.path, base)
	if err := os.WriteFile(filepath.Join(loser.path, "retained-output"), []byte("retained\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rerun := systemLand(t, root, home, tally, trees, ready, release, loser, base)
	if rerun.code != 3 || !strings.Contains(rerun.stderr, "command-registry:worktree") {
		t.Fatalf("rerun land = (%d, %q, %q)", rerun.code, rerun.stdout, rerun.stderr)
	}
	published := systemGitOutput(t, root, "rev-parse", "main")
	publishedTree := systemGitOutput(t, root, "rev-parse", published+"^{tree}")
	next := "bench worktree land --resume '" + published + "' --request <request> --base '" + base + "' --source-tip '" + loser.tip + "' --spec 'x' '" + loser.path + "'"
	rerunEnvelope := "landed{source_base=" + base + ",source_tip=" + loser.tip + ",destination_base=" + winnerCommit + ",published_commit=" + published + ",tree=" + publishedTree + ",worktree=incomplete:release,next=" + next + ",census=0}\n"
	if rerun.stdout != rerunEnvelope {
		t.Fatalf("rerun terminal envelope = %q, want %q", rerun.stdout, rerunEnvelope)
	}
	if strings.Contains(rerun.stdout, loser.request) {
		t.Fatalf("rerun terminal envelope echoed caller token: %q", rerun.stdout)
	}
	if err := os.Remove(filepath.Join(loser.path, "retained-output")); err != nil {
		t.Fatal(err)
	}
	resumed := systemSelected(t, root, systemLandEnv(root, home, tally, trees, ready, release), "worktree", "land", "--resume", published, "--request", loser.request, "--base", base, "--source-tip", loser.tip, "--spec", "x", loser.path)
	releasedEnvelope := "effects[2]{effect,result}:\n  refresh,not-applicable\n  cleanup,complete\nlanded{source_base=" + base + ",source_tip=" + loser.tip + ",destination_base=" + winnerCommit + ",published_commit=" + published + ",tree=" + publishedTree + ",worktree=released,census=0}\n"
	if resumed.code != 0 || resumed.stdout != releasedEnvelope || !strings.Contains(resumed.stderr, "command-registry:worktree") {
		t.Fatalf("resume land = (%d, %q, %q), want %q", resumed.code, resumed.stdout, resumed.stderr, releasedEnvelope)
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

func systemRetainLandingEvidence(t *testing.T, source, base string) string {
	t.Helper()
	recordtest.RetainSingleChunk(t, source, "specs/x/spec.md", base)
	return systemGitOutput(t, source, "rev-parse", "HEAD")
}
