package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// The idiom-setting table test for the module: pure logic, table-driven, no process
// boundary. Acceptance for the version line lives at the shell seam (the gate's
// version-routing contract); this pins the format one layer down.
func TestVersionLine(t *testing.T) {
	cases := []struct {
		v, goos, goarch, want string
	}{
		{"0.2.0", "linux", "amd64", "bench 0.2.0 (linux/amd64)"},
		{"dev", "darwin", "arm64", "bench dev (darwin/arm64)"},
		{"1.0.0", "linux", "arm64", "bench 1.0.0 (linux/arm64)"},
	}
	for _, c := range cases {
		if got := versionLine(c.v, c.goos, c.goarch); got != c.want {
			t.Errorf("versionLine(%q,%q,%q) = %q, want %q", c.v, c.goos, c.goarch, got, c.want)
		}
	}
}

func TestRunVersionExits0(t *testing.T) {
	if rc := run([]string{"version"}, nil, nil); rc != 0 {
		t.Errorf("run version exit = %d, want 0", rc)
	}
}

func TestRunUnknownExits2(t *testing.T) {
	if rc := run([]string{"nope"}, nil, nil); rc != 2 {
		t.Errorf("run nope exit = %d, want 2", rc)
	}
}

// TestResolveModelHarnessFlag drives the CLI's argument surface: --harness selects the
// column, and the retired --alias / --provider-model spellings are rejected rather than
// quietly resolving a model, so there is only one way to ask the binding a question.
func TestResolveModelHarnessFlag(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	binding := "BENCH_CODEX_TOP=gpt-5.6-sol\nBENCH_CODEX_MID=gpt-5.6-terra\nBENCH_CODEX_CHEAP=gpt-5.6-luna\n" +
		"BENCH_CLAUDE_TOP=fable\nBENCH_CLAUDE_MID=opus\nBENCH_CLAUDE_CHEAP=sonnet\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(binding), 0o644); err != nil {
		t.Fatal(err)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("BENCH_MODEL", "cheap")

	for _, tt := range []struct {
		name     string
		args     []string
		wantOut  string
		wantCode int
	}{
		{"codex column", []string{"--harness", "codex"}, "gpt-5.6-luna\n", 0},
		{"claude column", []string{"--harness", "claude"}, "sonnet\n", 0},
		{"unbound column", []string{"--harness", "opencode"}, "", 1},
		{"unknown harness", []string{"--harness", "gemini"}, "", 1},
		{"missing harness", nil, "", 2},
		{"retired alias flag", []string{"--alias"}, "", 2},
		{"retired provider-model flag", []string{"--provider-model"}, "", 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out, code := resolveModel(tt.args)
			if out != tt.wantOut || code != tt.wantCode {
				t.Fatalf("resolveModel(%v) = (%q, %d), want (%q, %d)", tt.args, out, code, tt.wantOut, tt.wantCode)
			}
		})
	}
}

// TestCheckAgentLineHarnessFlag pins the guard's own argument surface: it takes the same
// --harness flag, and a retired flag is a usage error rather than a silent allow.
func TestCheckAgentLineHarnessFlag(t *testing.T) {
	for _, tt := range []struct {
		name     string
		args     []string
		wantCode int
	}{
		{"missing harness", nil, 2},
		{"retired alias flag", []string{"--alias"}, 2},
		{"unknown harness", []string{"--harness", "gemini"}, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			code := checkAgentLine(tt.args, strings.NewReader(`{"tool_input":{"model":"opus"}}`), nil, &stderr)
			if code != tt.wantCode {
				t.Fatalf("checkAgentLine(%v) = %d, want %d (stderr=%q)", tt.args, code, tt.wantCode, stderr.String())
			}
		})
	}
}

func TestRunCanaryDispatchesToCommand(t *testing.T) {
	stderr := tempFile(t)
	if rc := run([]string{"canary", "one", "two"}, nil, stderr); rc != 2 {
		t.Fatalf("run canary usage exit = %d, want 2", rc)
	}
	got := readFile(t, stderr)
	if !strings.Contains(got, "usage: bench canary") {
		t.Fatalf("canary did not dispatch to command usage, stderr:\n%s", got)
	}
}

