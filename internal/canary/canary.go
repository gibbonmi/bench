// Package canary runs the gate against known-broken fixture roots and proves each
// fixture still triggers its targeted diagnostic.
package canary

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
)

const absentHarnessMessage = "canary harness absent — tests/canary/ has no fixtures; the gate cannot prove its own checks bite"

// PhaseEnv carries the one owning benchkit gate phase for a nested fixture.
// An absent value means the fixture must exercise the full inner gate.
const PhaseEnv = "BENCH_CANARY_PHASE"

// FixturePhase maps the canary directory convention to the phase that owns it.
// Legacy flat fixtures have tests/canary as their parent and keep the full gate.
func FixturePhase(family string) string {
	switch family {
	case "", "canary":
		return ""
	case "behavior-owned":
		return "contract"
	default:
		return "conformance"
	}
}

// expectFileName holds the diagnostic a fixture's inner gate must emit. Its presence
// directly under tests/canary/ doubles as the marker of a legacy flat fixture, since a
// family directory holds fixture directories and never an expectation of its own.
const expectFileName = "EXPECT"

// IsConformanceFamily reports whether dir, a directory directly under tests/canary/, is
// a conformance family holding fixtures rather than a fixture in its own right. This
// package owns the fixture-tree layout, so a caller grading the tree against the
// registry's family table asks here instead of rederiving either half: a second reading
// of the flat-fixture marker fails open, reporting every flat fixture as an unbound
// family the day the marker moves.
func IsConformanceFamily(dir string) bool {
	return FixturePhase(filepath.Base(dir)) == "conformance" && !isFlatFixture(dir)
}

// isFlatFixture reports whether dir is a legacy flat fixture — one living directly under
// tests/canary/ instead of inside a family.
func isFlatFixture(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, expectFileName))
	return err == nil
}

// RunCall is one inner gate invocation: Cwd is the materialized repo under grade,
// Gate is the real gate script from the root being checked, and FixtureDir names
// the source fixture or is empty for the vacuity baseline.
type RunCall struct {
	Cwd        string
	Gate       string
	FixtureDir string
	Env        []string
}

// RunResult captures the inner gate verdict and combined output.
type RunResult struct {
	ExitCode int
	Output   string
}

// Runner runs one inner gate invocation.
type Runner func(RunCall) RunResult

// Run is the `bench canary [root]` command.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 1 {
		fmt.Fprintln(stderr, "usage: bench canary [root]")
		return 2
	}
	root := ""
	if len(args) == 1 && args[0] != "" {
		root = args[0]
	} else {
		var err error
		root, err = git.Root()
		if err != nil {
			fmt.Fprintln(stderr, toon.NotInRepo())
			return 1
		}
	}
	if err := Sweep(root, defaultRunner); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, "canary ok (every fixture bit)")
	return 0
}

// Sweep runs the dev-tier canary harness for root — the sweep `bench canary` and the
// gate's canary phase perform. It returns all attributable harness failures it observes
// instead of stopping at the first one.
func Sweep(root string, runner Runner) error {
	return SweepTier(root, registry.Dev, runner)
}

// SweepShip runs the ship-tier fixtures for root through the real inner gate. It is the
// canary step of `bench prep-release`, which is the only surface that runs the ship-tier
// checks those fixtures grade.
func SweepShip(root string) error {
	return SweepTier(root, registry.Ship, defaultRunner)
}

// SweepTier runs the canary harness for the fixtures tier owns.
func SweepTier(root string, tier registry.Tier, runner Runner) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	all, err := fixtures(filepath.Join(root, "tests", "canary"))
	if err != nil {
		return err
	}
	fixtures, err := selectTier(all, tier)
	if err != nil {
		return err
	}
	if len(fixtures) == 0 {
		// Nothing to grade, and the vacuity baselines below are inner gate runs —
		// too expensive to pay for a tier that owns no fixtures.
		return nil
	}
	gate := filepath.Join(root, ".bench", "gate.sh")
	env := innerEnv(tier)

	baselines, err := scopeBaselines(fixtures, gate, env, runner)
	if err != nil {
		return err
	}

	errs := runFixtures(root, fixtures, baselines, gate, env, runner)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// selected is one fixture a tier sweeps, paired with the conformance check its inner
// run is scoped to. An empty scope runs the whole tier: contract and legacy flat
// fixtures grade surfaces no single conformance check owns.
type selected struct {
	dir   string
	scope string
}

