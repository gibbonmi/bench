package runtime

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRuntimeGateContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench gate repo-root cwd contract", testRuntimeGateRepoRootCWD)
	contract.RunParallel(t, "bench gate BENCH_GATE cwd contract", testRuntimeGateBenchGateCWD)
	contract.RunParallel(t, "bench gate resolution-order contract", testRuntimeGateResolutionOrder)
	contract.RunParallel(t, "bench status gate-cache write contract", testRuntimeStopHookGateCacheWrite)
	contract.RunParallel(t, "stop hook no-gate no-cache contract", testRuntimeStopHookNoGateNoCache)
	contract.RunParallel(t, "stop hook missing-core-binary fail-safe contract", testRuntimeStopHookMissingCoreBinary)
	contract.RunParallel(t, "bench gate missing-core-binary fail-safe contract", testRuntimeGateMissingCoreBinary)
	contract.RunParallel(t, "bench gate verdict-record contract", testRuntimeGateVerdictRecord)
	contract.RunParallel(t, "bench gate pin non-TTY refusal contract", testRuntimeGatePinNonTTYRefusal)
	contract.RunParallel(t, "bench symlinked kit-dir portability contract", testRuntimeSymlinkedKitDir)
	contract.RunParallel(t, "stop hook stop_hook_active contract", testRuntimeStopHookActive)
	contract.RunParallel(t, "stop hook missing-bench fail-open contract", testRuntimeStopHookMissingBenchFailOpen)
	contract.RunParallel(t, "stop hook intent refresh contract", testRuntimeStopHookIntentRefresh)
}

func testRuntimeStopHookIntentRefresh(t *testing.T) {
	contract.NoteContractFailure(t, "stop hook intent refresh contract failed")
	f := copiedCLIHookFixture(t, true)
	f.CommitAll("init")
	ledger := filepath.Join(gitDir(t, f), "bench-intent.json")
	contract.WriteFileAbs(t, ledger, `{"schema":1,"entries":[{"key":"gone","kind":"shift","objective":"done","created_at":"2026-07-11T00:00:00Z","worktree":"/definitely/missing"}]}`+"\n")
	runStopHook(t, f, nil, "{}\n").RequireExit(0)
	first := contract.ReadFileAbs(t, ledger)
	if !strings.Contains(first, `"entries":[]`) {
		t.Fatalf("unarmed Stop did not compact proven-done intent: %s", first)
	}
	runStopHook(t, f, nil, `{"stop_hook_active":true}`+"\n").RequireExit(0)
	if second := contract.ReadFileAbs(t, ledger); second != first {
		t.Fatalf("active Stop refresh changed bytes\nfirst=%s\nsecond=%s", first, second)
	}
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("intent-only Stop refresh forged gate cache")
	}
}

func testRuntimeGateRepoRootCWD(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("package.json", "{\"ok\":true}\n")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\n[ -f package.json ]\n")
	contract.Mkdir(t, filepath.Join(f.Root, "sub"))

	contract.RunAt(t, f, filepath.Join(f.Root, "sub"), nil, "bash", benchPath(t), "gate").RequireExit(0)
}

func testRuntimeGateBenchGateCWD(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("package.json", "{\"ok\":true}\n")
	f.WriteExecutable("gate-root.sh", "#!/usr/bin/env bash\n[ -f package.json ]\n")
	contract.Mkdir(t, filepath.Join(f.Root, "sub"))

	contract.RunAt(t, f, filepath.Join(f.Root, "sub"), map[string]string{"BENCH_GATE": "./gate-root.sh"}, "bash", benchPath(t), "gate").RequireExit(0)
}

func testRuntimeGateResolutionOrder(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.BenchEnv(map[string]string{"BENCH_GATE": "exit 1"}, "gate").RequireExit(0)

	contract.Remove(t, filepath.Join(f.Root, ".bench", "gate.sh"))
	f.WriteFile("package.json", "{\"private\":true}\n")
	f.BenchEnv(map[string]string{"BENCH_GATE": "exit 0"}, "gate").RequireExit(0)

	auto := f.Bench("gate")
	auto.RequireNotContains(auto.Stdout+auto.Stderr, "no gate found")
	if auto.ExitCode == 3 {
		t.Fatalf("package.json resolved to no-gate exit 3\nstdout:\n%s\nstderr:\n%s", auto.Stdout, auto.Stderr)
	}

	contract.Remove(t, filepath.Join(f.Root, "package.json"))
	commitAllowEmpty(t, f, "init")
	contract.Remove(t, filepath.Join(gitDir(t, f), "bench-last-gate"))
	noGate := f.Bench("gate")
	noGate.RequireExit(3)
	noGate.RequireContains(noGate.Stdout+noGate.Stderr, "no gate found")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("no-gate case recorded a verdict")
	}
}

