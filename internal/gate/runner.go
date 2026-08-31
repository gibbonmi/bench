package gate

// The execution engine for gate phases: the DAG scheduler both runners share, the
// per-phase process launch with optional-binary resolution, and the prefixed
// output plumbing. The phase table and the PhasesCommand surface live in
// phases.go.

import (
	"bufio"
	"bytes"
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

	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const processGroupCancelGrace = 2 * time.Second

type processGroupCancelGraceKey struct{}

func processGroupGrace(ctx context.Context) time.Duration {
	if grace, ok := ctx.Value(processGroupCancelGraceKey{}).(time.Duration); ok {
		return grace
	}
	return processGroupCancelGrace
}

type processGroupResult struct {
	Code     int
	StartErr error
	// Cancelled reports that the context, not the command, decided the outcome. It is
	// the only way to tell a killed command from one that chose the same exit code.
	Cancelled bool
}

func runProcessGroupCommand(ctx context.Context, cmd *exec.Cmd) processGroupResult {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	if err := cmd.Start(); err != nil {
		if ctx.Err() != nil {
			return processGroupResult{Code: 130, StartErr: err, Cancelled: true}
		}
		return processGroupResult{Code: 1, StartErr: err}
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		drainProcessGroup(cmd.Process.Pid)
		return processGroupResult{Code: processExitCode(cmd, err)}
	case <-ctx.Done():
		// Both flavors of cancellation get the same cascade: a catchable signal, the
		// grace, then SIGKILL. They differ only in which signal and which code.
		//
		// The deadline needs the grace most. It fires with no operator watching. So what
		// the child says on its way out is the only account of what the run was stuck
		// on. Opening with SIGKILL would take that account with it.
		notice, code := syscall.SIGINT, 130
		if errors.Is(context.Cause(ctx), errGateTimeout) {
			notice, code = syscall.SIGTERM, 124
		}
		_ = syscall.Kill(-cmd.Process.Pid, notice)
		select {
		case <-done:
			drainProcessGroup(cmd.Process.Pid)
		case <-time.After(processGroupGrace(ctx)):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
			drainProcessGroup(cmd.Process.Pid)
		}
		return processGroupResult{Code: code, Cancelled: true}
	}
}

func drainProcessGroup(pgid int) {
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	for {
		if err := syscall.Kill(-pgid, 0); errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
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
	Name string
	// Argv is the phase's declared command. The report reads it to tell a `go test`
	// phase, whose red stream has a classifier, from every other phase, whose red stream
	// is failure rows line by line.
	Argv     []string
	Code     int
	Skipped  bool
	StartErr error
	// SkippedBy names the need whose own red or skip kept this phase from launching.
	// Such a phase is never red on its own. The verdict belongs to the phase that
	// actually failed, so the red set a fix loop reads stays free of cascade noise.
	SkippedBy string
	// Interrupted marks a phase the cancellation caught mid-run, which is what the
	// straggler report names. A phase that chose exit 130 itself carries the same code
	// and is an ordinary red, so the code alone cannot identify the interrupted set.
	Interrupted bool
	// ElapsedMS is the phase's wall time, and it is the green table's third cell. It
	// carries the same measurement the `phase.finish` log record takes, so the table and
	// the progress log cannot disagree about one phase. A phase that never launched keeps
	// the zero: nothing ran, so there is no time to report.
	ElapsedMS int64
}

// green reports whether a settled phase satisfies a dependent's edge. Both skip flavors
// are excluded: a dependent needs its need's artifact, and neither a skip nor a red
// produced one.
func (r phaseResult) green() bool {
	return r.Code == 0 && !r.Skipped && r.SkippedBy == ""
}

func runPhases(ctx context.Context, root string, phases []Phase, stdout, stderr io.Writer) int {
	ctx, closeRecord := beginPhaseRecord(ctx)
	defer closeRecord()
	skipLog, cleanup, err := newSkipLog()
	if err != nil {
		fmt.Fprintf(stderr, "gate: cannot open the capability skip log: %v\n", err)
		fmt.Fprintln(stderr, "gate: red")
		return 1
	}
	defer cleanup()
	return runPhasesSerial(ctx, root, withSkipLog(phases, skipLog), skipLog, stdout, stderr)
}

func runPhasesSerial(ctx context.Context, root string, phases []Phase, skipLog string, stdout, stderr io.Writer) int {
	streams := newPhaseStreams(stderr)
	streams.retain(gateRunStreamFile(ctx))
	results, cancelled := schedule(ctx, root, phases, streams.open)
	return aggregateAndReport(results, cancelled, streams, stdout, stderr, func() ([]string, string, bool) {
		return reportCapabilitySkips(skipLog)
	}, cacheFootprintReport(ctx, os.Environ(), gocache.Measure, gocache.Bound))
}

// prefixedPhaseWriters is the outer phase output plumbing. The mutex keeps each
// phase's two streams coherent with diagnostics emitted by the scheduler.
func prefixedPhaseWriters(stdout, stderr io.Writer) func(Phase) (io.Writer, io.Writer, func()) {
	var writeMu sync.Mutex
	return func(phase Phase) (io.Writer, io.Writer, func()) {
		out := newPrefixWriter(&writeMu, stdout, phase.Name)
		errOut := newPrefixWriter(&writeMu, stderr, phase.Name)
		return out, errOut, func() { out.Close(); errOut.Close() }
	}
}

// reportStragglers names, in table order, the phases a cancellation caught mid-run.
// This is the one thing a killed gate can still tell an operator about where it was
// stuck. Reading it off the settled results is the only race-free way to ask. A
// snapshot taken in the launch loop would be stale by the time the reaping loop
// finished with it.
func reportStragglers(results []phaseResult, stderr io.Writer) {
	var names []string
	for _, result := range results {
		if result.Interrupted {
			names = append(names, result.Name)
		}
	}
	if len(names) == 0 {
		return
	}
	fmt.Fprintln(stderr, "gate: cancelled; still running: "+strings.Join(names, ", "))
}

// schedule runs phases in dependency order and returns their results in table order,
// plus whether the run was interrupted. A phase launches once every need present in
// the table has settled green. A need that settled red or skipped resolves the
// dependent as skipped-with-cause without launching it. So a red phase costs the run
// only the work that actually depended on it. Sequential caps the run at one phase in
// flight and takes the first ready phase in declaration order.
//
// A need naming a phase absent from the table is already satisfied. phasesForMode
// filters the table after the edges are declared, so an inner run legitimately
// carries edges to phases it does not execute.
func schedule(ctx context.Context, root string, phases []Phase, open func(Phase) (io.Writer, io.Writer, func())) ([]phaseResult, bool) {
	index := make(map[string]int, len(phases))
	for i, phase := range phases {
		index[phase.Name] = i
	}
	results := make([]phaseResult, len(phases))
	settled := make([]bool, len(phases))
	launched := make([]bool, len(phases))
	done := make(chan int, len(phases))
	inFlight := 0

	for {
		progressed := false
		// Cancellation stops the launch loop rather than the whole scheduler. Phases
		// already in flight still have to be reaped before their output writers close,
		// but nothing new may start behind an interrupt.
		for i, phase := range phases {
			if settled[i] || launched[i] || ctx.Err() != nil {
				continue
			}
			blocker, ready := edgeState(phase, index, settled, results)
			if blocker != "" {
				results[i] = phaseResult{Name: phase.Name, Argv: phase.Argv, SkippedBy: blocker}
				endPhaseSpan(startPhaseSpan(ctx, phase.Name), results[i])
				logGateEvent(ctx, gateLogRecord{Event: "phase.skip", Phase: phase.Name, Detail: blocker})
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
			logGateEvent(ctx, gateLogRecord{Event: "phase.start", Phase: phase.Name, Root: phase.Dir})
			go func() {
				started := time.Now()
				span := startPhaseSpan(ctx, phase.Name)
				out, errOut, closeWriters := open(phase)
				results[i] = runPhase(ctx, root, phase, out, errOut)
				endPhaseSpan(span, results[i])
				// One reading of the clock feeds both the result and the log record. Two
				// calls would answer two different numbers, and the report and the progress
				// log would then disagree about the same phase.
				elapsed := time.Since(started).Milliseconds()
				results[i].ElapsedMS = elapsed
				exit := results[i].Code
				logGateEvent(ctx, gateLogRecord{Event: "phase.finish", Phase: phase.Name, Exit: &exit, ElapsedMS: elapsed})
				closeWriters()
				done <- i
			}()
			break
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
	}

	// The launch loop stops for exactly two reasons, and they settle the leftovers
	// differently. Behind a stopped context a phase simply never got its turn.
	//
	// With the context live, no further progress means the table's own edges deadlock
	// it. This is a defect the loader refuses, but an injected table can still carry
	// it. A phase the run can never launch has to be red, or a graph that executes
	// nothing reports green.
	interrupted := ctx.Err() != nil
	for i, phase := range phases {
		if settled[i] {
			continue
		}
		need := firstUnsettledNeed(phase, index, settled)
		if !interrupted {
			results[i] = phaseResult{Name: phase.Name, Argv: phase.Argv, Code: 1, StartErr: fmt.Errorf("stuck behind unsatisfied need %s", need)}
			continue
		}
		results[i] = phaseResult{Name: phase.Name, Argv: phase.Argv, Skipped: true, SkippedBy: need, StartErr: errInterruptedBeforeLaunch}
	}
	// A deadline is not an interrupt: it keeps its own exit code and its summaries, so
	// only an operator's signal suppresses the run's report.
	return results, interrupted && !errors.Is(context.Cause(ctx), errGateTimeout)
}

// errInterruptedBeforeLaunch is the cause a phase carries when the run was stopped
// before its turn came. It is a skip rather than a red because nothing graded it. It
// is never the zero value, so no summary can call such a phase green.
var errInterruptedBeforeLaunch = errors.New("interrupted before launch")

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

func runPhase(ctx context.Context, root string, phase Phase, stdout, stderr io.Writer) phaseResult {
	result := phaseResult{Name: phase.Name, Argv: phase.Argv}
	if len(phase.Argv) == 0 || phase.Argv[0] == "" {
		result.Code = 1
		result.StartErr = fmt.Errorf("empty argv")
		return result
	}
	// This checks the working directory before the binary. chdir fails ENOENT exactly
	// as a missing binary does. So an optional phase whose dir is a typo would
	// otherwise report itself not installed and take its check off the gate.
	if phase.Dir != "" {
		if err := usableDir(phase.Dir); err != nil {
			result.Code = 1
			result.StartErr = err
			return result
		}
	}
	argv := phase.Argv
	if phase.Optional {
		resolved, absent := phaseToolAbsent(phase)
		if absent {
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
	if phase.Dir != "" {
		cmd.Dir = phase.Dir
	}
	var observed bytes.Buffer
	cmd.Stdout = stdout
	if len(phase.ExpectedRuns) != 0 {
		cmd.Stdout = io.MultiWriter(stdout, &observed)
	}
	cmd.Stderr = stderr
	baseEnv, err := gateEnv()
	if err != nil {
		// The phase reds before the child starts. A child launched without the entry
		// would write to the ambient cache, which is the state this refusal exists to
		// prevent.
		fmt.Fprintf(stderr, "%s cache environment unavailable: %v\n", phase.Name, err)
		result.Code = 1
		result.StartErr = err
		return result
	}
	cmd.Env = mergeEnv(baseEnv, phase.Env)
	run := runProcessGroupCommand(ctx, cmd)
	result.Interrupted = run.Cancelled
	if run.StartErr != nil {
		if run.Cancelled {
			result.Code = run.Code
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
	for _, name := range phase.ExpectedRuns {
		if strings.Contains(observed.String(), "=== RUN   "+name) {
			continue
		}
		fmt.Fprintf(stderr, "%s test did not run: %s\n", phase.Name, name)
		result.Code = 1
	}
	return result
}

// usableDir reports why a phase's working directory cannot be entered, or nil when it
// can. This is the one place a dir defect is named as itself. Leaving it to exec
// hands back a bare ENOENT that reads like a missing binary.
func usableDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("unusable working directory %s: %v", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("unusable working directory %s: not a directory", dir)
	}
	return nil
}

// mergeEnv applies overrides over base strip-then-set. Every base entry for an
// overridden key is dropped before the override is appended, so the child is handed
// one value per key. Plain appending leaves both, and a phase may not be built on
// which one a program's getenv answers with. withSkipLog rides the same path; the
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

// phaseToolAbsent reports that phase is optional and the tool it invokes is not on this
// host, handing back the resolved path when it is. Absent is the one way a phase
// settles green having graded nothing at all. So the per-component scoping withholds
// that component's evidence on this same answer.
//
// Two spellings of the condition would let the runner skip a component the scoping
// had already credited with a slot. The component would then skip forever, because
// its declared inputs do not move when the tool is installed.
//
// A required phase is never absent here. Its missing binary is that phase's own red,
// so the run does grade the component, in the direction that costs no evidence.
func phaseToolAbsent(phase Phase) (resolved string, absent bool) {
	if !phase.Optional || len(phase.Argv) == 0 || phase.Argv[0] == "" {
		return "", false
	}
	path, present := resolveOnPath(phase.Argv[0])
	return path, !present
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

// beginPhaseRecord attaches the seam record of the gate run that started this process.
// The phases run in their own process, so the repository and the trace come from the
// environment the run's parent composed. A process outside a recorded run records
// nothing: TracerFrom then answers a no-op tracer, and the phase spans cost nothing.
func beginPhaseRecord(ctx context.Context) (context.Context, func()) {
	root := os.Getenv(otelRootEnv)
	if root == "" {
		return ctx, func() {}
	}
	provider := otelrecord.NewProvider("", root)
	ctx = otelrecord.WithTracer(ctx, provider.Tracer())
	if parent := os.Getenv(otelTraceparentEnv); parent != "" {
		ctx = propagation.TraceContext{}.Extract(ctx, propagation.MapCarrier{"traceparent": parent})
	}
	return ctx, func() { _ = provider.Shutdown(context.WithoutCancel(ctx)) }
}

// startPhaseSpan opens one phase's span. The phase name is the span name, so a reader
// derives the per-phase time from the record without a second attribute for it.
func startPhaseSpan(ctx context.Context, name string) trace.Span {
	_, span := otelrecord.TracerFrom(ctx).Start(ctx, name,
		trace.WithAttributes(attribute.String(otelrecord.AttrSeam, otelGatePhaseSeam)))
	return span
}

// endPhaseSpan closes a settled phase's span with its exit, and names the need that
// skipped it. A cascade skip stays attributed to the phase that actually failed.
func endPhaseSpan(span trace.Span, result phaseResult) {
	span.SetAttributes(attribute.String(otelrecord.AttrOutcome, phaseSpanOutcome(result)))
	if result.SkippedBy != "" {
		span.SetAttributes(attribute.String(otelrecord.AttrOutcomeBlocker, result.SkippedBy))
	}
	span.End()
}

func phaseSpanOutcome(result phaseResult) string {
	if result.Skipped || result.SkippedBy != "" {
		return otelrecord.OutcomeSkipped
	}
	return otelrecord.ExitOutcome(result.Code)
}
