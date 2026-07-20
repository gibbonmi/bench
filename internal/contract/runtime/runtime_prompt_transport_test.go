package runtime

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

// Seam-B sentinel contracts: the prompt transport from the shift loop to its adapter and
// on to the harness CLI, observed at the built binary rather than at a Go-level unit.
// Each assertion plants a marker in the prompt and looks for it on the kernel's own view
// of a process's argv (/proc/<pid>/cmdline) versus its stdin, so a regression back to an
// argv prompt is caught by the gate rather than by a process listing. Every dump refuses
// to be read as proof unless the child actually ran with the real prompt on stdin.
func TestRuntimePromptTransportContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "stub adapter cmdline carries no prompt marker", testStubAdapterCmdlineNoMarker)
	contract.RunParallel(t, "stub adapter stdin carries the prompt unaltered", testStubAdapterStdinCarriesPrompt)
	contract.RunParallel(t, "claude adapter forwards stdin to the harness CLI", testClaudeAdapterForwardsStdin)
	contract.RunParallel(t, "codex adapter forwards stdin to the harness CLI", testCodexAdapterForwardsStdin)
	contract.RunParallel(t, "opencode adapter passes the prompt positionally", testOpencodeAdapterPassesPositional)
	contract.RunParallel(t, "argv sentinel runs in the gate's contract phase", testArgvSentinelIsGateAttached)
}

// promptDumpFixture builds a shift-capable repo whose stub adapter records its own
// /proc/<pid>/cmdline and its full stdin to two absolute paths outside the repo, baked
// into the script rather than passed through the environment (a path read from the
// environment would itself be subject to the adapter's env passlist). The adapter makes
// no change, so the loop's honest taxonomy is no-op/4.
func promptDumpFixture(t *testing.T) (f contract.Fixture, cmdline, stdin string) {
	t.Helper()
	f = shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	cmdline = filepath.Join(t.TempDir(), "adapter-cmdline")
	stdin = filepath.Join(t.TempDir(), "adapter-stdin")
	f.WriteExecutable("agent", fmt.Sprintf("#!/usr/bin/env bash\ncat /proc/$$/cmdline > %q\ncat > %q\nexit 0\n", cmdline, stdin))
	f.CommitAll("seam-B adapter")
	return f, cmdline, stdin
}

// runShiftPrompt drives one capped shift iteration with objective as the prompt's
// operator text. The stub adapter writes its cmdline and stdin dumps as a side effect.
func runShiftPrompt(t *testing.T, f contract.Fixture, objective string) {
	t.Helper()
	home := t.TempDir()
	env := map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}
	f.BenchEnv(env, "shift", objective).RequireExit(4)
}

// requirePromptChildRan refuses to let an absence assertion pass unless the stub adapter
// actually launched with the real prompt: the stdin dump must exist and carry the loop's
// iteration-prompt sentinel. Without this, deleting the argument alone would green a
// cmdline-absence assertion against a subprocess that never received a prompt at all.
func requirePromptChildRan(t *testing.T, stdinPath string) string {
	t.Helper()
	data, err := os.ReadFile(stdinPath)
	if err != nil {
		t.Fatalf("adapter wrote no stdin dump at %s: %v — the sentinel proves nothing unless the child ran", stdinPath, err)
	}
	text := string(data)
	if !strings.Contains(text, "You are one iteration of a Bench shift") {
		t.Fatalf("adapter stdin does not carry the iteration prompt: %q — the child did not receive the prompt on stdin", text)
	}
	return text
}

func readDump(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read dump %s: %v", path, err)
	}
	return string(data)
}

// testStubAdapterCmdlineNoMarker is the story-6/7 loop-side claim: the prompt marker
// never appears on the adapter's argv as the kernel reports it. Red today because the
// loop passes the prompt as argv[1], which /proc/<pid>/cmdline exposes verbatim.
func testStubAdapterCmdlineNoMarker(t *testing.T) {
	f, cmdline, stdin := promptDumpFixture(t)
	marker := "ft88-prompt-marker-in-cmdline"
	runShiftPrompt(t, f, marker)
	requirePromptChildRan(t, stdin)
	if cm := readDump(t, cmdline); strings.Contains(cm, marker) {
		t.Fatalf("adapter /proc/<pid>/cmdline carries the prompt marker %q — the prompt is on argv:\n%q", marker, cm)
	}
}

