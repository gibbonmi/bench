package worktree

// This file owns the single-read census. A census entry is an exported function or
// method. The read set is each effects.go function, each package-qualified function that
// effects.go calls, and each exported gate function that reaches such a gate function. A
// kind is one package-qualified call that a read-set name reaches. An entry can read
// each kind once in its own body, outside each function literal and loop body.

import (
	"fmt"
	"go/ast"
	"go/token"
	"maps"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

// effectsFile is the effect boundary. The census reads it for the read set only.
const effectsFile = "effects.go"

// readFinding selects one of the three census messages.
type readFinding int

const (
	readBelowEntry readFinding = iota
	readSecondTime
	readThroughEntry
)

// singleReadReport renders one census message. A read below an entry names the source
// spelling, a second read names its kind, and a call names the entry and one kind.
func singleReadReport(file string, line int, declaration string, finding readFinding, name, kind string) string {
	where := fmt.Sprintf("%s:%d: %s", file, line, declaration)
	switch finding {
	case readSecondTime:
		return fmt.Sprintf("%s reads %s a second time", where, kind)
	case readThroughEntry:
		return fmt.Sprintf("%s calls %s, which reads %s", where, name, kind)
	}
	return fmt.Sprintf("%s reads %s below a census entry", where, name)
}

// qualifiedName spells expr as package.Name when expr selects from a package that
// imports names, and is empty otherwise.
func qualifiedName(expr ast.Expr, imports map[string]bool) string {
	if selector, ok := expr.(*ast.SelectorExpr); ok {
		if pkg, ok := selector.X.(*ast.Ident); ok && imports[pkg.Name] {
			return pkg.Name + "." + selector.Sel.Name
		}
	}
	return ""
}

// topLevelFuncs indexes each function that files declare by name. A method is not
// indexed, because the census follows calls by bare identifier.
func topLevelFuncs(files []*ast.File) map[string]*ast.FuncDecl {
	funcs := map[string]*ast.FuncDecl{}
	for _, file := range files {
		for _, node := range file.Decls {
			if decl, ok := node.(*ast.FuncDecl); ok && decl.Recv == nil && decl.Body != nil {
				funcs[decl.Name.Name] = decl
			}
		}
	}
	return funcs
}

// reachedFuncs adds name and each function of funcs that name reaches by bare calls.
func reachedFuncs(name string, funcs map[string]*ast.FuncDecl, reached map[string]bool) map[string]bool {
	if decl, ok := funcs[name]; ok && !reached[name] {
		reached[name] = true
		for _, callee := range directCallees(decl) {
			reachedFuncs(callee, funcs, reached)
		}
	}
	return reached
}

// readSet maps each read-set name to the sorted kinds that it reaches. A gate function
// takes the qualifier that effects.go spells for the gate package.
func readSet(dir, gateDir string) (map[string][]string, error) {
	files, _, _, err := parseGoFiles(dir, func(name string) bool { return name == effectsFile })
	if err != nil {
		return nil, err
	}
	reads, effects := map[string][]string{}, topLevelFuncs(files)
	for _, file := range files {
		imports := fileImportNames(file)
		for name := range effects {
			reads[name] = []string{}
			for reached := range reachedFuncs(name, effects, map[string]bool{}) {
				ast.Inspect(effects[reached].Body, func(node ast.Node) bool {
					if call, ok := node.(*ast.CallExpr); ok {
						if kind := qualifiedName(call.Fun, imports); kind != "" {
							reads[name], reads[kind] = append(reads[name], kind), []string{kind}
						}
					}
					return true
				})
			}
		}
	}
	gateFiles, err := parseSourceFiles(gateDir)
	if err != nil {
		return nil, err
	}
	gate := topLevelFuncs(gateFiles)
	for _, kind := range slices.Collect(maps.Keys(reads)) {
		pkg, target, _ := strings.Cut(kind, ".")
		for name := range gate {
			if gate[target] != nil && gateFiles[0].Name.Name == pkg && ast.IsExported(name) && reachedFuncs(name, gate, map[string]bool{})[target] {
				reads[pkg+"."+name] = append(reads[pkg+"."+name], kind)
			}
		}
	}
	for name, kinds := range reads {
		slices.Sort(kinds)
		reads[name] = slices.Compact(kinds)
	}
	return reads, nil
}

// readRef is one reference to a read-set name, with the kinds that it reaches, or to an
// entry. nested marks a reference in a function literal or a loop, which can run twice.
type readRef struct {
	pos    token.Pos
	name   string
	kinds  []string
	entry  bool
	nested bool
}

// declRefs returns each reference in root to a read-set name or to an entry, in source
// order. A selector names a package only when the file imports it. A field, and a
// struct literal key, named like a read is not the read.
func declRefs(root ast.Node, imports map[string]bool, reads map[string][]string, funcs map[string]*ast.FuncDecl) []readRef {
	var refs []readRef
	var walk func(ast.Node, bool)
	walk = func(root ast.Node, nested bool) {
		if root == nil {
			return
		}
		ast.Inspect(root, func(node ast.Node) bool {
			name := ""
			switch node := node.(type) {
			case *ast.FuncLit:
				walk(node.Body, true)
				return false
			case *ast.ForStmt:
				walk(node.Init, nested)
				walk(node.Cond, true)
				walk(node.Post, true)
				walk(node.Body, true)
				return false
			case *ast.RangeStmt:
				walk(node.X, nested)
				walk(node.Body, true)
				return false
			case *ast.KeyValueExpr:
				if _, ok := node.Key.(*ast.Ident); ok {
					walk(node.Value, nested)
					return false
				}
			case *ast.SelectorExpr:
				if name = qualifiedName(node, imports); name == "" {
					walk(node.X, nested)
					return false
				}
			case *ast.Ident:
				name = node.Name
			}
			if kinds, read := reads[name]; read || funcs[name] != nil && ast.IsExported(name) {
				refs = append(refs, readRef{pos: node.Pos(), name: name, kinds: kinds, entry: !read, nested: nested})
			}
			return name == ""
		})
	}
	walk(root, false)
	return refs
}

// singleReadCensus reports each read in the non-test files of dir, other than
// effectsFile, that is not the first read of its kind in an entry's own body. It also
// reports each reference from another declaration to an entry whose own body reads.
// gateDir is the gate package directory. The reports are sorted.
func singleReadCensus(dir, gateDir string) ([]string, error) {
	reads, err := readSet(dir, gateDir)
	if err != nil {
		return nil, err
	}
	files, fset, _, err := parseGoFiles(dir, func(name string) bool {
		return name != effectsFile && !strings.HasSuffix(name, "_test.go") && strings.HasSuffix(name, ".go")
	})
	if err != nil {
		return nil, err
	}
	funcs := topLevelFuncs(files)
	type graded struct {
		name  string
		entry bool
		refs  []readRef
	}
	var decls []graded
	entryKinds := map[string][]string{}
	for _, file := range files {
		imports := fileImportNames(file)
		for _, node := range file.Decls {
			if decl, ok := node.(*ast.FuncDecl); ok && decl.Body != nil {
				refs := declRefs(decl.Body, imports, reads, funcs)
				decls = append(decls, graded{decl.Name.Name, decl.Name.IsExported(), refs})
				for _, ref := range refs {
					if decl.Recv == nil {
						entryKinds[decl.Name.Name] = append(entryKinds[decl.Name.Name], ref.kinds...)
					}
				}
			} else if decl, ok := node.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					if value, ok := spec.(*ast.ValueSpec); ok {
						for _, expr := range value.Values {
							decls = append(decls, graded{specName(spec), false, declRefs(expr, imports, reads, funcs)})
						}
					}
				}
			}
		}
	}
	var reports []string
	for _, decl := range decls {
		seen := map[string]bool{}
		for _, ref := range decl.refs {
			position := fset.Position(ref.pos)
			report := func(finding readFinding, name, kind string) {
				reports = append(reports, singleReadReport(filepath.Base(position.Filename), position.Line, decl.name, finding, name, kind))
			}
			switch {
			case ref.entry && !decl.entry:
				for _, kind := range entryKinds[ref.name] {
					report(readThroughEntry, ref.name, kind)
				}
			case ref.entry:
				// The called entry's reads are its own, so an entry may call it.
			case ref.nested || !decl.entry:
				report(readBelowEntry, ref.name, "")
			default:
				for _, kind := range ref.kinds {
					if seen[kind] {
						report(readSecondTime, "", kind)
					}
					seen[kind] = true
				}
			}
		}
	}
	// An entry that reads one kind twice yields one helper-call report for that kind.
	sort.Strings(reports)
	return slices.Compact(reports), nil
}

