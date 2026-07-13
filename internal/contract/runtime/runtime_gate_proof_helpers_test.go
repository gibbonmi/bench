package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

const localManifest = `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`

func manifestProof(id, contents, reason string, write int) ft78ProofCase {
	return ft78ProofCase{id: id, driver: func(t *testing.T) {
		f := proofFixture(t)
		if write == 1 {
			f.WriteFile(".bench/gate-inputs.json", contents)
		}
		f.Bench("gate").RequireExit(0)
		got := gate.Inspect(f.Root)
		if got.State != gate.Ready || got.Reason != reason || got.ReusableGreen {
			t.Fatalf("inspection = %s/%q reusable=%v, want ready/%q reusable=false", got.State, got.Reason, got.ReusableGreen, reason)
		}
		assertRuns(t, f, 1)
	}}
}

func proofFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\necho run >> .git/ft78-runs\nexit 0\n")
	f.WriteFile(".gitignore", ".bench/gate.sh\n.bench/gate-inputs.json\nft78-*\ninputs/\ntools/\npnpm-lock.yaml\npackage.json\npyproject.toml\nCargo.toml\n")
	f.CommitAll("base")
	return f
}

func escapedManifestPathProof(t *testing.T) {
	f := proofFixture(t)
	external := filepath.Join(t.TempDir(), "outside")
	contract.WriteFileAbs(t, external, "outside\n")
	if err := os.Symlink(external, filepath.Join(f.Root, "ft78-link")); err != nil {
		t.Fatal(err)
	}
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":["ft78-link"],"tools":[]}`)
	f.Bench("gate").RequireExit(0)
	got := gate.Inspect(f.Root)
	if got.State != gate.Ready || got.Reason != "declared path unavailable" || got.ReusableGreen {
		t.Fatalf("inspection = %s/%q reusable=%v, want ready/declared path unavailable/false", got.State, got.Reason, got.ReusableGreen)
	}
	assertRuns(t, f, 1)
}

func manifestByteLimitProof(t *testing.T, size int, state gate.State, reason string) {
	f := proofFixture(t)
	contents := localManifest + strings.Repeat(" ", size-len(localManifest))
	if len(contents) != size {
		t.Fatalf("manifest bytes = %d, want %d", len(contents), size)
	}
	f.WriteFile(".bench/gate-inputs.json", contents)
	f.Bench("gate").RequireExit(0)
	got := gate.Inspect(f.Root)
	if got.State != state || got.Reason != reason {
		t.Fatalf("inspection = %s/%q, want %s/%q", got.State, got.Reason, state, reason)
	}
	assertRuns(t, f, 1)
}

type mutationKind int

const (
	mutationGateScript mutationKind = iota
	mutationGateInterpreter
	mutationManifest
	mutationToolContent
	mutationToolMode
	mutationToolTarget
	mutationDeclaredFile
	mutationDeclaredDirectory
	mutationPATH
	mutationAutoKind
)

func mutationProof(id string, kind mutationKind) ft78ProofCase {
	return ft78ProofCase{id: id, driver: func(t *testing.T) { subjectMutationProof(t, kind) }}
}

