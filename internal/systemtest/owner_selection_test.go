//go:build system

package systemtest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/canary"
)

func TestFocusedRunUsesSelectedVerdict(t *testing.T) {
	for _, args := range [][]string{
		{"-test.run=^TestOwnerSelectionChild$"},
		{"-test.run", "^TestOwnerSelectionChild$"},
		{"--test.run=^TestOwnerSelectionChild$"},
		{"-test.run=", "-test.run=^TestOwnerSelectionChild$"},
	} {
		for _, verdict := range []string{"pass", "fail"} {
			t.Run(strings.Join(args, " ")+"/"+verdict, func(t *testing.T) {
				result, _ := runOwnerSelectionChild(t, verdict, args...)
				wantCode := 0
				if verdict == "fail" {
					wantCode = 1
					if !strings.Contains(result.stdout, "selected failure diagnostic") {
						t.Fatalf("selected failure lost: %q", result.stdout)
					}
				}
				if result.code != wantCode || !strings.Contains(result.stdout, "selected test ran") || strings.Contains(result.stderr, "system owner verification:") {
					t.Fatalf("focused verdict = (%d, %q, %q)", result.code, result.stdout, result.stderr)
				}
			})
		}
	}
}

func TestFocusedRunStillCleansOwnerRoot(t *testing.T) {
	for _, verdict := range []string{"pass", "fail", "cleanup-failure"} {
		t.Run(verdict, func(t *testing.T) {
			result, root := runOwnerSelectionChild(t, verdict, "-test.run=^TestOwnerSelectionChild$")
			if verdict == "cleanup-failure" {
				if result.code == 0 || !strings.Contains(result.stderr, "system owner cleanup:") {
					t.Fatalf("cleanup failure = (%d, %q)", result.code, result.stderr)
				}
				return
			}
			assertOwnerRootRemoved(t, root)
		})
	}
}

func TestUnfilteredRunVerifiesLedger(t *testing.T) {
	// Go lists the live test inventory. Exclude each other test with test.skip
	// so the child reaches TestMain without a non-empty test.run value.
	listed := owner.runAt(owner.root, []string{"BENCH_KIT=" + owner.kit}, os.Args[0], "-test.list=.")
	if listed.code != 1 || !strings.Contains(listed.stderr, "system owner verification: no selected executable observations recorded") {
		t.Fatalf("unfiltered list baseline = (%d, %q, %q)", listed.code, listed.stdout, listed.stderr)
	}
	var excluded []string
	for _, name := range strings.Fields(listed.stdout) {
		if strings.HasPrefix(name, "Test") && name != "TestOwnerSelectionChild" {
			excluded = append(excluded, regexp.QuoteMeta(name))
		}
	}
	if len(excluded) == 0 {
		t.Fatal("the child listed no tests to exclude")
	}
	args := []string{"-test.skip=^(" + strings.Join(excluded, "|") + ")$"}
	for _, tc := range []struct {
		role, diagnostic string
	}{
		{"repositories", "repository count = 0, want 3"},
		{"processes", "no owned process starts recorded"},
		{"executable", "no selected executable observations recorded"},
		{"identity", "selected executable identity ledger diverged"},
		{"green", `terminal outcome "green" was not observed`},
		{"red", `terminal outcome "red" was not observed`},
		{"interrupt", `terminal outcome "interrupt" was not observed`},
		{"timeout", `terminal outcome "timeout" was not observed`},
		{"complete", ""},
		{"cleanup-failure", "system owner cleanup:"},
	} {
		t.Run(tc.role, func(t *testing.T) {
			result, root := runOwnerSelectionChild(t, tc.role, args...)
			wantCode := 1
			if tc.diagnostic == "" {
				wantCode = 0
			}
			if result.code != wantCode || !strings.Contains(result.stdout, "selected test ran") || !strings.Contains(result.stderr, tc.diagnostic) {
				t.Fatalf("unfiltered verdict = (%d, %q, %q), want (%d, %q)", result.code, result.stdout, result.stderr, wantCode, tc.diagnostic)
			}
			if tc.role != "cleanup-failure" {
				assertOwnerRootRemoved(t, root)
			}
		})
	}
	// An explicit empty value still requires the full ledger, even when a
	// preceding occurrence names a selection.
	result, root := runOwnerSelectionChild(t, "executable", append(args, "-test.run=^$", "-test.run=")...)
	if result.code != 1 || !strings.Contains(result.stderr, "system owner verification: no selected executable observations recorded") {
		t.Fatalf("empty selection = (%d, %q)", result.code, result.stderr)
	}
	assertOwnerRootRemoved(t, root)
}

