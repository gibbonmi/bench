# OpenAI model scorecard

Last incorporated landing: `harness-proof-closure` (`60269540fd31225fd3e29476060e6b2f81a6b2a5`, 2026-09-13).
Five Astra/high authors implemented one ticket each, and the Astra/high orchestrator integrated all five chunks.
Three Terra/high reviewers supplied Standards, Spec, and Coverage for every chunk in one user-capped round.
Coverage returned one accepted test gap; all other source findings were clear.
Native token counts, provider costs, and comparative review latency remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / high, medium | retained implementation and orchestration, latest five tickets | Astra/high completed all five FT120 tickets with biting restored probes. Four tickets needed no returned write repair; one received a test-only repair after Coverage found a missing refusal case. | Retained implementation on the user-approved model and effort. |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `harness-proof-closure`, Terra/high completed 15 axis reviews in one round and found one missing outside-Git refusal case. A Standards reviewer corrected one false finding after the coordinator showed the canonical journal. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held, and the repair-scoped re-review verified all seven predicates. | Low-cost writer for narrow slices under mandatory inspection |
| Sol / high | bounded implementation, latest five fixes | Five fixes reached green landings under coordinator mutation. One filename repair added modified, copied, and renamed inputs after review exposed their absence. | Small fixes with a named owner, complete input enumeration, and independent coordinator verification. |
| Sol / high | combined Standards, Spec, and Coverage review, latest ten fixes | The filename and full-output reviews found two coverage gaps and one changelog defect. All final reviews were clear after the retained authors repaired accepted findings. | One reviewer per fix for the user-requested trial, with separate axis reports. |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Roadmap filename fix | The combined reviewer found missing modified, copied, and renamed cases; the author added real Git statuses and a biting filter mutation. | delegate and reviewer | Enumerate the producer's input family before accepting a bounded fix. |
| Complete failure output | The combined reviewer found assertions that accepted stripped controls; the retained author added escaped-output expectations and a biting strip mutation. | implementer and reviewer | Distinct axes can find useful gaps in one context, but this does not establish equivalence with separate reviewers. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Workflow assessment | Astra/medium retained all authorship; Sol/high reviewed three chunks, and all 39 rows passed reconciliation and landing. | orchestrator | Preserve failed attempts and unknown measurements in the assessment. |
| Implementation continuation | Sol/high implemented two retained chunks; Terra/high found a Coverage gap in each, and both repairs gained biting omission checks. | orchestrator and reviewer | Keep repairs with the retained author and test the claimed consequence. |

## Current decisions

- Keep implementation, tests, probes, and repairs in the user-approved author session.
- Keep Standards, Spec, and Coverage in separate reviewer contexts for spec-backed chunks.
- Honor a user-set review-round cap; resolve accepted findings inside the opened round without a fresh fan-out.
- Keep the conditional review line from `craft-line` when the user has not overridden it.
- Bind every accepted review to the examined source and verify the tree stayed unchanged.
- Give coordinator mutations a different kind and site from the delegated author's probe.
- Return accepted findings to the original author and rebind evidence to the repaired source.
- Re-run author verification after a rebase when the source digest changes.
- Preserve unknown model identity, token counts, costs, and comparative latency as unknown.
- End completed author turns before allocating fresh reviewer contexts.
- Change general routing only after two comparable runs, one controlled comparison, or explicit user direction.
