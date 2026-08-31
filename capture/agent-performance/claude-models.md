# Claude model scorecard

Last incorporated landing: `citation-execution-proof` (`5b5daa51`, 2026-08-31) —
Fable/medium as orchestrator over seven Opus/medium ticket charges and one
Opus/medium repair charge. Three Opus/medium review axes and one Opus/medium
scoped re-review ran. Seven of seven ticket charges landed first-pass on
behavior. The review found 13 raw findings that collapsed to eight repair
targets; the re-review verified all eight predicates clean.

Thirty-nine completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 18 landings + implementer, 9 charges + reviewer, 3 specs | On `citation-execution-proof` at medium effort it probed each accepted diff at a distinct site and kind, and every probe bit; it attributed a composed-gate red to the right ticket and repaired it in one line. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 40 review axes, 12 of 18 orchestrated landings | On `citation-execution-proof` seven of seven ticket charges and the repair charge landed first-pass with biting self-probes; one charge flagged its out-of-fence registry edit instead of hiding it. The Spec and Coverage axes found the same subtest-anchor defect without seeing each other, and the Coverage axis proved a path-escape green with a live probe. | Medium for gate and conformance logic, guidance prose, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 1 landing | On `bounded-gate-output` it declared the line once for eleven charges, ran a fresh mutation probe distinct from each delegate's own before every commit, and caught a coverage-map amendment with no owning ticket file before it reached preflight. | New this landing; continue and compare after a second orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. | Standing tier for the three review axes before the Opus rule; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `exec-census` ticket 06 | Opus / high / implementer | The `census=<n>` key broke 35 pinned expectations in 11 files, nine outside the fence. The delegate applied the sweep plan-before-apply, weakened nothing, and reported the fence exception with the one part to revert. |
| `prospective-artifact-recovery` review | Opus / medium / reviewer | The Coverage axis registered a worktree through a symlinked path in a real repository, showed Git records the resolved spelling, and cited the comparison that could never match. The Spec axis traced the one branch that let an authored binary escape the bundle. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `native-runtime-evidence-reduction` review | Opus / medium / reviewer | The Spec and Coverage axes each found, without seeing the other, that the gate graded a workflow matrix consumer and never its producer. One token on the producer step restarted the runners the landing existed to remove, with the gate green. |
| `otel-seam-record` review | Opus / medium / reviewer | The Standards and Spec axes each found, without seeing the other, that the exec verb collapsed every nonzero child exit to 1 under a green gate. The Coverage axis planted the FIFO class the hostile-input checklist names and no test held. |

## Current decisions

- Change routing only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run (reviewer direction, 2026-08-26).
- Every subagent runs Opus at low or medium effort (reviewer direction, 2026-08-26).
  Ten landings now run under this rule. On `citation-execution-proof` seven of seven
  ticket charges landed first-pass. The rule stands.
- Opus at medium serves the review axes and the scoped re-reviews. On
  `citation-execution-proof` two axes found the same worst defect independently, which
  is the isolation working as designed.
- A Coverage finding on a hop the build ran live is refuted by that run before it
  becomes a repair.
- The review pickup file and every file a repair adds join the spec fence before the
  next preflight or landing. On `citation-execution-proof` the canary registry joined
  by amendment before the landing.
- The resolved pickup artifact deletes in the same commit that closes its last finding.
- A repair-scoped re-review runs after every repair charge.
- A prose file runs `bench gate-prose` before `bench commit`.
- The coordinator reads the whole census record before `bench worktree land`, because
  the release deletes it. On `citation-execution-proof` the full record was read.
- The coordinator probes every done-claim at a distinct site and a distinct kind. It
  takes a commit's path list from the checkout's status. It never merges into an
  integration worktree while a delegate edits there.
- Tickets that write one shared file run serially; independent tickets run in
  parallel in separate worktrees and fold through `bench worktree merge`. At most two
  delegates run tests at once, each at `-parallel 2`.
- A test whose fixture needs the system tag set pins `BENCH_KIT` to the fixture root.
  The gate exports an ambient kit, and that ambient kit flips the census under composition.
