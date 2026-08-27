package testreport

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestChangedPackageClosureAcrossAllGoEdges(t *testing.T) {
	root := changedSelectionModule(t)
	got, err := selectCurrentPackages(root, changedSelectionGraph(root), []changedPath{
		{path: "direct/direct.go"},
		{path: "embedprod/input.txt"},
		{path: "embedtest/input.txt"},
		{path: "embedxtest/input.txt"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"changedfixture/direct",
		"changedfixture/embedprod",
		"changedfixture/embedtest",
		"changedfixture/embedxtest",
		"changedfixture/production",
		"changedfixture/testedge",
		"changedfixture/xtestedge",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("selected packages = %v, want %v", got, want)
	}
}

func TestChangedPackageSelectionMetadataAndMixedUnion(t *testing.T) {
	root := changedSelectionModule(t)
	packages := changedSelectionGraph(root)
	got, err := selectCurrentPackages(root, packages, []changedPath{{path: "go.mod"}, {path: "direct/direct.go"}})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"changedfixture/direct", "changedfixture/embedprod", "changedfixture/embedtest",
		"changedfixture/embedxtest", "changedfixture/production", "changedfixture/testedge", "changedfixture/xtestedge",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("selected packages = %v, want %v", got, want)
	}
}

func TestChangedPackageSelectionRefusalMatrix(t *testing.T) {
	root := changedSelectionModule(t)
	previous := listCurrentPackages
	calls := 0
	listCurrentPackages = func(context.Context, string) ([]listedPackage, error) {
		calls++
		return changedSelectionGraph(root), nil
	}
	t.Cleanup(func() { listCurrentPackages = previous })
	unsafe := filepath.Join(root, "unsafe-link")
	if err := os.Symlink("direct/direct.go", unsafe); err != nil {
		t.Fatal(err)
	}
	fifo := filepath.Join(root, "unsafe-fifo")
	if err := syscall.Mkfifo(fifo, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"bad\x1bpath", "unsafe-link", "unsafe-fifo", "missing/package.go"} {
		if _, err := resolveChangedPackages(context.Background(), root, []string{path}); err == nil {
			t.Fatalf("resolveChangedPackages(%q) succeeded", path)
		}
	}
	if calls != 0 {
		t.Fatalf("go list calls = %d, want zero for refused inputs", calls)
	}
	if err := os.Remove(filepath.Join(root, "direct", "direct.go")); err != nil {
		t.Fatal(err)
	}
	got, err := resolveChangedPackages(context.Background(), root, []string{"direct/direct.go"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{"changedfixture/direct", "changedfixture/production", "changedfixture/testedge", "changedfixture/xtestedge"}) {
		t.Fatalf("deleted surviving-package selection = %v", got)
	}
}

func TestChangedNonGoSubjectRendersExplicitEmpty(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "README.md", "docs\n")
	installTestSelectionFactory(t, runbinary.Factory{TempRoot: t.TempDir(), Build: func(_ context.Context, _, output string) error {
		return os.WriteFile(output, []byte("selected"), 0o755)
	}, Verify: func(string, string) error { return nil }})
	output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
	if code != 0 {
		t.Fatalf("emptyReport code = %d\n%s", code, output)
	}
	for _, table := range []string{"packages[0]", "failures[0]", "skips[0]"} {
		if !strings.Contains(output, table) {
			t.Fatalf("empty output = %q, want %q", output, table)
		}
	}
}

func TestChangedPackageRunPattern(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "changed/changed.go", "package changed\n")
	installTestSelectionFactory(t, runbinary.Factory{TempRoot: t.TempDir(), Build: func(_ context.Context, _, output string) error {
		return os.WriteFile(output, []byte("selected"), 0o755)
	}, Verify: func(string, string) error { return nil }})
	marker := filepath.Join(t.TempDir(), "argv")
	t.Setenv("BENCH_TEST_MARKER", marker)
	output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip, "--run", "^TestSelected$"})
	if code != 0 {
		t.Fatalf("changed run = %d\n%s", code, output)
	}
	if got := readTestReportFile(t, marker); !strings.Contains(got, "-test.run=^TestSelected$") {
		t.Fatalf("test argv = %q, want unchanged run pattern", got)
	}
}

func TestChangedRunsWriteNoGateOwnedRecords(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "changed/changed.go", "package changed\n")
	installTestSelectionFactory(t, runbinary.Factory{TempRoot: t.TempDir(), Build: func(_ context.Context, _, output string) error {
		return os.WriteFile(output, []byte("selected"), 0o755)
	}, Verify: func(string, string) error { return nil }})
	record := filepath.Join(root, ".git", "bench-last-gate")
	if err := os.WriteFile(record, []byte("before"), 0o644); err != nil {
		t.Fatal(err)
	}
	if output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip}); code != 0 {
		t.Fatalf("changed run = %d\n%s", code, output)
	}
	if got := readTestReportFile(t, record); got != "before" {
		t.Fatalf("gate record = %q, want unchanged", got)
	}
}

