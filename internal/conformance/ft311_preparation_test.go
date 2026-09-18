package conformance

import (
	"slices"
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
	requireConformantDiagnostics(t, diags, preparedReviewDispatchAnchors())
}

// consumerContextPrefix opens every consumer-context diagnostic. One prefix lets this owner
// enumerate the family without a second registry.
const consumerContextPrefix = "consumer context: "

// consumerContextPhase is one canonical phase reader that states the consumer-context rules
// over its own consumers. Both readers carry the same two rules, so the expectation is
// written once here and bound to a phase rather than copied per phase.
type consumerContextPhase struct {
	// name titles the phase's subtests.
	name string
	// file is the canonical reader whose rows the phase owns.
	file string
	// peer names the other consumer of the same evidence. Its receipt is the cheapest
	// substitute for a fresh consumer's own retrieval, so each phase names its own peer.
	peer string
}

// consumerContextPhases is the one phase inventory. A phase that migrates to the bounded
// evidence path joins this table, and the rules below then grade it.
func consumerContextPhases() []consumerContextPhase {
	return []consumerContextPhase{
		{name: "build", file: ".agents/commands/bench-implement-spec.md", peer: "consumer"},
		{name: "review", file: ".agents/commands/bench-review-implementation.md", peer: "axis"},
	}
}

// wantedRows states the phase's two rules apart from the registry: reuse verifies the new
// manifest, and a fresh consumer runs its own retrieval. A reworded or deleted row bites.
func (p consumerContextPhase) wantedRows() []string {
	return []string{
		consumerContextPrefix + p.file + " permits reuse without verified membership, role, and requiredness",
		consumerContextPrefix + p.file + " permits another " + p.peer + "'s receipt or a final cursor to replace fresh required context",
	}
}

// TestEvidenceConsumerGuidance grades CE104 and CE105. The phase names are compared in
// document order, so a dropped phase cannot leave the two rules vacuously green. Each rule
// takes one Require row over its own reader, each row bites alone, and the shipped guidance
// satisfies every row today.
func TestEvidenceConsumerGuidance(t *testing.T) {
	phases := consumerContextPhases()
	var names, want []string
	for _, phase := range phases {
		names = append(names, phase.name)
		want = append(want, phase.wantedRows()...)
	}
	if !slices.Equal(names, []string{"build", "review"}) {
		t.Errorf("consumer context phases = %v, want [build review]", names)
	}
	family := anchorsWithDiagnosticPrefix(consumerContextPrefix)
	if got := len(family); got != len(want) {
		t.Errorf("consumer context anchor count = %d, want %d", got, len(want))
	}
	for _, phase := range phases {
		for _, diagnostic := range phase.wantedRows() {
			requireRegisteredRow(t, family, anchors.Require, phase.file, diagnostic)
		}
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
	requireConformantLiveTree(t, family)
}

// nativeHandoffPrefix opens every native-handoff diagnostic.
const nativeHandoffPrefix = "native handoff: "

// TestEvidenceHandoffGuidance grades CE107. The capable-harness handoff carries the trusted
// expected identity and the exact retrieval action, so a handoff that carries only the
// originating checkout path bites here.
func TestEvidenceHandoffGuidance(t *testing.T) {
	const reviewPhase = ".agents/commands/bench-review-implementation.md"
	want := []string{
		nativeHandoffPrefix + reviewPhase + " drops the trusted evidence identity or the exact retrieval command from the capable-harness handoff",
	}
	family := anchorsWithDiagnosticPrefix(nativeHandoffPrefix)
	if got := len(family); got != len(want) {
		t.Errorf("native handoff anchor count = %d, want %d", got, len(want))
	}
	for _, diagnostic := range want {
		requireRegisteredRow(t, family, anchors.Require, reviewPhase, diagnostic)
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
	requireConformantLiveTree(t, family)
}

// unchangedRoutePrefix opens every unchanged-route diagnostic.
const unchangedRoutePrefix = "unchanged route: "

// retainedWriteSpecDiagnostics names the registry rows that already state write-spec's
// decision-source and author-fork contract. CE113 keeps that contract, so this owner cites
// those rows rather than restating their prose behind a second pin.
func retainedWriteSpecDiagnostics() []string {
	return []string{
		".agents/commands/bench-write-spec.md dropped the exactly-one Decision source contract",
		".agents/commands/bench-write-spec.md dropped the conversation fork for spec authoring and ticket slicing",
	}
}

// TestEvidenceUnchangedRoutes grades CE112 and CE113. The implementation phase keeps its
// own full-run control, which names no preflight command and is therefore no charge route.
// Write-spec keeps its decision-source and author-fork contract, and it gains no charge
// retrieval of its own.
func TestEvidenceUnchangedRoutes(t *testing.T) {
	const (
		buildPhase = ".agents/commands/bench-implement-spec.md"
		writeSpec  = ".agents/commands/bench-write-spec.md"
	)
	want := []struct {
		kind       anchors.Kind
		file       string
		diagnostic string
	}{
		{anchors.Require, buildPhase, unchangedRoutePrefix + buildPhase + " dropped the phase-level full-run control"},
		{anchors.Forbid, writeSpec, unchangedRoutePrefix + writeSpec + " adds charge retrieval to the write-spec phase"},
	}
	family := anchorsWithDiagnosticPrefix(unchangedRoutePrefix)
	if got := len(family); got != len(want) {
		t.Errorf("unchanged route anchor count = %d, want %d", got, len(want))
	}
	for _, expected := range want {
		requireRegisteredRow(t, family, expected.kind, expected.file, expected.diagnostic)
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
	requireConformantLiveTree(t, family)

	registered := anchors.Entries()
	for _, diagnostic := range retainedWriteSpecDiagnostics() {
		if !slices.ContainsFunc(registered, func(a anchors.Anchor) bool { return a.Diagnostic == diagnostic }) {
			t.Errorf("no registry row states the retained write-spec rule %q", diagnostic)
		}
	}
	requireConformantLiveTree(t, rowsWithDiagnostics(registered, retainedWriteSpecDiagnostics()))
}

// rowsWithDiagnostics selects the registered rows the diagnostics name. An absent
// diagnostic selects nothing; the caller reports that absence separately.
func rowsWithDiagnostics(registered []anchors.Anchor, diagnostics []string) []anchors.Anchor {
	var found []anchors.Anchor
	for _, anchor := range registered {
		if slices.Contains(diagnostics, anchor.Diagnostic) {
			found = append(found, anchor)
		}
	}
	return found
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
