# A commitment deliverable closes or refuses, and a commitment plan is canonical

Status: staged

Roadmap: FT390

Decision source: `roadmap/FT390.md`, a named reviewed artifact from drain `d-0a44235dc225`.

Verification log: 2 iteration(s) to accept — Sonnet high ran the review. Iteration 1 found B1: the unscoped BuildPlan predicate refused a retained legacy binding. It also found S1, FD17 not red-capable, and S2, no identity seam. S3 was that FD21 could not tell byte order from numeric order. It also found N1 to N3. Iteration 2 accepted, and the fold added the FD10 inactive-milestone case, a corrected row list, and two story mappings.

## Problem

The commitment plan and the delivery record share one owner, `internal/commitment`. Three defects in that owner block or corrupt real delivery work.

First, a plan can approve a deliverable with no obligation for an outcome that owns sources. `validateDeliverables` in `parse.go` accepts an empty obligation list. At the completion landing, `Deliver` in `delivery.go` then returns the policy unchanged and no error. The landing closes the deliverable and records no delivery fact. The ADR26 landing `12e7c9aa` did this: its policy at `df45b6ab` bound `specs/adr-0026-retire-dc` with no obligation, and its published policy holds no delivery.

Second, a plan that removes a deliverable still needs that deliverable on `main`. `BuildPlan` in `authority.go` adds each unsettled source and deliverable of the current policy to the plan sources. `Store.Plan` then reads each one at `main` through `validateSources`. A deleted deliverable folder therefore refuses every plan, and the ADR26 repair took three landings to restore the folder first.

Third, the plan identity hashes the plan sources in traversal order. `BuildPlan` builds the source list from the policy order and appends the current-policy extras. A change to that traversal changes the plan identity while every bound fact stays the same.

## Solution

`bench commitment plan` refuses a proposal that adds or changes a deliverable with no obligation for an outcome that owns sources. The refusal names the outcome and the deliverable. A legacy policy that already holds such a binding stays readable. A plan that keeps that binding byte for byte still plans, so an unrelated plan and a repair plan both work. A completion landing of that legacy binding refuses before it publishes, and its refusal names `bench commitment plan --input <file>`.

A plan binds the content that its proposal keeps open, and each unsettled outcome source that its proposal drops. It does not bind a deliverable that the proposal no longer approves. A plan that drops a deleted deliverable therefore plans, and its approval commits.

The plan sorts its sources once, in one canonical order, and hashes that sorted list. The plan receipt lists the sources in that same order.

## User stories

Line: opus / high.
Implementation-line reason: FD-C1 is the hardest chunk, because its landing row drives a real `bench worktree land` fixture and must refuse before the gate. The spec pins each predicate. Each seam has a read precedent, and each row reds on the cheapest wrong result.
Harder chunks: FD-C1.

An obligation-free binding is a deliverable binding that lists no obligation, in an outcome that owns at least one source.

### A plan refuses a deliverable that can close nothing

1. As a reviewer, I want `bench commitment plan` to refuse an obligation-free binding, so that no approved build lands without a delivery fact.
2. As a reviewer, I want that refusal to name the outcome and the deliverable, so that I can correct the proposal without a search.
3. As a reviewer, I want an absent and an empty obligation list to refuse alike, so that the JSON spelling cannot bypass the rule.
4. As a reviewer, I want a rowless outcome to keep a deliverable with no obligation, so that the deliverable stays its own obligation.
5. As a reviewer, I want a refused plan to record no plan receipt, so that no approval can name it later.

### A completion landing refuses a delivery that records no fact

6. As a reviewer, I want `bench worktree land --spec` of an obligation-free binding to refuse before it publishes, so that no landing closes a deliverable without a fact.
7. As a reviewer, I want that landing refusal to name `bench commitment plan --input <file>`, so that I know the repair route.
8. As a reviewer, I want the closure owner to return the same refusal, so that the landing, the publication admission, and the gate oracle agree.
9. As a maintainer of a linked repository, I want a legacy obligation-free binding to stay readable, so that a plan can repair it.
10. As a worker, I want a landing of an unapproved deliverable to record no fact and pass, so that an unbound landing keeps working.

