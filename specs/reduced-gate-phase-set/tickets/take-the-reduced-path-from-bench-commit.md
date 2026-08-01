# Take the reduced path from bench commit

Blocked by: Reduce the gate run and inherit from a full-green ancestor

Ownership fence: `internal/commit`, `internal/contract/runtime/runtime_commit_reduced_test.go`
Assumptions: the reduced execution path and its announcement already exist, so this ticket owns only route selection from the staged set

## What to build

`bench commit` takes the reduced path automatically when its staged set is confined
to the allowlist. No new flag: an opt-in fast path is one an agent under time
pressure reaches for in cases the rule does not cover, and the changeset already
carries every fact the decision needs. `bench gate --fresh` stays the escape to a
real whole-tree run.

Confinement is *all*, not *any*. A staged set mixing an allowlisted path with an
unlisted one runs the full gate — that is the most likely coding error here and the
one whose consequence is an ungraded code change riding along with a doc edit.

## Acceptance

- [ ] [R22] An allowlist-confined staged set commits through the reduced path, observable in the announced phase list.
- [ ] [R23] A staged set mixing an allowlisted and an unlisted path runs the full gate.
