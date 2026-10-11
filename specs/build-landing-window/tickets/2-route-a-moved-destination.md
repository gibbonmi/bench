# Route a moved landing destination to the fold and the rerun

Blocked by: 1-hold-the-landing-lease.md
Writes: .bench/structure.budgets, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/commit/commit.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/authorization/authorization.go, internal/gate/authorization/authorization_test.go, internal/gate/checkpoint.go, internal/gate/complete_checkpoint.go, internal/gate/gate.go, internal/gate/run_transaction.go, internal/gate/verdict.go, internal/landing/completion_evidence_test.go, internal/landing/gitexec.go, internal/landing/landing.go, internal/landing/landing_reviewed_test.go, internal/landing/landing_test.go, internal/landing/merge.go, internal/landing/state_test.go, internal/refusalroute/faces_land.go, internal/refusalroute/registry.go, internal/systemtest/owner_artifact_recovery_test.go, internal/systemtest/owner_land_race_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/worktree/identity_component_test.go, internal/worktree/land.go, internal/worktree/land_identity_test.go, internal/worktree/land_journey_test.go, internal/worktree/land_moved_destination_test.go (new), internal/worktree/land_refusal.go, internal/worktree/land_refusal_fixture_test.go, internal/worktree/land_surface_test.go, internal/worktree/merge.go, internal/worktree/merge_refusal.go, internal/worktree/parallel_census_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/refusal_route_test.go
Covers: BW17, BW18, BW19, BW20, BW21, BW22, BW23, BW24

## What to build

A landing over a moved destination refuses with one face, `land-destination-moved`, and its route names the fold and the rerun, never `bench doctor`.
Three refusals take the new face:

- the post-gate recheck at `internal/landing/landing.go:283-285`, when the destination ref no longer equals the destination base, with its detail text unchanged
- a `TipMovedError` from the compare-and-swap at publish (`internal/landing/gitexec.go:22-31`)
- an authorization result with the completion proof cause, when the destination base is not an ancestor of the source tip

A post-gate checkout change with an unmoved ref keeps `land-handback`.
A completion proof cause over a source that holds the destination base also routes to `land-handback`.

The face has agent authority and the sentence `landing destination moved past the reviewed source`.
Its route has three steps: `bench worktree merge --from <branch> <label>`, the instruction `for a spec landing, record the fold, its review, the merged-test verification, and the completion at the new tip`, and the caller's rerun at the repaired tip.
Add the face to `internal/refusalroute/faces_land.go` and its constant to `internal/worktree/land_refusal.go`, after ticket 1's lease face.
Add one producing fixture to `internal/worktree/land_refusal_fixture_test.go`, and make the route walk carry the instruction step.

The gate carries the proof cause as a type.
`funnelFace` already maps `completionProofError` to its face (`internal/gate/checkpoint.go:152-153`).
Outside a verb context, the gate result records that face, and `authorization.Result` copies it into a typed cause that the landing reads.
No reason text decides the route.
Prefer a field on `gate.Result` to one on `Inspection`, because `internal/gate/verdict.go` is already over its budget.

Several owners have little headroom: `internal/worktree/land.go` about 9 lines, `internal/landing/landing.go` about 8 lines, and `internal/commit/commit.go` none.
Split new logic into a new file before any of these grows past its budget.
Change `.bench/structure.budgets` only with a reviewer grant.
`internal/systemtest/owner_artifact_recovery_test.go` and `internal/worktree/land_journey_test.go` are over budget, so a change there adds no net line.

BW20 extends `TestWorktreeLandPublicRaceAndRerun` in `internal/systemtest/owner_land_race_test.go`, a system-tagged file.
Run it with `BENCH_KIT` set through `bench test --check system --run TestWorktreeLandPublicRaceAndRerun`.
`TestLandSuppliesTheDeclaredBenchHomeItself` keeps its `bench doctor` route, because its cause is infrastructure.
`TestLandCommandRefusesPostGateUnknownIgnoredMutation` keeps `land-handback`.
Raise `worktreeTestCount` in `internal/worktree/parallel_census_test.go` in place for each new top-level test.

## Acceptance

- [ ] A completion landing over a destination with one new commit refuses with face `land-destination-moved`.
- [ ] That refusal's `next=` holds `bench worktree merge --from main` and ends with the caller's rerun.
- [ ] That refusal's output holds no `bench doctor`.
- [ ] A raw commit on the destination during the loser's gate gives a `next=` that holds `bench worktree merge --from main`.
- [ ] A destination that moves inside the publication gap gives a `TipMovedError` refusal with face `land-destination-moved`.
- [ ] An ignored-file change in the destination checkout during the gate, with an unmoved ref, keeps face `land-handback`.
- [ ] A completion proof fault over a source that holds the destination base gives face `land-handback`.
- [ ] `authorization.Result` carries the completion proof cause as a typed field when the gate refuses a completion proof.
