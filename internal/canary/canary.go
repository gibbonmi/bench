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
		// Nothing to grade, and the vacuity baseline below is a full inner gate run —
		// too expensive to pay for a tier that owns no fixtures.
		return nil
	}
	gate := filepath.Join(root, ".bench", "gate.sh")
	env := innerEnv(tier)

	baselineDir, err := os.MkdirTemp("", "bench-canary-empty-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(baselineDir)
	_ = gitInit(baselineDir)
	baseline := runner(RunCall{Cwd: baselineDir, Gate: gate, Env: env})

	errs := runFixtures(root, fixtures, baseline.Output, gate, env, runner)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// checkFileName optionally binds a fixture to the conformance check it grades. Absent,
// the fixture grades a check the dev gate runs — which is all but two of them, so the
// binding is written only where it changes the answer.
const checkFileName = "CHECK"

// fixtureTier reports which tier sweeps a fixture. The fixture never states a tier of
// its own: it names its check, and the tier is read from the registry entry, so a check
// that is retiered takes its fixtures with it and the two cannot disagree. A fixture
// that names a check the registry has since renamed away is an error rather than a
// silent demotion to dev, where it would report "did not bite" forever.
func fixtureTier(fx string) (registry.Tier, error) {
	data, err := os.ReadFile(filepath.Join(fx, checkFileName))
	if errors.Is(err, os.ErrNotExist) {
		return registry.Dev, nil
	}
	if err != nil {
		return "", err
	}
	name := strings.TrimSpace(string(data))
	for _, check := range registry.Checks {
		if check.Name == name {
			return check.Tier, nil
		}
	}
	return "", fmt.Errorf("canary fixture '%s' names check %q, which the conformance registry does not carry", filepath.Base(fx), name)
}

// selectTier keeps the fixtures tier sweeps. Membership is tier equality, not the
// registry's RunsAt superset: the tiers have to partition the harness so that every
// fixture is swept by exactly one of them, and a fixture belonging to neither is the
// unswept rot the canary exists to catch.
func selectTier(fixtures []string, tier registry.Tier) ([]string, error) {
	var out []string
	for _, fx := range fixtures {
		owner, err := fixtureTier(fx)
		if err != nil {
			return nil, err
		}
		if owner == tier {
			out = append(out, fx)
		}
	}
	return out, nil
}

func runFixtures(root string, fixtures []string, baselineOutput, gate string, env []string, runner Runner) []string {
	errs := make([]string, len(fixtures))
	jobs := make(chan int)
	workers := fixtureWorkers(runtime.GOMAXPROCS(0), len(fixtures))

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				errs[idx] = runFixture(root, fixtures[idx], baselineOutput, gate, env, runner)
			}
		}()
	}
	for idx := range fixtures {
		jobs <- idx
	}
	close(jobs)
	wg.Wait()

	out := errs[:0]
	for _, err := range errs {
		if err != "" {
			out = append(out, err)
		}
	}
	return out
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

func runFixture(root, fx, baselineOutput, gate string, env []string, runner Runner) string {
	name := filepath.Base(fx)
	expectPath := filepath.Join(fx, "EXPECT")
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
	fixtureEnv := append([]string(nil), env...)
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
		if _, err := os.Stat(filepath.Join(familyDir, "EXPECT")); err == nil {
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

// innerEnv is the environment every inner gate of a tier's sweep runs under. The tier
// is pinned rather than inherited: a fixture grades a check its own tier runs, so an
// inner gate at any other tier skips that check and the fixture reports "did not bite"
// forever. Every variable set here is scrubbed from the inherited environment first —
// a strip without its matching set, or a set without its matching strip, hands an
// ambient export control of what the sweep grades.
func innerEnv(tier registry.Tier) []string {
	env := make([]string, 0, len(os.Environ())+3)
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") || strings.HasPrefix(kv, "BENCH_CANARY_INNER=") || strings.HasPrefix(kv, PhaseEnv+"=") || strings.HasPrefix(kv, registry.ConformanceTierEnv+"=") || strings.HasPrefix(kv, "GOMAXPROCS=") {
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
