package roadmap

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// absentOutlook is the commitment source of a repository with no published policy.
func absentOutlook(string) commitment.Outlook {
	return commitment.Project(nil, intent.CommitmentState{})
}

// TestRoadmapRendersSuppliedOutlook holds both roadmap forms to the outlook that the
// adapter supplies. Each form renders the outlook row and the blocker row verbatim,
// beside the recommended sequence that stays the board's own text.
func TestRoadmapRendersSuppliedOutlook(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("# Roadmap\n\n**FT1 — A**\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. A — `/bench-write-spec`\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	supplied := func(got string) commitment.Outlook {
		if got != root {
			t.Errorf("outlook root = %q, want %q", got, root)
		}
		return commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Deliverable: "specs/b/spec.md", Blocked: []intent.OutcomeBlocker{{Outcome: "A", Reason: "vendor fix"}}, Waiting: []string{"C"}, Operation: "start", Command: "bench commitment start --outcome B"}
	}
	const want = "commitment_outlook[1]{state,active_milestone,next_outcome,deliverable,active,blocked,waiting,command}:\n  eligible,M1,B,specs/b/spec.md,\"\",A,C,bench commitment start --outcome B\n" +
		"commitment_blockers[1]{outcome,reason}:\n  A,vendor fix\n"
	for name, run := range map[string]func() (string, int){
		"board": func() (string, int) { return RoadmapCommand(nil, supplied) },
		"context": func() (string, int) {
			return ContextCommand([]string{"--context"}, func(string) GateCacheFact { return GateCacheFact{} }, supplied)
		},
	} {
		out, code := run()
		if code != 0 || !strings.Contains(out, want) {
			t.Fatalf("%s = (%d):\n%s\nwant the outlook blocks:\n%s", name, code, out, want)
		}
	}
}
