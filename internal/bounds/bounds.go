// Package bounds owns the fixed resource policy for Bench's Go runtime.
// The test switch removes verdict windows, including the 45-minute gate timeout.
package bounds

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/subprocess"
)

// The const block is the production policy registry. Callers name these entries. They do
// not redeclare values or locally reimplement classification.
//
// ControlRecordLimit is the one bound every control-record read applies: the journal,
// capture/IDEAS.md, ROADMAP.md, a decision map, and a spec, no matter which command reads
// it. A record read under two different bounds is the divergence this entry forbids: the
// same file would render rows on one surface and `unknown` on another, and a reader could
// not tell which answer is the repository's. The value bounds a hand-maintained markdown
// file read whole into memory, so it stays far below ModelReadLimit, which bounds an
// untrusted provider response instead.
const (
	ProviderTimeout             = 10 * time.Second
	EnvironmentDiscoveryTimeout = 2000 * time.Millisecond
	GitRefreshTimeout           = 30 * time.Second
	RefCheckTimeout             = 2 * time.Second
	HandoffLockTimeout          = 2 * time.Second
	CaptureLockTimeout          = 2 * time.Second
	IntentLockTimeout           = 2 * time.Second
	// Cancellation graces stay finite through FixedWindow.
	ProcessGroupCancelGrace = 2 * time.Second
	BuilderCancelGrace      = 2 * time.Second
	PublicationStepGrace    = 2 * time.Second
	// ProcessShutdownGrace is the grace of a catchable shutdown path that has no grace of its
	// own. ProcessFinalWindow bounds every teardown obligation after forced termination.
	ProcessShutdownGrace = 2 * time.Second
	ProcessFinalWindow   = 3 * time.Second
	// Git 2.43.0 can block while reading malformed worktree admin files; this
	// backstop retires when upstream bounds those reads itself.
	WorktreeListTimeout = 15 * time.Second
	GuardScanTimeout    = 5 * time.Second
	// EvidenceReadyTimeout bounds the release-evidence readiness hold. A run whose observer
	// never releases the hold refuses after this window instead of waiting without bound. It is
	// a wall limit, so the owner applies it through FixedWindow.
	EvidenceReadyTimeout = 30 * time.Second
	GateTimeout          = 45 * time.Minute
	// GoTestTimeout is the -timeout of each test binary in a Bench-owned go test command.
	// GateTimeout bounds the whole gate, which is a different fact.
	GoTestTimeout = 10 * time.Minute
	// ReleaseProbeTimeout bounds one release-preflight run that a conformance probe drives. It
	// stays below GoTestTimeout, so the process-group owner stops the probe's group before the
	// test binary dies and leaves that group behind.
	ReleaseProbeTimeout = 5 * time.Minute
	// RecordClockSkew is how far a gate record time may sit after the reader's now. A
	// host wall clock can step back by up to a second under load, between the write and
	// a later read.
	RecordClockSkew = 2 * time.Second
	// PackageLoadTimeout bounds one Go package-loader invocation. The loader opens every
	// file it walks, so one FIFO anywhere under a loaded package tree blocks it in open(2)
	// with no deadline of its own. The value sits far above a cold expansion of `./...` on
	// a large tree, so only a loader that will never return reaches it.
	PackageLoadTimeout              = 30 * time.Second
	ModelReadLimit            int64 = 5 << 20
	OutlineFileLimit          int64 = 2 << 20
	ControlRecordLimit        int64 = 2 << 20
	IterationMin                    = 1
	IterationMax                    = 100
	MainIterationsDefault           = 12
	RefactorIterationsDefault       = 4
	MaxWall                         = 24 * time.Hour
	LeaseStale                      = time.Minute
	AssignmentStale                 = 7 * 24 * time.Hour
	// CoverageRowStories caps how many stories one acceptance-coverage row may reference;
	// a row spanning more is an outcome family no single test can go red on.
	CoverageRowStories = 4
	TestDeadlineFloor  = 20 * time.Second
	// PreviewRuneLimit caps how many code points a preview of operator-influenced text
	// keeps before the byte-count suffix replaces the rest. A `bench test` failure line
	// and a guard refusal both carry a preview, and each must stay readable inside one
	// bounded response. The value holds a full commit subject or an objective, and it
	// stays short enough that two previews and their surrounding row fit one screen.
	PreviewRuneLimit = 240
	// RecordSegmentLimit bounds the live segment of one repository's seam record. An
	// append that would take the segment past it seals the segment first.
	// RecordSegmentsRetained is how many sealed segments the record keeps. Together they
	// bound the record's disk use; the reviewer owns both sizes. The limit is 16 MiB.
	RecordSegmentLimit     int64 = 1 << 24
	RecordSegmentsRetained       = 8
	// RecordMemoryRetained is how many retained shift-notes files the record keeps beside
	// its segments; the reviewer owns the size.
	RecordMemoryRetained = 64
	// ResponseLines and ResponseBytes are the line value and the byte value of one public
	// Bench response. A response with more lines or more bytes prints its head, one spill
	// line, and its tail, and a private file holds the complete output. The response owner
	// is the one reader of these values.
	ResponseLines = 10
	ResponseBytes = 4096
	// SubjectFirstLineLimit bounds the read of a gate subject's first line.
	SubjectFirstLineLimit int64 = 4096
	// ExecWaitDelay bounds how long `bench worktree exec` waits for its child's output
	// after the child exits, so a descendant that holds a pipe open cannot hang the verb.
	// An interrupted child gets the same delay to exit before the verb kills its group.
	ExecWaitDelay = 3 * time.Second
)

