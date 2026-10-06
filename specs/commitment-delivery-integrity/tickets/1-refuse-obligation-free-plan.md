# 1. Refuse an obligation-free binding at plan time

Blocked by: none
Writes: internal/commitment/authority.go, internal/commitment/delivery.go, internal/commitment/authority_test.go, internal/commitment/command_test.go, internal/commitment/parse_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: FD1, FD2, FD3, FD4, FD5, FD9, FD23, FD24, FD26

## What to build

This ticket is the first ticket of review chunk FD-C1. It delivers the plan
refusal of an obligation-free binding, through `bench commitment plan` and
through the `BuildPlan` recheck at the commit.

First, verify the premise of the decision source against the code.
`validateDeliverables` in `parse.go` accepts an empty obligation list, and the
current `BuildPlan` returns a plan for an obligation-free binding.

Add one unexported predicate to `delivery.go`. The predicate is true when a
binding lists no obligation and its outcome owns at least one source. The
existing inline condition in `Deliver` calls this predicate, so the package
holds one copy of the rule. `Deliver` keeps its current behavior in this
ticket; ticket 2 changes it.

`BuildPlan` refuses each obligation-free binding of the proposed policy that is
new or changed against the current policy. The check reads every milestone,
not only the active milestone. A binding is retained under two conditions.
First, the current policy holds the same outcome with an equal binding.
Second, that binding is already obligation-free in the current policy. The
comparison uses the whole binding.

A retained binding plans. The refusal text names the outcome identifier, the
binding identifier, and the words `names no obligation`.

`Store.Plan` calls `BuildPlan` before it writes a receipt, so a refused plan
records no receipt. `approvedTransition` calls `BuildPlan` with the policy at
`main` as the current policy, so a retained binding stays retained at the
commit. Neither repository owner needs an edit.

Do not put the rule in `Validate`. `Parse` calls `Validate` on each policy
read, and the plan that repairs a legacy policy reads that policy first.

Each test policy that binds a deliverable in an outcome with sources lists an
obligation, unless the test pins the refusal or the retained case.

## Acceptance

- [ ] `TestCommitmentPlanRefusesObligationFreeBinding` in `internal/commitment/authority_test.go` shows that `BuildPlan` refuses a new binding with no `obligations` key in an outcome that owns `FT1` (FD1).
- [ ] The same test shows that the refusal text names `names no obligation`, the outcome, and the binding (FD3).
- [ ] The same test shows that `BuildPlan` plans a rowless outcome whose deliverable lists no obligation (FD4).
- [ ] The same test shows that `BuildPlan` plans a proposal that keeps a legacy obligation-free binding byte for byte and adds outcome `C` (FD23).
- [ ] The same test shows that `BuildPlan` refuses the legacy binding at a changed identity (FD24).
- [ ] The same test shows that `BuildPlan` refuses a new obligation-free binding in an inactive milestone (FD26).
- [ ] `TestCommitmentPlanCommandRefusesEmptyObligations` in `internal/commitment/command_test.go` shows that `bench commitment plan --input` exits 1 for `"obligations": []` when the staged spec exists at `main` (FD2).
- [ ] The same test shows that the intent ledger holds no commitment receipt after that refusal (FD5).
- [ ] `TestCommitmentParseKeepsObligationFreeBinding` in `internal/commitment/parse_test.go` shows that `commitment.Parse` accepts a policy with an obligation-free binding (FD9).
- [ ] `Deliver` has the same results as before this ticket, and `bench test --package ./internal/commitment` passes.
