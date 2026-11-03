# Run the commit chain before publication

Blocked by: 01-grade-the-candidate-commit.md, 02-check-before-publication.md
Writes: cmd/bench/commit_chain.go, cmd/bench/commit_chain_test.go, internal/systemtest/commit_chain_test.go (new), internal/anchors/registry_plan_commit.go (new), internal/anchors/registry_plan_commit_test.go (new), internal/anchors/registry_data.go, CHANGELOG.md, .bench/commitment.json, internal/commit/grammar.go, internal/commit/check_test.go
Covers: PP17, PP18, PP19, PP20, PP21, PP22, PP23, PP24, PP25, PP26, PP27, PP28, PP29, PP30, PP31

## What to build

`bench commit --preflight-build <slug> [--base <commit>]` runs its chain steps before publication. `commitChainCommand` calls `commit.RunChecked` (ticket 02) with a check when the flag names a slug.
The check runs the worktree build first. If the build is green, it runs `preflight.BuildVerdictAt` (ticket 01) with the root, the slug, the base, and the candidate. It prints each step's own response and records each step state.
A red build or a red preflight makes the check return an error, so the commit publishes nothing and exits 1.

The chain line keeps its fields and their order: `commit-chain{commit=<sha|none>,build=<state>,preflight=<state>}`. A lane refusal prints every step `skipped`. A publication that does not reconcile keeps exit 3, and the line names the commit with both steps green.
`publishedStepRedExit` and exit 4 retire. Rewrite `chainHelp` so that it states that the build and the preflight run before publication, at the candidate commit, and that a red step publishes nothing.
`chainSteps.Commit` takes the check, and `chainSteps.Preflight` takes the root, the slug, the base, and the candidate. Change the stub commit of `runCommitChain` so that it calls the check with a fixed candidate. Rewrite the inverse assertions of the chain tests that rows PP17 to PP21 and PP24 to PP26 cite, and retire the local `wantPublishedStepRedExit`.

Add `internal/systemtest/commit_chain_test.go` after the linked-project precedent of `internal/systemtest/otel_verbs_test.go:40`. A red plan commit exits 1 with the branch ref unchanged, and its repaired retry exits 0 with the chain line that names the new branch tip.
Run the system suite as `AGENTS.md` states, so that the suite gets its BENCH_KIT.

Add `internal/anchors/registry_plan_commit.go` with two Require anchors on `.agents/commands/bench-implement-spec.md`, and register the slice in the `registry` chain of `internal/anchors/registry_data.go:63`. The needles are "A plan commit that a dispatch follows uses plain `bench commit`" and "a plan commit that no dispatch follows uses `bench commit --preflight-build <slug>`". Add the test mirror after the precedent of `internal/anchors/registry_check_floor_test.go`. The guidance text does not change.

Add one `CHANGELOG.md` entry under `[Unreleased]`. It states the new step order, the exit-1 refusal, and the retired exit 4.
The spec-stage close writes the deliverable binding in `.bench/commitment.json`. This ticket names that path so that the plan-only preflight authorizes the binding commit.

## Acceptance

- [ ] PP17: a green chain prints the build, the preflight, and the commit lines in that order, then the chain line, at exit 0.
- [ ] PP18: the preflight step receives the stub candidate, not the published commit.
- [ ] PP19: a chain that names a base gives that base to the preflight step.
- [ ] PP20: a red build exits 1, calls no preflight, and prints `commit-chain{commit=none,build=red,preflight=skipped}`.
- [ ] PP21: a red preflight at exit 1, 2, or 3 makes the chain exit 1 and print `commit-chain{commit=none,build=green,preflight=red}`.
- [ ] PP22: a lane refusal before the check exits 1 and prints `commit-chain{commit=none,build=skipped,preflight=skipped}`.
- [ ] PP23: a green chain line names the published commit.
- [ ] PP24: a publication that does not reconcile exits 3 and names the commit with both steps green.
- [ ] PP25: the commit help holds no `exit 4` text.
- [ ] PP26: the chain help line starts with the `--preflight-build` group and states that the steps run before publication.
- [ ] PP27: `TestTicketStagesTheRetainedProbe` stays green.
- [ ] PP28: in a linked-project worktree, the real verb refuses a red plan commit with the branch ref unchanged, and it publishes the repaired retry.
- [ ] PP29: the anchor test reds when the dispatch half of the route sentence is removed.
- [ ] PP30: the anchor test reds when the chain half of the route sentence is removed.
- [ ] PP31: the changelog entry states that the chain steps run before publication and that exit 4 retired.
