package runtime

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestRuntimeFreshInstall(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "fresh-install PATH-stripped gate resolves local CLI", testRuntimeFreshInstallGateResolves)
}

func TestRuntimeAgentEntryContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "bench shift adapter preflight contract", testShiftAdapterPreflight)
	contract.RunParallel(t, "bench shift adapter stdin-transport contract", testShiftAdapterStdinTransport)
	contract.RunParallel(t, "reference adapter files contract", testReferenceAdapterFiles)
}

// testRuntimeFreshInstallGateResolves is the only test that installs cold with no global
// bench on PATH — the class the kit otherwise can't see. It links and inits a throwaway
// repo, then runs the scaffolded .bench/gate.sh, all under a PATH stripped to
// /usr/bin:/bin (git and bash present, no global bench). The PATH strip is load-bearing:
// a smoke that leaves a global bench on PATH false-greens against the old bare-`bench`
// scaffold, so it is asserted as part of the fixture, not left implicit.
func testRuntimeFreshInstallGateResolves(t *testing.T) {
	repo := t.TempDir()
	f := contract.NewFixtureAt(t, repo, contract.IsolatedEnv(t, repo))
	f.Git("init", "-q")

	stripped := map[string]string{"PATH": "/usr/bin:/bin"}
	contract.RunAt(t, f, repo, selectedBenchEnv(t, stripped), benchPath(t), "link").RequireExit(0)
	contract.RunAt(t, f, repo, selectedBenchEnv(t, stripped), benchPath(t), "init").RequireExit(0)

	// The fresh install must land the local CLI and the copied binary the scaffolded gate
	// resolves, so the resolution assertion below tests resolution — not mere absence.
	requireFreshFile(t, filepath.Join(repo, ".bench", "bin", "bench.sh"), "fresh link did not install .bench/bin/bench.sh")
	requireFreshFile(t, filepath.Join(repo, ".bench", "dist", "bench"), "fresh link did not copy .bench/dist/bench")

	stripped[runbinary.Env] = filepath.Join(repo, ".bench", "dist", "bench")
	gate := contract.RunAt(t, f, repo, stripped, "bash", filepath.Join(repo, ".bench", "gate.sh"))
	out := gate.Stdout + gate.Stderr
	// The scaffolded gate is intentionally red on its sentinel; the contract is that it
	// RESOLVES the local CLI and reaches canary, not that it goes green. A bare-`bench`
	// scaffold emits "bench: command not found" and never runs the canary under this PATH.
	contract.RequireNotContains(t, out, "command not found")
	contract.RequireNotContains(t, out, "canary sweep failed")

	// story 5: the profile templates the fresh install points at are reachable in the
	// resolved kit — the ships-regression itself is owned by the package-surface pin.
	requireFreshFile(t, filepath.Join(contract.SubjectRoot(t), "projects", "gl-axi.md"), "resolved kit dir lacks projects/ profile templates")
}

