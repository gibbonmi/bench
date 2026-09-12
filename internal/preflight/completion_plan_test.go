package preflight

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// planFence renders the bench-completion-plan section the checkpoint reader
// parses, as one chunk over tickets. Every seeded spec carries it, so a green
// fixture states the plan the landing later requires. The section heading
// closes the ownership-fence section above it, so no payload byte reads as a
// fence token.
func planFence(tickets ...string) string {
	plan := reviewrecord.Plan{
		Version: 1,
		Chunks: []reviewrecord.PlannedChunk{{
			ID:           "c1",
			Tickets:      tickets,
			Verification: []reviewrecord.Requirement{{ID: "tests", Command: "go test ./..."}},
		}},
		FinalVerification: []reviewrecord.Requirement{{ID: "acceptance", Command: "go test ./..."}},
	}
	data, err := json.Marshal(plan)
	if err != nil {
		panic(err)
	}
	return "\n## Completion plan\n\n```bench-completion-plan\n" + string(data) + "\n```\n"
}

// specWithoutPlan is the seeded spec body with its completion plan removed, the
// state every spec written before the checkpoint reader is in.
func specWithoutPlan(slug string) string {
	return strings.TrimSuffix(specBody(slug), planFence("one.md"))
}

// replanSpec rewrites the seeded spec so its completion plan names tickets, the
// graph the fixture writes. A plan that omits a blocker names a dependency it
// never orders, so a fixture that grows its graph states that graph here.
func replanSpec(t *testing.T, slug string, tickets ...string) {
	t.Helper()
	mustWriteFile(t, "specs/"+slug+"/spec.md", specWithoutPlan(slug)+planFence(tickets...))
}

// planRow is the rendered completion-plan row of one verdict table.
func planRow(t *testing.T, out string) string {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "completion-plan,") {
			return strings.TrimSpace(line)
		}
	}
	t.Fatalf("output carries no completion-plan row:\n%s", out)
	return ""
}

// TestCommandBuildRedsAbsentCompletionPlan is the regression at the highest
// observable seam. A staged spec with no completion plan cannot pass the
// landing's checkpoint, so the build refuses before the work starts rather
// than at the checkpoint after it.
func TestCommandBuildRedsAbsentCompletionPlan(t *testing.T) {
	_, slug := seedConformant(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specWithoutPlan(slug))
	runGit(t, "commit", "-q", "-a", "-m", "drop the completion plan")

	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	row := planRow(t, out)
	if !strings.HasPrefix(row, "completion-plan,red,") || !strings.Contains(row, "bench-completion-plan") {
		t.Errorf("completion-plan row = %q, want a red naming the fence", row)
	}
}

// TestCommandBuildGreenCompletionPlan is that regression's sibling: the same
// tree with the plan in place answers green and exits 0.
func TestCommandBuildGreenCompletionPlan(t *testing.T) {
	_, slug := seedConformant(t)

	out, code := Command([]string{"build", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	if row := planRow(t, out); row != "completion-plan,green,\"\",\"\"" {
		t.Errorf("completion-plan row = %q, want green", row)
	}
}

// TestCommandReviewRedsAbsentCompletionPlan is the review-mode counterpart. The
// review phase reads the same plan the checkpoint does, so the missing fence
// reds there too.
func TestCommandReviewRedsAbsentCompletionPlan(t *testing.T) {
	_, slug := seedConformant(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specWithoutPlan(slug))
	runGit(t, "commit", "-q", "-a", "-m", "drop the completion plan")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	row := planRow(t, out)
	if !strings.HasPrefix(row, "completion-plan,red,") || !strings.Contains(row, "bench-completion-plan") {
		t.Errorf("completion-plan row = %q, want a red naming the fence", row)
	}
}

// TestCommandReviewGreenCompletionPlan is the review-mode green sibling.
func TestCommandReviewGreenCompletionPlan(t *testing.T) {
	_, slug := seedConformant(t)

	out, code := Command([]string{"review", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	if row := planRow(t, out); row != "completion-plan,green,\"\",\"\"" {
		t.Errorf("completion-plan row = %q, want green", row)
	}
}

// TestDecideCompletionPlanRow grades the row against the gathered facts alone:
// green with no reader error, red with one, and not-applicable in a build with
// no tickets/ directory, the same gate every ticket-reading row shares.
func TestDecideCompletionPlanRow(t *testing.T) {
	f := baseFacts()
	f.SourceTip = "cafe"
	f.CompletionPlanDigest = "d1"
	if c, _ := checkRow(Decide(f), "completion-plan"); c.Verdict != verdictGreen || c.Detail != "" {
		t.Errorf("completion-plan row over a parsed plan = %+v, want green", c)
	}

	f.CompletionPlanDigest = ""
	f.CompletionPlanError = "missing or invalid completion plan: no record"
	c, ok := checkRow(Decide(f), "completion-plan")
	want := "spec carries no valid bench-completion-plan fence at cafe: " +
		"missing or invalid completion plan: no record; see .bench/BENCH-reference.md, bench gate --checkpoint"
	if !ok || c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("completion-plan row = %+v, want a red detailed %q", c, want)
	}
	// WF35: no command answers an absent plan, so the row states its detail
	// and stops.
	if c.Next != "" {
		t.Errorf("completion-plan red carries Next = %q, want empty", c.Next)
	}

	f.Mode = modeBuild
	f.TicketsDirExists = false
	if c, _ := checkRow(Decide(f), "completion-plan"); c.Verdict != verdictNA {
		t.Errorf("completion-plan row in a fresh build = %+v, want not-applicable", c)
	}
}
