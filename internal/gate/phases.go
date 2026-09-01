package gate

// The gate's command surfaces: its phase table and attended gate-pin workflow.
// Phase execution itself lives in runner.go.

import (
	"bufio"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/terminal"
	"github.com/gibbonmi/bench/internal/toon"
)

// Phase is one benchkit gate check. Argv is passed directly to exec.Command; callers
// must not smuggle shell-interpolated command strings through it. Needs names the
// phases that must complete green before this one starts. A need that ends red or
// skipped skips this phase too, because it would grade an artifact its need never
// produced. Dir is the phase's absolute working directory; empty means the runner's
// root.
//
// Anchoring a declared directory is the producer's job. The graded root and the
// runner's root are different trees in a linked repo. Only the producer knows which
// one a path was written against.
type Phase struct {
	Name         string
	Argv         []string
	Env          []string
	Optional     bool
	Needs        []string
	Dir          string
	ExpectedRuns []string
}

var benchkitPhasesForCommand = BenchkitPhases

// BenchkitPhases selects the kit phase table for root. kit owns the Go tests and wrapper
// scripts; the execution owner selects the Bench executable before it runs this table.
func BenchkitPhases(root, kit string) []Phase {
	phases := toolchainPhases(root, kit)
	if SystemSuiteRuns(root, kit) {
		operands, env := SystemSuite(root)
		phases = append(phases, Phase{
			Name: SystemPhaseName,
			Argv: BaseTestArgv("", operands...),
			Env:  env,
			Dir:  kit,
		})
	}
	return append(phases, Phase{
		Name:     "shellcheck",
		Argv:     shellcheckArgv(kit),
		Optional: true,
	})
}

// SystemPhaseName names the system suite on both of its surfaces: the gate's phase
// table and the `bench test --check` grammar. One name serves both because both run the
// same suite.
const SystemPhaseName = "system"

// SystemRootEnv carries the tree under grade to the system suite's owner.
const SystemRootEnv = "BENCH_SYSTEM_ROOT"

// SystemSuite selects the system-suite test operands and its graded-root environment.
// Callers compose this result into their own execution and do not run a second suite.
func SystemSuite(root string) (operands, env []string) {
	return []string{"-tags=system", "./internal/systemtest"}, []string{SystemRootEnv + "=" + root}
}

// SystemSuiteRuns reports whether root is the kit checkout that the system suite may
// grade. Linked repositories do not select or execute the suite.
func SystemSuiteRuns(root, kit string) bool {
	return sameDirectory(root, kit)
}

// toolchainPhases are the Go steps that grade source rather than the built binary.
// Each materializes only when the graded root carries what the step grades. So a
// linked repo never reds on a check the kit wrote for itself. Gofmt, vet, and test
// need a Go module and a toolchain to run it. The kit-only race step additionally
// probes for its sentinel declaration, rather than assuming a colliding package path
// carries the same tests.
func toolchainPhases(root, kit string) []Phase {
	goModule, err := goModuleToolchain(root)
	if !goModule || err != nil {
		return nil
	}
	phases := []Phase{
		{Name: canary.PhaseGofmt, Argv: gatePhaseGoArgv("gofmt", root)},
		// vet needs no gate-go wrapper: it exits nonzero on its own findings and carries
		// no policy the argv would have to encode.
		{Name: canary.PhaseVet, Argv: []string{"go", "-C", root, "vet", trimPath, "./..."}},
		{Name: canary.PhaseTest, Argv: BaseTestArgv("", "./..."), Dir: root, Env: rootConformanceEnv(root, kit)},
	}
	if sameDirectory(root, kit) && declaresRaceTest(root) {
		phases = append(phases, Phase{Name: canary.PhaseRace, Argv: raceDriverArgv(), Dir: root, ExpectedRuns: raceTestNames()})
	}
	return phases
}

func goModuleToolchain(root string) (bool, error) {
	if !isRegularFile(filepath.Join(root, "go.mod")) {
		return false, nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return true, fmt.Errorf("gate: Go executable go is required by %s but is not on PATH", root)
	}
	return true, nil
}

// rootConformanceEnv points the ordinary test phase at the tree under grade. This way
// the conformance registry's checks run inside the oracle, instead of skipping for an
// unset root. There is no separate conformance phase or driver. Go owns package
// scheduling inside the one ordinary test driver. This is what makes that driver
// grade the live tree rather than only the fixtures.
//
// It materializes on the same terms as the race and system phases. It requires the
// graded root to be the kit, and requires that kit to actually declare the entry
// test. So a linked repo, or a root whose package path merely collides, never
// inherits a variable its test binaries cannot honor.
//
// The tier is pinned rather than left to its unset default. The phase inherits the
// operator's environment, and an ambient ship tier would otherwise widen what the dev
// gate grades.
func rootConformanceEnv(root, kit string) []string {
	if !sameDirectory(root, kit) {
		return nil
	}
	if !declaresTest(filepath.Join(root, filepath.FromSlash(conformancePackagePath)), registry.RootConformanceTest) {
		return nil
	}
	return []string{
		registry.ConformanceRootEnv + "=" + root,
		registry.ConformanceTierEnv + "=" + string(registry.Dev),
	}
}

