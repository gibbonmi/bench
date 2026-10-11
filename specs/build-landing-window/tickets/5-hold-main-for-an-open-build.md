# Hold main for an open chunked build

Blocked by: 1-hold-the-landing-lease.md, 4-review-a-fold-resolution.md
Writes: .bench/structure-accept, DATA_HANDLING.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_commitment.go, internal/commitment/admission.go, internal/commitment/admission_test.go, internal/commitment/concurrent_test.go, internal/commitment/continuation_test.go, internal/commitment/repository/admission.go, internal/commitment/repository/closure.go, internal/commitment/repository/lease.go (new), internal/commitment/repository/light_path_test.go, internal/commitment/repository/publication.go, internal/commitment/repository/publication_lock_test.go, internal/commitment/repository/publication_test.go, internal/commitment/repository/readiness.go, internal/commitment/repository/sources.go, internal/commitment/repository/window.go (new), internal/commitment/repository/window_test.go (new), internal/commitment/store_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/commitment_record_test.go, internal/conformance/subcommand_routing_table_test.go, internal/intent/ledger/commitment.go, internal/intent/ledger/validate.go, internal/intent/ledger_aliases.go, internal/intent/transaction.go, internal/refusalroute/faces_commitment.go, internal/refusalroute/faces_land.go, internal/systemtest/adoption_test.go, internal/systemtest/otel_verbs_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/worktree/commitment_landing_fixture_test.go, internal/worktree/commitment_light_landing_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/identity_component_test.go, internal/worktree/land.go, internal/worktree/land_bench_home_test.go, internal/worktree/land_broker_notice_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_local_capture_test.go, internal/worktree/land_marker_fixture_test.go, internal/worktree/land_prunes_landed_siblings_test.go, internal/worktree/land_refusal.go, internal/worktree/land_refusal_fixture_test.go, internal/worktree/land_specless_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_tickets_only_test.go, internal/worktree/land_trace_test.go, internal/worktree/land_window_test.go (new), internal/worktree/parallel_census_test.go, internal/worktree/pool_root_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/refusal_route_test.go, internal/worktree/tree_target_test.go, tests/canary/data-handling-derivation/undocumented-passlist-var
Covers: BW46, BW47, BW48, BW49, BW50, BW51, BW52, BW53, BW54, BW55, BW56, BW57, BW58, BW59, BW94, BW96, BW97

## What to build

While a chunked build is open, the landing refuses each landing that is not the build's own, unless the build is in a gap.
The refusal names the build and the gap to wait for.

Put the open-build derivation and the window decision in `internal/commitment/repository/window.go`.
An open build is a claim `(milestone, outcome)` with at least one qualifying binding of that outcome.
A qualifying binding names a staged-spec deliverable and an assignment that is active in `ledger.Assignments`.
The staged-spec test reuses `validateDeliverable` (`internal/commitment/repository/sources.go:90`); write no second classification.

The build's assignments are the active bound assignments of that outcome, a `--from` sibling included.
The build's own landing lands the bound deliverable from a bound assignment.
A spec-less or tickets-only landing from a bound assignment is not the build's own landing.

A publication passes the window when, for each open build, one of these holds:

- The publication is the build's own landing.
- The build's gap record has `Open` true.
- The ledger holds an override for the publication assignment.

Ticket 7 writes the override record; this ticket reads it.
Decide the window inside ticket 1's acquire decision, before the lease, and again in `PublishAdmitted`.
The publish decision is the authoritative one, because it holds the intent lock across the ref update.

The window refusal prints `refused{detail=landing destination refs/heads/<branch> is held by open build <outcome> (<deliverable>); wait for <gap>,next=...}`.
The `<gap>` text derives from the build's gap record:

- No record gives `its first chunk checkpoint`.
- A closed record with `Last` false gives `the checkpoint of chunk <chunk>`.
- A closed record with `Last` true gives `its completion landing`.

When several open builds refuse, the refusal names the first in ledger claim order.
Its face is `land-build-window` with agent authority.
Its route is the instruction `wait for the gap that the refusal names, or ask the reviewer for an override`, then the caller's rerun.
Add a producing fixture for the face.

`ReconcileDelivered` and `SetBlocker` drop the gap record of the outcome whose claim they drop.
`bench commitment start` drops a stale gap record of the outcome that it claims.

Posture change: `specLessLandingFixture` and `ticketsOnlyLandingFixture` (`internal/worktree/land_fixtures_test.go:110,124`) bind a staged spec, so their spec-less and tickets-only landings now refuse.
Change those fixtures so that each test still grades its own subject, and keep each existing assertion.
The callers in this `Writes:` line follow from that change.
`internal/systemtest/adoption_test.go` and `internal/systemtest/otel_verbs_test.go` are system-tagged, so run them with `BENCH_KIT` set through `bench test --check system --run <name>`.
`internal/worktree/land.go` has about 9 lines of headroom after ticket 1; change `.bench/structure-accept` only with a reviewer grant.
Raise `worktreeTestCount` in `internal/worktree/parallel_census_test.go` in place for each new top-level test.

## Acceptance

- [ ] After `bench commitment start` binds a staged spec, an unbound landing refuses with face `land-build-window`, and the refusal names the outcome, the deliverable, and `its first chunk checkpoint`.
- [ ] A closed gap record for chunk `C2` that is not last gives `the checkpoint of chunk C2`, and a closed record with `Last` true gives `its completion landing`.
- [ ] A landing of the bound spec from a `--from` sibling passes the window, and a spec-less landing from that sibling refuses.
- [ ] After the completion landing publishes, after `bench commitment block`, and after the release of every bound assignment, an unbound landing publishes.
- [ ] With one of two bound assignments released, an unbound landing still refuses.
- [ ] A started outcome with a tickets-only deliverable refuses no landing.
- [ ] A build that starts during an unrelated landing's gate refuses that landing at publish, and `main` stays unchanged.
- [ ] A light-path landing with a tickets-only `--spec` refuses with `land-build-window` while a build is outside a gap.
- [ ] With two open builds, a landing inside one build's gap and outside the other's refuses and names the other build.
- [ ] With an open build outside a gap and a live lease holder, the refusal is the window refusal.
- [ ] After a block during an open gap, an unblock, and a new start, an unbound landing refuses.
- [ ] After a block, and after the completion landing, the ledger holds no gap record for the outcome.
