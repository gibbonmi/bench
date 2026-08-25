// The dashboard's gate row against the fast lane: a lane run is not a gate verdict.
package status

import (
	"context"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// OG20: a lane run writes its own record class and never the gate cache, so the gate row
// keeps reporting the last gate verdict. A dashboard that read the lane file would show a
// lane pass where the operator reads the oracle's truth.
func TestGateRowKeepsTheLastGateVerdictAfterALaneRun(t *testing.T) {
	root := initRepo(t)
	tree := treeOf(t, root, map[string]string{"f.txt": "x\n"})
	// The lane materializes the graded tree beside HEAD, so the fixture needs a commit.
	gitRun(t, root, "commit", "-q", "--allow-empty", "-m", "base")
	writeFullGateCache(t, root, tree, "green")
	before := GateVerdict(root)
	if !before.Present || before.Status != "green" {
		t.Fatalf("verdict = %#v, want the recorded green gate verdict", before)
	}

	result, err := gate.RunLane(context.Background(), gate.LaneRequest{
		Root: root, Tree: tree, Lane: "fixture",
		Checks: []gate.Phase{{Name: "unit", Argv: []string{"true"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed() {
		t.Fatalf("lane result = %#v, want a pass", result)
	}

	if after := GateVerdict(root); !reflect.DeepEqual(after, before) {
		t.Fatalf("verdict after the lane = %#v, want the gate verdict %#v", after, before)
	}
}
