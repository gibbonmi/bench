// Package registry is the conformance check inventory. It names which checks exist, which
// tier runs each one, and the timing contract for the ship rehearsal.
//
// Package registry imports nothing from internal/conformance. That package's test files
// import internal/canary. Anything internal/canary needs to read must sit below that
// edge, or the conformance test binary refuses to build.
package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Tier names an oracle surface. Dev is `bench gate` and the surfaces that stand in for
// it. Ship is the release rehearsal.
type Tier string

const (
	Dev  Tier = "dev"
	Ship Tier = "ship"
)

// Valid reports whether the source names one gate-owned input derivation.
func (source InputSource) Valid() bool {
	switch source {
	case InputCatchAll, InputGoSource, InputGoAndDataHandling, InputGateEntry, InputOfflineSmoke,
		InputBenchRoutes, InputDecisionDocuments, InputRoadmapBoard, InputBenchkitProfile,
		InputCaptureRetros:
		return true
	default:
		return false
	}
}

// ConformanceRootEnv names the tree the explicit conformance entry point grades. It is
// the only way that entry point learns a root. A surface that means to grade a root must
// set it; unset, the entry point skips.
const ConformanceRootEnv = "BENCH_CONFORMANCE_ROOT"

// ConformanceTierEnv selects the tier the explicit conformance entry point runs at. When
// absent, the surface grades the dev tier. The default stays un-overridable by accident,
// because a writer that means ship must set it explicitly.
const ConformanceTierEnv = "BENCH_CONFORMANCE_TIER"

// ConformanceChecksEnv transports the executed half of the gate-authored ordered
// ordinary-check partition.
const ConformanceChecksEnv = "BENCH_CONFORMANCE_CHECKS"

// ConformanceInheritedEnv transports the half covered by retained evidence.
const ConformanceInheritedEnv = "BENCH_CONFORMANCE_INHERITED"

// TierFor resolves a tier name to a Tier. Anything but the ship name is the dev tier. A
// surface that means ship must say so, and an unset or misspelled value can never quietly
// widen what a run grades.
func TierFor(name string) Tier {
	if name == string(Ship) {
		return Ship
	}
	return Dev
}

// Check is one conformance check's identity. Its position in Checks fixes both the
// execution order and the timing-line index, so the order is part of the contract.
type Check struct {
	Name           string
	Implementation string
	Tier           Tier
	Meta           bool
	Subject        Subject
	Inputs         InputSource
}

// InputSource names the gate-owned derivation that supplies one check's inputs.
type InputSource string

const (
	// InputCatchAll conservatively binds a check to the complete subject tree.
	InputCatchAll InputSource = "catch-all"
	// InputGoSource binds checks that inspect only compiled Go source.
	InputGoSource InputSource = "go-source"
	// InputGoAndDataHandling binds Go enforcement and its DATA_HANDLING advertisement.
	InputGoAndDataHandling InputSource = "go-source+data-handling"
	// InputGateEntry binds the gate entry script inspected by its contract check.
	InputGateEntry InputSource = "gate-entry"
	// InputOfflineSmoke binds the executable offline proof inspected by its check.
	InputOfflineSmoke InputSource = "offline-smoke"
	// InputBenchRoutes binds the wrapper dispatch inspected by its routing check.
	InputBenchRoutes InputSource = "bench-routes"
	// InputDecisionDocuments binds active and compiled decision-map trees.
	InputDecisionDocuments InputSource = "decision-documents"
	// InputRoadmapBoard binds the split roadmap board, the ROADMAP.md index and the roadmap/
	// detail owners its rows name.
	InputRoadmapBoard InputSource = "roadmap-board"
	// InputCaptureRetros binds the pending retrospective capture directory inspected by
	// the improvement-marker check.
	InputCaptureRetros InputSource = "capture-retros"
	// InputBenchkitProfile binds the checks a profile-owned table drives, each reading that
	// table and the subjects it names.
	InputBenchkitProfile InputSource = "benchkit-profile"
)

