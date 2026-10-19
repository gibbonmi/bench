# Grade each plan commit before the next dispatch

Status: staged

Decision source: named reviewed artifact `roadmap/FT380.md` with the five approved `FT380-plan-commit` criteria in `.bench/commitment.json` (reviewer decision 2026-10-18, commit 66a3da81d), plus the reviewer-confirmed late answer of 2026-10-19 that splits the plan-commit rule at the chunk checkpoint

Audience: this repository, and every repository that links the kit, because the guidance and the build preflight ship in the kit payload

Verification log: 3 iteration(s) to accept — the spec round took 2 iterations: Opus/high found a range-gap row, an existing-test seam that could not red, and a system fixture outside the PC-C1 plan; the ticket round accepted in 1

## Problem

A plan commit changes the coverage map, a ticket fence, or the completion plan during a build.
Today no guidance tells the orchestrator to grade a plan commit before the next dispatch.
Plan commits after a chunk freeze caused refreezes, stale reviews, and chain-gap refusals in many builds (see the occurrence ledger of `roadmap/FT380.md`).

The roadmap premise "the Markdown lane does not run the coverage parser" is stale.
Since commit `5288e0273` (2026-10-07), the `markdown` lane class selects `docs-currency-workflow` (`internal/gate/lane_select.go:145-149`).
That check runs `coverage --check` on every `specs/*/spec.md` (`checkCoverageMaps`, `internal/conformance/docs_workflow_checks_test.go:331-358`).
The 2026-10-10 build-directions occurrence, where the lane refused a citation to a later test, confirms that the lane runs the parser.

Three gaps remain.
First, no build preflight row compares the `## Implementation chunks` table with the coverage map, so a plan amendment can drop a row from every chunk (the FT376 GT-C1 occurrence).
Second, a red `check-floor` row refuses the review charge in code (`internal/preflight/review.go:207-212`), but no test shows it.
Third, the guidance still says that comment-only corrections follow the chunk tip (`.agents/commands/bench-implement-spec.md:60`).
The reviewer decided on 2026-10-18 that any source change after a chunk freeze moves the chunk tip and needs a refreeze.

## Solution

The orchestrator commits each plan commit with `bench commit --preflight-build <slug>`, and it dispatches nothing until that chain reads green.
The build preflight gains one `chunk-rows` row.
That row compares the chunk table with the coverage map in both directions.
A red `check-floor` row refuses the review charge as well as the build charge, and a test shows each refusal.

The implement and review phases state the plan-commit rules.
Before a chunk's green checkpoint, each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze.
After the green checkpoint, a plan commit lands in the chunk gap and joins the successor's review delta.
The last chunk has no successor, so each source change after its freeze, except a record commit or a recorded fold, takes a refreeze.

The delivery closes no roadmap row.
The check before publication, the commit-lane floor, and lane enforcement of the parser stay open on FT380.

## User stories

Line: `opus` / high.
Implementation-line reason: PC-C1 adds a table reader and a ticket-gated preflight row whose fixture blast radius crosses several packages; the spec fixes the predicate exactly and each row has a red-capable package test, but the census of spec fixtures is a known seam risk; PC-C2 is kit guidance, so the leverage override of `craft-line` sets high effort.
Harder chunks: PC-C1.

### The chunk table matches the coverage map

1. As an orchestrator, I want the build preflight to red a coverage row that no chunk-table row lists, so that a plan amendment cannot drop a row from every chunk.
2. As an orchestrator, I want the build preflight to red an acceptance-cell ID that the coverage map does not declare, so that a chunk cannot claim a deleted or misspelled row.
3. As a spec author, I want a hyphen range and an en dash range to expand to the declared rows between their endpoints, so that the staged tables need no rewrite.
4. As a spec author, I want a malformed range to give a fault that names it, so that a bad range cannot hide a row.
5. As a spec author, I want a token that is neither a row ID nor a range to give a fault that names it, so that free text cannot pass as coverage.
6. As a spec author, I want a row ID that two chunk rows list to stay green, so that a cross-chunk row keeps its current form.
7. As a spec author, I want the reader to find the acceptance column by its name in any letter case and at any position, so that the staged header variants still read.
8. As an orchestrator, I want a spec with no chunk table, or with no acceptance column, to red the row, so that an absent table is not a silent pass.
9. As a ticket author, I want the row to stay not-applicable before the tickets folder exists, so that the spec stage keeps the current ticket gate.
10. As an orchestrator, I want a red `chunk-rows` row to refuse the build charge, so that no author starts on a plan whose table and map disagree.
11. As a maintainer, I want the chunk-table reader and the coverage-map reader to split pipe cells through one function, so that two table readers cannot drift.
12. As a maintainer, I want the shared preflight spec fixture to render its chunk table from the rows that it declares, so that the green fixtures stay green.
13. As a maintainer, I want every ticketed staged spec to stay green on the new row, so that no queued build stops at its entry preflight.

