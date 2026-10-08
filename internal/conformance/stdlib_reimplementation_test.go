package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/analysis/passes/modernize"
	"golang.org/x/tools/go/packages"
)

// stdlibReimplementationAnalyzers are the modernize categories that find a hand-written copy
// of a standard-library function. The other modernize categories are idiom, not duplicated
// knowledge, so the check leaves them out.
var stdlibReimplementationAnalyzers = []*analysis.Analyzer{
	modernize.SlicesContainsAnalyzer,
	modernize.MapsLoopAnalyzer,
	modernize.MinMaxAnalyzer,
	modernize.StringsCutAnalyzer,
	modernize.StringsCutPrefixAnalyzer,
}

// stdlibReimplementationLoadMode type-checks the module's own packages from source. None of
// the analyzers exchanges facts, so dependencies load from export data and not from syntax.
const stdlibReimplementationLoadMode = packages.LoadSyntax

// checkStdlibReimplementation reports each production site that one of the analyzers flags
// in the packages that moduleGoFiles names. A load error or a package error is a diagnostic,
// because a package that does not load gives no analysis and would otherwise read as a clean
// pass. A site whose position is outside root is a diagnostic that names its file. The load
// leaves out test packages, so no test file reaches the analyzers. The load uses the host
// build context, so a file that a build constraint excludes on the host is not graded.
func checkStdlibReimplementation(root string) []string {
	// A root without a go.mod is not a Go module, so it carries no module source to grade.
	if !exists(filepath.Join(root, "go.mod")) {
		return nil
	}
	patterns := moduleSourcePatterns(root)
	if len(patterns) == 0 {
		return nil
	}
	cfg := &packages.Config{Mode: stdlibReimplementationLoadMode, Dir: root, Tests: false}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return []string{"stdlib-reimplementation cannot load the module packages: " + err.Error()}
	}
	var diags []string
	for _, pkg := range pkgs {
		for _, pkgErr := range pkg.Errors {
			diags = append(diags, fmt.Sprintf("stdlib-reimplementation %s %s: %s", packageErrorVerb(pkgErr.Kind), pkg.PkgPath, pkgErr))
		}
	}
	if len(diags) > 0 {
		return uniqueSorted(diags)
	}
	graph, err := checker.Analyze(stdlibReimplementationAnalyzers, pkgs, nil)
	if err != nil {
		return []string{"stdlib-reimplementation cannot run the analyzers: " + err.Error()}
	}
	for _, action := range graph.Roots {
		if action.Err != nil {
			diags = append(diags, fmt.Sprintf("stdlib-reimplementation analyzer %s failed on %s: %s", action.Analyzer.Name, action.Package.PkgPath, action.Err))
		}
		for _, diag := range action.Diagnostics {
			pos := action.Package.Fset.Position(diag.Pos)
			prefix, file := "", slashRel(root, pos.Filename)
			if file == ".." || strings.HasPrefix(file, "../") || filepath.IsAbs(file) {
				prefix, file = "stdlib-reimplementation reports a site outside the root at ", filepath.ToSlash(pos.Filename)
			}
			diags = append(diags, fmt.Sprintf("%s%s:%d:%d: %s: %s", prefix, file, pos.Line, pos.Column, action.Analyzer.Name, diag.Message))
		}
	}
	return uniqueSorted(diags)
}

