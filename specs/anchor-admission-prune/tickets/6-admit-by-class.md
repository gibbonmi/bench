# Admit each anchor by class and retire the step kind

Blocked by: 5-prune-hand-pins.md
Writes: internal/anchors, internal/prose, cmd/bench/anchors_command.go, cmd/bench/anchor_help_test.go, cmd/bench/anchors_dir_test.go, internal/gate/lane_select_test.go, internal/gate/lane_test.go, CHANGELOG.md, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP01, AP09, AP11, AP12, AP13, AP14, AP15, AP16, AP17, AP33, AP34, AP35, AP36, AP37, AP38

## What to build

The `Anchor` type gains a class field and a reader string. Each Require row declares one of four classes: LABEL, ROUTE, INVARIANT, or COMMAND. A Forbid kind implies the FORBID class. Each LABEL row names the Go parser or the template that reads it.

An admission test in an existing anchors test file walks `Entries()` and applies the shape rules of the spec. It refuses an unclassified Require row. It refuses an INVARIANT row outside the core sections of `.bench/BENCH.md`. It refuses a LABEL row with no reader, a ROUTE needle over six words, and a COMMAND needle with no `bench ` span and argument. It refuses a needle over `prose.MaxSentenceWords` under the `internal/prose` word count. Export that counter from `internal/prose` if it is not exported.

Remove the nine step-scoped rows and their fixtures, then retire the step-scoped kind, the step field, and the step narrowing. The `bench anchors` table loses its `step` column, and `CHANGELOG.md` records that change. `bench anchors` must print a definitive empty state for a file with no rows. A non-Markdown file with no rows must take the unknown lane class.

## Acceptance

- [ ] The admission test fails for an unclassified row and for each shape-rule table case, and it passes on the live registry.
- [ ] The anchors package declares no step-scoped kind, and the ticket record shows each step symbol gone.
- [ ] `bench anchors` on a file with no rows prints a zero-row table at exit 0.
- [ ] `bench test --package ./internal/anchors`, `bench test --package ./cmd/bench`, `bench test --package ./internal/gate`, and `bench test --check docs-currency-workflow` pass.
