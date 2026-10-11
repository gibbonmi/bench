# Measure the guidance prose budgets in bytes

Blocked by: 5-move-the-batch-approval-rule-beside-the-sign-off-stop.md
Writes: .agents/skills/bench-craft-spec/SKILL.md, CHANGELOG.md, CONTEXT.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/checks_test.go, internal/conformance/prose_budget_test.go, internal/conformance/skill_description_budgets_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/tier_live_tree_test.go, projects/benchkit.md, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/skill-description-budgets/budget-table-missing, tests/canary/skill-description-budgets/description-folded, tests/canary/skill-description-budgets/description-missing, tests/canary/skill-description-budgets/over-budget-command, tests/canary/skill-description-budgets/over-budget-description, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route, tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening, tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary, tests/canary/workflow-guidance-anchors/craft-spec-empty-or-invalid-fence, tests/canary/workflow-guidance-anchors/craft-spec-exact-literal-fence, tests/canary/workflow-guidance-anchors/craft-spec-fence-after-slice, tests/canary/workflow-guidance-anchors/craft-spec-further-notes-bare, tests/canary/workflow-guidance-anchors/craft-spec-reader-sweep-name, tests/canary/workflow-guidance-anchors/craft-spec-red-signal-grammar-retired, tests/canary/workflow-guidance-anchors/craft-spec-template-verification-log, tests/canary/workflow-guidance-anchors/dg-10, tests/canary/workflow-guidance-anchors/dg-11, tests/canary/workflow-guidance-anchors/dg-12, tests/canary/workflow-guidance-anchors/dg-13, tests/canary/workflow-guidance-anchors/dg-14, tests/canary/workflow-guidance-anchors/dg-15, tests/canary/workflow-guidance-anchors/dg-15-executable-red, tests/canary/workflow-guidance-anchors/dg-15-mandatory-red, tests/canary/workflow-guidance-anchors/dg-15-mandatory-red-capitalized, tests/canary/workflow-guidance-anchors/dg-15-mandatory-red-emphasis, tests/canary/workflow-guidance-anchors/dg-9, tests/canary/workflow-guidance-anchors/ticket-cross-pointers-anchor, tests/canary/workflow-guidance-anchors/write-spec-acceptance-round-exit, tests/canary/workflow-guidance-anchors/write-spec-approval-table-emit, tests/canary/workflow-guidance-anchors/write-spec-contract-rederivation, tests/canary/workflow-guidance-anchors/write-spec-degenerate-cheapness, tests/canary/workflow-guidance-anchors/write-spec-fence-approval, tests/canary/workflow-guidance-anchors/write-spec-identified-coverage-default, tests/canary/workflow-guidance-anchors/write-spec-materiality-exit, tests/canary/workflow-guidance-anchors/write-spec-ownership-fences, tests/canary/workflow-guidance-anchors/write-spec-promise-inflation-guard, tests/canary/workflow-guidance-anchors/write-spec-unflagged-addition-removal, tests/canary/workflow-guidance-anchors/write-spec-unique-row-id
Covers: GW38, GW39, GW40, GW41, GW42, GW43, GW44, GW45, GW46, GW47, GW48, GW49, GW50, GW51, GW53

## What to build

The `guidance-prose-budgets` check grades each subject by its size in bytes, so a long line costs what it weighs.
Every row of `### Guidance prose budgets` in `projects/benchkit.md` limits its subject in bytes, and each row starts tight.
`AGENTS.md` gets its own row, and the `.bench/BENCH.md` row bounds the guide at its size after the moves.
Each red names its unit, and the over-budget canary bites in bytes.

This ticket reads four contracts from the earlier chunks:

