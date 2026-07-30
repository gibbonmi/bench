// Package testreport runs a fresh Go test invocation and renders its observed package terminals.
package testreport

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var grammar = usage.Grammar{
	Cmd:     "bench test [--full] [package]",
	Help:    "usage: bench test [--full] [package]",
	Flags:   []usage.Flag{{Name: "--full"}},
	MaxArgs: 1,
}

// Command runs Go from root and renders one stable row for each observed result.
func Command(root string, args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	packageExpr := "./..."
	if len(parsed.Positionals) == 1 {
		packageExpr = parsed.Positionals[0]
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cmd := exec.Command("go", "test", "-json", "-count=1", packageExpr)
	cmd.Dir = root
	cmd.Env = capability.WithoutEnvironment(os.Environ(), capability.LogEnv)
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
			if runnerLine(line) {
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

func (r *report) test(packageName, name string) *testResult {
	key := packageName + "\x00" + name
	if r.tests[key] == nil {
		r.tests[key] = &testResult{packageName: packageName, test: name}
	}
	return r.tests[key]
}

func runnerLine(line string) bool {
	return strings.HasPrefix(line, "bench-skip ") || strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "=== RUN") || strings.HasPrefix(line, "=== PAUSE") || strings.HasPrefix(line, "=== CONT") || strings.HasPrefix(line, "=== NAME") || strings.HasPrefix(line, "--- FAIL:") || strings.HasPrefix(line, "--- PASS:") || strings.HasPrefix(line, "--- SKIP:") || line == "FAIL" || strings.HasPrefix(line, "FAIL ") || strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") || strings.HasPrefix(line, "? ") || strings.HasPrefix(line, "?\t")
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
		if reason == "" || strings.HasSuffix(reason, ":") {
			reason = "reason not emitted"
		}
		rows = append(rows, []string{test.packageName, test.test, diagnosticCell(reason, full)})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0] || rows[i][0] == rows[j][0] && rows[i][1] < rows[j][1]
	})
	return rows
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
