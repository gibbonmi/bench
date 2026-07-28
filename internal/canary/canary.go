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
	"slices"
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

// SubjectRootEnv names the tree a contract test grades instead of the kit checkout it was
// compiled from. It lives beside the phase names for the same reason they do: the sweep,
// the gate's contract phase, and internal/contract's helper must agree on one spelling,
// and this package is the only one all three can read — internal/gate and
// internal/contract each import it, and neither edge runs the other way.
const SubjectRootEnv = "BENCH_CONTRACT_ROOT"

// PhaseManifestPath is where a graded root declares a phase table of its own, relative
// to that root. Declaring one replaces the built-in table outright.
const PhaseManifestPath = ".bench/phases.json"

// FixturePhase maps the canary directory convention to the phase that owns it.
// Legacy flat fixtures have tests/canary as their parent and keep the full gate.
func FixturePhase(family string) string {
	switch family {
	case "", "canary":
		return ""
	case "behavior-owned":
		return PhaseContract
	}
	if phaseFamilies[family] {
		return family
	}
	return "conformance"
}

// The gate phases the fixture router names, held here because the router and the phase
// table must agree on every one of them and only one of the two can own the string. This
// package is the owner: internal/gate imports it to build the table, and the edge does
// not run the other way, so a name typed a second time in the table could disagree with
// the routing forever.
//
// PhaseContract sits apart from the rest: no family is named after it — behavior-owned
// routes there — and it is the one phase whose fixtures spawn no gate at all, being
// graded by the compiled test binary of the contract package that owns the diagnostic.
const (
	PhaseBuild            = "build"
	PhaseGofmt            = "gofmt"
	PhaseVet              = "vet"
	PhaseTest             = "test"
	PhaseRace             = "race"
	PhaseConformanceSuite = "conformance-suite"
	PhaseContract         = "contract"
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

// filesDirName holds the tree a fixture materializes over its base. It is the only
// directory a fixture carries, which is what makes a fixture a leaf of the walk.
const filesDirName = "files"

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

// RunKind names the three things a sweep runs. A behavior-owned fixture's EXPECT is a
// contract test's own failure message, so that family is graded by compiling the owning
// package's test binary once and invoking it per subject tree; every other family needs a
// whole gate around its mutated tree and spawns one.
type RunKind int

const (
	RunGate RunKind = iota
	RunCompile
	RunBite
)

// RunCall is one thing the sweep runs. Kind says which, and the fields the other kinds
// need stay zero: a gate spawn carries Gate, a compile carries Package and Binary, and a
// bite invocation carries the Binary it runs and the Package that produced it. FixtureDir
// names the source fixture, or is empty for a vacuity baseline.
//
// Cwd is the working directory, which differs by kind: a gate spawn and a compile both
// run over a tree (the materialized repo under grade, and the swept root), while a bite
// invocation runs in its package's source directory — where `go test` runs a test binary —
// and carries the tree it grades in the environment instead.
type RunCall struct {
	Kind       RunKind
	Cwd        string
	Gate       string
	FixtureDir string
	Package    string
	Binary     string
	Env        []string
}

// RunResult captures the verdict and combined output of one run.
type RunResult struct {
	ExitCode int
	Output   string
}

// Runner runs one call of any kind.
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
		// Nothing to grade, and the compiles and vacuity baselines below are real runs —
		// too expensive to pay for a tier that owns no fixtures.
		return nil
	}
	base := sweepEnv()
	run := sweepRun{
		root:    root,
		gate:    filepath.Join(root, ".bench", "gate.sh"),
		base:    base,
		gateEnv: gateEnv(base, tier),
	}

	binDir, binaries, err := compileContractPackages(root, fixtures, base, runner)
	if binDir != "" {
		// Every return below this point is an exit path the binaries must not survive,
		// including the compile's own failure, which reports a directory it has already
		// half-filled.
		defer os.RemoveAll(binDir)
	}
	if err != nil {
		return err
	}
	run.binaries = binaries

	baselines, err := scopeBaselines(fixtures, run, runner)
	if err != nil {
		return err
	}

	errs := runFixtures(fixtures, baselines, run, runner)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// sweepRun is everything a graded run needs that does not vary from fixture to fixture:
