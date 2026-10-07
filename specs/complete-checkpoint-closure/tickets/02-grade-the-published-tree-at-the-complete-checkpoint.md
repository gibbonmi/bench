# Grade the published tree at the complete checkpoint

Blocked by: 01-move-fixture-witnesses-to-the-common-directory.md
Writes: internal/gate/gate.go, internal/gate/checkpoint.go, internal/gate/engine.go, internal/gate/complete_checkpoint.go (new), internal/gate/complete_checkpoint_test.go (new), internal/gate/review_checkpoint_commits_test.go, internal/gate/prospective_owner_test.go, .bench/BENCH-reference.md, CHANGELOG.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_retained_workflow.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/agents-handoff-section-rule, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: CC01, CC02, CC03, CC04, CC05, CC06, CC07, CC08, CC09, CC10, CC11, CC12, CC13, CC14, CC15, CC16, CC17, CC18, CC19, CC20, CC21, CC22, CC23, CC24, CC25, CC26, CC27, CC28, CC29, CC30, CC31, CC32

## What to build

`bench gate --checkpoint <spec> --complete` grades the tree that the landing publishes from the committed source tip.
Every other gate run, the chunk checkpoint included, keeps its current evaluation.

When the parsed checkpoint selects `--complete`, `RunCommand` takes the completion route.
The route does these steps in this order:

1. It refuses a checkout whose working-tree hash differs from the tree of `HEAD`. The refusal exits 1 before the oracle runs. Its text contains `complete checkpoint requires a clean checkout: commit or remove each uncommitted change`. The review record has no exemption.
2. It resolves `HEAD` as the source tip.
3. It composes the published tree with `published.Tree` from the tip's tree, the spec path, and the tip. A compose error refuses with exit 1 and its own reason, before the oracle runs.
4. It grades that tree through the prospective execution that `ExecuteTree` uses, with `WithCompletion` for the spec and the tip.

The route keeps the caller's run mode, so `--fresh` still forces execution.
It keeps the signal arm that the CLI path installs today.
`published.Tree` is the one closure derivation; the route adds no second one.
A change to `executeTreeWithOwner` updates its four call sites in `internal/gate/prospective_owner_test.go`.

Promote `TestCollisionCompletionCheckpointGradesTheClosedTree` from the `landing-collisions` worktree into `internal/gate/complete_checkpoint_test.go`, with no build constraint.
Add the planned tests that the coverage map names in that file.
Each new expectation derives from the fixture inputs; CC05 and CC31 derive the expected tree through `published.Tree`.

The `complete` case of `TestReviewCheckpointCommentOnlyGap` commits its completion record before `--complete`, and its assertions stay unchanged.

The reference guide gets one sentence beside the `--complete` text: `--complete` grades the published tree of the committed source.
The build's `CHANGELOG.md` entry states the published-tree grade and the clean-checkout refusal, and it records Git 2.31 as the minimum Git version.

## Acceptance

- [ ] The complete checkpoint on the promoted repro source exits 9 and prints `decisions/x.md cites missing roadmap/FT1.md`.
- [ ] An ordinary gate run and a chunk checkpoint on the promoted repro source exit 0.
- [ ] After a green complete checkpoint, `InspectTreeContext` reports a reusable green for the `published.Tree` result of `HEAD` under `WithCompletion`.
- [ ] After a green complete checkpoint, `.gate-record-during` exists in the Git common directory.
- [ ] A repeated complete checkpoint on an unchanged source reuses its green, and `--fresh` adds exactly one run.
- [ ] A complete checkpoint leaves `HEAD` and an empty `git status --porcelain` in the checkout.
- [ ] A tracked edit, an untracked file, or an uncommitted review record refuses with the clean-checkout text before the oracle runs, and before the stale-evidence refusal.
- [ ] A spec with no `Status: staged` line refuses with `spec has no Status: staged line` before the oracle runs.
- [ ] `TestReviewCheckpointReuse`, `TestReviewCheckpointKeepsStrictEvidence`, `TestReviewCheckpointCommentOnlyGap`, `TestCommitmentExactTransform`, `TestGateCheckpointRoute`, and `TestLandCommandPublicRealGitJourney` pass.
- [ ] `internal/gate/complete_checkpoint_test.go` carries no build constraint.
