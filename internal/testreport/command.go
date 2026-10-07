package testreport

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/diff"
	benchenv "github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// Usage is the one owner of the bench test grammar line.
const Usage = "bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name> | bench test [--full] --check system --run <go-regex> | bench test --check <name> --fixtures | bench test --checks"

var grammar = usage.Grammar{
	Cmd:  Usage,
	Help: "usage: " + Usage,
	Flags: []usage.Flag{
		{Name: "--full"},
		{Name: "--package", HasValue: true, NoEmptyValue: true},
		{Name: "--run", HasValue: true, NoEmptyValue: true},
		{Name: "--changed"},
		{Name: "--base", HasValue: true, NoEmptyValue: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true},
		{Name: "--check", HasValue: true, NoEmptyValue: true},
		{Name: "--fixtures"},
		{Name: "--checks"},
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
	fixtures    bool
	checks      bool
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
	_, fixtures := parsed.Flags["--fixtures"]
	_, checks := parsed.Flags["--checks"]
	if checks && (len(parsed.Flags) != 1 || len(parsed.Positionals) != 0) {
		return focusedRequest{}, toon.Usage(grammar.Cmd, "--checks"), 2
	}
	if fixtures && (!hasCheck || len(parsed.Flags) != 2) {
		return focusedRequest{}, toon.Usage(grammar.Cmd, "--fixtures"), 2
	}
	if len(parsed.Positionals) > 0 {
		if explicit || changed || hasCheck {
			return focusedRequest{}, toon.Usage(grammar.Cmd, parsed.Positionals[0]), 2
		}
	}
	if changed && explicit {
		return focusedRequest{}, toon.Usage(grammar.Cmd, "--changed"), 2
	}
	if hasCheck && (explicit || changed || parsed.Flags["--run"] != "" && check != gate.SystemPhaseName) {
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
	if hasCheck && !isNamedCheck(check) {
		return focusedRequest{}, unknownCheck(check), 2
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
		fixtures:    fixtures,
		checks:      checks,
	}, "", 0
}

func testGrammar() usage.Grammar {
	withInventory := grammar
	withInventory.Help = grammar.Help + "\nnotes:\n  " + PackageExpressionNote + "\n" + namedCheckInventory()
	return withInventory
}

func runFocusedRequest(root string, request focusedRequest) (Outcome, string, int) {
	if request.fixtures {
		return checkFixtures(root, request.check)
	}
	if request.checks {
		return checksInventory(root)
	}
	// The refusal precedes the run-owner selection, which builds a Bench executable with
	// Go. A root the suite may not grade therefore starts no child at all.
	if request.check == gate.SystemPhaseName && !gate.SystemSuiteRuns(root, testBenchSource(root)) {
		return refusedOutcome(toon.Errorf("system check unavailable", "the system suite grades the kit checkout only")+"\n", 1)
	}
	if request.check == proseCheckName {
		return runProseCheck(root, request.full)
	}
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	selection, err := selectRunBinary(ctx, testBenchSource(root))
	if err != nil {
		return refusedOutcome(toon.Errorf("Bench executable selection failed", err.Error())+"\n", 1)
	}
	defer selection.Close()
	if request.changed {
		subject, kind, hint := diff.ResolveChangedSubject(root, request.base, request.sourceTip)
		if kind != "" {
			return refusedOutcome(toon.Errorf("changed selection failed", kind+": "+hint)+"\n", 1)
		}
		changedEnv, err := selectedRunEnvironment(os.Environ(), selection)
		if err != nil {
			return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
		}
		packages, err := resolveChangedPackagesWithEnvironment(ctx, root, subject.Paths, changedEnv)
		if err != nil {
			return refusedOutcome(toon.Errorf("changed selection failed", err.Error())+"\n", 1)
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
		return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
	}
	return runGoTest(ctx, root, request, focusedTestArgv(operands...), env)
}

// focusedTestArgv is the `bench test` invocation over one operand list. It takes its
// flag pair from the gate's one test-argv producer, so the focused run shares the gate's
// build cache entries instead of writing a second, path-keyed set.
func focusedTestArgv(operands ...string) []string {
	return gate.BaseTestArgv("", append([]string{"-json"}, operands...)...)
}

func runGoTest(ctx context.Context, root string, request focusedRequest, argv, env []string) (outcome Outcome, out string, code int) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	// The focused run holds the shared cache lock for its span, so a clean cannot remove an
	// archive this run is writing or reading. A lock it cannot take refuses the run before
	// the Go child starts, because a compile that ran unlocked would write archives a clean
	// is free to remove. An environment that names no home is the one exception, the same
	// one the gate keeps: it derives no Bench cache for a clean to reach.
	holder, err := gocache.Hold(env)
	if err != nil && gocache.Declared(env) {
		return refusedOutcome(gocache.Refusal(env, err)+"\n", 1)
	}
	defer holder.Release()
	if gate.SystemSuiteRuns(root, environmentValue(env, "BENCH_KIT")) {
		run, err := benchenv.OpenKitTestRun(env)
		if err != nil {
			return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
		}
		defer func() {
			if err := run.Close(); err != nil {
				outcome, out, code = refusedOutcome(toon.Errorf("kit test cleanup failed", err.Error())+"\n", 1)
			}
		}()
		env = append(env, run.Entries()...)
	}
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stream, err := cmd.StdoutPipe()
	if err != nil {
		return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
	}
	if err := cmd.Start(); err != nil {
		return refusedOutcome(toon.Errorf("go test failed to start", err.Error())+"\n", 1)
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
		return interruptedOutcome(toon.Errorf("go test interrupted", goChildGroupCancelled)+"\n", 1)
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
		return interruptedOutcome(toon.Errorf("go test interrupted", goChildGroupCancelled)+"\n", 1)
	}
	report, decodeErr := result.report, result.err
	if decodeErr != nil {
		return refusedOutcome(toon.Errorf("go test output malformed", decodeErr.Error())+"\n", 1)
	}
	if waitErr != nil {
		report.markNonzeroFailures()
	}
	if !report.terminal {
		return refusedOutcome(toon.Errorf("go test reported no packages", "no package terminal event")+"\n", 1)
	}
	if incomplete := report.incompletePackages(); len(incomplete) != 0 {
		return refusedOutcome(toon.Errorf("go test reported incomplete packages", strings.Join(incomplete, ", "))+"\n", 1)
	}
	outcome = report.outcome()
	if request.run != "" && outcome.Kind == OutcomeNoTestRun {
		return Outcome{Kind: OutcomeNoTestRun}, toon.Errorf("go test reported no test runs", "run pattern matched no tests") + "\n", 1
	}
	out, renderErr := report.render(request.full)
	if renderErr != nil {
		return refusedOutcome(toon.RenderError(renderErr)+"\n", 1)
	}
	if waitErr != nil {
		return outcome, out, 1
	}
	return outcome, out, 0
}

func cancelGoProcessGroup(cmd *exec.Cmd, completed <-chan struct{}) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
	select {
	case <-completed:
	case <-time.After(bounds.FixedWindow(runbinary.BuilderCancelGrace)):
	}
	drainGoProcessGroup(cmd.Process.Pid)
}

func drainGoProcessGroup(pgid int) {
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
}

func emptyReport(full bool) (Outcome, string, int) {
	empty := newReport()
	out, err := empty.render(full)
	if err != nil {
		return refusedOutcome(toon.RenderError(err)+"\n", 1)
	}
	return empty.outcome(), out, 0
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
