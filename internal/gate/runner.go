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
	// Cancelled reports that the context, not the command, decided the outcome — the
	// only way to tell a killed command from one that chose the same exit code.
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
		// Both flavors of cancellation get the same cascade — a catchable signal, the
		// grace, then SIGKILL — and differ only in which signal and which code. The
		// deadline needs the grace most: it fires with no operator watching, so what
		// the child says on its way out is the only account of what the run was stuck
		// on, and opening with SIGKILL takes that account with it.
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
	Name     string
	Code     int
	Skipped  bool
	StartErr error
	// SkippedBy names the need whose own red or skip kept this phase from launching.
	// Such a phase is never red on its own: the verdict belongs to the phase that
	// actually failed, so the red set a fix loop reads stays free of cascade noise.
	SkippedBy string
	// Interrupted marks a phase the cancellation caught mid-run, which is what the
	// straggler report names. A phase that chose exit 130 itself carries the same code
	// and is an ordinary red, so the code alone cannot identify the set.
	Interrupted bool
}

// green reports whether a settled phase satisfies a dependent's edge. Both skip flavors
// are excluded: a dependent needs its need's artifact, and neither a skip nor a red
// produced one.
func (r phaseResult) green() bool {
	return r.Code == 0 && !r.Skipped && r.SkippedBy == ""
}

func runPhases(ctx context.Context, root string, phases []Phase, stdout, stderr io.Writer) int {
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
	results, cancelled := schedule(ctx, root, phases, prefixedPhaseWriters(stdout, stderr))
	return aggregateAndReport(results, cancelled, stdout, stderr, func() bool {
		return reportCapabilitySkips(skipLog, stdout, stderr)
	})
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

// aggregateAndReport is the one verdict tail every settled schedule reports through:
// the per-phase summaries, any extra red-reporting checks (capability skips, the
// stripped-subject skip posture), and the `gate: red` / `gate: green` line. This is the
// operator's view of one command whichever schedule produced the results, so an edit to
// the reported shape lands everywhere at once.
//
// An interrupt is not a verdict, so a cancelled run publishes neither summaries nor a
// gate line — reporting one would grade phases that never got to answer. Naming the
// stragglers is not a verdict either: it says what the run was doing.
func aggregateAndReport(results []phaseResult, cancelled bool, stdout, stderr io.Writer, redReports ...func() bool) int {
	if cancelled {
		reportStragglers(results, stderr)
		return 130
	}
	red := false
	for _, result := range results {
		fmt.Fprintln(stdout, phaseSummary(result))
		if result.Code != 0 {
			red = true
		}
	}
	for _, report := range redReports {
		if report() {
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

// reportStragglers names, in table order, the phases a cancellation caught mid-run —
// the one thing a killed gate can still tell an operator about where it was stuck.
// Reading it off the settled results is the only race-free way to ask, since a snapshot
// taken in the launch loop would be stale by the time the reaping loop finished with it.
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
// plus whether the run was interrupted. A phase launches once every need present in the
// table has settled green; a need that settled red or skipped resolves the dependent as
// skipped-with-cause without launching it, so a red phase costs the run only the work
// that actually depended on it. Sequential caps the run at one phase in flight and
// takes the first ready phase in declaration order.
//
// A need naming a phase absent from the table is already satisfied: phasesForMode
// filters the table after the edges are declared, so an inner run legitimately carries
// edges to phases it does not execute.
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
		// Cancellation stops the launch loop rather than the whole scheduler: phases
		// already in flight still have to be reaped before their output writers can
		// be closed, but nothing new may start behind an interrupt.
		for i, phase := range phases {
			if settled[i] || launched[i] || ctx.Err() != nil {
				continue
			}
			blocker, ready := edgeState(phase, index, settled, results)
			if blocker != "" {
				results[i] = phaseResult{Name: phase.Name, SkippedBy: blocker}
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
				out, errOut, closeWriters := open(phase)
				results[i] = runPhase(ctx, root, phase, out, errOut)
				exit := results[i].Code
				logGateEvent(ctx, gateLogRecord{Event: "phase.finish", Phase: phase.Name, Exit: &exit, ElapsedMS: time.Since(started).Milliseconds()})
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
	// differently. Behind a stopped context a phase simply never got its turn; with the
	// context live, no further progress means the table's own edges deadlock it — a
	// defect the loader refuses but an injected table can still carry — and a phase the
	// run can never launch has to be red, or a graph that executes nothing reports green.
	interrupted := ctx.Err() != nil
	for i, phase := range phases {
		if settled[i] {
			continue
		}
		need := firstUnsettledNeed(phase, index, settled)
		if !interrupted {
			results[i] = phaseResult{Name: phase.Name, Code: 1, StartErr: fmt.Errorf("stuck behind unsatisfied need %s", need)}
			continue
		}
		results[i] = phaseResult{Name: phase.Name, Skipped: true, SkippedBy: need, StartErr: errInterruptedBeforeLaunch}
	}
	// A deadline is not an interrupt: it keeps its own exit code and its summaries, so
	// only an operator's signal suppresses the run's report.
	return results, interrupted && !errors.Is(context.Cause(ctx), errGateTimeout)
}

// errInterruptedBeforeLaunch is the cause a phase carries when the run was stopped
// before its turn came. It is a skip rather than a red because nothing graded it, and it
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

func phaseSummary(result phaseResult) string {
	if result.SkippedBy != "" {
		return "phase " + result.Name + ": skipped (needs " + result.SkippedBy + ")"
	}
	if result.Skipped {
		if result.StartErr != nil {
			return fmt.Sprintf("phase %s: skipped (%v)", result.Name, result.StartErr)
		}
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
	// A working directory the run cannot enter is checked before the binary is: chdir
	// fails ENOENT exactly as a missing binary does, so an optional phase whose dir is a
	// typo would otherwise report itself not installed and take its check off the gate.
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
	cmd.Env = mergeEnv(gateEnv(), phase.Env)
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
// can. The check is the one place a dir defect is named as itself; leaving it to exec
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

// phaseToolAbsent reports that phase is optional and the tool it invokes is not on this
// host, handing back the resolved path when it is. Absent is the one way a phase settles
// green having graded nothing at all, so the per-component scoping withholds that
// component's evidence on this same answer: two spellings of the condition would let the
// runner skip a component the scoping had already credited with a slot, and the component
// would then skip forever — its declared inputs do not move when the tool is installed.
//
// A required phase is never absent here. Its missing binary is that phase's own red, so the
// run does grade the component, in the direction that costs no evidence.
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
