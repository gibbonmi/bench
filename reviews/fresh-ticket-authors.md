# Fresh ticket authors review record

## FA ticket 1 author evidence

Author session: `claude:bench-writer/fta-t1-author`. Line: opus / high / cap 3 attempts. Source tip at start: `42fcfb131425e074a2bf8b536ad00767bcfe87bc`. Ticket commit: `41f19bce59bd91357c7c1a94672ed1577ba538c1`.

The author added the FA rows to the anchor registries before the guidance edit. Then the author ran the rows against the unchanged guidance, and each row went red. After the guidance edit, each row went green.

### Done claims

The coordinator fills the label cell after its own probe.

| Row | Status | Confidence | Label |
|---|---|---|---|
| FA1 | verified | 9 | |
| FA2 | verified | 9 | |
| FA3 | verified | 9 | |
| FA4 | verified | 9 | |
| FA5 | verified | 9 | |
| FA6 | verified | 9 | |
| FA7 | verified | 9 | |
| FA8 | verified | 9 | |
| FA9 | verified | 9 | |
| FA10 | verified | 9 | |
| FA11 | verified | 9 | |
| FA12 | verified | 9 | |
| FA13 | verified | 9 | |
| FA14 | verified | 9 | |

### Red then green

The red column quotes the diagnostic from `bench test --full --package ./internal/conformance --run TestFreshTicketAuthors` on the unchanged guidance. The FA6 red comes from `TestEvidenceBoundedActionSentencesArePinned/build` on the unchanged guidance. The green column is the same test after the guidance edit.

| Row | Seam | Red before the edit | Green after the edit |
|---|---|---|---|
| FA1 | Require, `registry_ft311_review_dispatch.go` | operating guide dropped the fresh author session for each ticket on the declared line | pass |
| FA2 | Require, `registry_ft311_review_dispatch.go` | operating guide dropped the limit-1 version 2 delegate plan | pass |
| FA3 | Require, `registry_retained_workflow.go` | operating guide dropped the concurrent authors and the full tier range of `--delegate` | pass |
| FA4 | Forbid, `registry_ft311_review_dispatch.go` | operating guide restored the retained implementation session | pass |
| FA5 | Require, `registry_ft311_preparation.go` | `bench-implement-spec.md` dropped the narrow author charge | pass |
| FA6 | `TestEvidenceBoundedActionGuidance`, build delivery `a narrow author read` | the verified-delivery sentence is pinned by 0 Require rows | pass |
| FA7 | Forbid, `registry_ft311_preparation.go` | `bench-implement-spec.md` restored the full evidence retrieval | pass |
| FA8 | Forbid, `registry_ft311_preparation.go` | `bench-implement-spec.md` restored the verified-delivery build rule | pass |
| FA9 | Require, `registry_ft311_preparation.go` | `bench-implement-spec.md` dropped the orchestrator's read bound | pass |
| FA10 | Require, `registry_ft311_preparation.go` | `bench-implement-spec.md` dropped the handoff refresh at each chunk checkpoint | pass |
| FA11 | Require, `registry_ft311_review_dispatch.go` | operating guide dropped the orchestrator's final reconciliation | pass |
| FA12 | Require, `registry_retained_workflow.go` | operating guide dropped the fresh repair session for each affected ticket | pass |
| FA13 | Forbid, `registry_retained_workflow.go` | operating guide restored the ticket-author repair ownership | pass |
| FA14 | Require, `registry_retained_workflow.go` | reference dropped the version 2 plan amendment or its verification for each ticket | pass |

### Probe verdicts

Each probe ran through `bench probe` after `bench worktree build`. Each restore reads `yes`.

| Row | File | Mutation | Check | Verdict | Diagnostic |
|---|---|---|---|---|---|
| FA4 | `.bench/BENCH.md` | swap: the retired authorship sentence returns | `docs-currency-workflow` | bit | operating guide restored the retained implementation session |
| FA7 | `.agents/commands/bench-implement-spec.md` | swap: the retired full-retrieval sentence returns | `docs-currency-workflow` | bit | restored the full evidence retrieval |
| FA8 | `.agents/commands/bench-implement-spec.md` | swap: the retired verified-delivery sentence returns | `docs-currency-workflow` | bit | restored the verified-delivery build rule |
| FA13 | `.bench/BENCH.md` | swap: the retired repair sentence returns | `docs-currency-workflow` | bit | operating guide restored the ticket-author repair ownership |
| FA6 | `internal/conformance/charge_evidence_guidance_test.go` | swap: the build delivery returns to `verified delivery` | `TestEvidenceBoundedActionGuidance` | bit | no registry row states the verified-delivery diagnostic |
| FA1 | `.bench/BENCH.md` | omission: the fresh author sentence leaves | `docs-currency-workflow` | bit | operating guide dropped the fresh author session for each ticket on the declared line |

