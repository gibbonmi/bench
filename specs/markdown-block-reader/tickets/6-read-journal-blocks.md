# Read learning and retrospective lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/learnings/learnings.go, internal/learnings/entry.go, internal/learnings/entry_test.go, internal/retros/retros.go, internal/retros/recommendations.go, internal/retros/recommendations_test.go, tests/canary/package-core-guard/bounds-classify-limit-restated, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB41, MB42, MB43, MB44, MB45, MB46, MB47

## What to build

Make `learnings.Parse` and `hasAnyHeading` read the reader's H2 headings, so a
fenced dated heading stays in its entry body. The body keeps its raw lines.
`Parse` adds one malformed record for the reader's fault, at the fault line. The
reasons are `unterminated fenced block`, `unterminated HTML comment`, and
`unterminated frontmatter block`. The prose mechanics check prints each record
through its existing malformed-record path.

Make `retros.Parse` and `Recommendations` read the reader. A fenced required
heading does not count. A fenced `## ` line does not end the improvements
section. A fence marker line and each fenced line are no recommendation unit,
and each one ends the open unit.

Render each heading through the reader's H2 constant, so neither package keeps a
`## ` literal. `internal/learnings/learnings.go` stays within its line budget.

## Acceptance

- [ ] An entry body with a fenced dated heading gives one entry whose body holds the fenced line.
- [ ] An unterminated fence, comment, and frontmatter block each give one malformed record with its reason and line.
- [ ] A retrospective whose only `## Outcome` line is fenced fails `Parse` with the missing-heading error.
- [ ] An item after a fenced `## Notes` line counts, and a fenced `- item` line gives no recommendation.
- [ ] `TestRecommendationsKeepsImprovementParagraphsAndListItemsSeparate` stays green without an edit.
