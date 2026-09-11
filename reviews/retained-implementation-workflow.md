# Retained implementation workflow review pickup

Frozen delta: `57fb766d58954a603b51d521ddcd203db4764441..9a3409e1b05e45583ceb41cca05dfc9533914ba0`.

## Standards

Finding count: 2. Worst issue: P2.

1. [P2] `auto-fix` — Reuse the shared anchor fixture harness. `internal/conformance/retained_workflow_test.go:49` duplicates the lifecycle in `internal/conformance/docs_workflow_helpers_test.go:450` and omits its Forbid and duplicate-subject protections. Extend the shared helper with section-aware fixture text.
2. [P2] `auto-fix` — Reference the implementation recommendation contract. `.agents/commands/bench-write-spec.md:48` restates factors and harder-chunk guidance owned by `.agents/skills/bench-craft-spec/SKILL.md:28`.

## Spec

Finding count: 1. Worst issue: P2.

1. [P2] `auto-fix` — Clarify review-rule precedence. `.agents/skills/bench-craft-line/SKILL.md:68` sends kit guidance through mid/high in every stage, while line 73 sends a Sol implementation review to top/high. Restrict the leverage override or make the conditional review line authoritative for review.

## Coverage

Finding count: 1. Worst issue: P2.

1. [P2] `auto-fix` — Pin the policy predicates independently. `internal/conformance/retained_workflow_test.go:67` derives healthy fixture text from the registry predicate. A review probe changed the default route from mid to top, and `TestRetainedWorkflow` stayed green. Add independent expectations for W12's reason factors, both W14 branches, and W15's switch and retained-effort boundary.
