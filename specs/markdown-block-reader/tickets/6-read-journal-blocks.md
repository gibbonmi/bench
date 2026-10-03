# Read learning and retrospective lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/learnings/learnings.go, internal/learnings/entry.go, internal/learnings/learnings_test.go, internal/retros/retros.go, internal/retros/retros_test.go, internal/retros/recommendations.go, internal/retros/recommendations_test.go, tests/canary/package-core-guard/bounds-classify-limit-restated, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB41, MB42, MB43, MB44, MB45, MB46, MB47, MB84

## What to build

Make `learnings.Parse` and `hasAnyHeading` read the reader's H2 headings on raw
text, so a fenced dated heading stays in its entry body. The body keeps its raw
lines. `Parse` adds one malformed record for the reader's fault, at the fault
line. The reasons are `unterminated fenced block`, `unterminated HTML comment`,
and `unterminated frontmatter block`. The prose mechanics check prints each
record through its existing malformed-record path.

One exception stays: `unaccountedRegion` reads raw lines, with no fence rule,
because its rule is that nothing belongs below the entries marker. It tests the
H2 prefix through the reader's exported constant. DL31 stays green without an
edit.

Make `retros.Parse` and `Recommendations` read the reader on raw text. A fenced
required heading does not count. A fenced `## ` line does not end the
improvements section. A fence marker line and each fenced line are no
recommendation unit, and each one ends the open unit.

Render each heading through the reader's H2 constant, so neither package keeps a
`## ` literal. The learnings rows go in `learnings_test.go`, because
`entry_test.go` tests only the formatter. The `retros.Parse` row goes in
`retros_test.go`.

## Acceptance

- [ ] An entry body with a fenced dated heading gives one entry whose body holds the fenced line.
- [ ] An unterminated fence, comment, and frontmatter block each give one malformed record with its reason and line.
- [ ] `TestParseReportsUnaccountedContentBelowTheEntriesMarker` stays green without an edit.
- [ ] A retrospective whose only `## Outcome` line is fenced fails `Parse` with the missing-heading error.
- [ ] An item after a fenced `## Notes` line counts, and a fenced `- item` line gives no recommendation.
- [ ] `TestRecommendationsKeepsImprovementParagraphsAndListItemsSeparate` stays green without an edit.
- [ ] `learnings.Parse` of the primary journal and `retros.Parse` of each file under `capture/retros/` give the same result at the base and at the tip.
- [ ] `internal/learnings` and `internal/retros` hold no literal `## `. The sites today are `learnings.go:87`, `:107`, `:146`, `:182`, `:197`, `:218`, `:252`, `:266`, `entry.go:26`, and `recommendations.go:22`, `:23`, and `:39`.
- [ ] No file in `Writes:` grows past 400 lines.
