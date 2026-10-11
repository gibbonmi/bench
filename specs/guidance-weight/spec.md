# Shrink the always-loaded guidance and budget it in bytes

Status: staged

Roadmap: FT100

Decision source: specs/guidance-weight/decisions/guidance-weight.md (ready compiled map), decisions #9, #10, #14, and #15 and then #3, #6, #7, and #11, as decision #16 splits the outcome

Audience: every repository that links the kit, for `.bench/BENCH.md` and the phase files, and this repository, for `AGENTS.md`, the profile, and the checks

Verification log: <n> iteration(s) to accept — <note>

## Problem

Every session loads `.bench/BENCH.md` and `AGENTS.md`, so each of their bytes costs every session.
`.bench/BENCH.md` is 15,797 bytes and 181 lines at main `8c8a6ca5`, and `AGENTS.md` is 4,800 bytes and 86 lines.
Much of that text acts only in some phases, or it repeats a rule that an on-demand owner already states (`specs/guidance-weight/decisions/guidance-weight/assets/clause-inventory.md`).

The build-orchestration clauses B11 and B26 to B33 hold about 3,450 bytes, and only the build, review, and close phases act on them.
Six places state a light-path route, and two of them repeat the drain dispatch.
The batch-approval rule B43 acts only at spec sign-off.
Three clauses repeat a rule that an owner keeps: the B10 staleness-pass exception, B12, and `AGENTS.md` A17.

The budget check cannot see this weight.
The `guidance-prose-budgets` check counts lines (`internal/conformance/prose_budget_test.go:204-211`), so a long line costs nothing.
Nine of the eleven exact rows equal their file's current line count today, and `AGENTS.md` has no row (`projects/benchkit.md:488-501`).

## Solution

Each phase-only clause moves word for word, with its anchors, into a file that its phases load.
The build rules move into `.agents/skills/bench-implement-spec/references/build-rules.md`, and `.bench/BENCH.md` keeps a one-line pointer.
The light-path statements merge into one block under the right-size table, and the repeated drain sentences go.
The batch-approval rule moves beside the sign-off stop of `/bench-write-spec`.
The repeated copies go, and the B10 exception moves into the staleness pass, which is its only acting phase.

Then every row of the budget table limits its subject in bytes.
Each row starts at its file's size after the moves, rounded up to a multiple of 100 bytes.
`AGENTS.md` gets its own row, and `.bench/BENCH.md` gets its bound of about 11,600 bytes.

## User stories

Line: opus / xhigh
Implementation-line reason: GW1 is the hardest material chunk. It moves anchors, forbid rows, pointers, mirrors, and canaries across owner seams, so a repair round costs over 100k tokens. The scorecard routes such a ticket to Opus at xhigh, and the leverage override covers each guidance chunk. GW2 to GW5 are exact edits at known seams, so their authors can lower effort in session under `craft-line`.
Harder chunks: GW1, GW5.

### The build rules move (decision #9)

1. As a session outside a build, I want the build clauses gone from `.bench/BENCH.md`, so that I pay only for rules I use.
2. As a build orchestrator, I want B11 and B26 to B33 word for word in `build-rules.md`, so that the rules stay the same.
3. As any session, I want `.bench/BENCH.md` to keep one pointer to `build-rules.md`, so that I can find the build rules.
4. As a build orchestrator, I want `build-rules.md` to open with a lead that names its phases, so that I know when it applies.
5. As a maintainer, I want the 21 moved require anchors to bind `build-rules.md`, so that a dropped rule still reds the gate.
6. As a maintainer, I want the continuation row to name a section of `build-rules.md`, so that the move does not red it.
7. As a maintainer, I want the five build-rule forbid rows to bind `build-rules.md`, so that a retired sentence cannot return there.
8. As a reader of a phase file, I want each of the 28 build-rule pointers to name `build-rules.md`, so that none misleads me.
9. As a maintainer, I want each anchor that quotes a repointed line to quote its new text, so that it stays green and bites.
10. As a maintainer, I want the conformance mirrors of the moved rows to name `build-rules.md`, so that the independent expectations still match their rows.
11. As a maintainer, I want the six delegated canaries to mutate `build-rules.md`, so that each canary still bites its moved row.
12. As a maintainer, I want the three pointer canaries to quote the new pointer text, so that each mutation still finds its target once.
13. As a reader of a red, I want each moved row's diagnostic to name the build rules, so that it names the right file.
14. As any session, I want the predicate lead-in to count the bullets that stay, so that the guide promises no lost bullet.

### The repeated copies (decision #14)

15. As any session, I want B12 gone from `.bench/BENCH.md`, so that the owned-red rule loads only where `craft-line` states it.
16. As a maintainer, I want the B12 anchor to retire, because the `craft-line` copy keeps its own anchor.
17. As a hand runner of the system suite, I want `projects/benchkit.md` to keep the route, so that removing A17 loses no rule.
18. As a maintainer, I want the A17 anchor, mirror, canary, and fixture plant to retire together, so that no check expects the sentence.
19. As the staleness-pass orchestrator, I want the B10 exception in `staleness-pass.md` with its flag duty, so that I decide and flag each finding.
20. As a session outside the staleness pass, I want the B10 bullet to end at "Ask about a behavioral one.", so that I ask.

### The light-path block (decisions #10 and #4)

21. As any session, I want the light-path rules in one block under the right-size table, so that I read one place for the light path.
22. As any session, I want the table to keep its inline route word for word, so that a reviewer-asked change still goes inline.
23. As a session with a learning or drain fix, I want the block to route it to a write delegate, not inline.
24. As a build orchestrator, I want the block to keep the build-hold pointer, so that a build dispatches no light-path fix.
25. As any session, I want the fix-don't-park split in the block, without the repeated drain clause, so that the split reads once.
26. As any session, I want the Capture section without its repeated drain sentence, so that the drain dispatch has one statement.
27. As a maintainer, I want each anchored needle of the block kept word for word, the B40 lead included, so that each anchor binds.
28. As a maintainer, I want the commitment mirrors to find the kept B40 lead once, so that the outside-a-build mirror still bites.
29. As a maintainer, I want the three shared-rule markers to stay in the Workflow section, so that the single-source and section checks stay green.

### The batch-approval rule (decision #15)

30. As a spec stage, I want B43 beside the sign-off stop of `/bench-write-spec`, so that the batch approval loads where sign-off happens.
31. As any session, I want B43 gone from `.bench/BENCH.md`, so that a session that never signs off a spec does not pay for it.
32. As a slicer, I want craft-tickets to point to `/bench-write-spec` for the batch-approval carve-out, so that the pointer names the owner.
33. As a drain, I want `bench-drain.md` to name its own diff approval in its own words, so that "batch approval" keeps one meaning.
34. As a maintainer, I want the drain and craft-tickets rows to quote the new text, so that both stay green and bite.
35. As a maintainer, I want `bench-write-spec.md` to stay inside its line budget after B43 lands, so that the move does not red the budget check.

### The byte budget (decisions #3, #6, #7, and #11)

36. As a maintainer, I want the budget check to count the bytes of each subject, so that a long line costs what it weighs.
37. As a maintainer, I want a subject one byte over its limit to red the check, so that no growth passes unseen.
38. As a maintainer, I want a subject at its limit to pass, so that the limit is the largest allowed size.
39. As a maintainer, I want the final newline to count as one byte, so that the count equals `wc -c`.
40. As a maintainer, I want the diagnostic and the limit error to name bytes, so that a red names its unit.
41. As a maintainer, I want each exact row to start at its post-move size, rounded up to 100, so that rows start tight.
42. As a maintainer, I want the glob row to start at its largest tracked glob-only file, rounded up, so that new skills arrive budgeted.
43. As a maintainer, I want `AGENTS.md` to get an exact row, so that the agreement cannot grow without a table edit.
44. As a maintainer, I want the `.bench/BENCH.md` row at its final size, rounded up, so that later growth needs a reviewed edit.
45. As a maintainer, I want the table prose to state that each limit is a size in bytes, so that a reader knows the unit.
46. As a maintainer, I want the over-budget canary to bite in bytes, so that the fixture proves the new count.
47. As a spec author, I want the craft-spec budget sentence to name the current size, so that the rule matches the byte measure.

### Reviewed exclusions

48. As the reviewer, I want no rule with no other copy deleted, so that a deletion waits for an FT231 result.
49. As the reviewer, I want no description content phrase trimmed, so that a trim waits for an FT231 trial.
50. As the reviewer, I want each guidance edit recorded in one `CHANGELOG.md` entry, so that a linked repository sees the change.

## Implementation decisions

