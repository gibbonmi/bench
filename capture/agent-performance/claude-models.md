# Claude model scorecard

Last incorporated landing: `ft311-diagnostics` (`2578c9856448c93ccadadd853491066cc85c336f`, 2026-09-09). Fable, at medium effort, orchestrated the build.

Opus/medium implemented tickets 2, 3, 4, 5, and 7 and ran the three axes and the re-review. Opus/high folded the guidance. Sonnet/low implemented ticket 1.

Sixty-five completed landings are recorded. Routing follows the harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 32 landings + implementer, 9 charges + reviewer, 3 specs | On `ft311-diagnostics` at medium effort it ran eight write charges and four review passes, probed every return at a distinct site and kind, and collapsed 15 findings to 9 targets in one round. It extended the ticket fences eight times in range and did not run the build preflight after the first commit, so one charge failed on a removed fixture path. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 guidance and Go-seam charges | On `ft311-diagnostics` it folded the reference, the glossary, and the changelog first-pass, kept every pinned sentence byte-identical, and audited the landed diff for out-of-scope changes. | High for process lifecycle, cleanup authority, destructive commands, anchored guidance, and foundational Go seams. |
| Opus / medium, low | implementer, orchestrator, and reviewer combined | On `ft311-diagnostics` five medium write charges landed, three first-pass, and two stopped at the fence instead of a write outside it. The three medium axes returned 15 cited findings, one of them a behavior defect in the help's own example, and the re-review was clean. | Medium for gates, conformance, guidance, canaries, repair, Coverage, and triage. Low for exact tickets and Standards or Spec review. |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing |
| Sonnet / high, xhigh | reviewer, 3 axes on 13 landings + 12 scoped re-reviews + 1 spec round | On `structural-refactor-pass` one xhigh round over the spec and nine tickets resolved all 56 cited test names, verified four decisions against the code, and returned one blocking Coverage finding: a moved scan would drop its active-state filter in silence. | Spec-and-tickets review round when the reviewer names it; the review axes stay with Opus |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft311-diagnostics` Coverage axis | Opus / medium / reviewer | The axis built the tip and ran the help's own `bench gate-prose . --staged` example, which refused on every relative root, and it named the unbounded index blob read. |
| FT311 native axes and triage | Opus / medium / reviewer | Three axes found duplicate rendering, weak cell assertions, a missing evidence consumer, and missing durable live evidence. DP24 kept incomplete evidence explicit. |
| `ft311-diagnostics` ticket 3 | Opus / medium / implementer | The delegate stopped diff-ready at its fence when eleven command-level tests did not fit one file, and it named the exact mechanical move. |
| `bench-probe` Coverage axis | Opus / medium / reviewer | The axis combined a control-byte subject with a failed restore and exposed a missing copy identity. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The delegate found that one environment change crossed three packages and stopped at the fence. |

## Current decisions

- Route changes only after two comparable runs or one controlled model comparison.
- The top tier implements only when the reviewer names it for the run.
- Use Opus/medium for exact-spec Go tickets and for the three independent review axes and the repair-scoped re-review.
- Use Opus/high for guidance prose under the leverage override.
- Use Sonnet/low for an exact ticket only when the coordinator probes the row's named mutation before it commits.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Run the build preflight after every ticket commit on the integration source and before the next charge.
- Extend a ticket fence in range for a test split or a one-source seam, commit the amendment, and flag it for veto.
- Run a whole-tree gate after the last repair commit and before landing.
- Read the assignment census before landing releases a worktree.
- Probe every done-claim at a distinct site and mutation kind.
- Require each probe to compile, bite, and restore before its verdict counts.
- Route material acceptance shortfalls back to the reviewer.
