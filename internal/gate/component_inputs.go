package gate

// What each scoped gate component reads. A component may inherit an ancestor's evidence
// only for a changeset its declared inputs do not contain, so an entry that under-reports
// its inputs buys a skip over work nobody graded. Every entry here therefore names the
// derivation it came from and computes its set from that derivation on the spot: a copied
// path list would agree with the tree only until the tree moved.
//
// The declaration and its provenance are halves of one claim, so they live together, and
// the fields are reached through accessors — a consumer holding a copy cannot edit the
// source out from under the next reader.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/packagesurface"
)

// Source names the derivation an entry's inputs were computed from. The profile's rendered
// table and the derivation-source conformance check both read this one value, so a
// component's provenance is stated once.
type Source string

const (
	// SourceModuleTestClosure is the module-wide `go list -deps -test ./...` closure plus
	// the testdata/ contents of the packages it lists, plus the module manifest: go.mod
	// always, and go.sum when the module has one. `go list` reports neither file, so a
	// dependency version bump that leaves the closure's file set byte-identical would
	// otherwise leave this declaration unmoved by a change that can red every component
	// that carries it.
	SourceModuleTestClosure Source = "module-test-closure+manifest"
	// SourceModuleTestClosureWithConsumerDocuments adds the consumer-inventory documents
	// lifecycle contracts consume.
	SourceModuleTestClosureWithConsumerDocuments Source = "module-test-closure+manifest+consumer-document-inventory"
	// SourceShellcheckArgv is shellcheckArgv's own file enumeration — the exact argument
	// list the shellcheck phase lints, read rather than restated.
	SourceShellcheckArgv Source = "shellcheck-argv"
	// SourceHandDeclared marks the one entry with no derivable source: canary's surfaces
	// are named directly rather than computed from a listing.
	SourceHandDeclared Source = "hand-declared"
)

// ComponentInputs is one component's declared input set: the repository-relative,
// slash-separated, sorted paths it reads, plus any input whose content has no file in the
// tree, and the derivation all of it came from.
type ComponentInputs struct {
	component string
	source    Source
	paths     []string
	digests   []string
}

// Component is the phase name this declaration answers for.
func (c ComponentInputs) Component() string { return c.component }

// Source is the derivation the declaration was computed from.
func (c ComponentInputs) Source() Source { return c.source }

// Paths are the declared repository-relative input paths, sorted and deduplicated.
func (c ComponentInputs) Paths() []string { return slices.Clone(c.paths) }

// Digests are declared inputs that name content rather than a path.
func (c ComponentInputs) Digests() []string { return slices.Clone(c.digests) }

// componentInputDeclaration binds a component to the derivation that answers for it. The
// resolver is a method value rather than a path list so the binding cannot drift from what
// the derivation actually returns, and so a component whose inputs have no derivable
// source can join this table later without changing its shape.
type componentInputDeclaration struct {
	component string
	source    Source
	resolve   func(*inputResolver) (paths []string, digests []string, err error)
}

// componentInputDeclarations is the registry: one entry per scoped component whose inputs
// this package declares. It is a function rather than a package variable so no consumer
// can rewrite the family it enumerates.
//
// The toolchain components take the module-wide closure, not the binary's: `./cmd/bench`
// excludes the conformance and contract packages they grade, and a component blind to the
// package it grades would skip on a change that reds it.
func componentInputDeclarations() []componentInputDeclaration {
	return []componentInputDeclaration{
		{canary.PhaseGofmt, SourceModuleTestClosure, (*inputResolver).moduleClosure},
		{canary.PhaseVet, SourceModuleTestClosure, (*inputResolver).moduleClosure},
		{canary.PhaseTest, SourceModuleTestClosure, (*inputResolver).moduleClosure},
		{canary.PhaseRace, SourceModuleTestClosure, (*inputResolver).moduleClosure},
		{canary.PhaseConformanceSuite, SourceModuleTestClosure, (*inputResolver).moduleClosure},
		{canary.PhaseContract, SourceModuleTestClosureWithConsumerDocuments, (*inputResolver).contractInputs},
		{"shellcheck", SourceShellcheckArgv, (*inputResolver).shellcheckInputs},
		{"canary", SourceHandDeclared, (*inputResolver).canaryInputs},
	}
}

// ComponentInputSource is one registry entry's component/source pair, with no filesystem
// resolution attached — for a consumer that binds a profile's rendering of the registry to
// the registry itself and has no use for the resolved paths.
type ComponentInputSource struct {
	Component string
	Source    Source
}

