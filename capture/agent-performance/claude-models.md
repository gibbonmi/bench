# Claude model scorecard

Last incorporated landing: `craft-research-skill` (`28e47254`, 2026-09-06).
Fable/low ran as orchestrator across one integration source and one review
worktree. Opus at high ran the two guidance tickets, and Sonnet at low ran the
review repair as the reviewer named. Opus at low and medium ran the three review
axes and the scoped re-review. Every charge landed first-pass on behavior, and
every delegate probe bit.

Sixty-three completed landings are recorded. Routing follows the
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / low–high | orchestrator, 31 landings + implementer, 9 charges + reviewer, 3 specs | On `craft-research-skill` at low effort it ran three charges and five review passes, probed every return at a distinct site and kind, and collapsed eight findings to five repair targets in one round. It omitted the headroom route from the first charge and started to delete a pickup with an open item, and each slip cost one correction. | Coordination of a parallel build and adversarial spec review; it implements nothing unless the reviewer names it |
| Opus / high | implementer, latest 10 charges (Go-seam rewrites, lifecycle, guidance prose) | On `craft-research-skill` the skill charge wrote a 121-line skill that carried all nine anchored sentences verbatim, passed the prose lane first time, and listed the silences it filled and left open. The callers charge kept every pre-existing needle byte-identical and cited the diff file list for the untouched skills. | High for process-lifecycle, cleanup-authority, destructive-command, anchored guidance prose, and foundational Go-seam rewrites |
| Opus / medium, low | implementer, orchestrator, and reviewer combined; latest 10 medium implementer charges, 100 review axes, 13 of 19 orchestrated landings | On `craft-research-skill` the Coverage axis at medium ran seven probes and found that four assessment anchors also sit in the frontmatter, so a body-only rewrite stays green. The low axes returned three cited findings, two accepted, and the low re-review verified all six predicates and stayed inside its blocking scope. | Medium for gate and conformance logic, guidance prose, canary fixtures, repair charges, the Coverage axis, and orchestration; low for a ticket from an exact spec at a known seam under a covering gate, and for the Standards and Spec axes |
| Sonnet / high | orchestrator, 3 landings | On `roadmap-light-path-fixes-2` it ran ten ticket charges and two review rounds (six axes) across two shared worktrees, caught a read-only delegate leaving the integration worktree dirty before the next commit, and routed two material acceptance shortfalls to the reviewer instead of silently resolving them. | Continues to hold at high effort; compare again after a fourth orchestrated build |
| Sonnet / low–medium | implementer, latest 10 of 78 ticket-sized charges | On `craft-research-skill` one low review repair closed five prose and registry findings first-pass, kept every anchored sentence byte-identical, added a biting fixture, and ran an omission probe that bit. On `git-admin-readers` three of eight charges landed first-pass, and the misses were fence steps, over-budget growth, and a re-derived fact. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing |
| Sonnet / high, xhigh | reviewer, 3 axes on 13 landings + 12 scoped re-reviews + 1 spec round | On `structural-refactor-pass` one xhigh round over the spec and nine tickets resolved all 56 cited test names, verified four decisions against the code, and returned one blocking Coverage finding: a moved scan would drop its active-state filter in silence. | Spec-and-tickets review round when the reviewer names it; the review axes stay with Opus |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `structural-refactor-pass` Coverage axis | Opus / medium / reviewer | The axis planted twenty lines in an over-budget file, ran `bench commit --dry-run` through the built binary, and showed the new lane check pass on the growth. It then rebuilt the lane's detached checkout in a fresh repository and showed the `base..HEAD` query empty, which named the one-line fix. |
| `roadmap-light-path-fixes-3` Coverage axis | Opus / medium / reviewer | The axis wrote four throwaway probes (a symlink loop, an aliased import, a tab in a field name, a `//` inside a string), observed each escape, removed every probe, and cited the spec line each edge sits outside. |
| `bench-probe` Coverage axis | Opus / medium / reviewer | The axis composed two rows the spec kept apart, a control byte in the subject name and a failed restore, ran the built verb over it, and showed the render error at exit 1 with the mutation left on disk and the copy unnamed. The finding became row PB47 and the one repair. |
| `worktree-test-floor` ticket 12 | Opus / high / implementer | The charge asked the verb to set the gate child's environment. The delegate traced the closure, showed the value hashes into the subject identity, and stopped with the seam gap instead of editing three packages outside its fence. |
| `ledger-settle-policy` Coverage axis | Opus / medium / reviewer | The axis fed the settle policy a base-only stage for a union path and showed the adapter render `union content not text` for content that had no content. It also showed the tolerant read refuse an undecodable entries field, a behavior the spec had recorded with no row. Both became rows and repairs. |

