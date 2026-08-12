# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `fd6ac31`, 1 dirty path, 14 unpushed commits
Spec: `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `2bee65b` — stale, work tree `8f172d6`

## State

`axi-coherent-diff` is implemented and retired. This maintenance batch drains
its retro into FT113's final-landing CLI contract, FT162's exact repair-review
subject, and FT164's repair/probe guidance. FT173 now names only the staged
`axi-query-disclosure` capstone and is first in the refreshed sequence. There
are no parked ideas or open learnings. The coherent-diff implementation commits
remain unpushed; no foreign worktree or branch was cleaned.

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
