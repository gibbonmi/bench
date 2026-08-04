# Bound the abandoned repair child and grade the process tree

Blocked by: none
Ownership fence: `bin/bench-repair-binary.mjs`, `internal/contract/command.go`, `internal/contract/helper_test.go`, `internal/contract/surface/binary_repair_hardened_test.go`, `internal/contract/runtime/runtime_gate_partial_proof_test.go`
Contracts: the process-group-drained predicate crosses `internal/contract`→`internal/contract/surface/binary_repair_hardened_test.go`, asserted by RO2 and RO3 against a real abandoned `bench.sh repair` process group; the release deadline crosses `bin/bench-repair-binary.mjs`→`internal/contract/surface/binary_repair_hardened_test.go` as `BENCH_TEST_REPAIR_RELEASE_DEADLINE_MS`, asserted by RO1
Assumptions: `waitForTestRelease` is reached only under a `BENCH_TEST_` marker variable, so its deadline is test-path-only and needs no production default change; `bin/bench.sh` runs node in the foreground, so a drained group after `cmd.Wait()` is the normal-path expectation

## What to build

`bench repair`'s node child must not outlive whatever started it, and the binary-repair
contract must grade that.

`waitForTestRelease` in `bin/bench-repair-binary.mjs` polls at 25ms with no bound. It ends
only when a Go test deletes the marker file. When the test binary dies first — a `go test`
timeout, an interrupted gate, a cancelled CI job — the node process spins forever: one was
observed alive 14h36m, reparented to `/init`, with its `t.TempDir()` fixture still on disk.
`testRepairLosingRacerPreservesWinner` puts that child in its own process group via
`Setpgid` and then only ever `cmd.Wait()`s the direct child, so nothing signals the group
on an early `t.Fatal` either.

Bound the wait with a deadline in the shape the file already uses at the fetch path
(`testNumber("BENCH_TEST_REPAIR_DEADLINE_MS", FETCH_DEADLINE_MS)`), under its own variable
so the fetch deadline the earliest-interrupt test already sets to 5s does not also shorten
the handshake. Expiry throws, so the existing catch removes the temp file and exits 1.

Both hardened repair tests then start their child through one shared helper that registers
a `t.Cleanup` group signal, so a failure before the test's own teardown cannot leak the
group it created.

The contract itself is the other half. `the loser left nothing behind` is currently graded
by a `.bench-*.tmp` glob alone — filesystem residue standing in for the whole invariant,
the same defect class as the assertion that cited residue instead of what it named. Grade
the process group too, and add the case that can actually go red: a repair whose owner dies
mid-handshake must still drain.

The drained-group predicate already exists twice — `processGroupAlive` in
`internal/contract/helper_test.go` and `waitR17ProcessGroupExit` in
`internal/contract/runtime/runtime_gate_partial_proof_test.go`. Land one exported helper in
`internal/contract` and move both onto it rather than a third derivation of the same fact.

## Acceptance

- [ ] [RO1] A repair parked in the test-release handshake exits on its own once the release deadline passes, instead of polling forever, and reports the refusal on stderr.
- [ ] [RO2] The losing-racer contract fails when the repair leaves any process in the group it was started in, not only when it leaves a `.bench-*.tmp` file; the existing temp-file assertion still fails on a leaked temp.
- [ ] [RO3] A repair whose owning process is killed mid-handshake leaves no surviving process in its group.
- [ ] [RO4] Every hardened repair test signals the process group it created even when the test fails before reaching its own teardown.
- [ ] [RO5] One exported `internal/contract` helper is the only derivation of the group-drained predicate; the two prior copies call it and grade what they graded before.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RO1 | restore the unbounded `while (true)` poll in `waitForTestRelease` | the abandoned-child contract test | revert the deadline, run `go test ./internal/contract/surface -run 'TestBinaryRepairContracts/repair_abandoned-child'`, expect the surviving-group failure |
| RO2 | leave a lingering member in the repair process group (`sleep 60 &` before the node call in `bin/bench.sh`) | the losing-racer contract test | add the background call, run `go test ./internal/contract/surface -run 'TestBinaryRepairContracts/repair_losing-racer'`, expect the surviving-group failure after the temp glob assertion has already passed |
| RO3 | remove the release deadline's `throw` and return instead | the abandoned-child contract test | swallow the expiry, run the abandoned-child test, expect the surviving-group failure |
| RO4 | delete the `t.Cleanup` group signal from the shared start helper | the repro loop `repro-orphan.sh` | force an early `t.Fatal` after the markers, run the subtest, expect an orphaned `bench-repair-binary.mjs` in `ps -g <pgid>` |
| RO5 | signal only the direct child in the fixture runner's `terminateProcessGroup` (`Kill(pid)` for `Kill(-pid)`) | `TestFixtureRunnerReapsSpawnedProcessGroup`, the first migrated caller | narrow the signal, run `go test ./internal/contract -run TestFixtureRunnerReapsSpawnedProcessGroup`, expect the surviving-group failure naming the escaped child. Mutating the shared predicate itself proves nothing — an always-drained answer is invisible to a suite whose subjects drain — so the migration is graded through each caller's own subject. |
