package canary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// TestSweepRefusesUnboundFamilyBeforeAnyRun grades the assertion's position as much as
// its verdict: an unbound family costs a full unscoped inner gate per fixture, so the
// sweep has to refuse before it spends any of them. The adopting-repo half is the
// scoping the assertion lives or dies by — a repo whose families the kit table will
// never carry must sweep exactly as it does today.
func TestSweepRefusesUnboundFamilyBeforeAnyRun(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, "unbound-family", "orphan"), "")

	t.Run("kit", func(t *testing.T) {
		t.Setenv("BENCH_KIT", root)
		calls, err := countedSweep(t, root)
		if err == nil {
			t.Fatal("SweepTier err = nil, want the unbound family refused")
		}
		want := `canary conformance family "unbound-family" is bound to no conformance check`
		if !strings.Contains(err.Error(), want) {
			t.Errorf("SweepTier err = %v, want a diagnostic containing %q", err, want)
		}
		if calls != 0 {
			t.Errorf("sweep ran %d inner gates before refusing, want none", calls)
		}
	})

	t.Run("adopting repo", func(t *testing.T) {
		t.Setenv("BENCH_KIT", t.TempDir())
		calls, err := countedSweep(t, root)
		if err != nil {
			t.Fatalf("SweepTier err = %v, want an adopting repo's unbound family swept unscoped", err)
		}
		if calls == 0 {
			t.Error("sweep ran no inner gates, want the fixture and its baseline graded")
		}
	})
}

// TestSweepBindingAssertionSkipsWithoutTheWrapperExport covers the shapes under which
// the kit's own assertion silently does not run. The wrapper exports BENCH_KIT as an
// absolute path; anything else is some other caller, so the sweep falls through to the
// adopting-repo behavior rather than refusing a tree whose families a kit-owned table
// will never carry. Every other test here sets the variable to a real kit, so this is
// the side no other assertion reaches.
func TestSweepBindingAssertionSkipsWithoutTheWrapperExport(t *testing.T) {
	root := t.TempDir()
	fixture(t, canaryFixture(root, "unbound-family", "orphan"), "")

	for name, kit := range map[string]string{
		"empty":    "",
		"relative": filepath.Base(root),
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("BENCH_KIT", kit)
			calls, err := countedSweep(t, root)
			if err != nil {
				t.Fatalf("SweepTier err = %v, want the unbound family swept rather than refused", err)
			}
			if calls == 0 {
				t.Error("sweep ran no inner gates, want the fixture and its baseline graded")
			}
		})
	}
}

// TestSweepBindingAssertionResolvesSymlinks pins the comparison both sides go through.
// bin/bench.sh derives BENCH_KIT with a physical cd while the sweep normalizes its root
// with filepath.Abs alone, so a raw string compare makes a symlinked kit checkout read
// as an adopting repo and skip its own assertion silently.
func TestSweepBindingAssertionResolvesSymlinks(t *testing.T) {
	physical, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	fixture(t, canaryFixture(physical, "unbound-family", "orphan"), "")
	link := filepath.Join(t.TempDir(), "kit-link")
	if err := os.Symlink(physical, link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
	t.Setenv("BENCH_KIT", physical)

	calls, err := countedSweep(t, link)
	if err == nil {
		t.Fatal("SweepTier err = nil, want the symlinked kit checkout to assert its own binding")
	}
	if calls != 0 {
		t.Errorf("sweep ran %d inner gates before refusing, want none", calls)
	}
}

// countedSweep sweeps root at the dev tier and reports how many inner gates a sweep that
// got past the assertion would have run.
func countedSweep(t *testing.T, root string) (int, error) {
	t.Helper()
	var mu sync.Mutex
	calls := 0
	err := SweepTier(root, registry.Dev, func(call RunCall) RunResult {
		mu.Lock()
		calls++
		mu.Unlock()
		if result, done := stubToolchain(call); done {
			return result
		}
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline noise\n"}
		}
		return RunResult{ExitCode: 1, Output: "target-" + filepath.Base(call.FixtureDir) + "\n"}
	})
	return calls, err
}
