# Extract lifecycle policy decisions

Blocked by: 02-add-explicit-effect-inputs.md
Writes: internal/worktree/, internal/worktree/lifecyclepolicy/ (new)

## What to build

Extract ownership, lease, eligibility, age, ignored-output, preservation, and
action decisions into a pure child package. Parent adapters continue to own
native registration, locking, recovery, and removal effects.

Move lifecycle matrices below the typed seam. Retain representative real-Git
journeys and focused fact-adapter coverage.

## Acceptance

- [ ] LC1: Typed lifecycle tables vary every named fact and preserve current verdicts.
- [ ] FA1: Focused real-Git adapter tests prove lifecycle fact translation.
- [ ] RJ1: Native create, remove, registration, lock, and recovery journeys remain serial and green.
