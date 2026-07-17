package gate

// The gate's command surfaces: its phase table and attended gate-pin workflow.
// Phase execution itself lives in runner.go.

import (
	"bufio"
	"context"
	"fmt"
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
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/terminal"
	"github.com/gibbonmi/bench/internal/toon"
)

type phaseMode int

const (
	outerMode phaseMode = iota
	innerMode
)

// Phase is one independent benchkit gate check. Argv is passed directly to
// exec.Command; callers must not smuggle shell-interpolated command strings through it.
// A Serial phase runs to completion before any other phase starts, and its red stops
// the run: the phases behind it would grade the artifact it failed to produce.
type Phase struct {
	Name     string
	Argv     []string
	Env      []string
	Optional bool
	Serial   bool
}

var benchkitPhasesForCommand = BenchkitPhases

// BenchkitPhases is the real phase table for the kit gate. root is the tree under
// grade; kit is the checkout that owns the Go tests and wrapper scripts.
//
// The serial build phase owns the only write to root's dist/bench during a gate run.
// The contract and canary phases exec and copy that binary, and `go build` replaces a
// stale output non-atomically — a concurrent rebuild hands readers a stale or
// partially-written binary. Conformance's own build check goes to a throwaway path
// for the same reason.
func BenchkitPhases(root, kit string) []Phase {
	var phases []Phase
	buildHelper := filepath.Join(root, "scripts", "go-build.sh")
	if isRegularFile(buildHelper) && isRegularFile(filepath.Join(root, "go.mod")) {
		phases = append(phases, Phase{
			Name:   "build",
			Argv:   []string{"bash", buildHelper, root, filepath.Join(root, "dist", "bench")},
			Serial: true,
		})
	}
	return append(phases, []Phase{
		{
			Name: "conformance",
			Argv: goTestArgv(kit, "./internal/conformance", "-run", "^TestRootConformance$"),
			Env:  []string{"BENCH_CONFORMANCE_ROOT=" + root},
		},
		{
			Name: "contract",
			Argv: goTestArgv(kit, "./internal/contract/..."),
			Env:  []string{"BENCH_CONTRACT_ROOT=" + root},
		},
		{
			Name:     "shellcheck",
			Argv:     shellcheckArgv(kit),
			Optional: true,
		},
		{
			Name: "canary",
			Argv: []string{"bash", filepath.Join(kit, "bin", "bench.sh"), "canary", root},
		},
	}...)
}

func goTestArgv(kit string, args ...string) []string {
	argv := []string{"go"}
	if kit != "" {
		argv = append(argv, "-C", kit)
	}
	return append(append(argv, "test", "-count=1"), args...)
}

// PhasesCommand is the `bench gate-phases [root]` plumbing command. It intentionally
// does not record the verdict cache; `gate-run` owns resolve-run-record for the public
// `bench gate` path.
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	return runPhases(ctx, kit, phasesForMode(benchkitPhasesForCommand(root, kit), mode), mode, stdout, stderr)
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
	owner := os.Getenv(canary.PhaseEnv)
	if owner != "conformance" && owner != "contract" {
		owner = ""
	}
	filtered := make([]Phase, 0, len(phases))
	for _, phase := range phases {
		if phase.Name == "canary" {
			continue
		}
		if owner != "" && phase.Name != owner {
			continue
		}
		filtered = append(filtered, phase)
	}
	return filtered
}
