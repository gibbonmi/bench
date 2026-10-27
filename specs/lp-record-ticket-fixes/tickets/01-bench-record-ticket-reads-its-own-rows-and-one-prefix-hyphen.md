# Bench record ticket reads its own rows and writes one prefix hyphen

Blocked by: none
Writes: internal/reviewrecord/plan.go, internal/reviewrecord/ticket.go, internal/reviewrecord/recordcmd/ticket.go, internal/reviewrecord/recordcmd/planned_cell.go, internal/reviewrecord/recordcmd/planned_cell_test.go, internal/reviewrecord/recordcmd/chunk_test.go, .agents/commands/bench-implement-spec.md, CHANGELOG.md
Covers: none

## What to build

Today `bench record ticket` refuses when any coverage row of the chunk holds a planned seam cell at the freeze tip.
In a chunk with more than one ticket, a ticket before the last ticket cannot record until the last ticket cites its rows.
The reviewer decided on 2026-10-25 that the refusal reads only the rows that the `Covers:` line of the recording ticket names (roadmap row FT452).
The FT445 build and the FT450 build each hit this refusal.

The plan parser keeps the rows of each ticket of a chunk beside the chunk rows.
`bench record ticket` reads the planned cells of the rows of its own ticket only.
`bench record chunk` and `bench record completion` keep their scope.

`bench record ticket --id-prefix GR-` writes the evidence ID `GR--tests`, because `EvidenceID` adds a hyphen after the prefix.
The FT314 build hit this doubled hyphen.
`EvidenceID` now removes each trailing hyphen from the prefix before it adds one hyphen.
This ticket trims the prefix and does not refuse it.
No reader of the record reads the raw prefix back: each reader compares the full evidence ID only.
The trim in `EvidenceID` also keeps the fixture writer and the verb on one rule.

The `/bench-implement-spec` command text states that a cited cell lets `bench record ticket` accept the ticket.
A `CHANGELOG.md` entry under `[Unreleased]` records the two changes.

## Acceptance

- [ ] In a chunk with two tickets, the first ticket records when its own row is cited and the row of the second ticket is planned.
- [ ] In the same chunk, the first ticket refuses when its own row is planned, and the refusal names that row only.
- [ ] `bench record ticket --id-prefix GR-` writes the same evidence IDs as `--id-prefix GR`.
- [ ] Each existing test of `./internal/reviewrecord/...` stays green.
- [ ] `CHANGELOG.md` has one entry under `[Unreleased]` that states the two changes.
