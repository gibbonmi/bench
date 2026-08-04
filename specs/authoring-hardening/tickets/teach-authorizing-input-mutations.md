# Teach input mutations for authorizing values in craft-tickets

Blocked by: teach-size-split-signal.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`
Contracts: none crosses

## What to build

The red-mutations field bullet in `craft-tickets` gains one rule: a ticket
touching a value that authorizes an action — a fingerprint, a digest, a token —
carries at least one mutation row that changes the value's inputs (revert the
constant, drop a hashed field) rather than its code path. Control-flow
mutations grade the plumbing; only an input mutation grades what the value
commits to.

The evidence: every mutation row across eight recovery-discard tickets was
control-flow, and both properties that shipped unasserted — the recovery
fingerprint's domain tag and the reclamation fingerprint's disposition
commitment — were input-shaped, invisible to that entire mutation class until
coordinator probes caught them.

Blocked by the size-signal ticket only because both edit the same skill file;
the coordinator lands them serially.

## Acceptance

- [ ] [IM1] the red-mutations bullet requires at least one input mutation for a ticket touching an authorizing value, naming the revert-the-constant and drop-a-hashed-field shapes.
- [ ] [IM2] the rule sits beside the existing mutate-the-subject rule and the skill's frontmatter and structure checks stay green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| IM1 | revert the rule | the reviewer, at review | remove the input-mutation sentence, re-read the bullet against `specs/authoring-hardening/spec.md` story 3, expect the promised rule to be absent |
| IM2 | move the rule into a new standalone section | the reviewer, at review | relocate it away from the red-mutations bullet, re-read against the spec's lands-inside-owning-sections decision, expect the decision violated |
