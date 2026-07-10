package gate

// The gate's phase table and the PhasesCommand surface: which checks exist, their
// argv/env/serial shape, and mode filtering. The execution engine that runs them
// lives in runner.go.

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/git"
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
			Argv: []string{"go", "test", "-count=1", "./internal/conformance", "-run", "^TestRootConformance$"},
			Env:  []string{"BENCH_CONFORMANCE_ROOT=" + root},
		},
		{
			Name: "contract",
			Argv: []string{"go", "test", "-count=1", "./internal/contract/..."},
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
