# OpenAI model scorecard

Last incorporated landing: `record-completion` (`6281fc3f6614cd14040b4dd1d621222f15acf639`, 2026-10-02).
An unknown Codex model coordinated the resumed landing. Sol/high sessions implemented the completion writer and prepared the phase-close capture.

The completion writer landed green after focused checks, a routing check, and an independent mutation. The original landing also stayed green.
Provider usage and charges remain unknown.

## Cost assumptions

The current planning input prices Luna tokens at 0.2x Terra, so Luna reaches dollar break-even near 5x Terra's token use. The FT311 benchmark arm recorded 99 million input tokens with 98 percent cached and no dollar figure.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Astra / xhigh, ultra, high, medium, low | decider, 8 consultations; retained implementation, repair, verification, and coordination | On `ref-inventory` Astra/xhigh gave exact row text and decisions in seven consultations, and each decision held through the confirming rounds. The eighth run hung with no output and was stopped. | Reviewer decisions by user direction when the route answers; retain the user-approved author through probes and repairs. | Repair forecast: 0.065 over 4 batch pairs; 0 abstentions |
| Unknown Codex model / high or unknown | inline implementation, repair, and orchestration; 1 spec build, 3 light tickets, and 1 resumed landing | The earlier high-effort author landed three light tickets with passing gates and omission probes. The current unknown-effort coordinator preserved five checks, reconciled 86 rows, and isolated three host-tool failures. | Retain user-directed inline authorship for bounded tasks; keep exact model identity and effort unknown when the transcript does not supply them. | Repair forecast: 0.49 over 1 prior pair; current landing unscored |
| Terra / high, medium, low | research and seam inspection, latest 5 assignments; semantic review, latest 10 axis passes | On `bounded-repair-policy`, all three Terra/high axes passed before a cross-harness reviewer found two policy defects. Each axis verified the repairs; Coverage corrected one source-identity claim after the coordinator checked Git. | Standards, Spec, and Coverage review in separate contexts; independently census production callers. | unknown |
| Luna / max, medium | prose implementation and bounded repairs | Luna preserved the Ticket 4 prose pass. Review found owner-identity and instruction-shape defects that required an Astra repair and refreshed adoption evidence. | Use Luna for narrow prose changes after an owner census and before independent review. | unknown |
| Sol / high | retained implementation, adoption, and phase-close work | One Sol/high writer landed the completion writer after focused checks and an independent mutation. A second Sol/high writer prepared the retro and both scorecards; usage and charges are unknown. | Use Sol/high for exact specification chunks under independent review and coordinator probes. | unknown; this repair has 1 abstention |
| Sol / high | independent review, latest 10 of 39 axis returns | The latest ten axes identified the import-resolution gap and confirmed its repairs. Earlier clean axes missed the Git guard fixture omission that the checkpoint exposed. | Keep separate axes, compile proposed examples, and census production-window callers independently. | 0.04 over 1 labeled finding; 9 abstentions |

## Representative evidence

| task | result | attribution | routing signal |
| --- | --- | --- | --- |
| Bounded repair policy | Astra closed two cross-harness findings in two repair cycles, and all 21 rows passed final reconciliation. Terra missed both initial defects and verified their repair. | spec/ticket, reviewer, and orchestrator | Keep independent cross-harness review for shared workflow guidance. |
| Test determinism | Sol/high found five retained defects across 39 axes; all current axes and the final gate pass after repair. | spec/ticket for the incomplete wait-family and switch-fixture inventories | Apply current coverage rules to older specs and retain independent caller checks. |
| FT311 recoverable-reset candidate | A Fable/high round found three behavior defects after the candidate's medium-tier review. | delegate and reviewer | Keep independent adversarial verification when authority or destructive behavior crosses boundaries. |
| Repair collection pilot | Sol/high implemented three chunks; Astra/medium found gaps in each, then passed every repaired source and final composition. | delegate | Keep one retained author across bounded repair rounds and bind every review to its source. |
| Worktree seam reduction close | The resumed coordinator found the missing completion writer after the original implementation passed. Sol/high added the writer, and an independent mutation changed covered evidence to pending. | tree/tooling and delegate | Dogfood completion before source release, and keep host-tool failures apart from author quality. |

## Current decisions

- Preserve the user's authorship and delegation choices across resumed sessions.
- Keep the coordinator's exact model unknown when the transcript supplies only the provider.
- Attribute the completion-writer gap to tree/tooling, outside the original ticket authors and spec.
- Keep the Node and npm artifact failure separate from model quality.
- Treat bounded inline repairs as behavior-preservation evidence, not a model comparison.
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
