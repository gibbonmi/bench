# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — gate-pipeline map closed (`0b89710` + asset-move commit),
**unpushed**; `git push` is hook-blocked for the worker, so the push is the
reviewer's own command.

## State

- **`decisions/gate-pipeline.md` is closed** — all nine tickets resolved, the
  `## Handoff` section is written, `bench maps` shows no gate-pipeline rows.
  The fixture inventory asset lives at
  `decisions/assets/gate-pipeline-fixture-inventory.md`.
- **Recommended slicing (reviewer's call, recorded in the map's Handoff):**
  slice A — canary check-scoping prerequisite (#5 + stray CHECK files),
  buildable on today's gate; slice B — manifest + DAG runner; slice C —
  `checkGoCore` split + fixture migration + parity, after A and B.
- **Decisions that stay closed:** six-field manifest
  (`name/argv/env/needs/optional/dir`), `.bench/phases.json`,
  absent=built-in / empty+malformed=red; many-to-one family→check binding
  with CHECK overrides for the nine strays; phase-named canary families with
  new-phase fixtures authored in the migration pass; unweighted width; single
  gate deadline naming stragglers; fixture-backed parity; diff-scoped gating
  stays ruled unsound; cross-language incrementality stays behind FT91's
  revive trigger.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows, not handoff state.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees, work
  verified present in `main`) remains proposed, not executed — reviewer's call.

## Next command

`/bench-write-spec` on a fresh mid-tier session, seeded with
`decisions/gate-pipeline.md` (its Handoff carries the seams; start with
slice A per the recorded dependency order). Reviewer-run
`git push origin main` still pending.

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
