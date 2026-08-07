package specbuild

import (
	"path/filepath"
	"strings"
	"testing"
)

// Regression for the conformance-harness-scope miss: a producer ticket's
// Integration surfaces names a sibling basename as the dependent consumer of
// its exported value, but that consumer declares "Blocked by: none" and
// contracts against the raw input instead. Assign must refuse the consumer
// until the reciprocal edge exists.
func TestAssignRefusesConsumerMissingReciprocalDependencyEdge(t *testing.T) {
	root := repo(t)
	producer := "# Expose resolved checks\n\n" +
		"Blocked by: none\n" +
		"Ownership fence: `internal/canary`\n" +
		"Integration surfaces: exported resolved check→consumer.md + C1 after lifecycle refresh\n" +
		"Contracts: resolved check name crosses `internal/canary`→the refreshed consumer assignment, asserted by P1 against the real inventory\n" +
		"Closure: P1/check-override\n\n" +
		"- [ ] [P1] (covers local) the inventory exposes the resolved check\n\n" +
		"## Red mutations\n\n" +
		"| criterion | mutation | owner | operation sequence |\n" +
		"|---|---|---|---|\n" +
		"| P1/check-override | ignore a present fixture CHECK | the canary regression | apply the omission, run the canary regression, require the fallback |\n"
	consumer := "# Scope fixture bites\n\n" +
		"Blocked by: none\n" +
		"Ownership fence: `internal/conformance`\n" +
		"Integration surfaces: canary fixture inventory→`internal/conformance` + C1\n" +
		"Contracts: resolved check name crosses the canary inventory→`internal/conformance`, asserted by C1 against the real inventory\n" +
		"Closure: C1/scoped-run\n\n" +
		"- [ ] [C1] (covers local) each fixture journey runs exactly its resolved check\n\n" +
		"## Red mutations\n\n" +
		"| criterion | mutation | owner | operation sequence |\n" +
		"|---|---|---|---|\n" +
		"| C1/scoped-run | widen the journey back to the full table | the fixture-bite regression | apply the widening, run the journey, require the timing identity to differ |\n"
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "producer.md"), producer)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "consumer.md"), consumer)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "stage producer and consumer tickets")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	_, _, err := service.Assign(t.Context(), "build demo", "consumer.md", "consumer request")
	if err == nil {
		t.Fatal("assign accepted a consumer named as dependent by producer.md while declaring Blocked by: none")
	}
	if !strings.Contains(err.Error(), "Blocked by") {
		t.Fatalf("assign refusal = %v, want the missing reciprocal Blocked by edge named", err)
	}
	// The repaired shape assigns: the consumer names the producer back, and the
	// producer itself stays assignable because the consumer's mention of its
	// basename rides the blocker direction, not a dependent claim.
	repaired := strings.Replace(consumer, "Blocked by: none", "Blocked by: producer.md", 1)
	repaired = strings.Replace(repaired, "canary fixture inventory→`internal/conformance` + C1", "canary fixture inventory→producer.md + C1", 1)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "consumer.md"), repaired)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "repair the reciprocal edge")
	if _, err := service.Promote(t.Context(), "build demo"); err != nil {
		t.Fatalf("recomposing Promote: %v", err)
	}
	if _, _, err := service.Assign(t.Context(), "build demo", "consumer.md", "consumer request"); err != nil {
		t.Fatalf("Assign repaired consumer: %v", err)
	}
	if _, _, err := service.Assign(t.Context(), "build demo", "producer.md", "producer request"); err != nil {
		t.Fatalf("Assign producer: %v", err)
	}
}

// The mid-run variant of the same miss: a repair ticket staged after the
// consumer's assign names the consumer's ticket as a dependent, so the
// preserved assignment's metadata is stale by construction — refresh must
// refuse instead of carrying the stale contract onto the repaired candidate.
func TestRefreshRefusesStaleDependencyMetadata(t *testing.T) {
	repair := "# Repair landing\n\n" +
		"Blocked by: none\n" +
		"Ownership fence: `internal/landing`\n" +
		"Integration surfaces: repaired landing value→one.md + R90\n" +
		"Contracts: none crosses\n" +
		"Closure: R90/landing-fix\n\n" +
		"- [ ] [R90] repair the prerequisite defect\n\n" +
		"## Red mutations\n\n" +
		"| criterion | mutation | owner | operation sequence |\n" +
		"|---|---|---|---|\n" +
		"| R90/landing-fix | omit the landing fix | the landing regression | apply the omission, run the landing regression, require red |\n"
	fixture := newCleanRefreshFixtureWithRepair(t, repair)
	_, _, err := fixture.service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, debugReceiptFor(t, fixture)))
	if err == nil {
		t.Fatal("refresh carried a stale dependency contract past a repair that names its ticket as a dependent")
	}
	if !strings.Contains(err.Error(), "Blocked by") {
		t.Fatalf("refresh refusal = %v, want the missing reciprocal Blocked by edge named", err)
	}
}
