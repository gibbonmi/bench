package commitment_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestCommitmentOutlook decides the reader projection from one policy and runtime snapshot.
// Each case writes its expected outlook by hand, so a projection that ranked by another
// rule than admission names another outcome. Each case also asks Admit about every
// remaining outcome: admission accepts exactly the next and the active outcomes.
func TestCommitmentOutlook(t *testing.T) {
	independent := policy([]commitment.Milestone{milestone("M1", "A", "B", "C")}, "M1")
	dependent := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	dependent.Milestones[0].Outcomes[1].Dependencies = []string{"A"}
	parallel := policy([]commitment.Milestone{milestone("M1", "A", "B", "C")}, "M1")
	parallel.ParallelGrants = []commitment.ParallelGrant{{Outcomes: []string{"A", "B"}}}
	inactive := policy([]commitment.Milestone{milestone("M1", "A")}, "")
	delivered := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	delivered.Milestones[0].Outcomes[0].Deliverables = []commitment.DeliveryBinding{approved("A.spec", "specs/a/spec.md")}
	delivered.Deliveries = []commitment.DeliveryFact{{Milestone: "M1", Outcome: "A", Binding: "A.spec", Identity: commitment.Identity([]byte("specs/a/spec.md")), Source: "source", Evidence: "evidence"}}
	sole := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	sole.Milestones[0].Outcomes[1].Deliverables = []commitment.DeliveryBinding{approved("B.spec", "specs/b/spec.md")}
	several := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	several.Milestones[0].Outcomes[1].Deliverables = []commitment.DeliveryBinding{approved("B.spec", "specs/b/spec.md"), approved("B.tickets", "specs/b-light")}
	partial := several
	partial.Deliveries = []commitment.DeliveryFact{{Milestone: "M1", Outcome: "B", Binding: "B.spec", Identity: commitment.Identity([]byte("specs/b/spec.md")), Source: "source", Evidence: "evidence"}}
	blockedA := intent.CommitmentState{Blockers: []intent.OutcomeBlocker{{Outcome: "A", Reason: "waits on a vendor fix"}}}
	activeA := intent.CommitmentState{Claims: []intent.OutcomeClaim{{Milestone: "M1", Outcome: "A"}}}

	for _, tc := range []struct {
		name   string
		policy *commitment.Policy
		state  intent.CommitmentState
		want   commitment.Outlook
	}{
		{name: "absent policy", want: commitment.Outlook{State: "adoption-required", Operation: "plan"}},
		{name: "no active milestone", policy: &inactive, want: commitment.Outlook{State: "no-active-milestone", Operation: "plan"}},
		{name: "first outcome", policy: &independent, want: commitment.Outlook{State: "eligible", Milestone: "M1", Next: "A", Waiting: []string{"B", "C"}, Operation: "start"}},
		{name: "blocked predecessor", policy: &independent, state: blockedA, want: commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Blocked: blockedA.Blockers, Waiting: []string{"C"}, Operation: "start"}},
		{name: "blocked dependency", policy: &dependent, state: blockedA, want: commitment.Outlook{State: "all-blocked", Milestone: "M1", Blocked: blockedA.Blockers, Waiting: []string{"B"}, Operation: "plan"}},
		{name: "active outcome", policy: &independent, state: activeA, want: commitment.Outlook{State: "active", Milestone: "M1", Active: []string{"A"}, Waiting: []string{"B", "C"}}},
		{name: "parallel grant", policy: &parallel, state: activeA, want: commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Active: []string{"A"}, Waiting: []string{"C"}, Operation: "start"}},
		{name: "delivered milestone", policy: &delivered, want: commitment.Outlook{State: "delivered", Milestone: "M1", Operation: "verify"}},
		{name: "sole deliverable", policy: &sole, state: blockedA, want: commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Deliverable: "specs/b/spec.md", Blocked: blockedA.Blockers, Operation: "start"}},
		{name: "several deliverables", policy: &several, state: blockedA, want: commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Blocked: blockedA.Blockers, Operation: "start"}},
		{name: "one deliverable left", policy: &partial, state: blockedA, want: commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Deliverable: "specs/b-light", Blocked: blockedA.Blockers, Operation: "start"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := commitment.Project(tc.policy, tc.state)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("outlook = %#v, want %#v", got, tc.want)
			}
			if tc.policy == nil {
				return
			}
			for _, id := range commitment.Remaining(*tc.policy) {
				binding := intent.DeliveryBinding{Assignment: "assignment", Milestone: "M1", Outcome: id}
				_, err := commitment.Admit(*tc.policy, tc.state, binding)
				if admitted := err == nil; admitted != (id == got.Next || slices.Contains(got.Active, id)) {
					t.Errorf("Admit(%s) = %v, but the outlook names next %q and active %q", id, err, got.Next, got.Active)
				}
			}
		})
	}
}

func approved(id, path string) commitment.DeliveryBinding {
	return commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: id, Path: path, Identity: commitment.Identity([]byte(path))}}
}

// TestCommitmentOutlookCells pins the reader row: each list cell joins its identities in
// milestone order, and the blocker table carries each reason.
func TestCommitmentOutlookCells(t *testing.T) {
	outlook := commitment.Outlook{State: "eligible", Milestone: "M1", Next: "C", Deliverable: "specs/c/spec.md", Active: []string{"A"}, Blocked: []intent.OutcomeBlocker{{Outcome: "B", Reason: "vendor"}}, Waiting: []string{"D", "E"}, Command: "bench commitment start"}
	if got, want := outlook.Cells(), []string{"eligible", "M1", "C", "specs/c/spec.md", "A", "B", "D E", "bench commitment start"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cells = %q, want %q", got, want)
	}
	if got, want := outlook.BlockerCells(), [][]string{{"B", "vendor"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("blocker cells = %q, want %q", got, want)
	}
}
