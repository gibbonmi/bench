# Fresh ticket authors

Status: staged

Decision source: the reviewer-confirmed current conversation, 2026-09-24.

Verification log: 2 iteration(s) to accept — iteration 1 rejected on 2 blockers and 4 minor findings (R1–R6), and the author folded them. Iteration 2 accepted with no finding above the blocking bar.

## Problem

A spec-backed build keeps one retained implementation session for every ticket, repair, and reconciliation. The FT71 build made 1,126 model requests at an average context of about 445k tokens, so 501M cached-read tokens. The context length, not the work, was the cost driver. Each ticket charge also made the author read the full spec and three skill documents again, page by page.

Fresh narrow review axes cost about 20k to 330k tokens each, because each axis starts small and reads only its delta. The build has no equal rule. The operating guide, the phase commands, the skills, the field guide, the README, and ADR 0021 state the single-author contract. About 20 anchors pin it.

## Solution

Every spec-backed build gives each ticket a fresh author session. The build records the authors in the existing version 2 delegate plan with an author limit of 1. So the record checks accept one author for each ticket with no code change. The authors stay on the declared line, and a tier move still asks the reviewer. An explicit `--delegate` adds concurrent authors and the full tier range.

The author charge is narrow: the ticket, its coverage rows, its `Writes:` fence, the evidence identity, and the declared line. The author binds the evidence, fetches the metadata and ticket pages, and reads the rest with targeted reads. The orchestrator keeps a small context, refreshes `bench handoff` at each chunk checkpoint, and reconciles the final source. A post-review repair goes to a fresh session that the plan records as a new assignment.

## User stories

Line: opus / high.
Implementation-line reason: each chunk rewrites shared guidance that steers every later build, so the leverage override routes it high. The rows are exact anchor needles, the seams exist, and the anchor registry reds each dropped or retired sentence.
Harder chunks: none.

Fresh authors:

1. As a reviewer, I want each ticket of a spec-backed build to get a fresh author session, so that the author context stays small.
2. As a reviewer, I want the build to record its authors in the limit-1 delegate plan, so that the record checks name each ticket's author.
3. As a reviewer, I want each fresh author to stay on the declared line, so that a tier move still asks me first.
4. As a reviewer, I want an explicit `--delegate` to add concurrent authors and the full tier range, so that the opt-in keeps its meaning.
5. As a reviewer, I want the rule to cover every spec-backed build, with or without `--full`, so that no build keeps a long context.
6. As a reviewer, I want the retired single-author sentence refused, so that it cannot return.

Narrow author charge:

7. As a ticket author, I want a charge with only the ticket, its rows, its fence, and the evidence, so that I start small.
8. As a ticket author, I want to fetch only the metadata and ticket pages, so that I read the spec and skills with targeted reads.
9. As a reviewer, I want the retired full-retrieval build sentences refused, so that no guidance sends an author through every evidence page.

Small orchestrator:

10. As an orchestrator, I want to read manifests, returns, and verdicts, not code, so that my context stays small.
11. As an orchestrator, I want to refresh `bench handoff` at each chunk checkpoint, so that a later orchestrator can resume from the tree.
12. As a reviewer, I want the orchestrator to reconcile the final acceptance and integration, so that the version 2 integration role owns completion.

Repairs:

13. As a reviewer, I want each post-review repair in a fresh session recorded as a new assignment, so that the plan names the repair author.
14. As a reviewer, I want the retired repair-ownership sentence refused, so that the old rule cannot return.

Skills and phase commands:

15. As a ticket author, I want the continuation policy of `craft-line` to govern each ticket author, so that the progress and stop rules still apply.
16. As a slicer, I want a ticket sized to one fresh author context, so that the ticket fits its session.
17. As a spec author, I want the write-spec handoff to recommend the line for fresh ticket authors, so that the next phase starts right.
18. As a reviewer, I want accepted review findings routed to a fresh repair author, so that the review phase agrees with the build phase.
19. As a reader of the field guide and the README, I want the current authorship rule stated, so that the public docs match the tree.
20. As a reader of the ADRs, I want the decision recorded as the current state, so that ADR 0021 no longer states the single-author contract.

Reviewed exclusions:

21. As a reviewer, I want the light path to stay in the main session, so that a one-ticket change pays no dispatch.
22. As a reviewer, I want `internal/reviewrecord` unchanged, so that this spec ships no record-schema change.
23. As a reviewer, I want the small-output evidence command left to the CLI research row, so that ADR 0022 stays in force here.

