# Bind a new deliverable from its publishing tree

Blocked by: none
Writes: .bench/commitment.json, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/commitment/commitcmd/command.go, internal/commitment/commitmenttest/repo.go, internal/commitment/delivery.go, internal/commitment/head_binding_test.go (new), internal/commitment/repository/admission.go, internal/commitment/repository/candidate.go, internal/commitment/repository/closure.go, internal/commitment/repository/head_binding_test.go (new), internal/commitment/repository/publication.go, internal/commitment/repository/readiness.go, internal/commitment/repository/repository.go, internal/commitment/repository/sources.go, internal/commitment/store_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: SB01, SB02, SB03, SB04, SB05, SB06, SB07, SB08, SB09, SB10, SB11, SB12, SB13, SB14, SB15, SB16, SB20, SB21, SB22

## What to build

A spec author binds a staged spec that the worktree HEAD holds and main does not hold.
`bench commitment plan` and `bench commitment approve` read each new deliverable source at the checkout HEAD.
`bench commit` and `bench worktree land` read each new deliverable source at the candidate tree.
Every other source keeps the default-branch tip, and start, readiness, and the outlook keep the default-branch policy.
So a binding counts for start, readiness, and status only after it lands.

The spec terms apply here.
A deliverable of the proposed policy names a new deliverable source.
No deliverable of the current default-branch policy names that id, path, and identity.
The checkout HEAD is the commit that `git rev-parse HEAD` names in the store root.
The candidate tree is the tree that `authorizeCandidate` grades: the composed candidate of `bench commit`, or the composed tree of a landing.

Add one owner that decides which tree each source reads, so plan, approve, and the approval check cannot drift.
The owner receives the current policy, the source, the publishing tree, and the default-branch revision.
It returns the publishing tree for a new deliverable source and the default-branch revision for every other source.
Put the owner in `internal/commitment/repository/sources.go` or `internal/commitment/repository/candidate.go`.
If the owner needs a policy rule of `internal/commitment`, put that rule beside `Unsettled` in `internal/commitment/delivery.go`.

Change these four sites through the owner:

- `Store.Plan` validates each unsettled deliverable of the proposed policy and each plan source, with the checkout HEAD as the publishing tree.
- `Store.Approve` and its replay validate each plan source, with the checkout HEAD as the publishing tree.
- `Store.Inventory` lists each staged spec that the checkout HEAD holds, reads its identity at that HEAD, and keeps the `spec.Facts` join for its scope.
- `approvedTransition` receives the candidate tree from `authorizeCandidate`, and validates each source of the expected plan with the candidate tree as the publishing tree.

The inventory roadmap rows keep the default-branch tip.
The read at a publishing tree reuses `sourceAt`, the staged check of `validateDeliverable`, and the tickets-only rule of `ticketsOnlyAt` at that tree.
`validateDeliverable` keeps its default-tip behavior for `Store.Start` (`admission.go:41`) and `Store.readyFor` (`readiness.go:60`).
A new deliverable source that the publishing tree does not hold refuses with `commitment source "<id>" refused:`, also when main holds it.
The receipt match, the plan identity, and the refusal `candidate policy has no exact approval` stay unchanged.
The commitment grammar, its help rows, its refusal faces, and `internal/commit/commit.go` stay unchanged.

`internal/commitment/repository/repository.go` holds 400 of 400 lines, so keep it at 400 lines or fewer.
Each edit there removes as many lines as it adds.
`sources.go` holds 106 lines, and `candidate.go` holds 292 lines.
Keep every file that this ticket writes at 400 lines or fewer.
The lane runs `bench structure --growth`, and a file over 400 lines that gains a line turns it red.

Put SB01-SB11 and SB21 in `internal/commitment/head_binding_test.go` (new), in the planned tests that the spec coverage map names.
Put SB12-SB15 in `internal/commitment/repository/head_binding_test.go` (new), in the planned tests that the spec coverage map names.
Each test commits a staged spec in a linked worktree from `commitmenttest.Assignment`, and runs plan and approve there.
SB12 and SB13 grade the worktree candidate, and SB14 and SB15 grade a composed tree that the test builds.
Each identity comes from `commitment.Identity` over the fixture bytes, and each path comes from the fixture writes.
No test swaps a package variable, and no expectation calls the code under test.

SB16, SB20, and SB22 are current pins, and each one stays green with no edit:

- SB16 is `TestCommitmentProposalIdentity` in `internal/commit/commitment_test.go`.
- SB20 is `TestCommitmentDeliverableTypes` in `internal/commitment/deliverable_test.go`.
- SB22 is `TestCommitmentNextProjection` in `cmd/bench/commitment_test.go`.