### Terms

- **Move** puts a clause, word for word and with its anchors, into a file that its acting phases load.
- **Repoint** changes a pointer line, an anchor `File`, a mirror, or a canary path from `.bench/BENCH.md` to the new owner, with no other text change.
- **Build rules** names `.agents/skills/bench-implement-spec/references/build-rules.md`. Avoid: build reference, orchestration file.

### Closed decisions

The decisions #2, #3, #4, #6, #7, #9, #10, #11, #14, #15, and #16 of the map stay closed.
The reviewer's consultant (Fable/high) answered one question and graded thirteen calls on 2026-10-11, by the reviewer's direction while the reviewer was away.
Each answer is a decision of record under Flagged additions, open to the reviewer's veto at sign-off.

### The build rules file

The file `build-rules.md` holds a heading, a lead, and the moved clauses in their current order.
The lead sentence is: "`/bench-implement-spec`, `/bench-review-implementation`, and `/bench-final-check` load these rules for each spec-backed build."
The clauses keep their bold labels and their reviewer voice ("me", "my") word for word.
Each clause group takes an H2: `## Acceptance shortfall` (B11), `## Ticket authors` (B26), `## Chunks and review` (B27 and B28), `## Repairs` (B29), `## Plan expansion` (B30 and B31), and `## Delegated runs` (B32 and B33).

The B11 bullet becomes this paragraph: "If a build cannot meet an acceptance row, that is a **material acceptance shortfall**: the build exits and reports. It does not land a silent partial."
`.bench/BENCH.md` keeps one pointer in `## Workflow`, after the numbered phase list.
The pointer reads: "`.agents/skills/bench-implement-spec/references/build-rules.md` holds the build rules for authors, chunks, repairs, plan expansion, and delegated runs."

### The anchors and their mirrors

Each require row of B11 and B26 to B33 changes its `File` to the build rules and keeps its needle.
The rows are `internal/anchors/registry_data.go:377`, `registry_ft311_review_dispatch.go:39`, `:40`, `:45` to `:48`, `:51`, and `:52`.
They also include `registry_retained_workflow.go:213` to `:216`, `:237` to `:242`, `:294`, and `:341`.
The row at `registry_retained_workflow.go:341` is `RequireInSection` with section `Workflow`, so its section becomes `Ticket authors`.

The forbid rows with `File` `.bench/BENCH.md` at `registry_ft311_review_dispatch.go:41`, `:49`, `:50` and `registry_retained_workflow.go:217`, `:270` repoint to the build rules.
The forbid rows at `registry_retained_workflow.go:295` to `:302` keep their phase files.
No new forbid row binds `.bench/BENCH.md` for this text, so GW01 stays review-owned. No canary overlay holds B11 or B26 to B33, so a forbid row there stays a reviewer option.

A diagnostic that names the operating guide as the owner of a moved rule names the new owner instead.
The new owner is the build rules for GW1 and `/bench-write-spec` for GW4.
The rule covers moved rows and unmoved rows, such as `registry_retained_workflow.go:46`, `:245`, `:295-303`, `:345`, and `:30-31`.
The B11 row `registry_data.go:377` names `.bench/BENCH.md` as the file that dropped its text, so its diagnostic names the build rules too.

The mirrors repoint with their rows: `internal/conformance/retained_workflow_test.go:118-137`, `:151-166`, `:183-202`, and `internal/conformance/implementation_continuation_test.go:81-84`.
The `wantDiagnostics` list of `retained_workflow_test.go:46-95` takes the renamed diagnostics.
Two new require rows pin the new pointer in `.bench/BENCH.md` and the lead of the build rules.
A planned test proves the absence and lead diagnostics of the build rules, as `TestReferenceFileAnchorsRedOnAbsence` (`internal/anchors/registry_data_test.go:571`) proves them for its three precedents.

The structure lane reds a Go file that is over its cap and gains a line (`internal/structure/structure.go:308-331`, cap 400 lines).
`internal/anchors/registry_data.go` (478 lines) and `internal/anchors/registry_data_test.go` (1,286 lines) are over the cap with no grant.
So each ticket keeps those two files at or below their base line counts.
The new require rows go in `internal/anchors/registry_chunk_chain.go`, which already binds the staleness pass, and the new tests go in `internal/anchors/registry_chunk_chain_test.go`.

### The pointer lines

Each of these lines names the build rules by its full path in place of `.bench/BENCH.md`, on the same physical line:

- `.agents/commands/bench-implement-spec.md:21`, `:23`, `:33`, `:54`, `:60`, `:90`
- `.agents/commands/bench-review-implementation.md:34`, `:44`, `:50`, `:52`
- `.agents/commands/bench-final-check.md:194`, `.agents/commands/bench-write-spec.md:50`, `.agents/commands/bench-debug.md:157`
- `.agents/skills/bench-craft-spec/SKILL.md:68`, `:70`, `.agents/skills/bench-craft-delegate/SKILL.md:13`, `:122`, `.agents/skills/bench-craft-gate/SKILL.md:28`
- `.agents/skills/bench-craft-tickets/SKILL.md:67`, and `:95`, which says "the operating guide's delegated policy"
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:10`, `:39`, `:94`, `:105`
- `.agents/skills/bench-implement-spec/references/staleness-pass.md:52`
- `docs/field-guide.html:1333` and the sentence of `.bench/BENCH-reference.md` that says "under the repair rule of `.bench/BENCH.md`"

The anchors that quote a repointed line take the new text.
They are `registry_retained_workflow.go:47`, `:212`, `:218` to `:220`, `:280`, `:303`, `:307`, `:310`, and `:316`.
They also include `registry_data.go:136`, `:140`, `:442`, `registry_chunk_chain.go:62`, and `registry_ft311_review_dispatch.go:42`.
The mirrors that quote them take the new text: `retained_workflow_test.go:138-141`, `:146-150`, `:167-178`. The line `.agents/commands/bench-drain.md:72` names the light-path dispatch rule, which stays, so it and its mirrors `recurrence_maintenance_contract_test.go:91` and `:217` stay unchanged.
The anchor harness mirrors take it too: `internal/anchors/registry_data_test.go:349` and `:892`, and `internal/anchors/registry_chunk_chain_test.go:51`.

The canaries `review-repair-ticket-owner`, `review-repair-ticket-covers`, and `delegated-serial-ticket-checkpoint` quote the new text in `MUTATE.json`.
The six canaries `delegated-chunk-tip-review`, `delegated-green-predecessor`, `delegated-opt-in-entry`, `delegated-pending-predecessor-stop`, `delegated-prerequisite-checkpoint`, and `delegated-repair-ownership` repoint `BASE` and the `MUTATE.json` path, and each `EXPECT` takes the renamed diagnostic.

### The predicate lead-in

`.bench/BENCH.md:59` says "Three predicates ride with them:".
The GW1 ticket moves B11 and changes the lead-in to "Two predicates ride with them:".
The GW2 ticket removes B12 and changes it to "One predicate rides with them:".
No anchor or check pins the lead-in.

### The repeated copies and the B10 move

B12 (`.bench/BENCH.md:66-67`) goes, and its row `registry_data.go:378` retires.
The `craft-line` copy (`.agents/skills/bench-craft-line/SKILL.md:80-81`) keeps its row `registry_data.go:191`.
A17 (`AGENTS.md:85-86`) goes, and its row `registry_data.go:383` retires with three dependents.
The dependents are the `AGENTS.md` half of `TestSystemSuiteRouteAnchorsRedOnRemoval` (`registry_data_test.go:1240-1255`), the canary `agents-system-suite-route`, and the plant `cmd/bench/testdata/anchors/fixture-repo/AGENTS.md:11` that `cmd/bench/anchor_help_test.go` reads.
The other readers of the plant, `cmd/bench/anchors_dir_test.go` and `cmd/bench/command_registry_test.go`, count from the registry and need no edit.
The `projects/benchkit.md:676-677` copy keeps its row `registry_data.go:385` and its canary `benchkit-system-suite-route`.

The B10 exception moves into `staleness-pass.md`.
`.bench/BENCH.md:62-63` then ends the bullet at "Ask about a behavioral one."
The last sentence of `staleness-pass.md:45-46` becomes two sentences.
The first reads: "It decides each finding by the spec-contradiction predicate of `.bench/BENCH.md`, a behavioral one included, and flags it."
The second reads: "It does not ask the reviewer."

The row `registry_data.go:376` keeps its needle, which sits in the kept first sentence.
A new require row pins the moved sentence in `staleness-pass.md`.

### The light-path block

The block sits directly under the right-size table in `## Workflow`, before the paragraph that starts "**Deliver only the committed outcome.**".
It holds four paragraphs, because a bold label counts as one sentence and the bound is six sentences (`internal/prose/prose.go:18`).

