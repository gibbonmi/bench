package conformance

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess"
)

var conformanceFamilies = []string{
	"load-validity-metadata",
	"skills-index-command-adapters",
	"docs-currency-token-diet",
	"workflow-guidance-anchors",
	"coverage-map-validation",
	"line-routing",
	"package-core-guard",
	"compliance-hardening",
}

func RunConformance(root, kitRoot string) []string {
	var diags []string
	diags = append(diags, checkConformanceCanaryFamilies(kitRoot)...)
	diags = append(diags, checkKitCompliance(kitRoot)...)
	if os.Getenv("BENCH_CANARY_INNER") == "1" && exists(filepath.Join(root, ".bench-compliance-canary")) {
		diags = append(diags, checkKitCompliance(root)...)
	}
	diags = append(diags, checkLoadValidityMetadata(root)...)
	diags = append(diags, checkSkillsIndexAndCommandAdapters(root)...)
	diags = append(diags, checkDocsCurrencyAndWorkflow(root, kitRoot)...)
	diags = append(diags, checkLineRouting(root)...)
	diags = append(diags, checkPackageCoreAndGuards(root)...)
	diags = append(diags, checkBenchShRoutes(root)...)
	return diags
}

// benchShRoutes are the top-level bin/bench.sh case labels that must reach the Go core so
// every shipped surface (kit CLI, linked by-path CLI, hooks) hits one implementation. The
// route-anchor: dropping a route sends a shipped command to a dead key, and this fires.
var benchShRoutes = []string{"commit", "spec", "resume-clean", "worktree-hook"}

// checkBenchShRoutes asserts bin/bench.sh carries a case route for each command in
// benchShRoutes. It bites when a route is removed (the `<name>)` label disappears).
func checkBenchShRoutes(root string) []string {
	bench := readIfExists(filepath.Join(root, "bin", "bench.sh"))
	if bench == "" {
		return nil
	}
	var diags []string
	for _, route := range benchShRoutes {
		if !regexp.MustCompile(`(?m)^  ` + regexp.QuoteMeta(route) + `\)\s`).MatchString(bench) {
			diags = append(diags, fmt.Sprintf("bin/bench.sh has no route for '%s' (a shipped command with no case label reaches a dead key)", route))
		}
	}
	return diags
}

// TestBenchShRouteAnchorBites is the recorded bite proof for checkBenchShRoutes (per
// craft-gate): a bin/bench.sh with a route present passes; removing that route's case
// label makes the anchor fire. It runs against a synthetic script, not the repo tree.
func TestBenchShRouteAnchorBites(t *testing.T) {
	root := t.TempDir()
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(body string) {
		if err := os.WriteFile(filepath.Join(binDir, "bench.sh"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("case \"${1:-help}\" in\n  commit)   route_porcelain \"$@\" ;;\n  spec)     route_porcelain \"$@\" ;;\n  resume-clean) route_porcelain \"$@\" ;;\n  worktree-hook) route_porcelain \"$@\" ;;\nesac\n")
	if diags := checkBenchShRoutes(root); len(diags) != 0 {
		t.Fatalf("both routes present: want no diagnostics, got %v", diags)
	}

	write("case \"${1:-help}\" in\n  spec)     route_porcelain \"$@\" ;;\n  resume-clean) route_porcelain \"$@\" ;;\n  worktree-hook) route_porcelain \"$@\" ;;\nesac\n")
	diags := checkBenchShRoutes(root)
	if len(diags) != 1 || !strings.Contains(diags[0], "no route for 'commit'") {
		t.Fatalf("dropped commit route: want a single commit diagnostic, got %v", diags)
	}
}

func checkConformanceCanaryFamilies(kitRoot string) []string {
	var diags []string
	for _, family := range conformanceFamilies {
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		entries, err := os.ReadDir(familyDir)
		if err != nil {
			diags = append(diags, fmt.Sprintf("canary conformance family %q has no fixture directories under %s", family, filepath.ToSlash(filepath.Join("tests", "canary", family))))
			continue
		}
		count := 0
		for _, entry := range entries {
			if entry.IsDir() {
				count++
			}
		}
		if count == 0 {
			diags = append(diags, fmt.Sprintf("canary conformance family %q has no fixture directories under %s", family, filepath.ToSlash(filepath.Join("tests", "canary", family))))
		}
	}
	return diags
}

func isConformanceFamily(family string) bool {
	for _, candidate := range conformanceFamilies {
		if family == candidate {
			return true
		}
	}
	return false
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
	r := subprocess.Capture(cmd)
	return &Probe{Args: append([]string(nil), args...), ExitCode: r.ExitCode, Stdout: r.Stdout, Stderr: r.Stderr, Err: r.Err}
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
	hasNpmCache := false
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_CONFORMANCE_ROOT=") {
			continue
		}
		if strings.HasPrefix(kv, "NPM_CONFIG_CACHE=") && strings.TrimPrefix(kv, "NPM_CONFIG_CACHE=") != "" {
			hasNpmCache = true
		}
		env = append(env, kv)
	}
	if !hasNpmCache {
		env = append(env, "NPM_CONFIG_CACHE="+filepath.Join(os.TempDir(), "bench-npm-cache"))
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