// Subject names the tree a conformance check grades.
type Subject string

const (
	SubjectRoot           Subject = "root"
	SubjectKitRoot        Subject = "kitRoot"
	SubjectRootAndKitRoot Subject = "root+kitRoot"
)

// Checks is the conformance inventory in execution order.
var Checks = []Check{
	{Name: "conformance-meta", Implementation: "checkConformanceMeta", Tier: Dev, Meta: true, Subject: SubjectKitRoot},
	{Name: "conformance-canary-families", Implementation: "checkConformanceCanaryFamilies", Tier: Dev, Meta: true, Subject: SubjectKitRoot},
	{Name: "kit-compliance", Implementation: "checkKitCompliance", Tier: Dev, Subject: SubjectKitRoot, Inputs: InputCatchAll},
	{Name: "canary-fixture-compliance", Implementation: "checkCanaryFixtureCompliance", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "load-validity-metadata", Implementation: "checkLoadValidityMetadata", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "skills-index-command-adapters", Implementation: "checkSkillsIndexAndCommandAdapters", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "docs-currency-workflow", Implementation: "checkDocsCurrencyAndWorkflow", Tier: Dev, Subject: SubjectRootAndKitRoot, Inputs: InputCatchAll},
	{Name: "gate-entry-contract", Implementation: "checkGateEntryContract", Tier: Dev, Subject: SubjectRoot, Inputs: InputGateEntry},
	{Name: "ordinary-build-census", Implementation: "checkOrdinaryBuildCensus", Tier: Dev, Subject: SubjectKitRoot, Inputs: InputCatchAll},
	{Name: "offline-smoke-proof", Implementation: "checkOfflineSmokeProof", Tier: Dev, Subject: SubjectRoot, Inputs: InputOfflineSmoke},
	{Name: "handoff-shape-single-source", Implementation: "checkHandoffShape", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "harness-prefix-single-source", Implementation: "checkHarnessPrefix", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoSource},
	{Name: "package-shipped-surface", Implementation: "checkPackageShippedSurface", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "line-routing", Implementation: "checkLineRouting", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "package-core-guard", Implementation: "checkPackageCoreAndGuards", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "release-evidence-probe", Implementation: "checkReleaseEvidenceProbe", Tier: Ship, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "bench-sh-routes", Implementation: "checkBenchShRoutes", Tier: Dev, Subject: SubjectRoot, Inputs: InputBenchRoutes},
	{Name: "default-branch-single-source", Implementation: "checkDefaultBranchSingleSource", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoSource},
	{Name: "data-handling-derivation", Implementation: "checkDataHandlingDerivation", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoAndDataHandling},
	{Name: "single-control-escaper", Implementation: "checkSingleControlEscaper", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoSource},
	{Name: "bounds-policy", Implementation: "checkBoundsPolicy", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "marker-wait-deadlines", Implementation: "checkMarkerWaitDeadlines", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoSource},
	{Name: "subcommand-routing", Implementation: "checkSubcommandRouting", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoSource},
	{Name: "axi-query-registry", Implementation: "checkAXIQueryRegistry", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "skip-ownership", Implementation: "checkSkipOwnership", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoSource},
	{Name: "decision-map-integrity", Implementation: "ValidateDecisionMapTree", Tier: Dev, Subject: SubjectRoot, Inputs: InputDecisionDocuments},
	{Name: "injected-port-registry", Implementation: "checkInjectedPortRegistry", Tier: Dev, Subject: SubjectRoot, Inputs: InputGoSource},
	{Name: "guidance-prose-budgets", Implementation: "checkGuidanceProseBudgets", Tier: Dev, Subject: SubjectRoot, Inputs: InputBenchkitProfile},
	{Name: "profile-lane-table", Implementation: "checkProfileLaneTable", Tier: Dev, Subject: SubjectRoot, Inputs: InputBenchkitProfile},
	{Name: "roadmap-detail-integrity", Implementation: "ValidateRoadmapTree", Tier: Dev, Subject: SubjectRoot, Inputs: InputRoadmapBoard},
	{Name: "structure-accept-currency", Implementation: "ValidateAcceptGrants", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "retro-improvement-markers", Implementation: "ValidateImprovementMarkers", Tier: Dev, Subject: SubjectRoot, Inputs: InputCaptureRetros},
	{Name: "row-next-grammar", Implementation: "checkRowNextGrammar", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "prose-mechanics", Implementation: "checkProseMechanics", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
	{Name: "entry-point-parity", Implementation: "checkEntryPointParity", Tier: Dev, Subject: SubjectRoot, Inputs: InputCatchAll},
}

// familyChecks binds each canary conformance family directory to the check whose
// diagnostics its fixtures grade. The binding follows the emitting code, not the
// directory's name. Three doc families share docs-currency-workflow. compliance-hardening
// grades canary-fixture-compliance against the immutable fixture tree, rather than kit-
// compliance against the live kit root.
//
// familyChecks is unexported, because map iteration order is nondeterministic and the
// family list feeds a diagnostic. Families is the ordered way in.
var familyChecks = map[string]string{
	"package-core-guard":            "package-core-guard",
	"line-routing":                  "line-routing",
	"load-validity-metadata":        "load-validity-metadata",
	"skills-index-command-adapters": "skills-index-command-adapters",
	"data-handling-derivation":      "data-handling-derivation",
	"docs-currency-token-diet":      "docs-currency-workflow",
	"workflow-guidance-anchors":     "docs-currency-workflow",
	"coverage-map-validation":       "docs-currency-workflow",
	"compliance-hardening":          "canary-fixture-compliance",
	"decision-map-integrity":        "decision-map-integrity",
	"injected-ports":                "injected-port-registry",
	"guidance-prose-budgets":        "guidance-prose-budgets",
	"roadmap-detail-integrity":      "roadmap-detail-integrity",
	"retro-improvement-markers":     "retro-improvement-markers",
	"row-next-grammar":              "row-next-grammar",
	"prose-mechanics":               "prose-mechanics",
	"entry-point-parity":            "entry-point-parity",
}

// Families lists the family names this table binds, in sorted order. They are the table's
// own keys, not a reading of tests/canary/. A family directory the table omits therefore
// does not appear here. Grading the tree against the table is a caller's job.
func Families() []string {
	names := make([]string, 0, len(familyChecks))
	for family := range familyChecks {
		names = append(names, family)
	}
	sort.Strings(names)
	return names
}

// FamilyCheck names the check a conformance family's fixtures grade. It reports false for
// a family the table does not bind. An unbound family is a caller's error to raise, never
// a silent unscoped run.
func FamilyCheck(family string) (string, bool) {
	check, bound := familyChecks[family]
	return check, bound
}

// CanaryFamilies lists the canary families owned by check in sorted order.
func CanaryFamilies(check string) []string {
	var families []string
	for family, owner := range familyChecks {
		if owner == check {
			families = append(families, family)
		}
	}
	sort.Strings(families)
	return families
}

// Find returns the registry row for a check name.
func Find(name string) (Check, bool) {
	for _, check := range Checks {
		if check.Name == name {
			return check, true
		}
	}
	return Check{}, false
}

// RunsAt reports whether tier executes the check. Ship is a superset of Dev. The release
// rehearsal must reprove everything dev green claims. The stress-tagged cross-compile
// matrix only a release runs sits inside a dev check.
func (c Check) RunsAt(tier Tier) bool { return c.Tier == Dev || tier == Ship }

// Names lists, in execution order, the checks tier runs.
func Names(tier Tier) []string {
	var names []string
	for _, check := range Checks {
		if check.RunsAt(tier) {
			names = append(names, check.Name)
		}
	}
	return names
}

// OrdinaryNames lists, in execution order, the non-meta checks tier runs.
func OrdinaryNames(tier Tier) []string {
	var names []string
	for _, check := range Checks {
		if !check.Meta && check.RunsAt(tier) {
			names = append(names, check.Name)
		}
	}
	return names
}

// CanonicalOrdinarySelection validates names as a duplicate-free subset of tier's
// ordinary checks and returns that subset in registry order.
func CanonicalOrdinarySelection(tier Tier, names []string) ([]string, error) {
	want := make(map[string]bool, len(names))
	for _, name := range names {
		check, found := Find(name)
		if !found || check.Meta || !check.RunsAt(tier) || want[name] {
			return nil, fmt.Errorf("invalid ordinary %s check %q", tier, name)
		}
		want[name] = true
	}
	ordered := make([]string, 0, len(names))
	for _, check := range Checks {
		if want[check.Name] {
			ordered = append(ordered, check.Name)
		}
	}
	return ordered, nil
}

// RootConformanceTest is the explicit ship rehearsal entry point.
const RootConformanceTest = "TestRootConformance"

const timingFileName = "bench-conformance-timing"

// TimingPath is the per-check timing file for a graded root, or "" when that root carries
// no git dir. Resolution reads <root>/.git alone. It never ascends and never consults the
// working directory. Canary fixtures grade temp roots concurrently while cwd sits in the
// host checkout. An ascending lookup would have every one of those runs truncating the
// host's single file.
func TimingPath(root string) string {
	if root == "" {
		return ""
	}
	dotGit := filepath.Join(root, ".git")
	info, err := os.Stat(dotGit)
	if err != nil {
		return ""
	}
	if info.IsDir() {
		return filepath.Join(dotGit, timingFileName)
	}
	// A worktree or submodule checkout carries .git as a file naming the real dir.
	data, err := os.ReadFile(dotGit)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "gitdir:") {
			continue
		}
		dir := strings.TrimSpace(strings.TrimPrefix(line, "gitdir:"))
		if dir == "" {
			return ""
		}
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(root, dir)
		}
		return filepath.Join(dir, timingFileName)
	}
	return ""
}

