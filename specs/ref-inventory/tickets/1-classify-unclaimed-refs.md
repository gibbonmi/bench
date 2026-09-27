# 1. Classify each unclaimed ref and narrow the bulk sweep

Blocked by: none
Writes: internal/worktree/clean_classes.go (new), internal/worktree/clean_classes_test.go (new), internal/worktree/clean_unclaimed.go, internal/worktree/clean_unclaimed_test.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_outcomes_test.go, internal/worktree/clean_set_wiring_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, CONTEXT.md, internal/anchors/registry_data.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary
Covers: RI1, RI2, RI3, RI4, RI5, RI6, RI7, RI8, RI9, RI10, RI11, RI12, RI13, RI14, RI15, RI16, RI17, RI18, RI19, RI20, RI21, RI22, RI23, RI55, RI57

## What to build

Chunk: RI-C1a.

Add the class function beside the unclaimed planner.
It takes the protected set, the default branch, and the sorted unrecorded refs.
It returns one class and one holder per ref.
It calls `git.LandedInDefault` for each unrecorded ref.
A ref whose tip does not resolve to a commit produces an error row with no action, and the set fingerprint stays empty.

Among refs with an equal tip, the lexically first full ref name is the root.
A ref is subsumed when its tip equals the tip of a lexically earlier unrecorded ref.
A ref is also subsumed when its tip is a strict ancestor of a protected branch, a landed ref, or another unrecorded ref.
The holder is the lexically first ref that satisfies the rule.
Protected branches and landed refs order before unrecorded refs in that search.

The plan row keeps the seven cleanup columns and the `tracked` cell `unclaimed`.
The detail cell starts with `class=<class>`, continues with ` holder=<ref>` for a subsumed row, and ends with the existing removal text for a removing row.
A unique row carries the action `retain`, and its detail ends with `bench worktree clean --discard-branch --target <assignment id>`.
The fingerprint version becomes `bench-unclaimed-assignment-branches/v2` and binds each row's ref, tip, class, and holder.
The apply loop, under a fingerprint and under `--apply-current`, skips a row whose action does not remove.
The plan prints the apply help action only when a row removes.

Two existing tests plant a unique branch and expect its removal.
Rewrite them so that a landed branch is the removal case and the unique case is retained.
Repair the glossary term `unclaimed ref` to name `refs/heads/bench/shift-`.

## Acceptance

- [ ] A plan over the seven ref shapes of the spec prints the class and holder of the spec for each row.
- [ ] A branch ref that points at a blob prints action `error`, an empty fingerprint, and no apply help action.
- [ ] An apply with the fingerprint, and an `--apply-current` run, each remove the landed and subsumed refs and keep the unique ref.
- [ ] After each apply, every deleted tip is reachable from `main` or from a surviving ref.
- [ ] A plan, one new commit on a subsumed ref, then the old fingerprint refuses as stale and prints `bench worktree clean --discard-branch --unclaimed`.
- [ ] A plan over unique rows only prints no apply help action, and a plan over no Bench-namespace branch prints the empty table and exits 0.
- [ ] Two consecutive plans print identical rows and `for-each-ref refs/bench/` is unchanged between them.
- [ ] `TestLandCommandPrunesSquashFoldedSiblingBranch` still passes.
- [ ] The glossary term `unclaimed ref` names `refs/heads/bench/shift-`.
