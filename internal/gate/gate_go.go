package gate

// The two toolchain steps whose argv is policy rather than a plain command. Gofmt
// belongs to the dev phase table; the filtered core test set belongs only to the
// explicit ship workflow.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/racetests"
	"github.com/gibbonmi/bench/internal/toon"
)

// raceTests is the gate view of the authoritative race-test registry.
var raceTests = racetests.Tests

const gateGoUsage = "usage: bench gate-go <gofmt|test> [root]"
const disableBuildVCS = "-buildvcs=false"

// trimPath is the one owner of the path-trimming flag for every Bench-owned Go argv
// that is not a test: the vet phase and the kit lane's vet and build checks. Without it
// Go hashes the package directory into each compile action ID, so every worktree writes
// its own archive set. It stays separate from the release build flags in
// internal/releaseevidence, because those state the shipped binary's evidence contract
// rather than gate policy.
const trimPath = "-trimpath"
const runBinaryArgvToken = "<bench-run-binary>"

// GateGoCommand is the `bench gate-go <step> [root]` plumbing command. Exit 0 is a
// green step, 1 a red one, 2 a usage error, and 3 an omitted root that no git worktree
// resolves. An unrecognized step is never a silent success, because a typo in a phase
// manifest would otherwise grade nothing while reporting green.
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
	default:
		fmt.Fprintf(stderr, "bench gate-go: unknown step %q\n%s\n", step, gateGoUsage)
		return 2
	}
}

// GateGoArgv is the release-only entry for the core test step. The release workflow
// retains an independent `go run` because it grades a distinct release subject.
func GateGoArgv(kit, step, root string) []string {
	argv := []string{"go"}
	if kit != "" {
		argv = append(argv, "-C", kit)
	}
	return append(argv, "run", disableBuildVCS, "./cmd/bench", "gate-go", step, root)
}

func gatePhaseGoArgv(step, root string) []string {
	return []string{runBinaryArgvToken, "gate-go", step, root}
}

// CoreTestPackages enumerates the packages the core `go test` step runs at tier: the
// module's own package list, less the exclusions the registry owns. It returns the
// `go list` output alongside, so a caller that fails can attribute the failure to the
// enumeration rather than to the test run.
func CoreTestPackages(root string, tier registry.Tier) ([]string, string, error) {
	listed, combined, err := stepOutput(root, "go", "list", disableBuildVCS, "./...")
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

// gofmtStep reds on the files `gofmt -l` names, which `gofmt -l` cannot do on its own:
// the listing is a green exit with output. The verdict goes to stderr, where every
// other red in this file writes. So a caller reading a step's stdout for the tool's own
// output never finds a policy line folded into it.
func gofmtStep(root string, stdout, stderr io.Writer) int {
	listed, combined, err := stepOutput(root, "gofmt", "-l", ".")
	if err != nil {
		fmt.Fprintln(stderr, "gofmt failed: "+err.Error())
		fmt.Fprintln(stderr, combined)
		return 1
	}
	files := listedFiles(listed)
	if len(files) == 0 {
		return 0
	}
	fmt.Fprintln(stderr, "gofmt: unformatted Go files: "+strings.Join(files, " "))
	return 1
}

// listedFiles splits a one-path-per-line tool listing. Splitting on whitespace instead
// would report a path containing a space as two paths that do not exist. Neither would
// name the file a reader has to go fix. A final line with no newline after it still
// names a file, so the split is on the separator rather than on a terminator.
func listedFiles(listed string) []string {
	var files []string
	for _, line := range strings.Split(listed, "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files
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
	return runStep(root, coreTestArgv(packages), stdout, stderr)
}

// coreTestArgv is the release core test step's argv over the enumerated packages.
func coreTestArgv(packages []string) []string {
	return BaseTestArgv("", packages...)
}

func raceTestFilter() string {
	names := make([]string, 0, len(raceTests))
	for _, test := range raceTests {
		names = append(names, regexp.QuoteMeta(test.Name))
	}
	return "^(" + strings.Join(names, "|") + ")$"
}

func raceTestNames() []string {
	names := make([]string, 0, len(raceTests))
	for _, test := range raceTests {
		names = append(names, test.Name)
	}
	return names
}

func raceDriverArgv() []string {
	argv := BaseTestArgv("", "-race", "-v")
	for _, test := range raceTests {
		if !contains(argv, test.PackagePath) {
			argv = append(argv, test.PackagePath)
		}
	}
	return append(argv, "-run", raceTestFilter())
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func declaresRaceTest(root string) bool {
	for _, test := range raceTests {
		if declaresTest(filepath.Join(root, filepath.FromSlash(test.PackagePath)), test.Name) {
			return true
		}
	}
	return false
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
	// A tool that never started (ProcessState nil) wrote to neither stream. Without this
	// check the phase reds while carrying no account of itself, the one shape a reader
	// cannot diagnose from the output.
	if err != nil && cmd.ProcessState == nil {
		fmt.Fprintf(stderr, "%s failed to start: %v\n", argv[0], err)
		return 1
	}
	if code := processExitCode(cmd, err); code != 0 {
		return 1
	}
	return 0
}

// stepOutput runs a step whose output is read rather than streamed. It returns stdout
// on its own, because the tools here report their findings there, and a folded-in
// stderr line would parse as one more finding. It also returns both streams together,
// for the diagnostic a failing invocation carries.
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
