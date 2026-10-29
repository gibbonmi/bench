package gate

// The lane's short pass: a pass keeps each check's report in the lane log and prints no
// check line, and each run starts the log again.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestLanePassPrintsNoCheckReport is BC49, BC50, and BC52. Each run is a pass over a
// prose check and a docs-currency-workflow check. The response holds no line of either
// check, and the lane log holds only the run's own lines. The first run writes more bytes
// than the second, so a log that the second run does not truncate keeps a tail of the first.
func TestLanePassPrintsNoCheckReport(t *testing.T) {
	root := outcomeFixture(t)
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	gitdir, err := benchgit.AdminDir(root)
	if err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(gitdir, "bench-lane.log")
	for _, marker := range []string{"first run, which writes the longer lines", "second run"} {
		var stdout, stderr bytes.Buffer
		result, err := RunLane(context.Background(), LaneRequest{
			Root: root, Tree: tree, Stdout: &stdout, Stderr: &stderr,
			Checks: []Phase{
				{Name: "prose", Argv: []string{"sh", "-c", "echo prose " + marker}},
				{Name: "docs-currency-workflow", Argv: []string{"sh", "-c", "echo docs " + marker + " >&2"}},
			},
		})
		if err != nil || !result.Passed() {
			t.Fatalf("%s: lane = %+v, %v; want a pass", marker, result, err)
		}
		for _, line := range strings.Split(stdout.String()+stderr.String(), "\n") {
			if strings.HasPrefix(line, "[prose]") || strings.HasPrefix(line, "[docs-currency-workflow]") {
				t.Errorf("%s: response line %q, want no check line on a pass", marker, line)
			}
		}
		log, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("%s: lane log: %v", marker, err)
		}
		// The two checks run concurrently, so the log's line order is not pinned.
		want := []string{"[docs-currency-workflow] docs " + marker + "\n", "[prose] prose " + marker + "\n"}
		if got := strings.SplitAfter(string(log), "\n"); !slices.Equal(slices.Sorted(slices.Values(got)), append([]string{""}, want...)) {
			t.Errorf("%s: lane log = %q, want only the lines %q", marker, log, want)
		}
	}
}