Each FA4, FA7, FA8, and FA13 probe adds the retired sentence and keeps every other needle. So only the Forbid row can make the check red.

### Verification

| Check | Verdict | Elapsed |
|---|---|---|
| `bench test --check docs-currency-workflow` | pass | 781 ms |
| `bench test --package ./internal/conformance --run TestRootConformance` | pass | 4954 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 13119 ms |
| `bench test --package ./internal/...` | pass, 11 capability skips | 36726 ms for the conformance package |
| `bench gate-prose . -- <three edited Markdown files>` | pass | not reported |
| `go vet ./...` | pass | not reported |
| `bench commit` lane | pass | not reported |

### Flags for review

- The ticket text puts the verification split in the build phase too. The build phase points to `.bench/BENCH-reference.md` for the amendment contents, so that fact has one source.
- The plan-expansion paragraph of `.bench/BENCH.md` now names the orchestrator, not the retained author.
- The `.bench/BENCH.md` line budget is 185 lines. To stay in budget, the author joined the wrapped landing-shape sentence onto one line.
- The new rows use the diagnostic prefix `fresh ticket author: `. `TestFreshTicketAuthors` makes each row bite and reads the live guidance.

## FA ticket 2 author evidence

Author session: `claude:bench-writer/fta-t2-author`. Line: opus / high / cap 3 attempts. Source tip at start: `48ac4380df5fd2a46e1cd662f70678d0f72ef0af`. Ticket commit: `5950b823`.

The author added the FA rows to the anchor registries before the skill edit. Then the author ran `bench anchors` on each unchanged skill file, and each row went red. After the skill edit, each row went green.

### Done claims

The coordinator fills the label cell after its own probe.

| Row | Status | Confidence | Label |
|---|---|---|---|
| FA15 | verified | 9 | |
| FA16 | verified | 9 | |
| FA17 | verified | 9 | |
| FA18 | verified | 9 | |
| FA19 | verified | 9 | |
| FA26 | verified | 9 | |

### Red then green

The red column quotes the diagnostic from `bench anchors <skill file>` on the unchanged skill, after `bench worktree build`. The green column is `bench test --check docs-currency-workflow` and `TestFreshTicketAuthors` after the skill edit.

| Row | Seam | Red before the edit | Green after the edit |
|---|---|---|---|
| FA15 | RequireInSection "Retained implementation continuation", `registry_retained_workflow.go` | craft-line dropped the continuation rules for each ticket author | pass |
| FA16 | Require, `registry_retained_workflow.go` | craft-line dropped the reviewer stop before a tier move of a fresh author | pass |
| FA17 | Require, `registry_data.go` | craft-delegate dropped the owner pointer for ticket author sessions | pass |
| FA18 | Require, `registry_data.go` | craft-tickets dropped the ticket size of one fresh author context | pass |
| FA19 | Forbid, `registry_data.go` | craft-tickets restored the retained frontier session | pass |
| FA26 | Forbid, `registry_retained_workflow.go` | craft-delegate restored the repair return to the retained session | pass |

### Probe verdicts

Each probe ran through `bench probe` after `bench worktree build`. Each restore reads `yes`.

| Row | File | Mutation | Check | Verdict | Diagnostic |
|---|---|---|---|---|---|
| FA19 | `.agents/skills/bench-craft-tickets/SKILL.md` | swap: the retired frontier sentence returns | `docs-currency-workflow` | bit | craft-tickets restored the retained frontier session |
| FA26 | `.agents/skills/bench-craft-delegate/SKILL.md` | swap: the retired repair sentence returns | `docs-currency-workflow` | bit | craft-delegate restored the repair return to the retained session |
| FA16 | `.agents/skills/bench-craft-line/SKILL.md` | omission: the tier-move sentence leaves | `docs-currency-workflow` | bit | craft-line dropped the reviewer stop before a tier move of a fresh author |