// checkFileName optionally binds a fixture to the conformance check it grades. Absent,
// the fixture grades a check the dev gate runs — which is all but two of them, so the
// binding is written only where it changes the answer.
const checkFileName = "CHECK"

// fixtureCheck reports which tier sweeps a fixture and which check its CHECK file
// names. The fixture never states a tier of its own: it names its check, and the tier
// is read from the registry entry, so a check that is retiered takes its fixtures with
// it and the two cannot disagree. One read answers both questions for the same reason
// — the tier and the inner run's scope must not be able to disagree about what the
// file says. A fixture that names a check the registry has since renamed away is an
// error rather than a silent demotion to dev, where it would report "did not bite"
// forever. Only absence means dev and no name: a file present but holding no name is
// an error of its own, since dev is what deleting the file asks for and a blank file
// is far likelier to be a truncated write than an intent.
func fixtureCheck(fx string) (registry.Tier, string, error) {
	data, err := os.ReadFile(filepath.Join(fx, checkFileName))
	if errors.Is(err, os.ErrNotExist) {
		return registry.Dev, "", nil
	}
	if err != nil {
		return "", "", err
	}
	name := strings.TrimSpace(string(data))
	if name == "" {
		return "", "", fmt.Errorf("canary fixture '%s' has an empty %s file, which names no check; delete the file to sweep the fixture at the dev tier", filepath.Base(fx), checkFileName)
	}
	check, carried := registry.Find(name)
	if !carried {
		return "", "", fmt.Errorf("canary fixture '%s' names check %q, which the conformance registry does not carry", filepath.Base(fx), name)
	}
	return check.Tier, check.Name, nil
}

// fixtureScope names the one check a fixture's inner gate runs, or nothing for a
// fixture whose bite needs the full tier. A conformance family's fixtures all grade
// the check the registry binds the family to, and a fixture's own CHECK file overrides
// that for the strays living in a family they do not grade.
//
// A family the table does not bind falls back to the full inner gate. This sweep runs
// against every adopting repo, whose family names a kit-owned table will never carry,
// so the fallback is the only correct answer at this layer. The kit's own unbound
// family is a red the conformance layer raises instead, where the table and the tree
// it describes are both in scope.
func fixtureScope(fx, checkName string) string {
	family := filepath.Base(filepath.Dir(fx))
	if FixturePhase(family) != "conformance" {
		return ""
	}
	if checkName != "" {
		return checkName
	}
	scope, _ := registry.FamilyCheck(family)
	return scope
}

// selectTier keeps the fixtures tier sweeps, each with its resolved scope. Membership
// is tier equality, not the registry's RunsAt superset: the tiers have to partition the
// harness so that every fixture is swept by exactly one of them, and a fixture
// belonging to neither is the unswept rot the canary exists to catch.
func selectTier(fixtures []string, tier registry.Tier) ([]selected, error) {
	var out []selected
	for _, fx := range fixtures {
		owner, checkName, err := fixtureCheck(fx)
		if err != nil {
			return nil, err
		}
		if owner != tier {
			continue
		}
		out = append(out, selected{dir: fx, scope: fixtureScope(fx, checkName)})
	}
	return out, nil
}

// scopeBaselines runs one empty-tree gate per scope group the sweep grades, and
// returns each group's output for the vacuity comparison. A fixture's EXPECT is
// compared against a run of its own shape: another group's baseline executes different
// checks, so it would both miss a genuinely vacuous EXPECT and flag a sound one. The
// key is the resolved check alone, which keeps every unscoped fixture on the single
// full baseline they share today.
//
// A group whose baseline prints nothing is an error rather than a group that runs on:
// the vacuity test asks whether the baseline output already contains the EXPECT, and
// an empty output contains none of them, so every fixture in that group would clear
// the check unguarded while the other groups stayed graded. Errors for all such groups
// are reported together, matching how the sweep reports its fixture failures.
func scopeBaselines(fixtures []selected, gate string, env []string, runner Runner) (map[string]string, error) {
	var scopes []string
	seen := map[string]bool{}
	for _, fx := range fixtures {
		if !seen[fx.scope] {
			seen[fx.scope] = true
			scopes = append(scopes, fx.scope)
		}
	}

	dirs := make([]string, len(scopes))
	defer func() {
		for _, dir := range dirs {
			if dir != "" {
				os.RemoveAll(dir)
			}
		}
	}()
	for idx := range scopes {
		dir, err := os.MkdirTemp("", "bench-canary-empty-*")
		if err != nil {
			return nil, err
		}
		dirs[idx] = dir
		_ = gitInit(dir)
	}

	outputs := make([]string, len(scopes))
	eachIndex(len(scopes), func(idx int) {
		outputs[idx] = runner(RunCall{Cwd: dirs[idx], Gate: gate, Env: scopedEnv(env, scopes[idx])}).Output
	})

	baselines := make(map[string]string, len(scopes))
	var empty []string
	for idx, scope := range scopes {
		if outputs[idx] == "" {
			empty = append(empty, fmt.Sprintf("canary baseline for %s produced no output, so it can grade no EXPECT as vacuous", groupLabel(scope)))
			continue
		}
		baselines[scope] = outputs[idx]
	}
	if len(empty) > 0 {
		return nil, errors.New(strings.Join(empty, "\n"))
	}
	return baselines, nil
}

