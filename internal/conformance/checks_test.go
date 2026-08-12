package conformance

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// checkBinding is the executable half of a registry row. The map below deliberately repeats
// only the facts an independent mutation oracle needs: the name-to-function binding, tier,
// and subject. Registry order, meta membership, and inputs stay single-sourced in
// registry.Checks. Keeping these three facts independent is what makes the named CM5/CM6/CM7
// mutations red when functions are swapped or the advertised tier or subject drifts while
// the executable binding is unchanged.
type checkBinding struct {
	implementation any
	tier           registry.Tier
	subject        registry.Subject
}

var conformanceChecks map[string]checkBinding

func init() {
	conformanceChecks = map[string]checkBinding{
		"conformance-meta":              {checkConformanceMeta, registry.Dev, registry.SubjectKitRoot},
		"conformance-canary-families":   {checkConformanceCanaryFamilies, registry.Dev, registry.SubjectKitRoot},
		"kit-compliance":                {checkKitCompliance, registry.Dev, registry.SubjectKitRoot},
		"canary-fixture-compliance":     {checkCanaryFixtureCompliance, registry.Dev, registry.SubjectRoot},
		"load-validity-metadata":        {checkLoadValidityMetadata, registry.Dev, registry.SubjectRoot},
		"skills-index-command-adapters": {checkSkillsIndexAndCommandAdapters, registry.Dev, registry.SubjectRoot},
		"docs-currency-workflow":        {checkDocsCurrencyAndWorkflow, registry.Dev, registry.SubjectRootAndKitRoot},
		"gate-entry-contract":           {checkGateEntryContract, registry.Dev, registry.SubjectRoot},
		"ordinary-build-census":         {checkOrdinaryBuildCensus, registry.Dev, registry.SubjectKitRoot},
		"offline-smoke-proof":           {checkOfflineSmokeProof, registry.Dev, registry.SubjectRoot},
		"handoff-shape-single-source":   {checkHandoffShape, registry.Dev, registry.SubjectRoot},
		"harness-prefix-single-source":  {checkHarnessPrefix, registry.Dev, registry.SubjectRoot},
		"package-shipped-surface":       {checkPackageShippedSurface, registry.Dev, registry.SubjectRoot},
		"line-routing":                  {checkLineRouting, registry.Dev, registry.SubjectRoot},
		"package-core-guard":            {checkPackageCoreAndGuards, registry.Dev, registry.SubjectRoot},
		"release-evidence-probe":        {checkReleaseEvidenceProbe, registry.Ship, registry.SubjectRoot},
		"bench-sh-routes":               {checkBenchShRoutes, registry.Dev, registry.SubjectRoot},
		"default-branch-single-source":  {checkDefaultBranchSingleSource, registry.Dev, registry.SubjectRoot},
		"data-handling-derivation":      {checkDataHandlingDerivation, registry.Dev, registry.SubjectRoot},
		"single-control-escaper":        {checkSingleControlEscaper, registry.Dev, registry.SubjectRoot},
		"bounds-policy":                 {checkBoundsPolicy, registry.Dev, registry.SubjectRoot},
		"marker-wait-deadlines":         {checkMarkerWaitDeadlines, registry.Dev, registry.SubjectRoot},
		"subcommand-routing":            {checkSubcommandRouting, registry.Dev, registry.SubjectRoot},
		"skip-ownership":                {checkSkipOwnership, registry.Dev, registry.SubjectRoot},
		"decision-map-integrity":        {maps.ValidateDecisionMapTree, registry.Dev, registry.SubjectRoot},
		"injected-port-registry":        {checkInjectedPortRegistry, registry.Dev, registry.SubjectRoot},
		"guidance-prose-budgets":        {checkGuidanceProseBudgets, registry.Dev, registry.SubjectRoot},
	}
}

func (b checkBinding) identity() string {
	fn := runtime.FuncForPC(reflect.ValueOf(b.implementation).Pointer())
	if fn == nil {
		return ""
	}
	name := fn.Name()
	return name[strings.LastIndex(name, ".")+1:]
}

func (b checkBinding) runsAt(tier registry.Tier) bool {
	return b.tier == registry.Dev || tier == registry.Ship
}

func (b checkBinding) run(root, kitRoot string, tier registry.Tier) []string {
	subject := root
	if b.subject == registry.SubjectKitRoot {
		subject = kitRoot
	}
	switch run := b.implementation.(type) {
	case func(string) []string:
		return run(subject)
	case func(string, string) []string:
		return run(root, kitRoot)
	case func(string, registry.Tier) []string:
		return run(subject, tier)
	default:
		return []string{"conformance check carries an unsupported executable binding"}
	}
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
	return RunConformanceSelection(root, kitRoot, tier, scope, nil, nil)
}

