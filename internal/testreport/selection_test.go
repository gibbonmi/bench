package testreport

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
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
	want := []string{
		"changedfixture/direct", "changedfixture/embedprod", "changedfixture/embedtest",
		"changedfixture/embedxtest", "changedfixture/production", "changedfixture/spaceglob", "changedfixture/testedge", "changedfixture/xtestedge",
	}
	for _, metadata := range []string{"go.mod", "go.sum", "go.work", "go.work.sum"} {
		t.Run(metadata, func(t *testing.T) {
			got, err := selectCurrentPackages(root, packages, []changedPath{{path: metadata}, {path: "direct/direct.go"}})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("selected packages = %v, want %v", got, want)
			}
		})
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
	if err := os.Symlink("direct/direct.go", filepath.Join(root, "unsafe-live-link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing-target", filepath.Join(root, "unsafe-dangling-link")); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(root, "unsafe-fifo"), 0o644); err != nil {
		t.Fatal(err)
	}
	socket, err := net.Listen("unix", filepath.Join(root, "unsafe-socket"))
	hasSocket := err == nil
	if hasSocket {
		t.Cleanup(func() { _ = socket.Close() })
	} else {
		t.Logf("socket case skipped: cannot create Unix socket: %v", err)
	}
	device := filepath.Join(root, "unsafe-device")
	if err := syscall.Mknod(device, syscall.S_IFCHR|0o600, 0); err != nil {
		t.Logf("device case skipped: cannot create character device: %v", err)
		device = ""
	}
	paths := []string{
		"bad\x1bpath", "bad\x07path", "tab\tpath", "newline\npath", "return\rpath",
		"unsafe-live-link", "unsafe-dangling-link", "unsafe-fifo", "missing/package.go",
	}
	if hasSocket {
		paths = append(paths, "unsafe-socket")
	}
	if device != "" {
		paths = append(paths, "unsafe-device")
	}
	for _, path := range paths {
		t.Run(strings.ReplaceAll(path, "\n", " newline "), func(t *testing.T) {
			if _, err := resolveChangedPackages(context.Background(), root, []string{path}); err == nil {
				t.Fatalf("resolveChangedPackages(%q) succeeded", path)
			}
		})
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
	if err := os.Remove(filepath.Join(root, "embedprod", "input.txt")); err != nil {
		t.Fatal(err)
	}
	if _, err := resolveChangedPackages(context.Background(), root, []string{"embedprod/input.txt"}); err == nil {
		t.Fatal("deleted embed selection succeeded")
	}
}

func TestChangedSpaceAndGlobPathsReachResolverAndCommand(t *testing.T) {
	root, _, base := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
	paths := []string{"space [glob]/package file.go", "space [glob]/embed input[*].txt"}
	for _, path := range paths {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte("package spaceglob\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runChangedGit(t, root, "add", ".")
	runChangedGit(t, root, "commit", "-m", "add space and glob inputs")
	tip := strings.TrimSpace(runChangedGit(t, root, "rev-parse", "HEAD"))
	subject, kind, hint := diff.ResolveChangedSubject(root, base, tip)
	if kind != "" {
		t.Fatalf("changed subject = (%q, %q)", kind, hint)
	}
	for _, path := range paths {
		if !slices.Contains(subject.Paths, path) {
			t.Fatalf("changed paths = %v, want %q", subject.Paths, path)
		}
	}

	goDir := t.TempDir()
	writeChangedSubjectGo(t, filepath.Join(goDir, "go"), filepath.Join(t.TempDir(), "list-environment"), filepath.Join(t.TempDir(), "test-environment"), []listedPackage{{
		Dir:        filepath.Join(root, "space [glob]"),
		ImportPath: "changedcommand/spaceglob",
		Match:      []string{currentPackagePattern},
		EmbedFiles: []string{"embed input[*].txt"},
	}}, "changedcommand/spaceglob")
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installTestSelectionFactory(t, runbinary.Factory{TempRoot: t.TempDir(), Build: func(_ context.Context, _, output string) error {
		return os.WriteFile(output, []byte("selected"), 0o755)
	}, Verify: func(string, string) error { return nil }})

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			got, err := resolveChangedPackages(context.Background(), root, []string{path})
			if err != nil {
				t.Fatal(err)
			}
			if want := []string{"changedcommand/spaceglob"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("selected packages = %v, want %v", got, want)
			}
		})
	}
	output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
	if code != 0 || !strings.Contains(output, "changedcommand/spaceglob") {
		t.Fatalf("changed command = (%d, %q), want selected special-path package", code, output)
	}
}

