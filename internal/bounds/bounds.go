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

// The const block is the production policy registry. Callers name these entries;
// they do not redeclare values or locally reimplement classification.
//
// ControlRecordLimit is the one bound every control-record read applies — the journal,
// IDEAS.md, ROADMAP.md, a decision map, a spec — no matter which command reads it. One
// record read under two bounds is the divergence this entry exists to forbid: the same
// file would render rows on one surface and `unknown` on another, and a reader has no
// way to tell which answer is the repository's. The value bounds a hand-maintained
// markdown file read whole into memory, so it is deliberately far below ModelReadLimit,
// which bounds an untrusted provider response instead.
const (
	ProviderTimeout                 = 10 * time.Second
	GitRefreshTimeout               = 30 * time.Second
	GuardScanTimeout                = 5 * time.Second
	GateTimeout                     = 45 * time.Minute
	ModelReadLimit            int64 = 5 << 20
	OutlineFileLimit          int64 = 2 << 20
	ControlRecordLimit        int64 = 2 << 20
	OutlineRowLimit                 = 200
	IterationMin                    = 1
	IterationMax                    = 100
	MainIterationsDefault           = 12
	RefactorIterationsDefault       = 4
	MaxWall                         = 24 * time.Hour
	LeaseStale                      = time.Minute
	AssignmentStale                 = 7 * 24 * time.Hour
	CanaryInnerWidth                = 2
	TestDeadlineFloor               = 20 * time.Second
)

// TestDeadline derives an outer test deadline from the inner bound that deadline has
// to contain. Half the bound plus a fixed floor keeps the result strictly greater than
// the input for every entry in the registry, so a wait can never expire at the same
// instant as the window it is waiting out — an outer deadline equal to its inner window
// is a coin flip, not a bound. A negative bound contains nothing and clamps to the floor.
// The one place strictly-greater is unreachable is the top of the duration range, where
// the sum saturates at math.MaxInt64: equality there is the honest answer, and the
// alternative is the wrapped negative deadline that would expire immediately.
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
	if err := parent.Err(); err != nil {
		return ProcessResult{Status: ProcessCanceled, Err: err}
	}
	ctx, cancel := context.WithTimeout(parent, limit)
	defer cancel()
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		return ProcessResult{Status: ProcessStart, Err: err}
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			return ProcessResult{Status: ProcessComplete, Output: output.Bytes()}
		}
		exit := 1
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() > 0 {
			exit = cmd.ProcessState.ExitCode()
		}
		return ProcessResult{Status: ProcessExit, Output: output.Bytes(), Err: err, Exit: exit}
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		status := ProcessTimeout
		if parent.Err() != nil {
			status = ProcessCanceled
		}
		return ProcessResult{Status: status, Output: output.Bytes(), Err: ctx.Err()}
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
