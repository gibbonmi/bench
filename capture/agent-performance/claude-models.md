# Claude model scorecard

Last incorporated landing: `calibrated-decisions` (`5139240fe02a22199f5114eff5c5eee0386970a7`, 2026-09-17).
Seven Opus/high authors landed seven tickets in six chunks under a Fable/low orchestrator.
Fable/medium and Opus/medium ran the first review pass, and Sonnet ran every later pass.
The landing recorded the first 61 labeled pairs, with a Brier mean of 0.104 and no abstention.
Token counts, provider costs, and comparative latency remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 35 landings + implementer, 10 charges + reviewer, 15 axes | On `calibrated-decisions` the low orchestrator ran seven delegated tickets, derived the checkpoint chain rule from one refusal, and refuted two Coverage findings from the pre-repair file. Its medium first-pass axes recorded 14 labeled findings, and six of them were refuted; its line declaration expected one repair round and the build took four. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.198 over 14 pairs; 0 abstained |
| Fable / high | reviewer, 3 axes on 1 landing | On `ft311-recoverable-reset` the three axes found the below-path collision, the ignore-rule drift that deleted bytes, the hidden index flags, and the primary-side checkpoint resolution, each with an executed probe, and the Standards axis enumerated every one-source duplicate with its callers. | Review axes over a candidate another provider built, when the reviewer names the tier | unknown |
| Opus / high | implementer, latest 10 guidance and Go-seam charges; cross-harness reviewer, latest 2 passes | On `calibrated-decisions` seven authors returned 36 done-claim rows, and every row held; three tickets landed first-pass and ticket 7 took two repair cycles for a second owner of the step narrowing and two silent probes. One author raised its confidences after return, and the coordinator refused the raise. | High for process lifecycle, cleanup authority, destructive commands, anchored guidance, and foundational Go seams. | 0.015 over 36 pairs; 0 abstained |
| Opus / medium, low | implementer, orchestrator, and reviewer combined | On `calibrated-decisions` medium first-pass axes over CD3, CD4, and CD5 raised 18 findings; the eight labeled pairs held three times and were refuted five times, with three refutations at confidence 6 to 8. Every accepted finding named its site and closed in one repair. | Medium for gates, conformance, guidance, canaries, repair, Coverage, triage, and the repair-scoped re-review. Low for exact tickets and Standards or Spec review. | 0.236 over 8 pairs; 0 abstained |
| Sonnet / high | orchestrator, 3 landings + later-pass reviewer, 14 axes | On `calibrated-decisions` the later passes at high reaffirmed every repaired chunk at its exact tip with zero findings, and the Standards pass recomputed the Brier total from six sources. One xhigh Coverage pass on CD2 raised two findings and one was refuted at confidence 8. | Later review passes after a repair, and orchestration; compare again after a fourth orchestrated build | high 0 labeled pairs; xhigh 0.40 over 2 pairs |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `calibrated-decisions` CD2b Coverage axis | Fable / medium / reviewer | The axis added silent probes, then described the post-probe tree as production behavior; the coordinator refuted both findings from the pre-repair file. |
| `calibrated-decisions` ticket 7 | Opus / high / implementer | The author folded a step-scoped anchor kind into the shared scope walk, and two repair cycles removed a second owner and added the biting tests. |
| `delegate-boot-cost` spec review | Opus / medium / reviewer | The round measured every command description, found the six the author missed, and named the cheapest wrong implementation per chunk with the row that would catch it. |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope, which no row had decided. |
| `ft311-landing-completion` Coverage axis | Opus / medium / reviewer | The axis probed a state file whose fence never closes with a scratch test, observed the next run refuse the document, and named the row that should exist. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- The top tier implements only when the reviewer names it for the run.
- Use Opus/high as the delegated ticket author, because 36 of 36 done-claim rows held on this run.
- Use Opus/medium for the first review pass and Sonnet/high for every later pass under the configured line.
- Use Fable/high to review another provider's candidate when the reviewer names the tier.
- Keep Opus/high for shared guidance and its standing cross-harness review.
- Use Sonnet/low for an exact ticket only when the coordinator probes the named mutation before commit.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Freeze a stated confidence at return time, and refuse a post-hoc raise.
- Read the pre-repair tree before you accept a Coverage finding that follows a silent probe.
- Apply the repair allowance to the chunk, and retain optional advice outside blocking findings.
- Extend an ownership fence only under the approved plan-expansion rule, and flag the change for veto.
- Land plan commits and main merges before the ticket merge, and land record commits after the chunk tip.
- Run the whole-tree gate after the last repair and before landing.
- Read the assignment census before landing releases a worktree.
- Return material acceptance shortfalls and undecided behavior to the reviewer.
- Treat the calibration cell as one input to the two-run rule, and move no tier on this single run.
- Preserve unknown costs and usage, and do not change model defaults from this single run.
