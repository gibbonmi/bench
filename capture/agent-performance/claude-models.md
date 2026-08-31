# Claude model scorecard

Last incorporated landing: `otel-seam-record` (`2c4d263b`, 2026-08-31) —
Fable/medium as orchestrator over thirteen Opus ticket charges (twelve medium,
one low) and one Opus/medium repair charge. Three Opus/medium review axes and
one Opus/medium scoped re-review ran. Twelve of thirteen ticket charges landed
first-pass on behavior. The review found 12 raw findings that collapsed to nine
repair targets; the re-review verified all nine predicates and found one stale
comment claim.

Thirty-eight completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 17 landings + implementer, 9 charges + reviewer, 3 specs | On `otel-seam-record` at medium effort it probed each accepted diff at a distinct site, and every probe bit. Two ownership-fence gaps still reached the review preflight and the landing before it amended the spec. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium + 10 low implementer charges, 40 review axes, 12 of 18 orchestrated landings | On `otel-seam-record` twelve of thirteen ticket charges and the repair charge landed first-pass with biting self-probes; one charge flagged its out-of-fence edit instead of hiding it. The three medium axes returned 12 cited findings with nine repair targets, and one axis found an exec exit regression the green gate could not see. | Medium for gate and conformance logic, guidance prose, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
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
  Nine landings now run under this rule. On `otel-seam-record` twelve of thirteen
  ticket charges landed first-pass. The rule stands.
- Opus at low and medium serves the review axes and the scoped re-reviews. On
  `otel-seam-record` the axes found the exec exit regression and the FIFO hang class.
- A Coverage finding on a hop the build ran live is refuted by that run before it
  becomes a repair.
- The review pickup file and every file a repair adds join the spec fence before the
  next preflight or landing. On `otel-seam-record` the fence reddened twice.
- The resolved pickup artifact deletes in the same commit that closes its last finding.
- A repair-scoped re-review runs after every repair charge. On `otel-seam-record` it
  found a stale one-source claim the repair itself introduced.
- A prose file runs `bench gate-prose` before `bench commit`.
- The coordinator reads the whole census record before `bench worktree land`, because
  the release deletes it. On `otel-seam-record` only the head lines were retained.
- The coordinator probes every done-claim at a distinct site and kind. It takes a
  commit's path list from the checkout's status. It never cleans a worktree while a
  gate runs.
- Tickets that write one shared file run serially; independent tickets run in
  parallel in separate worktrees and fold through `bench worktree merge`. At most two
  delegates run tests at once, each at `-parallel 2`.
- The landing's `--base` is the assignment's recorded start or a descendant of it on
  `main`. The review preflight runs from inside the integration worktree.
- A spec with an open review pickup is not retired by the build session; the reviewer
  decides the pickup first.
- A cross-family top-tier reviewer (codex Sol/high) stays available for a large new
  package; compare it against Opus axes on the next one.