Before the first production edit, run the three commitment packages.
Record the SB01 red at `Store.Plan`, and list each other red as a posture change.
SB01, SB02, SB03, SB04, SB05, SB06, SB08, and SB09 fail before the edit, because the current code reads each source at the default-branch tip.
SB09 asserts the sentence `deliverable "<path>" is not staged`, and the current code gives `commitment source "<id>" refused:` there.

After the edit, record one red through `bench probe` for each row that passes before the edit, and for the central probe:

- Central probe: swap the candidate tree for the default-branch revision in the new-deliverable branch of `approvedTransition`. SB12 and SB14 fail.
- Swap the candidate tree for the checkout HEAD in the same branch. SB13 fails.
- Omit the identity check of a new deliverable source in `approvedTransition`. SB13 and SB15 fail.
- Make the owner return the publishing tree for every source. SB10 and SB11 fail.
- Read the new deliverable source from the working tree of the store root in place of the checkout HEAD. SB07 fails.
- Omit the receipt match in `approvedTransition`, so any approved receipt passes. SB16 fails.
- Make `validateDeliverable` refuse each deliverable at the default-branch tip. SB20 fails.
- Make `Store.Start` read the policy at the checkout HEAD in place of the default-branch policy. SB21 fails.
- Fill the `--request <request>` placeholder of `withCommand` with a fixed value. SB22 fails.

Each probe reports a restored subject.

The `Writes:` line names `.bench/commitment.json` under decision D5, and no other ticket of this spec names it.
The spec-stage close commits the binding inside the range that the plan-only preflight grades.
This ticket does not change the policy file.

The binding registry binds the five command registry files to `internal/commitment`, so they join this `Writes:` line.
This ticket does not change them, and four of them hold more than 400 lines.
`admission.go`, `readiness.go`, `closure.go`, `publication.go`, and `commitcmd/command.go` call a function that this ticket can change.
Each one changes only when a signature changes.
`commitmenttest/repo.go` and `delivery.go` are possible helper sites.
The comment at `store_test.go:146` says that the default branch does not hold a spec, and it can name the checkout HEAD instead.

This ticket gives ticket 2 this contract.
In a linked planning worktree, plan, approve, and `bench commit` accept a binding of a staged spec that HEAD holds and main does not hold.
`AdmitPublication` admits the composed tree that holds the spec and its approved binding.
Start, readiness, and the outlook read only the default-branch policy.

## Acceptance

- [ ] In a linked worktree whose HEAD holds a staged spec that main lacks, `Store.Plan` accepts a policy that binds that spec.
- [ ] `Store.Approve` of that plan writes the bound policy into the worktree's `.bench/commitment.json`.
- [ ] In that worktree, `Store.Inventory` lists a `deliverable` row for the HEAD-held spec.
- [ ] The inventory row of the HEAD-held spec carries the spec's `Roadmap:` value as its scope.
- [ ] The inventory identity of the HEAD-held spec equals the identity that `Store.Plan` accepts for that spec.
- [ ] In a linked worktree whose HEAD holds a tickets-only folder that main lacks, `Store.Plan` accepts a policy that binds that folder.
- [ ] `Store.Plan` in the worktree refuses a staged spec that only the working tree holds.
- [ ] `Store.Plan` in a worktree whose HEAD lacks a spec that main holds refuses a new binding of that spec.
- [ ] `Store.Plan` in the worktree refuses a HEAD-held spec with no `Status: staged` line, with the sentence `deliverable "<path>" is not staged`.
- [ ] `Store.Plan` in a worktree that commits a changed roadmap row refuses a plan that names the changed row identity.
- [ ] `Store.Plan` in a worktree whose HEAD predates another outcome's binding on main accepts a policy that keeps that binding.
- [ ] In the worktree, `Store.AuthorizeCandidate` admits a candidate tree that holds the staged spec and its approved binding.
- [ ] `Store.AuthorizeCandidate` refuses a candidate tree whose spec bytes differ from the approved identity while HEAD keeps the approved bytes.
- [ ] `Store.AdmitPublication` admits a spec-less publication whose composed tree holds the staged spec and its approved binding.
- [ ] `Store.AdmitPublication` refuses a composed tree whose spec bytes differ from the approved identity.
- [ ] `bench commit` refuses a policy candidate whose transition has no exact approved plan.
- [ ] `bench commitment start` admits a staged spec that the default-branch policy binds.
- [ ] Before the landing, `Store.Start` from the planning worktree refuses the HEAD-bound spec with `is not approved for outcome`.
- [ ] The bound start command keeps the `--request <request>` placeholder.
- [ ] `internal/commitment/repository/repository.go` holds 400 lines or fewer.
