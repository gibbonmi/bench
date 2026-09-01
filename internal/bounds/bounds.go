// Package bounds owns the fixed resource policy for Bench's Go runtime.
package bounds

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math"
	"os"
	"os/exec"
	"syscall"
	"time"
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
	// Git 2.43.0 can block while reading malformed worktree admin files; this
	// backstop retires when upstream bounds those reads itself.
	WorktreeListTimeout = 15 * time.Second
	GuardScanTimeout    = 5 * time.Second
	GateTimeout         = 45 * time.Minute
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

func Offline() bool { return os.Getenv("BENCH_OFFLINE") == "1" }

func Context(parent context.Context, limit time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, limit)
}

func ContextCause(parent context.Context, limit time.Duration, cause error) (context.Context, context.CancelFunc) {
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

func Run(parent context.Context, limit time.Duration, cmd *exec.Cmd) ProcessResult {
	var output bytes.Buffer
	result := run(parent, limit, cmd, &output, &output)
	result.Output = output.Bytes()
	return result
}

// RunOutput preserves command stdout in stdout while retaining stderr in Result.Output.
func RunOutput(parent context.Context, limit time.Duration, cmd *exec.Cmd, stdout io.Writer) ProcessResult {
	var stderr bytes.Buffer
	result := run(parent, limit, cmd, stdout, &stderr)
	result.Output = stderr.Bytes()
	return result
}

func run(parent context.Context, limit time.Duration, cmd *exec.Cmd, stdout, stderr io.Writer) ProcessResult {
	if err := parent.Err(); err != nil {
		return ProcessResult{Status: ProcessCanceled, Err: err}
	}
	ctx, cancel := context.WithTimeout(parent, limit)
	defer cancel()
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Start(); err != nil {
		return ProcessResult{Status: ProcessStart, Err: err}
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			return ProcessResult{Status: ProcessComplete}
		}
		exit := 1
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() > 0 {
			exit = cmd.ProcessState.ExitCode()
		}
		return ProcessResult{Status: ProcessExit, Err: err, Exit: exit}
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		status := ProcessTimeout
		if parent.Err() != nil {
			status = ProcessCanceled
		}
		return ProcessResult{Status: status, Err: ctx.Err()}
	}
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
