# Claude model scorecard

Last incorporated landing: `ticket-grammar` (`9f341d53`, 2026-08-31) —
Fable/medium as orchestrator over seven Opus/medium ticket charges, one
Opus/medium repair charge, and one Opus/low light-path migration. Three
Opus/medium review axes and one Opus/medium scoped re-review ran. Six of
seven ticket charges landed first-pass on behavior. The review found 23 raw
findings that collapsed to four repair targets; the re-review found one
repair-induced defect, fixed in one scoped round.

Forty completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 19 landings + implementer, 9 charges + reviewer, 3 specs | On `ticket-grammar` at medium effort it caught a self-probe that mutated no source and redid it, attributed a landing-gate red to a concurrent staged spec outside every fence, and unblocked it through a light-path migration landing. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 43 review axes, 12 of 18 orchestrated landings | On `ticket-grammar` six of seven medium charges landed first-pass with biting self-probes, and three charges reported honest blockers — live closure reds, two unbuildable canary classes, a compiled-binding fixture limit — instead of weakened rules. The Standards and Coverage axes found the same enumeration divergence without seeing each other. The Spec axis refuted one impossibility claim with the mutation-harness mechanics, and the repair delegate then refuted the refutation with a live probe. One low charge landed the migration first-pass. | Medium for gate and conformance logic, guidance prose, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 1 landing | On `bounded-gate-output` it declared the line once for eleven charges, ran a fresh mutation probe distinct from each delegate's own before every commit, and caught a coverage-map amendment with no owning ticket file before it reached preflight. | New this landing; continue and compare after a second orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. | Standing tier for the three review axes before the Opus rule; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ticket-grammar` canary ticket | Opus / medium / implementer | The charge named nine canary classes; two could not red through the harness. The delegate proved both limits — one with a live probe that showed 61 identical diagnostics in the mutated and restored runs — and stopped with the evidence instead of shipping a fixture that could not bite. |
| `prospective-artifact-recovery` review | Opus / medium / reviewer | The Coverage axis registered a worktree through a symlinked path in a real repository, showed Git records the resolved spelling, and cited the comparison that could never match. The Spec axis traced the one branch that let an authored binary escape the bundle. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `native-runtime-evidence-reduction` review | Opus / medium / reviewer | The Spec and Coverage axes each found, without seeing the other, that the gate graded a workflow matrix consumer and never its producer. One token on the producer step restarted the runners the landing existed to remove, with the gate green. |
| `otel-seam-record` review | Opus / medium / reviewer | The Standards and Spec axes each found, without seeing the other, that the exec verb collapsed every nonzero child exit to 1 under a green gate. The Coverage axis planted the FIFO class the hostile-input checklist names and no test held. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26).
- Every subagent runs Opus at low or medium effort (reviewer direction, 2026-08-26).
  Eleven landings now run under this rule. On `ticket-grammar` six of seven ticket
  charges landed first-pass and the seventh paid one composed-tree round. The rule stands.
- Opus at medium serves the review axes and the scoped re-reviews. On `ticket-grammar`
  two axes found the same enumeration divergence independently, which is the isolation
  working as designed.
- A delegate impossibility claim gets one adversarial pass before it closes a spec
  class. On `ticket-grammar` the Spec axis refuted one claim, and the repair delegate
  then disproved the refutation with a live probe; the recorded evidence decided it.
- The review pickup file and every file a repair adds join the spec fence before the
  next preflight or landing.
- The resolved pickup artifact deletes with its closing commit, or it retires with the
  spec when only reviewer-pending items remain and the retro records them.
- A repair-scoped re-review runs after every repair charge, and a repair-induced fix
  gets its own scoped check.
- A prose file runs `bench gate-prose` before `bench commit`.
- The coordinator reads the whole census record before `bench worktree land`, because
  the release deletes it.
- The coordinator probes every done-claim at a distinct site and a distinct kind, and
  confirms a probe mutation changed bytes before it trusts the verdict. On
  `ticket-grammar` one sed probe matched nothing and returned a false green. It takes
  a commit's path list from the checkout's status. It never merges into an integration
  worktree while a delegate edits there.
- The first merge into the integration source pays a composed gate early, because
  focused package checks cannot see cross-package fixture reds. On `ticket-grammar`
  the merge gate caught four systemtest reds that seven green package suites hid.
- Verdict-row changes are verified through `go run ./cmd/bench`, because the installed
  `dist/bench` renders the old rows until a landing publishes the new binary.
- Tickets that write one shared file run serially; independent tickets run in
  parallel in separate worktrees and fold through `bench worktree merge`. At most two
  delegates run tests at once, each at `-parallel 2`.
- A test whose fixture needs the system tag set pins `BENCH_KIT` to the fixture root.
  The gate exports an ambient kit, and that ambient kit flips the census under composition.
