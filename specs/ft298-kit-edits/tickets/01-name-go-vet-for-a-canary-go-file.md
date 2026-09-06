# Name go vet for a canary go file

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md
Covers: none

## What to build

Add one rule to `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` under `In the charge`. The rule reads: A charge that writes a Go file under `tests/canary/` lists `go vet` in its focused checks. The fixture-bite test only parses the overlay, and `go vet` compiles every Go file in the tree. Write it in ASD-STE100, in the shape of its neighbouring items, and add no other sentence. Source: roadmap row FT298.

## Acceptance

- [ ] the rule sits as one list item under `## In the charge`.
- [ ] the item names `tests/canary/` and `go vet`.
- [ ] `bench gate-prose` passes on the edited file.
