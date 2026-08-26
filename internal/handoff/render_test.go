package handoff

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestRenderPathAbbreviates pins the `~` form, including exactly at $HOME where the
// remainder is empty and the abbreviation is the whole path.
func TestRenderPathAbbreviates(t *testing.T) {
	cases := []struct{ root, home, want string }{
		{"/home/a", "/home/a", "~"},
		{"/home/a/", "/home/a", "~"},
		{"/home/a", "/home/a/", "~"},
		{"/home/a/workspace/bench", "/home/a", "~/workspace/bench"},
		{"/home/a/x", "/home/a", "~/x"},
		{"/", "/", "~"},
		{"/srv/x", "/", "~/srv/x"},
	}
	for _, tc := range cases {
		if got := renderPath(tc.root, tc.home); got != tc.want {
			t.Fatalf("renderPath(%q, %q) = %q, want %q", tc.root, tc.home, got, tc.want)
		}
	}
}

// TestRenderPathOutsideHome pins the other side of the boundary. A prefix match on the
// raw string would turn /home/abc into ~bc: a path that resolves nowhere.
func TestRenderPathOutsideHome(t *testing.T) {
	cases := []struct{ root, home, want string }{
		{"/home/abc", "/home/a", "/home/abc"},
		{"/home/abc/deep", "/home/a", "/home/abc/deep"},
		{"/srv/checkouts/bench", "/home/a", "/srv/checkouts/bench"},
		{"/home/a", "", "/home/a"},
		{"/home/a", "relative/home", "/home/a"},
		{"/home", "/home/a", "/home"},
	}
	for _, tc := range cases {
		if got := renderPath(tc.root, tc.home); got != tc.want {
			t.Fatalf("renderPath(%q, %q) = %q, want %q", tc.root, tc.home, got, tc.want)
		}
	}
}

func TestCommandUsesCleanBoardRouteFallback(t *testing.T) {
	got := runHandoffCommand(t, nil, map[string]string{"ROADMAP.md": "# Roadmap\n"})
	if !strings.Contains(string(got), "## Next command\n\n`bench roadmap`") {
		t.Fatalf("handoff next command = %q, want bench roadmap", got)
	}
}

func TestCommandUsesCodexRouteForPhase(t *testing.T) {
	got := runHandoffCommand(t, []string{"--harness", "codex"}, map[string]string{
		"capture/IDEAS.md": "- 2026-08-18  pending\n",
	})
	if !strings.Contains(string(got), "## Next command\n\n`$bench-drain`") {
		t.Fatalf("handoff next command = %q, want $bench-drain", got)
	}
}

// HC22. `--harness none` is accepted, and the routed shell command is the one a cold
// session without a harness reads. The board is led by `git push`: a clean tree one commit
// ahead of its upstream.
func TestCommandUsesShellRouteForNoHarness(t *testing.T) {
	got := runHandoffCommand(t, []string{"--harness", "none"}, map[string]string{"ROADMAP.md": "# Roadmap\n"}, unpushedCommit)
	if !strings.Contains(string(got), "## Next command\n\n`git push`") {
		t.Fatalf("handoff next command = %q, want git push", got)
	}
}

// unpushedCommit gives the repository an upstream it is one commit ahead of, which is the
// `git push` board.
func unpushedCommit(t *testing.T, root string) {
	t.Helper()
	remote := filepath.Join(t.TempDir(), "remote.git")
	for _, args := range [][]string{
		{"init", "-q", "--bare", remote},
		{"-C", root, "remote", "add", "origin", remote},
		{"-C", root, "push", "-q", "-u", "origin", "HEAD"},
		{"-C", root, "commit", "-q", "--allow-empty", "-m", "ahead"},
	} {
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
}

func runHandoffCommand(t *testing.T, commandArgs []string, files map[string]string, setup ...func(*testing.T, string)) []byte {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{{"init", "-q", root}, {"-C", root, "config", "user.email", "t@example.com"}, {"-C", root, "config", "user.name", "t"}} {
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	files[".bench/keep"] = "\n"
	for path, content := range files {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if out, err := exec.Command("git", "-C", root, "add", "-A").CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, out)
	}
	if out, err := exec.Command("git", "-C", root, "commit", "-qm", "base").CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v: %s", err, out)
	}
	for _, step := range setup {
		step(t, root)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	if out, code := Command(commandArgs); code != 0 {
		t.Fatalf("handoff = (%q, %d), want exit 0", out, code)
	}
	got, err := os.ReadFile(filepath.Join(root, "capture", "session-handoff.md"))
	if err != nil {
		t.Fatal(err)
	}
	return got
}
