package conformance

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess"
)

// Harness is the shared test-only conformance fixture. Root holds the tree under grading.
// KitRoot holds the real Bench kit checkout that owns these helpers.
type Harness struct {
	t       testing.TB
	Root    string
	KitRoot string
}

type Probe struct {
	Args     []string
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func NewHarness(t testing.TB) Harness {
	t.Helper()
	kitRoot, err := findKitRoot()
	if err != nil {
		t.Fatalf("resolve kit root: %v", err)
	}
	root, err := resolveGradedRoot()
	if err != nil {
		t.Fatalf("resolve graded root: %v", err)
	}
	return Harness{t: t, Root: root, KitRoot: kitRoot}
}

func (h Harness) RootPath(elem ...string) string {
	return join(h.Root, elem...)
}

func (h Harness) KitPath(elem ...string) string {
	return join(h.KitRoot, elem...)
}

func (h Harness) ReadRootFile(elem ...string) string {
	h.t.Helper()
	data, err := os.ReadFile(h.RootPath(elem...))
	if err != nil {
		h.t.Fatalf("read %s: %v", filepath.Join(elem...), err)
	}
	return string(data)
}

func (h Harness) RequireExecutable(elem ...string) {
	h.t.Helper()
	path := h.RootPath(elem...)
	info, err := os.Stat(path)
	if err != nil {
		h.t.Fatalf("stat executable %s: %v", path, err)
	}
	if info.IsDir() || info.Mode()&0o111 == 0 {
		h.t.Fatalf("%s is not an executable file", path)
	}
}

func (h Harness) Run(args ...string) Probe {
	h.t.Helper()
	if len(args) == 0 {
		h.t.Fatal("Run requires at least one argv element")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = h.Root
	r := subprocess.Capture(cmd)
	stderr := r.Stderr
	// Surface a spawn failure (e.g. command not found) that wrote nothing itself.
	var exitErr *exec.ExitError
	if r.Err != nil && !errors.As(r.Err, &exitErr) && stderr == "" {
		stderr = r.Err.Error()
	}
	return Probe{
		Args:     append([]string(nil), args...),
		ExitCode: r.ExitCode,
		Stdout:   r.Stdout,
		Stderr:   stderr,
		Err:      r.Err,
	}
}

func RequireSubstring(t testing.TB, got, want, label string) {
	t.Helper()
	if want == "" {
		t.Fatalf("%s: empty expected substring", label)
	}
	if !strings.Contains(got, want) {
		t.Fatalf("%s: missing substring %q in %q", label, want, got)
	}
}

func resolveGradedRoot() (string, error) {
	if root := os.Getenv("BENCH_CONFORMANCE_ROOT"); root != "" {
		return filepath.Abs(root)
	}
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func findKitRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime caller unavailable")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod not found above conformance helpers")
		}
		dir = parent
	}
}

func join(root string, elem ...string) string {
	parts := append([]string{root}, elem...)
	return filepath.Join(parts...)
}

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

// throwawayRoot is the shared throwaway-tree builder for this package's fixtures. It owns
// the knowledge every per-check builder used to repeat: a temp root, a write loop that
// creates each parent directory, and the git init a root needs before a timing file has
// anywhere to live. A caller supplies only the content its own check grades.
type throwawayRoot struct {
	// gitInit runs `git init -q` in the root.
	gitInit bool
	// files maps a slash-relative path to the content planted there.
	files map[string]string
	// plants maps a slash-relative path to a planter the builder hands the absolute path,
	// for a fixture whose bytes a hostile-input helper writes rather than this builder.
	// Two readers naming one path collapse to a single plant, as the path is the key.
	plants map[string]func(*testing.T, string)
}

func (spec throwawayRoot) build(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if spec.gitInit {
		runGit(t, root, "init", "-q")
	}
	for rel, content := range spec.files {
		if err := os.WriteFile(rootFilePath(t, root, rel), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for rel, plant := range spec.plants {
		plant(t, rootFilePath(t, root, rel))
	}
	return root
}

// rootFilePath resolves a slash-relative path under root and creates its parent directory.
func rootFilePath(t *testing.T, root, rel string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}