- Tickets 1 to 5 moved, repointed, or removed text in `.bench/BENCH.md`, `AGENTS.md`, and the budgeted command and skill files. Each subject has its final size only after ticket 5.
- Ticket `5-move-the-batch-approval-rule-beside-the-sign-off-stop.md` raised the line row of `.agents/commands/bench-write-spec.md`. This ticket replaces that row with a byte row.
- Ticket `2-repoint-the-build-rule-pointers.md` changed only the first sentence of `.agents/skills/bench-craft-spec/SKILL.md:68`. This ticket changes only the third sentence of that line.
- Ticket `1-move-the-build-rules-into-their-own-file.md` added the section `### Always-loaded guidance weight` to `CHANGELOG.md`. Tickets 3 to 5 added the GW2 to GW4 bullets after the GW1 bullet.

Line citations below are at base `bc2ea0ff`.
If a landing moves a line, find the text by its quotation.

**The count.** In `checkGuidanceProseBudgets`, compare `len(subject.Data)` with the row limit by `>` (`internal/conformance/prose_budget_test.go:99`).
The diagnostic at `:100` becomes `prose-budget exceeded: %s is %d bytes, over its %d-byte budget`.
Remove `proseBudgetLineCount` (`:204-211`), whose only caller is `:99`, and remove the `bytes` import, which only that function uses.
The limit error at `:147` ends with "which is not a positive byte count".
The header stays `| subject | limit |`, because the parser requires it (`:132-135`).

The count reads the bytes that `bounds.ClassifyNoFollow` returns at `:81`, so each refusal at `:82-98` stays as it is.
Do not change the shape of `proseBudgetPolicy`, `limitFor`, `profileBudgetTable`, or `proseBudgetRows`.
`internal/conformance/skill_description_budgets_test.go:44`, `:102`, `:162-220` and `internal/conformance/checks_test.go:576` read them.
The byte count needs no sentence splitter, so this ticket does not call `internal/prose`.

**The test helpers.** Rename `proseBudgetLines` (`:230-241`) to `proseBudgetBytes(size int, trailing bool) string`, because a helper that renders bytes must not say lines.
It returns exactly `size` bytes on one physical line, and with `trailing` its last byte is the final newline.
The filler holds one line, so a line count stays at 1 and never reaches a fixture limit.
That is how each byte expectation reds a line count.

`healthyProseBudgetFiles` (`:255-261`) renders each subject one byte under its limit, at 149, 99, and 119 bytes.
`proseBudgetRows` (`:224-228`) keeps its values, which are now byte limits.
Each other call of the old helper keeps its number as a byte size: `:282`, `:333`, `:385`, `:428-429`, and `:500-502`.
Each fixture size comes from the bytes that the test writes, not from the check.

**The expectations.** Each expectation below names bytes:

| site | new expectation |
|---|---|
| `:276` | `prose-budget exceeded: .agents/skills/bench-craft-tickets/SKILL.md is 99 bytes, over its 90-byte budget` |
| `:284` | `over its 100-byte budget` |
| `:315` | the whole limit error, which ends with `, which is not a positive byte count` |
| `:335` | `prose-budget exceeded: .agents/skills/bench-craft-new/SKILL.md is 121 bytes, over its 120-byte budget` |
| `:506-508` | the three diagnostics at 151, 101, and 121 bytes, over their 150-, 100-, and 120-byte budgets |

The case at `:315` now names the unit, so a stale limit error reds it.
Correct the doc comments at `:230-231`, `:253-254`, `:349-350`, and `:516-518`, so that none of them says lines.
The doc comment at `:204-205` goes with its function.

**The boundary test.** `TestGuidanceProseBudgetBoundaryCountsEveryByte` replaces `TestGuidanceProseBudgetBoundaryIsNewlineIndifferent` at `:516-539`.
It grades the craft-tickets row of 100 bytes in four cases:

| case | subject | result |
|---|---|---|
| at the limit, no final newline | `proseBudgetBytes(100, false)` | no diagnostic (GW40) |
| at the limit, final newline | `proseBudgetBytes(100, true)` | no diagnostic (GW40, GW41) |
| one byte over by its final newline | `proseBudgetBytes(101, true)` | the diagnostic at 101 bytes, over its 100-byte budget (GW39, GW41) |
| one byte over, no final newline | `proseBudgetBytes(101, false)` | the diagnostic at 101 bytes, over its 100-byte budget (GW39) |

