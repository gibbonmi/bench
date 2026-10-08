# Start the next outcome at phase close

Blocked by: none
Writes: .bench/BENCH.md, .agents/commands/bench-drain.md, internal/anchors/registry_commitment.go, internal/conformance/commitment_guidance_test.go, CHANGELOG.md
Covers: none

## What to build

A phase close handed the reviewer the `commitment_outlook` command. The reviewer then had to run it, or approve each preparation step. The reviewer directed on 2026-10-08 that the agent never hands a `bench` command to the reviewer. The next eligible outcome in the approved order starts without a new approval, because the closing phase's approval covers the start.

The operating guide gains one rule. At a phase close, the agent does the preparation and runs `bench commitment start` from the outcome's own worktree. If the outcome has no bound deliverable and its row says `Next: ticket`, the agent writes and lands that ticket folder. It then binds the folder through `bench commitment plan` and `bench commitment approve`, and it cites the closing phase's approval.

That plan changes only the binding. A row that needs a spec routes to `/bench-write-spec`. A reorder or a new admission still needs the reviewer's direction.

The drain exit routes to that rule. The commitment anchor family pins the two new guide sentences and the drain sentence.

## Acceptance

- [ ] The operating guide states that the agent starts the next outcome at a phase close and hands no `bench` command to the reviewer.
- [ ] The operating guide limits the binding plan to a missing ticket deliverable.
- [ ] The drain exit tells the coordinator to start the outcome under that rule, and it no longer names the command as a recommendation.
- [ ] The commitment anchor family refuses the removal of each new sentence, and `TestCommitmentGuidance` records each red.
- [ ] `CHANGELOG.md` has one entry for the changed behavior.
