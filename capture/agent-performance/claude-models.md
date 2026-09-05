# Claude model scorecard

Last incorporated landing: `ledger-settle-policy` (`3ed9140f`, 2026-09-05).
Fable/low ran as orchestrator across one integration source and ten sibling
worktrees. Opus ran the eight ticket charges, two at low and six at medium.
Opus at medium also ran the two repair charges, the three review axes, and the
scoped re-review. Five ticket charges landed first-pass on behavior, and nine of the
coordinator's eleven probes bit.

Fifty-nine completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 27 landings + implementer, 9 charges + reviewer, 3 specs | On `ledger-settle-policy` at low effort it ran ten tickets in pairs across ten sibling worktrees, caught a census premise the spec review had accepted, and added an engagement rule before the reason ticket started. Two folds ran under a delegate's test load and cost one gate each. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `harness-capability-seam` the record package and the parity check each landed first-pass on behavior; the parity charge found and joined one fixture registry the ticket did not name. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium implementer charges, 89 review axes, 13 of 19 orchestrated landings | On `ledger-settle-policy` eight ticket charges and two repair charges ran; five tickets and both repairs landed first-pass on behavior, every delegate probe bit, and three delegates reported a behavior edge or a test-shape limit instead of hiding it. The misses were a raw `t.Skip`, two narrating comments, and two union journeys deleted against the ticket's own line. The four review passes returned seventeen findings, of which seven became repairs. | Medium for gate and conformance logic, guidance prose, canary fixtures, repair charges, review axes, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–medium | implementer, latest 10 of 77 ticket-sized charges | On `git-admin-readers` seven ticket charges and one review repair ran; three landed first-pass on behavior, and every probe the delegates ran bit. The misses were a fence step outside the ticket, two over-budget files grown, a comment widened to hold a budget, a test that ranged over the production list it graded, and a repair that re-derived a fact beside its owner. | Low for an exact-spec ticket at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing |
| Sonnet / high, xhigh | reviewer, 3 axes on 13 landings + 12 scoped re-reviews + 1 spec round | On `structural-refactor-pass` one xhigh round over the spec and nine tickets resolved all 56 cited test names, verified four decisions against the code, and returned one blocking Coverage finding: a moved scan would drop its active-state filter in silence. | Spec-and-tickets review round when the reviewer names it; the review axes stay with Opus |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `structural-refactor-pass` Coverage axis | Opus / medium / reviewer | The axis planted twenty lines in an over-budget file, ran `bench commit --dry-run` through the built binary, and showed the new lane check pass on the growth. It then rebuilt the lane's detached checkout in a fresh repository and showed the `base..HEAD` query empty, which named the one-line fix. |
| `roadmap-light-path-fixes-3` Coverage axis | Opus / medium / reviewer | The axis wrote four throwaway probes (a symlink loop, an aliased import, a tab in a field name, a `//` inside a string), observed each escape, removed every probe, and cited the spec line each edge sits outside. |
| `agent-push-guard` Spec and Coverage axes | Opus / low / reviewer | The Spec axis traced `git -C /other push origin main` through the scanner and showed the push graded against the wrong repository. The Coverage axis probed `xargs git push`, `@`, and `heads/main` with throwaway tests and quoted each observed allow. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `ledger-settle-policy` Coverage axis | Opus / medium / reviewer | The axis fed the settle policy a base-only stage for a union path and showed the adapter render `union content not text` for content that had no content. It also showed the tolerant read refuse an undecodable entries field, a behavior the spec had recorded with no row. Both became rows and repairs. |

## Current decisions

- Routing changes only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run.
- Every subagent runs Opus at low or medium effort unless the reviewer names another
  tier for the run. On `ledger-settle-policy` Opus served every charge, and five of
  eight tickets landed first-pass, so the Opus default holds.
- A reviewer-named `--reviewer` override sets the model for the review round it names.
- A whole-tree gate runs on the integration source after the last repair commit and
  before `bench worktree land`, because a lane-only repair reached the landing gate red.
- A fold or a whole-tree gate runs only when no delegate test run is live on the
  machine. Two landing tests refused as `infrastructure` under load and proved green in
  isolation.
- A test that ranges over the production list it grades is a coordinator catch. The
  omission probe stays silent, and the test enumerates the list itself.
- A light-path ticket from an exact ticket file runs Opus at low. It runs at medium
  when it adds a conformance check, a canary fixture, or CLI output.
- Opus at medium serves the research censuses, the repair charges, the review axes, and
  the scoped re-reviews. The axes run at low when the reviewer names it for the run.
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
  probe that comes back silently green is a missing row. If the same mutation is
  silent at the base, it is an inherited gap, and the retro records it.
- The coordinator probes a delegate's fixture as well as its production code. A test
  that derives its expectation from the thing it grades stays green through a swap.
- A ticket's fence names the file the ticket's own change will move. It also names
  the fence closure paths the preflight names for a bound package or a pinned file.
- A row that widens a forbidden-import pattern names the enumeration of that path's
  current importers across every graded package. The census premise on
  `ledger-settle-policy` passed a review round and reddened the first fold.
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
  after its ticket folds.
- A ticket that adds a fast-lane check proves it through the real lane over a composed
  tree. A shell stand-in for the check proved the wiring and hid a check that never ran.
- The coordinator records every dogfood run in the spec before the repair-scoped
  re-review starts, because an unrecorded run is a blocking finding.
- A spec that changes a rendered message enumerates the exact-match tests on that
  text. A closure-headroom claim of "no edit" hid one exact match.
- A probe that fails to compile proves nothing. The coordinator replaces it with a
  swap that compiles before it reads a verdict.
