# 3. Bind only the content that a plan keeps open

Blocked by: 1-refuse-obligation-free-plan.md, 2-refuse-obligation-free-delivery.md
Writes: internal/commitment/authority.go, internal/commitment/authority_test.go, internal/commitment/removed_source_test.go, internal/commitment/repository/plan_after_delivery_test.go, internal/commitment/commitmenttest/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: FD12, FD13, FD14, FD15, FD16, FD17, FD18, FD19, FD25

## What to build

This ticket is the first ticket of review chunk FD-C2. That chunk starts after
the FD-C1 checkpoint, because both chunks write `BuildPlan`.

Each fixture in this ticket obeys the plan refusal of ticket 1. A binding in an
outcome with sources lists an obligation, or the current policy already holds
that binding unchanged.

`BuildPlan` binds the unsettled sources and the unsettled deliverables of the
proposed policy. From the current policy, it adds only each unsettled outcome
source that the proposed policy does not hold. It no longer adds a deliverable
of the current policy. The comparison uses the whole binding, as it does today.

The removed outcome's roadmap row stays bound, so the approval of a removal
still refuses when that row changes. A deliverable with a recorded delivery
fact stays unbound, as `Unsettled` defines it today. A kept unsettled
deliverable stays bound, so `Store.Approve` refuses when a later `main` commit
deletes it. A plan can approve a changed staged spec at its new identity.

`approvedTransition` recomputes `BuildPlan` at the commit, so an approved drop
or removal passes `Store.AuthorizeCandidate` without the deliverable on `main`.
No repository owner needs an edit. A shared fixture setup that two tests need
goes in `internal/commitment/commitmenttest/`.

## Acceptance

- [ ] `TestCommitmentPlanSourcesOmitDroppedDeliverable` in `internal/commitment/authority_test.go` shows that the `BuildPlan` sources omit a deliverable that the current policy approves and the proposal drops (FD13).
- [ ] `TestCommitmentPlanSurvivesRemovedDeliverable` in `internal/commitment/repository/plan_after_delivery_test.go` shows that `Store.Plan` plans a drop of a tickets-only binding whose folder a later `main` commit deleted (FD12).
- [ ] The same test shows that `Store.Plan` plans the removal of outcome `B` when a later `main` commit deleted the folder that `B` approves (FD14).
- [ ] The same test shows that `Store.Plan` plans a changed staged spec at its new `main` identity (FD18).
- [ ] The same test shows that `Store.AuthorizeCandidate` of the planning checkout's committed tree returns no error after the FD12 approval (FD19).
- [ ] The same test shows that `Store.AuthorizeCandidate` returns no error after the FD14 approval (FD25).
- [ ] `TestCommitmentApprovalBindsKeptDeliverable` in `internal/commitment/removed_source_test.go` shows that approval refuses a plan that keeps a tickets-only binding when a later `main` commit deletes that folder (FD17).
- [ ] `TestCommitmentRemovalBindsRemovedSource` stays green with no edit (FD15).
- [ ] `TestCommitmentPlanAfterDelivery` stays green with no edit (FD16).
