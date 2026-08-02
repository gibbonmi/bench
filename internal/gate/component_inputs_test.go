package gate

// The component input declarations are graded against the kit-shaped fixture, which spans
// the two Go closures on purpose: a derivation that reads only the binary's closure passes
// every assertion about the packages ./cmd/bench imports and is blind to the rest.

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
)

// moduleClosureComponents are the components declared to read the module-wide closure.
// Taking them from the registry rather than listing them keeps a component that joins that
// derivation graded by these tests from the moment it joins.
func moduleClosureComponents() []string {
	var names []string
	for _, declaration := range componentInputDeclarations() {
		if declaration.source == SourceModuleTestClosure || declaration.source == SourceModuleTestClosureWithSealAndAgentMarkdown {
			names = append(names, declaration.component)
		}
	}
	return names
}

func mustResolveComponentInputs(t *testing.T, root string) map[string]ComponentInputs {
	t.Helper()
	sets, err := ResolveComponentInputs(root)
	if err != nil {
		t.Fatalf("ResolveComponentInputs = %v, want the fixture's declarations", err)
	}
	return sets
}

func componentEntry(t *testing.T, sets map[string]ComponentInputs, name string) ComponentInputs {
	t.Helper()
	entry, declared := sets[name]
	if !declared {
		t.Fatalf("no input declaration for %q", name)
	}
	return entry
}

// [PC8a] Every module-closure component declares the test files of a package outside
// ./cmd/bench's closure. Those files are what the test and race components grade, so a
// declaration blind to them buys a skip while the tree's suite no longer builds.
func TestToolchainInputsCoverTestFilesOutsideTheBinaryClosure(t *testing.T) {
	fixture := newKitShapedFixture(t)
	sets := mustResolveComponentInputs(t, fixture.root)

	const outside = "internal/canary/canary_test.go"
	for _, name := range moduleClosureComponents() {
		if paths := componentEntry(t, sets, name).Paths(); !slices.Contains(paths, outside) {
			t.Fatalf("%s inputs omit %q, which lives outside ./cmd/bench's closure: %v", name, outside, paths)
		}
	}
}

// [PC8b] The declarations carry the testdata/ contents of listed packages. `go list` names
// no file under testdata, so a derivation that trusts the listing alone reports a package's
// fixtures unmoved when a test that reads them starts failing.
func TestToolchainInputsCoverListedTestdata(t *testing.T) {
	fixture := newKitShapedFixture(t)
	const fixtureCase = "internal/canary/testdata/case.txt"
	writeGateTestFile(t, fixture.root, fixtureCase, "graded fixture\n", 0o644)

	sets := mustResolveComponentInputs(t, fixture.root)
	for _, name := range moduleClosureComponents() {
		if paths := componentEntry(t, sets, name).Paths(); !slices.Contains(paths, fixtureCase) {
			t.Fatalf("%s inputs omit %q: %v", name, fixtureCase, paths)
		}
	}
}

// [PS9] The build component's set is the freshness digest's own build inputs, element for
// element. That derivation is what decides whether the published binary is stale; a second
// list here would answer that question differently the day either one moved.
// The comparison is made twice, across a source the binary gains in between: a hand-written
// list can match one tree, so agreeing on both is what separates the derivation from a
// literal that happens to be right today.
func TestBuildInputsAreNotRestated(t *testing.T) {
	fixture := newKitShapedFixture(t)
	assertBuildInputsMatchTheDigest(t, fixture.root)

	writeGateTestFile(t, fixture.root, "cmd/bench/extra.go", "package main\n\nvar extra = 1\n", 0o644)
	assertBuildInputsMatchTheDigest(t, fixture.root)
}

func assertBuildInputsMatchTheDigest(t *testing.T, root string) {
	t.Helper()
	want, err := benchfreshness.BuildInputs(root)
	if err != nil {
		t.Fatalf("benchfreshness.BuildInputs = %v", err)
	}
	got := componentEntry(t, mustResolveComponentInputs(t, root), canary.PhaseBuild).Paths()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("build inputs = %v, want freshness.BuildInputs' %v", got, want)
	}
}