// conformancePackagePath is where the entry test lives inside the kit, relative to the
// graded root.
const conformancePackagePath = "internal/conformance"

func withRunBinary(phases []Phase, selection *runbinary.Selection) []Phase {
	selected := make([]Phase, len(phases))
	for i, phase := range phases {
		selected[i] = phase
		selected[i].Argv = append([]string(nil), phase.Argv...)
		if len(selected[i].Argv) > 0 && selected[i].Argv[0] == runBinaryArgvToken {
			selected[i].Argv[0] = selection.Path
		}
		selected[i].Env = mergeEnv(phase.Env, []string{
			runbinary.Env + "=" + selection.Path,
			"BENCH_KIT=" + selection.SourceRoot,
		})
	}
	return selected
}

// baselinePolicyEnv carries the landing destination's root into the prospective gate.
// The phase schedule comes from that baseline, never from the candidate tree under
// grade, so a candidate manifest cannot omit the checks that grade it.
const baselinePolicyEnv = "BENCH_GATE_BASELINE"

// declaresTest reports whether dir holds a test file declaring a top-level func named
// name. Both filtered steps ask this rather than looking at the directory, because a
// `-run` filter that matches nothing exits 0. The did-it-run guard behind it would
// then red a repo that never asked for the check.
//
// The answer comes from parsed syntax. A byte scan counts the name inside a comment
// or a string literal. That materializes a phase that can only red, precisely the
// harm the probe exists to prevent. Parsing one directory is cheap enough to pay
// before every gate run.
func declaresTest(dir, name string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	fset := token.NewFileSet()
	for _, entry := range entries {
		// Only a regular file is opened. A FIFO named like a test file would block the
		// read forever, wedging the phase table before the gate starts.
		if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.SkipObjectResolution)
		if err != nil {
			continue
		}
		for _, decl := range file.Decls {
			fn, isFunc := decl.(*ast.FuncDecl)
			if isFunc && fn.Recv == nil && fn.Name.Name == name {
				return true
			}
		}
	}
	return false
}

// BaseTestArgv selects a Bench-owned go test command. It adds -C for a named kit and
// always adds -trimpath and -count=1 before the caller's operands. The flags prevent
// checkout paths and cached test results from changing execution evidence.
func BaseTestArgv(kit string, args ...string) []string {
	argv := []string{"go"}
	if kit != "" {
		argv = append(argv, "-C", kit)
	}
	return append(append(argv, "test", trimPath, "-count=1"), args...)
}

// kitRoot resolves the wrapper-selected kit before an exported entry starts its work.
// The resolved value is passed through the run so a later phase cannot change oracles.
func kitRoot(root string) string {
	if kit := os.Getenv("BENCH_KIT"); kit != "" {
		return kit
	}
	return root
}

// PhasesCommand executes the gate phase table for its root argument or the current
// repository root. It selects a manifest table when present, otherwise the kit table.
// It stops its phase process on cancellation and does not record verdict evidence.
// gate-run owns that boundary.
func PhasesCommand(args []string, stdout, stderr io.Writer) int {
	var root string
	if len(args) > 0 && args[0] != "" {
		root = args[0]
	} else {
		r, err := git.Root()
		if err != nil {
			fmt.Fprintln(stderr, toon.NotInRepo())
			return 3
		}
		root = r
	}
	return phasesCommandAtKit(root, kitRoot(root), stdout, stderr)
}

func phasesCommandAtKit(root, kit string, stdout, stderr io.Writer) int {
	return phasesCommandAtKitWithContext(context.Background(), root, kit, stdout, stderr)
}

