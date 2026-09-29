# 1. Declare the scope of each public command leaf

Blocked by: none
Writes: cmd/bench/tree_scope.go (new), cmd/bench/tree_scope_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/main.go, cmd/bench/worktree_leaves.go, tests/canary/package-core-guard/unrouted-subcommand, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/subcommand_routing_test.go, internal/conformance/command_registry_parse_test.go (new), internal/conformance/command_scope_test.go (new)
Covers: TT1, TT2, TT3, TT4, TT5, TT6, TT7

## What to build

Chunk: TT-C1.

Add one scope type with the values `tree` and `repository` in `cmd/bench/tree_scope.go`. Its zero value means undeclared. Add the field to `commandDefinition` and to `commandLeaf`.
Declare the scope inline on each public registry line in `main.go` and on each `worktreeLeaves` row, from the classification table in the spec. A family definition and a plumbing definition declare none.

Make the `subcommand-routing` check grade the four declaration rules with the four exact diagnostics in the spec. Move `parseCommandRegistry` out of `subcommand_routing_test.go` into `command_registry_parse_test.go` with the same name and signature. Add a reader for the `worktreeLeaves` table beside it.

A repository-scoped definition that is not a family refuses `--in` as its first argument. It prints `toon.Usage("bench <name>", "--in")` on stdout at exit 2, before the verb runs.

`command_registry.go` holds 389 lines against a budget of 400. Put new logic in `tree_scope.go`. `main.go` and `subcommand_routing_test.go` are over budget, so keep each at or under its current line count.

## Acceptance

- [ ] The root conformance test passes on the live tree with every public leaf declared.
- [ ] A planted public definition with no scope reports `command "x" declares no scope`.
- [ ] A planted `worktreeLeaves` row with no scope reports `worktree leaf "y" declares no scope`.
- [ ] A planted plumbing definition with a scope reports `plumbing command "z" declares a scope`.
- [ ] A planted family definition with a scope reports `command family "f" declares a scope`.
- [ ] `bench version --in primary` exits 2 with `usage: bench version (unknown argument: --in)` and prints no version line.
- [ ] `bench idea --in primary x` exits 2 and writes no `capture/IDEAS.md`.