Each FA19 and FA26 probe adds the retired sentence and keeps every other needle. So only the Forbid row can make the check red.

### Verification

| Check | Verdict | Elapsed |
|---|---|---|
| `bench test --check docs-currency-workflow` | pass | 940 ms |
| `bench test --package ./internal/conformance --run TestRootConformance` | pass | 7696 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 18398 ms |
| `bench test --check guidance-prose-budgets` | pass | 5 ms |
| `bench test --package ./internal/anchors` | pass | 883 ms |
| `bench gate-prose . -- <four edited Markdown files>` | pass | not reported |
| `go vet ./...` | pass | not reported |
| `bench commit` lane | pass | not reported |

### Flags for review

- The rewritten frontier sentence keeps a Require row beside the FA19 Forbid row. That Require row is in `registry_retained_workflow.go`, not in `registry_data.go`, because `registry_data.go` is over its structure budget and must not grow.
- The `delegated-per-ticket-author` canary now mutates the fresh author sentence of craft-delegate, and its diagnostic moved to the `fresh ticket author: ` family.
- The craft-line continuation rows keep their needles in the "Retained implementation continuation" section, but the needles now name the ticket author. Three continuation diagnostics dropped the word "retained".
- The new `user-directed` transfer entry in the delegation discipline has no anchor row, because the spec plans none.
- Retained-author wording remains outside this fence in `.agents/commands/bench-review-implementation.md`: "A clean chunk review hands its frozen pair back to the retained author."

## FA ticket 3 author evidence

Author session: `claude:bench-writer/fta-t3-author`. Line: opus / high / cap 3 attempts. Source tip at start: `164a3c689958d02e9600e4bc7a8f0eabe79697a9`. Ticket commit: `0b947924`.

The author added the FA rows to the anchor registries before the guidance edit. Then the author edited one file at a time and ran `TestFreshTicketAuthors` after each edit. Each row went red on the unchanged file and green after its edit.

### Done claims

The coordinator fills the label cell after its own probe.

| Row | Status | Confidence | Label |
|---|---|---|---|
| FA20 | verified | 9 | |
| FA21 | verified | 9 | |
| FA22 | verified | 9 | |
| FA23 | verified | 9 | |
| FA24 | claimed | 8 | |
| FA25 | verified | 9 | |

### Red then green

The red column quotes the first diagnostic of `TestFreshTicketAuthors` on the unchanged file. The green column is `bench test --check docs-currency-workflow` and `TestFreshTicketAuthors` after the edit.

| Row | Seam | Red before the edit | Green after the edit |
|---|---|---|---|
| FA20 | Require, `registry_data.go` | spec authoring dropped the line for fresh ticket authors on one integration source | pass |
| FA21 | Require, `registry_ft311_review_dispatch.go` | review phase dropped the fresh repair author for accepted findings | pass |
| FA22 | Require, `registry_data.go` | field guide dropped the fresh author session for each ticket | pass |
| FA23 | Require, `registry_retained_workflow.go` | README dropped the fresh author session for each ticket | pass |

FA24 is review-owned and has no anchor. FA25 keeps its existing Require row byte for byte, and this commit does not change `.agents/commands/bench-drain.md`.

### Probe verdicts

Each probe ran through `bench probe` after `bench worktree build`. Each restore reads `yes`.

| Row | File | Mutation | Check | Verdict | Diagnostic |
|---|---|---|---|---|---|
| FA20 | `.agents/commands/bench-write-spec.md` | swap: the retired retained-session sentence returns | `docs-currency-workflow` | bit | spec authoring dropped the line for fresh ticket authors on one integration source |
| FA21 | `.agents/commands/bench-review-implementation.md` | swap: the retired retained-session return returns | `docs-currency-workflow` | bit | review phase dropped the fresh repair author for accepted findings |
| FA22 | `docs/field-guide.html` | swap: the retired retained-authorship sentence returns | `docs-currency-workflow` | bit | field guide dropped the fresh author session for each ticket |
| FA23 | `README.md` | omission: the fresh author sentence leaves | `docs-currency-workflow` | bit | README dropped the fresh author session for each ticket |

### Verification

