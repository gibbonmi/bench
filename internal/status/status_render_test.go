// Tests for the rendered board and the worktree rows it carries.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// A clean, committed tree with no signals renders the clean message and nothing else.
func TestRenderClean(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	out := render(root, false)
	if out != "bench: clean — nothing pending\n" {
		t.Errorf("clean render = %q", out)
	}
}

// A dirty tree leads with the git action; the capture-drain row is present but outranked.
func TestRenderDirtyLeadsGitOverDrainRow(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	// Make the tree dirty (uncommitted change) so the git signal fires.
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Park an idea so the drain signal fires.
	if err := os.WriteFile(filepath.Join(root, "capture/IDEAS.md"), []byte("- 2026-07-03  an idea\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := render(root, false)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if !strings.HasPrefix(lines[0], "▶ /bench-final-check  (git)") {
		t.Errorf("lead line = %q, want git action lead", lines[0])
	}
	if !strings.Contains(out, "1 idea(s), 0 open learning(s)") || !strings.Contains(out, "/bench-drain") {
		t.Errorf("drain row missing from:\n%s", out)
	}
}

// A working roadmap alone is not pending capture: no drain row, the board stays clean.
func TestRenderWorkingRoadmapAloneIsClean(t *testing.T) {
	root := initRepo(t)
	content := "# Roadmap\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n"
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	if out := render(root, false); out != "bench: clean — nothing pending\n" {
		t.Errorf("roadmap-only render = %q, want clean board", out)
	}
}

func TestAppendWorktreeIgnoresUnownedBranchPrefix(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	gitRun(t, root, "branch", "worktree-agent-orphan")

	if rows := appendWorktree(nil, root); len(rows) != 0 {
		t.Fatalf("branch prefix created worktree status without ownership evidence: %#v", rows)
	}
}

// appendWorktree must surface a discovery failure as a visible row, never as silence that
// a reader mistakes for "no worktree signals". This filesystem refusal reaches common-dir
// resolution before porcelain, while the PATH-stub fixtures cover typed and generic routing.
func TestAppendWorktreeSurfacesClassifyFailure(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	gitDir := filepath.Join(root, ".git")
	if err := os.Chmod(gitDir, 0o000); err != nil {
		t.Fatalf("chmod .git unreadable: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(gitDir, 0o755) })

	rows := appendWorktree(nil, root)
	if len(rows) == 0 {
		t.Fatal("appendWorktree dropped the classify failure instead of surfacing a row")
	}
	if !strings.Contains(rows[0].detail, "git common directory") || rows[0].action.render() != "bench worktree list" {
		t.Errorf("row = %#v, want typed resolution refusal", rows[0])
	}
}

func TestAppendWorktreeKeepsTypedAndPorcelainFailureActionsDistinct(t *testing.T) {
	for _, tc := range []struct {
		mode, detail, action string
	}{
		{"fail-rev-parse", "rev-parse", "bench worktree list"},
		{"fail-worktree", "git worktree list failed", "git worktree list"},
	} {
		t.Run(tc.mode, func(t *testing.T) {
			root := initRepo(t)
			gittest.StubGit(t, root, tc.mode, filepath.Join(t.TempDir(), "argv"))
			rows := appendWorktree(nil, root)
			if len(rows) != 1 || !strings.Contains(rows[0].detail, tc.detail) || rows[0].action.render() != tc.action {
				t.Fatalf("%s row = %#v", tc.mode, rows)
			}
		})
	}
}

func TestAppendWorktreeRendersBoundExpiryAsTypedFailure(t *testing.T) {
	restore := git.SetWorktreeListTimeoutForTest(100 * time.Millisecond)
	t.Cleanup(restore)
	root := initRepo(t)
	gittest.StubGit(t, root, "block-worktree", filepath.Join(t.TempDir(), "argv"))
	rows := appendWorktree(nil, root)
	if len(rows) != 1 || !strings.Contains(rows[0].detail, "worktree list") || rows[0].action.render() != "bench worktree list" {
		t.Fatalf("bound row = %#v", rows)
	}
}

func TestAppendWorktreeRendersTypedAdminRefusal(t *testing.T) {
	root := initRepo(t)
	gittest.FIFOWorktreeAdmin(t, root, "typed")
	rows := appendWorktree(nil, root)
	if len(rows) != 1 || !strings.Contains(rows[0].detail, "worktrees/typed/gitdir") || !strings.Contains(rows[0].detail, "fifo") || rows[0].action.render() != "bench worktree list" {
		t.Fatalf("typed row = %#v", rows)
	}
}

// TestExpandCensusSignalsEscapesTheLabel proves a control byte in a label is replaced
// before the fixed-width row renders, so no label splits a row. (Coverage row EC25.)
func TestExpandCensusSignalsEscapesTheLabel(t *testing.T) {
	t.Parallel()
	root, home := initRepo(t), t.TempDir()
	id := strings.Repeat("b", 32)
	seedAssignment(t, root, id, "alpha\x1b[31m", intent.StateActive)
	seedCensus(t, home, root, id, 2)

	signals := expandCensusSignals(root, home, censusSignals(t, root, home))
	if len(signals) != 1 || strings.ContainsRune(signals[0].Detail, 0x1b) {
		t.Fatalf("expanded census signals = %#v, want one row with no raw control byte", signals)
	}
	if want := "alpha" + sanitize.Controls("\x1b") + "[31m 2 raw calls"; signals[0].Detail != want {
		t.Fatalf("expanded detail = %q, want %q", signals[0].Detail, want)
	}
}

// censusSignals renders the summed census row as the Signal list the --all expanders
// consume.
func censusSignals(t *testing.T, root, home string) []Signal {
	t.Helper()
	rows := appendCensus(nil, root, home)
	out := make([]Signal, len(rows))
	for i, r := range rows {
		out[i] = newSignal(r.sev, r.signal, r.detail, r.action)
	}
	return out
}

// TestRenderAllExpandsTheCensusRowPerWorktree proves --all names each worktree, while
// the default board keeps the one summed row inside its five-row budget.
// (Coverage rows EC34 and EC17.)
func TestRenderAllExpandsTheCensusRowPerWorktree(t *testing.T) {
	root, home := initRepo(t), t.TempDir()
	t.Setenv("BENCH_HOME", home)
	first, second := strings.Repeat("b", 32), strings.Repeat("c", 32)
	seedAssignment(t, root, first, "alpha", intent.StateActive)
	seedAssignment(t, root, second, "beta", intent.StateActive)
	seedCensus(t, home, root, first, 2)
	seedCensus(t, home, root, second, 1)

	board := render(root, false)
	if !strings.Contains(board, "3 raw calls across 2 worktrees") {
		t.Fatalf("default board missing the summed census row:\n%s", board)
	}
	all := render(root, true)
	if !strings.Contains(all, "alpha 2 raw calls") || !strings.Contains(all, "beta 1 raw call") {
		t.Fatalf("--all board missing a per-worktree census row:\n%s", all)
	}
	if strings.Contains(all, "3 raw calls across 2 worktrees") {
		t.Fatalf("--all board kept the summed census row:\n%s", all)
	}
}

// TestRenderStatesOneRawCall proves the singular row reaches the rendered board.
func TestRenderStatesOneRawCall(t *testing.T) {
	root, home := initRepo(t), t.TempDir()
	t.Setenv("BENCH_HOME", home)
	id := strings.Repeat("b", 32)
	seedAssignment(t, root, id, "alpha", intent.StateActive)
	seedCensus(t, home, root, id, 1)

	if board := render(root, false); !strings.Contains(board, "1 raw call across 1 worktree") {
		t.Fatalf("board missing the singular census row:\n%s", board)
	}
}
