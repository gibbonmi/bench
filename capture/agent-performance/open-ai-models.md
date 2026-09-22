# OpenAI model scorecard

Last incorporated landing: `ft323-project-conventions` (`333bc7ce38821638c5d352a65ac896d2642f943e`, 2026-09-22 UTC).
Astra/ultra retains implementation, coordination, verification, and one post-review repair cycle.
Sol/high supplies three independent review axes. Sol/medium supplies one read-only host-fact consultation.

The final axes have zero blockers. The retained-source and prospective landing gates pass.
Run tokens, context capacity, provider costs, and comparative latency remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Astra / ultra, high, medium, low | retained implementation, repair, verification, and coordination; latest FT323 ticket | Astra/ultra completes FT323 in one post-review repair cycle. Review removes extra policy and makes the local opener exception explicit. | Use medium for Coverage and debug validation. Use low for bounded coding when the user selects it. | Repair forecast: 0.09 over 1 labeled pair; 0 abstentions |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. | unknown |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `bounded-repair-policy`, all three Terra/high axes passed before a cross-harness reviewer found two policy defects. Each axis verified the repairs; Coverage corrected one source-identity claim after the coordinator checked Git. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. | unknown |
| Luna / max, medium | prose implementation and bounded repairs | Luna preserved the Ticket 4 prose pass. Review found owner-identity and instruction-shape defects that required an Astra repair and refreshed adoption evidence. | Use Luna for narrow prose changes after an owner census and before independent review. | unknown |
| Sol / high | retained implementation and adoption work | Sol completed substantive guidance and fresh adoption work. Coordinator checks found stale source identity and one final parser defect outside the prose. | Use Sol/high for exact specification chunks under independent review and coordinator probes. | unknown |
| Sol / high, medium | independent review axes and read-only host diagnosis | Sol/high finds two FT323 repair targets and confirms their closure in three independent axes. Sol/medium traces the exact host pins and retained opener evidence. | Use separate Standards and Spec contexts with an exact source range and an explicit review cap. | unknown |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy | Astra closed two cross-harness findings in two repair cycles, and all 21 rows passed final reconciliation. Terra missed both initial defects and verified their repair. | spec/ticket, reviewer, and orchestrator | Keep independent cross-harness review for shared workflow guidance. |
| Project conventions | One retained Astra/ultra author resolves two review targets in one cycle. Three Sol/high axes confirm the repaired source before the green landing. | retained author | Keep host facts local and give host exceptions explicit precedence over shared command steps. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Repair collection pilot | Sol/high implemented three chunks; Astra/medium found gaps in each, then passed every repaired source and final composition. | delegate | Keep one retained author across bounded repair rounds and bind every review to its source. |
| Debug loop guidance | Astra, Sol, and Luna completed two tickets and the DG15 debug repair. Six Ticket 5 repair rounds ended with three passing review axes and a green landing gate. | delegate, reviewer, tree/tooling, and orchestrator | Keep source identity exact, retain one debug author, and require independent final axes. |

## Current decisions

- Keep implementation, tests, probes, and repairs in the user-approved author session.
- Apply the review and repair allowance across fresh and resumed sessions.
- Use Astra/medium for Coverage review when the user names that line.
- Use separate Sol/high contexts for Standards and Spec review when the user names that line.
- Use Astra/low for bounded coding only when the user names that line.
- Use Luna for narrow prose changes after the author identifies every affected owner.
- Bind each review to its exact source and verify every returned source claim against Git.
- Attach each omission probe to the check that independently detects the omission.
- Serialize gate operations and repeat author verification after the source digest changes.
- Run preflight after ticket metadata, fixture, registry, or owner changes.
- Require an explicit bounded extension after the recorded repair limit is exhausted.
- Preserve unknown model identity, token counts, cache benefit, costs, and comparative latency as unknown.
- Change general routing only after two comparable runs, one controlled comparison, or explicit user direction.
