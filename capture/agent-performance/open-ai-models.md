# OpenAI model scorecard

Last incorporated landing: `durable-broker-shim` (`4757bfb4edf04daf955ca49b92fd3d3be9c53dca`, 2026-09-28).
One inline author completed three light tickets without delegation, as the user requested.
Its model identifier is unknown, and its effort was high.

Two Git reader extractions passed behavior tests and three omission probes.
The guidance ticket passed its prose and anchor checks; all three landing gates passed.
No formal independent review ran for these light tickets.
Provider usage, costs, and comparative latency remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Astra / xhigh, ultra, high, medium, low | decider, 8 consultations; retained implementation, repair, verification, and coordination | On `ref-inventory` Astra/xhigh gave exact row text and decisions in seven consultations, and each decision held through the confirming rounds. The eighth run hung with no output and was stopped. | Reviewer decisions by user direction when the route answers; retain the user-approved author through probes and repairs. | Repair forecast: 0.065 over 4 batch pairs; 0 abstentions |
| Unknown Codex model / high | inline implementation, repair, and orchestration; 1 spec build and 3 light tickets | Three light tickets landed with passing focused checks and full gates, without a post-review repair cycle. Omission probes rejected missing dirtiness, commit peeling, and caller flags; installed-command checks verified the shim restoration sequence. | Retain user-directed inline authorship for bounded tasks; formal reviews keep independent contexts. | Repair forecast: 0.49 over 1 prior pair; this drain is unscored |
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
| Drain Git readers and broker guidance | One inline author landed three light tickets with passing gates and three demonstrated omission probes. | orchestrator | Preserve caller behavior in shared readers and verify guidance against its real command sequence. |

## Current decisions

- Preserve the user's authorship and delegation choices across resumed sessions.
- Treat these light tickets as behavior-preservation evidence, not a comparison of models.
- Keep independent Standards, Spec, and Coverage contexts for formal review.
- Apply current input-family rules to older staged specs before implementation.
- Compile reviewer examples before accepting their mutation evidence.
- Include composition fixtures when changed production windows affect their callers.
- Retain a demonstrated omission or mutation for each independent expectation.
- Bind verification to its source and preserve evidence before worktree release.
- Follow final-check's broker-rehearsal sequence when candidate repair changes the installed shim.
- Attribute upstream inventory omissions separately from author defects.
- Keep model identity, usage, costs, and unscored claims explicit when evidence is unavailable.
- Change routing only after comparable evidence or explicit user direction.

## Specification evidence

On 2026-10-02, the invoking session authored the staged CLI and desktop consistency spec under the user's explicit authorship choice.
The declared authoring line was Astra / high; provider usage, cost, and comparative latency are unknown.
The draft defines 64 acceptance rows across three serial tickets.
Mechanical author checks found planned-test citation and ownership-closure defects, which the author corrected before sign-off.
The reviewer approved the specification and ticket graph on 2026-10-02 in one sign-off round.
Implementation qualification remains pending.

The approved implementation line is Sol / high, with repair and live qualification identified as the harder chunks.
This specification supplies no new implementation performance evidence and does not change the routing table.