// testStubAdapterStdinCarriesPrompt pins that the argument was not merely dropped: the
// full multi-line prompt arrives on stdin byte-faithfully, carrying the operator marker
// verbatim. Folded in is the edge that a prompt with a newline and a leading dash arrives
// unaltered — the objective is dash-leading and the composed prompt is inherently
// multi-line. Red against an implementation that removes argv without wiring stdin (the
// dump would be empty), and against one that line-buffers or truncates the prompt.
func testStubAdapterStdinCarriesPrompt(t *testing.T) {
	f, cmdline, stdin := promptDumpFixture(t)
	marker := "-ft88-prompt-marker-leading-dash"
	runShiftPrompt(t, f, marker)
	text := requirePromptChildRan(t, stdin)
	// The objective marker, leading dash and all, appears verbatim on its own prompt line.
	if !strings.Contains(text, "Objective: "+marker+"\n") {
		t.Fatalf("adapter stdin does not carry the dash-leading objective marker verbatim:\n%q", text)
	}
	// The whole multi-line prompt arrived, not just its first line.
	if !strings.Contains(text, "decides if it counts") {
		t.Fatalf("adapter stdin is truncated — the full multi-line prompt did not arrive:\n%q", text)
	}
	if !strings.Contains(text, "\n") {
		t.Fatalf("adapter stdin carries no newline — the multi-line prompt was flattened:\n%q", text)
	}
	if cm := readDump(t, cmdline); strings.Contains(cm, marker) {
		t.Fatalf("adapter cmdline carries the prompt marker %q — the prompt leaked to argv:\n%q", marker, cm)
	}
}

// stubCLIRun runs the real shipped adapter directly with a prompt on stdin and a stub
// harness CLI first on PATH. The stub records its own /proc/<pid>/cmdline and the prompt
// it actually resolves. It emulates the real CLIs' contract that a positional prompt
// argument overrides stdin — the argument after a `--` wins over piped input, exactly as
// `claude -p`/`codex exec`/`opencode run` behave. That emulation is load-bearing: it is
// what makes "adapter left reading $1" red, since such an adapter passes an (empty)
// positional prompt that a stub blindly reading stdin would fail to notice. A stub
// `bench` answers `resolve-model` empty so the adapter takes its no-model branch; cwd is a
// throwaway non-repo directory so the wrapper resolver falls through to `bench` on PATH.
func stubCLIRun(t *testing.T, adapter, cli, prompt string) (cliCmdline, cliResolved string) {
	t.Helper()
	kit := contract.KitRoot(t)
	adapterPath := filepath.Join(kit, ".bench", "adapters", adapter)
	stubDir := t.TempDir()
	cmdlineFile := filepath.Join(t.TempDir(), "cli-cmdline")
	resolvedFile := filepath.Join(t.TempDir(), "cli-resolved")
	stubBody := fmt.Sprintf(`#!/usr/bin/env bash
cat /proc/$$/cmdline > %q
_p=""
_saw=0
_pos=0
for _a in "$@"; do
  if [ "$_saw" = 1 ]; then _p="$_a"; _pos=1; break; fi
  if [ "$_a" = "--" ]; then _saw=1; fi
done
if [ "$_pos" = 1 ]; then
  printf '%%s' "$_p" > %q
else
  cat > %q
fi
exit 0
`, cmdlineFile, resolvedFile, resolvedFile)
	writeStubExec(t, filepath.Join(stubDir, cli), stubBody)
	writeStubExec(t, filepath.Join(stubDir, "bench"), "#!/usr/bin/env bash\nexit 0\n")

	cwd := t.TempDir()
	cmd := exec.Command("bash", adapterPath)
	cmd.Dir = cwd
	cmd.Env = []string{"PATH=" + stubDir + ":/usr/bin:/bin", "HOME=" + t.TempDir()}
	cmd.Stdin = strings.NewReader(prompt)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("real %s adapter failed against the stub CLI: %v\n%s", adapter, err, out.String())
	}
	// The stub CLI must have run — otherwise an absence assertion proves nothing.
	if _, err := os.Stat(resolvedFile); err != nil {
		t.Fatalf("stub %s CLI never ran (no resolved-prompt dump): %v\n%s", cli, err, out.String())
	}
	return readDump(t, cmdlineFile), readDump(t, resolvedFile)
}

func writeStubExec(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write stub %s: %v", path, err)
	}
}

// cliProbePrompt is a multi-line, dash-leading prompt: the two properties the removed
// `--` argv guard existed for, now asserted to survive the stdin hop unaltered.
const cliProbePrompt = "-leading-dash first line\nsecond line ft88-cli-marker\n"

