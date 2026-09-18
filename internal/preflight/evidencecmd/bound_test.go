package evidencecmd

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// unboundedKinds names every operation kind that produces no evidence response. It is the
// independent half of the guard: the registry supplies the kinds, this list supplies the
// exemption, so a kind registered later carries the bound or reds here by itself.
var unboundedKinds = []Kind{KindVerdict, KindProposal}

// TestEvidenceFinalGuard proves that every evidence form obeys the shared response bound and
// that the bound turns an oversized response into one bounded refusal.
func TestEvidenceFinalGuard(t *testing.T) {
	for _, op := range operations {
		exempt := false
		for _, kind := range unboundedKinds {
			if op.Kind == kind {
				exempt = true
			}
		}
		if op.Bounded == exempt {
			t.Errorf("operation %q declares Bounded=%t, want %t for its kind", op.usageLine(), op.Bounded, !exempt)
		}
	}
	out, code := Bound(strings.Repeat("x", preflighttest.ResponseBudget+1), 0)
	if code != 1 || len(out) > preflighttest.ResponseBudget || !strings.Contains(out, "response bound exceeded") {
		t.Fatalf("guard = (%d, %d bytes)", code, len(out))
	}
}
