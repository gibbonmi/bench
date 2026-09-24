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
