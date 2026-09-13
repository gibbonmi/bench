# OpenAI model scorecard

Last incorporated phase: `bounded-repair-policy` specification (review basis `c7255226ce92d4299b808d0ff374526d9cdbe805`, 2026-09-13).
Astra authored the spec and reviewed the slice.
Sol/high reviewed the spec, then sliced it in the same agent context.
The user capped each review at one round and pre-approved landing.
Token counts, provider costs, and comparative latency remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / high, medium | retained implementation, orchestration, and latest specification phase | Astra folded two spec-authoring gaps found by Sol. Its slice review caught a probe attached to the wrong check and redundant full gating. | Retained authorship and source-based acceptance on the user-approved line. |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `harness-proof-closure`, Terra/high completed 15 axis reviews in one round and found one missing outside-Git refusal case. A Standards reviewer corrected one false finding after the coordinator showed the canonical journal. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held, and the repair-scoped re-review verified all seven predicates. | Low-cost writer for narrow slices under mandatory inspection |
| Sol / high | bounded implementation, latest five fixes | Five fixes reached green landings under coordinator mutation. One filename repair added modified, copied, and renamed inputs after review exposed their absence. | Small fixes with a named owner, complete input enumeration, and independent coordinator verification. |
| Sol / high | latest spec review and ticket-slicing assignment | Sol found the optional-advice representation gap and missing mandatory-standard coverage. Its one-ticket slice needed two returned passes for a probe attachment, gate cost, and registry ownership closure. | Spec review and slicing under the user-approved choreography, with independent slice review. |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy spec | Sol found two blocking semantic gaps, and Astra caught the slice's probe attachment defect. Both reviews stayed within one round. | spec/ticket and delegate | Keep review evidence tied to the exact failure the named check can detect. |
| Complete failure output | The combined reviewer found assertions that accepted stripped controls; the retained author added escaped-output expectations and a biting strip mutation. | implementer and reviewer | Distinct axes can find useful gaps in one context, but this does not establish equivalence with separate reviewers. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Workflow assessment | Astra/medium retained all authorship; Sol/high reviewed three chunks, and all 39 rows passed reconciliation and landing. | orchestrator | Preserve failed attempts and unknown measurements in the assessment. |
| Implementation continuation | Sol/high implemented two retained chunks; Terra/high found a Coverage gap in each, and both repairs gained biting omission checks. | orchestrator and reviewer | Keep repairs with the retained author and test the claimed consequence. |

## Current decisions

- Keep implementation, tests, probes, and repairs in the user-approved author session.
- Honor explicit review choreography and round caps, and fold accepted findings within the opened round.
- Keep the conditional review line from `craft-line` when the user has not overridden it.
- Bind reviews to the examined source and verify the returned paths before accepting a delegate.
- Attach each omission probe to the check that independently detects that omission.
- Give coordinator mutations a different kind and site from the author's probe.
- Re-run author verification after integration when the source digest changes.
- Preserve unknown model identity, token counts, costs, and comparative latency as unknown.
- Change general routing only after two comparable runs, one controlled comparison, or explicit user direction.
