# Claude model scorecard

Last incorporated landing: `one-change-one-grade` (`9c1580b9`, 2026-08-25) —
Fable/high as orchestrator and as the guidance implementer, Opus/medium on three
Go charges, Sonnet/low on two ticket charges, and Sonnet/high across four review
passes.

Twenty-five completed landings are recorded. Routing follows the harness-to-tier
binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high, medium | orchestrator, 7 landings + implementer, 9 charges + reviewer, 2 specs | On `one-change-one-grade` it ran three tickets in parallel worktrees, folded each diff by patch, and every one of its seven probes bit at a site distinct from the delegate's. Its guidance charge landed prose-clean first-pass. | Top tier for guidance prose, for coordination of a parallel build, and for adversarial spec review |
| Opus / medium, low | implementer, latest 10 medium charges + 11 low charges | On `one-change-one-grade` two of three medium charges landed first-pass, and the third took one repair round for a raw `os/exec` fixture the census forbids. The lane charge flagged a record-class validator order the ticket did not name. | Medium for gate and conformance logic; low for a ticket from an exact spec at a known seam under a covering gate |
| Opus / high | implementer, latest 10 charges (lifecycle, guidance prose, Go-seam rewrites) | All landed first-pass. On `ft234-pool-key-reclaim` the destructive-apply charge showed its own acceptance row's seam was unreachable instead of writing a test that could not fail. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium | orchestrator, latest 10 of 14 landings | On `learnings-dated-line-visibility` it found untested marker leniency no review axis raised and pinned it. It reproduces every accepted finding and probes every done-claim at a distinct site and kind. | Continue; the distinct-site probe discipline holds |
| Sonnet / low–high | implementer, latest 10 of 57 ticket-sized charges | On `one-change-one-grade` both low charges landed first-pass: the `gate-prose` verb with its grader, and the two review-repair tests. The verb charge flagged the routing registry it could not edit instead of guessing. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor |
| Sonnet / low–high | reviewer, 3 axes on 11 landings + 9 scoped re-reviews | On `one-change-one-grade` (high) the axes returned 10 raw findings, all cited, and the coordinator collapsed them to 2 repair targets. The Coverage axis traced six edges and refuted four by code reading before it reported. | Standing tier for the three review axes; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft228-debug-restoration` IP10 | Opus / medium / orchestrator | A coverage row required a description to carry a trigger and omit an anti-trigger. The check graded only the omission, and the delegate's own probe could not see the gap. A coordinator probe written from the row's text set the description to neither, and the gate stayed green. |
| `module-size-split` batch 1 | Opus / low vs Sonnet / medium / implementer | The same two file-split tickets ran on both lines. Every automated row passed on all four trees, so only the grouping axis separated them. Opus/low placed each overflow symbol with its owner group; Sonnet/medium parked one in the fixtures file, at about twice the tokens. |
| `learnings-dated-line-visibility` round 1 | Fable / high / reviewer | The spec anchored a parser rule on a literal that a scaffold generator also emits, above the first `## ` line. The generator prints that literal below its own `## ` example, so the rule could never open on a fresh journal. Both files had been read; neither was walked against the other. |
| `one-change-one-grade` ticket 01 | Sonnet / low / implementer | The delegate's package tests passed while the live-root routing check would have gone red, because the exhaustive registry runs only under the conformance root env. The delegate named the registry it could not edit; the coordinator's live-root run confirmed the red and added the one row. |
| `one-change-one-grade` ticket 03 | Opus / medium / implementer | The charge asked for two owner constructors; the delegate kept one `New()` with two accept sets so `land.go` outside the fence stayed untouched, and it named the deviation. Its one repair round was a raw `os/exec` fixture where `gittest` and `internal/git` were the seams. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- Review axes route to Sonnet (reviewer's direction, 2026-08-19). Eleven landings held
  the citation standard at three axes each with no re-run.
- Opus/low beats Sonnet/medium on exact-spec Go file moves: equal correctness, better
  grouping, half the tokens. A charge for a file split names every check that reads
  the file by path.
- Route a ticket from the decision table at charge time, not from its story's line. An
  exact spec, a known seam, and a covering gate make the cheap row.
- Opus/high is the routing for process-lifecycle, signal, cleanup-authority, and
  foundational Go-seam work until a medium-effort sample contradicts it.
- Independent tickets run in parallel worktrees off `main` under disjoint fences; the
  coordinator folds each diff by patch in `Blocked by:` order and probes each fold.
- A charge whose test fixture shells out names the census's allowed process seams
  (`gittest`, `internal/git`); a raw `os/exec` fixture costs one repair round.
- The coordinator probes every done-claim at a distinct site and kind. At least one
  probe per acceptance row comes from the row's own text.
- A red is attributed or it is not resolved. A delegate's red-then-green report is an
  attribution claim, not a flake.
- A charge names `bench worktree exec` as the only command form. It names every
  state that reaches a branch, not only the state the acceptance row tests.
- Run the live-root conformance test on a delegate's tree before `bench commit`
  whenever the diff adds a dispatch name, an anchor, or a profile row.
