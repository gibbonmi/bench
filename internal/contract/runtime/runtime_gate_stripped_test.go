package runtime

// The stripped-subject construction through the built binary: a full `bench gate-phases`
// run must grade its excludable phases against a materialized worktree with the declared
// allowlist absent, while the included phases keep seeing the real root. Each contract
// plants exactly the dependency the declaration forbids — a hard read, and the
// soft-skip degradation that is this tree's idiom for a missing subject file — and
// requires the gate to notice.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

// allowlistedProbePath is the declared path every fixture plants its dependency on. A
// literal rather than a read of the declaration: the assertion below keeps the fixture
// honest if the allowlist ever drops it, without the test deriving its input from the
// code under test.
const allowlistedProbePath = "ROADMAP.md"

// [R08] An excludable phase that hard-reads an allowlisted path must red a full gate:
// the stripped construction is what removes the file from the tree that phase grades.
func TestExcludablePhaseCannotReadAllowlistedPath(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "stripped-subject hard-read contract failed")
	probe := runStrippedGate(t, map[string]string{
		"conformance":    "test -f " + allowlistedProbePath,
		canary.PhaseTest: "cat " + allowlistedProbePath + " >/dev/null",
	})
	if probe.ExitCode == 0 {
		t.Fatalf("full gate stayed green while an excludable phase read %s:\n%s\n%s",
			allowlistedProbePath, probe.Stdout, probe.Stderr)
	}
	requireOutput(t, probe, "phase "+canary.PhaseTest+": red")
	// The included phase saw the real root on the same run, so the red belongs to the
	// planted read alone.
	requireOutput(t, probe, "phase conformance: green")
}

// [R09] An excludable phase that soft-skips on the missing path must red a full gate.
// The planted phase mimics skipIfSubjectFileMissing exactly — a structured
// kind=environment line, then exit 0 — because that degradation is the kit's dominant
// idiom for an absent subject file, and stripping alone converts it into a permanent
// silent green.
func TestExcludablePhaseCannotCapabilitySkip(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "stripped-subject soft-skip contract failed")
	probe := runStrippedGate(t, map[string]string{
		"conformance": "test -f " + allowlistedProbePath,
		canary.PhaseTest: softSkipScript(t, capability.Skip{
			Kind:   capability.KindEnvironment,
			Reason: "subject root has no " + allowlistedProbePath,
		}),
	})
	if probe.ExitCode == 0 {
		t.Fatalf("full gate stayed green while an excludable phase skipped on the stripped %s:\n%s\n%s",
			allowlistedProbePath, probe.Stdout, probe.Stderr)
	}
	requireOutput(t, probe, "fatal against the stripped subject")
}

// [R09 boundary] Only the environment kind is fatal against the stripped subject. A
// kind=capability skip is a host limitation — a filesystem without FIFOs, say — that
// stripping cannot induce and that real hosts emit, so it must leave a full gate green
// exactly as the dev tier does; the environment kind on the same fixture must red. The
// two cases differ only in the rendered skip line, so the kind boundary itself is the
// thing under test.
func TestHostCapabilitySkipDoesNotRedStrippedRun(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "stripped-subject skip-kind boundary contract failed")
	for _, tc := range []struct {
		name    string
		skip    capability.Skip
		wantRed bool
	}{
		{"capability kind stays informational", capability.Skip{
			Kind: capability.KindCapability, Class: capability.Fifo,
			Reason: "FIFOs unavailable on this filesystem",
		}, false},
		{"environment kind is fatal", capability.Skip{
			Kind:   capability.KindEnvironment,
			Reason: "subject root has no " + allowlistedProbePath,
		}, true},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			probe := runStrippedGate(t, map[string]string{
				"conformance":    "test -f " + allowlistedProbePath,
				canary.PhaseTest: softSkipScript(t, tc.skip),
			})
			if tc.wantRed {
				if probe.ExitCode == 0 {
					t.Fatalf("environment-kind skip left the full gate green:\n%s\n%s", probe.Stdout, probe.Stderr)
				}
				requireOutput(t, probe, "fatal against the stripped subject")
				return
			}
			probe.RequireExit(0)
			requireOutput(t, probe, "phase "+canary.PhaseTest+": green")
		})
	}
}

