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
type Phase struct {
	Name     string
	Argv     []string
	Env      []string
	Optional bool
}

type phaseResult struct {
	Name     string
	Code     int
	Skipped  bool
	StartErr error
}

var benchkitPhasesForCommand = BenchkitPhases

// BenchkitPhases is the real four-phase table for the kit gate. root is the tree under
// grade; kit is the checkout that owns the Go tests and wrapper scripts.
func BenchkitPhases(root, kit string) []Phase {
	return []Phase{
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
	}
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
	return argv
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
	results := make([]phaseResult, len(phases))
	var wg sync.WaitGroup
	for idx, phase := range phases {
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

func phaseSummary(result phaseResult) string {
	if result.Skipped {
		return "phase " + result.Name + ": skipped"
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
	if phase.Optional && !commandAvailable(phase.Argv[0]) {
		result.Skipped = true
		return result
	}

	cmd := exec.Command(phase.Argv[0], phase.Argv[1:]...)
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
		if phase.Optional && optionalStartSkip(run.StartErr) {
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

func commandAvailable(name string) bool {
	if strings.ContainsRune(name, os.PathSeparator) {
		info, err := os.Stat(name)
		return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
	}
	_, err := exec.LookPath(name)
	return err == nil
}

func optionalStartSkip(err error) bool {
	return errors.Is(err, os.ErrNotExist) ||
		errors.Is(err, os.ErrPermission) ||
		errors.Is(err, syscall.ENOENT) ||
		errors.Is(err, syscall.EACCES)
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
