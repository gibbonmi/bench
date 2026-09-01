# Claude model scorecard

Last incorporated landing: `citation-phase-package-scope` (`75781592`,
2026-09-01) — Sonnet/high as orchestrator over one whole-spec ticket charge,
two Opus/medium finish-and-addition charges, and two Opus/medium repair
charges. Two full three-axis Opus/medium review rounds ran. The initial
review found 13 raw findings that collapsed to 9 repair targets. The
repair-scoped re-review found 6 more, one a real bug the repair itself
introduced. Both rounds landed clean on the second re-check.

Forty-one completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 19 landings + implementer, 9 charges + reviewer, 3 specs | On `ticket-grammar` at medium effort it caught a self-probe that mutated no source and redid it, attributed a landing-gate red to a concurrent staged spec outside every fence, and unblocked it through a light-path migration landing. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 49 review axes, 12 of 18 orchestrated landings | On `citation-phase-package-scope` a whole-spec implementer charge reached 9 of 12 rows first-pass and needed two follow-up charges to close the rest; two full review rounds each found real defects the prior pass missed, including a `bounds` policy redeclare invisible to the fast lane and a cross-architecture cgo bug the repair itself introduced. | Medium for gate and conformance logic, guidance prose, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 2 landings | On `citation-phase-package-scope` it independently re-ran every delegate's tests, vet, and gofmt before trusting a done-claim; caught its own premature spec-status flip before it compounded; and ran two full three-axis review rounds that found 13 findings, most cited by more than one axis independently, before landing. | Continues to hold at high effort; compare again after a third orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. | Standing tier for the three review axes before the Opus rule; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ticket-grammar` canary ticket | Opus / medium / implementer | The charge named nine canary classes; two could not red through the harness. The delegate proved both limits — one with a live probe that showed 61 identical diagnostics in the mutated and restored runs — and stopped with the evidence instead of shipping a fixture that could not bite. |
| `prospective-artifact-recovery` review | Opus / medium / reviewer | The Coverage axis registered a worktree through a symlinked path in a real repository, showed Git records the resolved spelling, and cited the comparison that could never match. The Spec axis traced the one branch that let an authored binary escape the bundle. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `citation-phase-package-scope` repair-scoped re-review | Opus / medium / reviewer | The repair-scoped re-review found a real bug the first repair round had itself introduced: a cross-architecture phase left `CgoEnabled` at the host's value, a false green a fresh three-axis pass would have missed entirely without a scoped second look. |
| `otel-seam-record` review | Opus / medium / reviewer | The Standards and Spec axes each found, without seeing the other, that the exec verb collapsed every nonzero child exit to 1 under a green gate. The Coverage axis planted the FIFO class the hostile-input checklist names and no test held. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26).
- Every subagent runs Opus at low or medium effort (reviewer direction, 2026-08-26).
  Forty-one landings now run under this rule; `citation-phase-package-scope` ran six
  Opus/medium charges across implementation, timeout addition, and two repair rounds.
- Opus at medium serves the review axes and the scoped re-reviews. On
  `citation-phase-package-scope` two full three-axis rounds found 13 findings, several
  cited by more than one axis independently, without one axis seeding another.
- `bench worktree land --spec <slug>` flips a spec's `Status` on any landing under that
  spec's fence, not only the ticket landing. On `citation-phase-package-scope` a
  stand-alone fence-amendment landing carried `--spec` and flipped `Status` to
  `implemented` before the ticket's code had landed. The next `bench preflight build`
  refusal caught it. `--spec` now goes only on the landing that completes a spec's
  final ticket, never on an earlier landing under the same fence.
- A repair round can itself introduce a new defect the original review never saw. On
  `citation-phase-package-scope` the repair-scoped re-review found a cross-architecture
  cgo bug the first repair round's own fix had introduced. The scoped re-review is not
  optional even when the repair looks small.
- The review pickup file and every file a repair adds join the spec fence before the
  next preflight or landing. `citation-phase-package-scope` needed two such amendments
  (a test file, then the `bounds` policy registry), each its own small worktree landed
  without `--spec`.
- The resolved pickup artifact deletes with its closing commit, or it retires with the
  spec when only reviewer-pending items remain and the retro records them.
- A repair-scoped re-review runs after every repair charge, and a repair-induced fix
  gets its own scoped check.
- A prose file runs `bench gate-prose` before `bench commit`; `bench commit` itself
  also runs a prose lane and refuses on a violation.
- The coordinator reads the whole census record before `bench worktree land`, because
  the release deletes it. On `citation-phase-package-scope` the coordinator missed this
  step and could not recover the per-verb breakdown after release. The rule stands,
  and now has a concrete cost attached.
- The coordinator probes every done-claim at a distinct site and a distinct kind, and
  confirms a probe mutation changed bytes before it trusts the verdict. It takes
  a commit's path list from the checkout's status. It never merges into an integration
  worktree while a delegate edits there.
- The first merge into the integration source pays a composed gate early, because
  focused package checks cannot see cross-package fixture reds. On
  `citation-phase-package-scope` `bench worktree merge`'s own prospective gate caught a
  `bounds` policy redeclare that every package-scoped test run had missed.
- Verdict-row changes are verified through `go run ./cmd/bench`, because the installed
  `dist/bench` renders the old rows until a landing publishes the new binary.
- Tickets that write one shared file run serially; independent tickets run in
  parallel in separate worktrees and fold through `bench worktree merge`. At most two
  delegates run tests at once, each at `-parallel 2`.
- A test whose fixture needs the system tag set pins `BENCH_KIT` to the fixture root.
  The gate exports an ambient kit, and that ambient kit flips the census under
  composition. `citation-phase-package-scope`'s repair round found and fixed exactly
  this leak, in an existing fixture the original ticket had not touched.