func subjectMutationProof(t *testing.T, kind mutationKind) {
	f := proofFixture(t)
	pathA := filepath.Join(f.Root, "ft78-path-a")
	pathB := filepath.Join(f.Root, "ft78-path-b")
	for _, dir := range []string{pathA, pathB} {
		contract.Mkdir(t, dir)
		contract.WriteExecutableAbs(t, filepath.Join(dir, "ft78-tool"), "#!/bin/sh\nexit 0\n")
	}
	f.WriteFile("inputs/file", "one\n")
	f.WriteFile("inputs/dir/entry", "one\n")
	f.WriteExecutable("tools/target-a", "#!/bin/sh\nexit 0\n")
	f.WriteExecutable("tools/target-b", "#!/bin/sh\nexit 0\n")
	if err := os.Symlink("target-a", filepath.Join(f.Root, "tools", "selected")); err != nil {
		t.Fatal(err)
	}
	manifest := `{"schema":1,"closure":"local","environment":[],"paths":["inputs/file","inputs/dir"],"tools":["tools/selected"]}`
	f.WriteFile(".bench/gate-inputs.json", manifest)
	env := map[string]string{"PATH": pathA + ":/usr/bin:/bin"}
	if kind == mutationAutoKind {
		contract.Remove(t, filepath.Join(f.Root, ".bench", "gate.sh"))
		f.WriteFile("pnpm-lock.yaml", "lock\n")
		writeToolShim(t, pathA, "bash")
		writeToolShim(t, pathA, "pnpm")
		f.WriteFile(".bench/gate-inputs.json", localManifest)
	}
	first := f.BenchEnv(env, "gate")
	first.RequireExit(0)
	before := readVerdict(t, f)
	switch kind {
	case mutationGateScript:
		f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\n# changed\necho run >> .git/ft78-runs\nexit 0\n")
	case mutationGateInterpreter:
		interpreter := filepath.Join(f.Root, "ft78-interpreter")
		contract.WriteExecutableAbs(t, interpreter, "#!/bin/sh\nexec /bin/sh \"$@\"\n")
		f.WriteExecutable(".bench/gate.sh", "#!"+interpreter+"\necho run >> .git/ft78-runs\nexit 0\n")
		f.BenchEnv(env, "gate").RequireExit(0)
		before = readVerdict(t, f)
		contract.WriteExecutableAbs(t, interpreter, "#!/bin/sh\n# changed\nexec /bin/sh \"$@\"\n")
	case mutationManifest:
		f.WriteFile(".bench/gate-inputs.json", strings.Replace(manifest, `"closure":"local"`, `"closure":"remote"`, 1))
	case mutationToolContent:
		f.WriteExecutable("tools/target-a", "#!/bin/sh\n# changed\nexit 0\n")
	case mutationToolMode:
		if err := os.Chmod(filepath.Join(f.Root, "tools", "target-a"), 0o700); err != nil {
			t.Fatal(err)
		}
	case mutationToolTarget:
		if err := os.Remove(filepath.Join(f.Root, "tools", "selected")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("target-b", filepath.Join(f.Root, "tools", "selected")); err != nil {
			t.Fatal(err)
		}
	case mutationDeclaredFile:
		f.WriteFile("inputs/file", "two\n")
	case mutationDeclaredDirectory:
		f.WriteFile("inputs/dir/entry", "two\n")
	case mutationPATH:
		env["PATH"] = pathB + ":/usr/bin:/bin"
	case mutationAutoKind:
		contract.Remove(t, filepath.Join(f.Root, "pnpm-lock.yaml"))
		f.WriteFile("package.json", "{}\n")
		writeToolShim(t, pathA, "npm")
	}
	f.BenchEnv(env, "gate").RequireExit(0)
	after := readVerdict(t, f)
	if before.Tree != after.Tree || before.Oracle == after.Oracle {
		t.Fatalf("subject mutation = tree %q→%q oracle %q→%q, want stable tree and changed oracle", before.Tree, after.Tree, before.Oracle, after.Oracle)
	}
}

func repositoryRootMutationProof(t *testing.T) {
	a, b := proofFixture(t), proofFixture(t)
	a.WriteFile(".bench/gate-inputs.json", localManifest)
	b.WriteFile(".bench/gate-inputs.json", localManifest)
	a.Bench("gate").RequireExit(0)
	b.Bench("gate").RequireExit(0)
	av, bv := readVerdict(t, a), readVerdict(t, b)
	if av.Tree != bv.Tree || av.Oracle == bv.Oracle || a.Root == b.Root {
		t.Fatalf("root mutation = roots %q/%q tree %q/%q oracle %q/%q", a.Root, b.Root, av.Tree, bv.Tree, av.Oracle, bv.Oracle)
	}
}

func ignoredInputRerunProof(t *testing.T) {
	f := ignoredInputFixture(t)
	f.Bench("gate").RequireExit(0)
	f.WriteFile("ft78-ignored", "red\n")
	f.Bench("gate").RequireExit(19)
	assertRuns(t, f, 2)
}

func ignoredInputRefusalProof(t *testing.T) {
	f := ignoredInputFixture(t)
	f.Bench("gate").RequireExit(0)
	f.WriteFile("ft78-ignored", "red\n")
	f.WriteFile("work.txt", "changed\n")
	before := strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
	probe := f.Bench("commit", "-m", "must refuse", "work.txt")
	if probe.ExitCode == 0 || strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout) != before {
		t.Fatalf("declared ignored-input mutation did not refuse commit: exit=%d", probe.ExitCode)
	}
	assertRuns(t, f, 2)
}