// ComponentInputSources exposes the registry's component/source pairs, in registration
// order. A caller wanting the resolved paths and digests calls ResolveComponentInputs
// instead; this accessor exists so a profile-binding check can enumerate the component
// family and each one's declared source without paying for every derivation to run.
func ComponentInputSources() []ComponentInputSource {
	declarations := componentInputDeclarations()
	sources := make([]ComponentInputSource, len(declarations))
	for i, d := range declarations {
		sources[i] = ComponentInputSource{Component: d.component, Source: d.source}
	}
	return sources
}

// ResolveComponentInputs resolves every declared component's inputs against root, keyed by
// phase name. A derivation that fails fails the whole resolution: a partial set names
// fewer inputs than the component reads, which is the shape that buys a wrong skip.
//
// Each underlying derivation runs at most once per call, and the toolchain components share
// one module-wide listing — the resolution happens inside every gate decision, so paying
// for that listing per component would be paying for it five times over.
func ResolveComponentInputs(root string) (map[string]ComponentInputs, error) {
	resolver := &inputResolver{root: root}
	sets := map[string]ComponentInputs{}
	for _, declaration := range componentInputDeclarations() {
		paths, digests, err := declaration.resolve(resolver)
		if err != nil {
			return nil, fmt.Errorf("derive %s inputs: %w", declaration.component, err)
		}
		sets[declaration.component] = ComponentInputs{
			component: declaration.component,
			source:    declaration.source,
			paths:     paths,
			digests:   digests,
		}
	}
	return sets, nil
}

// inputResolver memoizes each derivation for the span of one resolution, error included:
// a listing that failed for one component fails identically for the next, and re-running it
// would only pay the failure again.
type inputResolver struct {
	root string

	modulePaths []string
	moduleErr   error
	moduleDone  bool
}

func (r *inputResolver) moduleClosure() ([]string, []string, error) {
	if !r.moduleDone {
		r.modulePaths, r.moduleErr = moduleTestClosure(r.root)
		if r.moduleErr == nil {
			r.modulePaths, r.moduleErr = withModuleManifest(r.root, r.modulePaths)
		}
		r.moduleDone = true
	}
	return r.modulePaths, nil, r.moduleErr
}

// withModuleManifest adds the module manifest to a module-test closure's derived paths:
// go.mod always, and go.sum only when the module has one. `go list` names neither file, so
// a dependency version bump that leaves every listed source byte-identical would otherwise
// leave the toolchain and contract components blind to a change that can red all of them —
// the same gap freshness.BuildInputs already closes for build by adding the same two files
// beside its own derived closure. The addition is bounded to exactly these two paths; the
// derivation still supplies every source path.
func withModuleManifest(root string, closure []string) ([]string, error) {
	if _, err := os.Lstat(filepath.Join(root, "go.mod")); err != nil {
		return nil, fmt.Errorf("required module manifest %q: %w", "go.mod", err)
	}
	paths := make(map[string]struct{}, len(closure)+2)
	for _, path := range closure {
		paths[path] = struct{}{}
	}
	paths["go.mod"] = struct{}{}
	if _, err := os.Lstat(filepath.Join(root, "go.sum")); err == nil {
		paths["go.sum"] = struct{}{}
	}
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		ordered = append(ordered, path)
	}
	sort.Strings(ordered)
	return slices.Compact(ordered), nil
}

// shellcheckInputs takes its set from shellcheckFiles — the same enumeration
// shellcheckArgv builds its invocation on top of — so the component that lints a file and
// the declaration that says it reads that file cannot disagree: a script the enumeration
// gains is a script the declaration gains, and a flag added to the invocation prefix
// cannot shift where the linted paths begin, because neither side restates that offset.
func (r *inputResolver) shellcheckInputs() ([]string, []string, error) {
	paths := append([]string(nil), shellcheckFiles(r.root)...)
	sort.Strings(paths)
	return slices.Compact(paths), nil, nil
}

// canaryInputs is the registry's one hand-declared entry: canary's inputs have no
// derivable source, so they are named directly. internal/canary/ and tests/canary/ are
// the sources and fixtures the phase grades; bin/bench.sh and .bench/lib/canary-run.sh
// are the wrapper scripts its phase wiring execs; the agent-guidance roots join them
// because canary fixtures seed from the kit tree, so a guidance edit can move a sweep's
// expected diagnostics.
//
// The published binary's digest is deliberately absent, a recorded tripwire narrowing
// rather than an oversight: canary execs the binary but skips on an ordinary Go edit
// anyway, so a binary digest here is never a missing input, it is the scoping this
// feature buys. The profile records the ruling behind the narrowing.
func (r *inputResolver) canaryInputs() ([]string, []string, error) {
	paths := []string{
		".bench/lib/canary-run.sh",
		"bin/bench.sh",
		"internal/canary/",
		"tests/canary/",
	}
	paths = append(paths, agentMarkdownDirectories...)
	sort.Strings(paths)
	return paths, nil, nil
}

