# Delegate every light-path fix in a drain

Blocked by: none
Writes: .bench/BENCH.md, .agents/commands/bench-drain.md, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/conformance/commitment_guidance_test.go, CHANGELOG.md
Covers: none

## What to build

The **Delegate each light-path fix** paragraph in `.bench/BENCH.md` is the one owner of the drain's light-path route. Today it routes only a light-path fix that a drain verdict keeps. A light-path fix from another source in the run has no stated route. Examples of such a source are a reconcile finding, a worktree-closure finding, a retro recommendation, and a delegate's report.

The paragraph states that `/bench-drain` implements no light-path fix inline. The drain dispatches every light-path fix in its run, whatever its source and whatever the active commitment is. The dispatch line, the fresh write delegate, and the exception for a learning fix that the active outcome needs stay as they are.

`.agents/commands/bench-drain.md` refers to that owner for a light-path change from any source and does not restate the rule. The drain text also agrees with the owner and with `craft-delegate` on the landing. The delegate commits the ticket. The coordinator verifies the done-claim and then lands the ticket.

The commitment anchor that pins the old drain dispatch sentence moves to the new sentence, and its conformance case moves with it. The anchor that pins the old delegate-lands sentence moves to the new delegate sentence.

## Acceptance

- [ ] `.bench/BENCH.md` states that `/bench-drain` implements no light-path fix inline and dispatches every light-path fix in its run, whatever its source.
- [ ] `.agents/commands/bench-drain.md` routes a light-path change from any source in the run to the light-path fix rule in `.bench/BENCH.md`, with no second copy of the rule.
- [ ] `.agents/commands/bench-drain.md` states that the delegate commits the ticket and the coordinator lands it.
- [ ] The commitment guidance check goes red when the new drain dispatch sentence is removed from `.bench/BENCH.md`.
- [ ] `CHANGELOG.md` records the change.
