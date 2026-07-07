# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

**FT38 (MED) — `bench commit` stages named deletions and renames reliably.**
Kit defect reproduced through the sanctioned command (git 2.43): the per-path
literal-pathspec staging fails for some paths absent from the worktree — hit on
a delete-plus-rename diff, while plain named deletions commit fine — forcing a
raw-git fallback the command exists to prevent. Fix the staging path (an absent
tracked path stages its removal) and pin both the deletion and rename cases;
the command's own doc comment already claims deletion support.

**FT37 (MED) — deflake the gate under concurrent load.** Transient contract or
conformance reds recurred repeatedly when multiple full gates ran concurrently
(delegate worktrees plus the main checkout), every one green on immediate
re-run with an identical tree; the earlier Watch note saw the same under a
single gate's phase concurrency. Make the worktree concurrent-acquire contract
(and any load-sensitive sibling) deterministic under parallel gate runs, or
serialize its acquire window.

**FT12 (LOW, kit discipline) — repro a defect claim through the accused command
before draining it.** FT11 was minted from a learning that quoted a raw `git add`
run by hand; the real `bench commit` path already staged deletions, so the row
described a defect that did not exist. Tighten `/bench-what-next` step 3 (and
`bench-debug`'s repro discipline) so a defect-shaped learning becomes a roadmap
row only after its red signal reproduces through the sanctioned command, not a
lookalike. Built later under the `craft-synthesis` discipline.

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

**FT22 (LOW, parked) — `bench spec history <slug>`.** Fold the duplicated
`git log --grep=spec-retire` recovery incantation into the CLI (FT9 pattern).
Parked from the artifact-lifecycle build's out-of-scope list.

**FT24 (LOW, parked) — Codex agent-line guard parity.** `check-agent-line` on
the secondary harness, pending research on whether Codex hooks support an
Agent matcher. Parked from the claude-hook-conformance build.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Recommended sequence

1. FT38 `bench commit` deletion/rename staging — `/bench-debug` (defect with a
   recorded repro)
2. FT37 gate deflake — `/bench-debug`
