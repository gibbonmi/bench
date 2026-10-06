# 3. Publish an unbound light-path landing that names its folder

Blocked by: 2-admit-light-path-commit.md
Writes: internal/commitment/repository/light_path.go (new), internal/commitment/repository/light_path_test.go (new), internal/commitment/repository/candidate.go, internal/commitment/repository/publication.go, internal/spec/tickets_only.go, internal/commitment/commitmenttest/, internal/worktree/commitment_light_landing_test.go, internal/worktree/commitment_landing_fixture_test.go, internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: LP3, LP4, LP5, LP6, LP24, LP25, LP26, LP50, LP52, LP53

## What to build

This ticket is the third ticket of review chunk LP-C1. It delivers publication
mode: `bench worktree land --spec <slug>` publishes an unbound light-path change
and closes its folder. It calls the tree reader and the predicate from ticket 2,
and it adds no second copy of the qualify rule or the cover rule.

`admitPublication` grades in publication mode. `published.Tree` removes the
`--spec` folder from the graded tree before admission, so the graded tree cannot
supply the ticket. Publication mode therefore reads only the folder that
`Delivery.Spec` names, at the reviewed source commit `Delivery.Source`. The
folder qualifies through the reader of ticket 2. Its ticket must cover every
production path of the graded tree. The Writes refusal and the span refusal
keep their words.

A publication with no delivery is a spec-less landing, and it is never
light-path work. When the composed tree holds a qualifying folder that covers
every production path, its refusal is the `readyFor` refusal, then
`; land the light-path change with --spec %q`. The operand is the folder slug.
Each other spec-less case keeps its current refusal.

A bound assignment stays admitted through `readyFor`, whatever folders the tree
holds. A delivery that names a staged spec keeps the current rules.
`Store.Closure` returns no fact for a folder that no milestone approves, so the
landing leaves `.bench/commitment.json` unchanged.

The landing rows share one new top-level test after the precedent of
`TestCommitmentTicketsOnlyClosure`, so `worktreeTestCount` in
`internal/worktree/parallel_census_test.go` moves by one. Put any new landing
route in `internal/worktree/commitment_landing_fixture_test.go`. Ticket 4
writes the `CHANGELOG.md` entry for the whole change.

## Acceptance

- [ ] `TestCommitmentLightPathLanding` in `internal/worktree/commitment_light_landing_test.go` shows that `bench worktree land --spec <slug>` of an unbound light-path assignment moves the `main` ref to the landing commit (LP3).
- [ ] The same test shows that the published tree holds no `specs/<slug>` entry (LP4).
- [ ] The same test shows that the published `.bench/commitment.json` equals its bytes at the landing base (LP5).
- [ ] The same test shows that a spec-less landing of that change exits nonzero with output that contains `--spec` (LP6).
- [ ] `TestLightPathPublication` in `internal/commitment/repository/light_path_test.go` shows that `Store.AdmitPublication` with a delivery that names the folder at its source admits a composed tree that lacks the folder (LP24).
- [ ] The same test shows that an uncovered production path refuses with the path named (LP25).
- [ ] The same test shows that a publication with no delivery refuses a covered light-path tree with `--spec` in the refusal (LP26).
- [ ] The same test shows that the refusal for the folder slug `a b` contains `--spec "a b"` (LP53).
- [ ] The same test shows that a delivery that names a two-ticket folder refuses with `assignment has no current delivery binding` (LP52).
- [ ] The same test shows admission of a bound spec delivery whose composed tree holds a qualifying folder that does not cover its paths (LP50).
- [ ] `bench test --package ./internal/commitment/repository` and `bench test --package ./internal/worktree` pass.
