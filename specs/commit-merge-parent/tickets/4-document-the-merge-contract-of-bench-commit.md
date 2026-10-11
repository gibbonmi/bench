# Document the merge contract of bench commit

Blocked by: 1-read-the-pending-merge-through-one-owner.md, 2-refuse-an-incomplete-merge-name-set.md, 3-record-each-merge-head-and-end-the-merge.md
Writes: .bench/BENCH-reference.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/agents-handoff-section-rule, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route
Covers: MP55

## What to build

A teammate who walks in cold reads the merge contract of `bench commit` in the working reference.
The first commit after a hand merge then names every merged path, and no refusal has to teach that rule.

Tickets 2 and 3 supply the behavior that the new sentences state.
Add these two sentences to the Command Notes of `.bench/BENCH-reference.md`, after the fast-lane note:

> During a pending merge, a worktree `bench commit` requires every path that the merge changed. It records each merge head as a further parent and then ends the merge.

Keep that paragraph at six sentences or fewer, and keep every anchor of the file intact.
The Command Notes needle `runs the fast lane on a private checkout of the composed snapshot, and a lane pass` sits in that paragraph.
The route needle `the reviewer merges the destination into the source worktree with raw Git` sits later in the same section.
The MP2 review reads the two sentences and the `bench anchors .bench/BENCH-reference.md` output for MP55.

The anchor registry files name `.bench/BENCH-reference.md`, and the binding registry binds the command registry files to `internal/anchors`.
The canary units on this `Writes:` line pin `.bench/BENCH-reference.md`.
Each unit replaces one sentence outside the fast-lane paragraph, or it puts a short copy in place of the live file.
The new sentences therefore change no unit.
All these files join this `Writes:` line, and this ticket does not change them.

## Acceptance

- [ ] The Command Notes of `.bench/BENCH-reference.md` state that a worktree `bench commit` requires every merged path during a pending merge.
- [ ] The Command Notes state that the commit records each merge head as a further parent and then ends the merge.
- [ ] The two sentences follow the fast-lane note in the Command Notes.
