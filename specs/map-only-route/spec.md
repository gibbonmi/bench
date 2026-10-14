# Map-only route

Status: staged

Decision source: reviewer-confirmed conversation, 2026-10-13.

Closed decisions: the reviewer chose an opt-in map-only route and kept the full route as the default.
The reviewer reopened the scorecard condition "repeat on a code-seam map first", and queued this spec after FT168-probe-batch.

Verification log: 3 iteration(s) to accept — the spec review took 2 iterations, because iteration 1 found two fence orders.
The ticket review accepted in 1 iteration. The existing rows-owned and rows-membership preflight checks enforce MR9.

## Problem

On the `guidance-weight` map, the full spec route cost about 6.6M subagent tokens before the build started.
On the `guidance-weight` map, one map-only slicer matched the full route's tickets at about one twelfth of the tokens.
It also used about one quarter of the wall time.
That slicer also found one build defect that the full route missed.
The pipeline still needs a spec file, a coverage map, a completion plan, and a fence, so map-only tickets cannot enter `/bench-implement-spec` today.

## Solution

The reviewer can select a map-only route for a ready compiled map whose decisions are exact predicates.
On that route, one fresh slicer writes the stories and coverage rows from the map, cuts the tickets, and records a gap log.
The coordinator then completes a thin spec, runs the closure proposal and one ticket review, and presents the usual approval table.
The full route stays the default, and the route stops and returns to the reviewer when the map is not spec-grade.

## User stories

Line: opus / high.
Implementation-line reason: the only chunk is guidance prose with new anchors, and the spec states each sentence.
The leverage override applies, because this guidance steers every later spec stage.
Harder chunks: none.

### Selection

1. As the reviewer, I want to select a map-only route for a ready compiled map, so that a spec-grade map skips the research cost.
2. As the reviewer, I want the full route to stay the default, so that a spec stage without my selection runs as it does today.
3. As the coordinator, I want the route to accept only a ready compiled map with no `bench maps` diagnostic, so that other sources never take it.
4. As the reviewer, I want the guidance to say that no check grades predicate quality, so that I know my selection is the eligibility gate.
5. As the coordinator, I want the route to compile and move the map under the Entry contract, so that the map lifecycle stays one procedure.

### Route steps

6. As the coordinator, I want one fresh slicer on `craft-line`'s ticket slicing line to author the route, so that one context holds the map.
7. As the coordinator, I want the route to omit the research, the spec fork, the spec review, and the slicing readers, so that cost falls.
8. As the slicer, I want to write the stories and coverage rows before I cut tickets, so that each `Covers:` line names a declared row.
9. As the slicer, I want each ticket to name its rows in `Covers:`, so that every row has an owner ticket.
10. As the slicer, I want to do my own tree research under the slicing checks' `Fence-expansion evidence` classes, so that each `Writes:` line keeps its closure.
11. As the slicer, I want to record a gap log in the thin spec's `Further notes`, so that every map gap is visible at sign-off.
12. As the slicer, I want each gap-log entry to carry one of four classes, so that the reviewer finds each reviewer-owned decision.
13. As the coordinator, I want the chunk table and the completion plan after the tickets, under an interim fence, so that the proposal runs.
14. As the coordinator, I want to run `--propose-writes` for each ticket on the committed thin spec, so that every closure file joins `Writes:`.
15. As the coordinator, I want one ticket review on `craft-line`'s ticket review line with a cap of 2, so that one independent grade remains.
16. As the coordinator, I want the fence, the plan-only preflight, and the approval table after the ticket review, so that build entry is unchanged.
17. As the reviewer, I want each class (d) gap-log entry as an approval-table decision row, so that I decide each one at sign-off.
18. As the coordinator, I want the verification log to count only the ticket round and name the route, so that it states what ran.
19. As the reviewer, I want the route to keep the spec-stage close steps unchanged, so that the binding and the landing stay one procedure.

### Stop and return

