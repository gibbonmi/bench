// Package preprelease holds the `bench prep-release` surface contracts: the ship-tier
// rehearsal driven end to end in throwaway repos, never against the real kit, whose one
// real run costs a four-platform artifact matrix plus the ~372 s release-evidence probe.
//
// Two kinds of fixture appear here, and the difference is which half is under test.
// A shipRepo pairs a throwaway graded tree with a throwaway kit, so what the command
// orchestrates — the argv it hands each script, the evidence that lands, the exit code —
// is observed cheaply and in full. The tier rows drive the real kit's conformance entry
// point instead, because tier membership is a fact about that driver and nothing else
// can answer for it.
package preprelease

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// Recording paths every fixture script writes through, all under the gitignored dist/
// tree: a fixture that recorded into a tracked path would dirty the subject and falsify
// the very verdict the command read before it started.
const (
	artifactsArgvFile = "dist/prep-artifacts-argv"
	preflightArgvFile = "dist/prep-preflight-argv"
	conformanceFile   = "dist/prep-conformance"
	evidenceIndexPath = "dist/preflight/release-index.json"
	evidenceArtifacts = "dist/artifacts"

	// stubPreflightSleepEnv holds the throwaway preflight script open between staging
	// and promotion, which is the window the interrupt row signals into.
	stubPreflightSleepEnv = "BENCH_STUB_PREFLIGHT_SLEEP"
)

// shipRepo is a throwaway repo `bench prep-release` runs to completion in, paired with
// the throwaway kit its conformance step compiles against.
type shipRepo struct {
	contract.Fixture
	t   *testing.T
	Kit string
}

// runEnv is the environment every surface drives this fixture through: the kit the
// conformance step compiles in, plus whatever a row seeds on top.
func (r shipRepo) runEnv(extra map[string]string) map[string]string {
	env := map[string]string{
		"BENCH_KIT":   r.Kit,
		runbinary.Env: contract.SelectedBench(r.t).Path,
	}
	for k, v := range extra {
		env[k] = v
	}
	return env
}

func (r shipRepo) benchScript() string {
	r.t.Helper()
	return filepath.Join(contract.SubjectRoot(r.t), "bin", "bench.sh")
}

// prepRelease runs the command through the kit CLI, the surface a maintainer types.
func (r shipRepo) prepRelease(extra map[string]string) contract.Probe {
	r.t.Helper()
	return r.RunEnv(r.runEnv(extra), "bash", r.benchScript(), "prep-release")
}

// newShipRepo materializes the throwaway tree: the three scripts prep-release shells
// out to, the plan they read, a dev gate that greens, and a gitignored dist/ so nothing
// a step writes can change the tree hash the verdict is bound to.
func newShipRepo(t *testing.T) shipRepo { return newShipRepoAt(t, t.TempDir()) }

func newShipRepoAt(t *testing.T, root string) shipRepo {
	t.Helper()
	f := contract.NewExecFixtureAt(t, root)
	f.Git("init", "-q")

	f.WriteFile(".gitignore", "dist/\n")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	// The declared-input manifest is what closes the gate subject; without one the
	// verdict is never reusable and the precondition could never be satisfied by any
	// implementation, so every fixture carries the smallest honest declaration.
	f.WriteFile(".bench/gate-inputs.json", gateInputsJSON)
	// A canary harness holding one dev-tier fixture, so the ship sweep has a harness to
	// read and selects nothing from it. A ship-tier fixture here would cost a full inner
	// gate, and whether the two real probe-derived ones bite is graded on the real tree.
	f.WriteFile("tests/canary/dev-only/EXPECT", "a dev-tier fixture the ship sweep passes over\n")
	f.WriteFile("scripts/release-plan.json", planJSON)
	f.WriteExecutable("scripts/build-artifacts.sh", buildArtifactsScript)
	f.WriteExecutable("scripts/release-preflight.sh", releasePreflightScript)
	f.CommitAll("ship fixture")

	repo := shipRepo{Fixture: f, t: t, Kit: newStubKit(t)}
	repo.greenTheGate()
	return repo
}

// hostileRoot is the repo-root shape the path row names: a directory whose name carries
// a space and a glob character, so an unquoted expansion anywhere in the four scripts
// the command hands the root to fails here rather than in a maintainer's checkout.
func hostileRoot(t *testing.T, parent string) string {
	t.Helper()
	root := filepath.Join(parent, "ship repo [v1]")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("create hostile-path root: %v", err)
	}
	return root
}

// greenTheGate records the dev-green verdict the command's precondition reads, through
// the real `bench gate` rather than a hand-written cache record: the precondition is a
// claim about what the gate wrote, so seeding it any other way would assert against this
// test's idea of the format instead of the gate's.
func (r shipRepo) greenTheGate() {
	r.t.Helper()
	r.RunEnv(r.runEnv(nil), "bash", r.benchScript(), "gate").RequireExit(0)
}

// newStubKit is a throwaway checkout standing where BENCH_KIT points: a module whose
// conformance entry point records the tier it was asked for and whether the stress tag
// reached the compiler. The selected binary is inherited separately, so this kit stays
// source-only. The real kit cannot stand in for the module —
// its ship tier is the ~372 s probe this whole split exists to keep off the common path.
func newStubKit(t *testing.T) string {
	t.Helper()
	kit := filepath.Join(t.TempDir(), "stub-kit")
	pkg := filepath.Join(kit, "internal", "conformance")
	contract.WriteFileAbs(t, filepath.Join(kit, "go.mod"), "module benchstubkit\n\ngo 1.25\n")
	contract.WriteFileAbs(t, filepath.Join(pkg, "entry_test.go"), stubEntrySource)
	contract.WriteFileAbs(t, filepath.Join(pkg, "stress_test.go"), stubStressSource)
	contract.WriteFileAbs(t, filepath.Join(pkg, "nostress_test.go"), stubNoStressSource)
	contract.WriteFileAbs(t, filepath.Join(kit, "cmd", "bench", "main.go"), stubGateGoSource)

	return kit
}

