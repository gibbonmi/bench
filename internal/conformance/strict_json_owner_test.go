package conformance

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"golang.org/x/mod/modfile"
)

// strictJSONOwnerDir is the one package whose code walks JSON tokens recursively to refuse a
// duplicate key. strictJSONOwnerPath is its import path.
const (
	strictJSONOwnerDir        = "internal/jsonfile"
	strictJSONOwnerPath       = benchModule + strictJSONOwnerDir
	strictJSONOwnerDiagPrefix = "strict JSON owner "
)

// strictJSONCallers names each migrated caller package and the strict readers in it that must
// reach the owner. The walker rule grades each package and the packages beneath it. A reader
// outside this table, such as a stream reader or a partial projection, keeps its own decoder.
var strictJSONCallers = []struct {
	dir     string
	readers []string
}{
	{"internal/releasepreflight", []string{"ValidateVulnerabilityPolicy"}},
	{"internal/gate", []string{"strictJSON"}},
	{"internal/reviewrecord", []string{"decode"}},
	{"internal/releaseevidence", []string{"decodeStrict", "decodeDocument", "decodeExactObject", "validatePackageEvidence"}},
}

// checkStrictJSONOwner refuses a recursive JSON token walker in the production Go beneath each
// caller package, because only the owner walks object keys. It also refuses each named strict
// reader that reaches no owner call. An absent caller package is an absent optional surface.
// A root whose module is not the Bench kit has no caller package.
func checkStrictJSONOwner(root string) []string {
	if !benchModuleRoot(root) {
		return nil
	}
	var diags []string
	for _, caller := range strictJSONCallers {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(caller.dir))); errors.Is(err, fs.ErrNotExist) {
			continue
		}
		sources, sourceDiags := guardOwnerSources(root, caller.dir, strictJSONOwnerDiagPrefix)
		diags = append(diags, sourceDiags...)
		packages := map[string]packageFuncs{}
		diags = append(diags, guardOwnerParse(root, strictJSONOwnerDiagPrefix, sources, func(fset *token.FileSet, rel string, file *ast.File) []string {
			dir := path.Dir(rel)
			if packages[dir] == nil {
				packages[dir] = packageFuncs{}
			}
			packages[dir].add(fset, rel, file)
			return nil
		})...)
		for _, dir := range slices.Sorted(maps.Keys(packages)) {
			diags = append(diags, strictJSONWalkers(packages[dir])...)
		}
		if functions := packages[caller.dir]; functions != nil {
			diags = append(diags, strictJSONBypasses(functions, caller.dir, caller.readers)...)
		}
	}
	return diags
}

// benchModuleRoot reports whether the go.mod of root names the Bench kit module.
func benchModuleRoot(root string) bool {
	return modfile.ModulePath([]byte(readIfExists(filepath.Join(root, "go.mod"))))+"/" == benchModule
}

// strictJSONCalls reports whether the body of fn holds a call that match accepts.
func strictJSONCalls(fn *packageFunc, match func(*ast.CallExpr) bool) bool {
	found := false
	ast.Inspect(fn.body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		found = found || ok && match(call)
		return !found
	})
	return found
}

// strictJSONWalkers reports each function that reads a token and reaches itself, directly or through other callees.
func strictJSONWalkers(functions packageFuncs) []string {
	var diags []string
	readsToken := func(call *ast.CallExpr) bool {
		selector, ok := call.Fun.(*ast.SelectorExpr)
		return ok && selector.Sel.Name == "Token" && len(call.Args) == 0
	}
	for _, fns := range functions {
		for _, fn := range fns {
			if strictJSONCalls(fn, readsToken) && functions.reaches(fn, func(next *packageFunc) bool { return next == fn }) {
				diags = append(diags, fmt.Sprintf("%s%s:%d: %s walks JSON tokens recursively outside %s", strictJSONOwnerDiagPrefix, fn.rel, fn.line, fn.name, strictJSONOwnerDir))
			}
		}
	}
	slices.Sort(diags)
	return diags
}

// strictJSONBypasses reports each named strict reader of the package at dir that is absent or reaches no owner call.
func strictJSONBypasses(functions packageFuncs, dir string, readers []string) []string {
	owner := func(fn *packageFunc) bool {
		return strictJSONCalls(fn, func(call *ast.CallExpr) bool {
			return guardProjectionImported(fn.file, call, strictJSONOwnerPath) != ""
		})
	}
	var diags []string
	for _, name := range readers {
		named := functions.named(name)
		switch {
		case len(named) == 0:
			diags = append(diags, fmt.Sprintf("%s%s: strict reader %s is not defined", strictJSONOwnerDiagPrefix, dir, name))
		case !owner(named[0]) && !functions.reaches(named[0], owner):
			diags = append(diags, fmt.Sprintf("%s%s: strict reader %s does not reach %s", strictJSONOwnerDiagPrefix, dir, name, strictJSONOwnerDir))
		}
	}
	return diags
}

