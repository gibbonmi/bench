# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, 20 unpushed commits ahead of `origin/main`
Spec: `specs/ft91-gate-phase-split.md` (Status: implemented, deliberately unretired)
Gate: green at `08482c6`; every commit since is doc-only

## State

- **The drain is committed (`93dc1a0`) and both inboxes are empty.** FT146 closed
  and was removed — its script half shipped in `08482c6`, its contract-test half
  was a mis-attribution with nothing to build, and the on-disk residue is gone.
  Two new rows: FT148, worktree orphan retirement, whose fix shape you signed off
  2026-07-27 and whose cause is confirmed in code; and FT149, the branch-delete
  guard label that quotes `-D` for a `-d`. The open learning folded onto FT126.
- **FT148 is the next build and is ready for a spec.** Cause: `bench worktree
  release` matches only the exact plaintext request that created the assignment,
  the ledger stores a one-way digest, and the harness hook derives it from the
  session id — so a dead session's worktrees are structurally unreleasable.
  Decided and closed: orphans route to `bench worktree clean`; a
  request-derivation override for `release` is rejected. The row carries the
  three feeders and the (a)/(b)/(c) fix split.
- **The "preserved" wall at session start is expected, not a new failure.** The
  pool itself is drained — `git worktree list` shows only this checkout and
  `main` is the only branch — but 17 tree-missing assignment rows and 20 retained
  recovery refs survive it, because nothing compacts an active record with a
  missing tree and FT98's landed proof misses reshaped commits. Both are on
  FT148. One ref, `refs/bench/recovery/incident-20260712-ambient-probe`, matches
  no assignment and still wants a by-hand look.
- **`ft91-gate-phase-split` stays unretired on purpose** — retiring it destroys
  your veto surface on stories 4, 5, and 9.
- **Push needs `bench gate pin` first** — interactive TTY, so it is yours.

## Next command

`/bench-write-spec`

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