// planTargetCount is the artifact count the evidence assertion is graded against, read
// from the plan the build script reads. Deriving it here rather than writing a number
// is what makes a stub that exits 0 without building anything fail.
func planTargetCount(t *testing.T, root string) int {
	t.Helper()
	var plan struct {
		Targets []struct {
			Goos   string `json:"goos"`
			Goarch string `json:"goarch"`
		} `json:"targets"`
	}
	data, err := os.ReadFile(filepath.Join(root, "scripts", "release-plan.json"))
	if err != nil {
		t.Fatalf("read release plan: %v", err)
	}
	if err := json.Unmarshal(data, &plan); err != nil {
		t.Fatalf("decode release plan: %v", err)
	}
	return len(plan.Targets)
}

func countFiles(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	return len(entries)
}

// pathWithout is a PATH holding every required tool but one, built by symlinking the
// resolved tools into a fresh directory. Emptying PATH outright would prove nothing:
// the command has to name the tool it cannot find, not fail on the first thing missing.
func pathWithout(t *testing.T, omit string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "path")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create path dir: %v", err)
	}
	for _, tool := range []string{"bash", "git", "go", "node"} {
		if tool == omit {
			continue
		}
		resolved, err := exec.LookPath(tool)
		if err != nil {
			t.Fatalf("resolve %s: %v", tool, err)
		}
		if err := os.Symlink(resolved, filepath.Join(dir, tool)); err != nil {
			t.Fatalf("link %s: %v", tool, err)
		}
	}
	return dir
}

func requireContains(t *testing.T, what, body, want string) {
	t.Helper()
	if !strings.Contains(body, want) {
		t.Fatalf("%s = %q, want it to contain %q", what, body, want)
	}
}

const gateInputsJSON = `{
  "schema": 1,
  "closure": "local",
  "environment": ["HOME"],
  "paths": ["scripts/release-plan.json"],
  "tools": ["bash", "git"]
}
`

const planJSON = `{
  "schema_version": 1,
  "targets": [
    {"goos":"darwin","goarch":"arm64"},
    {"goos":"darwin","goarch":"amd64"},
    {"goos":"linux","goarch":"arm64"},
    {"goos":"linux","goarch":"amd64"}
  ]
}
`

// buildArtifactsScript stands in for scripts/build-artifacts.sh at its published
// interface — <source-root> <output-dir> — and promotes the finished set by rename, so
// a re-run replaces the directory rather than accumulating into it.
const buildArtifactsScript = `#!/usr/bin/env bash
set -euo pipefail
source_root="${1:?usage: build-artifacts.sh <source-root> <output-dir>}"
output="${2:?usage: build-artifacts.sh <source-root> <output-dir>}"
mkdir -p "$source_root/dist"
printf '%s\n%s\n' "$source_root" "$output" > "$source_root/dist/prep-artifacts-argv"
parent="$(dirname "$output")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.artifacts.XXXXXX")"
node -e '
const fs = require("fs");
const plan = JSON.parse(fs.readFileSync(process.argv[1], "utf8"));
for (const target of plan.targets) {
  fs.writeFileSync(process.argv[2] + "/bench-" + target.goos + "-" + target.goarch + ".tgz", "artifact\n");
}
' "$source_root/scripts/release-plan.json" "$stage"
rm -rf "$output"
mv "$stage" "$output"
`

// releasePreflightScript stands in for scripts/release-preflight.sh, resolving its own
// root the way the real script does and promoting the evidence directory by rename. The
// optional sleep is the interrupt row's window: it holds the run open between staging
// and promotion, which is exactly where a direct write would leave a partial index.
const releasePreflightScript = `#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$root/dist"
printf '%s\n' "$@" > "$root/dist/prep-preflight-argv"
stage="$(mktemp -d "$root/dist/.preflight.XXXXXX")"
printf '{"mode":"verify","status":"green"}\n' > "$stage/release-index.json"
sleep "${BENCH_STUB_PREFLIGHT_SLEEP:-0}"
rm -rf "$root/dist/preflight"
mv "$stage" "$root/dist/preflight"
`

// stubEntrySource is the throwaway kit's conformance entry point. It records the tier it
// was asked for and whether the stress tag reached the compiler, which is what makes a
// forgotten -tags stress observable without paying for the real matrix.
const stubEntrySource = `package conformance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootConformance(t *testing.T) {
	root := os.Getenv("BENCH_CONFORMANCE_ROOT")
	if root == "" {
		t.Fatal("BENCH_CONFORMANCE_ROOT not set")
	}
	record := "tier=" + os.Getenv("BENCH_CONFORMANCE_TIER") + " stress=" + stressTag + "\n"
	if err := os.WriteFile(filepath.Join(root, "dist", "prep-conformance"), []byte(record), 0o644); err != nil {
		t.Fatal(err)
	}
}
`

const stubStressSource = "//go:build stress\n\npackage conformance\n\nconst stressTag = \"on\"\n"

const stubNoStressSource = "//go:build !stress\n\npackage conformance\n\nconst stressTag = \"off\"\n"

// stubGateGoSource stands in for the `go run ./cmd/bench gate-go test` the ship-tier core
// test step drives. These rows grade orchestration — that the sequence runs and leaves its
// evidence — and the real step would enumerate and test the throwaway module's packages,
// which is cost with nothing to observe. What the step actually asks for is pinned at the
// Steps seam by TestStepsRunReleaseOnlyPackages, not here.
const stubGateGoSource = "package main\n\nfunc main() {}\n"
