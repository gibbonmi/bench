package runbinary

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
)

const missingGoDiagnostic = "Go is absent from PATH; prepend an executable Go toolchain directory to PATH and retry"

func TestSourceWrapperRecoversGoForPrivateBuild(t *testing.T) {
	fixture := newBootstrapFixture(t)
	result := fixture.run(t, fixture.harnessPath, fixture.goExecutable+"\n")
	if result.code != 0 {
		t.Fatalf("partial environment exit = %d, want 0; output=%q", result.code, result.output)
	}
}

func TestSourceWrapperRecoveryPreservesHarnessPath(t *testing.T) {
	fixture := newBootstrapFixture(t)
	result := fixture.run(t, fixture.harnessPath, fixture.goExecutable+"\n")
	if result.code != 0 {
		t.Fatalf("partial environment exit = %d, want 0; output=%q", result.code, result.output)
	}
	wantPath := filepath.Dir(fixture.goExecutable) + string(os.PathListSeparator) + fixture.harnessPath
	paths := fixture.recordedPaths(t)
	if len(paths) == 0 {
		t.Fatal("recovered Go did not record a child PATH")
	}
	for _, got := range paths {
		if got != wantPath {
			t.Fatalf("recovered PATH = %q, want %q", got, wantPath)
		}
	}
	if _, err := os.Stat(filepath.Join(fixture.toolDir, "harness-tool")); err != nil {
		t.Fatalf("unrelated harness tool after recovery: %v", err)
	}
}

func TestSourceWrapperRejectsMissingAndUnsafeGoDiscovery(t *testing.T) {
	for _, test := range []struct {
		name   string
		output func(*bootstrapFixture) string
	}{
		{name: "missing", output: func(*bootstrapFixture) string { return "" }},
		{name: "relative", output: func(*bootstrapFixture) string { return "relative/go\n" }},
		{name: "nonexistent", output: func(f *bootstrapFixture) string { return filepath.Join(f.root, "missing", "go") + "\n" }},
		{name: "multiline", output: func(f *bootstrapFixture) string { return f.goExecutable + "\n" + f.goExecutable + "\n" }},
		{name: "control", output: func(f *bootstrapFixture) string { return f.goExecutable + "\x1bunsafe\n" }},
		{name: "windows", output: func(f *bootstrapFixture) string { return f.windowsGoExecutable + "\n" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newBootstrapFixture(t)
			result := fixture.run(t, fixture.harnessPath, test.output(fixture))
			if result.code == 0 || !strings.Contains(result.output, missingGoDiagnostic) {
				t.Fatalf("unsafe discovery result = (%d, %q), want actionable Go/PATH refusal", result.code, result.output)
			}
			if strings.Contains(result.output, "export PATH=") {
				t.Fatalf("unsafe discovery emitted a recovered assignment: %q", result.output)
			}
			if paths := fixture.recordedPaths(t); len(paths) != 0 {
				t.Fatalf("unsafe discovery executed recovered Go with PATH values %q", paths)
			}
		})
	}
}

func TestSourceWrapperBoundsCleanLoginDiscovery(t *testing.T) {
	fixture := newBootstrapFixture(t)
	started := time.Now()
	result := fixture.runWith(t, fixture.harnessPath, fixture.goExecutable+"\n", []string{"BOOTSTRAP_DISCOVERY_HANG=1"})
	if elapsed := time.Since(started); elapsed >= 3*time.Second {
		t.Fatalf("clean-login discovery elapsed = %s, want less than 3s", elapsed)
	}
	if result.code == 0 || !strings.Contains(result.output, missingGoDiagnostic) {
		t.Fatalf("bounded discovery result = (%d, %q), want actionable Go/PATH refusal", result.code, result.output)
	}
}

func TestSourceWrapperLeavesHealthyPathBytesUnchanged(t *testing.T) {
	fixture := newBootstrapFixture(t)
	healthyPath := filepath.Dir(fixture.goExecutable) + string(os.PathListSeparator) + fixture.harnessPath
	result := fixture.run(t, healthyPath, fixture.goExecutable+"\n")
	if result.code != 0 {
		t.Fatalf("healthy environment exit = %d, want 0; output=%q", result.code, result.output)
	}
	paths := fixture.recordedPaths(t)
	if len(paths) == 0 {
		t.Fatal("healthy Go did not record a child PATH")
	}
	for _, got := range paths {
		if got != healthyPath {
			t.Fatalf("healthy PATH = %q, want unchanged bytes %q", got, healthyPath)
		}
	}
	if data, err := os.ReadFile(fixture.discoveryLog); err == nil || !os.IsNotExist(err) {
		t.Fatalf("healthy wrapper ran clean-login discovery: data=%q err=%v", data, err)
	}
}

type bootstrapFixture struct {
	root                string
	home                string
	toolDir             string
	harnessPath         string
	goExecutable        string
	windowsGoExecutable string
	pathLog             string
	discoveryLog        string
	realGo              string
	realBash            string
	loginHome           string
}

type bootstrapResult struct {
	code   int
	output string
}

