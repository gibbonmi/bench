# Land the spec and its binding together

Blocked by: 1-bind-a-new-deliverable-from-its-publishing-tree.md
Writes: CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/commitment/commitmenttest/repo.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/land_specless_test.go, internal/worktree/parallel_census_test.go, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: SB17, SB18, SB19

## What to build

One spec-less `bench worktree land` of a planning worktree publishes a staged spec and its deliverable binding together.
After that landing, `bench commitment show` in the primary checkout prints the start command with the spec path.
Before that landing, the command keeps `--deliverable <path>`, because the outlook reads only the default-branch policy.

Ticket 1 supplies the contract that the last paragraph of its What to build states.
This ticket reads that contract through the real commit and landing owners, and it changes no file of those owners.

Add the planned `TestCommitmentSpecStageLandsSpecAndBinding` to `internal/worktree/land_specless_test.go`.
The journey runs these steps:

1. Make the fixture with `planningLandingFixture` and an adopted policy.
2. In the planning worktree, write a staged spec with a `Roadmap:` line, and commit it through `benchCommit`.
3. Add one deliverable binding of that spec to the outcome that `bench commitment show` names, with the identity from `commitment.Identity` over the spec bytes.
4. Run `bench commitment plan` and `bench commitment approve` in the planning worktree.
5. Commit `.bench/commitment.json` and each other path that approve staged, through `benchCommit`.
6. Run `bench commitment show` in the primary checkout and in the planning worktree.
7. Land the planning worktree spec-less once.
8. Run `bench commitment show` in the primary checkout again.

`TestCommitmentPlanningBootstrap` is the precedent for the steps that commit and land.
The journey reads each spec path from its fixture writes and each outcome id from the fixture policy.
No expectation calls the code under test.

These checks prove the rows:

- SB17: the landing exits zero and releases the worktree, and main holds the spec bytes and the bound policy.
- SB18: after the landing, the primary checkout prints `--deliverable specs/<slug>/spec.md` in the start command.
- SB19: before the landing, the primary checkout prints `--deliverable <path>` in the start command.

The journey also requires `--deliverable <path>` from the planning worktree before the landing, because the outlook keeps the default-branch policy.
That check gives SB19 its red: the primary checkout holds main before the landing, so only the worktree read can show a worktree policy.

The journey calls `t.Parallel`, as every test of this file does, so `worktreeSerialCeiling` stays at 45.
Raise `worktreeTestCount` at `internal/worktree/parallel_census_test.go:1061` by one for the new top-level test.
That file holds 1096 lines, so the edit changes one line and adds none.
`land_specless_test.go` holds 300 lines, so keep it at 400 lines or fewer.
Put a new helper in `commitmenttest/repo.go` only when the journey needs a policy helper that the package does not supply.

Ticket 1 lands first, so the journey passes at once.
Record each red through `bench probe`:

- Run the central probe of ticket 1. SB17 and SB18 fail.
- Make `Store.Outlook` read the policy at the checkout HEAD in place of the default-branch policy. The SB19 check in the planning worktree fails.

Each probe reports a restored subject.
If the landing refuses for a reason that ticket 1 does not own, stop and report the refusal with its detail.
Do not edit a landing owner outside this `Writes:` line.

The `Writes:` line names `CHANGELOG.md` as an upper bound, because a research return named it for SB1.
This ticket does not write the `CHANGELOG.md` entry, because the spec gives the one entry to the last ticket.
The anchor registry files `internal/anchors/registry_data.go` and `internal/anchors/registry_data_test.go` name `CHANGELOG.md`, and two canary units pin it.
So these four files join the `Writes:` line, and this ticket does not change them.

The binding registry binds the five command registry files to `internal/commitment`, `internal/anchors`, and `internal/worktree`, so they join this `Writes:` line.
This ticket does not change them, and four of them hold more than 400 lines.

## Acceptance

- [ ] A spec-less `bench worktree land` of the planning worktree publishes the staged spec and its binding in one landing.
- [ ] After that landing, `bench commitment show` in the primary checkout prints the start command with `--deliverable specs/<slug>/spec.md`.
- [ ] Before that landing, `bench commitment show` in the primary checkout prints the start command with `--deliverable <path>`.
- [ ] Before that landing, `bench commitment show` in the planning worktree prints the start command with `--deliverable <path>`.
- [ ] `worktreeTestCount` equals the count of top-level tests in `internal/worktree`.
