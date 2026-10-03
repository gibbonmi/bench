# Forbid a block-rule copy outside the block reader

Blocked by: 2-read-spec-and-coverage-blocks.md, 3-read-roadmap-blocks.md, 4-read-field-scan-blocks.md, 5-read-handoff-blocks.md, 6-read-journal-blocks.md, 7-read-review-record-blocks.md, 8-read-anchor-blocks.md, 9-read-skills-frontmatter.md, 10-read-agents-marker-blocks.md
Writes: internal/conformance/markdown_block_owner_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/checks.go, internal/conformance/tier_test.go, internal/preflight/preflighttest/fixture.go, projects/benchkit.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/skill-description-budgets/budget-table-missing, tests/canary/skill-description-budgets/description-folded, tests/canary/skill-description-budgets/description-missing, tests/canary/skill-description-budgets/over-budget-command, tests/canary/skill-description-budgets/over-budget-description, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: MB62, MB63, MB64, MB65, MB66, MB67, MB68, MB69

## What to build

Add the `markdown-block-owner` conformance check with the `go-source` input. It
follows `checkGitPlumbingOwner`. It parses each non-test Go file under `cmd/`
and `internal/`, outside `internal/markdown`, and it reads each string literal.

The check gives one diagnostic for each literal that starts with three backticks
or three tildes. It also gives one for each literal that is exactly `---`,
`<!--`, `-->`, or `## `. Each diagnostic names the file, the line, and the
literal kind.

Register the check in `registry/checks.go` and `checks_test.go`. Add its row to
the conformance table in `projects/benchkit.md`. Render the completion-plan fence
in `preflighttest/fixture.go` through the reader's fence constant. Put the bite
proof and the live-tree test in the new check file.

## Acceptance

- [ ] The bite proof plants each of the six literal kinds and gets one diagnostic for each.
- [ ] The same literals in a `_test.go` file and in `internal/markdown` give no diagnostic.
- [ ] `bench test --check markdown-block-owner` is green on the live tree.
- [ ] The profile table row and the registry row name the same check and input.
