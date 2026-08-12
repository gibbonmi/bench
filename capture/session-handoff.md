# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `543d855`, 1 dirty path, 34 unpushed commits
Spec: `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `87299b5` — stale, work tree `a10f583`

## State

The repaired implementation was re-reviewed on exact candidate
`c2542a97..543d8551`: Standards 3, Spec 2, Coverage 3, collapsing to six
deterministic repair targets recorded in `reviews/axi-query-disclosure.md`.
The formerly accepted seven targets are resolved. Keep the approved surface,
default-extraction scope, exact compatibility contract, production-registry
ownership, and the reviewer-selected call-site provenance decision closed.
No push has run. The eight pre-existing foreign worktree records remain untouched.

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
