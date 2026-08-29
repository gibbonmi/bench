package consumers

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"sort"
	"testing"
)

// fixturePkg is one package of in-process fixture source: its import path and its files.
// The core's whole input is typed packages, so a test builds them with go/parser plus
// go/types and never spawns a process. The loader seam's real-path tests drive the real
// loader instead.
type fixturePkg struct {
	path  string
	files map[string]string
}

// importerFunc adapts a closure to types.Importer.
type importerFunc func(string) (*types.Package, error)

func (f importerFunc) Import(path string) (*types.Package, error) { return f(path) }

// typecheckFixture type-checks pkgs in the given order and returns them in the core's
// Package shape. A fixture package resolves against the packages already built; anything
// else falls through to the standard-library importer.
func typecheckFixture(t *testing.T, pkgs []fixturePkg) []*Package {
	t.Helper()
	fset := token.NewFileSet()
	built := map[string]*types.Package{}
	imp := importerFunc(func(path string) (*types.Package, error) {
		if p, ok := built[path]; ok {
			return p, nil
		}
		return importer.Default().Import(path)
	})
	var out []*Package
	for _, fp := range pkgs {
		names := make([]string, 0, len(fp.files))
		for name := range fp.files {
			names = append(names, name)
		}
		sort.Strings(names)
		var files []*ast.File
		for _, name := range names {
			f, err := parser.ParseFile(fset, name, fp.files[name], parser.AllErrors)
			if err != nil {
				t.Fatalf("parse %s: %v", name, err)
			}
			files = append(files, f)
		}
		info := &types.Info{
			Uses:       map[*ast.Ident]types.Object{},
			Defs:       map[*ast.Ident]types.Object{},
			Selections: map[*ast.SelectorExpr]*types.Selection{},
		}
		conf := types.Config{Importer: imp}
		tp, err := conf.Check(fp.path, fset, files, info)
		if err != nil {
			t.Fatalf("typecheck %s: %v", fp.path, err)
		}
		built[fp.path] = tp
		out = append(out, &Package{PkgPath: fp.path, Fset: fset, Files: files, Types: tp, Info: info})
	}
	return out
}

// referenceFixture plants the reference shapes the core must find: a call through a
// renamed import, a value use inside a package-level var, and a use inside a method.
// The line numbers are load-bearing, so the sources are laid out to be read with them.
var referenceFixture = []fixturePkg{
	{path: "example.com/target", files: map[string]string{"/repo/target/target.go": "" +
		/* 1 */ "package target\n" +
		/* 2 */ "\n" +
		/* 3 */ "func Symbol() {}\n" +
		/* 4 */ "\n" +
		/* 5 */ "type T struct{ Run func() }\n" +
		/* 6 */ "\n" +
		/* 7 */ "type Count int\n"}},
	{path: "example.com/consumer", files: map[string]string{"/repo/consumer/consumer.go": "" +
		/* 1 */ "package consumer\n" +
		/* 2 */ "\n" +
		/* 3 */ "import tg \"example.com/target\"\n" +
		/* 4 */ "\n" +
		/* 5 */ "func Direct() { tg.Symbol() }\n" +
		/* 6 */ "\n" +
		/* 7 */ "var someRegistry = []tg.T{{Run: tg.Symbol}}\n" +
		/* 8 */ "\n" +
		/* 9 */ "type Holder struct{}\n" +
		/* 10 */ "\n" +
		/* 11 */ "func (h Holder) Use() { tg.Symbol() }\n" +
		/* 12 */ "\n" +
		/* 13 */ "func Convert(n int) tg.Count { return tg.Count(n) }\n" +
		/* 14 */ "\n" +
		/* 15 */ "func use(f func()) { _ = f }\n" +
		/* 16 */ "\n" +
		/* 17 */ "func Pass() { use(tg.Symbol) }\n"}},
}

