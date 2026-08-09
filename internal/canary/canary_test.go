package canary

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestSweepPreservesOneExactSelectedBenchPath(t *testing.T) {
	selected := filepath.Join(t.TempDir(), "selected bench")
	t.Setenv(runbinary.Env, selected)

	env := gateEnv(sweepEnv(), registry.Dev)
	if values := envValues(env, runbinary.Env); !reflect.DeepEqual(values, []string{selected}) {
		t.Fatalf("inner gate selected paths = %q, want exact inherited path", values)
	}
}

// TestRunPrintsCanaryOkOnCleanSweep pins the `bench canary` success feedback,
// mirroring `structure ok`: a clean sweep must say so on stdout instead of
// exiting 0 silently.
func TestRunPrintsCanaryOkOnCleanSweep(t *testing.T) {
	root := t.TempDir()
	fixture := canaryFixture(root, mappedFamily(t), "mybreak")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "boom detected\n")
	write(t, filepath.Join(fixture, "files", "marker.txt"), "x\n")
	// The phase line prints before the branch, as the real gate's do: a gate silent on
	// an empty tree yields an empty vacuity baseline, which grades nothing and is a
	// sweep error rather than a clean sweep.
	write(t, filepath.Join(root, ".bench", "gate.sh"),
		"#!/bin/sh\necho \"gate: running\"\nif [ -f marker.txt ]; then echo \"boom detected\"; exit 1; else exit 0; fi\n")

	var stdout, stderr bytes.Buffer
	code := Run([]string{root}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("Run = %d, want 0; stderr:\n%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "canary ok") {
		t.Fatalf("stdout = %q, want it to contain %q", stdout.String(), "canary ok")
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty on a clean sweep", stderr.String())
	}
}

func TestSweepRejectsMissingAndEmptyHarness(t *testing.T) {
	root := t.TempDir()
	runner := &recordingRunner{}

	err := Sweep(root, runner.Run)
	if err == nil || !strings.Contains(err.Error(), "canary harness absent") {
		t.Fatalf("missing harness err = %v, want absent harness", err)
	}

	if err := os.MkdirAll(filepath.Join(root, "tests", "canary"), 0o755); err != nil {
		t.Fatal(err)
	}
	err = Sweep(root, runner.Run)
	if err == nil || !strings.Contains(err.Error(), "canary harness absent") {
		t.Fatalf("empty harness err = %v, want absent harness", err)
	}
}

func TestSweepRejectsMalformedFixtures(t *testing.T) {
	root := t.TempDir()
	missingExpect := canaryFixture(root, mappedFamily(t), "missing-expect")
	missingFiles := canaryFixture(root, mappedFamily(t), "missing-files")
	mkdir(t, filepath.Join(missingExpect, "files"))
	mkdir(t, missingFiles)
	write(t, filepath.Join(missingFiles, "EXPECT"), "target\n")

	err := Sweep(root, (&recordingRunner{}).Run)
	if err == nil {
		t.Fatal("Sweep err = nil, want malformed fixture errors")
	}
	want := strings.Join([]string{
		"canary fixture 'missing-expect' has no EXPECT file",
		"canary fixture 'missing-files' has no files/ tree",
	}, "\n")
	if err.Error() != want {
		t.Fatalf("Sweep err:\n%s\nwant:\n%s", err, want)
	}
}

func TestSweepRejectsVacuousExpect(t *testing.T) {
	root := t.TempDir()
	fixture := canaryFixture(root, mappedFamily(t), "generic")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "generic failure\n")

	runner := &recordingRunner{outputs: map[string]RunResult{
		"baseline": {ExitCode: 1, Output: "generic failure\n"},
	}}
	err := Sweep(root, runner.Run)
	want := "canary 'generic' EXPECT is vacuous (also matches an empty fixture)"
	if err == nil || err.Error() != want {
		t.Fatalf("Sweep err = %v, want %s", err, want)
	}
}

