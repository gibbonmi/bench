# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `b90d8a4`, clean tree, 29 unpushed commits
Spec: `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `4440ce1` — stale, work tree `210c21c`

## State

Implementation review is recorded at `reviews/axi-query-disclosure.md` in
`b90d8a4c`: Standards 2, Spec 5, Coverage 2, collapsing to seven repair targets.
Six are deterministic auto-fixes. One reviewer decision remains: remove the
test-fitted literal `unknown` rejection and move guessed-value falsification to
call-site derivation tests (recommended), or redesign the public action type with
a new enforceable provenance contract. Keep the approved query surface,
compatibility deltas, production-registry ownership, and Fable findings closed.
No repair has started, no push has run, and the eight pre-existing foreign
worktree records remain untouched.

## Next command

`$bench-implement-spec specs/axi-query-disclosure/spec.md --full`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
