// Package registry is the conformance check inventory: which checks exist, which
// tier runs each one, which tests the gate's filtered inner run skips, and the
// contract for the timing file the conformance driver writes and the gate runner
// prints.
//
// It imports nothing from internal/conformance. That package's test files import
// internal/canary, so anything internal/canary needs to read has to sit below that
// edge or the conformance test binary refuses to build.
package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Tier names an oracle surface: Dev is `bench gate` and the surfaces that stand in
// for it, Ship is the release rehearsal.
type Tier string

const (
	Dev  Tier = "dev"
	Ship Tier = "ship"
)

// ConformanceTierEnv selects the tier an inner grading surface runs at: the
// conformance entry point reads it, and so does the gate a canary sweep drives.
// Absent, the surface grades the dev tier — the default stays un-overridable by
// accident because a writer that means ship has to set it explicitly.
const ConformanceTierEnv = "BENCH_CONFORMANCE_TIER"

// ConformanceCheckEnv scopes an inner grading surface to a single check: the
// conformance entry point reads it, and a canary sweep sets it per fixture. Absent,
// the surface runs the whole tier. Any value naming no check the tier runs is a red
// diagnostic rather than a fallback, so a stale or misspelled scope can never read as
// green silence.
const ConformanceCheckEnv = "BENCH_CONFORMANCE_CHECK"

// Check is one conformance check's identity. Its position in Checks fixes both the
// execution order and the timing-line index, so the order is part of the contract.
type Check struct {
	Name string
	Tier Tier
}

// Checks is the conformance inventory in execution order.
var Checks = []Check{
	{Name: "conformance-canary-families", Tier: Dev},
	{Name: "kit-compliance", Tier: Dev},
	{Name: "canary-inner-compliance", Tier: Dev},
	{Name: "load-validity-metadata", Tier: Dev},
	{Name: "skills-index-command-adapters", Tier: Dev},
	{Name: "docs-currency-workflow", Tier: Dev},
	{Name: "line-routing", Tier: Dev},
	{Name: "package-core-guard", Tier: Dev},
	{Name: "release-evidence-probe", Tier: Ship},
	{Name: "bench-sh-routes", Tier: Dev},
	{Name: "default-branch-single-source", Tier: Dev},
	{Name: "data-handling-derivation", Tier: Dev},
	{Name: "single-control-escaper", Tier: Dev},
	{Name: "bounds-policy", Tier: Dev},
	{Name: "marker-wait-deadlines", Tier: Dev},
	{Name: "subcommand-routing", Tier: Dev},
	{Name: "skip-ownership", Tier: Dev},
}

// familyChecks binds each canary conformance family directory to the check whose
// diagnostics its fixtures grade. The binding follows the emitting code, not the
// directory's name: three doc families share docs-currency-workflow, and
// compliance-hardening grades canary-inner-compliance — the check that runs against
// the fixture tree — rather than the similarly named kit-compliance, which grades the
// kit root instead and would scope those fixtures away from their emitter.
//
// It is unexported because map iteration order is nondeterministic and the family
// list feeds a diagnostic; Families is the ordered way in.
var familyChecks = map[string]string{
	"package-core-guard":            "package-core-guard",
	"line-routing":                  "line-routing",
	"load-validity-metadata":        "load-validity-metadata",
	"skills-index-command-adapters": "skills-index-command-adapters",
	"data-handling-derivation":      "data-handling-derivation",
	"docs-currency-token-diet":      "docs-currency-workflow",
	"workflow-guidance-anchors":     "docs-currency-workflow",
	"coverage-map-validation":       "docs-currency-workflow",
	"compliance-hardening":          "canary-inner-compliance",
}

// Families lists the canary conformance family directories in sorted order.
func Families() []string {
	names := make([]string, 0, len(familyChecks))
	for family := range familyChecks {
		names = append(names, family)
	}
	sort.Strings(names)
	return names
}

// FamilyCheck names the check a conformance family's fixtures grade, and reports
// false for a family the table does not bind — an unbound family is a caller's error
// to raise, never a silent unscoped run.
func FamilyCheck(family string) (string, bool) {
	check, bound := familyChecks[family]
	return check, bound
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

// RunsAt reports whether tier executes the check. Ship is a superset of Dev: the
// release rehearsal has to reprove everything dev green claims, and the
// stress-tagged cross-compile matrix only a release runs sits inside a dev check.
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

// InnerSkipTests names the conformance tests the gate's filtered inner run excludes.
// TestRootConformance is the outer run's own entry point, so running it inside the
// run that implements it is the recursion this list exists to make impossible.
var InnerSkipTests = []string{"TestRootConformance"}

// InnerSkipPattern is the `go test -skip` argument built from InnerSkipTests.
func InnerSkipPattern() string {
	return "^(" + strings.Join(InnerSkipTests, "|") + ")$"
}

const timingFileName = "bench-conformance-timing"

// TimingPath is the per-check timing file for a graded root, or "" when that root
// carries no git dir. Resolution reads <root>/.git alone: it never ascends and never
// consults the working directory, because canary fixtures grade temp roots
// concurrently while cwd sits in the host checkout, and an ascending lookup would
// have every one of those runs truncating the host's single file.
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

// TimingLine is the one timing format: a zero-padded position, the check name
// verbatim, and a duration. The index carries the ordering, so the line sequence
// stays byte-stable while the durations vary; the absence of prose keeps a line
// from swallowing a canary fixture's expectation and rendering it vacuous.
func TimingLine(index int, name string, elapsed time.Duration) string {
	return fmt.Sprintf("%02d %s %s", index, name, elapsed.Round(time.Millisecond))
}

// ReadTimingLines returns the timing lines recorded for a graded root, and nothing
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

// TimingWriter appends one line per executed check, numbering them in call order.
// A root with no git dir yields a writer that counts but discards: the driver grades
// trees that are not repositories at all.
type TimingWriter struct {
	path  string
	count int
}

// ClearTiming empties whatever timing lines a root's file already holds, so no
// reader can attribute an earlier run's lines to the current one. Whoever starts a
// conformance run clears at the run boundary; a reader afterwards then sees that
// run's lines or none at all. A root with no git dir clears nothing and reports no
// error: the driver grades trees that are not repositories at all.
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
