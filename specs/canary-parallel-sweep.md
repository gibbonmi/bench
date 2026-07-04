# Parallelize the canary sweep

Status: staged

Source map: `decisions/canary-parallel-sweep.md` (closed; Handoff complete).

## Problem

The canary sweep runs its 56 fixtures plus the vacuity baseline sequentially —
57 inner gate runs, ~35s wall on a 16-core machine. The fixtures are independent,
the sweep dominates gate wall-clock, and the gate runs after every shift
iteration, so the waste compounds across every build loop.

## Solution

The sweep keeps its exact external contract — same `Sweep` signature, same CLI
exit codes, byte-identical red output — but runs fixture inner gates through a
worker pool bounded at `runtime.NumCPU()`. The vacuity baseline still completes
first (its output feeds every fixture's vacuity check), and failures are
collected per fixture and emitted in sorted fixture order, so a red parallel
sweep prints exactly what the sequential sweep printed.

## User stories

1. As a kit developer running the gate, I want fixture inner gates to run
   concurrently, so that the sweep stops dominating gate wall-clock on every
   shift iteration.
   Line: claude-sonnet-5 / high. The seam is pre-agreed and unchanged, every
   pool behavior is asserted by fake-runner tests the gate runs, and the one
   gate-invisible risk (data races) is closed by the required race-detector run
   in the testing decisions, so the cheap model at careful effort suffices.
2. As a kit developer, I want the worker pool bounded at `runtime.NumCPU()`
   with no configuration knob, so that concurrent inner runs saturate the
   machine without oversubscribing it.
   Line: claude-sonnet-5 / high. This is the same build as story 1 and is
   asserted by the same fake-runner test, so it inherits that routing.
3. As a kit developer, I want the vacuity baseline run to complete before any
   fixture worker starts, so that each fixture's vacuity comparison keeps its
   data dependency on the baseline output.
   Line: claude-sonnet-5 / high. The ordering is a single synchronization point
   that a fake-runner test observes directly, so cheap-tier work is safe here.
4. As a reviewer reading a red sweep, I want per-fixture failures collected and
   emitted in sorted fixture order with today's exact message texts, so that a
   red parallel sweep is indistinguishable from the sequential one.
   Line: claude-sonnet-5 / high. The assertion is a byte-equality check on the
   joined error string under a scripted shuffled completion order, which the
   gate fully observes.
5. As a reviewer, I want the sweep to keep reporting all attributable failures
   instead of stopping at the first, so that one broken fixture cannot mask
   others.
   Line: claude-sonnet-5 / medium. Existing multi-error tests already pin this
   behavior and only need to stay green under the rewrite.
6. As a kit developer, I want per-fixture validation failures — missing EXPECT,
   missing files/ tree, vacuous EXPECT, temp-dir or materialize setup failure —
   reported with today's exact texts under concurrency, so that fixture
   authoring mistakes stay as diagnosable as they are today.
   Line: claude-sonnet-5 / medium. The existing validation tests carry the
   assertions unchanged; the work is keeping them green through the extraction.
7. As a kit developer, I want each worker to create and remove its own temp
   work dir on every path out of the worker, so that a sweep of any color
   leaves no residue.
   Line: claude-sonnet-5 / medium. Cleanup is observable through the fake
   runner's recorded working directories, a mechanical assertion.
8. As a CLI user, I want `bench canary [root]` exit codes (0/1/2) and the
   `Sweep(root, runner) error` signature unchanged, so that nothing outside
   `internal/canary` changes.
   Line: claude-sonnet-5 / low. This is a no-op contract pinned by existing
   tests; the story exists so the breadth is explicit, not because new work is
   needed.
9. As the reviewer, I want a one-line before/after wall-clock measurement of
   the gate's sweep on this machine reported at final check, so that the
   Handoff's memory-pressure uncertainty flag closes on data instead of
   remaining open.
   Line: claude-sonnet-5 / low. It is two timed gate runs and a sentence; the
   expected result is ~35s dropping to ~5–8s, and only a surprise escalates.

## Implementation decisions

- All changes stay inside `internal/canary`. `Sweep` keeps its signature and
  becomes the pool orchestrator; the per-fixture pipeline (validate, vacuity
  check, materialize, run, judge, clean) extracts into one worker function.
- Workers = `runtime.NumCPU()`, fixed policy, no env knob. Rejected in the map:
  `BENCH_CANARY_JOBS`, a conservative fixed cap, and overlapping the baseline
  with fixture runs (the last breaks the vacuity data dependency).
- The baseline runs to completion sequentially before the pool starts. Workers
  read `baseline.Output` only after that point, so the read is of immutable
  data and needs no locking.
- Results collect into a pre-sized per-fixture-index slot (at most one error
  message per fixture, as today); after all workers finish, non-empty slots
  join in index order. The fixture list is already sorted at discovery, so
  index order reproduces today's emission order exactly.
- No `errgroup` and no cancellation: the sweep's contract is to report all
  attributable failures, never to stop at the first (map watch-out 9). A plain
  `sync.WaitGroup` plus an index channel is the whole mechanism.
- Every error message text is contract per the map; the newline join separator
  (today's join, verified in the current code) is additionally pinned because
  row 4's byte-equality assertion depends on it.
- Test fakes become concurrency-safe: the recording runner guards its state
  with a mutex and identifies the baseline call by `FixtureDir == ""` instead
  of by call order, which parallel completion invalidates.

## Testing decisions

- A good test here drives `Sweep` with a fake `Runner` and asserts observable
  behavior — call overlap, the in-flight bound, baseline-before-fixtures
  ordering, the exact joined error string, recorded working directories — never
  pool internals. Prior art: the `recordingRunner` tests in the package.
- The tested seam is the injected `Runner` at `Sweep` (existing seam, no new
  seam introduced). Real-subprocess behavior (exit codes, the real gate) stays
  covered by the CLI tests and by the gate's own canary layer, which must still
  go green with every fixture biting.
- Gate command: `bench gate`. The package tests run inside the gate via the
  root-conformance `go test ./...` probe, so every coverage row below is
  gate-observable.
- Required one-off verification the gate does not run: the implementer runs
  `go test -race -count=1 ./internal/canary` before final check and reports the
  result. The gate has no `-race` mode, and this is the only pool risk the gate
  cannot see; a red race run blocks completion like a red gate.
- Concurrency choreography in tests uses the fake runner itself — no sleeps as
  synchronization. Overlap is proven by holding an early call open until a
  second arrives; the bound is a passive high-water mark of in-flight calls
  recorded atomically (no absence-detection timeout); sorted emission is proven
  by completing calls in reverse-sorted order. The overlap choreography
  requires a pool, so it cannot run against the sequential baseline, and it
  must guard `runtime.NumCPU()==1` (single-worker pool) rather than deadlock —
  the guard is a declared limitation in the edge inventory, not a silent skip.
- Temp-dir cleanup is observed by redirecting `TMPDIR` to a test-owned
  directory (`t.Setenv`) and asserting it is empty after `Sweep` returns —
  this observes every exit path, including validation and materialize failures
  that never reach the runner. **Deliberate deviation from Handoff item 4**,
  which placed cleanup under "via real sweep": a real-sweep check cannot
  distinguish pre-existing temp entries and misses pre-run early returns;
  flagged for reviewer sign-off in the approval table.

### Seam diagram

    trigger: bench gate (outer mode) → bench canary → Run → Sweep(root, defaultRunner)
        │
        ▼
    tests/canary/* (sorted) ──▶ [ Sweep                                ] ──▶ nil (green), or
    .bench/gate.sh path     ──▶ [   baseline run — sequential, first   ]     error: per-fixture
                                [   NumCPU workers, each: validate →   ]     messages joined in
                                [   vacuity → materialize → run →      ]     sorted fixture order
                                [   judge → clean                      ]
                        ◀ tests attach here: fake Runner injected at Sweep —
                          records calls/overlap/env, scripts completion order,
                          returns canned RunResults; assertions on the joined
                          error text and recorded temp dirs

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | with >1 fixture and NumCPU>1, at least two fixture runs are in flight simultaneously | `Sweep` + fake `Runner` | new overlap test, run red against today's sequential loop before any pool code lands | a sequential loop never has two calls in flight, so the test cannot pass by accident |
| 2 | with fixtures > NumCPU, high-water mark of in-flight fixture runs ≤ `runtime.NumCPU()` | `Sweep` + fake `Runner` | cannot start red against the sequential baseline (max in-flight is trivially 1); a regression pin recorded as a passive high-water mark, red against a grossly unbounded intermediate | a rewrite that drops the bound exceeds NumCPU under a fixture surplus and the high-water mark records it; a one-off transient excess is beyond this row's detection claim |
| 3 | the baseline call completes before the first fixture call starts | `Sweep` + fake `Runner` | starts green on today's code — a behavior pin written before the rewrite, red under a naive fully-parallel implementation | losing the ordering breaks the vacuity check's data dependency and the pin catches it immediately |
| 4 | scripted shuffled completion order still yields an error string byte-equal to the sequential sweep's | `Sweep` + fake `Runner` | cannot run against the sequential baseline — the reordering choreography needs a pool; the red step is deliberate: build the pool with naive completion-order collection first, observe this test red, then switch to indexed slots | nondeterministic output is exactly what completion-order collection produces, and byte equality rejects it |
| 5 | all failing fixtures are reported, not only the first | `Sweep` + fake `Runner` | already covered — existing malformed-fixture test asserts two messages; must stay green unmodified | a short-circuiting collector drops the second message and fails the existing assertion |
| 6 | missing EXPECT / missing files/ / vacuous EXPECT / setup failure each reported with today's exact text under concurrency | `Sweep` + fake `Runner` | already covered — existing validation tests carry the texts; must stay green with assertion texts unmodified | text drift or a validation lost in the worker extraction fails the existing tests |
| 7 | every worker temp dir is removed by the time `Sweep` returns, on green and red paths including pre-run early returns (validation and materialize failures) | `Sweep` with `TMPDIR` redirected to a test-owned dir | starts green on today's code — pin asserting the redirected temp root is empty after `Sweep`, red against a worker that leaks on any exit path | any leaked dir, including one from a path that never reaches the runner, survives in the redirected root and the emptiness check fails |
| 8 | `bench canary [root]` exits 0/1/2 as today; `Sweep` signature unchanged | `Run` / CLI tests | already covered by existing CLI and shim tests | any signature or exit-code drift breaks compilation or the existing tests |
| 9 | sweep wall-clock measured before and after on this machine | manual | not TDD-able — a one-off timing of two gate runs, reported at final check | closes Handoff uncertainty flag 7 with data; a surprising number escalates per the map |

### Edge inventory

Walked per behavior: error path, empty/absent, boundary, malformed,
interrupted/partial, re-run idempotency, hostile environment.

- Error path: baseline exits nonzero — normal by design (its output feeds only
  the vacuity comparison); per-fixture `MkdirTemp`/materialize failure → story
  6 row. Runner process-start failure → existing `defaultRunner` fallback,
  unchanged.
- Empty/absent: missing `tests/canary/` or zero fixtures → absent-harness
  error before the pool spins up; already covered by existing tests.
- Boundary: exactly one fixture (pool degenerates, overlap assertion must not
  false-fail) → handled inside the story 1 test design; fixtures > NumCPU →
  story 2 row.
- Malformed: missing EXPECT / files/, vacuous EXPECT → story 6 row.
- **Won't handle:** fixture names with path separators — impossible, names are
  single directory entries; other odd characters flow through the `MkdirTemp`
  pattern unchanged and are safe by construction.
- **Won't handle:** SIGINT mid-sweep — worker temp dirs are throwaway under the
  OS temp root and no repo state is touched; existing contract carried over.
- **Won't handle:** re-run idempotency — the sweep is stateless with fresh temp
  dirs per run; safe by construction.
- **Won't handle:** forcing NumCPU=1 in tests — the bound is asserted relative
  to `runtime.NumCPU()`, which covers the degenerate pool without faking the
  runtime.
- **Won't handle:** proving overlap on a 1-CPU machine — the pool policy makes
  overlap impossible there, so the overlap test guards on `runtime.NumCPU()>1`
  instead of deadlocking; the guard is declared here so it is a decision, not a
  skip-as-pass hole (this repo's machines are multi-core, so the assertion runs
  everywhere the kit is developed).
- Hostile environment: constrained-memory machines under NumCPU concurrent
  `go test` compiles → story 9 measurement; escalate only if the number
  surprises (Handoff flag 7).

## Out of scope

- Live per-fixture progress output while the sweep runs (streaming results as
  fixtures finish) — a separate UX capability with its own output-format
  decisions, not part of making the existing contract fast; ~8 edits, ~3 gate
  runs to build later.
