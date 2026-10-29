package worktree

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"syscall"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/shellcommand"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var worktreeExecGrammar = usage.Grammar{
	Cmd: usage.CommandName(usage.WorktreeExec),
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
	exit, _ := ExecCommandResolving(root, home, args, stdin, stdout, stderr)
	return exit
}

// ExecCommandResolving is ExecCommand that also answers the assignment its target resolved
// to, or the empty string when no target resolved. The dispatcher keys the census output
// record by this assignment, so it runs no second resolution of the target.
func ExecCommandResolving(root, home string, args []string, stdin io.Reader, stdout, stderr io.Writer) (int, string) {
	parsed, line, code := usage.Parse(worktreeExecGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code, ""
	}
	if !parsed.EndedFlags || parsed.PositionalsBeforeTerminator != 1 || len(parsed.Positionals) < 2 {
		fmt.Fprintln(stderr, worktreeExecGrammar.Help)
		return 2, ""
	}
	values, refusal := childEnvValues(parsed)
	if refusal != "" {
		fmt.Fprintln(stderr, refusal)
		return 2, ""
	}
	// The record opens after the grammar answers. The verb resolves its assignment inside
	// the span, so a target refusal records the verb with no subject.
	var assignment string
	finishSpan := beginVerbSpan(home, root, otelExecSeam)
	exit := execAttributed(&assignment, parsed, root, home, newAmbient(home, stderr).environ, values, stdin, stdout, stderr)
	// The verb returns the child's own code, which the grammar help promises. The span
	// takes the outcome instead: a 3 here is the child speaking rather than a Bench
	// publication, so only a zero reads green.
	spanExit := exit
	if spanExit != 0 {
		spanExit = 1
	}
	finishSpan(spanExit, assignment)
	return exit, assignment
}

// AssignmentActive reports whether the ledger at root holds the assignment id in the
// active state. It reads the ledger and applies the landing's active check, and it
// resolves no target. A ledger that cannot be read holds no active assignment.
func AssignmentActive(root, id string) bool {
	a, err := assignmentByID(nil, root, id)
	return err == nil && landingActiveState(a.State)
}

// execAttributed is the exec verb's own work, with the assignment that owns the child's
// tree written to assignment. The ledger is read once, through the one target authority,
// so the record names the same assignment the child ran in.
func execAttributed(assignment *string, parsed usage.Result, root, home string, environ, values []string, stdin io.Reader, stdout, stderr io.Writer) int {
	selected, err := resolveAssignment(root, parsed.Positionals[0])
	if err != nil {
		return printTargetRefusal(stderr, worktreeExecGrammar.Cmd, err)
	}
	*assignment = selected.ID
	return runWorktreeChild(parsed.Positionals[1:], selected.Worktree, home, environ, values, stdin, stdout, stderr)
}

func runWorktreeChild(argv []string, dir, home string, environ, extraEnv []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return nameWorktree(stderr, dir, runChild(worktreeExecGrammar.Cmd, argv, dir, home, environ, extraEnv, stdin, stdout, stderr))
}

// RunTreeChild runs argv in dir as the one child of a tree-scoped call, and answers the
// child's exit. The child takes the exec child environment. It prints no worktree line,
// because the child names its own tree in its response. verb names the call in a start
// failure.
func RunTreeChild(verb string, argv []string, dir, home string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runChild(verb, argv, dir, home, baseEnviron(), nil, stdin, stdout, stderr)
}

// execPolicy keeps the child's group after a normal exit and waits ExecWaitDelay for
// output that a descendant still holds open. An interruption sends INT, and KILL follows
// the same delay.
var execPolicy = subprocess.Policy{
	Normal:      subprocess.Preserve,
	Grace:       bounds.ProcessWindow(bounds.FixedWindow(bounds.ExecWaitDelay)),
	Final:       bounds.ProcessWindow(bounds.FixedWindow(bounds.ProcessFinalWindow)),
	OutputDelay: bounds.ProcessWindow(bounds.FixedWindow(bounds.ExecWaitDelay)),
	Poll:        bounds.ProcessPoll,
}

