package evidencecmd

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestEvidenceFinalGuard proves that every evidence form obeys the shared response bound and
// that the bound turns an oversized response into one bounded refusal.
func TestEvidenceFinalGuard(t *testing.T) {
	for _, op := range operations {
		if (op.Kind == KindPrepareEvidence || op.Kind == KindReadEvidence) && !op.Bounded {
			t.Errorf("operation %q bypasses the final guard", op.usage)
		}
	}
	out, code := Bound(strings.Repeat("x", preflighttest.ResponseBudget+1), 0)
	if code != 1 || len(out) > preflighttest.ResponseBudget || !strings.Contains(out, "response bound exceeded") {
		t.Fatalf("guard = (%d, %d bytes)", code, len(out))
	}
}