func TestCurrentPackageGraphLoadsAllEmbedClasses(t *testing.T) {
	packages, err := currentPackages(context.Background(), currentPackageGraphModule(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 1 {
		t.Fatalf("current packages = %d, want one ordinary package", len(packages))
	}
	got := packages[0]
	if got.ImportPath != "currentgraph" {
		t.Fatalf("package import path = %q, want currentgraph", got.ImportPath)
	}
	if !reflect.DeepEqual(got.EmbedFiles, []string{"production.txt"}) {
		t.Fatalf("production embed files = %v, want [production.txt]", got.EmbedFiles)
	}
	if !reflect.DeepEqual(got.TestEmbedFiles, []string{"test.txt"}) {
		t.Fatalf("test embed files = %v, want [test.txt]", got.TestEmbedFiles)
	}
	if !reflect.DeepEqual(got.XTestEmbedFiles, []string{"external-test.txt"}) {
		t.Fatalf("external-test embed files = %v, want [external-test.txt]", got.XTestEmbedFiles)
	}
}

func changedCommandRepository(t *testing.T, changedPath, changedSource string) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	write := func(path, source string) {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module changedcommand\n\ngo 1.25\n")
	write("changed/changed.go", "package changed\n\nconst Version = \"base\"\n")
	write("changed/changed_test.go", `package changed

import (
	"os"
	"strings"
	"testing"
)

func TestSelected(t *testing.T) {
	if marker := os.Getenv("BENCH_TEST_MARKER"); marker != "" {
		if err := os.WriteFile(marker, []byte(strings.Join(os.Args, "\n")), 0o644); err != nil { t.Fatal(err) }
	}
}
`)
	runChangedGit(t, root, "init")
	runChangedGit(t, root, "config", "user.email", "test@example.com")
	runChangedGit(t, root, "config", "user.name", "Test")
	runChangedGit(t, root, "add", ".")
	runChangedGit(t, root, "commit", "-m", "base")
	base := strings.TrimSpace(runChangedGit(t, root, "rev-parse", "HEAD"))
	write(changedPath, changedSource)
	runChangedGit(t, root, "add", ".")
	runChangedGit(t, root, "commit", "-m", "changed")
	tip := strings.TrimSpace(runChangedGit(t, root, "rev-parse", "HEAD"))
	return root, base, tip
}

func runChangedGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
	return string(output)
}

func changedSelectionGraph(root string) []listedPackage {
	packageAt := func(dir, importPath string) listedPackage {
		return listedPackage{Dir: filepath.Join(root, dir), ImportPath: importPath}
	}
	direct := packageAt("direct", "changedfixture/direct")
	production := packageAt("production", "changedfixture/production")
	production.Imports = []string{direct.ImportPath}
	testEdge := packageAt("testedge", "changedfixture/testedge")
	testEdge.TestImports = []string{production.ImportPath}
	xTestEdge := packageAt("xtestedge", "changedfixture/xtestedge")
	xTestEdge.XTestImports = []string{testEdge.ImportPath}
	embedProd := packageAt("embedprod", "changedfixture/embedprod")
	embedProd.EmbedFiles = []string{"input.txt"}
	embedTest := packageAt("embedtest", "changedfixture/embedtest")
	embedTest.TestEmbedFiles = []string{"input.txt"}
	embedXTest := packageAt("embedxtest", "changedfixture/embedxtest")
	embedXTest.XTestEmbedFiles = []string{"input.txt"}
	return []listedPackage{direct, production, testEdge, xTestEdge, embedProd, embedTest, embedXTest}
}

func changedSelectionModule(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"go.mod":                      "module changedfixture\n\ngo 1.25\n",
		"direct/direct.go":            "package direct\n",
		"production/production.go":    "package production\n\nimport _ \"changedfixture/direct\"\n",
		"testedge/testedge.go":        "package testedge\n",
		"testedge/testedge_test.go":   "package testedge\n\nimport _ \"changedfixture/production\"\n",
		"xtestedge/xtestedge.go":      "package xtestedge\n",
		"xtestedge/xtestedge_test.go": "package xtestedge_test\n\nimport _ \"changedfixture/testedge\"\n",
		"embedprod/embed.go":          "package embedprod\n\nimport \"embed\"\n\nvar _ embed.FS\n\n//go:embed input.txt\nvar input string\n",
		"embedprod/input.txt":         "production\n",
		"embedtest/embed.go":          "package embedtest\n",
		"embedtest/embed_test.go":     "package embedtest\n\nimport \"embed\"\n\nvar _ embed.FS\n\n//go:embed input.txt\nvar input string\n",
		"embedtest/input.txt":         "test\n",
		"embedxtest/embed.go":         "package embedxtest\n",
		"embedxtest/embed_test.go":    "package embedxtest_test\n\nimport \"embed\"\n\nvar _ embed.FS\n\n//go:embed input.txt\nvar input string\n",
		"embedxtest/input.txt":        "external test\n",
	}
	for path, source := range files {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func currentPackageGraphModule(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"go.mod":                 "module currentgraph\n\ngo 1.25\n",
		"embed.go":               "package currentgraph\n\nimport \"embed\"\n\nvar _ embed.FS\n\n//go:embed production.txt\nvar production string\n",
		"production.txt":         "production\n",
		"embed_test.go":          "package currentgraph\n\nimport \"embed\"\n\nvar _ embed.FS\n\n//go:embed test.txt\nvar testInput string\n",
		"test.txt":               "test\n",
		"embed_external_test.go": "package currentgraph_test\n\nimport \"embed\"\n\nvar _ embed.FS\n\n//go:embed external-test.txt\nvar externalTest string\n",
		"external-test.txt":      "external test\n",
	}
	for path, source := range files {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
