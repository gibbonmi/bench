package canary

import (
	"strings"
	"testing"
)

func TestBenchkitGateCallsCanarySubcommandOnlyInOuterMode(t *testing.T) {
	gate := read(t, kitPath(t, ".bench", "gate.sh"))
	if strings.Contains(gate, "lib/canary-run.sh") {
		t.Fatalf("benchkit gate still sources canary-run.sh")
	}
	if !strings.Contains(gate, `BENCH_CANARY_INNER`) || !strings.Contains(gate, `bin/bench.sh" canary "$root"`) {
		t.Fatalf("benchkit gate does not guard an outer-mode bench canary call:\n%s", gate)
	}
}
