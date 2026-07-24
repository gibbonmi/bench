package surface

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

// setup_plan_test.go covers FT76 story 2: the plan preview lists every
// inferred fact with its consequence, on stdout, before any write.

func TestSetupPlanContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench setup --plan previews a Go-module repo without writing", testSetupPlanPreviewsGoModule)
	contract.RunParallel(t, "bench setup --plan previews a zero-signal repo without writing", testSetupPlanPreviewsZeroSignal)
	contract.RunParallel(t, "bench setup --plan names an unreadable package.json instead of dropping it silently", testSetupPlanIgnoresUnreadablePackageJSON)
	contract.RunParallel(t, "bench setup --plan names a malformed package.json instead of dropping it silently", testSetupPlanIgnoresMalformedPackageJSON)
	contract.RunParallel(t, "bench setup --plan names an unreadable Makefile instead of dropping it silently", testSetupPlanIgnoresUnreadableMakefile)
}

// skipIfRoot skips a permission-denied fixture on a root-run test process: root can
// read a 0o000 file regardless of its mode bits, so the fixture would not exercise the
// read-error path it is named for.
func skipIfRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "running as root - permission-denied fixtures don't apply")
	}
}

// testSetupPlanIgnoresUnreadablePackageJSON pins C3: "nothing is acted on silently" -
// a package.json bench setup cannot read must still surface as a named preview fact,
// not just vanish from the gate-inference table with no trace.
func testSetupPlanIgnoresUnreadablePackageJSON(t *testing.T) {
	skipIfRoot(t)
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.WriteFile("package.json", `{"scripts":{"test":"echo ok"}}`+"\n")
	pkgPath := filepath.Join(f.Root, "package.json")
	if err := os.Chmod(pkgPath, 0o000); err != nil {
		t.Fatalf("chmod package.json: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(pkgPath, 0o644) })

	probe := f.Bench("setup", "--plan")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "package.json unreadable")
	probe.RequireNotContains(probe.Stdout, "npm test")
}

func testSetupPlanIgnoresMalformedPackageJSON(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.WriteFile("package.json", "{not valid json")

	probe := f.Bench("setup", "--plan")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "package.json malformed JSON")
	probe.RequireNotContains(probe.Stdout, "npm test")
}

func testSetupPlanIgnoresUnreadableMakefile(t *testing.T) {
	skipIfRoot(t)
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.WriteFile("Makefile", "test:\n\techo ok\n")
	makePath := filepath.Join(f.Root, "Makefile")
	if err := os.Chmod(makePath, 0o000); err != nil {
		t.Fatalf("chmod Makefile: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(makePath, 0o644) })

	probe := f.Bench("setup", "--plan")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "Makefile unreadable")
	probe.RequireNotContains(probe.Stdout, "make test")
}

func testSetupPlanPreviewsGoModule(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")

	probe := f.Bench("setup", "--plan")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "go.mod")
	probe.RequireContains(probe.Stdout, "go test ./...")
	probe.RequireContains(probe.Stdout, "AGENTS.md")
	probe.RequireContains(probe.Stdout, "CLAUDE.md")
	if f.Exists("AGENTS.md") || f.Exists(".bench") || f.Exists("CLAUDE.md") {
		t.Fatal("bench setup --plan wrote files for a Go-module fixture")
	}
}

func testSetupPlanPreviewsZeroSignal(t *testing.T) {
	f := contract.NewFixture(t)

	probe := f.Bench("setup", "--plan")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "no build system detected")
	probe.RequireNotContains(probe.Stdout, "go test ./...")
	probe.RequireNotContains(probe.Stdout, "npm test")
	if f.Exists("AGENTS.md") || f.Exists(".bench") || f.Exists("CLAUDE.md") {
		t.Fatal("bench setup --plan wrote files for a zero-signal fixture")
	}
}
