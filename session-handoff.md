# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `440266e`, main checkout clean, 9 unpushed commits
Spec: none staged
Gate: green at `edc15e9`, reported plain-stale — the drain's doc-only commit is
the only drift.

## State

- **The board is drained and reconciled.** 46 rows, both capture sources empty,
  `bench roadmap --context` parses with zero failures. FT96 shipped complete on
  2026-07-26 and its row is gone; the four rows that cross-referenced it (FT103,
  FT123, FT131, FT136) now say what actually landed.
- **Two drain calls are open to post-hoc veto.** The self-contradicting-spec
  learning was folded into FT107 as a fourth clause rather than dismissed — the
  batch-approval rule covers an AFK reviewer and genuinely says nothing about a
  spec whose own sections disagree. And FT133 was placed third in the sequence
  ahead of FT107; a hole in the coverage oracle outranks prose, but not by much.
- **FT133 gained a reproduced second instance.** On `main`, `go test
  ./internal/conformance -run ^TestRootConformance$` without
  `BENCH_CONFORMANCE_ROOT` prints `ok … 0.002s` and skips invisibly. That widens
  the row from *does the citation resolve* to *does it actually execute*.
- **The nine commits are unpushed, and the push is the reviewer's.** Nothing in
  the tree waits on it.
- Known ambient facts, unchanged: 17 worktrees remain from earlier sessions and
  were left untouched; the structure budget violations and the conformance-phase
  long pole stand where they were.

## Next command

`/bench-shape-idea` — the cost-follows-project-size complaint, shaped once
across its three angles (FT91's remaining gate arms, FT101's scoped surfaces,
FT136's delegate slicing). Gate wall-clock is the reviewer's stated dominant
cost, and FT91's first step — timing the fifteen conformance checks — is the
session's cheapest evidence request.

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
