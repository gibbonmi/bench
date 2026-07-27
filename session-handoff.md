# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `c311946`, 11 unpushed commits
Spec: `specs/ft91-phase-manifest-dag.md` — staged, reviewer-approved 2026-07-27.
Gate: green at `ad7607c` — stale only by doc commits (spec + handoff).

## State

- **Slice B spec is staged and approved (2026-07-27):**
  `specs/ft91-phase-manifest-dag.md` — the phase manifest + DAG runner,
  compiled from `decisions/gate-pipeline.md` (#1/#4/#9). The falsification
  pass (codex `gpt-5.6-sol`, xhigh) blocked the first draft; all six
  findings were verified and folded, and the reviewer approved the result
  including the spec's 13-item map-silent veto inventory (items a–m).
  Decisions that closed with approval: the single gate deadline stays
  `gate-run`'s `bounds.GateTimeout` with SIGTERM-grace straggler naming
  (no second timer, no new bounds constant), strict manifest decoding, and
  the #9 timing reading (status quo — no new per-phase timing format).
- **Build routing:** stories 8, 9, and 12 are concurrent cancellation code —
  the spec fails `craft-line`'s venue-routing test for an unattended
  `bench shift`; build interactively at the mid tier.
- **After slice B builds:** slice C (`checkGoCore` split + fixture
  migration + parity, #3/#6/#7), whose spec also carries FT143's kit-root
  family→check binding assertion (reviewer-bundled). FT143's roadmap row
  stays until that ships. FT144's workflow decision remains the
  reviewer's, unmade.
- The FT91 interim-defect pass is committed at `19d1e5b`, gate-green; both
  defects can leave the FT91 row at the next drain.
- Codex CLI note: `codex exec` must run with stdin closed (`</dev/null`) or
  it blocks reading the pipe forever.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows, not handoff state.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees) remains
  proposed, not executed — reviewer's call.

## Next command

`/bench-implement-spec specs/ft91-phase-manifest-dag.md` — fresh mid-tier
session, interactive build at the pre-agreed seams.

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
