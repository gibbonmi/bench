package conformance

import (
	"path"
	"slices"
	"strings"
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

// consumerContextLead opens the consumer-context paragraph in each canonical reader. The
// pinned-paragraph rule selects a paragraph by its lead, so this one lead binds both readers'
// paragraphs to their own rows.
const consumerContextLead = "reuses exact available source bytes"

// consumerRows returns the family's consumer-context rows. The reader comes from the
// family's own bounded action rows, so no second inventory names a guidance path.
func (f boundedActionFamily) consumerRows(guidance string) []anchors.Anchor {
	var found []anchors.Anchor
	for _, anchor := range anchorsWithDiagnosticPrefix(consumerContextPrefix) {
		if anchor.File == guidance {
			found = append(found, anchor)
		}
	}
	return found
}

// wantedConsumerRows states the family's two consumer rules apart from the registry: reuse
// verifies the new manifest, and a fresh consumer runs its own retrieval. A reworded or
// deleted row bites.
func (f boundedActionFamily) wantedConsumerRows(guidance string) []wantedRow {
	reader := path.Base(guidance)
	return []wantedRow{
		{anchors.Require, consumerContextPrefix + reader + " permits reuse without verified membership, role, and requiredness"},
		{anchors.Require, consumerContextPrefix + reader + " permits another " + f.peer + "'s receipt or a final cursor to replace fresh required context"},
	}
}

// TestEvidenceConsumerGuidance grades CE104 and CE105 over the one family inventory. Each
// rule takes one Require row on its own reader, each row bites alone, the shipped guidance
// satisfies every row, and the whole consumer paragraph sits inside those rows. The family
// count keeps a dropped family from leaving the rules vacuously green.
func TestEvidenceConsumerGuidance(t *testing.T) {
	families := boundedActionFamilies()
	for _, family := range families {
		t.Run(family.name, func(t *testing.T) {
			guidance := family.guidance(t, family.anchors())
			rows := family.consumerRows(guidance)
			want := family.wantedConsumerRows(guidance)
			if got := len(rows); got != len(want) {
				t.Errorf("%s consumer context anchor count = %d, want %d", family.name, got, len(want))
			}
			for _, expected := range want {
				requireRegisteredRow(t, rows, expected.kind, guidance, expected.diagnostic)
			}
			runAnchorBites(t, rows, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
			requireConformantLiveTree(t, rows)
			requirePinnedParagraph(t, rows, guidance, family.consumerSentences(liveGuidanceText(t, guidance)))
		})
	}
	registered := anchorsWithDiagnosticPrefix(consumerContextPrefix)
	if got, want := len(registered), 2*len(families); got != want {
		t.Errorf("consumer context anchor count = %d, want %d", got, want)
	}
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

// retainedWriteSpecDiagnostics names the registry rows that already state write-spec's
// decision-source and author-fork contract. CE113 keeps that contract, so this owner cites
// those rows rather than restating their prose behind a second pin.
func retainedWriteSpecDiagnostics() []string {
	return []string{
		".agents/commands/bench-write-spec.md dropped the exactly-one Decision source contract",
		".agents/commands/bench-write-spec.md dropped the conversation fork for spec authoring and ticket slicing",
	}
}

// TestEvidenceUnchangedRoutes grades CE112 and CE113. It grades the two unchanged-route
// rows, and it requires the registered rows that state write-spec's decision-source and
// author-fork contract. The registry comment beside those rows owns the why.
func TestEvidenceUnchangedRoutes(t *testing.T) {
	build := boundedActionFamilyNamed(t, "build")
	buildPhase := build.guidance(t, build.anchors())
	const writeSpec = ".agents/commands/bench-write-spec.md"
	want := []struct {
		row  wantedRow
		file string
	}{
		{wantedRow{anchors.Require, unchangedRoutePrefix + path.Base(buildPhase) + " dropped the phase-level full-run control"}, buildPhase},
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

// requirePinnedParagraph binds one guidance paragraph to one row set. Exactly one Require
// needle covers each sentence of the pinned paragraph, so a needle that keeps only its lead
// clause, a needle that repeats another row's text, and an unpinned sentence added to that
// paragraph all bite here. One parser reads the guidance and the needle, so a reflow that
// changes no word changes nothing here.
func requirePinnedParagraph(t *testing.T, rows []anchors.Anchor, guidance string, sentences []string) {
	t.Helper()
	needles := map[string]int{}
	for _, anchor := range rows {
		if anchor.Kind == anchors.Require {
			needles[anchor.Needle]++
		}
	}
	for needle, count := range needles {
		if count != 1 {
			t.Errorf("%d Require rows share the needle %q, want one row each", count, needle)
		}
	}
	if len(sentences) == 0 {
		t.Fatalf("%s states no sentence of the pinned paragraph", guidance)
	}
	for _, sentence := range sentences {
		pinned := 0
		for needle := range needles {
			if strings.Contains(needleText(needle), sentence) {
				pinned++
			}
		}
		if pinned != 1 {
			t.Errorf("guidance sentence %q is pinned by %d Require rows, want one", sentence, pinned)
		}
	}
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
