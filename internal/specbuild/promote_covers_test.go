package specbuild

import (
	"path/filepath"
	"testing"
)

// coversMapRows renders one opted-in coverage map row per id, so a case declares
// only the map identities its composition has to account for.
func coversMapRows(ids ...string) []string {
	rows := make([]string, len(ids))
	for i, id := range ids {
		rows[i] = "| " + id + " | 1 | behavior | seam | red | why |"
	}
	return rows
}

// coversTicket renders the assigned ticket carrying rows. Tests that pin a digest
// share this one rendering with the fixture, so the digest they compute is of a
// ticket text the parser really produces.
func coversTicket(rows string) string {
	return "# One\n\nOwnership fence: internal/specbuild\nAssumptions: receipt contract\n\n" + rows
}

// coveredPromotionFixture composes a reviewed run ready to promote under an
// opted-in map: the spec declares ids and the assigned ticket charges rows. The
// spec and the ticket are the real files on disk because promote resolves
// totality from them, not from anything the run record carries.
func coveredPromotionFixture(t *testing.T, ids []string, rows string, configure ...func(string)) checkpointFixture {
	t.Helper()
	seed := func(root string) {
		write(t, filepath.Join(root, "specs", "build demo", "spec.md"), coversSpec(optedInCoversHeader, coversMapRows(ids...)...))
		write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), coversTicket(rows))
	}
	return reviewedPromotionFixture(t, append([]func(string){seed}, configure...)...)
}

// zeroAssignmentPromotionFixture composes a reviewed run under an opted-in map
// that never assigns a ticket at all: requireCoversMapping never runs (assign
// is never called), so promote's own totality check is the only thing that can
// answer for the map, whether the gap is a coverage shortfall or the map's own
// grammar failing to validate.
func zeroAssignmentPromotionFixture(t *testing.T, mapRows ...string) checkpointFixture {
	t.Helper()
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "spec.md"), coversSpec(optedInCoversHeader, mapRows...))
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "covers fixture")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	run := loadRun(t, service)
	receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatal(err)
	}
	return checkpointFixture{root: root, service: service, run: loadRun(t, service)}
}

// [RG4] A composition that never assigns a single ticket still owes every
// declared map row a claimant: totality names every one of them, not just the
// ones a ticket happened to charge.
func TestPromoteRefusesAZeroAssignmentOptedInRunNamingEveryMapID(t *testing.T) {
	fixture := zeroAssignmentPromotionFixture(t, coversMapRows("AB1", "AB2")...)
	requirePromotionRefusal(t, fixture, "spec build promote requires every coverage map row covered by an integrated assignment's ticket, but nothing covers AB1, AB2")
}

// [RG5] requireCoversTotality's own invalid-map branch: nothing else in this
// run's path (no ticket was ever assigned) could have already refused the
// map's duplicate row id, so reaching this refusal proves totality checks the
// map's violations itself rather than trusting a prior assign to have done so.
func TestPromoteRefusesAnOptedInMapThatFailsIDValidationAtTotality(t *testing.T) {
	fixture := zeroAssignmentPromotionFixture(t, coversMapRows("AB1", "AB1")...)
	requirePromotionRefusal(t, fixture, "spec build promote requires the spec's opted-in coverage map to validate, but it reports coverage map row 2 has a duplicate row id 'AB1' (first used at row 1)")
}

func TestPromoteRefusesUncoveredMapRowsBeforeTheGate(t *testing.T) {
	fixture := coveredPromotionFixture(t, []string{"AB1", "AB2", "AB3"},
		"- [ ] [R10] (covers AB1) mapped\n- [ ] [R11] (covers local) ticket-time repair\n")
	requirePromotionRefusal(t, fixture, "spec build promote requires every coverage map row covered by an integrated assignment's ticket, but nothing covers AB2, AB3")
}

