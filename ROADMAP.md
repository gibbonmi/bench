# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain.

## Features, in priority order

Rows FT13–FT21 are the 2026-07-06 assessment drain (`ASSESSMENT.md` carries the
finding-level disposition). Each is staged and spec'd; build order within the
batch follows row order, with one coupling: FT19 and FT15 touch the same diff
parser — build FT19 first or land them together.

**FT13 (HIGH) — own the artifact-lifecycle backward path.** Staged:
`specs/artifact-lifecycle.md` — `bench spec retire`, review-pickup ownership in
implement/review/final-check phases, status orphan signal.

**FT14 (HIGH) — conformance-gate Claude hook wiring.** Staged:
`specs/claude-hook-conformance.md` — Stop/Bash/SessionStart checks mirroring
the Codex standard, plus a canary needle.

**FT15 (MED-HIGH) — review-after-merge diff surface.** Staged:
`specs/review-after-merge.md` — `bench diff --commit <sha>` + the step-1
fallback in the review phase.

**FT16 (MED) — roadmap row for shipped work signal.** Staged:
`specs/roadmap-reconcile.md` — `bench status` cross-checks rows naming spec
paths against the tree.

**FT17 (MED) — guards report wiring, not file presence.** Staged:
`specs/guards-wiring.md` — `wired` field from harness configs; state-aware
pre-push `--describe`.

**FT18 (MED) — one-source collapses.** Staged:
`specs/one-source-collapses.md` — not-in-repo phrase, coverage-map schema
owner, bench.sh header roster, profile seam-list pointer.

**FT19 (MED) — CLI contract accuracy.** Staged:
`specs/cli-contract-accuracy.md` — `bench commit -m` docs, coverage `--check`
pass line + canonical errors, roadmap not-in-repo posture, diff arg
attribution + false-empty + control-byte posture, hidden-flag help. Its build
deletes `reviews/ft9.md` (findings drained there).

**FT20 (MED-LOW) — shellcheck covers enforcement shell.** Staged:
`specs/shellcheck-coverage.md` — adapters + gate.sh + embedded pre-push asset;
loud skip; EACCES is red.

**FT21 (LOW) — docs drift pass.** Staged: `specs/docs-drift.md` — ADR 0001
rewrite, Hook Layers pin entry, CONTEXT.md skill definition, plumbing
demotion, craft-cli trigger, edge-inventory pointers.

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

**FT22 (LOW, parked) — `bench spec history <slug>`.** Fold the duplicated
`git log --grep=spec-retire` recovery incantation into the CLI (FT9 pattern).
Parked from `specs/artifact-lifecycle.md` out-of-scope.

**FT23 (LOW, parked) — model-invocable spec-authoring skill.** A `craft-spec`
skill owning coverage-map discipline; structural root of the schema
duplication. Parked from `specs/one-source-collapses.md`; build under
`craft-synthesis` after FT18 lands.

**FT24 (LOW, parked) — Codex agent-line guard parity.** `check-agent-line` on
the secondary harness, pending research on whether Codex hooks support an
Agent matcher. Parked from `specs/claude-hook-conformance.md`.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Watch

- `bench worktree concurrent-acquire` contract test failed once under
  full-gate load, then passed 3/3 in isolation and on rerun — likely a timing
  flake surfaced by gate phase concurrency. Journal it if it recurs.

## Recommended sequence

1. FT13 artifact lifecycle (`specs/artifact-lifecycle.md`) — fresh mid-tier
   session, `/bench-implement-spec`
2. FT14 Claude hook-wiring conformance (`specs/claude-hook-conformance.md`) —
   `/bench-implement-spec`
3. FT19 CLI contract accuracy (`specs/cli-contract-accuracy.md`) —
   `/bench-implement-spec` (before FT15; shared diff parser)
