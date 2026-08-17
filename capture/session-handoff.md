# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `bdab7c3`, 6 dirty paths, 1 unpushed commit
Spec: none staged.
Gate: green at `ce01d2a` — stale, work tree `5537bc2`

## State

`$bench-deepen` reverified the current architecture at `main@1c86bdf9`,
wrote the visual report to
`capture/architecture-review-20260817T104714.html`, and committed it as
`bdab7c3b`. Three opportunities remain: FT216 worktree eligibility, FT217
adopt lifecycle planning, and FT218 named Git readers.

`decisions/deepening-2026-08.md` is now the one current ready map for those
three. FT216 and FT217 are spec-ready; FT218 takes one light-path ticket per
reader, beginning with the six-site private administration-directory fact.
The map preserves exact cleanup precedence, keeps adopt dry-run behavior out of
scope, and retains `git.Output`/`Raw`/`OK` as plumbing. The roadmap now
matches those closed decisions. `bench maps` must omit this map before pickup.

An open learning in `capture/learnings.md` asks `/bench-deepen` to refresh a
ready map, replace duplicate roadmap decision claims with pointers, and rewrite
the handoff automatically after a new survey proves the frontier is empty.

## Next command

`$bench-write-spec`

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