// An unassigned file in the tickets directory is unverified by construction —
// nothing checkpointed it, integrated it, or reviewed it — so a promote that
// counted it would accept a map satisfied entirely on paper.
func TestPromoteIgnoresAnUnassignedDecoyTicket(t *testing.T) {
	fixture := coveredPromotionFixture(t, []string{"AB1", "AB2"}, "- [ ] [R10] (covers AB1) mapped\n", func(root string) {
		write(t, filepath.Join(root, "specs", "build demo", "tickets", "decoy.md"), "# Decoy\n\nOwnership fence: internal/specbuild\n\n- [ ] [R90] (covers AB2) never assigned\n")
	})
	requirePromotionRefusal(t, fixture, "spec build promote requires every coverage map row covered by an integrated assignment's ticket, but nothing covers AB2")
}

// Promote refuses a covers edit made after the evidence that graded it. The edit
// cannot be a live file write: promote's clean-checkout precondition would refuse
// first, and committing it would move the working tip into recomposition. So the
// run records the digest of the covers text that was assigned while the file on
// disk holds the edited one — the same mismatch, reached through promote.
func TestPromoteRefusesAPostIntegrationCoversEdit(t *testing.T) {
	fixture := coveredPromotionFixture(t, []string{"AB1"}, "- [ ] [R10] (covers AB1) mapped\n")
	updatePromotionRun(t, fixture, func(run *record) {
		for key, assigned := range run.Assignments {
			assigned.TicketDigest = digest(coversTicket("- [ ] [R10] (covers local) mapped\n"))
			run.Assignments[key] = assigned
		}
	})
	requirePromotionRefusal(t, fixture, "spec build promote requires every integrated assignment's ticket unchanged since assign: checkpoint ticket drifted")
}

func TestPromoteAcceptsFullCoverageIncludingOverCoverage(t *testing.T) {
	fixture := coveredPromotionFixture(t, []string{"AB1", "AB2"},
		"- [ ] [R10] (covers AB1) mapped\n- [ ] [R11] (covers AB1) also mapped\n- [ ] [R12] (covers AB2) mapped\n- [ ] [R13] (covers local) ticket-time repair\n")
	owner := &promotionGate{accept: true}
	fixture.service.gate = owner
	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.State != "terminal" {
		t.Fatalf("Promote = %#v, %v", status, err)
	}
	if owner.executions != 1 {
		t.Fatalf("prospective gate executions = %d, want 1", owner.executions)
	}
}

// [AN2] A `(covers ...)` mention that is not bracket-adjacent is prose, not a
// claimant: a composition whose only mention of a map ID sits in a row's
// trailing prose refuses to name that ID covered. The row still carries a
// bracket-adjacent `(covers local)` so assign accepts it; only promote's
// totality check can see that AB2 has no real claimant.
func TestPromoteRefusesAMapRowClaimedOnlyByAProseMention(t *testing.T) {
	fixture := coveredPromotionFixture(t, []string{"AB1", "AB2"},
		"- [ ] [R10] (covers AB1) mapped\n- [ ] [R11] (covers local) see the example `(covers AB2)` in the docs.\n")
	requirePromotionRefusal(t, fixture, "spec build promote requires every coverage map row covered by an integrated assignment's ticket, but nothing covers AB2")
}

// `local` is accepted machinery, not coverage: a map row whose only claimant
// declares itself unmapped is still unproven.
func TestPromoteRefusesAMapRowCoveredOnlyByLocalRows(t *testing.T) {
	fixture := coveredPromotionFixture(t, []string{"AB1", "AB2"},
		"- [ ] [R10] (covers AB1) mapped\n- [ ] [R11] (covers local) stands in for AB2\n")
	requirePromotionRefusal(t, fixture, "spec build promote requires every coverage map row covered by an integrated assignment's ticket, but nothing covers AB2")
}

// The refusal has to name the operation the reviewer would re-run, not borrow a
// sibling's wording.
func TestPromoteTotalityRefusalNamesPromote(t *testing.T) {
	fixture := coveredPromotionFixture(t, []string{"AB1", "AB2"}, "- [ ] [R10] (covers AB1) mapped\n")
	owner := &promotionGate{accept: true}
	fixture.service.gate = owner
	_, err := fixture.service.Promote(t.Context(), "build demo")
	if err == nil {
		t.Fatal("Promote accepted an uncovered map row")
	}
	if word := operationWordIn(t, err.Error()); word != string(mutationPromote) {
		t.Fatalf("operation word = %q, want %q", word, mutationPromote)
	}
}
