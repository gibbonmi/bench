# State the chunk chain and the record order

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, internal/anchors/registry_chunk_chain.go (new), internal/anchors/registry_chunk_chain_test.go (new), internal/anchors/registry_data.go
Covers: none

## What to build

This ticket is the FT319 kit edit. The implement and review phase files left the chunk
chain partly unstated. The bounded-charge-evidence and slicing-closure builds paid for
each gap with a refused checkpoint or an added review round.

The implement phase states these rules:

- Build preflight runs through the worktree after each ticket commit.
- The author commits the verification and probe record before the axis dispatch.
- Plan commits and `main` merges land before the ticket merge, and only record commits follow the chunk tip.
- The reconciliation commit joins the review delta of the last chunk.
- The ordinary assessment record comes before the landing step.

The review phase states these rules:

- A later plan commit is never a chunk base.
- A chunk that ends on a repair takes one confirming round of all three axes at its final tip.
- On one shared tree, only one axis runs tests or probes.
- A review worktree moves to the record commit, and the frozen pair still names the source tip.

The bounded repair policy states the hardening cap. The reviewer set the cap at one
cycle on 2026-09-18. Each new sentence gets one anchor row and one independent test
expectation.

## Acceptance

- [ ] The implement phase states the five rules above, and each rule has an anchor row.
- [ ] The review phase states the four rules above, and each rule has an anchor row.
- [ ] The bounded repair policy states the cap of one hardening cycle, and the cap has an anchor row.
- [ ] If you delete the hardening-cap sentence, `bench test --check docs-currency-workflow` goes red.
- [ ] If you delete the shared-tree sentence, `bench test --check docs-currency-workflow` goes red.
- [ ] The prose lane and the gate pass.