// the swept tree, the gate script a spawning family runs, the two environments, and the
// compiled test binary of each contract package.
type sweepRun struct {
	root     string
	gate     string
	base     []string
	gateEnv  []string
	binaries map[string]string
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
const contractGroupPrefix = PhaseContract + ":"

// group names the baseline this fixture's EXPECT is graded against. Fixtures of one
// group run the same checks, so another group's baseline is output no run in this group
// can produce; what the comparison does and does not establish is stated on
// scopeBaselines.
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

// assertContractScopes refuses a sweep while a behavior-owned fixture cannot be bound to
// one contract package, ahead of the first run. A fixture whose binding names no package,
// or a package the swept tree does not hold, has no test binary to be graded by at all, so
// the sweep says which fixture and stops rather than reaching a compile that would name
// the package without naming the fixture that asked for it. Every defect is reported,
// matching how the sweep reports its fixture failures.
func assertContractScopes(root string, fixtures []selected) error {
	var errs []string
	for _, fx := range fixtures {
		if FixturePhase(fx.family) != PhaseContract {
			continue
		}
		name := filepath.Base(fx.dir)
		switch {
		case fx.pkg == "":
			errs = append(errs, fmt.Sprintf("canary fixture '%s' sits directly under tests/canary/%s/, which binds it to no contract package; move it to tests/canary/%s/<package path>/%s/", name, fx.family, fx.family, name))
		case !isDir(contractPackageDir(root, fx.pkg)):
			errs = append(errs, fmt.Sprintf("canary fixture '%s' is bound to contract package %q, which does not exist under internal/contract/", name, fx.pkg))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// contractPackageDir is where a contract package's source lives under a swept tree. The
// binding a fixture's directory nesting states is resolved against the swept root, so the
// package the sweep compiles and the package the structural refusal accepts are the same
// directory by construction.
func contractPackageDir(root, pkg string) string {
	return filepath.Join(root, "internal", "contract", filepath.FromSlash(pkg))
}

// contractPackages lists, sorted, each contract package the swept fixtures bind to. One
// entry is one compile: every fixture of a package reuses that package's binary, which is
// the whole saving over a process tree per fixture.
func contractPackages(fixtures []selected) []string {
	seen := map[string]bool{}
	var out []string
	for _, fx := range fixtures {
		if fx.pkg == "" || seen[fx.pkg] {
			continue
		}
		seen[fx.pkg] = true
		out = append(out, fx.pkg)
	}
	sort.Strings(out)
	return out
}

// compileContractPackages builds each bound contract package's test binary into one
// sweep-owned temporary directory, and returns that directory alongside the binary path
// of every package. The directory is returned even when the compile fails, because the
// caller owns removing it on every exit path and a failed compile still leaves files.
//
// A compile that exits nonzero and a compile that exits zero having written nothing are
// both errors naming the package. The second is not hypothetical: `go test -c` on a
// package holding no test files succeeds and writes no binary, so an exit-code-only check
// would go on to invoke a path that does not exist and surface an exec error naming
// nothing an author can act on.
func compileContractPackages(root string, fixtures []selected, base []string, runner Runner) (string, map[string]string, error) {
	pkgs := contractPackages(fixtures)
	if len(pkgs) == 0 {
		return "", nil, nil
	}
	dir, err := os.MkdirTemp("", "bench-canary-bin-*")
	if err != nil {
		return "", nil, err
	}
	binaries := make(map[string]string, len(pkgs))
	for _, pkg := range pkgs {
		binaries[pkg] = filepath.Join(dir, contractBinaryName(pkg))
	}

	errs := make([]string, len(pkgs))
	eachIndex(len(pkgs), func(idx int) {
		pkg := pkgs[idx]
		result := runner(RunCall{Kind: RunCompile, Cwd: root, Package: pkg, Binary: binaries[pkg], Env: base})
		switch {
		case result.ExitCode != 0:
			errs[idx] = fmt.Sprintf("canary could not compile contract package %q (exit %d): %s", pkg, result.ExitCode, strings.TrimSpace(result.Output))
		case !regularFile(binaries[pkg]):
			errs[idx] = fmt.Sprintf("canary compiled contract package %q to no binary, so none of its fixtures can be graded; a package holding no test files compiles to nothing", pkg)
		}
	})

	var failed []string
	for _, msg := range errs {
		if msg != "" {
			failed = append(failed, msg)
		}
	}
	if len(failed) > 0 {
		return dir, nil, errors.New(strings.Join(failed, "\n"))
	}
	return dir, binaries, nil
}

// contractBinaryName is the file name one package path compiles to. The escape character
// is encoded before the separator is, which is what makes the encoding injective: without
// that first pass a path already containing the escape would collide with a differently
// nested one, and the loser of the collision would grade the wrong package's tests.
func contractBinaryName(pkg string) string {
	return strings.ReplaceAll(strings.ReplaceAll(pkg, "_", "_u"), "/", "_s") + ".test"
}

// scopeBaselines runs one empty-tree run per scope group the sweep grades, and
// returns each group's output for the baseline comparison. A fixture's EXPECT is
// compared against a run of its own shape; the key is the group each fixture resolves
// to, which keeps every unscoped fixture on the single full baseline they share today.
//
// For a contract group the comparison is a collision screen, not a vacuity check. The
// empty tree is a degenerate tree, not an unmutated twin of any fixture's: a baseline
// match means the EXPECT collides with the infrastructure noise a test binary emits
// over a tree it cannot grade — missing files, exit-status text — and a miss says
// nothing about whether the fixture's mutation is what makes the EXPECT appear.
// Mutation-specificity would take a per-fixture unmutated twin, which no group has.
// The screen stays because it costs one run per group and rejects an EXPECT sloppy
// enough to match that noise. Nothing reads a baseline's exit code, and nothing may
// start to: the groups share no exit shape — one whose tests all skip over a
// subjectless tree exits green, one whose tests reach the tree reds loudly — so any
// invariant keyed on it is broken on arrival.
//
// A group whose baseline prints nothing is an error rather than a group that runs on:
// the baseline check asks whether the baseline output already contains the EXPECT, and
// an empty output contains none of them, so every fixture in that group would clear
// the check unguarded while the other groups stayed graded. Errors for all such groups
// are reported together, matching how the sweep reports its fixture failures.
//
// Each group is represented by one of its own fixtures, so a baseline call is built by
// the same function that builds its fixtures': a baseline of a different scope, or a gate
// spawn where the group's fixtures invoke a test binary, grades its group's EXPECTs
// against output no fixture in the group can produce.
//
// The phase pin is the one axis where the two deliberately part — a baseline runs the
// whole inner gate even where its group's fixtures are each pinned to a single phase —
// because the two directions fail differently. A baseline narrower than the runs it
// grades passes every vacuous EXPECT in its group in silence; a wider one at worst calls
// a sound EXPECT vacuous, which is a red someone reads.
func scopeBaselines(fixtures []selected, run sweepRun, runner Runner) (map[string]string, error) {
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
		outputs[idx] = runner(subjectCall(scopes[idx], dirs[idx], "", run, noPhasePin)).Output
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

func runFixtures(fixtures []selected, baselines map[string]string, run sweepRun, runner Runner) []string {
	errs := make([]string, len(fixtures))
	eachIndex(len(fixtures), func(idx int) {
		errs[idx] = runFixture(fixtures[idx], baselines[fixtures[idx].group()], run, runner)
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
// and returns once all of them have finished. Compiles, baselines, and fixtures share the
// one budget, and run as three sequenced stages: a group's fixtures cannot start until
// its package is compiled and its baseline has an output to compare against.
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

func runFixture(fixture selected, baselineOutput string, run sweepRun, runner Runner) string {
	fx := fixture.dir
	name := filepath.Base(fx)
	expectPath := filepath.Join(fx, expectFileName)
	filesDir := filepath.Join(fx, filesDirName)

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
	if err := materializeMutationFixture(run.root, fx, work); err != nil {
		return fmt.Sprintf("canary '%s' setup failed: %v", name, err)
	}
	_ = gitInit(work)
	result := runner(subjectCall(fixture, work, fx, run, pinFamilyPhase))
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
		if FixturePhase(name) == PhaseContract {
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
			leaf, err := fixtureLeaf(child)
			if err != nil {
				return err
			}
			if !leaf {
				return fmt.Errorf("canary directory '%s' holds an %s and further directories; a fixture is a leaf carrying only %s/, so every fixture below this %s goes unswept", path.Join(family, pkg, ent.Name()), expectFileName, filesDirName, expectFileName)
			}
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

// fixtureLeaf reports whether dir carries nothing but its files/ tree and its own marker
// files. The walk stops at the first EXPECT it finds, so a package directory that also
// holds one takes the place of every fixture beneath it — a harness that grades those
// fixtures no longer, and says nothing about it.
func fixtureLeaf(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, ent := range entries {
		if ent.IsDir() && ent.Name() != filesDirName {
			return false, nil
		}
	}
	return true, nil
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

// phasePin says whether a subject call narrows its inner gate to the one phase the
// fixture's family names. Only a mutated tree is graded that narrowly: a baseline pinned
// to a phase prints a fraction of what the empty tree can produce, so an EXPECT the full
// run already emits goes unflagged and every fixture in the group clears the vacuity
// check unguarded.
type phasePin bool

const (
	pinFamilyPhase phasePin = true
	noPhasePin     phasePin = false
)

// subjectCall is the call that grades one tree — a fixture's materialized mutation, or a
// group's empty baseline, which is why fixtureDir and the pin are parameters rather than
// read from the fixture. Both shapes come from here so a group's baseline can never drift
// from the runs it is the yardstick for.
//
// A fixture bound to a contract package invokes that package's compiled test binary over
// the tree; every other fixture spawns the inner gate around it.
func subjectCall(fixture selected, subjectRoot, fixtureDir string, run sweepRun, pin phasePin) RunCall {
	if fixture.pkg != "" {
		return RunCall{
			Kind:       RunBite,
			Cwd:        contractPackageDir(run.root, fixture.pkg),
			FixtureDir: fixtureDir,
			Package:    fixture.pkg,
			Binary:     run.binaries[fixture.pkg],
			Env:        biteEnv(run.base, subjectRoot),
		}
	}
	env := scopedEnv(run.gateEnv, fixture)
	if phase := FixturePhase(fixture.family); pin == pinFamilyPhase && phase != "" {
		env = append(env, PhaseEnv+"="+phase)
	}
	return RunCall{Cwd: subjectRoot, Gate: run.gate, FixtureDir: fixtureDir, Env: env}
}

// scopedEnv is env plus the conformance check the fixture's family binds. A run without
// one carries no variable at all rather than an empty one, which names nothing and reds
// the inner gate. The copy is what keeps concurrent runs from appending into one shared
// backing array.
func scopedEnv(env []string, fixture selected) []string {
	out := append([]string(nil), env...)
	if fixture.scope != "" {
		out = append(out, registry.ConformanceCheckEnv+"="+fixture.scope)
	}
	return out
}

// biteEnv is the environment a compiled contract test binary is invoked under: the tree it
// grades, and nothing else the gate would have read. No inner-gate marker, no phase pin,
// and no conformance scope — there is no gate in this run to read any of them, and a
// binary carrying them claims to be a nested gate to whatever reads them next.
//
// The width pin stays, because the sweep's worker budget still divides the machine by the
// inner width: unpinned binaries running at full width in every worker would oversubscribe
// the machine exactly as the nested gates did.
func biteEnv(base []string, subjectRoot string) []string {
	out := append([]string(nil), base...)
	return append(out, SubjectRootEnv+"="+subjectRoot, innerWidthPin())
}

// innerWidthPin holds the budget arithmetic in one place: the divisor fixtureWorkers uses
// and the width each run is pinned to have to be the same number, or the sweep sizes its
// worker pool against a width nothing enforces.
func innerWidthPin() string {
	return fmt.Sprintf("GOMAXPROCS=%d", bounds.CanaryInnerWidth)
}

// sweepEnvKeys are the variables the sweep decides the value of for every call it makes.
// The whole list is stripped from the inherited environment before any of it is set again:
// Go's exec environment has no duplicate-key precedence, so a set without its matching
// strip hands an ambient export control of what the sweep grades, and a strip without its
// matching set leaves a run reading a value the sweep never chose.
var sweepEnvKeys = []string{
	"BENCH_KIT",
	"BENCH_WRAPPER",
	"BENCH_CANARY_INNER",
	PhaseEnv,
	registry.ConformanceTierEnv,
	registry.ConformanceCheckEnv,
	"GOMAXPROCS",
	SubjectRootEnv,
}

// sweepEnv is the inherited environment with every sweep-controlled variable removed. It
// is the base each call kind sets its own variables onto.
func sweepEnv() []string {
	env := make([]string, 0, len(os.Environ()))
	for _, kv := range os.Environ() {
		if slices.ContainsFunc(sweepEnvKeys, func(key string) bool { return strings.HasPrefix(kv, key+"=") }) {
			continue
		}
		env = append(env, kv)
	}
	return env
}

// gateEnv is the environment every inner gate of a tier's sweep runs under. The tier is
// pinned rather than inherited: a fixture grades a check its own tier runs, so an inner
// gate at any other tier skips that check and the fixture reports "did not bite" forever.
// The phase and the check scope vary per run, so they are left to the caller that knows
// the run's fixture and group.
func gateEnv(base []string, tier registry.Tier) []string {
	out := append([]string(nil), base...)
	return append(out, "BENCH_CANARY_INNER=1", registry.ConformanceTierEnv+"="+string(tier), innerWidthPin())
}

// runnerCommand is what a call of each kind actually executes. It is separate from the
// runner that runs it so the dispatch can be graded without spawning anything: every other
// assertion about the sweep runs through an injected fake, so a change that emitted the
// right call metadata while still spawning a gate would satisfy all of them.
func runnerCommand(call RunCall) *exec.Cmd {
	switch call.Kind {
	case RunCompile:
		// -C rather than a working directory, so the package argument resolves against
		// the swept module the same way the structural refusal resolves the binding.
		cmd := exec.Command("go", "-C", call.Cwd, "test", "-c", "-o", call.Binary, contractPackagePrefix+call.Package)
		cmd.Env = call.Env
		return cmd
	case RunBite:
		cmd := exec.Command(call.Binary)
		cmd.Dir = call.Cwd
		cmd.Env = call.Env
		return cmd
	default:
		cmd := exec.Command("bash", call.Gate)
		cmd.Dir = call.Cwd
		cmd.Env = call.Env
		return cmd
	}
}

// contractPackagePrefix is the import path a bound package's slash path hangs off, which
// is the same directory contractPackageDir resolves to under the swept root.
const contractPackagePrefix = "./internal/contract/"

func defaultRunner(call RunCall) RunResult {
	cmd := runnerCommand(call)
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