func requireFreshFile(t *testing.T, path, msg string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func testShiftAdapterPreflight(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	unset := f.BenchEnvSpec(contract.Env{"BENCH_AGENT": nil, "BENCH_HOME": strPtr(home)}, "shift", "probe")
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

// testShiftAdapterStdinTransport pins the loop → adapter contract: the prompt reaches
// the adapter on stdin with no positional argument, BENCH_SHIFT armed. The probe records
// its own argc separately from its stdin, so a regression that put the prompt back on
// argv fails the argc=0 assertion even while the stdin assertions still pass.
func testShiftAdapterStdinTransport(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("adapter", `#!/usr/bin/env bash
{
  printf 'argc=%s\n' "$#"
  printf 'shift_env=%s\n' "${BENCH_SHIFT:-unset}"
} >> "$BENCH_TEST_RECORD"
{ cat; printf '\n@@@@\n'; } >> "$BENCH_TEST_RECORD"
`)
	f.CommitAll("adapter")
	home := t.TempDir()
	record := filepath.Join(t.TempDir(), "record.txt")

	// The probe adapter only records its invocation; it never mutates the tree, so the
	// honest taxonomy is no-op/4, not complete/0.
	f.BenchEnv(map[string]string{"BENCH_TEST_RECORD": record, "BENCH_AGENT": filepath.Join(f.Root, "adapter"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "adapter-arg-probe").RequireExit(4)

	data, err := os.ReadFile(record)
	if err != nil {
		t.Fatalf("adapter was never invoked: %v", err)
	}
	text := string(data)
	for _, needle := range []string{
		"argc=0",
		"shift_env=1",
		"adapter-arg-probe",
		"You are one iteration of a Bench shift",
		"decides if it counts",
	} {
		if !strings.Contains(text, needle) {
			t.Fatalf("adapter record missing %q:\n%s", needle, text)
		}
	}
	if strings.Contains(text, "argc=1") {
		t.Fatalf("loop still passes the prompt as a positional argument:\n%s", text)
	}
}

func testReferenceAdapterFiles(t *testing.T) {
	root := contract.KitRoot(t)
	for _, adapter := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", adapter)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("reference adapter missing: .bench/adapters/%s", adapter)
		}
		if info.Mode()&0o111 == 0 {
			t.Fatalf("reference adapter not executable: .bench/adapters/%s", adapter)
		}
		probe := contract.NewFixture(t, contract.WithNoRepo()).Run("bash", "-n", path)
		probe.RequireExit(0)
		text, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read adapter %s: %v", adapter, err)
		}
		if !regexp.MustCompile(`(?m)^exec `).Match(text) {
			t.Fatalf("reference adapter %s does not exec its harness", adapter)
		}
	}
	// claude and codex read the prompt from the stdin the loop pipes them: no positional
	// prompt survives on the final hop, so the prompt never reaches the harness CLI's argv.
	requireFileNotContains(t, filepath.Join(root, ".bench", "adapters", "claude"), `"$1"`, "claude adapter still passes the prompt on argv instead of reading stdin")
	requireFileNotContains(t, filepath.Join(root, ".bench", "adapters", "codex"), `"$1"`, "codex adapter still passes the prompt on argv instead of reading stdin")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "claude"), `exec claude -p --model "$model"`, "routed claude adapter does not preserve the model while reading the prompt from stdin")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "claude"), "exec claude -p\n", "unrouted claude adapter does not map stdin to claude -p")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "codex"), `exec codex exec --sandbox workspace-write -m "$model"`, "routed codex adapter does not select workspace-write while preserving the model")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "codex"), "exec codex exec --sandbox workspace-write\n", "unrouted codex adapter does not select workspace-write")
	// opencode's CLI documents only a positional prompt, so its adapter reads stdin into a
	// variable and passes it positionally — the documented residual final-hop exposure.
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `model="$("$_cmd" resolve-model --harness opencode)"`, "opencode adapter does not resolve its own harness column")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `opencode run --model "$model" -- "$prompt"`, "routed opencode adapter does not pass the stdin-read prompt positionally")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `opencode run -- "$prompt"`, "opencode adapter does not pass the stdin-read prompt positionally")
}

// The keys are pinned empty so discovery is hermetic: the API sources short-circuit to
// unavailable without a network call, and the contracts assert only argv posture.
var modelsHermeticEnv = map[string]string{"OPENAI_API_KEY": "", "ANTHROPIC_API_KEY": ""}

func TestRuntimeModelsContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench models unknown-arg contract", testRuntimeModelsUnknownArg)
	contract.RunParallel(t, "bench models discovery-tolerance contract", testRuntimeModelsDiscoveryTolerance)
}

// An unknown argument is a misuse: reject it with a usage line at exit 2, the sibling
// porcelain norm, instead of silently printing the inventory at exit 0.
func testRuntimeModelsUnknownArg(t *testing.T) {
	f := contract.NewFixture(t)
	probe := f.BenchEnv(modelsHermeticEnv, "models", "bogus")
	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage")
}

// The regression guard: the argv guard must not leak into the no-arg path — a plain
// invocation still emits the inventory tables at exit 0 even when every source is
// unreachable (discovery tolerance is a separate contract).
func testRuntimeModelsDiscoveryTolerance(t *testing.T) {
	f := contract.NewFixture(t)
	probe := f.BenchEnv(modelsHermeticEnv, "models")
	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "model_sources[")
	probe.RequireContains(probe.Stdout, "models[")
}