// [PS10] The contract component execs the published CLI, so the seal's source digest is one
// of its inputs and moves when the binary's sources do. Declaring only the file set would
// let a contract skip inherit evidence taken against a different binary.
func TestContractInputsCarryTheSealSourceDigest(t *testing.T) {
	fixture := newKitShapedFixture(t)
	before := componentEntry(t, mustResolveComponentInputs(t, fixture.root), canary.PhaseContract).Digests()
	if len(before) == 0 {
		t.Fatalf("contract digests = %v, want the seal's source digest", before)
	}

	writeGateTestFile(t, fixture.root, "cmd/bench/main.go", "package main\n\nfunc main() { _ = 1 }\n", 0o644)
	sealKitShapedBinary(t, fixture.root)

	after := componentEntry(t, mustResolveComponentInputs(t, fixture.root), canary.PhaseContract).Digests()
	if reflect.DeepEqual(before, after) {
		t.Fatalf("contract digests = %v after republishing, want the seal's source digest to move", after)
	}
}

func TestContractInputsCoverDeclaredAgentMarkdown(t *testing.T) {
	fixture := newKitShapedFixture(t)
	const markdown = ".agents/skills/new-skill/SKILL.md"
	const structured = ".agents/skills/new-skill/agents/openai.yaml"
	writeGateTestFile(t, fixture.root, markdown, "# New skill\n", 0o644)
	writeGateTestFile(t, fixture.root, structured, "interface:\n", 0o644)

	paths := componentEntry(t, mustResolveComponentInputs(t, fixture.root), canary.PhaseContract).Paths()
	if !slices.Contains(paths, markdown) {
		t.Fatalf("contract inputs omit declared agent Markdown %q: %v", markdown, paths)
	}
	if slices.Contains(paths, structured) {
		t.Fatalf("contract inputs include non-Markdown agent asset %q through the Markdown derivation", structured)
	}
}

// [PS11] Every entry names the derivation it came from. The profile table and the
// derivation-source conformance check read that name to tell a computed set from a
// hand-written one.
func TestComponentInputsReportANamedDerivationSource(t *testing.T) {
	fixture := newKitShapedFixture(t)
	named := []Source{
		SourceBuildClosure, SourceModuleTestClosure, SourceModuleTestClosureWithSealAndAgentMarkdown,
		SourceShellcheckArgv, SourceHandDeclared,
	}

	for name, entry := range mustResolveComponentInputs(t, fixture.root) {
		if !slices.Contains(named, entry.Source()) {
			t.Fatalf("%s source = %q, want one of the named derivations %v", name, entry.Source(), named)
		}
		if entry.Component() != name {
			t.Fatalf("entry keyed %q answers for %q", name, entry.Component())
		}
	}
}

