# OpenAI model scorecard

Last incorporated landing: `completion-evidence` (`ea90a8500405a289a8ff125e49dd5a7b7e682282`, 2026-09-11).
Astra/high retained implementation across three chunks. Sol/high supplied nine initial
axis returns and nine clean repair returns. Token use and provider-dollar cost remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / high | retained implementer and orchestrator, 3 chunks on 1 landing | On `completion-evidence`, Astra retained all code, tests, probes, and repairs and closed all 38 acceptance rows. Review exposed missing evidence predicates and system fixture consumers; the author repaired them before each required checkpoint. | User-approved retained implementation across source identity, gate reuse, and publication boundaries. |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `implementation-continuation`, eight Terra/high review passes found two Coverage gaps. One gap protected continued progress, and one required an available diagnostic fallback; both second passes accepted the repairs. | Standards, Spec, and Coverage review in separate contexts; independently census production callers |
| Luna / max, medium | implementation, 7 bounded tickets/repairs + reviewer, 4 axis passes on `stable-owner-landing` | Implementation: 2/7 first-pass with all terminal gates green. Review: three initial axes returned 9 raw findings and 8 de-duplicated targets with the citation standard held, and the repair-scoped re-review verified all seven predicates. | Low-cost writer for narrow slices under mandatory inspection |
| Sol / low | implementation, bounded tickets and repairs on 1 landing | The delegates returned focused tests and mutation probes. The retirement repair reproduced the FT94 ledger red, changed one owner, proved that restoring the old value made the test red, and restored green. | Exact ticket seams and small repairs under coordinator verification |
| Sol / high | independent review, latest 10 of 18 axis returns; retained implementation on 1 earlier landing | On `completion-evidence`, separate axes found source-chain, verification, fixture, and one-source defects; all nine final repair returns were clean. The author verified each result against the frozen source and current tests. | Separate Standards, Spec, and Coverage contexts under the approved line; retained implementation at related seams. |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| FT311 recoverable-reset candidate | Astra/medium built five tickets green under Sol/high; a Fable/high round then found three behavior defects the Terra/medium review missed. | delegate and reviewer | A same-provider medium review of a medium build needs one higher-tier round before a landing. |
| FT311 ticket 2 repair | Sol/high unified charge and proposal orchestration and completed every DP10-DP14 partition. Exact-baseline tests and two restored mutations bit. | delegate | Use Sol/high when bounded repair crosses command retry, authority, and graph semantics. |
| FT311 ticket 4 | Terra/high produced build and triage guidance. Review found that active guidance omitted the required triage inputs. | delegate | Trace each spec input into final phase text before returning a guidance ticket. |
| Completion evidence | Astra/high implemented three chunks; Sol/high returned 19 initial findings and nine clean repair results. All 38 rows and both final gates passed. | implementer, reviewer, and orchestrator | Use current evidence at chunk boundaries and independently enumerate binary-driven consumers. |
| Implementation continuation | Sol/high implemented two retained chunks, and Terra/high found one Coverage gap in each chunk. Both repairs gained live omission consequences before the green landing. | orchestrator and reviewer | Keep one retained author, and use independent Coverage review to test prose consequences. |

## Current decisions

- Use the user-approved Astra/high line for retained work across source identity and publication boundaries.
- Use Sol/high for separate Standards, Spec, and Coverage review of that work.
- Apply the conditional review line from `craft-line` after each implementation chunk.
- Keep all implementation, tests, probes, and repairs in the approved retained session.
- Include binary-driven system journeys when enumerating consumers of a changed gate obligation.
- Bind native review outcomes and author verification to the source each performer examined.
- Require every axis to cover the current repair before a completion checkpoint.
- Use the full ordinary suite and the separate system suite when changed-package selection excludes tagged owners.
- Keep the installed broker in authority until publication and the sanctioned rebuild finish.
- Change general routing only after two comparable runs, one controlled comparison, or explicit reviewer direction.
