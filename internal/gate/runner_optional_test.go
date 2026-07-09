package gate

// Phase start behavior: an Optional phase whose binary is truly absent skips green
// and says why, a present-but-broken one surfaces its real exec failure as red, and
// a required phase that cannot start is red.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunnerShellcheckAbsentSkips(t *testing.T) {
	root := t.TempDir()
	phase := Phase{Name: "shellcheck", Argv: []string{"definitely-not-installed-shellcheck-for-bench-test"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, want skip to stay green; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if strings.Contains(out, "gate: red") || strings.Contains(out, "shellcheck reported issues") {
		t.Fatalf("optional missing shellcheck looked red:\n%s", out)
	}
	if !strings.Contains(out, "phase shellcheck: skipped") {
		t.Fatalf("missing shellcheck skip summary:\n%s", out)
	}
}

func TestRunnerOptionalBrokenSymlinkSkips(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	interpreters := filepath.Join(root, "interpreters")
	if err := os.Mkdir(interpreters, 0o755); err != nil {
		t.Fatalf("mkdir interpreters: %v", err)
	}
	brokenInterpreter := filepath.Join(interpreters, "sh")
	if err := os.Symlink(filepath.Join(root, "missing-sh"), brokenInterpreter); err != nil {
		t.Fatalf("symlink broken interpreter: %v", err)
	}
	shellcheck := filepath.Join(bin, "shellcheck")
	if err := os.WriteFile(shellcheck, []byte("#!"+brokenInterpreter+"\n"), 0o755); err != nil {
		t.Fatalf("write shellcheck shim: %v", err)
	}
	t.Setenv("PATH", bin)
	phase := Phase{Name: "shellcheck", Argv: []string{"shellcheck"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("runPhases rc = %d, want broken optional executable skipped; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if out := stdout.String() + stderr.String(); !strings.Contains(out, "phase shellcheck: skipped") || strings.Contains(out, "gate: red") {
		t.Fatalf("broken optional executable did not produce a green skip:\n%s", out)
	}
}

func TestRunnerShellcheckAbsentSkipVerdictNamesNotInstalled(t *testing.T) {
	root := t.TempDir()
	phase := Phase{Name: "shellcheck", Argv: []string{"definitely-absent-shellcheck-binary-for-bench-test"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 0 {
		t.Fatalf("absent optional binary went red rc=%d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	if out := stdout.String() + stderr.String(); !strings.Contains(out, "phase shellcheck: skipped (not installed)") {
		t.Fatalf("absent shellcheck skip verdict is silent about why (want 'skipped (not installed)'):\n%s", out)
	}
}

func TestRunnerOptionalUnexecutableStubGoesRed(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	stub := filepath.Join(bin, "shellcheck")
	if err := os.WriteFile(stub, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o000); err != nil {
		t.Fatalf("write unexecutable stub: %v", err)
	}
	t.Setenv("PATH", bin)
	phase := Phase{Name: "shellcheck", Argv: []string{"shellcheck"}, Optional: true}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("present-but-unexecutable shellcheck did not go red rc=%d; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if strings.Contains(out, "skipped") {
		t.Fatalf("present-but-unexecutable shellcheck was masked as a skip:\n%s", out)
	}
	if !strings.Contains(out, "phase shellcheck: red") || !strings.Contains(out, "permission denied") {
		t.Fatalf("red verdict does not name the exec failure:\n%s", out)
	}
}

func TestRunnerRequiredStartFailureRed(t *testing.T) {
	root := t.TempDir()
	phase := Phase{Name: "required", Argv: []string{filepath.Join(root, "missing-required")}}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)
	if rc != 1 {
		t.Fatalf("runPhases rc = %d, want required start failure red; stdout=%q stderr=%q", rc, stdout.String(), stderr.String())
	}
	out := stdout.String() + stderr.String()
	if !strings.Contains(out, "phase required: red") || !strings.Contains(out, "gate: red") {
		t.Fatalf("required start failure did not stay red:\n%s", out)
	}
}