func phasesCommandAtKitWithContext(base context.Context, root, kit string, stdout, stderr io.Writer) int {
	base, closeLog := inheritGateRunLog(base, stderr)
	defer closeLog()
	selection, err := runbinary.Inherit(kit)
	if err != nil {
		fmt.Fprintf(stderr, "gate: selected Bench executable refused: %v\n", err)
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	return phasesCommandAtKitWithSelection(base, root, kit, selection, stdout, stderr)
}

func phasesCommandAtKitWithSelection(base context.Context, root, kit string, selection *runbinary.Selection, stdout, stderr io.Writer) int {
	phases, err := phaseTable(root, kit)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if decision := Decide(DecisionInput{Subject: "phase-table", Resolution: Resolution{Kind: GateSh}, Phases: decisionPhases(phases)}); !decision.Accepted {
		fmt.Fprintf(stderr, "gate: phase schedule refused: %s\n", decision.Refusal)
		return 1
	}
	phases = withRunBinary(phases, selection)
	ctx, stop := subprocess.NotifyCancel(base)
	defer stop()
	return runPhases(ctx, kit, phases, stdout, stderr)
}

func shellcheckArgv(root string) []string {
	return append([]string{"shellcheck", "-S", "warning"}, shellcheckFiles(root)...)
}

// shellcheckFiles is the exact file list the shellcheck phase lints, apart from the
// invocation itself. shellcheckArgv builds its argv on top of this, and so does the
// shellcheck component's input declaration. Both read the same enumeration, so neither
// can drift from where the linted paths actually start.
func shellcheckFiles(root string) []string {
	files := []string{"bin/bench.sh"}
	for _, dir := range []string{".bench/hooks", ".bench/lib"} {
		files = append(files, shellFilesIn(root, dir)...)
	}
	// This list also carries enforcement shell that suffix-scanning misses by
	// extension or location. The extensionless shift adapters carry no .sh by
	// contract, so this list names them explicitly rather than discovering them by
	// suffix. It also names the gate entry script and the embedded pre-push hook
	// asset. A missing file here is a conformance concern, not a shellcheck one, so
	// shellcheckFiles lints only the files present.
	for _, named := range []string{
		".bench/adapters/claude",
		".bench/adapters/codex",
		".bench/adapters/opencode",
		".bench/gate.sh",
		"internal/adopt/prepush.sh",
	} {
		if isRegularFile(filepath.Join(root, filepath.FromSlash(named))) {
			files = append(files, named)
		}
	}
	return files
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

const pinFileName = "bench-gate-pin"

// PinCommand records HEAD's committed .bench tree for the managed pre-push hook. It
// refuses non-interactive input and requires the caller to confirm before it writes.
func PinCommand(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return pinCommand(args, stdin, stdout, stderr, terminal.IsTerminal)
}

func pinCommand(args []string, stdin io.Reader, stdout, stderr io.Writer, isTerminal func(io.Reader) bool) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "usage: bench gate pin")
		return 2
	}
	if !isTerminal(stdin) {
		fmt.Fprintln(stderr, "error: bench gate pin requires an interactive TTY")
		return 1
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	if !showPinReview(root, stdout, stderr) {
		return 1
	}
	fmt.Fprint(stdout, "Type 'pin .bench' to update the gate pin: ")
	reader := bufio.NewReader(stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Fprintln(stderr, "error: could not read confirmation")
		return 1
	}
	if strings.TrimSpace(line) != "pin .bench" {
		fmt.Fprintln(stderr, "bench gate pin declined; no pin written")
		return 1
	}
	if err := writePinFromHead(root); err != nil {
		fmt.Fprintln(stderr, "error: "+err.Error())
		return 1
	}
	fmt.Fprintln(stdout, "bench gate pin updated")
	return 0
}

func showPinReview(root string, stdout, stderr io.Writer) bool {
	if dirtyBench(root) {
		fmt.Fprintln(stderr, "warning: .bench has uncommitted changes; pinning HEAD's committed .bench tree")
	}
	commit := existingPinnedCommit(root)
	if commit == "" {
		fmt.Fprintln(stdout, "Initial gate pin for HEAD:.bench")
		return true
	}
	fmt.Fprintf(stdout, "Diff since pinned gate commit %s:\n", commit)
	cmd := exec.Command("git", "-C", root, "diff", commit+"..HEAD", "--", ".bench")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run() == nil
}

func dirtyBench(root string) bool {
	out, err := git.Output("-C", root, "status", "--porcelain", "--", ".bench")
	return err == nil && out != ""
}

func existingPinnedCommit(root string) string {
	data, err := os.ReadFile(pinPath(root))
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return ""
	}
	return strings.TrimSpace(lines[1])
}

func writePinFromHead(root string) error {
	tree, err := git.Output("-C", root, "rev-parse", "HEAD:.bench")
	if err != nil || tree == "" {
		return fmt.Errorf("cannot resolve HEAD:.bench")
	}
	commit, err := git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil || commit == "" {
		return fmt.Errorf("cannot resolve HEAD")
	}
	content := fmt.Sprintf("%s\n%s\n%s\n", tree, commit, time.Now().UTC().Format(time.RFC3339))
	return os.WriteFile(pinPath(root), []byte(content), 0o644)
}

func pinPath(root string) string {
	gitdir, err := git.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil || gitdir == "" {
		return filepath.Join(root, ".git", pinFileName)
	}
	return filepath.Join(gitdir, pinFileName)
}

func shellFilesIn(root, relDir string) []string {
	entries, err := os.ReadDir(filepath.Join(root, relDir))
	if err != nil {
		return nil
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sh") {
			continue
		}
		files = append(files, filepath.ToSlash(filepath.Join(relDir, entry.Name())))
	}
	sort.Strings(files)
	return files
}
