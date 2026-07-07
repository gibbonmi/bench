# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain.

## Features, in priority order

**FT12 (LOW, kit discipline) — repro a defect claim through the accused command
before draining it.** FT11 was minted from a learning that quoted a raw `git add`
run by hand; the real `bench commit` path already staged deletions, so the row
described a defect that did not exist. Tighten `/bench-what-next` step 3 (and
`bench-debug`'s repro discipline) so a defect-shaped learning becomes a roadmap
row only after its red signal reproduces through the sanctioned command, not a
lookalike. Built later under the `craft-synthesis` discipline.

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

1. FT12 repro-before-drain discipline — `/bench-shape-idea` (kit-discipline edit,
   built under `craft-synthesis`)
2. FT10 doctor installs the kit repo's pre-push guard — `/bench-implement-spec`