func (r *inputResolver) contractInputs() ([]string, []string, error) {
	paths, _, err := r.moduleClosure()
	if err != nil {
		return nil, nil, err
	}
	documents, err := packagesurface.ContractDocumentInputs(r.root)
	if err != nil {
		return nil, nil, fmt.Errorf("derive consumer inventory documents: %w", err)
	}
	paths = append(paths, documents...)
	sort.Strings(paths)
	return slices.Compact(paths), nil, nil
}

// agentMarkdownDirectories is canary's portable-guidance root. The list lives beside its
// declaration rather than the reduced-scope declaration, which never covers these paths.
var agentMarkdownDirectories = []string{".agents/"}

// listedTestPackage is the `go list -json` shape the module-wide closure reads. The file
// groups a test-augmented listing carries are a superset of the build-input groups the
// freshness digest reads: these components grade test sources, ignored-by-constraint
// sources, and testdata, none of which link into the binary.
type listedTestPackage struct {
	Dir               string
	GoFiles           []string
	CgoFiles          []string
	CFiles            []string
	CXXFiles          []string
	MFiles            []string
	HFiles            []string
	FFiles            []string
	SFiles            []string
	SwigFiles         []string
	SwigCXXFiles      []string
	SysoFiles         []string
	EmbedFiles        []string
	TestGoFiles       []string
	XTestGoFiles      []string
	TestEmbedFiles    []string
	XTestEmbedFiles   []string
	IgnoredGoFiles    []string
	IgnoredOtherFiles []string
}

func (p listedTestPackage) files() []string {
	var files []string
	for _, group := range [][]string{
		p.GoFiles, p.CgoFiles, p.CFiles, p.CXXFiles, p.MFiles, p.HFiles, p.FFiles, p.SFiles,
		p.SwigFiles, p.SwigCXXFiles, p.SysoFiles, p.EmbedFiles, p.TestGoFiles, p.XTestGoFiles,
		p.TestEmbedFiles, p.XTestEmbedFiles, p.IgnoredGoFiles, p.IgnoredOtherFiles,
	} {
		files = append(files, group...)
	}
	return files
}

// moduleTestClosure returns the repository-relative paths of every file in root's
// module-wide test closure, plus the testdata/ contents of the packages that closure lists.
// Packages outside root — the standard library, the module cache, and the generated test
// mains go builds under its cache — contribute nothing: they are not what a changeset under
// root can move.
func moduleTestClosure(root string) ([]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	command := exec.Command("go", "list", "-buildvcs=false", "-json", "-deps", "-test", "./...")
	command.Dir = root
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("resolve the module test closure: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	paths := map[string]struct{}{}
	directories := map[string]struct{}{}
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var pkg listedTestPackage
		if err := decoder.Decode(&pkg); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decode the module test closure: %w", err)
		}
		if pkg.Dir == "" || !withinRoot(root, pkg.Dir) {
			continue
		}
		directories[pkg.Dir] = struct{}{}
		for _, name := range pkg.files() {
			// A generated source is listed by absolute path into the build cache; joining
			// it onto Dir would manufacture a path that names no file at all.
			if filepath.IsAbs(name) {
				continue
			}
			if path := filepath.Join(pkg.Dir, name); withinRoot(root, path) {
				paths[path] = struct{}{}
			}
		}
	}
	// The listing names a package's directory several times over — once for the package,
	// once for its test-augmented variant, once for its external test package — so testdata
	// is walked per directory rather than per listing.
	for dir := range directories {
		if err := collectTestdata(root, dir, paths); err != nil {
			return nil, err
		}
	}
	return relativeSorted(root, paths)
}

// collectTestdata adds every regular file under dir's testdata directory. `go list` names
// no file inside testdata — the toolchain hands the directory to the tests wholesale — so
// a component that grades those tests reads it without any listing saying so.
func collectTestdata(root, dir string, paths map[string]struct{}) error {
	testdata := filepath.Join(dir, "testdata")
	err := filepath.WalkDir(testdata, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() && withinRoot(root, path) {
			paths[path] = struct{}{}
		}
		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("walk %s: %w", testdata, err)
	}
	return nil
}

func relativeSorted(root string, paths map[string]struct{}) ([]string, error) {
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, err
		}
		ordered = append(ordered, filepath.ToSlash(rel))
	}
	sort.Strings(ordered)
	return slices.Compact(ordered), nil
}

func withinRoot(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
