package conformance

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestHarnessUsesBenchConformanceRootAsGradedRoot(t *testing.T) {
	gradedRoot := t.TempDir()
	marker := filepath.Join(gradedRoot, "graded.txt")
	if err := os.WriteFile(marker, []byte("fixture root\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_CONFORMANCE_ROOT", gradedRoot)

	h := NewHarness(t)

	if h.Root != gradedRoot {
		t.Fatalf("Root = %q, want BENCH_CONFORMANCE_ROOT %q", h.Root, gradedRoot)
	}
	if h.KitRoot == "" {
		t.Fatal("KitRoot is empty")
	}
	if h.KitRoot == h.Root {
		t.Fatalf("KitRoot = Root = %q; kit root must stay separate from graded root", h.Root)
	}
	if got := h.ReadRootFile("graded.txt"); got != "fixture root\n" {
		t.Fatalf("ReadRootFile = %q, want fixture marker", got)
	}
	if _, err := os.Stat(h.KitPath("go.mod")); err != nil {
		t.Fatalf("KitPath(go.mod) did not resolve inside the kit root: %v", err)
	}
}

func TestHarnessDefaultsToCurrentGitRoot(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_CONFORMANCE_ROOT", "")
	t.Chdir(nested)

	h := NewHarness(t)

	if h.Root != root {
		t.Fatalf("Root = %q, want current git root %q", h.Root, root)
	}
	if h.KitRoot == "" || h.KitRoot == h.Root {
		t.Fatalf("KitRoot = %q, Root = %q; want distinct kit and graded roots", h.KitRoot, h.Root)
	}
}

func TestHarnessHelpersOperateAtGradedRoot(t *testing.T) {
	root := t.TempDir()
	t.Setenv("BENCH_CONFORMANCE_ROOT", root)
	if err := os.WriteFile(filepath.Join(root, "note.txt"), []byte("needle in root\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(root, "ok.sh")
	if err := os.WriteFile(script, []byte("#!/usr/bin/env bash\nprintf 'probe:%s\\n' \"$PWD\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	h := NewHarness(t)

	if got := h.RootPath("note.txt"); got != filepath.Join(root, "note.txt") {
		t.Fatalf("RootPath = %q, want root-relative path", got)
	}
	RequireSubstring(t, h.ReadRootFile("note.txt"), "needle", "root read")
	h.RequireExecutable("ok.sh")
	probe := h.Run("./ok.sh")
	if probe.ExitCode != 0 {
		t.Fatalf("probe ExitCode = %d, want 0; stderr=%q", probe.ExitCode, probe.Stderr)
	}
	RequireSubstring(t, probe.Stdout, "probe:"+root, "probe stdout")
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