## Current decisions

- Routing changes only after two comparable runs or one controlled model comparison.
- The top tier implements nothing, code or guidance prose, unless the reviewer names it
  for the run.
- Every subagent runs Opus at low or medium effort unless the reviewer names another
  tier for the run. On `craft-research-skill` the reviewer named Sonnet at low for the
  review repair, and it landed first-pass.
- Guidance prose runs Opus at high under the leverage override. On `craft-research-skill`
  both high charges kept every anchor and passed the prose lane first time.
- A charge that will touch an over-budget file names its headroom route. The route is
  a fenced file with room, or a sibling file on an existing precedent. The first charge
  here omitted it and paid one continuation.
- A material shortfall that needs a seam outside the spec fence takes a fence extension
  recorded under the spec's Build decisions, not a narrowed row.
- The landing rebuilds the broker; a hand rebuild of the primary binary before the
  landing breaks the manifest digest.
- A reviewer-named `--reviewer` override sets the model for the review round it names.
- A whole-tree gate runs on the integration source after the last repair commit and
  before `bench worktree land`.
- A fold or a whole-tree gate runs only when no delegate test run is live on the
  machine.
- A test that ranges over the production list it grades is a coordinator catch.
- A light-path ticket from an exact ticket file runs Opus at low. It runs at medium
  when it adds a conformance check, a canary fixture, or CLI output.
- Opus at medium serves the research censuses, the repair charges, and the Coverage
  axis. The Standards and Spec axes run at low; on `craft-research-skill` the low axes
  found the two accepted Standards repairs.
- A review finding that restates a decision the spec records under Further notes or
  Build decisions is a no-op. The coordinator cites the decision line.
- A review finding that names a missing capability is refuted or accepted against the
  package's exports, never against one file's grep.
- The coordinator writes the repair ticket that cites the amended rows before the
  repair-scoped re-review starts, and closes its registry fence through the preflight.
- A kit-guidance diff takes the standing Codex falsification pass at the mid tier
  beside the three axes. The kit-guidance set is `.agents/` and `.bench/BENCH.md`.
- `bench worktree land --spec <slug>` goes only on the landing that completes a
  spec's final ticket. Its `--base` is the `main` tip the source folded.
- A repair-scoped re-review runs after every repair charge, even a one-line repair. An
  undecided edge it finds is disposed of in the spec, not charged as a second round.
- The coordinator reads every census record before `bench worktree land`,
  because the release deletes them.
- The coordinator probes every done-claim through `bench probe` at a distinct site and
  a distinct kind from the delegate's probe. A probe the verb reports as `invalid` is
  replaced with a swap that compiles before a verdict is read.
- The coordinator probes a delegate's fixture as well as its production code.
- A ticket that adds a line to an over-budget file moves its headroom in the same
  ticket. The lane grades growth against the current tip.
- Independent tickets run in parallel in separate worktrees created from the
  integration tip; each merges back through `bench worktree merge` by label. At most
  two delegates run tests at once.
- A material acceptance shortfall a ticket surfaces mid-build routes to the reviewer.
  Under a `--full` batch approval, the build amends the row, records the reason in the
  spec's Build decisions, and proceeds.
- The Coverage review axis runs in its own worktree when it writes throwaway probes;
  the Standards and Spec axes read the retained source.
- A non-blocking review finding is reported in the pickup, not repaired. A pickup that
  holds an open reviewer item survives the repair commit with that item alone.
- A fixture-pinned path in a ticket Writes line names its fixture. A repair that adds
  a fixture updates every ticket that names the pinned path.
- A spec that changes a rendered message enumerates the exact-match tests on that
  text.
