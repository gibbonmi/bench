# Claude model scorecard

Last incorporated landing: `path-aware-lane` (`fc3a4046`, 2026-08-30) —
Fable/medium as orchestrator over five Opus ticket charges (one low, four
medium) and one Opus/medium repair charge with one continuation. Three
Opus/medium review axes and one Opus/medium scoped re-review ran. All five ticket charges
landed first-pass on behavior. The review found 17 raw findings that collapsed
to six repair targets and five reviewer decisions; the re-review found one
citation defect.

Thirty-seven completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 16 landings + implementer, 9 charges + reviewer, 3 specs | On `path-aware-lane` at medium effort it probed five returned diffs at distinct sites and kinds, and every probe bit. It refuted the Coverage axis's worst finding with two live runs. Its own review artifact broke the prose lane on four sentences, and two fence gaps reached preflight and the landing before it amended the spec. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 37 review axes, 11 of 17 orchestrated landings | On `path-aware-lane` five of five ticket charges and the repair charge landed first-pass on behavior with biting self-probes; two charges reported an out-of-fence file instead of an edit. The three medium axes returned 17 cited findings with six repair targets, the Coverage axis's worst finding was refuted by live evidence, and the re-review caught a renamed test still cited by a coverage row. | Medium for gate and conformance logic, guidance prose, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate, with the red-first rule enforced |
| Sonnet / high | orchestrator, 1 landing | On `bounded-gate-output` it declared the line once for eleven charges, ran a fresh mutation probe distinct from each delegate's own before every commit, and caught a coverage-map amendment with no owning ticket file before it reached preflight. Two rounds of parallel charges on disjoint files landed with no collision. | New this landing; continue and compare after a second orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. One ran `cd` into the pool path against its charge and reported it. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. The re-review verified both predicates and the delta. | Standing tier for the three review axes; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `exec-census` ticket 06 | Opus / high / implementer | The `census=<n>` key broke 35 pinned expectations in 11 files, nine outside the fence. The delegate applied the sweep plan-before-apply, weakened nothing, and reported the fence exception with the one part to revert. |
| `prospective-artifact-recovery` review | Opus / medium / reviewer | The Coverage axis registered a worktree through a symlinked path in a real repository, showed Git records the resolved spelling, and cited the comparison that could never match. The Spec axis traced the one branch that let an authored binary escape the bundle. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure to `internal/gate/subject.go`, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `worktree-native-forms` re-review | Opus / medium / reviewer | The scoped re-review ran the WF3 row 500 times, found the TOON encoder quotes an id that starts with `0` and a digit, enumerated the two fixture sites at risk, and named the one-source fix through `toon.Table`. |
| `native-runtime-evidence-reduction` review | Opus / medium / reviewer | The Spec and Coverage axes each found, without seeing the other, that the gate graded a workflow matrix consumer and never its producer. One token on the producer step restarted the runners the landing existed to remove, with the gate green. |
| `native-runtime-evidence-reduction` ticket 05 | Opus / high / implementer | The charge listed four consumers of a changed upload root. A fifth existed. The delegate finished its own scope, then stopped and reported the gap as a spec contradiction rather than editing a job no ticket owned. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26).
- Every subagent runs Opus at low or medium effort (reviewer direction, 2026-08-26).
  Eight landings now run under this rule. On `path-aware-lane` five of five ticket
  charges landed first-pass. The rule stands.
- Opus at low and medium serves the review axes and the scoped re-reviews. On
  `path-aware-lane` the axes found an imprecise guidance claim, an unbound class
  table, and four undecided prefix edges; the re-review found a stale test citation.
- A Coverage finding on a hop the build ran live is refuted by that run before it
  becomes a repair. On `path-aware-lane` two live runs of `bench test --check`
  refuted the axis's worst finding.
- A promise in a spec's edge inventory gets a coverage row. On `path-aware-lane` the
  PL5 glob clause shipped without a row and became repair P1.
- The review pickup file and every file a repair adds join the spec fence before the
  next preflight or landing. On `path-aware-lane` the fence reddened twice.
- A repair that renames a test updates every coverage row that cites it, because
  `bench coverage --check` grades structure and not names.
- The coordinator reads the census record before `bench worktree land`, because the
  release deletes it.
- A prose file runs `bench gate-prose` before `bench commit`. On `path-aware-lane`
  four sentences in the pickup artifacts refused.
- A sibling worktree folds from the integration source only, or is created before
  `main` moves. Independent tickets run in parallel in separate worktrees and fold
  through `bench worktree merge`. Dependent tickets share the retained source. At most
  two delegates run tests at once, each at `-parallel 2`.
- The coordinator probes every done-claim at a distinct site and kind. A row that
  admits one mutation keeps the delegate's red proof as its evidence.
- A hand-verified acceptance row is not closed. A red is attributed or it is not resolved.
- A charge names `bench worktree exec` as the only command form and forbids a
  background test run. It names the package that owns each fact and every registry
  the family already appears in.
- The coordinator takes a commit's path list from the checkout's status, never by hand.
  It never creates or cleans a worktree while a gate runs.
- A repair-scoped re-review runs after every repair charge, because a repair induces
  defects. On `path-aware-lane` it found the renamed-test citation.
- The landing's `--base` is the assignment's recorded start or a descendant of it on
  `main`. The review preflight runs from inside the integration worktree.
- A spec with an open review pickup is not retired by the build session; the reviewer
  decides the pickup first.
- A cross-family top-tier reviewer (codex Sol/high) stays available for a large new
  package; compare it against Opus axes on the next one.
