# Cut the subsumed repair tests

Blocked by: none
Writes: internal/adopt/repairtest/repair_test.go, specs/landing-test-efficiency/spec.md
Covers: TC2, TC3, TC4, TC5, TC6, TC7, TC8, TC9, TC43

## What to build

The repair suite loses three tests that a broader test subsumes, and one merged test replaces three boundary tests. The spec rows T1-3, T1-4, and T1-5 own this ticket.

- Remove `TestCompatibilityUndoCreated` and `TestCompatibilityManagedRepair`. `TestCompatibilityUndo` stays unchanged and is the kept test for both. The spec-stage probes bit on it with no move.
- Replace `TestCompatibilityPolicyBoundary`, `TestCompatibilityInterruptionBoundary`, and `TestCompatibilityRuntimeBoundary` with `TestCompatibilityRecoveryBoundaries`. The new test runs one linked session and one `doctor --fix`. It plants the combined `config.toml` and the `trust.json` that the spec's Merges section states, and it starts the supervised process before the repair. After the repair it keeps every assertion of the three removed tests.
- In `specs/landing-test-efficiency/spec.md`, point the seam cells of LTE4 and LTE9 at `TestCompatibilityUndo`, and the seam cells of LTE12, LTE13, and LTE14 at `TestCompatibilityRecoveryBoundaries`. Change no other byte of that spec. If the drain retired that spec before this ticket, skip this step.
- In `specs/subsumed-test-cuts/spec.md`, rewrite the `planned` seam cells of TC4 to TC7 and TC43 into the citation form.
- Remove each helper that only a removed test used.

The build changes no production file. Run each probe through `bench probe` with `--expect` on the kept test, and read the `restored` cell of each answer.

## Acceptance

- [ ] `TestCompatibilityUndoCreated`, `TestCompatibilityManagedRepair`, and the three boundary tests no longer exist, and `TestCompatibilityRecoveryBoundaries` exists.
- [ ] The completion-plan probes `tc2` and `tc3` read `bit` on `TestCompatibilityUndo`.
- [ ] The completion-plan probes `tc4`, `tc5`, `tc6`, `tc7`, and `tc43` read `bit` on `TestCompatibilityRecoveryBoundaries`.
- [ ] The merged test asserts that the supervised process answers after the repair (TC8, review-owned).
- [ ] `bench test --check docs-currency-workflow` passes with the re-pointed LTE cells.
- [ ] `bench test --package ./internal/adopt/repairtest` passes.
