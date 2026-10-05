package commitment_test

import (
	"reflect"
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
