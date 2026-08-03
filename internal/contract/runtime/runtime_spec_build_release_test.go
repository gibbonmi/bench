package runtime

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

type runtimeReleaseState struct {
	Candidate    string                              `json:"candidate"`
	CandidateTip string                              `json:"candidate_tip"`
	Assignments  map[string]runtimeReleaseAssignment `json:"assignments"`
}

type runtimeReleaseAssignment struct {
	ID, Integrated           string
	CleanupPending, Released bool
}

func TestRuntimeSpecBuildIntegrateResumesExactLiveRelease(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	contract.RequireFreshBench(t)
	f := contract.NewFixture(t)
	f.WriteFile(".bench/gate.sh", "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(filepath.Join(f.Root, ".bench", "gate.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	f.WriteFile(".gitignore", ".bench-contract-env/\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/one.md", "# One\n\nOwnership fence: internal/demo\n\n- [ ] [R1] exact live release\n")
	f.CommitAll("staged spec")
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)
	assignment := assignRuntimeBuild(t, f, "one.md", "live-release-request")
	contract.WriteFileAbs(t, filepath.Join(assignment.Path, "internal", "demo", "change.go"), "package demo\n")
	requireDirtyAssignmentAtBase(t, f, assignment)
	receipt := runtimeCheckpointReceipt(t, f, assignment, []string{"internal/demo/change.go"})
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", assignment.ID, "--evidence", receipt).RequireExit(0)

	drift := filepath.Join(assignment.Path, "post-checkpoint.txt")
	contract.WriteFileAbs(t, drift, "drift\n")
	failed := f.Bench("spec", "build", "integrate", "demo", "--assignment", assignment.ID)
	failed.RequireExit(1)
	contract.RequireContains(t, failed.Stdout, "release integrated assignment")
	afterFailure := readRuntimeReleaseState(t, f)
	stored := releaseAssignmentByID(t, afterFailure, assignment.ID)
	if stored.Integrated == "" || stored.Integrated != afterFailure.CandidateTip || !stored.CleanupPending || stored.Released {
		t.Fatalf("pending release state = %#v candidate=%s", stored, afterFailure.CandidateTip)
	}
	commits := strings.TrimSpace(f.Git("rev-list", "--count", afterFailure.Candidate).Stdout)
	f.WriteFile("specs/demo/tickets/one.md", "ticket dirt after durable integration\n")
	if err := os.Remove(drift); err != nil {
		t.Fatal(err)
	}

	f.Bench("spec", "build", "integrate", "demo", "--assignment", assignment.ID).RequireExit(0)
	afterResume := readRuntimeReleaseState(t, f)
	stored = releaseAssignmentByID(t, afterResume, assignment.ID)
	if afterResume.CandidateTip != afterFailure.CandidateTip || stored.Integrated != afterFailure.CandidateTip || stored.CleanupPending || !stored.Released {
		t.Fatalf("resumed release state = %#v candidate before=%s after=%s", stored, afterFailure.CandidateTip, afterResume.CandidateTip)
	}
	if got := strings.TrimSpace(f.Git("rev-list", "--count", afterResume.Candidate).Stdout); got != commits {
		t.Fatalf("resume candidate commits = %s, want %s", got, commits)
	}
	if _, err := os.Stat(assignment.Path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("released assignment remains: %v", err)
	}
	if refs := strings.TrimSpace(f.Git("for-each-ref", "--format=%(refname)", "refs/bench/recovery/").Stdout); refs != "" {
		t.Fatalf("exact live release created recovery refs: %s", refs)
	}
}

func readRuntimeReleaseState(t *testing.T, f contract.Fixture) runtimeReleaseState {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(gitDir(t, f), "bench", "specbuild", "*.json"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("spec build state paths = %v, %v", paths, err)
	}
	var state runtimeReleaseState
	data, err := os.ReadFile(paths[0])
	if err != nil || json.Unmarshal(data, &state) != nil {
		t.Fatalf("read spec build release state: %v", err)
	}
	return state
}

func releaseAssignmentByID(t *testing.T, state runtimeReleaseState, id string) runtimeReleaseAssignment {
	t.Helper()
	for _, assignment := range state.Assignments {
		if assignment.ID == id {
			return assignment
		}
	}
	t.Fatalf("assignment %s missing from durable release state", id)
	return runtimeReleaseAssignment{}
}
