# Contract legacy fixtures and enforce the census

Blocked by: 01-expose-branch-native-command-decisions.md, 02-assemble-branch-native-gate-drivers.md, 03-own-bounded-system-journeys.md, 04-move-canary-and-stripped-proofs.md, 05-enforce-ordinary-adapter-budgets.md
Ownership fence: `internal/`, `projects/benchkit.md`, `CHANGELOG.md`, `specs/branch-native-build-test-architecture/spec.md`
Integration surfaces: constructor and argv exception inventories→`internal/conformance/`; retired fixture symbols and effects→every package under `internal/`; current gate description→`projects/benchkit.md`; user-visible replacement entry→`CHANGELOG.md`; final state→`specs/branch-native-build-test-architecture/spec.md`
Contracts: the exact repository, process, build, ordinary-driver, system-driver, race-driver, stripped-journey, and retired-effect inventories cross `internal/`→`internal/conformance/` in stable path order, asserted by AR1 against the complete syntax census; retired fixture effects cross their former packages→the deletion test with absence meaning successfully contracted, asserted by DL1
Closure: AR1/one-build, AR1/one-ordinary-driver, AR1/zero-decision-repositories, AR1/zero-decision-processes, AR1/git-repository-exception, AR1/gate-process-exception, AR1/zero-nested-go, AR1/zero-inner-gates, AR1/three-system-repositories, AR1/four-dev-repositories, AR1/one-stripped-journey, DL1/fixture-bench, DL1/fixture-bench-wrapper, DL1/commit-all, DL1/copied-kit, DL1/per-test-repository, DL1/duplicated-phase

## What to build

Delete the general fixture framework and process-backed ordinary suites after their behaviors have explicit direct, adapter, system, ship-tier, or intentional-deletion dispositions. Replace the permissive build census with the exact branch-native architecture census, update current-state gate documentation and changelog, and mark the spec implemented only in the committed green composition.

## Acceptance

- [x] [AR1] (covers RG1) the default conformance syntax and argv census enforces every exact build, driver, repository, process, nested-Go, canary, system, dev-gate, and stripped-journey budget from the spec.
- [x] [DL1] (covers RG2) general `Fixture.Bench*`, `CommitAll`, copied-kit constructors, per-test repository constructors, duplicated phase constructors, and renamed wrappers around those effects are absent.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AR1/one-build | add a second selected Bench build | architecture census mutation test | add the constructor, call the owning check, expect exact build-count diagnostic |
| AR1/one-ordinary-driver | add a second ordinary Go driver | architecture census mutation test | add the argv producer, call the owning check, expect exact ordinary-driver diagnostic |
| AR1/zero-decision-repositories | add a repository constructor to a decision test | architecture census mutation test | add the constructor, call the owning check, expect package diagnostic |
| AR1/zero-decision-processes | add a process constructor to a command test | architecture census mutation test | add the constructor, call the owning check, expect package diagnostic |
| AR1/git-repository-exception | add a second `internal/git` repository constructor | architecture census mutation test | add the constructor, call the owning check, expect exact exception diagnostic |
| AR1/gate-process-exception | add a second `internal/gate` process-group constructor | architecture census mutation test | add the constructor, call the owning check, expect exact exception diagnostic |
| AR1/zero-nested-go | add a nested Go test or run command | architecture census mutation test | add the command, call the owning check, expect nested-Go diagnostic |
| AR1/zero-inner-gates | add an inner canary gate | architecture census mutation test | add the call, call the owning check, expect inner-gate diagnostic |
| AR1/three-system-repositories | add a fourth system repository site | architecture census mutation test | add the constructor, call the owning check, expect system budget diagnostic |
| AR1/four-dev-repositories | add a fifth dev-gate repository site | architecture census mutation test | add the constructor, call the owning check, expect total budget diagnostic |
| AR1/one-stripped-journey | add a second stripped journey marker | architecture census mutation test | add the marker, call the owning check, expect stripped-count diagnostic |
| DL1/fixture-bench | reintroduce a process-backed `Fixture.Bench` helper | deletion test | add the helper, run the focused check, expect retired fixture API diagnostic |
| DL1/fixture-bench-wrapper | reintroduce a process-backed wrapper helper under a new name | syntax census | add the effect, run the focused check, expect ordinary process diagnostic |
| DL1/commit-all | reintroduce `CommitAll` | deletion test | add the symbol, run the focused check, expect retired fixture API diagnostic |
| DL1/copied-kit | reintroduce a copied-kit test constructor | syntax census | add the effect, run the focused check, expect copied-kit diagnostic |
| DL1/per-test-repository | add a repository constructor outside the owners | syntax census | add the effect, run the focused check, expect package diagnostic |
| DL1/duplicated-phase | add a second contract or conformance phase constructor | phase argv and deletion test | add the phase, run focused checks, expect duplicated phase diagnostic |
