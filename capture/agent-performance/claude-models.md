# Claude model scorecard

Last incorporated landing: `worktree-test-floor` (`cc408adf`, 2026-08-25) —
Fable/high as orchestrator and spec amender, Opus/high on two Go-seam charges,
and Opus/medium on nine charges. Opus/low ran one charge, Sonnet/low ran ten
ticket charges, and Sonnet/high ran five review passes.

Twenty-six completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high, medium | orchestrator, 8 landings + implementer, 9 charges + reviewer, 2 specs | On `worktree-test-floor` it ran up to seven delegates in parallel worktrees, folded 19 diffs by patch, and every probe bit at a site distinct from the delegate's. It measured the wall after eleven tickets, found the spec premise false, and amended the spec three times. It caused two red gates itself: a hand-typed path list and a worktree created during a gate. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `worktree-test-floor` ticket 12 stopped before any code on a gate-closure seam the ticket had wrong, and landed first-pass after the re-scope. Ticket 13c moved seventeen injectables into one per-call value across 47 files and landed first-pass, and it attributed a one-off red instead of calling it transient. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, latest 10 medium charges + 11 low charges | On `worktree-test-floor` nine medium charges landed; one took a repair round for a census row the coordinator's probe found, and three stopped correctly on a defect outside their fence. The low charge threaded a root through four verbs and eleven test files first-pass. | Medium for gate and conformance logic; low for a ticket from an exact spec at a known seam under a covering gate |
| Opus / medium | orchestrator, latest 10 of 14 landings | On `learnings-dated-line-visibility` it found untested marker leniency no review axis raised and pinned it. It reproduces every accepted finding and probes every done-claim at a distinct site and kind. | Continue; the distinct-site probe discipline holds |
| Sonnet / low–high | implementer, latest 10 of 67 ticket-sized charges | On `worktree-test-floor` ten low charges marked tests parallel and converted binds under the census; all landed first-pass. Two stopped their turn on a background test run and needed a resume, and one called a one-off red transient. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges forbid background runs |
| Sonnet / low–high | reviewer, 3 axes on 12 landings + 11 scoped re-reviews | On `worktree-test-floor` (high) the axes returned 6 raw findings, all cited, collapsed to 3 repair targets; the Coverage axis found a census blind spot by reading the AST predicate. The scoped re-review probed six AST shapes and found one repair-induced gap. | Standing tier for the three review axes; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure to `internal/gate/subject.go`, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `worktree-test-floor` ticket 13 | Opus / medium / implementer | The fixtures stopped binding `BENCH_HOME`, and 47 stub-swapping tests collided under `-race`. The delegate enumerated the injectables, priced two fixes, and stopped for the reviewer instead of adding locks at 93 call sites. |
| `worktree-test-floor` review | Sonnet / high / reviewer | The Coverage axis read the census's left-hand-side match, found that `cryptorand.Reader = …` through a selector was invisible, and showed the two live callers were serial only by an incidental bind. |
| `ft228-debug-restoration` IP10 | Opus / medium / orchestrator | A coverage row required a description to carry a trigger and omit an anti-trigger. The check graded only the omission, and the delegate's own probe could not see the gap. A coordinator probe written from the row's text set the description to neither, and the gate stayed green. |
| `one-change-one-grade` ticket 01 | Sonnet / low / implementer | The delegate's package tests passed while the live-root routing check would have gone red, because the exhaustive registry runs only under the conformance root env. The delegate named the registry it could not edit; the coordinator's live-root run confirmed the red and added the one row. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26). Guidance prose routes to Opus/high.
- Review axes route to Sonnet (reviewer's direction, 2026-08-19). Twelve landings held
  the citation standard at three axes each with no re-run.
- Opus/low beats Sonnet/medium on exact-spec Go file moves: equal correctness, better
  grouping, half the tokens.
- Route a ticket from the decision table at charge time, not from its story's line. An
  exact spec, a known seam, and a covering gate make the cheap row.
- Opus/high is the routing for process-lifecycle, signal, cleanup-authority, and
  foundational Go-seam work. On `worktree-test-floor` both such charges landed
  first-pass or stopped correctly on a seam gap.
- Independent tickets run in parallel worktrees under disjoint fences; the coordinator
  folds each diff by patch in `Blocked by:` order and probes each fold. At most two
  delegates run tests at once, each at `-parallel 2`, until the test-thread pool lands.
- A performance spec measures its wall after the first slice; a premise about where
  the cost sits is a claim until `go test -json` shows it.
- The coordinator probes every done-claim at a distinct site and kind. At least one
  probe per acceptance row comes from the row's own text.
- A red is attributed or it is not resolved. A one-off red under parallel load is
  reported to review with its reproduction attempts, never called transient.
- A charge names `bench worktree exec` as the only command form, forbids a background
  test run, and names every registry the family already appears in.
- The coordinator takes a commit's path list from the checkout's status, never by hand,
  and never creates or cleans a worktree while a gate runs.
