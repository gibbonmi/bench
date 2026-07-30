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
	if !strings.Contains(gate, `exec env BENCH_KIT="$kit" "$bench" gate-phases "$root"`) {
		t.Fatalf("benchkit gate does not exec gate-phases:\n%s", gate)
	}
	freshness := `go run ./internal/freshness/check "$kit" "$bench"`
	phase := `exec env BENCH_KIT="$kit" "$bench" gate-phases "$root"`
	if strings.Index(gate, freshness) < 0 || strings.Index(gate, freshness) > strings.Index(gate, phase) {
		t.Fatalf("benchkit gate does not verify %q before %q:\n%s", freshness, phase, gate)
	}
	for _, retired := range []string{
		`bin/bench.sh" canary "$root"`,
	} {
		if strings.Contains(gate, retired) {
			t.Fatalf("benchkit gate still carries retired inline canary orchestration %q:\n%s", retired, gate)
		}
	}
	if strings.Contains(gate, `BENCH_CANARY_INNER`) {
		t.Fatalf("benchkit gate still owns the inner-mode env boundary:\n%s", gate)
	}

	canaryRunner := read(t, kitPath(t, "internal", "canary", "canary.go"))
	if !strings.Contains(canaryRunner, `"BENCH_CANARY_INNER=1"`) {
		t.Fatalf("canary runner no longer marks inner gate runs with BENCH_CANARY_INNER")
	}
	phaseRunner := read(t, kitPath(t, "internal", "gate", "phases.go"))
	if !strings.Contains(phaseRunner, `os.Getenv("BENCH_CANARY_INNER") == "1"`) {
		t.Fatalf("phase runner no longer owns the BENCH_CANARY_INNER mode branch")
	}
}
