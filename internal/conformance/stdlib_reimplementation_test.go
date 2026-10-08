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

// excludedSite is one production file that only an excluded modernize analyzer reports. A
// dormant site cannot fire under the check's load at the kit's Go version, and dormant says why.
type excludedSite struct {
	analyzer      *analysis.Analyzer
	body, dormant string
}

// stdlibExcludedSites gives one site for each modernize analyzer that the check leaves out:
// each modernize.Suite member outside the run list, and each exported analyzer outside the suite.
func stdlibExcludedSites(t *testing.T) []excludedSite {
	suiteMember := func(name string) *analysis.Analyzer {
		index := slices.IndexFunc(modernize.Suite, func(analyzer *analysis.Analyzer) bool { return analyzer.Name == name })
		if index < 0 {
			t.Fatalf("modernize.Suite has no analyzer %s", name)
		}
		return modernize.Suite[index]
	}
	return []excludedSite{
		{analyzer: modernize.AnyAnalyzer, body: "package idiom\n\nfunc Show(v interface{}) { _ = v }\n"},
		{analyzer: modernize.AppendClippedAnalyzer, body: "package idiom\n\nfunc Clone(s []int) []int { return append([]int{}, s...) }\n"},
		{analyzer: modernize.AtomicTypesAnalyzer, body: "package idiom\n\nimport \"sync/atomic\"\n\nfunc Hits() int32 {\n\tvar n int32\n\tatomic.AddInt32(&n, 1)\n\treturn atomic.LoadInt32(&n)\n}\n"},
		{analyzer: modernize.BLoopAnalyzer, body: "package idiom\n\nimport \"testing\"\n\nfunc BenchmarkLoop(b *testing.B) {\n\tfor range b.N {\n\t}\n}\n"},
		{analyzer: modernize.EmbedLitAnalyzer, dormant: "embedlit needs Go 1.27", body: "package idiom\n\ntype Inner struct{ V int }\n\ntype Outer struct{ Inner }\n\nvar Literal = Outer{Inner: Inner{V: 1}}\n"},
		{analyzer: modernize.ErrorsAsTypeAnalyzer, dormant: "errorsastype needs Go 1.26", body: "package idiom\n\nimport (\n\t\"errors\"\n\t\"os\"\n)\n\nfunc PathOf(err error) string {\n\tvar pathErr *os.PathError\n\tif errors.As(err, &pathErr) {\n\t\treturn pathErr.Path\n\t}\n\treturn \"\"\n}\n"},
		{analyzer: modernize.FmtAppendfAnalyzer, body: "package idiom\n\nimport \"fmt\"\n\nfunc Line(n int) []byte { return []byte(fmt.Sprintf(\"%d\", n)) }\n"},
		{analyzer: modernize.ForVarAnalyzer, body: "package idiom\n\nfunc Each(values []int, f func(*int)) {\n\tfor _, v := range values {\n\t\tv := v\n\t\tf(&v)\n\t}\n}\n"},
		{analyzer: suiteMember("importcomment"), body: "package idiom // import \"example.com/fixture/internal/idiom\"\n"},
		{analyzer: modernize.NewExprAnalyzer, dormant: "newexpr needs Go 1.26", body: "package idiom\n\nfunc IntPtr(i int) *int { return &i }\n"},
		{analyzer: modernize.OmitZeroAnalyzer, body: "package idiom\n\ntype Config struct {\n\tInner struct{ A int } `json:\",omitempty\"`\n}\n"},
		{analyzer: modernize.PlusBuildAnalyzer, body: "//go:build !fixtureoff\n// +build !fixtureoff\n\npackage idiom\n"},
		{analyzer: modernize.RangeIntAnalyzer, body: "package idiom\n\nfunc Count(n int) int {\n\ttotal := 0\n\tfor i := 0; i < n; i++ {\n\t\ttotal += i\n\t}\n\treturn total\n}\n"},
		{analyzer: suiteMember("reflecttypeassert"), body: "package idiom\n\nimport \"reflect\"\n\nfunc Text(v reflect.Value) (string, bool) {\n\ts, ok := v.Interface().(string)\n\treturn s, ok\n}\n"},
		{analyzer: modernize.ReflectTypeForAnalyzer, body: "package idiom\n\nimport \"reflect\"\n\nvar IntType = reflect.TypeOf(0)\n"},
		{analyzer: suiteMember("slicesbackward"), body: "package idiom\n\nfunc Reverse(s []int) []int {\n\tvar out []int\n\tfor i := len(s) - 1; i >= 0; i-- {\n\t\tout = append(out, s[i])\n\t}\n\treturn out\n}\n"},
		{analyzer: suiteMember("slicesclip"), body: "package idiom\n\nfunc Clip(s []int) []int { return s[:len(s):len(s)] }\n"},
		{analyzer: modernize.SlicesDeleteAnalyzer, body: "package idiom\n\nfunc Drop(s []int, i int) []int { return append(s[:i], s[i+1:]...) }\n"},
		{analyzer: modernize.SlicesSortAnalyzer, body: "package idiom\n\nimport \"sort\"\n\nfunc Sort(s []int) { sort.Slice(s, func(i, j int) bool { return s[i] < s[j] }) }\n"},
		{analyzer: modernize.StdIteratorsAnalyzer, body: "package idiom\n\nimport \"go/types\"\n\nfunc Names(tuple *types.Tuple) []string {\n\tvar out []string\n\tfor i := 0; i < tuple.Len(); i++ {\n\t\tout = append(out, tuple.At(i).Name())\n\t}\n\treturn out\n}\n"},
		{analyzer: modernize.StringsBuilderAnalyzer, body: "package idiom\n\nfunc Join(parts []string) string {\n\tout := \"\"\n\tfor _, part := range parts {\n\t\tout += part\n\t}\n\treturn out\n}\n"},
		{analyzer: modernize.StringsSeqAnalyzer, body: "package idiom\n\nimport \"strings\"\n\nfunc Width(s string) int {\n\tn := 0\n\tfor _, field := range strings.Split(s, \",\") {\n\t\tn += len(field)\n\t}\n\treturn n\n}\n"},
		{analyzer: modernize.TestingContextAnalyzer, body: "package idiom\n\nimport (\n\t\"context\"\n\t\"testing\"\n)\n\nfunc TestContext(t *testing.T) {\n\tctx, cancel := context.WithCancel(context.Background())\n\tdefer cancel()\n\t_ = ctx\n}\n"},
		{analyzer: suiteMember("unsafefuncs"), body: "package idiom\n\nimport \"unsafe\"\n\nfunc Next(p unsafe.Pointer) unsafe.Pointer { return unsafe.Pointer(uintptr(p) + 1) }\n"},
		{analyzer: modernize.WaitGroupGoAnalyzer, body: "package idiom\n\nimport \"sync\"\n\nfunc Run(f func()) {\n\tvar wg sync.WaitGroup\n\twg.Add(1)\n\tgo func() {\n\t\tdefer wg.Done()\n\t\tf()\n\t}()\n\twg.Wait()\n}\n"},
	}
}

