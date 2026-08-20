# Warn before removing the binary answering bench

Blocked by: none
Writes: internal/worktree

## What to build

The residue guard treats an ignored `dist/bench` as removable, and removing it
in a self-hosting checkout disables the CLI, the git guard that CLI backs, and
the `BENCH_RUN_BINARY` the gate requires. Before removing a `dist/bench` the
guard warns, and the warning names the `bash scripts/go-build.sh` invocation
that rebuilds it rather than plain `go build`. The predicate
is identity rather than path: the guard compares the candidate against the
binary the wrapper's own resolution selects, so residue in an unrelated checkout
is removed without the warning.

## Acceptance

- [ ] a `dist/bench` matching the resolved binary produces the warning before removal, naming the `scripts/go-build.sh` invocation (H25).
- [ ] a `dist/bench` that is not the resolved binary is removed with no warning (H26).
