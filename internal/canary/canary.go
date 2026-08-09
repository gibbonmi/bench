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
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
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
// it lives, the family that routes its inner gate, the resolved check it grades, and the
// contract package that owns its EXPECT where one does. An empty Family is a legacy flat
// fixture, which belongs to no family and runs the full inner gate.
type Fixture struct {
	Dir     string
	Family  string
	Package string
	Check   string
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
		_, checkName, err := fixtureCheck(fx.dir)
		if err != nil {
			return nil, err
		}
		out[filepath.Base(fx.dir)] = Fixture{
			Dir:     fx.dir,
			Family:  fx.family,
			Package: fx.pkg,
			Check:   fixtureScope(fx.family, checkName),
		}
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

// RunKind names what a sweep runs. A behavior-owned fixture's EXPECT is a contract test's
// own failure message, so that family is graded by compiling the owning package's test
// binary once and invoking it per subject tree; every other family needs a whole gate
// around its mutated tree and spawns one. RunList asks a compiled binary which tests it
// carries, which is how a fixture's named owner is graded against the binary that would
// run it rather than against a second reading of the package source.
type RunKind int

const (
	RunGate RunKind = iota
	RunCompile
	RunBite
	RunList
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
//
// Test is the one contract test or named subtest a bite runs, and is empty for a bite that
// runs its package whole.
type RunCall struct {
	Kind       RunKind
	Cwd        string
	Gate       string
	FixtureDir string
	Package    string
	Binary     string
	Test       string
	Env        []string
}

// RunResult captures the verdict, combined output, and termination of one run.
type RunResult struct {
	ExitCode    int
	Termination subprocess.Termination
	Output      string
}

// Runner runs one call of any kind.
type Runner func(RunCall) RunResult

// grammar is the declared argument shape usage.Parse enforces, so help flags and
// unknown flags answer as invocations instead of resolving as fixture roots.
var grammar = usage.Grammar{
	Cmd:     "bench canary",
	Help:    "usage: bench canary [root]",
	MaxArgs: 1,
}

// Run is the `bench canary [root]` command.
func Run(args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
		} else {
			fmt.Fprintln(stderr, line)
		}
		return code
	}
	root := ""
	if len(parsed.Positionals) == 1 {
		root = parsed.Positionals[0]
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
	canaryDir := filepath.Join(root, "tests", "canary")
	all, err := fixtures(canaryDir)
	if err != nil {
		return err
	}
	fixtures, err := selectTier(all, tier)
	if err != nil {
		return err
	}
	selectedFamilies, scopedFamilies, err := gateSelectedFamilies(tier)
	if err != nil {
		return err
	}
	if scopedFamilies {
		fixtures = slices.DeleteFunc(fixtures, func(fixture selected) bool { return !selectedFamilies[fixture.family] })
	}
	if err := assertContractScopes(root, fixtures); err != nil {
		return err
	}
	if err := assertTestBindings(canaryDir, all); err != nil {
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

	if err := assertDeclaredOwners(fixtures, run, runner); err != nil {
		return err
	}

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
// behavior-owned fixture, the one test inside that package the EXPECT belongs to, and the
// conformance check for a fixture whose family binds one. The family and package are
// carried rather than re-derived from the path at each use, because a nested package puts
// more than one directory between the family and the fixture and a parent-basename read
// would silently name the wrong one.
//
// An empty scope and an empty package together run the whole tier: legacy flat fixtures
// grade surfaces no single check or package owns. Only a group's vacuity baseline runs a
// bound package whole, which is why the narrowing is a parameter of the call rather than a
// state a fixture can be in.
type selected struct {
	dir    string
	family string
	pkg    string
	test   string
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

// testFileName binds a behavior-owned fixture to the one contract test whose failure message
// its EXPECT quotes. Every such fixture carries one: the binding is what narrows its bite to
// the test that owns the expectation, so a fixture without it would pay for its whole package
// while the tree around it claims otherwise.
const testFileName = "TEST"

// readMarker reports the trimmed name a fixture's marker file carries and whether the
// file is there at all. Absence is a separate answer from a present-but-empty file: for
// every marker, deleting the file is how a fixture asks for the default, while a blank
// file is far likelier to be a truncated write than an intent, so only the caller can say
// what each one costs.
//
// The path is typed before it is opened, because the sweep discovers marker files rather
// than being handed them: an open of a FIFO planted at one blocks until a writer that
// never comes, and the whole harness would hang with no diagnostic naming the path that
// stopped it. Anything but a regular file is refused here instead.
func readMarker(fx, marker string) (string, bool, error) {
	path := filepath.Join(fx, marker)
	read := bounds.Classify(path, bounds.ControlRecordLimit)
	switch read.State {
	case bounds.StateAbsent:
		return "", false, nil
	case bounds.StateEmpty, bounds.StateParsed:
		return strings.TrimSpace(string(read.Data)), true, nil
	}
	return "", false, fmt.Errorf("canary fixture marker %s cannot be read: %s", path, read.Reason)
}

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
	name, present, err := readMarker(fx, checkFileName)
	if err != nil {
		return "", "", err
	}
	if !present {
		return registry.Dev, "", nil
	}
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

// selectTier keeps the fixtures tier sweeps, each with its resolved scope and the test
// that owns its EXPECT. It reads the binding rather than grading it: assertTestBindings
// answers for the whole tree, including the fixtures this tier leaves behind, so a defect
// caught here would be reported for a subset of the tree it belongs to.
//
// Membership is tier equality, not the registry's RunsAt superset: the tiers have to
// partition the harness so that every fixture is swept by exactly one of them, and a
// fixture belonging to neither is the unswept rot the canary exists to catch.
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
		// Only a bound package's compiled binary can be narrowed to one test; a fixture
		// graded by a spawned gate has no such notion, so its marker is never read here.
		if fx.pkg != "" {
			test, _, err := readMarker(fx.dir, testFileName)
			if err != nil {
				return nil, err
			}
			fx.test = test
		}
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

// assertTestBindings refuses a sweep while any owner binding in the harness tree is one no
// run would read, ahead of the first compile. Two shapes cost the harness its scoping: a
// behavior-owned fixture naming no owner — the file absent, or present and blank — falls back
// to running every test its package carries, which is the cost the binding exists to remove
// and is paid in silence; and a TEST anywhere no fixture reads it states an owner nothing
// applies, which reads as a binding to the next author and is a lie in the tree. Every defect
// is reported, matching how the sweep reports its fixture failures.
//
// The whole tree is graded rather than the tier's own fixtures: an unread binding misleads
// wherever it sits, including under a family this sweep does not select.
func assertTestBindings(canaryDir string, all []selected) error {
	owned := map[string]bool{}
	var errs []string
	for _, fx := range all {
		if FixturePhase(fx.family) != PhaseContract {
			continue
		}
		owned[fx.dir] = true
		name, present, err := readMarker(fx.dir, testFileName)
		marker := filepath.Join(fx.dir, testFileName)
		switch {
		case err != nil:
			errs = append(errs, err.Error())
		case !present:
			errs = append(errs, fmt.Sprintf("canary fixture '%s' has no %s file at %s, so it names no owning test; name the one contract test whose failure message its EXPECT quotes", filepath.Base(fx.dir), testFileName, marker))
		case name == "":
			errs = append(errs, fmt.Sprintf("canary fixture '%s' has an empty %s file at %s, which names no test; name the one contract test whose failure message its EXPECT quotes", filepath.Base(fx.dir), testFileName, marker))
		}
	}
	unread, err := unreadTestBindings(canaryDir, owned)
	if err != nil {
		return err
	}
	errs = append(errs, unread...)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// unreadTestBindings reports every TEST in the harness tree that sits where no run resolves
// it: outside the behavior-owned family, whose fixtures are the only ones a compiled binary
// can be narrowed for, and above the fixtures inside it, where a family or package directory
// is a step of the walk rather than something graded. Each names the offending path, because
// the file itself is what has to move or go.
//
// A fixture's files/ tree is skipped: it is payload materialized into a graded repo rather
// than harness metadata, and a repo under grade is free to carry a file of that name.
func unreadTestBindings(canaryDir string, owned map[string]bool) ([]string, error) {
	var errs []string
	err := filepath.WalkDir(canaryDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == filesDirName {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() != testFileName || owned[filepath.Dir(path)] {
			return nil
		}
		rel, err := filepath.Rel(canaryDir, path)
		if err != nil {
			return err
		}
		family, _, _ := strings.Cut(filepath.ToSlash(rel), "/")
		if FixturePhase(family) == PhaseContract {
			errs = append(errs, fmt.Sprintf("canary %s file %s sits above the fixtures it would bind, where no run reads it; an owner is named in the fixture directory that grades it", testFileName, path))
			return nil
		}
		errs = append(errs, fmt.Sprintf("canary %s file %s sits outside the behavior-owned family, where no run reads it; only a fixture graded by a compiled contract binary names an owning test", testFileName, path))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return errs, nil
}

// assertDeclaredOwners refuses a sweep while a fixture names an owning test its package's
// compiled binary does not carry, after the compiles and ahead of the first graded run. The
// binary is asked what it holds rather than the package source read a second way, because the
// names a run can be narrowed to are exactly the ones the binary carries. Every defect is
// reported: a renamed owner otherwise surfaces as a did-not-bite per fixture, which is an
// archaeology session ending at the marker this refusal names outright.
func assertDeclaredOwners(fixtures []selected, run sweepRun, runner Runner) error {
	pkgs := contractPackages(fixtures)
	if len(pkgs) == 0 {
		return nil
	}

	carried := make([]map[string]bool, len(pkgs))
	listErrs := make([]string, len(pkgs))
	eachIndex(len(pkgs), func(idx int) {
		pkg := pkgs[idx]
		result := runner(RunCall{
			Kind:    RunList,
			Cwd:     contractPackageDir(run.root, pkg),
			Package: pkg,
			Binary:  run.binaries[pkg],
			Env:     binaryEnv(run.base),
		})
		if result.ExitCode != 0 {
			listErrs[idx] = fmt.Sprintf("canary could not list the tests of contract package %q (exit %d): %s", pkg, result.ExitCode, strings.TrimSpace(result.Output))
			return
		}
		carried[idx] = listedTests(result.Output)
	})

	var errs []string
	for _, msg := range listErrs {
		if msg != "" {
			errs = append(errs, msg)
		}
	}
	// A package whose membership is unknown can grade none of its fixtures' owners, so the
	// list failures answer alone rather than joined by every owner they made unresolvable.
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	membership := make(map[string]map[string]bool, len(pkgs))
	for idx, pkg := range pkgs {
		membership[pkg] = carried[idx]
	}
	for _, fx := range fixtures {
		owner := strings.SplitN(fx.test, "/", 2)[0]
		if fx.pkg == "" || membership[fx.pkg][owner] {
			continue
		}
		errs = append(errs, fmt.Sprintf("canary fixture '%s' names owning test %q, which the compiled binary of contract package %q does not carry", filepath.Base(fx.dir), fx.test, fx.pkg))
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// listedTests is the set of test names a -test.list output carries. The flag lists every
// runnable name the binary holds — benchmarks, fuzz targets, and examples among them — and
// none of those can own a fixture's EXPECT or be reached by a -test.run filter, so only the
// Test names are membership.
func listedTests(output string) map[string]bool {
	out := map[string]bool{}
	for _, line := range strings.Split(output, "\n") {
		if name := strings.TrimSpace(line); strings.HasPrefix(name, "Test") {
			out[name] = true
		}
	}
	return out
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
		outputs[idx] = runner(subjectCall(scopes[idx], dirs[idx], "", run, wideBaseline)).Output
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
	result := runner(subjectCall(fixture, work, fx, run, narrowToFixture))
	owner := fixtureOwnerFor(fixture)
	// An aborted test binary cannot prove a bite, even if it printed EXPECT first.
	if result.Termination.Aborted() {
		return fmt.Sprintf("canary '%s' process abort in %s", name, sanitize.Preview(owner.label))
	}
	if result.ExitCode != 0 {
		if testOwner, marker, aborted := goTestAbort(result.Output); aborted {
			if testOwner == "" {
				testOwner = owner.unknownTest()
			}
			return fmt.Sprintf("canary '%s' inner test abort in %s: %s", name, sanitize.Preview(testOwner), sanitize.Preview(marker))
		}
	}
	if result.ExitCode == 0 || !strings.Contains(result.Output, expect) {
		return fmt.Sprintf("canary '%s' did not bite (want red + %q; got exit %d)", name, expect, result.ExitCode)
	}
	return ""
}

type fixtureOwner struct {
	label string
}

// fixtureOwnerFor is the single package/scope/gate precedence used by both abort classes;
// each class adds only its diagnostic wording.
func fixtureOwnerFor(fixture selected) fixtureOwner {
	switch {
	case fixture.pkg != "":
		return fixtureOwner{label: fmt.Sprintf("contract package %q", fixture.pkg)}
	case fixture.scope != "":
		return fixtureOwner{label: fmt.Sprintf("conformance check %q", fixture.scope)}
	default:
		return fixtureOwner{label: "the inner gate"}
	}
}

func (owner fixtureOwner) unknownTest() string {
	return "unknown test in " + owner.label
}

func goTestAbort(output string) (string, string, bool) {
	owners := map[string]string{}
	for _, rawLine := range strings.Split(output, "\n") {
		phase, line := outerPhaseLine(rawLine)
		if header := goTestFailureHeader(line); header != "" {
			owners[phase] = header
			continue
		}
		if strings.HasPrefix(line, "panic:") || strings.HasPrefix(line, "fatal error:") {
			return owners[phase], line, true
		}
	}
	return "", "", false
}

func outerPhaseLine(line string) (string, string) {
	if !strings.HasPrefix(line, "[") {
		return "", line
	}
	end := strings.Index(line, "] ")
	if end < 1 {
		return "", line
	}
	return line[:end+2], line[end+2:]
}

func goTestFailureHeader(line string) string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "--- FAIL: Test") {
		return ""
	}
	name, _, found := strings.Cut(strings.TrimPrefix(line, "--- FAIL: "), " ")
	if !found {
		return ""
	}
	for _, r := range name {
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return ""
		}
	}
	return name
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

// narrowing says whether a subject call is scoped down to what one fixture needs — the
// phase its family names, and the contract test its EXPECT belongs to. Only a mutated tree
// is graded that narrowly: a baseline scoped by either axis prints a fraction of what the
// empty tree can produce, so an EXPECT the wide run already emits goes unflagged and every
// fixture in the group clears the vacuity check unguarded.
type narrowing bool

const (
	narrowToFixture narrowing = true
	wideBaseline    narrowing = false
)

// subjectCall is the call that grades one tree — a fixture's materialized mutation, or a
// group's empty baseline, which is why fixtureDir and the narrowing are parameters rather
// than read from the fixture. Both shapes come from here so a group's baseline can never
// drift from the runs it is the yardstick for.
//
// A fixture bound to a contract package invokes that package's compiled test binary over
// the tree; every other fixture spawns the inner gate around it.
func subjectCall(fixture selected, subjectRoot, fixtureDir string, run sweepRun, narrow narrowing) RunCall {
	if fixture.pkg != "" {
		call := RunCall{
			Kind:       RunBite,
			Cwd:        contractPackageDir(run.root, fixture.pkg),
			FixtureDir: fixtureDir,
			Package:    fixture.pkg,
			Binary:     run.binaries[fixture.pkg],
			Env:        biteEnv(run.base, subjectRoot),
		}
		if narrow == narrowToFixture {
			call.Test = fixture.test
		}
		return call
	}
	env := scopedEnv(run.gateEnv, fixture)
	if phase := FixturePhase(fixture.family); narrow == narrowToFixture && phase != "" {
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

// binaryEnv is the environment every invocation of a compiled contract test binary runs
// under, and nothing else the gate would have read. No inner-gate marker, no phase pin,
// and no conformance scope — there is no gate in this run to read any of them, and a
// binary carrying them claims to be a nested gate to whatever reads them next.
//
// The width pin stays, because the sweep's worker budget still divides the machine by the
// inner width: unpinned binaries running at full width in every worker would oversubscribe
// the machine exactly as the nested gates did.
func binaryEnv(base []string) []string {
	out := append([]string(nil), base...)
	return append(out, innerWidthPin())
}

// biteEnv is binaryEnv plus the tree the bite grades. A list call takes the environment
// without it: it asks the binary which tests it carries rather than grading anything, so
// naming a tree there would state a subject its answer does not depend on.
func biteEnv(base []string, subjectRoot string) []string {
	return append(binaryEnv(base), SubjectRootEnv+"="+subjectRoot)
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
	runbinary.Env,
	"BENCH_KIT",
	"BENCH_WRAPPER",
	"BENCH_CANARY_INNER",
	PhaseEnv,
	registry.ConformanceTierEnv,
	registry.ConformanceCheckEnv,
	registry.ConformanceChecksEnv,
	registry.ConformanceInheritedEnv,
	FamilySelectionEnv,
	FamilySelectionOwnerEnv,
	FamilySelectionAuthorityEnv,
	"GOMAXPROCS",
	SubjectRootEnv,
}

const (
	// FamilySelectionEnv carries the gate-derived conformance-family subset.
	FamilySelectionEnv = "BENCH_CANARY_FAMILIES"
	// FamilySelectionOwnerEnv distinguishes phase-owned selection from ambient input.
	FamilySelectionOwnerEnv = "BENCH_CANARY_FAMILIES_OWNER"
	// FamilySelectionAuthorityEnv names the inherited descriptor proving the selection
	// came from the gate phase runner rather than from ambient public-canary input.
	FamilySelectionAuthorityEnv   = "BENCH_CANARY_FAMILIES_FD"
	familySelectionAuthorityLimit = 4096
)

func gateSelectedFamilies(tier registry.Tier) (map[string]bool, bool, error) {
	if tier == registry.Ship || os.Getenv(FamilySelectionOwnerEnv) != "gate" {
		return nil, false, nil
	}
	raw := os.Getenv(FamilySelectionEnv)
	if raw == "" {
		return nil, false, nil
	}
	fd, err := strconv.Atoi(os.Getenv(FamilySelectionAuthorityEnv))
	if err != nil || fd < 3 {
		return nil, false, nil
	}
	authority := os.NewFile(uintptr(fd), "bench-canary-family-selection")
	if authority == nil {
		return nil, false, nil
	}
	proof, readErr := readCanaryAuthority(authority, familySelectionAuthorityLimit)
	closeErr := authority.Close()
	if readErr != nil || closeErr != nil || string(proof) != raw {
		return nil, false, fmt.Errorf("invalid gate-owned canary family authority")
	}
	names := strings.Split(raw, ",")
	known := registry.Families()
	want := make(map[string]bool, len(names))
	for _, name := range names {
		if !slices.Contains(known, name) || want[name] {
			return nil, false, fmt.Errorf("invalid gate-owned canary family selection %q", raw)
		}
		want[name] = true
	}
	ordered := make([]string, 0, len(names))
	for _, name := range known {
		if want[name] {
			ordered = append(ordered, name)
		}
	}
	if !slices.Equal(ordered, names) {
		return nil, false, fmt.Errorf("invalid gate-owned canary family selection %q", raw)
	}
	return want, true, nil
}

// sweepEnv is the inherited environment with every sweep-controlled variable removed. It
// is the base each call kind sets its own variables onto.
func sweepEnv() []string {
	selected, inherited := os.LookupEnv(runbinary.Env)
	env := make([]string, 0, len(os.Environ()))
	for _, kv := range os.Environ() {
		if slices.ContainsFunc(sweepEnvKeys, func(key string) bool { return strings.HasPrefix(kv, key+"=") }) {
			continue
		}
		env = append(env, kv)
	}
	if inherited {
		env = runbinary.WithEnv(env, selected)
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
		cmd := exec.Command(call.Binary, biteArgs(call.Test)...)
		cmd.Dir = call.Cwd
		cmd.Env = call.Env
		return cmd
	case RunList:
		// The pattern matches every name the binary carries: a list call asks what the
		// package holds, and narrowing it would answer with a subset of the membership the
		// answer is compared against.
		cmd := exec.Command(call.Binary, "-test.list", ".*")
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

// biteArgs narrows a bite to the one test its fixture names, and runs the package whole
// for a fixture naming none. The name is quoted and anchored because -test.run takes an
// unanchored regexp: raw, a name matches every test it is a substring of, and one carrying
// a metacharacter matches a superset of itself — either way a test the fixture does not
// own can satisfy the bite.
func biteArgs(test string) []string {
	if test == "" {
		return nil
	}
	parts := strings.Split(test, "/")
	for i, part := range parts {
		parts[i] = "^" + regexp.QuoteMeta(part) + "$"
	}
	return []string{"-test.run", strings.Join(parts, "/")}
}

// contractPackagePrefix is the import path a bound package's slash path hangs off, which
// is the same directory contractPackageDir resolves to under the swept root.
const contractPackagePrefix = "./internal/contract/"

func defaultRunner(call RunCall) RunResult {
	cmd := runnerCommand(call)
	r := subprocess.CaptureMerged(cmd)
	return RunResult{ExitCode: r.ExitCode, Termination: r.Termination, Output: r.Stdout}
}

func gitInit(dir string) error {
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	return cmd.Run()
}
