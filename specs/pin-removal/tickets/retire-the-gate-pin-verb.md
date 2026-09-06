# Retire the gate pin verb

Blocked by: none
Writes: internal/gate/phases.go, internal/gate/gate.go, internal/gate/command_test.go, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, internal/conformance/entry_point_parity_test.go, bin/bench.sh, internal/testreport/check_test.go, internal/testreport/testreport_test.go, tests/canary/package-core-guard/unrouted-subcommand, tests/canary/injected-ports/unregistered-port, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/docs-currency-token-diet/stale-skill-cli-reference, tests/canary/load-validity-metadata/extensionless-gate-ref, tests/canary/package-core-guard/bounds-duplicate-owner, tests/canary/package-core-guard/reintroduced-bare-skip
Covers: PR6, PR7, PR8, PR9, PR10, PR11, PR12, PR13, PR14

## What to build

Verify the premise first: read `PinCommand`, `pinCommand`, `showPinReview`,
`dirtyBench`, `existingPinnedCommit`, `writePinFromHead`, `pinPath`, and
`pinFileName` in internal/gate/phases.go. Read `Command` and `commandUsage`
in internal/gate/gate.go. Read the `gate` and `gate-pin` entries in
cmd/bench/main.go, the help fixture in `TestHelpInventoryIsComplete`, and
`TestRunGateRejectsBriefUsage` in cmd/bench/main_test.go. Read
`TestCommandDispositionsAreComplete` in cmd/bench/command_registry_test.go.
Read `gate_usage` and `gate_command` in bin/bench.sh.

Read the `gate-pin` rows in internal/conformance/subcommand_routing_test.go
and internal/conformance/entry_point_parity_test.go. Read the
`checkSubcommandRouting` clause that reds a registry row the dispatch no
longer names. Read the gate-owned record maps in
internal/testreport/check_test.go and internal/testreport/testreport_test.go.

Delete the pin command and its helpers from the gate package, and update the
file's leading comment. Drop the `pin` case from `Command` and set
`commandUsage` to `usage: bench gate [--fresh]`. Remove unused imports as
the compiler names them. Drop the `gate-pin` registry entry and the pin help
row from cmd/bench/main.go, and drop `gate-pin` from the disposition census.

Drop the `gate-pin` rows from the routing and parity tables. Drop the pin
route and the pin usage from the wrapper. After that edit, `gate_command`
keeps the bare run, `--fresh`, and the help spellings. Every other shape
prints the usage and exits 2. Drop the `.git/bench-gate-pin` row from both
test-report record maps.

Replace the `pin extra argument` row in the gate command test with a `pin`
row that expects exit 2 and `commandUsage` on stderr. Update the usage
string in `TestRunGateRejectsBriefUsage` and the help fixture.

## Acceptance

- [ ] `bench gate pin` exits 2 and prints `usage: bench gate [--fresh]` on stderr.
- [ ] `bench gate --brief` prints `usage: bench gate [--fresh]` on stderr.
- [ ] `bench help` prints `bench gate [--fresh]` and no line that holds `bench gate pin`.
- [ ] The system-attachment census holds no `gate-pin`.
- [ ] The wrapper's `gate pin` exits 2 and prints the usage line.
- [ ] The subcommand-routing check stays green with no `gate-pin` row.
- [ ] The entry-point-parity exemption table holds no `gate-pin`.
- [ ] A search for `bench-gate-pin` and `PinCommand` over the Go sources returns no line.
- [ ] Both test-report record maps hold no `bench-gate-pin` row.
- [ ] Self-probe: keep the `gate-pin` routing row after the registry entry leaves, and report the subcommand-routing check red.
