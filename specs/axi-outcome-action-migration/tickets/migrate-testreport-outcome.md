# Migrate testreport outcomes

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/testreport`
Integration surfaces: returned-output adapter→migrate-output-adapter.md; shared carrier→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-outcome-action-routes.md
Contracts: domain result kind string, exact exit integer, payload bytes, and empty/refusal disposition cross `internal/testreport`→shared outcome carrier, membership is every testreport command result, ordering is owner-carry-render, and absence retains the current empty/refusal form, asserted by TES1
Closure: TES1/kind, TES1/exit, TES1/bytes, TES1/bypass

## What to build

testreport constructs a shared outcome before its existing renderer with exact bytes and exit.

## Acceptance

- [ ] [TES1] (covers OA3) testreport constructs a shared outcome before its existing renderer with exact bytes and exit.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TES1/kind | replace one owner kind in testreport | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| TES1/exit | infer the testreport exit from payload bytes | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| TES1/bytes | normalize one testreport renderer byte | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| TES1/bypass | restore direct local result rendering in testreport | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

