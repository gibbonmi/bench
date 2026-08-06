# Anchor the covers annotation to the row-ID bracket

Blocked by: single-source-covers-grammar.md
Ownership fence: `internal/specbuild`
Integration surfaces: taught-example grammar→existing internal/conformance/example_agreement_test.go path, unchanged (its annotations are bracket-adjacent); assign refusal path→existing assign_covers_test.go rows exercised unchanged
Contracts: the per-row covers mapping's position rule (bracket-adjacent only) inside `internal/specbuild`, asserted by AN1 and AN2 against real ticket files

## What to build

The review demonstrated that a `(covers <ID>)` appearing anywhere in a row's
trailing prose — for example inside backticks documenting the grammar — parses
as a real annotation and can falsely satisfy promote's totality. The reviewer
decided the spec's wording ("after its row-ID bracket") means bracket-adjacent:
an annotation counts only when it immediately follows the row-ID bracket with
only whitespace between; anything later in the row is prose and parses as
unannotated. The multiples rule keeps its intent at the anchor: a second
`(covers ...)` chained immediately after the first is malformed and the row
parses as unannotated, while a distant prose mention is simply ignored — the
existing multiple-annotation test's input and expectation move to this rule if
its occurrences were prose-separated. Parse stays permissive; assign still owns
every refusal.

## Acceptance

- [ ] [AN1] `(covers <ID>)` and `(covers local)` annotate only when immediately after the row-ID bracket; a mention later in the row's prose parses as unannotated
- [ ] [AN2] a promote composition whose only claimant for a map ID is a prose mention refuses naming that ID
- [ ] [AN3] bracket-adjacent annotations, including the taught example's forms, parse exactly as before, and a chained second annotation still parses as unannotated

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AN1 | scan the row's whole trailing text instead of anchoring at the bracket | the prose-mention parse test | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the prose-annotated failure |
| AN2 | allow arbitrary text between the bracket and the annotation | the prose-mention promote fixture | apply, run `go test ./internal/specbuild -run TestPromote`, expect the false-coverage failure |
| AN3 | require the annotation at the start of the whole line | the bracket-adjacent capture tests | apply, run `go test ./internal/specbuild -run TestParseTicket`, expect the valid-form failure |
