# OpenAI model scorecard

Last incorporated landing: `bounded-repair-policy` (`c89922c93d75d34840863ca1a4d58cb217c3b6e4`, 2026-09-13).
The user approved Astra/high for implementation and orchestration in this session.
Three Terra/high contexts reviewed Standards, Spec, and Coverage before and after repair.
Token counts, provider costs, and comparative latency remain unknown.

Latest spec observation: `repair-collection-pilot` used Astra/high authorship and two Sol/high review passes.
Sol found missing collection and report contracts, then narrowed the remaining targets to storage and interval evidence.
The final author fold remains outside independent review; the spec needs reviewer sign-off.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / high, medium | retained implementation and orchestration, latest guidance landing | Astra closed two findings from the cross-harness review and retained all 21 acceptance outcomes. It used two repair cycles, corrected one reviewer source claim, and caused three review-worktree merge refusals through concurrent dispatch. | Retained authorship and source-based acceptance on the user-approved line. |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `bounded-repair-policy`, all three Terra/high axes passed before a cross-harness reviewer found two policy defects. Each axis verified the repairs; Coverage corrected one source-identity claim after the coordinator checked Git. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held, and the repair-scoped re-review verified all seven predicates. | Low-cost writer for narrow slices under mandatory inspection |
| Sol / high | bounded implementation, latest five fixes | Five fixes reached green landings under coordinator mutation. One filename repair added modified, copied, and renamed inputs after review exposed their absence. | Small fixes with a named owner, complete input enumeration, and independent coordinator verification. |
| Sol / high | latest spec review and ticket-slicing assignment | Sol found the optional-advice representation gap and missing mandatory-standard coverage. Its one-ticket slice needed two returned passes for a probe attachment, gate cost, and registry ownership closure. | Spec review and slicing under the user-approved choreography, with independent slice review. |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy | Astra closed two cross-harness findings in two repair cycles, and all 21 rows passed final reconciliation. Terra missed both initial defects and verified their repair. | spec/ticket, reviewer, and orchestrator | Keep independent cross-harness review for shared workflow guidance. |
| Complete failure output | The combined reviewer found assertions that accepted stripped controls; the retained author added escaped-output expectations and a biting strip mutation. | implementer and reviewer | Distinct axes can find useful gaps in one context, but this does not establish equivalence with separate reviewers. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Workflow assessment | Astra/medium retained all authorship; Sol/high reviewed three chunks, and all 39 rows passed reconciliation and landing. | orchestrator | Preserve failed attempts and unknown measurements in the assessment. |
| Implementation continuation | Sol/high implemented two retained chunks; Terra/high found a Coverage gap in each, and both repairs gained biting omission checks. | orchestrator and reviewer | Keep repairs with the retained author and test the claimed consequence. |

## Current decisions

- Keep implementation, tests, probes, and repairs in the user-approved author session.
- Apply the chunk allowance across fresh reviews and resumed sessions.
- Keep the conditional review line from `craft-line` when the user has not overridden it.
- Keep the standing cross-harness pass for shared kit guidance.
- Bind reviews to the examined source and check returned source claims against Git.
- Attach each omission probe to the check that independently detects that omission.
- Use fresh read-only worktrees for parallel reviews; serialize operations that run a gate.
- Re-run author verification after integration when the source digest changes.
- Preserve unknown model identity, token counts, costs, and comparative latency as unknown.
- Change general routing only after two comparable runs, one controlled comparison, or explicit user direction.
