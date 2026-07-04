package canary

import (
	"os"
	"path/filepath"
	"strings"
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
	mkdir(t, filepath.Join(root, "tests", "canary", "missing-expect", "files"))
	mkdir(t, filepath.Join(root, "tests", "canary", "missing-files"))
	write(t, filepath.Join(root, "tests", "canary", "missing-files", "EXPECT"), "target\n")

	err := Sweep(root, (&recordingRunner{}).Run)
	if err == nil {
		t.Fatal("Sweep err = nil, want malformed fixture errors")
	}
	got := err.Error()
	for _, want := range []string{
		"canary fixture 'missing-expect' has no EXPECT file",
		"canary fixture 'missing-files' has no files/ tree",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Sweep err missing %q:\n%s", want, got)
		}
	}
}

func TestSweepRejectsVacuousExpect(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "generic")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "generic failure\n")

	runner := &recordingRunner{outputs: map[string]RunResult{
		"baseline": {ExitCode: 1, Output: "generic failure\n"},
	}}
	err := Sweep(root, runner.Run)
	if err == nil || !strings.Contains(err.Error(), "canary 'generic' EXPECT is vacuous") {
		t.Fatalf("Sweep err = %v, want vacuous EXPECT", err)
	}
}

func TestSweepMaterializesFixtureAndRequiresTargetedBite(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "dot-restore")
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
	fixture := filepath.Join(root, "tests", "canary", "weak")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "target\n")

	runner := &recordingRunner{outputs: map[string]RunResult{
		"baseline": {ExitCode: 1, Output: "other\n"},
		"weak":     {ExitCode: 0, Output: "target\n"},
	}}
	err := Sweep(root, runner.Run)
	if err == nil || !strings.Contains(err.Error(), `canary 'weak' did not bite`) {
		t.Fatalf("Sweep err = %v, want did-not-bite", err)
	}
}

type recordingRunner struct {
	calls           []RunCall
	outputs         map[string]RunResult
	sawMaterialized map[string]bool
}

func (r *recordingRunner) Run(call RunCall) RunResult {
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
	if len(r.calls) > 1 {
		key = filepath.Base(call.FixtureDir)
	}
	if out, ok := r.outputs[key]; ok {
		return out
	}
	return RunResult{ExitCode: 1, Output: "other failure\n"}
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
