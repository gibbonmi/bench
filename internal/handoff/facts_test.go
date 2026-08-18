package handoff

import (
	"testing"

	"github.com/gibbonmi/bench/internal/status"
)

func TestApplyRouteUsesNoCommandState(t *testing.T) {
	route := status.Route([]status.Signal{{Name: "gate", Detail: "locked-pending"}}, status.HarnessClaude)
	var f facts
	applyRoute(&f, route)
	if f.Action != "" || f.Signal != "gate" || !f.NoInvocable {
		t.Fatalf("facts = %#v, want route's empty command and no-command state", f)
	}
}

func TestApplyRouteUsesCleanFallback(t *testing.T) {
	var f facts
	applyRoute(&f, status.RouteFor(t.TempDir(), nil, status.HarnessClaude))
	if got := nextField(f); got != "`/bench-drain` — the board's leading invocable signal (`clean`)." {
		t.Fatalf("next field = %q, want clean fallback", got)
	}
}
