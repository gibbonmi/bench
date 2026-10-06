package testreport

import (
	"context"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/toon"
)

const proseCheckName = "prose"

func unknownCheck(check string) string {
	return "unknown check: " + check + "\n" + namedCheckInventory()
}

func namedCheckInventory() string {
	checks := namedChecks()
	return "checks:\n  " + strings.Join(checks, "\n  ")
}

func namedChecks() []string {
	return append(registry.Names(registry.Dev), gate.SystemPhaseName, proseCheckName)
}

func isNamedCheck(check string) bool {
	for _, name := range namedChecks() {
		if name == check {
			return true
		}
	}
	return false
}

func runNamedCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (Outcome, string, int) {
	if request.check == gate.SystemPhaseName {
		return runSystemCheck(ctx, root, request, selection)
	}
	argv := focusedTestArgv("./internal/conformance", "-run", namedCheckRunPattern())
	env, err := conformanceEnvironment(os.Environ(), root, request.check, selection)
	if err != nil {
		return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
	}
	return runGoTest(ctx, selection.SourceRoot, request, argv, env)
}

// runProseCheck grades sentences rather than a Go test, so its outcome reads each
// finding as a failure row and an empty grade as a pass.
func runProseCheck(root string) (Outcome, string, int) {
	findings := prose.Grade(root)
	if len(findings) == 0 {
		return Outcome{Kind: OutcomePassed}, "", 0
	}
	return Outcome{Kind: OutcomeFailed, FailedTests: len(findings)}, strings.Join(findings, "\n") + "\n", 1
}

// runSystemCheck runs the gate's system phase as a focused run. It reads the phase's
// operands and environment from the gate's producer, and it sets no conformance
// variable, because the system suite is a build-tagged package rather than a
// conformance scope.
func runSystemCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (Outcome, string, int) {
	operands, suiteEnv := gate.SystemSuite(root)
	env, err := selectedRunEnvironment(os.Environ(), selection)
	if err != nil {
		return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
	}
	return runGoTest(ctx, root, request, focusedTestArgv(operands...), append(env, suiteEnv...))
}
