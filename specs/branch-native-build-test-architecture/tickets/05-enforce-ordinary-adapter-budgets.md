# Enforce ordinary adapter budgets

Blocked by: 01-expose-branch-native-command-decisions.md
Ownership fence: `internal/gate/`, `internal/git/`
Integration surfaces: repository constructor inventory→`internal/git/`; repository census consumer→06-contract-legacy-fixtures-and-enforce-census.md; controlled process-group constructor inventory→`internal/gate/`; process census consumer→06-contract-legacy-fixtures-and-enforce-census.md
Contracts: the one ordinary repository constructor crosses `internal/git/`→its package-wide adapter test owner with absent meaning zero repositories, asserted by AB1; the one controlled process-group constructor crosses `internal/gate/`→its teardown adapter test owner with absent meaning zero process groups, asserted by AB1
Closure: AB1/git-one-repository, AB1/gate-one-process-group, AB1/git-no-full-gate, AB1/gate-no-full-gate, AB1/git-no-go-suite, AB1/gate-no-go-suite

## What to build

Keep exactly two ordinary side-effect exceptions: one Git representation repository for the complete `internal/git` package run and one controlled descendant process group for `internal/gate`. Both remain narrow adapters and invoke no full gate or Go package suite.

## Acceptance

- [x] [AB1] (covers SY4) `internal/git` owns at most one ordinary repository and `internal/gate` owns at most one ordinary controlled process group; both invoke zero full gates and zero Go package suites.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AB1/git-one-repository | add a second repository constructor | adapter budget test | mutate the package test, run it, expect repository budget exceeded |
| AB1/gate-one-process-group | add a second controlled process-group constructor | adapter budget test | mutate the package test, run it, expect process budget exceeded |
| AB1/git-no-full-gate | call the gate from the Git adapter test | architecture census | add the call, run the census, expect forbidden full-gate diagnostic |
| AB1/gate-no-full-gate | call the full gate from the process adapter test | architecture census | add the call, run the census, expect forbidden full-gate diagnostic |
| AB1/git-no-go-suite | start `go test` from the Git adapter test | architecture census | add the command, run the census, expect nested-Go diagnostic |
| AB1/gate-no-go-suite | start `go test` from the process adapter test | architecture census | add the command, run the census, expect nested-Go diagnostic |
