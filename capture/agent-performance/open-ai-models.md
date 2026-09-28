# OpenAI model scorecard

Last incorporated landing: `test-determinism` (`76ad133c6b852472abf7ea39cdd28a22df1c2b38`, 2026-09-28).
One inline author completed eight tickets and ten repair rounds; its model identifier is unknown.
Thirty-nine Sol/high review axes returned, and two additional dispatches failed before startup.
All 51 acceptance rows passed reconciliation, and the implementation landing gate passed.
Provider usage, costs, and comparative latency remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Astra / xhigh, ultra, high, medium, low | decider, 8 consultations; retained implementation, repair, verification, and coordination | On `ref-inventory` Astra/xhigh gave exact row text and decisions in seven consultations, and each decision held through the confirming rounds. The eighth run hung with no output and was stopped. | Reviewer decisions by user direction when the route answers; retain the user-approved author through probes and repairs. | Repair forecast: 0.065 over 4 batch pairs; 0 abstentions |
| Unknown Codex model / high | inline implementation, repair, and orchestration; 1 build with 8 tickets | The author completed all 51 rows through ten repair rounds and preserved native red and green evidence. The final gate passed; review and checkpoint checks exposed both implementation defects and upstream inventory omissions. | Keep the explicit user-directed authorship route and independent source-bound review. | Repair forecast: 0.49 over 1 pair; 0 abstentions |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `bounded-repair-policy`, all three Terra/high axes passed before a cross-harness reviewer found two policy defects. Each axis verified the repairs; Coverage corrected one source-identity claim after the coordinator checked Git. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. | unknown |
| Luna / max, medium | prose implementation and bounded repairs | Luna preserved the Ticket 4 prose pass. Review found owner-identity and instruction-shape defects that required an Astra repair and refreshed adoption evidence. | Use Luna for narrow prose changes after an owner census and before independent review. | unknown |
| Sol / high | retained implementation and adoption work | Five Sol/high writers delivered the learning fixes. Coordinator feedback preserved plan authority and conditional consultation deadlines; focused checks and landing gates passed. | Use Sol/high for exact specification chunks under independent review and coordinator probes. | unknown; this repair has 1 abstention |
| Sol / high | independent review, latest 10 of 39 axis returns | The latest ten axes identified the import-resolution gap and confirmed its repairs. Earlier clean axes missed the Git guard fixture omission that the checkpoint exposed. | Keep separate axes, compile proposed examples, and census production-window callers independently. | 0.04 over 1 labeled finding; 9 abstentions |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy | Astra closed two cross-harness findings in two repair cycles, and all 21 rows passed final reconciliation. Terra missed both initial defects and verified their repair. | spec/ticket, reviewer, and orchestrator | Keep independent cross-harness review for shared workflow guidance. |
| Test determinism | Sol/high found five retained defects across 39 axes; all current axes and the final gate pass after repair. | spec/ticket for the incomplete wait-family and switch-fixture inventories | Apply current coverage rules to older specs and retain independent caller checks. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Repair collection pilot | Sol/high implemented three chunks; Astra/medium found gaps in each, then passed every repaired source and final composition. | delegate | Keep one retained author across bounded repair rounds and bind every review to its source. |
| Drain learning fixes | Five Sol/high authors passed focused checks and landing gates after coordinator corrections to two guidance changes. | delegate | Verify preserved authority and timeout conditions against the canonical owner. |

## Current decisions

- Keep the user-directed author and repair policy explicit across resumed sessions.
- Keep independent Standards, Spec, and Coverage contexts for formal review.
- Apply current coverage inventory rules when an older spec begins implementation.
- Compile Go examples before accepting their claimed mutation evidence.
- Census callers of changed production windows and include their composition fixtures in required verification.
- Require each independent expectation to reject its named omission or mutation.
- Bind review and verification evidence to the exact source and retain failed attempts.
- Serialize gates and source writes, and preserve local evidence before releasing a worktree.
- Reconcile repair totals from ticket causes and keep upstream omissions separate from author defects.
- Keep the installed shim on a durable checkout after broker rehearsal.
- Reserve time for retirement, the combined retro and scorecard commit, and handoff.
- Keep unknown model identity, usage, and costs explicit.
- Change routing only after comparable evidence or explicit user direction.