// testClaudeAdapterForwardsStdin proves the final hop for claude: the prompt reaches
// `claude -p` on its stdin, byte for byte, and never on argv. Red against an adapter left
// reading $1 (the CLI's stdin would be empty) and against one that leaves the prompt
// positional (the marker would appear in the CLI's cmdline).
func testClaudeAdapterForwardsStdin(t *testing.T) {
	cmdline, resolved := stubCLIRun(t, "claude", "claude", cliProbePrompt)
	if resolved != cliProbePrompt {
		t.Fatalf("claude adapter did not deliver the prompt to `claude -p` via stdin byte for byte:\ngot  %q\nwant %q", resolved, cliProbePrompt)
	}
	if strings.Contains(cmdline, "ft88-cli-marker") {
		t.Fatalf("the harness CLI's own cmdline carries the prompt marker — the prompt leaked to argv:\n%q", cmdline)
	}
}

// testCodexAdapterForwardsStdin is the same final-hop proof for codex exec, whose help
// documents reading instructions from stdin when none are given as an argument.
func testCodexAdapterForwardsStdin(t *testing.T) {
	cmdline, resolved := stubCLIRun(t, "codex", "codex", cliProbePrompt)
	if resolved != cliProbePrompt {
		t.Fatalf("codex adapter did not deliver the prompt to `codex exec` via stdin byte for byte:\ngot  %q\nwant %q", resolved, cliProbePrompt)
	}
	if strings.Contains(cmdline, "ft88-cli-marker") {
		t.Fatalf("the harness CLI's own cmdline carries the prompt marker — the prompt leaked to argv:\n%q", cmdline)
	}
}

// testOpencodeAdapterPassesPositional pins the documented residual: opencode's CLI takes
// only a positional prompt, so the adapter reads stdin and passes it to `opencode run`
// positionally. This asserts the known limitation deliberately — red against an adapter
// changed to a stdin form the CLI does not document. Command substitution strips the
// prompt's trailing newline, so the positional is compared minus that byte.
func testOpencodeAdapterPassesPositional(t *testing.T) {
	cmdline, resolved := stubCLIRun(t, "opencode", "opencode", cliProbePrompt)
	want := strings.TrimRight(cliProbePrompt, "\n")
	if resolved != want {
		t.Fatalf("opencode adapter did not pass the prompt to `opencode run` positionally:\ngot  %q\nwant %q", resolved, want)
	}
	// The residual: the prompt DOES appear on the CLI's argv. Pinned so a change either way
	// (to a stdin form the CLI does not document, or a silent drift) is a deliberate edit.
	if !strings.Contains(cmdline, "ft88-cli-marker") {
		t.Fatalf("opencode adapter no longer passes the prompt on argv — the documented residual changed silently:\n%q", cmdline)
	}
}

// testArgvSentinelIsGateAttached is the story-9 claim that the transport proof is the
// gate's, not a human's: the file it lives in must sit under a package the gate's contract
// phase actually executes. It resolves this file's own package rather than naming it, so
// moving the sentinel out of the phase's reach turns it red.
func testArgvSentinelIsGateAttached(t *testing.T) {
	if _, err := os.Stat(filepath.Join(sentinelPackageDir(t), "runtime_prompt_transport_test.go")); err != nil {
		t.Fatalf("argv sentinel file is not in the resolved sentinel package: %v", err)
	}
	kit := contract.KitRoot(t)
	pkg, err := filepath.Rel(kit, sentinelPackageDir(t))
	if err != nil {
		t.Fatalf("locate the sentinel package under the kit root: %v", err)
	}
	pkg = filepath.ToSlash(pkg)
	for _, phase := range gate.BenchkitPhases(kit, kit) {
		if phase.Name != "contract" {
			continue
		}
		for _, arg := range phase.Argv {
			if packagePatternCovers(arg, pkg) {
				return
			}
		}
		t.Fatalf("the gate's contract phase does not run the prompt-transport sentinel's package %q: argv %#v", pkg, phase.Argv)
	}
	t.Fatal("the gate has no contract phase, so the prompt-transport sentinel is not gate-attached")
}

// requireFileNotContains fails when path contains needle — the negative of
// requireFileContains, used to assert an adapter no longer carries the argv prompt form.
func requireFileNotContains(t *testing.T, path, needle, msg string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if strings.Contains(string(data), needle) {
		t.Fatal(msg)
	}
}
