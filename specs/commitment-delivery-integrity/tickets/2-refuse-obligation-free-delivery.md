# 2. Refuse an obligation-free delivery at the completion landing

Blocked by: 1-refuse-obligation-free-plan.md
Writes: internal/commitment/delivery.go, internal/commitment/delivery_test.go, internal/commitment/repository/closure_test.go, internal/worktree/commitment_light_landing_test.go, internal/worktree/commitment_landing_fixture_test.go, internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: FD6, FD7, FD8, FD10, FD11

## What to build

This ticket is the second ticket of review chunk FD-C1. It delivers the landing
refusal of a legacy obligation-free binding. It calls the predicate that
ticket 1 adds to `delivery.go`, and it adds no second copy of that rule.

`Deliver` returns an error for an approved obligation-free binding of the
active milestone. The error names the deliverable path, the outcome, the words
`names no obligation`, and the command `bench commitment plan --input <file>`.
`Deliver` keeps the unchanged return with no error for a path that the active
milestone does not approve. A binding that only an inactive milestone approves
is such a path. Update the `Deliver` doc comment to state the refusal.

`Store.Closure` returns the `Deliver` error unchanged. Its callers
`published.Tree`, `closedPolicy`, and `completionClosure` then refuse, and none
of them needs an edit. In the landing, `published.Tree` runs before the
admission check and before the gate authorization. The landing therefore
refuses before the gate and before the publication, and the `main` ref stays
at its base. No candidate-controlled executable runs before this refusal.

The legacy fixture is the closure spec that `SeedTicketsOnly` approves with no
obligation. `bench commitment start` still admits that binding, so the landing
fixture can bind its assignment to the spec. The landing test drives a
`deliveryRoute` with that spec as its deliverable and `SeedTicketsOnly` as its
seed. Put any new route in `internal/worktree/commitment_landing_fixture_test.go`.

## Acceptance

- [ ] `TestCommitmentLandingRefusesObligationFreeDelivery` in `internal/worktree/commitment_light_landing_test.go` shows that `bench worktree land --spec` of that spec exits nonzero and leaves the `main` ref at its base (FD6).
- [ ] The same test shows that the landing output contains `bench commitment plan --input <file>` (FD7).
- [ ] `TestCommitmentClosureRefusesObligationFreeBinding` in `internal/commitment/repository/closure_test.go` shows that `Store.Closure` returns an error that contains `names no obligation` for that spec (FD8).
- [ ] `TestCommitmentDeliver` shows that `Deliver` returns no error and no new fact for `specs/other/spec.md`, and for an obligation-free binding that only an inactive milestone approves (FD10).
- [ ] `TestCommitmentDeliverRowless` shows that `Deliver` of the obligation-free binding of an outcome with sources returns an error that contains `names no obligation` (FD11).
- [ ] The existing tickets-only, rowless, and legacy landing tests stay green, and `bench test --package ./internal/worktree` passes.
