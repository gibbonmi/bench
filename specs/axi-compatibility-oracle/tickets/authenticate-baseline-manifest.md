# Authenticate the baseline manifest

Blocked by: none
Ownership fence: `internal/axi/compatibility`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: manifest reader and fixture format→`internal/axi/compatibility`; paired capture→capture-pinned-baseline.md
Contracts: baseline subject string, seal preimage fields, fixture file kind, stable case ID, and stdout/stderr/exit/acceptance presence cross `specs/axi-compatibility-oracle/testdata`→`internal/axi/compatibility`, with fixed field order and refusal on any absence, asserted by BM1 against the real reader
Closure: BM1/subject, BM1/seal-preimage, BM1/regular-file, BM1/non-symlink, BM1/unique-id, BM1/four-observations

## What to build

the manifest accepts only the pinned subject, matching builder-seal preimage, regular non-symlink fixtures, unique IDs, and complete four-observation records.

## Acceptance

- [ ] [BM1] (covers CO1) the manifest accepts only the pinned subject, matching builder-seal preimage, regular non-symlink fixtures, unique IDs, and complete four-observation records.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BM1/subject | replace the pinned baseline subject | manifest reader test | load the manifest and require the subject mismatch refusal |
| BM1/seal-preimage | drop the source-subject field from the builder-seal preimage | seal authority test | build under a named timeout and require the computed seal to differ |
| BM1/regular-file | replace one fixture with a FIFO | hostile manifest test | load under a named timeout and require non-regular refusal before reading |
| BM1/non-symlink | replace one fixture with a symlink | hostile manifest test | load and require symlink refusal before target resolution |
| BM1/unique-id | duplicate one case ID | manifest closure test | load and require the duplicate ID refusal |
| BM1/four-observations | remove the stderr observation field | manifest closure test | load and require the missing-field refusal |

