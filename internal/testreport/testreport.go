// Package testreport runs a fresh Go test invocation and renders its observed package terminals.
package testreport

import (
	"encoding/json"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testlines"
	"github.com/gibbonmi/bench/internal/toon"
)

type event struct {
	Action     string
	Package    string
	ImportPath string
	Test       string
	Output     string
	Elapsed    float64
}

type testResult struct {
	packageName string
	test        string
	lines       []string
	failed      bool
	skipped     bool
	structured  string
}

type report struct {
	statuses   map[string]string
	elapsedMS  map[string]int64
	seen       map[string]bool
	tests      map[string]*testResult
	packageLog map[string][]string
	terminal   bool
	// ranTests holds, for each package, one key for each test that emitted a run event,
	// so the packages row count, the count of tests that ran, and the fact that any ran
	// are one observation.
	ranTests map[string]map[string]bool
}

func newReport() *report {
	return &report{statuses: map[string]string{}, elapsedMS: map[string]int64{}, seen: map[string]bool{}, tests: map[string]*testResult{}, packageLog: map[string][]string{}, ranTests: map[string]map[string]bool{}}
}

func decode(stream io.Reader) (*report, error) {
	report := newReport()
	decoder := json.NewDecoder(stream)
	for {
		var e event
		err := decoder.Decode(&e)
		if err == io.EOF {
			return report, nil
		}
		if err != nil {
			return nil, err
		}
		if e.Package == "" && (e.Action == "build-output" || e.Action == "build-fail") {
			e.Package = e.ImportPath
		}
		if e.Package == "" {
			continue
		}
		report.seen[e.Package] = true
		if e.Action == "run" && e.Test != "" {
			if report.ranTests[e.Package] == nil {
				report.ranTests[e.Package] = map[string]bool{}
			}
			report.ranTests[e.Package][e.Test] = true
		}
		if e.Test == "" && strings.Contains(e.Output, "[no test files]") {
			report.statuses[e.Package] = "no-tests"
		}
		if e.Action == "pass" || e.Action == "fail" || e.Action == "skip" {
			if e.Test == "" {
				report.terminal = true
				// The terminal package event is the one carrier of the package's own
				// wall time. Every other event's Elapsed belongs to a single test.
				report.elapsedMS[e.Package] = int64(math.Round(e.Elapsed * 1000))
				if e.Action == "fail" || report.statuses[e.Package] != "no-tests" {
					report.statuses[e.Package] = e.Action
				}
			}
			if e.Test != "" {
				test := report.test(e.Package, e.Test)
				test.failed = e.Action == "fail"
				test.skipped = e.Action == "skip"
			}
		}
		if e.Action != "output" && e.Action != "build-output" {
			continue
		}
		for _, line := range strings.Split(e.Output, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			structured, isStructured := capability.ParseLine(line)
			if isStructured && e.Test != "" {
				reason := string(structured.Kind)
				if structured.Kind == capability.KindCapability {
					reason += ": " + string(structured.Class)
				}
				report.test(e.Package, e.Test).structured = reason + ": " + structured.Reason
				continue
			}
			if testlines.RunnerLine(line) {
				continue
			}
			if e.Test == "" {
				report.packageLog[e.Package] = append(report.packageLog[e.Package], line)
				continue
			}
			test := report.test(e.Package, e.Test)
			test.lines = append(test.lines, line)
		}
	}
}

func (r *report) markNonzeroFailures() {
	for pkg := range r.seen {
		if r.statuses[pkg] == "" {
			r.statuses[pkg] = "fail"
		}
		r.terminal = true
	}
}

func (r *report) incompletePackages() []string {
	packages := make([]string, 0)
	for pkg := range r.seen {
		if r.statuses[pkg] == "" {
			packages = append(packages, pkg)
		}
	}
	sort.Strings(packages)
	return packages
}

func (r *report) test(packageName, name string) *testResult {
	key := packageName + "\x00" + name
	if r.tests[key] == nil {
		r.tests[key] = &testResult{packageName: packageName, test: name}
	}
	return r.tests[key]
}

