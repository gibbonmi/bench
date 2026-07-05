package conformance

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func RunConformance(root, kitRoot string) []string {
	var diags []string
	diags = append(diags, checkLoadValidityMetadata(root)...)
	diags = append(diags, checkSkillsIndexAndCommandAdapters(root)...)
	diags = append(diags, checkDocsCurrencyAndWorkflow(root, kitRoot)...)
	diags = append(diags, checkLineRouting(root)...)
	diags = append(diags, checkPackageCoreAndGuards(root)...)
	return diags
}

func containsDiagnostic(diags []string, want string) bool {
	for _, diag := range diags {
		if strings.Contains(diag, want) {
			return true
		}
	}
	return false
}

func frontmatterField(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	fence := 0
	prefix := key + ":"
	for _, line := range strings.Split(string(data), "\n") {
		if line == "---" {
			fence++
			continue
		}
		if fence == 1 && strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
		if fence > 1 {
			return ""
		}
	}
	return ""
}

func readIfExists(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func anyContains(values []string, needle string) bool {
	for _, value := range values {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

func slashRel(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}

func uniqueSorted(values []string) []string {
	sort.Strings(values)
	if len(values) == 0 {
		return values
	}
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}

// runProbe captures stdout and stderr separately: probes like the npm pack
// JSON parse read stdout alone, and subprocess stderr chatter (npm's update
// notifier, warnings) must not corrupt it.
func runProbe(cmd *exec.Cmd, args []string) *Probe {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		if cmd.ProcessState != nil {
			exitCode = cmd.ProcessState.ExitCode()
		}
	}
	return &Probe{Args: append([]string(nil), args...), ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
}

func runAt(dir string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	return runProbe(cmd, args)
}

func runAtCleanEnv(dir string, args ...string) *Probe {
	return runAtEnv(dir, conformanceSubprocessEnv(), args...)
}

func runAtEnv(dir string, env []string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = env
	return runProbe(cmd, args)
}

func runWithInput(dir, input string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)
	return runProbe(cmd, args)
}

func runWithInputEnv(dir string, env []string, input string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = strings.NewReader(input)
	return runProbe(cmd, args)
}

func conformanceSubprocessEnv() []string {
	env := make([]string, 0, len(os.Environ()))
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_CONFORMANCE_ROOT=") {
			continue
		}
		env = append(env, kv)
	}
	return env
}

func wrapperStubDir(realBench string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "bench-wrapper-stub-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	content := "#!/usr/bin/env bash\nexec " + shellQuote(realBench) + " \"$@\"\n"
	if err := os.WriteFile(filepath.Join(dir, "bench"), []byte(content), 0o755); err != nil {
		cleanup()
		return "", func() {}, err
	}
	return dir, cleanup, nil
}

func adapterStubDir(realBench string) (string, func(), error) {
	dir, cleanup, err := wrapperStubDir(realBench)
	if err != nil {
		return "", cleanup, err
	}
	for _, name := range []string{"claude", "codex", "opencode"} {
		content := "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\"\n"
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o755); err != nil {
			cleanup()
			return "", func() {}, err
		}
	}
	return dir, cleanup, nil
}

func tempGitRepoWithLines(linesEnv string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "bench-line-repo-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	if probe := runAt(dir, "git", "init", "-q"); probe == nil || probe.ExitCode != 0 {
		cleanup()
		return "", func() {}, fmt.Errorf("git init failed")
	}
	if err := os.MkdirAll(filepath.Join(dir, ".bench"), 0o755); err != nil {
		cleanup()
		return "", func() {}, err
	}
	if err := os.WriteFile(filepath.Join(dir, ".bench", "lines.env"), []byte(linesEnv), 0o644); err != nil {
		cleanup()
		return "", func() {}, err
	}
	return dir, cleanup, nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
