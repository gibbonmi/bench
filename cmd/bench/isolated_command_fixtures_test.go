package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

func commandFixtureKitRoot(t *testing.T) string {
	t.Helper()
	if kit := os.Getenv("BENCH_KIT"); kit != "" {
		root, err := filepath.Abs(kit)
		if err != nil {
			t.Fatalf("resolve BENCH_KIT: %v", err)
		}
		return root
	}
	root, err := git.Root()
	if err != nil {
		t.Fatalf("resolve Bench kit root: %v", err)
	}
	return root
}

func TestRunStatusRouteEmitsOneNextRow(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	t.Chdir(root)

	stdout := tempFile(t)
	if code := (Command{Stdout: stdout}).Run([]string{"status", "--route"}); code != 0 {
		t.Fatalf("status --route exit = %d, want 0", code)
	}
	if got := readFile(t, stdout); !strings.HasPrefix(got, "next[1]{state,why,command}:\n") {
		t.Fatalf("status --route = %q, want one next row", got)
	}
}

func TestGuardsQueryReportsFollowOnGuardThroughCommandSeam(t *testing.T) {
	kitRoot := commandFixtureKitRoot(t)
	root := gittest.RepoOnBranch(t, "main")
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	copyExecutable(t,
		filepath.Join(kitRoot, ".bench", "hooks", "block-bench-follow-on.sh"),
		filepath.Join(root, ".bench", "hooks", "block-bench-follow-on.sh"),
	)
	for _, path := range []string{".claude/settings.json", ".codex/hooks.json"} {
		writeAXIFixture(t, filepath.Join(root, filepath.FromSlash(path)), readPath(t, filepath.Join(kitRoot, filepath.FromSlash(path))))
	}
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"guards"})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("bench guards = stdout=%q stderr=%q exit=%d, want stdout/0", stdout.String(), stderr.String(), code)
	}
	for _, want := range []string{
		"block-bench-follow-on",
		"PreToolUse:Bash",
		"Bench shell follow-ons",
		"claude,codex",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("bench guards = %q, want %q", stdout.String(), want)
		}
	}
}
