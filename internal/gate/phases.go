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
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/terminal"
	"github.com/gibbonmi/bench/internal/toon"
)

// Phase is one benchkit gate check. Argv is passed directly to exec.Command; callers
// must not smuggle shell-interpolated command strings through it. Needs names the
// phases that must complete green before this one starts — a need that ends red or
// skipped skips this phase too, because it would grade an artifact its need never
// produced. Dir is the phase's absolute working directory; empty means the runner's
// root. Anchoring a declared directory is the producer's job — the graded root and the
// runner's root are different trees in a linked repo, and only the producer knows which
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

// BenchkitPhases is the real phase table for the kit gate. root is the tree under
// grade; kit is the checkout that owns the Go tests and wrapper scripts. The run owner
// selects the Bench executable before this table is constructed, so phases consume it
// and no phase authors another Bench binary.
func BenchkitPhases(root, kit string) []Phase {
	phases := toolchainPhases(root, kit)
	if sameDirectory(root, kit) {
		phases = append(phases, Phase{
			Name: "system",
			Argv: goTestArgv("", "-tags=system", "./internal/systemtest"),
			Env:  []string{"BENCH_SYSTEM_ROOT=" + root},
			Dir:  kit,
		})
	}
	return append(phases, Phase{
		Name:     "shellcheck",
		Argv:     shellcheckArgv(kit),
		Optional: true,
	})
}

// toolchainPhases are the Go steps that grade source rather than the built binary. Each
// materializes only when the graded root carries what the step grades, so a linked repo
// never reds on a check the kit wrote for itself. Gofmt, vet, and test need a Go module
// and a toolchain to run it; the kit-only race step additionally probes for its sentinel
// declaration rather than assuming a colliding package path carries the same tests.
func toolchainPhases(root, kit string) []Phase {
	if !isRegularFile(filepath.Join(root, "go.mod")) {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return nil
	}
	phases := []Phase{
		{Name: canary.PhaseGofmt, Argv: gatePhaseGoArgv("gofmt", root)},
		// vet needs no gate-go wrapper: it exits nonzero on its own findings and carries
		// no policy the argv would have to encode.
		{Name: canary.PhaseVet, Argv: []string{"go", "-C", root, "vet", "./..."}},
		{Name: canary.PhaseTest, Argv: []string{"go", "test", "-count=1", "./..."}, Dir: root},
	}
	if sameDirectory(root, kit) && declaresRaceTest(root) {
		phases = append(phases, Phase{Name: canary.PhaseRace, Argv: raceDriverArgv(), Dir: root, ExpectedRuns: raceTestNames()})
	}
	return phases
}

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

// declaresTest reports whether dir holds a test file declaring a top-level func named
// name. Both filtered steps ask this rather than looking at the directory, because a
// `-run` filter that matches nothing exits 0 and the did-it-run guard behind it would
// then red a repo that never asked for the check. The answer comes from parsed syntax:
// a byte scan counts the name inside a comment or a string literal, which materializes
// a phase that can only red — precisely the harm the probe exists to prevent. Parsing
// one directory is cheap enough to pay before every gate run.
func declaresTest(dir, name string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	fset := token.NewFileSet()
	for _, entry := range entries {
		// Only a regular file is opened: a FIFO named like a test file would block the
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

func goTestArgv(kit string, args ...string) []string {
	argv := []string{"go"}
	if kit != "" {
		argv = append(argv, "-C", kit)
	}
	return append(append(argv, "test", "-count=1"), args...)
}

// kitRoot resolves the wrapper-selected kit before an exported entry starts its work.
// The resolved value is passed through the run so a later phase cannot change oracles.
func kitRoot(root string) string {
	if kit := os.Getenv("BENCH_KIT"); kit != "" {
		return kit
	}
	return root
}

// PhasesCommand is the `bench gate-phases [root]` plumbing command. Its table comes
// from the graded root's phase manifest, or the built-in kit table when the root
// declares none. It intentionally does not record the verdict cache; `gate-run` owns
// resolve-run-record for the public `bench gate` path.
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
// shellcheck component's input declaration — the one enumeration both read, so neither
// can drift from where the linted paths actually start.
func shellcheckFiles(root string) []string {
	files := []string{"bin/bench.sh"}
	for _, dir := range []string{".bench/hooks", ".bench/lib"} {
		files = append(files, shellFilesIn(root, dir)...)
	}
	// Load-bearing enforcement shell that suffix-scanning misses by extension or
	// location: the extensionless shift adapters (named explicitly, not discovered
	// by suffix — they carry no .sh by contract), the gate entry script, and the
	// embedded pre-push hook asset. A missing file here is a conformance concern,
	// not a shellcheck one, so only present files are linted.
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

// PinCommand is the human-attended `bench gate pin` porcelain. It refuses non-TTY
// stdin before doing any write, then records HEAD's committed .bench tree beside the
// gate cache for the managed pre-push hook to verify.
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