// RunConformanceSelection accepts the inner-canary single-check control and the outer
// gate's exact executed/inherited ordinary-check partition as distinct authorities.
func RunConformanceSelection(root, kitRoot string, tier registry.Tier, scope string, selected, inheritedControl *string) []string {
	// The writer clears the root's timing file, so it is established before the scope
	// postures return: a run that executes nothing still has to leave the file empty,
	// or a reader attributes the previous run's lines to this one.
	timing := registry.NewTimingWriter(root)
	if scope != "" && (selected != nil || inheritedControl != nil) {
		return []string{"conformance selection carries both inner and outer controls"}
	}
	if scope != "" {
		check, found := registry.Find(scope)
		if !found {
			return []string{fmt.Sprintf("conformance scope %q names no registered check", scope)}
		}
		if !check.RunsAt(tier) {
			return []string{fmt.Sprintf("conformance scope %q names a check the %s tier does not run", scope, tier)}
		}
	}
	// Ship is a lifecycle-final rehearsal, so only the canary-owned singular control may
	// narrow it. An outer ordered set is ignored and the complete ship tier runs.
	if tier == registry.Ship {
		selected = nil
		inheritedControl = nil
	}
	selection, executed, inherited, selectionDiags := orderedConformancePartition(tier, selected, inheritedControl)
	var diags []string
	diags = append(diags, selectionDiags...)
	for _, check := range registry.Checks {
		run, bound := conformanceChecks[check.Name]
		if !bound {
			diags = append(diags, "conformance check "+check.Name+" is registered with no bound function")
			continue
		}
		if !run.runsAt(tier) || scope != "" && check.Name != scope || selection != nil && !selection[check.Name] && !check.Meta {
			continue
		}
		start := time.Now()
		if check.Name == "conformance-meta" && selection != nil {
			diags = append(diags, checkConformanceMetaForPartition(kitRoot, tier, executed, inherited)...)
		} else {
			diags = append(diags, run.run(root, kitRoot, tier)...)
		}
		timing.Record(check.Name, time.Since(start))
	}
	return diags
}

func orderedConformancePartition(tier registry.Tier, selected, inherited *string) (map[string]bool, []string, []string, []string) {
	if selected == nil && inherited == nil {
		return nil, registry.Names(tier), nil, nil
	}
	if selected == nil || inherited == nil {
		return nil, registry.Names(tier), nil, []string{"conformance ordered partition must carry both executed and inherited sets"}
	}
	executedSet, _, selectedDiags := orderedConformanceSelection(tier, selected)
	_, inheritedNames, inheritedDiags := orderedConformanceSelection(tier, inherited)
	if len(selectedDiags) > 0 || len(inheritedDiags) > 0 {
		return nil, registry.Names(tier), nil, append(selectedDiags, inheritedDiags...)
	}
	executed := make([]string, 0, len(registry.Checks))
	for _, check := range registry.Checks {
		if check.RunsAt(tier) && (check.Meta || executedSet[check.Name]) {
			executed = append(executed, check.Name)
		}
	}
	return executedSet, executed, inheritedNames, nil
}

func orderedConformanceSelection(tier registry.Tier, selected *string) (map[string]bool, []string, []string) {
	if selected == nil {
		return nil, nil, nil
	}
	want := map[string]bool{}
	if *selected == "" {
		return want, nil, nil
	}
	names := strings.Split(*selected, ",")
	ordered, err := registry.CanonicalOrdinarySelection(tier, names)
	if err != nil || !slices.Equal(ordered, names) {
		return nil, nil, []string{fmt.Sprintf("conformance ordered selection %q is invalid for the %s tier", *selected, tier)}
	}
	for _, name := range ordered {
		want[name] = true
	}
	return want, ordered, nil
}

// checkCanaryFixtureCompliance grades kit-compliance rules directly against a marked
// immutable fixture tree. An ordinary live root has no marker and pays no duplicate
// check; a mutation test supplies the marker and calls this owner without another gate.
func checkCanaryFixtureCompliance(root string) []string {
	if !exists(filepath.Join(root, ".bench-compliance-canary")) {
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
	// The scan lives in the canary package, which owns the fixture-tree layout and asserts
	// the same binding at the top of its sweep. One derivation serves both callers, so the
	// check and the sweep cannot disagree about which family is unbound.
	return append(diags, canary.UnboundConformanceFamilies(kitRoot)...)
}

// familyIsBound reports whether the registry's family table binds family to a check,
// which is a separate question from whether the directory is a family at all — the
// canary package answers that one.
func familyIsBound(family string) bool {
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
		// The scrub is symmetric across every conformance control var: any one
		// leaking into a probe subprocess is the recursive-cascade shape.
		if strings.HasPrefix(kv, "BENCH_CONFORMANCE_ROOT=") ||
			strings.HasPrefix(kv, registry.ConformanceTierEnv+"=") ||
			strings.HasPrefix(kv, registry.ConformanceChecksEnv+"=") ||
			strings.HasPrefix(kv, registry.ConformanceInheritedEnv+"=") {
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
	content := "#!/usr/bin/env bash\nexec " + sanitize.ShellQuote(realBench) + " \"$@\"\n"
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
