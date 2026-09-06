# Make headroom in the three registry files

Blocked by: none
Writes: cmd/bench/main.go, cmd/bench/anchors_command.go (new), cmd/bench/main_test.go, cmd/bench/help_inventory_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, internal/conformance/subcommand_routing_table_test.go (new), internal/tickets/registry_data.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: PB32, PB33

## What to build

Verify the premise first. Run `bench structure` and confirm that `cmd/bench/main.go`,
`cmd/bench/main_test.go`, and `internal/conformance/subcommand_routing_test.go` sit
over the 400-line budget. Read `Growth` in internal/structure/structure.go, which
reds an over-budget source file that gained lines since the base. Read `dispatchFile`
in internal/conformance/subcommand_routing_test.go, which keeps naming
`cmd/bench/main.go`. Read `commandRegistries` in internal/tickets/registry_data.go.

Move three cohesive blocks, each as a pure move with no edit inside the block:

- `anchorsGrammar`, `anchorsCommand`, `anchorQueryPath`, and `anchorKindName` from
  cmd/bench/main.go to the new file cmd/bench/anchors_command.go.
- `TestHelpRendersPublicCommandRegistryRows` and `TestHelpInventoryIsComplete` from
  cmd/bench/main_test.go to the new file cmd/bench/help_inventory_test.go.
- The `subcommandRouting` table and the why constants it uses from
  internal/conformance/subcommand_routing_test.go to the new file
  internal/conformance/subcommand_routing_table_test.go.

In `commandRegistries`, replace `cmd/bench/main_test.go` with
`cmd/bench/help_inventory_test.go` and `internal/conformance/subcommand_routing_test.go`
with `internal/conformance/subcommand_routing_table_test.go`. Add the row
`{Prefix: "internal/probe", Files: commandRegistries}` in sorted position. Read
the tests in internal/tickets that grade the binding table before you edit it.

Run `gofmt` and `go vet` on the two packages. Run the moved tests from their new
files. Run `bench structure --growth <base>` from the worktree and record the ok line.
Stop and report when a check pins one of the moved names by file.

## Acceptance

- [ ] `bench structure --growth <base>` prints the ok line, and each of the three files counts fewer lines than at the base.
- [ ] `go test ./cmd/bench -run 'Help'` and `go test ./internal/conformance -run 'Routing'` pass from the new files.
- [ ] `go test ./internal/tickets` passes with the two renamed registry files and the `internal/probe` row.