### A red check floor blocks dispatch, not publication

14. As an orchestrator, I want a red `check-floor` row to refuse the review charge, so that no review starts on a plan that omits a floor check.
15. As an orchestrator, I want a red `check-floor` row to refuse the build charge, so that no author starts on a plan that omits a floor check.
16. As an orchestrator, I want a red preflight after a published plan commit to exit 4 and name the commit, so that the floor blocks dispatch and not publication.

### The phases state the plan-commit rules

17. As an orchestrator, I want the implement phase to require `bench commit --preflight-build <slug>` for each plan commit before the next dispatch, so that the build preflight grades the plan before an author reads it.
18. As an orchestrator, I want the implement phase to state that a source change before the green checkpoint moves the chunk tip and takes a refreeze, so that I do not spend a review on a stale pair.
19. As an orchestrator, I want the implement phase to state that a plan commit after the green checkpoint joins the successor's review delta, so that a successor-only plan commit waits for that checkpoint.
20. As an orchestrator, I want the implement phase to state that a source change after the last chunk's freeze takes a refreeze, so that the completion checkpoint grades a reviewed pair.
21. As a review axis, I want the review phase to state that the successor base is the last recorded source tip of its predecessor, so that no axis grades a range from a record commit.
22. As an orchestrator, I want the implement phase to state that a coverage-cell citation and its cited test land in one commit, so that the lane does not refuse the plan commit.
23. As an orchestrator, I want the implement phase to batch a review's prose-only spec corrections before the repair's refreeze, so that the repair session verifies once.
24. As an orchestrator, I want the implement phase to state that a red `check-floor` row blocks the next charge and the review charge, so that I repair the plan before dispatch.
25. As an orchestrator, I want the chain-gap refusal to state the current rule, so that the refusal does not advise the retired comment-only exception.
26. As an orchestrator, I want the repair policy to state that a comment-only correction before the green checkpoint still takes a refreeze, so that the evidence-only rule and the refreeze rule agree.
27. As a spec author, I want the completion-plan reference to name the `chunk-rows` row, so that I know that the build preflight grades the table.
28. As a maintainer, I want the implement phase growth to stay inside a recorded profile budget, so that the growth is a decision and not drift.
29. As the reviewer, I want the delivery to close no roadmap row, so that the check before publication and the commit-lane floor stay open on FT380.
30. As an orchestrator, I want the recorded-fold rule of the Land section to stay, so that a `main` fold in a chunk gap needs no refreeze.

## Implementation decisions

### The `chunk-rows` row

The coverage package owns the chunk-table reader, because it owns the coverage-map parse and its pipe-cell splitter (inside `processMapLine`, `internal/coverage/rows.go:178-205`, at lines 183-195).
The ticket extracts that splitter into one function, and the coverage-map reader and the chunk-table reader both call it.
The repository has no `internal/markdown` package today, so no production block reader exists to compose.
The test-only chunk-table shape check at `internal/conformance/retained_workflow_test.go:135-179` stays a fixture check; the staged `markdown-block-reader` spec can later move both readers onto its block reader.
The reader reads the first pipe table under the `## Implementation chunks` heading, before the next `##` or `###` heading.
It finds the acceptance column by the lowercased header cell `acceptance rows`, at any column position.

The reader splits each acceptance cell on commas.
It trims spaces and backticks from each token, and it skips an empty token, so a trailing comma adds no fault.
An empty cell lists no row and gives no fault; the unlisted-row check then reds each row that only that cell could list.
A token that matches the row-ID grammar `tickets.RowIDPattern` names one row.
A token `<A>-<B>` or `<A>–<B>` is a range when both endpoints match that grammar with one tag.
A range names its two endpoints and each declared row with that tag whose number lies between the endpoint numbers.
A range with an undeclared endpoint, two tags, or a first number above its last number is a fault that names the token.
Any other token is a fault that names the token.

The gatherer already holds the spec bytes that it resolved (`content`, `internal/preflight/gather.go:37`), and it passes those bytes to the chunk-table reader.
`coverage.ParseSpec` keeps its signature and its own read (`internal/coverage/coverage.go:262-268`), so its callers do not change.
`Facts` carries the chunk-table row IDs and the chunk-table faults.
`Decide` appends `ticketRow(f, "chunk-rows", chunkRowsCheck)` directly after `rows-membership`.
The row is ticket-gated through the existing `ticketRow` gate, so build mode without a tickets folder reads it not-applicable.
The `ticketRow` gate applies only to build mode (`internal/preflight/decision.go:267-271`), so review mode and the current-evidence check grade the row too, and a red row also refuses the review charge.
The row is red when the table is absent, when the acceptance column is absent, when the reader gives a fault, when a declared row is unlisted, or when a listed ID is undeclared.
An absent table gives the detail `no Implementation chunks table`, and an absent column gives the detail `chunk table has no acceptance rows column`.
Otherwise the detail names each fault, each unlisted row, and each undeclared ID.