## Implementation decisions

- `.bench/BENCH.md` owns the authorship rule. Its "Retain implementation authorship" paragraph becomes a fresh-ticket-author paragraph, and its "Delegate a full run only on my request" paragraph keeps only what `--delegate` adds. Every other file points to that owner and restates no rule.
- The build plan uses the version 2 delegate form. Before the first dispatch, the orchestrator writes the `execution` block. The block names mode `delegate`, an author limit of 1, the run id, its own session, and an empty history for each ticket. That plan edit takes the ordinary plan amendment.
- The same amendment sets the plan version to 2 and splits each chunk verification into one verification for each ticket. Version 2 requires each chunk verification to name one of that chunk's tickets, and each ticket to own one.
- `internal/reviewrecord` allows one session for one ticket. So a repair session joins that ticket's history as a replacement with the trigger `user-directed`, the stopped-writer evidence, and the preserved source.
- A finding that touches several tickets takes one repair session for each affected ticket. The repair session becomes that ticket's verifier, so it reruns that ticket's verification.
- At final reconciliation, the orchestrator does not repair. A finding there goes to a fresh repair session for the ticket whose `Writes:` line holds the path. The orchestrator then reruns the final verification.
- A finding on a path that no ticket's `Writes:` line holds is a material acceptance shortfall, so the build stops and reports.
- `craft-line` keeps its section name "Retained implementation continuation". So ticket 1 keeps the continuation-policy pointer sentences in `.bench/BENCH.md` and `.agents/commands/bench-implement-spec.md` byte for byte, and their section anchors stay valid.
- The build delivery rule follows the narrow review read that landed at `5f1453fc`. The charge-evidence guidance test names the build family's delivery prerequisite `a narrow author read`, beside the review family's `a narrow axis read`.
- A fresh author works in the one integration worktree, serially, in `Blocked by:` order. Only `--delegate` with a limit above 1 runs authors at the same time.
- The Claude Code auto-mode classifier refuses an agent edit to `.bench/BENCH.md` as a self-modification, because `CLAUDE.md` imports it. Ticket 1 therefore names a reviewer approval step: the reviewer makes that edit or grants a permission rule for it before the author starts.
- ADR 0023 records the decision. ADR 0021's first, second, third, and fifth consequences change to the decided state and point to ADR 0023.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| FA / `1-route-each-ticket-to-a-fresh-author.md`, `2-align-the-line-and-delegation-skills.md`, `3-align-the-phase-commands-and-docs.md` | Every spec-backed build gives each ticket a fresh author under a narrow charge, and every guidance file, doc, and ADR states that rule. | FA1, FA2, FA3, FA4, FA5, FA6, FA7, FA8, FA9, FA10, FA11, FA12, FA13, FA14, FA15, FA16, FA17, FA18, FA19, FA20, FA21, FA22, FA23, FA24, FA25, FA26, FA27, FA28, FA29, FA30 | `bench test --check docs-currency-workflow`, `bench test --package ./internal/conformance --run TestRootConformance` | no |

## Testing decisions

- A good test reads the live guidance through the anchor registry. A dropped sentence reds a Require row, and a returned sentence reds a Forbid row.
- The anchor registries under `internal/anchors` and the conformance tests `retained_workflow_test.go`, `charge_evidence_guidance_test.go`, `implementation_continuation_test.go`, and `docs_workflow_helpers_test.go` are the precedents. The narrow review change at `5f1453fc` is the nearest precedent.
- The `docs-currency-workflow` check and the root conformance test observe the feature. Each changed canary under `tests/canary/workflow-guidance-anchors` keeps its mutation red.

