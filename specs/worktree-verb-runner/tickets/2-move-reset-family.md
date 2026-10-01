# Move the reset family onto the verb runner

Blocked by: 1-add-verb-runner.md
Writes: internal/worktree/verb_fixture_test.go (new), internal/worktree/reset_apply_test.go, internal/worktree/reset_fingerprint_test.go, internal/worktree/reset_plan_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/reset_repair_test.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/reset_restore_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR23, VR24, VR25

## What to build

Move every reset verb call in the listed files onto the verb runner. A stubbed reset call passes its joins value through the runner. Read each fingerprint with `mustFingerprint`.

Delete `runReset`, `runResetWith`, `resetFingerprint`, and `restoreFingerprint`. Add `verb_fixture_test.go` with the named value that `restoreFixture` returns, and update its callers. The `mergeFixture` call in `reset_repair_test.go` stays for ticket 7.

Keep every test name and every assertion. Record the tuple scan program in the review record.

## Acceptance

- [ ] The VR23 command prints no line.
- [ ] The tuple scan omits `restoreFixture`.
- [ ] The verb form command over the listed reset files prints no line.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