Every red row already refuses the build and review charges through `verdict.Red` (`internal/preflight/charge.go:16-22`, `internal/preflight/review.go:207-212`), so the new row needs no charge code.
The new row does not join the `completion-plan` and `check-floor` set that `.agents/commands/bench-implement-spec.md:19` routes to the staleness pass.
The census below shows every ticketed staged spec green, so a red at build entry is a plan defect that stops the phase.

### The review-charge refusal

The review-charge refusal on a red `check-floor` row exists in code.
The ticket adds the missing test beside `TestCheckFloorRefusesTheBuildCharge`, through the same `seedFloorRepo` fixture.

### The guidance

The implement phase Land section states the plan-commit rules in one new paragraph after the serial-commit paragraph.
It replaces the clause "Only record commits, recorded folds, and comment-only corrections follow the chunk tip." with the moved-tip rule.
The new text avoids the Forbid needles "review delta of the last chunk" and "the last chunk refreezes" (`internal/anchors/registry_chunk_chain.go`).
The guidance ticket adds one Forbid row to `internal/anchors/registry_chunk_chain.go` for the implement phase with the exact needle `comment-only corrections follow the chunk tip`.
A plan commit in a chunk gap can meet a moved default tip, because the chunk gap admits other landings (`.bench/BENCH-reference.md`, Command Notes).
The new paragraph therefore tells the orchestrator to run the chunk-gap fold of the Land section before that plan commit, so the `base-current` row of the chained preflight stays green.
A plan commit that changes the completion plan still takes the review-record amendment that `.bench/BENCH.md` already names for an enabling plan commit.
The recorded-fold sentence of the Land section stays, because a fold brings default-branch bytes that the fold record and its resolution review cover; a fold is not a build source change.

The review phase step 1 states that the successor base is the last recorded source tip of its predecessor, never a record commit.
That sentence rewrites the base sentence at `.agents/commands/bench-review-implementation.md:97`, which keeps its fold clause: the base is the last recorded fold after that tip, or that tip when no fold follows it.
It also states that a plan commit after the predecessor's green checkpoint sits in the successor's delta.
The repair policy keeps its evidence-only rule for a comment-only correction and adds that such a correction before the green checkpoint still takes a refreeze.
The completion-plan reference names the `chunk-rows` row beside the `check-floor` row.

The chain-gap refusal at `internal/reviewrecord/coverage.go:102` states the retired rule.
The ticket replaces the clause `only record commits, recorded folds, and comment-only corrections follow a chunk tip` with the exact clause `only record commits and recorded folds follow a chunk tip until its checkpoint`.
The one exact-text test `TestReviewCheckpointChainGapNamesTheExpectedBase` moves with it.
The code that admits a comment-only gap at the checkpoint stays unchanged; the reviewer answer of 2026-10-19 asks for no code change to the chain.
So the guidance is stricter than the checkpoint code for a comment-only correction, and the flagged additions record that gap.

The implement phase sits at its budget of 93 lines (`projects/benchkit.md:498`).
The new paragraph is one physical line plus one blank separator line, so it adds two lines.
The guidance ticket raises the budget row to 95 in the same commit and states the reason in its commit message: the plan-commit rules are the phase's own build rules.

