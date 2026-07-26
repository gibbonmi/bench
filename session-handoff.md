# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — FT91 merged and green at `880357f`, **15 commits unpushed**.
Uncommitted: the 2026-07-26 `/bench-what-next` drain, staged as one batch diff
**awaiting reviewer approval** — do not commit it without that approval.

## State

- **The drain is done but uncommitted.** Both capture sources are empty (seven
  ideas, two learnings), the FT91 spec and review are retired, and the board
  now carries the pass: FT91 rewritten around the pipeline-refactor arm with
  the canary check-scoping prerequisite and two interim defects; FT116 carries
  the attributed `guards.Scan` goroutine leak; FT120 gained the self-host
  teardown-race flake; FT107 gained the doc-only gate-anchored-surface clause;
  FT141 (gate pin red verdicts) and FT142 (FT91 review residuals) are new.
- **On approval:** run the dev gate, then commit the batch as the drain commit
  plus a `spec-retire: ft91-gate-tier-split` commit, then push `main`.
- **The next build work is the gate pipeline refactor** (FT91's next arm,
  reviewer's stated top priority). Inputs are the FT91 row and
  `decisions/gate-concurrency.md`'s watch-outs; the check-scoping levers are a
  prerequisite slice, not a follow-up.
- **`bench prep-release` stays shelved** — blocked by FT116's race and FT142's
  ship-track findings; both are board rows now, not handoff state.
- **Decisions that stay closed:** ship is a superset of dev; `internal/conformance`
  is excluded from the unfiltered inner run at both tiers; the release-only
  package tests are owned by the ship-tier conformance run; diff-scoped gating
  stays ruled unsound; no check weakening for wall-clock.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees, work
  verified present in `main`) remains proposed, not executed — reviewer's call.

## Next command

Approve or adjust the staged drain diff, then, in a fresh session:
`/bench-shape-idea` on FT91's pipeline arm — the gate becomes a true pipeline.

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
