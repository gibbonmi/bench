package testreport

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var grammar = usage.Grammar{
	Cmd:  "bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name>",
	Help: "usage: bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name>",
	Flags: []usage.Flag{
		{Name: "--full"},
		{Name: "--package", HasValue: true, NoEmptyValue: true},
		{Name: "--run", HasValue: true, NoEmptyValue: true},
		{Name: "--changed"},
		{Name: "--base", HasValue: true, NoEmptyValue: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true},
		{Name: "--check", HasValue: true, NoEmptyValue: true},
	},
	MaxArgs: 1,
}

var selectRunBinary = runbinary.ReuseOrOwn

const goChildGroupCancelled = "child process group cancelled"

type focusedRequest struct {
	packageExpr string
	packages    []string
	full        bool
	run         string
	changed     bool
	base        string
	sourceTip   string
	check       string
}

// Command runs Go from root and renders one stable row for each observed result.
func Command(root string, args []string) (string, int) {
	request, line, code := parseFocusedRequest(root, args)
	if line != "" {
		return line + "\n", code
	}
	return runFocusedRequest(root, request)
}

func parseFocusedRequest(root string, args []string) (focusedRequest, string, int) {
	parsed, line, code := usage.Parse(testGrammar(), args)
	if line != "" {
		return focusedRequest{}, line, code
	}
	_, changed := parsed.Flags["--changed"]
	_, explicit := parsed.Flags["--package"]
	check, hasCheck := parsed.Flags["--check"]
	base, hasBase := parsed.Flags["--base"]
	sourceTip, hasSourceTip := parsed.Flags["--source-tip"]
	if len(parsed.Positionals) > 0 {
		if explicit || changed || hasCheck {
			return focusedRequest{}, toon.Usage(grammar.Cmd, parsed.Positionals[0]), 2
		}
	}
	if changed && explicit {
		return focusedRequest{}, toon.Usage(grammar.Cmd, "--changed"), 2
	}
	if hasCheck && (explicit || changed || parsed.Flags["--run"] != "") {
		return focusedRequest{}, toon.Usage(grammar.Cmd, "--check"), 2
	}
	if (hasBase || hasSourceTip) && !changed {
		flag := "--base"
		if hasSourceTip {
			flag = "--source-tip"
		}
		return focusedRequest{}, toon.Usage(grammar.Cmd, flag), 2
	}
	if hasSourceTip && !hasBase {
		return focusedRequest{}, toon.Usage(grammar.Cmd, "--source-tip"), 2
	}
	if hasCheck && check != gate.SystemPhaseName {
		registered, found := registry.Find(check)
		if !found || !registered.RunsAt(registry.Dev) {
			return focusedRequest{}, unknownCheck(check), 2
		}
	}
	packageOperand := strings.Join(parsed.Positionals, "")
	if explicit, ok := parsed.Flags["--package"]; ok {
		packageOperand = explicit
	}
	_, full := parsed.Flags["--full"]
	return focusedRequest{
		packageExpr: packagePattern(root, packageOperand),
		full:        full,
		run:         parsed.Flags["--run"],
		changed:     changed,
		base:        base,
		sourceTip:   sourceTip,
		check:       check,
	}, "", 0
}

func testGrammar() usage.Grammar {
	withInventory := grammar
	withInventory.Help = grammar.Help + "\n" + namedCheckInventory()
	return withInventory
}

func unknownCheck(check string) string {
	return "unknown check: " + check + "\n" + namedCheckInventory()
}

func namedCheckInventory() string {
	checks := append(registry.Names(registry.Dev), gate.SystemPhaseName)
	return "checks:\n  " + strings.Join(checks, "\n  ")
}

func runFocusedRequest(root string, request focusedRequest) (string, int) {
	// The refusal precedes the run-owner selection, which builds a Bench executable with
	// Go. A root the suite may not grade therefore starts no child at all.
	if request.check == gate.SystemPhaseName && !gate.SystemSuiteRuns(root, testBenchSource(root)) {
		return toon.Errorf("system check unavailable", "the system suite grades the kit checkout only") + "\n", 1
	}
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	selection, err := selectRunBinary(ctx, testBenchSource(root))
	if err != nil {
		return toon.Errorf("Bench executable selection failed", err.Error()) + "\n", 1
	}
	defer selection.Close()
	if request.changed {
		subject, kind, hint := diff.ResolveChangedSubject(root, request.base, request.sourceTip)
		if kind != "" {
			return toon.Errorf("changed selection failed", kind+": "+hint) + "\n", 1
		}
		changedEnv, err := selectedRunEnvironment(os.Environ(), selection)
		if err != nil {
			return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
		}
		packages, err := resolveChangedPackagesWithEnvironment(ctx, root, subject.Paths, changedEnv)
		if err != nil {
			return toon.Errorf("changed selection failed", err.Error()) + "\n", 1
		}
		if len(packages) == 0 {
			return emptyReport(request.full)
		}
		request.packages = packages
	}
	if request.check != "" {
		return runNamedCheck(ctx, root, request, selection)
	}
	operands := []string{}
	if request.run != "" {
		operands = append(operands, "-run", request.run)
	}
	if len(request.packages) != 0 {
		operands = append(operands, request.packages...)
	} else {
		operands = append(operands, request.packageExpr)
	}
	env, err := selectedRunEnvironment(os.Environ(), selection)
	if err != nil {
		return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
	}
	return runGoTest(ctx, root, request, focusedTestArgv(operands...), env)
}

func runNamedCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (string, int) {
	if request.check == gate.SystemPhaseName {
		return runSystemCheck(ctx, root, request, selection)
	}
	argv := focusedTestArgv("./internal/conformance", "-run", "^"+registry.RootConformanceTest+"$")
	env, err := conformanceEnvironment(os.Environ(), root, request.check, selection)
	if err != nil {
		return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
	}
	return runGoTest(ctx, root, request, argv, env)
}

// runSystemCheck runs the gate's system phase as a focused run. It reads the phase's
// operands and environment from the gate's producer, and it sets no conformance
// variable, because the system suite is a build-tagged package rather than a
// conformance scope.
func runSystemCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (string, int) {
	operands, suiteEnv := gate.SystemSuite(root)
	env, err := selectedRunEnvironment(os.Environ(), selection)
	if err != nil {
		return toon.Errorf("go test failed to start", err.Error()) + "\n", 1
	}
	return runGoTest(ctx, root, request, focusedTestArgv(operands...), append(env, suiteEnv...))
}

// focusedTestArgv is the `bench test` invocation over one operand list. It takes its
// flag pair from the gate's one test-argv producer, so the focused run shares the gate's
// build cache entries instead of writing a second, path-keyed set.
func focusedTestArgv(operands ...string) []string {
	return gate.BaseTestArgv("", append([]string{"-json"}, operands...)...)
}

func conformanceEnvironment(base []string, root, scope string, selection *runbinary.Selection) ([]string, error) {
	env, err := selectedRunEnvironment(base, selection)
	if err != nil {
		return nil, err
	}
	return append(env,
		registry.ConformanceRootEnv+"="+root,
		registry.ConformanceTierEnv+"="+string(registry.Dev),
		registry.ConformanceScopeEnv+"="+scope,
	), nil
}

func selectedRunEnvironment(base []string, selection *runbinary.Selection) ([]string, error) {
	env, err := testEnvironment(base, selection.Path)
	if err != nil {
		return nil, err
	}
	return append(env, "BENCH_KIT="+selection.SourceRoot), nil
}

// testEnvironment returns the environment the focused run's Go child carries: the
// caller's, without the inherited conformance and capability entries, with the selected
// Bench executable, and with the Bench build cache entry so a focused run warms the
// archives a gate reads.
func testEnvironment(base []string, binary string) ([]string, error) {
	return gocache.Apply(runbinary.WithEnv(withoutConformanceEnvironment(base), binary))
}

func withoutConformanceEnvironment(base []string) []string {
	env := base
	for _, name := range []string{
		registry.ConformanceRootEnv,
		registry.ConformanceTierEnv,
		registry.ConformanceScopeEnv,
		registry.ConformanceChecksEnv,
		registry.ConformanceInheritedEnv,
		capability.LogEnv,
		"BENCH_KIT",
	} {
		env = capability.WithoutEnvironment(env, name)
	}
	return env
}

func runGoTest(ctx context.Context, root string, request focusedRequest, argv, env []string) (string, int) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	cmd.Env = env
	// The focused run holds the shared cache lock for its span, so a clean cannot remove an
	// archive this run is writing or reading. A lock it cannot take never fails the run.
	if holder, err := gocache.Hold(env); err == nil {
		defer holder.Release()
	}
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
	decodedDone := make(chan struct{})
	go func() {
		report, err := decode(stream)
		decodedResult <- decoded{report: report, err: err}
		close(decodedDone)
	}()
	var result decoded
	select {
	case result = <-decodedResult:
	case <-ctx.Done():
		cancelGoProcessGroup(cmd, decodedDone)
		result = <-decodedResult
		_ = cmd.Wait()
		return toon.Errorf("go test interrupted", goChildGroupCancelled) + "\n", 1
	}
	completed := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		completed <- cmd.Wait()
		close(done)
	}()
	var waitErr error
	select {
	case waitErr = <-completed:
	case <-ctx.Done():
		cancelGoProcessGroup(cmd, done)
		<-completed
		return toon.Errorf("go test interrupted", goChildGroupCancelled) + "\n", 1
	}
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
	if request.run != "" && !report.ranTest {
		return toon.Errorf("go test reported no test runs", "run pattern matched no tests") + "\n", 1
	}
	out, renderErr := report.render(request.full)
	if renderErr != nil {
		return toon.RenderError(renderErr) + "\n", 1
	}
	if waitErr != nil {
		return out, 1
	}
	return out, 0
}

func cancelGoProcessGroup(cmd *exec.Cmd, completed <-chan struct{}) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
	select {
	case <-completed:
	case <-time.After(runbinary.BuilderCancelGrace):
	}
	drainGoProcessGroup(cmd.Process.Pid)
}

func drainGoProcessGroup(pgid int) {
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
}

func emptyReport(full bool) (string, int) {
	out, err := (&report{statuses: map[string]string{}, seen: map[string]bool{}, tests: map[string]*testResult{}, packageLog: map[string]string{}}).render(full)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
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
