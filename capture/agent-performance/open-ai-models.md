# OpenAI model scorecard

Last incorporated landing: `ref-inventory` (`ab886f5b5ee04867110e79e4e1879a920a50434f`, 2026-09-27).
Astra/xhigh served as the reviewer's decider through `codex exec` for seven consultations on that Claude Code build.
It decided D1 to D10, including the dropped holder kind, the branch-path route, and the faulted count line.
The eighth consultation produced no output in over three hours, and the reviewer moved later consultations to Fable/high.

At capture time, FT337 has an accepted draft, but human approval and its staging gate remain pending.
The spec author and coordinator's exact model metadata and provider costs are unknown.
The debug author's invalid confidence remains an abstention.
Tokens and comparative latency remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Astra / xhigh, ultra, high, medium, low | decider, 8 consultations; retained implementation, repair, verification, and coordination | On `ref-inventory` Astra/xhigh gave exact row text and decisions in seven consultations, and each decision held through the confirming rounds. The eighth run hung with no output and was stopped. | Reviewer decisions by user direction when the route answers; retain the user-approved author through probes and repairs. | Repair forecast: 0.065 over 4 batch pairs; 0 abstentions |
| Terra / low, medium, high | implementation and fresh-session adoption | Terra/low interprets all six guidance scenarios as intended. This exercise proves comprehension only. | Use fresh readers against committed owner bytes. | unknown |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `bounded-repair-policy`, all three Terra/high axes passed before a cross-harness reviewer found two policy defects. Each axis verified the repairs; Coverage corrected one source-identity claim after the coordinator checked Git. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. | unknown |
| Luna / max, medium | prose implementation and bounded repairs | Luna preserved the Ticket 4 prose pass. Review found owner-identity and instruction-shape defects that required an Astra repair and refreshed adoption evidence. | Use Luna for narrow prose changes after an owner census and before independent review. | unknown |
| Sol / high | retained implementation and adoption work | Sol/high repaired the FT337 guidance conflict in one fix cycle after daemon recovery. Independent omission probes and the full landing gate passed. | Use Sol/high for exact specification chunks under independent review and coordinator probes. | unknown; this repair has 1 abstention |
| Sol / high, medium | independent review and admission diagnosis | Sol/high identifies six accepted repair targets across two tickets and one refuted request. All four tickets pass final native review and their landing gates. | Keep independent contexts for each axis and exact source identity. | unknown |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy | Astra closed two cross-harness findings in two repair cycles, and all 21 rows passed final reconciliation. Terra missed both initial defects and verified their repair. | spec/ticket, reviewer, and orchestrator | Keep independent cross-harness review for shared workflow guidance. |
| Map retirement debug | Astra/high removes six obsolete maps and proves the ownership check with a reversible omission probe. Landing remains subject to reviewer approval. | tree/tooling | Separate structural map ownership from semantic closure evidence. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Repair collection pilot | Sol/high implemented three chunks; Astra/medium found gaps in each, then passed every repaired source and final composition. | delegate | Keep one retained author across bounded repair rounds and bind every review to its source. |
| FT337 review-guidance debug | Sol/high removed stale retrieval instructions and preserved the full control. The coordinator confirmed omission probes and the complete landing gate. | tree/tooling | Keep Process links to canonical Entry orientation beside explicit contradiction guards. |

## Current decisions

- Retain the user-directed author through an authorized debug repair and its verification.
- Carry repair counts and limits across resumed sessions.
- Keep independent Standards, Spec, and Coverage contexts for formal review.
- Bind review and adoption evidence to exact committed source bytes.
- Check the current tree before accepting a retained source for a later landing.
- Serialize full gates and landings, with no repository mutation during a gate.
- Reserve time for the paired capture commit, restoration, and handoff.
- Keep unknown model usage and costs explicit, and count invalid confidence as an abstention.
- Stop a `codex exec` consultation that produces no output past its expected window, and reroute the question.
- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
