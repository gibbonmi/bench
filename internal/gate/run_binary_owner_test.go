package gate

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestGateOwnerBuildsOncePropagatesAndCleansEveryVerdict(t *testing.T) {
	for _, test := range []struct {
		name string
		exit string
		want int
	}{
		{name: "green", exit: "0", want: 0},
		{name: "red", exit: "7", want: 7},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("BENCH_KIT", "")
			root := gateTestRepo(t, "#!/usr/bin/env bash\nprintf '%s\\n' \"$BENCH_RUN_BINARY\" > .git/selected-path\nexit "+test.exit+"\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
			builds := 0
			var selectedDir string
			factory := runbinary.Factory{
				TempRoot: t.TempDir(),
				Build: func(_ context.Context, source, output string) error {
					builds++
					if source != root {
						t.Fatalf("owner source = %q, want %q", source, root)
					}
					selectedDir = filepath.Dir(output)
					return os.WriteFile(output, []byte("selected"), 0o755)
				},
				Verify: func(string, string) error { return nil },
			}
			evaluation := newWorkingTreeEvaluationAtKit(root, root)
			result := executeSubjectWithRunBinary(context.Background(), root, root, io.Discard, io.Discard, productionGateEngine{}, nil, forceRun, evaluation, factory.Own)
			if result.ActionExit != test.want || result.GateExit != test.want {
				t.Fatalf("result = %+v, want exit %d", result, test.want)
			}
			if builds != 1 {
				t.Fatalf("builds = %d, want one", builds)
			}
			selected := strings.TrimSpace(string(mustRead(t, filepath.Join(root, ".git", "selected-path"))))
			if filepath.Dir(selected) != selectedDir {
				t.Fatalf("gate selected %q, want path in %q", selected, selectedDir)
			}
			if _, err := os.Stat(selectedDir); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("run directory after %s = %v, want removed", test.name, err)
			}
		})
	}
}

func TestGateOwnerDrainsDescendantsBeforeRemovingSelection(t *testing.T) {
	t.Setenv("BENCH_KIT", "")
	root := gateTestRepo(t, `#!/usr/bin/env bash
(
  trap '' INT TERM
  exec >/dev/null 2>&1
  while test -e "$BENCH_RUN_BINARY"; do sleep .01; done
  printf premature > .git/selection-removed-first
) &
echo $! > .git/selected-child
exit 0
`, `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	var selectedDir string
	factory := runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			selectedDir = filepath.Dir(output)
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
	evaluation := newWorkingTreeEvaluationAtKit(root, root)
	ctx := withProcessGroupCancelGrace(context.Background(), fastProcessGroupCancelGrace)
	var stderr bytes.Buffer
	result := executeSubjectWithRunBinary(ctx, root, root, io.Discard, &stderr, productionGateEngine{}, nil, forceRun, evaluation, factory.Own)
	if result.ActionExit != 0 {
		t.Fatalf("result = %+v, want green; stderr=%q", result, stderr.String())
	}
	child := waitForPIDFile(t, filepath.Join(root, ".git", "selected-child"))
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	waitForProcessExit(t, child)
	if _, err := os.Stat(filepath.Join(root, ".git", "selection-removed-first")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("descendant observed premature selection removal: %v", err)
	}
	if _, err := os.Stat(selectedDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("run directory after descendant drain = %v, want removed", err)
	}
}
