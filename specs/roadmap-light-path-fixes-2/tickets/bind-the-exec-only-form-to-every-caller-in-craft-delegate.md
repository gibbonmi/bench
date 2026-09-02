# Bind the exec-only form to every caller in craft-delegate

Blocked by: state-the-fence-order-and-the-claim-words-in-craft-spec.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/benchguard/benchguard_test.go, tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller (new), tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package (new), projects/benchkit.md, tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green, tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer, tests/canary/workflow-guidance-anchors/delegate-model-id-escalation, tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface, tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor, tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance, tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents, tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row, tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor, tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor, tests/canary/workflow-guidance-anchors/shared-worktree-path-pin, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership
Covers: LP16, LP17

## What to build

Ship two sentences in the five-part precedent shape, on the integration source
after the blocker ticket lands. In the Isolation section, rewrite the exec-only
sentence. `bench worktree exec "<label>" -- <command>` is the one command form
for every caller into an assignment worktree. The rule covers the coordinator,
and it covers a read or a write. A shell loop inside the pool path is the same
bypass. In the charge section, beside the registry sentence, say a cap-change
charge's search list names the closest pinning package.

Read `TestPoolReferenceRefusesAPoolPathOutsideTheExecVerb` first. Add a
guard-table row in `internal/benchguard/benchguard_test.go` only if the guard
accepts a `for` loop whose body names the pool path. Report what you found.
Keep `SKILL.md` inside its 122-line budget row.

## Acceptance

- [ ] The registry test reds a synthetic tree that drops either sentence and stays silent on the live root.
- [ ] Both fixtures bite through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] The prose budget check passes on the worktree.
- [ ] The report names whether the pool-reference guard refuses a shell loop over the pool path today.
- [ ] Self-probe: restore the old charge-only wording, and report the registry test red.