// strictJSONOwnerCheck is the registry name of the check that the registered-path test runs.
const strictJSONOwnerCheck = "strict-json-owner"

// strictJSONOwnerImport imports the owner in a fixture source. Each fixture packs its declarations
// tightly, because only its syntax nodes matter.
const strictJSONOwnerImport = `"` + strictJSONOwnerPath + `"`

// strictJSONOwnerCompliant holds one caller tree in which each strict reader reaches the owner.
// The gate and release readers keep their local adapters: a typed decode that keeps precedence,
// an untyped decode refusal, and a partial identity projection.
var strictJSONOwnerCompliant = map[string]string{
	"internal/releasepreflight/vulnerability.go": `package releasepreflight; import ` + strictJSONOwnerImport + `
func ValidateVulnerabilityPolicy(findings []string, data []byte) error { var exceptions []struct{ ID string }; return jsonfile.DecodeDocument(data, &exceptions) }`,
	"internal/gate/verdict.go": `package gate; import ("bytes"; "encoding/json"; "errors"; ` + strictJSONOwnerImport + `)
func strictJSON(data []byte, dst any) error { dec := json.NewDecoder(bytes.NewReader(data)); dec.DisallowUnknownFields(); if err := dec.Decode(dst); err != nil { return err }
	if jsonfile.RejectTrailing(dec) != nil { return errors.New("trailing JSON") }; return jsonfile.DecodeDocument(data, new(json.RawMessage)) }`,
	"internal/reviewrecord/parse.go": `package reviewrecord; import ("bytes"; "encoding/json"; "errors"; ` + strictJSONOwnerImport + `)
func decode(data []byte, value any) error { err := jsonfile.DecodeDocument(data, value)
	if errors.Is(err, jsonfile.ErrTrailingData) { decoder := json.NewDecoder(bytes.NewReader(data)); decoder.DisallowUnknownFields(); return decoder.Decode(value) }; return err }`,
	"internal/releaseevidence/json_validation.go": `package releaseevidence; import ("encoding/json"; "errors"; ` + strictJSONOwnerImport + `)
func decodeStrict(data []byte, value any) error { if !json.Valid(data) { return errors.New("invalid JSON") }; return decodeDocument(data, value) }
func decodeDocument(data []byte, value any) error { if err := jsonfile.DecodeDocument(data, value); err != nil { return err }; return json.Unmarshal(data, new(any)) }`,
	"internal/releaseevidence/component_manifest.go": `package releaseevidence; import "encoding/json"
func decodeExactObject(data []byte, fields []string) (object map[string]json.RawMessage, err error) { err = decodeDocument(data, &object); return object, err }`,
	"internal/releaseevidence/package_artifact.go": `package releaseevidence; import "encoding/json"
func validatePackageEvidence(data []byte, version string) bool { var identity struct{ Version string }
	return decodeDocument(data, new(json.RawMessage)) == nil && json.Unmarshal(data, &identity) == nil && identity.Version == version }`,
}

// strictJSONOwnerKitMod is the go.mod of a fixture root whose module is the Bench kit.
var strictJSONOwnerKitMod = "module " + strings.TrimSuffix(benchModule, "/") + "\n"

// strictJSONOwnerRoot builds a kit root from the compliant caller tree, with sources and plants laid over it.
func strictJSONOwnerRoot(t *testing.T, sources map[string]string, plants map[string]func(*testing.T, string)) string {
	t.Helper()
	files := maps.Clone(strictJSONOwnerCompliant)
	files["go.mod"] = strictJSONOwnerKitMod
	maps.Copy(files, sources)
	return throwawayRoot{files: files, plants: plants}.build(t)
}

// strictJSONOwnerRun runs the bound strict-json-owner check on root.
func strictJSONOwnerRun(t *testing.T, root string) []string {
	t.Helper()
	binding, bound := conformanceChecks[strictJSONOwnerCheck]
	if !bound {
		t.Fatalf("%s conformance owner is not bound", strictJSONOwnerCheck)
	}
	return binding.run(root, root, registry.Dev)
}