1. B39 stays word for word: "**A light-path change needs no commitment.**" and its three sentences.
2. "**Delegate each light-path fix.** Outside a build, a `bench learning` entry's light-path fix ships at once only when the active outcome needs it; `/bench-implement-spec` states the build rule. Every light-path fix from a `bench learning` entry or a drain goes to a write delegate. Only a change that I ask for in the session goes inline under the table. `/bench-drain` implements no light-path fix inline: it dispatches every light-path fix in its run, whatever its source and whatever the active commitment is."
3. The three remaining B40 sentences stay word for word. They run from "Dispatch that fix" to "as the light-path row states."
4. B42 stays word for word, except that ", and the drain delegates its light-path fix" goes from its second sentence.

The Capture sentence "A drain implements a light-path idea under the light-path fix rule instead." goes.
The second sentence of paragraph 2 is the B40 lead word for word, so its timing rule stays.
The third and fourth sentences of paragraph 2 add the delegated route of #4.
Each sentence that `TestCommitmentGuidance` quotes with a leading space stays in the middle of a physical line (`internal/conformance/commitment_guidance_test.go:49-52` and `:54`).
The row `registry_commitment.go:26` and the mirror `commitment_guidance_test.go:53` keep their needle and their diagnostic.

### The batch-approval move

B43 (`.bench/BENCH.md:156-159`) moves word for word into `.agents/commands/bench-write-spec.md`.
It sits as its own paragraph in `## Exit handoff`, directly after the paragraph that holds "Stop for sign-off.", and that paragraph stays byte-identical.
That paragraph runs over lines 15 and 16 and ends at "`bench commitment start`.", so B43 starts after line 16.
The craft-tickets sentence at `.agents/skills/bench-craft-tickets/SKILL.md:37` becomes this text:

```text
For spec-backed builds, this is the only route onto the frontier; the batch-approval AFK carve-out in `/bench-write-spec` is the sole no-round-trip exception.
```

The row `registry_data.go:207` takes that sentence as the first half of its needle, and its second sentence at `:38` stays byte-identical.
The second sentence of `bench-drain.md:291` becomes "That approval is the reviewer's sign-off on this drain's proposed diff.", and the row `registry_retained_workflow.go:29` takes it.
The forbid rows `registry_retained_workflow.go:28`, `:30`, `:31` and `registry_data.go:182` stay.

### The prose budgets during the moves

Decision #16 puts every move before the byte switch, so the line rows grade each move ticket.
`.agents/commands/bench-write-spec.md` sits at 87 of 87 lines and `.agents/commands/bench-implement-spec.md` at 93 of 93 (`projects/benchkit.md:491-492`).
Each pointer edit stays on its physical line, so it adds no line.
The GW4 ticket raises the `bench-write-spec.md` line row by the lines that B43 adds, in the same commit.
At base `bc2ea0ff` that is four lines and one blank line, so the row goes from 87 to 92.
The `guidance-prose-budgets` lane runs only on a profile edit, so each move ticket runs `bench test --check guidance-prose-budgets` itself.

### The byte budget

`checkGuidanceProseBudgets` compares `len(subject.Data)` with the row limit, and `proseBudgetLineCount` goes.
The diagnostic becomes "prose-budget exceeded: %s is %d bytes, over its %d-byte budget".
The limit error of `parseProseBudgetPolicy` becomes "not a positive byte count".
The table header stays `| subject | limit |`, because the parser requires it (`prose_budget_test.go:132-135`).
The table prose of `projects/benchkit.md:483-486` gains the sentence "Each limit is a size in bytes."

The rows take their sizes from the tree after GW1 to GW4, at the commit of GW5.
A row's limit is the smallest multiple of 100 that is not less than the size.
The glob row takes the largest tracked `SKILL.md` that no exact row names: `bench-craft-cli` at 7,252 bytes today, so 7,300.
`AGENTS.md` takes a new exact row, and `.bench/BENCH.md` takes its size after this outcome, about 11,600.
`.agents/skills/bench-craft-spec/SKILL.md:68` says "current size" in place of "current line count".

The tests change their unit: `TestGuidanceProseBudgetBoundaryIsNewlineIndifferent` becomes `TestGuidanceProseBudgetBoundaryCountsEveryByte`, and `proseBudgetLines` becomes `proseBudgetBytes`, and it and `healthyProseBudgetFiles` build byte sizes.

`internal/conformance/prose_budget_test.go` has 579 lines, over the 400-line structure cap with no grant, so the GW5 ticket keeps it at or below its base line count.
The rename replaces the old test, so the file gains no test.
`internal/conformance/checks_test.go:576` reads `proseBudgetRows`, so it changes only if that helper changes shape.
The over-budget canary's fixture skill and `MUTATE.json` describe the mutation in lines, and both take byte wording.
`CONTEXT.md` defines "prose budget" as "the per-surface line-count ceiling", so the GW5 ticket makes it a byte ceiling.

The expected diagnostics at `prose_budget_test.go:276`, `:284`, `:335`, `:506-508`, and `:534` name bytes, and the `invalid limit` case at `:315` matches the whole new limit error.
The canary `tests/canary/guidance-prose-budgets/over-budget-skill` sets its fixture row in bytes, and its mutation adds bytes past it.
`skill_description_budgets_test.go` reuses `proseBudgetPolicy` and `limitFor` but has its own parse loop, so it changes nothing.

### The changelog entry

One `### Always-loaded guidance weight` section under `## [Unreleased]` of `CHANGELOG.md` holds one bullet for each chunk.
Each chunk has one bullet. A later ticket of the same chunk extends that bullet, so the tickets that write `CHANGELOG.md` run serially.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| GW1 / `1-move-the-build-rules-into-their-own-file.md`, `2-repoint-the-build-rule-pointers.md` | The build rules live in `build-rules.md`, with their anchors, mirrors, canaries, and every pointer, and `.bench/BENCH.md` keeps one pointer | GW01-GW15, GW51, GW52 | the `internal/anchors`, `internal/conformance`, and `cmd/bench` package tests, and the `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, and `root-conformance` checks | yes |
| GW2 / `3-remove-the-repeated-copies.md` | The repeated copies B12 and A17 go with their anchors, and the B10 exception lives in the staleness pass | GW14, GW16-GW22, GW51 | the `internal/anchors`, `internal/conformance`, and `cmd/bench` package tests, and the `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, and `root-conformance` checks | no |
| GW3 / `4-merge-the-light-path-rules-into-one-block.md` | The light-path rules read as one block under the table, and the repeated drain sentences go | GW23-GW31, GW51 | the `internal/anchors`, `internal/conformance`, and `cmd/bench` package tests, and the `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, and `root-conformance` checks | no |
| GW4 / `5-move-the-batch-approval-rule-beside-the-sign-off-stop.md` | The batch-approval rule sits beside the sign-off stop, and the drain names its own approval | GW32-GW37, GW51 | the `internal/anchors`, `internal/conformance`, and `cmd/bench` package tests, and the `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, and `root-conformance` checks | no |
| GW5 / `6-measure-the-guidance-budgets-in-bytes.md` | Every budget row limits its subject in bytes, `AGENTS.md` has a row, and `.bench/BENCH.md` has its bound | GW38-GW51, GW53 | the `internal/anchors`, `internal/conformance`, and `cmd/bench` package tests, and the `guidance-prose-budgets`, `skill-description-budgets`, `prose-mechanics`, and `root-conformance` checks | yes |

GW1 to GW4 write `.bench/BENCH.md`, `CHANGELOG.md`, and the same anchor registries, so they run serially in the order of the table.
Each chunk review also grades GW51, because each ticket writes its own changelog bullet.
GW5 follows GW4, because its rows take the sizes after every move.

### Completion plan

