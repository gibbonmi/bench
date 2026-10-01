# Move the landing effect and resume tests onto the verb runner

Blocked by: 8-name-landing-fixtures.md
Writes: internal/worktree/land_effects_cleanup_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_empty_sibling_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_census_test.go, internal/worktree/land_census_output_test.go, internal/worktree/land_trace_test.go, internal/worktree/land_journey_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR39

## What to build

Move every `land` and `land-resume` verb call in the listed files onto the verb runner, together with each joins form call. Each stubbed landing passes its joins value through the runner, so each effect stub still runs. `resumeLandArgs` builds an argument list and stays.

The folded-sibling landing in `land_effects_cleanup_test.go` checks that stdout carries no cleanup fingerprint with an inline 64-hex match. That check becomes `mustNoFingerprint`, and its `--apply` check stays. The VR-C4 review logs the narrower check as an accepted drop under VR46.

Keep every test name and every assertion. The over-budget file `land_journey_test.go` stays at or below its base line count.

## Acceptance

- [ ] The VR39 command over the listed files prints no line.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
- [ ] No over-budget test file grows past its base line count.
