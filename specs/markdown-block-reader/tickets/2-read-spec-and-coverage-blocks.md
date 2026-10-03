# Read spec and coverage lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/spec/spec.go, internal/spec/fences.go, internal/spec/fences_test.go, internal/coverage/coverage.go, internal/coverage/coverage_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB24, MB25, MB26, MB27, MB28, MB29, MB30, MB80, MB83

## What to build

Make `LiveSpecSlugs`, `metadata`, `deriveImplemented`, and `FenceTokens` in
`internal/spec` read one fence classification from the block reader. Each one
matches its grammar on the raw text of the unfenced lines. The trim
inconsistency between `LiveSpecSlugs` and `metadata` then closes.
`deriveImplemented` writes its output from the source bytes by the reader's line
offsets, so a CRLF spec flips and keeps each carriage return.

`FenceTokens` opens its section on the reader's H2 heading with the title
`Ownership fences`, at column zero only. It keeps its own rule that any `#{2,} `
heading ends the section, applied to unfenced lines only.

Make `coverage.parse` read the block reader on raw text. A fenced line is no
story number, no `Not covered:` line, and no map row. The historical marker stays
a raw-text match, because the marker is itself a comment. The reader's H2 heading
replaces the `## ` prefix test, so a fenced heading does not end the story list.

The MB25 fixture puts its `Status: implemented` line at column zero inside the
indented fence, because `metadata` reads that label only at column zero. The new
spec rows go in `internal/spec/fences_test.go`, because
`internal/spec/spec_test.go` is over the line cap. Each exported signature stays
the same.

## Acceptance

- [ ] A spec with a real `Status: staged` line and a later tilde-fenced `Status: implemented` line reads as staged.
- [ ] An indented fence that holds a spec path and a column-zero status line hides both from `LiveSpecSlugs` and from `Facts`.
- [ ] `Implemented` flips only the real `Status: staged` line when a fenced copy exists, and it keeps each carriage return of a CRLF spec.
- [ ] A fenced example inside `## Ownership fences` gives no token, and a fenced heading opens no section.
- [ ] A fenced map row and a fenced story line add nothing, and a story after a fenced `## ` line counts.
- [ ] The historical marker keeps its exemption.
- [ ] `bench coverage --check` gives the same result on each staged spec at the base and at the tip.
- [ ] `internal/spec` and `internal/coverage` hold no string literal with a run of three backticks, and no literal `## `. The sites today are `spec.go:71`, `spec.go:93`, and `coverage.go:200`.
- [ ] `internal/spec/spec.go` and `internal/coverage/coverage.go` do not grow.
