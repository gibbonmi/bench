# Migrate publication outcomes and actions

Blocked by: none
Ownership fence: `internal/publication`
Integration surfaces: shared carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-outcome-action-routes.md
Contracts: operation enum, durable result kind, exact exit, ordered fixed action tokens, and record absence cross `internal/publication` state machine→renderer, membership is prepare/submit/promote/rollback/status, order follows durable record transition, and absence is explicit no-record, asserted by PB1
Closure: PB1/five-operations, PB1/kind, PB1/exit, PB1/action-input, PB1/record-authority, PB1/bypass

## What to build

all five publication operations route durable outcome/action facts without moving record authority.

## Acceptance

- [ ] [PB1] (covers OA6) all five publication operations route durable outcome/action facts without moving record authority.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PB1/five-operations | bypass prepare | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| PB1/kind | replace one durable kind | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| PB1/exit | infer exit from next action | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| PB1/action-input | drop version from a fixed action input | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| PB1/record-authority | derive next action outside the record | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| PB1/bypass | restore local status result rendering | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

