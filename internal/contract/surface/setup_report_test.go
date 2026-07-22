package surface

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// setup_report_test.go covers FT76 story 10: setup ends by running doctor and
// printing that report, the harness reload instruction, and the exact next action,
// with the exit code tracking the doctor verdict — never a fabricated print.

func TestSetupReportContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench setup converged green run prints doctor report, reload, and next action at exit 0", testSetupReportGreen)
	contract.RunParallel(t, "bench setup zero-signal run prints the full report and next action at partial exit", testSetupReportZeroSignalPartial)
}

func testSetupReportGreen(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(0)
	out := probe.Stdout
	probe.RequireContains(out, "ok: AGENTS.md")
	probe.RequireContains(out, "ok: CLAUDE.md")
	probe.RequireContains(out, "ok: gate")
	probe.RequireContains(strings.ToLower(out), "reload")
	probe.RequireContains(strings.ToLower(out), "next:")
	probe.RequireContains(out, "/bench-setup-repo")
}

func testSetupReportZeroSignalPartial(t *testing.T) {
	f := contract.NewFixture(t)

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(3)
	out := probe.Stdout
	probe.RequireContains(out, "red: .bench/gate.sh")
	probe.RequireContains(strings.ToLower(out), "reload")
	probe.RequireContains(strings.ToLower(out), "next:")
	probe.RequireContains(out, "configure")
}
