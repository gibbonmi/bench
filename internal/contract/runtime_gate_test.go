package contract

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRuntimeGateContracts(t *testing.T) {
	skipIfSubjectBenchMissing(t)
	t.Parallel()
	runParallel(t, "bench gate repo-root cwd contract", testRuntimeGateRepoRootCWD)
	runParallel(t, "bench gate BENCH_GATE cwd contract", testRuntimeGateBenchGateCWD)
	runParallel(t, "bench gate resolution-order contract", testRuntimeGateResolutionOrder)
	runParallel(t, "bench status gate-cache write contract", testRuntimeStopHookGateCacheWrite)
	runParallel(t, "stop hook no-gate no-cache contract", testRuntimeStopHookNoGateNoCache)
	runParallel(t, "stop hook missing-core-binary fail-safe contract", testRuntimeStopHookMissingCoreBinary)
	runParallel(t, "bench gate missing-core-binary fail-safe contract", testRuntimeGateMissingCoreBinary)
	runParallel(t, "bench gate verdict-record contract", testRuntimeGateVerdictRecord)
	runParallel(t, "bench symlinked kit-dir portability contract", testRuntimeSymlinkedKitDir)
	runParallel(t, "stop hook stop_hook_active contract", testRuntimeStopHookActive)
	runParallel(t, "stop hook missing-bench fail-open contract", testRuntimeStopHookMissingBenchFailOpen)
}

func testRuntimeGateRepoRootCWD(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("package.json", "{\"ok\":true}\n")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\n[ -f package.json ]\n")
	mkdir(t, filepath.Join(f.Root, "sub"))

	runAt(t, f, filepath.Join(f.Root, "sub"), nil, "bash", benchPath(t), "gate").RequireExit(0)
}

func testRuntimeGateBenchGateCWD(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("package.json", "{\"ok\":true}\n")
	f.WriteExecutable("gate-root.sh", "#!/usr/bin/env bash\n[ -f package.json ]\n")
	mkdir(t, filepath.Join(f.Root, "sub"))

	runAt(t, f, filepath.Join(f.Root, "sub"), map[string]string{"BENCH_GATE": "./gate-root.sh"}, "bash", benchPath(t), "gate").RequireExit(0)
}

func testRuntimeGateResolutionOrder(t *testing.T) {
	f := NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.BenchEnv(map[string]string{"BENCH_GATE": "exit 1"}, "gate").RequireExit(0)

	remove(t, filepath.Join(f.Root, ".bench", "gate.sh"))
	f.WriteFile("package.json", "{\"private\":true}\n")
	f.BenchEnv(map[string]string{"BENCH_GATE": "exit 0"}, "gate").RequireExit(0)

	auto := f.Bench("gate")
	auto.RequireNotContains(auto.Stdout+auto.Stderr, "no gate found")
	if auto.ExitCode == 3 {
		t.Fatalf("package.json resolved to no-gate exit 3\nstdout:\n%s\nstderr:\n%s", auto.Stdout, auto.Stderr)
	}

	remove(t, filepath.Join(f.Root, "package.json"))
	commitAllowEmpty(t, f, "init")
	remove(t, filepath.Join(gitDir(t, f), "bench-last-gate"))
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
	data := readFileAbs(t, cache)
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
	f := NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	writeFileAbs(t, cache, "green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n")
	f.Bench("gate").RequireExit(0)
	requireContains(t, readFileAbs(t, cache), "green "+strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout))
	requireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	commitAllowEmpty(t, f, "same-tree")
	requireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
	if p := f.Bench("gate"); p.ExitCode == 0 {
		t.Fatal("red gate run exited zero")
	}
	requireContains(t, readFileAbs(t, cache), "red "+strings.TrimSpace(f.Bench("tree-hash").Stdout))
}

func testRuntimeSymlinkedKitDir(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	binDir := filepath.Join(tmp, "bin")
	shim := filepath.Join(tmp, "shim")
	mkdir(t, repo)
	mkdir(t, binDir)
	mkdir(t, shim)
	if err := os.Symlink(benchPath(t), filepath.Join(binDir, "bench")); err != nil {
		t.Fatalf("symlink bench: %v", err)
	}
	writeExecutableAbs(t, filepath.Join(shim, "readlink"), "#!/usr/bin/env bash\nif [ \"${1:-}\" = \"-f\" ]; then exit 1; fi\n/usr/bin/readlink \"$@\"\n")
	f := Fixture{t: t, Root: repo, Env: isolatedEnv(t, repo)}
	f.Git("init", "-q")
	runAt(t, f, repo, map[string]string{"PATH": shim + ":/usr/bin:/bin"}, filepath.Join(binDir, "bench"), "link").RequireExit(0)
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
	f := NewFixture(t)
	commitAllowEmpty(t, f, "init")
	probe := runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1", "PATH": "/usr/bin:/bin"}, "{}\n")
	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout+probe.Stderr, "bench")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("missing bench forged a gate cache")
	}
}

func copiedCLIHookFixture(t *testing.T, withCore bool) Fixture {
	t.Helper()
	f := NewFixture(t)
	mkdir(t, filepath.Join(f.Root, ".bench"))
	mkdir(t, filepath.Join(f.Root, "bin"))
	if withCore {
		mkdir(t, filepath.Join(f.Root, "dist"))
		copyRuntimeFile(t, filepath.Join(SubjectRoot(t), "dist", "bench"), filepath.Join(f.Root, "dist", "bench"), 0o755)
	}
	matches, err := filepath.Glob(filepath.Join(SubjectRoot(t), "bin", "*.sh"))
	if err != nil {
		t.Fatalf("glob bin scripts: %v", err)
	}
	for _, src := range matches {
		copyRuntimeFile(t, src, filepath.Join(f.Root, "bin", filepath.Base(src)), 0o755)
	}
	return f
}

func runStopHook(t *testing.T, f Fixture, env map[string]string, stdin string) Probe {
	t.Helper()
	return runFixtureCommand(t, f, f.Root, env, stdin, "bash", filepath.Join(SubjectRoot(t), ".bench", "hooks", "stop.sh"))
}
