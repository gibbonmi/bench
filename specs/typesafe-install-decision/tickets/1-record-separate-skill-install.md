# Record the separate install of the TypeSafe skill

Blocked by: none
Writes: .gitignore
Covers: none

## What to build

The ignore file states the decided state for the `typesafe-ai` skill. Each user installs the skill with `npx skills`, and the kit does not bundle or maintain it.

## Acceptance

- [ ] The `.gitignore` comment above the skill paths states that each user installs the skill and that the kit does not bundle it.
- [ ] Git still ignores the skill paths and `skills-lock.json`.
