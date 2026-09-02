# State the census read, the CHANGELOG rule, and the review base

Blocked by: state-the-helper-return-rule-in-craft-tdd.md
Writes: .agents/commands/bench-final-check.md, .agents/commands/bench-review-implementation.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/final-check-census-read-before-land (new), tests/canary/workflow-guidance-anchors/final-check-light-path-changelog-heading (new), tests/canary/workflow-guidance-anchors/review-base-merged-main-tip (new), tests/canary/docs-currency-token-diet/introduces-undeclared-command, tests/canary/docs-currency-token-diet/stale-command-reference, tests/canary/workflow-guidance-anchors/final-check-bare-leftover-clean-retired, tests/canary/workflow-guidance-anchors/final-check-landed-worktree-sweep, tests/canary/workflow-guidance-anchors/final-check-scratch-branch-clean, tests/canary/workflow-guidance-anchors/coverage-axis-anchor, tests/canary/workflow-guidance-anchors/review-falsification-accept-routing, tests/canary/workflow-guidance-anchors/review-falsification-dispositions, tests/canary/workflow-guidance-anchors/review-kit-guidance-set, tests/canary/workflow-guidance-anchors/review-persistence-anchor, tests/canary/workflow-guidance-anchors/review-preflight-explicit-base, tests/canary/workflow-guidance-anchors/review-repair-ticket-covers, tests/canary/workflow-guidance-anchors/review-repair-ticket-owner, tests/canary/workflow-guidance-anchors/review-standing-falsification, tests/canary/workflow-guidance-anchors/review-universal-claim-bar, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LP19, LP20, LP21, LP22, LP23

## What to build

Ship three sentences in the five-part precedent shape, on the integration
source after the blocker ticket lands. In `bench-final-check.md`, before the
landing step, add the census rule. The phase close reads the assignment census
record before `bench worktree land` removes it, and it carries the per-verb
breakdown into the close. In the same file, add the CHANGELOG rule. A
light-path fix lands before a spec's final merge only when its `CHANGELOG.md`
entry sits under a heading no sibling touches.

In `bench-review-implementation.md`, in the pin-the-diff step, name the frozen
base. It is the `main` tip merged into the source before the landing, so the
range holds the spec diff alone. Neither file has a budget row.

This ticket is the last one to touch the anchor registry and the fixture
directory. It therefore carries the whole-family invariant: every fixture this
spec added bites, and the guidance token sweep passes.

## Acceptance

- [ ] The registry test reds a synthetic tree that drops each of the three sentences and stays silent on the live root.
- [ ] The census sentence sits before the sentence that reads `census=<n>` from the landed record.
- [ ] Every fixture this spec added bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] The guidance token sweep and the prose budget check pass on the worktree.
- [ ] Self-probe: move the census sentence after the landed-record sentence, and report which check reds it or that none does.
