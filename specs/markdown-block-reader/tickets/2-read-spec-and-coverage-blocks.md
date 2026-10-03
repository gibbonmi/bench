# Read spec and coverage lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/spec/spec.go, internal/spec/fences.go, internal/spec/fences_test.go, internal/coverage/coverage.go, internal/coverage/coverage_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB24, MB25, MB26, MB27, MB28, MB29, MB30

## What to build

Make `LiveSpecSlugs`, `metadata`, `deriveImplemented`, and `FenceTokens` in
`internal/spec` read one fence classification from the block reader. The trim
inconsistency between `LiveSpecSlugs` and `metadata` then closes. `FenceTokens`
opens its section on the reader's H2 heading with the title `Ownership fences`.
It keeps its own rule that any `#{2,} ` heading ends the section, applied to
unfenced lines only.

Make `coverage.parse` read the block reader. A fenced line is no story number,
no `Not covered:` line, and no map row. The reader's H2 heading replaces the
`## ` prefix test, so a fenced heading does not end the story list.

`internal/spec/spec.go` and `internal/coverage/coverage.go` are over their line
budgets, so neither file grows. Each exported signature stays the same.

## Acceptance

- [ ] A spec with a real `Status: staged` line and a later tilde-fenced `Status: implemented` line reads as staged.
- [ ] An indented fence that holds a spec path and a status line hides both from `LiveSpecSlugs` and from `Facts`.
- [ ] `Implemented` flips only the real `Status: staged` line when a fenced copy exists.
- [ ] A fenced example inside `## Ownership fences` gives no token, and a fenced heading opens no section.
- [ ] A fenced map row and a fenced story line add nothing, and a story after a fenced `## ` line counts.
- [ ] `bench coverage --check` gives the same result on each staged spec at the base and at the tip.
