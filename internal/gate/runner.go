package gate

// The execution engine for gate phases: the DAG scheduler both runners share, the
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
	// SkippedBy names the need whose own red or skip kept this phase from launching.
	// Such a phase is never red on its own: the verdict belongs to the phase that
	// actually failed, so the red set a fix loop reads stays free of cascade noise.
	SkippedBy string
}

// green reports whether a settled phase satisfies a dependent's edge. Both skip flavors
// are excluded: a dependent needs its need's artifact, and neither a skip nor a red
// produced one.
func (r phaseResult) green() bool {
	return r.Code == 0 && !r.Skipped && r.SkippedBy == ""
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
	results, cancelled := schedule(ctx, root, phases, true, func(Phase) (io.Writer, io.Writer, func()) {
		return stdout, stderr, func() {}
	})
	if cancelled {
		return 130
	}
	for _, result := range results {
		if result.Code != 0 {
			fmt.Fprintln(stderr, "gate: red")
			return 1
		}
	}
	fmt.Fprintln(stdout, "gate: green")
	return 0
}

func runPhasesConcurrent(ctx context.Context, root string, phases []Phase, skipLog string, stdout, stderr io.Writer) int {
	var writeMu sync.Mutex
	results, cancelled := schedule(ctx, root, phases, false, func(phase Phase) (io.Writer, io.Writer, func()) {
		out := newPrefixWriter(&writeMu, stdout, phase.Name)
		err := newPrefixWriter(&writeMu, stderr, phase.Name)
		return out, err, func() { out.Close(); err.Close() }
	})
	// An interrupt is not a verdict, so a cancelled run publishes neither summaries
	// nor a gate line — reporting one would grade phases that never got to answer.
	if cancelled {
		return 130
	}

	red := false
	for _, result := range results {
		fmt.Fprintln(stdout, phaseSummary(result))
		if result.Code != 0 {
			red = true
		}
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

// schedule runs phases in dependency order and returns their results in table order,
// plus whether the run was interrupted. A phase launches once every need present in the
// table has settled green; a need that settled red or skipped resolves the dependent as
// skipped-with-cause without launching it, so a red phase costs the run only the work
// that actually depended on it. Sequential caps the run at one phase in flight and
// takes the first ready phase in declaration order, which is the topological order
// inner mode's pinned byte shape needs.
//
// A need naming a phase absent from the table is already satisfied: phasesForMode
// filters the table after the edges are declared, so an inner run legitimately carries
// edges to phases it does not execute.
func schedule(ctx context.Context, root string, phases []Phase, sequential bool, open func(Phase) (io.Writer, io.Writer, func())) ([]phaseResult, bool) {
	index := make(map[string]int, len(phases))
	for i, phase := range phases {
		index[phase.Name] = i
	}
	results := make([]phaseResult, len(phases))
	settled := make([]bool, len(phases))
	launched := make([]bool, len(phases))
	done := make(chan int, len(phases))
	inFlight, cancelled := 0, false

	for {
		progressed := false
		// Cancellation stops the launch loop rather than the whole scheduler: phases
		// already in flight still have to be reaped before their output writers can
		// be closed, but nothing new may start behind an interrupt.
		for i, phase := range phases {
			if settled[i] || launched[i] || cancelled || ctx.Err() != nil {
				continue
			}
			blocker, ready := edgeState(phase, index, settled, results)
			if blocker != "" {
				results[i] = phaseResult{Name: phase.Name, SkippedBy: blocker}
				settled[i] = true
				progressed = true
				continue
			}
			if !ready {
				continue
			}
			launched[i] = true
			inFlight++
			progressed = true
			i, phase := i, phase
			go func() {
				out, errOut, closeWriters := open(phase)
				results[i] = runPhase(ctx, root, phase, out, errOut)
				closeWriters()
				done <- i
			}()
			if sequential {
				break
			}
		}
		if inFlight == 0 {
			if !progressed {
				break
			}
			continue
		}
		i := <-done
		inFlight--
		settled[i] = true
		if results[i].Code == 130 {
			cancelled = true
		}
	}

	// Anything still unsettled is either behind an interrupt or on a cycle the loader
	// should have refused. Naming its first unmet need keeps a never-launched phase
	// from reporting the zero value, which reads as green.
	for i, phase := range phases {
		if !settled[i] {
			results[i] = phaseResult{Name: phase.Name, SkippedBy: firstUnsettledNeed(phase, index, settled)}
		}
	}
	return results, cancelled
}

// edgeState reports the need blocking phase permanently, if any, and otherwise whether
// every need has settled green. A need that settled non-green wins over one still
// running: the outcome is already decided, so the dependent need not wait to learn it.
func edgeState(phase Phase, index map[string]int, settled []bool, results []phaseResult) (blocker string, ready bool) {
	ready = true
	for _, need := range phase.Needs {
		i, present := index[need]
		if !present {
			continue
		}
		if !settled[i] {
			ready = false
			continue
		}
		if !results[i].green() {
			return need, false
		}
	}
	return "", ready
}

func firstUnsettledNeed(phase Phase, index map[string]int, settled []bool) string {
	for _, need := range phase.Needs {
		if i, present := index[need]; present && !settled[i] {
			return need
		}
	}
	return ""
}

func phaseSummary(result phaseResult) string {
	if result.SkippedBy != "" {
		return "phase " + result.Name + ": skipped (needs " + result.SkippedBy + ")"
	}
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

	if phase.Name == conformancePhaseName {
		// The run boundary is here, not at the read below: clearing first is what makes
		// the print answer for this run, so a file some earlier gate, a killed run, or a
		// different invocation left in the git dir cannot be printed as if it were ours.
		_ = registry.ClearTiming(root)
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	if phase.Dir != "" {
		cmd.Dir = phase.Dir
		if !filepath.IsAbs(phase.Dir) {
			cmd.Dir = filepath.Join(root, phase.Dir)
		}
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = mergeEnv(gateEnv(), phase.Env)
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

// mergeEnv applies overrides over base strip-then-set: every base entry for an
// overridden key is dropped before the override is appended, so the child is handed one
// value per key. Plain appending leaves both, and which one a program's getenv answers
// with is not something a phase may be built on. withSkipLog rides the same path — the
// capability log variable is an override like any other.
func mergeEnv(base, overrides []string) []string {
	if len(overrides) == 0 {
		return base
	}
	overridden := make(map[string]bool, len(overrides))
	for _, entry := range overrides {
		overridden[envKey(entry)] = true
	}
	merged := make([]string, 0, len(base)+len(overrides))
	for _, entry := range base {
		if !overridden[envKey(entry)] {
			merged = append(merged, entry)
		}
	}
	return append(merged, overrides...)
}

func envKey(entry string) string {
	if idx := strings.Index(entry, "="); idx >= 0 {
		return entry[:idx]
	}
	return entry
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
