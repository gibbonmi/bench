# Add pre-lock premise reads

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md, internal/anchors, internal/conformance, tests/canary/workflow-guidance-anchors
Covers: none

## What to build

Add four premise checks to the existing pre-lock exploration guidance. Keep the guidance in its current owner and do not expand model routing.

## Acceptance

- [ ] The author validates each new cross-package import edge with `go list` before the map locks.
- [ ] The author reads the first operand guard before the map locks.
- [ ] The author enumerates test callers of each changed export before the map locks.
- [ ] The author exercises reachability against an adversarial fixture before the map locks.
