package testreport

import (
	"context"
	"os"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// inheritedRunBinary captures the operator's BENCH_RUN_BINARY before any test clears it.
// Package variables initialize before TestMain and before every t.Setenv.
var inheritedRunBinary, inheritedRunBinaryPresent = os.LookupEnv(runbinary.Env)

// packageRunBinary is the one run binary that the kit-copy tests of a package run share.
// Each focused run validates the shared binary against its own kit copy's sources.
var packageRunBinary = &runbinary.Shared{Raw: inheritedRunBinary, Present: inheritedRunBinaryPresent}

func TestMain(m *testing.M) {
	code := m.Run()
	_ = packageRunBinary.Close()
	os.Exit(code)
}

func TestOwnedSelectionBuildsOnce(t *testing.T) {
	builds := 0
	previous := packageRunBinary
	packageRunBinary = &runbinary.Shared{
		Factory: runbinary.Factory{
			TempRoot: t.TempDir(),
			Build: func(_ context.Context, _, output string) error {
				builds++
				return os.WriteFile(output, []byte("selected"), 0o755)
			},
			Verify: func(string, string) error { return nil },
		},
	}
	t.Cleanup(func() {
		_ = packageRunBinary.Close()
		packageRunBinary = previous
	})
	selected := map[string]bool{}
	for _, name := range []string{"first kit copy", "second kit copy"} {
		t.Run(name, func(t *testing.T) {
			installOwnedSelection(t, t.TempDir())
			selected[os.Getenv(runbinary.Env)] = true
		})
	}
	if builds != 1 || len(selected) != 1 || selected[""] {
		t.Fatalf("two kit copies ran the builder %d times and selected %v; want one build and one selection", builds, selected)
	}
}