// TimingLine is the one timing format: a zero-padded position, the check name verbatim,
// and a duration. The index carries the ordering, so the line sequence stays byte-stable
// while the durations vary. The absence of prose keeps a line from swallowing a canary
// fixture's expectation and rendering it vacuous.
func TimingLine(index int, name string, elapsed time.Duration) string {
	return fmt.Sprintf("%02d %s %s", index, name, elapsed.Round(time.Millisecond))
}

// ReadTimingLines returns the timing lines recorded for a graded root. It returns nothing
// at all when the root has no git dir or no timing file. The print decorates a gate
// verdict and must never be the thing that reds one.
func ReadTimingLines(root string) []string {
	path := TimingPath(root)
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

// TimingWriter appends one line per executed check, numbering them in call order. A root
// with no git dir yields a writer that counts but discards. The driver grades trees that
// are not repositories at all.
type TimingWriter struct {
	path  string
	count int
}

// ClearTiming empties whatever timing lines a root's file already holds. No reader can
// then attribute an earlier run's lines to the current one. Whoever starts a conformance
// run clears at the run boundary. A reader afterward sees that run's lines or none at
// all. A root with no git dir clears nothing and reports no error, because the driver
// grades trees that are not repositories at all.
func ClearTiming(root string) error {
	path := TimingPath(root)
	if path == "" {
		return nil
	}
	return os.WriteFile(path, nil, 0o644)
}

// NewTimingWriter clears the root's timing file so each run stands alone.
func NewTimingWriter(root string) *TimingWriter {
	path := TimingPath(root)
	if path != "" {
		if err := ClearTiming(root); err != nil {
			path = ""
		}
	}
	return &TimingWriter{path: path}
}

// Record writes the timing line for one executed check.
func (w *TimingWriter) Record(name string, elapsed time.Duration) {
	w.count++
	if w.path == "" {
		return
	}
	file, err := os.OpenFile(w.path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()
	fmt.Fprintln(file, TimingLine(w.count, name, elapsed))
}
