# 8. List the fixtures that a check owns

Blocked by: 6-filter-system-suite.md
Writes: internal/testreport/, internal/canary/inventory.go, internal/canary/inventory_test.go, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: TP27, TP28, TP29, TP30, TP33, TP40, TP52, TP53

## What to build

Chunk: TP-C3.

Add `bench test --check <name> --fixtures`. It prints `fixtures[N]{family,fixture,path}`, sorted by path, with one row for each fixture whose `canary.Fixtures` owner is the check.
The face reads `tests/canary` under the graded repository root. It selects no run binary and starts no Go child.

Export one no-fixtures sentinel error from `internal/canary`, and make `Fixtures` and `FixturePins` both use it, with no line growth in `inventory.go`.
An absent or empty `tests/canary` directory gives the empty table at exit 0, through `errors.Is` on that sentinel. Each other inventory error exits 1 with the inventory diagnostic.
`--fixtures` requires `--check` and accepts no other flag. An unknown check gives the unknown-check refusal.

Add the form `bench test --check <name> --fixtures` to the grammar text and to the `bench help` row.
This ticket supplies the owner filter that ticket 9 counts. Rows TP31, TP32, and TP39 name both faces, so ticket 9 owns them.

## Acceptance

- [ ] Two fixtures of one family print two rows with repo-relative paths.
- [ ] A `CHECK`-marked fixture prints under its marked check and not under its family owner.
- [ ] `--check system --fixtures` prints the empty table at exit 0.
- [ ] The face writes no canned `go` marker and calls no run binary builder.
- [ ] A tree with no `tests/canary` directory prints the empty table at exit 0.
- [ ] A present and empty `tests/canary` directory prints the empty table at exit 0.
- [ ] `FixturePins` still returns an empty map for a tree with no fixture.
- [ ] An owned fixture whose name holds U+0001 exits 1 with the render error.
- [ ] A `CHECK` file that names `no-such-check` exits 1 with `names unknown check`.
- [ ] `--fixtures` alone, and `--fixtures` with `--full` or `--run`, exit 2 with usage.
- [ ] `--check not-registered --fixtures` gives the unknown-check refusal.
