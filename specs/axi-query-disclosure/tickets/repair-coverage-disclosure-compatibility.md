# Repair coverage disclosure compatibility

Blocked by: repair-action-argument-provenance.md
Writes: `internal/coverage/`, `cmd/bench/command_registry_test.go`, `.agents/skills/bench-craft-cli/SKILL.md`, `internal/conformance/axi_query_registry_test.go`

## What to build

Restore extraction-mode exit compatibility, make a mapped zero-row table an honest terminal empty, and describe the shipped map-row behavior without claiming an unimplemented checked-state classifier.

## Acceptance

- [ ] [RC1] (covers QD6) extraction mode preserves its pre-change exit 0 for readable specs, including malformed mapped input, while appending the state-derived retry action where repair is possible.
- [ ] [RC2] (covers QD1) a canonical mapped table with zero data rows returns its primary zero-row table, `help[0]{cmd,why}:`, empty stderr, and exit 0; the seven-member envelope matrix exercises that exact partition.
- [ ] [RC3] (covers QD4) `craft-cli` and its conformance fixture describe one check action per mapped row, without promising an `unchecked` row filter.
