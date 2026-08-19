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
	if rc := (Command{}).Run([]string{"version"}); rc != 0 {
		t.Errorf("run version exit = %d, want 0", rc)
	}
}

func TestRunUnknownExits2(t *testing.T) {
	if rc := (Command{}).Run([]string{"nope"}); rc != 2 {
		t.Errorf("run nope exit = %d, want 2", rc)
	}
}

func TestRunStatusRouteEmitsOneNextRow(t *testing.T) {
	stdout := tempFile(t)
	if code := (Command{Stdout: stdout}).Run([]string{"status", "--route"}); code != 0 {
		t.Fatalf("status --route exit = %d, want 0", code)
	}
	if got := readFile(t, stdout); !strings.HasPrefix(got, "next[1]{state,why,command}:\n") {
		t.Fatalf("status --route = %q, want one next row", got)
	}
}

func TestHelpRendersPublicCommandRegistryRows(t *testing.T) {
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	commandRegistry = []commandDefinition{
		{Name: "help", Inventory: publicInventory(), Kind: commandHelp},
		{
			Name:      "status",
			Inventory: publicInventory(helpRow{Order: 1, Description: "prove the public command owns this row"}),
		},
	}

	for _, spelling := range []string{"help", "--help", "-h"} {
		t.Run(spelling, func(t *testing.T) {
			var stdout bytes.Buffer
			if code := (Command{Stdout: &stdout}).Run([]string{spelling}); code != 0 {
				t.Fatalf("%s exit = %d, want 0", spelling, code)
			}
			if want := helpInventoryTitle + "\n  bench status               prove the public command owns this row\n"; stdout.String() != want {
				t.Fatalf("%s stdout = %q, want registry rows %q", spelling, stdout.String(), want)
			}
		})
	}
}

func TestRootAndHelpAlignWrapperAndBinary(t *testing.T) {
	var directRoot bytes.Buffer
	if code := (Command{Stdout: &directRoot}).Run(nil); code != 0 {
		t.Fatalf("in-process root exit = %d, want 0", code)
	}
	if !strings.HasPrefix(directRoot.String(), "next[1]{state,why,command}:\n") {
		t.Fatalf("in-process root = %q, want next route table", directRoot.String())
	}

	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), root, binary)
	cleanEnv := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	build.Env = cleanEnv
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build bench: %v\n%s", err, out)
	}

	command := func(path string, args ...string) *exec.Cmd {
		t.Helper()
		cmd := exec.Command(path, args...)
		cmd.Env = cleanEnv
		if path != binary {
			cmd.Env = append(capability.WithoutEnvironment(runbinary.WithEnv(cleanEnv, binary), "BENCH_KIT"), "BENCH_KIT="+root)
		}
		return cmd
	}
	run := func(path string, args ...string) string {
		t.Helper()
		cmd := command(path, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", path, args, err, out)
		}
		return string(out)
	}

	binaryRoot := run(binary)
	if !strings.HasPrefix(binaryRoot, "next[1]{state,why,command}:\n") {
		t.Fatalf("binary root = %q, want next route table", binaryRoot)
	}
	if binaryRoot != directRoot.String() {
		t.Fatalf("binary root = %q, in-process root = %q", binaryRoot, directRoot.String())
	}
	wrapper := filepath.Join(root, "bin", "bench.sh")
	if wrapperRoot := run(wrapper); wrapperRoot != binaryRoot {
		t.Fatalf("wrapper root = %q, binary root = %q", wrapperRoot, binaryRoot)
	}

	binaryHelp := run(binary, "help")
	if !strings.HasPrefix(binaryHelp, "bench — Pocock"+" pipeline") {
		t.Fatalf("binary help = %q, want inventory", binaryHelp)
	}
	if !strings.Contains(binaryHelp, "bench canary [root]        validate fixture inventory") {
		t.Fatalf("binary help missing canary inventory wording:\n%s", binaryHelp)
	}
	if strings.Contains(binaryHelp, "run the gate against known-broken fixtures") {
		t.Fatalf("binary help retained stale canary execution wording:\n%s", binaryHelp)
	}
	for _, spelling := range []string{"help", "--help", "-h"} {
		if spellingHelp := run(binary, spelling); spellingHelp != binaryHelp {
			t.Errorf("binary %s = %q, binary help = %q", spelling, spellingHelp, binaryHelp)
		}
		if wrapperHelp := run(wrapper, spelling); wrapperHelp != binaryHelp {
			t.Errorf("wrapper %s = %q, binary help = %q", spelling, wrapperHelp, binaryHelp)
		}
	}

	cmd := command(wrapper, "help", "extra")
	out, err := cmd.CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 2 || string(out) != "usage: bench help (unknown argument: extra)\n" {
		t.Fatalf("wrapper help extra = (output %q, error %v), want help usage and exit 2", out, err)
	}
}

func TestRunHelpRejectsTrailingArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"help", "extra"})
	if code != 2 || stdout.String() != "usage: bench help (unknown argument: extra)\n" || stderr.Len() != 0 {
		t.Fatalf("help extra = (stdout %q, stderr %q, exit %d), want help usage on stdout and exit 2", stdout.String(), stderr.String(), code)
	}
}

func TestHelpKeepsStatusPublicRoute(t *testing.T) {
	t.Run("inventory", func(t *testing.T) {
		var stdout bytes.Buffer
		if code := (Command{Stdout: &stdout}).Run([]string{"help"}); code != 0 {
			t.Fatalf("help exit = %d, want 0", code)
		}
		const row = "  bench status               ambient dashboard: what needs attention + the next action\n"
		if !strings.Contains(stdout.String(), row) {
			t.Fatalf("help omitted independently required public status row %q", row)
		}
	})
	t.Run("dispatch", func(t *testing.T) {
		var stdout bytes.Buffer
		if code := (Command{Stdout: &stdout}).Run([]string{"status", "--help"}); code != 0 {
			t.Fatalf("status --help exit = %d, want 0", code)
		}
		if !strings.HasPrefix(stdout.String(), "usage: bench status") {
			t.Fatalf("status --help = %q, want status grammar", stdout.String())
		}
	})
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

func TestRunCanaryRetainsPositionalGrammar(t *testing.T) {
	for _, help := range []string{"help", "--help", "-h"} {
		stdout, stderr := tempFile(t), tempFile(t)
		if rc := (Command{Stdout: stdout, Stderr: stderr}).Run([]string{"canary", help}); rc != 0 {
			t.Errorf("canary %s exit = %d, want 0", help, rc)
		}
		if got := readFile(t, stdout); got != "usage: bench canary [root]\n" {
			t.Errorf("canary %s stdout = %q, want exact usage", help, got)
		}
	}

	stderr := tempFile(t)
	if rc := (Command{Stderr: stderr}).Run([]string{"canary", "one", "two"}); rc != 2 {
		t.Fatalf("run canary too-many-arguments exit = %d, want 2", rc)
	}
	if got := readFile(t, stderr); !strings.Contains(got, "usage: bench canary") || !strings.Contains(got, "unknown argument: two") {
		t.Fatalf("canary too-many-arguments stderr = %q, want usage and offending argument", got)
	}

	stderr = tempFile(t)
	if rc := (Command{Stderr: stderr}).Run([]string{"canary", filepath.Join(t.TempDir(), "missing")}); rc != 1 {
		t.Fatalf("run canary invalid root exit = %d, want 1", rc)
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
	rc := (Command{Stdout: stdout, Stderr: stderr}).Run([]string{"gate-phases", "/tmp/root"})

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
