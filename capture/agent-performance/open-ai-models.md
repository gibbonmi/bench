# OpenAI model scorecard

Last incorporated landing: `debug-loop-guidance` (`ae5182b1178603fb3f5ad7cc734c3c890ef00f46`, 2026-09-16 UTC).
Astra, Sol, and Luna implemented or reviewed the final two tickets and the DG15 debug repair.
All final review axes passed, and the prospective landing gate passed every available phase.
Token counts, cache benefit, provider costs, and comparative latency remain unknown.

The research report records two user-supplied Codex issue reports about repeated waits and context work.
The cause and quota effect in this run remain unverified.
Do not attribute those repeated checks to the fork workflow without separate evidence.

Latest shaping observation: `debug-loop-guidance` dispatched Sol/high for one primary-source research question.
The coordinator reopened all three papers and retained their transfer limits.
The return misidentified its inherited model; runtime identity remains unknown beyond the explicit dispatch setting.
This research observation supplies no implementation or review-quality result and changes no routing decision.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Astra / high, medium, low | implementation, debug repair, Coverage review, and coordination | Astra/medium found the final DG15 bypass. The retained debug author needed three attempts and one approved Standards-only extension before all axes passed. | Use medium for Coverage and debug validation. Use low for bounded coding when the user selects it. |
| Terra / medium, high | implementation, latest 10 bounded tickets and repairs | Terra/high implemented FT311 ticket 4 with live evidence open. Review found the missing triage input contract, and repair added its five inputs. | Medium for exact one-seam tickets under coordinator mutation. High when one fact crosses multiple policy consumers. |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `bounded-repair-policy`, all three Terra/high axes passed before a cross-harness reviewer found two policy defects. Each axis verified the repairs; Coverage corrected one source-identity claim after the coordinator checked Git. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. |
| Luna / max, medium | prose implementation and bounded repairs | Luna preserved the Ticket 4 prose pass. Review found owner-identity and instruction-shape defects that required an Astra repair and refreshed adoption evidence. | Use Luna for narrow prose changes after an owner census and before independent review. |
| Sol / high | retained implementation and adoption work | Sol completed substantive guidance and fresh adoption work. Coordinator checks found stale source identity and one final parser defect outside the prose. | Use Sol/high for exact specification chunks under independent review and coordinator probes. |
| Sol / high | Standards and Spec review | Separate Sol/high axes found Ticket 4 and Ticket 5 defects, then passed the final range through `6265ac69`. | Use separate Standards and Spec contexts with an exact source range and an explicit review cap. |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy | Astra closed two cross-harness findings in two repair cycles, and all 21 rows passed final reconciliation. Terra missed both initial defects and verified their repair. | spec/ticket, reviewer, and orchestrator | Keep independent cross-harness review for shared workflow guidance. |
| Shared preparation spec | Two sequential author forks produce three serial tickets. One Sol/high pass finds two blockers, which the author repairs under coordinator verification. | spec/ticket, reviewer, and coordinator | Match parsed authority to the ticket union and assign combined behavior to its final consumer. This phase supplies no comparative cache or cost evidence. |
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
