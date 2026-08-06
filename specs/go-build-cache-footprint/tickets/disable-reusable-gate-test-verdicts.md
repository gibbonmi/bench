# Disable reusable gate test verdicts

Blocked by: none
Ownership fence: `internal/gate/gate_go.go`, `internal/gate/gate_go_test.go`
Integration surfaces: core package derivation and executed argv→`internal/gate/gate_go.go`; filtered conformance registry pattern and executed argv→`internal/gate/gate_go.go`; fake-Go argv, environment, cleanup, sentinel, retained-evidence, package-set, output, and red-exit controls→`internal/gate/gate_go_test.go`; conformance package and skip membership→existing `internal/conformance/registry` + GV2 and GV4
Contracts: the ordered Go argv slice and inherited `GOCACHE` path cross `internal/gate/gate_go.go`→the spawned Go process; GV1-GV3 assert argv begins `go test -count=1`, retains the existing derived package or registry-owned skip suffix order, treats absent conformance as no invocation, and passes an absent or present caller-selected cache environment through unchanged

## What to build

Whenever component selection executes the core Go tests or filtered conformance suite, the owning argv constructor disables successful test-result reuse with `-count=1`. The owners preserve package derivation, the registry skip pattern, caller-selected compilation cache, streaming, red exits, and retained component evidence; they do not clean or replace the cache.

## Acceptance

- [ ] [GV1] Selected core tests execute as `go test -count=1 <enumerated packages>` at dev and ship tiers without changing package derivation or release-only exclusions.
- [ ] [GV2] Selected filtered conformance executes as `go test -count=1 ./internal/conformance -skip <registry pattern>` while an absent declared suite still produces no invocation.
- [ ] [GV3] Both test owners inherit the caller-selected `GOCACHE` unchanged, invoke no cache cleanup, and leave a fake-Go-created `# test log` sentinel in that cache after returning.
- [ ] [GV4] Valid retained evidence may still skip unchanged components, and existing package, filter, output-streaming, and red-exit behavior remains green.
- [ ] [GV5] The representative implemented argv produces zero reusable `# test log` records in a disposable cache while a warmed `-x -run '^$'` probe executes zero compiler commands.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GV1 | omit `-count=1` from `coreTestStep` | the core-owner argv contract | run the focused core contract and expect its owner-specific missing-argument failure |
| GV2 | omit `-count=1` from `ConformanceSuiteArgv` | the conformance-owner argv contract | run the focused conformance contract and expect its owner-specific missing-argument failure |
| GV3 | set `GOCACHE=off`, invoke `go clean`, or delete the planted log after return | the fake-Go environment and filesystem recorder | run each owner against a disposable cache and expect the exact inherited path, no cleanup call, and surviving sentinel |
| GV4 | drop a package, change the registry skip pattern, swallow output, soften a red exit, or bypass retained evidence | the existing gate-go and component-decision controls | run the unchanged focused suites and expect the corresponding package, marker, stream, exit, or skip regression |
| GV5 | disable or discard the compilation cache together with test-result reuse | the disposable-cache comparison | run the implemented workload twice, inspect test-log records and `-x` compiler commands, and expect zero of both |
