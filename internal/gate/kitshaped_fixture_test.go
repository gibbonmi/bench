package gate

// The kit-shaped fixture root. The synthetic root in reduced_run_test.go declares two
// phases over a tree with no Go module, so nothing about the build phase, the freshness
// seal, or the canary surfaces can be observed against it. This root carries the tree
// shape BenchkitPhases reads — a module with a ./cmd/bench main, the build helper beside
// its auxiliary input manifest, the wrapper script the canary phase execs, and the canary
// source and fixture directories — and seals a published dist/bench, so those components
// have somewhere to be graded. Both roots stand; neither replaces the other.
//
// Execution stays observable the same way: the resolved gate script appends to
// .git/full-runs and every phase appends its own name to .git/phase-runs, so the executed
// set is read from a durable marker rather than from a return value.

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	kitpayload "github.com/gibbonmi/bench"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	// The package's own freshness constant owns the bare name here, so the seal package
	// is reached through an alias.
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/racetests"
)

// fixtureModulePath is the fixture module's import prefix, so a package of its own can be
// named apart from the standard library entries both closures carry.
const fixtureModulePath = "benchfixture"

// outsideBinaryClosurePackage is the fixture package ./cmd/bench never imports. Its test
// file is what carries it into the module-wide `go list -deps -test ./...` closure, which
// is the gap the toolchain components' input declaration has to cover: a derivation taken
// from the binary's closure alone is blind to every package like this one.
const outsideBinaryClosurePackage = fixtureModulePath + "/internal/canary"

// fixtureGateScript is the resolved gate. Two things about it are load-bearing and pull in
// different directions, so both are stated.
//
// The closing exec is the `gate-phases` hand-off phaseTableGate reads to decide whether the
// root may be offered a reduced run at all; a stand-in binary keeps it inert, exactly as
// the synthetic root's does. The loop above it is what stands in for the phase table on the
// resolved path, and it discovers the phase scripts rather than listing them — so the
// manifest stays the only place the fixture's table is written down, and a manifest phase
// whose script is missing shows up as an absent marker instead of as a silently shorter run.
const fixtureGateScript = `#!/usr/bin/env bash
set -uo pipefail
echo full >> .git/full-runs
for script in .bench/phase-*.sh; do
  bash "$script" || exit 1
done
exec true gate-phases "$PWD"
`

// kitShapedFixture is a constructed root paired with the phase table that root resolves.
// Tests take their expected phase set from phases; a second literal list would disagree
// with the table the run is actually made of.
type kitShapedFixture struct {
	root   string
	phases []Phase
}

var kitShapedTemplateState struct {
	once sync.Once
	path string
	dir  string
	err  error
}

func TestMain(m *testing.M) {
	if err := provideDefaultGoCache(); err != nil {
		fmt.Fprintln(os.Stderr, "resolve the default Go build cache:", err)
		os.Exit(1)
	}
	code := m.Run()
	if err := os.RemoveAll(kitShapedTemplateState.dir); err != nil && code == 0 {
		fmt.Fprintln(os.Stderr, "remove kit-shaped fixture template:", err)
		code = 1
	}
	os.Exit(code)
}

func provideDefaultGoCache() error {
	if os.Getenv("GOCACHE") != "" {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return nil
	}
	output, err := exec.Command("go", "env", "GOCACHE").Output()
	if err != nil {
		return err
	}
	return os.Setenv("GOCACHE", strings.TrimSpace(string(output)))
}

// newKitShapedFixture builds the root and returns it with its resolved table.
func newKitShapedFixture(t *testing.T) kitShapedFixture {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		capability.Environment(t, "go toolchain absent; the kit-shaped fixture needs a real module")
	}
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeKitShapedTree(t, root)
	sealInitialKitShapedBinary(t, root)
	// The manifest is generated from BenchkitPhases' own answer for this tree, so the
	// executed table carries the names the kit table materializes here and nothing else:
	// a tree that stops satisfying a phase's shape stops declaring that phase.
	writeKitShapedManifest(t, root, BenchkitPhases(root, root))
	phases, err := phaseTable(root, root)
	if err != nil {
		t.Fatalf("resolve the fixture phase table: %v", err)
	}
	return kitShapedFixture{root: root, phases: phases}
}

