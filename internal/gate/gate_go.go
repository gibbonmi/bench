package gate

// The four toolchain steps whose argv is policy rather than a plain command, behind
// one plumbing subcommand so a phase manifest stays declarative. Each of the four
// exists here because the bare command cannot red on its own: `gofmt -l` exits 0
// while naming the files it rejects, the test step's package set is exclusion policy,
// a `-run` filter that matches nothing exits 0, and the filtered conformance suite
// carries a skip pattern whose single source is the registry.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// cleanupRaceTest is the one test the race step runs: the concurrent-cleanup
// transaction is the only place a data race would be silent under the ordinary suite.
const cleanupRaceTest = "TestConcurrentCleanupRecordsOneTransaction"

const gateGoUsage = "usage: bench gate-go <gofmt|test|race|conformance-suite> [root]"

// GateGoCommand is the `bench gate-go <step> [root]` plumbing command. Exit 0 is a
// green step, 1 a red one, and 2 a usage error — an unrecognized step is never a
// silent success, because a typo in a phase manifest would otherwise grade nothing
// while reporting green.
func GateGoCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || len(args) > 2 {
		fmt.Fprintln(stderr, gateGoUsage)
		return 2
	}
	step := args[0]
	var root string
	if len(args) == 2 && args[1] != "" {
		root = args[1]
	} else {
		r, err := git.Root()
		if err != nil {
			fmt.Fprintln(stderr, toon.NotInRepo())
			return 3
		}
		root = r
	}
	switch step {
	case "gofmt":
		return gofmtStep(root, stdout, stderr)
	case "test":
		return coreTestStep(root, stdout, stderr)
	case "race":
		return raceStep(root, stdout, stderr)
	case "conformance-suite":
		return conformanceSuiteStep(root, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "bench gate-go: unknown step %q\n%s\n", step, gateGoUsage)
		return 2
	}
}

// GateGoArgv is how any caller invokes one of these steps — a phase in the table, a
// step in the release step list. It runs through `go run` rather than dist/bench
// because these steps declare no dependency on the build phase: the build phase owns
// the only write to dist/bench and `go build` replaces it non-atomically, so a reader
// that overlaps it can exec a partially written binary. The Go build cache backs
// `go run`, so the compile is paid once. An empty kit leaves the working directory to
// the caller.
func GateGoArgv(kit, step, root string) []string {
	argv := []string{"go"}
	if kit != "" {
		argv = append(argv, "-C", kit)
	}
	return append(argv, "run", "./cmd/bench", "gate-go", step, root)
}

// CoreTestPackages enumerates the packages the core `go test` step runs at tier: the
// module's own package list, less the exclusions the registry owns. It returns the
// `go list` output alongside, so a caller that fails can attribute the failure to the
// enumeration rather than to the test run.
func CoreTestPackages(root string, tier registry.Tier) ([]string, string, error) {
	listed, combined, err := stepOutput(root, "go", "list", "./...")
	if err != nil {
		return nil, combined, err
	}
	var packages []string
	for _, pkg := range strings.Fields(listed) {
		if registry.IsExcludedTestPackage(pkg, tier) {
			continue
		}
		packages = append(packages, pkg)
	}
	return packages, combined, nil
}

// ConformanceSuiteArgv is the argv for the filtered conformance run, the one that
// keeps the conformance package's own suite in the oracle after the core enumeration
// drops it. It is nil for a root carrying no such package — linked repos and the
// minimal fixtures with a go.mod — where the invocation would report a failure about
// a package that was never there.
func ConformanceSuiteArgv(root string) []string {
	if !isDir(filepath.Join(root, filepath.FromSlash(registry.ConformancePackage))) {
		return nil
	}
	return []string{"go", "test", "./" + registry.ConformancePackage, "-skip", registry.InnerSkipPattern()}
}

// gofmtStep keeps the label the conformance check emitted before the step became a
// phase, so a fixture's expectation and a reader's vocabulary survive the move.
func gofmtStep(root string, stdout, stderr io.Writer) int {
	listed, combined, err := stepOutput(root, "gofmt", "-l", ".")
	if err != nil {
		fmt.Fprintln(stderr, "gofmt failed: "+err.Error())
		fmt.Fprintln(stderr, combined)
		return 1
	}
	files := strings.Fields(listed)
	if len(files) == 0 {
		return 0
	}
	fmt.Fprintln(stdout, "gofmt: unformatted Go files: "+strings.Join(files, " "))
	return 1
}

func coreTestStep(root string, stdout, stderr io.Writer) int {
	packages, output, err := CoreTestPackages(root, registry.TierFor(os.Getenv(registry.ConformanceTierEnv)))
	if err != nil {
		fmt.Fprintln(stderr, "go list failed: "+err.Error())
		fmt.Fprintln(stderr, output)
		return 1
	}
	if len(packages) == 0 {
		return 0
	}
	return runStep(root, append([]string{"go", "test"}, packages...), stdout, stderr)
}

// raceStep reds when the target test did not execute, not only when it failed: the
// `-run` filter exits 0 when it matches nothing, so the `=== RUN` line is the only
// thing separating a pass from a test that was never there. Only stdout is captured,
// where `go test -v` writes that line: tapping both streams would hand exec.Cmd two
// distinct writers over one buffer, and its per-stream copying goroutines would race
// on it.
func raceStep(root string, stdout, stderr io.Writer) int {
	var seen bytes.Buffer
	argv := []string{"go", "test", "-race", "-count=1", "-v", "./internal/worktree", "-run", "^" + cleanupRaceTest + "$"}
	code := runStep(root, argv, io.MultiWriter(stdout, &seen), stderr)
	if code != 0 {
		return code
	}
	if !strings.Contains(seen.String(), "=== RUN   "+cleanupRaceTest) {
		fmt.Fprintln(stderr, "worktree cleanup race test did not run: "+cleanupRaceTest)
		return 1
	}
	return 0
}

func conformanceSuiteStep(root string, stdout, stderr io.Writer) int {
	argv := ConformanceSuiteArgv(root)
	if argv == nil {
		return 0
	}
	return runStep(root, argv, stdout, stderr)
}

// runStep executes one step's argv in root, streaming the tool's own output. The argv
// is an exec array with root already anchored as the working directory, so a root
// path containing a space needs no quoting anywhere.
func runStep(root string, argv []string, stdout, stderr io.Writer) int {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if code := processExitCode(cmd, err); code != 0 {
		return 1
	}
	return 0
}

// stepOutput runs a step whose output is read rather than streamed. It returns stdout
// on its own — the tools here report their findings there, and a stderr line folded
// in would be parsed as one more finding — and both streams together for the
// diagnostic a failing invocation carries.
func stepOutput(root string, argv ...string) (out, combined string, err error) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	combined = stdout.String()
	if combined != "" && stderr.Len() > 0 {
		combined += "\n"
	}
	return stdout.String(), combined + stderr.String(), err
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
