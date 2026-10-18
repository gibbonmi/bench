# Admit each anchor by class

Blocked by: 6a-root-parameter-grammars.md
Writes: internal/anchors, internal/prose, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP01, AP09, AP11, AP12, AP13, AP14, AP15, AP16, AP17

## What to build

The `Anchor` type gains a class field and a reader string. Each Require row declares one of four classes: LABEL, ROUTE, INVARIANT, or COMMAND. A Forbid kind implies the FORBID class. Each LABEL row names the Go parser or the template that reads it. Put the class on each row's existing line, so that no registry file grows.

An admission test walks `Entries()` and applies the shape rules of the spec. Put it in an existing anchors test file with line headroom, such as `match_test.go` or `locate_test.go`. The test refuses each of these rows:

- a Require row with no class;
- an INVARIANT row outside the core sections of `.bench/BENCH.md`;
- a LABEL row with no reader;
- a ROUTE needle over six words, or with no route token;
- a COMMAND needle with no `bench ` span and argument;
- a needle over `prose.MaxSentenceWords` under the `internal/prose` word count.

State the exact token rule for a COMMAND argument in the test. It must accept each COMMAND row of `assets/classified.md` and refuse each argument-free verb sentence there. Export the one-sentence word counter from `internal/prose` if it is not exported. No touched file may grow.

## Acceptance

- [ ] The admission test fails for an unclassified row and for each shape-rule table case, and it passes on the live registry.
- [ ] The admission test fails on the live tree when one unclassified Require row is added.
- [ ] `bench test --package ./internal/anchors`, `bench test --package ./internal/prose`, `bench test --package ./internal/conformance`, and `bench test --check docs-currency-workflow` pass.
