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
	"github.com/gibbonmi/bench/internal/shellcommand"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var worktreeExecGrammar = usage.Grammar{
	Cmd: "bench worktree exec",
	// The help carries the stdin form and the exit-2 rule because neither reads off the
	// grammar line. A child inherits this process's stdin, so a heredoc is the shape an
	// agent reaches for; and the child's own exit 2 looks like a grammar refusal until
	// the reader knows which prefix marks the refusal.
	Help: "usage: " + usage.WorktreeExec + "\n" +
		"stdin: bench worktree exec <target> -- python3 - <<'EOF'\n" +
		`exit 2: a line starting "usage: bench worktree exec" is a grammar refusal; any other exit 2 is the child's own`,
	Flags:                               []usage.Flag{{Name: "--env", HasValue: true, NoEmptyValue: true, Repeatable: true}},
	MinArgs:                             2,
	MaxArgs:                             -1,
	ReservedPositionalsBeforeTerminator: 1,
	ChildArgvAfterTerminator:            true,
}

// childEnvValues returns the --env assignments in argv order, or the usage line for the
// first malformed one. The check runs before the target resolves, so a typo names itself
// and no child starts on a half-read environment.
func childEnvValues(parsed usage.Result) ([]string, string) {
	values := parsed.Repeated["--env"]
	for _, value := range values {
		if !shellcommand.IsAssignment(value) {
			return nil, toon.Usage(worktreeExecGrammar.Cmd, value)
		}
	}
	return values, ""
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
	values, refusal := childEnvValues(parsed)
	if refusal != "" {
		fmt.Fprintln(stderr, refusal)
		return 2
	}
	// The record opens after the grammar answers. The verb resolves its assignment inside
	// the span, so a target refusal records the verb with no subject.
	var assignment string
	finishSpan := beginVerbSpan(home, root, otelExecSeam)
	exit := execAttributed(&assignment, parsed, root, home, values, stdin, stdout, stderr)
	finishSpan(exit, assignment)
	return exit
}

// execAttributed is the exec verb's own work, with the assignment that owns the child's
// tree written to assignment. The ledger is read once, through the one target authority,
// so the record names the same assignment the child ran in.
func execAttributed(assignment *string, parsed usage.Result, root, home string, values []string, stdin io.Reader, stdout, stderr io.Writer) int {
	selected, err := resolveAssignment(root, parsed.Positionals[0])
	if err != nil {
		return printTargetRefusal(stderr, "bench worktree exec", err)
	}
	*assignment = selected.ID
	return execSpanExit(runWorktreeChild(parsed.Positionals[1:], selected.Worktree, home, values, stdin, stdout, stderr))
}

// execSpanExit maps the child's exit onto the outcome derivation the worktree verbs
// share. The exit this verb returns is the child's own, and a 3 there is the child
// speaking rather than a Bench publication, so only a zero reads green.
func execSpanExit(code int) int {
	if code == 0 {
		return 0
	}
	return 1
}

func runWorktreeChild(argv []string, dir, home string, extraEnv []string, stdin io.Reader, stdout, stderr io.Writer) int {
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = dir, stdin, stdout, stderr
	cmd.Env = execEnv(dir, home, extraEnv)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(stderr, "bench worktree exec: %v\n", err)
		return nameWorktree(stderr, dir, 1)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return nameWorktree(stderr, dir, childExitCode(cmd, err))
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return nameWorktree(stderr, dir, 130)
	}
}

// nameWorktree returns code, and names the tree the child ran in when that code is a
// failure. Every exit path returns through here, so the line prints one time for one
// child run and never on a green run. It follows the child's own stderr, which the verb
// passes through unchanged, so the reader gets the failure first and then the path the
// recovery command needs.
func nameWorktree(stderr io.Writer, dir string, code int) int {
	if code == 0 {
		return code
	}
	fmt.Fprintf(stderr, "worktree: %s\n", dir)
	return code
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
//
// The order is the caller's environment, then the --env values, then the values this
// verb owns. A named routing variable therefore meets the same strip an inherited one
// meets, and --env cannot repoint the child's pool or its kit.
func execEnv(dir, home string, extra []string) []string {
	base := withHome(env.WithoutWrapperRouting(append(os.Environ(), extra...), runbinary.Env), home)
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
