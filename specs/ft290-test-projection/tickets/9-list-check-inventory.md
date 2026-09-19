# 9. List the check inventory

Blocked by: 8-list-check-fixtures.md
Writes: internal/testreport/, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: TP31, TP32, TP34, TP35, TP36, TP37, TP38, TP39, TP50

## What to build

Chunk: TP-C3.

Add `bench test --checks`. It prints `checks[N]{name,kind,families}` at exit 0, in the order of the help check list.
The producer calls `namedChecks()`, the one name source of the help text.
The `kind` cell reads the kind owner of ticket 3. The `families` cell counts the distinct family names in which the check owns one fixture or more, through the owner filter of ticket 8.

A fixture with no family adds nothing to the count. An absent `tests/canary` directory gives `families` 0 on each row.
`--checks` accepts no other flag and no operand, selects no run binary, and starts no Go child.

Complete the grammar text and the `bench help` row with the form `bench test --checks`, so that they equal the text in the spec.

## Acceptance

- [ ] The `name` cells equal the help check list in order.
- [ ] The rows print `conformance`, `system`, and `prose` for the three kinds.
- [ ] A check with fixtures in two families prints `families` 2.
- [ ] A check that owns only one `CHECK`-marked fixture prints `families` 1.
- [ ] A fixture directly under `tests/canary` prints an empty `family` cell and adds nothing to `families`.
- [ ] `--checks --full` and `--checks --check prose` exit 2 with usage.
- [ ] Both faces write no canned `go` marker, and both give the empty answer for a tree with no `tests/canary` directory.
- [ ] `bench help` and `bench test --help` hold the exact grammar text of the spec.
