# Rebuild the integrity fixture family on the split shape

Blocked by: parse-ticket-files-beside-the-inline-shape.md
Writes: tests/canary/decision-map-integrity/, internal/conformance/decision_map_integrity_test.go
Covers: DS15, DS41

## What to build

The base fixture `tests/canary/decision-map-integrity/decision-map.md` becomes an index plus four ticket files under the fixture's compiled path. Each existing `MUTATE.json` targets the file that now holds the line it mutates, and its `EXPECT` keeps the same diagnostic text. The new-diagnostic fixtures already exist from the parse ticket, and the `inline-ticket-heading` fixture waits for the contract ticket.

The fixture builder in `TestDecisionMapIntegrityCheckValidatesEveryCandidate` writes split maps by hand, because the template still renders the inline shape until the contract lands. This ticket runs beside the migration ticket: their `Writes:` do not overlap, and the expand parser reads both shapes.

## Acceptance

- [ ] [F1] Every fixture under `tests/canary/decision-map-integrity/` bites under the fixture-bite test.
- [ ] [F2] Every existing `MUTATE.json` targets a file that exists in the split base, and the inventory test still passes.
- [ ] [F3] `bench test --check decision-map-integrity` passes over the repository tree.
