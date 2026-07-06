// Package subprocess is the one seam for running an external command and
// capturing its outcome. Conformance probes, the test harness, and the canary
// runner all cross it instead of hand-rolling capture: the exit-code
// derivation lives here once, and the two capture modes differ only in whether
// stdout and stderr stay separate or interleave.
package subprocess

import (
	"bytes"
	"os/exec"
)

// Result is the outcome of running a command. ExitCode is 0 on success, the
// process's own code when it ran and failed, or 1 when it never started. Err is
// the raw error from the run. Capture fills Stdout and Stderr separately;
// CaptureMerged interleaves both into Stdout and leaves Stderr empty.
type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

// Capture runs cmd with stdout and stderr buffered separately. Callers that
// parse stdout (JSON, machine output) need this so subprocess stderr chatter
// cannot corrupt the parse.
func Capture(cmd *exec.Cmd) Result {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return Result{ExitCode: exitCode(cmd, err), Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
}

// CaptureMerged runs cmd with stdout and stderr interleaved into one stream, in
// the order a terminal would show them. EXPECT-style matching that asserts on
// ordering across both streams needs this mode.
func CaptureMerged(cmd *exec.Cmd) Result {
	out, err := cmd.CombinedOutput()
	return Result{ExitCode: exitCode(cmd, err), Stdout: string(out), Err: err}
}

// exitCode derives a subprocess exit code from a completed command: 0 when it
// succeeded, the process's own code when it ran and failed, and 1 when it never
// started (a spawn failure, where ProcessState is nil).
func exitCode(cmd *exec.Cmd, err error) int {
	if err == nil {
		return 0
	}
	if cmd.ProcessState != nil {
		return cmd.ProcessState.ExitCode()
	}
	return 1
}
