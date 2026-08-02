# Document recomposition in the lifecycle

Blocked by: nothing
Ownership fence: `.bench/BENCH-reference.md` (the Spec-build lifecycle lookup
section), `.agents/commands/bench-implement-spec.md`, and the conformance anchor
that pins both
Assumptions: the lookup table in `.bench/BENCH-reference.md` is an
operation-to-purpose table and explicitly not a second workflow — it names
`.agents/commands/bench-implement-spec.md` as the canonical phase. The skills
index conformance check already reads this file, so verify at pickup which
existing check owns it before adding an anchor, and extend rather than
duplicate. `bench-implement-spec.md` is the source that `/bench-implement-spec`
loads in Claude Code and `$bench-implement-spec` loads in Codex; one edit serves
every harness.

## What to build

The word "recomposition" appears in no kit document. That absence is what let a
session believe the repair round was walkable, spend two commits discovering it
was not, and finish a ten-ticket build as light-path work with the candidate
applied as a patch.

Document the step the lifecycle has always had but never described: when the
branch tip moves, `promote` is the operation that recomposes the run onto the
new tip, and recomposition discards the review. So the repair round is

    repair → promote → review → assign … integrate → review → promote

Write it for the teammate who just walked in — the current decided state, not
the history of the deadlock that motivated it. Two places carry it: the lookup
table's `promote` row gains recomposition as part of that operation's purpose,
and the implement-spec command gains the round itself where it already
describes repair tickets re-entering assign, checkpoint, and integrate.

**Do not restate the round in both files.** The lookup table is a lookup; the
command is the workflow. One sentence in the table pointing at the operation,
the round itself in the command. A second full copy drifts.

**Pin the sentences.** A documented step that no check protects is a step the
next edit silently deletes. The conformance anchor mechanism is how this repo
holds documented sentences in place; anchor on the load-bearing phrases in both
files so deleting either goes red.

## Acceptance

- [ ] DR1 — the `promote` row of the spec-build lifecycle lookup states that a moved tip recomposes through `promote`.
- [ ] DR2 — `.agents/commands/bench-implement-spec.md` states that recomposition discards the review and gives the repair round in order.
- [ ] DR3 — a conformance check goes red when either documented sentence is deleted.
- [ ] DR4 — the repair round appears in full in exactly one of the two files.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DR1 | document the round in the command only | conformance anchor over `.bench/BENCH-reference.md` | delete the recomposition phrase from the `promote` row, run the conformance check, expect a diagnostic |
| DR2 | state that recomposition preserves the review | conformance anchor over `.agents/commands/bench-implement-spec.md` | delete the round from the command file, run the conformance check, expect a diagnostic |
| DR3 | add the sentences without an anchor | the same conformance check | delete each sentence in turn and assert the check reports each one |
| DR4 | paste the full round into both files | review of the two files | read both files and assert the ordered round appears once |
