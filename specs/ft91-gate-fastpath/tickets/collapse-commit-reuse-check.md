# Collapse bench commit's reuse check into the gate home

Blocked by: Plumb bench gate --fresh through the wrapper

## What to build

Story 11 of `specs/ft91-gate-fastpath/spec.md`: `bench commit`'s private
pre-gate `ReusableGreen` check (`internal/commit/commit.go:118`) collapses
into the single gate-package home — commit calls the gate unconditionally, the
reuse line is printed by the gate emitter, and verdict-reuse policy has
exactly one source. Observable behavior unchanged: reuse on fresh green,
refusal on red, exactly one gate run tallied. The runtime commit contract's
reuse-line expectation moves to the new emitter in the same change.

The blocker is a file fence, not logic: this ticket and its blocker both edit
`internal/gate/gate.go` and the runtime contract suite.

## Acceptance

- [ ] `bench commit` reuses a fresh green: exactly one gate run tallied across
      gate-then-commit (`runtime_commit_test.go` tally contract).
- [ ] `bench commit` still refuses on red through the collapsed path.
- [ ] The reuse line the operator reads survives, now emitted by the gate
      home; no `ReusableGreen` consultation remains in `internal/commit`.