func testRuntimeStopHookGateCacheWrite(t *testing.T) {
	f := copiedCLIHookFixture(t, true)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")

	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n")
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	data := contract.ReadFileAbs(t, cache)
	if !regexp.MustCompile(`^(green|red) [0-9a-f]+ [0-9T:Z-]+$`).MatchString(strings.TrimSpace(data)) {
		t.Fatalf("gate cache not <status> <tree> <iso8601>: %q", data)
	}
}

func testRuntimeStopHookNoGateNoCache(t *testing.T) {
	f := copiedCLIHookFixture(t, true)
	f.CommitAll("init")

	probe := runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n")
	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout+probe.Stderr, "no gate found")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("armed no-gate stop recorded a gate cache")
	}
}

func testRuntimeStopHookMissingCoreBinary(t *testing.T) {
	f := copiedCLIHookFixture(t, false)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")

	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("missing core binary forged a gate verdict")
	}
}

func testRuntimeGateMissingCoreBinary(t *testing.T) {
	f := copiedCLIHookFixture(t, false)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")

	f.Run("bash", filepath.Join(f.Root, "bin", "bench.sh"), "gate")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("gate_record forged a verdict with no core binary")
	}
}

func testRuntimeGateVerdictRecord(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	contract.WriteFileAbs(t, cache, "green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n")
	f.Bench("gate").RequireExit(0)
	contract.RequireContains(t, contract.ReadFileAbs(t, cache), "green "+strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout))
	contract.RequireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	commitAllowEmpty(t, f, "same-tree")
	contract.RequireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
	if p := f.Bench("gate"); p.ExitCode == 0 {
		t.Fatal("red gate run exited zero")
	}
	contract.RequireContains(t, contract.ReadFileAbs(t, cache), "red "+strings.TrimSpace(f.Bench("tree-hash").Stdout))
}

func testRuntimeGatePinNonTTYRefusal(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")
	pin := filepath.Join(gitDir(t, f), "bench-gate-pin")

	probe := contract.RunAtWithInput(t, f, f.Root, nil, "pin .bench\n", "bash", benchPath(t), "gate", "pin")
	if probe.ExitCode == 0 {
		t.Fatal("bench gate pin accepted non-TTY stdin")
	}
	probe.RequireContains(probe.Stderr, "interactive TTY")
	if _, err := os.Stat(pin); err == nil {
		t.Fatal("bench gate pin wrote a pin file after non-TTY refusal")
	}
}

func testRuntimeSymlinkedKitDir(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	binDir := filepath.Join(tmp, "bin")
	shim := filepath.Join(tmp, "shim")
	contract.Mkdir(t, repo)
	contract.Mkdir(t, binDir)
	contract.Mkdir(t, shim)
	if err := os.Symlink(benchPath(t), filepath.Join(binDir, "bench")); err != nil {
		t.Fatalf("symlink bench: %v", err)
	}
	contract.WriteExecutableAbs(t, filepath.Join(shim, "readlink"), "#!/usr/bin/env bash\nif [ \"${1:-}\" = \"-f\" ]; then exit 1; fi\n/usr/bin/readlink \"$@\"\n")
	f := contract.NewFixtureAt(t, repo, contract.IsolatedEnv(t, repo))
	f.Git("init", "-q")
	contract.RunAt(t, f, repo, map[string]string{"PATH": shim + ":/usr/bin:/bin"}, filepath.Join(binDir, "bench"), "link").RequireExit(0)
	if _, err := os.Stat(filepath.Join(repo, ".bench", "BENCH.md")); err != nil {
		t.Fatal("symlinked bench did not resolve kit dir without readlink -f")
	}
}

func testRuntimeStopHookActive(t *testing.T) {
	f := copiedCLIHookFixture(t, true)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
	f.CommitAll("init")
	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{\"stop_hook_active\":true}\n").RequireExit(0)
	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n").RequireExit(2)
}

func testRuntimeStopHookMissingBenchFailOpen(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	probe := runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1", "PATH": "/usr/bin:/bin"}, "{}\n")
	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout+probe.Stderr, "bench")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("missing bench forged a gate cache")
	}
}

func copiedCLIHookFixture(t *testing.T, withCore bool) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	contract.Mkdir(t, filepath.Join(f.Root, ".bench"))
	contract.Mkdir(t, filepath.Join(f.Root, "bin"))
	if withCore {
		contract.Mkdir(t, filepath.Join(f.Root, "dist"))
		copyRuntimeFile(t, filepath.Join(contract.SubjectRoot(t), "dist", "bench"), filepath.Join(f.Root, "dist", "bench"), 0o755)
	}
	matches, err := filepath.Glob(filepath.Join(contract.SubjectRoot(t), "bin", "*.sh"))
	if err != nil {
		t.Fatalf("glob bin scripts: %v", err)
	}
	for _, src := range matches {
		copyRuntimeFile(t, src, filepath.Join(f.Root, "bin", filepath.Base(src)), 0o755)
	}
	return f
}

func runStopHook(t *testing.T, f contract.Fixture, env map[string]string, stdin string) contract.Probe {
	t.Helper()
	return contract.RunAtWithInput(t, f, f.Root, env, stdin, "bash", filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "stop.sh"))
}
