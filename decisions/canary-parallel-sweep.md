# Parallelize the canary sweep

Fixtures are independent; the sweep in `internal/canary` runs them sequentially.
Measured baseline: 57 inner runs (56 fixtures + vacuity baseline), ~35s wall on a
16-core machine, ~0.6s per fixture run. The sweep dominates gate wall-clock, and
the gate runs after every shift iteration, so the save compounds.

## #1: What bounds the worker pool?

Blocked by: —
Type: Grill

### Question
Each worker spawns an inner gate that itself runs `go test` plus shell contract
fragments, so unbounded fixture-count concurrency would oversubscribe. Fixed
policy, or an operator knob?

### Answer
Workers = `runtime.NumCPU()`, no knob. Inner runs are subprocess-wait-heavy, so
NumCPU saturates without thrashing. An env override (`BENCH_CANARY_JOBS`) and a
conservative fixed cap were both rejected — add a knob only when a real
constrained environment shows up.

## #2: Is there shared mutable state across concurrent inner runs?

Blocked by: —
Type: Research

### Answer
No race surface found. Each inner run gets its own temp work dir with its own git
dir. Inner gate runs invoke `gate.sh` directly and never write the verdict cache
(only `internal/gate.Record` on the wrapper path does, keyed to the outer repo's
git dir). Concurrent `go test` invocations share GOCACHE, which is
concurrency-safe by design.

## #3: How do ordering and output stay deterministic?

Blocked by: —
Type: Grill

### Answer
Two constraints. (a) The vacuity check compares each fixture's EXPECT against the
baseline run's output, so the baseline runs to completion sequentially before any
fixture worker starts. (b) Failure output is collected per fixture and emitted in
sorted fixture order regardless of completion order, so a red sweep prints
identically to today's sequential sweep.

## Handoff

1. **Module boundaries.** All inside `internal/canary`: `Sweep` keeps its
   signature and becomes the pool orchestrator; per-fixture work (validate,
   materialize, run, judge, clean) extracts into one worker function. Nothing
   outside the package changes — `Run`, the gate line, and the CLI surface are
   untouched.
2. **Contracts.** `Sweep(root, runner) error` unchanged: nil on green; on red, an
   error joining per-fixture messages in sorted fixture order with today's exact
   message texts. `bench canary [root]` exit codes unchanged (0/1/2).
3. **Deep vs thin.** `Sweep` is the deep unit hiding pool mechanics behind the
   unchanged signature. The injected `Runner` stays the seam for tests — a fake
   runner that records concurrency and timing needs no real subprocesses.
4. **Black-box assertables.** Via a fake `Runner`: overlapping in-flight calls
   prove parallelism; max in-flight ≤ NumCPU proves the bound; baseline call
   completes before the first fixture call proves ordering; shuffled completion
   still yields sorted error output proves determinism. Via real sweep: exit
   codes and temp-dir cleanup.
5. **Gate attachment.** The existing canary layer itself is the end-to-end net —
   the gate's sweep still must go green on the kit and every fixture must still
   bite. Pool-specific behavior attaches at the `Runner` seam in `go test`
   (package tests the gate already runs). No seam is gate-invisible.
6. **Hostile-input owners.** Fixture names with spaces/globs → temp-dir naming
   in the worker (already exercised by `MkdirTemp` pattern). Missing
   EXPECT/files → per-fixture validation, must still report under concurrency.
   SIGINT mid-sweep → workers' temp dirs are throwaway; no repo state touched
   (existing contract, carried over).
7. **Uncertainty flags.** Whether NumCPU workers each spawning a `go test`
   compile step causes memory pressure on small machines — the spec should
   measure once on this machine (expect ~35s → ~5–8s) and note the result;
   escalate only if the measurement surprises.
8. **Rejected alternatives.** `BENCH_CANARY_JOBS` env knob; conservative fixed
   cap (e.g. 8); parallelizing the baseline with the fixtures (breaks the
   vacuity check's data dependency).
9. **Domain watch-outs.** `errgroup`/channel collection must not short-circuit:
   the sweep's contract is to report *all* attributable failures, not stop at
   the first.

Dependency order: n/a — single spec.
