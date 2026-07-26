package gate

// The execution engine for gate phases: the serial-then-concurrent runners, the
// per-phase process launch with optional-binary resolution, and the prefixed
// output plumbing. The phase table and the PhasesCommand surface live in
// phases.go.

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

const processGroupCancelGrace = 2 * time.Second

type processGroupResult struct {
	Code     int
	StartErr error
}

func runProcessGroupCommand(ctx context.Context, cmd *exec.Cmd) processGroupResult {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	if err := cmd.Start(); err != nil {
		if ctx.Err() != nil {
			return processGroupResult{Code: 130, StartErr: err}
		}
		return processGroupResult{Code: 1, StartErr: err}
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return processGroupResult{Code: processExitCode(cmd, err)}
	case <-ctx.Done():
		if errors.Is(context.Cause(ctx), errGateTimeout) {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
			return processGroupResult{Code: 124}
		}
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		select {
		case <-done:
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		case <-time.After(processGroupCancelGrace):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		return processGroupResult{Code: 130}
	}
}

func processExitCode(cmd *exec.Cmd, err error) int {
	if err == nil {
		return 0
	}
	if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() > 0 {
		return cmd.ProcessState.ExitCode()
	}
	return 1
}

type phaseResult struct {
	Name     string
	Code     int
	Skipped  bool
	StartErr error
}

func runPhases(ctx context.Context, root string, phases []Phase, mode phaseMode, stdout, stderr io.Writer) int {
	if mode == innerMode {
		// The inner gate grades a canary's deliberately mutated tree, so its skips
		// describe that fixture rather than the host. It launches its phases with no
		// skip log at all — gateEnv strips any value the outer run put in the
		// environment — which keeps the fixture's lines out of the outer tally.
		return runPhasesSequential(ctx, root, phases, stdout, stderr)
	}
	skipLog, cleanup, err := newSkipLog()
	if err != nil {
		fmt.Fprintf(stderr, "gate: cannot open the capability skip log: %v\n", err)
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	defer cleanup()
	return runPhasesConcurrent(ctx, root, withSkipLog(phases, skipLog), skipLog, stdout, stderr)
}

func runPhasesSequential(ctx context.Context, root string, phases []Phase, stdout, stderr io.Writer) int {
	// splitSerialPhases is the one source of the serial-first ordering for both
	// runners; without it a serial phase in a non-first table position would
	// fail-fast in outer mode but not here.
	serial, concurrent := splitSerialPhases(phases)
	red := false
	for _, phase := range serial {
		result := runPhase(ctx, root, phase, stdout, stderr)
		if result.Code == 130 {
			return 130
		}
		if result.Code != 0 {
			red = true
			break
		}
	}
	if !red {
		for _, phase := range concurrent {
			result := runPhase(ctx, root, phase, stdout, stderr)
			if result.Code == 130 {
				return 130
			}
			if result.Code != 0 {
				red = true
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

func runPhasesConcurrent(ctx context.Context, root string, phases []Phase, skipLog string, stdout, stderr io.Writer) int {
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
			reportCapabilitySkips(skipLog, stdout, stderr)
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
	if reportCapabilitySkips(skipLog, stdout, stderr) {
		red = true
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
	printConformanceTiming(root, phase, stdout)
	return result
}

// printConformanceTiming emits the conformance driver's per-check timing lines. The
// phase runs under non-verbose `go test`, which swallows a passing test binary's own
// stdout, so the file the driver leaves under the graded root's git dir is the only
// way that timing reaches gate output. Printing here covers both runners, both
// verdicts, and both modes — the canary's vacuity guard grades inner-mode output, so
// an outer-only print would leave it blind to the format. A root with no git dir or
// no timing file prints nothing at all.
func printConformanceTiming(root string, phase Phase, stdout io.Writer) {
	if phase.Name != conformancePhaseName {
		return
	}
	for _, line := range registry.ReadTimingLines(root) {
		fmt.Fprintln(stdout, line)
	}
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
