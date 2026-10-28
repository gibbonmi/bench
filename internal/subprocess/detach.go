package subprocess

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// detachEnv marks a detach helper: a run of this executable that starts the program of its
// arguments in a new session and exits.
const detachEnv = "BENCH_SUBPROCESS_DETACH"

// Detached is the start of one program in a new session. Cmd is a detach helper, which is this
// executable again, and a caller runs it under RunInherited with Streams. The helper starts the
// program with Setsid, gives it the null device for each standard stream and no inherited
// protection vector, writes its pid, and exits. So the program is in no group and in no session
// of the caller, it outlives the caller, and its parent is not the caller, so the caller keeps
// no zombie. The program keeps the directory and the environment of the command it came from.
type Detached struct {
	Cmd     *exec.Cmd
	Streams Streams
	program string
}

// Detach answers the detached start of the program of cmd.
func Detach(cmd *exec.Cmd) (Detached, error) {
	if cmd.Err != nil {
		return Detached{}, cmd.Err
	}
	self, err := gateExecutable()
	if err != nil {
		return Detached{}, err
	}
	helper := exec.Command(self, append([]string{cmd.Path}, cmd.Args...)...)
	helper.Dir, helper.Env = cmd.Dir, append(cmd.Environ(), detachEnv+"=1")
	return Detached{Cmd: helper, Streams: Streams{Stdout: Output{Capture: true}, Stderr: Output{Capture: true}}, program: cmd.Path}, nil
}

// PID answers the pid of the program from the outcome and the error of the helper's run.
func (d Detached) PID(outcome Outcome, err error) (int, error) {
	switch {
	case err != nil:
		return 0, err
	case outcome.Process.StartErr != nil:
		return 0, outcome.Process.StartErr
	case outcome.Process.Code != 0:
		return 0, fmt.Errorf("start %s: %s", d.program, strings.TrimSpace(string(outcome.Streams.Stderr.Data)))
	}
	if err := outcome.Uncertified(); err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(outcome.Streams.Stdout.Data)))
	if err != nil || pid <= 1 {
		return 0, errors.New("subprocess: the detach helper named no pid")
	}
	return pid, nil
}

// runDetach is the detach helper. Its arguments are the program path and then the program
// argv. It answers 1 when the program does not start.
func runDetach() int {
	if len(os.Args) < 3 {
		return 1
	}
	env := slices.DeleteFunc(os.Environ(), func(entry string) bool {
		return strings.HasPrefix(entry, detachEnv+"=") || strings.HasPrefix(entry, protection.Env+"=")
	})
	program := &exec.Cmd{Path: os.Args[1], Args: os.Args[2:], Env: env, SysProcAttr: &syscall.SysProcAttr{Setsid: true}}
	if err := program.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println(program.Process.Pid)
	return 0
}
