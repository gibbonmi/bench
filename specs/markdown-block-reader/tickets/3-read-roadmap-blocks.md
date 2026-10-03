# Read roadmap lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/roadmap/tree.go, internal/roadmap/tree_helpers_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB31, MB32, MB33

## What to build

Replace the three fence toggles in `internal/roadmap/tree.go` with the block
reader's classification. `ParseDocument` skips a fenced line for the row grammar
and for the section title. `parseSequence` reads the reader's H2 headings, so a
fenced heading neither starts nor ends the sequence. `rowNextDiagnostics` reads
only unfenced lines for the `Next:` marker. It keeps its own column-zero and
separator grammar.

`internal/roadmap` is a crowded directory, so the ticket adds no file. The new
cases go in `tree_helpers_test.go`.

## Acceptance

- [ ] A tilde-fenced row example gives no row and no malformed-row failure.
- [ ] A tilde-fenced `## Recommended sequence` line opens no sequence section.
- [ ] A tilde-fenced `Next:` line after a real one gives no duplicate diagnostic.
- [ ] `bench roadmap` prints the same output on the live tree at the base and at the tip.
