# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, 10 unpushed commits ahead of `origin/main`
Spec: `specs/ft91-gate-phase-split.md` (Status: implemented)
Gate: green at HEAD — all ten phases

## State

- **FT91 slice C is built, reviewed, repaired, and gate-green — awaiting your
  push.** `checkGoCore`'s seven serial toolchain steps are now first-class gate
  phases: `gofmt`, `vet`, `test`, `race`, `conformance-suite`, plus the existing
  `build`. Four policy-carrying steps sit behind a new `bench gate-go` plumbing
  subcommand; the residual conformance check keeps only the absent-toolchain
  diagnostic and the ship cross-compile matrix. FT143 (kit-root family→check
  binding) closed with it.
- **Push needs `bench gate pin` first.** `.bench/` changed, so the pre-push hook
  wants a fresh pin. That command requires an interactive TTY, so it is yours.
- **One spec deviation is open for your veto.** Stories 4 and 5 shipped as
  *probed phases* — `race` materializes only when the graded root declares
  `TestConcurrentCleanupRecordsOneTransaction`, `conformance-suite` only when it
  declares `TestRootConformance` — instead of the kit-owned `.bench/phases.json`
  stories 4, 5, and 9 specified. **Story 9 is unbuilt.** Its acceptance row is
  unsatisfiable: the manifest loader forces `Dir` onto every declared phase while
  the built-in table leaves it empty, the `build` argv carries absolute paths a
  committed file cannot hold, and `shellcheck`'s argv is discovered at run time by
  scanning `.bench/hooks`. The spec is left unretired so the veto surface survives.
- **The slice's duration premise did not hold, and this is the finding that
  matters.** Measured on this host: `package-core-guard` 1m52.8s → 3.3s and the
  `conformance` phase 117.4s → 8.5s, but the whole gate is unchanged (~4m51s).
  Conformance was never the critical path — `internal/contract/surface/artifact`
  (~207s) and `surface` (~178s) are, and they are untouched. Any future gate-speed
  work belongs there, not in conformance.
- **A phase-table change does not take effect until the *next* gate run.**
  `gate-phases` resolves the table in a process that started before the build
  phase rewrote `dist/bench`, so the first run after a table edit silently grades
  with the old table. Observed here: a green run showing five phases, then ten on
  re-run with no code change. Not yet parked — decide whether it is a defect.
- **A worktree checkout was destroyed mid-session.** `scripts/build-offline-archives.sh`
  deletes whatever directory it is handed as `<output-dir>`, validating only that
  it is a directory; driving the artifact contract test with the graded root set
  to a live checkout replaced that checkout with release tarballs. No commits were
  lost. Parked in `IDEAS.md`. **A stale pool entry remains** at
  `~/.bench/worktrees/bench-2826441890/220aa857…-72b9811f…` holding ~47 MB of
  tarballs and no git repo; both `bench worktree release` and `bench worktree clean`
  fail closed on it, so removing it is a manual `rm -rf` and your call.
- **`bench status`'s git row reads `git state unavailable`** because of that broken
  pool entry — one unreadable worktree blinds the whole aggregate. Parked.
- **Two review findings were accepted as residual risk, then fixed anyway at your
  instruction**, so nothing is outstanding: the cross-compile stress assertion now
  executes at ship and has a dev-tier tripwire on its call site, and no rebuilt
  fixture EXPECT depends on Go toolchain wording except the race detector's own
  `WARNING: DATA RACE`.
- **Queued, with your ownership decision already made:** codify the
  `## Next command` field — a conformance check (not `bench handoff` write-time
  validation), because the field is most often hand-written. It closes the one open
  `.bench/learnings.md` entry and is small enough for the lighter path.
- `IDEAS.md` holds 3 parked items and `.bench/learnings.md` 1 open entry, so a
  drain is pending. `internal/gate/` is further over its granted file budget; it
  surfaces on its own.

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