// [R09 boundary] The environment kind alone is not the signature. Excludable phases emit
// ordinary kind=environment skips on every real host — an absent non-declared subject
// file, an unmaterialized fixture, a release plan with no target for this host — none of
// which stripping can induce, because the declaration removes only its own paths. Reading
// the kind as the whole signature reds every full gate on a normal host, so each case here
// carries a reason naming nothing the declaration strips and must leave the gate green.
func TestOrdinaryEnvironmentSkipDoesNotRedStrippedRun(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "stripped-subject ordinary-environment-skip contract failed")
	scope := gate.ReducedScope()
	for _, reason := range []string{
		"subject root has no bin/bench.sh",
		"authoritative native-proof canary fixture is not materialized",
		"package contract tests require release plan target for host linux/amd64",
		"BENCH_CONFORMANCE_ROOT not set",
	} {
		reason := reason
		t.Run(reason, func(t *testing.T) {
			t.Parallel()
			for _, field := range strings.Fields(reason) {
				if scope.Member(field) {
					t.Fatalf("fixture reason %q names the declared path %q; repoint the fixture", reason, field)
				}
			}
			probe := runStrippedGate(t, map[string]string{
				"conformance": "test -f " + allowlistedProbePath,
				canary.PhaseTest: softSkipScript(t, capability.Skip{
					Kind:   capability.KindEnvironment,
					Reason: reason,
				}),
			})
			probe.RequireExit(0)
			requireOutput(t, probe, "phase "+canary.PhaseTest+": green")
		})
	}
}

// softSkipScript plants skipIfSubjectFileMissing's shape as a phase script: when the
// probe path is absent, deliver one structured line to the gate's skip log and exit
// green. The line itself comes from capability.Render so the fixture cannot drift from
// the shape the collector parses. Render's terminating newline is stripped and re-added
// by printf: %q would escape it into a literal backslash-n, which the shell passes
// through verbatim and which then rides on the end of the parsed reason.
func softSkipScript(t *testing.T, skip capability.Skip) string {
	t.Helper()
	line, err := capability.Render(skip)
	if err != nil {
		t.Fatalf("render fixture skip line: %v", err)
	}
	return fmt.Sprintf(`if [ ! -f %s ]; then printf '%%s\n' %q >> "$%s"; exit 0; fi
exit 0
`, allowlistedProbePath, strings.TrimSuffix(line, "\n"), capability.LogEnv)
}

// [R10] Included phases still find the allowlisted paths on the same run: stripping is
// scoped to the excludable set, and over-stripping would break exactly the checks that
// exist to grade these files. The excludable probe asserts absence, so this contract
// also pins that the stripped construction ran at all.
func TestIncludedPhasesSeeAllowlistedPaths(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "included-phase real-root contract failed")
	probe := runStrippedGate(t, map[string]string{
		"conformance":    "test -f " + allowlistedProbePath,
		canary.PhaseTest: "test ! -f " + allowlistedProbePath,
	})
	probe.RequireExit(0)
	requireOutput(t, probe, "phase conformance: green")
	requireOutput(t, probe, "phase "+canary.PhaseTest+": green")
}

