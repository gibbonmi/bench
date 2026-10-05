# A commitment deliverable closes or refuses, and a commitment plan is canonical

Status: staged

Roadmap: FT390

Decision source: `roadmap/FT390.md`, a named reviewed artifact from drain `d-0a44235dc225`.

## Problem

The commitment plan and the delivery record share one owner, `internal/commitment`. Three defects in that owner block or corrupt real delivery work.

First, a plan can approve a deliverable with no obligation for an outcome that owns sources. `validateDeliverables` in `parse.go` accepts an empty obligation list. At the completion landing, `Deliver` in `delivery.go` then returns the policy unchanged and no error. The landing closes the deliverable and records no delivery fact. The ADR26 landing `12e7c9aa` did this: its policy at `df45b6ab` bound `specs/adr-0026-retire-dc` with no obligation, and its published policy holds no delivery.

Second, a plan that removes a deliverable still needs that deliverable on `main`. `BuildPlan` in `authority.go` adds each unsettled source and deliverable of the current policy to the plan sources. `Store.Plan` then reads each one at `main` through `validateSources`. A deleted deliverable folder therefore refuses every plan, and the ADR26 repair took three landings to restore the folder first.

Third, the plan identity hashes the plan sources in traversal order. `BuildPlan` builds the source list from the policy order and appends the current-policy extras. A change to that traversal changes the plan identity while every bound fact stays the same.

## Solution

`bench commitment plan` refuses a proposal that approves a deliverable with no obligation for an outcome that owns sources. The refusal names the outcome and the deliverable. A legacy policy that already holds such a binding stays readable, so a plan can repair it. A completion landing of that legacy binding refuses before it publishes, and its refusal names `bench commitment plan --input <file>`.

A plan binds the content that its proposal keeps open, and each unsettled outcome source that its proposal drops. It does not bind a deliverable that the proposal no longer approves. A plan that drops a deleted deliverable therefore plans, and its approval commits.

The plan sorts its sources in one canonical order before it hashes them. The plan receipt lists the sources in that same order.

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
15. As a reviewer, I want a kept unsettled deliverable to stay bound at `main`, so that a changed or deleted one still refuses the plan.
16. As a reviewer, I want a plan to re-approve a changed deliverable at its new identity, so that it needs no old bytes.
17. As a worker, I want the commit of an approved removal to pass candidate authorization without the removed deliverable, so that the approval can land.

### The plan identity is canonical

18. As a reviewer, I want the plan identity to hash its sources in one canonical order, so that a traversal change cannot change that identity.
19. As a reviewer, I want the plan to list its sources in that same canonical order, so that the receipt and the identity agree.
20. As a reviewer, I want two sources with the same identifier to order by path, so that the order is total across outcomes.

### Reviewed exclusions

21. As a reviewer, I want `bench commitment start` to keep admitting a legacy obligation-free binding, so that this spec changes only the plan and the closure.
22. As a reviewer, I want a re-pinned roadmap row source to keep its plan behavior, so that this spec does not reopen the row-identity rule.

## Implementation decisions

### One obligation predicate, two enforcement points

One unexported predicate in `internal/commitment` decides whether a binding is obligation-free: the binding lists no obligation, and its outcome owns at least one source. Both enforcement points call that predicate. Rows FD1 to FD11 reach it.

`BuildPlan` refuses a proposed policy that holds an obligation-free binding. `BuildPlan` is the one plan derivation: `Store.Plan` calls it before it records a receipt, and `approvedTransition` calls it again at the commit. The refusal names the outcome identifier, the binding identifier, and the words `names no obligation`. Rows FD1, FD2, FD3, and FD5 reach it through `BuildPlan` and the plan command.

`Deliver` refuses an approved obligation-free binding with an error. The error names the deliverable path, the outcome, the words `names no obligation`, and `bench commitment plan --input <file>`. The current early return that hands back the policy unchanged goes away for that case. `Deliver` keeps the unchanged return for a path that the active milestone does not approve. `Store.Closure` returns the `Deliver` error, so `published.Tree`, `closedPolicy`, and `completionClosure` each refuse. Rows FD6 to FD10 reach it.

