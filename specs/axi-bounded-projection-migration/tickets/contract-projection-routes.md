# Contract legacy projection routes

Blocked by: migrate-sanitize-projection.md, migrate-roadmap-projection.md, migrate-worktree-display-projection.md, migrate-outline-projection.md
Ownership fence: `internal/conformance`, `internal/axi`, `internal/sanitize`, `internal/roadmap`, `internal/worktree`, `internal/outline`, `projects/benchkit.md`
Integration surfaces: four migration basenames→their production owners; compatibility oracle→implemented prerequisite `axi-compatibility-oracle`
Contracts: owner identity, route identity, policy identity, and legacy symbol census cross production owners→`internal/conformance`, membership is the four helper-census owners, order is census order, and missing route or residual derivation refuses, asserted by PJ1
Closure: PJ1/four-owners, PJ1/routes, PJ1/policies, PJ1/residue

## What to build

all four owners reach their route and superseded projection derivations are absent.

## Acceptance

- [ ] [PJ1] (covers BP5) all four owners reach their route and superseded projection derivations are absent.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PJ1/four-owners | omit one helper-census owner | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| PJ1/routes | restore one owner bypass | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| PJ1/policies | bind one owner to another policy identity | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| PJ1/residue | restore one exact legacy selector/count symbol | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |

