package worktree

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/testcensus"
)

// ambientReadPattern names the process-context reads only the effect boundary may
// perform. Lower owners receive home, root, time, the child base environment, and
// repository facts explicitly.
var ambientReadPattern = regexp.MustCompile(`\b(os\.Getenv|os\.LookupEnv|os\.Environ|os\.Getwd|time\.Now)\(`)

// TestEffectBoundaryCensus fails when any production file below the boundary reads
// the environment, the current directory, or the clock. The one named boundary
// adapter file, effects.go, owns every ambient read; commands resolve there once
// and pass values down. (Coverage row EI1.)
func TestEffectBoundaryCensus(t *testing.T) {
	t.Parallel()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !isSourceFile(name) {
			continue
		}
		if name == effectsFile {
			continue // the named effect boundary adapter file
		}
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			code := line
			if idx := strings.Index(code, "//"); idx >= 0 {
				code = code[:idx]
			}
			if match := ambientReadPattern.FindString(code); match != "" {
				t.Errorf("%s:%d: ambient read %s below the effect boundary; resolve it in %s and pass the value down", name, i+1, match, effectsFile)
			}
		}
	}
}

// gitPackagePath is the import path of the package that declares git.Run.
var gitPackagePath = reflect.TypeFor[git.Run]().PkgPath()

// runConstructionSites names each run construction in the production files of dir that
// sits outside newAmbient and the exported functions, as file:line: owner. A run
// construction is a git.Run composite literal, such as &git.Run{}, or a new(git.Run) call.
// The owner is the function that holds the construction, or the package scope. A method is
// never an exempt entry, so a run that a method builds is named too.
func runConstructionSites(dir string) ([]string, error) {
	files, fset, names, err := parseGoFiles(dir, isSourceFile)
	if err != nil {
		return nil, err
	}
	var sites []string
	for i, file := range files {
		gitName := ""
		for name, path := range fileImportNames(file) {
			if path == gitPackagePath {
				gitName = name
			}
		}
		if gitName == "" {
			continue
		}
		for _, decl := range file.Decls {
			owner := "package scope"
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Recv == nil && (fn.Name.IsExported() || fn.Name.Name == "newAmbient") {
					continue
				}
				owner = fn.Name.Name
			}
			ast.Inspect(decl, func(node ast.Node) bool {
				if runConstruction(node, gitName) {
					sites = append(sites, fmt.Sprintf("%s:%d: %s", names[i], fset.Position(node.Pos()).Line, owner))
				}
				return true
			})
		}
	}
	return sites, nil
}

// runConstruction reports whether node builds a git.Run, where gitName is the file's name
// for the git package.
func runConstruction(node ast.Node, gitName string) bool {
	isRun := func(expr ast.Expr) bool {
		sel, ok := expr.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Run" {
			return false
		}
		pkg, ok := sel.X.(*ast.Ident)
		return ok && pkg.Name == gitName
	}
	switch n := node.(type) {
	case *ast.CompositeLit:
		return isRun(n.Type)
	case *ast.CallExpr:
		fun, ok := n.Fun.(*ast.Ident)
		return ok && fun.Name == "new" && len(n.Args) == 1 && isRun(n.Args[0])
	}
	return false
}

// plantedRunSource is one production file that builds a run in an exported function and in
// the unexported plantedHelper. Its operands are the git import path and the run form, twice.
const plantedRunSource = "package worktree\n\nimport %q\n\nfunc Entry() *git.Run { return %s }\n\nfunc plantedHelper() *git.Run { return %s }\n"

// TestRunConstructionCensus proves that an unexported production helper builds no run. A
// verb's one run comes from newAmbient or from an exported entry, and each helper takes it
// as a parameter. Each planted case builds one run of its form in plantedHelper and one in
// an exported function, and the census names the helper's file and line alone. (Coverage
// rows GQ38, GQ39.)
func TestRunConstructionCensus(t *testing.T) {
	t.Parallel()
	sites, err := runConstructionSites(".")
	if err != nil {
		t.Fatalf("census the package: %v", err)
	}
	for _, site := range sites {
		t.Errorf("%s builds a run outside newAmbient and the exported functions; take the verb's run as a parameter", site)
	}
	for _, form := range []string{"&git.Run{}", "new(git.Run)"} {
		dir := testcensus.Plant(t, map[string]string{"planted.go": fmt.Sprintf(plantedRunSource, gitPackagePath, form, form)})
		planted, err := runConstructionSites(dir)
		if err != nil {
			t.Fatalf("census the planted %s: %v", form, err)
		}
		if want := "planted.go:7: plantedHelper"; len(planted) != 1 || planted[0] != want {
			t.Errorf("census of a planted %s = %q, want exactly [%q]", form, planted, want)
		}
	}
}

// TestHarnessCensusRefusesChildStartOutsideHarness proves the harness effect
// census still refuses a child-process start in a test file outside the harness
// file set. The census walks the package directory itself, so this test drives
// its two sources — the effect pattern and the harness file set — over a
// synthetic line instead. The subject text is assembled from parts, because the
// census reads this file too and a literal would report this line.
// (Coverage row WF13.)
func TestHarnessCensusRefusesChildStartOutsideHarness(t *testing.T) {
	t.Parallel()
	line := "\tcmd := " + "exec" + "." + "Command(\"git\", \"status\")"
	if journeyHarnessFiles["clean_test.go"] {
		t.Fatal("clean_test.go is not a harness file")
	}
	if match := outsideHarnessEffect.FindString(line); match == "" {
		t.Fatalf("the harness census accepts a child start outside the harness: %q", line)
	}
}

// TestPoolAtExplicitHome proves the pool path derives from the injected home value
// alone, with no environment read. (Coverage row EI1.)
func TestPoolAtExplicitHome(t *testing.T) {
	t.Parallel()
	root := filepath.Join(string(filepath.Separator), "repos", "example")
	a := poolAt(filepath.Join(string(filepath.Separator), "home-a"), root)
	b := poolAt(filepath.Join(string(filepath.Separator), "home-b"), root)
	if !strings.HasPrefix(a, filepath.Join(string(filepath.Separator), "home-a", "worktrees")) {
		t.Fatalf("poolAt ignored the injected home: %s", a)
	}
	if filepath.Base(a) != filepath.Base(b) {
		t.Fatalf("pool key must depend only on the root: %s vs %s", a, b)
	}
}

// TestPoolKeysDirAtExplicitHome proves the pool parent derives from the injected home.
func TestPoolKeysDirAtExplicitHome(t *testing.T) {
	t.Parallel()
	home := filepath.Join(string(filepath.Separator), "elsewhere", ".bench")
	if got, want := poolKeysDirAt(home), filepath.Join(home, "worktrees"); got != want {
		t.Fatalf("poolKeysDirAt = %s, want %s", got, want)
	}
}

// TestLeaseLineExplicitTime proves the lease record carries the injected instant,
// not an ambient clock read.
func TestLeaseLineExplicitTime(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 10, 30, 0, 0, time.FixedZone("plus2", 2*3600))
	line := string(leaseLine(now))
	if !strings.Contains(line, now.UTC().Format(leaseTimeLayout)) {
		t.Fatalf("lease line %q does not carry the injected UTC instant", line)
	}
}
