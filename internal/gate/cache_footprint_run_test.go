package gate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gocache"
)

// C12: runPhases hands os.Environ() down, gocache.FromEnv reads the entry gocache.Apply
// wrote there, and the run records one cache.footprint event at that directory.
func TestGateRunnerRecordsOneCacheFootprintEventForItsProcessEnv(t *testing.T) {
	applyProcessBuildCache(t)
	root := newLoggingPruneRoot(t)
	var stdout, stderr bytes.Buffer

	ctx, finish := beginGateRunLog(context.Background(), root, &stderr, "dev")
	if strings.Contains(stderr.String(), "progress logging unavailable") {
		t.Fatalf("the run opened no record, so this asserts nothing: %q", stderr.String())
	}
	record := gateRunLogFrom(t, ctx).file.Name()
	code := runPhases(ctx, root, []Phase{{Name: "trivial", Argv: []string{"true"}}}, &stdout, &stderr)
	finish(Result{GateExit: code, ActionExit: code})

	if code != 0 {
		t.Fatalf("run exit = %d, want a green run; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	want, err := gocache.Dir(os.Environ())
	if err != nil {
		t.Fatal(err)
	}
	events := footprintEvents(t, record)
	if len(events) != 1 {
		t.Fatalf("run recorded %d cache.footprint events, want exactly 1", len(events))
	}
	if events[0].Path != want {
		t.Errorf("recorded path = %q, want the process environment's %q", events[0].Path, want)
	}
}

// applyProcessBuildCache puts the process environment in the shape a gate runs under: an
// absolute HOME, and the one GOCACHE entry gocache.Apply derives from it. A `go test`
// process inherits neither, so without this the runner would read an environment no gate
// ever runs with.
func applyProcessBuildCache(t *testing.T) {
	t.Helper()
	if !filepath.IsAbs(os.Getenv("HOME")) {
		t.Setenv("HOME", t.TempDir())
	}
	applied, err := gocache.Apply(os.Environ())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv(gocache.Env, gocache.FromEnv(applied))
}
