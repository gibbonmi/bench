//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/toon"
)

func TestCompatibilityMissingPath(t *testing.T) {
	result, repo, wrapper := runCompatibilityDoctor(t)
	if result.code != 1 {
		t.Fatalf("doctor --compat without global Bench = %d, want 1; stdout=%q stderr=%q", result.code, result.stdout, result.stderr)
	}
	for _, want := range []string{"global-bench,failed", wrapper, filepath.Clean(repo)} {
		if !strings.Contains(result.stdout, want) {
			t.Fatalf("doctor --compat output missing %q: %q", want, result.stdout)
		}
	}
}

func TestCompatibilityPaths(t *testing.T) {
	result, repo, _ := runCompatibilityDoctor(t)
	expected, err := toon.Table("context", []string{"field", "value", "source"}, [][]string{{"repository", filepath.Clean(repo), "git root"}})
	if err != nil {
		t.Fatal(err)
	}
	_, row, found := strings.Cut(expected, "\n")
	if !found || !strings.Contains(result.stdout, row) {
		t.Fatalf("doctor --compat did not reach repository %q: (%d, %q, %q)", repo, result.code, result.stdout, result.stderr)
	}
}

func runCompatibilityDoctor(t *testing.T) (processResult, string, string) {
	t.Helper()
	repo := owner.repos[0]
	if result := owner.runSelected(repo, "link", "copy"); result.code != 0 {
		t.Fatalf("link fixture = %d; stdout=%q stderr=%q", result.code, result.stdout, result.stderr)
	}
	working := filepath.Join(repo, "path with spaces [and] glob*")
	if err := os.MkdirAll(working, 0o755); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	configHome := filepath.Join(home, "configuration home [desktop]*")
	if err := os.MkdirAll(configHome, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configHome, "config.toml"), []byte("model = \"personal\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(repo, ".bench", "bin", "bench.sh")
	overrides := []string{
		"BENCH_HOME=" + filepath.Join(home, ".bench"),
		"BENCH_KIT=" + owner.kit,
		"BENCH_RUN_BINARY=" + owner.selected.path,
		"CODEX_HOME=" + configHome,
		"HOME=" + home,
		"PATH=" + privateToolPath(t, "git", "bash", "uname", "dirname", "basename", "readlink", "tr"),
	}
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runAt(working, overrides, "bash", wrapper, "doctor", "--compat", "codex-cli"), repo, wrapper
}