func ignoredInputFixture(t *testing.T) contract.Fixture {
	f := proofFixture(t)
	f.WriteFile("work.txt", "base\n")
	f.WriteFile("ft78-ignored", "green\n")
	f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\necho run >> .git/ft78-runs\n[ \"$(cat ft78-ignored)\" = green ] || exit 19\n")
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":["ft78-ignored"],"tools":[]}`)
	f.Git("add", "work.txt")
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "work")
	return f
}

func passlistProof(id, name, value, expected string, declaredAbsent bool) ft78ProofCase {
	return ft78ProofCase{id: id, driver: func(t *testing.T) {
		f := proofFixture(t)
		f.WriteExecutable(".bench/gate.sh", fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' \"${%s-unset}\" > .git/ft78-env\n", name))
		manifest := localManifest
		if declaredAbsent {
			manifest = `{"schema":1,"closure":"local","environment":["FT78_DECLARED_ABSENT"],"paths":[],"tools":[]}`
		}
		f.WriteFile(".bench/gate-inputs.json", manifest)
		if name == "BENCH_KIT" {
			value = contract.SubjectRoot(t)
		}
		if name == "BENCH_WRAPPER" {
			value = benchPath(t)
		}
		env := contract.Env{name: &value}
		if declaredAbsent {
			env[name] = nil
		}
		f.BenchEnvSpec(env, "gate").RequireExit(0)
		if got := strings.TrimSpace(contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "ft78-env"))); got != expected {
			t.Fatalf("%s observed %q, want literal %q", name, got, expected)
		}
		if declaredAbsent {
			got := gate.Inspect(f.Root)
			if got.State != gate.Ready || got.Reason != "declared environment unavailable" || got.ReusableGreen {
				t.Fatalf("inspection = %s/%q", got.State, got.Reason)
			}
		}
	}}
}

type resolverKind int

const (
	resolverGateSh resolverKind = iota
	resolverBenchGate
	resolverPnpm
	resolverNpm
	resolverPython
	resolverCargo
)

func resolverPathProof(id, lockfile, marker string) ft78ProofCase {
	return ft78ProofCase{id: id, driver: func(t *testing.T) {
		f, env := resolverFixture(t, map[string]string{lockfile: "fixture\n"}, "")
		f.BenchEnv(env, "gate").RequireExit(0)
		assertToolMarker(t, f, marker)
		if got := gate.Inspect(f.Root); got.State != gate.Ready || got.Status != "green" {
			t.Fatalf("inspection = %+v, want ready green", got)
		}
	}}
}

func resolverMutationProof(id string, kind resolverKind, tool string) ft78ProofCase {
	return ft78ProofCase{id: id, driver: func(t *testing.T) {
		files, benchGate := resolverFiles(kind)
		f, env := resolverFixture(t, files, benchGate)
		if kind == resolverGateSh {
			interpreter := filepath.Join(f.Root, "ft78-bin", "gate-interpreter")
			contract.WriteExecutableAbs(t, interpreter, "#!/bin/sh\nexec /bin/sh \"$@\"\n")
			f.WriteExecutable(".bench/gate.sh", "#!"+interpreter+"\necho gate.sh >> .git/ft78-tools\n")
		}
		f.BenchEnv(env, "gate").RequireExit(0)
		before := readVerdict(t, f)
		path := filepath.Join(f.Root, "ft78-bin", tool)
		if tool == "gate.sh" {
			path = filepath.Join(f.Root, ".bench", "gate.sh")
		}
		contract.WriteExecutableAbs(t, path, contract.ReadFileAbs(t, path)+"# changed\n")
		f.BenchEnv(env, "gate").RequireExit(0)
		after := readVerdict(t, f)
		if before.Tree != after.Tree || before.Oracle == after.Oracle {
			t.Fatalf("%s mutation did not change only oracle", tool)
		}
	}}
}

