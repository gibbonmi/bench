# Show one chunk with its required base, findings, and chain breaks

Blocked by: 02-show-the-whole-plan-status.md
Writes: internal/reviewrecord/recordcmd/status.go (new), internal/reviewrecord/recordcmd/status_chunk_test.go (new), internal/reviewrecord/recordtest/status.go (new), CHANGELOG.md
Covers: RP15, RP16, RP17, RP18, RP19, RP20, RP21

## What to build

Review chunk: RP-C1.

Add `--chunk <id>` to `bench record status`. With it, the form prints `chunk[1]`, `tickets[N]`, `verification[N]`, `findings[N]`, and `chain[N]` with the fields that the spec's output table fixes. The `required_base` cell comes from the accessor that ticket 01 exports from `chunkBase`. The `findings` table lists every finding that the RP-C1 collector yields for that chunk, each with the repair command that a record refusal prints for its `Repair`. The `chain` table lists each commit that the collector names on the gap range or the uncovered range, with the reason `unrecorded fold` for a default-branch merge that no fold covers. A proved comment-only Go commit stays off the list.

An unplanned chunk ID refuses at exit 1 and names each planned chunk ID, as the ticket refusal names the chunk tickets (`internal/reviewrecord/ticket.go:46`). Put the new tests in the new `status_chunk_test.go`, so that `status_test.go` stays inside its line budget. Add one `CHANGELOG.md` entry under `## [Unreleased]` for RP-C1.

## Acceptance

- [ ] The chunk view prints the `required_base` that `chunkBase` derives for a chunk whose predecessor has a recorded fold (RP15).
- [ ] The chunk view prints only the tickets, requirements, and covered rows of the named chunk (RP16).
- [ ] With two stale entries in one chunk, `findings` holds both, each with its repair command (RP17).
- [ ] A non-comment commit between the predecessor tip and the chunk base appears in `chain` (RP18).
- [ ] An unrecorded default-branch merge appears with the reason `unrecorded fold` (RP19).
- [ ] A proved comment-only Go commit on the gap range does not appear in `chain` (RP20).
- [ ] An unplanned `--chunk` refuses at exit 1 and names each planned chunk ID (RP21).

## Verification

Run `bench test --package ./internal/reviewrecord/...` whole. Observe each planned row red before its production edit, and record the route of each red. Self-probe: omit the `chain` rows from the chunk view, and report that RP18 and RP19 go red.
