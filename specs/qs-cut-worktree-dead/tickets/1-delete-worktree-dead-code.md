# Delete the dead worktree code and the eligibility forwarders

Blocked by: none
Writes: internal/worktree/resume.go, internal/worktree/build_outputs.go, internal/worktree/effects.go, internal/worktree/land.go, cmd/bench/worktree_leaves.go, internal/worktree/eligibility.go, internal/worktree/subshell.go, internal/worktree/classifier.go, internal/worktree/clean_landed.go, internal/worktree/eligibility_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_census_output_test.go, internal/worktree/land_census_test.go, internal/worktree/land_effects_cleanup_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_empty_sibling_test.go, internal/worktree/land_facts_test.go, internal/worktree/land_flags_test.go, internal/worktree/land_folded_base_test.go, internal/worktree/land_freshness_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_journey_test.go, internal/worktree/land_local_capture_test.go, internal/worktree/land_prunes_landed_siblings_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_reauthorize_operand_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_specless_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_tickets_only_test.go, internal/worktree/land_trace_test.go
Covers: none

## What to build

The quality survey of 2026-09-29 lists small certain cuts in `internal/worktree`. This ticket removes the dead code that the survey found, and it changes no CLI output and no exit code.

These three functions have no caller in production code or in tests:

- `ConservativeCleanup` in `resume.go`. The resume path calls `conservativeCleanupAt` directly.
- `ignoredWithinDeclaredOutputs` in `build_outputs.go`. The landing allowance uses `ignoredWithinLandingAllowance`.
- `landingAlreadyRebuilt` in `effects.go`. The landing never rebuilds, so the `rebuiltLandingEnv` constant in `land.go` has no other reader.

Delete the three functions and the `rebuiltLandingEnv` constant. Move the useful text of the `ConservativeCleanup` comment onto `conservativeCleanupAt`.

`landAttributed` ignores its executable parameter. Remove the parameter from `LandCommand`, `landWith`, and `landAttributed`. Remove the argument from the `bench worktree land` leaf in `cmd/bench` and from each test call.

The file `eligibility.go` holds unexported aliases and two one-line forwarders to `internal/worktree/lifecyclepolicy`. These aliases give no new name and no narrower type. Other files in the package already call `lifecyclepolicy` directly. Replace each alias and each forwarder with its `lifecyclepolicy` name at every use. Keep `explicitOutcome` and `automaticPreservationVerdict`, because each one projects a `CleanupPlan` into a policy input.

## Acceptance

- [ ] No Go source names `ConservativeCleanup`, `ignoredWithinDeclaredOutputs`, `landingAlreadyRebuilt`, or `rebuiltLandingEnv`.
- [ ] `LandCommand`, `landWith`, and `landAttributed` take no executable parameter, and the `bench worktree land` leaf passes none.
- [ ] `eligibility.go` declares no alias and no forwarder to `lifecyclepolicy`, and it keeps `explicitOutcome` and `automaticPreservationVerdict`.
- [ ] `go vet ./...` passes.
- [ ] The `internal/worktree` and `cmd/bench` package tests pass.
