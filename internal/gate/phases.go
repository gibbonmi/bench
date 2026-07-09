package gate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"

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

type phaseResult struct {
	Name     string
	Code     int
	Skipped  bool
	StartErr error
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
	filtered := make([]Phase, 0, len(phases))
	for _, phase := range phases {
		if phase.Name == "canary" {
			continue
		}
		filtered = append(filtered, phase)
	}
	return filtered
}

func runPhases(ctx context.Context, root string, phases []Phase, mode phaseMode, stdout, stderr io.Writer) int {
	if mode == innerMode {
		return runPhasesSequential(ctx, root, phases, stdout, stderr)
	}
	return runPhasesConcurrent(ctx, root, phases, stdout, stderr)
}

func runPhasesSequential(ctx context.Context, root string, phases []Phase, stdout, stderr io.Writer) int {
	red := false
	for _, phase := range phases {
		result := runPhase(ctx, root, phase, stdout, stderr)
		if result.Code == 130 {
			return 130
		}
		if result.Code != 0 {
			red = true
			if phase.Serial {
				break
			}
		}
	}
	if red {
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	fmt.Fprintln(stdout, "gate: green")
	return 0
}

func runPhasesConcurrent(ctx context.Context, root string, phases []Phase, stdout, stderr io.Writer) int {
	var writeMu sync.Mutex
	serial, concurrent := splitSerialPhases(phases)
	for _, phase := range serial {
		out := newPrefixWriter(&writeMu, stdout, phase.Name)
		err := newPrefixWriter(&writeMu, stderr, phase.Name)
		result := runPhase(ctx, root, phase, out, err)
		out.Close()
		err.Close()
		if result.Code == 130 {
			return 130
		}
		if result.Code != 0 {
			fmt.Fprintln(stdout, phaseSummary(result))
			fmt.Fprintln(stderr, "gate: red")
			return 1
		}
		fmt.Fprintln(stdout, phaseSummary(result))
	}

	results := make([]phaseResult, len(concurrent))
	var wg sync.WaitGroup
	for idx, phase := range concurrent {
		idx, phase := idx, phase
		wg.Add(1)
		go func() {
			defer wg.Done()
			out := newPrefixWriter(&writeMu, stdout, phase.Name)
			err := newPrefixWriter(&writeMu, stderr, phase.Name)
			results[idx] = runPhase(ctx, root, phase, out, err)
			out.Close()
			err.Close()
		}()
	}
	wg.Wait()

	cancelled := false
	red := false
	for _, result := range results {
		if result.Code == 130 {
			cancelled = true
		}
		if result.Code != 0 {
			red = true
		}
	}
	if cancelled {
		return 130
	}

	for _, result := range results {
		fmt.Fprintln(stdout, phaseSummary(result))
	}
	if red {
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	fmt.Fprintln(stdout, "gate: green")
	return 0
}

func splitSerialPhases(phases []Phase) (serial, concurrent []Phase) {
	for _, phase := range phases {
		if phase.Serial {
			serial = append(serial, phase)
		} else {
			concurrent = append(concurrent, phase)
		}
	}
	return serial, concurrent
}

func phaseSummary(result phaseResult) string {
	if result.Skipped {
		return "phase " + result.Name + ": skipped (not installed)"
	}
	if result.Code == 0 {
		return "phase " + result.Name + ": green"
	}
	if result.StartErr != nil {
		return fmt.Sprintf("phase %s: red (%v)", result.Name, result.StartErr)
	}
	return fmt.Sprintf("phase %s: red (exit %d)", result.Name, result.Code)
}

func runPhase(ctx context.Context, root string, phase Phase, stdout, stderr io.Writer) phaseResult {
	result := phaseResult{Name: phase.Name}
	if len(phase.Argv) == 0 || phase.Argv[0] == "" {
		result.Code = 1
		result.StartErr = fmt.Errorf("empty argv")
		return result
	}
	argv := phase.Argv
	if phase.Optional {
		resolved, present := resolveOnPath(argv[0])
		if !present {
			// Truly absent from PATH: skip. The summary states "not installed" so the
			// defense's absence is a fact on the record, not silence.
			result.Skipped = true
			return result
		}
		// Exec the resolved path directly so a present-but-unexecutable binary
		// surfaces its real exec error (EACCES) instead of being masked as
		// not-found by exec.LookPath's exec-bit filter. Only an exec-not-found
		// (a missing interpreter, say) still counts as absent below; every other
		// exec failure is red.
		argv = append([]string{resolved}, argv[1:]...)
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = append(gateEnv(), phase.Env...)
	run := runProcessGroupCommand(ctx, cmd)
	if run.StartErr != nil {
		if run.Code == 130 {
			result.Code = 130
			return result
		}
		if phase.Optional && errors.Is(run.StartErr, os.ErrNotExist) {
			result.Skipped = true
			return result
		}
		result.Code = run.Code
		result.StartErr = run.StartErr
		return result
	}
	result.Code = run.Code
	return result
}

// resolveOnPath reports whether a file with the given command name exists on PATH,
// ignoring the executable bit, and returns its resolved path. A bare name is searched
// across PATH entries; a name with a separator is checked directly. Ignoring the exec
// bit is deliberate: a present-but-unexecutable binary must reach exec so its real
// failure surfaces, rather than being silently classified as absent.
func resolveOnPath(name string) (string, bool) {
	if strings.ContainsRune(name, os.PathSeparator) {
		if info, err := os.Stat(name); err == nil && !info.IsDir() {
			return name, true
		}
		return "", false
	}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" {
			dir = "."
		}
		candidate := filepath.Join(dir, name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, true
		}
	}
	return "", false
}

type prefixWriter struct {
	mu     *sync.Mutex
	dst    io.Writer
	prefix string
	buf    []byte
}

func newPrefixWriter(mu *sync.Mutex, dst io.Writer, name string) *prefixWriter {
	return &prefixWriter{mu: mu, dst: dst, prefix: "[" + name + "] "}
}

func (w *prefixWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	for {
		idx := bytesIndexByte(w.buf, '\n')
		if idx < 0 {
			return len(p), nil
		}
		line := append([]byte(nil), w.buf[:idx+1]...)
		w.buf = w.buf[idx+1:]
		if err := w.writeLine(line); err != nil {
			return 0, err
		}
	}
}

func (w *prefixWriter) Close() {
	if len(w.buf) == 0 {
		return
	}
	_ = w.writeLine(w.buf)
	w.buf = nil
}

func (w *prefixWriter) writeLine(line []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	bw := bufio.NewWriter(w.dst)
	if _, err := bw.WriteString(w.prefix); err != nil {
		return err
	}
	if _, err := bw.Write(line); err != nil {
		return err
	}
	return bw.Flush()
}

func bytesIndexByte(b []byte, c byte) int {
	for i, v := range b {
		if v == c {
			return i
		}
	}
	return -1
}
