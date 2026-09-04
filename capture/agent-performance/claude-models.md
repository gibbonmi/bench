# Claude model scorecard

Last incorporated landing: `handoff-sections` (`dae9a77e`, 2026-09-04).
Fable/low ran as orchestrator over nine Opus ticket charges, one Opus
continuation, and five Opus repair charges across fifteen sibling
worktrees. Three Opus review axes ran at medium effort, and two Opus
repair-scoped re-reviews ran at medium and low. Seven of nine ticket
charges landed first-pass, and two delegates reported a fence gap instead
of an edit outside it. The coordinator ran fifteen probes; thirteen bit,
and two vacuous ones were replaced and bit.

Fifty-six completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 24 landings + implementer, 9 charges + reviewer, 3 specs | On `handoff-sections` at low effort it ran nine tickets in four waves across fifteen worktrees, one review round, and two repair-scoped re-reviews. It found the legacy document refusal by a dogfood run of the new verb, routed the lock residue and the dead write-time check to repair tickets, and recorded five build decisions in the spec for veto. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium implementer charges, 81 review axes, 13 of 19 orchestrated landings | On `handoff-sections` fourteen ticket and repair charges verified their premise with citations; one found the premise wrong, two stopped at a fence gap and reported it, and one caught its own vacuous row when its self-probe stayed green. The three review axes at medium returned fourteen cited findings, and the Coverage axis probed an open fence that swallowed sibling sections. The first repair-scoped re-review found a write-time check the parser made unreachable. | Medium for gate and conformance logic, guidance prose, canary fixtures, repair charges, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate, and for the review axes when the reviewer names it |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high | reviewer, 3 axes on 13 landings + 12 scoped re-reviews | On `exec-census` the three axes returned one cited finding each, collapsed to two repair targets; the Coverage axis probed a two-id text with a throwaway test and found the undecided edge. | Standing tier for the three review axes before the Opus rule; medium or high where the charge names a concern to settle by measurement |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `roadmap-light-path-fixes-3` cache ticket | Opus / medium / implementer | Charged to refuse every hold error, the delegate found the unconditional refusal reds fifteen gate tests whose fixtures declare no HOME, added one predicate at all four sites, and reported the narrowing as a reviewer decision with the fixture cause traced. |
| `roadmap-light-path-fixes-3` Coverage axis | Opus / medium / reviewer | The axis wrote four throwaway probes (a symlink loop, an aliased import, a tab in a field name, a `//` inside a string), observed each escape, removed every probe, and cited the spec line each edge sits outside. |
| `agent-push-guard` Spec and Coverage axes | Opus / low / reviewer | The Spec axis traced `git -C /other push origin main` through the scanner and showed the push graded against the wrong repository. The Coverage axis probed `xargs git push`, `@`, and `heads/main` with throwaway tests and quoted each observed allow. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `handoff-sections` repair-scoped re-review | Opus / medium / reviewer | The re-review enumerated both branches of a new write-time State check and showed each unreachable, because the store parses before the mutator and the verb never rewrites State. The check was deleted and two rows were restated as parser rows. |

## Current decisions

- Routing changes only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run.
- Every subagent runs Opus at low or medium effort. Fifty-six landings now run under
  this rule.
- A light-path ticket from an exact ticket file runs Opus at low. It runs at medium
  when it adds a conformance check, a canary fixture, or CLI output.
- Opus at medium serves the repair charges and the scoped re-reviews. The review axes
  run at medium by default, and at low when the reviewer names it for the run.
- The coordinator writes the repair ticket that cites the amended rows before the
  repair-scoped re-review starts. The review preflight is then green on its first run.
- A kit-guidance diff takes the standing Codex falsification pass at the mid tier
  beside the three axes. The kit-guidance set is `.agents/` and `.bench/BENCH.md`; a
  diff that touches neither takes no such pass.
- `bench worktree land --spec <slug>` goes only on the landing that completes a
  spec's final ticket, never on an earlier landing under the same fence.
- A repair-scoped re-review runs after every repair charge, even a one-line repair.
- The coordinator reads every census record before `bench worktree land`,
  because the release deletes them.
- The coordinator probes every done-claim at a distinct site and a distinct kind. It
  confirms with `cmp` that the probe changed bytes before it reads a verdict. A
  probe that comes back silently green is a missing row, and the delegate adds it.
- The coordinator probes a delegate's fixture as well as its production code. All
  three silently-green probes on `binary-freshness` were tests that derived their
  expectation from the thing they claimed to grade.
- A ticket's fence names the file the ticket's own change will move — a captured
  snapshot, a literal row count, a registry row. Three repair rounds on
  `binary-freshness` went to that omission alone.
- Independent tickets run in parallel in separate worktrees created from the
  integration tip; each merges back through `bench worktree merge`. At most two
  delegates run tests at once.
- A material acceptance shortfall a ticket surfaces mid-build routes to the reviewer.
  Under a `--full` batch approval, the build records it in the spec's decision line
  for veto and proceeds.
- The Coverage review axis runs in its own worktree, because it writes throwaway
  probes; the Standards and Spec axes read the retained source.
- A non-blocking review finding is reported, not repaired.
- A hand-resolved merge conflict produces a single-parent commit, so the sibling
  branch is not recorded as merged. The coordinator runs the whole-project gate by
  hand after such a merge.
- The coordinator runs a new verb over the real artifact at the first phase boundary
  after its ticket folds. The fixture never holds the legacy shape; the real file did.
- A charge that adds a write-time check proves the check reachable from a production
  input before the review. A parser upstream can make the check dead code.
- A probe that fails to compile proves nothing. The coordinator replaces it with a
  swap that compiles before it reads a verdict.
