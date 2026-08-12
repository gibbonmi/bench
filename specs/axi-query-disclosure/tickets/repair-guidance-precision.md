# Repair guidance precision

Blocked by: none
Writes: `.agents/skills/bench-craft-cli/SKILL.md`, `internal/conformance/axi_query_registry_test.go`

## What to build

Describe coverage disclosure at the transition that actually emits it: a successful extraction query appends a repair-then-check retry when the mapped rows are repairable, while `--check` refusals retain their existing error contract.

## Acceptance

- [ ] [GP1] (covers QD4) guidance names the successful extraction transition precisely and does not claim that the retry is appended after a refusal; its conformance mutation turns the inaccurate wording red.

