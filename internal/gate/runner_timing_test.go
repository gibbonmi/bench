package gate

// The conformance phase's per-check timing print: where the runner reads it from,
// and the roots where it must stay silent.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestRunnerPrintsConformanceTiming(t *testing.T) {
	const timing = "01 load-validity-metadata 2ms\n02 package-core-guard 41s\n"

	seedGitDir := func(t *testing.T, root string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	seedStaleTiming := func(t *testing.T, root string) {
		t.Helper()
		seedGitDir(t, root)
		if err := os.WriteFile(registry.TimingPath(root), []byte(timing), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// A phase that writes the timing file itself, the way the real conformance driver
	// does. Seeding the file before the phase instead would prove nothing here: the
	// runner clears at the run boundary, so only what the phase writes can print.
	writesTiming := func(dest, verdict string) string {
		return "printf %s '" + timing + "' > " + dest + "; " + verdict
	}

	for _, mode := range []struct {
		name   string
		mode   phaseMode
		prefix string
	}{
		{name: "outer", mode: outerMode, prefix: "[conformance] "},
		{name: "inner", mode: innerMode},
	} {
		for _, verdict := range []struct {
			name   string
			script string
		}{
			{name: "green", script: "true"},
			{name: "red", script: "exit 1"},
		} {
			t.Run(mode.name+"/"+verdict.name, func(t *testing.T) {
				root := t.TempDir()
				seedGitDir(t, root)
				phases := []Phase{fakePhase(conformancePhaseName, writesTiming(registry.TimingPath(root), verdict.script))}

				var stdout, stderr bytes.Buffer
				runPhases(context.Background(), root, phasesForMode(phases, mode.mode), mode.mode, &stdout, &stderr)

				out := stdout.String()
				for _, line := range strings.Split(strings.TrimSpace(timing), "\n") {
					if !strings.Contains(out, mode.prefix+line) {
						t.Fatalf("output missing timing line %q:\n%s", mode.prefix+line, out)
					}
				}
			})
		}
	}

	t.Run("stale file from an earlier run", func(t *testing.T) {
		root := t.TempDir()
		seedStaleTiming(t, root)

		var stdout, stderr bytes.Buffer
		runPhases(context.Background(), root, []Phase{fakePhase(conformancePhaseName, "true")}, outerMode, &stdout, &stderr)

		if got := stdout.String() + stderr.String(); strings.Contains(got, "load-validity-metadata") {
			t.Fatalf("this run's conformance phase wrote no timing, so the file an earlier run left behind must not be printed as if it were current:\n%s", got)
		}
	})

	t.Run("no git dir under a git-bearing ancestor", func(t *testing.T) {
		ancestor := t.TempDir()
		seedGitDir(t, ancestor)
		root := filepath.Join(ancestor, "graded")
		if err := os.Mkdir(root, 0o755); err != nil {
			t.Fatal(err)
		}
		// The phase writes the ancestor's file during the run, so an implementation
		// that ascended would find fresh lines rather than lines the run boundary
		// had already cleared.
		phase := fakePhase(conformancePhaseName, writesTiming(registry.TimingPath(ancestor), "true"))

		var stdout, stderr bytes.Buffer
		runPhases(context.Background(), root, []Phase{phase}, outerMode, &stdout, &stderr)

		if got := stdout.String() + stderr.String(); strings.Contains(got, "load-validity-metadata") {
			t.Fatalf("the graded root has no git dir, so the ancestor's timing file must stay unread:\n%s", got)
		}
	})
}
