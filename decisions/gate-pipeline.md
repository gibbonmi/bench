# Gate pipeline (FT91, pipeline arm)

Status: shaping

## Destination

The dev gate becomes a local pipeline of first-class phases with declared
dependencies. `checkGoCore`'s seven serial toolchain steps overlap instead of
running inside one test function. Conformance grades structure only. Every
canary fixture pays for only the check its `CHECK` file names. Phase
definitions become project-owned data (the runner ships in the kit, shaped
against regroup-app). This cuts dev wall-clock — today 4m36s, long pole
`package-core-guard` ~86 s — without changing what green means.

## Notes

## Decisions so far

- [What does the phase manifest declare, and who owns its failure modes?](gate-pipeline/tickets/1.md): Resolved 2026-07-26.
- [How much of `internal/gate`'s existing runner survives as the pipeline runner?](gate-pipeline/tickets/2.md): Extension, not rewrite.
- [Where exactly does `checkGoCore` split?](gate-pipeline/tickets/3.md): Resolved 2026-07-26.
- [What are the DAG's execution semantics?](gate-pipeline/tickets/4.md): Resolved 2026-07-26.
- [How does a canary fixture's inner run scope to its named check, and what shape is the dedup?](gate-pipeline/tickets/5.md): Resolved 2026-07-26.
- [What do the package-core-guard fixtures migrate to?](gate-pipeline/tickets/6.md): Resolved 2026-07-26; full inventory in `decisions/gate-pipeline/assets/gate-pipeline-fixture-inventory.md` (per-claim citations, spot-checked).
- [What proves parity — that restaging lost no check?](gate-pipeline/tickets/7.md): Resolved 2026-07-26.
- [Does the manifest hold against regroup-app?](gate-pipeline/tickets/8.md): Resolved 2026-07-26 against `~/workspace/regroup-app` (Python backend via uv + TS frontend).
- [What happens to the gate's timing and output contract?](gate-pipeline/tickets/9.md): Resolved at bootstrap — continuity, decided by the tier split.

## Not yet specified

- Ship-tier phases joining the manifest (`prep-release` as a manifest
  consumer).
- How linked repos receive runner upgrades once phase definitions are theirs
  (`bench upgrade` semantics for a project-owned manifest).
- A cross-language capability-skip surface: a documented skip-line contract,
  or a `bench skip` helper a non-Go phase can exec to report
  skipped-not-green. This is #8's gap; `optional` covers absent binaries
  today.

## Out of scope

- Cross-language incrementality (phases declare input globs, runner hashes
  them) — a small build system; deferred behind FT91's standing revive
  trigger, never built speculatively on this map.
- Removing `-count=1` and cache infrastructure — same trigger;
  oracle-semantics decisions, reviewer-led.
- Diff-scoped gating in any form — ruled unsound; the ruling stands.
- Weakening or dropping any check to buy wall-clock — green keeps meaning
  the same thing.
- Outer conformance/contract width capping — dormant unless contention
  flakes persist.
- Removing canary nesting — clause transferred to
  `decisions/gate-critical-path.md`, where it was
  reopened and ruled
  (2026-07-28): behavior-owned bites move to the owning contract test.
- The two interim defects on the FT91 row (`BENCH_CONFORMANCE_TIER` scrub
  symmetry, probe-output spill) — shape already decided there; build work,
  not map fog.

## Spec-writer discretion

## Sources
