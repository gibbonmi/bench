# Cost follows project size

Status: shaping

## Destination

One complaint drives this map: gate, context, and delegation cost grow with
the tree, while the value of a small change does not. This map decides that
complaint across three angles: FT91 (gate wall-clock), FT101 (ambient-surface
scope), and FT136 (delegate slicing). Each angle leaves this map spec-ready or
explicitly deferred with a revive trigger. The roadmap rows stay separate
because the owners differ.

## Notes

## Decisions so far

- [Which FT91 arms are in scope for the next build cycle?](cost-follows-project-size/tickets/1.md): Conformance arm only (reviewer, 2026-07-26).
- [How does conformance wall clock distribute across the fifteen checks?](cost-follows-project-size/tickets/2.md): Measured 2026-07-26 (throwaway probe on this tree).
- [Given the timings, does the parallelization spec proceed, and what ships with it?](cost-follows-project-size/tickets/3.md): No-go on the fan-out; the gate splits into two tiers instead (reviewer.
- [Does the FT136 slicing rule wait for the cheap-tier retest?](cost-follows-project-size/tickets/4.md): Rule lands now; retest runs separately (reviewer, 2026-07-26).
- [Where does the FT136 rule land in the kit?](cost-follows-project-size/tickets/5.md): Three surfaces, one source (reviewer, 2026-07-26).
- [How much of FT101 does this map decide now?](cost-follows-project-size/tickets/7.md): Guardrails now, build deferred (reviewer, 2026-07-26).
- [Is byte-reproducibility of release artifacts a dev-gate property?](cost-follows-project-size/tickets/8.md): Move it to the ship tier (reviewer, 2026-07-27).
- [Does the dev tier prove the generator at full four-platform breadth?](cost-follows-project-size/tickets/9.md): Host-only in dev (reviewer, 2026-07-27).
- [Must the dev gate build without network egress?](cost-follows-project-size/tickets/10.md): Yes — share the module cache (reviewer, 2026-07-27).

## Not yet specified

- `-count=1` freshness semantics — same trigger; oracle decision, reviewer-led.
  Now carries a measured price: uncached `internal/preflight` alone is 10+ min,
  so blanket `-count=1` on the inner suite is off the table without the split.
- FT101 docs-half and profile-half design — dim until the revive trigger fires.

## Out of scope

- Diff-scoped gating in any form — ruled unsound (no file→test map for contract
  and canary); the ruling stands.
- Reviving the outer conformance/contract concurrency cap — dormant unless
  contention flakes persist.
- Gate scope as a speed lever — FT101 guardrail; never reopened for wall-clock.
- Weakening any check to buy wall clock — green keeps meaning what it means.

## Spec-writer discretion

## Sources
