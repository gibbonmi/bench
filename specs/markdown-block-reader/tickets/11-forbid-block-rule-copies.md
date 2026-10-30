# Forbid a block-rule copy outside the block reader

Blocked by: 2-read-spec-and-coverage-blocks.md, 3-read-roadmap-blocks.md, 4-read-field-scan-blocks.md, 5-read-handoff-blocks.md, 6-read-journal-blocks.md, 7-read-review-record-blocks.md, 8-read-anchor-blocks.md, 9-read-skills-frontmatter.md, 10-read-agents-marker-blocks.md
Writes: internal/conformance/markdown_block_owner_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/checks.go, internal/conformance/tier_test.go, internal/conformance/tier_live_tree_test.go, internal/conformance/skill_description_budgets_test.go, internal/conformance/skills_index_checks_test.go, internal/conformance/docs_workflow_checks_test.go, internal/preflight/preflighttest/fixture.go, internal/preflight/ownership/inventorytest/inventory.go, projects/benchkit.md, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: MB62, MB63, MB64, MB65, MB66, MB67, MB68, MB69, MB98, MB101

## What to build

Add the `markdown-block-owner` conformance check with the `go-source` input. It
follows `checkGitPlumbingOwner` (`git_plumbing_owner_test.go:27`). It parses each non-test Go file under `cmd/`
and `internal/`, outside `internal/markdown`, and it reads each string literal.

The check gives one diagnostic for each literal that contains a run of three
backticks or three tildes anywhere in its value. It also gives one for each literal that is exactly `---`,
`<!--`, `-->`, or `## `. Each diagnostic names the file, the line, and the
literal kind.

Move three conformance readers onto the block reader: `skillDescriptionFolded`
(`skill_description_budgets_test.go:140`), `frontmatterHasKey`
(`skills_index_checks_test.go:249`), and `h2Headings`
(`docs_workflow_checks_test.go:743`).
Each one is a test file that the check exempts, but each one copies a block rule.

Register the check in `registry/checks.go` and `checks_test.go`. The
registration lines stay within the `checks_test.go` budget in
`.bench/structure.budgets`. Add its row to
the conformance table in `projects/benchkit.md`. Render the completion-plan fence
at `preflighttest/fixture.go:111` and the planning-facts fence at
`inventorytest/inventory.go:23` through the reader's fence constant. Put the bite
proof and `TestMarkdownBlockOwnerHoldsOnTheLiveTree` in the new check file.

When `testChecksLiveTree` (`tier_test.go:261`) detects the assertion's body, add its
entry to `classifiedLiveTreeTests` (`tier_live_tree_test.go:19`). Keep the hidden-reader
metadata check and `TestClassifiedLiveTreeInventoryNamesDetectedTests` intact.
The former rejects an omitted classification, and the latter rejects a stale one.

## Acceptance

- [ ] The bite proof plants each of the six literal kinds and gets one diagnostic for each.
- [ ] The same literals in a `_test.go` file and in `internal/markdown` give no diagnostic.
- [ ] `bench test --check markdown-block-owner` is green on the live tree.
- [ ] The detected `TestMarkdownBlockOwnerHoldsOnTheLiveTree` assertion has its entry in `classifiedLiveTreeTests`.
- [ ] Removing that required entry gives the unregistered-live-tree assertion diagnostic through `conformance-meta`.
- [ ] `bench test --check conformance-meta` and `TestClassifiedLiveTreeInventoryNamesDetectedTests` stay green.
- [ ] The profile table row and the registry row name the same check and input.
- [ ] `skillDescriptionFolded`, `frontmatterHasKey`, and `h2Headings` read the block reader, and `skill-description-budgets`, `skills-index-command-adapters`, and `docs-currency-workflow` stay green.
- [ ] `internal/preflight/preflighttest/fixture.go` and `internal/preflight/ownership/inventorytest/inventory.go` hold no string literal that contains a run of three backticks.
- [ ] `skill_description_budgets_test.go`, `skills_index_checks_test.go`, and `docs_workflow_checks_test.go` do not grow.
- [ ] `checks_test.go` stays within its declared budget, and the new check file stays under the structure file limit.
