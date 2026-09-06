# Rebuild the integrity fixture family on the split shape

Blocked by: parse-ticket-files-beside-the-inline-shape.md
Writes: tests/canary/decision-map-integrity/, internal/conformance/decision_map_integrity_test.go
Covers: DS15, DS41, DS42

## What to build

The base fixture `tests/canary/decision-map-integrity/decision-map.md` becomes an index plus four ticket files under the fixture's compiled path. Each existing `MUTATE.json` targets the file that now holds the line it mutates, and its `EXPECT` keeps the same diagnostic text. Add one fixture per new diagnostic that the expand parser already reports: `ticket-basename`, `tickets-absent`, `tickets-empty`, `orphan-tickets-folder`, `notes-missing`, `decisions-so-far-missing`, `gist-missing`, `gist-unresolved`, `gist-missing-file`, and `gist-malformed`. The `inline-ticket-heading` fixture waits for the contract ticket.

The fixture inventory in `decision_map_integrity_test.go` names each new fixture. The fixture builder in `TestDecisionMapIntegrityCheckValidatesEveryCandidate` writes split maps by hand, because the template still renders the inline shape until the contract lands.

## Acceptance

- [ ] [F1] Every fixture under `tests/canary/decision-map-integrity/` bites under the fixture-bite test.
- [ ] [F2] The inventory names every fixture directory, and deleting any one reds the inventory test.
- [ ] [F3] `bench test --check decision-map-integrity` passes over the repository tree.
