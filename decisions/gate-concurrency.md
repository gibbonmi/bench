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
— (open)

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

- Capping the outer conformance/contract `go test` width — graduates only if
  #2's residual-load numbers show contention left after the canary cap.

## Out of scope

- Removing `-count=1` / Go test-result caching — a separate FT91 arm.
- Shared hermetic build cache, and caching keyed on the pinned gate subject —
  separate FT91 arms.
- Scoped verdicts of any kind; diff-scoped gating stays ruled unsound
  (contract/canary are behavior contracts with no file→test map).
- Weakening any check to buy wall-clock — green must keep meaning the same
  thing.
