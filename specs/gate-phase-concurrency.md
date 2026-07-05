# Gate phase-level concurrency

Status: implemented

## Problem

`bench gate` runs its four phases — conformance, contract subtree, shellcheck,
and canary sweep — strictly sequentially in `.bench/gate.sh`. The gate is the
oracle: it runs after every shift iteration, on pre-push, and on every manual
`bench gate`, so its wall time compounds through every loop. On the 12-core dev
box with a warm build cache the sequential gate is 47.1s — conformance 8.1s
(mostly serial), contract subtree 5.9s (already parallel across packages),
shellcheck 0.3s, canary sweep 32.4s (~9× parallel internally); the phases are
independent and CPU headroom exists, so the serialization is pure waste. The
numbers are one 12-core box; the win scales with cores and shrinks on smaller
runners — accepted, not a blocker.

## Solution

Move the four-phase orchestration out of the shell and into a new
consumer-invisible Go plumbing subcommand (`gate-phases`) that runs the phases
concurrently. `.bench/gate.sh` shrinks to root-derivation plus one `exec` into
that subcommand, staying the resolved oracle entrypoint so gate resolution, the
verdict cache, hooks, the shift loop, and pre-push are all untouched — they keep
reaching the gate through the same `gate.sh`. Outer runs stream phase-prefixed
lines and a per-phase verdict summary; the final `gate: green` / `gate: red`
lines and the exit codes (0/1/3) are preserved verbatim. Inner mode
(`BENCH_CANARY_INNER=1`) keeps today's byte-shape exactly — sequential,
unprefixed, sweep skipped — so the canary sweep's substring matching is
unaffected. Measured ceiling: 35.3s all-concurrent, a 25% win that scales with
cores.

## User stories

1. As an agent running `bench gate`, I want the four phases to run concurrently
   in outer mode, so that the gate's wall time drops ~25% and every shift
   iteration, pre-push, and manual run is faster.
   Line: claude-opus-4-8 / medium. The concurrent executor is the genuinely
   uncertain oracle logic and correctness of the oracle outranks speed, so it
   takes the mid tier at medium effort.
2. As an agent reading gate output, I want each streamed outer-mode line prefixed
   with its phase (`[contract] …`) and a one-line per-phase verdict summary at
   the end, so that concurrent output stays legible instead of interleaving into
   noise.
   Line: claude-opus-4-8 / medium. Line-safe multiplexing under concurrency is a
   real hazard the gate can only partly observe, so it stays on the mid tier.
3. As an agent, hook, and cache consumer, I want the final `gate: green` /
   `gate: red` line and the exit codes (0 green, 1 red, 3 not-in-a-repo) preserved
   verbatim, so that the verdict cache, the Stop hook, pre-push, and humans keep
   working with no change.
   Line: claude-opus-4-8 / low. The bytes are a fixed contract and the change is
   mechanical preservation, but it sits on the oracle surface so it stays mid.
4. As an agent, I want the run-all-and-aggregate failure posture kept — one red
   phase exits 1 but every phase still completes — so that I see all failures in
   a single run rather than only the first.
   Line: claude-opus-4-8 / medium. Aggregation across concurrent phases is easy to
   get subtly wrong (early-exit on first red), so it takes medium effort on mid.
5. As the canary sweep, I want inner mode (`BENCH_CANARY_INNER=1`) to run phases
   sequentially, unprefixed, with the sweep skipped, so that EXPECT substring
   matching stays byte-compatible and I don't spawn ~180 concurrent test
   processes when 61 fixture gates run.
   Line: claude-opus-4-8 / medium. Inner byte-compatibility is the load-bearing
   constraint the whole canary layer depends on, so it stays on the mid tier.
6. As an agent who hits Ctrl-C mid-gate, I want SIGINT to kill the whole phase
   process group, so that no orphaned `go test` keeps running after the gate dies.
   Line: claude-opus-4-8 / medium. Signal and process-group handling is
   correctness-critical; it reuses the existing `RunContext` precedent but must be
   verified, so mid at medium effort.
7. As the kit maintainer, I want the four-phase table pinned by a Go unit test and
   the thin `gate.sh` exec line pinned by the re-targeted conformance anchor, so
   that a phase cannot be silently dropped and the oracle-file → Go handoff cannot
   be rerouted.
   Line: claude-opus-4-8 / low. Re-targeting a gate check is conformance-adjacent
   and correctness matters, but the edit is a mechanical substring swap, so mid at
   low effort.