func TestChangedGraphUsesOneClosedSelectionEnvironment(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
	listEnvironment := filepath.Join(t.TempDir(), "list-environment")
	testEnvironment := filepath.Join(t.TempDir(), "test-environment")
	goDir := t.TempDir()
	writeChangedSubjectGo(t, filepath.Join(goDir, "go"), listEnvironment, testEnvironment, []listedPackage{{
		Dir:        filepath.Join(root, "changed"),
		ImportPath: "changedcommand/changed",
		Match:      []string{currentPackagePattern},
	}}, "changedcommand/changed")
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	for name, value := range map[string]string{
		registry.ConformanceRootEnv:      "/ambient/root",
		registry.ConformanceTierEnv:      "ship",
		registry.ConformanceScopeEnv:     "ambient-scope",
		registry.ConformanceChecksEnv:    "ambient-checks",
		registry.ConformanceInheritedEnv: "ambient-inherited",
		capability.LogEnv:                filepath.Join(t.TempDir(), "capability-log"),
	} {
		t.Setenv(name, value)
	}

	builds := 0
	var selected, source string
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, sourceRoot, output string) error {
			builds++
			source, selected = sourceRoot, output
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})

	output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
	if code != 0 {
		t.Fatalf("changed command = %d\n%s", code, output)
	}
	if builds != 1 {
		t.Fatalf("selection builds = %d, want one", builds)
	}
	for _, environment := range []string{listEnvironment, testEnvironment} {
		got := readTestReportFile(t, environment)
		for _, name := range []string{
			registry.ConformanceRootEnv,
			registry.ConformanceTierEnv,
			registry.ConformanceScopeEnv,
			registry.ConformanceChecksEnv,
			registry.ConformanceInheritedEnv,
			capability.LogEnv,
		} {
			if strings.Contains(got, name+"=") {
				t.Errorf("%s retained ambient %s:\n%s", environment, name, got)
			}
		}
		for name, want := range map[string]string{"BENCH_KIT": source, runbinary.Env: selected} {
			if !strings.Contains(got, name+"="+want+"\n") {
				t.Errorf("%s missing %s=%q:\n%s", environment, name, want, got)
			}
		}
	}
}