// [RB2] A root that is not the kit runs its full gate unsplit, materializing no stripped
// worktree. The construction enforces the kit's own reduced-scope declaration, and a
// linked repo gating through the shipped binary — BENCH_KIT naming the kit checkout, the
// graded root a different tree — never declared the allowlist being enforced. The planted
// phase reads an allowlisted path, exactly the dependency the split reds on the kit, and
// must stay green here because every phase graded the real tree.
func TestForeignRootRunsUnsplit(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "foreign-root unsplit contract failed")
	probe := runStrippedGateWithKit(t, map[string]string{
		"conformance":    "test -f " + allowlistedProbePath,
		canary.PhaseTest: "cat " + allowlistedProbePath + " >/dev/null",
	}, t.TempDir())
	probe.RequireExit(0)
	requireOutput(t, probe, "phase "+canary.PhaseTest+": green")
	requireOutput(t, probe, "phase conformance: green")
	if strings.Contains(probe.Stdout+probe.Stderr, "stripped subject") {
		t.Fatalf("foreign root materialized a stripped subject:\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
}

// [R11] The materialized stripped tree is a real git repository: a contract test stages
// its subject through `git ls-files`, which fails outright against a plain directory
// copy. The probe asserts the stripped root first so it cannot pass by running against
// the real repository.
func TestStrippedSubjectIsGitRepository(t *testing.T) {
	t.Parallel()
	contract.NoteContractFailure(t, "stripped-subject git-repository contract failed")
	probe := runStrippedGate(t, map[string]string{
		"conformance": "test -f " + allowlistedProbePath,
		canary.PhaseTest: "test ! -f " + allowlistedProbePath +
			" && git ls-files --error-unmatch graded.txt >/dev/null",
	})
	probe.RequireExit(0)
	requireOutput(t, probe, "phase "+canary.PhaseTest+": green")
}

// runStrippedGate drives the built binary's full phase runner over a fixture repository
// whose .bench/phases.json declares one bash script per named phase. Phase names carry
// the declaration's meaning — "conformance" is included, the test phase is excludable —
// so the fixture guards both memberships before relying on them.
func runStrippedGate(t *testing.T, phases map[string]string) contract.Probe {
	t.Helper()
	return runStrippedGateWithKit(t, phases, "")
}

// runStrippedGateWithKit is runStrippedGate with the kit identity chosen by the case:
// empty means the fixture is its own kit (the split-eligible shape), while a non-empty
// path plays a linked repo's wrapper naming a kit checkout elsewhere.
func runStrippedGateWithKit(t *testing.T, phases map[string]string, kit string) contract.Probe {
	t.Helper()
	contract.RequireFreshBench(t)
	scope := gate.ReducedScope()
	if !scope.Member(allowlistedProbePath) {
		t.Fatalf("fixture probe path %s is no longer declared; repoint the fixture", allowlistedProbePath)
	}
	if !scope.Excludable(canary.PhaseTest) || scope.Excludable("conformance") {
		t.Fatal("fixture phase names no longer match the declaration; repoint the fixture")
	}
	root := t.TempDir()
	entries := make([]string, 0, len(phases))
	for _, name := range []string{"conformance", canary.PhaseTest} {
		script, declared := phases[name]
		if !declared {
			continue
		}
		rel := ".bench/phase-" + name + ".sh"
		contract.WriteFileAbs(t, filepath.Join(root, filepath.FromSlash(rel)), "set -u\n"+script+"\n")
		entries = append(entries, fmt.Sprintf(`{"name":%q,"argv":["bash",%q]}`, name, rel))
	}
	manifest := `{"phases":[` + strings.Join(entries, ",") + `]}` + "\n"
	contract.WriteFileAbs(t, filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath)), manifest)
	contract.WriteFileAbs(t, filepath.Join(root, allowlistedProbePath), "allowlisted fixture content\n")
	contract.WriteFileAbs(t, filepath.Join(root, "graded.txt"), "tracked fixture content\n")
	fixtureGit(t, root, "init", "-q")
	fixtureGit(t, root, "add", "-A")
	fixtureGit(t, root, "commit", "-q", "-m", "fixture")

	f := contract.NewExecFixtureAt(t, root)
	if kit == "" {
		kit = root
	}
	run := contract.Env{
		"BENCH_KIT":                  &kit,
		"BENCH_CANARY_INNER":         nil,
		"BENCH_REQUIRE_CAPABILITIES": nil,
		capability.LogEnv:            nil,
	}
	return f.RunEnvSpec(run, filepath.Join(contract.SubjectRoot(t), "dist", "bench"), "gate-phases", root)
}

// fixtureGit runs git in the fixture repository with an inline identity, so the fixture
// commits on hosts with no git user configured.
func fixtureGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root,
		"-c", "user.name=bench-fixture", "-c", "user.email=fixture@bench.invalid"}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func requireOutput(t *testing.T, probe contract.Probe, want string) {
	t.Helper()
	if !strings.Contains(probe.Stdout+probe.Stderr, want) {
		t.Fatalf("gate output missing %q:\nstdout:\n%s\nstderr:\n%s", want, probe.Stdout, probe.Stderr)
	}
}