`bench commit --preflight-build` passes `--source-tip` and no `--base` (`cmd/bench/commit_chain.go:58`), so its `base-current` row reads the implicit range against the default tip.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| PC-C1 / 01-grade-the-chunk-table-against-the-map.md, 02-refuse-the-review-charge-on-a-red-floor.md | The build preflight grades the chunk table against the coverage map, and a red floor refuses both charges | PC1–PC21, PC37, PC38 | coverage, preflight, gate, and diff package tests, the command tests, the system suite, and root conformance | yes |
| PC-C2 / 03-state-the-plan-commit-rules.md | The implement and review phases, the repair policy, the completion-plan reference, and the chain-gap refusal state the plan-commit rules | PC22–PC36, PC39 | anchors, reviewrecord, gate, and command package tests, `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, and root conformance | no |

PC-C2 follows PC-C1, because the completion-plan reference names the row that PC-C1 adds.
Each chunk is one behavior outcome with one review checkpoint.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft380-plan-commit-full-20261019","orchestrator_session":"claude:ft380-orchestrator-20261019","author_limit":1,"assignments":{"01-grade-the-chunk-table-against-the-map.md":[{"session":"claude:ft380-t01-author-20261019","assignment":"41d892992aff31a551f2052e3cae8f6d","model":"opus","effort":"high","source":"8e6a8b2c94ed3b53bd376d1d7d2e36500c1efb53","native_ref":"claude-agent:bench-writer-high:ft380-t01-author-20261019"}],"02-refuse-the-review-charge-on-a-red-floor.md":[{"session":"claude:ft380-t02-author-20261019","assignment":"41d892992aff31a551f2052e3cae8f6d","model":"opus","effort":"high","source":"2e97bcf1590e16d425a46168b0455a416e32b98c","native_ref":"claude-agent:bench-writer-high:ft380-t02-author-20261019"},{"session":"claude:ft380-t02-repair1-20261019","assignment":"41d892992aff31a551f2052e3cae8f6d","model":"opus","effort":"high","source":"148b02da29b37fe4949e984e4ab4d15ba84455f7","native_ref":"claude-agent:bench-writer-high:ft380-t02-repair1-20261019","predecessor":"claude:ft380-t02-author-20261019","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:ft380-t02-author-20261019","preserved":"5daa190037d8cc2503ff1e3d01802c9e7cf886f5"}],"03-state-the-plan-commit-rules.md":[]}},"chunks":[{"id":"PC-C1","tickets":["01-grade-the-chunk-table-against-the-map.md","02-refuse-the-review-charge-on-a-red-floor.md"],"verification":[{"id":"t1-coverage-package","command":"bench test --package ./internal/coverage","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-preflight-packages","command":"bench test --package ./internal/preflight/...","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-command-package","command":"bench test --package ./cmd/bench","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-gate-package","command":"bench test --package ./internal/gate","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-diff-package","command":"bench test --package ./internal/diff","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-recordtest-package","command":"bench test --package ./internal/reviewrecord/recordtest","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-system","command":"bench test --check system","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-root-conformance","command":"bench test --check root-conformance","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t1-unlisted-row-omission","command":"bench test --package ./internal/preflight","probe":"drop the declared-row-to-cell direction in chunkRowsCheck; TestChunkRowsRedsAnUnlistedRow must fail","ticket":"01-grade-the-chunk-table-against-the-map.md"},{"id":"t2-preflight-packages","command":"bench test --package ./internal/preflight/...","ticket":"02-refuse-the-review-charge-on-a-red-floor.md"},{"id":"t2-review-charge-omission","command":"bench test --package ./internal/preflight","probe":"remove the verdict.Red refusal of the review charge; TestCheckFloorRefusesTheReviewCharge must fail","ticket":"02-refuse-the-review-charge-on-a-red-floor.md"},{"id":"t2-root-conformance","command":"bench test --check root-conformance","ticket":"02-refuse-the-review-charge-on-a-red-floor.md"}]},{"id":"PC-C2","tickets":["03-state-the-plan-commit-rules.md"],"verification":[{"id":"t3-anchors-package","command":"bench test --package ./internal/anchors","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-reviewrecord-package","command":"bench test --package ./internal/reviewrecord/...","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-gate-package","command":"bench test --package ./internal/gate","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-command-package","command":"bench test --package ./cmd/bench","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-guidance","command":"bench test --check docs-currency-workflow","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-budgets","command":"bench test --check guidance-prose-budgets","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-prose","command":"bench test --check prose-mechanics","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-root-conformance","command":"bench test --check root-conformance","ticket":"03-state-the-plan-commit-rules.md"},{"id":"t3-forbid-clause-restore","command":"bench test --check docs-currency-workflow","probe":"restore the clause comment-only corrections follow the chunk tip in the implement phase; the new Forbid row must fail","ticket":"03-state-the-plan-commit-rules.md"}]}],"final_verification":[{"id":"guidance","command":"bench test --check docs-currency-workflow"},{"id":"root-conformance","command":"bench test --check root-conformance"}]}
```

## Testing decisions

- A good test drives the build preflight verdict through `Decide` or through the `bench preflight` command, and it reads the `chunk-rows` row and its detail.
- The reader tests drive the coverage package entry with exact spec bytes and read the row IDs and the faults.
- The charge tests run the real command with charge arguments over the seeded fixture repository and read the refusal line.
- Prior art: `TestCheckFloorRefusesTheBuildCharge` for a charge refusal, `TestDecideAllGreen` for the row order, and the `rows-membership` cases of `decision_test.go` for a row detail.
- Guidance rows are review-owned, except the Forbid anchor row and the budget row, which the docs and budget checks grade.
- The gate observes the feature through the preflight and coverage package tests, `docs-currency-workflow`, and `guidance-prose-budgets`.

