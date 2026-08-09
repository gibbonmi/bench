# Contract legacy outcome and action routes

Blocked by: migrate-output-adapter.md, migrate-learnings-outcome.md, migrate-maps-outcome.md, migrate-diff-outcome.md, migrate-coverage-outcome.md, migrate-structure-outcome.md, migrate-models-outcome.md, migrate-testreport-outcome.md, migrate-guards-outcome.md, migrate-outline-outcome.md, migrate-roadmap-outcome.md, migrate-ambient-outcomes-actions.md, migrate-specbuild-outcomes-actions.md, migrate-publication-outcomes-actions.md, migrate-worktree-outcomes-actions.md, migrate-shift-outcome.md
Ownership fence: `internal/conformance`, `internal/axi`, `cmd/bench`, `internal/learnings`, `internal/maps`, `internal/diff`, `internal/coverage`, `internal/structure`, `internal/models`, `internal/testreport`, `internal/guards`, `internal/outline`, `internal/roadmap`, `internal/status`, `internal/handoff`, `internal/dashboard`, `internal/specbuild`, `internal/publication`, `internal/worktree`, `internal/shift`, `projects/benchkit.md`
Integration surfaces: all sibling migration basenames→their fenced production routes; compatibility oracle→implemented prerequisite `axi-compatibility-oracle`
Contracts: production member name, required route enum, observed route identity, and legacy symbol census cross production declarations→`internal/conformance`, membership is every source-census member, order is declaration order, and absence of either observation or contraction refuses, asserted by CT1
Closure: CT1/all-members, CT1/all-routes, CT1/legacy-symbols, CT1/zero-consumer-exports

## What to build

every declared route is observed and every superseded local derivation or zero-consumer export is removed.

## Acceptance

- [ ] [CT1] (covers OA9) every declared route is observed and every superseded local derivation or zero-consumer export is removed.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CT1/all-members | omit one source-census member from conformance | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| CT1/all-routes | restore one named bypass | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| CT1/legacy-symbols | restore one exact moved local carrier symbol | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| CT1/zero-consumer-exports | leave one exported shared symbol with no consumer | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

