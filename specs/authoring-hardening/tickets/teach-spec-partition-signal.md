# Teach the disjoint-partition split signal in craft-spec

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-spec/SKILL.md`
Contracts: none crosses

## What to build

The story-sizing section in `craft-spec` gains a partition check: when a spec's
stories partition into disjoint package sets connected by no shared seam or
contract, surface the split to the reviewer at spec time. A deliberate bundle
can still be chosen — but chosen, not defaulted, and recorded in the spec.

The evidence: the recovery-discard spec bundled recovery discard
(`internal/worktree`) and spec-build reclamation (`internal/specbuild`) —
disjoint packages, two disjoint ticket blocker chains, a shared fixture package
that even split by file — and the narrower capability could not ship on its own
gate. The authoring-hardening spec itself records its reviewer-chosen bundle
under this rule, and the paragraph may cite that shape without naming file
paths that rot.

## Acceptance

- [ ] [PS1] the story-sizing section surfaces a disjoint story partition as a split signal for the reviewer at spec time, with the bundle-must-be-chosen posture.
- [ ] [PS2] the paragraph lands inside the existing sizing section and the skill's frontmatter and structure checks stay green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PS1 | revert the partition check | the reviewer, at review | remove the paragraph, re-read the sizing section against `specs/authoring-hardening/spec.md` story 5, expect the promised signal to be absent |
| PS2 | rewrite the signal as an automatic split | the reviewer, at review | change the wording to mandate splitting, re-read against the chosen-not-defaulted decision, expect the reviewer's bundling authority to be removed |
