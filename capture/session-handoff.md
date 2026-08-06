# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `4523c7f`, one tracked handoff edit, 23 unpushed commits
Spec: `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: the atomic ticket-closure commit passed the full gate in its exact prospective worktree; the main-worktree marker is stale after fast-forward

## State

**Phase reached: FT187 implementation composed provisionally; recomposition pending.**

The spec-build lifecycle is active on subject `52401e58eb49ffd29394452e0cef2c03ffb74456`.
All FT187 assignments and review repairs are integrated and released. Main advanced to
`4523c7f3eba36df9950a9f656aa0d608c8ced413` for the atomic ticket-closure checker and
FT187 ticket migration, so promotion must first recompose the provisional candidate.
No final review evidence is retained.

Closed decisions stay closed: uniform prose anchors compose with FT156's
`internal/anchors` registry; the structured-phase grammar remains bespoke; the prose
cut carries five fixed reader cases and six measured shrinking passages. The required
cross-harness semantic review remains `opus / low` and runs after recomposition. The
assignment-time closure checker is implemented and full-gate green; it is not deferred
to a later ticket.

## Next command

`GOMAXPROCS=2 GOFLAGS=-p=2 bench spec build promote ft187-communication-surface-cut`

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
