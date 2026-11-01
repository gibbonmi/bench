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
func checkStrictJSONOwner(root string) []string {
	var diags []string
	for _, caller := range strictJSONCallers {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(caller.dir))); errors.Is(err, fs.ErrNotExist) {
			continue
		}
		sources, sourceDiags := guardOwnerSources(root, caller.dir, strictJSONOwnerDiagPrefix)
		diags = append(diags, sourceDiags...)
		graphs := map[string]strictJSONGraph{}
		diags = append(diags, guardOwnerParse(root, strictJSONOwnerDiagPrefix, sources, func(fset *token.FileSet, rel string, file *ast.File) []string {
			dir := path.Dir(rel)
			if graphs[dir] == nil {
				graphs[dir] = strictJSONGraph{}
			}
			graphs[dir].add(fset, rel, file)
			return nil
		})...)
		for _, dir := range slices.Sorted(maps.Keys(graphs)) {
			diags = append(diags, graphs[dir].walkers()...)
		}
		if graph := graphs[caller.dir]; graph != nil {
			diags = append(diags, graph.bypasses(caller.dir, caller.readers)...)
		}
	}
	return diags
}

// strictJSONFunc is one function, method, or closure of a package. Its body calls a Token method
// when tokens is true, and the owner when owner is true. calls holds a graph key for each call.
type strictJSONFunc struct {
	name, at      string
	tokens, owner bool
	calls         []any
}

// strictJSONGraph indexes the functions of one package by call key: a function name, a method
// name after a dot, or the variable object of a closure, so two closures of one name stay apart.
type strictJSONGraph map[any][]*strictJSONFunc

// add indexes each function, method, and assigned closure of file.
func (g strictJSONGraph) add(fset *token.FileSet, rel string, file *ast.File) {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		key := any(fn.Name.Name)
		if fn.Recv != nil {
			key = "." + fn.Name.Name
		}
		g[key] = append(g[key], strictJSONNode(fset, rel, file, fn.Name.Name, fn.Pos(), fn.Body))
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			if assign, ok := node.(*ast.AssignStmt); ok && len(assign.Lhs) == len(assign.Rhs) {
				for i, lhs := range assign.Lhs {
					ident, named := lhs.(*ast.Ident)
					literal, closure := assign.Rhs[i].(*ast.FuncLit)
					if named && closure && ident.Obj != nil {
						g[ident.Obj] = append(g[ident.Obj], strictJSONNode(fset, rel, file, ident.Name, literal.Pos(), literal.Body))
					}
				}
			}
			return true
		})
	}
}

// strictJSONNode reads the calls of body. A call through an import qualifier other than the owner
// adds no key, because it cannot reach a function of this package.
func strictJSONNode(fset *token.FileSet, rel string, file *ast.File, name string, pos token.Pos, body *ast.BlockStmt) *strictJSONFunc {
	node := &strictJSONFunc{name: name, at: fmt.Sprintf("%s:%d", rel, fset.Position(pos).Line)}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if guardProjectionImported(file, call, strictJSONOwnerPath) != "" {
			node.owner = true
			return true
		}
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			if fun.Obj != nil && fun.Obj.Kind == ast.Var {
				node.calls = append(node.calls, fun.Obj)
			} else {
				node.calls = append(node.calls, fun.Name)
			}
		case *ast.SelectorExpr:
			node.tokens = node.tokens || fun.Sel.Name == "Token" && len(call.Args) == 0
			if qualifier, ok := fun.X.(*ast.Ident); !ok || qualifier.Obj != nil {
				node.calls = append(node.calls, "."+fun.Sel.Name)
			}
		}
		return true
	})
	return node
}

// reaches reports whether a function that from calls, directly or through other calls, is done.
func (g strictJSONGraph) reaches(from *strictJSONFunc, done func(*strictJSONFunc) bool) bool {
	seen := map[*strictJSONFunc]bool{}
	for stack := []*strictJSONFunc{from}; len(stack) > 0; {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, key := range node.calls {
			for _, next := range g[key] {
				if done(next) {
					return true
				}
				if !seen[next] {
					seen[next] = true
					stack = append(stack, next)
				}
			}
		}
	}
	return false
}

// walkers reports each function that reads a token and calls itself, directly or through other calls.
func (g strictJSONGraph) walkers() []string {
	var diags []string
	for _, nodes := range g {
		for _, node := range nodes {
			recursive := g.reaches(node, func(next *strictJSONFunc) bool { return next == node })
			if node.tokens && recursive {
				diags = append(diags, fmt.Sprintf("%s%s: %s walks JSON tokens recursively outside %s", strictJSONOwnerDiagPrefix, node.at, node.name, strictJSONOwnerDir))
			}
		}
	}
	slices.Sort(diags)
	return diags
}

// bypasses reports each named strict reader of the package at dir that is absent or reaches no owner call.
func (g strictJSONGraph) bypasses(dir string, readers []string) []string {
	var diags []string
	for _, name := range readers {
		nodes := g[name]
		switch {
		case len(nodes) == 0:
			diags = append(diags, fmt.Sprintf("%s%s: strict reader %s is not defined", strictJSONOwnerDiagPrefix, dir, name))
		case !nodes[0].owner && !g.reaches(nodes[0], func(next *strictJSONFunc) bool { return next.owner }):
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

// strictJSONOwnerRoot builds a root from the compliant caller tree, with sources and plants laid over it.
func strictJSONOwnerRoot(t *testing.T, sources map[string]string, plants map[string]func(*testing.T, string)) string {
	t.Helper()
	files := maps.Clone(strictJSONOwnerCompliant)
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
// caller package or beneath it is refused. A strict reader that reaches no owner call is refused.
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
	t.Run("absent and empty packages", func(t *testing.T) {
		for name, files := range map[string]map[string]string{
			"absent": {"internal/other/walk.go": strictJSONOwnerFunctionWalker},
			"empty":  {"internal/gate/README.md": "gate\n", "internal/reviewrecord/keys_test.go": strictJSONOwnerFunctionWalker, "internal/releaseevidence/testdata/keys.go": strictJSONOwnerClosureWalker},
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
	for _, caller := range strictJSONCallers {
		if !exists(filepath.Join(h.KitRoot, filepath.FromSlash(caller.dir))) {
			t.Errorf("caller package %s is absent from the kit root", caller.dir)
		}
	}
	if diags := strictJSONOwnerRun(t, h.KitRoot); len(diags) != 0 {
		t.Fatalf("a strict JSON copy survives outside the owner:\n%s", strings.Join(diags, "\n"))
	}
}
