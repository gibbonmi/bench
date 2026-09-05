# Construct the lane owner once

Blocked by: none
Writes: internal/landing/landing.go, internal/landing/lane_test.go, internal/commit/commit.go, internal/worktree/merge.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SR15, SR16, SR17, SR18, SR19

## What to build

The line is opus / low. Give `internal/landing` one constructor that takes a
resolved lane and a base and answers an owner. A nil lane answers the gate
owner. A lane answers the lane owner with the lane's Checks, Kit, and
Selective values and the caller's base. A pure function owns the mapping, and
the constructor composes it.

`bench commit` and `bench worktree merge` each call the constructor once.
Neither verb builds the lane authority by hand any more. The two accepted-kind
lists stay as they are, so `Owner.Land` and `Owner.LandReviewed` keep their
kinds and the reviewed landing keeps refusing a lane pass.

The exit proof for this ticket is the pre-existing suite, green with its test
logic unchanged. A mechanical rename is permitted. A needed assertion change
stops the ticket and reports.

## Acceptance

- [ ] The pure mapping from a lane and a base answers a lane authority with the lane's Checks, Kit, and Selective values and that base.
- [ ] `bench commit` under a declared lane runs the lane and publishes.
- [ ] `bench worktree merge` runs the declared lane on the composed tree.
- [ ] `bench consumers landing.NewLane` lists one caller, inside `internal/landing`.
- [ ] A root with no lane keeps the gate commit, and a merge with no lane publishes the merge tree.
- [ ] The reviewed landing refuses a lane pass and names its kind.
- [ ] Self-probe: drop the Selective value in the mapping, and report the observed red.