### A plan survives a settled or removed deliverable

11. As a reviewer, I want a plan that drops a deliverable to plan without it on `main`, so that a deleted folder blocks no plan.
12. As a reviewer, I want an outcome removal to plan without that outcome's deliverable on `main`, so that a removal binds only its row.
13. As a reviewer, I want a removed outcome's roadmap row to stay bound, so that a removal approval binds the exact obligation that it drops.
14. As a reviewer, I want a deliverable with a recorded delivery fact to stay unbound, so that a settled deliverable never blocks a later plan.
15. As a reviewer, I want a kept unsettled deliverable to stay bound, so that a changed or deleted one refuses the plan or its approval.
16. As a reviewer, I want a plan to re-approve a changed deliverable at its new identity, so that it needs no old bytes.
17. As a worker, I want an approved drop or removal to pass candidate authorization without the deliverable on `main`, so that the approval can land.

### The plan identity is canonical

18. As a reviewer, I want the plan identity to hash its sources in one canonical order, so that a traversal change cannot change that identity.
19. As a reviewer, I want the plan to list its sources in that same canonical order, so that the receipt and the identity agree.
20. As a reviewer, I want two sources with the same identifier to order by path, so that the order is total across outcomes.

### A plan keeps a legacy binding and checks every milestone

21. As a reviewer, I want a plan that keeps a legacy obligation-free binding unchanged to plan, so that a legacy policy blocks no plan.
22. As a reviewer, I want a plan to refuse a new obligation-free binding in an inactive milestone, so that no activation makes it landable.

### Reviewed exclusions

23. As a reviewer, I want `bench commitment start` to keep admitting a legacy obligation-free binding, so that this spec changes only the plan and the closure.
24. As a reviewer, I want a re-pinned roadmap row source to keep its plan behavior, so that this spec does not reopen the row-identity rule.

## Implementation decisions

### One obligation predicate, two enforcement points

One unexported predicate in `internal/commitment` decides whether a binding is obligation-free: the binding lists no obligation, and its outcome owns at least one source. Both enforcement points call that predicate. Rows FD1 to FD11, FD23, FD24, and FD26 reach it.

`BuildPlan` refuses each obligation-free binding of the proposed policy that is new or changed against the current policy. A binding is retained when the current policy holds the same outcome with an equal binding that is already obligation-free. A retained binding plans, and the check covers every milestone.

`BuildPlan` is the one plan derivation: `Store.Plan` calls it before it records a receipt, and `approvedTransition` calls it again at the commit. The refusal names the outcome identifier, the binding identifier, and the words `names no obligation`. Rows FD1, FD2, FD3, FD5, FD23, FD24, and FD26 reach it through `BuildPlan` and the plan command.

`Deliver` refuses an approved obligation-free binding with an error. The error names the deliverable path, the outcome, the words `names no obligation`, and `bench commitment plan --input <file>`. The current early return that hands back the policy unchanged goes away for that case. `Deliver` keeps the unchanged return for a path that the active milestone does not approve. `Store.Closure` returns the `Deliver` error, so `published.Tree`, `closedPolicy`, and `completionClosure` each refuse. Rows FD6 to FD10 reach it.

Contestable call, for reviewer veto: the source offers a plan refusal or a landing refusal. This spec takes both, through one predicate, because the kit ships to linked repositories. A plan refusal alone leaves a legacy binding to land silently. A landing refusal alone lets a plan approve a build that cannot land.

Contestable call, for reviewer veto: the plan refusal is scoped to new or changed bindings. An unscoped check refuses every later plan on a legacy policy, the repair plan included. `Deliver` stays the landing-time guard for a retained binding, so a retained binding still cannot land without a fact.

`Validate` does not take the rule. `Parse` calls `Validate` on every policy read, so a parse-time rule makes a legacy policy unreadable. The plan that repairs it reads that policy first. Row FD9 pins this.

### Plan sources

`BuildPlan` binds the unsettled sources and the unsettled deliverables of the proposed policy. From the current policy it adds only each unsettled outcome source that the proposed policy does not hold. The comparison uses the whole binding, as it does today. It no longer adds a current-policy deliverable. Rows FD12 to FD19 and FD25 reach this.

