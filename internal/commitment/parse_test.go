package commitment_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
)

func validPolicy() commitment.Policy {
	return commitment.Policy{
		Version:         1,
		ActiveMilestone: "M1",
		Milestones: []commitment.Milestone{{
			ID: "M1",
			Outcomes: []commitment.Outcome{{
				ID:       "A",
				Criteria: []commitment.Criterion{{ID: "done", Text: "Outcome A is delivered"}},
			}},
		}},
	}
}

func TestCommitmentEmptyPolicy(t *testing.T) {
	_, err := commitment.Parse(nil)
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("Parse(nil) error = %v, want empty-policy refusal", err)
	}
}

func TestCommitmentMalformedPolicy(t *testing.T) {
	valid, err := json.Marshal(validPolicy())
	if err != nil {
		t.Fatal(err)
	}
	duplicateField := strings.Replace(string(valid), `"version":1`, `"version":1,"version":1`, 1)
	unknownField := strings.Replace(string(valid), `"version":1`, `"version":1,"surprise":true`, 1)
	duplicateID := strings.Replace(string(valid), `"milestones":[`, `"milestones":[{"id":"M1","outcomes":[{"id":"B","criteria":[],"sources":[]}]} ,`, 1)
	cyclePolicy := validPolicy()
	cyclePolicy.Milestones[0].Outcomes[0].Dependencies = []string{"A"}
	cycle, err := json.Marshal(cyclePolicy)
	if err != nil {
		t.Fatal(err)
	}

	for name, input := range map[string]string{
		"duplicate field": duplicateField,
		"unknown field":   unknownField,
		"duplicate ID":    duplicateID,
		"cycle":           string(cycle),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := commitment.Parse([]byte(input)); err == nil {
				t.Fatalf("Parse() accepted %s input", name)
			}
		})
	}
}

func TestCommitmentPolicyAuthorityInvariants(t *testing.T) {
	t.Run("active criteria", func(t *testing.T) {
		p := validPolicy()
		p.Milestones[0].Outcomes[0].Criteria = nil
		if err := commitment.Validate(p); err == nil {
			t.Fatal("active milestone without criteria was accepted")
		}
	})
	t.Run("source ownership", func(t *testing.T) {
		p := validPolicy()
		p.Milestones[0].Outcomes[0].Sources = []commitment.SourceBinding{{ID: "FT1", Path: "obligation.md", Identity: commitment.Identity([]byte("body"))}}
		p.Milestones[0].Outcomes = append(p.Milestones[0].Outcomes, commitment.Outcome{ID: "B", Sources: []commitment.SourceBinding{{ID: "FT2", Path: "obligation.md", Identity: commitment.Identity([]byte("body"))}}})
		if err := commitment.Validate(p); err == nil {
			t.Fatal("one source path gained two outcome owners")
		}
	})
}