// implementsFixture plants an interface, a second interface that structurally satisfies
// it, a value-receiver implementer, a pointer-receiver implementer, a type that
// implements nothing, and one plain reference to the interface name. An interface query
// answers with the two implementers and the reference, and with nothing else.
var implementsFixture = []fixturePkg{
	{path: "example.com/target", files: map[string]string{"/repo/target/target.go": "" +
		/* 1 */ "package target\n" +
		/* 2 */ "\n" +
		/* 3 */ "type Runner interface{ Run() }\n" +
		/* 4 */ "\n" +
		/* 5 */ "type Other interface {\n" +
		/* 6 */ "\tRun()\n" +
		/* 7 */ "\tStop()\n" +
		/* 8 */ "}\n"}},
	{path: "example.com/consumer", files: map[string]string{"/repo/consumer/consumer.go": "" +
		/* 1 */ "package consumer\n" +
		/* 2 */ "\n" +
		/* 3 */ "import \"example.com/target\"\n" +
		/* 4 */ "\n" +
		/* 5 */ "type Value struct{}\n" +
		/* 6 */ "\n" +
		/* 7 */ "func (Value) Run() {}\n" +
		/* 8 */ "\n" +
		/* 9 */ "type Pointer struct{}\n" +
		/* 10 */ "\n" +
		/* 11 */ "func (p *Pointer) Run() {}\n" +
		/* 12 */ "\n" +
		/* 13 */ "type Nope struct{}\n" +
		/* 14 */ "\n" +
		/* 15 */ "func (Nope) Walk() {}\n" +
		/* 16 */ "\n" +
		/* 17 */ "var Check target.Runner = Value{}\n"}},
}

// aliasFixture plants one alias declaration and a consumer spelled each way, so the two
// spellings of one declaration can be compared byte for byte.
var aliasFixture = []fixturePkg{
	{path: "example.com/target", files: map[string]string{"/repo/target/target.go": "" +
		/* 1 */ "package target\n" +
		/* 2 */ "\n" +
		/* 3 */ "type Origin struct{}\n" +
		/* 4 */ "\n" +
		/* 5 */ "type Alias = Origin\n"}},
	{path: "example.com/consumer", files: map[string]string{"/repo/consumer/consumer.go": "" +
		/* 1 */ "package consumer\n" +
		/* 2 */ "\n" +
		/* 3 */ "import \"example.com/target\"\n" +
		/* 4 */ "\n" +
		/* 5 */ "func UseAlias() { var a target.Alias; _ = a }\n" +
		/* 6 */ "\n" +
		/* 7 */ "func UseOrigin() { var o target.Origin; _ = o }\n"}},
}

// genericFixture plants a generic type with a method and a generic struct field, each
// consumed through two instantiations. Row identity is the generic origin, so the two
// instantiations of one declaration must read as one declaration.
var genericFixture = []fixturePkg{
	{path: "example.com/target", files: map[string]string{"/repo/target/target.go": "" +
		/* 1 */ "package target\n" +
		/* 2 */ "\n" +
		/* 3 */ "type Box[T any] struct{ First T }\n" +
		/* 4 */ "\n" +
		/* 5 */ "func (Box[T]) Get() T { var z T; return z }\n"}},
	{path: "example.com/consumer", files: map[string]string{"/repo/consumer/consumer.go": "" +
		/* 1 */ "package consumer\n" +
		/* 2 */ "\n" +
		/* 3 */ "import \"example.com/target\"\n" +
		/* 4 */ "\n" +
		/* 5 */ "func UseInt() int { return target.Box[int]{}.Get() }\n" +
		/* 6 */ "\n" +
		/* 7 */ "func UseString() string { return target.Box[string]{}.Get() }\n" +
		/* 8 */ "\n" +
		/* 9 */ "func FieldInt() int { return target.Box[int]{}.First }\n" +
		/* 10 */ "\n" +
		/* 11 */ "func FieldString() string { return target.Box[string]{}.First }\n"}},
}
