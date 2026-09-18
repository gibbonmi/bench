# Claude model scorecard

Last incorporated landing: `bounded-charge-evidence` (`2e30834b78eb828d19a5fca6bd41450ab1489643`, 2026-09-18).
Opus/high authors landed seven tickets in seven chunks, and every chunk needed at least one repair cycle.
Opus/medium ran every review axis and most repairs, and the last coordinator session ran on Fable.
The landing recorded 5 labeled pairs, with a Brier mean of 0.27 and 11 abstentions.
Token counts, provider costs, and comparative latency remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 36 landings + implementer, 10 charges + reviewer, 15 axes | On `bounded-charge-evidence` the orchestrator reproduced each accepted finding with its own probe, cleared two checkpoint refusals with one more review round each, and produced both transport records. It ran the first Codex decoder itself, which the Spec axis returned, and its charges asked for no stated confidence. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.198 over 14 pairs; 0 abstained |
| Fable / high | reviewer, 3 axes on 1 landing | On `ft311-recoverable-reset` the three axes found the below-path collision, the ignore-rule drift that deleted bytes, the hidden index flags, and the primary-side checkpoint resolution, each with an executed probe, and the Standards axis enumerated every one-source duplicate with its callers. | Review axes over a candidate another provider built, when the reviewer names the tier | unknown |
| Opus / high | implementer, latest 10 guidance and Go-seam charges; cross-harness reviewer, latest 2 passes | On `bounded-charge-evidence` the ticket 7 author returned five rows that all held, with nine recorded reds, and its first return still carried a second phase inventory that one repair cycle removed. One high Coverage pass on a light path found a rule that contradicted the chunk policy. | High for process lifecycle, cleanup authority, destructive commands, anchored guidance, and foundational Go seams. | 0.017 over 37 pairs; 2 abstained |
| Opus / medium, low | implementer, orchestrator, and reviewer combined | On `bounded-charge-evidence` the medium axes of CE-C3 and CE-C4 raised eight findings after the first repair, and the coordinator reproduced all eight. One Spec axis reported existing tests as absent, because it searched one package. One medium repair closed 13 targets with 16 recorded reds and left the worktree seal red, and both Opus line declarations expected one repair round against two and four. | Medium for gates, conformance, guidance, canaries, repair, Coverage, triage, and the repair-scoped re-review. Low for exact tickets and Standards or Spec review. | 0.262 over 12 pairs; 9 abstained |
| Sonnet / high | orchestrator, 3 landings + later-pass reviewer, 14 axes | On `calibrated-decisions` the later passes at high reaffirmed every repaired chunk at its exact tip with zero findings, and the Standards pass recomputed the Brier total from six sources. One xhigh Coverage pass on CD2 raised two findings and one was refuted at confidence 8. | Later review passes after a repair, and orchestration; compare again after a fourth orchestrated build | high 0 labeled pairs; xhigh 0.40 over 2 pairs |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `calibrated-decisions` CD2b Coverage axis | Fable / medium / reviewer | The axis added silent probes, then described the post-probe tree as production behavior; the coordinator refuted both findings from the pre-repair file. |
| `calibrated-decisions` ticket 7 | Opus / high / implementer | The author folded a step-scoped anchor kind into the shared scope walk, and two repair cycles removed a second owner and added the biting tests. |
| `bounded-charge-evidence` CE-C4 Coverage axis | Opus / medium / reviewer | The axis narrowed one needle to its first sentence, saw no red, and traced the hole to a paragraph pin that selected only two leads. |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope, which no row had decided. |
| `ft311-landing-completion` Coverage axis | Opus / medium / reviewer | The axis probed a state file whose fence never closes with a scratch test, observed the next run refuse the document, and named the row that should exist. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- The top tier implements only when the reviewer names it for the run.
- Use Opus/high as the delegated ticket author, and use a fresh writer when a fork would inherit another tier.
- Use Opus/medium for the review axes of a chunk, and continue the same axis session for a confirming round.
- Use Fable/high to review another provider's candidate when the reviewer names the tier.
- Keep Opus/high for shared guidance and its standing cross-harness review.
- Use Sonnet/low for an exact ticket only when the coordinator probes the named mutation before commit.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Ask for a stated confidence in every write charge and every axis charge, and freeze it at return time.
- Reproduce a finding with a probe before you charge its repair.
- Read the pre-repair tree before you accept a Coverage finding that follows a silent probe.
- Apply the repair allowance to the chunk, and retain optional advice outside blocking findings.
- Freeze a chunk base at the close commit of the predecessor chunk.
- Join the final reconciliation commit to the last chunk review before the completion checkpoint.
- Run the whole-tree gate after the last repair and before landing.
- Read the assignment census before landing releases a worktree.
- Return material acceptance shortfalls and undecided behavior to the reviewer.
- Preserve unknown costs and usage, and do not change model defaults from this single run.
