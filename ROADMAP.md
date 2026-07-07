# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

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