8. As the kit maintainer, I want the existing canary bite to still turn the inner
   gate red when a phase's check is broken, so that a vanished phase is caught
   end-to-end and not just by the two static pins.
   Line: claude-opus-4-8 / low. This is verification that existing fixtures still
   bite against the thinned gate; almost no new code, so mid at low effort.
9. As the kit, I want a new `gate-phases` plumbing subcommand routed through
   `bin/bench.sh` and `cmd/bench/main.go`, with `.bench/gate.sh` reduced to
   root-derivation plus one `exec`, so that the phase orchestration lives in the Go
   core alongside `gate-run` and `tree-hash`.
   Line: claude-sonnet-5 / low. This is shell + dispatch plumbing at a known seam
   with an existing precedent (`gate-run`, `tree-hash`), so it routes cheap.

## Implementation decisions

- **New Go unit in `internal/gate`.** Add a phase runner (working file
  `internal/gate/phases.go`) owning: the phase table, the concurrent executor, the
  line-buffered output multiplexer, and process-group signal handling. Everything
  else in `internal/gate` — the resolution chain, `RunContext`, and the verdict
  cache `Record` — is untouched; the runner **reuses** the existing
  `RunContext`-style process-group kill (SIGINT then SIGKILL after a grace, exit
  130) rather than inventing a second signal mechanism.
- **Injectable phase list is the seam.** The executor takes an explicit
  `[]Phase` (each phase: a display name, an argv, and an `optional` flag for
  shellcheck's best-effort skip). Tests inject fake phases that run
  `bash -c 'echo …; exit N'` / sleeps to exercise concurrency, aggregation,
  multiplexing, and signals without a real Go toolchain. The real benchkit
  four-phase table is a separate constant the runner assembles per mode.
- **Subcommand entry follows the `gate-run` precedent.** Add
  `gate.PhasesCommand(args, stdout, stderr)` and dispatch it from
  `cmd/bench/main.go`'s raw-stdio `switch` (it needs live streaming, not a
  buffered string return, so it joins the `switch` beside `gate-run`, not the
  `commands` map). Add one `gate-phases) route_binary "$@" ;;` arm in
  `bin/bench.sh` beside `tree-hash` and `worktree-lease-file` — no other
  `bin/bench.sh` change.
- **Thin `gate.sh` derives root, then execs.** `.bench/gate.sh` keeps its
  `git rev-parse --show-toplevel` root-derivation (preserving the `gate: not in a
  git repo` message and `exit 3` fast path) and then `exec`s the subcommand
  through `bin/bench.sh`, passing `$root`. The subcommand re-derives root itself so
  it is also correct when invoked directly as plumbing. gate.sh stays the resolved
  oracle entrypoint.
- **Modes.** Outer: conformance, contract, shellcheck, and canary sweep run
  concurrently; output is `[phase]`-prefixed; a per-phase verdict summary precedes
  the verbatim `gate:` line. Inner (`BENCH_CANARY_INNER=1`): the same phases minus
  the canary sweep run sequentially with unprefixed output — byte-identical to
  today. The contract phase runs `./internal/contract/...` exactly once in both
  modes.
- **Re-target the conformance anchor.**
  `internal/conformance/gate_entry_test.go`
  (`TestGateEntryRunsGoConformanceAndBehaviorContracts`) stops matching the three
  inline `go test` / canary command lines (those facts move into the Go phase
  table, newly pinned by the phase-table unit test) and instead pins the thin
  `gate.sh`'s exec line verbatim. The separate
  `internal/canary/gate_entry_test.go` pin of the `BENCH_CANARY_INNER` guard is
  preserved.

## Testing decisions

- **A good test here** drives the phase runner at its `[]Phase` seam with fake
  phases and observes external behavior — stdout bytes, exit code, per-phase
  summary, child-process state — never the runner's internals. The pin tests read
  `gate.sh` text and the phase-table constant as black-box facts.
- **Seams tested:** (1) the `gate-phases` runner (the deep module); (2) the phase
  table constant; (3) the re-targeted conformance anchor over `gate.sh`; (4) the
  existing canary sweep, unchanged, as the end-to-end bite. The thin `gate.sh`
  one-liner, the `bin/bench.sh` case arm, and the `main.go` dispatch entry are
  pass-throughs with no seam of their own — they are exercised whenever any runner
  test or the full `bench gate` run routes through them.