func resolverFiles(kind resolverKind) (map[string]string, string) {
	switch kind {
	case resolverGateSh:
		return map[string]string{}, ""
	case resolverBenchGate:
		return map[string]string{}, "echo bench-gate >> .git/ft78-tools"
	case resolverPnpm:
		return map[string]string{"pnpm-lock.yaml": "lock\n"}, ""
	case resolverNpm:
		return map[string]string{"package.json": "{}\n"}, ""
	case resolverPython:
		return map[string]string{"pyproject.toml": "[tool]\n"}, ""
	default:
		return map[string]string{"Cargo.toml": "[package]\n"}, ""
	}
}

func resolverFixture(t *testing.T, files map[string]string, benchGate string) (contract.Fixture, map[string]string) {
	f := proofFixture(t)
	contract.Remove(t, filepath.Join(f.Root, ".bench", "gate.sh"))
	f.WriteFile(".bench/gate-inputs.json", localManifest)
	for name, contents := range files {
		f.WriteFile(name, contents)
	}
	bin := filepath.Join(f.Root, "ft78-bin")
	contract.Mkdir(t, bin)
	for _, tool := range []string{"bash", "pnpm", "npm", "mypy", "pytest", "ruff", "cargo", "rustc", "clippy-driver"} {
		writeToolShim(t, bin, tool)
	}
	return f, map[string]string{"PATH": bin + ":/usr/bin:/bin", "BENCH_GATE": benchGate}
}

func writeToolShim(t *testing.T, dir, tool string) {
	path := filepath.Join(dir, tool)
	if tool == "bash" {
		contract.WriteExecutableAbs(t, path, "#!/bin/sh\necho bash >> .git/ft78-tools\nexec /bin/bash \"$@\"\n")
		return
	}
	contract.WriteExecutableAbs(t, path, "#!/bin/sh\necho "+tool+" >> .git/ft78-tools\nexit 0\n")
}

func fullResolverPrecedenceProof(t *testing.T) {
	f, env := resolverFixture(t, map[string]string{"pnpm-lock.yaml": "x\n", "package.json": "{}\n", "pyproject.toml": "x\n", "Cargo.toml": "x\n"}, "echo bench-gate >> .git/ft78-tools; exit 1")
	f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\necho gate.sh >> .git/ft78-tools\nexit 0\n")
	f.BenchEnv(env, "gate").RequireExit(0)
	if got := contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "ft78-tools")); got != "gate.sh\n" {
		t.Fatalf("precedence marker = %q, want literal gate.sh", got)
	}
}

func noGateResolverProof(t *testing.T) {
	f, env := resolverFixture(t, map[string]string{}, "")
	contract.Remove(t, filepath.Join(gitDir(t, f), "bench-last-gate"))
	f.BenchEnv(env, "gate").RequireExit(3)
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); !os.IsNotExist(err) {
		t.Fatalf("no-gate record stat error = %v, want absent", err)
	}
}

func assertToolMarker(t *testing.T, f contract.Fixture, marker string) {
	t.Helper()
	lines := contract.NonEmptyLines(contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "ft78-tools")))
	for _, line := range lines {
		if line == marker {
			return
		}
	}
	t.Fatalf("tool markers = %#v, want literal %q", lines, marker)
}

type proofVerdict struct{ Tree, Oracle string }

func readVerdict(t *testing.T, f contract.Fixture) proofVerdict {
	t.Helper()
	var record proofVerdict
	if err := json.Unmarshal([]byte(contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"))), &record); err != nil {
		t.Fatal(err)
	}
	return record
}

func assertRuns(t *testing.T, f contract.Fixture, want int) {
	t.Helper()
	got := len(contract.NonEmptyLines(contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "ft78-runs"))))
	if got != want {
		t.Fatalf("gate runs = %d, want literal %d", got, want)
	}
}