| Check | Verdict | Elapsed |
|---|---|---|
| `bench test --check docs-currency-workflow` | pass | 1027 ms |
| `bench test --package ./internal/conformance --run TestRootConformance` | pass | 7844 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 17805 ms |
| `bench test --check guidance-prose-budgets` | pass | 8 ms |
| `bench test --package ./internal/anchors` | pass | 921 ms |
| `bench gate-prose . -- <six edited Markdown files>` | pass | not reported |
| `go vet ./...` | pass | not reported |
| `bench commit` lane | pass | not reported |

### Flags for review

- The FA23 row is in `registry_retained_workflow.go`, not in `registry_data.go`, because `registry_data.go` is over its structure budget and must not grow. FA20 and FA22 replace the old rows in place, so `registry_data.go` keeps its length.
- The FA20 row replaces the `workflow integration source: ` write-spec row. So the integration-source family count in `docs_workflow_helpers_test.go` goes from 10 to 9, and its stale-claim scan no longer reads `.agents/commands/bench-write-spec.md`.
- The FA22 row replaces the `retained workflow: ` field-guide row, and its two entries in `retained_workflow_test.go` leave.
- The review convergence contract in `docs_workflow_helpers_test.go` now requires "the orchestrator reconciles overall acceptance and integration before landing".
- The `write-spec-frozen-base-and-tip-review` canary copy now carries the FA20 sentence, so its mutation still reds only its own row.
- The author also changed sentences that no row pins. These are the write-spec authorship sentence, two field-guide sentences, and the README repair node. They also include the final-check retro item, two review-phase sentences, and one registry comment.
- ADR 0021 names its first outcome "ticket implementation", not "retained implementation".

## FA chunk review, round 1

The frozen pair is base `d23694e925c6e0ffd15d6eda355daef4497acbea` and tip `202a9196038e43a652147a59ec30aef9406ff8be`. The shared evidence is `sha256:a793bbcc5ef4271fa7a68cefd0ac1557ca3240c566b18917f7e4ab4a17a6f373`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort. Only the Coverage axis ran probes on the shared tree, and it left the tree clean.

The raw finding count is 22: Standards 9, Spec 7, and Coverage 6. After the orchestrator merges the findings that name the same fix, 13 repair targets remain. The reviewer decided three questions on 2026-09-24:

- Add Forbid rows for the three retired sentences that Coverage named.
- Fix the retired wording in `projects/benchkit.md` and `CONTEXT.md` inside the ticket 3 repair.
- Treat the restatements that approved rows require as `no-op`.

## Standards

Findings: 9. The worst issue is the tier-move stop in five places.