func TestSweepMaterializesFixtureAndRequiresTargetedBite(t *testing.T) {
	root := t.TempDir()
	fixture := canaryFixture(root, mappedFamily(t), "dot-restore")
	mkdir(t, filepath.Join(fixture, "files", "dot-bench", "hooks"))
	mkdir(t, filepath.Join(fixture, "files", "nested", "dot-codex"))
	write(t, filepath.Join(fixture, "EXPECT"), "targeted regression\n")
	write(t, filepath.Join(fixture, "files", "dot-bench", "hooks", "x.sh"), "#!/bin/sh\n")
	write(t, filepath.Join(fixture, "files", "nested", "dot-codex", "hooks.json"), "{}\n")

	runner := &recordingRunner{outputs: map[string]RunResult{
		"baseline":    {ExitCode: 1, Output: "other failure\n"},
		"dot-restore": {ExitCode: 1, Output: "targeted regression\n"},
	}}
	if err := Sweep(root, runner.Run); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if len(runner.calls) != 2 {
		t.Fatalf("runner calls = %d, want baseline + fixture", len(runner.calls))
	}
	fixtureCall := runner.calls[1]
	if fixtureCall.Cwd == "" {
		t.Fatal("fixture call did not record cwd")
	}
	for _, want := range []string{".bench/hooks/x.sh", "nested/.codex/hooks.json"} {
		if !runner.sawMaterialized[want] {
			t.Fatalf("runner did not observe materialized %s; saw %#v", want, runner.sawMaterialized)
		}
	}
	if !envHas(fixtureCall.Env, "BENCH_CANARY_INNER=1") {
		t.Fatalf("inner env missing BENCH_CANARY_INNER=1: %#v", fixtureCall.Env)
	}
	if !envHas(fixtureCall.Env, "BENCH_CANARY_PHASE=conformance") {
		t.Fatalf("inner env missing fixture phase: %#v", fixtureCall.Env)
	}
	if envHasPrefix(fixtureCall.Env, "BENCH_KIT=") || envHasPrefix(fixtureCall.Env, "BENCH_WRAPPER=") {
		t.Fatalf("inner env leaked wrapper routing: %#v", fixtureCall.Env)
	}
}

func TestSweepReportsDidNotBite(t *testing.T) {
	root := t.TempDir()
	fixture := canaryFixture(root, mappedFamily(t), "weak")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "target\n")

	runner := &recordingRunner{outputs: map[string]RunResult{
		"baseline": {ExitCode: 1, Output: "other\n"},
		"weak":     {ExitCode: 0, Output: "target\n"},
	}}
	err := Sweep(root, runner.Run)
	want := `canary 'weak' did not bite (want red + "target"; got exit 0)`
	if err == nil || err.Error() != want {
		t.Fatalf("Sweep err = %v, want %s", err, want)
	}
}

// TestSweepReportsDidNotBiteWhenARedRunOmitsItsExpect is the other half of the bite
// condition: the test above grades a green run whose output matched, this one a red run
// whose output did not. A fixture whose EXPECT stops appearing while its run still fails
// is the failure a red-alone check reports as a bite forever, and it is reached through
// the compiled path because that is where a behavior-owned fixture's output comes from.
func TestSweepReportsDidNotBiteWhenARedRunOmitsItsExpect(t *testing.T) {
	root := t.TempDir()
	contractFixture(t, root, "axi", "axi-fx")

	err := Sweep(root, func(call RunCall) RunResult {
		if result, done := stubToolchain(call); done {
			return result
		}
		return RunResult{ExitCode: 1, Output: "unrelated failure\n"}
	})

	want := `canary 'axi-fx' did not bite (want red + "target-axi-fx"; got exit 1)`
	if err == nil || err.Error() != want {
		t.Fatalf("Sweep err = %v, want %s", err, want)
	}
}

func TestSweepAcceptsLegacyFlatSeedCanary(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "example")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "example check\n")
	t.Setenv(PhaseEnv, "contract")

	var fixtureCalls []RunCall
	err := Sweep(root, func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 0, Output: "baseline\n"}
		}
		fixtureCalls = append(fixtureCalls, call)
		return RunResult{ExitCode: 1, Output: "example check\n"}
	})
	if err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if len(fixtureCalls) != 1 || fixtureCalls[0].FixtureDir != fixture {
		t.Fatalf("fixture calls = %#v, want only legacy fixture %q", fixtureCalls, fixture)
	}
	if envHasPrefix(fixtureCalls[0].Env, PhaseEnv+"=") {
		t.Fatalf("legacy flat fixture inherited targeted phase: %#v", fixtureCalls[0].Env)
	}
}

func TestSweepUsesLiteralFixturePathWithSpacesAndGlobCharacters(t *testing.T) {
	root := t.TempDir()
	// A legacy flat fixture is where a hostile name still fits: family directories
	// name conformance checks, so the spaces and glob metacharacters live in the
	// fixture's own base name and the sweep resolves it to no scope.
	fixture := filepath.Join(root, "tests", "canary", "fixture * with spaces [abc]")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "literal path check\n")

	var fixtureCalls []string
	err := Sweep(root, func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 0, Output: "baseline\n"}
		}
		fixtureCalls = append(fixtureCalls, call.FixtureDir)
		return RunResult{ExitCode: 1, Output: "literal path check\n"}
	})
	if err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if len(fixtureCalls) != 1 || fixtureCalls[0] != fixture {
		t.Fatalf("fixture calls = %#v, want exact fixture path %q", fixtureCalls, fixture)
	}
}

