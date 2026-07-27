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
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// checkFunc is the uniform shape every registered check is bound through. Only the
// package-core check reads the tier; the rest ignore the arguments they do not need.
type checkFunc func(root, kitRoot string, tier registry.Tier) []string

// conformanceChecks binds each registry.Checks row to the function that runs it.
// The registry owns the names, the tiers, and the order; this map owns only the
// binding, and TestRegistryBindsEveryCheck asserts the two halves match in both
// directions so tier metadata and executable checks cannot drift apart.
var conformanceChecks = map[string]checkFunc{
	"conformance-canary-families":   func(_, kitRoot string, _ registry.Tier) []string { return checkConformanceCanaryFamilies(kitRoot) },
	"kit-compliance":                func(_, kitRoot string, _ registry.Tier) []string { return checkKitCompliance(kitRoot) },
	"canary-inner-compliance":       func(root, _ string, _ registry.Tier) []string { return checkCanaryInnerCompliance(root) },
	"load-validity-metadata":        func(root, _ string, _ registry.Tier) []string { return checkLoadValidityMetadata(root) },
	"skills-index-command-adapters": func(root, _ string, _ registry.Tier) []string { return checkSkillsIndexAndCommandAdapters(root) },
	"docs-currency-workflow": func(root, kitRoot string, _ registry.Tier) []string {
		return checkDocsCurrencyAndWorkflow(root, kitRoot)
	},
	"line-routing":                 func(root, _ string, _ registry.Tier) []string { return checkLineRouting(root) },
	"package-core-guard":           func(root, _ string, tier registry.Tier) []string { return checkPackageCoreAndGuards(root, tier) },
	"release-evidence-probe":       func(root, _ string, _ registry.Tier) []string { return checkReleaseEvidenceProbe(root) },
	"bench-sh-routes":              func(root, _ string, _ registry.Tier) []string { return checkBenchShRoutes(root) },
	"default-branch-single-source": func(root, _ string, _ registry.Tier) []string { return checkDefaultBranchSingleSource(root) },
	"data-handling-derivation":     func(root, _ string, _ registry.Tier) []string { return checkDataHandlingDerivation(root) },
	"single-control-escaper":       func(root, _ string, _ registry.Tier) []string { return checkSingleControlEscaper(root) },
	"bounds-policy":                func(root, _ string, _ registry.Tier) []string { return checkBoundsPolicy(root) },
	"marker-wait-deadlines":        func(root, _ string, _ registry.Tier) []string { return checkMarkerWaitDeadlines(root) },
	"subcommand-routing":           func(root, _ string, _ registry.Tier) []string { return checkSubcommandRouting(root) },
	"skip-ownership":               func(root, _ string, _ registry.Tier) []string { return checkSkipOwnership(root) },
}

// RunConformance grades root against the checks tier runs, timing each one. An empty
// scope runs the whole tier; otherwise it names the single check to run. Callers name
// a tier and learn nothing about which check belongs to which.
//
// A scope the tier will not execute is a diagnostic and runs nothing at all, because
// the alternative — falling back to the full tier or to zero checks in silence — hands
// a stale binding a green verdict. All three postures live here so no entry point has
// to restate them.
func RunConformance(root, kitRoot string, tier registry.Tier, scope string) []string {
	// The writer clears the root's timing file, so it is established before the scope
	// postures return: a run that executes nothing still has to leave the file empty,
	// or a reader attributes the previous run's lines to this one.
	timing := registry.NewTimingWriter(root)
	if scope != "" {
		check, found := registry.Find(scope)
		if !found {
			return []string{fmt.Sprintf("conformance scope %q names no registered check", scope)}
		}
		if !check.RunsAt(tier) {
			return []string{fmt.Sprintf("conformance scope %q names a check the %s tier does not run", scope, tier)}
		}
	}
	var diags []string
	for _, check := range registry.Checks {
		if !check.RunsAt(tier) {
			continue
		}
		if scope != "" && check.Name != scope {
			continue
		}
		run, bound := conformanceChecks[check.Name]
		if !bound {
			diags = append(diags, "conformance check "+check.Name+" is registered with no bound function")
			continue
		}
		start := time.Now()
		diags = append(diags, run(root, kitRoot, tier)...)
		timing.Record(check.Name, time.Since(start))
	}
	return diags
}

// checkCanaryInnerCompliance grades the kit-compliance rules against the fixture tree
// itself, which only the canary's own inner gate asks for.
func checkCanaryInnerCompliance(root string) []string {
	if os.Getenv("BENCH_CANARY_INNER") != "1" || !exists(filepath.Join(root, ".bench-compliance-canary")) {
		return nil
	}
	return checkKitCompliance(root)
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
	for _, family := range registry.Families() {
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
	return append(diags, unboundCanaryFamilies(kitRoot)...)
}

// unboundCanaryFamilies reports the kit's conformance family directories the registry
// table binds to no check. The sweep scopes a fixture's inner gate by its family, and
// falls back to a full run for a family it cannot resolve — correct for an adopting
// repo, whose families this table will never carry, but in the kit it is the whole
// per-fixture cost the scoping removes, paid in silence. Reading the tree is what
// catches it: the family-presence loop above iterates the table, so a family the table
// omits is invisible from that side.
func unboundCanaryFamilies(kitRoot string) []string {
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	entries, err := os.ReadDir(canaryDir)
	if err != nil {
		return nil
	}
	var diags []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() || canary.FixturePhase(name) != "conformance" {
			continue
		}
		// A directory carrying its own EXPECT is a legacy flat fixture, not a family.
		if _, err := os.Stat(filepath.Join(canaryDir, name, "EXPECT")); err == nil {
			continue
		}
		if !isConformanceFamily(name) {
			diags = append(diags, fmt.Sprintf("canary conformance family %q is bound to no conformance check; add it to the registry family table so its fixtures run scoped", name))
		}
	}
	return diags
}

func isConformanceFamily(family string) bool {
	_, bound := registry.FamilyCheck(family)
	return bound
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
	// Under the shipped prompt-on-stdin contract, claude and codex receive the prompt on
	// their stdin, so their stubs echo argv (which still carries the routed model flag) and
	// then their stdin (the prompt). opencode's CLI documents only a positional prompt, so
	// its adapter reads stdin and re-emits it positionally — its stub echoes argv alone.
	bodies := map[string]string{
		"claude":   "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\"\ncat\n",
		"codex":    "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\"\ncat\n",
		"opencode": "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\"\n",
	}
	for _, name := range []string{"claude", "codex", "opencode"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(bodies[name]), 0o755); err != nil {
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
