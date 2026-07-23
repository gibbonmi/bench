package runtime

import (
	"encoding/json"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"os/exec"
	"path/filepath"
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
	contract.RunParallel(t, "bench gate rebuilt self-host contract", testRuntimeGateRebuiltSelfHost)
	contract.RunParallel(t, "oracle-bound gate verdict contract", testRuntimeOracleBoundVerdict)
	contract.RunParallel(t, "fail-closed gate verdict persistence contract", testRuntimePendingBeforeGate)
	contract.RunParallel(t, "bench gate pin non-TTY refusal contract", testRuntimeGatePinNonTTYRefusal)
	contract.RunParallel(t, "bench symlinked kit-dir portability contract", testRuntimeSymlinkedKitDir)
	contract.RunParallel(t, "stop hook stop_hook_active contract", testRuntimeStopHookActive)
	contract.RunParallel(t, "stop hook missing-bench fail-open contract", testRuntimeStopHookMissingBenchFailOpen)
	contract.RunParallel(t, "stop hook intent refresh contract", testRuntimeStopHookIntentRefresh)
	contract.RunParallel(t, "bench gate owner and signal cleanup contract", testRuntimeGateOwnerAndSignalCleanup)
}

func testRuntimeGateRebuiltSelfHost(t *testing.T) {
	contract.NoteContractFailure(t, "bench gate rebuilt self-host contract failed")
	root := contract.SubjectRoot(t)
	f := contract.NewFixture(t)
	wrapper := filepath.Join(f.Root, "bin", "bench.sh")
	copyRuntimeFile(t, filepath.Join(root, "bin", "bench.sh"), wrapper, 0o755)
	copyRuntimeFile(t, filepath.Join(root, ".bench", "gate-inputs.json"), filepath.Join(f.Root, ".bench", "gate-inputs.json"), 0o644)
	copyRuntimeFile(t, filepath.Join(root, "package.json"), filepath.Join(f.Root, "package.json"), 0o644)
	f.WriteFile(".gitignore", ".bench-contract-env/\n")
	f.WriteExecutable(".bench/gate.sh", `#!/usr/bin/env bash
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
test -n "$(go env GOPATH)"
exec "$root/bin/bench.sh" version
`)
	driver := `wrapper="$1"; kit="$2"; set --; source "$wrapper" >/dev/null; printf '%s\n%s\n' "$(package_version "$kit")" "$(platform_suffix)"`
	query := contract.RunAt(t, f, f.Root, nil, "bash", "-c", driver, "bench-wrapper-query", wrapper, f.Root)
	query.RequireExit(0)
	parts := contract.NonEmptyLines(query.Stdout)
	if len(parts) != 2 {
		t.Fatalf("wrapper cache query returned %d lines, want version and platform suffix:\n%s", len(parts), query.Stdout)
	}
	benchHome, home := f.Env["BENCH_HOME"], f.Env["HOME"]
	cacheBinary := filepath.Join(benchHome, "cache", "bin", parts[0], parts[1], "bench")
	contract.Mkdir(t, filepath.Dir(cacheBinary))
	build := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), root, cacheBinary)
	build.Dir = root
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build current bench binary: %v\n%s", err, output)
	}
	f.CommitAll("self-host fixture")

	f.RunEnv(map[string]string{"BENCH_HOME": benchHome, "HOME": home}, "bash", wrapper, "gate").RequireExit(0)
	if _, err := os.Stat(filepath.Join(home, ".bench")); !os.IsNotExist(err) {
		t.Fatalf("default HOME cache exists after cache-only gate run: %v", err)
	}
}

func testRuntimeOracleBoundVerdict(t *testing.T) {
	contract.NoteContractFailure(t, "oracle-bound gate verdict contract failed")
	f := contract.NewFixture(t)
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.WriteFile("work.txt", "work\n")
	f.CommitAll("base")
	f.WriteFile("work.txt", "changed\n")
	f.BenchEnv(map[string]string{"BENCH_GATE": "echo run >> .git/gate-runs; exit 0"}, "gate").RequireExit(0)
	before := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	probe := f.BenchEnv(map[string]string{"BENCH_GATE": "echo run >> .git/gate-runs; exit 23"}, "commit", "-m", "must refuse", "work.txt")
	if probe.ExitCode == 0 {
		t.Fatal("commit trusted green from a different exact gate command")
	}
	if got := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout); got != before {
		t.Fatal("oracle-mismatched commit moved HEAD")
	}
	if runs := len(contract.NonEmptyLines(contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "gate-runs")))); runs != 2 {
		t.Fatalf("oracle mutation gate runs = %d, want 2", runs)
	}
}

func testRuntimePendingBeforeGate(t *testing.T) {
	contract.NoteContractFailure(t, "fail-closed gate verdict persistence contract failed")
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", `#!/usr/bin/env bash
record=""
[ ! -f .git/bench-last-gate ] || record="$(<.git/bench-last-gate)"
case "$record" in
  *'"state":"pending"'*) printf 'pending\n' > .git/saw-pending ;;
  *) printf 'old-or-absent\n' > .git/saw-pending ;;
esac
exit 0
`)
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.CommitAll("base")
	f.Bench("gate").RequireExit(0)
	if got := contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "saw-pending")); got != "pending\n" {
		t.Fatalf("gate started before durable pending replacement: %q", got)
	}
}

func testRuntimeStopHookIntentRefresh(t *testing.T) {
	contract.NoteContractFailure(t, "stop hook intent refresh contract failed")
	f := copiedCLIHookFixture(t, true)
	f.CommitAll("init")
	ledger := filepath.Join(gitDir(t, f), "bench-intent.json")
	contract.WriteFileAbs(t, ledger, `{"schema":1,"entries":[{"key":"gone","kind":"shift","created_at":"2026-07-11T00:00:00Z","worktree":"/definitely/missing"}]}`+"\n")
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
	var record map[string]any
	if err := json.Unmarshal([]byte(data), &record); err != nil || record["schema"] != float64(1) || record["state"] != "ready" || record["status"] != "green" {
		t.Fatalf("gate cache is not a schema-1 ready green: %q (%v)", data, err)
	}
	if info, err := os.Stat(cache); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("gate cache mode = %v, %v; want 0600", info, err)
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
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.CommitAll("init")
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	contract.WriteFileAbs(t, cache, "green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n")
	f.Bench("gate").RequireExit(0)
	var record map[string]any
	if err := json.Unmarshal([]byte(contract.ReadFileAbs(t, cache)), &record); err != nil {
		t.Fatalf("decode ready verdict: %v", err)
	}
	if record["state"] != "ready" || record["status"] != "green" || record["tree"] != strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout) {
		t.Fatalf("unexpected green verdict: %#v", record)
	}
	contract.RequireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	commitAllowEmpty(t, f, "same-tree")
	contract.RequireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
	if p := f.Bench("gate"); p.ExitCode == 0 {
		t.Fatal("red gate run exited zero")
	}
	if err := json.Unmarshal([]byte(contract.ReadFileAbs(t, cache)), &record); err != nil {
		t.Fatalf("decode red verdict: %v", err)
	}
	if record["status"] != "red" || record["tree"] != strings.TrimSpace(f.Bench("tree-hash").Stdout) {
		t.Fatalf("unexpected red verdict: %#v", record)
	}
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