func runOwnerSelectionChild(t *testing.T, role string, args ...string) (processResult, string) {
	t.Helper()
	report := filepath.Join(t.TempDir(), "owner-root")
	result := owner.runAt(owner.root, []string{
		"BENCH_KIT=" + owner.kit,
		"BENCH_SYSTEM_SELECTION_CHILD=" + role,
		"BENCH_SYSTEM_SELECTION_REPORT=" + report,
	}, os.Args[0], append([]string{"-test.v"}, args...)...)
	data, err := os.ReadFile(report)
	if err != nil {
		t.Fatalf("child root report: %v; child = (%d, %q, %q)", err, result.code, result.stdout, result.stderr)
	}
	root := string(data)
	t.Cleanup(func() {
		if err := os.RemoveAll(root); err != nil {
			t.Error(err)
		}
	})
	return result, root
}

func assertOwnerRootRemoved(t *testing.T, root string) {
	t.Helper()
	if _, err := os.Stat(root); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("owner root remains after TestMain: %q: %v", root, err)
	}
}

func TestOwnerSelectionChild(t *testing.T) {
	role := os.Getenv("BENCH_SYSTEM_SELECTION_CHILD")
	if role == "" {
		return
	}
	if err := os.WriteFile(os.Getenv("BENCH_SYSTEM_SELECTION_REPORT"), []byte(owner.root), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Log("selected test ran")
	switch role {
	case "pass":
		return
	case "fail":
		t.Fatal("selected failure diagnostic")
	case "repositories":
		owner.repos = nil
	case "processes":
		owner.starts = 0
	case "executable":
		return
	case "cleanup-failure":
		// A NUL makes removal fail on every host, including a root process.
		// The parent retains the valid root path and removes it afterward.
		owner.root += "\x00"
	}
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	if role == "identity" {
		owner.seen[0].inode++
	}
	for _, outcome := range []string{"green", "red", "interrupt", "timeout"} {
		if role != outcome {
			owner.markTerminal(outcome)
		}
	}
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

// TestChildEnvironmentDefaultsBenchHomeWhenUnnamed checks that overrides
// without BENCH_HOME use the owner's private home, never the operator's home.
func TestChildEnvironmentDefaultsBenchHomeWhenUnnamed(t *testing.T) {
	env := owner.childEnvironment([]string{"BENCH_COMMAND_OBSERVE=1"})
	if got := envValue(env, benchhome.Env); got != owner.home {
		t.Fatalf("BENCH_HOME = %q, want the owner's private home %q", got, owner.home)
	}
}

// TestChildEnvironmentKeepsAnExplicitBenchHome checks the exception: a caller
// that names its own BENCH_HOME keeps it instead of the owner's private home.
func TestChildEnvironmentKeepsAnExplicitBenchHome(t *testing.T) {
	want := t.TempDir()
	env := owner.childEnvironment([]string{benchhome.Env + "=" + want})
	if got := envValue(env, benchhome.Env); got != want {
		t.Fatalf("BENCH_HOME = %q, want the caller's own home %q", got, want)
	}
}

// envValue returns the value entries assigns key, or "" when entries never names it.
func envValue(entries []string, key string) string {
	for _, entry := range entries {
		if entryKey, value, found := strings.Cut(entry, "="); found && entryKey == key {
			return value
		}
	}
	return ""
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
