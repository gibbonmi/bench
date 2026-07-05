# gate-phase-concurrency

Status: comparison record — not staged for build (see SPEC-COMPARISON.md;
the spec of record is specs/gate-phase-concurrency.md)

<!-- coverage-map: historical -->

## Problem

`bench gate` runs its four phases — conformance, contract subtree, shellcheck,
canary sweep — sequentially inside `.bench/gate.sh`. The gate runs after every
shift iteration, so its wall time compounds through every loop of every build.
Measured 2026-07-05 on a 12-core dev box with a warm build cache: conformance
8.1s (mostly serial), contract subtree 5.9s (already parallel across packages),
shellcheck 0.3s, canary sweep 32.4s (already ~9× parallel internally). Run
end to end they cost 47.1s. The phases are independent — nothing in one reads
the output of another — yet the script pays for them one after the next.
These numbers are one 12-core box with a warm cache; the win scales with core
count and shrinks on smaller runners, which is an accepted characteristic of
the approach, not a blocker.

## Solution

Move the phase orchestration into Go. A new consumer-invisible plumbing
subcommand (working name `gate-phases`) hard-codes benchkit's four phases and
runs them concurrently; measured all-concurrent wall time is 35.3s, a 25% win,
with the ceiling at CPU saturation rather than orchestration. `.bench/gate.sh`
shrinks to root derivation plus one exec of that subcommand, staying the
resolved oracle entrypoint so gate resolution, the verdict cache, and the canary
runner are untouched.

Outer runs stream a live multiplexed output, each line prefixed with its phase
(`[contract] …`), plus a one-line per-phase verdict summary, and preserve the
`gate: green` / `gate: red` final line verbatim. Failure posture is unchanged:
run-all-and-aggregate, so one red phase exits 1 but every phase still completes.
Inner mode (`BENCH_CANARY_INNER=1`) keeps today's behavior exactly — phases
sequential, output unprefixed, sweep skipped — because the canary sweep matches
`EXPECT` substrings against inner-gate output byte-for-byte.

The phase list is pinned by three layers of existing machinery: the re-targeted
`gate_entry` conformance anchor pins the thin `gate.sh` exec line, a Go unit test
pins the phase table, and the existing canary sweep proves a dropped phase stops
biting.

## User stories

1. As a shift loop, I want the gate's four phases to run concurrently on an outer
   run, so each iteration's wall time drops without changing what the gate proves.
   Line: claude-sonnet-5 / low. The concurrent executor is plumbing at a
   pre-agreed seam whose exit-code and output contracts the gate fully observes.
2. As an agent reading gate output, I want each streamed line prefixed with its
   phase and a per-phase verdict summary at the end, so a red phase is legible in
   interleaved output. Line: claude-sonnet-5 / low. The multiplexer is a
   deterministic output shape the runtime contract asserts directly.
3. As the gate owner, I want a red phase to exit 1 while every other phase still
   runs to completion, so the aggregate verdict reports all failures at once.
   Line: claude-sonnet-5 / low. Run-all-and-aggregate is today's posture and is
   checked by exit code plus completion of every phase.
4. As the canary sweep, I want inner mode (`BENCH_CANARY_INNER=1`) to stay
   sequential, unprefixed, and sweep-skipped, so fixture `EXPECT` substring
   matches keep passing byte-for-byte. Line: claude-sonnet-5 / low. Inner-mode
   output shape is a hard byte-compatibility contract the canary bite observes.
5. As an agent interrupting a run, I want SIGINT to kill every concurrent
   phase's process group, so no orphaned `go test` survives a run with four
   children in flight. Line: claude-sonnet-5 / low. The single-process
   cancellation shape is proven by the existing `RunContext` pattern; fanning
   it out to four concurrent children is new code, but the resulting behavior
   (no surviving child process) is still a plain black-box assertion the gate
   fully observes, so cheap holds.
6. As the oracle owner, I want the phase list pinned by the conformance anchor,
   a phase-table unit test, and the canary bite, so a phase cannot be silently
   dropped. Line: claude-sonnet-5 / low. All three pins are existing machinery
   re-pointed at the new entry line and phase table, fully gate-observable.
7. As the maintainer, I want `bench gate`, `gate-run`, the shift loop, pre-push,
   and canary inner runs to all resolve the same `gate.sh` into the same
   subcommand, so every surface shares one phase orchestrator. Line:
   claude-sonnet-5 / low. Every surface reaches the gate through the same
   resolved entrypoint, so a single behavior contract covers them all.

## Implementation decisions

- One new Go unit lives in or beside `internal/gate` and owns the phase table,
  the concurrent executor, the output multiplexer, and signal handling. It
  exposes a `gate.PhasesCommand(args []string, stdout, stderr io.Writer) int`
  entrypoint following the established `gate.RunCommand` subcommand signature.
