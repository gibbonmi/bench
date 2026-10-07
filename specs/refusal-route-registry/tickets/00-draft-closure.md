# Draft closure probe

Blocked by: none
Writes: internal/refusalroute/registry.go (new), internal/refusalroute/route.go (new), internal/refusalroute/command.go (new), internal/refusalroute/registry_test.go (new), internal/refusalroute/route_test.go (new), internal/refusalroute/command_test.go (new), internal/conformance/refusal_route_guard_test.go (new), internal/conformance/refusal_route_bypass_test.go (new), internal/worktree/land_refusal.go, internal/worktree/land.go, internal/worktree/land_identity.go, internal/worktree/land_resume.go, internal/worktree/land_rerun.go, internal/worktree/merge.go, internal/worktree/build.go, internal/worktree/reset.go, internal/worktree/reset_apply.go, internal/worktree/reset_restore.go, internal/worktree/path.go, internal/worktree/classifier.go, internal/worktree/identity_component.go, internal/worktree/refusal_route_test.go (new), internal/worktree/merge_route_test.go (new), internal/worktree/commit_route_test.go (new), internal/worktree/identity_component_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_tickets_only_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_folded_base_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_journey_test.go, internal/worktree/merge_test.go, internal/landing/landing.go, internal/landing/merge.go, internal/landing/landing_reviewed_test.go, internal/commit/commit.go, internal/commit/landing_test.go, internal/commit/dry_run_test.go, internal/commit/refusal_route_test.go (new), internal/gate/checkpoint.go, internal/gate/run_transaction.go, internal/gate/gate.go, internal/gate/run_outcomes_test.go, internal/gate/review_checkpoint_test.go, internal/gate/refusal_route_test.go (new), internal/commitment/commitcmd/command.go, internal/commitment/commitcmd/admission.go, internal/commitment/commitcmd/refusal_route_test.go (new), internal/commitment/verification_test.go, cmd/bench/main.go, cmd/bench/commitment_test.go, cmd/bench/help_inventory_test.go, .bench/BENCH-reference.md, CHANGELOG.md
Covers: none

## What to build

This draft ticket probes the write closure. The slice replaces it.

## Acceptance

- [ ] The probe lists the closure.
