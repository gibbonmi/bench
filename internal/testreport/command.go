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
	parsed, line, code := usage.Parse(grammar, args)
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
	if hasCheck {
		check, found := registry.Find(check)
		if !found || !check.RunsAt(registry.Dev) {
			return focusedRequest{}, toon.Usage(grammar.Cmd, "--check"), 2
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

func runFocusedRequest(root string, request focusedRequest) (string, int) {
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	if request.changed {
		subject, kind, hint := diff.ResolveChangedSubject(root, request.base, request.sourceTip)
		if kind != "" {
			return toon.Errorf("changed selection failed", kind+": "+hint) + "\n", 1
		}
		packages, err := resolveChangedPackages(ctx, root, subject.Paths)
		if err != nil {
			return toon.Errorf("changed selection failed", err.Error()) + "\n", 1
		}
		if len(packages) == 0 {
			return emptyReport(request.full)
		}
		request.packages = packages
	}
	selection, err := selectRunBinary(ctx, testBenchSource(root))
	if err != nil {
		return toon.Errorf("Bench executable selection failed", err.Error()) + "\n", 1
	}
	defer selection.Close()
	if request.check != "" {
		return runNamedCheck(ctx, root, request, selection)
	}
	argv := []string{"test", "-json", "-count=1"}
	if request.run != "" {
		argv = append(argv, "-run", request.run)
	}
	if len(request.packages) != 0 {
		argv = append(argv, request.packages...)
	} else {
		argv = append(argv, request.packageExpr)
	}
	return runGoTest(ctx, root, request, argv, selectedRunEnvironment(os.Environ(), selection))
}

func runNamedCheck(ctx context.Context, root string, request focusedRequest, selection *runbinary.Selection) (string, int) {
	argv := []string{"test", "-json", "-count=1", "./internal/conformance", "-run", "^" + registry.RootConformanceTest + "$"}
	return runGoTest(ctx, root, request, argv, conformanceEnvironment(os.Environ(), root, request.check, selection))
}

func conformanceEnvironment(base []string, root, scope string, selection *runbinary.Selection) []string {
	env := selectedRunEnvironment(base, selection)
	return append(env,
		registry.ConformanceRootEnv+"="+root,
		registry.ConformanceTierEnv+"="+string(registry.Dev),
		registry.ConformanceScopeEnv+"="+scope,
	)
}

func selectedRunEnvironment(base []string, selection *runbinary.Selection) []string {
	env := runbinary.WithEnv(withoutConformanceEnvironment(base), selection.Path)
	return append(env, "BENCH_KIT="+selection.SourceRoot)
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
	cmd := exec.Command("go", argv...)
	cmd.Dir = root
	cmd.Env = env
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
