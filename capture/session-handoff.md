# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main`

## State

`/bench-drain` reconciled the board and emptied every capture source. Pinned
pre-commit HEAD: `f297ec4d`.

- `ft229-hygiene-batch` verified shipped (`5a9f3a54`) and removed from the
  board; its `## Next` audit-portfolio section is gone (A1, A2, A3, A6, FT228,
  FT229 all landed). FT174's audit-gated rewrite (drop the orphaned-ticket
  half now that FT229's close step shipped) applied per
  `docs/audits/2026-08-bench-capability/.../roadmap-dispositions.yaml`.
- 4 ideas and both pending retros drained into existing rows (FT213, FT214,
  FT223, FT225, FT233, FT169) and one new row, **FT238** (worktree-path
  ergonomics, `bench commit --dry-run`, the heredoc guard gap, a run-binary
  glossary term).
- 6 open learnings plus one untracked legacy entry verdicted and removed; all
  merged into the same existing rows above — none stood alone.
- Repair-attribution tally from both retros: 13 tickets, 5 one-shots, 10
  repair rounds (6 `spec-row`, 3 `other`, 1 `tree-drift`).
- `## Recommended sequence` refreshed: FT225 stays #1 (decision-required
  landing-amendment blocker); FT233 (7 occurrences) replaces the shipped
  FT229 at #2, ahead of FT224 by occurrence-count tie-break — both MEDIUM,
  both actionable, no stated dependency between them.
- Flagged, not applied: FT169/FT224/FT233/FT225 all edit the landing-refusal
  surface and are a restructure candidate for a future `--restructure` pass.
  FT223's fourth occurrence (a reused-cache gate verdict) is a thematic
  stretch onto a row about `bench commit`'s specific message, not `bench
  gate`'s — flagged for veto rather than moved to FT213.
- Eight landed worktrees remain from FT229, all `retain` under `bench worktree
  clean --landed` (uncommitted tracked changes). Untouched by this drain;
  still need per-path resolution by whoever owns them.
- `capture/agent-performance/claude-models.md` was already refreshed and
  uncommitted by `/bench-final-check`; included in this commit as-is.

## Next command

`/bench-shape-idea`

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