func newBootstrapFixture(t *testing.T) *bootstrapFixture {
	t.Helper()
	realGo, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("fixture requires ambient Go: %v", err)
	}
	realGo, err = filepath.Abs(realGo)
	if err != nil {
		t.Fatal(err)
	}
	realBash, err := exec.LookPath("bash")
	if err != nil {
		t.Fatal(err)
	}
	// The repository root comes from git, not from the compiled source file name. Go
	// strips that name under -trimpath, so a runtime.Caller root resolves to a relative
	// module path and every lstat of it fails.
	gitRoot, err := git.Root()
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	root, err := filepath.EvalSymlinks(gitRoot)
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	toolDir := filepath.Join(t.TempDir(), "harness tools")
	goDir := filepath.Join(t.TempDir(), "Go SDK [linux]", "bin")
	if err := os.MkdirAll(toolDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(goDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pathLog := filepath.Join(t.TempDir(), "go-paths")
	discoveryLog := filepath.Join(t.TempDir(), "discoveries")
	goExecutable := filepath.Join(goDir, "go")
	goScript := "#!/bin/sh\nprintf '%s\\n' \"$PATH\" >> \"$BOOTSTRAP_PATH_LOG\"\nexec \"$BOOTSTRAP_REAL_GO\" \"$@\"\n"
	if err := os.WriteFile(goExecutable, []byte(goScript), 0o755); err != nil {
		t.Fatal(err)
	}
	windowsGoExecutable := filepath.Join(goDir, "go.exe")
	if err := os.WriteFile(windowsGoExecutable, []byte(goScript), 0o755); err != nil {
		t.Fatal(err)
	}
	bashScript := "#!/bin/sh\nif [ \"${1:-}\" = -lc ]; then\n  printf discovered >> \"$BOOTSTRAP_DISCOVERY_LOG\"\n  if [ \"${BOOTSTRAP_DISCOVERY_HANG:-}\" = 1 ]; then while :; do :; done; fi\n  printf '%s' \"$BOOTSTRAP_DISCOVERY_OUTPUT\"\n  exit \"${BOOTSTRAP_DISCOVERY_EXIT:-0}\"\nfi\nexec \"$BOOTSTRAP_REAL_BASH\" \"$@\"\n"
	if err := os.WriteFile(filepath.Join(toolDir, "bash"), []byte(bashScript), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(toolDir, "harness-tool"), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"basename", "dirname", "env", "git", "mkdir", "mktemp", "mv", "node", "readlink", "rm", "sh", "timeout", "tr", "uname"} {
		path, err := exec.LookPath(name)
		if err != nil {
			t.Fatalf("fixture tool %s: %v", name, err)
		}
		if err := os.Symlink(path, filepath.Join(toolDir, name)); err != nil {
			t.Fatal(err)
		}
	}
	harnessPath := toolDir
	fixture := &bootstrapFixture{
		root: root, home: home, toolDir: toolDir, harnessPath: harnessPath,
		goExecutable: goExecutable, windowsGoExecutable: windowsGoExecutable,
		pathLog: pathLog, discoveryLog: discoveryLog, realGo: realGo, realBash: realBash,
		loginHome: os.Getenv("HOME"),
	}
	fixture.installSeedBinary(t)
	return fixture
}

func (f *bootstrapFixture) installSeedBinary(t *testing.T) {
	t.Helper()
	seed := filepath.Join(t.TempDir(), "bench")
	if err := Build(context.Background(), f.root, seed); err != nil {
		t.Fatalf("build seed Bench executable: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(f.root, "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil || pkg.Version == "" {
		t.Fatalf("read package version: %v", err)
	}
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	target := filepath.Join(f.home, "cache", "bin", pkg.Version, runtime.GOOS+"-"+arch, "bench")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(seed)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, bytes, 0o755); err != nil {
		t.Fatal(err)
	}
}

func (f *bootstrapFixture) run(t *testing.T, path, discoveryOutput string) bootstrapResult {
	return f.runWith(t, path, discoveryOutput, nil)
}

func (f *bootstrapFixture) runWith(t *testing.T, path, discoveryOutput string, extra []string) bootstrapResult {
	t.Helper()
	cmd := exec.Command(f.realBash, filepath.Join(f.root, "bin", "bench.sh"), "test", "./internal/modelid")
	cmd.Dir = f.root
	overrides := []string{
		"BENCH_HOME=" + f.home,
		"BENCH_KIT=" + f.root,
		"BENCH_RUN_BINARY",
		"CGO_ENABLED=0",
		"BOOTSTRAP_DISCOVERY_EXIT=0",
		"BOOTSTRAP_DISCOVERY_LOG=" + f.discoveryLog,
		"BOOTSTRAP_DISCOVERY_OUTPUT=" + discoveryOutput,
		"BOOTSTRAP_PATH_LOG=" + f.pathLog,
		"BOOTSTRAP_REAL_BASH=" + f.realBash,
		"BOOTSTRAP_REAL_GO=" + f.realGo,
		"ENVMAN_LOAD=loaded",
		"HOME=" + f.loginHome,
		"PATH=" + path,
	}
	overrides = append(overrides, extra...)
	cmd.Env = bootstrapEnvironment(os.Environ(), overrides)
	output, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		code = 1
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		}
	}
	return bootstrapResult{code: code, output: string(output)}
}

func (f *bootstrapFixture) recordedPaths(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(f.pathLog)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	text := strings.TrimSuffix(string(data), "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func bootstrapEnvironment(base, overrides []string) []string {
	values := make(map[string]string, len(overrides))
	for _, item := range overrides {
		key, _, _ := strings.Cut(item, "=")
		values[key] = item
	}
	environment := make([]string, 0, len(base)+len(overrides))
	for _, item := range base {
		key, _, _ := strings.Cut(item, "=")
		if _, replace := values[key]; !replace {
			environment = append(environment, item)
		}
	}
	for _, item := range overrides {
		if strings.Contains(item, "=") {
			environment = append(environment, item)
		}
	}
	return environment
}
