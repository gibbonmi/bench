# Claude model scorecard

Last incorporated landing: `prospective-artifact-recovery` (`0ec709aa`,
2026-08-28) — Fable/medium as orchestrator over four Opus write charges, six
Opus review axes, and three scoped re-review axes. Every charge landed
first-pass on behavior; the review found two real defects in the prior
session's ticket.

Thirty-two completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 12 landings + implementer, 9 charges + reviewer, 3 specs | On `prospective-artifact-recovery` at medium effort it refuted a delegate's false absence claim, attributed two environmental reds before it believed them, and bit four distinct probe sites. It charged two gate-package test files without the branch-native census, and the landing went red once for that omission. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 25 review axes, 11 of 17 orchestrated landings | On `prospective-artifact-recovery` all four charges landed first-pass on behavior with biting self-probes, and one reported an unreachable spec seam instead of working around it. One charge duplicated a record shape across three test files. Six review axes found 15 raw findings in 4 targets, and the Coverage axis proved the worst one with a real repository run. | Medium for gate and conformance logic, a new package's algorithm, guidance prose, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 1 landing | On `bounded-gate-output` it declared the line once for eleven charges, ran a fresh mutation probe distinct from each delegate's own before every commit, and caught a coverage-map amendment with no owning ticket file before it reached preflight. Two rounds of parallel charges on disjoint files landed with no collision. | New this landing; continue and compare after a second orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. One ran `cd` into the pool path against its charge and reported it. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. The re-review verified both predicates and the delta. | Standing tier for the three review axes; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `exec-census` ticket 06 | Opus / high / implementer | The `census=<n>` key broke 35 pinned expectations in 11 files, nine outside the fence. The delegate applied the sweep plan-before-apply, weakened nothing, and reported the fence exception with the one part to revert. |
| `prospective-artifact-recovery` review | Opus / medium / reviewer | The Coverage axis registered a worktree through a symlinked path in a real repository, showed Git records the resolved spelling, and cited the comparison that could never match. The Spec axis traced the one branch that let an authored binary escape the bundle. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure to `internal/gate/subject.go`, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `ft228-debug-restoration` IP10 | Opus / medium / orchestrator | A coverage row required a description to carry a trigger and omit an anti-trigger. The check graded only the omission, and the delegate's own probe could not see the gap. A coordinator probe written from the row's text set the description to neither, and the gate stayed green. |
| `native-runtime-evidence-reduction` review | Opus / medium / reviewer | The Spec and Coverage axes each found, without seeing the other, that the gate graded a workflow matrix consumer and never its producer. One token on the producer step restarted the runners the landing existed to remove, with the gate green. |
| `native-runtime-evidence-reduction` ticket 05 | Opus / high / implementer | The charge listed four consumers of a changed upload root. A fifth existed. The delegate finished its own scope, then stopped and reported the gap as a spec contradiction rather than editing a job no ticket owned. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26).
- Every subagent runs Opus at low or medium effort (reviewer direction, 2026-08-26).
  Four landings now run under this rule. On `native-runtime-evidence-reduction` six of
  eight charges landed first-pass, and both misses came from the spec. The rule stands.
- Opus at low and medium serves the review axes as well as the write charges. The
  Sonnet review routing is superseded.
- The top tier reviews a spec and its tickets as one round. On
  `native-runtime-evidence-reduction` that round blocked on three findings, and two of
  them were readers of a changed value that the spec never listed.
- A gate check on a derived or indirected value names both ends. A check that asserts a
  consumer's reference to an output proves nothing about what fills that output. Two
  independent axes found this hole on `native-runtime-evidence-reduction`, and the
  coordinator had missed it twice.
- Before the coverage map locks, sweep the whole tree for readers of every value the
  change derives. Four of five repair rounds on `native-runtime-evidence-reduction`
  traced to unlisted readers.
- A Standards finding that deletes a test is a coverage question first. One such finding
  removed the only assertion binding a plan's proven set, and the scoped re-review caught
  the loss.
- A charge that names an exact verification command must name one the coordinator has
  run. A ship-tier probe command missing its root variable skips in milliseconds and
  reports ok.
- Route a ticket from the decision table at charge time, not from its story's line. An
  exact spec, a known seam, and a covering gate make the cheap row.
- A charge names the package that owns each fact the ticket touches, and it names every
  registry the family already appears in.
- A charge that moves a test or reflows anchored prose names the canary fixtures that
  anchor on it. Its focused checks run the fixture-bite test.
- Independent tickets run in parallel under disjoint fences; dependent tickets share the
  retained integration source. At most two delegates run tests at once, each at
  `-parallel 2`.
- The coordinator probes every done-claim at a distinct site and kind, through the
  writer the test captures. A probe that stays green is a finding, not a pass.
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
- A spec amendment that adds acceptance rows during a repair pass needs a ticket
  file before the next `bench preflight review`.
- A charge that adds a test under an architecture-owned package names the
  branch-native census and runs it in the focused checks. Focused tests were green
  on `prospective-artifact-recovery`, and the landing went red on that census.
- A review axis refutes a strong finding with a real run before it reports. One
  repository probe made the symlinked-root defect undeniable.