The removal binding of `TestCommitmentRemovalBindsRemovedSource` stays: a removed outcome's row is still bound. A deliverable is the work product, and the row is the obligation, so only the row binds a removal.

A settled deliverable is one with a recorded delivery fact, as `Unsettled` defines it today. This spec does not change that definition. Row FD16 pins it.

### Canonical order

One new unexported function in `authority.go`, `boundSources`, owns the canonical order. It sorts a copy of the plan sources once and returns the sorted list and the identity of its encoding. The order compares the identifier, then the path, then the identity, each as a Go string, which is byte order. `BuildPlan` stores the sorted list in `Plan.Sources` and puts the returned identity into the bound plan input. Today no separate identity function exists, because `BuildPlan` encodes the sources inline. Rows FD20, FD21, and FD22 reach `boundSources`.

Two causes move an existing plan identity: the canonical order, and the smaller source set of the plan-sources decision. A receipt approved before this change and not yet committed can therefore fail to match at the commit. The commit then refuses with its existing `candidate policy has no exact approval; run bench commitment plan --input <file>` message, and the reviewer plans again.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| FD-C1 / `1-refuse-obligation-free-plan.md`, `2-refuse-obligation-free-delivery.md` | A plan refuses an obligation-free binding, and a completion landing of a legacy one refuses before it publishes. | FD1, FD2, FD3, FD4, FD5, FD6, FD7, FD8, FD9, FD10, FD11, FD23, FD24, FD26 | `bench test --package ./internal/commitment`, `bench test --package ./internal/commitment/repository`, `bench test --package ./internal/worktree` | yes |
| FD-C2 / `3-bind-only-open-plan-sources.md`, `4-order-plan-sources-canonically.md` | A plan binds only what its proposal keeps open and the removed rows, in one canonical order. | FD12, FD13, FD14, FD15, FD16, FD17, FD18, FD19, FD20, FD21, FD22, FD25 | `bench test --package ./internal/commitment`, `bench test --package ./internal/commitment/repository` | no |

Both chunks write `internal/commitment/authority.go`, so FD-C2 starts after FD-C1 commits green. The four tickets run in series on one integration source.

| Ticket | Blocked by | Delivered coverage |
| --- | --- | --- |
| [1. Refuse an obligation-free binding at plan time](tickets/1-refuse-obligation-free-plan.md) | none | FD1, FD2, FD3, FD4, FD5, FD9, FD23, FD24, FD26 |
| [2. Refuse an obligation-free delivery at the completion landing](tickets/2-refuse-obligation-free-delivery.md) | 1-refuse-obligation-free-plan.md | FD6, FD7, FD8, FD10, FD11 |
| [3. Bind only the content that a plan keeps open](tickets/3-bind-only-open-plan-sources.md) | 1-refuse-obligation-free-plan.md, 2-refuse-obligation-free-delivery.md | FD12, FD13, FD14, FD15, FD16, FD17, FD18, FD19, FD25 |
| [4. Order plan sources in one canonical order](tickets/4-order-plan-sources-canonically.md) | 3-bind-only-open-plan-sources.md | FD20, FD21, FD22 |

Ticket 1 adds the obligation predicate, and ticket 2 calls it. Ticket 3 waits for the FD-C1 checkpoint, and its fixtures obey the plan refusal of ticket 1. Ticket 4 orders the source set that ticket 3 defines.

## Testing decisions

- A good test drives the real plan owner or the real closure owner with exact policy bytes and a real fixture repository. It observes the plan, the refusal text, the receipt ledger, or the published `main` ref.
- The pure rules attach at `BuildPlan` and `Deliver`, after the precedent of `TestCommitmentRemovalEffects` and `TestCommitmentDeliver`.
- The repository rules attach at `Store.Plan`, `Store.Closure`, and `Store.AuthorizeCandidate`, after the precedent of `TestCommitmentPlanAfterDelivery`, `TestCommitmentLightClosure`, and `TestCommitmentRemovalBindsRemovedSource`.
- The command rows attach at `commitcmd.Command`, after the precedent of `TestCommitmentLiteralInput`.
- The landing rows attach at `bench worktree land` through the `deliveryRoute` fixture, after the precedent of `TestCommitmentTicketsOnlyClosure`.
- The canonical-identity row is an internal test of `boundSources`, because no public input changes the traversal order and keeps every bound fact.
- The package tests run in the gate's `test` phase.