// groupLabel names a scope group in a diagnostic. The unscoped group's key is the
// empty string, which reads as a missing name rather than as the group every fixture
// needing the full inner gate shares, so it is named by what it runs.
func groupLabel(scope string) string {
	if scope == "" {
		return "the unscoped group"
	}
	return fmt.Sprintf("scope group %q", scope)
}

func runFixtures(root string, fixtures []selected, baselines map[string]string, gate string, env []string, runner Runner) []string {
	errs := make([]string, len(fixtures))
	eachIndex(len(fixtures), func(idx int) {
		errs[idx] = runFixture(root, fixtures[idx], baselines[fixtures[idx].scope], gate, env, runner)
	})

	out := errs[:0]
	for _, err := range errs {
		if err != "" {
			out = append(out, err)
		}
	}
	return out
}

// eachIndex runs body over every index below count under the sweep's worker budget,
// and returns once all of them have finished. Baselines and fixtures share the one
// budget: an inner gate costs the same either way, and the fixtures cannot start until
// the baseline of their group has an output to compare against.
func eachIndex(count int, body func(int)) {
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range fixtureWorkers(runtime.GOMAXPROCS(0), count) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				body(idx)
			}
		}()
	}
	for idx := range count {
		jobs <- idx
	}
	close(jobs)
	wg.Wait()
}

// fixtureWorkers floors at one so a small budget still makes progress, caps at
// fixtureCount so idle workers never outnumber the work, and takes budget as a
// parameter rather than reading the machine so it stays pure and testable.
func fixtureWorkers(budget, fixtureCount int) int {
	workers := budget / bounds.CanaryInnerWidth
	if workers < 1 {
		workers = 1
	}
	if workers > fixtureCount {
		workers = fixtureCount
	}
	return workers
}

func runFixture(root string, fixture selected, baselineOutput, gate string, env []string, runner Runner) string {
	fx := fixture.dir
	name := filepath.Base(fx)
	expectPath := filepath.Join(fx, expectFileName)
	filesDir := filepath.Join(fx, "files")

	expBytes, err := os.ReadFile(expectPath)
	if err != nil {
		return fmt.Sprintf("canary fixture '%s' has no EXPECT file", name)
	}
	basePath := filepath.Join(fx, "BASE")
	if info, err := os.Stat(filesDir); (err != nil || !info.IsDir()) && !regularFile(basePath) {
		return fmt.Sprintf("canary fixture '%s' has no files/ tree", name)
	}
	expect := trimExpectation(expBytes)
	if strings.Contains(baselineOutput, expect) {
		return fmt.Sprintf("canary '%s' EXPECT is vacuous (also matches an empty fixture)", name)
	}

	work, err := os.MkdirTemp("", "bench-canary-"+name+"-*")
	if err != nil {
		return fmt.Sprintf("canary '%s' setup failed: %v", name, err)
	}
	defer os.RemoveAll(work)
	if err := materializeMutationFixture(root, fx, work); err != nil {
		return fmt.Sprintf("canary '%s' setup failed: %v", name, err)
	}
	_ = gitInit(work)
	fixtureEnv := scopedEnv(env, fixture.scope)
	if phase := FixturePhase(filepath.Base(filepath.Dir(fx))); phase != "" {
		fixtureEnv = append(fixtureEnv, PhaseEnv+"="+phase)
	}
	result := runner(RunCall{Cwd: work, Gate: gate, FixtureDir: fx, Env: fixtureEnv})
	if result.ExitCode == 0 || !strings.Contains(result.Output, expect) {
		return fmt.Sprintf("canary '%s' did not bite (want red + %q; got exit %d)", name, expect, result.ExitCode)
	}
	return ""
}