- **Prior art:** table-driven plain-stdlib tests (`internal/gate/gate_test.go`'s
  `TestResolvePrecedence`, FS-injection pattern); the runner-injection pattern in
  `internal/canary`'s `Sweep(root, runner)`; the substring pins in the two
  `gate_entry_test.go` files.
- **Gate command:** the project gate, `.bench/gate.sh`.

### Seam diagram

**Seam 1 — the `gate-phases` phase runner (deep module).**

    trigger: `.bench/gate.sh` exec
      (also reached via: bench gate → gate-run → gate.sh; the shift loop;
       pre-push; and the canary sweep's inner gate runs)
        │
        ▼
    mode env (BENCH_CANARY_INNER)  ──▶  [ gate-phases                          ]  ──▶  [phase]-prefixed live stream (outer)
    injected []Phase (in tests)    ──▶  [   phase table (4; contract ×1)       ]  ──▶  per-phase verdict summary
                                        [   concurrent executor  (outer)       ]  ──▶  `gate: green` / `gate: red` (verbatim)
                                        [   sequential+unprefixed (inner)      ]  ──▶  exit 0 / 1 / 3
                                        [   line-buffered output multiplexer    ]
                                        [   process-group signal handler        ]
                                            ◀ tests attach here: call the runner with fake phases
                                              (echo / false / sleep) and capture stdout + exit code
                                              + child-process state

**Seam 2 — the three-layer pin (the oracle graded against itself).**

    trigger: `bench gate`
        │
        ▼
    `.bench/gate.sh` text     ──▶  [ conformance: gate_entry anchor ]  ──▶  red if exec line absent or stale needles present
    phase-table constant      ──▶  [ phase-table unit test          ]  ──▶  red if ≠ 4 phases or contract ≠ exactly once
    broken-check fixture      ──▶  [ canary sweep (existing)        ]  ──▶  red if a dropped phase stops biting
                                       ◀ tests attach here: read gate.sh text; assert the table
                                         shape; run the inner gate on a fixture and require
                                         exit ≠ 0 + the EXPECT substring

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | all four phases run concurrently in outer mode | gate-phases runner | new `TestRunnerRunsPhasesConcurrently`: inject four fake phases that each emit a marker then sleep; assert every marker present and wall-time below the serial sum — fails before the concurrent executor exists | sequential execution or a dropped phase fails the timing/marker assertion |
| 2 | every outer line prefixed `[<phase>] `; no line split or garbled across phases | runner multiplexer | new `TestRunnerPrefixesAndKeepsLinesIntact`: two phases emit interleaved multi-line output; assert each output line begins with its own `[phase] ` and equals the source line | a naive byte-copy mux interleaves partial lines, producing a line with no clean prefix |
| 3 | final `gate: green`/`gate: red` verbatim; exit 0 green / 1 red / 3 not-in-repo | runner + thin gate.sh | new `TestRunnerFinalLineAndExitCodes` (green and red cases) plus a shell check that outside a repo `gate.sh` exits 3 — fails before the runner emits the verbatim lines and codes | cache and hooks parse these exact bytes and codes; a changed line or code fails the assertion |
| 4 | one red phase exits 1 but all phases still complete (aggregate posture) | runner | new `TestRunnerAggregatesAllPhases`: one failing plus three passing fake phases; assert exit == 1 AND all four produced output | an early-exit-on-first-red drops later phases' output, failing the completeness assertion |
| 5 | inner mode is sequential, unprefixed, sweep excluded | runner inner branch + canary sweep | new `TestRunnerInnerModeByteShape`: run with `BENCH_CANARY_INNER=1`; assert no `[` prefix on any line and no canary phase in the table; plus the full `bench canary` sweep stays green — fails before the inner branch suppresses prefixing | the canary's EXPECT check is a raw substring match against combined output, so the sweep stays green even if inner output gains a `[phase] ` prefix or reorders — this direct shape assertion is the only pin that catches that divergence |
| 6 | SIGINT kills the whole phase process group; no orphaned `go test` | runner signal handler | new `TestRunnerCancelKillsGroup`: start a long sleep phase, cancel the context, assert the child process group is reaped and exit == 130 — fails before `Setpgid` + group-kill are wired | without the process-group kill the child outlives the runner, failing the reaped-child assertion |
| 7 | phase table has exactly four phases; contract runs `./internal/contract/...` exactly once | phase-table unit test | new `TestPhaseTable`: assert `len == 4` and exactly one entry targets the contract subtree — fails before the table constant exists | a silently dropped or duplicated phase changes the count or uniqueness, failing the assertion |
| 7 | thin `gate.sh` exec line pinned; retired inline `go test`/canary needles gone | conformance gate_entry anchor | re-targeted `TestGateEntryRunsGoConformanceAndBehaviorContracts`: assert `gate.sh` contains the exec line and NOT the old inline needles — fails until `gate.sh` is thinned and the anchor re-targeted | rerouting `gate.sh` away from the subcommand drops the exec line, failing the anchor |
| 8 | a broken conformance/contract check still turns the inner gate red under the new runner | canary sweep (existing) | already covered — existing canary fixtures run the now-thinned inner gate and require exit ≠ 0 plus the EXPECT substring; not new TDD | if a phase vanished from the runner its fixture would stop biting, turning the canary red |
| 9 | `gate-phases` reachable via `bin/bench.sh` + `main.go` dispatch; `gate.sh` execs it | integration (gate.sh → gate-phases) | already covered — every runner test and the full `bench gate` run route through the subcommand; a missing case arm makes the gate fail to run at all | a broken route makes the entire gate unrunnable, failing every downstream check |
| edge of 1 | root path containing a space still runs every phase | runner | new `TestRunnerRootWithSpace`: run in a temp dir whose path contains a space; assert phases run and the exit code is correct | shell interpolation of an unquoted path would break on the space; `exec.Command` argv would not — the test pins argv handling |
| edge of 5 | shellcheck absent → phase skipped, not failed | runner | new `TestRunnerShellcheckAbsentSkips`: run with a PATH lacking shellcheck; assert exit is unaffected and no shellcheck failure line appears | a hard dependency would fail the gate on machines without shellcheck; the exit assertion catches the regression |

### Edge inventory

Edge classes walked per behavior against the profile's shell-CLI hostile-input
checklist; each lands as a coverage row above or a **Won't handle** line here.

- **paths/dirs with spaces or globs** — coverage row (edge of 1); the thin
  `gate.sh` quotes `"$root"` and the Go executor passes argv unsheltered by a
  shell.
- **required tool missing (shellcheck)** — coverage row (edge of 5); best-effort
  skip, unchanged.
- **interrupt (SIGINT) mid-run** — coverage row (story 6); process-group kill, no
  orphaned `go test`.
- **invocation through every shipped surface** — coverage row (story 9) plus the
  canary inner runs (story 5); `bench gate`, `gate-run`, the shift loop, pre-push,
  and canary all resolve the same `gate.sh` → same subcommand.
- **empty/absent input, trailing-newline, absent-vs-empty file** — **Won't
  handle** — the runner writes no new files and reads no new file formats; there is
  no absent/empty/newline surface to distinguish.
- **invocation through a symlink** — **Won't handle** — `gate.sh`'s
  `${BASH_SOURCE[0]}` path resolution is unchanged by this feature; not its
  surface.
- **go toolchain missing** — **Won't handle** — the affected phase fails red with
  the toolchain's own honest error, exactly as today; no new masking.
- **re-run idempotency** — **Won't handle** — the gate writes no state of its own
  (the verdict cache is out of scope, unchanged); a re-run recomputes the same
  verdict with nothing to reconcile.
- **cwd deeper than repo root** — **Won't handle** here as a new case — the
  subcommand derives root via `git rev-parse` and works from any subdirectory,
  same as the existing `gate-run`; covered by that shared derivation.
- **concurrent inner mode** — **Won't handle** — rejected in the map (#4); inner
  mode stays sequential to avoid oversubscription for zero measured win. An
  in-scope caller (the canary sweep) still exercises the full feature in inner mode
  under this exclusion.

## Out of scope

None as deferred cuts. The map's rejected alternatives — bash background jobs in
`gate.sh`, a generic phase-runner DSL for linked-repo reuse, front-phases-only
overlap, buffered capture-and-replay output, concurrent inner mode, a
`--describe` phase manifest, and a verbatim Go-source grep pin — are recorded
**rejections** in `decisions/gate-phase-concurrency.md`, not scope cuts to build
later. The one-time before/after wall-time measurement (the #1 numbers, re-taken
during the build) is part of this build, not a deferral: it is not
gate-observable and leaves no permanent seam.
