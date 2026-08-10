# Repair the nondigest retained-state finding

Blocked by: none
Ownership fence: `internal/specbuild/state.go`, `internal/specbuild/runs_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: NR1/nondigest-name-diagnostic

## What to build

Close the accepted Spec and Coverage findings (P1-nondigest-state-name,
C1-nondigest-state-name) from the Terra/xhigh review of candidate
`8309ddb08a80034867238f47364f52eb547b1dec` (receipt
`b98fb10f30d6e08f4aaa1d4334db6850f8851309a4e7d22ddc651ad852161392`):
`Service.Runs` accepts any `*.json` entry that decodes into an internally
valid record as a healthy retained run, without checking that the entry's
filename is that record's own canonical digest name (`statePath` always
writes a slug's state at `digest(slug)+".json"`). A well-formed record
planted under any other `*.json` filename therefore renders as an ordinary
healthy row instead of its own named diagnostic, violating SB4's "any other
entry... is its own named diagnostic row" requirement. Add a distinct
diagnostic class for this case — a `*.json` entry that decodes to a valid
record but whose filename does not equal `digest(record.Slug)+".json"` is
diagnosed and excluded from the healthy rows, exactly like every other named
hostile class the enumerator already carries.

## Acceptance

- [ ] [NR1] (covers SB4) (P1-nondigest-state-name, C1-nondigest-state-name) a
  well-formed, internally valid retained record stored under a `*.json`
  filename other than its own canonical digest name is reported as its own
  named diagnostic row and never counted among the healthy runs; the same
  record under its correct canonical filename renders healthy, unaffected.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NR1/nondigest-name-diagnostic | drop the filename-identity check (treat any decodable, internally valid `*.json` entry as healthy regardless of name) | focused `Runs()` test | write a valid record's bytes under a non-canonical `*.json` name alongside the same record under its canonical name and require the non-canonical entry to surface its own named diagnostic while the canonical entry stays healthy |
