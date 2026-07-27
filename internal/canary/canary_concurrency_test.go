package canary

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

func TestSweepRunsFixturesConcurrently(t *testing.T) {
	if bound := fixtureWorkers(runtime.GOMAXPROCS(0), 2); bound < 2 {
		capability.Capability(t, capability.CPU, fmt.Sprintf("derived worker bound %d makes overlap impossible by policy", bound))
	}
	root := t.TempDir()
	for _, name := range []string{"a", "b"} {
		fixture := canaryFixture(root, mappedFamily(t), name)
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

func TestSweepBoundsFixtureConcurrencyAtDerivedWorkerBound(t *testing.T) {
	root := t.TempDir()
	budget := runtime.GOMAXPROCS(0)
	fixtureCount := budget + 3
	for i := 0; i < fixtureCount; i++ {
		name := fmt.Sprintf("fx-%02d", i)
		fixture := canaryFixture(root, mappedFamily(t), name)
		mkdir(t, filepath.Join(fixture, "files"))
		write(t, filepath.Join(fixture, "EXPECT"), "target-"+name+"\n")
	}
	bound := fixtureWorkers(budget, fixtureCount)

	release := make(chan struct{})
	reachedBound := make(chan struct{})
	var reachedOnce sync.Once
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
		if inFlight == bound {
			reachedOnce.Do(func() { close(reachedBound) })
		}
		mu.Unlock()
		<-release
		mu.Lock()
		inFlight--
		mu.Unlock()
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}

	// A worker must not both signal reaching the bound and block on the settle,
	// so the settle-and-release lives on its own goroutine. A budget mismatch
	// between this test's GOMAXPROCS read and runFixtures' would otherwise leave
	// every worker blocked on <-release forever; the timeout turns that into a
	// reported failure instead of a hang.
	go func() {
		select {
		case <-reachedBound:
			time.Sleep(150 * time.Millisecond)
		case <-time.After(2 * time.Second):
			t.Errorf("in-flight count never reached the derived bound %d", bound)
		}
		close(release)
	}()

	if err := Sweep(root, runner); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if highWater != bound {
		t.Fatalf("fixture concurrency high-water = %d, want == derived bound %d", highWater, bound)
	}
}

// TestSweepCompletesEachGroupsBaselineBeforeItsFixtures grades the barrier per scope
// group, which is the only form of it that holds: a fixture is graded vacuous against
// its own group's baseline alone, so a run that starts before that baseline finishes
// compares against an absent output and clears the vacuity check unguarded. A single
// tracked flag would pass on any baseline finishing before any fixture, which a sweep
// that raced a later group's baseline against that group's fixtures still satisfies —
// hence the three groups and the per-group flags.
func TestSweepCompletesEachGroupsBaselineBeforeItsFixtures(t *testing.T) {
	root := t.TempDir()
	first, second := mappedFamilies(t)
	for _, family := range []string{first, second, "behavior-owned"} {
		for _, name := range []string{"a", "b"} {
			fixture := canaryFixture(root, family, family+"-"+name)
			mkdir(t, filepath.Join(fixture, "files"))
			write(t, filepath.Join(fixture, "EXPECT"), "target-"+family+"-"+name+"\n")
		}
	}

	var mu sync.Mutex
	baselineDone := map[string]bool{}
	runner := func(call RunCall) RunResult {
		group := callGroup(t, call)
		mu.Lock()
		defer mu.Unlock()
		if call.FixtureDir == "" {
			baselineDone[group] = true
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		if !baselineDone[group] {
			t.Errorf("fixture %s started before the baseline of its own group %q completed", filepath.Base(call.FixtureDir), group)
		}
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	}

	if err := Sweep(root, runner); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if len(baselineDone) != 3 {
		t.Fatalf("sweep ran baselines for %d groups, want 3 distinct groups", len(baselineDone))
	}
}

// callGroup names the scope group a call belongs to, which is the key the sweep pairs a
// fixture with its baseline under. More than one scope on a call means the pairing is
// ambiguous, so it fails rather than picking one.
func callGroup(t *testing.T, call RunCall) string {
	t.Helper()
	scopes := scopeValues(call.Env)
	switch len(scopes) {
	case 0:
		return ""
	case 1:
		return scopes[0]
	default:
		t.Fatalf("call %q carried scopes %v, want at most one", call.FixtureDir, scopes)
		return ""
	}
}

func TestSweepReportsErrorsInSortedFixtureOrder(t *testing.T) {
	if bound := fixtureWorkers(runtime.GOMAXPROCS(0), 2); bound < 2 {
		capability.Capability(t, capability.CPU, fmt.Sprintf("derived worker bound %d makes reverse completion impossible by policy", bound))
	}
	root := t.TempDir()
	for _, name := range []string{"alpha", "bravo"} {
		fixture := canaryFixture(root, mappedFamily(t), name)
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
	fixture := canaryFixture(root, mappedFamily(t), "valid")
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
	valid := canaryFixture(root, mappedFamily(t), "valid")
	mkdir(t, filepath.Join(valid, "files"))
	write(t, filepath.Join(valid, "EXPECT"), "target-valid\n")
	vacuous := canaryFixture(root, mappedFamily(t), "vacuous")
	mkdir(t, filepath.Join(vacuous, "files"))
	write(t, filepath.Join(vacuous, "EXPECT"), "vacuous\n")
	brokenLink := canaryFixture(root, mappedFamily(t), "broken-link")
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

func TestFixtureWorkers(t *testing.T) {
	tests := []struct {
		name         string
		budget       int
		fixtureCount int
		want         int
	}{
		{"negative budget floors at one", -4, 1000, 1},
		{"budget 0 floors at one", 0, 1000, 1},
		{"budget 1 floors at one", 1, 1000, 1},
		{"budget 2 divides evenly", 2, 1000, 1},
		{"budget 3 floors the division", 3, 1000, 1},
		{"budget 16 divides by the width", 16, 1000, 8},
		{"fixture count 1 caps at one", 16, 1, 1},
		{"fixture count 3 caps at three", 16, 3, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixtureWorkers(tt.budget, tt.fixtureCount)
			if got != tt.want {
				t.Fatalf("fixtureWorkers(%d, %d) = %d, want %d", tt.budget, tt.fixtureCount, got, tt.want)
			}
		})
	}
}

func TestSweepPinsSingleGOMAXPROCSInInnerEnv(t *testing.T) {
	tests := []struct {
		name string
		val  string
	}{
		{"present", "32"},
		{"present but empty", ""},
		{"non-numeric", "abc"},
		{"zero", "0"},
		{"negative", "-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GOMAXPROCS", tt.val)
			assertSweepPinsGOMAXPROCS(t)
		})
	}
	t.Run("unset", func(t *testing.T) {
		t.Setenv("GOMAXPROCS", "32")
		if err := os.Unsetenv("GOMAXPROCS"); err != nil {
			t.Fatal(err)
		}
		if _, ok := os.LookupEnv("GOMAXPROCS"); ok {
			t.Fatal("GOMAXPROCS still present in the environment after Unsetenv")
		}
		assertSweepPinsGOMAXPROCS(t)
	})
}

func assertSweepPinsGOMAXPROCS(t *testing.T) {
	t.Helper()
	outerVal, outerOK := os.LookupEnv("GOMAXPROCS")
	outerRuntime := runtime.GOMAXPROCS(0)

	root := t.TempDir()
	fixture := canaryFixture(root, mappedFamily(t), "only")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "target-only\n")

	var mu sync.Mutex
	var calls []RunCall
	runner := func(call RunCall) RunResult {
		mu.Lock()
		calls = append(calls, call)
		mu.Unlock()
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline\n"}
		}
		return RunResult{ExitCode: 1, Output: "target-only\n"}
	}

	if err := Sweep(root, runner); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}

	if gotVal, gotOK := os.LookupEnv("GOMAXPROCS"); gotOK != outerOK || gotVal != outerVal {
		t.Fatalf("outer GOMAXPROCS env = (%q, %v), want (%q, %v) unchanged by Sweep", gotVal, gotOK, outerVal, outerOK)
	}
	if got := runtime.GOMAXPROCS(0); got != outerRuntime {
		t.Fatalf("outer runtime.GOMAXPROCS(0) = %d, want %d unchanged by Sweep", got, outerRuntime)
	}

	if len(calls) != 2 {
		t.Fatalf("recorded %d calls, want 2 (baseline + one fixture)", len(calls))
	}

	baselineCall, fixtureCall := calls[0], calls[1]
	if baselineCall.FixtureDir != "" {
		t.Fatalf("calls[0] is not the baseline call")
	}
	if fixtureCall.FixtureDir == "" {
		t.Fatalf("calls[1] is not a fixture call")
	}

	assertSinglePinnedGOMAXPROCS(t, "baseline", baselineCall.Env)
	assertSinglePinnedGOMAXPROCS(t, "fixture", fixtureCall.Env)

	if !slices.Contains(fixtureCall.Env, PhaseEnv+"=conformance") {
		t.Fatalf("fixture call Env = %v, want a %s entry", fixtureCall.Env, PhaseEnv)
	}
}

func assertSinglePinnedGOMAXPROCS(t *testing.T, label string, env []string) {
	t.Helper()
	want := fmt.Sprintf("GOMAXPROCS=%d", bounds.CanaryInnerWidth)
	count := 0
	for _, kv := range env {
		if strings.HasPrefix(kv, "GOMAXPROCS=") {
			count++
			if kv != want {
				t.Fatalf("%s call GOMAXPROCS entry = %q, want %q", label, kv, want)
			}
		}
	}
	if count != 1 {
		t.Fatalf("%s call has %d GOMAXPROCS= entries, want exactly 1", label, count)
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
