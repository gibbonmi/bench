# Move the cleanup verbs onto the verb runner

Blocked by: 4-name-pool-and-residue-fixtures.md
Writes: internal/worktree/verb_fixture_test.go (new), internal/worktree/clean_classes_test.go, internal/worktree/clean_discard_test.go, internal/worktree/clean_discard_transaction_test.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_landed_hostile_test.go, internal/worktree/clean_landed_test.go, internal/worktree/clean_operand_test.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_command_test.go, internal/worktree/clean_set_outcomes_test.go, internal/worktree/clean_set_refusal_test.go, internal/worktree/clean_set_test.go, internal/worktree/clean_set_wiring_test.go, internal/worktree/clean_unclaimed_test.go, internal/worktree/landed_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_test.go, internal/worktree/live_binary_test.go, internal/worktree/orphan_render_test.go, internal/worktree/orphan_test.go, internal/worktree/pool_reclaim_test.go, internal/worktree/release_inside_test.go, internal/worktree/response_spill_test.go, internal/worktree/resume_reconcile_test.go, internal/worktree/resume_test.go, internal/worktree/subshell_test.go, internal/worktree/unlanded_route_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR28, VR29, VR30

## What to build

Move every `release`, `clean`, `resume-clean`, `reclaim`, `list`, and `path` verb call in the listed files onto the verb runner. A stubbed call passes its joins value through the runner. Read each fingerprint with `mustFingerprint` and each examined table with `mustRows`.

Delete `runCleanup`, `runCleanupWith`, `runDiscard`, `planAndApply`, `runResume`, `runResumeAt`, `mustResumeClean`, `mustReclaim`, `cleanupRowFingerprint`, `reclaimFingerprint`, and the test-side rows reader `cleanupRows`. Each `cleanupRows` caller reads the `worktree_cleanup` table with `mustRows`. Change `landedSetFixture`, `retainedMemberFixture`, `removableSetFixture`, and `refusedUnlandedRelease` to return one named value each, declared in `verb_fixture_test.go`.

Replace each inline fingerprint match in these files:

- `clean_set_test.go`: the `setFingerprint` variable and every use of it in the clean files
- `clean_landed_test.go`: the three inline 64-hex matches
- `pool_reclaim_test.go`: the apply-invocation match inside `reclaimFingerprint`
- `orphan_render_test.go`: the aggregate-row match in `planReclaimableCount`, which reads the count with `mustRows`
- `clean_operand_test.go`: the inline 64-hex match

A positive read becomes `mustFingerprint`. A check that an error plan carries no fingerprint becomes `mustNoFingerprint`. The `subshellAt` calls in `subshell_test.go` stay, as the spec's Won't handle states. Keep every test name and every assertion. The over-budget files `lifecycle_test.go`, `pool_reclaim_test.go`, and `resume_test.go` stay at or below their base line counts.

## Acceptance

- [ ] The VR28 command prints no line.
- [ ] The tuple scan omits the four cleanup-family builders.
- [ ] The verb form command over the listed files prints no line.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
- [ ] No over-budget test file grows past its base line count.
