# Own the run binary in an exec child

Blocked by: none
Writes: internal/worktree/exec.go, internal/worktree/exec_test.go

## What to build

`bench worktree exec` strips the caller's kit, wrapper, and run-binary variables
from its child, correctly — the caller's selection was built for the caller's
run. But the gate's owner lookup reads the absence of both routing variables as
"nobody owns this run" and returns nothing, so the gate plan's environment never
gains a selection and the gate entry refuses. A child invoked as the profile
directs, through the worktree's own `./dist/bench`, cannot reach a gate.

After the strip, point the run-binary variable at the child worktree's own
`dist/bench` when a regular file is present there. The worktree path arrives
already absolute and cleaned — the assignment ledger refuses one that is not —
so the artifact path inherits both properties rather than re-deriving them.

Exec's added knowledge is exactly one predicate: a regular file exists at the
artifact path. Every trust condition stays where it already lives, in the gate
entry: absolute, executable, not a symbolic link, physical path equal to the
given path, and fresh against the kit. Do not restate any of them here. A live
symbolic link at the artifact path is therefore emitted, not filtered, and the
gate entry refuses it — that split is the point, and the acceptance rows pin it.
A component that names a path cannot authenticate that path, so the validator
that launches the binary stays the independent one.

When no regular file sits at the artifact path — absent, a directory, a FIFO, a
device, a socket, or a dangling symbolic link — leave the variable unset and the
child behaves as it does today.

This ticket does not stop every abandonment of `bench worktree exec`. A child
can still meet the inherited-verdict refusal that FT223 tracks, which also
misreports a partial verdict as a failure. FT223 owns that trigger and its own
competing fixes; this ticket removes the other one.

Evidence comes from running a real child through the exec path and reading the
assignments it received, plus one composed run against a real worktree.

## Acceptance

- [ ] A child in a worktree holding a regular-file artifact receives the run-binary variable naming that worktree's artifact path.
- [ ] The emitted value is an absolute path.
- [ ] A child never receives the caller's run-binary value.
- [ ] A child receives neither the caller's kit nor the caller's wrapper variable.
- [ ] A child receives every unrelated variable the caller set.
- [ ] A child in a worktree with no artifact path receives no run-binary variable.
- [ ] A child in a worktree whose artifact path is a directory receives no run-binary variable.
- [ ] A child in a worktree whose artifact path is a FIFO receives no run-binary variable.
- [ ] A child in a worktree whose artifact path is a dangling symbolic link receives no run-binary variable.
- [ ] A child in a worktree whose artifact path is a live symbolic link to a regular file receives the variable, and the gate entry refuses that value.
- [ ] A child's environment differs from today's by the run-binary assignment alone.
- [ ] `bench worktree exec <target> -- ./dist/bench gate` over a clean worktree of this kit reports green.
