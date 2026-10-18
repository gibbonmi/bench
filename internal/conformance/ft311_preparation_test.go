package conformance

import (
	"path"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

const preparedBuildGuidanceDiagnosticPrefix = "prepared build guidance: "

func preparedBuildGuidanceAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(preparedBuildGuidanceDiagnosticPrefix)
}

// TestPreparedBuildGuidanceAnchorsBiteIndependently grades the build phase's charge
// preparation command. Its one row is a Require row, so it bites when its clause leaves the file.
func TestPreparedBuildGuidanceAnchorsBiteIndependently(t *testing.T) {
	buildAnchors := preparedBuildGuidanceAnchors()
	if got, want := len(buildAnchors), 1; got != want {
		t.Fatalf("prepared-build guidance anchor count = %d, want %d", got, want)
	}
	runAnchorBites(t, buildAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

const preparedReviewDispatchDiagnosticPrefix = "prepared review dispatch: "

func preparedReviewDispatchAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(preparedReviewDispatchDiagnosticPrefix)
}

// TestPreparedReviewGuidance grades the review phase's prepared-evidence commands. A
// Require row must bite when its command line leaves the file. A Forbid row must bite when
// a retired retrieval, same-family CLI, or inline-axis route returns to the file.
func TestPreparedReviewGuidance(t *testing.T) {
	dispatchAnchors := preparedReviewDispatchAnchors()
	if got, want := len(dispatchAnchors), 8; got != want {
		t.Fatalf("prepared-review dispatch anchor count = %d, want %d", got, want)
	}
	var required, forbidden int
	for _, anchor := range dispatchAnchors {
		switch anchor.Kind {
		case anchors.Require:
			required++
		case anchors.Forbid:
			forbidden++
		default:
			t.Errorf("prepared-review dispatch anchor %q uses an unsupported kind", anchor.Diagnostic)
		}
	}
	if required != 2 || forbidden != 6 {
		t.Fatalf("prepared-review dispatch kinds = %d Require and %d Forbid, want 2 and 6", required, forbidden)
	}
	runAnchorBites(t, dispatchAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

// TestPreparedReviewGuidanceHoldsOnTheLiveTree grades the shipped guidance files. The
// synthetic trees above prove that each row can bite. Only the live tree proves that
// the canonical readers carry the finished rules today.
func TestPreparedReviewGuidanceHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	diags := checkWorkflowAnchors(h.KitRoot)
	requireConformantDiagnostics(t, diags, preparedReviewDispatchAnchors())
}

// nativeHandoffPrefix opens every native-handoff diagnostic.
const nativeHandoffPrefix = "native handoff: "

// TestEvidenceHandoffGuidance grades CE107. The capable-harness handoff carries the trusted
// expected identity and the exact retrieval action, so a handoff that carries only the
// originating checkout path bites here.
func TestEvidenceHandoffGuidance(t *testing.T) {
	reviewPhase := boundedActionFamilyNamed(t, "review")
	guidance := reviewPhase.guidance(t, reviewPhase.anchors())
	want := []wantedRow{
		{anchors.Require, nativeHandoffPrefix + path.Base(guidance) + " drops the trusted evidence identity or the exact retrieval command from the capable-harness handoff"},
	}
	family := anchorsWithDiagnosticPrefix(nativeHandoffPrefix)
	if got := len(family); got != len(want) {
		t.Errorf("native handoff anchor count = %d, want %d", got, len(want))
	}
	for _, expected := range want {
		requireRegisteredRow(t, family, expected.kind, guidance, expected.diagnostic)
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
	requireConformantLiveTree(t, family)
}

// unchangedRoutePrefix opens every unchanged-route diagnostic.
const unchangedRoutePrefix = "unchanged route: "

// TestEvidenceUnchangedRoutes grades CE112. It grades the unchanged-route row.
func TestEvidenceUnchangedRoutes(t *testing.T) {
	const writeSpec = ".agents/commands/bench-write-spec.md"
	want := []struct {
		row  wantedRow
		file string
	}{
		{wantedRow{anchors.Forbid, unchangedRoutePrefix + path.Base(writeSpec) + " adds charge retrieval to the write-spec phase"}, writeSpec},
	}
	family := anchorsWithDiagnosticPrefix(unchangedRoutePrefix)
	if got := len(family); got != len(want) {
		t.Errorf("unchanged route anchor count = %d, want %d", got, len(want))
	}
	for _, expected := range want {
		requireRegisteredRow(t, family, expected.row.kind, expected.file, expected.row.diagnostic)
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
	requireConformantLiveTree(t, family)
}

// requireRegisteredRow reports the one family row the diagnostic names, with its expected
// kind and its expected canonical reader. A row that moves to a second file, a row that
// changes kind, and a deleted row each report here.
func requireRegisteredRow(t *testing.T, family []anchors.Anchor, kind anchors.Kind, file, diagnostic string) {
	t.Helper()
	for _, anchor := range family {
		if anchor.Diagnostic != diagnostic {
			continue
		}
		if anchor.Kind != kind || anchor.File != file {
			t.Errorf("row %q is a %v row over %s, want a %v row over %s", diagnostic, anchor.Kind, anchor.File, kind, file)
		}
		return
	}
	t.Errorf("no registry row states %q", diagnostic)
}

// requireConformantLiveTree grades the shipped guidance files. The synthetic trees in
// runAnchorBites prove that each row can bite. Only the live tree proves that the canonical
// readers carry the finished rules today.
func requireConformantLiveTree(t *testing.T, family []anchors.Anchor) {
	t.Helper()
	requireConformantDiagnostics(t, checkWorkflowAnchors(NewHarness(t).KitRoot), family)
}

// requireConformantDiagnostics reports every family row that the collected diagnostics
// contradict. It is the one comparison the guidance owners share, so a caller that reads the
// kit root itself keeps that read and still grades its family here.
func requireConformantDiagnostics(t *testing.T, diags []string, family []anchors.Anchor) {
	t.Helper()
	for _, anchor := range family {
		if containsDiagnostic(diags, anchor.Diagnostic) {
			t.Errorf("live guidance is not conformant: %s", anchor.Diagnostic)
		}
	}
}