// strictJSONOwnerWant gives the expected diagnostics in order, each with the check prefix.
func strictJSONOwnerWant(diags ...string) []string {
	want := make([]string, len(diags))
	for i, diag := range diags {
		want[i] = strictJSONOwnerDiagPrefix + diag
	}
	return want
}

// Each walker below reads object keys and calls itself. The closure is the retired release
// walker, renamed, on line 3. The function and the method value are on line 2.
const (
	strictJSONOwnerClosureWalker = `package releaseevidence; import ("bytes"; "encoding/json")
func uniqueKeys(data []byte) error { decoder := json.NewDecoder(bytes.NewReader(data)); var visit func() error
	visit = func() error { token, err := decoder.Token(); if delim, ok := token.(json.Delim); err == nil && ok && delim == '{' { for decoder.More() { if _, err = decoder.Token(); err == nil { err = visit() }; if err != nil { return err } }; _, err = decoder.Token() }; return err }
	return visit() }`
	strictJSONOwnerFunctionWalker = `package reviewrecord; import "encoding/json"
func keysOnce(d *json.Decoder) error { token, err := d.Token(); if delim, ok := token.(json.Delim); err == nil && ok && delim == '{' { for d.More() { if _, err = d.Token(); err == nil { err = keysOnce(d) }; if err != nil { return err } }; _, err = d.Token() }; return err }`
	strictJSONOwnerMethodWalker = `package cache; import "encoding/json"; type walker struct{ d *json.Decoder }
func (w walker) value() error { token, err := w.d.Token(); if delim, ok := token.(json.Delim); err == nil && ok { return w.members(delim) }; return err }
func (w walker) members(delim json.Delim) error { for w.d.More() { if err := w.value(); err != nil { return err } }; return nil }`
)

// strictJSONOwnerBypass is the release decodeStrict with its walker deleted and no owner call. A
// comment and a string name the owner, and decodeDocument still reaches it. strictJSONOwnerErrorOnly
// is a gate strictJSON that reads an owner error value but calls no owner function.
const (
	strictJSONOwnerBypass = `package releaseevidence; import ("bytes"; "encoding/json"; "errors"; "io"; ` + strictJSONOwnerImport + `)
// decodeStrict once returned jsonfile.DecodeDocument(data, value).
func decodeStrict(data []byte, value any) error { if !json.Valid(data) { return errors.New("invalid JSON: jsonfile.DecodeDocument") }; decoder := json.NewDecoder(bytes.NewReader(data)); decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil { return err }; if decoder.Decode(new(any)) != io.EOF { return errors.New("trailing JSON") }; return nil }
func decodeDocument(data []byte, value any) error { return jsonfile.DecodeDocument(data, value) }`
	strictJSONOwnerErrorOnly = `package gate; import ("bytes"; "encoding/json"; "errors"; ` + strictJSONOwnerImport + `)
func strictJSON(data []byte, dst any) error { dec := json.NewDecoder(bytes.NewReader(data)); dec.DisallowUnknownFields()
	if err := dec.Decode(dst); errors.Is(err, jsonfile.ErrDuplicateKey) { return errors.New("duplicate name") }; return nil }`
)

