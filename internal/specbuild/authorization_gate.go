package specbuild

import (
	"context"

	gateauth "github.com/gibbonmi/bench/internal/gate/authorization"
)

// AuthorizationGate adapts the gate-authorization owner to the lifecycle ports.
type AuthorizationGate struct{}

func (AuthorizationGate) Bootstrap(_ context.Context, root, branch, tip, expected string) error {
	return gateauth.Bootstrap(root, branch, tip, expected)
}

func (AuthorizationGate) AdvanceMarker(ctx context.Context, root, branch, destination, expected string) error {
	return gateauth.AdvanceMarker(ctx, root, branch, destination, expected)
}

func (AuthorizationGate) CheckMarker(ctx context.Context, root, branch, destination, expected string) error {
	return gateauth.CheckMarker(ctx, root, branch, destination, expected)
}

func (AuthorizationGate) Execute(ctx context.Context, root, tree string) (GateOutcome, error) {
	result := gateauth.Authorize(ctx, root, tree)
	outcome := GateOutcome{Evidence: result.Evidence}
	switch result.Kind {
	case gateauth.Green:
		outcome.Green = true
	case gateauth.Candidate:
		outcome.Disposition = GateCandidate
	case gateauth.Inherited:
		outcome.Disposition = GateInherited
	default:
		outcome.Disposition = GateInfrastructure
	}
	return outcome, nil
}

func (AuthorizationGate) Validate(_ context.Context, root, tree, evidence string) (bool, error) {
	return gateauth.Validate(root, tree, evidence), nil
}
