// The gate authority: attribution of an unpublished tree, beside the lane authority.
package authorization

import (
	"bytes"
	"context"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// An infrastructure attribution carries the gate's own open reason, so the refusal that
// reads this result can name the cause instead of the kind alone. A fixture with no gate
// script leaves the subject open, which is exactly the operational outcome the kind
// covers.
func TestAuthorizeCarriesTheInfrastructureReason(t *testing.T) {
	root, tree := laneFixture(t)
	var stdout, stderr bytes.Buffer
	result := AuthorizeWithWriters(context.Background(), root, tree, &stdout, &stderr)
	if result.Kind != Infrastructure {
		t.Fatalf("kind = %q, want %q", result.Kind, Infrastructure)
	}
	if result.Reason == "" {
		t.Fatalf("reason = %q, want the gate's open reason", result.Reason)
	}
	t.Logf("infrastructure reason = %q", result.Reason)
}

// An infrastructure attribution that the gate left without a reason still names what the
// gate observed, so a refusal with the kind alone has a cause on stderr. A gate reason
// needs no second line, because the gate printed that reason itself.
func TestInfrastructureDiagnosticNamesTheSilentOutcomes(t *testing.T) {
	for _, tc := range []struct {
		name       string
		execution  gate.Result
		inspection gate.EvidenceInspection
		want       string
	}{
		{"green without reusable evidence", gate.Result{}, gate.EvidenceInspection{Reason: "no retained verdict"},
			"gate: infrastructure: the gate passed but left no reusable green evidence (no retained verdict)\n"},
		{"stop without a reason", gate.Result{GateExit: 130, ActionExit: 130}, gate.EvidenceInspection{},
			"gate: infrastructure: the gate stopped with gate exit 130 and action exit 130 and recorded no reason\n"},
		{"gate reason", gate.Result{ActionExit: 1, Inspection: gate.Inspection{Reason: "gate lock unavailable"}}, gate.EvidenceInspection{},
			""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := infrastructureDiagnostic(tc.execution, tc.inspection); got != tc.want {
				t.Fatalf("infrastructureDiagnostic = %q, want %q", got, tc.want)
			}
		})
	}
}
