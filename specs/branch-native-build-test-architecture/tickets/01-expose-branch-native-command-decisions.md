# Expose branch-native command decisions

Blocked by: none
Ownership fence: `cmd/bench/`, `internal/adopt/`, `internal/canary/`, `internal/freshness/`, `internal/gate/`, `internal/preflight/`, `internal/specbuild/`
Integration surfaces: top-level process-attachment disposition registry→`cmd/bench/`; disposition census consumer→06-contract-legacy-fixtures-and-enforce-census.md; immutable decision values→the six named package fences; gate-plan consumer→02-assemble-branch-native-gate-drivers.md; selected command implementation observation→03-own-bounded-system-journeys.md
Contracts: the ordered 28-member command disposition registry crosses `cmd/bench/`→the structural census in 06-contract-legacy-fixtures-and-enforce-census.md, asserted by DC1 against the real dispatch names; immutable accepted values, refusal values, plans, lifecycle events, canary selections, and freshness selections cross their named package APIs→production command adapters, asserted by DD1 against the real producers
Closure: DC1/version, DC1/worktree, DC1/resume-clean, DC1/session-inspect, DC1/shift, DC1/commit, DC1/spec, DC1/gate-go, DC1/guard-git, DC1/check-agent-line, DC1/setup, DC1/link, DC1/init, DC1/doctor, DC1/unlink, DC1/upgrade, DC1/worktree-hook, DC1/gate, DC1/gate-run, DC1/gate-pin, DC1/gate-phases, DC1/freshness-check, DC1/freshness-publish, DC1/canary, DC1/stop-verdict, DC1/release-preflight, DC1/prep-release, DC1/release, DZ1/zero-repositories, DZ1/zero-processes, DZ1/zero-commits, DD1/gate, DD1/adopt, DD1/preflight, DD1/specbuild, DD1/canary, DD1/freshness

## What to build

Make the production command implementations and six named decision packages the ordinary in-process test surface. The disposition registry records exactly the 28 process-attachment decisions without replacing the existing registry that governs every CLI route.

## Acceptance

- [x] [DC1] (covers BN1) all 28 named commands call production command implementations directly in ordinary tests and carry exactly one direct, system, or ship process-attachment disposition.
- [x] [DZ1] (covers BN2) decision and command tests construct zero Git repositories, start zero operating-system processes, and commit zero fixture trees.
- [x] [DD1] (covers BN5) the exact six decision domains consume immutable domain values and express accepted, refused, empty, hostile, error, and rerun outcomes through public production APIs.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DC1/version | remove `version` from direct-only disposition | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/worktree | move `worktree` to system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the wrong-disposition diagnostic |
| DC1/resume-clean | remove `resume-clean` from direct-only disposition | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/session-inspect | remove `session-inspect` from direct-only disposition | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/shift | move `shift` to system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the wrong-disposition diagnostic |
| DC1/commit | remove `commit` from direct-only disposition | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/spec | remove `spec` from direct-only disposition | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/gate-go | move `gate-go` to system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the wrong-disposition diagnostic |
| DC1/guard-git | remove `guard-git` from direct-only disposition | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/check-agent-line | remove `check-agent-line` from direct-only disposition | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/setup | move `setup` to direct-only attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the wrong-disposition diagnostic |
| DC1/link | remove `link` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/init | remove `init` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/doctor | remove `doctor` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/unlink | remove `unlink` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/upgrade | remove `upgrade` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/worktree-hook | remove `worktree-hook` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/gate | move `gate` to direct-only attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the wrong-disposition diagnostic |
| DC1/gate-run | remove `gate-run` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/gate-pin | remove `gate-pin` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/gate-phases | remove `gate-phases` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/freshness-check | move `freshness-check` to direct-only attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the wrong-disposition diagnostic |
| DC1/freshness-publish | remove `freshness-publish` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/canary | remove `canary` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/stop-verdict | remove `stop-verdict` from system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/release-preflight | move `release-preflight` to system attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the wrong-disposition diagnostic |
| DC1/prep-release | remove `prep-release` from ship attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DC1/release | remove `release` from ship attachment | command disposition registry test | mutate the production registry, run its focused package test, expect the missing-member diagnostic |
| DZ1/zero-repositories | add a repository constructor to a direct command test | architecture census test | add the constructor, run the census, expect the package and constructor diagnostic |
| DZ1/zero-processes | add `exec.Command` to a decision test | architecture census test | add the process entry, run the census, expect the package and process diagnostic |
| DZ1/zero-commits | call a fixture commit helper from a command test | architecture census test | add the call, run the census, expect the retired-effect diagnostic |
| DD1/gate | drop refusal from gate scheduling decision | gate decision table test | mutate the decision result, run the direct test, expect the refusal mismatch |
| DD1/adopt | reorder a link plan operation | adopt plan table test | mutate the operation order, run the direct test, expect the literal plan mismatch |
| DD1/preflight | accept malformed evidence | preflight decision table test | mutate the evidence input, run the direct test, expect refusal |
| DD1/specbuild | apply an invalid lifecycle event | lifecycle transition table test | mutate the event, run the direct test, expect unchanged state and refusal |
| DD1/canary | omit one selected owning check | canary selection table test | mutate the selected input, run the direct test, expect missing ownership |
| DD1/freshness | accept a mismatched source digest | freshness decision table test | mutate the digest value, run the direct test, expect refusal |