- `.agents/skills/bench-craft-delegate/SKILL.md:19` restates the routing rule of `.bench/BENCH.md`. The ticket 2 text and a planned needle require the sentence. `no-op` by reviewer decision. Confidence 6.
- `.agents/skills/bench-craft-tickets/SKILL.md:93` paraphrases the routing rule. FA18 requires the sentence. `no-op` by reviewer decision. Confidence 5.
- `.bench/BENCH.md:122` repeats the tier-move stop that `craft-line` owns. FA1 and FA16 require both copies. `no-op` by reviewer decision. Confidence 6.
- `.agents/commands/bench-implement-spec.md:26` restates the version 2 `execution` block before its pointer to `.bench/BENCH-reference.md`. Target R4. `auto-fix`. Confidence 5.
- `docs/adr/0021-benchmark-workflow-orchestration.md:9-13` restates facts that ADR 0023 owns. Target R12. `auto-fix`. Confidence 5.
- `internal/anchors/registry_retained_workflow.go:38` holds the FA23 row, but the FA23 coverage row names `registry_data.go`. The orchestrator amends the coverage row. Target R14. `auto-fix`. Confidence 6.
- `.agents/commands/bench-implement-spec.md:24,45` keeps the full-retrieval words beside the narrow author read. Target R5. `auto-fix`. Confidence 5.
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:99` can read as a replacement of the repair policy itself. Target R8. `auto-fix`. Confidence 6.
- `.bench/BENCH.md:122` keeps the "retained implementation continuation" section name, as FA15 requires. `no-op`. Confidence 6.

## Spec

Findings: 7. The worst issue is `docs/field-guide.html:1325`, which still states the retired single-author rule.

- `docs/field-guide.html:1325` states "the approved implementation line retains one author through production changes, tests, probes, repairs, and chunk reviews." Target R9. `auto-fix`, blocking. Confidence 9.
- `projects/benchkit.md:590` states that one session retains the build. Target R13. `ask-user`, and the reviewer approved the fix. Confidence 8.
- `CONTEXT.md:357` limits the orchestrator to a delegated run. Target R13. `ask-user`, and the reviewer approved the fix. Confidence 7.
- `internal/anchors/registry_data.go:348` with `internal/conformance/docs_workflow_helpers_test.go:524`: the integration-source scan no longer reads `bench-write-spec.md`. Target R10. `auto-fix`, because the repair restores the earlier guarantee. Confidence 7.
- `.agents/commands/bench-implement-spec.md:42` keeps "A fresh consumer runs its own retrieval". Target R5. `auto-fix`. Confidence 5.
- `docs/adr/0021-benchmark-workflow-orchestration.md:9` states the tier-move stop with no `--delegate` condition. Target R12. `auto-fix`. Confidence 5.
- `.bench/BENCH.md:134` gives no `--delegate` scope to concurrent chunk authors. The limit-1 plan already stops concurrent authors. `no-op`. Confidence 5.

## Coverage

Findings: 6. The worst issue is the FA10 needle, which a negation does not break.

- `internal/anchors/registry_ft311_preparation.go:542`: the FA10 needle stays green after "never refreshes". Probe P1 was silent. Target R1. `auto-fix`. Confidence 9.
- `internal/conformance/docs_workflow_helpers_test.go:498-524`: the integration-source family no longer scans `bench-write-spec.md`. Probe P3 was silent. Target R10. `auto-fix`. Confidence 7.
- `.bench/BENCH.md` and `.agents/commands/bench-review-implementation.md`: three retired sentences have no Forbid row. Targets R3 and R11. `ask-user`, and the reviewer approved the rows. Confidence 8.
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:99`: the `user-directed` transfer entry has no anchor row. Target R7. `auto-fix`. Confidence 9.
- `.bench/BENCH.md:128`: the material-shortfall sentence after the FA12 needle has no pin. Target R2. `auto-fix`. Confidence 8.
- `.agents/commands/bench-implement-spec.md:28,42`: the consumer retrieval sentences disagree with the narrow author read. Target R5. `auto-fix`. Confidence 6.

## FA repair routing

Each affected ticket takes one fresh repair session. The session reruns that ticket's verification.

| ticket | targets |
| --- | --- |
| `1-route-each-ticket-to-a-fresh-author.md` | R1, R2, R3, R4, R5 |
| `2-align-the-line-and-delegation-skills.md` | R7, R8 |
| `3-align-the-phase-commands-and-docs.md` | R9, R10, R11, R12, R13 |
| orchestrator spec amendment | R14, and the coverage rows for R2, R3, R7, and R11 |

No repair edits `.bench/BENCH.md`, so no repair needs the reviewer's permission rule.

## FA ticket 1 repair evidence

The fresh repair session `claude:bench-writer/fta-t1-repair` ran on opus at high effort with a cap of 3 attempts. It started from tip `b2f800ac66ce53301521b861a2e0d96d7ba23998` and bound evidence `sha256:2b325c99fe948c59b102e2ccfcb6b79f106688dda8d4bc39fcd6b45d434ae713` as current. The repair commit is `da89b791dee350574c10c6126a4df1ea0d656fc8`. This repair is the first repair cycle of chunk FA. It used one attempt.

| target | row | change | status |
| --- | --- | --- | --- |
| R1 | FA10 | The needle now starts at the read bound: "not code, and it refreshes `bench handoff` at each chunk checkpoint." | done |
| R2 | FA12 | The FA12 Require needle now also holds "A finding on a path that no `Writes:` line holds is a material acceptance shortfall." | done |
| R3 | FA27, FA28 | Two Forbid rows over `.bench/BENCH.md` in `registry_ft311_review_dispatch.go`. The family runner in `TestFreshTicketAuthors` bites each row. | done |
| R4 | none | The build phase keeps only the pointer to the plan amendment in `.bench/BENCH-reference.md`. | done |
| R5 | none | The context sentence and the consumer sentences now name the narrow author read. The guarantee stays: a receipt, a cursor, or another consumer's delivery never replaces what this session reads itself. | done |

