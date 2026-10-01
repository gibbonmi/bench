// Package subprocess runs an external command and captures its outcome, and it
// holds the termination signals a Bench command traps. Capture and CaptureMerged
// share one exit-code derivation and differ only in whether stdout and stderr stay
// separate or interleave. The package does not start, cancel, or kill a process
// group. Each owner that detaches a child into its own group does that work itself.
package subprocess

import (
	"bytes"
	"os/exec"
)

// Termination describes whether a command produced a numeric exit code or
// stopped before one could be observed.
type Termination uint8

const (
	// TerminationCompleted means the command produced a numeric exit code.
	TerminationCompleted Termination = iota
	// TerminationSpawnFailed means the command could not be started.
	TerminationSpawnFailed
	// TerminationSignaled means the command stopped because of a signal.
	TerminationSignaled
)

// Aborted reports whether the command stopped without a normal numeric exit.
func (termination Termination) Aborted() bool {
	return termination == TerminationSpawnFailed || termination == TerminationSignaled
}

// Result is the outcome of running a command. ExitCode is 0 on success, the
// process's own code when it completed, -1 when it was signaled, or 1 when it
// never started. Err is the raw error from the run. Termination records whether
// the exit was normal. Capture fills Stdout and Stderr separately;
// CaptureMerged interleaves both into Stdout and leaves Stderr empty.
type Result struct {
	ExitCode    int
	Termination Termination
	Stdout      string
	Stderr      string
	Err         error
}

// Capture runs cmd with stdout and stderr buffered separately. Callers that
// parse stdout (JSON, machine output) need this so subprocess stderr chatter
// cannot corrupt the parse.
func Capture(cmd *exec.Cmd) Result {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode, termination := outcome(cmd, err)
	return Result{ExitCode: exitCode, Termination: termination, Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
}

// CaptureMerged runs cmd with stdout and stderr interleaved into one stream, in
// the order a terminal would show them. EXPECT-style matching that asserts on
// ordering across both streams needs this mode.
func CaptureMerged(cmd *exec.Cmd) Result {
	out, err := cmd.CombinedOutput()
	exitCode, termination := outcome(cmd, err)
	return Result{ExitCode: exitCode, Termination: termination, Stdout: string(out), Err: err}
}

// outcome is the single translation from exec's process state to this package's
// result contract, so the capture modes cannot disagree about a process abort.
func outcome(cmd *exec.Cmd, err error) (int, Termination) {
	if err == nil {
		return 0, TerminationCompleted
	}
	if cmd.ProcessState == nil {
		return 1, TerminationSpawnFailed
	}
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode < 0 {
		return exitCode, TerminationSignaled
	}
	return exitCode, TerminationCompleted
}
