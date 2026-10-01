# Move the cleanup verbs onto the verb runner

Blocked by: 2-name-shared-fixtures.md
Writes: internal/worktree/verb_fixture_test.go (new), internal/worktree/clean_classes_test.go, internal/worktree/clean_discard_test.go, internal/worktree/clean_discard_transaction_test.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_landed_hostile_test.go, internal/worktree/clean_landed_test.go, internal/worktree/clean_operand_test.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_command_test.go, internal/worktree/clean_set_outcomes_test.go, internal/worktree/clean_set_refusal_test.go, internal/worktree/clean_set_test.go, internal/worktree/clean_set_wiring_test.go, internal/worktree/clean_unclaimed_test.go, internal/worktree/landed_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_test.go, internal/worktree/live_binary_test.go, internal/worktree/orphan_render_test.go, internal/worktree/orphan_test.go, internal/worktree/pool_reclaim_test.go, internal/worktree/release_inside_test.go, internal/worktree/response_spill_test.go, internal/worktree/resume_reconcile_test.go, internal/worktree/resume_test.go, internal/worktree/subshell_test.go, internal/worktree/unlanded_route_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR19, VR20

## What to build

Move every `release`, `clean`, `resume-clean`, and `reclaim` verb call in the listed files onto the verb runner. A stubbed call passes its joins value through the runner. Read each fingerprint through the verb result's fingerprint reader, and read each table through the rows reader where the test examines rows.

Delete `runCleanup`, `runCleanupWith`, `runDiscard`, `planAndApply`, `runResume`, `runResumeAt`, `mustResumeClean`, `mustReclaim`, `cleanupRowFingerprint`, and `reclaimFingerprint`. Change `landedSetFixture`, `retainedMemberFixture`, `removableSetFixture`, and `refusedUnlandedRelease` to return one named value each, declared in `verb_fixture_test.go`.

The `subshellAt` calls in `subshell_test.go` stay, as the spec's Won't handle states. Keep every test name and every assertion.

## Acceptance

- [ ] The ticket's run wrapper and extractor `rg` prints no line.
- [ ] The tuple scan omits the four cleanup-family builders.
- [ ] The listed files call no verb form directly, except `subshellAt`.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
- [ ] No over-budget test file grows past its base line count.
