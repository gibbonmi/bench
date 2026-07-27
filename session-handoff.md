# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, 12 unpushed commits ahead of `origin/main`
Spec: `specs/ft91-gate-phase-split.md` (Status: implemented, deliberately unretired)
Gate: green as of the last code commit; every commit since has been doc-only

## State

- **The drain is committed and both inboxes are empty.** FT143 left the board
  (shipped with slice C), FT146 and FT145 opened from `IDEAS.md`, FT147 from the
  journal, and FT131/FT142/FT144 absorbed what belonged to rows that already
  existed. 51 rows, no parse failures.
- **FT91 is retargeted, and this is the finding worth carrying forward.** Slice C
  did what it said — `package-core-guard` 1m52.8s → 3.3s, the `conformance` phase
  117.4s → 8.5s — and the whole gate did not move (~4m51s). Conformance was never
  the critical path. `internal/contract/surface/artifact` (~207 s) and
  `internal/contract/surface` (~178 s) are, and no arm has touched them. The row's
  entry is now `/bench-shape-idea`, because `decisions/gate-pipeline.md` is closed
  on a premise the measurement killed.
- **One spec deviation is still open for your veto, which is why the spec is
  unretired.** Stories 4 and 5 shipped as *probed* phases rather than the kit-owned
  `.bench/phases.json` they named, and story 9 — the manifest itself — is unbuilt
  as unsatisfiable. Retiring the spec destroys the veto surface, so the drain left
  it alone.
- **Push needs `bench gate pin` first.** `.bench/` changed back in slice C and the
  pre-push hook wants a fresh pin. That command needs an interactive TTY, so it is
  yours.
- **A destroyed worktree left residue you have to clear by hand.** A pool entry at
  `~/.bench/worktrees/bench-2826441890/220aa857…-72b9811f…` holds ~47 MB of release
  tarballs and no git repo; `bench worktree release` and `bench worktree clean`
  both fail closed on it, so it is a manual `rm -rf`. Until it goes,
  `bench status`'s git row reads `git state unavailable` — that blinding is FT145,
  and the script that caused the destruction is FT146, the top of the sequence.

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