// TestStrictJSONOwnerBites is the bite proof of the bound check. A renamed recursive walker in a
// caller package or beneath it is refused, also when it recurses through a function or method value.
// A strict reader that reaches no owner call is refused.
func TestStrictJSONOwnerBites(t *testing.T) {
	const outside = " walks JSON tokens recursively outside internal/jsonfile"
	for _, tc := range []struct {
		name    string
		sources map[string]string
		plants  map[string]func(*testing.T, string)
		want    []string
	}{
		{"renamed function walker", map[string]string{"internal/reviewrecord/keys.go": strictJSONOwnerFunctionWalker}, nil, strictJSONOwnerWant("internal/reviewrecord/keys.go:2: keysOnce" + outside)},
		{"closure walker", map[string]string{"internal/releaseevidence/keys.go": strictJSONOwnerClosureWalker}, nil, strictJSONOwnerWant("internal/releaseevidence/keys.go:3: visit" + outside)},
		{"method walker beneath a caller package", map[string]string{"internal/gate/cache/walk.go": strictJSONOwnerMethodWalker}, nil, strictJSONOwnerWant("internal/gate/cache/walk.go:2: value" + outside)},
		{"function value walker", map[string]string{"internal/reviewrecord/keys.go": strings.Replace(strictJSONOwnerFunctionWalker, "err = keysOnce(d)", "next := keysOnce; err = next(d)", 1)}, nil,
			strictJSONOwnerWant("internal/reviewrecord/keys.go:2: keysOnce" + outside)},
		{"method value walker", map[string]string{"internal/gate/cache/walk.go": strings.Replace(strictJSONOwnerMethodWalker, "return w.members(delim)", "next := w.members; return next(delim)", 1)}, nil,
			strictJSONOwnerWant("internal/gate/cache/walk.go:2: value" + outside)},
		{"readers that bypass the owner", map[string]string{"internal/gate/verdict.go": strictJSONOwnerErrorOnly, "internal/releaseevidence/json_validation.go": strictJSONOwnerBypass}, nil,
			strictJSONOwnerWant("internal/gate: strict reader strictJSON does not reach internal/jsonfile", "internal/releaseevidence: strict reader decodeStrict does not reach internal/jsonfile")},
		{"undefined reader", map[string]string{"internal/releasepreflight/vulnerability.go": "package releasepreflight; func validatePolicy() {}"}, nil,
			strictJSONOwnerWant("internal/releasepreflight: strict reader ValidateVulnerabilityPolicy is not defined")},
		{"unparsable source", map[string]string{"internal/gate/broken.go": "package gate; func ("}, nil, strictJSONOwnerWant(guardOwnerUnparsed + "internal/gate/broken.go")},
		{"unreadable source", nil, map[string]func(*testing.T, string){"internal/gate/locked.go": func(t *testing.T, path string) {
			writeFixtureFile(t, path, "package gate\n")
			_ = guardProjectionStrip(t, path)
		}}, strictJSONOwnerWant(guardOwnerUnparsed + "internal/gate/locked.go")},
		{"special source", nil, map[string]func(*testing.T, string){"internal/gate/fifo.go": hostileSkillPlanters["fifo"]}, strictJSONOwnerWant("internal/gate/fifo.go" + guardOwnerIrregular)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if diags := strictJSONOwnerRun(t, strictJSONOwnerRoot(t, tc.sources, tc.plants)); !slices.Equal(diags, tc.want) {
				t.Fatalf("diagnostics =\n%s\nwant\n%s", strings.Join(diags, "\n"), strings.Join(tc.want, "\n"))
			}
		})
	}
	// A linked module's packages only share a name with the caller packages, so the check grades none.
	t.Run("absent empty and linked packages", func(t *testing.T) {
		for name, files := range map[string]map[string]string{
			"absent": {"go.mod": strictJSONOwnerKitMod, "internal/other/walk.go": strictJSONOwnerFunctionWalker},
			"empty":  {"go.mod": strictJSONOwnerKitMod, "internal/gate/README.md": "gate\n", "internal/reviewrecord/keys_test.go": strictJSONOwnerFunctionWalker, "internal/releaseevidence/testdata/keys.go": strictJSONOwnerClosureWalker},
			"linked": {"go.mod": "module example.com/linked\n", "internal/gate/gate.go": "package gate; func Run() {}", "internal/reviewrecord/keys.go": strictJSONOwnerFunctionWalker},
		} {
			if diags := strictJSONOwnerRun(t, throwawayRoot{files: files}.build(t)); len(diags) != 0 {
				t.Errorf("%s packages: diagnostics = %q, want none", name, diags)
			}
		}
	})
}

// strictJSONOwnerTypedReader refuses an unknown field and a suffix but accepts a duplicate key, as
// the freshness seal and prospective record readers do.
func strictJSONOwnerTypedReader(pkg, name string) string {
	return fmt.Sprintf(`package %s; import ("bytes"; "encoding/json"; "errors")
func %s(data []byte, record any) error { decoder := json.NewDecoder(bytes.NewReader(data)); decoder.DisallowUnknownFields(); if err := decoder.Decode(record); err != nil { return err }
	if decoder.More() { return errors.New("trailing record data") }; return nil }`, pkg, name)
}

// strictJSONOwnerStreamReader reads a stream of JSON values, as the govulncheck and go list readers do.
func strictJSONOwnerStreamReader(pkg, name string) string {
	return fmt.Sprintf(`package %s; import ("bytes"; "encoding/json"; "io")
func %s(stream []byte) (count int, err error) { decoder := json.NewDecoder(bytes.NewReader(stream))
	for { var event struct{ ID string }; if err = decoder.Decode(&event); err == io.EOF { return count, nil } else if err != nil { return count, err }; count++ } }`, pkg, name)
}