Contestable call, for reviewer veto: the source offers a plan refusal or a landing refusal. This spec takes both, through one predicate, because the kit ships to linked repositories. A plan refusal alone leaves a legacy binding to land silently. A landing refusal alone lets a plan approve a build that cannot land.

`Validate` does not take the rule. `Parse` calls `Validate` on every policy read, so a parse-time rule makes a legacy policy unreadable. The plan that repairs it reads that policy first. Row FD9 pins this.

### Plan sources

`BuildPlan` binds the unsettled sources and the unsettled deliverables of the proposed policy. From the current policy it adds only each unsettled outcome source that the proposed policy does not hold. The comparison uses the whole binding, as it does today. It no longer adds a current-policy deliverable. Rows FD12 to FD19 reach this.

The removal binding of `TestCommitmentRemovalBindsRemovedSource` stays: a removed outcome's row is still bound. A deliverable is the work product, and the row is the obligation, so only the row binds a removal.

A settled deliverable is one with a recorded delivery fact, as `Unsettled` defines it today. This spec does not change that definition. Row FD16 pins it.

### Canonical order

One unexported function in `authority.go` returns the plan sources in canonical order. The order compares the identifier, then the path, then the identity, each as a Go string, which is byte order. `BuildPlan` stores that order in `Plan.Sources`. The plan identity derivation sorts its own input through the same function before it hashes, so the identity cannot depend on the caller's order. Rows FD20, FD21, and FD22 reach it.

The change moves each plan identity whose sources were out of canonical order. A receipt approved before this change and not yet committed no longer matches at the commit. The commit then refuses with its existing `candidate policy has no exact approval; run bench commitment plan --input <file>` message, and the reviewer plans again.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| FD-C1 / tickets pending the slicing delegate | A plan refuses an obligation-free binding, and a completion landing of a legacy one refuses before it publishes. | FD1, FD2, FD3, FD4, FD5, FD6, FD7, FD8, FD9, FD10, FD11 | `bench test --package ./internal/commitment`, `bench test --package ./internal/commitment/repository`, `bench test --package ./internal/worktree` | yes |
| FD-C2 / tickets pending the slicing delegate | A plan binds only what its proposal keeps open and the removed rows, in one canonical order. | FD12, FD13, FD14, FD15, FD16, FD17, FD18, FD19, FD20, FD21, FD22 | `bench test --package ./internal/commitment`, `bench test --package ./internal/commitment/repository` | no |

Both chunks write `internal/commitment/authority.go`, so FD-C2 starts after FD-C1 commits green. The slicing delegate writes the ticket basenames, the ticket graph, and the completion plan.

## Testing decisions

