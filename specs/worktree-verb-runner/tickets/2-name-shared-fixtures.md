# Name the shared assignment and pool fixtures

Blocked by: 1-add-verb-runner.md
Writes: internal/worktree/verb_fixture_test.go (new), internal/worktree/admin_readers_test.go, internal/worktree/build_test.go, internal/worktree/classifier_shape_test.go, internal/worktree/clean_branch_test.go, internal/worktree/clean_operand_test.go, internal/worktree/eligibility_test.go, internal/worktree/exec_pwd_test.go, internal/worktree/exec_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_reauthorize_operand_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/lifecycle_acquire_test.go, internal/worktree/lifecycle_facts_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_test.go, internal/worktree/list_actions_test.go, internal/worktree/live_binary_test.go, internal/worktree/orphan_render_test.go, internal/worktree/orphan_test.go, internal/worktree/ownership_test.go, internal/worktree/path_identifier_test.go, internal/worktree/pool_reclaim_facts_test.go, internal/worktree/pool_reclaim_test.go, internal/worktree/pool_root_test.go, internal/worktree/recovery_retry_test.go, internal/worktree/release_inside_test.go, internal/worktree/release_registration_test.go, internal/worktree/request_token_test.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_fingerprint_test.go, internal/worktree/reset_plan_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/reset_repair_test.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/reset_restore_test.go, internal/worktree/response_spill_test.go, internal/worktree/resume_test.go, internal/worktree/show_test.go, internal/worktree/unlanded_route_test.go, internal/worktree/worktree_test.go
Covers: VR18

## What to build

Change each shared assignment and pool fixture builder to return one named fixture value. The builders are `newOwnedAssignment`, `newPendingAssignment`, `newOwnedSubmoduleAssignment`, `newResidueGuardFixture`, `unprovableLandedAssignment`, `newReclaimPool`, and `poolRootFixture`. Declare each value type in `verb_fixture_test.go`, and keep each builder in its present file. Update every call site to read the named field it uses.

This ticket moves no verb call. It changes no assertion and no test name. A call site that reads the reset family's runner calls from ticket 1 keeps those calls and reads the new fields.

## Acceptance

- [ ] The tuple scan omits each of the seven builders.
- [ ] Each fixture value type that these builders return is declared in `verb_fixture_test.go`.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`.
- [ ] No over-budget test file grows past its base line count.
