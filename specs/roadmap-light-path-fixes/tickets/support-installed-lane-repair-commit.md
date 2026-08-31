# Support an installed-lane repair commit

Blocked by: verify-done-claim-owners.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go
Covers: LF11

## What to build

When an installed lane cannot commit its repair, run the same commit core from
the candidate tree. Grade the composed snapshot and require the sanctioned
rebuild after landing.

## Acceptance

- [ ] The fallback uses the ordinary commit core from the candidate tree.
- [ ] The composed snapshot receives the normal grade.
- [ ] Guidance requires a sanctioned post-landing rebuild.