### Seam diagram

    trigger: bench commitment plan --input <file>
        │
        ▼
    proposal bytes  ──▶  [ Store.Plan → BuildPlan: obligation predicate,
                           plan sources, canonical order, identity ]  ──▶  plan receipt or refusal
                      ◀ tests attach here: BuildPlan with exact policies; Store.Plan
                        and commitcmd.Command over a fixture repository

    trigger: bench worktree land --spec <deliverable>
        │
        ▼
    reviewed source  ──▶  [ published.Tree → Store.Closure → Deliver ]  ──▶  published tree or refusal
                      ◀ tests attach here: Deliver with exact policies; Store.Closure
                        over a fixture repository; the deliveryRoute landing fixture

    trigger: bench commit of an approved policy
        │
        ▼
    candidate tree  ──▶  [ AuthorizeCandidate → approvedTransition → BuildPlan ]  ──▶  admitted or refused
                      ◀ tests attach here: Store.AuthorizeCandidate in the planning checkout

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| FD1 | 1 | `BuildPlan` returns an error that contains `names no obligation` for a proposal whose outcome owns `FT1` and approves a deliverable with no `obligations` key, against a current policy without that binding | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | The current `BuildPlan` returns a plan for this proposal. |
| FD2 | 3 | `bench commitment plan --input` exits 1 for a proposal whose binding holds `"obligations": []` and whose staged spec exists at `main` | planned TestCommitmentPlanCommandRefusesEmptyObligations in internal/commitment/command_test.go | A check that tests only for an absent key lets the empty list plan. |
| FD3 | 2 | The FD1 refusal text contains the outcome identifier and the binding identifier | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | A generic refusal leaves the reviewer to search the policy for the binding. |
| FD4 | 4 | `BuildPlan` returns a plan for a proposal whose rowless outcome approves a deliverable with no obligation | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | A predicate that ignores the outcome sources refuses every rowless plan. |
| FD5 | 5 | After the FD2 refusal, the intent ledger holds no commitment receipt | planned TestCommitmentPlanCommandRefusesEmptyObligations in internal/commitment/command_test.go | A refusal placed after the receipt write leaves an approvable receipt. |
| FD6 | 6 | `bench worktree land --spec` of the obligation-free spec that `SeedTicketsOnly` approves exits nonzero and leaves the `main` ref at its base | planned TestCommitmentLandingRefusesObligationFreeDelivery in internal/worktree/commitment_light_landing_test.go | The current landing publishes the spec flip with no delivery fact. |
| FD7 | 7 | The FD6 landing output contains `bench commitment plan --input <file>` | planned TestCommitmentLandingRefusesObligationFreeDelivery in internal/worktree/commitment_light_landing_test.go | A refusal without the repair command leaves the reviewer with no route. |
| FD8 | 8 | `Store.Closure` returns an error that contains `names no obligation` for the obligation-free spec that `SeedTicketsOnly` approves | planned TestCommitmentClosureRefusesObligationFreeBinding in internal/commitment/repository/closure_test.go | The current closure returns no edit and no error, so each closure caller proceeds. |
| FD9 | 9 | `commitment.Parse` accepts a policy whose outcome owns `FT1` and approves a deliverable with no obligation | planned TestCommitmentParseKeepsObligationFreeBinding in internal/commitment/parse_test.go | A rule placed in `Validate` makes the legacy policy unreadable, and no plan can repair it. |
| FD10 | 10 | `Deliver` returns no error and no new fact for `specs/other/spec.md` and for an obligation-free binding that only an inactive milestone approves | `internal/commitment/delivery_test.go` (`TestCommitmentDeliver`) | A refusal widened to every no-fact delivery, or to every milestone, breaks an unbound landing. |
| FD11 | 6 | `Deliver` of the obligation-free binding of an outcome with sources returns an error that contains `names no obligation` | `internal/commitment/delivery_test.go` (`TestCommitmentDeliverRowless`) | The current `Deliver` returns the unchanged policy and no error. |
| FD12 | 11 | `Store.Plan` returns a plan for a proposal that drops a tickets-only binding whose folder a later `main` commit deleted | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | The current plan binds the current-policy deliverable and refuses at `main`. |
| FD13 | 11 | The `BuildPlan` sources omit a deliverable that the current policy approves and the proposal drops | planned TestCommitmentPlanSourcesOmitDroppedDeliverable in internal/commitment/authority_test.go | A union with the current-policy deliverables keeps the dropped binding. |
| FD14 | 12 | `Store.Plan` returns a plan for a proposal that removes outcome `B` when a later `main` commit deleted the folder that `B` approves | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | The current plan binds the removed outcome's deliverable and refuses at `main`. |
| FD15 | 13 | Approval of a removal refuses when the removed outcome's row changed after the plan | `internal/commitment/removed_source_test.go` (`TestCommitmentRemovalBindsRemovedSource`) | A source set that drops every current-policy source drops the removed row too. |
| FD16 | 14 | `Store.Plan` returns a plan after a delivery flips the delivered spec | `internal/commitment/repository/plan_after_delivery_test.go` (`TestCommitmentPlanAfterDelivery`) | A source set that binds settled deliverables refuses at the flipped spec. |
| FD17 | 15 | Approval refuses a plan that keeps a tickets-only binding when a `main` commit after the plan deletes that folder | planned TestCommitmentApprovalBindsKeptDeliverable in internal/commitment/removed_source_test.go | `Store.Approve` reads only `Plan.Sources`, so a fix that drops every deliverable from the plan sources lets the approval stage. |
| FD18 | 16 | `Store.Plan` returns a plan for a proposal that approves a changed staged spec at its new `main` identity | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | The current plan binds the old identity and refuses at `main`. |
| FD19 | 17 | After the FD12 plan is approved, `Store.AuthorizeCandidate` of the planning checkout's committed tree returns no error | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | `approvedTransition` recomputes `BuildPlan`, so a source set that differs at the commit refuses there. |
| FD20 | 18 | `boundSources` returns one sorted list and one identity for the inputs `FT9`, `FT1` and `FT1`, `FT9` | planned TestBoundSourcesIgnoresInputOrder in internal/commitment/authority_internal_test.go | A hash over the caller's order returns two identities. |
| FD21 | 19 | `BuildPlan` lists `FT10` before `FT9` when outcome `A` owns `FT9` and the later outcome `B` owns `FT10` | planned TestCommitmentPlanSourcesAreCanonical in internal/commitment/authority_test.go | Traversal order and numeric order each put `FT9` first, and only byte order puts `FT10` first. |
| FD22 | 20 | `BuildPlan` lists the deliverable `spec` at `specs/a/spec.md` before the deliverable `spec` at `specs/b/spec.md` when outcome `A` approves the `b` path | planned TestCommitmentPlanSourcesAreCanonical in internal/commitment/authority_test.go | A sort on the identifier alone leaves equal identifiers in traversal order. |
| FD23 | 21 | `BuildPlan` returns a plan for a proposal that keeps an obligation-free binding of the current policy byte for byte and adds outcome `C` | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | An unscoped check refuses every plan on a legacy policy. |
| FD24 | 1 | `BuildPlan` returns an error that contains `names no obligation` for a proposal that keeps the legacy binding at a changed identity | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | A retention test by binding identifier alone lets a changed binding plan. |
| FD25 | 17 | After the FD14 plan is approved, `Store.AuthorizeCandidate` of the planning checkout's committed tree returns no error | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | `approvedTransition` recomputes `BuildPlan`, so a removal that binds the deleted folder refuses at the commit. |
| FD26 | 22 | `BuildPlan` returns an error that contains `names no obligation` for a new obligation-free binding in an inactive milestone | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | A check copied from the active-milestone scope of `Deliver` lets the binding plan. |

