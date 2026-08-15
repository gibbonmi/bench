# Register verdict record classes once

Blocked by: characterize-gate-run-outcomes.md
Writes: internal/gate/verdict.go, internal/gate/verdict_test_records_test.go (deleted), internal/gate (one new registry file and one new guard-test file)

## What to build

Story 3, registry half. One ordered registry with exactly five rows — full,
partial, check-partial, combined-partial, pending — each carrying its name,
exact sorted field set, and validator. `validateRecordBytes` selects the row
whose field set the record carries and runs its validator; the four
`*ReadyFields` variables, the inline Pending field list, and the
`readyFieldClasses` map with its phantom comment references go;
`partialTestRecord`/`fullTestRecord` and their file are deleted.
`narrowVerdictReason` is untouched here (the sibling ticket consumes the
registry). The guard test's expectation is an independently authored literal
(ADR 0006): five names with their field sets. TDD order: stub the registry
empty, write the guard test red on the literal, make it green, then apply the
omission mutation (drop the `check-partial verdict` row) and record the observed
red. `Writes:` overlaps the extract/delete tickets on `verdict.go`, so this
ticket lands serially after them on the retained source even though it does not
depend on their code.

Return note (not acceptance): the omission mutation's observed red and each
row's validator no-op mutation's observed red.

## Acceptance

- [ ] Mixed-class records are `Invalid`; each of the four exact ready sets is `Ready` (covers GR1)
- [ ] Guard literal matches the registry (covers GR3)
- [ ] None of the retired identifiers or phantom strings occurs in `internal/gate` (covers GR4)
- [ ] Executed/skipped overlap, missing check evidence, and `owner_pid: 0` are each `Invalid` (covers GR5)
- [ ] Every characterized outcome is unchanged with test assertions unmodified (renames only, plus the deleted fixture file)
