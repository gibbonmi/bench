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
