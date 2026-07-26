# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `38cc1f5`, main checkout clean, 8 unpushed commits
Spec: none staged — FT96 shipped and retired
Gate: green at `edc15e9`, the exact tree of HEAD

## State

- **FT96 is built, gate-green, and retired.** Seven guidance clauses in
  `.agents/skills/bench-craft-delegate/SKILL.md`; a second destructive-git deny
  class for the working-tree-mutating `git stash` verbs; the conformance
  dead-pointer sweep widened from `.agents/commands` to the whole `.agents`
  markdown tree; two presence anchors on craft-delegate plus three canary
  fixtures proving the new bite paths. All ten stories shipped, all 31 coverage
  rows classified.
- **Nothing is uncommitted and nothing is half-done.** The gate ran green on the
  exact tree of HEAD as one subject, after the retirement commit — not a cached
  verdict from an earlier tree.
- **The eight commits are unpushed, and the push is the reviewer's.** Nothing in
  the tree waits on it.
- **The three-axis review already ran and its findings are fixed and landed.**
  Two mattered: the guard's first-free-argument scan let `git stash -m list` and
  `git stash -- list` through, because a flag value was read as the subcommand;
  and `bench shift` still advised `commit or stash first`, a command the new deny
  class refuses. No `reviews/` pickup file exists, by design — the findings were
  closed in the session that found them.
- **Two calls the coordinator made and flagged rather than settled.** The anchor
  canary fixtures went under `tests/canary/workflow-guidance-anchors/` instead of
  the family the spec's Implementation-decisions paragraph named, because the
  spec contradicted its own Prior-art line; and the spec text was left unamended
  where it described the superseded free-argument mechanism. Retirement removed
  both discrepancies, and the first is the open learnings entry.
- **Both capture sources hold one entry and want a drain.** `IDEAS.md`: coverage
  maps' conformance red signals omit `BENCH_CONFORMANCE_ROOT`, so the command as
  written capability-skips and reports a false `ok`. `.bench/learnings.md`: may a
  build correct a spec's internally-inconsistent instruction on its own when the
  two readings are functionally equivalent?
- **`ROADMAP.md` still carries FT96's row and six cross-references.** Left
  deliberately — `.bench/BENCH.md` gives the roadmap reconcile to the drain phase,
  and this is reconcile work, not retirement work. Expect it as the drain's first
  finding.
- Known ambient facts: 17 worktrees remain from earlier sessions, some flagged
  unmerged or uncertain at session start, so they were left untouched — only this
  session's six were released. The structure budget violations and the
  conformance-phase long pole are unchanged.

## Next command

`/bench-what-next` — the board needs a reconcile before anything else is built:
FT96's roadmap row is stale, and both capture sources have an entry awaiting a
verdict.

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
