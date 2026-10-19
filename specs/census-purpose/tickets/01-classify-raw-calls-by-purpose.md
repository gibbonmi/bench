# Classify census raw calls by purpose at the phase close

Blocked by: none
Writes: .agents/commands/bench-final-check.md, .agents/skills/bench-craft-delegate/SKILL.md, CHANGELOG.md
Covers: none

## What to build

The phase-close census duty of `/bench-final-check` names a Bench form for each raw-call head, or `none`. A census row holds only the time and the command head, so that duty gives no purpose. The FT366 close named `none` for 16 `rg` calls and proposed no sincere improvement for them.

The close now classifies each raw call by its purpose from the calls that the session saw. Each class routes to one kind of proposal: a guidance change for a missed form, a verb for a Bench-owned gap, no change for a sanctioned read, and no change for an unknown call. Each write delegate returns its raw calls with their purpose, so the close can classify calls that it did not make.

## Acceptance

- [ ] `/bench-final-check` states the four classes and the proposal that each class routes to.
- [ ] `/bench-final-check` forbids a Bench wrapper for a read that the repository rules sanction.
- [ ] The retro guidance records each census proposal as its own `### Bench CLI` item.
- [ ] `craft-delegate` asks each delegate to return its raw calls with their purpose, and it keeps the classification rule in `/bench-final-check` only.
- [ ] `CHANGELOG.md` records the changed close duty.
- [ ] The whole-project gate is green.
