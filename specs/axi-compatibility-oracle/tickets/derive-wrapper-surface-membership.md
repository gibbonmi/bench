# Derive the wrapper surface census

Blocked by: derive-root-registry-membership.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`
Integration surfaces: root census and case-index type→derive-root-registry-membership.md; wrapper-arm parser→`internal/axi/compatibility`; real wrapper behavior→`bin/bench.sh` exercised unchanged by every WS1 row through `cmd/bench/axi_compatibility_test.go`; routing observation→`cmd/bench/command_registry.go` (`BENCH_COMMAND_OBSERVE`) exercised unchanged by WS1/porcelain-arm and WS1/binary-arm; nested-grammar consumer→derive-nested-grammar-membership.md
Contracts: the wrapper census — one `{token, routing function}` pair per `case` arm in `bin/bench.sh` — crosses `bin/bench.sh`→`internal/axi/compatibility`; its type is an ordered slice of pairs, membership is every arm of the dispatch `case` plus the `${1-help}` default and the `*)` fall-through, order is the arm order in the file, and an arm the parser cannot resolve to a routing function is an error rather than an omission; asserted by WS1 against the really executed wrapper
Closure: WS1/no-argument, WS1/help-word, WS1/help-long-flag, WS1/help-short-flag, WS1/version-long-flag, WS1/version-short-flag, WS1/wrapper-repair, WS1/default-arm-fallthrough, WS1/gate-arm, WS1/adoption-arm, WS1/porcelain-arm, WS1/binary-arm

## What to build

The wrapper is a public surface the Go registry cannot see. This ticket derives
the wrapper half of the case index by parsing the dispatch `case` in
`bin/bench.sh` into `{token, routing function}` pairs, and `cmd/bench/axi_compatibility_test.go`
executes the real wrapper for each derived token and asserts the observed routing
matches.

The derived census must carry the surfaces the registry does not have: the
no-argument invocation (`case "${1-help}"` makes bare `bench` the help arm), the
help, `--help`, and `-h` arm itself, the `--version` and `-v` aliases that shift into
`route_porcelain version`, wrapper-only `repair`, and the `*)` fall-through that
hands an unrecognized token — including the registry-only `freshness-publish`
name that has no arm of its own — to `route_binary` so the Go dispatcher renders
its own `bench: unknown subcommand` message. It must also carry each arm's
routing function, because `gate_command`, `adoption_route`, `route_porcelain`,
and `route_binary` are four different public process identities.

Each derived wrapper surface lands in the case index under the stable case ID
`wrapper-<token>`, with `--`-prefixed tokens spelled `wrapper-flag-<token>`.

## Acceptance

- [ ] [WS1] (covers CO3) the derived wrapper census contains every dispatch arm of `bin/bench.sh` with its routing function, including the no-argument default, the help and version aliases, wrapper-only `repair`, and the unrecognized-token fall-through, and names any arm the executed wrapper reaches that the census omits.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WS1/no-argument | derive the census from the literal `case` arms only, so the `${1-help}` parameter default contributes no member | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/no_argument -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-no-argument` after running `bash bin/bench.sh` with no argument and observing the help body on stdout with exit 0; the child runs under a 30s `exec.CommandContext` |
| WS1/help-word | drop the `help` token from the parsed help, `--help`, and `-h` arm | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/help -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-help`; the `bash bin/bench.sh help` child runs under a 30s deadline |
| WS1/help-long-flag | drop the `--help` token from the parsed help, `--help`, and `-h` arm | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/help_long -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-flag---help`; the child runs under a 30s deadline |
| WS1/help-short-flag | drop the `-h` token from the parsed help, `--help`, and `-h` arm | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/help_short -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-flag--h`; the child runs under a 30s deadline |
| WS1/version-long-flag | drop the `--version` token from the parsed `--version` and `-v` arm | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/version_long -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-flag---version`, whose executed run emits the same bytes as `bench version` because the arm shifts into `route_porcelain version`; the child runs under a 30s deadline |
| WS1/version-short-flag | drop the `-v` token from the parsed `--version` and `-v` arm | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/version_short -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-flag--v`; the child runs under a 30s deadline |
| WS1/wrapper-repair | drop the `repair` arm, which has no `commandRegistry` entry, from the wrapper census | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/repair -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-repair` and reporting that no root census member covers it, since `repair` is routed to `repair_command` and never reaches the Go registry; the child runs under a 30s deadline |
| WS1/default-arm-fallthrough | drop the `*)` arm so an unrecognized token contributes no member | the wrapper-census test | run `go test ./cmd/bench -run TestDerivedWrapperCensusCoversEveryArm/fallthrough -timeout 300s`; it fails at the missing-member assertion naming case ID `wrapper-unrecognized-token`, whose executed run of `bash bin/bench.sh freshness-publish` reaches `route_binary` and the Go dispatcher's `bench: unknown subcommand` path; the child runs under a 30s deadline |
| WS1/gate-arm | record the `gate` arm's routing function as `route_porcelain` instead of `gate_command` | the wrapper-routing test | run `go test ./cmd/bench -run TestDerivedWrapperCensusRecordsRoutingFunction/gate -timeout 300s`; it fails at the routing-function equality assertion for token `gate`, reporting `route_porcelain` against the arm's `gate_command`; the child runs under a 30s deadline |
| WS1/adoption-arm | record the `doctor`, `setup`, `link`, `init`, and `upgrade` arms as `route_porcelain` instead of `adoption_route` | the wrapper-routing test | run `go test ./cmd/bench -run TestDerivedWrapperCensusRecordsRoutingFunction/adoption -timeout 300s`; it fails at the routing-function equality assertion for token `doctor`, reporting `route_porcelain` against the arm's `adoption_route`; the child runs under a 30s deadline |
| WS1/porcelain-arm | record the `status` arm as `route_binary` instead of `route_porcelain` | the wrapper-routing test | run `go test ./cmd/bench -run TestDerivedWrapperCensusRecordsRoutingFunction/porcelain -timeout 300s`; it fails at the routing-function equality assertion for token `status` after `BENCH_COMMAND_OBSERVE=1 bash bin/bench.sh status` reports `command-registry:status` through the porcelain path; the child runs under a 30s deadline |
| WS1/binary-arm | record the `commands` arm as `route_porcelain` instead of `route_binary` | the wrapper-routing test | run `go test ./cmd/bench -run TestDerivedWrapperCensusRecordsRoutingFunction/binary -timeout 300s`; it fails at the routing-function equality assertion for token `commands`, reporting `route_porcelain` against the arm's `route_binary`; the child runs under a 30s deadline |
