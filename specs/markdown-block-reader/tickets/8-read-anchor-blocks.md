# Read anchor sections and comments from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/anchors/locate.go, internal/anchors/match.go, internal/anchors/locate_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB52, MB53, MB54, MB55, MB56

## What to build

Make `stripCommentsMapped` in `internal/anchors` a projection of the reader's
removed comment spans. It keeps its rune-origin mapping, so `Locate` names the
same physical lines. Make `scopeRunesMapped` read the reader's fence classes
instead of its own backtick toggle. The section closer reads the reader's H2
heading. The section opener keeps its exact title match, and the step grammar
stays in `internal/anchors`.

Delete `commentOpen`, `commentClose`, `fenceMark`, and `headingMark`. Update the
comments in `match.go` that name the backtick fence rule.

`internal/anchors` is a crowded directory, so the ticket adds no file. The new
cases go in `locate_test.go`.

## Acceptance

- [ ] A needle after a tilde-fenced `## Other` line inside section `S` locates inside `S`.
- [ ] A tilde-fenced `## S` line before the real one gives a section count of 1.
- [ ] A tilde-fenced step opener gives no step body.
- [ ] `TestLocateStripsRejoinedComments` stays green without an edit.
- [ ] `bench test --check docs-currency-workflow` is green, so each live anchor keeps its verdict.
