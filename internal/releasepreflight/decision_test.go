package releasepreflight

import "testing"

func TestDecidePreflightFromImmutableEvidence(t *testing.T) {
	green := []Result{{Name: "gate", Status: StatusGreen}}
	redFailure := Failure{Kind: "phase", Message: "gate red"}
	for _, tc := range []struct {
		name  string
		input DecisionInput
		check func(*testing.T, Decision)
	}{
		{name: "verify accepted", input: DecisionInput{Mode: ModeVerify, Results: green}, check: func(t *testing.T, got Decision) {
			if !got.Accepted || got.Status != StatusGreen || got.Scope != ScopePreflight || got.Failure != nil {
				t.Fatalf("decision = %#v", got)
			}
		}},
		{name: "focused accepted", input: DecisionInput{Mode: ModeVerify, Focused: PhaseNames(ModeVerify)[0], Results: green}, check: func(t *testing.T, got Decision) {
			if !got.Accepted || got.Scope != ScopeFocused {
				t.Fatalf("decision = %#v", got)
			}
		}},
		{name: "red refused", input: DecisionInput{Mode: ModeVerify, Results: []Result{{Name: "gate", Status: StatusRed, Failure: &redFailure}}}, check: func(t *testing.T, got Decision) {
			if got.Accepted || got.Failure == nil || got.Failure.Message != redFailure.Message {
				t.Fatalf("decision = %#v", got)
			}
		}},
		{name: "empty refused", input: DecisionInput{Mode: ModeVerify}, check: func(t *testing.T, got Decision) {
			if got.Accepted {
				t.Fatalf("decision = %#v", got)
			}
		}},
		{name: "hostile mode refused", input: DecisionInput{Mode: Mode("verify\x1b"), Results: green}, check: func(t *testing.T, got Decision) {
			if got.Accepted || got.Failure == nil || got.Failure.Kind != "usage" {
				t.Fatalf("decision = %#v", got)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := Decide(tc.input)
			second := Decide(tc.input)
			if first.Accepted != second.Accepted || first.Status != second.Status || first.Scope != second.Scope {
				t.Fatalf("rerun changed decision: %#v then %#v", first, second)
			}
			tc.check(t, first)
		})
	}
}