// TestDeadline derives an outer test deadline from the inner bound that deadline has
// to contain. Half the bound plus a fixed floor keeps the result strictly greater than
// the input for every entry in the registry, so a wait can never expire at the same
// instant as the window it is waiting out. An outer deadline equal to its inner window
// is a coin flip, not a bound. A negative bound contains nothing and clamps to the floor.
// At the top of the duration range the sum saturates at math.MaxInt64, so the result
// equals the input there; the alternative is a wrapped negative deadline that expires
// immediately.
func TestDeadline(inner time.Duration) time.Duration {
	if inner < 0 {
		inner = 0
	}
	const ceiling = time.Duration(math.MaxInt64)
	half := inner / 2
	if half > ceiling-inner || TestDeadlineFloor > ceiling-inner-half {
		return ceiling
	}
	return inner + half + TestDeadlineFloor
}

// TestTimeoutVerdict renders what a wait prints when the window TestDeadline gave it runs
// out. The verdict names the wait and the exhausted window, so a red run reads as a
// timeout and not as a failed assertion, and the reader sees how long the wait allowed
// without opening the test file. The renderer takes no *testing.T, so this package stays
// free of the testing import and a helper process with no test of its own reports the
// same verdict.
func TestTimeoutVerdict(wait string, window time.Duration) string {
	return "timeout: waited " + window.String() + " for " + wait
}

func Offline() bool { return os.Getenv("BENCH_OFFLINE") == "1" }

// UnboundedWaitsEnv enables unbounded verdict waits only when its value is exactly 1.
const UnboundedWaitsEnv = "BENCH_TEST_UNBOUNDED_WAITS"

// Unbounded removes a wait's own deadline while preserving parent cancellation.
const Unbounded time.Duration = -1

// VerdictWindow selects the verdict policy before a caller applies a test override.
func VerdictWindow(policy time.Duration) time.Duration {
	if os.Getenv(UnboundedWaitsEnv) == "1" {
		return Unbounded
	}
	return policy
}

// FixedWindow preserves cancel graces, poll intervals, and operator wall limits.
func FixedWindow(window time.Duration) time.Duration { return window }

func Context(parent context.Context, limit time.Duration) (context.Context, context.CancelFunc) {
	if limit == Unbounded {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, limit)
}

func ContextCause(parent context.Context, limit time.Duration, cause error) (context.Context, context.CancelFunc) {
	if limit == Unbounded {
		return context.WithCancel(parent)
	}
	return context.WithTimeoutCause(parent, limit, cause)
}

type ProcessStatus string

const (
	ProcessComplete ProcessStatus = "complete"
	ProcessTimeout  ProcessStatus = "timeout"
	ProcessCanceled ProcessStatus = "canceled"
	ProcessExit     ProcessStatus = "exit"
	ProcessStart    ProcessStatus = "start"
)

type ProcessResult struct {
	Status ProcessStatus
	Output []byte
	Err    error
	Exit   int
}

// Run retains the merged stdout and stderr in Result.Output.
func Run(parent context.Context, limit time.Duration, cmd *exec.Cmd) ProcessResult {
	return run(parent, limit, cmd, nil)
}

