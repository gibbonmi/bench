package worktree

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// inheritedRunBinary captures the gate-supplied BENCH_RUN_BINARY before any test can
// mutate the process environment. Package variables initialize before TestMain and
// before every t.Setenv, so this is the operator's selection, never a fixture's.
var inheritedRunBinary, inheritedRunBinaryPresent = os.LookupEnv(runbinary.Env)

// testRunSelector owns the one Bench executable a worktree test run hands to every
// real-binary journey. It reuses an inherited selection when one is present, or builds
// and seals exactly one direct-run executable through internal/runbinary. A present
// invalid selection refuses; it is never permission to build a private second binary.
type testRunSelector struct {
	inherited func() (string, bool)
	build     runbinary.Builder
	verify    runbinary.Verifier
	tempRoot  string

	once      sync.Once
	selection *runbinary.Selection
	err       error

	// journeys is the identity log: one "journey path" line per successful selection
	// consumer. A refusal appends nothing, so its length is the journey-start counter.
	// Parallel journeys share one selector, so the mutex guards every append and every
	// read; journeyLog is the only reader.
	journeysMu sync.Mutex
	journeys   []string
}

func (s *testRunSelector) selectFor(journey string) (string, error) {
	s.once.Do(func() { s.selection, s.err = s.selectOnce() })
	if s.err != nil {
		return "", s.err
	}
	s.journeysMu.Lock()
	s.journeys = append(s.journeys, journey+" "+s.selection.Path)
	s.journeysMu.Unlock()
	return s.selection.Path, nil
}

// journeyLog returns a copy of the identity log.
func (s *testRunSelector) journeyLog() []string {
	s.journeysMu.Lock()
	defer s.journeysMu.Unlock()
	return append([]string(nil), s.journeys...)
}

func (s *testRunSelector) selectOnce() (*runbinary.Selection, error) {
	source, err := benchSourceRoot()
	if err != nil {
		return nil, err
	}
	factory := runbinary.Factory{TempRoot: s.tempRoot, Build: s.build, Verify: s.verify}
	if raw, present := s.inherited(); present {
		return factory.Inherit(source, raw)
	}
	return factory.Own(context.Background(), source)
}

func (s *testRunSelector) close() {
	if s.selection != nil {
		_ = s.selection.Close()
	}
}

func benchSourceRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("resolve Bench source root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// packageRunBinary is the package-wide owner. Every journey that starts a real Bench
// child obtains its executable here; no journey builds its own.
var packageRunBinary = &testRunSelector{
	inherited: func() (string, bool) { return inheritedRunBinary, inheritedRunBinaryPresent },
}

func testRunBinary(t *testing.T) string {
	t.Helper()
	path, err := packageRunBinary.selectFor(t.Name())
	if err != nil {
		t.Fatalf("selected Bench test-run executable: %v", err)
	}
	return path
}

// TestDirectRunBuildsAndSealsOneExecutableForAllJourneys is SB1. A direct run — no
// inherited selection — must invoke the builder exactly once and hand every journey the
// same identity. A per-journey builder increments more than once, or yields a second
// path, and this test goes red.
func TestDirectRunBuildsAndSealsOneExecutableForAllJourneys(t *testing.T) {
	t.Parallel()
	builds := 0
	selector := &testRunSelector{
		inherited: func() (string, bool) { return "", false },
		tempRoot:  t.TempDir(),
		build: func(_ context.Context, _, output string) error {
			builds++
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		verify: func(string, string) error { return nil },
	}
	t.Cleanup(selector.close)
	paths := map[string]bool{}
	for _, journey := range []string{"journey-a", "journey-b", "journey-c"} {
		path, err := selector.selectFor(journey)
		requireTest(t, err == nil, "selectFor(%s): %v", journey, err)
		paths[path] = true
	}
	requireTest(t, builds == 1, "builder ran %d times across three journeys, want exactly 1", builds)
	requireTest(t, len(paths) == 1, "journeys observed %d executable identities %v, want 1", len(paths), paths)
	log := selector.journeyLog()
	requireTest(t, len(log) == 3, "identity log = %q, want three journey entries", log)
	for _, entry := range log {
		requireTest(t, strings.HasSuffix(entry, " "+selector.selection.Path), "identity log entry %q does not carry the selected path %q", entry, selector.selection.Path)
	}
}

// TestPackageSelectionIsSealedAndStable pins the real package owner: repeated
// consumption returns one identity, and the selected executable carries an adjacent
// freshness seal, whether inherited from the gate or built by the direct-run owner.
func TestPackageSelectionIsSealedAndStable(t *testing.T) {
	t.Parallel()
	first := testRunBinary(t)
	second := testRunBinary(t)
	requireTest(t, first == second, "package selection = %q then %q, want one identity", first, second)
	info, err := os.Stat(first + ".seal")
	requireTest(t, err == nil && info.Mode().IsRegular(), "selected executable seal %q = %v, %v, want a regular sidecar", first+".seal", info, err)
}

// TestInheritedSelectionReachesJourneysUnchangedWithZeroBuilds is SB2. A gate-supplied
// BENCH_RUN_BINARY must reach multiple journeys as the exact inherited path, with the
// private builder never invoked. A private fallback build changes the counter even when
// every journey stays green.
func TestInheritedSelectionReachesJourneysUnchangedWithZeroBuilds(t *testing.T) {
	t.Parallel()
	sealed := testRunBinary(t)
	builds := 0
	selector := &testRunSelector{
		inherited: func() (string, bool) { return sealed, true },
		tempRoot:  t.TempDir(),
		build: func(_ context.Context, _, output string) error {
			builds++
			return os.WriteFile(output, []byte("private"), 0o755)
		},
	}
	t.Cleanup(selector.close)
	for _, journey := range []string{"journey-a", "journey-b"} {
		path, err := selector.selectFor(journey)
		requireTest(t, err == nil, "selectFor(%s): %v", journey, err)
		requireTest(t, path == sealed, "journey %s received %q, want the inherited %q unchanged", journey, path, sealed)
	}
	requireTest(t, builds == 0, "inherited selection caused %d private builds, want 0", builds)
}

// TestInvalidSelectionRefusesBeforeAnyJourney is SB3. A missing, invalid, stale, or
// seal-less inherited executable refuses before any journey starts: the journey-start
// counter stays zero and the private builder stays cold.
func TestInvalidSelectionRefusesBeforeAnyJourney(t *testing.T) {
	t.Parallel()
	sealed := testRunBinary(t)
	sealedBytes, err := os.ReadFile(sealed)
	requireTest(t, err == nil, "read sealed executable: %v", err)
	sealBytes, err := os.ReadFile(sealed + ".seal")
	requireTest(t, err == nil, "read executable seal: %v", err)

	stale := filepath.Join(t.TempDir(), "bench")
	mustWrite(t, stale, append(append([]byte{}, sealedBytes...), '\n'), 0o755)
	mustWrite(t, stale+".seal", sealBytes, 0o644)

	sealless := filepath.Join(t.TempDir(), "bench")
	mustWrite(t, sealless, sealedBytes, 0o755)

	for _, tc := range []struct {
		name  string
		value string
	}{
		{"missing", filepath.Join(t.TempDir(), "absent")},
		{"invalid", "dist/bench"},
		{"stale", stale},
		{"seal-less", sealless},
	} {
		t.Run(tc.name, func(t *testing.T) {
			builds := 0
			selector := &testRunSelector{
				inherited: func() (string, bool) { return tc.value, true },
				tempRoot:  t.TempDir(),
				build: func(_ context.Context, _, output string) error {
					builds++
					return os.WriteFile(output, []byte("fallback"), 0o755)
				},
			}
			t.Cleanup(selector.close)
			_, err := selector.selectFor("journey-a")
			requireTest(t, err != nil, "selectFor with %s selection %q = nil, want refusal", tc.name, tc.value)
			requireTest(t, builds == 0, "%s selection caused %d fallback builds, want 0", tc.name, builds)
			log := selector.journeyLog()
			requireTest(t, len(log) == 0, "%s selection started %d journeys %q, want 0", tc.name, len(log), log)
		})
	}
}

// TestParallelJourneysRecordEverySelection is WF06. Two parallel journeys share one
// selector: the sync.Once still yields one executable, and the identity log keeps one
// line per journey. An unguarded append drops a line or races.
func TestParallelJourneysRecordEverySelection(t *testing.T) {
	t.Parallel()
	selector := &testRunSelector{
		inherited: func() (string, bool) { return "", false },
		tempRoot:  t.TempDir(),
		build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		verify: func(string, string) error { return nil },
	}
	t.Cleanup(selector.close)
	names := []string{"parallel-journey-a", "parallel-journey-b"}
	var wait sync.WaitGroup
	paths := make([]string, len(names))
	errs := make([]error, len(names))
	for i, journey := range names {
		wait.Add(1)
		go func() {
			defer wait.Done()
			paths[i], errs[i] = selector.selectFor(journey)
		}()
	}
	wait.Wait()
	for i, journey := range names {
		requireTest(t, errs[i] == nil, "selectFor(%s): %v", journey, errs[i])
		requireTest(t, paths[i] == paths[0], "journey %s received %q, want the one identity %q", journey, paths[i], paths[0])
	}
	log := selector.journeyLog()
	requireTest(t, len(log) == len(names), "identity log = %q, want one line per parallel journey (%d)", log, len(names))
	for _, journey := range names {
		found := false
		for _, entry := range log {
			if strings.HasPrefix(entry, journey+" ") {
				found = true
			}
		}
		requireTest(t, found, "identity log %q lacks the line for journey %s", log, journey)
	}
}