func TestSweepRejectsDuplicateFixtureBaseNames(t *testing.T) {
	root := t.TempDir()
	for _, family := range []string{mappedFamily(t), secondMappedFamily(t)} {
		fixture := canaryFixture(root, family, "duplicate")
		mkdir(t, filepath.Join(fixture, "files"))
		write(t, filepath.Join(fixture, "EXPECT"), "target\n")
	}

	err := Sweep(root, (&recordingRunner{}).Run)
	want := `canary fixture name "duplicate" appears in multiple families; base names must be globally unique`
	if err == nil || err.Error() != want {
		t.Fatalf("Sweep err = %v, want %s", err, want)
	}
}

type recordingRunner struct {
	mu              sync.Mutex
	calls           []RunCall
	outputs         map[string]RunResult
	sawMaterialized map[string]bool
}

func (r *recordingRunner) Run(call RunCall) RunResult {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, call)
	if result, done := stubToolchain(call); done {
		return result
	}
	if r.sawMaterialized == nil {
		r.sawMaterialized = map[string]bool{}
	}
	for _, rel := range []string{".bench/hooks/x.sh", "nested/.codex/hooks.json"} {
		if _, err := os.Stat(filepath.Join(call.Cwd, rel)); err == nil {
			r.sawMaterialized[rel] = true
		}
	}
	key := "baseline"
	if call.FixtureDir != "" {
		key = filepath.Base(call.FixtureDir)
	}
	if out, ok := r.outputs[key]; ok {
		return out
	}
	return RunResult{ExitCode: 1, Output: "other failure\n"}
}

// mappedFamily and secondMappedFamily are conformance families the registry's family
// table binds, read rather than written down so the synthetic trees follow a rename or
// a rebinding instead of pinning two names. Synthetic trees have to sit under real
// families: a conformance family the table does not bind runs unscoped, so an invented
// name silently drops the scoping these tests grade.
func mappedFamily(t *testing.T) string {
	t.Helper()
	first, _ := mappedFamilies(t)
	return first
}

func secondMappedFamily(t *testing.T) string {
	t.Helper()
	_, second := mappedFamilies(t)
	return second
}

// mappedFamilies picks two bound families whose checks differ. The scope-group tests
// need the pair to land in separate groups, and several families in the table share one
// check, so the distinct-check condition is what the pair is selected on.
func mappedFamilies(t *testing.T) (string, string) {
	t.Helper()
	families := registry.Families()
	for _, first := range families {
		firstCheck, _ := registry.FamilyCheck(first)
		for _, second := range families {
			if secondCheck, _ := registry.FamilyCheck(second); secondCheck != firstCheck {
				return first, second
			}
		}
	}
	t.Fatal("registry family table binds no two families to different checks")
	return "", ""
}

func canaryFixture(root, family, name string) string {
	return filepath.Join(root, "tests", "canary", family, name)
}

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func envHas(env []string, want string) bool {
	for _, got := range env {
		if got == want {
			return true
		}
	}
	return false
}

func envHasPrefix(env []string, prefix string) bool {
	for _, got := range env {
		if strings.HasPrefix(got, prefix) {
			return true
		}
	}
	return false
}

func TestRunUsage(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantCode int
		// wantOn is the stream the line must land on: stdout for help, stderr for misuse.
		helpOnStdout bool
	}{
		{name: "help flag", args: []string{"--help"}, wantCode: 0, helpOnStdout: true},
		{name: "short help", args: []string{"-h"}, wantCode: 0, helpOnStdout: true},
		{name: "bare help", args: []string{"help"}, wantCode: 0, helpOnStdout: true},
		{name: "unknown flag", args: []string{"--nope"}, wantCode: 2},
		{name: "excess arguments", args: []string{"a", "b"}, wantCode: 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(tc.args, &stdout, &stderr)
			if code != tc.wantCode {
				t.Fatalf("Run(%q) = %d, want %d; stderr:\n%s", tc.args, code, tc.wantCode, stderr.String())
			}
			got, other := stderr.String(), stdout.String()
			if tc.helpOnStdout {
				got, other = stdout.String(), stderr.String()
			}
			wantLine := "usage: bench canary"
			if tc.helpOnStdout {
				wantLine = "usage: bench canary [root]"
			}
			if !strings.Contains(got, wantLine) {
				t.Fatalf("Run(%q) output = %q, want it to contain %q", tc.args, got, wantLine)
			}
			if strings.Contains(other, "canary harness absent") || strings.Contains(other, "no fixtures") {
				t.Fatalf("Run(%q) reached the sweep; other stream = %q", tc.args, other)
			}
		})
	}
}