R3 adds no test code. `TestFreshTicketAuthors` gives each registered row of the family a synthetic bite, so a new row joins that test through its diagnostic prefix.

### Red then green

Each probe ran through `bench probe <file> ... --check docs-currency-workflow`. Each row shows the verdict before and after the repair.

| probe | file | kind | before | after | diagnostic after |
| --- | --- | --- | --- | --- | --- |
| R1 negation: "and it refreshes" becomes "and it never refreshes" | `.agents/commands/bench-implement-spec.md` | swap | silent | bit | handoff refresh at each chunk checkpoint |
| R2 omission of the material-shortfall sentence | `.bench/BENCH.md` | omit | silent | bit | fresh repair session for each affected ticket |
| FA27 restore of the retired sentence before "Repeat delegated review" | `.bench/BENCH.md` | swap | silent | bit | retained author's final reconciliation |
| FA28 restore of the retired sentence before "After the last chunk" | `.bench/BENCH.md` | swap | silent | bit | finding return to the retained author |

Each probe reported one failure and restored the file. No Require needle fired in the FA27 and FA28 probes.

### Verification

| command | exit | result |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestRootConformance` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run "TestFreshTicketAuthors\|TestEvidence\|TestRetainedWorkflow"` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/anchors/...` | 0 | pass, 0 failures, 0 skips |
| `bench test --check guidance-prose-budgets` | 0 | pass, 0 failures, 0 skips |
| `bench gate-prose . -- .agents/commands/bench-implement-spec.md` | 0 | pass |
| `go vet ./...` | 0 | no output |

The first `TestRootConformance` run failed on a diff-owned red: the reworded context sentence had 27 words, over the 25-word bound. The repair shortened the sentence and its pin, and the rerun passed.

## FA ticket 2 repair evidence

The fresh repair session `claude:bench-writer/fta-t2-repair` ran on opus at high effort with a cap of 3 attempts. It started from tip `80be05150e8a69caee2d9770d2868044edda20bd` and bound evidence `sha256:9e83978c431cec2afb37570cad37ed38940fc29acab3b7e87578db50c4f5bb85` as current. The repair commit is `d6612e4df8de1854d8d25bd6e4167e49c456304f`. This repair is part of the first repair cycle of chunk FA. It used one attempt.

| target | row | change | status |
| --- | --- | --- | --- |
| R7 | FA30 | A `RequireInSection` row in `registry_retained_workflow.go` pins the `user-directed` entry under "Delegated author transfer", with the prefix `fresh ticket author: `. | done |
| R8 | FA30 | The entry now reads "A post-review repair under the standing policy of `.bench/BENCH.md` permits a `user-directed` replacement of the author by a fresh repair session." This sentence is the R7 needle. | done |

R7 adds no test code. `TestFreshTicketAuthors` gives each registered row of the family a synthetic bite, so the new row joins that test through its diagnostic prefix.

### Red then green

The R7 row was added before the R8 rewording. `bench test --check docs-currency-workflow` then failed with one failure: "fresh ticket author: delegation discipline dropped the user-directed transfer for a post-review repair". After the rewording, the same check passed.

Each probe ran through `bench probe .agents/skills/bench-craft-delegate/references/delegation-discipline.md ... --check docs-currency-workflow` after the repair.

| probe | kind | verdict | diagnostic |
| --- | --- | --- | --- |
| "permits a `user-directed` replacement" becomes "never permits a `user-directed` replacement" | swap | bit | user-directed transfer for a post-review repair |
| omission of the whole entry | omit | bit | user-directed transfer for a post-review repair |

Each probe reported one failure and restored the file.

### Verification

| command | exit | result |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestRootConformance` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestFreshTicketAuthors` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/anchors/...` | 0 | pass, 0 failures, 0 skips |
| `bench test --check guidance-prose-budgets` | 0 | pass, 0 failures, 0 skips |
| `bench gate-prose . -- .agents/skills/bench-craft-delegate/references/delegation-discipline.md` | 0 | pass |
| `go vet ./...` | 0 | no output |

## FA ticket 3 repair evidence

