# Audit spec staleness before the build

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/
Covers: none

## What to build

A spec-backed build starts with a staleness pass before its first ticket dispatch. After the build preflight, the orchestrator charges a read-only delegate on the mid tier. The delegate audits the spec and its tickets against the current tree and guidance, and it returns its findings. The orchestrator commits the fixes as one spec amendment on the integration source. The preflight then reruns green, and one delegated review round confirms the amendment with no second round.

A red build preflight stops the phase only when a red row is outside the staleness class. The staleness class is each `*-closure` row, `fence-writes`, and `completion-plan`. The staleness pass takes a red in that class and fixes it in the amendment.

The pass corrects claims, citations, fences, landed work, and a preflight red that it took. It does not reopen the approach. The phase file points to the owners of each rule and restates none of them:

- `craft-spec` owns the current-code claim rules.
- `craft-line` owns the tiers and the conditional review line.
- `.bench/BENCH.md` owns the spec-contradiction predicate and the approved plan-expansion policy.

## Acceptance

- [ ] The "Declare the line, validate the tickets, route the venue" section charges a read-only mid-tier delegate for a staleness audit after the build preflight.
- [ ] The section commits the fixes as the enabling plan commit of `.bench/BENCH.md`'s approved plan-expansion policy, before the version 2 plan amendment.
- [ ] The section routes each finding by `.bench/BENCH.md`'s spec-contradiction predicate and does not reopen the approach.
- [ ] The section runs one delegated review round on the conditional review line, and the preflight reruns green before the version 2 plan amendment.
- [ ] A red preflight whose red rows are all `*-closure` rows, `fence-writes`, or `completion-plan` goes to the staleness pass, and any other red stops the phase.
- [ ] The step names no model and restates no staleness rule.
- [ ] The anchor registry pins the preflight route with its own needle, and the `implement-spec-red-preflight-route` canary reds when a mutation widens the route to any red.
- [ ] The `implement-spec-worktree-before-preflight` canary pins the new preflight sentence and still reds on the order swap.
- [ ] Every other pinned needle in the phase file stays intact, and `.agents/commands/bench-implement-spec.md` stays inside its prose budget.
