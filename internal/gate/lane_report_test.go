package gate

// The lane's short pass: a pass keeps each check's report in the lane log and prints no
// check line, and each run starts the log again.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// TestLanePassPrintsNoCheckReport is BC49, BC50, and BC52. Each run is a pass over a
// prose check and a docs-currency-workflow check. The response holds no line of either
// check, the lane log holds the run's lines, and the second run leaves no line of the
// first run in the log.
func TestLanePassPrintsNoCheckReport(t *testing.T) {
	root := outcomeFixture(t)
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
	gitdir, err := benchgit.AdminDir(root)
	if err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(gitdir, "bench-lane.log")
	runs := []struct{ marker, stale string }{{marker: "first run"}, {marker: "second run", stale: "first run"}}
	for _, run := range runs {
		var stdout, stderr bytes.Buffer
		result, err := RunLane(context.Background(), LaneRequest{
			Root: root, Tree: tree, Stdout: &stdout, Stderr: &stderr,
			Checks: []Phase{
				{Name: "prose", Argv: []string{"sh", "-c", "echo prose " + run.marker}},
				{Name: "docs-currency-workflow", Argv: []string{"sh", "-c", "echo docs " + run.marker + " >&2"}},
			},
		})
		if err != nil || !result.Passed() {
			t.Fatalf("%s: lane = %+v, %v; want a pass", run.marker, result, err)
		}
		for _, line := range strings.Split(stdout.String()+stderr.String(), "\n") {
			if strings.HasPrefix(line, "[prose]") || strings.HasPrefix(line, "[docs-currency-workflow]") {
				t.Errorf("%s: response line %q, want no check line on a pass", run.marker, line)
			}
		}
		log, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("%s: lane log: %v", run.marker, err)
		}
		for _, want := range []string{"[prose] prose " + run.marker + "\n", "[docs-currency-workflow] docs " + run.marker + "\n"} {
			if !strings.Contains(string(log), want) {
				t.Errorf("%s: lane log = %q, want %q", run.marker, log, want)
			}
		}
		if run.stale != "" && strings.Contains(string(log), run.stale) {
			t.Errorf("%s: lane log = %q, want no line of the %s", run.marker, log, run.stale)
		}
	}
}