// RunOutput preserves command stdout in stdout while retaining stderr in Result.Output.
// stdout receives the captured bytes after the stream's worker finishes.
func RunOutput(parent context.Context, limit time.Duration, cmd *exec.Cmd, stdout io.Writer) ProcessResult {
	return run(parent, limit, cmd, stdout)
}

// ProcessWindow builds a subprocess window of limit. An Unbounded limit never expires.
func ProcessWindow(limit time.Duration) subprocess.Window {
	return func() (<-chan time.Time, func() bool) {
		if limit == Unbounded {
			return nil, func() bool { return false }
		}
		timer := time.NewTimer(limit)
		return timer.C, timer.Stop
	}
}

// ProcessPoll paces group-absence observations inside a final cleanup window.
var ProcessPoll = ProcessWindow(FixedWindow(10 * time.Millisecond))

// processPolicy is the runner's policy: a normal exit keeps its group, and cancellation
// sends KILL at once. The final window stays fixed, so teardown is finite under the
// unbounded-verdict switch.
var processPolicy = subprocess.Policy{
	Normal: subprocess.Preserve,
	Final:  ProcessWindow(FixedWindow(ProcessFinalWindow)),
	Poll:   ProcessPoll,
}

// run returns the captured stderr, or the merged output when stdout is nil. The child's group
// registers on every resource this process inherited before it starts.
func run(parent context.Context, limit time.Duration, cmd *exec.Cmd, stdout io.Writer) ProcessResult {
	if err := parent.Err(); err != nil {
		return ProcessResult{Status: ProcessCanceled, Err: err}
	}
	ctx, cancel := Context(parent, limit)
	defer cancel()
	streams := subprocess.Streams{Stdout: subprocess.Output{Capture: true}, Stderr: subprocess.Output{Capture: true}, MergeStderr: stdout == nil}
	request := func(context.Context) (subprocess.Cause, syscall.Signal) {
		if parent.Err() != nil {
			return subprocess.CauseInterruption, syscall.SIGKILL
		}
		return subprocess.CauseDeadline, syscall.SIGKILL
	}
	outcome, err := subprocess.RunInherited(ctx, cmd, processPolicy, streams, request, subprocess.Registration{})
	if err != nil {
		return ProcessResult{Status: ProcessStart, Err: err}
	}
	process := outcome.Process
	output := outcome.Streams.Stderr.Data
	if stdout == nil {
		output = outcome.Streams.Stdout.Data
	} else if _, err := stdout.Write(outcome.Streams.Stdout.Data); err != nil && process.WaitErr == nil {
		process.WaitErr = err
	}
	switch {
	case process.StartErr != nil && !errors.Is(process.StartErr, subprocess.ErrNoStart):
		return ProcessResult{Status: ProcessStart, Err: process.StartErr}
	case outcome.Cancellation.Cause == subprocess.CauseInterruption:
		return ProcessResult{Status: ProcessCanceled, Output: output, Err: ctx.Err()}
	case outcome.Cancellation.Cause == subprocess.CauseDeadline:
		return ProcessResult{Status: ProcessTimeout, Output: output, Err: ctx.Err()}
	}
	err = process.WaitErr
	for _, stream := range []subprocess.Stream{outcome.Streams.Stdout, outcome.Streams.Stderr} {
		if err == nil && stream.Status != subprocess.StreamComplete && stream.Status != subprocess.StreamNone {
			err = fmt.Errorf("process stream %s: %v", stream.Status, stream.Err)
		}
	}
	if err == nil {
		return ProcessResult{Status: ProcessComplete, Output: output}
	}
	exit := 1
	if process.Code > 0 {
		exit = process.Code
	}
	return ProcessResult{Status: ProcessExit, Output: output, Err: err, Exit: exit}
}

type ReadStatus string

const (
	ReadComplete  ReadStatus = "complete"
	ReadOversized ReadStatus = "oversized"
	ReadFailed    ReadStatus = "failed"
)

type ReadResult struct {
	Status ReadStatus
	Data   []byte
	Err    error
}

func Read(reader io.Reader, limit int64) ReadResult {
	data, err := io.ReadAll(io.LimitReader(reader, limit+1))
	if err != nil {
		return ReadResult{Status: ReadFailed, Err: err}
	}
	if int64(len(data)) > limit {
		return ReadResult{Status: ReadOversized, Data: data[:limit], Err: errors.New("read limit exceeded")}
	}
	return ReadResult{Status: ReadComplete, Data: data}
}