func fixtures(dir string) ([]string, error) {
	families, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New(absentHarnessMessage)
	}
	var out []string
	seen := map[string]string{}
	addFixture := func(name, fixtureDir string) error {
		if first := seen[name]; first != "" {
			return fmt.Errorf("canary fixture name %q appears in multiple families; base names must be globally unique", name)
		}
		seen[name] = fixtureDir
		out = append(out, fixtureDir)
		return nil
	}
	for _, family := range families {
		if !family.IsDir() {
			continue
		}
		familyDir := filepath.Join(dir, family.Name())
		if isFlatFixture(familyDir) {
			if err := addFixture(family.Name(), familyDir); err != nil {
				return nil, err
			}
			continue
		}
		entries, err := os.ReadDir(familyDir)
		if err != nil {
			return nil, err
		}
		for _, ent := range entries {
			if !ent.IsDir() {
				continue
			}
			name := ent.Name()
			fixtureDir := filepath.Join(familyDir, name)
			if err := addFixture(name, fixtureDir); err != nil {
				return nil, err
			}
		}
	}
	if len(out) == 0 {
		return nil, errors.New(absentHarnessMessage)
	}
	sort.Strings(out)
	return out, nil
}

func trimExpectation(data []byte) string {
	return strings.TrimRight(string(data), "\n")
}

// MaterializeFixture copies a canary files/ tree into dst and restores dot- path
// segments to dot directories, matching the real canary sweep.
func MaterializeFixture(src, dst string) error {
	return materialize(src, dst)
}

func materialize(src, dst string) error {
	if err := copyTree(src, dst); err != nil {
		return err
	}
	return restoreDotSegments(dst)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

func restoreDotSegments(root string) error {
	var dirs []string
	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), "dot-") {
			dirs = append(dirs, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, old := range dirs {
		newPath := filepath.Join(filepath.Dir(old), "."+strings.TrimPrefix(filepath.Base(old), "dot-"))
		if err := os.Rename(old, newPath); err != nil {
			return err
		}
	}
	return nil
}

// scopedEnv is env plus the check an inner gate is scoped to. A run with no scope
// carries no scope variable at all rather than an empty one, which names no check and
// reds the inner gate. The copy is what keeps concurrent runs from appending into one
// shared backing array.
func scopedEnv(env []string, scope string) []string {
	out := append([]string(nil), env...)
	if scope == "" {
		return out
	}
	return append(out, registry.ConformanceCheckEnv+"="+scope)
}

// innerEnv is the environment every inner gate of a tier's sweep runs under. The tier
// is pinned rather than inherited: a fixture grades a check its own tier runs, so an
// inner gate at any other tier skips that check and the fixture reports "did not bite"
// forever. Every variable an inner gate may carry is scrubbed from the inherited
// environment here — a strip without its matching set, or a set without its matching
// strip, hands an ambient export control of what the sweep grades. The phase and the
// check scope vary per run, so they are stripped here and set by the caller that knows
// the run's fixture and group.
func innerEnv(tier registry.Tier) []string {
	env := make([]string, 0, len(os.Environ())+3)
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") || strings.HasPrefix(kv, "BENCH_CANARY_INNER=") || strings.HasPrefix(kv, PhaseEnv+"=") || strings.HasPrefix(kv, registry.ConformanceTierEnv+"=") || strings.HasPrefix(kv, registry.ConformanceCheckEnv+"=") || strings.HasPrefix(kv, "GOMAXPROCS=") {
			continue
		}
		env = append(env, kv)
	}
	return append(env, "BENCH_CANARY_INNER=1", registry.ConformanceTierEnv+"="+string(tier), fmt.Sprintf("GOMAXPROCS=%d", bounds.CanaryInnerWidth))
}

func defaultRunner(call RunCall) RunResult {
	cmd := exec.Command("bash", call.Gate)
	cmd.Dir = call.Cwd
	cmd.Env = call.Env
	r := subprocess.CaptureMerged(cmd)
	output := r.Stdout
	// A spawn failure (ProcessState nil) writes nothing, so append the error.
	if r.Err != nil && cmd.ProcessState == nil {
		output += r.Err.Error()
	}
	return RunResult{ExitCode: r.ExitCode, Output: output}
}

func gitInit(dir string) error {
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	return cmd.Run()
}
