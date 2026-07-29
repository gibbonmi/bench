package contract

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/testrepo"
)

// ReleaseTarget is one entry of a staged release plan's target matrix. Its fields carry
// the plan's own tags, so a narrowed matrix re-encodes to the same bytes the plan shipped.
type ReleaseTarget struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	GOOS   string `json:"goos"`
	GOArch string `json:"goarch"`
	Runner string `json:"runner"`
}

// ReleasePlanTargets is the matrix NarrowReleasePlan hands the caller choosing the
// narrowed breadth: All is the staged plan's full target list, Host holds the first
// entry matching the GOOS/GOArch this toolchain builds for, empty when the plan has none.
type ReleasePlanTargets struct {
	All    []ReleaseTarget
	Host   []ReleaseTarget
	GOOS   string
	GOArch string
}

// NarrowReleasePlan commits root's working tree to a throwaway origin, clones it, and
// rewrites the clone's staged scripts/release-plan.json to the targets narrow returns,
// committing the rewrite. It returns the clone, ready to be handed to a build script as
// its source root.
//
// narrow chooses the breadth, so a caller can layer extra targets onto the host one or
// skip the test when the plan cannot serve it; returning an empty matrix is a caller bug
// rather than a narrowing, because a plan with no targets builds nothing and would leave
// every count-derived assertion trivially satisfied.
func NarrowReleasePlan(t testing.TB, root string, narrow func(ReleasePlanTargets) []ReleaseTarget) string {
	t.Helper()
	return NarrowReleasePlanIn(t, t.TempDir(), root, narrow)
}

// NarrowReleasePlanIn stages the narrowed clone under directory, whose lifetime
// remains owned by the caller.
func NarrowReleasePlanIn(t testing.TB, directory, root string, narrow func(ReleasePlanTargets) []ReleaseTarget) string {
	t.Helper()
	origin := filepath.Join(directory, "committed origin [*]")
	if err := testrepo.CommitWorkingTree(root, origin); err != nil {
		t.Fatal(err)
	}
	clone := filepath.Join(directory, "fresh source clone [*]")
	if output, err := exec.Command("git", "clone", "-q", origin, clone).CombinedOutput(); err != nil {
		t.Fatalf("clone committed source: %v\n%s", err, output)
	}
	planPath := filepath.Join(clone, "scripts", "release-plan.json")
	data, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatalf("read staged release plan: %v", err)
	}
	var plan map[string]json.RawMessage
	if err := json.Unmarshal(data, &plan); err != nil {
		t.Fatalf("parse staged release plan: %v", err)
	}
	matrix := ReleasePlanTargets{}
	if err := json.Unmarshal(plan["targets"], &matrix.All); err != nil {
		t.Fatalf("parse staged release targets: %v", err)
	}
	matrix.GOOS, matrix.GOArch = GoEnvPair(t, "GOOS", "GOARCH")
	for _, target := range matrix.All {
		if target.GOOS == matrix.GOOS && target.GOArch == matrix.GOArch {
			matrix.Host = []ReleaseTarget{target}
			break
		}
	}
	selected := narrow(matrix)
	if len(selected) == 0 {
		t.Fatalf("narrowed release plan for host %s/%s selected no targets", matrix.GOOS, matrix.GOArch)
	}
	plan["targets"], err = json.Marshal(selected)
	if err != nil {
		t.Fatalf("encode staged release targets: %v", err)
	}
	data, err = json.MarshalIndent(plan, "", "  ")
	if err != nil {
		t.Fatalf("encode staged release plan: %v", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(planPath, data, 0o644); err != nil {
		t.Fatalf("write staged release plan: %v", err)
	}
	commit := exec.Command("git", "-C", clone, "-c", "user.email=artifact-test@example.invalid", "-c", "user.name=artifact-test", "commit", "-qam", "stage artifact test release plan")
	if output, err := commit.CombinedOutput(); err != nil {
		t.Fatalf("commit staged release plan: %v\n%s", err, output)
	}
	return clone
}
