# Claude model scorecard

Last incorporated landing: `worktree-exec-comfort` (`ee0ed9d5`, 2026-08-29) —
Fable/medium as orchestrator over six Opus ticket charges (five low, one
medium), two Opus/low repair charges, three Opus/medium review axes, and one
Opus/medium scoped re-review. Four ticket charges landed first-pass on
behavior; the guard charge took three continuations. The review found 18 raw
findings that collapsed to 11 repair targets, and the re-review was clean.

Thirty-five completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 14 landings + implementer, 9 charges + reviewer, 3 specs | On `worktree-exec-comfort` at medium effort it probed each of nine returned diffs at a distinct site and kind, and every probe bit; one throwaway probe found an unpinned operator case. Its own spec prose broke the prose lane twice, and a repair worktree it created from a moved `main` carried five foreign commits into the source. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 29 review axes, 11 of 17 orchestrated landings | On `worktree-exec-comfort` seven low charges landed first-pass on behavior with biting self-probes, and the missing-tree charge corrected the charge's landed predicate inside its scope. The medium guard charge stopped at its `Writes:` fence and reported the hook edit, then took a repair for an unpinned operator case and a prose fix. The three medium axes returned 18 cited findings with 11 repair targets, and the re-review verified all eight predicates. | Medium for gate and conformance logic, guidance prose, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate, with the red-first rule enforced |
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
  Six landings now run under this rule. On `worktree-exec-comfort` seven of eight
  write charges landed first-pass. The guard charge's three continuations came
  from a `Writes:` gap, an unpinned row, and prose bounds. The rule stands.
- Opus at low and medium serves the review axes as well as the write charges. On
  `worktree-exec-comfort` the three medium axes found a `'` that broke a paste line.
  They also found a `--help` at exit 2 and a promised `lineSafe` check the code
  lacked. The Sonnet review routing is superseded.
- A prose file runs `bench gate-prose` before `bench commit`. On
  `worktree-exec-comfort` three lanes refused on prose bounds, one from a delegate
  and two from the orchestrator.
- A sibling repair worktree folds from the integration source only, or is created
  before `main` moves. On `worktree-exec-comfort` a worktree created from a moved
  `main` carried five foreign commits into the source, and the landing needed the
  moved review base.
- The top tier reviews a spec and its tickets as one round.
- A gate check on a derived or indirected value names both ends.
- Before the coverage map locks, sweep the whole tree for readers of every value the
  change derives.
- A Standards finding that deletes a test is a coverage question first.
- A charge that names an exact verification command must name one the coordinator has
  run.
- Route a ticket from the decision table at charge time, not from its story's line.
- A charge names the package that owns each fact the ticket touches, and it names every
  registry the family already appears in.
- A charge whose collapse crosses its Writes list inside the spec fence gets a fence
  extension in a continuation. On `worktree-exec-comfort` the guard charge stopped at
  `cmd/bench/main.go` and reported the edit, and the continuation delivered it with a
  hook-seam assertion.
- A charge that moves a test or reflows anchored prose names the canary fixtures that
  anchor on it.
- Independent tickets run in parallel in separate worktrees and fold through
  `bench worktree merge`; dependent tickets share the retained integration source. On
  `worktree-exec-comfort` three worktrees carried the frontier with no collision. At
  most two delegates run tests at once, each at `-parallel 2`.
- The coordinator probes every done-claim at a distinct site and kind. A probe that
  stays green is a finding, not a pass. On `worktree-exec-comfort` nine probes bit, and
  one throwaway probe found an unpinned case the delegate's rows missed.
- A hand-verified acceptance row is not closed. The coordinator asks for the test name
  before it folds a done-claim.
- A red is attributed or it is not resolved.
- A charge names `bench worktree exec` as the only command form. It forbids a
  background test run and a build into the worktree root.
- The coordinator takes a commit's path list from the checkout's status, never by hand.
  It never creates or cleans a worktree while a gate runs.
- A spec amendment that adds acceptance rows during a repair pass needs a ticket
  file before the next `bench preflight review`.
- A row whose named seam cannot reach the state amends its seam column as a recorded
  decision. On `worktree-merge` WM17 moved to a helper seam this way.
- A repair-scoped re-review runs after every repair charge, because a repair induces
  defects. On `worktree-exec-comfort` it verified eight predicates and found four
  non-blocking notes.
- The landing's `--base` is the assignment's recorded start or a descendant of it on
  `main`. The spec-review base is a source commit, and the landing refuses it. On
  `worktree-exec-comfort` the moved `main` tip was the base.
- A cross-family top-tier reviewer (codex Sol/high) stays available for a large new
  package; compare it against Opus axes on the next one.
