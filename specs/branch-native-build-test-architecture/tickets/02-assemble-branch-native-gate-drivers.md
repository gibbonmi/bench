# Assemble branch-native gate drivers

Blocked by: 01-expose-branch-native-command-decisions.md
Ownership fence: `.bench/gate.sh`, `.bench/gate-prospective.sh`, `bin/bench.sh`, `cmd/bench/`, `internal/gate/`, `internal/racetests/`, `scripts/go-build.sh`, `scripts/go-build.inputs`
Integration surfaces: ordinary, system, and race argv producers→`internal/gate/`; argv census consumer→06-contract-legacy-fixtures-and-enforce-census.md; selected executable owner→03-own-bounded-system-journeys.md; ship command disposition→01-expose-branch-native-command-decisions.md
Contracts: the stable phase plan with one ordinary driver, one tagged system driver, and one race driver crosses `internal/gate/`→the source-first wrappers and selected executable owner, asserted by GD1, GS1, GR1, and GT1 against assembled argv; the ordered race selector registry crosses `internal/racetests/`→the race phase, asserted by GR1 against the real producer
Closure: GD1/one-driver, GD1/exact-argv, GD1/no-package-loop, GD1/no-nested-go, GS1/one-system-driver, GS1/exact-system-argv, GS1/selected-binary, GR1/one-race-driver, GR1/registry-selectors, GR1/no-system, GT1/release-preflight, GT1/prep-release, GT1/release

## What to build

Make the Go gate plan authoritative. The source-first shell entry remains a bounded bootstrap, while one ordinary package-universe driver, one tagged system driver, and one targeted registry-derived race driver replace duplicated contract, conformance, and stripped phase schedules.

## Acceptance

- [x] [GD1] (covers BN3) the dev gate assembles exactly one `go test -count=1 ./...` ordinary driver and names no separate contract or conformance package phase.
- [x] [GD2] (covers BN4) Go owns ordinary package scheduling and neither the gate nor ordinary tests contain a package loop or nested `go test` or `go run` driver.
- [x] [GS1] (covers SY5) the dev gate assembles exactly one `go test -count=1 -tags=system ./internal/systemtest` driver with the inherited selected executable.
- [x] [GR1] (covers RC1) the gate assembles one `go test -race -count=1` invocation from the authoritative race registry and excludes the system package and journeys.
- [x] [GT1] (covers RG3) release, cross-target, reproducibility, and publication execution remain ship-tier only while all three ship command decisions retain direct coverage.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GD1/one-driver | append a second ordinary package phase | phase argv test | mutate the plan, run the focused gate test, expect duplicate ordinary-driver diagnostic |
| GD1/exact-argv | remove `-count=1` or replace `./...` | phase argv test | mutate the argv, run the focused gate test, expect exact argv mismatch |
| GD1/no-package-loop | enumerate packages into separate commands | architecture census | add the loop, run the census, expect forbidden package-loop diagnostic |
| GD1/no-nested-go | add a nested Go command to an ordinary test | architecture census | add the constructor, run the census, expect forbidden nested-Go diagnostic |
| GS1/one-system-driver | append a second system phase | phase argv test | mutate the plan, run the focused gate test, expect duplicate system-driver diagnostic |
| GS1/exact-system-argv | remove the `system` tag or change the package | phase argv test | mutate the argv, run the focused gate test, expect exact system argv mismatch |
| GS1/selected-binary | omit `BENCH_RUN_BINARY` from system environment | system phase test | mutate the phase environment, run the focused test, expect missing selected-binary diagnostic |
| GR1/one-race-driver | append a second race phase | race argv test | mutate the plan, run the focused gate test, expect duplicate race-driver diagnostic |
| GR1/registry-selectors | drop one authoritative selector | race execution sentinel | mutate the registry member transport, run the focused race test, expect missing sentinel |
| GR1/no-system | add `internal/systemtest` to race argv | race argv test | mutate the registry, run the focused gate test, expect forbidden system package diagnostic |
| GT1/release-preflight | add release-preflight workflow execution to dev phases | tier phase test | mutate the plan, run the focused test, expect ship-tier leakage diagnostic |
| GT1/prep-release | add prep-release workflow execution to dev phases | tier phase test | mutate the plan, run the focused test, expect ship-tier leakage diagnostic |
| GT1/release | add release workflow execution to dev phases | tier phase test | mutate the plan, run the focused test, expect ship-tier leakage diagnostic |
