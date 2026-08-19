# Stop leaking wrapper routing into a worktree child

Blocked by: none
Writes: internal/worktree/exec.go, internal/worktree/exec_test.go

## What to build

`bench worktree exec` runs its child in a different tree than the wrapper that
invoked it, so the invoking wrapper's routing internals must not reach that
child. Today `bin/bench.sh` exports `BENCH_KIT`, `BENCH_WRAPPER`, and
`BENCH_RUN_BINARY` before it execs the binary, and the exec child inherits the
whole environment. Inside a worktree of the kit itself the leaked `BENCH_KIT`
names the main checkout, so the gate's `sameDirectory(root, kit)` test fails:
the run silently drops the `race` and `system` phases and the
`BENCH_CONFORMANCE_ROOT` staging, then reports `gate: red` on the
environment-skip rule over a clean tree. The operator reads a red that names a
missing environment nobody asked them to set.

Strip those three variables from the exec child's environment, so the child's
own wrapper resolves its own kit — the same removal and the same reason as
`gateEnv` in `internal/gate/gate.go`, which already does this for the gate's
child and documents why. Everything else the operator set stays.

This closes the leak on the worktree-native invocation
(`bench worktree exec "<label>" -- bash bin/bench.sh gate`). A child invoked as
a bare `bench` still resolves the wrapper found on `PATH`, whose `kit_dir` is
that wrapper's own tree — a separate question about wrapper kit resolution, out
of scope here.

## Acceptance

- [ ] A child run by `bench worktree exec` sees no `BENCH_KIT`, `BENCH_WRAPPER`, or `BENCH_RUN_BINARY`, even when the caller's environment sets all three.
- [ ] A child run by `bench worktree exec` still sees every other variable the caller's environment carried.
- [ ] `bench worktree exec "<label>" -- bash bin/bench.sh gate --fresh` over a clean worktree of this kit runs the `race` and `system` phases and reports `gate: green` with `environment=0`.
