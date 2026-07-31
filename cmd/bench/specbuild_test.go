package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpecBuildStatusRendersDefinitiveEmptyProjection(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	if err := os.MkdirAll(filepath.Join(root, "specs", "demo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "demo", "spec.md"), []byte("Status: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	out, code := specBuildCommand([]string{"status", "demo"})
	if code != 0 || !strings.Contains(out, "spec_build[1]{slug,state,subject,next}:") || !strings.Contains(out, "demo,empty") {
		t.Fatalf("status = %q, %d", out, code)
	}
	full, fullCode := specBuildCommand([]string{"status", "demo", "--full"})
	if fullCode != 0 || !strings.Contains(full, "assignments[0]") || !strings.Contains(full, "review[0]") {
		t.Fatalf("full status = %q, %d", full, fullCode)
	}
}

func TestSpecBuildErrorCannotSplitStructuredOutput(t *testing.T) {
	out, code := buildError(errors.New("ignored\nerror"), "retry\nnow")
	if code != 1 || strings.Count(out, "\n") != 1 || !strings.Contains(out, `retry\nnow`) {
		t.Fatalf("structured error = %q, %d", out, code)
	}
}
