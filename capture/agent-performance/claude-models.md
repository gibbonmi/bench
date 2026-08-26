# Claude model scorecard

Last incorporated landing: `exec-census` (`415c99e9`, 2026-08-26) — Fable/low
as orchestrator, Opus/high on four charges, Opus/medium on two charges,
Sonnet/low on three ticket charges, and Sonnet/high on four review passes.

Twenty-seven completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low, high | orchestrator, 9 landings + implementer, 9 charges + reviewer, 2 specs | On `exec-census` at low effort it charged nine tickets, folded four one-source duplications the delegates left, and probed every fold at a distinct site and kind. It attributed one canary red to the right ticket, redid one invalid probe, and ran raw `cd` calls into the pool path itself. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `exec-census` the signal, the landing drop, the guidance prose, and the anchor repair each landed first-pass on behavior. The landing charge widened its fence for a 35-expectation sweep and reported it; the guidance charge reflowed one line a canary fixture anchored. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, latest 10 medium charges + 11 low charges | On `exec-census` the two prefactor charges landed on behavior first-pass, and each took one repair round to fold a derivation it had restated beside `poolkey`. | Medium for gate and conformance logic and a new package's algorithm; low for a ticket from an exact spec at a known seam under a covering gate |
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
| `one-change-one-grade` ticket 01 | Sonnet / low / implementer | The delegate's package tests passed while the live-root routing check would have gone red, because the exhaustive registry runs only under the conformance root env. The delegate named the registry it could not edit; the coordinator's live-root run confirmed the red and added the one row. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26). Guidance prose routes to Opus/high.
- Review axes route to Sonnet (reviewer's direction, 2026-08-19). Thirteen landings held
  the citation standard at three axes each with no re-run.
- Route a ticket from the decision table at charge time, not from its story's line. An
  exact spec, a known seam, and a covering gate make the cheap row.
- Opus/high is the routing for process-lifecycle, signal, cleanup-authority, and
  foundational Go-seam work. On `exec-census` all four such charges landed first-pass.
- A charge names the package that owns each fact the ticket touches. On `exec-census`
  four of nine charges restated a fact beside its owner, at every tier.
- A charge that moves a test or reflows anchored prose names the canary fixtures that
  anchor on it. Its focused checks run the fixture-bite test.
- Independent tickets run in parallel under disjoint fences; dependent tickets share the
  retained integration source. At most two delegates run tests at once, each at
  `-parallel 2`.
- The coordinator probes every done-claim at a distinct site and kind, through the
  writer the test captures. At least one probe per acceptance row comes from the row's
  own text.
- A red is attributed or it is not resolved. A one-off `infrastructure` refusal under
  full gate load is reported with its reproduction attempts, never called transient.
- A charge names `bench worktree exec` as the only command form. It forbids a
  background test run and a build into the worktree root. It names every registry the
  family already appears in.
- The coordinator takes a commit's path list from the checkout's status, never by hand,
  and never creates or cleans a worktree while a gate runs.
