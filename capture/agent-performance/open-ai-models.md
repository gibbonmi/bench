# OpenAI model scorecard

Last incorporated landing: `repair-collection-pilot` (`2675613a5473631dff45e665adf1f84a3bdf19ea`, 2026-09-13).
The user approved Sol/high for retained implementation and Astra/medium for native review.
The user excluded cross-harness review.
Token counts, provider costs, and comparative latency remain unknown.

Latest spec observation: `repair-collection-pilot` used Sol/high authorship and Astra/medium review across three chunks.
Astra found blocking gaps in each chunk. Sol closed them in one, two, and two repair rounds.
All 66 rows and the final landing gate passed.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / high, medium | retained implementation, orchestration, and native review; latest three pilot chunks | Astra/medium reviewed Standards, Spec, and Coverage for all three chunks. It found blocking gaps in each chunk and passed every repaired source, including the final composition. | Native review on the user-approved line, with separate contexts for each axis. |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `bounded-repair-policy`, all three Terra/high axes passed before a cross-harness reviewer found two policy defects. Each axis verified the repairs; Coverage corrected one source-identity claim after the coordinator checked Git. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held, and the repair-scoped re-review verified all seven predicates. | Low-cost writer for narrow slices under mandatory inspection |
| Sol / high | retained implementation, latest pilot and five bounded fixes | Sol/high implemented three tickets and 66 acceptance rows. It needed five total review repair rounds; all required mutations bit and the final gate passed. | Exact spec chunks under independent review and coordinator probes. |
| Sol / high | latest spec review and ticket-slicing assignment | Sol found the optional-advice representation gap and missing mandatory-standard coverage. Its one-ticket slice needed two returned passes for a probe attachment, gate cost, and registry ownership closure. | Spec review and slicing under the user-approved choreography, with independent slice review. |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy | Astra closed two cross-harness findings in two repair cycles, and all 21 rows passed final reconciliation. Terra missed both initial defects and verified their repair. | spec/ticket, reviewer, and orchestrator | Keep independent cross-harness review for shared workflow guidance. |
| Complete failure output | The combined reviewer found assertions that accepted stripped controls; the retained author added escaped-output expectations and a biting strip mutation. | implementer and reviewer | Distinct axes can find useful gaps in one context, but this does not establish equivalence with separate reviewers. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Repair collection pilot | Sol/high implemented three chunks; Astra/medium found gaps in each, then passed every repaired source and final composition. | delegate | Keep one retained author across bounded repair rounds and bind every review to its source. |
| Implementation continuation | Sol/high implemented two retained chunks; Terra/high found a Coverage gap in each, and both repairs gained biting omission checks. | orchestrator and reviewer | Keep repairs with the retained author and test the claimed consequence. |

## Current decisions

- Keep implementation, tests, probes, and repairs in the user-approved author session.
- Apply the chunk allowance across fresh reviews and resumed sessions.
- Keep the conditional review line from `craft-line` when the user has not overridden it.
- Keep the standing cross-harness pass for shared kit guidance unless the user excludes it.
- Use Astra/medium for separate native review axes when the user names that line.
- Bind each review to the examined source and check each returned source claim against Git.
- Attach each omission probe to the check that independently detects that omission.
- Use fresh read-only worktrees for parallel reviews and serialize operations that run a gate.
- Re-run author verification after integration when the source digest changes.
- Run review preflight after plan-affecting ticket metadata changes.
- Preserve unknown model identity, token counts, costs, and comparative latency as unknown.
- Change general routing only after two comparable runs, one controlled comparison, or explicit user direction.
