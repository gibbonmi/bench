# Review: roadmap-light-path-fixes-2

Base `70ccbe8ad978b6e65d5c72df49687067b925a73d`, tip
`85a788c1e8bd4cfdb2536cc548f3e0eec25f4d97`. Initial review, discovery pass.

## Standards

Count: 4. Worst: finding 2, a ticket's `Writes:` line documents a dependency
its `Blocked by:` line denies, on a fixture a later ticket created.

1. **Unordered `Writes:` overlap, sweep vs. craft-spec ticket** — `auto-fix`.
   `spec.md:93` states every `Writes:` overlap is a blocker edge. The
   fixture-closure fix at `85a788c1` widened `sweep-literal-wait-deadlines-in-tests.md`
   with seven paths `state-the-fence-order-and-the-claim-words-in-craft-spec.md`
   already names; neither ticket blocks the other. Both landed clean with no
   real byte conflict. The shared paths are fixture-closure pins on each
   side, not a competing edit to the same bytes. The graph text no longer
   matches the rule as stated.
2. **Overlap that inverts the real landing order** — `auto-fix`.
   `trim-over-length-prose-in-the-live-tree.md` writes
   `.agents/commands/bench-review-implementation.md` and pins
   `review-base-merged-main-tip`, a fixture
   `state-the-census-read-the-changelog-rule-and-the-review-base.md` creates.
   Trim landed at `83758684`; census landed later, at `b2a8a1aa`. The pin was
   backfilled after both landed, so the `Writes:` line names a fixture that
   did not exist when the ticket ran.
3. **`templateFields` derived twice, ungraded** — `ask-user`. The seven
   template field names live in `internal/prose/parse.go` and independently
   in `ste-prose.md:31`. Nothing checks the two against each other. The
   Coverage axis also found this list incomplete (see Coverage finding 2) —
   same repair target.
4. **Anchor diagnostic names a location it doesn't use** — `auto-fix`, low.
   `registry_data.go`'s new rows give `Section: "Exit handoff"` but their
   diagnostic strings say "post-merge tail" and "gate-then-commit path",
   neither of which is a real heading.

## Spec

Count: 3. Worst: finding 1, the spec's own amendment claims the Ownership
fences carry the full `Writes:` union, and they don't.

1. **Ownership fences omit one path** — `auto-fix`. `spec.md:224` claims the
   fence is the union of the tickets' `Writes:` lines. The repair ticket
   names `tests/canary/data-handling-derivation/undocumented-passlist-var`;
   no `data-handling-derivation/` prefix is in the fences.
2. **Stale ticket count** — `auto-fix`. `spec.md:224` still says "eight
   tickets"; nine ticket files exist since the LP5 repair.
3. **Three guidance sentences ship without a dedicated fixture, against the
   Implementation-decision's five-part promise** — non-behavioral
   contradiction, flagged for veto per the tree convention, not blocking.
   The coverage map's own row schema pins one fixture per row and lets a row
   cover two stories; the build matched the map. All three sentences are
   still pinned by their registry test.

## Coverage

Count: 3. Worst: finding 1, the new sweep check cannot bite the migrated-site
class its own self-probe was written to demonstrate on.

1. **The sweep check's spelling set has a blind spot on a site this diff
   itself migrated** — `ask-user`. `checkWaitDeadlineLiterals` recognizes
   only `time.After` and `time.Now().Add`.
   `internal/gate/prospective_owner_test.go:367` (migrated by this diff)
   spells its wait as `time.NewTimer(window)`; reverting that site to a
   literal stays green under the check. `context.WithTimeout`
   (`run_failure_outcomes_test.go:34`) and `time.AfterFunc`
   (`phases_test.go:103`) carry the same blind spot with an existing
   literal.
2. **`templateFields` omits two live field names** — `ask-user`.
   `decisions/*.md` uses `Supports:` and `Drift:` as terminated field lines
   in 67 places across 10 files; neither name is in the closed list.
   Demonstrated: adding two sentences to one `Supports:` line reds the
   6-sentence paragraph bound today. Latent, not yet triggered. Same repair
   target as Standards finding 3.
3. **`internal/models/models_test.go`'s changed wait is a serialization
   detector, not a hang guard** — `ask-user`. The site is fenced (the sweep
   ticket's mid-build widening added it), but it sits outside LP9's original
   fourteen-site enumeration. Deriving its window from 200ms to ~20s slows
   the concurrency-failure path from ~600ms to ~60s if `Inventory` ever
   serializes. That literal exists to catch a real serialization defect
   fast, arguably the sweep's own Won't-handle "subject under test" case.
   The fence-widening note that added this site did not cite it that way.
