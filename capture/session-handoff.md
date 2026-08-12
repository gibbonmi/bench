# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2d1bee6`, clean tree, 7 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `a126a86` — stale, work tree `40490c7`

## State

The `axi-coherent-diff` implementation breakdown is reviewer-approved and landed.
The staged spec's `CHANGELOG.md` ownership-fence amendment landed at `49bc1895`;
the single atomic tracer ticket
`specs/axi-coherent-diff/tickets/render-one-coherent-diff-snapshot.md` landed at
`2d1bee6b`. It owns CD1-CD9 because a thinner landing would expose a partial
public response before the paired compatibility oracle and every-response
`help[]` contract are complete. Terra then Fable each ran exactly one read-only
ticket review. Their accepted repairs preserve existing `git.Facts` callers via
an additive all-files facts path and close detached-HEAD, deep-cwd, stable-rerun,
and base-equals-HEAD cases. `bench preflight build axi-coherent-diff` and
`bench coverage --check specs/axi-coherent-diff/spec.md` are green. The next
phase is one fresh isolated write-delegate charge; no implementation edit has
started.

## Next command

`$bench-implement-spec specs/axi-coherent-diff/spec.md --full`

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
