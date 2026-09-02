# Claude model scorecard

Last incorporated landing: `kit-guidance-fold` (`6e51b1ec`, 2026-09-02) and the five
light-path fix landings of 2026-09-01 and 2026-09-02. Fable/low ran as orchestrator
over five Opus/medium spec write charges, five Opus fix charges, and one Opus/medium
repair charge. It also ran three Opus/medium review axes, two Opus scoped re-reviews,
and three Opus/low read delegates. Every spec write charge landed first-pass on
behavior with a biting self-probe, and the coordinator's heading-rename probe bit at
a distinct site each time.

Fifty-one completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 21 landings + implementer, 9 charges + reviewer, 3 specs | On `kit-guidance-fold` at low effort it ran a five-ticket chain, five parallel light-path fixes, one review round, and two scoped re-reviews. It classified two merge-gate reds by rerun and read, and it moved the landing base to the merged `main` tip when the fence refused the merged-in files. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 59 review axes, 12 of 18 orchestrated landings | On `kit-guidance-fold` five medium prose charges landed first-pass with a red-to-green log per row; two stopped at a fence and reworded instead of editing a fixture. The medium repair charge folded nine findings in one pass. One fix charge missed a pin in a package outside its search list. The three axes returned 19 cited findings that collapsed to nine targets, and the Coverage axis measured all 27 fixtures by script. | Medium for gate and conformance logic, guidance prose, canary fixtures, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 2 landings | On `citation-phase-package-scope` it independently re-ran every delegate's tests, vet, and gofmt before trusting a done-claim; caught its own premature spec-status flip before it compounded; and ran two full three-axis review rounds that found 13 findings, most cited by more than one axis independently, before landing. | Continues to hold at high effort; compare again after a third orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. | Standing tier for the three review axes before the Opus rule; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `kit-guidance-fold` Coverage axis | Opus / medium / reviewer | The axis scripted every one of the 27 new fixtures against its live file, found the one tuple with no fixture and no test arm, and cited the producer before it claimed a gap. |
| `kit-guidance-fold` Standards axis | Opus / medium / reviewer | The axis found a seven-sentence paragraph the prose check split at a label-shaped line, and it named the single edit that defeats the six-sentence bound. |
| `spec-authoring-discipline` reader-sweep ticket | Opus / medium / implementer | The charge asked for a byte-exact move of a pinned sentence. The delegate found two conformance tests outside the fence that pinned the bytes, stopped with both citations instead of editing out of fence, and finished first-pass once the coordinator placed the split arm after the pinned bytes. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `citation-phase-package-scope` repair-scoped re-review | Opus / medium / reviewer | The repair-scoped re-review found a real bug the first repair round had itself introduced: a cross-architecture phase left `CgoEnabled` at the host's value, a false green a fresh three-axis pass would have missed entirely without a scoped second look. |

## Current decisions

- Routing changes only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run.
- Every subagent runs Opus at low or medium effort. Fifty-one landings now run under
  this rule.
- A light-path ticket from an exact ticket file runs Opus at low effort. A light-path
  ticket that changes CLI output or a policy constant runs Opus at medium.
- Opus at medium serves the review axes and the scoped re-reviews.
- A kit-guidance diff takes the standing Codex falsification pass at the mid tier
  beside the three axes. On this landing it found the one term collision the axes
  missed.
- `bench worktree land --spec <slug>` goes only on the landing that completes a
  spec's final ticket, never on an earlier landing under the same fence.
- A repair-scoped re-review runs after every repair charge, even a one-line repair.
- The coordinator reads the whole census record before `bench worktree land`,
  because the release deletes it.
- The coordinator probes every done-claim at a distinct site and a distinct kind. It
  confirms with `cmp` that the probe changed bytes before it reads any verdict. A
  heading rename kept constant across a batch is a valid coordinator probe kind.
- A cap or policy constant change charges the delegate to sweep every package for the
  old literal, not only the packages that consume the constant.
- Independent tickets run in parallel in separate worktrees; dependent tickets branch
  from the integration tip. At most two delegates run tests at once.
- A spec landing after sibling light-path landings sets its review base to the merged
  `main` tip. The reviewed range then holds the spec diff alone.
