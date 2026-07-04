package canary

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
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
	fixture := filepath.Join(root, "tests", "canary", "generic")
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
	want := `canary 'weak' did not bite (want red + "target"; got exit 0)`
	if err == nil || err.Error() != want {
		t.Fatalf("Sweep err = %v, want %s", err, want)
	}
}

func TestSweepRunsFixturesConcurrently(t *testing.T) {
	if runtime.NumCPU() < 2 {
		t.Skip("NumCPU=1 makes overlap impossible by policy")
	}
	root := t.TempDir()
	for _, name := range []string{"a", "b"} {
		fixture := filepath.Join(root, "tests", "canary", name)
		mkdir(t, filepath.Join(fixture, "files"))
		write(t, filepath.Join(fixture, "EXPECT"), "target-"+name+"\n")
	}

	secondEntered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	var releaseOnce sync.Once
	var mu sync.Mutex
	inFlight := 0
	runner := func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		mu.Lock()
		inFlight++
		if inFlight == 2 {
			once.Do(func() { close(secondEntered) })
		}
		mu.Unlock()
		select {
		case <-secondEntered:
		case <-time.After(2 * time.Second):
			t.Errorf("second fixture run did not overlap first")
		}
		releaseOnce.Do(func() { close(release) })
		<-release
		mu.Lock()
		inFlight--
		mu.Unlock()
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}

	if err := Sweep(root, runner); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
}

func TestSweepBoundsFixtureConcurrencyAtNumCPU(t *testing.T) {
	root := t.TempDir()
	fixtureCount := runtime.NumCPU() + 3
	for i := 0; i < fixtureCount; i++ {
		name := fmt.Sprintf("fx-%02d", i)
		fixture := filepath.Join(root, "tests", "canary", name)
		mkdir(t, filepath.Join(fixture, "files"))
		write(t, filepath.Join(fixture, "EXPECT"), "target-"+name+"\n")
	}

	release := make(chan struct{})
	var releaseOnce sync.Once
	var mu sync.Mutex
	inFlight := 0
	highWater := 0
	runner := func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		mu.Lock()
		inFlight++
		if inFlight > highWater {
			highWater = inFlight
		}
		if inFlight == runtime.NumCPU() {
			releaseOnce.Do(func() { close(release) })
		}
		mu.Unlock()
		<-release
		mu.Lock()
		inFlight--
		mu.Unlock()
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}

	if err := Sweep(root, runner); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if highWater > runtime.NumCPU() {
		t.Fatalf("fixture concurrency high-water = %d, want <= NumCPU %d", highWater, runtime.NumCPU())
	}
}

func TestSweepCompletesBaselineBeforeStartingFixtures(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a", "b"} {
		fixture := filepath.Join(root, "tests", "canary", name)
		mkdir(t, filepath.Join(fixture, "files"))
		write(t, filepath.Join(fixture, "EXPECT"), "target-"+name+"\n")
	}

	baselineDone := false
	var mu sync.Mutex
	runner := func(call RunCall) RunResult {
		mu.Lock()
		defer mu.Unlock()
		if call.FixtureDir == "" {
			baselineDone = true
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		if !baselineDone {
			t.Errorf("fixture %s started before baseline completed", filepath.Base(call.FixtureDir))
		}
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}

	if err := Sweep(root, runner); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
}

func TestSweepReportsErrorsInSortedFixtureOrder(t *testing.T) {
	if runtime.NumCPU() < 2 {
		t.Skip("NumCPU=1 makes reverse completion impossible by policy")
	}
	root := t.TempDir()
	for _, name := range []string{"alpha", "bravo"} {
		fixture := filepath.Join(root, "tests", "canary", name)
		mkdir(t, filepath.Join(fixture, "files"))
		write(t, filepath.Join(fixture, "EXPECT"), "target-"+name+"\n")
	}

	release := map[string]chan struct{}{}
	for _, name := range []string{"alpha", "bravo"} {
		release[name] = make(chan struct{})
	}
	started := make(chan string, 2)
	runner := func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		name := filepath.Base(call.FixtureDir)
		started <- name
		<-release[name]
		return RunResult{ExitCode: 0, Output: "target-" + name + "\n"}
	}

	done := make(chan error, 1)
	go func() {
		done <- Sweep(root, runner)
	}()
	waitStarted(t, started, 2)
	close(release["bravo"])
	close(release["alpha"])
	err := <-done
	if err == nil {
		t.Fatal("Sweep err = nil, want did-not-bite errors")
	}
	want := strings.Join([]string{
		`canary 'alpha' did not bite (want red + "target-alpha"; got exit 0)`,
		`canary 'bravo' did not bite (want red + "target-bravo"; got exit 0)`,
	}, "\n")
	if err.Error() != want {
		t.Fatalf("Sweep err:\n%s\nwant:\n%s", err, want)
	}
}

func TestSweepRemovesTempWorkDirsOnGreenPath(t *testing.T) {
	tmpRoot := t.TempDir()
	t.Setenv("TMPDIR", tmpRoot)
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "valid")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "target-valid\n")

	runner := func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}
	if err := Sweep(root, runner); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	assertDirEmpty(t, tmpRoot)
}

func TestSweepRemovesTempWorkDirsOnRedPaths(t *testing.T) {
	tmpRoot := t.TempDir()
	t.Setenv("TMPDIR", tmpRoot)
	root := t.TempDir()
	valid := filepath.Join(root, "tests", "canary", "valid")
	mkdir(t, filepath.Join(valid, "files"))
	write(t, filepath.Join(valid, "EXPECT"), "target-valid\n")
	vacuous := filepath.Join(root, "tests", "canary", "vacuous")
	mkdir(t, filepath.Join(vacuous, "files"))
	write(t, filepath.Join(vacuous, "EXPECT"), "vacuous\n")
	brokenLink := filepath.Join(root, "tests", "canary", "broken-link")
	mkdir(t, filepath.Join(brokenLink, "files"))
	write(t, filepath.Join(brokenLink, "EXPECT"), "target-broken-link\n")
	if err := os.Symlink("missing-target", filepath.Join(brokenLink, "files", "broken")); err != nil {
		t.Fatal(err)
	}

	runner := func(call RunCall) RunResult {
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "vacuous\n"}
		}
		return RunResult{ExitCode: 0, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}
	if err := Sweep(root, runner); err == nil {
		t.Fatal("Sweep err = nil, want fixture errors")
	} else if !strings.Contains(err.Error(), "canary 'broken-link' setup failed:") {
		t.Fatalf("Sweep err = %v, want broken-link setup failure", err)
	}
	assertDirEmpty(t, tmpRoot)
}

func assertDirEmpty(t *testing.T, path string) {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Fatalf("temporary entries left behind: %v", names)
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

func waitStarted(t *testing.T, started <-chan string, want int) {
	t.Helper()
	seen := map[string]bool{}
	timeout := time.After(2 * time.Second)
	for len(seen) < want {
		select {
		case name := <-started:
			seen[name] = true
		case <-timeout:
			t.Fatalf("started fixtures = %v, want %d", seen, want)
		}
	}
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
