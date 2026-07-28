// Package canary runs the gate against known-broken fixture roots and proves each
// fixture still triggers its targeted diagnostic.
package canary

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
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

// ContractPackageEnv carries the one contract package a behavior-owned fixture's EXPECT
// belongs to, as a slash path relative to internal/contract. It narrows the contract
// phase's argv to that package alone; an absent value leaves the phase grading every
// contract package, which is what every other run wants.
const ContractPackageEnv = "BENCH_CANARY_CONTRACT_PACKAGE"

// PhaseManifestPath is where a graded root declares a phase table of its own, relative
// to that root. Declaring one replaces the built-in table outright, so a canary fixture
// whose tree carries it never reaches the narrowing the built-in contract phase applies.
const PhaseManifestPath = ".bench/phases.json"

// contractPhase runs the contract suite, and is the phase a behavior-owned fixture's
// EXPECT is emitted by. It is the one phase a fixture can narrow further, to the single
// contract package that owns its diagnostic.
const contractPhase = "contract"

// FixturePhase maps the canary directory convention to the phase that owns it.
// Legacy flat fixtures have tests/canary as their parent and keep the full gate.
func FixturePhase(family string) string {
	switch family {
	case "", "canary":
		return ""
	case "behavior-owned":
		return contractPhase
	}
	if phaseFamilies[family] {
		return family
	}
	return "conformance"
}

// The toolchain gate phases, named here because the fixture router and the phase table
// must agree on every one of them and only one of the two can own the string. This
// package is the owner: internal/gate imports it to build the table, and the edge does
// not run the other way, so a name typed a second time in the table could disagree with
// the routing forever.
const (
	PhaseBuild            = "build"
	PhaseGofmt            = "gofmt"
	PhaseVet              = "vet"
	PhaseTest             = "test"
	PhaseRace             = "race"
	PhaseConformanceSuite = "conformance-suite"
)

// phaseFamilies are the family names that name a gate phase rather than a conformance
// check. A fixture under one of them runs only the phase that owns its failure, instead
// of the every-non-canary-phase run a family the gate cannot attribute has to pay for.
var phaseFamilies = map[string]bool{
	PhaseBuild:            true,
	PhaseGofmt:            true,
	PhaseVet:              true,
	PhaseTest:             true,
	PhaseRace:             true,
	PhaseConformanceSuite: true,
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
	return FixturePhase(filepath.Base(dir)) == "conformance" && !holdsExpect(dir)
}

// UnboundConformanceFamilies reports the conformance family directories under kitRoot
// that the registry's family table binds to no check, one diagnostic apiece. The sweep
// scopes a fixture's inner gate by its family and falls back to a full run for a family
// it cannot resolve — correct for an adopting repo, whose families this table will never
// carry, but in the kit it is the whole per-fixture cost the scoping removes, paid in
// silence. Reading the tree is what catches it: a caller iterating the table instead sees
// nothing, because a family the table omits is invisible from that side.
func UnboundConformanceFamilies(kitRoot string) []string {
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	entries, err := os.ReadDir(canaryDir)
	if err != nil {
		return nil
	}
	var diags []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() || !IsConformanceFamily(filepath.Join(canaryDir, name)) {
			continue
		}
		if _, bound := registry.FamilyCheck(name); !bound {
			diags = append(diags, fmt.Sprintf("canary conformance family %q is bound to no conformance check; add it to the registry family table so its fixtures run scoped", name))
		}
	}
	return diags
}

// Fixture is one discovered canary fixture as a caller outside the sweep sees it: where
// it lives, the family that routes its inner gate, and the contract package that owns its
// EXPECT where one does. An empty Family is a legacy flat fixture, which belongs to no
// family and runs the full inner gate.
type Fixture struct {
	Dir     string
	Family  string
	Package string
}

// Fixtures maps each fixture's base name under canaryDir — a tests/canary tree — to what
// the sweep discovered about it. This package owns the layout, including the package
// nesting under the behavior-owned family, so a caller grading that tree against an
// inventory of its own asks here instead of walking it a second way: a second walk
// disagrees with the sweep the day the layout moves, and reports fixtures the sweep never
// runs — or reads a family from a parent directory that is a package path.
func Fixtures(canaryDir string) (map[string]Fixture, error) {
	found, err := fixtures(canaryDir)
	if err != nil {
		return nil, err
	}
	out := make(map[string]Fixture, len(found))
	for _, fx := range found {
		out[filepath.Base(fx.dir)] = Fixture{Dir: fx.dir, Family: fx.family, Package: fx.pkg}
	}
	return out, nil
}

