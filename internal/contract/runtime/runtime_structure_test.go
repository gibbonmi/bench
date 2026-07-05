package runtime

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/contract"
	"testing"
)

func TestRuntimeStructureContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench structure shell-file contract", testRuntimeStructureShellFile)
	contract.RunParallel(t, "bench structure budgets contract", testRuntimeStructureBudgets)
	contract.RunParallel(t, "bench structure path-with-spaces contract", testRuntimeStructurePathWithSpaces)
}

func testRuntimeStructureShellFile(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("big.sh", repeatLines(401, "x=\n"))
	f.Git("add", "big.sh")
	probe := f.Bench("structure")
	if probe.ExitCode == 0 {
		t.Fatalf("shell source over the line budget did not fail structure\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	probe.RequireContains(probe.Stdout+probe.Stderr, "FILE TOO LONG")
}

func testRuntimeStructureBudgets(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("big.sh", repeatLines(401, "x=\n"))
	f.WriteFile("mid.sh", repeatLines(200, "y=\n"))
	for i := 1; i <= 13; i++ {
		f.WriteFile(fmt.Sprintf("sub/f%d.sh", i), "z=1\n")
	}
	f.WriteFile(".bench/structure.budgets", "# reviewer grants\nbig.sh 500\nsub/ 20\nweird abc\nmid.sh 100")
	f.CommitAll("s")
	probe := f.Bench("structure")
	if probe.ExitCode == 0 {
		t.Fatal("override below the global cap did not fail structure")
	}
	out := probe.Stdout + probe.Stderr
	contract.RequireContains(t, out, "ignoring malformed line")
	contract.RequireNotContains(t, out, "big.sh")
	contract.RequireNotContains(t, out, "DIR CROWDED")
	contract.RequireContains(t, out, "200 lines (max 100)   mid.sh")
}

func testRuntimeStructurePathWithSpaces(t *testing.T) {
	f := contract.NewFixture(t)
	for i := 1; i <= 13; i++ {
		f.WriteFile(fmt.Sprintf("space dir/file%d.sh", i), fmt.Sprintf("x=%d\n", i))
	}
	f.Git("add", "space dir")
	probe := f.BenchEnv(map[string]string{"BENCH_MAX_DIR_FILES": "12"}, "structure")
	if probe.ExitCode == 0 {
		t.Fatalf("crowded path-with-spaces directory did not fail structure\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	out := probe.Stdout + probe.Stderr
	contract.RequireContains(t, out, "space dir/")
	contract.RequireNotContains(t, out, "   ./")
	contract.RequireNotContains(t, out, "   dir/")
}
