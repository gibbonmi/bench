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
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/terminal"
	"github.com/gibbonmi/bench/internal/toon"
)

// conformancePhaseName is the phase whose per-check timing the runner prints.
const conformancePhaseName = "conformance"

type phaseMode int

const (
	outerMode phaseMode = iota
	innerMode
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
	Name     string
	Argv     []string
	Env      []string
	Optional bool
	Needs    []string
	Dir      string
}

var benchkitPhasesForCommand = BenchkitPhases

// BenchkitPhases is the real phase table for the kit gate. root is the tree under
// grade; kit is the checkout that owns the Go tests and wrapper scripts.
//
// The build phase owns the only write to root's dist/bench during a gate run, which is
// why every other phase needs it. The contract and canary phases exec and copy that
// binary, and `go build` replaces a stale output non-atomically — a concurrent rebuild
// hands readers a stale or partially-written binary. A root with no Go build surface has
// no build phase and so no edges at all.
func BenchkitPhases(root, kit string) []Phase {
	var phases []Phase
	built := false
	buildHelper := filepath.Join(root, "scripts", "go-build.sh")
	if isRegularFile(buildHelper) && isRegularFile(filepath.Join(root, "go.mod")) {
		phases = append(phases, Phase{
			Name: canary.PhaseBuild,
			Argv: []string{"bash", buildHelper, root, filepath.Join(root, "dist", "bench")},
		})
		built = true
	}
	// Each downstream phase gets its own backing array: one shared slice would let an
	// edit to any phase's Needs rewrite every other phase's edges.
	needsBuild := func() []string {
		if !built {
			return nil
		}
		return []string{canary.PhaseBuild}
	}
	phases = append(phases, toolchainPhases(root, kit)...)
	return append(phases, []Phase{
		{
			Name:  conformancePhaseName,
			Argv:  goTestArgv(kit, "./internal/conformance", "-run", "^TestRootConformance$"),
			Env:   []string{"BENCH_CONFORMANCE_ROOT=" + root},
			Needs: needsBuild(),
		},
		{
			Name:  "contract",
			Argv:  goTestArgv(kit, "./internal/contract/..."),
			Env:   []string{"BENCH_CONTRACT_ROOT=" + root},
			Needs: needsBuild(),
		},
		{
			Name:     "shellcheck",
			Argv:     shellcheckArgv(kit),
			Optional: true,
			Needs:    needsBuild(),
		},
		{
			Name:  "canary",
			Argv:  []string{"bash", filepath.Join(kit, "bin", "bench.sh"), "canary", root},
			Needs: needsBuild(),
		},
	}...)
}

// toolchainPhases are the Go steps that grade source rather than the built binary. Each
// materializes only when the graded root carries what the step grades, so a linked repo
// never reds on a check the kit wrote for itself: gofmt, vet, and test need a Go module
// and a toolchain to run it, and race and the filtered conformance suite each need the
// test they filter for. Both of the latter probe for a declaration rather than a path,
// because a path is a name any repo can collide with while these test names are the
// kit's own.
//
// None of them declares a need on the build phase. That edge exists only to sequence the
// writers and readers of root's dist/bench, and none of these steps execs it — they run
// through `go run`, which the build cache backs. The absent edge is where the split's
// overlap comes from, so restoring it costs the whole win.
func toolchainPhases(root, kit string) []Phase {
	if !isRegularFile(filepath.Join(root, "go.mod")) {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return nil
	}
	phases := []Phase{
		{Name: canary.PhaseGofmt, Argv: GateGoArgv(kit, "gofmt", root)},
		// vet needs no gate-go wrapper: it exits nonzero on its own findings and carries
		// no policy the argv would have to encode.
		{Name: canary.PhaseVet, Argv: []string{"go", "-C", root, "vet", "./..."}},
		{Name: canary.PhaseTest, Argv: GateGoArgv(kit, "test", root)},
	}
	if declaresTest(filepath.Join(root, "internal", "worktree"), cleanupRaceTest) {
		phases = append(phases, Phase{Name: canary.PhaseRace, Argv: GateGoArgv(kit, "race", root)})
	}
	if declaresTest(conformancePackageDir(root), registry.RootConformanceTest) {
		phases = append(phases, Phase{Name: canary.PhaseConformanceSuite, Argv: GateGoArgv(kit, "conformance-suite", root)})
	}
	return phases
}

// conformancePackageDir is the directory the filtered conformance suite grades.
func conformancePackageDir(root string) string {
	return filepath.Join(root, filepath.FromSlash(registry.ConformancePackage))
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
	kit := os.Getenv("BENCH_KIT")
	if kit == "" {
		kit = root
	}
	mode := outerMode
	if os.Getenv("BENCH_CANARY_INNER") == "1" {
		mode = innerMode
	}
	phases, err := phaseTable(root, kit)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	return runPhases(ctx, kit, phasesForMode(phases, mode), mode, stdout, stderr)
}

func shellcheckArgv(root string) []string {
	argv := []string{"shellcheck", "-S", "warning", "bin/bench.sh"}
	for _, dir := range []string{".bench/hooks", ".bench/lib"} {
		argv = append(argv, shellFilesIn(root, dir)...)
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
			argv = append(argv, named)
		}
	}
	return argv
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

func phasesForMode(phases []Phase, mode phaseMode) []Phase {
	if mode != innerMode {
		return phases
	}
	// An inner gate never sweeps fixtures of its own, so the canary phase is gone before
	// the owner is read — which also keeps it from being an owner an ambient export could
	// name, leaving a run of no phases that greens on nothing.
	inner := make([]Phase, 0, len(phases))
	for _, phase := range phases {
		if phase.Name != "canary" {
			inner = append(inner, phase)
		}
	}
	// The resolved table is the only list of phase names here: an owner it carries names
	// a real phase and the fixture runs that one alone, while an owner it does not — a
	// phase this root lacks, or no owner at all — falls back to the full inner gate. A
	// second list of names would silently disagree with the table the run is made of.
	owner := os.Getenv(canary.PhaseEnv)
	if !carriesPhase(inner, owner) {
		return inner
	}
	filtered := make([]Phase, 0, 1)
	for _, phase := range inner {
		if phase.Name == owner {
			filtered = append(filtered, phase)
		}
	}
	return filtered
}

func carriesPhase(phases []Phase, name string) bool {
	for _, phase := range phases {
		if phase.Name == name {
			return true
		}
	}
	return false
}