// [PS11] A failed derivation fails the resolution. Returning whatever listed successfully
// would hand back a set smaller than what the component reads, and every path missing from
// it becomes a change the component skips.
func TestComponentInputsErrorOnDerivationFailure(t *testing.T) {
	fixture := newKitShapedFixture(t)
	if err := os.WriteFile(filepath.Join(fixture.root, "go.mod"), []byte("this is not a go.mod\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	sets, err := ResolveComponentInputs(fixture.root)
	if err == nil {
		t.Fatalf("ResolveComponentInputs = %v, want an error on the corrupted module", sets)
	}
	// Each derivation is driven on its own as well. The registry resolves in order and stops
	// at the first failure, so a component behind the first one would otherwise be reported
	// as refusing when nothing had asked it anything. This probes only the go-list-backed
	// derivations: shellcheck reads its argv's files directly and canary is hand-declared,
	// so neither one reads go.mod at all, and a corrupted module is not a failure either
	// could be expected to notice.
	for _, declaration := range componentInputDeclarations() {
		if declaration.source == SourceShellcheckArgv || declaration.source == SourceHandDeclared {
			continue
		}
		paths, _, err := declaration.resolve(&inputResolver{root: fixture.root})
		if err == nil {
			t.Fatalf("%s derivation = %v, want an error on the corrupted module", declaration.component, paths)
		}
		if len(paths) != 0 {
			t.Fatalf("%s reported %v alongside its error, want no paths", declaration.component, paths)
		}
	}
}

// [PS12] shellcheck's declared set is shellcheckFiles' own enumeration — the same one
// shellcheckArgv builds its invocation on top of: a script the enumeration gains is a
// script the declaration gains, with no edit to the derivation.
func TestShellcheckInputsFollowItsArgv(t *testing.T) {
	fixture := newKitShapedFixture(t)
	const added = ".bench/hooks/new.sh"
	writeGateTestFile(t, fixture.root, added, "#!/usr/bin/env bash\nexit 0\n", 0o755)

	paths := componentEntry(t, mustResolveComponentInputs(t, fixture.root), "shellcheck").Paths()
	if !slices.Contains(paths, added) {
		t.Fatalf("shellcheck inputs = %v, want the added script %q", paths, added)
	}
}

// [PS12b] shellcheck's declared set is shellcheckFiles' enumeration exactly — not
// shellcheckArgv sliced at a remembered offset. A flag added to shellcheck's invocation
// prefix must not drop or shift a linted path out of the declaration, which a
// fixed-offset slice of the argv would do the moment the prefix's length changed.
func TestShellcheckInputsMatchTheFileEnumerationExactly(t *testing.T) {
	fixture := newKitShapedFixture(t)
	want := shellcheckFiles(fixture.root)
	slices.Sort(want)

	got := componentEntry(t, mustResolveComponentInputs(t, fixture.root), "shellcheck").Paths()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("shellcheck inputs = %v, want shellcheckFiles' %v exactly", got, want)
	}
}

// [PS13] canary's declared set carries no binary digest: republishing a different binary
// over unchanged sources leaves the declaration unmoved.
func TestCanaryInputsExcludeTheBinary(t *testing.T) {
	fixture := newKitShapedFixture(t)
	before := componentEntry(t, mustResolveComponentInputs(t, fixture.root), "canary")
	if len(before.Digests()) != 0 {
		t.Fatalf("canary digests = %v, want none", before.Digests())
	}

	publishDifferentKitShapedBinary(t, fixture.root)

	after := componentEntry(t, mustResolveComponentInputs(t, fixture.root), "canary")
	if !reflect.DeepEqual(before.Paths(), after.Paths()) || !reflect.DeepEqual(before.Digests(), after.Digests()) {
		t.Fatalf("canary inputs moved after republishing the binary: before %v/%v, after %v/%v",
			before.Paths(), before.Digests(), after.Paths(), after.Digests())
	}
}

// publishDifferentKitShapedBinary republishes root's dist/bench with different executable
// bytes over the same source tree — sources unchanged, binary changed — the shape
// TestCanaryInputsExcludeTheBinary needs to tell "reads the binary" apart from "does not".
func publishDifferentKitShapedBinary(t *testing.T, root string) {
	t.Helper()
	executable := filepath.Join(root, "dist", "bench")
	original, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(root, "dist", "bench.restaged")
	if err := os.WriteFile(staged, append(append([]byte{}, original...), 0), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := benchfreshness.Publish(root, staged, executable); err != nil {
		t.Fatalf("publish a different binary: %v", err)
	}
}

// [PS14] canary is the registry's only hand-declared entry; every other entry names a
// derivation. A component that has no derivable source but reports one anyway would claim
// a provenance it never computed.
func TestOnlyCanaryIsHandDeclared(t *testing.T) {
	var handDeclared []string
	for _, declaration := range componentInputDeclarations() {
		if declaration.source == SourceHandDeclared {
			handDeclared = append(handDeclared, declaration.component)
		}
	}
	if want := []string{"canary"}; !reflect.DeepEqual(handDeclared, want) {
		t.Fatalf("hand-declared entries = %v, want exactly %v", handDeclared, want)
	}
}

// [PC8c] The toolchain and contract input sets carry the module manifest: go.mod always,
// and go.sum once the module has one. A dependency version bump edits both files while
// leaving the derived closure's file set byte-identical, so a declaration blind to them
// would call gofmt, vet, test, race, and contract unmoved by a change that can red all five.
func TestToolchainInputsCoverTheModuleManifest(t *testing.T) {
	fixture := newKitShapedFixture(t)

	sets := mustResolveComponentInputs(t, fixture.root)
	for _, name := range moduleClosureComponents() {
		if paths := componentEntry(t, sets, name).Paths(); !slices.Contains(paths, "go.mod") {
			t.Fatalf("%s inputs = %v, want go.mod", name, paths)
		}
	}

	writeGateTestFile(t, fixture.root, "go.sum",
		"example.com/dummy v1.0.0 h1:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=\n", 0o644)
	sets = mustResolveComponentInputs(t, fixture.root)
	for _, name := range moduleClosureComponents() {
		if paths := componentEntry(t, sets, name).Paths(); !slices.Contains(paths, "go.sum") {
			t.Fatalf("%s inputs = %v, want go.sum once the module has one", name, paths)
		}
	}
}

// [PC8d] Editing go.mod's content moves the toolchain and contract components' resolved
// identities even though it leaves the module-test closure's derived file set unchanged —
// go.mod is not itself a file `go list` names. A declaration that names the manifest in its
// source label without folding it into the resolved path set would leave the identity
// unmoved: the identity is computed from the declared paths' contents, so a path missing
// from that set is evidence the dependency bump never reaches.
func TestModuleManifestEditMovesToolchainInputs(t *testing.T) {
	fixture := newKitShapedFixture(t)
	before := mustResolveComponentIdentities(t, fixture.root)

	writeGateTestFile(t, fixture.root, "go.mod",
		"module "+fixtureModulePath+"\n\ngo 1.21\n\n// bumped dependency version\n", 0o644)

	after := mustResolveComponentIdentities(t, fixture.root)
	for _, name := range moduleClosureComponents() {
		if before[name] == after[name] {
			t.Fatalf("%s identity = %q before and after the go.mod edit, want it to move", name, before[name])
		}
	}
}

// [PS37] Every module-closure entry's source names both halves of its shape — the
// module-test closure and the manifest addition — so the derivation-source conformance
// check can tell this two-part derivation from a hand-copied path list.
func TestToolchainSourceNamesTheManifestAddition(t *testing.T) {
	fixture := newKitShapedFixture(t)
	sets := mustResolveComponentInputs(t, fixture.root)
	for _, name := range moduleClosureComponents() {
		source := string(componentEntry(t, sets, name).Source())
		if !strings.Contains(source, "module-test-closure") || !strings.Contains(source, "manifest") {
			t.Fatalf("%s source = %q, want it to name both the module-test closure and the manifest addition", name, source)
		}
	}
}

// [PS38] A module with no go.sum resolves without error and without a phantom go.sum entry:
// the addition is bounded to the manifest the module actually has, not a path added
// unconditionally.
func TestModuleWithoutGoSumResolves(t *testing.T) {
	fixture := newKitShapedFixture(t)
	if _, err := os.Lstat(filepath.Join(fixture.root, "go.sum")); err == nil {
		t.Fatalf("fixture unexpectedly carries go.sum; repoint the fixture")
	}

	sets := mustResolveComponentInputs(t, fixture.root)
	for _, name := range moduleClosureComponents() {
		if paths := componentEntry(t, sets, name).Paths(); slices.Contains(paths, "go.sum") {
			t.Fatalf("%s inputs = %v, want no go.sum entry for a module without one", name, paths)
		}
	}
}
