# Claude model scorecard

Last incorporated landing: `bounded-repair-policy` (`c89922c93d75d34840863ca1a4d58cb217c3b6e4`, 2026-09-13).
Opus/high supplied the standing cross-harness review in a fresh Claude session.
It found two policy defects after the three native axes passed, then verified both repairs in another fresh session.
Token counts, provider costs, and comparative latency remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 34 landings + implementer, 10 charges + reviewer, 6 axes | On `ft311-recoverable-reset` it repaired the candidate in its own context, proved every new test red under a named mutation, and folded two review rounds. It paid two landing refusals before it learned that a sibling-created assignment lands through an integration assignment, and it let the landing retire the frozen sibling without stating that first. | Coordination of a parallel build and adversarial spec review; it implements only when the reviewer names it |
| Fable / high | reviewer, 3 axes on 1 landing | On `ft311-recoverable-reset` the three axes found the below-path collision, the ignore-rule drift that deleted bytes, the hidden index flags, and the primary-side checkpoint resolution, each with an executed probe, and the Standards axis enumerated every one-source duplicate with its callers. | Review axes over a candidate another provider built, when the reviewer names the tier |
| Opus / high | implementer, latest 10 guidance and Go-seam charges; cross-harness reviewer, latest 2 passes | On `bounded-repair-policy`, Opus/high found a missing resume trigger and a glossary link unavailable to linked repositories. A fresh pass verified both repairs and kept optional advice outside the findings. | High for process lifecycle, cleanup authority, destructive commands, anchored guidance, and foundational Go seams. |
| Opus / medium, low | implementer, orchestrator, and reviewer combined | On the `delegate-boot-cost` spec one medium round read 19 enforcement files, found six over-budget command files the author's census had skipped, and named the hard-coded-limit degenerate that no row caught. All twelve findings held against the tree. | Medium for gates, conformance, guidance, canaries, repair, Coverage, triage, and the repair-scoped re-review. Low for exact tickets and Standards or Spec review. |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `bounded-repair-policy` | Opus / high / cross-harness reviewer | It found two linked-workflow defects after all native axes passed, then verified closure on the repaired source. |
| `delegate-boot-cost` spec review | Opus / medium / reviewer | The round measured every command description, found the six the author missed, and named the cheapest wrong implementation per chunk with the row that would catch it. |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope, which no row had decided. |
| `ft311-recoverable-reset` Coverage re-review | Opus / medium / reviewer | The axis followed the exit-3 record's own `next` command through its apply and found the drifted bytes one envelope away with no signal in the record. |
| `ft311-landing-completion` Coverage axis | Opus / medium / reviewer | The axis probed a state file whose fence never closes with a scratch test, observed the next run refuse the document, and named the row that should exist. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- The top tier implements only when the reviewer names it for the run.
- Use Opus/medium for exact-spec Go tickets and independent review axes under the configured line.
- Use Fable/high to review another provider's candidate when the reviewer names the tier.
- Keep Opus/high for shared guidance and its standing cross-harness review.
- Use Sonnet/low for an exact ticket only when the coordinator probes the named mutation before commit.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Verify the complete linked payload when a policy relies on another document.
- Apply the repair allowance to the chunk, and retain optional advice outside blocking findings.
- Extend an ownership fence only under the approved plan-expansion rule, and flag the change for veto.
- Run the whole-tree gate after the last repair and before landing.
- Read the assignment census before landing releases a worktree.
- Prove new omission expectations red while the production obligation is absent.
- Return material acceptance shortfalls and undecided behavior to the reviewer.
- Create no sibling worktree during a landing or its cleanup effects.
- Preserve unknown costs and usage, and do not change model defaults from this single run.