func (r *report) render(full bool) (string, error) {
	packages := make([]string, 0, len(r.statuses))
	for pkg := range r.statuses {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)
	packageRows := make([][]any, 0, len(packages))
	for _, pkg := range packages {
		// The name and the status are strings the encoder escapes; the last two cells are
		// counts, so they emit bare and stay integers on a round-trip.
		packageRows = append(packageRows, []any{pkg, r.statuses[pkg], r.elapsedMS[pkg], len(r.ranTests[pkg])})
	}
	packageBlock, err := toon.TableTyped("packages", []string{"package", "status", "elapsed_ms", "tests_run"}, packageRows)
	if err != nil {
		return "", err
	}
	failureBlock, err := toon.TableTyped("failures", []string{"package", "test", "line", "lines"}, r.failures(full))
	if err != nil {
		return "", err
	}
	skipBlock, err := toon.Table("skips", []string{"package", "test", "reason"}, r.skips(full))
	if err != nil {
		return "", err
	}
	rootSkips := r.rootConformanceSkips()
	if len(rootSkips) == 0 {
		return packageBlock + failureBlock + skipBlock, nil
	}
	rootBlock, err := toon.Table("root_conformance", []string{"package", "status", "route"}, rootSkips)
	if err != nil {
		return "", err
	}
	return packageBlock + failureBlock + skipBlock + rootBlock, nil
}

func (r *report) rootConformanceSkips() [][]string {
	rows := make([][]string, 0)
	for _, test := range r.tests {
		if test.test == registry.RootConformanceTest && test.skipped {
			rows = append(rows, []string{test.packageName, "skipped", "bench test --check <name>"})
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })
	return rows
}

func (r *report) skips(full bool) [][]string {
	rows := make([][]string, 0)
	for _, test := range r.tests {
		if !test.skipped {
			continue
		}
		reason := test.structured
		if reason == "" && len(test.lines) != 0 {
			reason = test.lines[len(test.lines)-1]
		}
		if reason == "" || goLocationOnly(reason) {
			reason = "reason not emitted"
		}
		rows = append(rows, []string{test.packageName, test.test, diagnosticCell(reason, full)})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0] || rows[i][0] == rows[j][0] && rows[i][1] < rows[j][1]
	})
	return rows
}

func goLocationOnly(reason string) bool {
	for start := 0; start < len(reason); {
		i := strings.IndexByte(reason[start:], ':')
		if i < 0 {
			return false
		}
		i += start
		end := i + 1
		for end < len(reason) && reason[end] >= '0' && reason[end] <= '9' {
			end++
		}
		if end > i+1 && end < len(reason) && reason[end] == ':' {
			return end == len(reason)-1
		}
		start = i + 1
	}
	return false
}

// failure is one failed test, or a package failure with an empty test name, with the
// diagnostic lines it emitted.
type failure struct {
	packageName string
	test        string
	lines       []string
}

// failed returns each failure sorted by package and then by test name.
func (r *report) failed() []failure {
	failed := make([]failure, 0)
	for _, test := range r.tests {
		if !test.failed || r.failedDescendant(test) && len(test.lines) == 0 {
			continue
		}
		failed = append(failed, failure{test.packageName, test.test, test.lines})
	}
	for pkg, status := range r.statuses {
		if status == "fail" && r.packageFailure(pkg) {
			failed = append(failed, failure{pkg, "", r.packageLog[pkg]})
		}
	}
	sort.Slice(failed, func(i, j int) bool {
		return failed[i].packageName < failed[j].packageName || failed[i].packageName == failed[j].packageName && failed[i].test < failed[j].test
	})
	return failed
}

func (r *report) failures(full bool) [][]any {
	rows := make([][]any, 0)
	for _, f := range r.failed() {
		rows = append(rows, f.rows(full)...)
	}
	return rows
}

// rows gives the first diagnostic line, or with full each line in emitted order. Each
// row carries the count of all the lines.
func (f failure) rows(full bool) [][]any {
	if len(f.lines) == 0 {
		return [][]any{{f.packageName, f.test, "no diagnostic emitted", 0}}
	}
	lines := f.lines
	if !full {
		lines = lines[:1]
	}
	rows := make([][]any, 0, len(lines))
	for _, line := range lines {
		rows = append(rows, []any{f.packageName, f.test, diagnosticCell(line, full), len(f.lines)})
	}
	return rows
}

func diagnosticCell(line string, full bool) string {
	if full {
		return sanitize.Controls(line)
	}
	return sanitize.Preview(line)
}

func (r *report) packageFailure(pkg string) bool {
	for _, test := range r.tests {
		if test.packageName == pkg && test.failed {
			return false
		}
	}
	return true
}

func (r *report) failedDescendant(parent *testResult) bool {
	prefix := parent.test + "/"
	for _, test := range r.tests {
		if test.packageName == parent.packageName && test.failed && strings.HasPrefix(test.test, prefix) {
			return true
		}
	}
	return false
}