20. As the slicer, I want to stop when a map decision is not an exact predicate, so that I do not invent a decision.
21. As the slicer, I want to stop when the map leaves a code seam open, so that seam design gets the full route.
22. As the reviewer, I want a stop to return to me with a recommendation for the full route, so that I choose the next route.

### Thin spec

23. As a build orchestrator, I want the thin spec to hold every section that the checks read, so that the build entry accepts it.
24. As a ticket author, I want the thin spec to name the compiled map as its decision source, so that I find the decisions.

### Pointer and enforcement

25. As a reader of `/bench-write-spec`, I want one sentence in `Who runs this phase` that names the route reference, so that I find the route from the phase.
26. As a maintainer, I want `/bench-write-spec` to stay at its 87-line budget, so that no budget row changes.
27. As a maintainer, I want every anchored full-route sentence and its line wraps to stay byte-identical, so that the existing anchors and canaries stay green.
28. As a maintainer, I want a `Require` anchor on each load-bearing route sentence, so that an omission turns the gate red.
29. As a maintainer, I want each anchor's test row written apart from the registry, so that a dropped rule cannot define itself green.

### Readers

30. As a maintainer, I want the spec-authoring entry of `projects/benchkit.md` to name the route, so that the routing prose stays true.
31. As a reader of the field guide, I want both write-spec passages to name the route, so that the guide stays true.
32. As a user, I want a `CHANGELOG.md` entry for the route, so that a release note shows the change.

### Reviewed exclusions

33. As the reviewer, I want `README.md` and ADR 0022 left as they are, so that their statements of the default stay true.
34. As the reviewer, I want no machine check of predicate quality in this spec, so that the route ships as a reviewer judgment.
35. As the reviewer, I want the route to stay opt-in, so that one data point does not change the default.
36. As the reviewer, I want the `craft-research` pair rule unchanged, so that the full route keeps its readers.
37. As the reviewer, I want the scorecard text to change at the spec-stage close, so that the scorecard keeps one update point.

## Implementation decisions

### Eligibility

The reviewer selects the route. The route accepts only a ready compiled map, which `bench maps` reports with no diagnostic.
The reviewer also judges that each decision is an exact predicate. The map validator does not read the answer text of a ticket, so no check grades that judgment (`internal/maps/validation.go:87-113`).
A reviewer-confirmed conversation and a named reviewed artifact never take the route.

### Placement

The route's procedure lives in one new reference, `.agents/skills/bench-craft-spec/references/map-only-route.md`.
A reference file is outside the line budget (`internal/conformance/prose_budget_test.go:35-37`), and the `craft-spec` skill file holds 155 of its 155 lines.
`/bench-write-spec` gains one pointer sentence at the end of the paragraph that starts "Each delegate returns its findings" in `Who runs this phase`.
That paragraph is one physical line, and no canary pins its bytes, so the file stays at 87 lines.
The proposed pointer sentence is: "On a reviewer-selected map-only route, `craft-spec`'s `references/map-only-route.md` replaces the research, the fork, the spec review, and the slicing research."
The edit adds no H2 to `/bench-write-spec`, and it changes no anchored sentence or line wrap.

### Route steps

The reference states the steps in this order.

1. Compile and move the map as the Entry contract states.
2. Charge one fresh slicer on `craft-line`'s ticket slicing line. The route has no research fan-out, no spec fork, no spec review, and no paired slicing readers.
3. The slicer writes the thin spec's stories and coverage rows from the decision predicates first, so the row IDs exist before the tickets (`internal/preflight/ticket_rows.go:41-74`).
4. The slicer cuts the tickets under `craft-tickets` and the slicing checks, and each `Covers:` line names declared rows.
5. The slicer records the gap log in the thin spec's `Further notes`. Each entry takes one class: (a) map discretion, (b) a seam or test choice, (c) a missing fact, or (d) a reviewer decision.
6. After every ticket exists, the coordinator completes the chunk table and the completion plan. Until step 8, the fence holds only the review pickup.
   Then the coordinator commits the thin spec and the tickets, and runs `--propose-writes` for each ticket.
   The proposal reads the committed thin spec (`internal/preflight/preparation.go:41`), and it tolerates only the closure-check reds and the `fence-writes` red (`internal/preflight/proposal.go:14`).
