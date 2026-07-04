package contract

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type FixtureOption func(*fixtureConfig)

type fixtureConfig struct {
	repo      bool
	spacePath bool
}

type Fixture struct {
	t    testing.TB
	Root string
	Env  map[string]string
}

type Env map[string]*string

type Probe struct {
	t        testing.TB
	ExitCode int
	Stdout   string
	Stderr   string
}

func WithNoRepo() FixtureOption {
	return func(cfg *fixtureConfig) { cfg.repo = false }
}

func WithSpacePath() FixtureOption {
	return func(cfg *fixtureConfig) { cfg.spacePath = true }
}

func KitRoot(t testing.TB) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve kit root: runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func SubjectRoot(t testing.TB) string {
	t.Helper()
	if root := os.Getenv("BENCH_CONTRACT_ROOT"); root != "" {
		abs, err := filepath.Abs(root)
		if err != nil {
			t.Fatalf("resolve BENCH_CONTRACT_ROOT: %v", err)
		}
		return abs
	}
	return KitRoot(t)
}

func NewFixture(t testing.TB, opts ...FixtureOption) Fixture {
	t.Helper()
	cfg := fixtureConfig{repo: true}
	for _, opt := range opts {
		opt(&cfg)
	}
	root := t.TempDir()
	if cfg.spacePath {
		root = filepath.Join(root, "space dir")
		if err := os.MkdirAll(root, 0o755); err != nil {
			t.Fatalf("create space-path fixture: %v", err)
		}
	}
	f := Fixture{t: t, Root: root}
	f.Env = isolatedEnv(t, f.Root)
	if cfg.repo {
		f.Git("init", "-q")
	}
	return f
}

func (f Fixture) Run(name string, args ...string) Probe {
	return f.RunEnv(nil, name, args...)
}

func (f Fixture) RunEnv(env map[string]string, name string, args ...string) Probe {
	f.t.Helper()
	typed := make(Env, len(env))
	for k, v := range env {
		value := v
		typed[k] = &value
	}
	return f.RunEnvSpec(typed, name, args...)
}

func (f Fixture) RunEnvSpec(env Env, name string, args ...string) Probe {
	f.t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = f.Root
	cmd.Env = mergeEnv(f.Env, env)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}
	return Probe{t: f.t, ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String()}
}

func (f Fixture) Git(args ...string) Probe {
	f.t.Helper()
	probe := f.Run("git", args...)
	probe.RequireExit(0)
	return probe
}

func (f Fixture) GitAllow(args ...string) Probe {
	f.t.Helper()
	return f.Run("git", args...)
}

func (f Fixture) CommitAll(message string) {
	f.t.Helper()
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "add", "-A")
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", message)
}

func (f Fixture) WriteFile(path, contents string) {
	f.t.Helper()
	full := filepath.Join(f.Root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(contents), 0o644); err != nil {
		f.t.Fatalf("write %s: %v", path, err)
	}
}

func (f Fixture) WriteExecutable(path, contents string) {
	f.t.Helper()
	full := filepath.Join(f.Root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(contents), 0o755); err != nil {
		f.t.Fatalf("write executable %s: %v", path, err)
	}
}

func (f Fixture) Exists(path string) bool {
	_, err := os.Stat(filepath.Join(f.Root, filepath.FromSlash(path)))
	return err == nil
}

func (f Fixture) ReadFile(path string) string {
	f.t.Helper()
	data, err := os.ReadFile(filepath.Join(f.Root, filepath.FromSlash(path)))
	if err != nil {
		f.t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func (f Fixture) Bench(args ...string) Probe {
	f.t.Helper()
	bench := filepath.Join(SubjectRoot(f.t), "bin", "bench.sh")
	return f.Run("bash", append([]string{bench}, args...)...)
}

func (f Fixture) BenchEnv(env map[string]string, args ...string) Probe {
	f.t.Helper()
	bench := filepath.Join(SubjectRoot(f.t), "bin", "bench.sh")
	return f.RunEnv(env, "bash", append([]string{bench}, args...)...)
}

func (f Fixture) BenchEnvSpec(env Env, args ...string) Probe {
	f.t.Helper()
	bench := filepath.Join(SubjectRoot(f.t), "bin", "bench.sh")
	return f.RunEnvSpec(env, "bash", append([]string{bench}, args...)...)
}

func noteContractFailure(t testing.TB, msg string) {
	t.Helper()
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(msg)
		}
	})
}

func skipIfSubjectBenchMissing(t testing.TB) {
	t.Helper()
	skipIfSubjectFileMissing(t, "bin/bench.sh")
}

func skipIfSubjectFileMissing(t testing.TB, rel string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(SubjectRoot(t), filepath.FromSlash(rel))); err != nil {
		if os.IsNotExist(err) {
			t.Skipf("subject root has no %s", rel)
		}
		t.Fatalf("stat subject %s: %v", rel, err)
	}
}

func (p Probe) RequireExit(want int) {
	if p.ExitCode != want {
		p.t.Helper()
		p.t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s", p.ExitCode, want, p.Stdout, p.Stderr)
	}
}

func (p Probe) RequireContains(haystack, needle string) {
	p.t.Helper()
	if !strings.Contains(haystack, needle) {
		p.t.Fatalf("missing %q\nstdout:\n%s\nstderr:\n%s", needle, p.Stdout, p.Stderr)
	}
}

func (p Probe) RequireNotContains(haystack, needle string) {
	p.t.Helper()
	if strings.Contains(haystack, needle) {
		p.t.Fatalf("unexpected %q\nstdout:\n%s\nstderr:\n%s", needle, p.Stdout, p.Stderr)
	}
}

func mergeEnv(base map[string]string, overrides Env) []string {
	env := os.Environ()
	for _, k := range []string{
		"HOME",
		"XDG_CONFIG_HOME",
		"XDG_CACHE_HOME",
		"XDG_STATE_HOME",
		"npm_config_cache",
		"GIT_CONFIG_NOSYSTEM",
		"BENCH_HOME",
	} {
		env = setEnv(env, k, "")
	}
	for k, v := range base {
		env = setEnv(env, k, v)
	}
	for k, v := range overrides {
		if v == nil {
			env = unsetEnv(env, k)
			continue
		}
		env = setEnv(env, k, *v)
	}
	return env
}

func setEnv(env []string, k, v string) []string {
	prefix := k + "="
	for i, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			env[i] = prefix + v
			return env
		}
	}
	return append(env, prefix+v)
}

func unsetEnv(env []string, k string) []string {
	prefix := k + "="
	out := env[:0]
	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			out = append(out, entry)
		}
	}
	return out
}

func isolatedEnv(t testing.TB, root string) map[string]string {
	t.Helper()
	base := filepath.Join(root, ".bench-contract-env")
	home := filepath.Join(base, "home")
	xdgConfig := filepath.Join(base, "xdg-config")
	xdgCache := filepath.Join(base, "xdg-cache")
	xdgState := filepath.Join(base, "xdg-state")
	npmCache := filepath.Join(base, "npm-cache")
	benchHome := filepath.Join(base, "bench-home")
	for _, dir := range []string{home, xdgConfig, xdgCache, xdgState, npmCache, benchHome} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("create isolated env dir %s: %v", dir, err)
		}
	}
	return map[string]string{
		"HOME":                home,
		"XDG_CONFIG_HOME":     xdgConfig,
		"XDG_CACHE_HOME":      xdgCache,
		"XDG_STATE_HOME":      xdgState,
		"npm_config_cache":    npmCache,
		"GIT_CONFIG_NOSYSTEM": "1",
		"BENCH_HOME":          benchHome,
	}
}
