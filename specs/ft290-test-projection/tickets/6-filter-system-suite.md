# 6. Filter the system suite with a run pattern

Blocked by: 3-prove-named-check-ran.md
Writes: internal/testreport/, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: TP19, TP20, TP21, TP22, TP23

## What to build

Chunk: TP-C2.

Accept `--run <regex>` with `--check system`, and append `-run <regex>` after the system suite operands.
`Request.Run()` returns the pattern for that request.
Today one case of `Request.Run()` returns `AllTests` for both the prose and the system check. Split that case, keep `AllTests` for the prose check, and update its comment.

Keep the usage refusal for `--run` with each other named check, before a Go child starts.
When the pattern matches no test, the refusal `go test reported no test runs` wins, and the zero-rule title does not print.

Add the form `bench test [--full] --check system --run <go-regex>` to the grammar text and to the `bench help` row, inside their current lines.
The rendered help pins in `command_registry_test.go` and `help_inventory_test.go` include the `[--in <label|primary>]` insertion that the registry renders. Keep that insertion. The `Suffix` and the `testreport` usage text have no `--in`.

## Acceptance

- [ ] The Go child of `--check system --run ^TestX$` receives `-run ^TestX$` after the suite operands.
- [ ] `Request.Run()` returns `^TestX$` for that request.
- [ ] A pattern with no match exits 1 with `go test reported no test runs` and without `named check ran nothing`.
- [ ] `--check line-routing --run ^TestX$` exits 2 and starts no Go child.
- [ ] `--check prose --run ^TestX$` exits 2.
- [ ] `Request.Run()` still returns `AllTests` for a prose request.
- [ ] `bench help` shows the system run form.
