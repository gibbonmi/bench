package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

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
// The module record gives each pass the file Go version that gates a suggestion.
const stdlibReimplementationLoadMode = packages.LoadSyntax | packages.NeedModule

// checkStdlibReimplementation reports each production site under root that one of the
// analyzers flags. A load error or a package error is a diagnostic, because a package that
// does not type-check gives no analysis and would otherwise read as a clean pass. The load
// leaves out test packages, so no test file reaches the analyzers.
func checkStdlibReimplementation(root string) []string {
	cfg := &packages.Config{Mode: stdlibReimplementationLoadMode, Dir: root, Tests: false}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return []string{"stdlib-reimplementation cannot load the module packages: " + err.Error()}
	}
	var diags []string
	for _, pkg := range pkgs {
		for _, pkgErr := range pkg.Errors {
			diags = append(diags, fmt.Sprintf("stdlib-reimplementation cannot type-check %s: %s", pkg.PkgPath, pkgErr))
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
			rel := slashRel(root, pos.Filename)
			if rel == ".." || strings.HasPrefix(rel, "../") || filepath.IsAbs(rel) {
				continue
			}
			diags = append(diags, fmt.Sprintf("%s:%d:%d: %s: %s", rel, pos.Line, pos.Column, action.Analyzer.Name, diag.Message))
		}
	}
	return uniqueSorted(diags)
}

// TestStdlibReimplementationBites builds a module with one production site for each kept
// category, one kept-category site in a test file, and one site of an excluded category. The
// check must report each production site and nothing else.
func TestStdlibReimplementationBites(t *testing.T) {
	kitRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	goModBytes, err := os.ReadFile(filepath.Join(kitRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goDirective := regexp.MustCompile(`(?m)^go [0-9.]+$`).FindString(string(goModBytes))
	if goDirective == "" {
		t.Fatal("the kit go.mod has no go directive")
	}
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\n"+goDirective+"\n")
	files := map[string]string{
		"contains.go":     "package fixture\n\nfunc Has(values []string, want string) bool {\n\tfor _, value := range values {\n\t\tif value == want {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}\n",
		"copy.go":         "package fixture\n\nfunc Merge(dst, src map[string]int) {\n\tfor key, value := range src {\n\t\tdst[key] = value\n\t}\n}\n",
		"minimum.go":      "package fixture\n\nfunc min(a, b int) int {\n\tif a < b {\n\t\treturn a\n\t}\n\treturn b\n}\n",
		"cut.go":          "package fixture\n\nimport \"strings\"\n\nfunc Key(line string) string {\n\ti := strings.Index(line, \"=\")\n\tif i >= 0 {\n\t\treturn line[:i]\n\t}\n\treturn line\n}\n",
		"prefix.go":       "package fixture\n\nimport \"strings\"\n\nfunc Strip(line string) string {\n\tif strings.HasPrefix(line, \"x-\") {\n\t\treturn strings.TrimPrefix(line, \"x-\")\n\t}\n\treturn line\n}\n",
		"count.go":        "package fixture\n\nfunc Count(n int) int {\n\ttotal := 0\n\tfor i := 0; i < n; i++ {\n\t\ttotal += i\n\t}\n\treturn total\n}\n",
		"fixture_test.go": "package fixture\n\nfunc hasInTest(values []string, want string) bool {\n\tfor _, value := range values {\n\t\tif value == want {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}\n",
	}
	for name, body := range files {
		writeFixtureFile(t, filepath.Join(root, name), body)
	}
	diags := checkStdlibReimplementation(root)
	for _, want := range []struct{ file, analyzer string }{
		{"contains.go", modernize.SlicesContainsAnalyzer.Name},
		{"copy.go", modernize.MapsLoopAnalyzer.Name},
		{"minimum.go", modernize.MinMaxAnalyzer.Name},
		{"cut.go", modernize.StringsCutAnalyzer.Name},
		{"prefix.go", modernize.StringsCutPrefixAnalyzer.Name},
	} {
		found := false
		for _, diag := range diags {
			found = found || (strings.HasPrefix(diag, want.file+":") && strings.Contains(diag, ": "+want.analyzer+": "))
		}
		if !found {
			t.Errorf("%s: want a %s diagnostic on that file, got %v", want.file, want.analyzer, diags)
		}
	}
	for _, unwanted := range []string{"fixture_test.go:", "count.go:"} {
		if containsDiagnostic(diags, unwanted) {
			t.Errorf("want no diagnostic on %s, got %v", strings.TrimSuffix(unwanted, ":"), diags)
		}
	}
	if len(diags) != 5 {
		t.Errorf("want exactly the five production sites, with no test-file or excluded-category site, got %v", diags)
	}
}
