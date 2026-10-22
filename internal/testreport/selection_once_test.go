package testreport

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// inheritedRunBinary captures the operator's BENCH_RUN_BINARY before any test clears it.
// Package variables initialize before TestMain and before every t.Setenv.
var inheritedRunBinary, inheritedRunBinaryPresent = os.LookupEnv(runbinary.Env)

// ownedRunBinary is the one run binary that the kit-copy tests of a package run share. An
// inherited selection wins, and a present invalid selection refuses instead of building.
// Each focused run validates the shared binary against its own kit copy's sources.
type ownedRunBinary struct {
	inherited func() (string, bool)
	factory   runbinary.Factory

	once      sync.Once
	selection *runbinary.Selection
	err       error
}

// path answers the shared run binary. The first call selects it for kit.
func (o *ownedRunBinary) path(kit string) (string, error) {
	o.once.Do(func() { o.selection, o.err = o.selectOnce(kit) })
	if o.err != nil {
		return "", o.err
	}
	return o.selection.Path, nil
}

func (o *ownedRunBinary) selectOnce(kit string) (*runbinary.Selection, error) {
	if raw, present := o.inherited(); present {
		return o.factory.Inherit(kit, raw)
	}
	return o.factory.Own(context.Background(), kit)
}

func (o *ownedRunBinary) close() {
	if o.selection != nil {
		_ = o.selection.Close()
	}
}

var packageRunBinary = &ownedRunBinary{
	inherited: func() (string, bool) { return inheritedRunBinary, inheritedRunBinaryPresent },
}

func TestMain(m *testing.M) {
	code := m.Run()
	packageRunBinary.close()
	os.Exit(code)
}

func TestOwnedSelectionBuildsOnce(t *testing.T) {
	builds := 0
	previous := packageRunBinary
	packageRunBinary = &ownedRunBinary{
		inherited: func() (string, bool) { return "", false },
		factory: runbinary.Factory{
			TempRoot: t.TempDir(),
			Build: func(_ context.Context, _, output string) error {
				builds++
				return os.WriteFile(output, []byte("selected"), 0o755)
			},
			Verify: func(string, string) error { return nil },
		},
	}
	t.Cleanup(func() {
		packageRunBinary.close()
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