// holdsExpect reports whether dir states an expectation of its own, which is what marks
// a directory as a fixture rather than a family or a package path above one.
func holdsExpect(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, expectFileName))
	return err == nil
}

// isDir reports whether path is a directory, following symlinks as the tools reading
// these trees do.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
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
	if err := assertFamilyBinding(root); err != nil {
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
	if err := assertContractScopes(root, fixtures); err != nil {
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

// assertFamilyBinding refuses a sweep of the kit's own tree while a conformance family
// there is bound to no check, ahead of the first inner gate — an unbound family costs a
// full unscoped run per fixture, so reporting it afterwards means paying for it first.
//
// BENCH_KIT is the discriminator: the wrapper always exports it as an absolute path, so
// equality means the kit is grading itself and the kit-owned family table is
// authoritative. Anything else — unset, empty, relative, or naming another tree — is not
// that export, so the sweep is grading a repo whose families a kit-owned table will never
// carry and the assertion stays out of the way. The unscoped fallback that repo already
// gets is the correct answer there, and refusing on it would red every adopting repo.
// Both sides resolve through EvalSymlinks first, because bin/bench.sh derives its default
// with a physical cd while this root is normalized with filepath.Abs alone; a raw compare
// would make a symlinked kit checkout read as an adopting repo and skip the kit's own
// assertion silently.
func assertFamilyBinding(root string) error {
	kit, err := filepath.EvalSymlinks(os.Getenv("BENCH_KIT"))
	if err != nil {
		return nil
	}
	swept, err := filepath.EvalSymlinks(root)
	if err != nil || swept != kit {
		return nil
	}
	if unbound := UnboundConformanceFamilies(root); len(unbound) > 0 {
		return errors.New(strings.Join(unbound, "\n"))
	}
	return nil
}

// selected is one fixture the sweep grades, with everything its inner run is scoped by:
// the family directory it lives under, the contract package that owns its EXPECT for a
// behavior-owned fixture, and the conformance check for a fixture whose family binds one.
// The family and package are carried rather than re-derived from the path at each use,
// because a nested package puts more than one directory between the family and the
// fixture and a parent-basename read would silently name the wrong one.
//
// An empty scope and an empty package together run the whole tier: legacy flat fixtures
// grade surfaces no single check or package owns.
type selected struct {
	dir    string
	family string
	pkg    string
	scope  string
}

// contractGroupPrefix keys a contract package's vacuity group, so a package path can
// never collide with a conformance check that happens to share its name.
const contractGroupPrefix = contractPhase + ":"

// group names the vacuity baseline this fixture's EXPECT is graded against. Fixtures of
// one group run the same checks, so a baseline of any other group both misses a genuinely
// vacuous EXPECT and flags a sound one.
func (s selected) group() string {
	if s.pkg != "" {
		return contractGroupPrefix + s.pkg
	}
	return s.scope
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
func fixtureScope(family, checkName string) string {
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
func selectTier(fixtures []selected, tier registry.Tier) ([]selected, error) {
	var out []selected
	for _, fx := range fixtures {
		owner, checkName, err := fixtureCheck(fx.dir)
		if err != nil {
			return nil, err
		}
		if owner != tier {
			continue
		}
		fx.scope = fixtureScope(fx.family, checkName)
		out = append(out, fx)
	}
	return out, nil
}

// assertContractScopes refuses a sweep while a behavior-owned fixture cannot be scoped to
// one contract package, ahead of the first inner gate. Each of these defects leaves the
// fixture grading every contract package instead of the one that owns its EXPECT — the
// whole per-fixture cost the scoping removes, paid in silence — so reporting them
// afterwards means paying for them first. Every defect is reported, matching how the
// sweep reports its fixture failures.
func assertContractScopes(root string, fixtures []selected) error {
	var errs []string
	for _, fx := range fixtures {
		if FixturePhase(fx.family) != contractPhase {
			continue
		}
		name := filepath.Base(fx.dir)
		switch {
		case fx.pkg == "":
			errs = append(errs, fmt.Sprintf("canary fixture '%s' sits directly under tests/canary/%s/, which binds it to no contract package; move it to tests/canary/%s/<package path>/%s/", name, fx.family, fx.family, name))
		case !isDir(filepath.Join(root, "internal", "contract", filepath.FromSlash(fx.pkg))):
			errs = append(errs, fmt.Sprintf("canary fixture '%s' is bound to contract package %q, which does not exist under internal/contract/", name, fx.pkg))
		case regularFile(filepath.Join(fx.dir, "files", dotEncodedPath(PhaseManifestPath))):
			errs = append(errs, fmt.Sprintf("canary fixture '%s' declares files/%s, and a declared phase table replaces the built-in one its contract scoping lives in; make it a legacy flat fixture at tests/canary/%s/ instead", name, dotEncodedPath(PhaseManifestPath), name))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// scopeBaselines runs one empty-tree gate per scope group the sweep grades, and
// returns each group's output for the vacuity comparison. A fixture's EXPECT is
// compared against a run of its own shape: another group's baseline executes different
// checks, so it would both miss a genuinely vacuous EXPECT and flag a sound one. The
// key is the group each fixture resolves to, which keeps every unscoped fixture on the
// single full baseline they share today.
//
// A group whose baseline prints nothing is an error rather than a group that runs on:
// the vacuity test asks whether the baseline output already contains the EXPECT, and
// an empty output contains none of them, so every fixture in that group would clear
// the check unguarded while the other groups stayed graded. Errors for all such groups
// are reported together, matching how the sweep reports its fixture failures.
//
// Each group is represented by one of its own fixtures, so the environment a baseline
// runs under is built by the same function that builds its fixtures': a baseline that
// ran a wider set of checks than its group grades EXPECTs against output no fixture in
// the group can produce.
func scopeBaselines(fixtures []selected, gate string, env []string, runner Runner) (map[string]string, error) {
	var scopes []selected
	seen := map[string]bool{}
	for _, fx := range fixtures {
		if key := fx.group(); !seen[key] {
			seen[key] = true
			scopes = append(scopes, fx)
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
		outputs[idx] = runner(RunCall{Cwd: dirs[idx], Gate: gate, Env: baselineEnv(env, scopes[idx])}).Output
	})

	baselines := make(map[string]string, len(scopes))
	var empty []string
	for idx, scope := range scopes {
		if outputs[idx] == "" {
			empty = append(empty, fmt.Sprintf("canary baseline for %s produced no output, so it can grade no EXPECT as vacuous", groupLabel(scope.group())))
			continue
		}
		baselines[scope.group()] = outputs[idx]
	}
	if len(empty) > 0 {
		return nil, errors.New(strings.Join(empty, "\n"))
	}
	return baselines, nil
}

// groupLabel names a scope group in a diagnostic. The unscoped group's key is the
// empty string, which reads as a missing name rather than as the group every fixture
// needing the full inner gate shares, so it is named by what it runs.
func groupLabel(group string) string {
	if group == "" {
		return "the unscoped group"
	}
	return fmt.Sprintf("scope group %q", group)
}

func runFixtures(root string, fixtures []selected, baselines map[string]string, gate string, env []string, runner Runner) []string {
	errs := make([]string, len(fixtures))
	eachIndex(len(fixtures), func(idx int) {
		errs[idx] = runFixture(root, fixtures[idx], baselines[fixtures[idx].group()], gate, env, runner)
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
	fixtureEnv := scopedEnv(env, fixture)
	if phase := FixturePhase(fixture.family); phase != "" {
		fixtureEnv = append(fixtureEnv, PhaseEnv+"="+phase)
	}
	result := runner(RunCall{Cwd: work, Gate: gate, FixtureDir: fx, Env: fixtureEnv})
	if result.ExitCode == 0 || !strings.Contains(result.Output, expect) {
		return fmt.Sprintf("canary '%s' did not bite (want red + %q; got exit %d)", name, expect, result.ExitCode)
	}
	return ""
}

// fixtures discovers every fixture under dir, the tests/canary tree: a legacy flat
// fixture directly under it, and otherwise the fixtures of each family directory. Base
// names are globally unique across the whole tree — two fixtures sharing a name collide
// in their diagnostics and in their temporary work directories wherever they live, so
// nesting is no escape from the check.
func fixtures(dir string) ([]selected, error) {
	families, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New(absentHarnessMessage)
	}
	var out []selected
	seen := map[string]bool{}
	addFixture := func(fixture selected) error {
		name := filepath.Base(fixture.dir)
		if seen[name] {
			return fmt.Errorf("canary fixture name %q appears in multiple families; base names must be globally unique", name)
		}
		seen[name] = true
		out = append(out, fixture)
		return nil
	}
	for _, family := range families {
		if !family.IsDir() {
			continue
		}
		name := family.Name()
		familyDir := filepath.Join(dir, name)
		if holdsExpect(familyDir) {
			if err := addFixture(selected{dir: familyDir}); err != nil {
				return nil, err
			}
			continue
		}
		if FixturePhase(name) == contractPhase {
			if err := addContractFixtures(familyDir, name, "", addFixture); err != nil {
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
			if err := addFixture(selected{dir: filepath.Join(familyDir, ent.Name()), family: name}); err != nil {
				return nil, err
			}
		}
	}
	if len(out) == 0 {
		return nil, errors.New(absentHarnessMessage)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].dir < out[j].dir })
	return out, nil
}

// addContractFixtures walks the family whose fixtures are bound to a contract package,
// where the directory holding EXPECT is the fixture and pkg accumulates every directory
// between the family and it. The descent is what lets a nested contract package
// (surface/artifact) bind its fixtures; a fixture reached with an empty pkg is bound to
// no package at all, which assertContractScopes reports.
func addContractFixtures(dir, family, pkg string, addFixture func(selected) error) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		child := filepath.Join(dir, ent.Name())
		if holdsExpect(child) {
			if err := addFixture(selected{dir: child, family: family, pkg: pkg}); err != nil {
				return err
			}
			continue
		}
		if err := addContractFixtures(child, family, path.Join(pkg, ent.Name()), addFixture); err != nil {
			return err
		}
	}
	return nil
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

// dotPrefix stands in for a leading dot in a fixture's stored files/ tree, so that a
// fixture's dot directories are visible to the tools that read the harness and are
// restored only in the materialized copy under grade.
const dotPrefix = "dot-"

// dotEncodedPath is the stored form of a slash path whose segments start with a dot —
// the shape a fixture's files/ tree holds, and the inverse of restoreDotSegments.
func dotEncodedPath(rel string) string {
	segments := strings.Split(rel, "/")
	for idx, segment := range segments {
		if bare, dotted := strings.CutPrefix(segment, "."); dotted {
			segments[idx] = dotPrefix + bare
		}
	}
	return filepath.Join(segments...)
}

func restoreDotSegments(root string) error {
	var dirs []string
	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), dotPrefix) {
			dirs = append(dirs, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, old := range dirs {
		newPath := filepath.Join(filepath.Dir(old), "."+strings.TrimPrefix(filepath.Base(old), dotPrefix))
		if err := os.Rename(old, newPath); err != nil {
			return err
		}
	}
	return nil
}

// scopedEnv is env plus the variables naming what one inner gate grades: the conformance
// check the fixture's family binds, and the contract package that owns a behavior-owned
// fixture's EXPECT. A run without one of them carries no variable at all rather than an
// empty one, which names nothing and reds the inner gate. The copy is what keeps
// concurrent runs from appending into one shared backing array.
func scopedEnv(env []string, fixture selected) []string {
	out := append([]string(nil), env...)
	if fixture.scope != "" {
		out = append(out, registry.ConformanceCheckEnv+"="+fixture.scope)
	}
	if fixture.pkg != "" {
		out = append(out, ContractPackageEnv+"="+fixture.pkg)
	}
	return out
}

// baselineEnv is the environment a group's empty-tree baseline runs under: the group's
// own scope variables, plus the phase pin for a contract group. The pin belongs to the
// baseline as much as to the fixture — a contract package narrows that one phase's argv,
// so a baseline running the full inner gate would grade its group's EXPECTs against the
// output of checks no fixture in the group runs.
func baselineEnv(env []string, group selected) []string {
	out := scopedEnv(env, group)
	if group.pkg != "" {
		out = append(out, PhaseEnv+"="+FixturePhase(group.family))
	}
	return out
}

// innerEnv is the environment every inner gate of a tier's sweep runs under. The tier
// is pinned rather than inherited: a fixture grades a check its own tier runs, so an
// inner gate at any other tier skips that check and the fixture reports "did not bite"
// forever. Every variable an inner gate may carry is scrubbed from the inherited
// environment here — a strip without its matching set, or a set without its matching
// strip, hands an ambient export control of what the sweep grades. The phase, the check
// scope, and the contract package vary per run, so they are stripped here and set by the
// caller that knows the run's fixture and group.
func innerEnv(tier registry.Tier) []string {
	env := make([]string, 0, len(os.Environ())+3)
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") || strings.HasPrefix(kv, "BENCH_CANARY_INNER=") || strings.HasPrefix(kv, PhaseEnv+"=") || strings.HasPrefix(kv, ContractPackageEnv+"=") || strings.HasPrefix(kv, registry.ConformanceTierEnv+"=") || strings.HasPrefix(kv, registry.ConformanceCheckEnv+"=") || strings.HasPrefix(kv, "GOMAXPROCS=") {
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
