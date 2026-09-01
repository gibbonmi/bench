# Claude model scorecard

Last incorporated landing: the three drain landings of 2026-09-01 — `canary-wrap-hint` (`35c30678`), `gate-prose-file-root` (`87599476`), and
`anchors-test-harness` (`76c0f412`). Fable/low ran as orchestrator over three
light-path ticket charges, two at Opus/low and one at Opus/medium, plus three
Opus/low read delegates. All three write charges landed first-pass on behavior
with a biting self-probe. The coordinator's independent swap probes bit at a
different site in each.

Forty-five completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 20 landings + implementer, 9 charges + reviewer, 3 specs | On `spec-authoring-discipline` at low effort it accepted a green gate under a probe that had changed no bytes, then caught the miss with `cmp` and reran the probe red; it resolved a fence block by reading the pinned bytes and placing the split arm after them. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 53 review axes, 12 of 18 orchestrated landings | On the three drain landings the two low charges landed exact-spec tickets first-pass with a red-to-green log per row. The medium charge collapsed 18 test closures into one harness with the test-name set unchanged and the same needle-deletion red before and after. | Medium for gate and conformance logic, guidance prose, canary fixtures, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 2 landings | On `citation-phase-package-scope` it independently re-ran every delegate's tests, vet, and gofmt before trusting a done-claim; caught its own premature spec-status flip before it compounded; and ran two full three-axis review rounds that found 13 findings, most cited by more than one axis independently, before landing. | Continues to hold at high effort; compare again after a third orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. | Standing tier for the three review axes before the Opus rule; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `spec-authoring-discipline` reader-sweep ticket | Opus / medium / implementer | The charge asked for a byte-exact move of a pinned sentence. The delegate found two conformance tests outside the fence that pinned the bytes, stopped with both citations instead of editing out of fence, and finished first-pass once the coordinator placed the split arm after the pinned bytes. |
| `prospective-artifact-recovery` review | Opus / medium / reviewer | The Coverage axis registered a worktree through a symlinked path in a real repository, showed Git records the resolved spelling, and cited the comparison that could never match. The Spec axis traced the one branch that let an authored binary escape the bundle. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `citation-phase-package-scope` repair-scoped re-review | Opus / medium / reviewer | The repair-scoped re-review found a real bug the first repair round had itself introduced: a cross-architecture phase left `CgoEnabled` at the host's value, a false green a fresh three-axis pass would have missed entirely without a scoped second look. |
| `spec-authoring-discipline` Coverage axis | Opus / medium / reviewer | The axis enumerated all 25 new fixture mutations against their live files with a script instead of a sample, and probed a colon-bearing wrapped Sources line that the new refusal branch does not reach. |

## Current decisions

- Routing changes only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run.
- Every subagent runs Opus at low or medium effort. Forty-five landings now run under
  this rule.
- A light-path ticket from an exact ticket file runs Opus at low effort.
- A test-only refactor over a large file runs Opus at medium effort.
- Opus at medium serves the review axes and the scoped re-reviews.
- `bench worktree land --spec <slug>` goes only on the landing that completes a
  spec's final ticket, never on an earlier landing under the same fence.
- A repair-scoped re-review runs after every repair charge, even a one-line repair.
- The coordinator reads the whole census record before `bench worktree land`,
  because the release deletes it.
- The coordinator probes every done-claim at a distinct site and a distinct kind.
  It confirms with `cmp` that the probe changed bytes before it reads any verdict.
- A live-tree anchor probe runs the gate entry test with `BENCH_CONFORMANCE_ROOT` set
  to the worktree. The conformance package alone reads only each fixture's pinned copy.
- A moved guidance sentence takes an `rg` over `tests/` and `internal/conformance` for
  its literal bytes before the fence locks. A pinning test outside the fence stops the
  delegate. The coordinator then resolves the placement.
- Independent tickets run in parallel in separate worktrees and fold through
  `bench worktree merge`; dependent tickets branch from the integration tip. At most
  two delegates run tests at once, each at `-parallel 2`.
- A ticket worktree stays retained after its merge into the integration source,
  because a merge is not a landing; `bench worktree clean --landed` sweeps them after
  the landing.
- Verdict-row changes are verified through `go run ./cmd/bench`, because the installed
  `dist/bench` renders the old rows until a landing publishes the new binary.
