# Trim over-length prose in the live tree

Blocked by: fix-the-label-line-rule-in-the-prose-grader.md
Writes: .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-review/SKILL.md, DATA_HANDLING.md, decisions/cost-follows-project-size.md, decisions/craft-research.md, decisions/diff-visual.md, decisions/gate-budget.md, decisions/gate-critical-path.md, docs/reporesident-distillation.md, tests/canary/workflow-guidance-anchors/coverage-axis-anchor, tests/canary/workflow-guidance-anchors/review-falsification-accept-routing, tests/canary/workflow-guidance-anchors/review-falsification-dispositions, tests/canary/workflow-guidance-anchors/review-kit-guidance-set, tests/canary/workflow-guidance-anchors/review-persistence-anchor, tests/canary/workflow-guidance-anchors/review-preflight-explicit-base, tests/canary/workflow-guidance-anchors/review-repair-ticket-covers, tests/canary/workflow-guidance-anchors/review-repair-ticket-owner, tests/canary/workflow-guidance-anchors/review-standing-falsification, tests/canary/workflow-guidance-anchors/review-universal-claim-bar, tests/canary/workflow-guidance-anchors/craft-review-coverage-row-projection, tests/canary/workflow-guidance-anchors/review-finding-discipline-pointer, tests/canary/data-handling-derivation/undocumented-passlist-var, tests/canary/workflow-guidance-anchors/distillation-reduced-schema-refactor-shape, tests/canary/workflow-guidance-anchors/review-base-merged-main-tip
Covers: LP5

## Why this ticket exists

The label-line ticket's own run of `TestProseMechanicsHoldsOnTheLiveTree` reds
on ten sites the spec's research did not anticipate. Each is a real
over-length sentence or paragraph the old field-line rule hid, because a
mid-sentence colon happened to land at a line wrap. LP5 requires the live tree
green under the corrected rule; this ticket carries the remaining delta. The
reviewer accepted this as a repair ticket rather than a fence widening on the
label-line ticket, 2026-09-02.

## What to build

Verify the premise first: run `TestProseMechanicsHoldsOnTheLiveTree` on the
integration source after the blocker ticket lands. Confirm the finding set
matches the ten sites below, no more and no fewer. Then, for each site, split
or trim the flagged sentence or paragraph. Make it satisfy the two rules — a
25-word descriptive sentence bound and a six-sentence paragraph bound —
without changing its meaning. Prefer splitting a long sentence into two over
deleting content.

Sites:

- `.agents/commands/bench-review-implementation.md:92` — paragraph of 7 sentences
- `.agents/skills/bench-craft-review/SKILL.md:94` — sentence of 30 words
- `DATA_HANDLING.md:5` — sentence of 26 words
- `decisions/cost-follows-project-size.md:197` — sentence of 26 words
- `decisions/craft-research.md:72` — sentence of 26 words
- `decisions/diff-visual.md:109` — sentence of 28 words
- `decisions/gate-budget.md:246` — sentence of 27 words
- `decisions/gate-budget.md:884` — paragraph of 7 sentences
- `decisions/gate-critical-path.md:18` — sentence of 28 words
- `docs/reporesident-distillation.md:188` — sentence of 31 words

`.agents/skills/bench-craft-review/SKILL.md` carries a 122-line budget row in
`projects/benchkit.md`; keep the edited file inside it. Every other file here
sits outside the reviewed prose-budget universe.

The `tests/canary/...` entries in `Writes:` are fixture-closure pins, not new
fixtures: each names a canary that pins other bytes inside one of the six
touched files above. Trim only the ten sites; leave every pinned byte those
canaries cover untouched, and run their fixture-bite tests unchanged to prove
it.

## Acceptance

- [ ] `TestProseMechanicsHoldsOnTheLiveTree` passes with zero findings.
- [ ] Each of the ten sites above no longer appears in the finding set.
- [ ] No new finding appears at a site not in the list above.
- [ ] The prose budget check passes on the worktree.
- [ ] Self-probe: revert one site's trim, and report the check red with that file and line named.
