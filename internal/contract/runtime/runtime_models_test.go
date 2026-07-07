package runtime

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// The keys are pinned empty so discovery is hermetic: the API sources short-circuit to
// unavailable without a network call, and the contracts assert only argv posture.
var modelsHermeticEnv = map[string]string{"OPENAI_API_KEY": "", "ANTHROPIC_API_KEY": ""}

func TestRuntimeModelsContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench models unknown-arg contract", testRuntimeModelsUnknownArg)
	contract.RunParallel(t, "bench models discovery-tolerance contract", testRuntimeModelsDiscoveryTolerance)
}

// An unknown argument is a misuse: reject it with a usage line at exit 2, the sibling
// porcelain norm, instead of silently printing the inventory at exit 0.
func testRuntimeModelsUnknownArg(t *testing.T) {
	f := contract.NewFixture(t)
	probe := f.BenchEnv(modelsHermeticEnv, "models", "bogus")
	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage")
}

// The regression guard: the argv guard must not leak into the no-arg path — a plain
// invocation still emits the inventory tables at exit 0 even when every source is
// unreachable (discovery tolerance is a separate contract).
func testRuntimeModelsDiscoveryTolerance(t *testing.T) {
	f := contract.NewFixture(t)
	probe := f.BenchEnv(modelsHermeticEnv, "models")
	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "model_sources[")
	probe.RequireContains(probe.Stdout, "models[")
}
