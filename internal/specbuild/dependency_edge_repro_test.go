package specbuild

import (
	"path/filepath"
	"strings"
	"testing"
)

// A sibling that names a dependent consumer in Integration surfaces requires
// the consumer's Blocked by: to name that sibling, even when the consumer
// otherwise contracts against raw input. Assign refuses the consumer until
// the reciprocal edge exists.
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

func TestAssignIgnoresNonDependentBasenameMentions(t *testing.T) {
	for _, tc := range []struct {
		name     string
		surfaces string
	}{
		{name: "near-name target", surfaces: "resolved value→stone.md + S1"},
		{name: "incidental prose", surfaces: "one.md remains documented→`internal/other` + S1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := repo(t)
			one := "# One\n\n" +
				"Blocked by: none\n" +
				"Ownership fence: `internal/specbuild`\n" +
				"Integration surfaces: one implementation→`internal/specbuild` + O1\n" +
				"Contracts: none crosses\n" +
				"Closure: O1/assign\n\n" +
				"- [ ] [O1] (covers local) one ticket assigns\n\n" +
				"## Red mutations\n\n" +
				"| criterion | mutation | owner | operation sequence |\n" +
				"|---|---|---|---|\n" +
				"| O1/assign | reject the ticket | the assign control | assign one.md, require success |\n"
			sibling := "# Sibling\n\n" +
				"Blocked by: none\n" +
				"Ownership fence: `internal/other`\n" +
				"Integration surfaces: " + tc.surfaces + "\n" +
				"Contracts: none crosses\n" +
				"Closure: S1/sibling\n\n" +
				"- [ ] [S1] (covers local) sibling ticket stays independent\n\n" +
				"## Red mutations\n\n" +
				"| criterion | mutation | owner | operation sequence |\n" +
				"|---|---|---|---|\n" +
				"| S1/sibling | reject the ticket | the assign control | assign sibling.md, require success |\n"
			write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), one)
			write(t, filepath.Join(root, "specs", "build demo", "tickets", "sibling.md"), sibling)
			git(t, root, "add", ".")
			git(t, root, "commit", "-qm", "stage non-dependent sibling mentions")
			service := New(root, greenGate{}, realOwner{})
			if _, err := service.Start(t.Context(), "build demo"); err != nil {
				t.Fatalf("Start: %v", err)
			}
			if _, _, err := service.Assign(t.Context(), "build demo", "one.md", "one request"); err != nil {
				t.Fatalf("Assign one.md: %v", err)
			}
		})
	}
}

// Refresh revalidates reciprocal edges against current sibling tickets. A
// staged repair that names a preserved consumer as dependent must refuse the
// refresh rather than carry stale dependency metadata onto the candidate.
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
