# Name the landing fixtures

Blocked by: 7-migrate-merge-and-reauthorize.md
Writes: internal/worktree/verb_fixture_test.go (new), internal/worktree/identity_component_test.go, internal/worktree/land_bench_home_test.go, internal/worktree/land_census_output_test.go, internal/worktree/land_census_test.go, internal/worktree/land_effects_cleanup_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_empty_sibling_test.go, internal/worktree/land_facts_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_folded_base_test.go, internal/worktree/land_freshness_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_journey_test.go, internal/worktree/land_local_capture_test.go, internal/worktree/land_prunes_landed_siblings_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_spec_amendment_test.go, internal/worktree/land_specless_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_tickets_only_test.go, internal/worktree/land_trace_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR36

## What to build

Change each landing fixture builder to return one named fixture value. The builders are `publicLandingFixture`, `publicLandingFixtureAtHome`, `specLessLandingFixture`, `foldedLandingFixture`, `ticketsOnlyLandingFixture`, `landingFixtureAtHome`, `redProspectiveGateLanding`, `brokerChangingLanding`, `brokerDestinationFixture`, `landSurface`, and `foldLandingSibling`. Declare each value type in `verb_fixture_test.go`, keep each builder in its present file, and update every call site.

This ticket moves no verb call. It changes no assertion and no test name. Tickets 9 and 10 then move the landing verb calls onto the runner with these values. The over-budget files `identity_component_test.go` and `land_journey_test.go` stay at or below their base line counts.

## Acceptance

- [ ] The tuple scan omits each of the eleven landing builders.
- [ ] Each landing fixture value type is declared in `verb_fixture_test.go`.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`.
- [ ] No over-budget test file grows past its base line count.