7. One ticket review runs on `craft-line`'s ticket review line, with an iteration cap of 2.
8. After the ticket review accepts, the coordinator writes the fence section, runs the plan-only preflight, and presents the approval table.
9. Each class (d) entry is a decision row in the approval table. These rows are not late clarification questions, so the two-question limit does not apply to them.
10. The verification log counts only the ticket round, and its note names the map-only route.

The spec-stage close in `/bench-write-spec`'s `Exit handoff` applies unchanged.

### Stop and return

The slicer stops when a map decision is not an exact predicate or when the map leaves a code seam open.
The slicer also stops when the map holds no decision ticket.
The stop returns to the reviewer with a recommendation for the full route, and the slicer writes no further ticket.

### Thin spec

A thin spec holds the parts that the coverage check and the build preflight read (`internal/preflight/gather.go:48-70`), and it holds the gap log for the sign-off.

- `Status: staged`, one `Decision source:` line that names the compiled map, a `Verification log:` line, and the `Line:` fields.
- `## User stories`, with at least one story for each map decision.
- `### Acceptance coverage map`, with the `row` column.
- `## Implementation chunks`, with the chunk table and the `bench-completion-plan` fence.
- `## Ownership fences`, which lists `reviews/<slug>.md`.
- `## Further notes`, with the gap log.

`Problem` and `Solution` stay short, and the thin spec omits a section that no check reads and the map already settles.

### Enforcement

New `RequireInSection` anchors pin each load-bearing route sentence in the reference, and one `RequireInSection` anchor pins the pointer under `Who runs this phase`.
The anchors join the write-spec slicing family in `internal/anchors/registry_ticket_passes.go`.
A new test function in `internal/anchors/registry_ticket_passes_test.go` writes each needle and diagnostic apart from the registry, as `TestWriteSpecChunkedSlicing` does.

### Readers

The spec-authoring entry in `projects/benchkit.md` gains one sentence that names the route.
The two write-spec passages in `docs/field-guide.html` gain one clause each.
`CHANGELOG.md` gains an Unreleased entry.
`README.md` and ADR 0022 stay as they are, because each one states the default route and stays true.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| MR-A / 1-add-the-map-only-route.md | The route reference, the write-spec pointer, the route anchors and their test rows, and the reader updates | MR1-MR34 | planned TestWriteSpecMapOnlyRoute, the live anchors check, the prose and budget checks | no |

