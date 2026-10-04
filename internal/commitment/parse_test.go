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
	t.Run("source ownership", func(t *testing.T) {
		p := validPolicy()
		p.Milestones[0].Outcomes[0].Sources = []commitment.SourceBinding{{ID: "FT1", Path: "obligation.md", Identity: commitment.Identity([]byte("body"))}}
		p.Milestones[0].Outcomes = append(p.Milestones[0].Outcomes, commitment.Outcome{ID: "B", Sources: []commitment.SourceBinding{{ID: "FT2", Path: "obligation.md", Identity: commitment.Identity([]byte("body"))}}})
		if err := commitment.Validate(p); err == nil {
			t.Fatal("one source path gained two outcome owners")
		}
	})
}

// A milestone needs outcome criteria before activation: the policy refuses an active
// milestone with no criteria and accepts the same milestone while it stays planned.
func TestCommitmentCriteriaBeforeActivation(t *testing.T) {
	p := validPolicy()
	p.Milestones[0].Outcomes[0].Criteria = nil
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := commitment.Parse(data); err == nil || !strings.Contains(err.Error(), `active milestone "M1" has no criteria`) {
		t.Fatalf("Parse(active milestone without criteria) = %v, want the no-criteria refusal", err)
	}
	p.ActiveMilestone = ""
	if err := commitment.Validate(p); err != nil {
		t.Fatalf("Validate(planned milestone without criteria) = %v, want valid", err)
	}
}

func TestCommitmentExactFieldNames(t *testing.T) {
	data, err := commitment.Bytes(validPolicy())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"version", "milestones", "outcomes", "id"} {
		t.Run(key, func(t *testing.T) {
			input := strings.Replace(string(data), `"`+key+`"`, `"`+strings.ToUpper(key)+`"`, 1)
			if _, err := commitment.Parse([]byte(input)); err == nil {
				t.Fatal("accepted case alias", key)
			}
		})
	}
}

// A parallel grant can name a legacy continuation's run beside one outcome. It still
// names an outcome and two members, and each named run is one distinct valid identity.
func TestCommitmentContinuationGrantShape(t *testing.T) {
	run, other := strings.Repeat("a", 32), strings.Repeat("c", 32)
	for _, row := range []struct {
		name  string
		grant commitment.ParallelGrant
		want  string
	}{
		{name: "beside-a-run", grant: commitment.ParallelGrant{Outcomes: []string{"A"}, Continuations: []string{run}}},
		{name: "no-outcome", grant: commitment.ParallelGrant{Continuations: []string{run, other}}, want: "names no outcome"},
		{name: "one-member", grant: commitment.ParallelGrant{Outcomes: []string{"A"}}, want: "at least two members"},
		{name: "invalid-run", grant: commitment.ParallelGrant{Outcomes: []string{"A"}, Continuations: []string{"not a run"}}, want: `invalid parallel grant continuation "not a run"`},
		{name: "duplicate-run", grant: commitment.ParallelGrant{Outcomes: []string{"A"}, Continuations: []string{run, run}}, want: `invalid parallel grant continuation "` + run + `"`},
	} {
		t.Run(row.name, func(t *testing.T) {
			p := validPolicy()
			p.ParallelGrants = []commitment.ParallelGrant{row.grant}
			err := commitment.Validate(p)
			if row.want == "" && err != nil || row.want != "" && (err == nil || !strings.Contains(err.Error(), row.want)) {
				t.Fatalf("Validate(%+v) = %v, want %q", row.grant, err, row.want)
			}
		})
	}
}

func TestCommitmentMultiOutcomeCycle(t *testing.T) {
	p := policy([]commitment.Milestone{milestone("M1", "A", "B")}, "M1")
	p.Milestones[0].Outcomes[0].Dependencies = []string{"B"}
	p.Milestones[0].Outcomes[1].Dependencies = []string{"A"}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := commitment.Parse(data); err == nil {
		t.Fatal("accepted two-outcome dependency cycle")
	}
}
