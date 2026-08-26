# Claude model scorecard

Last incorporated landing: `harness-capability-seam` (`1100d6bd`, 2026-08-26)
— Fable/low as orchestrator, Opus/high on two charges, Opus/medium on eight
ticket charges and five review passes, and Opus/low on two charges.

Twenty-eight completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low, high | orchestrator, 10 landings + implementer, 9 charges + reviewer, 2 specs | On `harness-capability-seam` at low effort it charged eleven tickets and probed every fold at a distinct site and kind. One vacuous wrapper-label probe exposed a missing test and a hand-kept check list, and a `main` diff against the frozen base caught a CI-script change before the landing. It ran eight worktree cleans in parallel and staled every plan. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, latest 10 medium charges + 10 low charges; reviewer, 5 passes | On `harness-capability-seam` eight medium charges and two low charges each landed first-pass on behavior, and two stopped correctly on a spec contradiction and an oracle narrowing. The three review axes returned seventeen cited findings that collapsed to eleven repair targets, and the re-review found one two-fault fixture. | Medium for gate and conformance logic, a new package's algorithm, guidance prose, and the review axes; low for a ticket from an exact spec at a known seam under a covering gate |
| Opus / medium | orchestrator, latest 10 of 14 landings | On `learnings-dated-line-visibility` it found untested marker leniency no review axis raised and pinned it. It reproduces every accepted finding and probes every done-claim at a distinct site and kind. | Continue; the distinct-site probe discipline holds |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. One ran `cd` into the pool path against its charge and reported it. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. The re-review verified both predicates and the delta. | Standing tier for the three review axes; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `exec-census` ticket 06 | Opus / high / implementer | The `census=<n>` key broke 35 pinned expectations in 11 files, nine outside the fence. The delegate applied the sweep plan-before-apply, weakened nothing, and reported the fence exception with the one part to revert. |
| `exec-census` review | Sonnet / high / reviewer | The Coverage axis wrote a throwaway test for a text that names two assignment ids, showed the second id uncounted, and cited the loop that returns on the first match. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure to `internal/gate/subject.go`, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `ft228-debug-restoration` IP10 | Opus / medium / orchestrator | A coverage row required a description to carry a trigger and omit an anti-trigger. The check graded only the omission, and the delegate's own probe could not see the gap. A coordinator probe written from the row's text set the description to neither, and the gate stayed green. |
| `harness-capability-seam` ticket 06 | Opus / medium / implementer | The AXI envelope fixture demanded an empty state and a repo refusal that a compiled-in record cannot give. The delegate stopped with the contradiction and two options instead of inventing behavior, and built the declared opt-out after the reviewer chose it. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26).
- Every subagent runs Opus at low or medium effort (reviewer direction, 2026-08-26).
  On `harness-capability-seam` this covered guidance prose and the review axes, and
  every charge landed first-pass on behavior. This is one run; the Sonnet review
  routing and the Opus/high prose routing stay recorded until a second run compares.
- Route a ticket from the decision table at charge time, not from its story's line. An
  exact spec, a known seam, and a covering gate make the cheap row.
- A charge names the package that owns each fact the ticket touches, and it names every
  registry the family already appears in. On `harness-capability-seam` ticket 06 named
  five of eight registries, and the two it missed cost two round trips.
- A charge that moves a test or reflows anchored prose names the canary fixtures that
  anchor on it. Its focused checks run the fixture-bite test.
- Independent tickets run in parallel under disjoint fences; dependent tickets share the
  retained integration source. At most two delegates run tests at once, each at
  `-parallel 2`.
- The coordinator probes every done-claim at a distinct site and kind, through the
  writer the test captures. A probe that stays green is a finding, not a pass: on
  `harness-capability-seam` one such probe exposed a missing test and a hand-kept list.
- A hand-verified acceptance row is not closed. The coordinator asks for the test name
  before it folds a done-claim.
- Before a landing, the coordinator compares `main` against the frozen base and grades
  every file a new oracle check reads.
- A red is attributed or it is not resolved. A delegate's live-root red from its own
  stale checkout is re-run on the source before it is believed.
- A charge names `bench worktree exec` as the only command form. It forbids a
  background test run and a build into the worktree root.
- The coordinator takes a commit's path list from the checkout's status, never by hand.
  It never creates or cleans a worktree while a gate runs, and it runs per-path
  worktree cleans serially.
