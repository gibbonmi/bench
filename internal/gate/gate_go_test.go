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

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/racetests"
	"github.com/gibbonmi/bench/internal/toon"
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

func writeGateGoExecutable(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "go")
	writeGateGoFile(t, path, `#!/bin/sh
if [ "$1" = list ]; then
		printf '%s\n' fixture/core
		exit 0
fi
printf 'GOCACHE=<%s>\n' "${GOCACHE-unset}" >> "$GATE_GO_RECORD"
for argument do
		printf 'arg=<%s>\n' "$argument" >> "$GATE_GO_RECORD"
done
if [ -n "${GOCACHE-}" ]; then
		printf '# test log\n' > "$GOCACHE/fake-test-log"
fi
`)
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("chmod %s: %v", path, err)
	}
	return path
}

func runGateGo(t *testing.T, step, root string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := GateGoCommand([]string{step, root}, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func writeRaceTestSources(t *testing.T, root string) {
	writeRaceTestSourcesFor(t, root, racetests.Tests)
}

func writeRaceTestSourcesFor(t *testing.T, root string, tests []racetests.Test) {
	t.Helper()
	for rel, source := range racetests.SyntheticSourcesFor(tests) {
		writeGateGoFile(t, filepath.Join(root, filepath.FromSlash(rel)), source)
	}
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

func TestCoreTestPackagesIgnoresMalformedAmbientVCSMetadata(t *testing.T) {
	parent := t.TempDir()
	writeGateGoFile(t, filepath.Join(parent, ".git"), "gitdir: missing\n")
	root := filepath.Join(parent, "module")
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(root, "cmd", "fixture", "main.go"), "package main\n\nfunc main() {}\n")

	packages, output, err := CoreTestPackages(root, registry.Dev)
	if err != nil {
		t.Fatalf("CoreTestPackages with malformed ambient VCS metadata: %v\n%s", err, output)
	}
	if len(packages) != 1 || packages[0] != "fixture/cmd/fixture" {
		t.Fatalf("CoreTestPackages = %v, want [fixture/cmd/fixture]", packages)
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

func TestGateGoCoreTestUsesFreshVerdict(t *testing.T) {
	root := t.TempDir()
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoExecutable(t, root)
	record := filepath.Join(root, "go-record")
	cache := filepath.Join(root, "cache")
	if err := os.Mkdir(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("GATE_GO_RECORD", record)
	t.Setenv("GOCACHE", cache)

	if code, stdout, stderr := runGateGo(t, "test", root); code != 0 {
		t.Fatalf("test rc = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if got, want := readGateGoFile(t, record), "GOCACHE=<"+cache+">\narg=<test>\narg=<-count=1>\narg=<fixture/core>\n"; got != want {
		t.Fatalf("core fake Go record = %q, want %q", got, want)
	}
	if !fileExists(filepath.Join(cache, "fake-test-log")) {
		t.Fatal("core step removed the fake Go test-log sentinel")
	}
	if err := os.Remove(record); err != nil {
		t.Fatal(err)
	}
	if err := os.Unsetenv("GOCACHE"); err != nil {
		t.Fatal(err)
	}
	if code, stdout, stderr := runGateGo(t, "test", root); code != 0 {
		t.Fatalf("test without GOCACHE rc = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if got, want := readGateGoFile(t, record), "GOCACHE=<unset>\narg=<test>\narg=<-count=1>\narg=<fixture/core>\n"; got != want {
		t.Fatalf("core fake Go record without GOCACHE = %q, want %q", got, want)
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
	writeGateGoFile(t, filepath.Join(unrun, "internal", "guards", "guards_test.go"),
		"package guards\n\nimport \"testing\"\n\nfunc TestSomethingElse(t *testing.T) {}\n")
	code, stdout, stderr := runGateGo(t, "race", unrun)
	if code != 1 {
		t.Fatalf("race rc = %d when the target test never ran, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	for _, test := range raceTests {
		if !strings.Contains(stderr, "race test did not run: "+test.PackagePath+" "+test.Name) {
			t.Fatalf("missing named-test diagnostic for %s; stdout=%q stderr=%q", test.Name, stdout, stderr)
		}
	}

	// The green side, with the target test emitting output of its own: the step taps
	// stdout for the `=== RUN` line and streams the rest, and driving this case under
	// -race is what proves the tap shares no buffer with the untouched stderr.
	ran := t.TempDir()
	writeGateGoFile(t, filepath.Join(ran, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeRaceTestSources(t, ran)
	code, stdout, stderr = runGateGo(t, "race", ran)
	if code != 0 {
		t.Fatalf("race rc = %d when the target test ran and passed, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	for _, test := range raceTests {
		if !strings.Contains(stdout, "=== RUN   "+test.Name) {
			t.Fatalf("race step did not run %s; stdout=%q stderr=%q", test.Name, stdout, stderr)
		}
	}
	if !strings.Contains(stdout, "race test noise") {
		t.Fatalf("race step swallowed the tool's output; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestGateGoConformanceSuiteUsesRegistrySkipPattern(t *testing.T) {
	kit := kitRootForTest(t)
	want := []string{"go", "test", "-count=1", "./internal/conformance", "-skip", registry.InnerSkipPattern()}
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

func TestGateGoConformanceSuitePreservesCache(t *testing.T) {
	root := t.TempDir()
	pkgDir := filepath.Join(root, "internal", "conformance")
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeGateGoFile(t, filepath.Join(pkgDir, "marker_test.go"), "package conformance\n\nimport \"testing\"\n\nfunc TestRootConformance(t *testing.T) {}\n")
	writeGateGoExecutable(t, root)
	record := filepath.Join(root, "go-record")
	cache := filepath.Join(root, "cache")
	if err := os.Mkdir(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("GATE_GO_RECORD", record)
	t.Setenv("GOCACHE", cache)

	if code, stdout, stderr := runGateGo(t, "conformance-suite", root); code != 0 {
		t.Fatalf("conformance-suite rc = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	want := "GOCACHE=<" + cache + ">\narg=<test>\narg=<-count=1>\narg=<./internal/conformance>\narg=<-skip>\narg=<" + registry.InnerSkipPattern() + ">\n"
	if got := readGateGoFile(t, record); got != want {
		t.Fatalf("conformance fake Go record = %q, want %q", got, want)
	}
	if !fileExists(filepath.Join(cache, "fake-test-log")) {
		t.Fatal("conformance step removed the fake Go test-log sentinel")
	}
}

func TestGateGoArgv(t *testing.T) {
	want := []string{"go", "-C", "/kit", "run", "-buildvcs=false", "./cmd/bench", "gate-go", "gofmt", "/root"}
	if got := GateGoArgv("/kit", "gofmt", "/root"); !reflect.DeepEqual(got, want) {
		t.Fatalf("GateGoArgv = %#v, want %#v", got, want)
	}
	want = []string{"go", "run", "-buildvcs=false", "./cmd/bench", "gate-go", "test", "/root"}
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

// TestGateGoWithoutRootOutsideARepo grades the exit code that separates "this step is
// red" from "this invocation never had a tree to grade". Both argument forms reach it —
// an omitted root and an empty one — because a manifest that interpolates an unset
// variable produces the second, and a caller that reads 1 as a red step would report a
// finding about a tree that was never resolved.
func TestGateGoWithoutRootOutsideARepo(t *testing.T) {
	outside := t.TempDir()
	if _, err := git.Output("-C", outside, "rev-parse", "--show-toplevel"); err == nil {
		capability.Environment(t, "the temp directory sits inside a git repository")
	}
	t.Chdir(outside)

	for name, args := range map[string][]string{
		"omitted root": {"gofmt"},
		"empty root":   {"gofmt", ""},
	} {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := GateGoCommand(args, &stdout, &stderr)
			if code != 3 {
				t.Fatalf("rc = %d outside a repo, want 3; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stderr.String(), toon.NotInRepo()) {
				t.Fatalf("stderr = %q, want the not-in-repo diagnostic", stderr.String())
			}
		})
	}
}

// TestGateGoStepReportsASpawnFailure covers the one red that otherwise carries no
// account of itself: a tool absent from PATH writes to neither stream, so the phase
// reds with empty output and nothing names the cause.
func TestGateGoStepReportsASpawnFailure(t *testing.T) {
	root := t.TempDir()
	writeGateGoFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeRaceTestSources(t, root)
	t.Setenv("PATH", "")

	code, stdout, stderr := runGateGo(t, "race", root)
	if code != 1 {
		t.Fatalf("race rc = %d with go off PATH, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "failed to start") {
		t.Fatalf("a step that never spawned reported nothing; stdout=%q stderr=%q", stdout, stderr)
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

func readGateGoFile(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(contents)
}
