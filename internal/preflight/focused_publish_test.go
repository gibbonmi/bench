package preflight

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuiltCommandFocusedPublishRunsDiagnosticWithoutAuthorizing(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(projectRoot(t), "scripts", "go-build.sh"), projectRoot(t), binary)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	root := preflightRepo(t)
	marker := filepath.Join(root, "focused-publish-ran")
	phase := filepath.Join(root, "focused-publish")
	if err := os.WriteFile(phase, []byte("#!/bin/sh\nprintf 'ran\\n' > \"$BENCH_FOCUSED_MARKER\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binary, "release-preflight", "--mode", "publish", "--profile", "public", "--phase", "gate")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "BENCH_PREFLIGHT_GATE="+phase, "BENCH_FOCUSED_MARKER="+marker)
	output, err := cmd.CombinedOutput()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		t.Fatalf("focused publish exit = %v, want non-authorizing exit 1\n%s", err, output)
	}
	if !strings.Contains(string(output), "focused publish runs cannot authorize publication") {
		t.Fatalf("focused publish output does not explain non-authorization:\n%s", output)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "ran\n" {
		t.Fatalf("focused publish diagnostic did not run: %q %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(root, "dist", "preflight", "release-index.json")); !os.IsNotExist(err) {
		t.Fatalf("focused publish created release evidence: %v", err)
	}
}
