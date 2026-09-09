package conformance

import (
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

const preparedBuildGuidanceDiagnosticPrefix = "prepared build guidance: "

func preparedBuildGuidanceAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(preparedBuildGuidanceDiagnosticPrefix)
}

// TestPreparedBuildGuidanceAnchorsBiteIndependently grades the build phase's prepared-charge
// authority. Every row is a Require row, so each one bites when its clause leaves the file.
func TestPreparedBuildGuidanceAnchorsBiteIndependently(t *testing.T) {
	buildAnchors := preparedBuildGuidanceAnchors()
	if got, want := len(buildAnchors), 4; got != want {
		t.Fatalf("prepared-build guidance anchor count = %d, want %d", got, want)
	}
	runAnchorBites(t, buildAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

const preparedReviewDispatchDiagnosticPrefix = "prepared review dispatch: "

func preparedReviewDispatchAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(preparedReviewDispatchDiagnosticPrefix)
}

// TestPreparedReviewGuidance grades the review phase's native-dispatch authority. A
// Require row must bite when its clause leaves the file. A Forbid row must bite when
// the retired same-family CLI route or inline-axis route returns to the file. The two
// directions together are what DP26 asks of the canonical review readers.
func TestPreparedReviewGuidance(t *testing.T) {
	dispatchAnchors := preparedReviewDispatchAnchors()
	if got, want := len(dispatchAnchors), 13; got != want {
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
	if required != 11 || forbidden != 2 {
		t.Fatalf("prepared-review dispatch kinds = %d Require and %d Forbid, want 11 and 2", required, forbidden)
	}
	runAnchorBites(t, dispatchAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

// TestPreparedReviewGuidanceHoldsOnTheLiveTree grades the shipped guidance files. The
// synthetic trees above prove that each row can bite. Only the live tree proves that
// the canonical readers carry the finished rules today.
func TestPreparedReviewGuidanceHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	diags := checkWorkflowAnchors(h.KitRoot)
	for _, anchor := range preparedReviewDispatchAnchors() {
		if containsDiagnostic(diags, anchor.Diagnostic) {
			t.Errorf("live guidance is not conformant: %s", anchor.Diagnostic)
		}
	}
}
