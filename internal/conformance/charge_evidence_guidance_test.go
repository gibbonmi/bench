package conformance

import (
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

const boundedBuildActionDiagnosticPrefix = "bounded build action: "

func boundedBuildActionAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(boundedBuildActionDiagnosticPrefix)
}

// TestEvidenceBuildGuidance grades the prerequisites the build phase requires before it
// acts on prepared evidence. The rows are CE101, CE102, CE103, CE140, and CE141. Each row
// is a Require row, so the phase that drops one prerequisite bites on that row alone.
// Approval and the supplement keep their own rows, because a phase that reads verified
// evidence as permission satisfies neither.
func TestEvidenceBuildGuidance(t *testing.T) {
	actionAnchors := boundedBuildActionAnchors()
	if got, want := len(actionAnchors), 5; got != want {
		t.Fatalf("bounded build action anchor count = %d, want %d", got, want)
	}
	for _, want := range []string{
		"permits build action without verified delivery",
		"permits build action without available required context",
		"permits build action without a current binding",
		"permits build action without reviewer approval",
		"permits build action without the complete task supplement",
	} {
		found := false
		for _, anchor := range actionAnchors {
			if anchor.Kind == anchors.Require && containsDiagnostic([]string{anchor.Diagnostic}, want) {
				found = true
			}
		}
		if !found {
			t.Errorf("no Require row states %q", want)
		}
	}
	runAnchorBites(t, actionAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}
