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

**S3 (LOW) — split three over-length files.** `bench structure` flags
`internal/adopt/link.go` (422), `internal/contract/runtime/runtime_status_test.go`
(410), and `internal/contract/surface/link_test.go` (515) over the 400-line
limit; split along responsibility, don't fragment to beat the number. Surfaced
every session via the status hook.
Next action: direct fix-and-gate with the craft-seams skill.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Watch

- `bench worktree concurrent-acquire` contract test failed once under
  full-gate load, then passed 3/3 in isolation and on rerun — likely a timing
  flake surfaced by gate phase concurrency. Journal it if it recurs.

## Recommended sequence

1. S2 unify subprocess-capture seams — direct fix-and-gate
2. S3 split three over-length files — direct fix-and-gate