A count that trims the final newline passes the third case, and a strict comparison reds the second case.
No file outside `prose_budget_test.go` names the old test, so no inventory changes.
The new test reads only a throwaway root, so it does not join `classifiedLiveTreeTests` (`internal/conformance/tier_live_tree_test.go:19-72`).
`TestGuidanceProseBudgetsHoldOnTheLiveTree` and `TestGuidanceProseBudgetCanaryFixtureBites` keep their names at `:30-31`.

**The profile rows.** Make the text edits of this ticket first, because the craft-spec edit changes the size of a subject.
Then measure each subject in the ticket worktree, and write each limit in `projects/benchkit.md:490-501` by these rules:

- A row's limit is the smallest multiple of 100 that is not less than the subject's size. For example, 7,252 bytes takes 7300, and 4,700 bytes takes 4700.
- Write each limit as a plain integer with no separator, because `strconv.Atoi` parses the cell (`internal/conformance/prose_budget_test.go:145`).
- Each of the eleven exact rows takes its subject's size by the first rule.
- Add the exact row `` | `AGENTS.md` | <limit> | `` directly after the `.bench/BENCH.md` row.
- The glob row `.agents/skills/*/SKILL.md` takes the largest tracked `SKILL.md` that no exact row names. At base that is `bench-craft-cli` at 7,252 bytes, so 7300.
- The spec estimates about 11,600 bytes for `.bench/BENCH.md`. Record its measured size in the return beside that estimate.

Measure with these two commands:

```sh
bench worktree exec "<label>" -- wc -c .bench/BENCH.md AGENTS.md .agents/commands/bench-implement-spec.md .agents/commands/bench-write-spec.md .agents/commands/bench-debug.md
bench worktree exec "<label>" -- sh -c 'git ls-files ".agents/skills/*/SKILL.md" | xargs wc -c'
```

The second command also gives the size of each skill that an exact row names.
The glob row sizes from tracked files only.
The git-ignored `.agents/skills/typesafe-ai/SKILL.md` exists only in a primary checkout, and the spec's Won't handle line keeps it out of the row.
A ticket worktree lacks that file, so the live-tree test there does not see it.

**The table prose.** Add the sentence "Each limit is a size in bytes." at the end of the paragraph at `projects/benchkit.md:483-486`.
That paragraph then holds four sentences.
No anchor row quotes the prose or the rows of this section.

**The craft-spec sentence.** On `.agents/skills/bench-craft-spec/SKILL.md:68`, "current line count" becomes "current size".
The first two sentences of the line stay byte-identical, and the line stays one physical line.
No anchor quotes the third sentence, because `rg -F "budget row equal"` finds it only in the skill.

**The glossary.** In the "prose budget" entry of `CONTEXT.md:399-401`, "the per-surface line-count ceiling" becomes "the per-surface byte ceiling".
The rest of the entry stays.
No anchor quotes the entry, because `rg -F "line-count ceiling"` finds it only in `CONTEXT.md`.

**The canary.** Edit these files of `tests/canary/guidance-prose-budgets/over-budget-skill`:

- `files/dot-agents/skills/bench-craft-fixture/SKILL.md:4` becomes "This fixture skill sits inside the byte budget its profile publishes.".
- `files/projects/benchkit.md:9` gives the glob row the fixture skill's size, rounded up to a multiple of 100.
- `MUTATE.json`: `old` takes the new line-4 text. `new` adds one sentence after it that carries the file past the limit, and its text names bytes, not lines.
- `EXPECT` becomes `prose-budget exceeded: .agents/skills/bench-craft-fixture/SKILL.md is <n> bytes, over its <limit>-byte budget`.

