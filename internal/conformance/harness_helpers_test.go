package conformance

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Harness is the shared test-only conformance fixture: Root is the tree being
// graded, while KitRoot is the real Bench kit checkout that owns these helpers.
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
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else if stderr.Len() == 0 {
			stderr.WriteString(err.Error())
		}
	}
	return Probe{
		Args:     append([]string(nil), args...),
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Err:      err,
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
