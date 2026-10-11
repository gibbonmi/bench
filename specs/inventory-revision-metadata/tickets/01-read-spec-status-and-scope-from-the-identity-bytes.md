# Read spec status and scope from the identity bytes

Blocked by: none
Writes: internal/commitment/repository/repository.go, internal/commitment/repository/sources.go, internal/commitment/store_test.go, internal/spec/spec.go, internal/spec/metadata.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, CHANGELOG.md, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go
Covers: none

## What to build

`bench commitment inventory` lists each staged spec as a `deliverable` row. The inventory reads the identity of that row from the spec bytes at the revision that plan validates. But it reads the `Status:` filter and the `Roadmap:` scope from the working tree. An uncommitted edit of one of these lines thus gives a row that does not describe the bytes that its identity binds. Plan then binds a row with an incorrect scope, or it cannot bind a row that the inventory omits.

The inventory now lists the specs that the revision holds. For each spec, it takes the status, the scope, and the identity from the same bytes. An uncommitted edit changes no row, and a spec that the revision does not hold stays unlisted. A spec that Git cannot read at the revision refuses the inventory with the current refusal sentence. The row fields, the row order, and the roadmap and run rows stay the same.

The spec package keeps one parser for the `Status:` and `Roadmap:` lines, and the inventory and `spec.Facts` both use it. A new exported function wraps that parser. The parser thus stays in `spec.go`, where the staged `markdown-block-reader` spec changes it. `spec.Facts` keeps its working-tree read for its other callers.

## Acceptance

- [ ] A repository commits a staged spec with `Roadmap: <row A>`, and the working tree changes that line to `Roadmap: <row B>`. The inventory row of that spec has `<row A>` as its scope.
- [ ] A repository commits a spec with `Status: staged`, and the working tree changes it to `Status: implemented`. The inventory lists that spec as a `deliverable` row with the identity of the committed bytes.
- [ ] A repository commits a spec with `Status: draft`, and the working tree changes it to `Status: staged`. The inventory does not list that spec.
- [ ] `TestCommitmentInventoryIdentity` stays green: an uncommitted spec stays unlisted, and each listed identity equals the identity that plan accepts.
- [ ] `CHANGELOG.md` has one `Fixed` entry for the inventory.
