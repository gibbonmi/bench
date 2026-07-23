package contract

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

type Env map[string]*string

type Probe struct {
	t        testing.TB
	ExitCode int
	Stdout   string
	Stderr   string
	TimedOut bool
}

func (f Fixture) Run(name string, args ...string) Probe {
	return f.RunEnv(nil, name, args...)
}

func (f Fixture) RunEnv(env map[string]string, name string, args ...string) Probe {
	f.t.Helper()
	return f.RunEnvSpec(envToSpec(env), name, args...)
}

func (f Fixture) RunEnvSpec(env Env, name string, args ...string) Probe {
	f.t.Helper()
	return runFixtureCommandSpec(f.t, f, f.Root, env, "", name, args...)
}

func envToSpec(env map[string]string) Env {
	out := make(Env, len(env))
	for k, v := range env {
		value := v
		out[k] = &value
	}
	return out
}

func RunAt(t testing.TB, f Fixture, dir string, env map[string]string, name string, args ...string) Probe {
	t.Helper()
	return runFixtureCommand(t, f, dir, env, "", name, args...)
}

func RunAtWithInput(t testing.TB, f Fixture, dir string, env map[string]string, stdin, name string, args ...string) Probe {
	t.Helper()
	return runFixtureCommand(t, f, dir, env, stdin, name, args...)
}

// RunAtWithTimeout drives a command that might hang (an uncapped loop, a stuck lock)
// under a wall-clock bound. On timeout it kills the whole process group — not just the
// direct child — so no descendant it spawned survives the test, and returns a Probe
// with TimedOut set rather than blocking the test run. A command that returns before
// the deadline reports its real exit code exactly like Run/RunEnv.
func RunAtWithTimeout(t testing.TB, f Fixture, dir string, env map[string]string, timeout time.Duration, name string, args ...string) Probe {
	t.Helper()
	return runFixtureCommandTimeout(t, f, dir, envToSpec(env), timeout, name, args...)
}

func runFixtureCommandTimeout(t testing.TB, f Fixture, dir string, env Env, timeout time.Duration, name string, args ...string) Probe {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnv(f.Env, env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case err := <-done:
		terminateProcessGroup(cmd)
		exitCode := 0
		if err != nil {
			exitCode = 1
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				exitCode = exitErr.ExitCode()
			}
		}
		return Probe{t: t, ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String()}
	case <-timer.C:
		terminateProcessGroup(cmd)
		<-done
		return Probe{t: t, ExitCode: -1, Stdout: stdout.String(), Stderr: stderr.String(), TimedOut: true}
	}
}

func runFixtureCommand(t testing.TB, f Fixture, dir string, env map[string]string, stdin, name string, args ...string) Probe {
	t.Helper()
	return runFixtureCommandSpec(t, f, dir, envToSpec(env), stdin, name, args...)
}

func runFixtureCommandSpec(t testing.TB, f Fixture, dir string, env Env, stdin, name string, args ...string) Probe {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnv(f.Env, env)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err == nil {
		err = cmd.Wait()
		terminateProcessGroup(cmd)
	}
	exitCode := 0
	if err != nil {
		exitCode = 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}
	return Probe{t: t, ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String()}
}

func terminateProcessGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
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
	t.Cleanup(func() { removeIsolatedEnv(t, base) })
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

func removeIsolatedEnv(t testing.TB, base string) {
	t.Helper()
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if err = os.RemoveAll(base); err == nil {
			return
		}
		if attempt < 2 {
			time.Sleep(time.Second)
		}
	}
	t.Errorf("remove isolated env %s: %v", base, err)
}

func IsolatedEnv(t testing.TB, root string) map[string]string {
	t.Helper()
	return isolatedEnv(t, root)
}

func ProcessEnv(base, overrides map[string]string) []string {
	env := os.Environ()
	for key, value := range base {
		env = setEnv(env, key, value)
	}
	for key, value := range overrides {
		env = setEnv(env, key, value)
	}
	return env
}

func NewExecFixtureAt(t testing.TB, root string) Fixture {
	t.Helper()
	f := NewFixtureAt(t, root, IsolatedEnv(t, t.TempDir()))
	f.Env["PATH"] = os.Getenv("PATH")
	for _, key := range []string{"GOCACHE", "GOMODCACHE"} {
		if value, err := exec.Command("go", "env", key).Output(); err == nil {
			f.Env[key] = strings.TrimSpace(string(value))
		}
	}
	return f
}
