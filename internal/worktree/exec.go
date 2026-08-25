package worktree

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/usage"
)

var worktreeExecGrammar = usage.Grammar{
	Cmd:                                 "bench worktree exec",
	Help:                                "usage: " + usage.WorktreeExec,
	MinArgs:                             2,
	MaxArgs:                             -1,
	ReservedPositionalsBeforeTerminator: 1,
}

// ExecCommand runs a direct child argv from one active Bench-owned worktree.
func ExecCommand(root, home string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(worktreeExecGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	if !parsed.EndedFlags || parsed.PositionalsBeforeTerminator != 1 || len(parsed.Positionals) < 2 {
		fmt.Fprintln(stderr, worktreeExecGrammar.Help)
		return 2
	}
	path, err := resolveWorktree(root, parsed.Positionals[0])
	if err != nil {
		return printTargetRefusal(stderr, "bench worktree exec", err)
	}
	return runWorktreeChild(parsed.Positionals[1:], path, home, stdin, stdout, stderr)
}

func runWorktreeChild(argv []string, dir, home string, stdin io.Reader, stdout, stderr io.Writer) int {
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = dir, stdin, stdout, stderr
	cmd.Env = execEnv(dir, home)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(stderr, "bench worktree exec: %v\n", err)
		return 1
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return childExitCode(cmd, err)
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return 130
	}
}

// execEnv is the environment a worktree child runs under, with that child owning its run.
// The child runs in a different tree than the wrapper that reached this call, so it has
// to resolve its own kit. An inherited BENCH_KIT would name the caller's checkout instead.
// The selected executable goes with them for the same reason — it was built for the
// caller's run, not this child's. Everything else the operator set stays.
//
// Stripping both routing variables would leave the child's gate with no owner for its
// run. It would then never select a binary, and it would refuse at the gate entry.
// BENCH_WRAPPER is the variable that already means "Bench rooted this run". Re-pointing it
// at the child's own wrapper makes the child's gate own the run. It builds one private
// binary from the tree it is about to grade.
//
// Naming the worktree's built executable instead would make the child inherit it. An
// inherited selection is verified against its own seal rather than against its source, so
// a stale artifact could grade the tree.
//
// Nothing here executes the value: the owner lookup tests it for non-emptiness, and the
// adoption doctor reports it as a path. exec authenticates no executable. The only
// predicate is that a regular file sits at the wrapper path. Content is not read: an
// empty wrapper is still the marker. dir arrives absolute and cleaned from the assignment
// ledger, so the joined path inherits both properties.
//
// BENCH_HOME comes from the verb's resolved home, not from the caller's process. A
// child that inherited the name would read a different pool than the verb that
// started it.
func execEnv(dir, home string) []string {
	base := withHome(env.WithoutWrapperRouting(os.Environ(), runbinary.Env), home)
	wrapper := filepath.Join(dir, "bin", "bench.sh")
	if !isRegularFile(wrapper) {
		return base
	}
	return append(base, "BENCH_WRAPPER="+wrapper)
}

// withHome puts the caller's resolved home on a child environment. The inherited
// assignment is dropped first, so the child reads exactly one value for the name.
func withHome(base []string, home string) []string {
	return append(capability.WithoutEnvironment(base, homeEnv), homeEnv+"="+home)
}

func childExitCode(cmd *exec.Cmd, err error) int {
	if err == nil {
		return 0
	}
	if cmd.ProcessState != nil {
		if exit := cmd.ProcessState.ExitCode(); exit >= 0 {
			return exit
		}
		if status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return 128 + int(status.Signal())
		}
	}
	return 1
}
