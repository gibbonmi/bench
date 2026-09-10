# Claude model scorecard

Last incorporated landing: `ft311-landing-completion` (`f50671018ec41711a89f8d73c0ae28df8af6fc71`, 2026-09-10). Fable, at medium effort, coordinated the build, the three-axis review, the repair, and the landing.

Opus/medium implemented five Go tickets and one repair ticket, and every self-probe and every coordinator probe bit. Four tickets needed a fence continuation for readers the spec sweep missed. Opus/high folded the guidance and moved one anchor first-pass, with one second canary fixture found by the fixture-bite test. Three Opus/medium axes returned 15 raw findings and 8 repair targets, and the scoped re-review was clean.

Seventy-one completed landings are recorded. Routing follows the harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 33 landings + implementer, 9 charges + reviewer, 3 specs | On `ft311-landing-completion` at medium effort it ran six write charges and four review passes, probed every return at a distinct site and kind, and ran the build preflight after every commit. It extended the fences seven times in range for readers the spec sweep missed, and it omitted the skip-ownership check from one charge, which the fold gate caught one round late. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 guidance and Go-seam charges | On `ft311-landing-completion` it folded the reference, the final-check guidance, the working agreement, and the changelog first-pass with zero prose findings, moved the anchor with its fixture, and stopped at the fence when a second canary fixture anchored on the replaced sentence. | High for process lifecycle, cleanup authority, destructive commands, anchored guidance, and foundational Go seams. |
| Opus / medium, low | implementer, orchestrator, and reviewer combined | On `ft311-landing-completion` six medium tickets landed with biting self-probes, and each stopped at its fence and named the exact out-of-fence edit instead of writing it. The three medium axes returned 15 raw findings and 8 targets, including a real text-sink gap, and the scoped re-review verified all six predicates. | Medium for gates, conformance, guidance, canaries, repair, Coverage, and triage. Low for exact tickets and Standards or Spec review. |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing |
| Sonnet / high, xhigh | reviewer, 3 axes on 13 landings + 12 scoped re-reviews + 1 spec round | On `structural-refactor-pass` one xhigh round over the spec and nine tickets resolved all 56 cited test names, verified four decisions against the code, and returned one blocking Coverage finding: a moved scan would drop its active-state filter in silence. | Spec-and-tickets review round when the reviewer names it; the review axes stay with Opus |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft311-diagnostics` Coverage axis | Opus / medium / reviewer | The axis built the tip and ran the help's own `bench gate-prose . --staged` example, which refused on every relative root, and it named the unbounded index blob read. |
| `ft311-landing-completion` Coverage axis | Opus / medium / reviewer | The axis probed a state file whose fence never closes with a scratch test, observed the next run refuse the document, and named the row that should exist. |
| `ft311-diagnostics` ticket 3 | Opus / medium / implementer | The delegate stopped diff-ready at its fence when eleven command-level tests did not fit one file, and it named the exact mechanical move. |
| `bench-probe` Coverage axis | Opus / medium / reviewer | The axis combined a control-byte subject with a failed restore and exposed a missing copy identity. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The delegate found that one environment change crossed three packages and stopped at the fence. |

## Current decisions

- Route changes only after two comparable runs or one controlled model comparison.
- The top tier implements only when the reviewer names it for the run.
- Use Opus/medium for exact-spec Go tickets, for the three independent review axes, and for the repair-scoped re-review.
- Use Opus/high for guidance prose under the leverage override.
- Use Sonnet/low for an exact ticket only when the coordinator probes the row's named mutation before it commits.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Run the build preflight after every ticket commit on the integration source and before the next charge.
- Extend a ticket fence in range for a changed record's reader, a test split, or a one-source seam, and flag it for veto.
- Tick the delegation discipline's In-the-charge list against the ticket's Writes before dispatch, so a skipping test names skip-ownership.
- Run a whole-tree gate after the last repair commit and before landing.
- Read the assignment census before landing releases a worktree.
- Probe every done-claim at a distinct site and mutation kind.
- Require each probe to compile, bite, and restore before its verdict counts.
- Route material acceptance shortfalls back to the reviewer.
- Create no sibling worktree while another session is inside a landing tail, or between a landing's publication and its effects.