### Seam diagram

    trigger: bench commit --preflight-build <slug>, or bench preflight build <slug> [--charge]
        │
        ▼
    spec bytes ──▶ [ coverage map reader + chunk-table reader (one cell splitter) ] ──▶ declared IDs, listed IDs, faults
                                          │
                                          ▼
                              [ Decide: chunk-rows after rows-membership ] ──▶ verdict row ──▶ charge refusal on red
                      ◀ tests attach here: coverage reader tests with exact bytes, Decide facts, and the real command

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PC1 | 1 | A declared row that no acceptance cell lists reds `chunk-rows`, and the detail names that row | `internal/preflight/chunk_rows_test.go` (`TestChunkRowsRedsAnUnlistedRow`) | A check of only the cell-to-map direction stays green on this fixture |
| PC2 | 2 | An acceptance-cell ID that the map does not declare reds `chunk-rows`, and the detail names that ID | `internal/preflight/chunk_rows_test.go` (`TestChunkRowsRedsAnUndeclaredID`) | A check of only the map-to-cell direction stays green on this fixture |
| PC3 | 3 | The reader expands the hyphen range `PF2-PF4` to PF2, PF3, and PF4 when the map declares all three | `internal/coverage/chunk_table_test.go` (`TestChunkTableExpandsRanges`) | A reader that keeps only the endpoints leaves PF3 unlisted |
| PC4 | 3 | The reader expands the en dash range `PF2–PF4` to PF2, PF3, and PF4 | `internal/coverage/chunk_table_test.go` (`TestChunkTableExpandsRanges`) | A reader that splits only on the hyphen reads the en dash range as one bad token |
| PC5 | 4 | A range whose endpoint the map does not declare gives a fault that names that endpoint | `internal/coverage/chunk_table_test.go` (`TestChunkTableRangeFaults`) | A reader that expands against the declared set alone hides the missing endpoint |
| PC6 | 4 | A range whose two endpoints carry different tags gives a fault that names the token | `internal/coverage/chunk_table_test.go` (`TestChunkTableRangeFaults`) | A reader that compares only the numbers accepts a range across two specs' tags |
| PC7 | 4 | A range whose first number is above its last number gives a fault that names the token | `internal/coverage/chunk_table_test.go` (`TestChunkTableRangeFaults`) | A reader that expands an empty interval lists only the endpoints and hides the error |
| PC8 | 5 | An acceptance-cell token that is neither a row ID nor a range gives a fault that names the token | `internal/coverage/chunk_table_test.go` (`TestChunkTableTokenFaults`) | A reader that skips unknown tokens lets free text pass as coverage |
| PC9 | 6 | A row ID that two chunk rows list reads `chunk-rows` green | `internal/preflight/chunk_rows_test.go` (`TestChunkRowsAcceptsARepeatedRow`) | A check that refuses duplicates reds the cross-chunk rows of staged specs |
| PC10 | 7 | A header `Acceptance rows` at the fifth position after a `blocked by` column gives the acceptance cells | `internal/coverage/chunk_table_test.go` (`TestChunkTableFindsTheAcceptanceColumn`) | A reader with a fixed column index or a case-sensitive header reads the wrong cell |
| PC11 | 8 | A spec with no chunk table reds `chunk-rows` with the detail `no Implementation chunks table` | `internal/preflight/chunk_rows_test.go` (`TestChunkRowsRedsAnAbsentTable`) | A row that treats an absent table as empty and green passes every spec |
| PC12 | 8 | A chunk table with no acceptance column reds `chunk-rows` with the detail `chunk table has no acceptance rows column` | `internal/preflight/chunk_rows_test.go` (`TestChunkRowsRedsAnAbsentTable`) | A row that reads a missing column as empty cells reds every row with a misleading detail |
| PC13 | 9 | Build mode without a tickets folder reads `chunk-rows` not-applicable with no detail | `internal/preflight/chunk_rows_test.go` (`TestChunkRowsNotApplicableWithoutTickets`) | A row that calls its check without the `ticketRow` gate reds a fresh spec before its tickets exist |
| PC14 | 9 | `Decide` lists `chunk-rows` directly after `rows-membership` in a green verdict | `internal/preflight/decision_test.go` (`TestDecideAllGreen`) | A missing or misplaced row breaks the pinned row order |
| PC15 | 10 | A red `chunk-rows` row refuses the build charge with `preflight required: chunk-rows` | `internal/preflight/chunk_rows_test.go` (`TestChunkRowsRefusesTheBuildCharge`) | A row outside `Decide` or marked exempt lets the charge pass |
| PC16 | 11 | The chunk-table reader and the coverage-map reader call one cell-split function | review-owned: the Standards axis compares the two readers in the frozen delta | A second splitter is a duplicated parse that the knowledge rule refuses |
| PC17 | 12 | The shared `specBody` fixture renders its chunk table from its declared rows, and the seeded build verdict reads `chunk-rows` green | `internal/preflight/command_build_test.go` (`TestCommandBuildSeedsAGreenChunkTable`) | A fixture without a table, or a table that keeps a row ID that a test rewrites, reds the seeded build verdict |
| PC18 | 13 | Each ticketed staged spec reads `chunk-rows` green in its plan-only preflight at the ticket tip | review-owned: the ticket author records one plan-only verdict for each ticketed staged spec | A stricter grammar than the staged tables use stops a queued build at entry |
| PC19 | 14 | A red `check-floor` row refuses the review charge with `preflight required: check-floor` | `internal/preflight/check_floor_command_test.go` (`TestCheckFloorRefusesTheReviewCharge`) | A review charge that skips the verdict lets a review start on a plan without its floor |
| PC20 | 15 | A red `check-floor` row refuses the build charge with `preflight required: check-floor` | `internal/preflight/check_floor_command_test.go` (`TestCheckFloorRefusesTheBuildCharge`) | A build charge that skips the verdict lets an author start on a plan without its floor |
| PC21 | 16 | A red preflight after a published commit exits 4 and the chain line names the commit | `cmd/bench/commit_chain_test.go` (`TestCommitChainReportsRedPreflight`) | A chain that returns exit 1 claims that nothing was committed |
| PC22 | 17 | The Land section requires `bench commit --preflight-build <slug>` for each plan commit before the next dispatch | review-owned: the Spec axis reads the Land section of the implement phase | A phase without the form lets a plan commit reach an author ungraded |
| PC23 | 18 | The Land section states that a source change before the green checkpoint, except a record commit or a recorded fold, moves the chunk tip and takes a refreeze | review-owned: the Spec axis reads the Land section of the implement phase | A phase that keeps the comment-only exception repeats the AP-C7 chain break |
| PC24 | 19 | The Land section states that a plan commit after the green checkpoint joins the successor's review delta and that a successor-only plan commit waits for that checkpoint | review-owned: the Spec axis reads the Land section of the implement phase | A phase without the split repeats the FT405 RT-CURRENT stall |
| PC25 | 20 | The Land section states that a source change after the last chunk's freeze, except a record commit or a recorded fold, takes a refreeze of that chunk | review-owned: the Spec axis reads the Land section of the implement phase | A phase without the rule repeats the BD-C4 completion refusal |
| PC26 | 21 | Review step 1 states that the successor base is the last recorded source tip of its predecessor, never a record commit | review-owned: the Spec axis reads step 1 of the review phase | A base at a record commit repeats the SS-C2 and SS-C3 chain gap |
| PC27 | 22 | The Land section states that a coverage-cell citation and its cited test land in one commit | review-owned: the Spec axis reads the Land section of the implement phase | A plan commit that cites a later test meets a lane refusal |
| PC28 | 23 | The Land section batches a review's prose-only spec corrections before the repair's refreeze, or defers them to the next plan commit before a ticket dispatch | review-owned: the Spec axis reads the Land section of the implement phase | A late prose commit stales the repair verification, as in SB1 |
| PC29 | 24 | The Land section states that a red `check-floor` row blocks the next charge and the review charge | review-owned: the Spec axis reads the Land section of the implement phase | An orchestrator without the rule retries the charge instead of the plan, as in FT379 |
| PC30 | 18 | The implement phase no longer contains `comment-only corrections follow the chunk tip` | Forbid row in `internal/anchors/registry_chunk_chain.go`, graded by `bench test --check docs-currency-workflow` | A restored clause reds the docs-currency check |
| PC31 | 25 | The chain-gap refusal contains the clause `only record commits and recorded folds follow a chunk tip until its checkpoint` | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointChainGapNamesTheExpectedBase`) | The old refusal text fails the updated exact-text assertion |
| PC32 | 26 | The repair policy states that a comment-only correction before the green checkpoint takes a refreeze | review-owned: the Spec axis reads the bounded repair policy | A policy that keeps only the evidence-only sentence contradicts the moved-tip rule |
| PC33 | 27 | The completion-plan reference names the `chunk-rows` row and its two directions | review-owned: the Spec axis reads the completion-plan reference | A reference without the row leaves the table grader unnamed |
| PC34 | 28 | The implement phase line count stays inside its profile budget row of 95 lines | `bench test --check guidance-prose-budgets` grades the phase against the profile table | A paragraph added without the budget edit reds the budget check |
| PC35 | 29 | The spec carries no `Roadmap:` line, and `roadmap/FT380.md` keeps its open clauses after the landing | review-owned: the Spec axis reads the spec header and the roadmap detail file | A `Roadmap:` line lets the spec retirement delete the open FT380 row |
| PC36 | 30 | The Land section keeps the recorded-fold sentence that begins `A build folds` | the existing Require row in `internal/anchors/registry_chunk_chain.go`, graded by `bench test --check docs-currency-workflow` | A rewrite that drops the fold sentence reds the docs-currency check |
| PC37 | 3 | The range `PF2–PF4` over a map that declares PF2 and PF4 and omits PF3 lists exactly PF2 and PF4, with no fault | `internal/coverage/chunk_table_test.go` (`TestChunkTableExpandsRanges`) | A reader that lists every number between the endpoints names the undeclared PF3 |
| PC38 | 5 | The cell `` `PF1`, PF2, `` lists PF1 and PF2 with no fault | `internal/coverage/chunk_table_test.go` (`TestChunkTableTokenFaults`) | A reader that keeps backticks or faults an empty token reds a staged cell form |
| PC39 | 19 | The Land section states that a plan commit in the chunk gap comes after the chunk-gap fold of `main` | review-owned: the Spec axis reads the Land section of the implement phase | A plan commit before the fold meets a moved default tip and reds `base-current` in the chained preflight |

### Edge inventory

The edge walk covers the absent-versus-empty pair for the chunk table: an absent section (PC11), an absent column (PC12), and an empty acceptance cell, which lists no row and so reds through PC1.
An empty token and a backticked token read as the staged cells use them (PC38).
Each refusal of the reader names its token (PC5 to PC8).
A repeated row stays green (PC9).
No test in this spec swaps a package variable.

- **Won't handle:** the `--preflight-build` chain passes no `--base`, so a reviewer landing override that moves `main` inside a chunk reds `base-current` — the post-ticket form `bench worktree exec <target> -- bench preflight build <slug> --base <chunk-base> --source-tip <tip>` stays the in-scope caller for that case.
- **Won't handle:** a linked repository's staged spec with no chunk table reds `chunk-rows` at build entry and stops the phase — the `craft-spec` template already requires the chunk table, so the spec author who adds the table is the in-scope caller.
- **Won't handle:** a chunk-table acceptance cell that names a row of another spec's tag — the reader reds it as undeclared, and no staged spec does it.
- **Won't handle:** a check that the chunk-table chunk IDs and tickets equal the completion-plan fence — `completion-plan` and `rows-owned` keep their current grading, and the criteria do not ask for it.
- **Won't handle:** a code admission of plan commits after a chunk tip — the reviewer answer of 2026-10-19 keeps the chain code unchanged, and the guidance routes such a commit to a refreeze or to the chunk gap.

## Out of scope

- The check before publication of a plan commit, which FT380 keeps open: 6 edits, 3 gate runs.
- The `check-floor` floor in the commit lane, which FT380 keeps open: 5 edits, 3 gate runs.
- Lane enforcement of the skipped parser beyond the current `docs-currency-workflow` selection, which graduates only on a `bench commit` repro: 4 edits, 2 gate runs.

## Further notes

### Census of staged specs

A read-only census on 2026-10-19 applied the two-way predicate with the range rule above to each staged spec at commit 781706f19.
All 17 specs with tickets passed with no unlisted row, no undeclared ID, and no fault.
That census included `degraded-guard-refusal`, which retired before this build, and it did not include this spec.
At build entry (commit 7b45186c9), 16 other staged specs have tickets, and the PC18 plan-only verdicts of ticket 01 grade them and this spec.
`session-context-efficiency` has no tickets and no chunk table, so the ticket gate keeps the row not-applicable for it.

### Flagged additions

- The chain-gap refusal text rewrite (PC31). The criteria name guidance only, but the refusal text states the retired rule to every operator who meets it.
- The repair-policy refreeze sentence (PC32). It reconciles `bounded-repair-policy.md:47` with the moved-tip criterion.
- The completion-plan reference sentence (PC33). It names the new row where the `check-floor` row already sits.
- The implement-phase budget raise from 93 to 95 (PC34).
- The review-charge test (PC19). The behavior exists in code; the criterion "blocks ... the review charge" had no test.
- The fold-before-plan-commit sentence (PC39). It keeps the chained preflight green in a chunk gap, where other landings can move the default tip.
- The guidance is stricter than the checkpoint code for a comment-only correction: the checkpoint still admits a comment-only gap at a chunk tip (`TestReviewCheckpointCommentOnlyGap`), and the reviewer answer of 2026-10-19 keeps that code.

### Source sentence to row

| source | rows |
|---|---|
| FT380.md:3-6 plan commit changes the map, a fence, or the plan | PC22 |
| FT380.md:3-4 the Markdown lane does not run the parser | stale premise, Problem section |
| FT380.md:5-6 a `Writes:` expansion needs the fence union and the closure | retired by FT293, no row |
| FT380.md:8-10 `--preflight-build` for each plan commit before the next dispatch | PC21, PC22 |
| FT380.md:10 a check before publication stays open | Out of scope, PC35 |
| FT380.md:10 a plan commit after a freeze joins the successor's delta | PC24, PC26 |
| FT380.md:10-11 lane enforcement graduates only on a repro | Out of scope |
| FT380.md:13 the build preflight checks the chunk table against the map | PC1–PC18 |
| FT380.md:13-14 a red `check-floor` blocks the next charge and the review charge, not the publication | PC19, PC20, PC21, PC29 |
| FT380.md:14-15 the floor in the commit lane stays open | Out of scope |
| FT380.md:17 the last chunk forces a refreeze | PC25 |
| FT380.md:17 a citation and its test land in one commit | PC27 |
| FT380.md:17 any source change after a freeze moves the tip and needs a refreeze | PC23, PC30, PC31, PC32 |
| FT380.md:17 the successor base is the last recorded source tip, never a record commit | PC26 |
| FT380.md:19-22 batch prose-only spec corrections | PC28 |
| criterion `.preflight-build` | PC22 |
| criterion `.successor-delta` | PC24, PC26 |
| criterion `.chunk-table` | PC1–PC20, PC29 |
| criterion `.open` | PC35 |
| criterion `.moved-tip` | PC23, PC25, PC26, PC30 |
| reviewer answer 2026-10-19 | PC23, PC24, PC25, PC36, PC39 |

### Pre-review proof checklist

- `Cited symbols`: `processMapLine`, `coverage.ParseSpec`, `tickets.RowIDPattern`, `ticketRow`, `Decide`, `verdict.Red`, `seedFloorRepo`, `TestCheckFloorRefusesTheBuildCharge`, `TestDecideAllGreen`, `TestCommitChainReportsRedPreflight`, `TestReviewCheckpointChainGapNamesTheExpectedBase`, and `checkCoverageMaps` each resolve at commit 781706f19.
- `Import edges`: none new; `internal/preflight` already imports `internal/coverage`, and `internal/coverage` already imports `internal/tickets`.
- `Source-row clauses and occurrences`: the source-sentence table above quotes each clause; the 17 occurrences map to PC1 (GT-C1), PC21 and PC22 (the DC-C8 plan commit), PC23 and PC30 (TF-C3, D-C4, AP-C7), PC24 (the DC-C8 seam correction, RT-CURRENT, durable-caller-migration), PC25 (BD-C4), PC26 (BD-C3, SS-C2), PC27 (FT370, build-directions), PC28 (SB1), and PC29 (PS-D, FT379); DC-C9 belongs to the closure that FT293 removed.
- `Promised field labels`: the row name `chunk-rows` and the refusal prefix `preflight required: chunk-rows`.
- `Changed-function callers`: `processMapLine` has one caller, `parse`; `coverage.ParseSpec` keeps its signature; `Decide` has its existing callers through the verdict, charge, and current-evidence paths; the slicer reruns `bench consumers` for each.
- `Copy survival`: PC16 refuses a second splitter.
- `Rendered-shape readers`: the preflight row list is asserted in `decision_test.go`, `charge_cells_test.go`, `source_tip_test.go` (`bareRowCount`), `binary_seal_test.go` (`rowIndex`), `command_build_test.go`, `command_bootstrap_test.go`, `command_review_test.go`, and `staleness/completeness_test.go`; the chain-gap text is asserted only in `review_checkpoint_commits_test.go:49`.
- `Pin operators`: `TestDecideAllGreen` compares the full ordered list by equality; `bareRowCount` compares a count by equality; the chain-gap assertion uses `strings.Contains`.
- `Entry reads`: none; the reader takes spec bytes.
- `Derived expectations`: the PC3 and PC4 expansions come from the fixture's declared rows, not from the reader.
- `Consolidated rules`: none.
- `Quantified obligations`: "every declared row" and "every listed ID" take PC1 and PC2; "every ticketed staged spec" takes PC18.
- `Workflow-step writes`: the plan commit writes the spec folder, which the source digest covers; the guidance rows state the refreeze that follows.
- `Input constructors`: each spec fixture that reaches a ticket-gated verdict must render a chunk table; the shared builders are `specBody` in `internal/preflight/preflighttest/fixture.go` and the plan writer of `internal/reviewrecord/recordtest/fixture.go:80-88`, which appends a plan to a caller's spec body with no chunk table. One confirmed reach is `internal/systemtest/charge_evidence_test.go:61`, whose inline spec has no chunk table and whose build charge at line 93 expects exit 0, so ticket 01 writes that file and PC-C1 runs `bench test --check system`. The other inline spec bodies sit in `charge_edges_test.go`, `command_bootstrap_test.go`, `command_review_test.go`, `cmd/bench/command_registry_test.go`, `cmd/bench/preflight_version_test.go`, `internal/systemtest/charge_evidence_test.go`, `internal/systemtest/owner_landing_fixture_test.go`, `internal/worktree/delegated_integration_test.go`, `internal/worktree/land_fixtures_test.go`, and `internal/preflight/ownership/ownershiptest/repo.go`; slicing research confirms which of them reach the row.
