# Grade the check floor per chunk

Blocked by: 03-record-entries-at-the-own-commit.md
Writes: internal/preflight/check_floor.go, internal/preflight/check_floor_test.go, internal/preflight/check_floor_command_test.go, internal/preflight/plan.go, internal/preflight/decision.go, internal/preflight/preflighttest/chunk_table.go, internal/reviewrecord/assignment.go, internal/reviewrecord/recordcmd/assignment_test.go
Covers: RC25, RC26, RC27, RC28, RC29, RC30

## What to build

Review chunk: RC-C3.

`check-floor` grades the union of a chunk's floor elements against the union of the chunk's requirements. A package that several tickets of one chunk write then needs one floor requirement in that chunk. A requirement of another chunk covers nothing, and a `--changed` run satisfies no element. RC-C2 and RC-C2B land first, because the chunk floor assumes that each entry is graded at its own commit.

Change `checkFloorCheck` (`internal/preflight/check_floor.go:127-151`) to grade each chunk once. Remove `ticketRequirements` (`check_floor.go:106-123`); its one caller is `check_floor.go:136`. `reviewrecord.TicketRequirements` (`internal/reviewrecord/ticket.go:39-58`) stays the one ownership derivation, and the chunk union calls no ownership filter. Each requirement still names one owner ticket, and `requirementsValid` keeps the rule that every ticket owes one requirement (`internal/reviewrecord/plan.go:193-196`).

The red detail names the chunk and the omitted element. Today `omittedElement` builds "<ticket>: <element> (add <command>)" and `floorRedDetail` adds the prefix "completion plan omits a floor check: " (`check_floor.go:153-161`). The detail format is a rendered shape, so move each pin in this ticket: the literal at `check_floor_test.go:205-207`, the helper calls at `check_floor_test.go:84,166,172,178,186,248-302,327,342,355`, and `check_floor_command_test.go:13,38,73`. `gatherCheckFloor` (`plan.go:32`) and `decision.go:220` read the facts; change them only if the fact shape changes. A two-ticket chunk fixture can use the builders in `internal/preflight/preflighttest/chunk_table.go`.

Rewrite `TestCheckFloorReadsOnlyTheOwnTicket` (`check_floor_test.go:330-343`) to RC25, because the chunk union turns its red green. Keep `TestCheckFloorReadsOnlyTheOwnChunk` and `TestCheckFloorRejectsAChangedRun`.

`check_floor_test.go` is at 383 of 400 lines, and the commit lane grades growth. Move an existing test to `check_floor_command_test.go` (106 lines) when the new rows need the room. Do not add a file to the crowded `internal/preflight/` directory.

RC30 already holds: `bench record assignment` refuses while a planned ticket file differs from HEAD (`internal/reviewrecord/assignment.go:81-85,105-117`), and the case "changed planned ticket" of `TestAssignmentSourceRefusals` (`internal/reviewrecord/recordcmd/assignment_test.go:182-187`) pins it. The repair order of ticket 05 depends on this refusal. Keep that case green; this ticket adds no RC30 test.

## Acceptance

- [ ] RC25: `check-floor` is green when one ticket's requirement covers a package that another ticket of the chunk writes.
- [ ] RC26: a `check-floor` red names the chunk and the omitted element.
- [ ] RC27: a requirement of another chunk covers no element of this chunk (`TestCheckFloorReadsOnlyTheOwnChunk`).
- [ ] RC28: a `--changed` run satisfies no floor element (`TestCheckFloorRejectsAChangedRun`).
- [ ] RC29: `bench consumers reviewrecord.TicketRequirements` lists the only ownership derivation, and the preflight package declares no `ticketRequirements`.
- [ ] RC30: `bench record assignment` refuses while a ticket file differs from HEAD.
- [ ] The probe "grade each ticket against its own requirements only" bites in `bench test --package ./internal/preflight`.
