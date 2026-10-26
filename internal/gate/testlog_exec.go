package gate

// The trace of the kit's whole test phase: the traced phase, the `bench gate-go exec` shim
// that cmd/go starts in place of each test binary, and the one key that names the Go test
// log of each test package.

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// testLogEnv names the directory that receives one Go test log for each test package of a
// traced test phase.
const testLogEnv = "BENCH_GATE_TESTLOG"

// gateGoExecStep is the gate-go step of the trace shim.
const gateGoExecStep = "exec"

// traceExecOperand is the `-exec` flag of the traced test phase. It stays one argv element,
// because cmd/go splits the flag's value into the program and its arguments.
var traceExecOperand = "-exec=" + strings.Join(gatePhaseGoArgv(gateGoExecStep), " ")

// wholeTestPackages is the package pattern of the whole test phase.
const wholeTestPackages = "./..."

// wholeTestPhase is the test phase over the whole root. Where RootConformanceRuns holds, the
// root is the kit, and each of its test binaries runs through the trace shim.
func wholeTestPhase(root, kit string) Phase {
	if !RootConformanceRuns(root, kit) {
		return testPhase(root, kit, wholeTestPackages)
	}
	return testPhase(root, kit, traceExecOperand, wholeTestPackages)
}

// withTestLog gives the log directory of this process to the traced test phase alone.
// gateEnv strips the directory from every phase child, so no other phase and no nested gate
// run writes into the logs of this run.
func withTestLog(phases []Phase) []Phase {
	dir := os.Getenv(testLogEnv)
	if dir == "" {
		return phases
	}
	for i, phase := range phases {
		if slices.Contains(phase.Argv, traceExecOperand) {
			phases[i].Env = mergeEnv(phase.Env, []string{testLogEnv + "=" + dir})
		}
	}
	return phases
}

// runBinaryArgv puts path in place of the run-binary token: as the program, and as the first
// word of an `-exec=` flag.
func runBinaryArgv(argv []string, path string) []string {
	selected := append([]string(nil), argv...)
	if len(selected) > 0 && selected[0] == runBinaryArgvToken {
		selected[0] = path
	}
	for i, arg := range selected {
		if strings.HasPrefix(arg, "-exec="+runBinaryArgvToken) {
			selected[i] = strings.Replace(arg, runBinaryArgvToken, goQuoted(path), 1)
		}
	}
	return selected
}

// goQuoted quotes path as one word of a cmd/go flag value. cmd/go reads a single-quoted or a
// double-quoted word whole and knows no escape, so a path that holds both quote marks has no
// spelling.
func goQuoted(path string) string {
	if strings.Contains(path, "'") {
		return `"` + path + `"`
	}
	return "'" + path + "'"
}

// gateGoExec is the trace shim: `bench gate-go exec <binary> <args...>`. With a log
// directory, it creates the log of the package in its working directory first, so a binary
// that exits before its testing package starts still leaves a log. It then replaces its
// image with the binary, with the log flag before the binary's own arguments.
func gateGoExec(args []string, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, gateGoUsage)
		return 2
	}
	binary, operands := args[0], args[1:]
	if dir := os.Getenv(testLogEnv); dir != "" {
		log, err := createTestLog(dir)
		if err != nil {
			fmt.Fprintf(stderr, "bench gate-go exec: %v\n", err)
			return 1
		}
		operands = append([]string{"-test.testlogfile=" + log}, operands...)
	}
	err := subprocess.Replace(binary, operands, os.Environ())
	fmt.Fprintf(stderr, "bench gate-go exec: %v\n", err)
	return 1
}

// createTestLog creates the empty log of the test package in the working directory. The
// package key is relative to the root-conformance root, which the traced phase always
// carries, because the trace applies only where root conformance runs.
func createTestLog(dir string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := testLogFile(dir, os.Getenv(registry.ConformanceRootEnv), cwd)
	if err != nil {
		return "", err
	}
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	return path, file.Close()
}

// testLogFile names, under logDir, the log of the test package in the directory pkg of the
// graded root. Both paths resolve their links first, so a root that a link reaches keys each
// package as its target does.
func testLogFile(logDir, root, pkg string) (string, error) {
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	resolvedPkg, err := filepath.EvalSymlinks(pkg)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(resolvedRoot, resolvedPkg)
	if err != nil {
		return "", err
	}
	return filepath.Join(logDir, testLogKey(filepath.ToSlash(rel))+".log"), nil
}

// testLogKey encodes a slash-separated package directory as one file name. The escape covers
// the slash and the escape mark itself, so two different directories never share a name.
func testLogKey(rel string) string {
	return url.PathEscape(rel)
}
