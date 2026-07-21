// Package bounds owns the fixed resource policy for Bench's Go runtime.
package bounds

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// The const block is the production policy registry. Callers name these entries;
// they do not redeclare values or locally reimplement classification.
const (
	ProviderTimeout                 = 10 * time.Second
	GitRefreshTimeout               = 30 * time.Second
	GuardScanTimeout                = 5 * time.Second
	GateTimeout                     = 45 * time.Minute
	ModelReadLimit            int64 = 5 << 20
	OutlineFileLimit          int64 = 2 << 20
	OutlineRowLimit                 = 200
	IterationMin                    = 1
	IterationMax                    = 100
	MainIterationsDefault           = 12
	RefactorIterationsDefault       = 4
	MaxWall                         = 24 * time.Hour
)

func Offline() bool { return os.Getenv("BENCH_OFFLINE") == "1" }

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