func TestRunGatePhasesDispatchesToCommand(t *testing.T) {
	old := gatePhasesCommand
	t.Cleanup(func() { gatePhasesCommand = old })
	var gotArgs []string
	gatePhasesCommand = func(args []string, stdout, stderr io.Writer) int {
		gotArgs = append([]string(nil), args...)
		return 37
	}

	stdout := tempFile(t)
	stderr := tempFile(t)
	rc := run([]string{"gate-phases", "/tmp/root"}, stdout, stderr)

	if rc != 37 {
		t.Fatalf("run gate-phases exit = %d, want injected exit 37", rc)
	}
	if want := []string{"/tmp/root"}; !reflect.DeepEqual(gotArgs, want) {
		t.Fatalf("gate-phases args = %#v, want %#v", gotArgs, want)
	}
}

func TestShellWrapperRoutesGatePhasesWithNonRootKit(t *testing.T) {
	root := t.TempDir()
	kit := filepath.Join(root, "kit")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))
	argvFile := filepath.Join(root, "argv")
	writeExecutable(t, filepath.Join(kit, "dist", "bench"), `#!/usr/bin/env bash
printf '%s\n' "$BENCH_KIT" "$@" > "$BENCH_TEST_ARGV"
`)

	cmd := exec.Command("bash", filepath.Join(kit, "bin", "bench.sh"), "gate-phases", "/tmp/repo root")
	cmd.Env = append(capability.WithoutEnvironment(os.Environ(), runbinary.Env),
		"BENCH_TEST_ARGV="+argvFile, "BENCH_HOME="+filepath.Join(root, "home"),
		"BENCH_KIT="+kit, "BENCH_WRAPPER=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bench.sh gate-phases failed: %v\n%s", err, out)
	}
	got := strings.Split(strings.TrimSpace(readPath(t, argvFile)), "\n")
	want := []string{kit, "gate-phases", "/tmp/repo root"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("wrapper routed argv = %#v, want %#v\noutput:\n%s", got, want, out)
	}
}

func TestGuardGitBlockAllow(t *testing.T) {
	var errb bytes.Buffer
	block := `{"tool_input":{"command":"git push"}}`
	if code := guardGit(nil, strings.NewReader(block), io.Discard, &errb); code != 2 {
		t.Errorf("block exit = %d, want 2", code)
	}
	if !strings.Contains(errb.String(), "BLOCKED:") {
		t.Errorf("block did not emit BLOCKED on stderr: %q", errb.String())
	}
	for _, in := range []string{`{"tool_input":{"command":"git status"}}`, "not json", `{"tool_input":{"command":""}}`} {
		if code := guardGit(nil, strings.NewReader(in), io.Discard, io.Discard); code != 0 {
			t.Errorf("allow exit for %q = %d, want 0", in, code)
		}
	}
}

func TestCaptureClaudeAgentIntentReplayIsByteIdempotent(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	envelope := []byte(`{"tool_name":"Agent","tool_use_id":"same-id","tool_input":{"description":"same objective"}}`)
	captureClaudeAgentIntent(envelope, io.Discard)
	path := filepath.Join(root, ".git", "bench-intent.json")
	before := readPath(t, path)
	captureClaudeAgentIntent(envelope, io.Discard)
	after := readPath(t, path)
	if after != before {
		t.Fatalf("replayed Claude capture changed ledger bytes:\nbefore=%s\nafter=%s", before, after)
	}
}

// panicReader forces guardGit's stdin read to panic, exercising the recover→exit-3
// rim so a crash can never masquerade as an exit-2 block.
type panicReader struct{}

func (panicReader) Read([]byte) (int, error) { panic("boom") }

func TestGuardGitRecoversToExit3(t *testing.T) {
	if code := guardGit(nil, panicReader{}, io.Discard, io.Discard); code != 3 {
		t.Errorf("panic mapped to exit %d, want 3", code)
	}
}

func tempFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "out-*")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func readFile(t *testing.T, f *os.File) string {
	t.Helper()
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func readPath(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func copyExecutable(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	writeExecutable(t, dst, string(data))
}

func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
