# Name the assignment fixtures

Blocked by: 2-move-reset-family.md
Writes: internal/worktree/verb_fixture_test.go (new), internal/worktree/admin_readers_test.go, internal/worktree/build_test.go, internal/worktree/classifier_shape_test.go, internal/worktree/clean_branch_test.go, internal/worktree/clean_operand_test.go, internal/worktree/eligibility_test.go, internal/worktree/exec_pwd_test.go, internal/worktree/exec_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_reauthorize_operand_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/lifecycle_acquire_test.go, internal/worktree/lifecycle_facts_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_test.go, internal/worktree/list_actions_test.go, internal/worktree/orphan_test.go, internal/worktree/ownership_test.go, internal/worktree/path_identifier_test.go, internal/worktree/recovery_retry_test.go, internal/worktree/release_inside_test.go, internal/worktree/release_registration_test.go, internal/worktree/request_token_test.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_fingerprint_test.go, internal/worktree/reset_plan_test.go, internal/worktree/reset_refusal_test.go, internal/worktree/reset_repair_test.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/reset_restore_test.go, internal/worktree/response_spill_test.go, internal/worktree/resume_test.go, internal/worktree/show_test.go, internal/worktree/unlanded_route_test.go, internal/worktree/worktree_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR26

## What to build

Change `newOwnedAssignment`, `newPendingAssignment`, and `newOwnedSubmoduleAssignment` to return one named fixture value each. Declare the value types in `verb_fixture_test.go`, keep the builders in `resume_test.go`, and update every call site to read the named field it uses. A fixture value may carry a method that builds a verb call value from its root and home.

This ticket moves no verb call. It changes no assertion and no test name. One value read replaces each tuple read. Thus each over-budget file in the spec's table for ticket 3 stays at or below its base line count.

## Acceptance

- [ ] The tuple scan omits the three assignment builders.
- [ ] Each value type that these builders return is declared in `verb_fixture_test.go`.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`.
- [ ] No over-budget test file grows past its base line count.