Not covered: story 23 — a reviewed exclusion; the edge inventory states why it is safe.
Not covered: story 24 — a reviewed exclusion; Out of scope prices it.

### Edge inventory

The hostile-input classes of `projects/benchkit.md` reach no new parse surface here. The sort and the predicate read only identifiers that `validateSource` already restricts to `^[A-Za-z][A-Za-z0-9._-]*$`, and paths that `repositoryPath` already cleans. The absent-versus-empty class is FD1 and FD2.

**Won't handle** — `bench commitment start` of a legacy obligation-free binding — only a legacy policy holds one, and FD6 refuses its landing with the repair command.

**Won't handle** — numeric order of row identifiers, so `FT10` sorts before `FT9` — byte order is one total order, and FD21 observes it.

The two checks have different milestone scopes, by decision. `Deliver` checks only the active milestone, because `activeDeliverable` resolves a path only there. A landing of an inactive-milestone deliverable records no fact and refuses nothing, as FD10 shows. `BuildPlan` checks every milestone, because a later activation makes an inactive binding landable. FD26 pins that scope.

**Won't handle** — an uncommitted receipt approved before this change — its commit refuses with the existing plan guidance, and FD19 shows that a fresh plan commits.

**Won't handle** — an undelivered binding of a closed outcome stays bound — a plan that drops that binding repairs it, and FD12 covers that plan.

