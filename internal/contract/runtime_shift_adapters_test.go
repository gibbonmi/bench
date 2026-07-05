package contract

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRuntimeShiftAdapterContracts(t *testing.T) {
	t.Parallel()
	runParallel(t, "bench shift adapter preflight contract", testShiftAdapterPreflight)
	runParallel(t, "bench shift adapter single-argument contract", testShiftAdapterSingleArgument)
	runParallel(t, "reference adapter files contract", testReferenceAdapterFiles)
}

func testShiftAdapterPreflight(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	unset := f.BenchEnvSpec(Env{"BENCH_AGENT": nil, "BENCH_HOME": strPtr(home)}, "shift", "probe")
	if unset.ExitCode == 0 {
		t.Fatal("shift with no BENCH_AGENT succeeded; should error")
	}
	unset.RequireContains(unset.Stderr, "BENCH_AGENT")
	if !regexp.MustCompile(`(?i)configure.*adapter|adapter.*configure`).MatchString(unset.Stderr) {
		t.Fatalf("unconfigured-adapter error is not a configure-your-adapter message:\n%s", unset.Stderr)
	}
	unset.RequireNotContains(unset.Stdout, "iteration 1/")

	empty := f.BenchEnv(map[string]string{"BENCH_AGENT": "", "BENCH_HOME": home}, "shift", "probe")
	if empty.ExitCode == 0 {
		t.Fatal("shift with empty BENCH_AGENT succeeded; should error")
	}
	empty.RequireContains(empty.Stderr, "BENCH_AGENT")
	if !regexp.MustCompile(`(?i)configure.*adapter|adapter.*configure`).MatchString(empty.Stderr) {
		t.Fatalf("empty-adapter error is not a configure-your-adapter message:\n%s", empty.Stderr)
	}
	empty.RequireNotContains(empty.Stdout, "iteration 1/")

	missing := f.BenchEnv(map[string]string{"BENCH_AGENT": "/no/such/adapter", "BENCH_HOME": home}, "shift", "probe")
	if missing.ExitCode == 0 {
		t.Fatal("shift with a missing adapter path succeeded; should error")
	}
	missing.RequireContains(missing.Stderr, "not executable")
	missing.RequireNotContains(missing.Stdout, "iteration 1/")

	keyword := f.BenchEnv(map[string]string{"BENCH_AGENT": "if", "BENCH_HOME": home}, "shift", "probe")
	if keyword.ExitCode == 0 {
		t.Fatal("shift with a shell-keyword adapter succeeded; should error")
	}
	keyword.RequireContains(keyword.Stderr, "not executable")
}

func testShiftAdapterSingleArgument(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("adapter", `#!/usr/bin/env bash
{
  printf 'argc=%s\n' "$#"
  printf 'shift_env=%s\n' "${BENCH_SHIFT:-unset}"
  printf '%s\n@@@@\n' "$1"
} >> "$BENCH_TEST_RECORD"
`)
	f.CommitAll("adapter")
	home := t.TempDir()
	record := filepath.Join(t.TempDir(), "record.txt")

	f.BenchEnv(map[string]string{"BENCH_TEST_RECORD": record, "BENCH_AGENT": filepath.Join(f.Root, "adapter"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "adapter-arg-probe").RequireExit(0)

	data, err := os.ReadFile(record)
	if err != nil {
		t.Fatalf("adapter was never invoked: %v", err)
	}
	text := string(data)
	for _, needle := range []string{
		"argc=1",
		"shift_env=1",
		"adapter-arg-probe",
		"You are one iteration of a Bench shift",
		"decides if it counts",
	} {
		if !strings.Contains(text, needle) {
			t.Fatalf("adapter record missing %q:\n%s", needle, text)
		}
	}
	if regexp.MustCompile(`(?m)^-p$`).MatchString(text) {
		t.Fatal("loop still passes the Claude-specific -p flag")
	}
}

func testReferenceAdapterFiles(t *testing.T) {
	root := KitRoot(t)
	for _, adapter := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", adapter)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("reference adapter missing: .bench/adapters/%s", adapter)
		}
		if info.Mode()&0o111 == 0 {
			t.Fatalf("reference adapter not executable: .bench/adapters/%s", adapter)
		}
		probe := NewFixture(t, WithNoRepo()).Run("bash", "-n", path)
		probe.RequireExit(0)
		text, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read adapter %s: %v", adapter, err)
		}
		if !regexp.MustCompile(`(?m)^exec `).Match(text) {
			t.Fatalf("reference adapter %s does not exec its harness", adapter)
		}
		if !strings.Contains(string(text), `"$1"`) {
			t.Fatalf("reference adapter %s does not pass the prompt as $1", adapter)
		}
	}
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "claude"), `claude -p -- "$1"`, "claude adapter does not map the prompt to claude -p")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "codex"), `codex exec -- "$1"`, "codex adapter does not map the prompt to codex exec")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `opencode run -- "$1"`, "opencode adapter does not map the prompt to opencode run")
}
