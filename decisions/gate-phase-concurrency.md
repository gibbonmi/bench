# Gate phase-level concurrency

`bench gate` runs its four phases (conformance, contract subtree, shellcheck,
canary sweep) sequentially in `.bench/gate.sh`. Parked 2026-07-04 pending a
remeasure after the test-layout splits; graduated 2026-07-05. The gate runs after
every shift iteration, so wall time compounds through every loop.

## #1: Is there a measurable win after the test-layout splits?

Blocked by: —
Type: Research

### Answer
Yes — 25%. Measured 2026-07-05 on the 12-core dev box, warm build cache:
conformance 8.1s wall (mostly serial), contract subtree 5.9s (already parallel
across packages), shellcheck 0.3s, canary sweep 32.4s (~9× parallel internally,
near core saturation). Sequential gate: 47.1s. All phases concurrent: 35.3s —
canary stretches only ~3s under contention. The ceiling is CPU saturation, not
orchestration. Go decision: proceed (reviewer, 2026-07-05).

## #2: Where does the phase orchestration live?

Blocked by: —
Type: Grill

### Answer
Fully in Go, resolved 2026-07-05. A new consumer-invisible plumbing subcommand
(working name `gate-phases`; plumbing group precedent: `tree-hash`,
`worktree-lease-file`) hard-codes benchkit's four phases and runs them
concurrently; `.bench/gate.sh` becomes a thin exec of it, staying the resolved
oracle entrypoint so gate resolution, the verdict cache, and the canary runner
are untouched. Rationale: no shell→Go phase-declaration protocol to invent, the
oracle's scheduling logic becomes unit-testable, and it matches the shell→Go
migration direction. Blast radius explicitly waived by the reviewer. Rejected:
bash background jobs in gate.sh (silent until done, manual signal plumbing,
untestable); a generic phase-runner the gate script drives (requires a
phase-declaration mini-DSL; linked-repo reuse is speculative — their gates are
`npm test`-shaped); overlapping only the front phases (~40.5s, smaller win).

## #3: What pins the phase list so a phase cannot be silently dropped?

Blocked by: #2
Type: Grill

### Answer
Three layers, all existing machinery; resolved 2026-07-05. (1) The `gate_entry`
conformance anchor re-targets to pin the thin `gate.sh`'s exec line verbatim, so
the oracle-file → Go handoff cannot be rerouted. (2) A Go unit test pins the
phase table itself: four phases, contract runs `./internal/contract/...` exactly
once. (3) The existing canary sweep proves dropped phases stop biting — a
fixture that breaks a conformance or contract check must turn the inner gate
red, which fails if that phase vanished. Rejected: a `--describe` phase manifest
(new surface for a fact two tests already pin); conformance grepping the Go
source verbatim (brittle to refactors).

## #4: What do concurrency, output, and failure look like at the surface?

Blocked by: #2
Type: Grill

### Answer
Resolved 2026-07-05. Output: live multiplexed stream, each line prefixed with
its phase (`[contract] …`), plus a one-line per-phase verdict summary at the
end; the final `gate: green` / `gate: red` lines are preserved verbatim.
Failure posture unchanged: run-all-and-aggregate — one red phase exits 1 but
every phase still completes. SIGINT kills the whole phase process group, no
orphaned `go test`. Inner mode (`BENCH_CANARY_INNER=1`) keeps today's behavior
exactly: phases sequential, output unprefixed, sweep skipped — byte-compatible
with what canary EXPECT substring matching sees today, and it avoids ~180
concurrent test processes when 61 fixture gates run under the sweep (the outer
measurement in #1 assumed sequential inner). Rejected: buffered
capture-and-replay (silent for ~35s); concurrent inner mode (oversubscription
for zero measured win).

## Handoff

1. **Module boundaries.** One new Go unit (in or beside `internal/gate`) owning
   the phase table, the concurrent executor, the output multiplexer, and signal
   handling; `bin/bench.sh` routes the new plumbing subcommand; `.bench/gate.sh`
   shrinks to root-derivation plus one exec. Outside: gate resolution and the
   verdict cache (`internal/gate`), the canary runner, the shift loop, and all
   hooks — unchanged, they reach the gate through the same resolved `gate.sh`.
2. **Contracts.** Subcommand: derives repo root itself (git rev-parse), exit 0
   green / 1 red / 3 not-in-a-repo (today's codes). Env: `BENCH_CANARY_INNER=1`
   → sequential, unprefixed, sweep skipped; shellcheck absent → phase skipped
   (best-effort, unchanged). Outer stdout: phase-prefixed lines + per-phase
   verdict summary + verbatim `gate: green`/`gate: red` final line.
3. **Deep vs thin.** Deep: the phase runner (hides concurrency, multiplexing,
   signals; seam at the subcommand CLI). Thin: the `gate.sh` one-liner and the
   `bench.sh` route — pass-throughs, no seam of their own.
4. **Black-box assertables.** Exit codes per #2 contract; final green/red line
   on each; all four phases still run when one is red (aggregate posture); outer
   output carries each phase's prefix; inner-mode output byte-shape matches
   today's sequential run; a broken-fixture inner gate still goes red with its
   EXPECT substring (canary sweep green).
5. **Gate attachment.** The three-layer pin from #3: re-targeted `gate_entry`
   anchor (conformance), phase-table unit test, existing canary bite. Not
   gate-observable: the wall-time win itself — one before/after measurement
   during the build (numbers in #1), no permanent seam.
6. **Hostile-input owners.** Paths with spaces/globs → the thin `gate.sh` quotes
   `"$root"`; Go exec passes argv unsheltered by a shell. Required tool missing
   → shellcheck skip arm (unchanged); missing go toolchain fails the phase with
   its own honest error. cwd deeper than root → subcommand's own root
   derivation. Interrupt mid-run → process-group kill (#4). Invocation through
   every surface → `bench gate`, `gate-run`, the shift loop, pre-push, and
   canary inner runs all resolve the same `gate.sh` → same subcommand. Absent vs
   empty / trailing-newline / symlink / re-run classes: n/a — no new file
   formats or state written.
7. **Uncertainty flags.** None — no seam needs escalation; spec on the mid tier
   per profile.
8. **Rejected alternatives.** Bash background jobs; generic phase-runner DSL;
   front-phases-only overlap; buffered output replay; concurrent inner mode;
   `--describe` manifest pin; verbatim Go-source grep pin (#2–#4 record why).
9. **Domain watch-outs.** Canary EXPECT matching is substring-based against
   inner gate output — inner mode must stay unprefixed or fixtures break. The
   `gate: green` / `gate: red` lines are load-bearing surface (humans, cache
   diagnostics) — preserve verbatim. The #1 numbers are one 12-core box; the win
   scales with cores and shrinks on small runners, which is acceptable, not a
   blocker.

Dependency order: n/a — single spec.
