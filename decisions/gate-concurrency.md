# Gate concurrency (FT91, first arm)

## Destination

Core-count-aware gate/phase concurrency: the canary phase's nested inner gates
must not oversubscribe the box, cutting gate wall-clock and the contention
symptoms (marker stalls, cleanup flakes) without changing what green means.

Measured baseline (2026-07-22, kit repo, 16 cores): gate 10–15 min, load
average ~123. Mechanism: `internal/canary/canary.go` fans 144 fixtures over
`runtime.NumCPU()` workers, each spawning a full inner gate whose `go test`
defaults to 16-wide, concurrent with the outer conformance/contract tests —
demand ~16× cores, uncoordinated.

## #1: What is the concurrency budget model?

Type: Grill

### Question
Cap only the canary worker count, cap only each inner gate's width, or budget
the product?

### Answer
Product budget. Each inner gate gets an explicit `GOMAXPROCS=k` in its env
(one lever: it caps both `go build` parallelism and the test binary), and
canary workers are capped at `max(1, budget / k)` so workers × inner width ≈
available cores. Outer phases stay as-is this arm — canary nesting is the
measured oversubscription source. The inner env must strip-and-set its own
`GOMAXPROCS` so an operator's outer override cannot leak through and defeat
the cap.

## #2: What inner width k minimizes gate wall-clock?

Type: Prototype

### Question
The optimum is unknown: each inner gate is one `go test` with a mostly warm
build cache, so width >2 may buy nothing — or width 1 may serialize compile
steps badly. Measure gate wall-clock and load at k ∈ {1, 2, 4} with workers =
`max(1, GOMAXPROCS(0)/k)`, on this repo, one gate run per candidate (~10–15
min each at today's clock). Record the residual load while canary runs so the
deferred outer-phase question (Not yet specified) gets its evidence in the
same runs. Deliverable: a short numbers asset; the reviewer picks k live from
it. The post-change measurement against the 2026-07-22 baseline is the arm's
ship evidence.

### Answer
k = 2. Measured 2026-07-24 on the kit repo (16 cores, full gate, warm compile
cache, idle box):

| k | workers | wall  | verdict | load peak | load mean |
|---|---------|-------|---------|-----------|-----------|
| 1 | 16      | 330 s | green   | 31.0      | 19.2      |
| 2 | 8       | 332 s | green   | 33.8      | 19.5      |
| 4 | 4       | 432 s | green   | 26.5      | 11.8      |

Baseline: 10–15 min, load ~123. k=1 and k=2 tie on wall; k=2 halves the
concurrent inner gates (memory/tmp pressure) at no wall cost, k=4 pays +30%
from under-parallelizing 144 fixtures. Gate wall-clock floor is now the
conformance phase (319 s), not canary. Chosen under the reviewer's batch
approval — **veto point**.

## #3: Does the budget need an operator override knob?

Type: Grill

### Question
A bench-specific env var (with its own validation and fail posture), or reuse
a standard lever?

### Answer
No knob. Budget = `runtime.GOMAXPROCS(0)` of the outer process — in Go 1.25 it
is cgroup/container-aware and already honors an explicit `GOMAXPROCS` env var,
so the standard Go lever is the escape hatch (`GOMAXPROCS=8 bench gate`).
`runtime.NumCPU()` is not the budget source anywhere in this arm.

## Not yet specified

- Capping the outer conformance/contract `go test` width — #2's evidence:
  load still peaks ~2× cores in bursts, but wall is now conformance-bound, so
  capping outer width buys nothing and could cost; dormant unless contention
  flakes persist after the canary cap ships.

## Out of scope

- Removing `-count=1` / Go test-result caching — a separate FT91 arm.
- Shared hermetic build cache, and caching keyed on the pinned gate subject —
  separate FT91 arms.
- Scoped verdicts of any kind; diff-scoped gating stays ruled unsound
  (contract/canary are behavior contracts with no file→test map).
- Weakening any check to buy wall-clock — green must keep meaning the same
  thing.

## Handoff

1. **Module boundaries.** `internal/canary` owns everything this arm changes:
   worker-count derivation in `runFixtures` and the inner-gate env pin in
   `innerEnv`. The width constant k lives in `internal/bounds` beside the
   other tunables. `internal/gate` (phase table, runner, outer phase env) is
   outside — untouched.
2. **Contracts.** `Sweep(root, Runner)` signature unchanged. Worker count =
   `max(1, runtime.GOMAXPROCS(0)/k)`, further capped at the fixture count.
   Every fixture `RunCall.Env` carries exactly one `GOMAXPROCS=k` entry, with
   any inherited outer `GOMAXPROCS` stripped first. No CLI surface, exit
   code, or output change anywhere.
3. **Deep vs thin.** `runFixtures` stays the deep unit hiding scheduling; the
   injected `Runner` is the seam. No new abstractions — the knobless constant
   replaces the prototype's env var.
4. **Black-box assertables.** Via a fake `Runner`: in-flight high-water ≤ the
   derived bound (the existing `TestSweepBoundsFixtureConcurrencyAtNumCPU`
   retargets from NumCPU to the bound); `RunCall.Env` contains the single
   pinned `GOMAXPROCS=k` even when the test sets an outer override; existing
   overlap, baseline-order, error-order, and temp-cleanup tests unchanged.
5. **Gate attachment.** The canary-package tests run inside the gate's
   conformance phase (its nested kit `go test`), so the gate sees the seam.
   The wall-clock/load outcome itself is not gate-assertable — ship evidence
   is the manual post-change measurement against the #2 table.
6. **Hostile-input owners.** No new parsed input (no knob). A hostile outer
   `GOMAXPROCS` value is handled by the Go runtime; the worker derivation
   clamps to ≥1 and the env pin strips the inherited value —
   `runFixtures`/`innerEnv` own both. Fixture counts 0/1 already handled.
7. **Uncertainty flags.** None — k is measured, the budget model and knob
   questions are closed.
8. **Rejected alternatives.** Worker-cap-only (leaves 4× oversubscription);
   global weighted semaphore across phases (machinery without measured need);
   bench-specific env knob; `runtime.NumCPU()` as budget source (not
   cgroup-aware, ignores operator override); k=1 (wall tie but double the
   concurrent inner gates); k=4 (+30% wall).
9. **Domain watch-outs.** The kit's own test suite runs nested inside the
   gate's conformance phase with the phase environment passed through — a
   canary test whose concurrency expectation assumes machine width deadlocks
   that nested run until its 600 s timeout and turns conformance red
   (observed during #2's first measurement pass). Go's exec env has no
   guaranteed duplicate-key precedence, so the inner `GOMAXPROCS` pin must
   strip-then-append, never append a duplicate. At k≤2 the gate's long pole
   is the conformance phase, not canary — further canary tuning cannot move
   wall-clock.

Dependency order: n/a — single spec.
