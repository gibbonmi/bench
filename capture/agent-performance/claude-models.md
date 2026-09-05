# Claude model scorecard

Last incorporated landing: `structural-refactor-pass` (`e86d344d`, 2026-09-05).
Fable/low ran as orchestrator across eleven sibling worktrees. It charged four
Opus censuses, one Opus ticket slicer, nine Opus ticket and split charges with
three continuations, and one Opus repair. Sonnet/xhigh reviewed the
spec and the tickets once, and three Opus review axes and one Opus re-review
ran at medium. Five of eight ticket charges landed first-pass, and three
delegates reported a spec or charge contradiction instead of an edit outside
it. The coordinator ran eleven probes; nine bit, and the two silent ones
became rows.

Fifty-seven completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 25 landings + implementer, 9 charges + reviewer, 3 specs | On `structural-refactor-pass` at low effort it wrote the spec, re-scoped it twice as the reviewer reopened the count decision, and ran nine tickets in waves of two across eleven worktrees. Its first dogfood run of the new growth verb caught the pass's one real growth, and two of its eleven probes came back silently green and became rows. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium implementer charges, 85 review axes, 13 of 19 orchestrated landings | On `structural-refactor-pass` nine ticket charges verified their premise with citations; three reported a spec or charge contradiction and offered the revert, and two named a pin the charge got wrong and followed the tree. Four medium research charges censused 104 files with path pins. The Coverage axis proved the lane's growth check a permanent no-op with a live dry run and an independent repro, and the repair charge fixed it with a detached-checkout test. | Medium for gate and conformance logic, guidance prose, canary fixtures, repair charges, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate, and for the review axes when the reviewer names it |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–high | implementer, latest 10 of 70 ticket-sized charges | On `exec-census` three low charges landed first-pass on behavior; two restated a table or a join that an existing package owned and took one repair round each. | Low for an exact-spec ticket at a known seam; medium or high for a behavior-preserving refactor; charges name the package that owns each shared fact |
| Sonnet / high, xhigh | reviewer, 3 axes on 13 landings + 12 scoped re-reviews + 1 spec round | On `structural-refactor-pass` one xhigh round over the spec and nine tickets resolved all 56 cited test names, verified four decisions against the code, and returned one blocking Coverage finding: a moved scan would drop its active-state filter in silence. | Spec-and-tickets review round when the reviewer names it; the review axes stay with Opus |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `structural-refactor-pass` Coverage axis | Opus / medium / reviewer | The axis planted twenty lines in an over-budget file, ran `bench commit --dry-run` through the built binary, and showed the new lane check pass on the growth. It then rebuilt the lane's detached checkout in a fresh repository and showed the `base..HEAD` query empty, which named the one-line fix. |
| `roadmap-light-path-fixes-3` Coverage axis | Opus / medium / reviewer | The axis wrote four throwaway probes (a symlink loop, an aliased import, a tab in a field name, a `//` inside a string), observed each escape, removed every probe, and cited the spec line each edge sits outside. |
| `agent-push-guard` Spec and Coverage axes | Opus / low / reviewer | The Spec axis traced `git -C /other push origin main` through the scanner and showed the push graded against the wrong repository. The Coverage axis probed `xargs git push`, `@`, and `heads/main` with throwaway tests and quoted each observed allow. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `handoff-sections` repair-scoped re-review | Opus / medium / reviewer | The re-review enumerated both branches of a new write-time State check and showed each unreachable, because the store parses before the mutator and the verb never rewrites State. The check was deleted and two rows were restated as parser rows. |

## Current decisions

- Routing changes only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run.
- Every subagent runs Opus at low or medium effort. Fifty-seven landings now run under
  this rule. A reviewer-named `--reviewer <tier> <effort>` override applies to the
  spec-and-tickets round alone.
- A light-path ticket from an exact ticket file runs Opus at low. It runs at medium
  when it adds a conformance check, a canary fixture, or CLI output.
- Opus at medium serves the research censuses, the repair charges, and the scoped
  re-reviews. The review axes run at medium by default, and at low when the reviewer
  names it for the run.
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
- The coordinator probes a delegate's fixture as well as its production code. A test
  that derives its expectation from the thing it grades stays green through a swap.
- A ticket's fence names the file the ticket's own change will move. It also names
  the fence closure paths the preflight names for a bound package or a pinned file.
- Independent tickets run in parallel in separate worktrees created from the
  integration tip; each merges back through `bench worktree merge`. At most two
  delegates run tests at once.
- A material acceptance shortfall a ticket surfaces mid-build routes to the reviewer.
  Under a `--full` batch approval, the build records it in the spec's decision line
  for veto and proceeds.
- The Coverage review axis runs in its own worktree, because it writes throwaway
  probes; the Standards and Spec axes read the retained source.
- A non-blocking review finding is reported, not repaired.
- The coordinator runs a new verb over the real artifact at the first phase boundary
  after its ticket folds. On `structural-refactor-pass` that run caught the one real
  growth of the pass.
- A ticket that adds a fast-lane check proves it through the real lane over a composed
  tree. A shell stand-in for the check proved the wiring and hid a check that never ran.
- The coordinator records every dogfood run in the spec before the repair-scoped
  re-review starts, because an unrecorded run is a blocking finding.
- A spec that inherits a closed decision about a bench signal quotes that signal's
  current value in its first exchange, before the decision locks.
- A probe that fails to compile proves nothing. The coordinator replaces it with a
  swap that compiles before it reads a verdict.
