package gate

import (
	"reflect"
	"testing"
)

func TestDecideGateFromImmutableValues(t *testing.T) {
	phases := []DecisionPhase{
		{Name: "test", Argv: []string{"go", "test", "./..."}},
		{Name: "system", Argv: []string{"go", "test", "-tags=system", "./internal/systemtest"}, Needs: []string{"test"}},
	}
	tests := []struct {
		name  string
		input DecisionInput
		want  Decision
	}{
		{name: "accepted", input: DecisionInput{Subject: "tree", Resolution: Resolution{Kind: GateSh}, Phases: phases}, want: Decision{Accepted: true, Scheduled: []string{"test", "system"}}},
		{name: "evidence", input: DecisionInput{Subject: "tree", Resolution: Resolution{Kind: GateSh}, Phases: phases, Evidence: []string{"test"}}, want: Decision{Accepted: true, Scheduled: []string{"system"}, Inherited: []string{"test"}}},
		{name: "external", input: DecisionInput{Subject: "tree", Resolution: Resolution{Kind: BenchGate, Command: "make verify"}}, want: Decision{Accepted: true}},
		{name: "empty", input: DecisionInput{}, want: Decision{Refusal: "gate subject is required"}},
		{name: "hostile", input: DecisionInput{Subject: "tree\nother", Resolution: Resolution{Kind: GateSh}}, want: Decision{Refusal: "gate subject is required"}},
		{name: "error", input: DecisionInput{Subject: "tree", Resolution: Resolution{Kind: Cargo + 1}}, want: Decision{Refusal: "gate resolution is required"}},
		{name: "missing command", input: DecisionInput{Subject: "tree", Resolution: Resolution{Kind: BenchGate}}, want: Decision{Refusal: "bench gate command is required"}},
		{name: "unknown evidence", input: DecisionInput{Subject: "tree", Resolution: Resolution{Kind: GateSh}, Phases: phases, Evidence: []string{"race"}}, want: Decision{Refusal: "invalid gate evidence"}},
		{name: "bad dependency", input: DecisionInput{Subject: "tree", Resolution: Resolution{Kind: GateSh}, Phases: []DecisionPhase{{Name: "system", Argv: []string{"go"}, Needs: []string{"test"}}}}, want: Decision{Refusal: "invalid gate dependency"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Decide(test.input)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Decide() = %#v, want %#v", got, test.want)
			}
			if rerun := Decide(test.input); !reflect.DeepEqual(rerun, got) {
				t.Fatalf("rerun = %#v, want %#v", rerun, got)
			}
		})
	}
}