This version-1 completion plan is a planning verification inventory, not a dispatch plan.
Its commands are future obligations, not executed proof.
Each slicer adds its chunk to this plan when it writes the chunk's tickets, because the plan names only ticket files that exist.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "GW1",
      "tickets": ["1-move-the-build-rules-into-their-own-file.md", "2-repoint-the-build-rule-pointers.md"],
      "verification": [
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
        {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
        {"id": "prose", "command": "bench test --check prose-mechanics"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"}
      ]
    },
    {
      "id": "GW2",
      "tickets": ["3-remove-the-repeated-copies.md"],
      "verification": [
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
        {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
        {"id": "prose", "command": "bench test --check prose-mechanics"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"}
      ]
    },
    {
      "id": "GW3",
      "tickets": ["4-merge-the-light-path-rules-into-one-block.md"],
      "verification": [
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
        {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
        {"id": "prose", "command": "bench test --check prose-mechanics"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"}
      ]
    },
    {
      "id": "GW4",
      "tickets": ["5-move-the-batch-approval-rule-beside-the-sign-off-stop.md"],
      "verification": [
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
        {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
        {"id": "prose", "command": "bench test --check prose-mechanics"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"}
      ]
    },
    {
      "id": "GW5",
      "tickets": ["6-measure-the-guidance-budgets-in-bytes.md"],
      "verification": [
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
        {"id": "descriptions", "command": "bench test --check skill-description-budgets"},
        {"id": "prose", "command": "bench test --check prose-mechanics"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"}
      ]
    }
  ],
  "final_verification": [
    {"id": "coverage", "command": "bench coverage --check specs/guidance-weight/spec.md"},
    {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
    {"id": "command-package", "command": "bench test --package ./cmd/bench"},
    {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
    {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
    {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
    {"id": "descriptions", "command": "bench test --check skill-description-budgets"},
    {"id": "prose", "command": "bench test --check prose-mechanics"},
    {"id": "skills-index", "command": "bench test --check skills-index-command-adapters"},
    {"id": "root-conformance", "command": "bench test --check root-conformance"}
  ]
}
```

## Testing decisions

- A good test reads the live tree for each moved, repointed, or retired sentence, and it reds when the sentence leaves its new owner.
- The anchor rows attach at the `docs-currency-workflow` check, which grades the real tree (`internal/conformance/docs_workflow_checks_test.go:20-36`). The anchor harness grades only a synthetic tree.
- An anchor omission probe runs through `bench test --check docs-currency-workflow`, by the scorecard rule.
- The mirrors of `retained_workflow_test.go`, `implementation_continuation_test.go`, and `commitment_guidance_test.go` are the permitted independent expectations, so each repointed mirror records its observed red.
- The canary fixtures attach at `TestEveryRetainedFixtureBitesThroughRegisteredOwner` (`internal/conformance/fixture_bite_test.go:21`).
- The byte count attaches at the existing unit tests of `prose_budget_test.go` and at its live-tree test.

### Seam diagram

    trigger: the landing gate, or a lane run of a named check
        │
        ▼
    live tree  ──▶  [ docs-currency-workflow: anchor rows over each File ]  ──▶  diagnostics
                      ◀ tests attach here: TestRootConformance, the mirrors, the canary bites
        │
        ▼
    profile table + subjects  ──▶  [ guidance-prose-budgets: len(bytes) > limit ]  ──▶  diagnostics
                      ◀ tests attach here: the prose_budget_test.go unit and live-tree tests

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GW01 | 1 | `.bench/BENCH.md` holds none of the needles of B11 and B26 to B33 | review-owned: the GW1 review runs `rg` for each moved needle over `.bench/BENCH.md` | No forbid row binds `.bench/BENCH.md` for this text, so only review sees a copy left behind. |
| GW02 | 2 | `build-rules.md` holds each moved clause word for word, in the current order | review-owned: the GW1 review compares the clauses of `.bench/BENCH.md` at the base with the new file | A require row pins only its needle, so a reworded sentence beside a needle passes every check. |
| GW03 | 3 | `.bench/BENCH.md` holds the one-line pointer to `build-rules.md` in `## Workflow` | planned TestBuildRulesPointerAnchorsRedOnRemoval in internal/anchors/registry_chunk_chain_test.go | The new require row reds when the pointer goes, so the always-loaded guide cannot lose its route. |
| GW04 | 4 | `build-rules.md` opens with its lead sentence | planned TestBuildRulesReferenceAnchorsRedOnAbsence in internal/anchors/registry_chunk_chain_test.go | The new test reds when the file is absent or its lead goes. |
| GW05 | 5 | Each of the 21 moved require rows binds `build-rules.md` and keeps its needle | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | A row that kept `.bench/BENCH.md` reds on the live tree, because the needle left that file. |
| GW06 | 6 | The continuation row names section `Ticket authors` of `build-rules.md` | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | A row that kept section `Workflow` reds, because the new file has no such section. |
| GW07 | 7 | The forbid row `registry_retained_workflow.go:217` binds `build-rules.md` | `internal/conformance/retained_workflow_test.go` (`TestRetainedWorkflow`) | The mirror at `:163-166` names the file, so the row left on `.bench/BENCH.md` reds the comparison. |
| GW52 | 7 | The four fresh-author forbid rows `registry_ft311_review_dispatch.go:41`, `:49`, `:50` and `registry_retained_workflow.go:270` bind `build-rules.md` | review-owned: the GW1 review reads the `File` of each of the four rows | `TestFreshTicketAuthors` plants each needle in the row's own `File`, so a row left on `.bench/BENCH.md` passes it. |
| GW08 | 8 | No line under `.agents`, `docs`, or `.bench/BENCH-reference.md` names `.bench/BENCH.md` as the owner of a moved rule | review-owned: the GW1 review runs `rg -F "BENCH.md"` and reads each hit against the 27-line list | Only 15 of the 27 pointer lines carry an anchor, so review sees a stale unanchored pointer. |
| GW09 | 9 | Each anchor that quotes a repointed line quotes its new text | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | An anchor that kept the old text reds, because the old pointer text left the line. |
| GW10 | 10 | The mirrors of `retained_workflow_test.go` name `build-rules.md` for each moved row | `internal/conformance/retained_workflow_test.go` (`TestRetainedWorkflow`) | The drift check compares File, Section, and Needle, so a stale mirror reds against the moved row. |
| GW11 | 11 | Each of the six delegated canaries mutates `build-rules.md` and bites its moved row | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A canary that kept `.bench/BENCH.md` finds no `old` text, so its materializer refuses and the test reds. |
| GW12 | 12 | The three pointer canaries quote the new pointer text in `MUTATE.json` | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A stale `old` string matches zero times, so the materializer refuses and the test reds. |
| GW13 | 13 | No diagnostic names the operating guide as the owner of a moved rule | review-owned: the GW1 review runs `rg "operating guide"` over the anchor registries | A diagnostic is a label that no check grades, so only review sees a stale subject. |
| GW14 | 14 | The predicate lead-in of `.bench/BENCH.md` counts the bullets under it | review-owned: the GW1 and GW2 reviews read the lead-in | No anchor pins the lead-in, so only review sees a wrong count. |
| GW15 | 2 | `build-rules.md` passes the prose bounds | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | The `prose-mechanics` check grades every tracked Markdown file, so a moved paragraph over a bound reds. |
| GW16 | 15 | `.bench/BENCH.md` no longer holds the B12 needle "only diff-owned reds count toward fix-loop convergence" | review-owned: the GW2 review runs `rg` for the needle | The row retires with the sentence, so no check reds a sentence left behind. |
| GW17 | 16 | The `craft-line` row `registry_data.go:191` still binds its kept copy | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | The row reds if the GW2 edit touches the kept copy. |
| GW18 | 17 | `projects/benchkit.md` keeps the hand-run route and its anchor | `internal/anchors/registry_data_test.go` (`TestSystemSuiteRouteAnchorsRedOnRemoval`) | The kept half of the test reds when the benchkit copy goes. |
| GW19 | 18 | No check expects the A17 sentence in `AGENTS.md` | `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsNeedleLines`) | A plant left in the fixture repo moves each later plant two lines down, so the help test reds. |
| GW20 | 18 | The canary `agents-system-suite-route` is gone with its row | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A canary whose owner row retired never sees its `EXPECT`, so the test reds. |
| GW21 | 19 | `staleness-pass.md` states the moved exception with its flag duty | planned TestStalenessExceptionAnchorRedsOnRemoval in internal/anchors/registry_chunk_chain_test.go | The new require row reds when the moved sentence or its flag duty goes. |
| GW22 | 20 | The B10 bullet of `.bench/BENCH.md` ends at "Ask about a behavioral one." | review-owned: the GW2 review reads the bullet | The row `registry_data.go:376` pins only the first sentence, so only review sees a kept exception. |
| GW23 | 21 | The light-path paragraphs sit directly under the right-size table | review-owned: the GW3 review reads the Workflow section order | Every block row is section-scoped or whole-file, so no row reds a scattered order. |
| GW24 | 22 | The table keeps its inline route needle | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | The require row `registry_data.go:187` reds when the table text changes. |
| GW25 | 23 | The block sends a learning or drain fix to a write delegate | review-owned: the GW3 review reads paragraph 2 for the route sentences of #4 | No row pins the new route sentences, so only review sees them missing. |
| GW26 | 24 | The block keeps "`/bench-implement-spec` states the build rule." | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | The `build rule pointer` case removes the sentence and requires its diagnostic. |
| GW27 | 25 | The fix-don't-park paragraph holds no drain clause | review-owned: the GW3 review runs `rg "drain delegates its light-path fix"` over `.bench/BENCH.md` | No row pins the removed clause, so only review sees it kept. |
| GW28 | 26 | The Capture section holds no "light-path fix rule instead" sentence | review-owned: the GW3 review runs `rg` for the sentence | No row pins the removed sentence, so only review sees it kept. |
| GW29 | 27 | Each block row keeps its needle and stays green | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | Each row of `registry_commitment.go:22-27` and `registry_retained_workflow.go:243-244` reds when its needle changes. |
| GW30 | 28 | Each commitment mirror finds its sentence exactly once in the live guide | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A sentence that moved to a line start loses its leading space, so its count is zero and the case reds. |
| GW31 | 29 | The right-size, fix-don't-park, and commitment markers stay in `## Workflow` of `.bench/BENCH.md` | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The canary `fix-dont-park-section-relocated` bites the section rule, and the section rows red a marker outside Workflow. |
| GW32 | 30 | `bench-write-spec.md` holds B43 word for word after the sign-off paragraph | review-owned: the GW4 review compares B43 at the base with the new paragraph | B43 has no anchor of its own, so only review sees a reworded rule. |
| GW33 | 31 | `.bench/BENCH.md` holds no B43 sentence | review-owned: the GW4 review runs `rg "batch approval covers"` over `.bench/BENCH.md` | No row pins B43, so only review sees a copy left behind. |
| GW34 | 32 | craft-tickets names `/bench-write-spec` as the home of the carve-out | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The canary `ticket-light-path-carve-out` bites the row `registry_data.go:207`, which takes the new text. |
| GW35 | 33 | `bench-drain.md` names the reviewer's sign-off on its own diff, not the batch approval rule | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | The require row `registry_retained_workflow.go:29` takes the new sentence, so the old pointer text reds. |
| GW36 | 34 | The forbid rows on `bench-drain.md` stay green after the edit | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | Each forbid row reds on the live tree when its needle returns to `bench-drain.md`, and no canary expects these diagnostics. |
| GW37 | 35 | `bench-write-spec.md` stays inside its line row after B43 lands | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | A line added past the row reds this live-tree check. |
| GW38 | 36 | The check grades each subject by `len` of its bytes | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsComeFromTheProfileTable`) | The test lowers a row below a fixture's byte size and requires the byte diagnostic, which a line count never raises. |
| GW39 | 37 | A subject one byte over its limit reds | planned TestGuidanceProseBudgetBoundaryCountsEveryByte in internal/conformance/prose_budget_test.go | A count that drops a byte, such as the trimmed final newline, passes the one-byte-over fixture. |
| GW40 | 38 | A subject exactly at its limit passes | planned TestGuidanceProseBudgetBoundaryCountsEveryByte in internal/conformance/prose_budget_test.go | A strict comparison reds the at-limit fixture. |
| GW41 | 39 | A final newline counts as one byte | planned TestGuidanceProseBudgetBoundaryCountsEveryByte in internal/conformance/prose_budget_test.go | A count that trims the final newline passes a fixture whose newline is its one excess byte. |
| GW42 | 40 | The diagnostic reads "is <n> bytes, over its <limit>-byte budget" | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetReportsEverySubjectInOneRun`) | The expected diagnostics name bytes, so a stale line format reds. |
| GW43 | 40 | A non-positive limit cell fails closed with "not a positive byte count" | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetPolicyFaultsFailClosed`) | The fault case requires the new text, so a stale message reds. |
| GW44 | 41 | Each exact row equals its file's size at the GW5 commit, rounded up to a multiple of 100 | review-owned: the GW5 review runs `wc -c` on each subject at the GW5 commit and checks the arithmetic | A row equal to any value below the size reds the live tree, but a loose row passes, so only review sees slack. |
| GW45 | 42 | The glob row equals the largest tracked `SKILL.md` that no exact row names, rounded up | review-owned: the GW5 review runs `git ls-files` and `wc -c` over `.agents/skills/*/SKILL.md` | A loose glob row passes every check, so only review sees slack. |
| GW46 | 43 | `AGENTS.md` has an exact row at its size rounded up | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | An exact row enters the universe, so a growth of `AGENTS.md` past its row reds. |
| GW47 | 44 | The `.bench/BENCH.md` row equals its size after this outcome, rounded up | review-owned: the GW5 review runs `wc -c .bench/BENCH.md` at the GW5 commit | A loose row passes the live tree, so only review sees slack. |
| GW48 | 45 | The table prose states that each limit is a size in bytes | review-owned: the GW5 review reads `projects/benchkit.md` | No check reads the table prose. |
| GW49 | 46 | The over-budget canary bites with a byte diagnostic | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetCanaryFixtureBites`) | A fixture left in lines raises a line diagnostic, which the byte `EXPECT` does not match. |
| GW50 | 47 | The craft-spec budget sentence names the current size | review-owned: the GW5 review reads `.agents/skills/bench-craft-spec/SKILL.md:68` | No anchor pins that sentence. |
| GW51 | 50 | `CHANGELOG.md` holds one section with one bullet for each chunk | review-owned: each chunk review reads the entry | `CHANGELOG.md` is excluded from prose grading, and no row pins a new entry. |
| GW53 | 45 | `CONTEXT.md` defines "prose budget" as a byte ceiling | review-owned: the GW5 review reads the "prose budget" entry of `CONTEXT.md` | No anchor pins the definition, so only review sees the old line-count wording. |

Not covered: story 48 — a reviewed exclusion that adds no behavior; GW01, GW16, and GW33 show that each removal has another copy.
Not covered: story 49 — a reviewed exclusion that adds no behavior; no ticket edits a `description:` line.

### Edge inventory

The profile's hostile-input checklist applies to the byte count and the table parse.
The read already refuses a symlink, a special file, an oversized file, and invalid UTF-8 (`prose_budget_test.go:81-98`), and existing tests pin each refusal (`:351-493`).
The absent-versus-empty pair holds for the skills root: an absent root and an unreadable root each raise their diagnostic (`:448-493`).
No test swaps a package variable.

- **Won't handle**: the git-ignored `.agents/skills/typesafe-ai/SKILL.md` in a primary checkout — the landing tree lacks it, and it exceeds the glob row before this outcome.
- **Won't handle**: a budget row for `build-rules.md` — reference files stay outside the budget universe (`prose_budget_test.go:35-37`), and its bytes leave every session.
- **Won't handle**: a subject with CRLF line endings — none exists, and a carriage return is a byte that the reader pays for.
- **Won't handle**: the 27 canary overlays that carry an old whole `BENCH.md` — each holds old text by design, and each still bites its own row.
- **Won't handle**: the line citations of `docs/research/roadmap-focus-and-completion.md:140` and `docs/research/aibadger-assessment.md:74`, `:219` — they are dated research records, not pointers.
- **Won't handle**: the audit records under `docs/audits/` that say "line budgets" — `.bench/prose-exclusions` excludes them as frozen records.
- **Won't handle**: a budget subject that grows in a later lane commit — the lane skips the budget check, but the landing gate reds it.

## Ownership fences

This fence is the union of the ticket Writes lines, which the plan-only preflight grades.
Each slicer extends it with its chunk, so the propose-writes preflight can grade each ticket during the slicing.
The coordinator owns the review pickup, which the coverage check requires in a declared fence section.

- `reviews/guidance-weight.md`
- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-drain.md`
- `.agents/commands/bench-final-check.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-gate/SKILL.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `.agents/skills/bench-implement-spec/references/build-rules.md`
- `.agents/skills/bench-implement-spec/references/staleness-pass.md`
- `.bench/BENCH-reference.md`
- `.bench/BENCH.md`
- `AGENTS.md`
- `CHANGELOG.md`
- `CONTEXT.md`
- `cmd/bench/anchor_help_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/testdata/anchors/fixture-repo/AGENTS.md`
- `docs/field-guide.html`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_check_floor_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_craft_research.go`
- `internal/anchors/registry_craft_research_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/commitment_guidance_test.go`
- `internal/conformance/implementation_continuation_test.go`
- `internal/conformance/prose_budget_test.go`
- `internal/conformance/recurrence_maintenance_contract_test.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/skill_description_budgets_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `internal/conformance/validity_checks_test.go`
- `projects/benchkit.md`
- `tests/canary/claude-agent-definitions/agent-unnamed-in-skill`
- `tests/canary/claude-agent-definitions/model-declared`
- `tests/canary/claude-agent-definitions/name-mismatch`
- `tests/canary/claude-agent-definitions/shell-tool-absent`
- `tests/canary/claude-agent-definitions/skill-names-missing-agent`
- `tests/canary/claude-agent-definitions/spawning-tool`
- `tests/canary/claude-agent-definitions/tools-absent`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/dogfood-referent-shipped`
- `tests/canary/docs-currency-token-diet/introduces-undeclared-command`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/docs-currency-token-diet/stale-command-reference`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/row-next-grammar/token-table-lacks-kit-edit`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/acceptance-coverage-anchor`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/agents-system-suite-route`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`
- `tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening`
- `tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion`
- `tests/canary/workflow-guidance-anchors/calibration-abstained-row`
- `tests/canary/workflow-guidance-anchors/calibration-abstention-final`
- `tests/canary/workflow-guidance-anchors/calibration-claim-schema`
- `tests/canary/workflow-guidance-anchors/calibration-coordinator-label`
- `tests/canary/workflow-guidance-anchors/calibration-delegate-pointer`
- `tests/canary/workflow-guidance-anchors/calibration-no-free-text`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-confidence`
- `tests/canary/workflow-guidance-anchors/calibration-pickup-step-move`
- `tests/canary/workflow-guidance-anchors/calibration-retro-aggregate-duty`
- `tests/canary/workflow-guidance-anchors/calibration-retro-table-duty`
- `tests/canary/workflow-guidance-anchors/calibration-status-definitions`
- `tests/canary/workflow-guidance-anchors/capture-sink-anchor`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/command-handoff-anchor`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/core-absent-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/core-absent-category-context`
- `tests/canary/workflow-guidance-anchors/core-absent-category-oracle`
- `tests/canary/workflow-guidance-anchors/core-absent-category-setup`
- `tests/canary/workflow-guidance-anchors/core-absent-category-work`
- `tests/canary/workflow-guidance-anchors/core-absent-gate-authority`
- `tests/canary/workflow-guidance-anchors/core-absent-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/core-absent-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/core-absent-retained-integration-source`
- `tests/canary/workflow-guidance-anchors/core-absent-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/core-absent-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/core-absent-skills-guidance`
- `tests/canary/workflow-guidance-anchors/core-absent-upgrade-route`
- `tests/canary/workflow-guidance-anchors/core-cli-bounded-complete-response`
- `tests/canary/workflow-guidance-anchors/core-cli-call-pair`
- `tests/canary/workflow-guidance-anchors/core-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/coverage-axis-anchor`
- `tests/canary/workflow-guidance-anchors/craft-gate-indirected-value-both-ends`
- `tests/canary/workflow-guidance-anchors/craft-gate-single-edit-defeat`
- `tests/canary/workflow-guidance-anchors/craft-spec-empty-or-invalid-fence`
- `tests/canary/workflow-guidance-anchors/craft-spec-exact-literal-fence`
- `tests/canary/workflow-guidance-anchors/craft-spec-fence-after-slice`
- `tests/canary/workflow-guidance-anchors/craft-spec-further-notes-bare`
- `tests/canary/workflow-guidance-anchors/craft-spec-reader-sweep-name`
- `tests/canary/workflow-guidance-anchors/craft-spec-red-signal-grammar-retired`
- `tests/canary/workflow-guidance-anchors/craft-spec-template-verification-log`
- `tests/canary/workflow-guidance-anchors/craft-tickets-slice-acceptance-row`
- `tests/canary/workflow-guidance-anchors/debug-archaeology-anchor`
- `tests/canary/workflow-guidance-anchors/debug-phase1-stop-gate-softened`
- `tests/canary/workflow-guidance-anchors/debug-red-commit`
- `tests/canary/workflow-guidance-anchors/debug-reproduction-economics-deleted`
- `tests/canary/workflow-guidance-anchors/decision-map-asset-path`
- `tests/canary/workflow-guidance-anchors/delegate-anchor-probe-owning-check`
- `tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package`
- `tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap`
- `tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge`
- `tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green`
- `tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer`
- `tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller`
- `tests/canary/workflow-guidance-anchors/delegate-grammar-fence-inventory`
- `tests/canary/workflow-guidance-anchors/delegate-live-tree-inventory-fence`
- `tests/canary/workflow-guidance-anchors/delegate-model-id-escalation`
- `tests/canary/workflow-guidance-anchors/delegate-out-of-fence-write`
- `tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface`
- `tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor`
- `tests/canary/workflow-guidance-anchors/delegate-probe-mutated-bytes`
- `tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance`
- `tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents`
- `tests/canary/workflow-guidance-anchors/delegate-root-conformance-command`
- `tests/canary/workflow-guidance-anchors/delegate-root-conformance-pass`
- `tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row`
- `tests/canary/workflow-guidance-anchors/delegate-serial-ceiling-fence`
- `tests/canary/workflow-guidance-anchors/delegate-skip-ownership-check`
- `tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor`
- `tests/canary/workflow-guidance-anchors/delegated-account-inventory`
- `tests/canary/workflow-guidance-anchors/delegated-account-reconciliation`
- `tests/canary/workflow-guidance-anchors/delegated-axis-exclusions`
- `tests/canary/workflow-guidance-anchors/delegated-cap-exhaustion-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review`
- `tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration`
- `tests/canary/workflow-guidance-anchors/delegated-entry-refusals`
- `tests/canary/workflow-guidance-anchors/delegated-final-verification`
- `tests/canary/workflow-guidance-anchors/delegated-green-predecessor`
- `tests/canary/workflow-guidance-anchors/delegated-lost-session-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-no-progress-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-opt-in-entry`
- `tests/canary/workflow-guidance-anchors/delegated-paid-comparison`
- `tests/canary/workflow-guidance-anchors/delegated-pending-predecessor-stop`
- `tests/canary/workflow-guidance-anchors/delegated-per-ticket-author`
- `tests/canary/workflow-guidance-anchors/delegated-prerequisite-checkpoint`
- `tests/canary/workflow-guidance-anchors/delegated-repair-ownership`
- `tests/canary/workflow-guidance-anchors/delegated-resumption-contents`
- `tests/canary/workflow-guidance-anchors/delegated-serial-ticket-checkpoint`
- `tests/canary/workflow-guidance-anchors/delegated-terminal-failure-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-writer-termination`
- `tests/canary/workflow-guidance-anchors/dg-1`
- `tests/canary/workflow-guidance-anchors/dg-1-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-10`
- `tests/canary/workflow-guidance-anchors/dg-11`
- `tests/canary/workflow-guidance-anchors/dg-12`
- `tests/canary/workflow-guidance-anchors/dg-13`
- `tests/canary/workflow-guidance-anchors/dg-14`
- `tests/canary/workflow-guidance-anchors/dg-15`
- `tests/canary/workflow-guidance-anchors/dg-15-executable-red`
- `tests/canary/workflow-guidance-anchors/dg-15-mandatory-red`
- `tests/canary/workflow-guidance-anchors/dg-15-mandatory-red-capitalized`
- `tests/canary/workflow-guidance-anchors/dg-15-mandatory-red-emphasis`
- `tests/canary/workflow-guidance-anchors/dg-17`
- `tests/canary/workflow-guidance-anchors/dg-18`
- `tests/canary/workflow-guidance-anchors/dg-19`
- `tests/canary/workflow-guidance-anchors/dg-2`
- `tests/canary/workflow-guidance-anchors/dg-2-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-20`
- `tests/canary/workflow-guidance-anchors/dg-21`
- `tests/canary/workflow-guidance-anchors/dg-22`
- `tests/canary/workflow-guidance-anchors/dg-23`
- `tests/canary/workflow-guidance-anchors/dg-25`
- `tests/canary/workflow-guidance-anchors/dg-26`
- `tests/canary/workflow-guidance-anchors/dg-29`
- `tests/canary/workflow-guidance-anchors/dg-29-verification-target`
- `tests/canary/workflow-guidance-anchors/dg-3`
- `tests/canary/workflow-guidance-anchors/dg-3-blanket-ban`
- `tests/canary/workflow-guidance-anchors/dg-3-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-30`
- `tests/canary/workflow-guidance-anchors/dg-31`
- `tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger`
- `tests/canary/workflow-guidance-anchors/dg-4`
- `tests/canary/workflow-guidance-anchors/dg-4-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-5`
- `tests/canary/workflow-guidance-anchors/dg-5-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-5-dirty`
- `tests/canary/workflow-guidance-anchors/dg-5-dirty-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-6`
- `tests/canary/workflow-guidance-anchors/dg-6-command`
- `tests/canary/workflow-guidance-anchors/dg-6-digest`
- `tests/canary/workflow-guidance-anchors/dg-6-dirty`
- `tests/canary/workflow-guidance-anchors/dg-6-owner-reference`
- `tests/canary/workflow-guidance-anchors/dg-6-surface`
- `tests/canary/workflow-guidance-anchors/dg-9`
- `tests/canary/workflow-guidance-anchors/drain-anchor`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-commit`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-per-spec-exception`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-route`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-row-fallback`
- `tests/canary/workflow-guidance-anchors/drain-implement-now-second-exception`
- `tests/canary/workflow-guidance-anchors/drain-roadmap-context-anchor`
- `tests/canary/workflow-guidance-anchors/drain-spec-history-anchor`
- `tests/canary/workflow-guidance-anchors/drain-split-board-detail-owner`
- `tests/canary/workflow-guidance-anchors/drain-split-board-retirement-pair`
- `tests/canary/workflow-guidance-anchors/drain-split-board-row-detail-owner`
- `tests/canary/workflow-guidance-anchors/edge-inventory-anchor`
- `tests/canary/workflow-guidance-anchors/field-guide-reduced-schema-map-card`
- `tests/canary/workflow-guidance-anchors/final-check-bare-leftover-clean-retired`
- `tests/canary/workflow-guidance-anchors/final-check-census-read-before-land`
- `tests/canary/workflow-guidance-anchors/final-check-landed-worktree-sweep`
- `tests/canary/workflow-guidance-anchors/final-check-light-path-changelog-heading`
- `tests/canary/workflow-guidance-anchors/final-check-scratch-branch-clean`
- `tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated`
- `tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness`
- `tests/canary/workflow-guidance-anchors/implement-spec-check-floor`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-planned-cell-rewrite`
- `tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/implementation-retro-drain-anchor`
- `tests/canary/workflow-guidance-anchors/light-path-inline`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `tests/canary/workflow-guidance-anchors/prepared-review-axis-returns`
- `tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff`
- `tests/canary/workflow-guidance-anchors/prepared-review-delegate-handoff-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points`
- `tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch`
- `tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability`
- `tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-triage-authority`
- `tests/canary/workflow-guidance-anchors/prepared-triage-bounds`
- `tests/canary/workflow-guidance-anchors/prepared-triage-cli-boundary`
- `tests/canary/workflow-guidance-anchors/prepared-triage-input-contract`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/review-base-merged-main-tip`
- `tests/canary/workflow-guidance-anchors/review-clean-terminal-result`
- `tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in`
- `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing`
- `tests/canary/workflow-guidance-anchors/review-falsification-dispositions`
- `tests/canary/workflow-guidance-anchors/review-persistence-anchor`
- `tests/canary/workflow-guidance-anchors/review-preflight-explicit-base`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner`
- `tests/canary/workflow-guidance-anchors/review-standing-falsification`
- `tests/canary/workflow-guidance-anchors/review-universal-claim-bar`
- `tests/canary/workflow-guidance-anchors/shared-worktree-path-pin`
- `tests/canary/workflow-guidance-anchors/slicing-in-write-spec`
- `tests/canary/workflow-guidance-anchors/spec-retire-roadmap-row`
- `tests/canary/workflow-guidance-anchors/staged-command-sweep-anchor`
- `tests/canary/workflow-guidance-anchors/story-line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-cross-pointers-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-decision-map-lifecycle-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-executable-route-pass`
- `tests/canary/workflow-guidance-anchors/ticket-light-path-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-light-path-carve-out`
- `tests/canary/workflow-guidance-anchors/ticket-lock-passes`
- `tests/canary/workflow-guidance-anchors/ticket-relocation-destinations`
- `tests/canary/workflow-guidance-anchors/ticket-seam-creating-slice`
- `tests/canary/workflow-guidance-anchors/ticket-skill-contract-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-source-clause-pass`
- `tests/canary/workflow-guidance-anchors/ticket-template-anchor`
- `tests/canary/workflow-guidance-anchors/write-spec-acceptance-round-exit`
- `tests/canary/workflow-guidance-anchors/write-spec-approval-table-emit`
- `tests/canary/workflow-guidance-anchors/write-spec-artifact-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-authorization-boundary`
- `tests/canary/workflow-guidance-anchors/write-spec-branch-split`
- `tests/canary/workflow-guidance-anchors/write-spec-contract-rederivation`
- `tests/canary/workflow-guidance-anchors/write-spec-conversation-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-conversation-fork`
- `tests/canary/workflow-guidance-anchors/write-spec-decision-source`
- `tests/canary/workflow-guidance-anchors/write-spec-degenerate-cheapness`
- `tests/canary/workflow-guidance-anchors/write-spec-fence-approval`
- `tests/canary/workflow-guidance-anchors/write-spec-frozen-base-and-tip-review`
- `tests/canary/workflow-guidance-anchors/write-spec-handoff-anchor`
- `tests/canary/workflow-guidance-anchors/write-spec-identified-coverage-default`
- `tests/canary/workflow-guidance-anchors/write-spec-late-uncertainty`
- `tests/canary/workflow-guidance-anchors/write-spec-map-sources`
- `tests/canary/workflow-guidance-anchors/write-spec-materiality-exit`
- `tests/canary/workflow-guidance-anchors/write-spec-ownership-fences`
- `tests/canary/workflow-guidance-anchors/write-spec-phase-ownership`
- `tests/canary/workflow-guidance-anchors/write-spec-post-slicing-handoff`
- `tests/canary/workflow-guidance-anchors/write-spec-promise-inflation-guard`
- `tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence`
- `tests/canary/workflow-guidance-anchors/write-spec-ready-map-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-review-made-conditional`
- `tests/canary/workflow-guidance-anchors/write-spec-review-tier-escalated`
- `tests/canary/workflow-guidance-anchors/write-spec-slicing-sign-off`
- `tests/canary/workflow-guidance-anchors/write-spec-ticket-approval-fields`
- `tests/canary/workflow-guidance-anchors/write-spec-ticket-breakdown-charge`
- `tests/canary/workflow-guidance-anchors/write-spec-unflagged-addition-removal`
- `tests/canary/workflow-guidance-anchors/write-spec-unique-row-id`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-learning`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-learning-contents`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-log`

## Out of scope

- A check that holds the kit's `CLAUDE.md` to the managed form: 2 edits, 1 gate run. The coordinator parks it with `bench idea` as a light-path ticket.
- The deletion of a rule with no other copy: 2 edits, 1 gate run for each rule, after an FT231 result. An example is `.bench/BENCH.md:133`.
- A trim of the description content phrases: about 20 edits, 1 gate run, after an FT231 trial of skill selection.
- The craft-skill library audit and the merge of overlapping craft skills: a later FT100 outcome with its own map.
- A lane rule that runs `guidance-prose-budgets` when a subject changes: 2 edits, 1 gate run.

## Further notes

Evidence status: source-backed proposal. Eight paired readers (Sonnet/high and Haiku/xhigh) read the tree at `8c8a6ca5`, and the pairs agreed on every load-bearing fact.
Remaining uncertainty: the red-to-green evidence of the build, and the overlaps below.

### Overlaps with staged specs

`specs/spec-stage-binding` (FT395, the committed outcome) ticket 3 writes `.bench/BENCH.md`, `AGENTS.md`, `bench-write-spec.md`, `commitment_guidance_test.go`, and the same anchor registries and canaries.
It inserts an eight-line block after the sign-off paragraph of `bench-write-spec.md` and pays for it with joins.
`specs/preflight-spec-staleness` ticket 05 writes `staleness-pass.md` and `bench-implement-spec.md`.
FT100 is not in the commitment, so its build follows both specs.
This spec cites the text that a landing can move by quotation, and the build-entry staleness pass re-cuts each citation.
The light-path ticket of decision #13 edits the descriptions of four command files and `.bench/BENCH-reference.md`, so the same rule applies to those files.

### Flagged additions

Each item goes beyond the literal decision source, and the reviewer can veto it at sign-off.
The consultant settled E1 to E12 on 2026-10-11, and the author chose E13 to E17.

- E1: the B10 exception moves into `staleness-pass.md` with its flag duty, because its kept copy lacks that duty (GW21, GW22). Decision #14 said the exception goes as a repeated copy.
- E2: the B12 row `registry_data.go:378` retires, because the `craft-line` copy keeps its own row (GW16, GW17).
- E3: the drain sentence "That approval is the reviewer's sign-off on this drain's proposed diff." and its row (GW35).
- E4: the craft-tickets carve-out names `/bench-write-spec`, and its row takes the text (GW34).
- E5: the light-path block spans four paragraphs (GW23, GW30). Paragraph 2 keeps the B40 lead word for word and adds the #4 route as two sentences (GW25, GW29).
- E6: the predicate lead-in counts its bullets, ending at "One predicate rides with them:" (GW14).
- E7: the five forbid rows of the build rules repoint to the new file (GW07, GW52).
- E8: 27 pointer lines repoint, not the 17 that decision #9 counted (GW08, GW09, GW12).
- E9: the GW4 ticket raises the `bench-write-spec.md` line row until the byte switch (GW37).
- E10: the glob row sizes from tracked files only (GW45).
- E11: "rounded up to the next multiple of 100" means the smallest multiple of 100 that is not less than the size (GW44).
- E12: the byte count counts every byte, the final newline included (GW39, GW41). This settles the first discretion item of the map. The map's discretion item on a `CLAUDE.md` check goes to Out of scope.
- E13: the new file is `build-rules.md`, its H2 sections, and its lead sentence (GW02, GW04). This settles the fourth discretion item of the map.
- E14: three new require rows pin the `.bench/BENCH.md` pointer, the lead of the build rules, and the moved B10 sentence (GW03, GW04, GW21).
- E15: each diagnostic that names the operating guide as the owner of a moved rule names the new owner, moved and unmoved rows alike (GW13).
- E16: each pointer names the build rules by its full path, on its own physical line (GW08).
- E17: the table prose names the byte unit, and craft-spec says "current size" (GW48, GW50).
- E18: the slicing research adds the `CONTEXT.md` "prose budget" entry to GW5, so the glossary names the byte unit (GW53).

### Source sentence to row

| source sentence | rows |
| --- | --- |
| #9: "B11 and B26–B33 move word for word into one new reference file beside the build phase" | GW01, GW02, GW15 |
| #9: "The 17 lines that point to these rules in `.bench/BENCH.md` point to the new file, and their anchors follow." | GW05, GW06, GW07, GW08, GW09, GW10, GW11, GW12, GW52 |
| #9: "`.bench/BENCH.md` keeps a one-line pointer." | GW03 |
| #14: "The B10 staleness-pass exception, B12, and A17 go from always-loaded guidance." | GW16, GW19, GW22 |
| #14: "Each anchor moves to the kept copy." and "the `AGENTS.md` anchor of A17 retires" | GW17, GW18, GW20, GW21 |
| #10: "One block under the light-path table. The table keeps the inline route for a change that the reviewer asks for." | GW23, GW24 |
| #10: "One paragraph under it states the commitment exemption, the delegated route of #4, the build hold, and the fix-don't-park split." | GW25, GW26, GW29, GW31 |
| #10: "The repeated drain sentences in fix-don't-park and in Capture go." | GW27, GW28 |
| #10: "Each anchored needle stays word for word, except the delegated-fix lead that #4 rewords." | GW29, GW30 |
| #15: "B43 moves beside the sign-off stop in `.agents/commands/bench-write-spec.md`, and craft-tickets points there." | GW32, GW33, GW34, GW37 |
| #15: "The line `.agents/commands/bench-drain.md:291` names the drain's own diff approval in its own words" | GW35, GW36 |
| #3: "Every row of the `### Guidance prose budgets` table in `projects/benchkit.md` limits its subject in bytes" | GW38, GW39, GW40, GW41, GW42, GW43, GW48, GW49 |
| #6: "Each exact row starts at its file's size in bytes when the measure changes, rounded up to the next multiple of 100 bytes." | GW44 |
| #6: "A glob row starts at the size of the largest file that only the glob row classifies, rounded up the same way." | GW45 |
| #7: "`AGENTS.md` gets an exact row in the budget table." | GW46 |
| #11: "The limit of `.bench/BENCH.md` is its size in bytes after the edits of this outcome." | GW47 |

### Pre-review proof checklist

- `Cited symbols`: `checkGuidanceProseBudgets`, `parseProseBudgetPolicy`, `proseBudgetLineCount`, `limitFor`, `proseBudgetSubjects`, `proseBudgetLines` (renamed `proseBudgetBytes`), and `healthyProseBudgetFiles` resolve in `internal/conformance/prose_budget_test.go`. `proseBudgetPolicy` and `profileBudgetTable` resolve there too, and `skill_description_budgets_test.go` reads them.
- `Import edges`: none. No ticket adds an import.
- `Source-row clauses and occurrences`: the source table above quotes each clause of the eight decisions.
- `Promised field labels`: the diagnostic "prose-budget exceeded: %s is %d bytes, over its %d-byte budget" and the error "not a positive byte count". The labels also include the six H2 sections of the build rules.
- `Changed-function callers`: `checkGuidanceProseBudgets` is called by the registry binding at `internal/conformance/checks_test.go:67` and by the tests of `prose_budget_test.go`. `parseProseBudgetPolicy` is called by `checkGuidanceProseBudgets`. `proseBudgetLineCount` goes, and its callers are `prose_budget_test.go:99` and its tests.
- `Copy survival`: GW01, GW16, GW22, GW27, GW28, and GW33 name the review search that fails if a removed copy survives.
- `Rendered-shape readers`: the old diagnostic needle "lines, over its" hits `prose_budget_test.go:100`, `:276`, `:284`, `:335`, `:506-508`, `:534`, and `tests/canary/guidance-prose-budgets/over-budget-skill/EXPECT`. Each hit joins the GW5 ticket.
- `Pin operators`: the budget check compares with `count > limit` (`prose_budget_test.go:99`). The commitment mirrors require `strings.Count(live, old) == 1` (`commitment_guidance_test.go:94`).
- `Entry reads`: none. The check reads the root that its harness passes.
- `Derived expectations`: each fixture size comes from the fixture bytes that the test writes, not from the check.
- `Consolidated rules`: the table below lists each consolidated light-path site.
- `Quantified obligations`: "each anchor" maps to GW05, GW07, GW09, GW17, GW52, GW21, GW29, and GW35. "Every pointer" maps to GW08 and GW09. "Every row" maps to GW44 to GW47.
- `Workflow-step writes`: none. The spec adds no workflow step.

| site | old rule | new rule | row |
| --- | --- | --- | --- |
| right-size table row | the inline light-path route | unchanged | GW24 |
| B40 lead | a learning fix ships at once only when the active outcome needs it | unchanged, plus the #4 route: a learning or drain fix goes to a write delegate, and a reviewer-asked change goes inline | GW25, GW29 |
| B42 second sentence | an unneeded defect goes to `bench learning`, and the drain delegates its fix | an unneeded defect goes to `bench learning` | GW27 |
| Capture paragraph | a drain implements a light-path idea under the light-path fix rule | removed | GW28 |

### Reader inventory

- The anchor registries that name `.bench/BENCH.md` by literal or constant: `registry_commitment.go`, `registry_data.go`, `registry_data_test.go`, `registry_retained_workflow.go`, `anchor_harness_diagnostics_test.go`, and `registry_ft311_review_dispatch.go`. `registry_chunk_chain.go` names the phase files by constant.
- The conformance readers of the moved text: `retained_workflow_test.go`, `implementation_continuation_test.go`, `commitment_guidance_test.go`, `recurrence_maintenance_contract_test.go`, and the shared-rule check `validity_checks_test.go:243-285`.
- The linked-repository surface: `.bench/consumer-payload.json` ships `.bench/BENCH.md` byte for byte and the `.agents/skills` tree, so `build-rules.md` ships with no manifest edit.
- The kit's `AGENTS.md` and `CLAUDE.md` do not ship (`internal/packagesurface/assets.go`).

### Not verified

- The exact physical lines that FT395 ticket 3 and the #13 light path leave in `bench-write-spec.md`, `bench-drain.md`, and `.bench/BENCH-reference.md`. The build-entry staleness pass re-cuts them.
- The diagnostic text of each moved row before its rename. The GW1 author reads each row.
- Whether `bench test` grades the working tree or a sealed export for the ignored `typesafe-ai` skill.

### Validation plan

The GW1 to GW4 tickets each run the named package tests, `bench test --check docs-currency-workflow`, `prose-mechanics`, and `guidance-prose-budgets`.
Each repointed mirror and each new row records one observed red through `bench probe` before its ticket commits.
GW5 records the red of each byte-boundary fixture before its production edit.
The landing gate grades the whole tree after the last chunk.