// syntheticEffectsFile mirrors effects.go: currentTime reads the clock, newAmbient the
// kit and the clock, and Home the Bench home. Its os import serves an added read.
const syntheticEffectsFile = `package worktree

import (
	"os"
	"time"
	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gate"
)

func currentTime() time.Time { return time.Now() }
func newAmbient(home string) ambient { return ambient{kit: gate.KitValue(), now: currentTime()} }
func Home() string { return benchhome.Dir() }
`

// syntheticGateFile is the gate directory of a synthetic package. KitDir calls
// KitValue, KitRoot reaches it through kitRoot, and LaneForCommitAtKit reaches no read.
const syntheticGateFile = `package gate

import "os"

func KitValue() string { return os.Getenv("BENCH_KIT") }
func KitDir() string { return KitValue() }
func KitRoot() string { return kitRoot() }
func kitRoot() string { return KitValue() }
func LaneForCommitAtKit(root, kit string) string { return kit }
`

// wantReadCensus runs the census over effects, syntheticGateFile, and a reader.go whose
// source begins on line 3, after the package clause. The reports must equal want.
func wantReadCensus(t *testing.T, effects, source string, want ...string) {
	t.Helper()
	dir := plantTestFiles(t, map[string]string{effectsFile: effects, "reader.go": "package worktree\n\n" + source})
	reports, err := singleReadCensus(dir, plantTestFiles(t, map[string]string{"gate.go": syntheticGateFile}))
	if err != nil {
		t.Fatalf("single-read census: %v", err)
	}
	if !slices.Equal(reports, want) {
		t.Fatalf("single-read census = %q, want exactly %q", reports, want)
	}
}

