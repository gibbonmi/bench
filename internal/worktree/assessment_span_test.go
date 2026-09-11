package worktree

import (
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/poolkey"
	"path/filepath"
	"testing"
)

func TestAssessmentAssignmentSpan(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	assignment := "11111111111111111111111111111111"
	root := filepath.Join(t.TempDir(), poolkey.AssignmentSegment(assignment, assignment))
	_, finish := beginLandingSpan(home, root)
	finish(0, landingMeasures{assignment: assignment, subject: "published", pathCount: 3, counted: true})
	spans, err := otelrecord.ReadSpans(home, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 1 || spans[0].Attributes[otelrecord.AttrAssignmentID] != assignment || spans[0].Attributes[otelrecord.AttrSubjectID] != "published" || spans[0].Attributes[otelrecord.AttrMeasurePathCount] != "3" {
		t.Fatalf("assignment provenance missing: %+v", spans)
	}
}
