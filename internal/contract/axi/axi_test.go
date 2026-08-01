package axi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestAXIQuerySurfaceContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI learnings two-entry contract", testAXILearningsTwoEntry)
	contract.RunParallel(t, "AXI learnings empty/template contract", testAXILearningsEmptyTemplate)
	contract.RunParallel(t, "AXI TOON field-escaping contract", testAXITOONFieldEscaping)
	contract.RunParallel(t, "AXI usage/exit-2 contract", testAXIUsageExit2)
	contract.RunParallel(t, "AXI subdirectory root-resolution contract", testAXISubdirectoryRootResolution)
	contract.RunParallel(t, "AXI path-with-spaces contract", testAXIPathWithSpaces)
	contract.RunParallel(t, "AXI learnings ascii-separator title contract", testAXILearningsASCIISeparatorTitle)
}

func testAXILearningsTwoEntry(t *testing.T) {
	contract.NoteContractFailure(t, "AXI learnings two-entry contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("capture/learnings.md", "## 2026-01-01 — first learning  [open]\n- body\n## 2026-02-02 — second learning  [open]\n- body\n")
	out := f.Bench("learnings")
	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "learnings[2]{date,title}:")
	requireIndentedRowCount(t, out.Stdout, 2)
	requireAXILine(t, out.Stdout, "  2026-01-01,first learning")
	requireAXILine(t, out.Stdout, "  2026-02-02,second learning")
}

func testAXILearningsEmptyTemplate(t *testing.T) {
	f := contract.NewFixture(t)
	for _, body := range []string{"", "## <date> — <short title>  [open]\n", "# Learnings — usage journal\n"} {
		if body != "" {
			f.WriteFile("capture/learnings.md", body)
		}
		out := f.Bench("learnings")
		out.RequireExit(0)
		requireAXILine(t, out.Stdout, "learnings[0]{date,title}:")
	}
}

func testAXITOONFieldEscaping(t *testing.T) {
	contract.NoteContractFailure(t, "AXI TOON field-escaping contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("capture/learnings.md", "## 2026-03-03 — a, \"b\"  [open]\n")
	out := f.Bench("learnings")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, `  2026-03-03,"a, \"b\""`)
}

func testAXIUsageExit2(t *testing.T) {
	f := contract.NewFixture(t)
	for _, command := range []string{"learnings", "maps"} {
		out := f.Bench(command, "bogusarg")
		out.RequireExit(2)
		requireContainsFold(t, out.Stdout, "usage")
	}
}

func testAXISubdirectoryRootResolution(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("capture/learnings.md", "## 2026-04-04 — sub check  [open]\n")
	subdir := filepath.Join(f.Root, "sub", "deeper")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	out := runBenchInDir(t, f, subdir, "learnings")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  2026-04-04,sub check")
}

func testAXIPathWithSpaces(t *testing.T) {
	f := contract.NewFixture(t, contract.WithSpacePath())
	f.WriteFile("capture/learnings.md", "## 2026-05-05 — spaced  [open]\n")
	f.WriteFile("decisions/s.md", activeMapDocument("shaping", "", "— (open)", "— (open)", "— (deferred)"))
	learnings := f.Bench("learnings")
	learnings.RequireExit(0)
	requireAXILine(t, learnings.Stdout, "  2026-05-05,spaced")
	maps := f.Bench("maps")
	maps.RequireExit(0)
	requireAXILine(t, maps.Stdout, "  s,First,Research,frontier,\"\"")
}

func testAXILearningsASCIISeparatorTitle(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("capture/learnings.md", "## 2026-01-01 - ascii title  [open]\n")
	out := f.Bench("learnings")
	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  2026-01-01,ascii title")
}