func TestSingleReadCensusAcceptsOneReadInAnEntry(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run(paths []string) {\n\tnow := currentTime()\n\tfor range paths {\n\t\t_ = now\n\t}\n\t_ = func() { _ = now }\n}\n")
}

func TestSingleReadCensusRefusesAReadInAnUnexportedFunction(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func stamp() { _ = currentTime() }\n",
		singleReadReport("reader.go", 3, "stamp", readBelowEntry, "currentTime", ""))
}

// TestSingleReadCensusRefusesAReadInAFunctionLiteral proves a literal is no shelter.
func TestSingleReadCensusRefusesAReadInAFunctionLiteral(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run() {\n\tstamp := func() { _ = currentTime() }\n\tstamp()\n}\n",
		singleReadReport("reader.go", 4, "Run", readBelowEntry, "currentTime", ""))
}

func TestSingleReadCensusRefusesAReadInALoopBody(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run(paths []string) {\n\tfor range paths {\n\t\t_ = Home()\n\t}\n}\n",
		singleReadReport("reader.go", 5, "Run", readBelowEntry, "Home", ""))
}

func TestSingleReadCensusRefusesASecondRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run() {\n\t_ = currentTime()\n\t_ = currentTime()\n}\n",
		singleReadReport("reader.go", 5, "Run", readSecondTime, "", "time.Now"))
}

func TestSingleReadCensusRefusesAHelperCallToAReadingEntry(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run() { _ = currentTime() }\n\nfunc helper() {\n\tRun()\n}\n",
		singleReadReport("reader.go", 6, "helper", readThroughEntry, "Run", "time.Now"))
}

// TestSingleReadCensusRefusesAFunctionValue proves a reference that is not a call reads.
func TestSingleReadCensusRefusesAFunctionValue(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func clocks() { _ = []any{currentTime} }\n",
		singleReadReport("reader.go", 3, "clocks", readBelowEntry, "currentTime", ""))
}

func TestSingleReadCensusRefusesAPackageLevelRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "var defaultHome = Home()\n",
		singleReadReport("reader.go", 3, "defaultHome", readBelowEntry, "Home", ""))
}

func TestSingleReadCensusDerivesTheReadSetFromEffects(t *testing.T) {
	t.Parallel()
	effects := syntheticEffectsFile + "func shellPath() string { return os.Getenv(\"SHELL\") }\n"
	wantReadCensus(t, effects, "func launch() { _ = shellPath() }\n",
		singleReadReport("reader.go", 3, "launch", readBelowEntry, "shellPath", ""))
}

func TestSingleReadCensusRefusesAQualifiedRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "import \"github.com/gibbonmi/bench/internal/benchhome\"\n\nfunc poolHome() string { return benchhome.Dir() }\n",
		singleReadReport("reader.go", 5, "poolHome", readBelowEntry, "benchhome.Dir", ""))
}

// TestSingleReadCensusRefusesAKitWrapperCall proves an exported gate function that
// calls the kit read is a read, and a gate form that takes the kit value is not.
func TestSingleReadCensusRefusesAKitWrapperCall(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "import \"github.com/gibbonmi/bench/internal/gate\"\n\nfunc kitDir(root, kit string) string {\n\t_ = gate.LaneForCommitAtKit(root, kit)\n\treturn gate.KitDir()\n}\n",
		singleReadReport("reader.go", 7, "kitDir", readBelowEntry, "gate.KitDir", ""))
}

func TestSingleReadCensusRefusesAnIndirectKitRead(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "import \"github.com/gibbonmi/bench/internal/gate\"\n\nfunc kitRoot() string { return gate.KitRoot() }\n",
		singleReadReport("reader.go", 5, "kitRoot", readBelowEntry, "gate.KitRoot", ""))
}

// TestSingleReadCensusCountsEachKindOfAConstructor proves newAmbient reads time.Now too.
func TestSingleReadCensusCountsEachKindOfAConstructor(t *testing.T) {
	t.Parallel()
	wantReadCensus(t, syntheticEffectsFile, "func Run(home string) {\n\ta := newAmbient(home)\n\t_, _ = a, currentTime()\n}\n",
		singleReadReport("reader.go", 5, "Run", readSecondTime, "", "time.Now"))
}

// TestSingleReadCensusOnTheLiveTree proves each read in the package sits at an entry.
func TestSingleReadCensusOnTheLiveTree(t *testing.T) {
	t.Parallel()
	reports, err := singleReadCensus(".", filepath.Join("..", "gate"))
	if err != nil {
		t.Fatalf("single-read census the package: %v", err)
	}
	if len(reports) != 0 {
		t.Fatalf("the single-read census reports %d reads:\n%s", len(reports), strings.Join(reports, "\n"))
	}
}
