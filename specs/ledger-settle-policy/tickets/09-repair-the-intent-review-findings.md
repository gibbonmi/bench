# Repair the intent review findings

Blocked by: 05-migrate-the-assignment-mutators.md
Writes: internal/intent/assignment.go, internal/intent/purge_test.go, internal/intent/admissionpolicy/admissionpolicy_test.go
Covers: LS40

## What to build

Verify the premise first. Read `readPathTolerant` in
internal/intent/transaction.go and `PurgeAssignments` in
internal/intent/assignment.go. Read the three tests in
internal/intent/purge_test.go. Read the comment above the reauthorization
closure in internal/intent/assignment.go, and the comment above the
"a refused reauthorization changes no field" case in
internal/intent/admissionpolicy/admissionpolicy_test.go.

Write the new test `TestPurgeAssignmentsRefusesAnUndecodableEntriesField` in
internal/intent/purge_test.go. It writes a ledger whose `entries` field is an
object, runs `PurgeAssignments`, and observes an error, unchanged file bytes,
and no lock file. This pins the decided tolerant-read behavior LS40 names.

Cut the comment sentence "This keeps the pre-transaction behavior." from the
reauthorization closure. Keep the sentence that says the compensation covers
the write alone.

Replace the provenance comment above the refused-reauthorization case with the
constraint it pins: a refused compare-and-swap changes no field of the
assignment. Name no package the case came from.

## Acceptance

- [ ] With a ledger whose `entries` field cannot decode, `PurgeAssignments` returns an error, the file bytes are unchanged, and no lock file remains.
- [ ] The reauthorization closure comment names the compensation's scope and narrates no migration.
- [ ] The refused-reauthorization case comment states its constraint and names no origin.
- [ ] The pre-existing `internal/intent` suite passes with its test logic unchanged.
- [ ] Self-probe: make the tolerant read ignore the entries decode error, and report the new test red.
