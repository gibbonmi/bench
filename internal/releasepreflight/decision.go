package releasepreflight

// DecisionInput is the immutable release evidence presented to preflight policy.
type DecisionInput struct {
	Mode    Mode
	Profile Profile
	Focused string
	Results []Result
}

// Decision is the branch-native preflight verdict and evidence scope.
type Decision struct {
	Accepted bool
	Status   Status
	Scope    Scope
	Failure  *Failure
}

// Decide validates a preflight request and classifies its completed evidence.
func Decide(input DecisionInput) Decision {
	if input.Mode != ModeVerify && input.Mode != ModePublish {
		return Decision{Failure: usageFailure()}
	}
	if input.Profile != "" && input.Profile != ProfilePublic && input.Profile != ProfileBank {
		return Decision{Failure: usageFailure()}
	}
	if input.Mode == ModePublish && input.Profile == "" {
		return Decision{Failure: usageFailure()}
	}
	if input.Focused != "" && !contains(PhaseNames(input.Mode), input.Focused) {
		return Decision{Failure: usageFailure()}
	}
	if len(input.Results) == 0 {
		return Decision{Scope: scopeFor(input.Focused), Failure: &Failure{Kind: "evidence", Message: "preflight evidence is empty"}}
	}
	decision := Decision{Status: terminalStatus(input.Results), Scope: scopeFor(input.Focused)}
	decision.Accepted = decision.Status == StatusGreen
	if !decision.Accepted {
		for _, result := range input.Results {
			if result.Failure != nil {
				failure := *result.Failure
				decision.Failure = &failure
				break
			}
		}
	}
	return decision
}