### Seam diagram

    trigger: `bench test --check docs-currency-workflow`, the gate's test phase
        │
        ▼
    guidance files  ──▶  [ anchor registry: Require and Forbid needles ]  ──▶  diagnostic or pass
                      ◀ tests attach here: add a needle row, and a canary mutation drops or restores its sentence

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| FA1 | 1, 3, 5 | `.bench/BENCH.md` requires a fresh author session for each ticket of every spec-backed build, on the declared line | planned Require needle in `internal/anchors/registry_ft311_review_dispatch.go` | A guide that keeps the single-author paragraph lacks the needle, so the check reds. |
| FA2 | 2 | `.bench/BENCH.md` requires the version 2 delegate plan with an author limit of 1 for the fresh authors | planned Require needle in `internal/anchors/registry_ft311_review_dispatch.go` | A guide that names no plan form lacks the needle, so the check reds. |
| FA3 | 4 | `.bench/BENCH.md` states that `--delegate` adds concurrent authors and the full tier range | planned Require needle in `internal/anchors/registry_retained_workflow.go` | A guide that keeps `--delegate` as the only route to per-ticket authors lacks the needle, so the check reds. |
| FA4 | 6 | `.bench/BENCH.md` does not contain "The retained implementation session writes production changes, tests, probes, and repairs." | planned Forbid needle in `internal/anchors/registry_ft311_review_dispatch.go` | A guide that keeps the retired sentence matches the Forbid needle, so the check reds. |
| FA5 | 7 | `.agents/commands/bench-implement-spec.md` requires the narrow author charge of ticket, rows, fence, evidence identity, and line | planned Require needle in `internal/anchors/registry_ft311_preparation.go` | A build phase that charges the author with the whole spec lacks the needle, so the check reds. |
| FA6 | 8 | `.agents/commands/bench-implement-spec.md` requires a narrow author read: a binding check, the metadata and ticket pages, and targeted reads | `internal/conformance/charge_evidence_guidance_test.go` (`TestEvidenceBoundedActionGuidance`) | A build family that keeps `verified delivery` fails the family prerequisite, so the test reds. |
| FA7 | 9 | `.agents/commands/bench-implement-spec.md` does not contain "Retrieve every required source with `bench preflight evidence <id>`, and follow each exact successor command until the stream ends." | planned Forbid needle in `internal/anchors/registry_ft311_preparation.go` | A build phase that keeps the full-retrieval sentence matches the Forbid needle, so the check reds. |
| FA8 | 9 | `.agents/commands/bench-implement-spec.md` does not contain "Build action requires verified delivery" | planned Forbid needle in `internal/anchors/registry_ft311_preparation.go` | A build phase that keeps the retired delivery rule matches the Forbid needle, so the check reds. |
| FA9 | 10 | `.agents/commands/bench-implement-spec.md` requires the orchestrator to read manifests, returns, and verdicts, not code | planned Require needle in `internal/anchors/registry_ft311_preparation.go` | A build phase with no orchestrator bound lacks the needle, so the check reds. |
| FA10 | 11 | `.agents/commands/bench-implement-spec.md` requires a `bench handoff` refresh at each chunk checkpoint | planned Require needle in `internal/anchors/registry_ft311_preparation.go` | A build phase that refreshes only at `--full` phase boundaries lacks the needle, so the check reds. |
| FA11 | 12 | `.bench/BENCH.md` requires the orchestrator to reconcile the final acceptance and integration | planned Require needle in `internal/anchors/registry_ft311_review_dispatch.go` | A guide that keeps the retained author as reconciler lacks the needle, so the check reds. |
| FA12 | 13 | `.bench/BENCH.md` requires one fresh repair session for each affected ticket, recorded as a new assignment with the trigger `user-directed`, that reruns that ticket's verification, and routes a final-reconciliation finding to a fresh repair session, not to the orchestrator, and treats a finding on a path that no `Writes:` line holds as a material acceptance shortfall | planned Require needle in `internal/anchors/registry_retained_workflow.go` | A guide that keeps the repair with the ticket author lacks the needle, so the check reds. |
| FA13 | 14 | `.bench/BENCH.md` does not contain "Production repairs stay with the recorded ticket author." | planned Forbid needle in `internal/anchors/registry_retained_workflow.go` | A guide that keeps the retired sentence matches the Forbid needle, so the check reds. |
| FA14 | 2 | `.bench/BENCH-reference.md` states that the plan amendment declares version 2 with a delegate execution block and splits each chunk verification into one verification for each ticket | planned Require needle in `internal/anchors/registry_retained_workflow.go` | A reference that still declares version 1, or keeps one shared chunk verification, lacks the needle, so the check reds. |
| FA15 | 15 | `.agents/skills/bench-craft-line/SKILL.md` applies its continuation section to each ticket author | planned RequireInSection needle in `internal/anchors/registry_retained_workflow.go` | A skill that keeps the retained-session wording lacks the needle, so the check reds. |
| FA16 | 3 | `.agents/skills/bench-craft-line/SKILL.md` requires a reviewer stop before a tier move of a fresh author outside `--delegate` | planned Require needle in `internal/anchors/registry_retained_workflow.go` | A skill that lets a limit-1 plan move tiers like `--delegate` lacks the needle, so the check reds. |
| FA17 | 1 | `.agents/skills/bench-craft-delegate/SKILL.md` points to `.bench/BENCH.md` as the owner of ticket author sessions | planned Require needle in `internal/anchors/registry_data.go` | A skill that keeps the retained-author pointer lacks the needle, so the check reds. |
| FA18 | 16 | `.agents/skills/bench-craft-tickets/SKILL.md` sizes a ticket to one fresh author context | planned Require needle in `internal/anchors/registry_data.go` | A skill that keeps the retained-session size lacks the needle, so the check reds. |
| FA19 | 6 | `.agents/skills/bench-craft-tickets/SKILL.md` does not contain "Spec-backed builds work the unblocked frontier in one retained implementation session." | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the retired sentence matches the Forbid needle, so the check reds. |
| FA20 | 17 | `.agents/commands/bench-write-spec.md` recommends the line for fresh ticket authors on one integration source | planned Require needle in `internal/anchors/registry_data.go` | A handoff that keeps "one retained session" lacks the needle, so the check reds. |
| FA21 | 18 | `.agents/commands/bench-review-implementation.md` routes accepted findings to a fresh repair author | planned Require needle in `internal/anchors/registry_ft311_review_dispatch.go` | A review phase that keeps the retained-session return lacks the needle, so the check reds. |
| FA22 | 19 | `docs/field-guide.html` states that each ticket gets a fresh author session | planned Require needle in `internal/anchors/registry_data.go` | A guide page that keeps the retained-authorship sentence lacks the needle, so the check reds. |
| FA23 | 19 | `README.md` states that each ticket gets a fresh author session | planned Require needle in `internal/anchors/registry_retained_workflow.go` | A README that keeps "Implementation stays with its retained author" lacks the needle, so the check reds. |
| FA24 | 20 | ADR 0023 records the fresh-ticket-author decision, and ADR 0021 points to it | review-owned: an ADR has no mechanical seam here | Review grades both ADRs against this spec's decisions. |
| FA25 | 21 | `.agents/commands/bench-drain.md` keeps the implement-now light path in the main session | `internal/anchors/registry_data.go` (existing Require needle "Write its one ticket file. Implement that ticket in the retained session under `craft-line`.") | A drain phase that moves the light path to a fresh author drops the needle, so the check reds. |
| FA26 | 14 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "Repairs return to the retained implementation session." | planned Forbid needle in `internal/anchors/registry_retained_workflow.go` | A skill that keeps the retired sentence matches the Forbid needle, so the check reds. |
| FA27 | 12 | `.bench/BENCH.md` does not contain "After the last chunk, the retained author reconciles overall acceptance and integration." | planned Forbid needle in `internal/anchors/registry_ft311_review_dispatch.go` | A guide that keeps the retired sentence matches the Forbid needle, so the check reds. |
| FA28 | 14 | `.bench/BENCH.md` does not contain "Return findings to the retained author and obtain current repair coverage." | planned Forbid needle in `internal/anchors/registry_ft311_review_dispatch.go` | A guide that keeps the retired sentence matches the Forbid needle, so the check reds. |
| FA29 | 18 | `.agents/commands/bench-review-implementation.md` does not contain "Accepted findings return to the retained `/bench-implement-spec` session." | planned Forbid needle in `internal/anchors/registry_ft311_review_dispatch.go` | A review phase that keeps the retired sentence matches the Forbid needle, so the check reds. |
| FA30 | 13 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` permits a `user-directed` replacement for a post-review repair under "Delegated author transfer" | planned Require needle in `internal/anchors/registry_retained_workflow.go` | A discipline that drops the entry lacks the needle, so the check reds. |

Not covered: story 22 — the reviewed exclusion changes no behavior, and the Won't handle lines record it.
Not covered: story 23 — the reviewed exclusion changes no behavior, and the Out of scope section prices it.

### Edge inventory

The canonical edge classes, walked at the anchor seam:

- A retired sentence that returns: each retired sentence takes its own Forbid row (FA4, FA7, FA8, FA13, FA19, FA26, FA27, FA28, FA29).
- A sentence reflowed across lines: the registry matches a needle across a line break, and `charge_evidence_guidance_test.go` already reads a wrapped sentence.
- A second spelling of the rule: every file points to `.bench/BENCH.md`, and review grades a paraphrase.
- A canary that mutates a changed sentence: each changed canary's `MUTATE.json` names the new sentence, so its mutation stays red.
- Hostile input: this spec changes guidance and anchor rows only, so no operator input reaches a new surface.

**Won't handle** lines:

- The light path — a one-ticket change stays in the main session, because a dispatch costs more than its context. FA25 survives.
- A serial execution mode in `internal/reviewrecord` — the version 2 delegate plan with an author limit of 1 records the same facts. FA2 survives.
- One author for two tickets — `internal/reviewrecord` refuses one session for two tickets, so each ticket takes its own session. FA1 survives.
- A classifier refusal of an agent edit to an imported instruction file — the reviewer makes that edit or grants the permission rule. FA1 survives.
- Retired specs, retros, and reviews that name a retained author — they are history, and no reader treats them as current. FA22 survives.
- The phrase "retained integration source" — the one integration worktree stays, so its anchors keep their bytes. FA1 survives.

## Ownership fences

- `reviews/fresh-ticket-authors.md`
- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `.agents/commands/bench-implement-spec.md`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/charge_evidence_guidance_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `tests/canary/docs-currency-token-diet/`
- `tests/canary/load-validity-metadata/`
- `tests/canary/skills-index-command-adapters/`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `internal/conformance/implementation_continuation_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `tests/canary/claude-agent-definitions/`
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-final-check.md`
- `README.md`
- `docs/field-guide.html`
- `docs/adr/0021-benchmark-workflow-orchestration.md`
- `docs/adr/0023-each-ticket-gets-a-fresh-author.md`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `projects/benchkit.md`
- `CONTEXT.md`
- `tests/canary/guidance-prose-budgets/`
- `tests/canary/line-routing/`
- `tests/canary/skill-description-budgets/`

Reviewer disposition: pending sign-off.

## Out of scope

- The small-output evidence command: `bench preflight evidence` prints a summary by default and spills each source on request. It crosses ADR 0022, so the CLI research row owns it — about 8 edits and 3 gate runs.
- A serial execution mode in `internal/reviewrecord` — about 5 edits and 2 gate runs. The delegate plan with an author limit of 1 makes it unnecessary today.

## Further notes

### Reviewer decisions

The reviewer closed these decisions on 2026-09-24:

1. Every spec-backed build gives each ticket a fresh author session on the declared line. The version 2 delegate plan with an author limit of 1 records the authors. A tier move still asks the reviewer.
2. An explicit `--delegate` adds concurrent authors and the full tier range.
3. The rule covers every spec-backed build, with or without `--full`.
4. The author charge is narrow, and no author retrieves every evidence page.
5. The orchestrator keeps a small context and refreshes `bench handoff` at each chunk checkpoint.
6. A post-review repair goes to a fresh session recorded as a new assignment with the trigger `user-directed`.
7. `internal/reviewrecord` stays unchanged.
8. The change supersedes ADR 0021's single-author decision by promotion.
9. For ticket 1, the reviewer grants a one-time permission rule for `.bench/BENCH.md` and removes it after the ticket commits.
10. The repair record uses the existing `user-directed` trigger with its stopped-writer evidence, and this spec adds no new trigger.

### Bootstrap

This spec's own build starts before its guidance lands. So the recommended build already uses the new shape under today's rules: `$bench-implement-spec --full --delegate specs/fresh-ticket-authors/spec.md` with an author limit of 1 on opus/high, and no tier escalation.

The staged completion plan below stays at version 1, because a version 2 plan needs the run id and the orchestrator session. Before the first dispatch, the orchestrator amends that plan to the version 2 shape:

- The plan version becomes 2. The `execution` block names mode `delegate`, an author limit of 1, the run id, and the orchestrator session. It also names an empty history for each ticket.
- Chunk FA's verification splits into one pair for each ticket. `1-docs` and `1-conformance` name `1-route-each-ticket-to-a-fresh-author.md`.
- `2-docs` and `2-conformance` name `2-align-the-line-and-delegation-skills.md`.
- `3-docs` and `3-conformance` name `3-align-the-phase-commands-and-docs.md`.
- Each `-docs` obligation runs `bench test --check docs-currency-workflow`, and each `-conformance` obligation runs `bench test --package ./internal/conformance --run TestRootConformance`.
- The final verification names no ticket, and it stays as the staged plan writes it.

### Reader sweep

- The sweep read every sentence that names a retained author or session outside `specs/`, `capture/`, `roadmap/`, and `reviews/` at `5f1453fc`. Each one takes a row, a fence entry, or a Won't handle line above.
- `projects/benchkit.md` and `.bench/BENCH-reference.md` name "one retained integration source". That source stays, so their anchors keep their bytes.
- The three named canaries `delegated-repair-ownership`, `delegated-per-ticket-author`, and `prepared-build-approval` mutate sentences that this spec changes. `delegated-opt-in-entry` mutates the `--delegate` paragraph. The fence holds the whole canary family directory.
- The sweep read `internal/reviewrecord/delegated.go`. Version 2 requires mode `delegate`, an author limit of 1 or more, and one session for each ticket. Each replacement needs a trigger and stopped-writer evidence.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"fresh-ticket-authors-full-20260924","orchestrator_session":"claude:session_01PzPVd5kMaFqKt7bLjSjtgN","author_limit":1,"assignments":{"1-route-each-ticket-to-a-fresh-author.md":[{"session":"claude:bench-writer/fta-t1-author","assignment":"fta-t1-author","model":"opus","effort":"high","source":"d23694e925c6e0ffd15d6eda355daef4497acbea","native_ref":"claude:agent/fta-t1-author-20260924@d23694e925c6e0ffd15d6eda355daef4497acbea"},{"session":"claude:bench-writer/fta-t1-repair","assignment":"fta-t1-repair","model":"opus","effort":"high","source":"879e757db6a177c04dcaf4ec7c4acebf0c6e9efb","native_ref":"claude:agent/fta-t1-repair-20260924@879e757db6a177c04dcaf4ec7c4acebf0c6e9efb","predecessor":"claude:bench-writer/fta-t1-author","trigger":"user-directed","stopped":"claude:agent/fta-t1-author-20260924@d23694e925c6e0ffd15d6eda355daef4497acbea returned its final report and holds no live tool call","preserved":"879e757db6a177c04dcaf4ec7c4acebf0c6e9efb"}],"2-align-the-line-and-delegation-skills.md":[{"session":"claude:bench-writer/fta-t2-author","assignment":"fta-t2-author","model":"opus","effort":"high","source":"6300a0a26a1bf6f2d8cc2ce3a8b89a3363b909c7","native_ref":"claude:agent/fta-t2-author-20260924@6300a0a26a1bf6f2d8cc2ce3a8b89a3363b909c7"},{"session":"claude:bench-writer/fta-t2-repair","assignment":"fta-t2-repair","model":"opus","effort":"high","source":"3308f8b971d88cc12c0f3e0bc143a63926d3ff17","native_ref":"claude:agent/fta-t2-repair-20260924@3308f8b971d88cc12c0f3e0bc143a63926d3ff17","predecessor":"claude:bench-writer/fta-t2-author","trigger":"user-directed","stopped":"claude:agent/fta-t2-author-20260924@6300a0a26a1bf6f2d8cc2ce3a8b89a3363b909c7 returned its final report and holds no live tool call","preserved":"3308f8b971d88cc12c0f3e0bc143a63926d3ff17"}],"3-align-the-phase-commands-and-docs.md":[{"session":"claude:bench-writer/fta-t3-author","assignment":"fta-t3-author","model":"opus","effort":"high","source":"efd372f138885da525cfe4efa9a6ba7a06523ae2","native_ref":"claude:agent/fta-t3-author-20260924@efd372f138885da525cfe4efa9a6ba7a06523ae2"}]}},"chunks":[{"id":"FA","tickets":["1-route-each-ticket-to-a-fresh-author.md","2-align-the-line-and-delegation-skills.md","3-align-the-phase-commands-and-docs.md"],"verification":[{"id":"1-docs","command":"bench test --check docs-currency-workflow","ticket":"1-route-each-ticket-to-a-fresh-author.md"},{"id":"1-conformance","command":"bench test --package ./internal/conformance --run TestRootConformance","ticket":"1-route-each-ticket-to-a-fresh-author.md"},{"id":"2-docs","command":"bench test --check docs-currency-workflow","ticket":"2-align-the-line-and-delegation-skills.md"},{"id":"2-conformance","command":"bench test --package ./internal/conformance --run TestRootConformance","ticket":"2-align-the-line-and-delegation-skills.md"},{"id":"3-docs","command":"bench test --check docs-currency-workflow","ticket":"3-align-the-phase-commands-and-docs.md"},{"id":"3-conformance","command":"bench test --package ./internal/conformance --run TestRootConformance","ticket":"3-align-the-phase-commands-and-docs.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/fresh-ticket-authors/spec.md"},{"id":"docs","command":"bench test --check docs-currency-workflow"},{"id":"conformance","command":"bench test --package ./internal/conformance --run TestRootConformance"}]}
```
