# Refuse bench link in the kit source checkout

Blocked by: none
Writes: internal/adopt/link.go, internal/adopt/adopt_test.go, projects/benchkit.md, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership
Covers: LQ4, LQ5, LQ6, LQ7

## What to build

Verify the premise first: `Link` in internal/adopt/link.go never calls
`kitSourceCheckout`. Then add the refusal directly after the line that reads
`kitDir()` for the plan, so the predicate and the plan share that one read.
Use the adopt `toon.Errorf` shape and exit 1. The title names the kit source
checkout, and the detail names `bench doctor --fix` as the kit-side route. Add
one sentence under the profile's cold-session notes that states the rule.

## Acceptance

- [ ] A new adopt test with `BENCH_KIT` at the repository root receives exit 1 and the kit-checkout title.
- [ ] `TestLinkOutsideGitRepoNamesGitRepository` still receives the git-repository message.
- [ ] `TestWrapperInstallFreshnessAndReloadJourneys` and `TestStrippedDistributionJourney` pass.
- [ ] The cold-session notes in `projects/benchkit.md` hold the rule sentence.
- [ ] Self-probe: replace the predicate with a `.bench/BENCH.md` marker test, and report the relink journey red.