One chunk holds the change, because every part is guidance prose and its anchors, and no part creates a seam for a sibling.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"map-only-route-build-20261013","orchestrator_session":"claude:mr-orchestrator-20261013","author_limit":1,"assignments":{"1-add-the-map-only-route.md":[{"session":"agy:mr_t01","assignment":"61be98876f5fd0ae0c1ba1b9ca95c354","model":"gemini-4-argon-high","effort":"high","source":"108a7288d42baa192a0c45d281bbec34a5aa9104","native_ref":"agy-print:mr_t01"},{"session":"claude:mr_t01_rep1","assignment":"61be98876f5fd0ae0c1ba1b9ca95c354","model":"opus","effort":"low","source":"717fad7651351ce70c7ccdc49050992cc829cc97","native_ref":"claude-agent:mr_t01_rep1","predecessor":"agy:mr_t01","trigger":"user-directed","stopped":"agy:mr_t01 ended after its ticket commit 6d0989b8 and its MR-A record commit fe033446","preserved":"no uncommitted work; ticket commit 6d0989b8 and records through 0fc73c29 on the map-only-route-build branch"}]}},"chunks":[{"id":"MR-A","tickets":["1-add-the-map-only-route.md"],"verification":[{"id":"guidance","command":"bench test --check docs-currency-workflow","ticket":"1-add-the-map-only-route.md"},{"id":"prose","command":"bench test --check prose","ticket":"1-add-the-map-only-route.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"1-add-the-map-only-route.md"},{"id":"anchors-package","command":"bench test --package ./internal/anchors","ticket":"1-add-the-map-only-route.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"1-add-the-map-only-route.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"1-add-the-map-only-route.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"1-add-the-map-only-route.md"},{"id":"central-property","command":"bench test --package ./internal/anchors","probe":"Delete the stop-and-return anchor from internal/anchors/registry_ticket_passes.go. TestWriteSpecMapOnlyRoute must turn red.","ticket":"1-add-the-map-only-route.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/map-only-route/spec.md"},{"id":"anchors-package","command":"bench test --package ./internal/anchors"},{"id":"command-package","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

- A good test reads the live guidance through the anchor harness and fails when a route sentence or the pointer is missing from its section.
- The seam is the anchor registry and its independent test rows, with `TestWriteSpecChunkedSlicing` as prior art.
- The gate observes the change through `bench test --check docs-currency-workflow`, which runs the live anchors, and through the prose and budget checks.
- The existing write-spec anchors and canaries are the regression guard for the full route, so this spec adds no canary for them.
- Rows on reader prose that no anchor pins are review-owned.

### Seam diagram

    trigger: the gate's docs-currency-workflow check and the anchors package tests
        │
        ▼
    live guidance files  ──▶  [ anchor registry: route needles ]  ──▶  pass, or a named diagnostic
                      ◀ tests attach here: TestWriteSpecMapOnlyRoute drives each needle on a synthetic tree, and the live check reads the tree

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| MR1 | 1 | The route reference states that the reviewer selects the route for a ready compiled map | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A reference without the selection sentence lets any session take the route. |
| MR2 | 2 | The route reference states that the full route stays the default | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A reference that drops the default sentence reads as a replacement of the full route. |
| MR3 | 3 | The route reference states that the map must be compiled, ready, and clean under `bench maps` | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A route without the readiness sentence accepts a shaping map. |
| MR4 | 4 | The route reference states that no check grades predicate quality and that the reviewer judges it | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A reference that implies a check misleads the reviewer about the gate. |
| MR5 | 5 | The route reference sends the map compile and move to the Entry contract of `/bench-write-spec` | review-owned | A second copy of the move procedure drifts from the Entry contract. |
| MR6 | 6 | The route reference charges one fresh slicer on `craft-line`'s ticket slicing line | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A route without the line sentence lets the coordinator slice on its own line. |
| MR7 | 7 | The route reference states that the route has no research fan-out, no spec fork, no spec review, and no paired slicing readers | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A reference that omits the list leaves a session to run the full route's steps by habit. |
| MR8 | 8 | The route reference orders the stories and coverage rows before the tickets | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | Tickets cut first carry `Covers: none`, and rows-owned turns red at the plan-only preflight. |
| MR9 | 9 | Each ticket's `Covers:` line names declared rows of the thin spec | `bench preflight build <slug> --plan-only` (rows-owned and rows-membership) | A ticket that names no row leaves a row without an owner. |
| MR10 | 10 | The route reference sends the slicer's tree research to the slicing checks' `Fence-expansion evidence` classes | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A slicer without the classes misses fixture and anchor closure files. |
| MR11 | 11 | The route reference puts the gap log in the thin spec's `Further notes` | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A gap log kept outside the spec never reaches the sign-off. |
| MR12 | 12 | The route reference names the four gap-log classes (a) to (d) | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A log without class (d) hides the decisions that only the reviewer makes. |
| MR13 | 13 | The route reference completes the chunk table and the completion plan after every ticket exists, under the interim fence | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A plan written before the tickets names missing files, and the completion-plan check turns red. |
| MR14 | 14 | The route reference runs `--propose-writes` for each ticket on the committed thin spec, with its chunk table and completion plan | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A proposal run on a partial spec refuses, and the closure files never join `Writes:`. |
| MR15 | 15 | The route reference runs one ticket review on `craft-line`'s ticket review line with an iteration cap of 2 | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A route without the review sentence ships tickets with no independent grade. |
| MR16 | 16 | The route reference writes the fence, runs the plan-only preflight, and presents the approval table after the ticket review accepts | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A route that skips the plan-only preflight stages a ticket graph that the build entry refuses. |
| MR17 | 17 | The route reference makes each class (d) entry a decision row of the approval table | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A route without the sentence routes those decisions through the two-question limit and drops some. |
| MR18 | 18 | The route reference makes the verification log count only the ticket round and name the route | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A log that sums a spec round that never ran states a false count. |
| MR19 | 19 | The route reference sends the spec-stage close to the `Exit handoff` of `/bench-write-spec` | review-owned | A second close procedure in the reference drifts from the binding steps. |
| MR20 | 20 | The route reference stops the slicer when a decision is not an exact predicate | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A slicer without the stop writes a decision that the reviewer never made. |
| MR21 | 21 | The route reference stops the slicer when the map leaves a code seam open | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A slicer without the stop designs a seam with no research or spec review. |
| MR22 | 22 | The route reference returns a stop to the reviewer with a recommendation for the full route | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A stop with no return leaves the spec stage idle with no decision request. |
| MR23 | 23 | The route reference lists the thin-spec parts that the coverage check and the build preflight read, and it names the gap log apart from them | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A thin spec without the fence or the plan fails the build preflight bootstrap. |
| MR24 | 24 | The route reference requires the thin spec's `Decision source:` line to name the compiled map | review-owned | A thin spec that names the conversation hides the map from the ticket author. |
| MR25 | 25 | `/bench-write-spec` holds the pointer sentence to `references/map-only-route.md` under `Who runs this phase` | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A route reference with no pointer is unreachable from the phase. |
| MR26 | 26 | `/bench-write-spec` holds 87 lines | `bench test --check guidance-prose-budgets` | A pointer on a new line breaks the 87-line budget. |
| MR27 | 27 | Every existing write-spec anchor and canary stays green | `bench test --check docs-currency-workflow` (the live anchors over `/bench-write-spec`) | A reworded full-route sentence or a moved line wrap reds its anchor or its canary. |
| MR28 | 28 | The registry holds one `RequireInSection` anchor for each route sentence of rows MR1-MR4, MR6-MR8, MR10-MR18, MR20-MR23, MR25, and MR34 | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A sentence with no anchor can leave the reference with no red. |
| MR29 | 29 | `TestWriteSpecMapOnlyRoute` states each needle and diagnostic as its own literal | review-owned | A test that reads the registry entries passes when a rule drops out of both. |
| MR30 | 30 | The spec-authoring entry of `projects/benchkit.md` names the map-only route | review-owned | A routing entry that names only the fork misstates the reviewer-selected route. |
| MR31 | 31 | Both write-spec passages of `docs/field-guide.html` name the map-only route | review-owned | A guide that names only the fan-out misstates the phase. |
| MR32 | 32 | `CHANGELOG.md` holds an Unreleased entry for the map-only route | review-owned | A missing entry hides the change from a release note. |
| MR33 | 33 | `README.md` and ADR 0022 hold no edit in the build diff | review-owned | An edit there restates the route in a second place. |
| MR34 | 22 | The route reference stops the slicer when the map holds no decision ticket | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecMapOnlyRoute`) | A slicer on an empty map writes a thin spec with zero stories, and no check turns red. |

Not covered: story 34 — a reviewed exclusion whose cost the Out of scope section prices.
Not covered: story 35 — a reviewed exclusion whose cost the Out of scope section prices.
Not covered: story 36 — a reviewed exclusion, because the full route keeps its readers unchanged.
Not covered: story 37 — a reviewed exclusion, because the spec-stage close owns the scorecard update.

### Edge inventory

The canonical edge classes at the anchor seam and the route procedure:

- A map with no decision ticket — the slicer has nothing to slice, so it stops and returns to the reviewer under story 22.
- A map where only some decisions are exact predicates — story 20 stops the route at the first such decision.
- A gap log with more than two class (d) entries — story 17 sends each one to the approval table, so no limit caps them.
- A ticket review that rejects twice — the cap of 2 stops the round, and the reviewer gets the stop, as on the full route.
- A pointer sentence placed inside a canary-pinned wrap — story 27 and row MR27 catch the moved bytes.

**Won't handle:** a slicer context past about 300k tokens — the route is for a map that one slicer can hold. The Process checkpoint rule still applies.
**Won't handle:** a hostile-input walk — the change is guidance prose and its anchors, and it adds no parser or input surface.

## Ownership fences

Ticket slicing wrote this fence. It equals the ticket `Writes:` union plus the review pickup.

- `.agents/commands/bench-write-spec.md`
- `.agents/skills/bench-craft-spec/references/map-only-route.md`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `docs/field-guide.html`
- `internal/anchors/registry_craft_research.go`
- `internal/anchors/registry_craft_research_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `projects/benchkit.md`
- `reviews/map-only-route.md`
- `tests/canary/guidance-prose-budgets/over-budget-skill/`
- `tests/canary/line-routing/line-binding-prose-drift/`
- `tests/canary/skill-description-budgets/`
- `tests/canary/skills-index-command-adapters/`
- `tests/canary/workflow-guidance-anchors/`

## Out of scope

- A machine check that a decision is an exact predicate: 6 edits, 2 gate runs. It needs a predicate grammar in `internal/maps` and a new `bench maps` diagnostic.
- The map-only route as the default for a ready map: 4 edits, 1 gate run. It needs the second-map comparison first.
- A change to the `craft-research` pair rule: 3 edits, 1 gate run. The 2026-10-13 learning entry owns its evidence.
- The scorecard text in `capture/agent-performance/claude-models.md`: the spec-stage close updates it.

## Further notes

### Flagged for reviewer veto

- D1 eligibility: the reviewer selects the route for a ready compiled map with exact-predicate decisions, and no check grades that judgment.
- D2 placement: the procedure lives in a new `craft-spec` reference, and `/bench-write-spec` gains one pointer sentence on an existing line, so no budget changes.
- D3 route order: the slicer writes rows before tickets, and the coordinator completes the thin spec before `--propose-writes`.
  This order reverses the experiment's tickets-first order, because rows-owned needs declared rows.
- D4 gap log: the four gap-log classes come from the experiment, and each class (d) entry is an approval-table decision row.
  A class (d) entry is not a late clarification question.
- D5 stop and return: a non-predicate decision or an open code seam stops the route with a recommendation for the full route.
- D6 thin spec: it holds only the parts that the coverage check and the build preflight read, plus the gap log.
  Its verification log counts only the ticket round, because no spec round runs.
- D7 enforcement: the route's anchors join the write-spec slicing family with a new independent test function.
- D8 readers: `README.md` and ADR 0022 stay unedited as statements of the default.

### Known overlap

The staged `specs/guidance-weight` spec (FT100) is not yet committed to an outcome. Its tickets 2 and 5 write `/bench-write-spec` and raise its line budget row.
This change lands first in the approved order. The build-entry staleness pass of the later build owns the rebase of either one.

### Source sentence to row

| source sentence | rows |
| --- | --- |
| Opt-in map-only route that the reviewer selects | MR1, MR3, MR4, MR25 |
| The full route stays the default | MR2, MR27 |
| One fresh slicer cuts tickets from the map | MR6, MR7, MR8, MR9, MR10 |
| A thin spec from the tickets: decision source, coverage map, chunk table, completion plan, fence | MR13, MR23, MR24 |
| One ticket review and the plan-only preflight | MR14, MR15, MR16 |
| Queue after FT168-probe-batch | the spec-stage close binding |