func TestChangedNonGoSubjectRendersExplicitEmpty(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "README.md", "docs\n", "README.md", "")
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
	root, base, tip := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
	other := filepath.Join(root, "other", "other_test.go")
	if err := os.MkdirAll(filepath.Dir(other), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(other, []byte("package other\n\nimport \"testing\"\n\nfunc TestSelected(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runChangedGit(t, root, "add", "other/other_test.go")
	runChangedGit(t, root, "commit", "-m", "second changed package")
	tip = strings.TrimSpace(runChangedGit(t, root, "rev-parse", "HEAD"))
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
	if !strings.Contains(output, "packages[2]") {
		t.Fatalf("changed run output = %q, want the complete two-package union", output)
	}
}

func TestChangedRenameHalvesSelectIndependently(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "changed/old.go", "package changed\n", "changed/old.go", "")
	newPath := filepath.Join(root, "changed", "new.go")
	if err := os.WriteFile(newPath, []byte("package changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runChangedGit(t, root, "add", "changed/new.go")
	runChangedGit(t, root, "commit", "-m", "add renamed path")
	tip = strings.TrimSpace(runChangedGit(t, root, "rev-parse", "HEAD"))
	subject, kind, hint := diff.ResolveChangedSubject(root, base, tip)
	if kind != "" {
		t.Fatalf("changed subject = (%q, %q)", kind, hint)
	}
	for _, path := range []string{"changed/old.go", "changed/new.go"} {
		t.Run(path, func(t *testing.T) {
			if !slices.Contains(subject.Paths, path) {
				t.Fatalf("rename paths = %v, want %q", subject.Paths, path)
			}
			got, err := resolveChangedPackages(context.Background(), root, []string{path})
			if err != nil {
				t.Fatal(err)
			}
			if want := []string{"changedcommand/changed"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("selected packages = %v, want %v", got, want)
			}
		})
	}
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
	if code != 0 || !strings.Contains(output, "changedcommand/changed") {
		t.Fatalf("changed rename = (%d, %q), want selected surviving package", code, output)
	}
}

func TestChangedRunsWriteNoGateOwnedRecords(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
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

func TestOrdinaryFocusedModesScrubConformanceEnvironment(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "environment")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	for name, value := range map[string]string{
		registry.ConformanceRootEnv:      "/ambient/root",
		registry.ConformanceTierEnv:      "ship",
		registry.ConformanceScopeEnv:     "ambient-scope",
		registry.ConformanceChecksEnv:    "ambient-checks",
		registry.ConformanceInheritedEnv: "ambient-inherited",
		capability.LogEnv:                filepath.Join(t.TempDir(), "capability-log"),
		"BENCH_KIT":                      t.TempDir(),
	} {
		t.Setenv(name, value)
	}
	var selected, source string
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, sourceRoot, output string) error {
			source = sourceRoot
			selected = output
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})

	for _, tc := range []struct {
		name string
		run  func() (string, int)
	}{
		{name: "default", run: func() (string, int) { return Command(focusedTestModule(t), nil) }},
		{name: "package", run: func() (string, int) { return Command(focusedTestModule(t), []string{"--package", "chosen"}) }},
		{name: "changed", run: func() (string, int) {
			root, base, tip := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
			changedGoDir := t.TempDir()
			writeChangedSubjectGo(t, filepath.Join(changedGoDir, "go"), filepath.Join(t.TempDir(), "list-environment"), marker, []listedPackage{{
				Dir:        filepath.Join(root, "changed"),
				ImportPath: "changedcommand/changed",
				Match:      []string{currentPackagePattern},
			}}, "changedcommand/changed")
			t.Setenv("PATH", changedGoDir+string(os.PathListSeparator)+os.Getenv("PATH"))
			return Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Remove(marker); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			if output, code := tc.run(); code != 0 {
				t.Fatalf("ordinary focused run = %d\n%s", code, output)
			}
			environment := readTestReportFile(t, marker)
			for _, name := range []string{
				registry.ConformanceRootEnv,
				registry.ConformanceTierEnv,
				registry.ConformanceScopeEnv,
				registry.ConformanceChecksEnv,
				registry.ConformanceInheritedEnv,
				capability.LogEnv,
			} {
				if strings.Contains(environment, name+"=") {
					t.Errorf("ordinary child retained %s:\n%s", name, environment)
				}
			}
			for name, want := range map[string]string{"BENCH_KIT": source, runbinary.Env: selected} {
				if !strings.Contains(environment, name+"="+want+"\n") {
					t.Errorf("ordinary child missing %s=%q:\n%s", name, want, environment)
				}
			}
		})
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

func changedCommandRepository(t *testing.T, basePath, baseSource, changedPath, changedSource string) (string, string, string) {
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
	if basePath != "" {
		write(basePath, baseSource)
	}
	runChangedGit(t, root, "init")
	runChangedGit(t, root, "config", "user.email", "test@example.com")
	runChangedGit(t, root, "config", "user.name", "Test")
	runChangedGit(t, root, "add", ".")
	runChangedGit(t, root, "commit", "-m", "base")
	base := strings.TrimSpace(runChangedGit(t, root, "rev-parse", "HEAD"))
	if changedSource == "" {
		if err := os.Remove(filepath.Join(root, changedPath)); err != nil {
			t.Fatal(err)
		}
	} else {
		write(changedPath, changedSource)
	}
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
	spaceGlob := packageAt("space [glob]", "changedfixture/spaceglob")
	spaceGlob.EmbedFiles = []string{"embed input[*].txt"}
	return []listedPackage{direct, production, testEdge, xTestEdge, embedProd, embedTest, embedXTest, spaceGlob}
}

func changedSelectionModule(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"go.mod":                          "module changedfixture\n\ngo 1.25\n",
		"direct/direct.go":                "package direct\n",
		"production/production.go":        "package production\n",
		"testedge/testedge.go":            "package testedge\n",
		"testedge/testedge_test.go":       "package testedge\n",
		"xtestedge/xtestedge.go":          "package xtestedge\n",
		"xtestedge/xtestedge_test.go":     "package xtestedge_test\n",
		"embedprod/embed.go":              "package embedprod\n",
		"embedprod/input.txt":             "production\n",
		"embedtest/embed.go":              "package embedtest\n",
		"embedtest/embed_test.go":         "package embedtest\n",
		"embedtest/input.txt":             "test\n",
		"embedxtest/embed.go":             "package embedxtest\n",
		"embedxtest/embed_test.go":        "package embedxtest_test\n",
		"embedxtest/input.txt":            "external test\n",
		"space [glob]/package file.go":    "package spaceglob\n",
		"space [glob]/embed input[*].txt": "embedded\n",
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

func writeChangedSubjectGo(t *testing.T, path, listEnvironment, testEnvironment string, packages []listedPackage, testPackage string) {
	t.Helper()
	listOutput, err := json.Marshal(packages[0])
	if err != nil {
		t.Fatal(err)
	}
	testOutput, err := json.Marshal(map[string]string{"Action": "pass", "Package": testPackage})
	if err != nil {
		t.Fatal(err)
	}
	source := "#!/usr/bin/env bash\ncase \"$1\" in\nlist)\nenv > " + sanitize.ShellQuote(listEnvironment) + "\nprintf '%s\\n' " + sanitize.ShellQuote(string(listOutput)) + "\n;;\ntest)\nenv > " + sanitize.ShellQuote(testEnvironment) + "\nprintf '%s\\n' " + sanitize.ShellQuote(string(testOutput)) + "\n;;\nesac\n"
	if err := os.WriteFile(path, []byte(source), 0o755); err != nil {
		t.Fatal(err)
	}
}
