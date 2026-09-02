# Claude model scorecard

Last incorporated landing: `roadmap-light-path-fixes-2` (`dbbec276`, spec-retire
`c04d3410`, 2026-09-02). Sonnet/high ran as orchestrator over ten Opus/medium
ticket charges (nine story tickets plus one repair-round ticket) and six
Opus/medium review axes (three initial, three repair-scoped). Six of ten
charges landed first-pass; three stopped and reported a material shortfall
instead of editing outside fence. The coordinator's ten independent mutation
probes each killed at a distinct site and kind from the delegate's own
self-probe.

Fifty-two completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 21 landings + implementer, 9 charges + reviewer, 3 specs | On `kit-guidance-fold` at low effort it ran a five-ticket chain, five parallel light-path fixes, one review round, and two scoped re-reviews. It classified two merge-gate reds by rerun and read, and it moved the landing base to the merged `main` tip when the fence refused the merged-in files. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 65 review axes, 12 of 18 orchestrated landings | On `roadmap-light-path-fixes-2` ten medium ticket charges self-diagnosed three material shortfalls rather than editing out of fence; one refuted the review's own claim with a hard trace (a cited literal was a named constant, not a numeric one) instead of accepting it. A repair-scoped Coverage axis found a real compile break, though it was the axis's own unrestored probe that caused it. | Medium for gate and conformance logic, guidance prose, canary fixtures, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
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
| `roadmap-light-path-fixes-2` sweep-widening ticket | Opus / medium / implementer | Charged to widen a check to three duration spellings, the delegate found two of the three would break a real test or contradict a closed Won't-handle, landed only the one that held, and named the review's own cited site as a false alarm with the traced reason why. |

## Current decisions

- Routing changes only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run.
- Every subagent runs Opus at low or medium effort. Fifty-two landings now run under
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
- A material acceptance shortfall a ticket surfaces mid-build routes to the reviewer
  as a question. It never gets a silent widening or a self-decided narrowing.
- A read-only review-axis delegate on a shared retained worktree can leave it dirty.
  The coordinator confirms `git status` clean before and after each axis dispatch,
  and restores any leftover probe itself before the next commit.
- `bench worktree merge --from <source> <target>` refuses when the source's tip
  carries an inherited gate red, even one the diff already has a ticket for. A
  dependent-but-independent ticket waits for that red to clear before it gets its own
  synced worktree.
