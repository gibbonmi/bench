package gate

import (
	"slices"
	"strings"
)

// DecisionPhase is the immutable scheduling value for one gate phase.
type DecisionPhase struct {
	Name  string
	Argv  []string
	Needs []string
}

// DecisionInput contains the accepted subject, resolved oracle, phase table, and reusable evidence.
type DecisionInput struct {
	Subject    string
	Resolution Resolution
	Phases     []DecisionPhase
	Evidence   []string
}

// Decision is the accepted gate plan or a refusal that leaves the input untouched.
type Decision struct {
	Accepted  bool
	Scheduled []string
	Inherited []string
	Refusal   string
}

// Decide validates a gate plan without repository or process access. It selects scheduled
// phases and inherited evidence, or returns a refusal that the execution owner reports.
func Decide(input DecisionInput) Decision {
	if invalidDecisionValue(input.Subject) || input.Subject == "" {
		return Decision{Refusal: "gate subject is required"}
	}
	if input.Resolution.Kind <= None || input.Resolution.Kind > Cargo {
		return Decision{Refusal: "gate resolution is required"}
	}
	if input.Resolution.Kind == BenchGate && strings.TrimSpace(input.Resolution.Command) == "" {
		return Decision{Refusal: "bench gate command is required"}
	}
	phaseNames := make(map[string]bool, len(input.Phases))
	for _, phase := range input.Phases {
		if invalidDecisionValue(phase.Name) || phase.Name == "" || len(phase.Argv) == 0 || phase.Argv[0] == "" || phaseNames[phase.Name] {
			return Decision{Refusal: "invalid gate phase"}
		}
		phaseNames[phase.Name] = true
	}
	inherited := make(map[string]bool, len(input.Evidence))
	for _, name := range input.Evidence {
		if !phaseNames[name] || inherited[name] {
			return Decision{Refusal: "invalid gate evidence"}
		}
		inherited[name] = true
	}
	decision := Decision{Accepted: true}
	for _, phase := range input.Phases {
		for _, need := range phase.Needs {
			if !phaseNames[need] || need == phase.Name {
				return Decision{Refusal: "invalid gate dependency"}
			}
		}
		if inherited[phase.Name] {
			decision.Inherited = append(decision.Inherited, phase.Name)
		} else {
			decision.Scheduled = append(decision.Scheduled, phase.Name)
		}
	}
	return decision
}

func decisionPhases(phases []Phase) []DecisionPhase {
	values := make([]DecisionPhase, len(phases))
	for i, phase := range phases {
		values[i] = DecisionPhase{Name: phase.Name, Argv: slices.Clone(phase.Argv), Needs: slices.Clone(phase.Needs)}
	}
	return values
}

func invalidDecisionValue(value string) bool {
	return strings.TrimSpace(value) != value || strings.ContainsAny(value, "\x00\r\n")
}
