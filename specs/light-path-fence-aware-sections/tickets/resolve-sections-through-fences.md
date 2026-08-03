# Resolve H2 sections through fenced code blocks

Blocked by: none
Ownership fence: `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/example_agreement_test.go`, `.agents/skills/bench-craft-tickets/SKILL.md`, `specs/light-path-fence-aware-sections/tickets/resolve-sections-through-fences.md`
Assumptions: no markdown in the tree opens a fence with more than three backticks or nests one fence inside another, so a plain three-backtick toggle resolves every document the conformance checks read; the example-agreement mutations heading already matches at any depth, so only its bite anchor is depth-bound; claims re-derived from the tree at pickup

## What to build

Section-scoped conformance anchors resolve past quoted markdown, so a skill can
teach a ticket's headings at the depth a real ticket file uses them instead of
writing them one level deeper to keep the resolver from cutting the section off.

## Acceptance

- [ ] [FA1] `markdownH2Sections` treats a `## ` line inside a fenced block as content, not a section boundary, and an unclosed fence runs the body to end of text.
- [ ] [FA2] The craft-tickets template, Good example, and Bad example carry `## What to build`, `## Acceptance`, and `## Red mutations`, and the skill no longer instructs authors to rewrite those headings at a different depth.
- [ ] [FA3] The example-agreement check grades the `## Red mutations` form, with its bite anchors and end-of-file and wrapped-field proofs still green.
- [ ] [FA4] No conformance needle or mutation row asserts a `### ` heading depth inside the craft-tickets template or examples.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FA1 | drop the fence-state tracking and restore the blind `## ` prefix scan | `TestMarkdownH2SectionsSkipsFencedHeadings` | restore the blind scan, run `go test ./internal/conformance -run '^TestMarkdownH2SectionsSkipsFencedHeadings$'`, expect the truncated-body failure |
| FA2 | demote one template heading back to `### What to build` | the `docs-currency-workflow` scoped-section anchors | demote the heading in a temp-tree copy, run `BENCH_CONFORMANCE_CHECK=docs-currency-workflow go test ./internal/conformance -run '^TestRootConformance$'`, expect the dropped-template-heading diagnostic |
| FA3 | rename the Good example's mutations heading to `## Notes` | `TestExampleAgreementParsesAuthoredLiterals` | rename inside the marked region, run `go test ./internal/conformance -run '^TestExampleAgreement'`, expect the missing Red-mutations-section diagnostic |
| FA4 | reintroduce a mutation row whose `old` string spans a `### ` template heading | `TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting` | add the row, run `go test ./internal/conformance -run '^TestSpecBuildCadenceAnchors'`, expect the anchor-count failure against the `## ` tree |