// excludedSiteFile gives the slash-relative fixture file that holds the site of the analyzer.
func excludedSiteFile(analyzer *analysis.Analyzer) string {
	return "internal/idiom/" + analyzer.Name + ".go"
}

// analyzerFindings runs the analyzers over the module at root under the check's load mode. It
// gives the slash-relative files that each analyzer reports, keyed by analyzer name.
func analyzerFindings(t *testing.T, root string, analyzers []*analysis.Analyzer) map[string][]string {
	t.Helper()
	pkgs, err := packages.Load(&packages.Config{Mode: stdlibReimplementationLoadMode, Dir: root}, moduleSourcePatterns(root)...)
	if err != nil {
		t.Fatal(err)
	}
	graph, err := checker.Analyze(analyzers, pkgs, nil)
	if err != nil {
		t.Fatal(err)
	}
	findings := map[string][]string{}
	for _, action := range graph.Roots {
		if action.Err != nil {
			t.Fatalf("analyzer %s failed on %s: %v", action.Analyzer.Name, action.Package.PkgPath, action.Err)
		}
		for _, diag := range action.Diagnostics {
			file := slashRel(root, action.Package.Fset.Position(diag.Pos).Filename)
			findings[action.Analyzer.Name] = append(findings[action.Analyzer.Name], file)
		}
	}
	return findings
}

