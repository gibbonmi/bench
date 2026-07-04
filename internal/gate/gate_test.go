package gate

import (
	"strings"
	"testing"
)

// TestResolvePrecedence pins the ordered chain as a pure function: `.bench/gate.sh`
// beats `$BENCH_GATE` beats the auto-detect lockfiles in their fixed order, and an
// absent gate resolves to None. A reordered chain would silently run the wrong oracle;
// no black-box assertion pins this cheaply, so the table is the genuinely-red signal.
func TestResolvePrecedence(t *testing.T) {
	// present names the set of paths the injected probes report as existing/executable.
	fs := func(present ...string) FS {
		set := map[string]bool{}
		for _, p := range present {
			set[p] = true
		}
		has := func(p string) bool { return set[p] }
		return FS{Executable: has, Exists: has}
	}
	const gateSh = "/r/.bench/gate.sh"

	cases := []struct {
		name      string
		benchGate string
		fs        FS
		want      Kind
	}{
		{"gate.sh beats BENCH_GATE and lockfiles", "echo hi", fs(gateSh, "/r/package.json"), GateSh},
		{"BENCH_GATE beats lockfiles", "echo hi", fs("/r/package.json"), BenchGate},
		{"pnpm beats npm", "", fs("/r/pnpm-lock.yaml", "/r/package.json"), Pnpm},
		{"package.json picks npm", "", fs("/r/package.json"), Npm},
		{"pyproject after npm", "", fs("/r/pyproject.toml"), Pyproject},
		{"cargo last", "", fs("/r/Cargo.toml"), Cargo},
		{"nothing resolves to None", "", fs(), None},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Resolve("/r", tc.benchGate, tc.fs)
			if got.Kind != tc.want {
				t.Errorf("Resolve = %v, want %v", got.Kind, tc.want)
			}
			if tc.want == BenchGate && got.Command != tc.benchGate {
				t.Errorf("BenchGate command = %q, want %q", got.Command, tc.benchGate)
			}
		})
	}
}

func TestGateEnvStripsWrapperRoutingInternals(t *testing.T) {
	t.Setenv("BENCH_KIT", "/wrong/kit")
	t.Setenv("BENCH_WRAPPER", "/wrong/wrapper")
	t.Setenv("BENCH_GATE", "echo ok")

	env := gateEnv()
	sawBenchGate := false
	for _, kv := range env {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") {
			t.Fatalf("gateEnv leaked wrapper-routing internal %q", kv)
		}
		if kv == "BENCH_GATE=echo ok" {
			sawBenchGate = true
		}
	}
	if !sawBenchGate {
		t.Fatal("gateEnv stripped BENCH_GATE, which is part of the project gate contract")
	}
}
