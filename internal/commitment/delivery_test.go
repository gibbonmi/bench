package commitment_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
)

const deliverySpec = "specs/a/spec.md"

// deliveryPolicy approves deliverySpec as the delivery of the named obligations of outcome
// A. Outcome A owns FT1 and FT3, and outcome B owns FT2.
func deliveryPolicy(obligations ...string) commitment.Policy {
	source := func(id string) commitment.SourceBinding {
		return commitment.SourceBinding{ID: id, Path: "roadmap/" + id + ".md", Identity: "sha256:" + id}
	}
	p := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	a, b := &p.Milestones[0].Outcomes[0], &p.Milestones[0].Outcomes[1]
	a.Sources, b.Sources = []commitment.SourceBinding{source("FT1"), source("FT3")}, []commitment.SourceBinding{source("FT2")}
	a.Deliverables = []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: "spec", Path: deliverySpec, Identity: "sha256:spec"}, Obligations: obligations}}
	return p
}

// A partial delivery records one fact, satisfies only its obligation, and leaves its
// outcome in the remaining sequence. A complete delivery removes the outcome. A second
// delivery of the same binding refuses, and an unapproved path changes nothing.
func TestCommitmentDeliver(t *testing.T) {
	partial, closed, err := commitment.Deliver(deliveryPolicy("FT1"), deliverySpec, "source", "evidence")
	if err != nil || len(closed) != 1 || closed[0].ID != "FT1" || len(partial.Deliveries) != 1 {
		t.Fatalf("partial Deliver = %+v, %+v, %v; want one fact closing FT1", partial.Deliveries, closed, err)
	}
	want := commitment.DeliveryFact{Milestone: "M1", Outcome: "A", Binding: "spec", Identity: "sha256:spec", Source: "source", Evidence: "evidence"}
	if partial.Deliveries[0] != want {
		t.Fatalf("delivery fact = %+v, want %+v", partial.Deliveries[0], want)
	}
	if satisfied := commitment.Satisfied(partial); len(satisfied) != 1 || !satisfied["FT1"] {
		t.Fatalf("partial Satisfied = %v, want FT1 alone", satisfied)
	}
	if remaining := commitment.Remaining(partial); !slices.Equal(remaining, []string{"A", "B"}) {
		t.Fatalf("partial Remaining = %v, want the open outcome A and B", remaining)
	}
	if !commitment.ScopeDelivered(partial, []string{"owned.txt", deliverySpec}) || commitment.ScopeDelivered(partial, []string{"owned.txt"}) {
		t.Fatal("ScopeDelivered must need at least one delivered deliverable in the scope")
	}
	if _, _, err := commitment.Deliver(partial, deliverySpec, "other", "other"); err == nil || !strings.Contains(err.Error(), "already delivered") {
		t.Fatalf("second Deliver = %v, want the already-delivered refusal", err)
	}
	complete, _, err := commitment.Deliver(deliveryPolicy("FT1", "FT3"), deliverySpec, "source", "evidence")
	if err != nil || !slices.Equal(commitment.Remaining(complete), []string{"B"}) {
		t.Fatalf("complete Remaining = %v, %v; want B alone", commitment.Remaining(complete), err)
	}
	unchanged, closed, err := commitment.Deliver(deliveryPolicy("FT1"), "specs/other/spec.md", "source", "evidence")
	if err != nil || closed != nil || len(unchanged.Deliveries) != 0 {
		t.Fatalf("unapproved Deliver = %+v, %+v, %v; want no change", unchanged.Deliveries, closed, err)
	}
}

