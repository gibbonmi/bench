package canary

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/capability"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSweepRunsFixturesConcurrently(t *testing.T) {
	if runtime.NumCPU() < 2 {
		capability.Capability(t, capability.CPU, "NumCPU=1 makes overlap impossible by policy")
	}
	root := t.TempDir()
	for _, name := range []string{"a", "b"} {
		fixture := canaryFixture(root, "test-family", name)
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
		fixture := canaryFixture(root, "test-family", name)
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
		fixture := canaryFixture(root, "test-family", name)
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
		capability.Capability(t, capability.CPU, "NumCPU=1 makes reverse completion impossible by policy")
	}
	root := t.TempDir()
	for _, name := range []string{"alpha", "bravo"} {
		fixture := canaryFixture(root, "test-family", name)
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
	fixture := canaryFixture(root, "test-family", "valid")
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
	valid := canaryFixture(root, "test-family", "valid")
	mkdir(t, filepath.Join(valid, "files"))
	write(t, filepath.Join(valid, "EXPECT"), "target-valid\n")
	vacuous := canaryFixture(root, "test-family", "vacuous")
	mkdir(t, filepath.Join(vacuous, "files"))
	write(t, filepath.Join(vacuous, "EXPECT"), "vacuous\n")
	brokenLink := canaryFixture(root, "test-family", "broken-link")
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