// runChild is the one child runner of the exec verb and the tree-scoped call. It answers
// the child's exit, 1 for a start failure or a failed stream or unfinished cleanup, and 130
// for an interrupt. A stream that the output delay forcibly closed is reported on stderr
// and keeps the child's exit. The child registers on its resource vector before it starts. A normal
// return releases no resource: a registered user that is still unresolved keeps its
// protection, and the enclosing owner observes it.
func runChild(verb string, argv []string, dir, home string, environ, extraEnv []string, stdin io.Reader, stdout, stderr io.Writer) int {
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	resources, err := childResources(dir, environ)
	if err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", verb, err)
		return 1
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	// The resource vector is a value the verb owns, so it lands after every --env value.
	cmd.Env = protection.Environ(execEnv(dir, home, environ, extraEnv), resources)
	streams := subprocess.Streams{Stdin: stdin, Stdout: subprocess.Output{Writer: stdout}, Stderr: subprocess.Output{Writer: stderr}}
	lifetime := subprocess.NewLifetime(cmd, execPolicy, streams, subprocess.Registration{Protocol: protection.Installed(), Resources: resources})
	if err := lifetime.Start(); err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", verb, err)
		return 1
	}
	cancel := context.AfterFunc(ctx, func() { lifetime.Cancel(subprocess.CauseInterruption, syscall.SIGINT) })
	defer cancel()
	outcome, _ := lifetime.Wait()
	switch {
	case outcome.Cancellation.Cause == subprocess.CauseInterruption:
		return 130
	case outcome.Process.StartErr != nil:
		fmt.Fprintf(stderr, "%s: %v\n", verb, outcome.Process.StartErr)
		return 1
	}
	code := lifetimeExitCode(outcome.Process)
	forced, unfinished := childUnfinished(outcome)
	if forced != nil {
		fmt.Fprintf(stderr, "%s: %v\n", verb, forced)
	}
	if code == 0 && unfinished != nil {
		fmt.Fprintf(stderr, "%s: %v\n", verb, unfinished)
		return 1
	}
	return code
}

// childUnfinished splits the child's incomplete state in two. forced names an owned
// stream that the output delay closed: exec reports it and still returns the child's code,
// so a descendant that holds the pipe cannot turn a child zero into a failure (BO69).
// unfinished names any other stream or cleanup obligation that did not complete. It
// leaves protection out: the enclosing owner consumes that disposition.
func childUnfinished(o subprocess.Outcome) (forced, unfinished error) {
	o.Protection = protection.Disposition{State: protection.Safe}
	for _, stream := range []*subprocess.Stream{&o.Streams.Stdin, &o.Streams.Stdout, &o.Streams.Stderr} {
		if stream.Status != subprocess.StreamClosed {
			continue
		}
		if forced == nil {
			// The outcome's own refusal names the stream, so the line keeps one format.
			none := subprocess.Stream{Status: subprocess.StreamNone}
			forced = subprocess.Outcome{Streams: subprocess.StreamSet{Stdin: *stream, Stdout: none, Stderr: none}, Protection: o.Protection}.Uncertified()
		}
		stream.Status = subprocess.StreamComplete
	}
	return forced, o.Uncertified()
}

// nameWorktree returns code, and names the tree the child ran in when that code is a
// failure. The exec verb returns every child exit through here, so the line prints one
// time for one child run and never on a green run. It follows the child's own stderr,
// which the verb passes through unchanged, so the reader gets the failure first and then
// the path the recovery command needs.
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
// caller's run, not this child's. PWD comes from the child's directory.
// Every unrelated variable the operator set stays.
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
// The order is the call's base environment, then the --env values, then the values this
// verb owns. A named routing variable therefore meets the same strip an inherited one
// meets, and --env cannot repoint the child's pool or its kit.
func execEnv(dir, home string, environ, extra []string) []string {
	base := withHome(env.WithoutWrapperRouting(slices.Concat(environ, extra), runbinary.Env, "PWD"), home)
	base = append(base, "PWD="+dir)
	wrapper := wrapperPath(dir)
	if !isRegularFile(wrapper) {
		return base
	}
	return append(base, env.WrapperEnv+"="+wrapper)
}

// withHome puts the caller's resolved home on a child environment. The inherited
// assignment is dropped first, so the child reads exactly one value for the name.
func withHome(base []string, home string) []string {
	return append(capability.WithoutEnvironment(base, homeEnv), homeEnv+"="+home)
}

// lifetimeExitCode is the shell exit code of a Lifetime child. A child whose one waiter never
// reported has no code of its own.
func lifetimeExitCode(p subprocess.Process) int {
	code := p.Code
	if !p.Waited {
		code = -1
	}
	return shellExitCode(code, p.Signal)
}

func childExitCode(cmd *exec.Cmd, err error) int {
	if err == nil {
		return 0
	}
	if cmd.ProcessState == nil {
		return 1
	}
	var signal syscall.Signal
	if status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok && status.Signaled() {
		signal = status.Signal()
	}
	return shellExitCode(cmd.ProcessState.ExitCode(), signal)
}

// shellExitCode is the one map from a child's end to the code that a worktree verb returns:
// 128 plus the signal that ended the child, else its own code, or 1 when it has no code.
func shellExitCode(code int, signal syscall.Signal) int {
	switch {
	case signal != 0:
		return 128 + int(signal)
	case code < 0:
		return 1
	}
	return code
}
