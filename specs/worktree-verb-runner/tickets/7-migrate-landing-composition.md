# Move the landing composition tests onto the verb runner

Blocked by: 6-name-landing-fixtures.md
Writes: internal/worktree/land_flags_test.go, internal/worktree/land_specless_test.go, internal/worktree/land_tickets_only_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_freshness_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_facts_test.go, internal/worktree/land_folded_base_test.go, internal/worktree/land_local_capture_test.go, internal/worktree/land_prunes_landed_siblings_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_reauthorize_operand_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/identity_component_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR25

## What to build

Move every `land`, `land-resume`, `release`, `reauthorize`, and `path` verb call in the listed files onto the verb runner. Each stubbed landing passes its joins value through the runner, so each landing stub still runs. Delete `landIn`. Move `interruptLandingAtMarker` and `landingFaceResume` onto the runner; they stay as scenario helpers because each one drives a scenario and not only a verb call.

Keep every test name and every assertion.

## Acceptance

- [ ] The ticket's `rg` for `landIn` prints no line.
- [ ] The listed files call no verb form directly.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
- [ ] No over-budget test file grows past its base line count.
