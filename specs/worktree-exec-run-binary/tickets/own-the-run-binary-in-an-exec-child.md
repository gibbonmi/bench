# Own the run binary in an exec child

Blocked by: none
Writes: internal/worktree/exec.go, internal/worktree/exec_test.go

## What to build

`bench worktree exec` strips the caller's kit, wrapper, and run-binary variables
from its child, correctly — the caller's selection was built for the caller's
run. The gate's owner lookup then reads the absence of both routing variables as
"no one owns this run" and returns nothing, so the gate never selects a binary
and the gate entry refuses. A child invoked as the profile directs, through the
worktree's own `./dist/bench`, cannot reach a gate.

After the strip, point `BENCH_WRAPPER` at the child worktree's own wrapper when
a regular file sits at that path. The gate's owner lookup then owns the run and
builds one private binary from the child's own source, which is what the profile
requires of every gate and `bench test` run.

Do not set the run-binary variable instead. That makes the child inherit, and
the inherited path verifies the binary's own seal rather than checking it
against its source, so an exec child's `bench test` or `bench commit` could
reuse an artifact that had fallen behind the tree it grades.

Discovery is one predicate: a regular file exists at the wrapper path. Do not
inspect its content — an empty wrapper still sets the marker, because the
marker's only readers are an owner lookup testing for non-emptiness and a doctor
reporting a path. Nothing executes the value, so this ticket authenticates no
executable and adds no third spelling of the launch-trust predicate.

The worktree path arrives absolute and cleaned, canonicalized through symbolic
links at creation and validated by the assignment ledger, so the wrapper path
inherits both properties rather than re-deriving them.

When no regular file sits at the wrapper path — absent, a directory, a FIFO, a
device, a socket, or a dangling symbolic link — leave the variable unset and the
child behaves as it does today. A worktree of a linked project repository, whose
wrapper sits under the vendored kit rather than at the worktree's own wrapper
path, falls in that unset case by design.

This ticket does not stop every abandonment of `bench worktree exec`. A child
can still meet the inherited-verdict refusal that FT223 tracks, which also
misreports a partial verdict as a failure. FT223 owns that trigger and its own
competing fixes; this ticket removes the other one.

Evidence comes from running a real child through the exec path and reading the
assignments it received, plus one composed run against a real worktree recorded
as ticket evidence.

## Acceptance

- [ ] A child in a worktree holding a regular-file wrapper receives the wrapper variable naming that worktree's wrapper path (WX1).
- [ ] The emitted value is an absolute path (WX2).
- [ ] A child's gate selects an owned binary rather than an inherited one (WX3).
- [ ] A child never receives the caller's wrapper value (WX4).
- [ ] A child receives neither the caller's kit nor the caller's run-binary variable (WX5).
- [ ] A child receives every unrelated variable the caller set (WX6).
- [ ] A child in a worktree with no wrapper path receives no wrapper variable (WX7).
- [ ] A child in a worktree whose wrapper path is a directory receives no wrapper variable (WX8).
- [ ] A child in a worktree whose wrapper path is a FIFO receives no wrapper variable (WX9).
- [ ] A child in a worktree whose wrapper path is a dangling symbolic link receives no wrapper variable (WX10).
- [ ] A child in a worktree whose wrapper path is an empty regular file receives the wrapper variable (WX11).
- [ ] A child running a verb that already owned its binary is unaffected (WX12).
- [ ] A child's environment differs from today's by the wrapper assignment alone (WX13).
- [ ] A child in a worktree whose wrapper path is a live symbolic link to a regular file receives the wrapper variable naming the link path (WX21).
- [ ] A child in a worktree whose path contains a space and a glob character receives the exact wrapper path (WX22).
- [ ] `bench worktree exec <target> -- ./dist/bench gate` over a clean worktree of this kit reports green, recorded as ticket evidence (WX20).