- The subcommand derives the repo root itself via `git rev-parse` (matching
  today's script) and returns today's exit codes: 0 green, 1 red, 3 not in a
  repo. It does not add a second writer of the verdict cache — `Record` in
  `internal/gate` stays the one writer; the aggregate verdict funnels through the
  same resolved `gate.sh` path that reaches it today.
- The phase table hard-codes benchkit's four phases: conformance
  (`go test ./internal/conformance -run '^TestRootConformance$'`), contract
  subtree (`go test ./internal/contract/...`, exactly once), shellcheck
  (best-effort, skipped when the tool is absent), and canary sweep
  (`bin/bench.sh canary "$root"`). The table is a Go value a unit test can pin.
  The executor accepts the phase table as an injectable value (mirroring the
  `FS`-style injection already used in `internal/gate/gate.go`), so a test can
  substitute a phase with a stub failing command to force a red phase without
  needing a broken fixture repo.
- `.bench/gate.sh` shrinks to root derivation plus one exec of the subcommand,
  quoting `"$root"`. It stays the resolved oracle entrypoint and the stable
  exit-code contract; gate resolution, the verdict cache, the canary runner, the
  shift loop, and all hooks are unchanged.
- `bin/bench.sh` routes the new plumbing subcommand with a single
  `gate-phases) route_binary "$@" ;;` arm alongside the other plumbing commands
  (no `--repair`, no second resolver). `cmd/bench/main.go` wires `gate-phases`
  into the streaming `switch` arm in `run()` (not the string-returning `commands`
  map), because it needs live multiplexed stdout/stderr.
- Outer runs prefix each streamed line with `[<phase>] `, append a one-line
  per-phase verdict summary, and print the `gate: green` / `gate: red` final line
  verbatim. Inner mode (`BENCH_CANARY_INNER=1`) runs phases sequentially, emits
  unprefixed output, and skips the sweep — byte-compatible with today.
- SIGINT kills the whole phase process group so no `go test` is orphaned, reusing
  the `RunContext` goroutine-plus-`select`-on-ctx-done cancellation pattern.
- Deep vs thin: the phase runner is the deep module (it hides concurrency,
  multiplexing, and signals behind the subcommand CLI). The `gate.sh` one-liner
  and the `bench.sh` route are thin pass-throughs with no seam of their own.

## Testing decisions

- A good test here is black-box against the subcommand's observable contract:
  exit codes, the final green/red line, phase completion under a red phase, outer
  prefixing, and inner-mode byte shape. The runtime/behavior contracts are the
  primary seam because every surface consumes the same resolved output.
- The phase-table pin is a plain unit test in `internal/gate`, `gate_test.go`
  style — table-driven, a fake for anything external, no subprocess fixtures —
  asserting four phases and that contract runs `./internal/contract/...` exactly
  once.
- The `gate_entry` conformance anchor in `internal/conformance/gate_entry_test.go`
  is re-targeted from its current loose substring checks to a verbatim match of
  `gate.sh`'s new Go-invocation exec line, so any drift in that line fails.
- The existing canary sweep is the third pin: a fixture that breaks a conformance
  or contract check must still turn the inner gate red, which fails if that phase
  vanished. This is already-present machinery, not a new test.
- Inner-mode byte compatibility needs a **direct** test, not just the canary
  bite: the canary's `EXPECT` check is a raw substring match against combined
  output, so a `[phase] ` prefix or reordered lines would still contain most
  `EXPECT` fragments and the sweep would stay green through exactly that
  divergence. A new behavior test runs the subcommand with
  `BENCH_CANARY_INNER=1` and asserts directly: no `[<phase>] ` prefixes appear,
  phases ran sequentially, and the sweep phase was skipped.
- Gate command: `.bench/gate.sh` (the default project gate).
- Not gate-observable: the wall-time win itself. It gets one before/after
  measurement during the build, not a permanent seam.

### Seam diagram

Seam 1 - phase runner (consumer-invisible orchestration contract):

    trigger: .bench/gate.sh exec bench gate-phases "$root"
        |
        v
    argv + BENCH_CANARY_INNER env + repo root (git rev-parse)
        -> [ internal/gate.PhasesCommand: phase table, concurrent
             executor, output multiplexer, process-group signals ]
        -> multiplexed prefixed stdout + per-phase summary
           + verbatim gate: green / gate: red, exit 0/1/3
              tests attach here: behavior/runtime contract on exit
              codes, prefixing, aggregate posture, inner-mode shape

Seam 2 - phase table (single source for the phase list):

    trigger: gate-phases builds its phase list
        |
        v
    hard-coded phase table value
        -> [ internal/gate phase-table unit test ]
        -> four phases, contract runs ./internal/contract/... once
              tests attach here: table-driven unit test in gate_test.go style

Seam 3 - gate.sh entry line (oracle-file to Go handoff pin):

    trigger: conformance reads .bench/gate.sh text
        |
        v
    gate.sh exec line
        -> [ internal/conformance gate_entry verbatim anchor ]
        -> pass only if the exec line matches exactly
              tests attach here: re-targeted TestGateEntry... verbatim match

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | outer `.bench/gate.sh` runs the four phases concurrently and exits 0 on all-green | phase runner behavior contract | new behavior test asserts an all-green outer run exits 0 through `bench gate-phases`; today no `gate-phases` subcommand exists so `bin/bench.sh gate-phases` errors | absence of the subcommand and its executor fails the invocation before any wall-time claim matters |
| 2 | outer stdout prefixes each line `[<phase>] ` and ends with a per-phase verdict summary | phase runner behavior contract | new test greps outer output for a `[contract] ` prefix and the summary line; today's sequential `gate.sh` emits neither | the multiplexer is absent, so prefixes and summary do not appear until it exists |
| 3 | one red phase exits 1 but every other phase still completes | phase runner behavior contract | new test injects a stub failing phase into the executor's phase table (test-only injection point) and asserts exit 1 while the other phases' output still appears | a fail-fast or short-circuit implementation would drop later phases and fail the completion assertion |
| 4a | inner mode (`BENCH_CANARY_INNER=1`) stays sequential, unprefixed, sweep-skipped | inner-mode direct behavior contract | new test runs the subcommand with `BENCH_CANARY_INNER=1` and asserts directly: no `[<phase>] ` prefix in output, phases complete in table order, sweep phase absent from output; today's script already does this but the new subcommand's inner-mode arm does not exist yet | this is the only test that can actually catch a broken inner mode — the canary's substring `EXPECT` match still passes even if inner output gains a prefix or reorders, so it cannot be the primary pin (see 4b) |
| 4b | canary sweep still bites on inner mode | canary sweep (secondary pin) | existing canary sweep: a fixture with a genuinely broken conformance/contract check must still turn inner gate red with its `EXPECT` substring | proves the phase still runs and still fails correctly, complementary to 4a's shape assertion — neither alone covers both correctness and byte-shape |
| 5 | SIGINT kills every concurrent phase's process group with no orphaned `go test` | phase runner behavior contract | new test starts a run with four children in flight, sends SIGINT, and asserts no surviving child process across all four; the new subcommand's multi-child fan-out doesn't exist yet, so the assertion has nothing to pass against | the single-process `RunContext` pattern already avoids orphans for one child today — the new work is fanning that out to four concurrent children, which is what this row actually tests |
| 6 | the phase list stays four phases with contract run exactly once | phase-table unit + conformance anchor + canary bite | `go test ./internal/gate/...` phase-table test and the re-targeted `gate_entry` verbatim check both fail if a phase or the exec line drifts; neither exists yet | three independent pins each fail on a dropped or altered phase, so silent removal cannot pass |
| 7 | every surface (`bench gate`, `gate-run`, shift loop, pre-push, canary inner) resolves the same `gate.sh` into the same subcommand | phase runner behavior contract | new test drives the gate through more than one surface and asserts identical resolution; today they resolve a sequential script with no subcommand | a second resolver or divergent entry would produce mismatched output and fail the shared-resolution assertion |

### Edge inventory

- **Error path — a phase fails.** One red phase exits 1 while all phases complete.
  Covered by story 3.
- **Error path — not in a repo.** The subcommand derives root itself and exits 3,
  today's code. Covered by the exit-code assertions under story 1.
- **Empty/absent input — shellcheck tool missing.** The shellcheck phase is
  skipped best-effort, unchanged from today. Covered under the phase-table and
  behavior contracts (a skipped phase is not a red phase).
- **Boundary — cwd deeper than the repo root.** The subcommand's own
  `git rev-parse` root derivation normalizes it. Covered by story 1's
  root-derivation path.
- **Malformed input — paths with spaces or globs.** The thin `gate.sh` quotes
  `"$root"`; the Go exec passes argv unsheltered by a shell. Covered by story 7's
  shared-resolution contract exercising a root path with a space.
- **Interrupted/partial state — SIGINT mid-run.** Process-group kill, no orphaned
  `go test`. Covered by story 5.
- **Malformed environment — go toolchain missing.** The affected phase fails with
  its own honest error and reddens the aggregate; no special-casing. Covered by
  story 3's aggregate posture.
- **Hostile environment — inner-mode byte compatibility.** Inner output must stay
  unprefixed and sequential or fixture `EXPECT` matches break. Covered by story 4.
- **Re-run idempotency.** N/a — the subcommand writes no new file format or state;
  the verdict cache stays the single existing writer. No coverage row needed.
- **Won't handle**: the wall-time win as a gate assertion — it is a one-off
  before/after build measurement, not a permanent seam (Handoff #5), and would be
  flaky to pin against core count and cache warmth.
- **Won't handle**: absent-vs-empty, trailing-newline, symlink, and file-format
  classes — n/a because no new file or state is written (Handoff #6).

## Out of scope

(none — the canary phase's internal fixture parallelism is a separate,
already-decided concern in `decisions/canary-parallel-sweep.md` that stacks with
this feature, not a cut from it; the wall-time win is a build-time measurement,
noted above, not a future capability.)