// An outcome with no source is its own obligation, so its delivery records one fact that
// satisfies no source and completes only that outcome. A binding with no obligation of an
// outcome that has sources records nothing. PathDelivered names only a delivered path.
func TestCommitmentDeliverRowless(t *testing.T) {
	const rowlessSpec = "specs/c/spec.md"
	p := deliveryPolicy()
	p.Milestones[0].Outcomes = append(p.Milestones[0].Outcomes, commitment.Outcome{ID: "C", Criteria: []commitment.Criterion{{ID: "c-done", Text: "C is delivered."}}, Sources: []commitment.SourceBinding{},
		Deliverables: []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: "rowless", Path: rowlessSpec, Identity: "sha256:rowless"}}}})
	delivered, closed, err := commitment.Deliver(p, rowlessSpec, "source", "evidence")
	want := commitment.DeliveryFact{Milestone: "M1", Outcome: "C", Binding: "rowless", Identity: "sha256:rowless", Source: "source", Evidence: "evidence"}
	if err != nil || len(closed) != 0 || len(delivered.Deliveries) != 1 || delivered.Deliveries[0] != want {
		t.Fatalf("rowless Deliver = %+v, %+v, %v; want the one fact %+v and no source", delivered.Deliveries, closed, err, want)
	}
	if remaining := commitment.Remaining(delivered); !slices.Equal(remaining, []string{"A", "B"}) {
		t.Fatalf("rowless Remaining = %v, want A and B open", remaining)
	}
	if len(commitment.Satisfied(delivered)) != 0 {
		t.Fatalf("rowless Satisfied = %v, want no source", commitment.Satisfied(delivered))
	}
	if !commitment.PathDelivered(delivered, rowlessSpec) || commitment.PathDelivered(delivered, deliverySpec) || commitment.PathDelivered(p, rowlessSpec) {
		t.Fatal("PathDelivered must name exactly the delivered deliverable")
	}
	unchanged, closed, err := commitment.Deliver(p, deliverySpec, "source", "evidence")
	if err != nil || closed != nil || len(unchanged.Deliveries) != 0 {
		t.Fatalf("obligation-free Deliver = %+v, %+v, %v; want no fact for an outcome with sources", unchanged.Deliveries, closed, err)
	}
}

// A rowless outcome with two approved deliverables stays open after the first delivery and
// completes only after the second, in either order.
func TestCommitmentDeliverRowlessPair(t *testing.T) {
	paths := []string{"specs/c/spec.md", "specs/d"}
	p := deliveryPolicy()
	c := commitment.Outcome{ID: "C", Criteria: []commitment.Criterion{{ID: "c-done", Text: "C is delivered."}}, Sources: []commitment.SourceBinding{}}
	for i, path := range paths {
		c.Deliverables = append(c.Deliverables, commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: "rowless" + string(rune('1'+i)), Path: path, Identity: "sha256:" + path}})
	}
	p.Milestones[0].Outcomes = append(p.Milestones[0].Outcomes, c)
	for _, order := range [][]string{paths, {paths[1], paths[0]}} {
		first, _, err := commitment.Deliver(p, order[0], "source", "evidence")
		if err != nil || !slices.Equal(commitment.Remaining(first), []string{"A", "B", "C"}) {
			t.Fatalf("first rowless Deliver of %s: Remaining = %v, %v; want C still open", order[0], commitment.Remaining(first), err)
		}
		second, _, err := commitment.Deliver(first, order[1], "source", "evidence")
		if err != nil || !slices.Equal(commitment.Remaining(second), []string{"A", "B"}) {
			t.Fatalf("second rowless Deliver of %s: Remaining = %v, %v; want C delivered", order[1], commitment.Remaining(second), err)
		}
	}
}

// Validate refuses each delivery fact that does not name exactly one approved binding with
// its approved identity and a reviewed source and evidence.
func TestCommitmentDeliveryFactValidation(t *testing.T) {
	delivered, _, err := commitment.Deliver(deliveryPolicy("FT1"), deliverySpec, "source", "evidence")
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct {
		name string
		edit func(*commitment.Policy)
	}{
		{"unresolved-binding", func(p *commitment.Policy) { p.Deliveries[0].Binding = "other" }},
		{"duplicate-fact", func(p *commitment.Policy) { p.Deliveries = append(p.Deliveries, p.Deliveries[0]) }},
		{"identity-mismatch", func(p *commitment.Policy) { p.Deliveries[0].Identity = "sha256:other" }},
		{"empty-source", func(p *commitment.Policy) { p.Deliveries[0].Source = "" }},
		{"empty-evidence", func(p *commitment.Policy) { p.Deliveries[0].Evidence = "" }},
	} {
		t.Run(row.name, func(t *testing.T) {
			p := delivered
			p.Deliveries = slices.Clone(delivered.Deliveries)
			row.edit(&p)
			if err := commitment.Validate(p); err == nil || !strings.Contains(err.Error(), "invalid delivery") {
				t.Fatalf("Validate = %v, want the invalid-delivery refusal", err)
			}
		})
	}
}
