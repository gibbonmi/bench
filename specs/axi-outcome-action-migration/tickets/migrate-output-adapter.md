# Migrate the returned-output adapter

Blocked by: none
Ownership fence: `cmd/bench/main.go`, `cmd/bench/command_registry.go`, `cmd/bench/main_test.go`
Integration surfaces: shared outcome package→implemented prerequisite `axi-carriers-and-registry`; query migrations→migrate-learnings-outcome.md, migrate-maps-outcome.md, migrate-diff-outcome.md, migrate-coverage-outcome.md, migrate-structure-outcome.md, migrate-models-outcome.md, migrate-testreport-outcome.md, migrate-guards-outcome.md, migrate-outline-outcome.md, migrate-roadmap-outcome.md, migrate-ambient-outcomes-actions.md, migrate-worktree-outcomes-actions.md
Contracts: command member name, domain kind, integer exit, payload bytes, and absence cross cmd-local producer→`cmd/bench/main.go`, membership is every cmd-local returned-output member, ordering is produce-carry-render, and absence preserves empty payload, asserted by AD1
Closure: AD1/members, AD1/kind, AD1/exit, AD1/bytes, AD1/bypass

## What to build

the adapter and every cmd-local returned-output member construct shared outcomes with exact bytes and exits.

## Acceptance

- [ ] [AD1] (covers OA1) the adapter and every cmd-local returned-output member construct shared outcomes with exact bytes and exits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AD1/members | leave one named cmd-local member on the legacy adapter | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AD1/kind | replace the owner kind with a universal success kind | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AD1/exit | infer exit from output emptiness | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AD1/bytes | normalize one returned newline | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| AD1/bypass | emit the pair without constructing the carrier | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

