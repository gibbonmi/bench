package canary

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

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
	missingExpect := canaryFixture(root, "test-family", "missing-expect")
	missingFiles := canaryFixture(root, "test-family", "missing-files")
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
	fixture := canaryFixture(root, "test-family", "generic")
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
	fixture := canaryFixture(root, "test-family", "dot-restore")
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
	if envHasPrefix(fixtureCall.Env, "BENCH_KIT=") || envHasPrefix(fixtureCall.Env, "BENCH_WRAPPER=") {
		t.Fatalf("inner env leaked wrapper routing: %#v", fixtureCall.Env)
	}
}

func TestSweepReportsDidNotBite(t *testing.T) {
	root := t.TempDir()
	fixture := canaryFixture(root, "test-family", "weak")
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

func TestSweepAcceptsLegacyFlatSeedCanary(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "example")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "example check\n")

	var fixtureCalls []string
	err := Sweep(root, func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 0, Output: "baseline\n"}
		}
		fixtureCalls = append(fixtureCalls, call.FixtureDir)
		return RunResult{ExitCode: 1, Output: "example check\n"}
	})
	if err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if len(fixtureCalls) != 1 || fixtureCalls[0] != fixture {
		t.Fatalf("fixture calls = %#v, want only legacy fixture %q", fixtureCalls, fixture)
	}
}

func TestSweepUsesLiteralFixturePathWithSpacesAndGlobCharacters(t *testing.T) {
	root := t.TempDir()
	fixture := canaryFixture(root, "family with spaces [abc]", "fixture * with spaces")
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
	for _, family := range []string{"alpha-family", "bravo-family"} {
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