- A good test drives the real plan owner or the real closure owner with exact policy bytes and a real fixture repository. It observes the plan, the refusal text, the receipt ledger, or the published `main` ref.
- The pure rules attach at `BuildPlan` and `Deliver`, after the precedent of `TestCommitmentRemovalEffects` and `TestCommitmentDeliver`.
- The repository rules attach at `Store.Plan`, `Store.Closure`, and `Store.AuthorizeCandidate`, after the precedent of `TestCommitmentPlanAfterDelivery`, `TestCommitmentLightClosure`, and `TestCommitmentRemovalBindsRemovedSource`.
- The command rows attach at `commitcmd.Command`, after the precedent of `TestCommitmentLiteralInput`.
- The landing rows attach at `bench worktree land` through the `deliveryRoute` fixture, after the precedent of `TestCommitmentTicketsOnlyClosure`.
- The canonical-identity row is an internal test in package `commitment`, because no public input can change the traversal order while it keeps every bound fact.
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
| FD1 | 1 | `BuildPlan` returns an error that contains `names no obligation` for a proposal whose outcome owns `FT1` and approves a deliverable with no `obligations` key | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | The current `BuildPlan` returns a plan for this proposal. |
| FD2 | 3 | `bench commitment plan --input` exits 1 for a proposal whose binding holds `"obligations": []` and whose staged spec exists at `main` | planned TestCommitmentPlanCommandRefusesEmptyObligations in internal/commitment/command_test.go | A check that tests only for an absent key lets the empty list plan. |
| FD3 | 2 | The FD1 refusal text contains the outcome identifier and the binding identifier | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | A generic refusal leaves the reviewer to search the policy for the binding. |
| FD4 | 4 | `BuildPlan` returns a plan for a proposal whose rowless outcome approves a deliverable with no obligation | planned TestCommitmentPlanRefusesObligationFreeBinding in internal/commitment/authority_test.go | A predicate that ignores the outcome sources refuses every rowless plan. |
| FD5 | 5 | After the FD2 refusal, the intent ledger holds no commitment receipt | planned TestCommitmentPlanCommandRefusesEmptyObligations in internal/commitment/command_test.go | A refusal placed after the receipt write leaves an approvable receipt. |
| FD6 | 6 | `bench worktree land --spec` of the obligation-free spec that `SeedTicketsOnly` approves exits nonzero and leaves the `main` ref at its base | planned TestCommitmentLandingRefusesObligationFreeDelivery in internal/worktree/commitment_light_landing_test.go | The current landing publishes the spec flip with no delivery fact. |
| FD7 | 7 | The FD6 landing output contains `bench commitment plan --input <file>` | planned TestCommitmentLandingRefusesObligationFreeDelivery in internal/worktree/commitment_light_landing_test.go | A refusal without the repair command leaves the reviewer with no route. |
| FD8 | 8 | `Store.Closure` returns an error that contains `names no obligation` for the obligation-free spec that `SeedTicketsOnly` approves | planned TestCommitmentClosureRefusesObligationFreeBinding in internal/commitment/repository/closure_test.go | The current closure returns no edit and no error, so each closure caller proceeds. |
| FD9 | 9 | `commitment.Parse` accepts a policy whose outcome owns `FT1` and approves a deliverable with no obligation | planned TestCommitmentParseKeepsObligationFreeBinding in internal/commitment/parse_test.go | A rule placed in `Validate` makes the legacy policy unreadable, and no plan can repair it. |
| FD10 | 10 | `Deliver` of a path that the active milestone does not approve returns no error and no new fact | `internal/commitment/delivery_test.go` (`TestCommitmentDeliver`) | A refusal widened to every no-fact delivery breaks an unbound landing. |
| FD11 | 6 | `Deliver` of the obligation-free binding of an outcome with sources returns an error that contains `names no obligation` | `internal/commitment/delivery_test.go` (`TestCommitmentDeliverRowless`) | The current `Deliver` returns the unchanged policy and no error. |
| FD12 | 11 | `Store.Plan` returns a plan for a proposal that drops a tickets-only binding whose folder a later `main` commit deleted | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | The current plan binds the current-policy deliverable and refuses at `main`. |
| FD13 | 11 | The `BuildPlan` sources omit a deliverable that the current policy approves and the proposal drops | planned TestCommitmentPlanSourcesOmitDroppedDeliverable in internal/commitment/authority_test.go | A union with the current-policy deliverables keeps the dropped binding. |
| FD14 | 12 | `Store.Plan` returns a plan for a proposal that removes outcome `B` when a later `main` commit deleted the folder that `B` approves | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | The current plan binds the removed outcome's deliverable and refuses at `main`. |
| FD15 | 13 | Approval of a removal refuses when the removed outcome's row changed after the plan | `internal/commitment/removed_source_test.go` (`TestCommitmentRemovalBindsRemovedSource`) | A source set that drops every current-policy source drops the removed row too. |
| FD16 | 14 | `Store.Plan` returns a plan after a delivery flips the delivered spec | `internal/commitment/repository/plan_after_delivery_test.go` (`TestCommitmentPlanAfterDelivery`) | A source set that binds settled deliverables refuses at the flipped spec. |
| FD17 | 15 | `Store.Plan` refuses a proposal that keeps a tickets-only binding whose folder a later `main` commit deleted | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | A fix that skips every deliverable at `main` lets a kept deleted folder plan. |
| FD18 | 16 | `Store.Plan` returns a plan for a proposal that approves a changed staged spec at its new `main` identity | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | The current plan binds the old identity and refuses at `main`. |
| FD19 | 17 | After the FD12 plan is approved, `Store.AuthorizeCandidate` of the planning checkout's committed tree returns no error | planned TestCommitmentPlanSurvivesRemovedDeliverable in internal/commitment/repository/plan_after_delivery_test.go | `approvedTransition` recomputes `BuildPlan`, so a source set that differs at the commit refuses there. |
| FD20 | 18 | The plan identity derivation returns one identity for the sources `FT9`, `FT1` and for the sources `FT1`, `FT9` | planned TestPlanIdentityIgnoresSourceOrder in internal/commitment/authority_internal_test.go | A hash over the caller's order returns two identities. |
| FD21 | 19 | `BuildPlan` lists `FT1` before `FT9` when outcome `A` owns `FT9` and the later outcome `B` owns `FT1` | planned TestCommitmentPlanSourcesAreCanonical in internal/commitment/authority_test.go | A sort applied only inside the hash leaves the receipt in traversal order. |
| FD22 | 20 | `BuildPlan` lists the deliverable `spec` at `specs/a/spec.md` before the deliverable `spec` at `specs/b/spec.md` when outcome `A` approves the `b` path | planned TestCommitmentPlanSourcesAreCanonical in internal/commitment/authority_test.go | A sort on the identifier alone leaves equal identifiers in traversal order. |

