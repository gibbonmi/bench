# Own resume spec-path derivation

Blocked by: none
Writes: internal/spec, internal/worktree

## What to build

Move slug-or-path normalization behind the existing spec owner and make resume
consume that owner rather than reimplementing the rule.

## Acceptance

- [ ] Resume accepts the same slug and explicit-path forms through the spec owner.
- [ ] No production duplicate of the spec-path derivation remains.
