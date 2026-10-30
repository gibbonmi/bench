# Cut the subsumed refusal-route tests

Blocked by: 04-cut-census-effects-and-freshness-tests.md
Writes: .bench/commitment.json, internal/worktree/identity_component_test.go, internal/worktree/refusal_route_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/land_refusal_fixture_test.go, internal/worktree/parallel_census_test.go
Covers: TC35, TC36, TC37, TC38, TC39, TC40

## What to build

The worktree package loses three refusal-route tests that the follow walk and the lease test subsume. The spec rows WT-1, WT-5, and WT-11 own this ticket. Ticket 04 also edits the count pin, so this ticket starts from its commit.

- Remove `TestLandingRefusalRegistryHasAProducingFixture` (WT-1) and `TestLandingRedRouteNamesTheRepair` (WT-11). Their kept tests already bit at the spec stage, and WT-11 needs no fold.
- WT-5 is a move row. First add the spec's list of the seven reviewer land faces to `TestLandingFacesFollowTheirRoutes` as an independent list, and assert that each listed face's printed route opens with the reviewer marker. Keep the list independent of the registry. Probe `tc37` on the walk and see `bit`, then remove `TestReviewerLandFacesOpenWithTheMarker`.
- Point the comment at `land_refusal_fixture_test.go:36` that names the WT-1 candidate at `TestLandingFacesFollowTheirRoutes`.
- Lower `worktreeTestCount` by 3 from its value at this ticket's base. Add no line to `parallel_census_test.go`.
- Remove each helper and import that only a removed test used. The `refusalroute` import of `identity_component_test.go` goes with the WT-1 test, because that test is its only user.
- `.bench/commitment.json` is in `Writes:` only because the spec-stage close commits the delivery binding inside the graded range. This ticket does not edit it.

The build changes no production file. Run each probe through `bench probe` with `--expect` on the kept test, and read the `restored` cell of each answer.

## Acceptance

- [ ] The three candidate tests no longer exist.
- [ ] The completion-plan probes `tc36`, `tc39`, and `tc40` read `bit` on `TestLandingFacesFollowTheirRoutes`.
- [ ] The completion-plan probe `tc37` reads `bit` on `TestLandingFacesFollowTheirRoutes` after the move.
- [ ] The completion-plan probe `tc38` reads `bit` on `TestLandReleasesTheLeaseOnRefusal`.
- [ ] `TestPackageTestCountPin` passes with `worktreeTestCount` equal to the live top-level count, and `parallel_census_test.go` does not grow.
- [ ] `bench test --package ./internal/worktree` passes.
