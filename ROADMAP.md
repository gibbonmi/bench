# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain.

## Ready to build

**S2 (LOW) — unify subprocess-capture seams.** Conformance `runProbe`,
`Harness.Run`, and canary `defaultRunner` each hand-roll subprocess capture;
collapse to one probe seam (canary's merged-stream EXPECT matching stays a
deliberate mode). Knowledge duplication, no behavior change.
Next action: direct fix-and-gate with the craft-seams skill.

**L1 (MED-LOW) — promote the shared-tree worktree rule.** From the learnings
journal (2026-07-05 shared-tree contention): when `git status` shows another
writer's in-flight edits, side-work goes to a bench worktree — or waits — so
every gate verdict is attributable to one diff. Kit edit under the
craft-synthesis discipline, gated as usual.
Next action: direct fix-and-gate (one rule sentence plus its anchor).

**FT11 (MED-LOW) — `bench worktree clean` + honest status actions.** Status
flags out-of-pool worktrees with action "resume or clean up (bench worktree)",
but the `worktree` case discards all args and opens a pool subshell —
`bench worktree clean` exits 0 and silently does the wrong thing. Add a
clean/prune verb (remove clean out-of-pool worktrees + `git worktree prune`),
reject unknown args instead of swallowing them, and make the status action
per-class honest. Next action: `/bench-write-spec`.

## Features, in priority order

**FT3 (MED-LOW) — `bench spec implemented` + `bench commit`.** Pair them:
commit could fold in the spec status flip. Replaces footgun prose in the
implement phase with two small wrappers over existing logic.

**FT4 (MED-LOW) — harness task list in `/bench-implement-spec`.** Per-harness
adapter (Claude hook + phase line; Codex native).

**FT9 (MED-LOW) — compiled git-context command.** One call bundling
status+diff+log+staged for agents, partitioned self vs other-writer changes.
Open fork: attribution via session-start baseline snapshot vs agent-passed
file list; tension with the shared-tree rule (L1 is the stronger fix), but
call-count value stands single-writer too. Needs its grill.

**FT10 (LOW) — doctor installs the kit repo's pre-push guard.** `bench guards`
already reports the missing guard; `bench doctor` should detect it on the kit
repo itself and offer the install (consumer repos get it via `bench link`).

**FT5 (LOW) — `bench outline`.** Marginal for this repo, real as a kit
affordance for large/polyglot linked repos. Needs its grill (languages,
on-demand vs committed, prose anchors).

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT7 (LOW) — dashboard.** Low priority by declaration.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Watch

- `bench worktree concurrent-acquire` contract test failed once under
  full-gate load, then passed 3/3 in isolation and on rerun — likely a timing
  flake surfaced by gate phase concurrency. Journal it if it recurs.

## Recommended sequence

1. Retire the implemented adversarial-gate-pinning spec — spec-retire
2. FT11 `bench worktree clean` — `/bench-write-spec`
