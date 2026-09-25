# Name bench worktree build in the binary-seal remedy

Blocked by: none
Writes: internal/freshness/freshness_verify.go, internal/preflight/gather.go, internal/preflight/decision.go, internal/preflight/decision_test.go
Covers: none

## What to build

When `bench preflight build` finds a stale `dist/bench`, the `binary-seal` row is red. The
row carries the refusal of the seal verifier whole, and that refusal names
`cd '<root>' && bash scripts/go-build.sh …` as the rebuild. In an assignment worktree the
root is under the pool path, and the hooks refuse a `cd` into the pool path. Thus the
operator cannot run the named remedy.

The seal verifier returns a typed refusal. The refusal keeps its current text, and it also
gives its reason without the rebuild command. When an active assignment owns the root, the
`binary-seal` row states that reason. Its `next` cell names `bench worktree build` with the
assignment id. The grammar of that verb comes from its one usage owner.

A root that no active assignment owns keeps the whole refusal. The hooks allow a `cd` into a checkout
outside the pool path.

## Acceptance

- [ ] In an assignment worktree, the red `binary-seal` row names `bench worktree build <assignment>` as its remedy, and no cell of the row names a `cd` or `scripts/go-build.sh`. `bench test --package ./internal/preflight` fails on the base and passes after the edit.
- [ ] In a root that no assignment owns, the red `binary-seal` row keeps the whole refusal of the seal verifier.
- [ ] A probe that swaps the remedy back to the refusal of the verifier turns the preflight test red, and the restore returns it to green.