// moduleSourcePatterns gives the package patterns that cover the files moduleGoFiles names:
// "." for a production file in the root package, and a recursive pattern for each top
// directory below it. The go tool leaves testdata directories out of a recursive pattern.
func moduleSourcePatterns(root string) []string {
	var patterns []string
	for _, path := range moduleGoFiles(root) {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		pattern := "."
		if top, _, nested := strings.Cut(slashRel(root, path), "/"); nested {
			pattern = "./" + top + "/..."
		}
		if !slices.Contains(patterns, pattern) {
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}

func packageErrorVerb(kind packages.ErrorKind) string {
	switch kind {
	case packages.TypeError:
		return "cannot type-check"
	case packages.ParseError:
		return "cannot parse"
	default:
		return "cannot load"
	}
}

// TestStdlibReimplementationBites builds a module with one production site for each kept
// function, spread over the root package and packages under cmd and internal. It adds a
// kept-category site in a test file, a site of an excluded category, a site in a canary
// payload package, and a site that a line directive moves outside the root. The check must
// report each production site by its slash-relative path, report the moved site as outside
// the root, and report nothing else.
func TestStdlibReimplementationBites(t *testing.T) {
	kitRoot, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	goModBytes, err := os.ReadFile(filepath.Join(kitRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	kitMod, err := modfile.Parse("go.mod", goModBytes, nil)
	if err != nil {
		t.Fatal(err)
	}
	if kitMod.Go == nil {
		t.Fatal("the kit go.mod has no go directive")
	}
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo "+kitMod.Go.Version+"\n")
	containsLoop := "func Has(values []string, want string) bool {\n\tfor _, value := range values {\n\t\tif value == want {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}\n"
	files := map[string]string{
		"contains.go":          "package fixture\n\n" + containsLoop,
		"count.go":             "package fixture\n\nfunc Count(n int) int {\n\ttotal := 0\n\tfor i := 0; i < n; i++ {\n\t\ttotal += i\n\t}\n\treturn total\n}\n",
		"fixture_test.go":      "package fixture\n\n" + strings.Replace(containsLoop, "Has", "hasInTest", 1),
		"cmd/tool/minimum.go":  "package tool\n\nfunc min(a, b int) int {\n\tif a < b {\n\t\treturn a\n\t}\n\treturn b\n}\n",
		"internal/x/copy.go":   "package x\n\nfunc Merge(dst, src map[string]int) {\n\tfor key, value := range src {\n\t\tdst[key] = value\n\t}\n}\n",
		"internal/x/cut.go":    "package x\n\nimport \"strings\"\n\nfunc Key(line string) string {\n\ti := strings.Index(line, \"=\")\n\tif i >= 0 {\n\t\treturn line[:i]\n\t}\n\treturn line\n}\n",
		"internal/x/prefix.go": "package x\n\nimport \"strings\"\n\nfunc Strip(line string) string {\n\tif strings.HasPrefix(line, \"x-\") {\n\t\treturn strings.TrimPrefix(line, \"x-\")\n\t}\n\treturn line\n}\n",
		"internal/x/suffix.go": "package x\n\nimport \"strings\"\n\nfunc Base(name string) string {\n\tif strings.HasSuffix(name, \".go\") {\n\t\treturn strings.TrimSuffix(name, \".go\")\n\t}\n\treturn name\n}\n",
		"internal/x/moved.go":  "package x\n\n//line /outside/moved.go:1\n" + containsLoop,
		"tests/canary/payload/files/internal/sample/sample.go": "package sample\n\n" + containsLoop,
	}
	for name, body := range files {
		writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(name)), body)
	}
	diags := checkStdlibReimplementation(root)
	wants := []struct{ prefix, analyzer string }{
		{"contains.go:", modernize.SlicesContainsAnalyzer.Name},
		{"cmd/tool/minimum.go:", modernize.MinMaxAnalyzer.Name},
		{"internal/x/copy.go:", modernize.MapsLoopAnalyzer.Name},
		{"internal/x/cut.go:", modernize.StringsCutAnalyzer.Name},
		{"internal/x/prefix.go:", modernize.StringsCutPrefixAnalyzer.Name},
		{"internal/x/suffix.go:", modernize.StringsCutPrefixAnalyzer.Name},
		{"stdlib-reimplementation reports a site outside the root at /outside/moved.go:", modernize.SlicesContainsAnalyzer.Name},
	}
	for _, want := range wants {
		if !slices.ContainsFunc(diags, func(diag string) bool {
			return strings.HasPrefix(diag, want.prefix) && strings.Contains(diag, ": "+want.analyzer+": ")
		}) {
			t.Errorf("want a diagnostic that starts %q and names %s, got %v", want.prefix, want.analyzer, diags)
		}
	}
	for _, unwanted := range []string{"fixture_test.go:", "count.go:", "tests/canary/"} {
		if containsDiagnostic(diags, unwanted) {
			t.Errorf("want no diagnostic on %s, got %v", strings.TrimSuffix(unwanted, ":"), diags)
		}
	}
	if len(diags) != len(wants) {
		t.Errorf("want exactly the %d wanted diagnostics, got %v", len(wants), diags)
	}
}

// TestStdlibReimplementationGradesOnlyModules gives two roots the same file, which does not
// type-check. The root without a go.mod sits inside a module, so a load would reach that file
// through the parent module. That root is not a Go module, so the check does not load it and
// reports nothing. The root with a go.mod reports the type-check error, not a clean pass.
func TestStdlibReimplementationGradesOnlyModules(t *testing.T) {
	const illTyped = "package broken\n\nvar Count int = \"many\"\n"
	module := t.TempDir()
	writeFixtureFile(t, filepath.Join(module, "go.mod"), "module example.com/broken\n")
	writeFixtureFile(t, filepath.Join(module, "broken.go"), illTyped)
	if diags := checkStdlibReimplementation(module); !containsDiagnostic(diags, "stdlib-reimplementation cannot type-check example.com/broken:") {
		t.Errorf("want a type-check diagnostic for a module that does not type-check, got %v", diags)
	}
	plain := filepath.Join(module, "plain")
	writeFixtureFile(t, filepath.Join(plain, "broken.go"), illTyped)
	if diags := checkStdlibReimplementation(plain); diags != nil {
		t.Errorf("want no diagnostic for a root without a go.mod, got %v", diags)
	}
}