The fresh repair session `claude:bench-writer/fta-t3-repair` ran on opus at high effort with a cap of 3 attempts. It started from tip `32285a56f2b809d8d3510b4649b1c73076f702d5` and bound evidence `sha256:403fe1a16decce1f1ce4602d3b249b04b4cc5df5956cea54dbc03059a33df7f9` as current. The repair commit is `05bbe62625756bacf7f5b35fd7187150ded28670`. This repair is part of the first repair cycle of chunk FA. It used one attempt.

| target | row | change | status |
| --- | --- | --- | --- |
| R9 | FA22 | The field-guide callout now states that each ticket gets a fresh author session, that a fresh repair session takes each accepted finding, and that the orchestrator reconciles before the landing. It points to `.bench/BENCH.md`. No other single-author sentence remains on the page. | done |
| R10 | none | The write-spec reconciliation row in `registry_data.go` now uses the prefix `workflow integration source: `. The family count goes back to 10, and the stale-claim scan reads `.agents/commands/bench-write-spec.md` again. The canary `write-spec-frozen-base-and-tip-review` expects the new diagnostic. `registry_data.go` does not grow. | done |
| R11 | FA29 | A Forbid row over `.agents/commands/bench-review-implementation.md` in `registry_ft311_review_dispatch.go`, with the prefix `fresh ticket author: `. | done |
| R12 | FA24 | ADR 0021 consequences 1 and 5 now point to ADR 0023 for ticket authorship, repair authorship, and the additions of a delegated run. They restate no ADR 0023 fact. | done |
| R13 | none | `projects/benchkit.md` and `CONTEXT.md` now state that every spec-backed build has an orchestrator and a fresh author for each ticket. The bytes "one retained integration source" do not change. No canary copy of `projects/benchkit.md` holds the changed sentence, so no canary changes. | done |

R11 adds no test code. `TestFreshTicketAuthors` gives each registered row of the family a synthetic bite, so the new row joins that test through its diagnostic prefix.

### Red then green

Each probe ran through `bench probe <file> ... --check docs-currency-workflow` after `bench worktree build`. The R10 before probe ran with the old diagnostic put back in the tree for the probe only.

| probe | file | kind | before | after | diagnostic after |
| --- | --- | --- | --- | --- | --- |
| R10 insertion of "Landing is the sole landing path." before "Then recommend the approved implementation line" | `.agents/commands/bench-write-spec.md` | swap | silent | bit | retains stale scalar or sole-path workflow claim "sole landing path" |
| FA29 insertion of the retired sentence after the FA21 Require needle | `.agents/commands/bench-review-implementation.md` | swap | not run; no Forbid row existed | bit | review phase restored the finding return to the retained session |
| self-probe: omission of the FA22 sentence | `docs/field-guide.html` | omit | not applicable | bit | field guide dropped the fresh author session for each ticket |

Each probe reported one failure and restored the file. The FA29 probe keeps every Require needle, so only the Forbid row can make the check red. No row pins the new R9 sentence, so the self-probe omits the FA22 sentence.

### Verification

| command | exit | result |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestRootConformance` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | 0 | pass, 0 failures, 0 skips |
| `bench test --full --package ./internal/conformance --run "TestFreshTicketAuthors\|TestIntegrationSourceWorkflowAnchorsBiteIndependently"` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance` | 0 | pass, 0 failures, 3 capability skips |
| `bench test --package ./internal/anchors/...` | 0 | pass, 0 failures, 0 skips |
| `bench test --check guidance-prose-budgets` | 0 | pass, 0 failures, 0 skips |
| `bench gate-prose . -- <three edited Markdown files>` | 0 | pass |
| `go vet ./...` | 0 | no output |
| `bench commit` lane | 0 | pass |

### Flags for review

- ADR 0021 no longer states "A model change requires reviewer direction." ADR 0023 owns the tier-move stop, and `craft-line` owns a change of model or session.
- The `CONTEXT.md` entry for "ticket author" still lists repairs among the author's work. A fresh repair session is a new author assignment for its ticket, so the entry stays.

## FA chunk review, round 2

This round confirms the first repair cycle. The frozen pair is base `d23694e925c6e0ffd15d6eda355daef4497acbea` and tip `b4ed0aa96f0de303821b45f7329755ca4d7a0616`, and the axes read the repair delta from `202a9196038e43a652147a59ec30aef9406ff8be`. The shared evidence is `sha256:87fd306037efcf863a59b90479ddac3f720f5e5e36ae77b0495efc796984f62a`. Each axis ran in a new `bench-reviewer` session on opus at high effort. Every fold of the first cycle landed.

