# Claude model scorecard

Last incorporated landing: `ft311-recoverable-reset` (`5201c9fc21f0dd85563edb704f401294816edeb7`, 2026-09-11). Fable repaired the frozen benchmark candidate in its own context at the reviewer's direction. It ran three Fable/high review axes and three Opus/medium re-review axes, and it landed through an integration assignment.

The Fable/high round returned 26 findings and 14 accepted targets over the candidate, including three behavior defects the candidate's own Terra/medium review had missed. The Opus/medium re-review verified every fold predicate, observed five mutations red, and returned 7 minor findings. No write delegate ran.

Seventy-two completed landings are recorded. Routing follows the harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 34 landings + implementer, 10 charges + reviewer, 6 axes | On `ft311-recoverable-reset` it repaired the candidate in its own context, proved every new test red under a named mutation, and folded two review rounds. It paid two landing refusals before it learned that a sibling-created assignment lands through an integration assignment, and it let the landing retire the frozen sibling without stating that first. | Coordination of a parallel build and adversarial spec review; it implements only when the reviewer names it |
| Fable / high | reviewer, 3 axes on 1 landing | On `ft311-recoverable-reset` the three axes found the below-path collision, the ignore-rule drift that deleted bytes, the hidden index flags, and the primary-side checkpoint resolution, each with an executed probe, and the Standards axis enumerated every one-source duplicate with its callers. | Review axes over a candidate another provider built, when the reviewer names the tier |
| Opus / high | implementer, latest 10 guidance and Go-seam charges | On `ft311-landing-completion` it folded the reference, the final-check guidance, the working agreement, and the changelog first-pass with zero prose findings, moved the anchor with its fixture, and stopped at the fence when a second canary fixture anchored on the replaced sentence. | High for process lifecycle, cleanup authority, destructive commands, anchored guidance, and foundational Go seams. |
| Opus / medium, low | implementer, orchestrator, and reviewer combined | On `ft311-recoverable-reset` three medium re-review axes verified every fold predicate with citations, ran five restored mutations, and found one real two-hop recovery gap that RR78 now pins. On `ft311-landing-completion` six medium tickets landed with biting self-probes. | Medium for gates, conformance, guidance, canaries, repair, Coverage, triage, and the repair-scoped re-review. Low for exact tickets and Standards or Spec review. |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope, which no row had decided. |
| `ft311-recoverable-reset` Coverage re-review | Opus / medium / reviewer | The axis followed the exit-3 record's own `next` command through its apply and found the drifted bytes one envelope away with no signal in the record. |
| `ft311-landing-completion` Coverage axis | Opus / medium / reviewer | The axis probed a state file whose fence never closes with a scratch test, observed the next run refuse the document, and named the row that should exist. |
| `ft311-diagnostics` ticket 3 | Opus / medium / implementer | The delegate stopped diff-ready at its fence when eleven command-level tests did not fit one file, and it named the exact mechanical move. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The delegate found that one environment change crossed three packages and stopped at the fence. |

## Current decisions

- Route changes only after two comparable runs or one controlled model comparison.
- The top tier implements only when the reviewer names it for the run.
- Use Opus/medium for exact-spec Go tickets, for the three independent review axes, and for the repair-scoped re-review.
- Use Fable/high for the review axes over a candidate another provider built, when the reviewer names the tier.
- Use Opus/high for guidance prose under the leverage override.
- Use Sonnet/low for an exact ticket only when the coordinator probes the row's named mutation before it commits.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Run the build preflight after every ticket commit on the integration source and before the next charge.
- Extend a ticket fence in range for a changed record's reader, a test split, or a one-source seam, and flag it for veto.
- Run a whole-tree gate after the last repair commit and before landing.
- Read the assignment census before landing releases a worktree.
- Probe every done-claim at a distinct site and mutation kind, and prove every new test red under a named mutation before its commit.
- Route material acceptance shortfalls and undecided behavior partitions back to the reviewer in one question with a recommendation each.
- Land a repair of a frozen sibling from an integration assignment created from main, and state before the landing that the landing retires the sibling.
- Declare a self-ignoring capture directory in the build-outputs file before the first landing that meets it.
- Create no sibling worktree while another session is inside a landing tail, or between a landing's publication and its effects.
