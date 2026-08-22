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

## #1: What is the concurrency budget model?

Blocked by: none
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

Blocked by: none
Type: Prototype

### Question
The optimum is unknown. Each inner gate is one `go test` with a warm build
cache, so width >2 may buy nothing, or width 1 may serialize compile steps
badly. Measure gate wall-clock and load at k ∈ {1, 2, 4}, with workers =
`max(1, GOMAXPROCS(0)/k)`. Run one gate per candidate on this repo, about
10–15 minutes each at today's clock.

Record the residual load while canary runs, so the deferred outer-phase
question (Not yet specified) gets its evidence in the same runs. Deliverable:
a short numbers asset; the reviewer picks k live from it. The post-change
measurement against the 2026-07-22 baseline is the arm's ship evidence.

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

Blocked by: none
Type: Grill

### Question
A bench-specific env var (with its own validation and fail posture), or reuse
a standard lever?

### Answer
No knob. Budget = `runtime.GOMAXPROCS(0)` of the outer process. In Go 1.25
this value is cgroup- and container-aware, and it already honors an explicit
`GOMAXPROCS` env var. The standard Go lever is therefore the escape hatch
(`GOMAXPROCS=8 bench gate`). `runtime.NumCPU()` is not the budget source
anywhere in this arm.

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