Take `<n>` from the fixture bytes: the `wc -c` of the base fixture skill, plus the bytes that the mutation adds.
The restored fixture sits inside its limit, so `TestGuidanceProseBudgetCanaryFixtureBites` sees the diagnostic go.

**The changelog.** Add one bullet for chunk GW5 to `### Always-loaded guidance weight` in `CHANGELOG.md`, after the GW4 bullet.
This is the proposed text:

```markdown
- Changed the `guidance-prose-budgets` check so that it measures each guidance file in bytes, not in lines. Each row of the budget table in `projects/benchkit.md` now starts at its file's size, rounded up to a multiple of 100 bytes. `AGENTS.md` now has its own row. A red now names the size in bytes and the byte budget. Before, a long line cost nothing against a line budget.
```

**The structure limits.** `internal/conformance/prose_budget_test.go` has 579 lines, over the 400-line cap with no grant (`internal/structure/structure.go:308-331`).
It ends at or below 579 lines: the new test replaces the old test, and `proseBudgetLineCount` and its import go.

`internal/conformance/checks_test.go` has 642 of its 709 lines, and `internal/conformance/skill_description_budgets_test.go` has 588 of its 592 lines (`.bench/structure.budgets:49`, `:64`).
Neither file changes, because the shared symbols keep their shape.
`internal/conformance/` has no file-count headroom, so add no file there.
No count pin reads a test count or a diagnostic of this check.

**The Writes line.** The line holds each research path that a design can touch, and the five bound files of `internal/anchors` (`internal/tickets/registry_data.go:21-27`, `:45`).
It holds the eight anchor registry files that hold a literal equal to an edited guidance path (`internal/anchors/references.go:31-59`).
It also holds the 53 other fixture directories that pin `projects/benchkit.md`, the craft-spec skill, `CHANGELOG.md`, or `CONTEXT.md`.
This design edits `projects/benchkit.md`, the craft-spec skill, `CONTEXT.md`, `CHANGELOG.md`, `prose_budget_test.go`, and the one canary.
The other paths are the approved upper bound.

This ticket drops one research path:

| path | evidence |
|---|---|
| `.bench/structure.budgets` | A grant there is the reviewer's named decision (`.bench/structure.budgets:1-3`). The spec holds `prose_budget_test.go` at or below its base line count (`specs/guidance-weight/spec.md:267-268`), so no design of this ticket needs a grant. |

## Acceptance

- [ ] `TestGuidanceProseBudgetsComeFromTheProfileTable` requires the byte diagnostic for a one-line fixture over its lowered row, which a line count never raises (GW38).
- [ ] `TestGuidanceProseBudgetBoundaryCountsEveryByte` reds a subject one byte over its limit, with and without a final newline (GW39).
- [ ] `TestGuidanceProseBudgetBoundaryCountsEveryByte` passes a subject exactly at its limit, with and without a final newline (GW40).
- [ ] `TestGuidanceProseBudgetBoundaryCountsEveryByte` reds a subject whose final newline is its one excess byte (GW41).
- [ ] `TestGuidanceProseBudgetReportsEverySubjectInOneRun` requires three diagnostics that read "is <n> bytes, over its <limit>-byte budget" (GW42).
- [ ] The case `invalid limit` of `TestGuidanceProseBudgetPolicyFaultsFailClosed` requires "which is not a positive byte count" (GW43).
- [ ] Each exact row equals its subject's size at this commit, rounded up to a multiple of 100, which the GW5 review checks with `wc -c` (GW44).
- [ ] The glob row equals the largest tracked `SKILL.md` that no exact row names, rounded up, which the GW5 review checks with `git ls-files` and `wc -c` (GW45).
- [ ] `AGENTS.md` has an exact row at its size rounded up, and `TestGuidanceProseBudgetsHoldOnTheLiveTree` passes and reds when that row falls below the size (GW46).
- [ ] The `.bench/BENCH.md` row equals the guide's size at this commit, rounded up, which the GW5 review checks with `wc -c .bench/BENCH.md` (GW47).
- [ ] The table prose of `projects/benchkit.md` states "Each limit is a size in bytes.", which the GW5 review reads (GW48).
- [ ] `TestGuidanceProseBudgetCanaryFixtureBites` sees the byte `EXPECT` on the mutated fixture and its absence on the restored fixture (GW49).
- [ ] `.agents/skills/bench-craft-spec/SKILL.md:68` says "current size", which the GW5 review reads (GW50).
- [ ] The "prose budget" entry of `CONTEXT.md` names a byte ceiling, which the GW5 review reads (GW53).
- [ ] The section `### Always-loaded guidance weight` of `CHANGELOG.md` holds the GW5 bullet after the GW4 bullet, which the GW5 review reads (GW51).
- [ ] `internal/conformance/prose_budget_test.go` has at most 579 lines, and `bench test --check skill-description-budgets` passes.