// TestStdlibReimplementationBites builds a module with one production site for each kept
// function, spread over the root package and packages under cmd and internal. It adds a
// kept-category site in a test file, a site for each excluded category, a site in a canary
// payload package, and a site that a line directive moves outside the root. The check must
// report each production site by its slash-relative path, report the moved site as outside
// the root, and report nothing else. Each live excluded site must be a finding of its own
// analyzer, so its silence is not vacuous, and each modernize.Suite member must be kept or
// excluded with a site.
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
		"fixture_test.go":      "package fixture\n\n" + strings.Replace(containsLoop, "Has", "hasInTest", 1),
		"cmd/tool/minimum.go":  "package tool\n\nfunc min(a, b int) int {\n\tif a < b {\n\t\treturn a\n\t}\n\treturn b\n}\n",
		"internal/x/copy.go":   "package x\n\nfunc Merge(dst, src map[string]int) {\n\tfor key, value := range src {\n\t\tdst[key] = value\n\t}\n}\n",
		"internal/x/cut.go":    "package x\n\nimport \"strings\"\n\nfunc Key(line string) string {\n\ti := strings.Index(line, \"=\")\n\tif i >= 0 {\n\t\treturn line[:i]\n\t}\n\treturn line\n}\n",
		"internal/x/prefix.go": "package x\n\nimport \"strings\"\n\nfunc Strip(line string) string {\n\tif strings.HasPrefix(line, \"x-\") {\n\t\treturn strings.TrimPrefix(line, \"x-\")\n\t}\n\treturn line\n}\n",
		"internal/x/suffix.go": "package x\n\nimport \"strings\"\n\nfunc Base(name string) string {\n\tif strings.HasSuffix(name, \".go\") {\n\t\treturn strings.TrimSuffix(name, \".go\")\n\t}\n\treturn name\n}\n",
		"internal/x/moved.go":  "package x\n\n//line /outside/moved.go:1\n" + containsLoop,
		"tests/canary/payload/files/internal/sample/sample.go": "package sample\n\n" + containsLoop,
	}
	excluded := stdlibExcludedSites(t)
	for _, site := range excluded {
		files[excludedSiteFile(site.analyzer)] = site.body
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
	for _, unwanted := range []string{"fixture_test.go:", "tests/canary/"} {
		if containsDiagnostic(diags, unwanted) {
			t.Errorf("want no diagnostic on %s, got %v", strings.TrimSuffix(unwanted, ":"), diags)
		}
	}
	if len(diags) != len(wants) {
		t.Errorf("want exactly the %d wanted diagnostics, got %v", len(wants), diags)
	}
	var excludedAnalyzers []*analysis.Analyzer
	for _, site := range excluded {
		excludedAnalyzers = append(excludedAnalyzers, site.analyzer)
	}
	findings := analyzerFindings(t, root, excludedAnalyzers)
	for _, site := range excluded {
		file := excludedSiteFile(site.analyzer)
		if containsDiagnostic(diags, file+":") {
			t.Errorf("want no diagnostic on the excluded %s site %s, got %v", site.analyzer.Name, file, diags)
		}
		switch fired := slices.Contains(findings[site.analyzer.Name], file); {
		case site.dormant == "" && !fired:
			t.Errorf("want %s to report its own site %s, so the silence of the check is not vacuous", site.analyzer.Name, file)
		case site.dormant != "" && fired:
			t.Errorf("%s now reports its dormant site %s; remove the dormant reason %q", site.analyzer.Name, file, site.dormant)
		}
	}
	for _, analyzer := range modernize.Suite {
		if !slices.Contains(stdlibReimplementationAnalyzers, analyzer) && !slices.Contains(excludedAnalyzers, analyzer) {
			t.Errorf("modernize analyzer %s has no run-list row and no excluded site", analyzer.Name)
		}
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
