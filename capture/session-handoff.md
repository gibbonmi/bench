# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `329a03ed`, clean, 5 unpushed commits
Spec: none staged.
Gate: green at `329a03ed` — current exact tree

## State

The 2026-08 deepening map lives at `decisions/deepening-2026-08.md` (re-homed after
the gate spec's retire deleted it). Landed from it: candidates 1+3 (gate spec, retired),
6 (`internal/gate/greenmarker`), 5 (`internal/shift/objective.go`). Remaining lanes in
map order: skills spec (candidate 8, ticket #13), adopt spec (candidate 4, #9), worktree
spec (candidate 2, #8) plus one light-path ticket per recurring reader for candidate 7
(#12) — FT189 has landed, so nothing blocks them. FT207 still holds its unresolved
decision; FT185 remains the next ready non-deepening build. One open learnings entry
awaits `/bench-what-next` (map placement on retire).

## Next command

`/bench-write-spec` for candidate 8 from `decisions/deepening-2026-08.md` ticket #13
(one shipped skill-frontmatter reader; kit edit under `craft-synthesis`, ADR 0006/0012
constraints).

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
