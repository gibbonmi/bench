package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

// The out-of-fence repair route, end to end through the compiled public
// porcelain: a consumer delegate with an in-fence partial edit is blocked by a
// prerequisite defect in a disjoint fence, returns a debug receipt instead of
// editing outside its fence, and the orchestrator lands the repair and refreshes
// the blocked assignment — no synthesized commits, refs, worktrees, or replay.
func TestRuntimeSpecBuildOutOfFenceRepairRefreshesBlockedAssignment(t *testing.T) {
	f := setupRuntimeBuildGate(t, "#!/bin/sh\nexit 0\n")
	f.WriteFile("specs/demo/spec.md", "# demo\n\nStatus: staged\n")
	f.WriteFile("specs/demo/tickets/consumer.md", "# Consumer\n\nOwnership fence: internal/consumer\n\n- [ ] [R70] consumer adoption\n")
	f.CommitAll("staged spec")
	f.Bench("gate").RequireExit(0)
	f.Bench("spec", "build", "start", "demo").RequireExit(0)

	consumer := assignRuntimeBuild(t, f, "consumer.md", "consumer-request")
	partial := filepath.Join(consumer.Path, "internal", "consumer", "partial.go")
	partialBody := "package consumer\n\nvar partial = true\n"
	contract.WriteFileAbs(t, partial, partialBody)

	// The repair ticket lands mid-run; until the run recomposes onto the new
	// tip, assign refuses — the insertion that used to force whole-run abandon.
	f.WriteFile("specs/demo/tickets/repair.md", "# Repair\n\nOwnership fence: internal/repairfence\n\n- [ ] [R71] prerequisite repair\n")
	f.CommitAll("stage repair ticket")
	f.Bench("gate").RequireExit(0)
	refused := f.Bench("spec", "build", "assign", "demo", "--ticket", "repair.md", "--request", "repair-request")
	refused.RequireExit(1)
	contract.RequireContains(t, refused.Stdout, "recomposition")
	f.Bench("spec", "build", "promote", "demo").RequireExit(0)

	repair := assignRuntimeBuild(t, f, "repair.md", "repair-request")
	writeAssignmentChange(t, repair.Path, "internal/repairfence/fix.go", "package repairfence\n")
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", repair.ID, "--evidence", runtimeCheckpointReceipt(t, f, repair, []string{"internal/repairfence/fix.go"})).RequireExit(0)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", repair.ID).RequireExit(0)

	state := readRuntimeBuildState(t, f)
	receipt := map[string]any{
		"version": 1, "run": state.Run, "assignment": consumer.ID, "base": consumer.Base,
		"repro": map[string]any{"command": "go test ./internal/consumer", "exit": 1, "output_digest": runtimeDigest("red"), "produced": time.Now().UTC().Add(time.Second).Format(time.RFC3339Nano)},
		"cause": "prerequisite defect in internal/repairfence", "required_fence": []string{"internal/repairfence"},
		"dirty_paths": []string{"internal/consumer/partial.go"}, "resumable": true,
	}
	forged := map[string]any{}
	for k, v := range receipt {
		forged[k] = v
	}
	forged["base"] = state.CandidateTip
	forgedProbe := f.Bench("spec", "build", "assign", "demo", "--ticket", "consumer.md", "--request", "consumer-request", "--refresh", writeRuntimeJSON(t, forged))
	forgedProbe.RequireExit(1)
	contract.RequireContains(t, forgedProbe.Stdout, "invalid spec build debug receipt")

	refreshed := f.Bench("spec", "build", "assign", "demo", "--ticket", "consumer.md", "--request", "consumer-request", "--refresh", writeRuntimeJSON(t, receipt))
	refreshed.RequireExit(0)
	contract.RequireContains(t, refreshed.Stdout, consumer.ID)
	contract.RequireContains(t, refreshed.Stdout, state.CandidateTip)
	if head := strings.TrimSpace(contract.RunAt(t, f, consumer.Path, nil, "git", "rev-parse", "HEAD").Stdout); head != state.CandidateTip {
		t.Fatalf("refreshed worktree HEAD = %s, want repaired candidate %s", head, state.CandidateTip)
	}
	surviving, err := os.ReadFile(partial)
	if err != nil || string(surviving) != partialBody {
		t.Fatalf("in-fence bytes did not survive refresh: %q, %v", surviving, err)
	}

	// The resumed delegate completes against the repaired base and the ordinary
	// checkpoint, integrate, review, and promote path lands the composition.
	writeAssignmentChange(t, consumer.Path, "internal/consumer/final.go", "package consumer\n")
	resumed := runtimeAssignmentByID(t, f, consumer.ID)
	f.Bench("spec", "build", "checkpoint", "demo", "--assignment", resumed.ID, "--evidence", runtimeCheckpointReceipt(t, f, resumed, []string{"internal/consumer/final.go", "internal/consumer/partial.go"})).RequireExit(0)
	f.Bench("spec", "build", "integrate", "demo", "--assignment", resumed.ID).RequireExit(0)
	f.Bench("spec", "build", "review", "demo", "--evidence", runtimeReviewReceipt(t, f)).RequireExit(0)
	f.Bench("spec", "build", "promote", "demo").RequireExit(0)
	for _, path := range []string{"internal/repairfence/fix.go", "internal/consumer/partial.go", "internal/consumer/final.go"} {
		f.Git("cat-file", "-e", "HEAD:"+path).RequireExit(0)
	}
}

func runtimeAssignmentByID(t *testing.T, f contract.Fixture, id string) runtimeBuildAssignment {
	t.Helper()
	for _, assignment := range readRuntimeBuildState(t, f).Assignments {
		if assignment.ID == id {
			return assignment
		}
	}
	t.Fatalf("assignment %s missing from durable state", id)
	return runtimeBuildAssignment{}
}
