package commitment_test

import (
	"cmp"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
)

func policy(milestones []commitment.Milestone, active string) commitment.Policy {
	return commitment.Policy{Version: 1, Milestones: milestones, ActiveMilestone: active}
}

func milestone(id string, outcomes ...string) commitment.Milestone {
	result := commitment.Milestone{ID: id}
	for _, id := range outcomes {
		result.Outcomes = append(result.Outcomes, commitment.Outcome{ID: id, Criteria: []commitment.Criterion{{ID: id + ".done", Text: "The outcome is delivered."}}})
	}
	return result
}

func TestCommitmentDisplacementEffects(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	proposed := policy([]commitment.Milestone{milestone("M1", "C", "A", "B")}, "M1")

	plan, err := commitment.BuildPlan(&current, commitment.Proposal{Policy: proposed})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"A", "B"}; !reflect.DeepEqual(plan.Effects.Delayed, want) {
		t.Fatalf("delayed = %v, want %v", plan.Effects.Delayed, want)
	}
}

func TestCommitmentRemovalEffects(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	proposed := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")

	plan, err := commitment.BuildPlan(&current, commitment.Proposal{Policy: proposed})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"B"}; !reflect.DeepEqual(plan.Effects.Removed, want) {
		t.Fatalf("removed = %v, want %v", plan.Effects.Removed, want)
	}
}

func TestCommitmentSwitchEffects(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A", "B"), milestone("M2", "C")}, "M1")
	proposed := policy([]commitment.Milestone{milestone("M1", "A", "B"), milestone("M2", "C")}, "M2")

	plan, err := commitment.BuildPlan(&current, commitment.Proposal{Policy: proposed})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"A", "B"}; !reflect.DeepEqual(plan.Effects.Delayed, want) {
		t.Fatalf("delayed = %v, want %v", plan.Effects.Delayed, want)
	}
	if want := []string{"A", "B"}; !reflect.DeepEqual(plan.Effects.Switched, want) {
		t.Fatalf("switched = %v, want %v", plan.Effects.Switched, want)
	}
}

// obligationFreeBinding names the binding that obligationFreeMilestone approves.
const obligationFreeBinding = "spec"

// obligationFreeMilestone returns bindObligationFree for the first outcome of milestone and
// the binding obligationFreeBinding.
func obligationFreeMilestone(milestone commitment.Milestone, row string) commitment.Milestone {
	return bindObligationFree(milestone, milestone.Outcomes[0].ID, row, obligationFreeBinding)
}

// bindObligationFree returns a copy of milestone whose outcome named id also approves binding
// with no obligation. The binding source derives from binding alone, so two outcomes that
// approve one name approve an equal binding. A nonempty row is a source that the outcome
// owns; an empty row adds no source.
func bindObligationFree(milestone commitment.Milestone, id, row, binding string) commitment.Milestone {
	milestone.Outcomes = slices.Clone(milestone.Outcomes)
	index := slices.IndexFunc(milestone.Outcomes, func(outcome commitment.Outcome) bool { return outcome.ID == id })
	outcome := &milestone.Outcomes[index]
	if row != "" {
		outcome.Sources = []commitment.SourceBinding{{ID: row, Path: "roadmap/" + row + ".md", Identity: commitment.Identity([]byte(row))}}
	}
	outcome.Deliverables = append(slices.Clone(outcome.Deliverables), commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: binding, Path: "specs/" + binding + "/spec.md", Identity: commitment.Identity([]byte(binding))}})
	return milestone
}

