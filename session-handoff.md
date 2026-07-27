# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, 19 unpushed commits ahead of `origin/main`
Spec: `specs/ft91-gate-phase-split.md` (Status: implemented, deliberately unretired)
Gate: green at `08482c6`; every commit since is doc-only

## State

- **The worktree-pool investigation is delivered and its decisions are closed
  (reviewer sign-off 2026-07-27).** Root cause: `bench worktree release` matches
  only the exact plaintext request string that created the assignment
  (`internal/worktree/ownership.go:344`); the ledger stores a one-way digest and
  the harness hook derives the request from the Claude session id, so a dead
  session's worktrees are structurally unreleasable. Confirmed in code and
  reproduced through the accused command. Decided: orphans route to
  `bench worktree clean` by design — a request-derivation override for `release`
  is rejected as voiding the ownership model.
- **The pool is drained and verified (2026-07-27).** All 20 registered
  worktrees, the stray directory, and all 20 branches are gone; `git worktree
  list` shows only the main checkout and `main` is the only branch. All 20
  branches were disposable: twelve landed clean, four were drafts of work main
  shipped reshaped (bench's `landed=false` on them is the FT98 patch-id gap,
  seen in the wild), the FT91 concurrency-budget arm was abandoned per the
  retarget. The temporary cleanup script has been deleted as agreed.
- **Ledger residue outlives the drain.** 17 assignment rows remain, all
  tree-missing: 16 `recovered` plus `72b9811f` (ft91-gate-phase-split) stuck
  `active` — nothing today compacts an active record with a missing tree;
  FT147(b) covers it. 20 recovery refs were retained by the tool's own plans
  (payloads it cannot prove landed — the same FT98 gap; six others retired
  cleanly during the run). The "preserved" wall at session start persists until
  they retire: `bench worktree recovery <ref>` to inspect,
  `--apply <fingerprint>` when its plan allows. One odd-shaped ref,
  `refs/bench/recovery/incident-20260712-ambient-probe`, matches no assignment
  and needs a by-hand look.
- **FT147 is signed off and awaits the drain.** The leak: kit prose orders
  worktree creation 12× for every retirement instruction (`release` is named in
  no guidance file), assignments have no timestamp/lease/reaper and the resume
  classifier hard-retains `active`, and FT98's landed proof misses reshaped
  commits. Fix shape as signed off: (a) prose — `craft-delegate` close-out duty
  (the coordinator releases each worktree it cut, at done-claim acceptance),
  `bench-implement-spec` stop-short names a retirement owner, BENCH.md inventory
  names the subcommands; (b) code — created-at timestamp plus an `orphaned`
  classifier verdict surfaced by resume as ready-to-run clean commands;
  (c) rides on FT98. Prose edits still go through `craft-synthesis` as a build.
  `/bench-what-next` should convert the parked release-refusal idea line into
  this row rather than promote it unverified.
- **Half of the FT146 row still needs your removal verdict** — the artifact
  contract-test half was a mis-attribution; nothing is left to build there.
- **`ft91-gate-phase-split` stays unretired on purpose** — retiring it destroys
  your veto surface on stories 4, 5, and 9.
- **Push needs `bench gate pin` first** — interactive TTY, so it is yours.
- Drain pending: 3 parked ideas, 1 open learning.

## Next command

`/bench-what-next`

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
