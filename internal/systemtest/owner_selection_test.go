//go:build system

package systemtest

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/canary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
