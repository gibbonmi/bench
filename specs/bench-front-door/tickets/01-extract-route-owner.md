# Extract the route owner from handoff into the status package

Blocked by: none
Writes: internal/status, internal/handoff

## What to build

Move the selection that turns the severity ladder into "the next command" — the
invocable predicate, the harness prefix table and its `translate`, the
first-invocable walk, and the next-action derivation — out of the handoff package
into the status package (which owns the ladder and is already imported by handoff).
Handoff reads them from there and keeps no copy; `bench handoff --harness` keeps
validating against the same table. Add the `git ` prefix to the predicate and make it
reject trailing prose after a command. Behavior of `bench handoff` is unchanged on
every existing fixture. Contract with 04: the owner exposes one function that returns
the lead (signal, detail, command), the runners-up, and the no-command flag for a
given harness, so `--route` renders and never re-selects.

Covers: R10, R25, R36

## Acceptance

- [ ] `TestFirstInvocable` (moved) passes with its existing cases plus: rejects `bench gate --fresh for a whole-tree verdict`, accepts `git push`; a table case renders `/bench-drain` for claude and `$bench-drain` for codex through the owner.
- [ ] `bench handoff` output is byte-identical to before on the existing render fixtures.
- [ ] handoff's no-invocable state is read from the owner, not recomputed.
- [ ] Gate green.