## Verification

Record each red in the return, with its command and its failure line.
Run `bench worktree build <target>` before each named-check run that follows a Go edit, because a named check compiles from the run binary's source.

- Write `proseBudgetBytes`, the byte expectations, and the new boundary test before the count edit. Run `bench test --package ./internal/conformance --run TestGuidanceProseBudget`. Each test that names bytes reds against the base code. Then change the count, and record the green run of those unit tests.
- With the new count and the old line rows, run `bench test --package ./internal/conformance --run TestGuidanceProseBudgetsHoldOnTheLiveTree`. It reds for each subject, because each line row sits far below the size of its subject in bytes. Then write the byte rows, and record the green run.
- With the new count and the old canary, run `bench test --package ./internal/conformance --run TestGuidanceProseBudgetCanaryFixtureBites`. It reds, because the mutated fixture raises a byte diagnostic that the line `EXPECT` does not match. Then edit the canary, and record the green run.

The central property is a count of every byte against an inclusive limit.
A count that drops a byte, such as the final newline, and a strict comparison each break it.
After the edit, record one observed red through `bench probe` for each executable row.
Each probe below runs with `--package ./internal/conformance` and the named `--run`:

- GW39 and GW41: in `prose_budget_test.go`, swap `len(subject.Data)` for `len(strings.TrimSuffix(string(subject.Data), "\n"))` with `--run TestGuidanceProseBudgetBoundaryCountsEveryByte`. The case that is one byte over by its final newline reds.
- GW40: in `prose_budget_test.go`, swap `> limit` for `>= limit` with the same `--run`. Both at-limit cases red.
- GW38: in `prose_budget_test.go`, swap `len(subject.Data)` for `strings.Count(string(subject.Data), "\n")` with `--run TestGuidanceProseBudgetsComeFromTheProfileTable`. The test reds, because a line count never reaches the lowered row.
- GW42: in `prose_budget_test.go`, swap `bytes, over its %d-byte budget` for `lines, over its %d-line budget` with `--run TestGuidanceProseBudgetReportsEverySubjectInOneRun`. The test reds.
- GW43: in `prose_budget_test.go`, swap `positive byte count` for `positive line count` with `--run TestGuidanceProseBudgetPolicyFaultsFailClosed`. The case `invalid limit` reds.
- GW46: in `projects/benchkit.md`, swap the limit of the `AGENTS.md` row for `100` with `--run TestGuidanceProseBudgetsHoldOnTheLiveTree`. The test reds with the `AGENTS.md` byte diagnostic.
- GW49: in the canary's `files/projects/benchkit.md`, swap the glob row limit for `100000` with `--run TestGuidanceProseBudgetCanaryFixtureBites`. The test reds, because the mutation no longer crosses the row.

Each probe reports a restored subject.
For GW44, GW45, and GW47, record each measured size and its derived limit in the return as a table.
Then run `bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench`.
Also run the checks `guidance-prose-budgets`, `skill-description-budgets`, `prose-mechanics`, and `root-conformance`.
