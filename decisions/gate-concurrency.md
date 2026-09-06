# Gate concurrency (FT91, first arm)

Status: shaping

## Destination

Core-count-aware gate/phase concurrency: the canary phase's nested inner gates
must not oversubscribe the box, cutting gate wall-clock and the contention
symptoms (marker stalls, cleanup flakes) without changing what green means.

Measured baseline (2026-07-22, kit repo, 16 cores): gate 10–15 min, load
average ~123. Mechanism: `internal/canary/canary.go` fans 144 fixtures over
`runtime.NumCPU()` workers. Each worker spawns a full inner gate whose
`go test` defaults to 16-wide. This runs concurrent with the outer
conformance/contract tests, so demand reaches roughly 16× the core count,
uncoordinated.

This arm is landed and stands as built. The outer layer belongs to
`decisions/gate-budget.md`. Its whole-run budget supersedes #1's budget model
and #3's `budget = runtime.GOMAXPROCS(0)`. The canary arithmetic here stops
computing from the box and draws from that pool instead.

## Notes

## Decisions so far

- [What is the concurrency budget model?](gate-concurrency/tickets/1.md): Product budget.
- [What inner width k minimizes gate wall-clock?](gate-concurrency/tickets/2.md): k = 2.
- [Does the budget need an operator override knob?](gate-concurrency/tickets/3.md): No knob.

## Not yet specified

- The outer-width question this section held is no longer fog here. The
  contention trigger it was dormant against has fired. `decisions/gate-budget.md`
  now owns the question.

## Out of scope

- Removing `-count=1` / Go test-result caching — a separate FT91 arm.
- Shared hermetic build cache, and caching keyed on the pinned gate subject —
  separate FT91 arms.
- Scoped verdicts of any kind; diff-scoped gating stays ruled unsound
  (contract/canary are behavior contracts with no file→test map).
- Weakening any check to buy wall-clock — green must keep meaning the same
  thing.

## Spec-writer discretion

## Sources
