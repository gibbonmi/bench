# 02 Shrink local adoption fixtures while retaining full journeys

Blocked by: 01-build-isolated-repair-kits.md
Writes: internal/adopt/link_transaction_test.go, internal/adopt/setup_prompt_test.go, internal/adopt/setup_test.go
Covers: LTE15, LTE16, LTE17

## What to build

Review chunk: LTE-C2.

Use the accepted selected-kit helper for local unlink and setup scenarios. Preserve all current exit-code, manifest, collision, dry-run, and input assertions.
Keep the existing full-kit default for `consumerRepo` where the relink journey in adopt_test.go needs it. Select minimal kits explicitly for local cases.
Keep full kit assets for `TestSetupGateRejectsIgnoredDeclaredInputs`, which executes the generated gate. Retain the complete system adoption journey.

The helper supplies canonical asset bytes and modes, not a second expected payload manifest. Keep the independent seeded gate-input expectation and demonstrate its omission proof.
Record the callers of consumerRepo, linkConsumer, setupPromptTestRepo, and setupBareTestRepo. Give each caller an explicit small-kit or full-kit disposition.

## Acceptance

- [ ] Local fixture kits contain only the assets needed for their recorded predicates.
- [ ] All unlink scenarios retain their original manifest, collision, dry-run, and status assertions.
- [ ] Operator gate-input bytes survive setup, and a second setup leaves them unchanged.
- [ ] Destination synchronization failure restores prior destinations and removes newly created destinations.
- [ ] The full adoption system journey still exercises the installed wrapper's sentinel failure and later recovery.
- [ ] The generated-gate ignored-input cases still use a complete executable kit.

## Verification

Run the focused link, unlink, and setup tests, including rollback. BENCH_KIT must identify the candidate kit for system verification. Run system tests only through `bench test --check system` during required integration verification. Preserve the seeded-input omission red and the existing full-payload journey.

