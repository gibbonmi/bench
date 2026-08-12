# Repair guidance precision

Blocked by: none
Writes: `.agents/skills/bench-craft-cli/SKILL.md`, `internal/conformance/axi_query_registry_test.go`

## What to build

Describe coverage disclosure at the transition that actually emits it: a successful extraction query appends a repair-then-check retry when the mapped rows are repairable, while `--check` refusals retain their existing error contract.

## Acceptance

- [ ] [GP1] (covers QD4) guidance says that a successful default `coverage` extraction with repairable mapped rows appends the exact repair-then-`bench coverage --check <spec>` retry. It does not claim that `coverage --check` or any refusal appends disclosure; its conformance mutation turns that inaccurate wording red.
