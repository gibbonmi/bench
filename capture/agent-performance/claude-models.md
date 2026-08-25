# Claude model scorecard

Last incorporated landing: `lifecycle-refusal-names-component` (`16bb950f`, 2026-08-25) —
Opus/medium as implementer on three tickets and Sonnet/high across four review
passes, under a Fable orchestrator.

Twenty-four completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high, medium | orchestrator, 6 landings + implementer, 8 Go charges at medium + reviewer, 2 specs | On `lifecycle-refusal-names-component` it probed every done-claim at a distinct site and kind, and all four bit. Its first charge offered a `cd` the delegate took, and its review artifact paid two red gates to its own prose bounds. | Top tier for guidance prose, for coordination when the spec's own claims are suspect, and for adversarial spec review |
| Opus / medium, low | implementer, latest 10 medium charges + 11 low charges | On `lifecycle-refusal-names-component` all three medium charges landed first-pass, and the ticket-02 charge flagged the precedence change its own composition caused. Low effort landed 11 of 11 first-pass on exact-spec tickets. | Medium for gate and conformance logic; low for a ticket from an exact spec at a known seam under a covering gate |
| Opus / high | implementer, latest 10 charges (lifecycle, guidance prose, Go-seam rewrites) | All landed first-pass. On `ft234-pool-key-reclaim` the destructive-apply charge showed its own acceptance row's seam was unreachable instead of writing a test that could not fail. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium | orchestrator, latest 10 of 14 landings | On `learnings-dated-line-visibility` it found untested marker leniency no review axis raised and pinned it. It reproduces every accepted finding and probes every done-claim at a distinct site and kind. | Continue; the distinct-site probe discipline holds |
| Sonnet / medium–high | implementer, latest 10 of 55 ticket-sized charges | On `module-size-split` (controlled A/B against Opus/low) both file-split charges landed first-pass, but one symbol went to the wrong group. Each charge cost about twice Opus/low's tokens. | Viable for behavior-preserving refactors at an effort that tracks the seam's risk; Opus/low for exact-spec file moves |
| Sonnet / low–high | reviewer, 3 axes on 10 landings + 8 scoped re-reviews | On `lifecycle-refusal-names-component` (high) the axes returned 6 raw findings and 5 repair targets, all cited, and the scoped re-review verified all five predicates. One axis marked an exact-predicate finding `ask-user`, and the coordinator re-disposed it as `auto-fix`. | Standing tier for the three review axes; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft228-debug-restoration` IP10 | Opus / medium / orchestrator | A coverage row required a description to carry a trigger and omit an anti-trigger. The check graded only the omission, and the delegate's own probe could not see the gap. A coordinator probe written from the row's text set the description to neither, and the gate stayed green. |
| `module-size-split` batch 1 | Opus / low vs Sonnet / medium / implementer | The same two file-split tickets ran on both lines. Every automated row passed on all four trees, so only the grouping axis separated them. Opus/low placed each overflow symbol with its owner group; Sonnet/medium parked one in the fixtures file, at about twice the tokens. |
| `learnings-dated-line-visibility` round 1 | Fable / high / reviewer | The spec anchored a parser rule on a literal that a scaffold generator also emits, above the first `## ` line. The generator prints that literal below its own `## ` example, so the rule could never open on a fresh journal. Both files had been read; neither was walked against the other. |
| `lifecycle-refusal-names-component` | Fable / medium / orchestrator | Every delegate claim was verified by `git status`, the focused tests, and a probe of a different kind and site; all four probes bit. The coordinator's own artifacts paid the avoidable cost: one `cd` offer in a charge and two red gates on prose bounds. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- Review axes route to Sonnet (reviewer's direction, 2026-08-19). Ten landings held
  the citation standard at three axes each with no re-run. Watch whether an axis
  under-reads a string-expectation seam.
- Opus/low beats Sonnet/medium on exact-spec Go file moves: equal correctness, better
  grouping, half the tokens. A charge for a file split names every check that reads
  the file by path.
- Route a ticket from the decision table at charge time, not from its story's line. An
  exact spec, a known seam, and a covering gate make the cheap row.
- Opus/high is the routing for process-lifecycle, signal, cleanup-authority, and
  foundational Go-seam work until a medium-effort sample contradicts it.
- The coordinator probes every done-claim at a distinct site and kind. At least one
  probe per acceptance row comes from the row's own text.
- A red is attributed or it is not resolved. A delegate's red-then-green report is an
  attribution claim, not a flake.
- A charge names `bench worktree exec` as the only command form. It names every
  state that reaches a branch, not only the state the acceptance row tests.
- Run the focused prose check before `bench commit` on authored Markdown. The
  orchestrator paid two red gates to its own bounds on each of the last three landings.