// strictJSONOwnerLocator reads array element offsets after one token read, with no recursion, as
// the canary pin line locator does.
func strictJSONOwnerLocator(pkg string) string {
	return fmt.Sprintf(`package %s; import ("bytes"; "encoding/json")
func elementLines(data []byte, count int) []int { lines := make([]int, count); decoder := json.NewDecoder(bytes.NewReader(data)); if _, err := decoder.Token(); err != nil { return lines }
	for i := range lines { start := int(decoder.InputOffset()); var element json.RawMessage; if err := decoder.Decode(&element); err != nil { return lines }; lines[i] = bytes.Count(data[:start], []byte("\n")) + 1 }
	return lines }`, pkg)
}

// strictJSONOwnerTextWalker is the source text of a recursive walker.
const strictJSONOwnerTextWalker = "func again(d *json.Decoder) error { d.Token(); return again(d) }"

// TestStrictJSONOwnerNegativeControls plants each named preservation exclusion beside the compliant
// caller tree. The bound check permits each one.
func TestStrictJSONOwnerNegativeControls(t *testing.T) {
	for name, sources := range map[string]map[string]string{
		"freshness parseSeal":               {"internal/freshness/seal.go": strictJSONOwnerTypedReader("freshness", "parseSeal")},
		"freshness buildInputs":             {"internal/freshness/inputs.go": strictJSONOwnerStreamReader("freshness", "buildInputs")},
		"prospectiveartifact readRecordAt":  {"internal/gate/prospectiveartifact/prospectiveartifact.go": strictJSONOwnerTypedReader("prospectiveartifact", "readRecordAt")},
		"releasepreflight findingIDs":       {"internal/releasepreflight/scan.go": strictJSONOwnerStreamReader("releasepreflight", "findingIDs")},
		"partial projection":                {"internal/gate/fields.go": `package gate; import "encoding/json"; func manifestFields(data []byte) (fields map[string]any, err error) { err = json.Unmarshal(data, &fields); return fields, err }`},
		"canary elementLines":               {"internal/canary/pins.go": strictJSONOwnerLocator("canary")},
		"token locator in a caller package": {"internal/releaseevidence/lines.go": strictJSONOwnerLocator("releaseevidence")},
		"local adapters":                    {},
		"tests comments and strings": {
			"internal/gate/walk_test.go": "package gate; import \"encoding/json\"\n" + strictJSONOwnerTextWalker,
			"internal/gate/notes.go":     "package gate\n// " + strictJSONOwnerTextWalker + "\nconst walker = `" + strictJSONOwnerTextWalker + "`\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if diags := strictJSONOwnerRun(t, strictJSONOwnerRoot(t, sources, nil)); len(diags) != 0 {
				t.Fatalf("diagnostics = %q, want none", diags)
			}
		})
	}
}

// TestStrictJSONOwnerRegisteredPath runs strict-json-owner through root conformance on a planted
// walker. It requires the exact planted diagnostic, and none after the walker is removed.
func TestStrictJSONOwnerRegisteredPath(t *testing.T) {
	const planted = "internal/releaseevidence/keys.go"
	root := strictJSONOwnerRoot(t, map[string]string{planted: strictJSONOwnerClosureWalker}, nil)
	want := strictJSONOwnerWant(planted + ":3: visit walks JSON tokens recursively outside internal/jsonfile")
	if diags := RunConformance(root, root, registry.Dev, strictJSONOwnerCheck); !slices.Equal(diags, want) {
		t.Fatalf("root conformance over the planted walker = %q, want %q", diags, want)
	}
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(planted))); err != nil {
		t.Fatal(err)
	}
	if diags := RunConformance(root, root, registry.Dev, strictJSONOwnerCheck); len(diags) != 0 {
		t.Fatalf("root conformance after the walker is removed = %q, want none", diags)
	}
}

// TestStrictJSONOwnerHoldsOnTheLiveTree is the live-tree assertion. Each caller package is present
// in the kit root, so the check grades it, and the check reports nothing there.
func TestStrictJSONOwnerHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	if !benchModuleRoot(h.KitRoot) {
		t.Errorf("the kit root module is not %s", benchModule)
	}
	for _, caller := range strictJSONCallers {
		if !exists(filepath.Join(h.KitRoot, filepath.FromSlash(caller.dir))) {
			t.Errorf("caller package %s is absent from the kit root", caller.dir)
		}
	}
	if diags := strictJSONOwnerRun(t, h.KitRoot); len(diags) != 0 {
		t.Fatalf("a strict JSON copy survives outside the owner:\n%s", strings.Join(diags, "\n"))
	}
}