The raw finding count is 2: Standards 1, Spec 0, and Coverage 1. Both findings name the same callout, so one repair target remains. Chunk FA has used one of its two repair cycles, and this target takes the second cycle. The reviewer chose the repair on 2026-09-24: trim the callout, and add one Forbid row for the retired field-guide sentence.

### Standards

Findings: 1. The worst issue is the field-guide callout.

- `docs/field-guide.html:1325-1328` restates the rule that line 1140 holds, and its tier-move clause has no `--delegate` condition. Target R15. `auto-fix`. Confidence 6.

### Spec

Findings: 0. Every fold is confirmed, and the three repair assignments are valid replacements.

### Coverage

Findings: 1. The worst issue is the unpinned field-guide callout.

- `docs/field-guide.html:1325` with `internal/anchors/registry_data.go:357`: the retired single-author sentence can return beside a green FA22 needle. Probe P1 was silent. Target R16. `ask-user`, and the reviewer approved one Forbid row. Confidence 8.

### Advice

- Spec: the plan names each repair's start as the tip before its assignment commit, and the evidence sections name the assignment commit.
- Standards: `CONTEXT.md:359` has a short reflowed line.
- Coverage: the FA10 needle does not stop a separate negating sentence after it.

### Repair routing

| ticket | targets |
| --- | --- |
| `3-align-the-phase-commands-and-docs.md` | R15, R16 |
| orchestrator spec amendment | row FA31 for R16 |

## FA ticket 3 repair evidence, cycle 2

The fresh repair session `claude:bench-writer/fta-t3-repair-2` ran on opus at high effort with a cap of 3 attempts. It started from tip `cfc68d9805f95d2ef100e8f74d1581b7038fbc11` and bound evidence `sha256:3518dfd8d5d21b47bad68bca37cbf7aa8e12f45bb26fea927e2026dd2d5fb568` as current. The repair commit is `e942e63a7fde9d1a7dcbe70d605d4cbd14aa90ec`. This repair is the second repair cycle of chunk FA, which is cycle 2 of 2. It used one attempt.

| target | row | change | status |
| --- | --- | --- | --- |
| R15 | FA22 | The field-guide callout now states only that `.bench/BENCH.md` owns ticket and repair authorship after ticket approval. It restates no fact of the "How it works" card or `.bench/BENCH.md`, and it has no tier-move clause. The FA22 sentence in the "How it works" card does not change. | done |
| R16 | FA31 | A Forbid row over `docs/field-guide.html` in `registry_ft311_review_dispatch.go`, with the needle "the approved implementation line retains one author through production changes, tests, probes, repairs, and chunk reviews." and the prefix `fresh ticket author: `. | done |

R16 adds no test code. `TestFreshTicketAuthors` gives each registered row of the family a synthetic bite, so the new row joins that test through its diagnostic prefix.

### Red then green

Each probe ran through `bench probe docs/field-guide.html ... --check docs-currency-workflow` after `bench worktree build`. The before probe ran on the trimmed callout before the Forbid row existed.

| probe | file | kind | before | after | diagnostic after |
| --- | --- | --- | --- | --- | --- |
| FA31 insertion of the retired sentence after the new pointer sentence, with the FA22 sentence in place | `docs/field-guide.html` | swap | silent | bit | field guide restored the one retained author on the approved implementation line |

Both probe runs restored the file. The after run reported one failure. No row pins the new pointer sentence, so no omission self-probe applies to it.

### Verification

| command | exit | result |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestRootConformance` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | 0 | pass, 0 failures, 0 skips |
| `bench test --full --package ./internal/conformance --run "TestFreshTicketAuthors\|TestIntegrationSourceWorkflowAnchorsBiteIndependently"` | 0 | pass, 0 failures, 0 skips |
| `bench test --package ./internal/anchors/...` | 0 | pass, 0 failures, 0 skips |
| `bench test --check guidance-prose-budgets` | 0 | pass, 0 failures, 0 skips |
| `go vet ./...` | 0 | no output |
| `bench commit` lane | 0 | pass |
