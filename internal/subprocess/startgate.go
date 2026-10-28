package subprocess

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

// gateEnv names the two inherited descriptors of a start gate: the open pipe that the gate
// reads and the exec-status pipe that it writes.
const gateEnv = "BENCH_SUBPROCESS_START_GATE"

// init makes a gated launch its start gate, and a detach helper its helper, before main or
// TestMain runs. Every actor that registers a launch links this package, so its own
// executable can serve as its gate and as its helper. A gate execs the helper with no gate
// variable, so the gate test comes first.
func init() {
	if fds, ok := os.LookupEnv(gateEnv); ok {
		os.Exit(runGate(fds))
	}
	if _, ok := os.LookupEnv(detachEnv); ok {
		os.Exit(runDetach())
	}
}

// runGate waits for the open signal and then executes the program in its own process and
// group. At end of file it exits and executes nothing. A failed exec writes its errno to the
// close-on-exec status pipe, so the actor tells it from a program exit.
func runGate(fds string) int {
	open, status, ok := strings.Cut(fds, ",")
	openFD, openErr := strconv.Atoi(open)
	statusFD, statusErr := strconv.Atoi(status)
	if !ok || openErr != nil || statusErr != nil || len(os.Args) < 3 {
		return 1
	}
	syscall.CloseOnExec(statusFD)
	signal := os.NewFile(uintptr(openFD), "start-gate-open")
	if n, _ := signal.Read(make([]byte, 1)); n != 1 {
		return 1
	}
	_ = signal.Close()
	env := slices.DeleteFunc(os.Environ(), func(entry string) bool { return strings.HasPrefix(entry, gateEnv+"=") })
	err := syscall.Exec(os.Args[1], os.Args[2:], env)
	var errno syscall.Errno
	if !errors.As(err, &errno) {
		errno = syscall.EINVAL
	}
	_, _ = syscall.Write(statusFD, []byte(strconv.Itoa(int(errno))))
	return 127
}

// startGate holds the actor's ends of one gated launch.
type startGate struct {
	program      string
	open, status *os.File
}

// startGated starts cmd's program behind a start gate in the program's new group, and it
// restores cmd's own fields after Start. A start error names the program, as it would
// without the gate.
func startGated(cmd *exec.Cmd, start func(*exec.Cmd) error) (*startGate, error) {
	self, err := gateExecutable()
	if err != nil {
		return nil, err
	}
	gateRead, open, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	status, statusWrite, err := os.Pipe()
	if err != nil {
		_, _ = gateRead.Close(), open.Close()
		return nil, err
	}
	path, args, files, env := cmd.Path, cmd.Args, cmd.ExtraFiles, cmd.Env
	argv := args
	if len(argv) == 0 {
		argv = []string{path}
	}
	fd := 3 + len(files)
	cmd.Path, cmd.Args = self, append([]string{path, path}, argv...)
	cmd.ExtraFiles = append(slices.Clip(files), gateRead, statusWrite)
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=%d,%d", gateEnv, fd, fd+1))
	err = start(cmd)
	cmd.Path, cmd.Args, cmd.ExtraFiles, cmd.Env = path, args, files, env
	_, _ = gateRead.Close(), statusWrite.Close()
	if err != nil {
		_, _ = open.Close(), status.Close()
		if pathErr := (*fs.PathError)(nil); errors.As(err, &pathErr) && pathErr.Path == self {
			pathErr.Path = path
		}
		return nil, err
	}
	return &startGate{program: path, open: open, status: status}, nil
}

// release sends the open signal and waits for the exec. A failed exec returns the start
// error that Start returns for the program.
func (g *startGate) release() error {
	_, _ = g.open.Write([]byte{1})
	data, _ := io.ReadAll(g.status)
	g.shut()
	if len(data) == 0 {
		return nil
	}
	errno, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("subprocess: start gate status %q: %w", data, err)
	}
	return &fs.PathError{Op: "fork/exec", Path: g.program, Err: syscall.Errno(errno)}
}

// shut closes the actor's ends. A gate that never received the open signal reads end of file.
func (g *startGate) shut() {
	if g == nil {
		return
	}
	for _, f := range []**os.File{&g.open, &g.status} {
		if *f != nil {
			_ = (*f).Close()
			*f = nil
		}
	}
}

// gateExecutable names the running executable. On Linux the process link stays valid after
// a rebuild replaces or removes the file.
func gateExecutable() (string, error) {
	if runtime.GOOS == "linux" {
		return "/proc/self/exe", nil
	}
	return os.Executable()
}