**Won't handle** — a plan that drops the deliverable of a bound running assignment — the existing run checks stay unchanged, and FD12 covers the drop.

## Ownership fences

- `internal/commitment/authority.go`
- `internal/commitment/authority_test.go`
- `internal/commitment/authority_internal_test.go`
- `internal/commitment/delivery.go`
- `internal/commitment/delivery_test.go`
- `internal/commitment/parse_test.go`
- `internal/commitment/removed_source_test.go`
- `internal/commitment/command_test.go`
- `internal/commitment/commitmenttest/`
- `internal/commitment/repository/closure_test.go`
- `internal/commitment/repository/plan_after_delivery_test.go`
- `internal/worktree/commitment_light_landing_test.go`
- `internal/worktree/commitment_landing_fixture_test.go`
- `internal/worktree/parallel_census_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `reviews/commitment-delivery-integrity.md`

Build preflight binds the commitment and worktree packages to the five command-registry files above, so each ticket names them. No ticket expects to edit them.

## Out of scope

- A re-pinned roadmap row source. Today a plan that moves a kept row to its new identity binds the old current-policy identity too, so the plan refuses at `main`. The fix needs its own decision about the order of a row edit and its re-pin. Estimate: 4 edits, 2 gate runs.
- A settled state for an undelivered binding whose outcome other deliveries closed. This changes `Unsettled` for every reader. Estimate: 6 edits, 2 gate runs.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| "A plan that binds a deliverable to an outcome with sources lists the obligations that the deliverable closes." | FD1, FD2, FD3, FD4, FD5, FD23, FD24, FD26 |
| "Today the plan accepts an empty obligation list." | FD1, FD2 |
| "The completion landing then records no delivery fact, because `Deliver` returns the policy unchanged." | FD6, FD8, FD11 |
| "Either the plan refuses that binding, or the completion landing refuses a completion that records no fact." | FD1 to FD11; the contestable call takes both sides |
| "A plan that settles or removes a deliverable does not need that deliverable on `main`." | FD12, FD13, FD14, FD16, FD19, FD25 |
| "Today `BuildPlan` adds each unsettled source of the current policy, and the source check reads each one at `main`." | FD13, FD15, FD17 |
| "A deleted deliverable folder then blocks every plan." | FD12, FD14 |
| "The plan identity sorts its policy sources in one canonical order." | FD20, FD21, FD22 |
| "Today it hashes them in traversal order, so a change to that order can change the identity while every bound fact stays the same." | FD20 |

Occurrences, each re-read in this session:

- 2026-10-04 roadmap-delivery-commitment build: the folded FT382 detail in commit `695fa4b5` states it.
- 2026-10-05 ADR26 landing `12e7c9aa`: its policy holds no delivery, and the binding at `df45b6ab` lists no obligation.
- 2026-10-05 ADR26 repair: commits `f99e825c`, `a7a70e5d`, and `62129649` restore the folder, bind `FT389`, and record the delivery.

### Reader sweep

- Readers of `Plan.Sources`: `Store.Plan` and `Store.Approve` through `validateSources`, `approvedTransition` through `validateSourceAt`, and the receipt payload in the intent ledger. Each reads the canonical order with no order assumption.
- Readers of the plan identity: the receipt lookup in `Store.Approve`, the receipt match in `approvedTransition`, and the `commitment_plan` table of `commitcmd`.
- Callers of `BuildPlan`: `Store.Plan`, `approvedTransition`, and the four tests in `authority_test.go`. `approvedTransition` passes the policy at `main` as the current policy, so a retained legacy binding stays retained at the commit.
- Callers of `Deliver`: `Store.Delivered` alone in production. The test callers are `delivery_test.go`, `retirePolicy` in `internal/spec/spec_test.go`, and `TestCommitmentContinuationOccupiesSlot` in `internal/commitment/continuation_test.go`. Each of the two extra callers delivers a binding that lists an obligation or a rowless binding, so neither reaches the new refusal.
- Callers of `Store.Closure`: `published.closeDelivery`, `closedPolicy` in `candidate.go`, and `completionClosure` in `internal/gate/completion.go`. Each returns the closure error unchanged, so none needs an edit.
- Input constructors of the plan predicate: each test builder that builds a binding with no obligation. `SeedTicketsOnly` is the one builder whose outcome owns sources; it writes the policy directly, and FD6 and FD8 use it as the legacy fixture. The rowless builders in `admission_test.go`, `outlook_test.go`, `deliverable_test.go`, `commitmenttest/repo.go`, and `delivery_test.go` stay valid under FD4.
- Shipped-surface claim words: ADR 0015 already says that a deliverable is "the complete delivery of named obligations", so no doc changes.

### Flagged additions

- FD18 lets a plan re-approve a changed deliverable. The source names only a settled or removed deliverable. The same `BuildPlan` line causes this refusal, so the fix removes it with no extra code.
- FD7 adds the repair command to the landing refusal. The source asks for a refusal only, and the operating guide requires a refusal to name the next action.

### Pre-review proof checklist

- `Cited symbols`: each symbol below resolves at `287777e5`:
  - `BuildPlan`, `policySources`, `Deliver`, `Unsettled`, `Validate`, and `validateDeliverables`
  - `Store.Plan`, `Store.Approve`, `approvedTransition`, `closedPolicy`, `Store.Closure`, and `Store.Delivered`
  - `published.Tree`, `closeDelivery`, `completionClosure`, `SeedTicketsOnly`, and each existing test in the map
- `Import edges`: none. Each change stays inside its package.
- `Source-row clauses and occurrences`: the source trace quotes each clause of `roadmap/FT390.md` and lists each occurrence.
- `Promised field labels`: the refusal words `names no obligation` and the command `bench commitment plan --input <file>`.
- `Changed-function callers`: the reader sweep lists each caller of `BuildPlan`, `Deliver`, and `Store.Closure`.
- `Copy survival`: none. No copy moves to a new owner.
- `Rendered-shape readers`: the final assertion of `TestCommitmentDeliverRowless` expects the silent unchanged return today, and FD11 changes it to the refusal. No other test matches the changed text.

### Source disclosure

The source names `capture/learnings.md` from drains `d-553013ef861d` and `d-0a44235dc225`. That file is git-ignored, and the drains closed those entries, so this session could not re-read their text. This spec re-reads the roadmap diffs of commits `695fa4b5` and `86ea7696`, the FT382 detail that `695fa4b5` folded, and the ADR26 commits named above.

### Completion plan

```bench-completion-plan
{"version":2,"chunks":[{"id":"FD-C1","tickets":["1-refuse-obligation-free-plan.md","2-refuse-obligation-free-delivery.md"],"verification":[{"id":"t1-commitment","command":"bench test --package ./internal/commitment","probe":"In BuildPlan, skip the refusal of a new obligation-free binding. TestCommitmentPlanRefusesObligationFreeBinding must fail and the restore must be exact.","ticket":"1-refuse-obligation-free-plan.md"},{"id":"t1-conformance","command":"bench test --package ./internal/conformance","ticket":"1-refuse-obligation-free-plan.md"},{"id":"t1-bench","command":"bench test --package ./cmd/bench","ticket":"1-refuse-obligation-free-plan.md"},{"id":"t2-commitment","command":"bench test --package ./internal/commitment","probe":"In Deliver, return the policy unchanged with no error for an obligation-free binding of an outcome with sources. TestCommitmentDeliverRowless must fail and the restore must be exact.","ticket":"2-refuse-obligation-free-delivery.md"},{"id":"t2-commitment-repository","command":"bench test --package ./internal/commitment/repository","ticket":"2-refuse-obligation-free-delivery.md"},{"id":"t2-worktree","command":"bench test --package ./internal/worktree","ticket":"2-refuse-obligation-free-delivery.md"},{"id":"t2-conformance","command":"bench test --package ./internal/conformance","ticket":"2-refuse-obligation-free-delivery.md"},{"id":"t2-bench","command":"bench test --package ./cmd/bench","ticket":"2-refuse-obligation-free-delivery.md"}]},{"id":"FD-C2","tickets":["3-bind-only-open-plan-sources.md","4-order-plan-sources-canonically.md"],"verification":[{"id":"t3-commitment","command":"bench test --package ./internal/commitment","probe":"In BuildPlan, add each unsettled deliverable of the current policy to the plan sources again. TestCommitmentPlanSourcesOmitDroppedDeliverable must fail and the restore must be exact.","ticket":"3-bind-only-open-plan-sources.md"},{"id":"t3-commitment-repository","command":"bench test --package ./internal/commitment/repository","ticket":"3-bind-only-open-plan-sources.md"},{"id":"t3-conformance","command":"bench test --package ./internal/conformance","ticket":"3-bind-only-open-plan-sources.md"},{"id":"t4-commitment","command":"bench test --package ./internal/commitment","probe":"In boundSources, return the copy without the sort. TestBoundSourcesIgnoresInputOrder must fail and the restore must be exact.","ticket":"4-order-plan-sources-canonically.md"},{"id":"t4-commitment-repository","command":"bench test --package ./internal/commitment/repository","ticket":"4-order-plan-sources-canonically.md"},{"id":"t4-conformance","command":"bench test --package ./internal/conformance","ticket":"4-order-plan-sources-canonically.md"}]}],"final_verification":[{"id":"coverage-check","command":"bench coverage --check specs/commitment-delivery-integrity/spec.md"},{"id":"commitment","command":"bench test --package ./internal/commitment"},{"id":"commitment-repository","command":"bench test --package ./internal/commitment/repository"},{"id":"worktree","command":"bench test --package ./internal/worktree"}],"execution":{"mode":"delegate","run_id":"fd-full-20261005","orchestrator_session":"claude:session_0128bLurgWw2evJv2LYjfDHH","author_limit":1,"assignments":{"1-refuse-obligation-free-plan.md":[{"session":"claude:fd_t1","assignment":"e2818e688bf9c87c5d08bf5241bb46af","model":"opus","effort":"high","source":"e57d65e9315ad7426443828377d06fa9c6ecfd09","native_ref":"claude-agent:fd_t1"},{"session":"claude:fd_t1_r1","assignment":"e2818e688bf9c87c5d08bf5241bb46af","model":"opus","effort":"high","source":"ed40f39a56494b0fa82e965a82d1602bfb2f2fba","native_ref":"claude-agent:fd_t1_r1","predecessor":"claude:fd_t1","trigger":"user-directed","stopped":"claude-agent:fd_t1 returned after record commit 65036774","preserved":"ed40f39a56494b0fa82e965a82d1602bfb2f2fba"}],"2-refuse-obligation-free-delivery.md":[{"session":"claude:fd_t2","assignment":"e2818e688bf9c87c5d08bf5241bb46af","model":"opus","effort":"high","source":"5b5cde7768f3f6a700c6933648884f099007c879","native_ref":"claude-agent:fd_t2"}],"3-bind-only-open-plan-sources.md":[],"4-order-plan-sources-canonically.md":[]}}}
```
