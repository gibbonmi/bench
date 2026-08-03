package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeOuterGateIgnoresAmbientConformanceCheck(t *testing.T) {
	for _, check := range []string{"line-routing", "no-such-check", "release-evidence-probe"} {
		t.Run(check, func(t *testing.T) {
			probe := runConformanceSelection(t, check, "", false)
			probe.RequireExit(0)
		})
	}
}

func TestRuntimeInnerGateRetainsCanaryConformanceCheck(t *testing.T) {
	const check = "line-routing"
	probe := runConformanceSelection(t, check, check, true)
	probe.RequireExit(0)
}

func runConformanceSelection(t *testing.T, check, want string, inner bool) contract.Probe {
	t.Helper()
	kit := t.TempDir()
	contract.WriteFileAbs(t, filepath.Join(kit, "go.mod"), "module benchconformanceselectionfixture\n\n"+subjectGoDirective(t)+"\n")
	contract.WriteFileAbs(t, filepath.Join(kit, "internal", "conformance", "root_test.go"), fmt.Sprintf(`package conformance

import (
	"os"
	"testing"
)

func TestRootConformance(t *testing.T) {
	if got := os.Getenv(%q); got != %q {
		t.Fatalf("conformance selector = %%q, want %%q", got, %q)
	}
}
`, registry.ConformanceCheckEnv, want, want))
	contract.WriteFileAbs(t, filepath.Join(kit, "internal", "contract", "probe", "probe_test.go"), "package probe\n\nimport \"testing\"\n\nfunc TestProbe(t *testing.T) {}\n")
	contract.WriteExecutableAbs(t, filepath.Join(kit, "bin", "bench.sh"), "#!/usr/bin/env bash\nexit 0\n")

	f := contract.NewExecFixtureAt(t, kit)
	graded := t.TempDir()
	run := contract.Env{
		"BENCH_KIT":                  &kit,
		registry.ConformanceCheckEnv: &check,
	}
	if inner {
		marker, phase := "1", "conformance"
		run["BENCH_CANARY_INNER"] = &marker
		run["BENCH_CANARY_PHASE"] = &phase
	}
	return f.RunEnvSpec(run, runtimeConformanceSelectionBench(t), "gate-phases", graded)
}

func runtimeConformanceSelectionBench(t *testing.T) string {
	t.Helper()
	if bench := os.Getenv("BENCH_RUNTIME_CONFORMANCE_SELECTION_BINARY"); bench != "" {
		return bench
	}
	contract.RequireFreshBench(t)
	return filepath.Join(contract.SubjectRoot(t), "dist", "bench")
}
