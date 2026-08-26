# Prove the verb journeys on the new shape

Blocked by: 04-print-the-green-phase-table.md, 05-write-the-complete-phase-stream-to-a-run-log-file.md
Writes: internal/gate/phases_test.go (new), internal/gate/run_outcomes_test.go, internal/commit/landing_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_journey_test.go, internal/systemtest/adoption_test.go

## What to build

`bench gate`, `bench commit`, and `bench worktree land` read the same. This
ticket proves the composition at three seams and changes no production code.

Read one tree fact first. The engine's report runs only when the resolved
gate script execs `bench gate-phases`. `phaseTableGate` demands that exec
line in `internal/gate/gate.go`, and `runCaptured` in
`internal/gate/run_transaction.go` relays any other script's stream. The exec
path also demands a sealed run binary through `runbinary.Inherit`, which no
unit test holds.

So the engine's own rows attach in-package. Call `phasesCommandAtKitWithSelection` with a constructed `runbinary.Selection`
and fixture phase scripts, and read its stdout, stderr, exit code, and
progress log. The run log starts one level up, so the test stands it up
itself through `beginGateRunLog` and the overridable `gateLogPathIgnored`
variable in `internal/gate/run_log.go`.

The verb journeys prove the relay, not the engine. Each fixture gate script
in `internal/commit/landing_test.go` and
`internal/worktree/land_fixtures_test.go` prints a canned bounded shape. That is a `failures[1]{phase,line}` table with `gate: red`, or a `phases[1]{...}` table
with `capability-skips: 0` and `gate: green`. The journey asserts the verb's
stdout carries those bytes unchanged, then the verb's own record. The commit
refusal record and the landing records stay exactly as they are.

Read the second tree fact. `internal/systemtest/adoption_test.go` pins
`gate: green` on stdout, and the scaffolded `.bench/gate.sh` in
`internal/adopt/init.go` prints that line itself. That script never reaches
the engine, so the pin holds unchanged; confirm it by running the test.

## Acceptance

- [ ] A red run of fixture phase scripts exits 1. (BG13)
- [ ] A green run of fixture phase scripts exits 0. (BG34)
- [ ] A cancelled run of fixture phase scripts exits 130. (BG36)
- [ ] A second run on an unchanged tree prints `gate: green (fresh verdict reused for this tree)` alone. (BG20)
- [ ] `bench commit` on a fixture gate that prints a canned red table relays that table, `gate: red`, and its own refusal record. (BG21)
- [ ] `bench worktree land` on a fixture gate that prints a canned green shape relays it before `landed{...}`. (BG22)
- [ ] A phase's `elapsed_ms` cell equals the `elapsed_ms` of its `phase.finish` record. (BG23)
- [ ] The run-outcome test that pins the reused line and the adoption test that pins `gate: green` pass unchanged. (BG28)
