package commitment_test

import (
	"reflect"
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

// obligationFreeMilestone returns milestone whose first outcome approves one binding with no
// obligation. A nonempty row is a source that the outcome owns; an empty row leaves the
// outcome rowless.
func obligationFreeMilestone(milestone commitment.Milestone, row string) commitment.Milestone {
	outcome := &milestone.Outcomes[0]
	if row != "" {
		outcome.Sources = []commitment.SourceBinding{{ID: row, Path: "roadmap/" + row + ".md", Identity: commitment.Identity([]byte(row))}}
	}
	outcome.Deliverables = []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: obligationFreeBinding, Path: "specs/" + outcome.ID + "/spec.md", Identity: commitment.Identity([]byte(obligationFreeBinding))}}}
	return milestone
}

// A plan refuses a new or changed binding that can close no row of its outcome, in any
// milestone. It plans a rowless deliverable and a legacy binding that it keeps unchanged.
func TestCommitmentPlanRefusesObligationFreeBinding(t *testing.T) {
	empty := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	legacy := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "FT1")}, "M1")
	changed := policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "FT1")}, "M1")
	changed.Milestones[0].Outcomes[0].Deliverables[0].Source.Identity = commitment.Identity([]byte("changed"))
	for _, row := range []struct {
		name              string
		current, proposed commitment.Policy
		refused           string
	}{
		{name: "new binding", current: empty, proposed: legacy, refused: "A"},
		{name: "rowless outcome", current: empty, proposed: policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A"), "")}, "M1")},
		{name: "kept legacy binding", current: legacy, proposed: policy([]commitment.Milestone{obligationFreeMilestone(milestone("M1", "A", "C"), "FT1")}, "M1")},
		{name: "changed legacy binding", current: legacy, proposed: changed, refused: "A"},
		{name: "inactive milestone", current: legacy, proposed: policy([]commitment.Milestone{legacy.Milestones[0], obligationFreeMilestone(milestone("M2", "B"), "FT2")}, "M1"), refused: "B"},
	} {
		t.Run(row.name, func(t *testing.T) {
			_, err := commitment.BuildPlan(&row.current, commitment.Proposal{Policy: row.proposed})
			if row.refused == "" && err != nil {
				t.Fatalf("BuildPlan() = %v, want a plan", err)
			}
			if row.refused != "" && (err == nil || !strings.Contains(err.Error(), "names no obligation") || !strings.Contains(err.Error(), `"`+row.refused+`"`) || !strings.Contains(err.Error(), `"`+obligationFreeBinding+`"`)) {
				t.Fatalf("BuildPlan() = %v, want the refusal of outcome %q binding %q", err, row.refused, obligationFreeBinding)
			}
		})
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