// A plan refuses a new or changed binding that can close no row of its outcome, in any
// outcome of any milestone. It plans a rowless deliverable. It also plans a legacy binding
// that the same outcome keeps unchanged when that binding is already obligation-free.
func TestCommitmentPlanRefusesObligationFreeBinding(t *testing.T) {
	empty := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	rowless := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "")}, "M1")
	legacy := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "FT1")}, "M1")
	changed := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "FT1")}, "M1")
	changed.Milestones[0].Outcomes[0].Deliverables[0].Source.Identity = commitment.Identity([]byte("changed"))
	second := policy([]commitment.Milestone{bindObligationFree(milestone("M1", "A", "B"), "B", "FT2", obligationFreeBinding)}, "M1")
	for _, row := range []struct {
		name              string
		current, proposed commitment.Policy
		refused, binding  string
	}{
		{name: "new binding", current: empty, proposed: legacy, refused: "A"},
		{name: "rowless outcome", current: empty, proposed: rowless},
		{name: "kept legacy binding", current: legacy, proposed: policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A", "C"), "FT1")}, "M1")},
		{name: "changed legacy binding", current: legacy, proposed: changed, refused: "A"},
		{name: "rowless binding gains a row", current: rowless, proposed: legacy, refused: "A"},
		{name: "legacy binding moved to another outcome", current: legacy, proposed: second, refused: "B"},
		{name: "new binding after a kept binding", current: legacy, proposed: policy([]commitment.Milestone{bindObligationFree(legacy.Milestones[0], "A", "", "draft")}, "M1"), refused: "A", binding: "draft"},
		{name: "second outcome", current: empty, proposed: second, refused: "B"},
		{name: "inactive milestone", current: legacy, proposed: policy([]commitment.Milestone{legacy.Milestones[0], obligationFreeMilestone(milestone("M2", "B"), "FT2")}, "M1"), refused: "B"},
	} {
		t.Run(row.name, func(t *testing.T) {
			binding := cmp.Or(row.binding, obligationFreeBinding)
			_, err := commitment.BuildPlan(&row.current, commitment.Proposal{Policy: row.proposed})
			if row.refused == "" && err != nil {
				t.Fatalf("BuildPlan() = %v, want a plan", err)
			}
			if row.refused != "" && (err == nil || !strings.Contains(err.Error(), "names no obligation") || !strings.Contains(err.Error(), `"`+row.refused+`"`) || !strings.Contains(err.Error(), `"`+binding+`"`)) {
				t.Fatalf("BuildPlan() = %v, want the refusal of outcome %q binding %q", err, row.refused, binding)
			}
		})
	}
}

// A plan binds only the content that its proposal keeps open. A deliverable that the
// proposal drops binds nothing, and the row that the outcome keeps stays bound.
func TestCommitmentPlanSourcesOmitDroppedDeliverable(t *testing.T) {
	current := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "FT1")}, "M1")
	proposed := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "FT1")}, "M1")
	proposed.Milestones[0].Outcomes[0].Deliverables = nil
	plan, err := commitment.BuildPlan(&current, commitment.Proposal{Policy: proposed})
	if err != nil {
		t.Fatal(err)
	}
	if want := current.Milestones[0].Outcomes[0].Sources; !reflect.DeepEqual(plan.Sources, want) {
		t.Fatalf("plan sources = %+v, want only the kept row %+v", plan.Sources, want)
	}
}

// A plan lists its sources in byte order of the identifier, then the path, whatever order
// the policy holds them in. Two deliverables with one identifier order by path.
func TestCommitmentPlanSourcesAreCanonical(t *testing.T) {
	owned := milestone("M1", "A", "B")
	for index, row := range []struct{ id, row, path string }{{"A", "FT9", "specs/b/spec.md"}, {"B", "FT10", "specs/a/spec.md"}} {
		source := commitment.SourceBinding{ID: row.row, Path: "roadmap/" + row.row + ".md", Identity: commitment.Identity([]byte(row.row))}
		deliverable := commitment.SourceBinding{ID: "spec", Path: row.path, Identity: commitment.Identity([]byte(row.path))}
		owned.Outcomes[index].Sources = []commitment.SourceBinding{source}
		owned.Outcomes[index].Deliverables = []commitment.DeliveryBinding{{Source: deliverable, Obligations: []string{row.row}}}
	}
	plan, err := commitment.BuildPlan(nil, commitment.Proposal{Policy: policy([]commitment.Milestone{owned}, "M1")})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, source := range plan.Sources {
		got = append(got, source.ID+" "+source.Path)
	}
	if want := []string{"FT10 roadmap/FT10.md", "FT9 roadmap/FT9.md", "spec specs/a/spec.md", "spec specs/b/spec.md"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("plan sources = %v, want %v", got, want)
	}
}

func TestCommitmentReorderDelaysPassedOutcomes(t *testing.T) {
	current := policy([]commitment.Milestone{milestone("M1", "A", "B", "C")}, "M1")
	proposed := policy([]commitment.Milestone{milestone("M1", "C", "A", "B")}, "M1")
	plan, err := commitment.BuildPlan(&current, commitment.Proposal{Policy: proposed})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"A", "B"}; !reflect.DeepEqual(plan.Effects.Delayed, want) {
		t.Fatalf("delayed=%v, want %v", plan.Effects.Delayed, want)
	}
}