func kitShapedTemplate(root string) (string, error) {
	kitShapedTemplateState.once.Do(func() {
		kitShapedTemplateState.dir, kitShapedTemplateState.err = os.MkdirTemp("", "bench-kitshaped-fixture-")
		if kitShapedTemplateState.err != nil {
			return
		}
		kitShapedTemplateState.path = filepath.Join(kitShapedTemplateState.dir, "bench.staged")
		kitShapedTemplateState.err = buildFixtureBinary(root, "./cmd/bench", kitShapedTemplateState.path)
	})
	return kitShapedTemplateState.path, kitShapedTemplateState.err
}

func (f kitShapedFixture) phaseNames() []string {
	names := make([]string, 0, len(f.phases))
	for _, phase := range f.phases {
		names = append(names, phase.Name)
	}
	return names
}

func (f kitShapedFixture) binaryPath() string { return filepath.Join(f.root, "dist", "bench") }

func writeKitShapedTree(t *testing.T, root string) {
	t.Helper()
	// dist/ stays out of the tree identity, as it does in the kit: a gate run that
	// republishes the binary would otherwise move the subject it was grading.
	writeGateTestFile(t, root, ".gitignore", "dist/\n", 0o644)
	writeGateTestFile(t, root, "go.mod", "module "+fixtureModulePath+"\n\ngo 1.21\n", 0o644)
	writeGateTestFile(t, root, "cmd/bench/main.go", "package main\n\nfunc main() {}\n", 0o644)
	writeGateTestFile(t, root, "internal/canary/canary.go",
		"package canary\n\n// Name is the surface this package's own test grades.\nfunc Name() string { return \"canary\" }\n", 0o644)
	writeGateTestFile(t, root, "internal/canary/canary_test.go",
		"package canary\n\nimport \"testing\"\n\nfunc TestName(t *testing.T) {\n\tif Name() != \"canary\" {\n\t\tt.Fatal(\"canary name moved\")\n\t}\n}\n", 0o644)
	writeGateTestFile(t, root, "internal/conformance/root_test.go",
		"package conformance\n\nimport \"testing\"\n\nfunc TestRootConformance(t *testing.T) {}\n", 0o644)
	writeConformanceCanaryOwners(t, root)
	writeGateTestFile(t, root, "tests/canary/fixture.txt", "canary fixture\n", 0o644)
	// BenchkitPhases materializes the build phase only when the build helper and go.mod
	// are both regular files, and the seal digest refuses a root whose auxiliary manifest is
	// absent or empty. Neither script is ever executed here — the phase manifest routes every
	// phase through a marker script instead — so both carry the inert body.
	writeGateTestFile(t, root, "scripts/go-build.sh", "#!/usr/bin/env bash\nexit 0\n", 0o755)
	writeGateTestFile(t, root, "scripts/go-build.inputs", "build_script=scripts/go-build.sh\n", 0o644)
	writeGateTestFile(t, root, "bin/bench.sh", "#!/usr/bin/env bash\nexit 0\n", 0o755)
	writeGateTestFile(t, root, ".bench/gate-inputs.json",
		`{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	writeGateTestFile(t, root, ".bench/gate.sh", fixtureGateScript, 0o755)
	// The race phase materializes only for a root that declares the registered race tests,
	// and a phase the table never carries has no identity to grade. The declarations come
	// from the registry's own generator, so the fixture keeps no second list of their names.
	for rel, source := range racetests.SyntheticSources() {
		writeGateTestFile(t, root, strings.TrimPrefix(rel, "./"), source, 0o644)
	}
	writeCaptureSurfaces(t, root)
	writeHandDeclaredSurfaces(t, root)
	writeConsumerInventoryFixture(t, root)
}

// writeConformanceCanaryOwners gives the fixture every implementation whose check owns a
// canary family. The registry remains the source of that membership, so a new owner joins
// the fixture without a copied function list and canary identity resolution stays honest.
func writeConformanceCanaryOwners(t *testing.T, root string) {
	t.Helper()
	owners := map[string]bool{}
	for _, check := range registry.Checks {
		if check.Meta || !check.RunsAt(registry.Dev) || len(registry.CanaryFamilies(check.Name)) == 0 {
			continue
		}
		owners[check.Implementation] = true
	}
	names := make([]string, 0, len(owners))
	for name := range owners {
		names = append(names, name)
	}
	sort.Strings(names)
	var source strings.Builder
	source.WriteString("package conformance\n\n")
	for _, name := range names {
		source.WriteString("func ")
		source.WriteString(name)
		source.WriteString("() {}\n")
	}
	writeGateTestFile(t, root, "internal/conformance/owners.go", source.String(), 0o644)
}

// writeConsumerInventoryFixture materializes the consumer inventory's fixture shape from
// the kit's one allowlist. A consumer row that joins the payload therefore joins every
// shared fixture without a copied path list beside it.
func writeConsumerInventoryFixture(t *testing.T, root string) {
	t.Helper()
	rows, err := kitpayload.PayloadRows()
	if err != nil {
		t.Fatalf("read consumer payload rows: %v", err)
	}
	payload, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, root, ".bench/consumer-payload.json", string(payload), 0o644)
	for _, row := range kitpayload.PayloadConsumerRows(rows) {
		path := filepath.Join(root, filepath.FromSlash(row.Source))
		if row.Tree {
			if err := os.MkdirAll(path, 0o755); err != nil {
				t.Fatal(err)
			}
			writeGateTestFile(t, root, row.Source+"/consumer-fixture.md", "# Consumer fixture\n", 0o644)
			continue
		}
		if isRegularFile(path) {
			continue
		}
		writeGateTestFile(t, root, row.Source, "consumer fixture asset\n", 0o644)
	}
}

// writeHandDeclaredSurfaces materializes whatever canary's declaration names and the tree
// above does not already carry. It is derived from that declaration rather than listed
// again here, so a surface the declaration gains lands on this root without an edit.
//
// Only the hand-declared entry needs this. Every other component's inputs are computed
// from a listing of the tree, so they can only name files the fixture already wrote —
// while a hand-declared surface the root lacks leaves canary unable to compute an identity
// at all, and every scoping assertion about it grading nothing.
func writeHandDeclaredSurfaces(t *testing.T, root string) {
	t.Helper()
	declared, _, err := (&inputResolver{root: root}).canaryInputs()
	if err != nil {
		t.Fatalf("resolve the hand-declared canary inputs: %v", err)
	}
	for _, path := range declared {
		if declaresDirectory(path) {
			// A directory entry is satisfied by any file beneath it, and the sources and
			// fixtures the canary phase grades are already written above.
			if !holdsRegularFile(filepath.Join(root, filepath.FromSlash(path))) {
				writeGateTestFile(t, root, path+"declared.txt", "declared canary surface\n", 0o644)
			}
			continue
		}
		// The declared file entries are the wrapper scripts the canary phase's wiring
		// execs. Nothing here ever runs them — the phase manifest routes every phase
		// through a marker script — so they carry the same inert body the other fixture
		// scripts do.
		if !isRegularFile(filepath.Join(root, filepath.FromSlash(path))) {
			writeGateTestFile(t, root, path, "#!/usr/bin/env bash\nexit 0\n", 0o755)
		}
	}
}

// holdsRegularFile reports whether dir holds a regular file at any depth.
func holdsRegularFile(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

// captureSurfacePaths derives one editable path per entry in the reduced-run declaration:
// each declared file, and one file under each declared directory. Deriving them is what
// keeps a capture-only changeset expressible on this root when the declaration gains an
// entry, and what makes the guard below able to catch a declaration that moved.
func captureSurfacePaths(scope Scope) []string {
	paths := scope.Files()
	for _, dir := range scope.Directories() {
		paths = append(paths, dir+"declared.md")
	}
	return paths
}

// writeCaptureSurfaces materializes the declared capture paths, guarding each one against
// the declaration it was derived from. The guard bites on a directory entry that lost its
// trailing slash or a file entry that stopped being repository-relative: either leaves a
// path this fixture writes but Confines refuses, and every reduced-run assertion built on
// the root would then quietly observe a full run instead.
func writeCaptureSurfaces(t *testing.T, root string) {
	t.Helper()
	scope := ReducedScope()
	if !scope.Excludable(canary.PhaseTest) || scope.Excludable(conformancePhaseName) {
		t.Fatal("fixture phase names no longer match the declaration; repoint the fixture")
	}
	for _, path := range captureSurfacePaths(scope) {
		if !scope.Member(path) {
			t.Fatalf("derived capture path %q is not declared; repoint the fixture", path)
		}
		writeGateTestFile(t, root, path, "declared capture surface\n", 0o644)
	}
}

func sealInitialKitShapedBinary(t *testing.T, root string) {
	t.Helper()
	staged := filepath.Join(root, "dist", "bench.staged")
	template, err := kitShapedTemplate(root)
	if err != nil {
		t.Fatalf("build the fixture binary template: %v", err)
	}
	if err := materializeFixtureBinary(template, staged); err != nil {
		t.Fatalf("materialize the fixture binary template: %v", err)
	}
	publishKitShapedBinary(t, root, staged)
}

// sealKitShapedBinary publishes dist/bench through the seal package's Publish, the only
// writer that produces a seal answering for the tree the binary was built from. Moving the
// bytes into place by hand leaves an executable no reader can verify.
func sealKitShapedBinary(t *testing.T, root string) {
	t.Helper()
	staged := filepath.Join(root, "dist", "bench.staged")
	buildFixtureBinaryTo(t, root, "./cmd/bench", staged)
	publishKitShapedBinary(t, root, staged)
}

func publishKitShapedBinary(t *testing.T, root, staged string) {
	t.Helper()
	if err := benchfreshness.Publish(root, staged, filepath.Join(root, "dist", "bench")); err != nil {
		t.Fatalf("publish the fixture binary: %v", err)
	}
}

func requireKitShapedBinaryFresh(t *testing.T, root string) {
	t.Helper()
	if err := benchfreshness.Verify(root, filepath.Join(root, "dist", "bench")); err != nil {
		t.Fatalf("fixture-only edit changed the synthetic Bench build inputs: %v", err)
	}
}

func materializeFixtureBinary(source, destination string) error {
	return materializeFixtureBinaryWithLink(source, destination, os.Link)
}

func materializeFixtureBinaryWithLink(source, destination string, link func(string, string) error) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	linkErr := link(source, destination)
	if linkErr == nil {
		return nil
	}
	if _, err := os.Lstat(destination); err == nil {
		return fmt.Errorf("link fixture binary: destination already exists: %w", linkErr)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect fixture binary destination: %w", err)
	}
	return copyFixtureBinary(source, destination)
}

func copyFixtureBinary(source, destination string) (err error) {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := input.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := output.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
		if err != nil {
			_ = os.Remove(destination)
		}
	}()
	if _, err := io.Copy(output, input); err != nil {
		return err
	}
	return os.Chmod(destination, info.Mode().Perm())
}

func makeFixtureBinaryPrivate(path string) (err error) {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".bench-private-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		if removeErr := os.Remove(temporaryPath); err == nil && removeErr != nil && !os.IsNotExist(removeErr) {
			err = removeErr
		}
	}()
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := copyFixtureBinary(path, temporaryPath); err != nil {
		return err
	}
	if err := os.Chmod(temporaryPath, 0o755); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

func TestKitShapedFixturesDetachIndependentTemplateLinks(t *testing.T) {
	t.Parallel()
	first := newKitShapedFixture(t)
	second := newKitShapedFixture(t)
	firstInfo, err := os.Stat(first.binaryPath())
	if err != nil {
		t.Fatal(err)
	}
	secondInfo, err := os.Stat(second.binaryPath())
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(firstInfo, secondInfo) {
		t.Fatal("initial fixture binaries do not share the immutable template inode")
	}
	if err := makeFixtureBinaryPrivate(first.binaryPath()); err != nil {
		t.Fatal(err)
	}
	privateInfo, err := os.Stat(first.binaryPath())
	if err != nil {
		t.Fatal(err)
	}
	if os.SameFile(privateInfo, secondInfo) {
		t.Fatal("detached fixture binary still shares the immutable template inode")
	}
	if err := os.WriteFile(first.binaryPath(), []byte("changed"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := benchfreshness.Verify(second.root, second.binaryPath()); err != nil {
		t.Fatalf("second fixture binary after mutating the first = %v, want independent published bytes", err)
	}
}

func TestMaterializeFixtureBinaryFallsBackToPrivateCopy(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	source := filepath.Join(root, "template")
	destination := filepath.Join(root, "fixture", "dist", "bench.staged")
	writeGateTestFile(t, root, "template", "fixture binary\n", 0o755)

	linkErr := fmt.Errorf("links unavailable")
	if err := materializeFixtureBinaryWithLink(source, destination, func(string, string) error {
		return linkErr
	}); err != nil {
		t.Fatal(err)
	}
	if got := string(mustRead(t, destination)); got != "fixture binary\n" {
		t.Fatalf("fallback bytes = %q, want the template bytes", got)
	}
	sourceInfo, err := os.Stat(source)
	if err != nil {
		t.Fatal(err)
	}
	destinationInfo, err := os.Stat(destination)
	if err != nil {
		t.Fatal(err)
	}
	if os.SameFile(sourceInfo, destinationInfo) {
		t.Fatal("copy fallback still shares the template inode")
	}
	if err := materializeFixtureBinaryWithLink(source, destination, func(string, string) error {
		return linkErr
	}); err == nil {
		t.Fatal("materialize over an existing destination succeeded")
	}
	if got := string(mustRead(t, destination)); got != "fixture binary\n" {
		t.Fatalf("existing destination after refused fallback = %q, want unchanged bytes", got)
	}
}

func TestKitShapedFixtureTemplateFailureIsSticky(t *testing.T) {
	t.Parallel()
	if os.Getenv("BENCH_KITSHAPED_TEMPLATE_FAILURE") != "1" {
		report := filepath.Join(t.TempDir(), "template-directory")
		command := exec.Command(os.Args[0], "-test.run", "^TestKitShapedFixtureTemplateFailureIsSticky$")
		command.Env = append(os.Environ(),
			"BENCH_KITSHAPED_TEMPLATE_FAILURE=1",
			"BENCH_KITSHAPED_TEMPLATE_REPORT="+report,
		)
		if out, err := command.CombinedOutput(); err != nil {
			t.Fatalf("template failure subprocess: %v\n%s", err, out)
		}
		contents, err := os.ReadFile(report)
		if err != nil {
			t.Fatalf("read the template lifetime report: %v", err)
		}
		dir := string(contents)
		if dir == "" {
			t.Fatal("template lifetime report is empty")
		}
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Fatalf("template directory after the test process = %v, want removed", err)
		}
		return
	}
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	writeKitShapedTree(t, root)
	writeGateTestFile(t, root, "cmd/bench/main.go", "package main\nfunc main( {\n", 0o644)
	_, first := kitShapedTemplate(root)
	_, second := kitShapedTemplate(root)
	if first == nil || second == nil || first != second {
		t.Fatalf("template failures = (%v, %v), want the same retained construction error", first, second)
	}
	if err := os.WriteFile(os.Getenv("BENCH_KITSHAPED_TEMPLATE_REPORT"), []byte(kitShapedTemplateState.dir), 0o644); err != nil {
		t.Fatalf("write the template lifetime report: %v", err)
	}
}

// writeKitShapedManifest declares one marker phase per name in declared, and writes the
// script each of them runs.
//
// The declared edges are carried through. They are what orders a writer of dist/bench ahead of
// its readers, so a fixture that dropped them would let a phase copying the artifact race the
// phase producing it — and would leave the build phase's skip untested against the edges it is
// claimed to satisfy trivially.
func writeKitShapedManifest(t *testing.T, root string, declared []Phase) {
	t.Helper()
	var doc manifestDoc
	for _, phase := range declared {
		script := ".bench/phase-" + phase.Name + ".sh"
		writeGateTestFile(t, root, script, "echo "+phase.Name+" >> .git/phase-runs\n", 0o644)
		doc.Phases = append(doc.Phases, manifestPhase{Name: phase.Name, Argv: []string{"bash", script}, Needs: phase.Needs})
	}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, root, canary.PhaseManifestPath, string(data)+"\n", 0o644)
}

// goListClosure returns the package set `go list` reports for args, keyed by the exact
// line so a test-binary entry stays distinguishable from the package it was built from.
func goListClosure(t *testing.T, root string, args ...string) map[string]bool {
	t.Helper()
	cmd := exec.Command("go", append([]string{"list", "-buildvcs=false"}, args...)...)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list %s: %v", strings.Join(args, " "), err)
	}
	closure := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			closure[line] = true
		}
	}
	return closure
}

// [PS5] The fixture's resolved table carries build and canary by name. Every later
// assertion about either component rests on this, so it is stated on its own: a tree shape
// that stops materializing the build phase leaves those assertions passing over a table
// that never declared what they claim to grade.
func TestKitShapedFixtureCarriesBuildAndCanary(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	for _, name := range []string{canary.PhaseBuild, "canary"} {
		if !carriesPhase(fixture.phases, name) {
			t.Fatalf("resolved table = %v, want a %q phase", fixture.phaseNames(), name)
		}
	}
}

// [PS4] The fixture gates green through the whole-tree path, and every phase in its
// resolved table leaves exactly one marker line. Comparing sorted sets is what pins both
// halves at once — a phase the gate script never reached is missing, and a phase run twice
// is duplicated.
func TestKitShapedFixtureGatesGreen(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	if got := fullRunCount(t, fixture.root); got != 1 {
		t.Fatalf("resolved gate runs = %d, want 1", got)
	}
	executed, want := phaseRunNames(t, fixture.root), fixture.phaseNames()
	sort.Strings(executed)
	sort.Strings(want)
	if !reflect.DeepEqual(executed, want) {
		t.Fatalf("executed phases = %v, want exactly the resolved table %v", executed, want)
	}
}

// [PS4b] The root narrows: a changeset confined to the declared capture surfaces produces a
// run that executes the unconditional phases and skips every evidence-covered component.
// Every later assertion about per-component scoping needs this to hold on this root — a
// fixture whose gate script carries no gate-phases hand-off pays a full run for every
// changeset, and a scoping test written against it would pass while observing nothing at all.
func TestKitShapedFixtureNarrowsACaptureOnlyChangeset(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	seeded := phaseRunNames(t, fixture.root)

	for _, path := range captureSurfacePaths(ReducedScope()) {
		writeGateTestFile(t, fixture.root, path, "capture-only edit\n", 0o644)
	}
	mustExecuteGreen(t, fixture.root, productionGateEngine{})

	if rec := slotRecord(t, fixture.root, time.Now().UTC()); rec.partition() == nil {
		t.Fatalf("capture-only record = %+v, want a partial verdict", rec)
	}
	if got := fullRunCount(t, fixture.root); got != 1 {
		t.Fatalf("resolved gate runs = %d, want 1 — the capture-only edit paid the full gate", got)
	}
	// The expected sides come from the same predicate the decision reads, so a phase added
	// to the fixture lands on whichever side that predicate puts it without an edit here.
	want, skipped := unconditionalPhaseNames(fixture.phases)
	executed := append([]string(nil), phaseRunNames(t, fixture.root)[len(seeded):]...)
	sort.Strings(executed)
	sort.Strings(want)
	if !reflect.DeepEqual(executed, want) {
		t.Fatalf("narrowed run executed %v, want the unconditional set %v (skipping %v)", executed, want, skipped)
	}
}

// unconditionalPhaseNames splits a resolved table by the skip predicate: the phases no
// evidence of any form can cover, and the component names some evidence can. The split reads
// the predicate rather than naming phases, so build joining the skippable side on its attested
// seal moves it here without an edit.
func unconditionalPhaseNames(table []Phase) (unconditional, skippable []string) {
	for _, phase := range table {
		if componentSkipsOnEvidence(phase.Name) {
			skippable = append(skippable, phase.Name)
			continue
		}
		unconditional = append(unconditional, phase.Name)
	}
	return unconditional, skippable
}

// [PS6] The fixture carries a package outside ./cmd/bench's build closure whose test files
// are inside the module-wide one. The toolchain components grade that package, so a
// derivation reading the binary's closure would call their inputs unmoved while the tree's
// test suite no longer builds — the hole this shape exists to expose.
func TestKitShapedFixtureHasAPackageOutsideTheBinaryClosure(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	binary := goListClosure(t, fixture.root, "-deps", "./cmd/bench")
	module := goListClosure(t, fixture.root, "-deps", "-test", "./...")
	if binary[outsideBinaryClosurePackage] {
		t.Fatalf("%s is inside ./cmd/bench's build closure; the fixture no longer spans the two closures", outsideBinaryClosurePackage)
	}
	if !module[outsideBinaryClosurePackage] {
		t.Fatalf("%s is absent from the module-wide closure: %v", outsideBinaryClosurePackage, sortedClosure(module))
	}
	if !module[outsideBinaryClosurePackage+".test"] {
		t.Fatalf("%s contributes no test binary, so its _test.go files are outside the module-wide closure: %v",
			outsideBinaryClosurePackage, sortedClosure(module))
	}
}

// [PS7] The published binary verifies against the tree it was built from, immediately
// after construction. Every build-skip decision rests on that seal, so a fixture whose
// dist/bench cannot be verified would make a skip look sound for the wrong reason.
func TestKitShapedFixtureBinaryIsSealed(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	if err := benchfreshness.Verify(fixture.root, fixture.binaryPath()); err != nil {
		t.Fatalf("benchfreshness.Verify = %v, want the published binary to verify", err)
	}
}

// [PS8] A second run over an unedited fixture reuses the whole-tree green: the resolved
// gate is not paid again and no phase leaves a new marker. A phase script that writes into
// the graded tree moves the subject under its own run, and every later test built on this
// root would then be measuring that drift rather than the decision it meant to observe.
func TestKitShapedFixtureReusesItsGreen(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	first := phaseRunNames(t, fixture.root)

	mustExecuteGreen(t, fixture.root, productionGateEngine{})
	if got := fullRunCount(t, fixture.root); got != 1 {
		t.Fatalf("resolved gate runs = %d, want 1 — the unedited tree did not reuse its green", got)
	}
	if got := phaseRunNames(t, fixture.root); !reflect.DeepEqual(got, first) {
		t.Fatalf("executed phases after the second run = %v, want the first run's %v unchanged", got, first)
	}
}

// [PS15d] Every component the input registry declares resolves an identity on this root.
// The family is read from the registry rather than listed here, so a component added to it
// is a component this root has to be able to grade — and a declaration naming a surface the
// fixture lacks refuses, which is a component that would run every time on the real kit.
func TestKitShapedFixtureResolvesEveryRegistryComponentIdentity(t *testing.T) {
	t.Parallel()
	fixture := newKitShapedFixture(t)
	identities := mustResolveComponentIdentities(t, fixture.root)
	for _, declaration := range componentInputDeclarations() {
		if identities[declaration.component] == "" {
			t.Errorf("no identity for the declared component %q; resolved %v", declaration.component, identities)
		}
	}
}

func sortedClosure(closure map[string]bool) []string {
	names := make([]string, 0, len(closure))
	for name := range closure {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
