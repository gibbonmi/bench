// Package testreport runs a fresh Go test invocation and renders its observed package terminals.
package testreport

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/testlines"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var grammar = usage.Grammar{
	Cmd:     "bench test [--full] [package]",
	Help:    "usage: bench test [--full] [package]",
	Flags:   []usage.Flag{{Name: "--full"}},
	MaxArgs: 1,
}

var selectRunBinary = runbinary.ReuseOrOwn

// Command runs Go from root and renders one stable row for each observed result.
func Command(root string, args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	packageExpr := packagePattern(root, strings.Join(parsed.Positionals, ""))

	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	selection, err := selectRunBinary(ctx, testBenchSource(root))
	if err != nil {
		return toon.Errorf("Bench executable selection failed", err.Error()) + "\n", 1
	}
	defer selection.Close()
	argv := focusedTestArgv(packageExpr)
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	childEnv, err := testEnvironment(os.Environ(), selection.Path)
	if err != nil {
		return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
	}
	cmd.Env = childEnv
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stream, err := cmd.StdoutPipe()
	if err != nil {
		return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
	}
	if err := cmd.Start(); err != nil {
		return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
	}
	type decoded struct {
		report *report
		err    error
	}
	decodedResult := make(chan decoded, 1)
	go func() {
		report, err := decode(stream)
		decodedResult <- decoded{report: report, err: err}
	}()
	var result decoded
	select {
	case result = <-decodedResult:
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		select {
		case result = <-decodedResult:
		case <-time.After(2 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			result = <-decodedResult
		}
		_ = cmd.Wait()
		return toon.Errorf("go test interrupted", "child process group cancelled") + "\n", 1
	}
	waitErr := cmd.Wait()
	report, decodeErr := result.report, result.err
	if decodeErr != nil {
		return toon.Errorf("go test output malformed", decodeErr.Error()) + "\n", 1
	}
	if waitErr != nil {
		report.markNonzeroFailures()
	}
	if !report.terminal {
		return toon.Errorf("go test reported no packages", "no package terminal event") + "\n", 1
	}
	if incomplete := report.incompletePackages(); len(incomplete) != 0 {
		return toon.Errorf("go test reported incomplete packages", strings.Join(incomplete, ", ")) + "\n", 1
	}
	_, full := parsed.Flags["--full"]
	out, renderErr := report.render(full)
	if renderErr != nil {
		return toon.RenderError(renderErr) + "\n", 1
	}
	if waitErr != nil {
		return out, 1
	}
	return out, 0
}

// testEnvironment returns the environment the focused run's Go child carries: the
// caller's, without the inherited capability log, with the selected Bench executable,
// and with the Bench build cache entry so a focused run warms the archives a gate reads.
func testEnvironment(base []string, binary string) ([]string, error) {
	return gocache.Apply(runbinary.WithEnv(capability.WithoutEnvironment(base, capability.LogEnv), binary))
}

// focusedTestArgv is the `bench test` invocation over one package expression. It takes
// its flag pair from the gate's one test-argv producer, so the focused run shares the
// gate's build cache entries instead of writing a second, path-keyed set.
func focusedTestArgv(packageExpr string) []string {
	return gate.BaseTestArgv("", "-json", packageExpr)
}

// packagePattern maps a bare directory-relative operand to a "./"-prefixed
// pattern so go test does not resolve it against std; anything else passes through.
func packagePattern(root, operand string) string {
	if operand == "" {
		return "./..."
	}
	if strings.HasPrefix(operand, "./") || strings.HasPrefix(operand, "../") || strings.HasPrefix(operand, "/") {
		return operand
	}
	dir := strings.TrimSuffix(operand, "/...")
	info, err := os.Stat(filepath.Join(root, dir))
	if err != nil || !info.IsDir() {
		return operand
	}
	return "./" + operand
}

func testBenchSource(root string) string {
	if kit := os.Getenv("BENCH_KIT"); kit != "" {
		return kit
	}
	return root
}

type event struct {
	Action     string
	Package    string
	ImportPath string
	Test       string
	Output     string
}

type testResult struct {
	packageName string
	test        string
	first       string
	failed      bool
	skipped     bool
	last        string
	structured  string
}

type report struct {
	statuses   map[string]string
	seen       map[string]bool
	tests      map[string]*testResult
	packageLog map[string]string
	terminal   bool
}

func decode(stream io.Reader) (*report, error) {
	report := &report{statuses: map[string]string{}, seen: map[string]bool{}, tests: map[string]*testResult{}, packageLog: map[string]string{}}
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
		if e.Test == "" && strings.Contains(e.Output, "[no test files]") {
			report.statuses[e.Package] = "no-tests"
		}
		if e.Action == "pass" || e.Action == "fail" || e.Action == "skip" {
			if e.Test == "" {
				report.terminal = true
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
				if report.packageLog[e.Package] == "" {
					report.packageLog[e.Package] = line
				}
				continue
			}
			test := report.test(e.Package, e.Test)
			if test.first == "" {
				test.first = line
			}
			test.last = line
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
	packageRows := make([][]string, 0, len(packages))
	for _, pkg := range packages {
		packageRows = append(packageRows, []string{pkg, r.statuses[pkg]})
	}
	failures := r.failures(full)
	packageBlock, err := toon.Table("packages", []string{"package", "status"}, packageRows)
	if err != nil {
		return "", err
	}
	failureBlock, err := toon.Table("failures", []string{"package", "test", "line"}, failures)
	if err != nil {
		return "", err
	}
	skipBlock, err := toon.Table("skips", []string{"package", "test", "reason"}, r.skips(full))
	if err != nil {
		return "", err
	}
	return packageBlock + failureBlock + skipBlock, nil
}

func (r *report) skips(full bool) [][]string {
	rows := make([][]string, 0)
	for _, test := range r.tests {
		if !test.skipped {
			continue
		}
		reason := test.structured
		if reason == "" {
			reason = test.last
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

func (r *report) failures(full bool) [][]string {
	rows := make([][]string, 0)
	for _, test := range r.tests {
		if !test.failed || r.failedDescendant(test) && test.first == "" {
			continue
		}
		line := test.first
		if line == "" {
			line = "no diagnostic emitted"
		}
		rows = append(rows, []string{test.packageName, test.test, diagnosticCell(line, full)})
	}
	for pkg, status := range r.statuses {
		if status == "fail" && r.packageFailure(pkg) {
			line := r.packageLog[pkg]
			if line == "" {
				line = "no diagnostic emitted"
			}
			rows = append(rows, []string{pkg, "", diagnosticCell(line, full)})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0] || rows[i][0] == rows[j][0] && rows[i][1] < rows[j][1]
	})
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
