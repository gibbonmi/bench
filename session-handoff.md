# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — drain and spec retirement committed (`44f898c`, `1deed3c`),
**unpushed**; `git push` is hook-blocked for the worker, so the push is the
reviewer's own command.

## State

- **The 2026-07-26 drain is committed.** Both capture sources are empty, the
  FT91 spec and review are retired, and the board carries the pass: FT91
  rewritten around the pipeline-refactor arm, FT141/FT142 opened, FT116/FT120/
  FT107 widened.
- **`decisions/gate-pipeline.md` is bootstrapped** — 8 open tickets, #9
  (timing continuity) resolved inline. Frontier, nothing blocked: #1 (manifest
  schema grill), #2 (runner-survey research, agent-alone), #5 (canary
  check-scoping grill — the prerequisite slice, buildable on today's gate).
  Everything else hangs off those.
- **`bench prep-release` stays shelved** — blocked by FT116's race and FT142's
  ship-track findings; both are board rows now, not handoff state.
- **`bench prep-release` stays shelved** — blocked by FT116's race and FT142's
  ship-track findings; both are board rows now, not handoff state.
- **Decisions that stay closed:** ship is a superset of dev; `internal/conformance`
  is excluded from the unfiltered inner run at both tiers; the release-only
  package tests are owned by the ship-tier conformance run; diff-scoped gating
  stays ruled unsound; no check weakening for wall-clock.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees, work
  verified present in `main`) remains proposed, not executed — reviewer's call.

## Next command

`/bench-shape-idea` resume on `decisions/gate-pipeline.md` #5 (with #2's
research runnable agent-alone alongside). Before or after: `git push origin
main` — reviewer-run, 18 commits waiting.

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