Not covered: story 21 — a reviewed exclusion; the edge inventory states why it is safe.
Not covered: story 22 — a reviewed exclusion; Out of scope prices it.

### Edge inventory

The hostile-input classes of `projects/benchkit.md` reach no new parse surface here. The sort and the predicate read only identifiers that `validateSource` already restricts to `^[A-Za-z][A-Za-z0-9._-]*$`, and paths that `repositoryPath` already cleans. The absent-versus-empty class is FD1 and FD2.

**Won't handle** — `bench commitment start` of a legacy obligation-free binding — only a legacy policy holds one, and FD6 refuses its landing with the repair command.

**Won't handle** — numeric order of row identifiers, so `FT10` sorts before `FT9` — byte order is one total order, and FD21 observes it.

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
- `internal/commitment/command_test.go`
- `internal/commitment/commitmenttest/`
- `internal/commitment/repository/closure_test.go`
- `internal/commitment/repository/plan_after_delivery_test.go`
- `internal/worktree/commitment_light_landing_test.go`
- `internal/worktree/commitment_landing_fixture_test.go`
- `reviews/commitment-delivery-integrity.md`

## Out of scope

- A re-pinned roadmap row source. Today a plan that moves a kept row to its new identity binds the old current-policy identity too, so the plan refuses at `main`. The fix needs its own decision about the order of a row edit and its re-pin. Estimate: 4 edits, 2 gate runs.
- A settled state for an undelivered binding whose outcome other deliveries closed. This changes `Unsettled` for every reader. Estimate: 6 edits, 2 gate runs.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| "A plan that binds a deliverable to an outcome with sources lists the obligations that the deliverable closes." | FD1, FD2, FD3, FD4, FD5 |
| "Today the plan accepts an empty obligation list." | FD1, FD2 |
| "The completion landing then records no delivery fact, because `Deliver` returns the policy unchanged." | FD6, FD8, FD11 |
| "Either the plan refuses that binding, or the completion landing refuses a completion that records no fact." | FD1 to FD11; the contestable call takes both sides |
| "A plan that settles or removes a deliverable does not need that deliverable on `main`." | FD12, FD13, FD14, FD16, FD19 |
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
- Callers of `BuildPlan`: `Store.Plan`, `approvedTransition`, and the four tests in `authority_test.go`.
- Callers of `Deliver`: `Store.Delivered` alone in production, and `delivery_test.go`.
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
