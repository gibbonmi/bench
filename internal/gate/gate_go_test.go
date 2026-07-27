package gate

// The bench gate-go toolchain steps: each is driven as the subcommand against a temp
// tree, asserting exit code and diagnostic, because the three fail-open shapes these
// steps exist to close (gofmt -l exiting 0, an empty package set, a -run filter that
// matches nothing) all present as green.

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
)

func kitRootForTest(t *testing.T) string {
	t.Helper()
	root, err := git.Root()
	if err != nil {
		t.Fatalf("resolve kit root: %v", err)
	}
	return root
}

func writeGateGoFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func runGateGo(t *testing.T, step, root string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := GateGoCommand([]string{step, root}, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestGateGoGofmt(t *testing.T) {
	root := t.TempDir()
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(root, "bad.go"), "package fixture\nfunc Bad()  {\nreturn\n}\n")

	code, stdout, stderr := runGateGo(t, "gofmt", root)
	if code != 1 {
		t.Fatalf("gofmt rc = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout+stderr, "gofmt: unformatted Go files: bad.go") {
		t.Fatalf("gofmt diagnostic missing the label and the file; stdout=%q stderr=%q", stdout, stderr)
	}

	writeGateGoFile(t, filepath.Join(root, "bad.go"), "package fixture\n\nfunc Bad() {}\n")
	code, stdout, stderr = runGateGo(t, "gofmt", root)
	if code != 0 {
		t.Fatalf("gofmt rc = %d for a formatted tree, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestGateGoTestPackageSet(t *testing.T) {
	kit := kitRootForTest(t)
	dev, _, err := CoreTestPackages(kit, registry.Dev)
	if err != nil {
		t.Fatalf("CoreTestPackages at dev: %v", err)
	}
	ship, _, err := CoreTestPackages(kit, registry.Ship)
	if err != nil {
		t.Fatalf("CoreTestPackages at ship: %v", err)
	}

	// internal/contract is the one core-adjacent consumer of dist/bench; its absence is
	// what makes the test phase safe without a build edge.
	for _, excluded := range []string{"/internal/contract", "/internal/contract/axi", "/internal/conformance"} {
		if hasPackageSuffix(dev, excluded) {
			t.Fatalf("dev package set included %s:\n%s", excluded, strings.Join(dev, "\n"))
		}
	}
	if !hasPackageSuffix(dev, "/internal/conformance/registry") {
		t.Fatalf("dev package set excluded the conformance registry leaf, which no filtered run grades:\n%s", strings.Join(dev, "\n"))
	}
	for _, releaseOnly := range registry.ReleaseOnlyPackages {
		if hasPackageSuffix(dev, "/"+releaseOnly) {
			t.Fatalf("dev package set included release-only %s:\n%s", releaseOnly, strings.Join(dev, "\n"))
		}
		if !hasPackageSuffix(ship, "/"+releaseOnly) {
			t.Fatalf("ship package set excluded release-only %s:\n%s", releaseOnly, strings.Join(ship, "\n"))
		}
	}
}

func TestGateGoTestReds(t *testing.T) {
	root := t.TempDir()
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(root, "core", "core_test.go"),
		"package core\n\nimport \"testing\"\n\nfunc TestCoreFails(t *testing.T) { t.Fatal(\"boom\") }\n")

	code, stdout, stderr := runGateGo(t, "test", root)
	if code != 1 {
		t.Fatalf("test rc = %d for a tree whose test fails, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout+stderr, "TestCoreFails") {
		t.Fatalf("test step did not stream the tool's own output; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestGateGoRaceRequiresTheTestToRun(t *testing.T) {
	absent := t.TempDir()
	writeGateGoFile(t, filepath.Join(absent, "go.mod"), "module fixture\n\ngo 1.25\n")
	if code, stdout, stderr := runGateGo(t, "race", absent); code != 1 {
		t.Fatalf("race rc = %d for a tree with no internal/worktree, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}

	// A -run filter that matches nothing exits 0, so the package existing and compiling
	// is not evidence the target test executed.
	unrun := t.TempDir()
	writeGateGoFile(t, filepath.Join(unrun, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(unrun, "internal", "worktree", "worktree_test.go"),
		"package worktree\n\nimport \"testing\"\n\nfunc TestSomethingElse(t *testing.T) {}\n")
	code, stdout, stderr := runGateGo(t, "race", unrun)
	if code != 1 {
		t.Fatalf("race rc = %d when the target test never ran, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}

	// The green side, with the target test emitting output of its own: the step taps
	// stdout for the `=== RUN` line and streams the rest, and driving this case under
	// -race is what proves the tap shares no buffer with the untouched stderr.
	ran := t.TempDir()
	writeGateGoFile(t, filepath.Join(ran, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(ran, "internal", "worktree", "worktree_test.go"),
		"package worktree\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"testing\"\n)\n\nfunc TestConcurrentCleanupRecordsOneTransaction(t *testing.T) {\n\tfmt.Fprintln(os.Stderr, \"cleanup noise\")\n}\n")
	code, stdout, stderr = runGateGo(t, "race", ran)
	if code != 0 {
		t.Fatalf("race rc = %d when the target test ran and passed, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "=== RUN   "+cleanupRaceTest) || !strings.Contains(stdout, "cleanup noise") {
		t.Fatalf("race step swallowed the tool's output; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestGateGoConformanceSuiteUsesRegistrySkipPattern(t *testing.T) {
	kit := kitRootForTest(t)
	want := []string{"go", "test", "./internal/conformance", "-skip", registry.InnerSkipPattern()}
	if got := ConformanceSuiteArgv(kit); !reflect.DeepEqual(got, want) {
		t.Fatalf("ConformanceSuiteArgv = %#v, want %#v", got, want)
	}
	if got := ConformanceSuiteArgv(t.TempDir()); got != nil {
		t.Fatalf("ConformanceSuiteArgv = %#v for a root with no conformance package, want none", got)
	}

	// The fixture's conformance package records which of its tests ran: the suite
	// member's marker proves the filtered run happened at all, and the entry point's
	// absent marker proves the skip pattern reached it.
	root := t.TempDir()
	pkgDir := filepath.Join(root, "internal", "conformance")
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(pkgDir, "marker_test.go"), `package conformance

import (
	"os"
	"testing"
)

func TestRootConformance(t *testing.T) {
	if err := os.WriteFile("ran-entry-point", nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFixtureSuiteMember(t *testing.T) {
	if err := os.WriteFile("ran-suite-member", nil, 0o644); err != nil {
		t.Fatal(err)
	}
}
`)

	code, stdout, stderr := runGateGo(t, "conformance-suite", root)
	if code != 0 {
		t.Fatalf("conformance-suite rc = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !fileExists(filepath.Join(pkgDir, "ran-suite-member")) {
		t.Fatal("the conformance package's own suite did not run; the filtered invocation is what keeps it in the oracle")
	}
	if fileExists(filepath.Join(pkgDir, "ran-entry-point")) {
		t.Fatal("the filtered run executed the entry-point test, so the inner run recurses into the outer one")
	}
}

func TestGateGoArgv(t *testing.T) {
	want := []string{"go", "-C", "/kit", "run", "./cmd/bench", "gate-go", "gofmt", "/root"}
	if got := GateGoArgv("/kit", "gofmt", "/root"); !reflect.DeepEqual(got, want) {
		t.Fatalf("GateGoArgv = %#v, want %#v", got, want)
	}
	want = []string{"go", "run", "./cmd/bench", "gate-go", "test", "/root"}
	if got := GateGoArgv("", "test", "/root"); !reflect.DeepEqual(got, want) {
		t.Fatalf("GateGoArgv with no kit = %#v, want %#v", got, want)
	}
}

func TestGateGoUnknownStep(t *testing.T) {
	root := t.TempDir()
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")

	code, stdout, stderr := runGateGo(t, "gofmtt", root)
	if code != 2 {
		t.Fatalf("unknown step rc = %d, want 2; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "usage: bench gate-go") {
		t.Fatalf("unknown step stderr = %q, want a usage diagnostic", stderr)
	}

	var out, errOut bytes.Buffer
	if code := GateGoCommand([]string{"gofmt", root, "extra"}, &out, &errOut); code != 2 {
		t.Fatalf("too many args rc = %d, want 2; stderr=%q", code, errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := GateGoCommand(nil, &out, &errOut); code != 2 {
		t.Fatalf("no step rc = %d, want 2; stderr=%q", code, errOut.String())
	}
}

func TestGateGoSpacedRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "spaced root")
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(root, "bad.go"), "package fixture\nfunc Bad()  {\nreturn\n}\n")

	code, stdout, stderr := runGateGo(t, "gofmt", root)
	if code != 1 {
		t.Fatalf("gofmt rc = %d for a spaced root, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout+stderr, "bad.go") {
		t.Fatalf("gofmt graded the wrong tree for a spaced root; stdout=%q stderr=%q", stdout, stderr)
	}
}

func hasPackageSuffix(packages []string, suffix string) bool {
	for _, pkg := range packages {
		if strings.HasSuffix(pkg, suffix) {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
