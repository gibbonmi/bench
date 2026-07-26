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

	seed := func(t *testing.T, root string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(registry.TimingPath(root), []byte(timing), 0o644); err != nil {
			t.Fatal(err)
		}
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
				seed(t, root)
				phases := []Phase{fakePhase(conformancePhaseName, verdict.script)}

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

	t.Run("no git dir under a git-bearing ancestor", func(t *testing.T) {
		ancestor := t.TempDir()
		seed(t, ancestor)
		root := filepath.Join(ancestor, "graded")
		if err := os.Mkdir(root, 0o755); err != nil {
			t.Fatal(err)
		}

		var stdout, stderr bytes.Buffer
		runPhases(context.Background(), root, []Phase{fakePhase(conformancePhaseName, "true")}, outerMode, &stdout, &stderr)

		if got := stdout.String() + stderr.String(); strings.Contains(got, "load-validity-metadata") {
			t.Fatalf("the graded root has no git dir, so the ancestor's timing file must stay unread:\n%s", got)
		}
	})
}
