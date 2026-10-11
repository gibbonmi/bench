# Hold one landing lease per destination from composition through publish

Blocked by: none
Writes: .bench/commitment.json, DATA_HANDLING.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_commitment.go, internal/commitment/commitmenttest/admission.go, internal/commitment/repository/admission.go, internal/commitment/repository/lease.go (new), internal/commitment/repository/lease_test.go (new), internal/commitment/repository/publication.go, internal/commitment/repository/publication_lock_test.go, internal/commitment/repository/publication_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/commitment_record_test.go, internal/conformance/data_handling_test.go, internal/conformance/ordinary_build_census_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/run_transaction.go, internal/gocache/lock.go, internal/intent/build_records_test.go (new), internal/intent/ledger/commitment.go, internal/intent/ledger/validate.go, internal/intent/ledger_aliases.go, internal/intent/ledger_aliases_test.go, internal/intent/transaction.go, internal/refusalroute/command_test.go, internal/refusalroute/faces_land.go, internal/refusalroute/registry.go, internal/refusalroute/registry_test.go, internal/systemtest/owner_land_lease_test.go (new), internal/systemtest/owner_land_race_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/systemtest/owner_test.go, internal/worktree/identity_component_test.go, internal/worktree/joins.go, internal/worktree/land.go, internal/worktree/land_lease.go (new), internal/worktree/land_lease_test.go (new), internal/worktree/land_refusal.go, internal/worktree/land_refusal_fixture_test.go, internal/worktree/land_surface_test.go, internal/worktree/parallel_census_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/refusal_route_test.go, tests/canary/data-handling-derivation/undocumented-passlist-var
Covers: BW1, BW2, BW3, BW4, BW5, BW6, BW7, BW8, BW9, BW10, BW11, BW12, BW13, BW14, BW15, BW16

## What to build

A landing holds one lease on its destination ref in the intent ledger from composition through publish.
A second landing on the same destination refuses before composition and names the owner's assignment id.
The refusal routes the operator to wait for that landing and then rerun.

First add the three record lists to `CommitmentState` under the spec's Records decision: `leases`, `gaps`, and `overrides`, each `omitempty`.
`ValidateCommitment` refuses each of these records:

- a second lease on one destination
- a second gap for one outcome
- a second override for one assignment
- an invalid assignment identity
- an empty reason or a reason with a control byte

The tolerant read keeps the three lists, `internal/intent/ledger_aliases.go` re-exports the three types, and `DATA_HANDLING.md` lists each leaf field between the commitment record markers.
This ticket adds the gap and override types only as validated records; tickets 5, 6, and 7 add their producers.

Put the lease owner in `internal/commitment/repository/lease.go`, beside `PublishAdmitted`.
The acquire runs after the admission proofs and before `j.landReviewed`, in one `intent.Transact` decision with a compensation.
Liveness is a kernel record lock on `<git-common-dir>/bench-landing-<first 16 hex of sha256(destination ref)>.lock`, at mode 0600, never removed, and opened without following a link.
An in-process holder map answers for the calling process, because a POSIX record lock belongs to its process.
Reuse one exported record-lock rule from `internal/gocache/lock.go`; add no fourth kill-0 probe and no second copy of the holder map.
A lock with no holder means a dead owner, and the new record replaces the old one.

The release has two sites.
`PublishAdmitted` removes the lease record in the same decision that publishes.
The pre-publish error branch at `internal/worktree/land.go:205-211` removes the record in its own transaction.
Each site then unlocks and closes the lock file.
`--resume` takes no lease.

`internal/worktree/land.go` has about 9 lines of headroom, so put the landing-side wiring in `internal/worktree/land_lease.go`.
Add the face `land-lease-held` with agent authority to `internal/refusalroute/faces_land.go`.
Its detail is `landing destination refs/heads/<branch> is leased by assignment <id>`, and its route is the instruction `wait until the landing of assignment <id> ends`, then the caller's rerun.
A lock file that cannot be opened or locked refuses with `landing lease unavailable: <reason>` under `land-infrastructure`.
Add a producing fixture for the new face to `internal/worktree/land_refusal_fixture_test.go`.

BW10 runs two real landing processes in `internal/systemtest/owner_land_lease_test.go`, a system-tagged file.
Run it with `BENCH_KIT` set through `bench test --check system --run TestTwoLandingProcessesShareOneLease`, which supplies the kit root.
`internal/worktree/parallel_census_test.go` pins the package's exact top-level test count; raise `worktreeTestCount` in place for each new top-level test, and add no line to that file.
`.bench/commitment.json` is on this `Writes:` line only for the spec-stage binding commit, and this ticket does not edit it.

## Acceptance

- [ ] A landing that reaches its gate holds a lease record for `refs/heads/main` with its own assignment id.
- [ ] A second landing during the first landing's gate refuses before composition with the detail `landing destination refs/heads/main is leased by assignment <first-id>`.
- [ ] That refusal's `next=` begins with `wait until the landing of assignment <first-id> ends` and ends with the caller's rerun.
- [ ] The second landing runs no gate, and its output holds no `bench doctor` and no `gate owner: pid`.
- [ ] After a published landing, the ledger holds no lease, and a competing acquire inside the publish decision sees no lease.
- [ ] After a landing refuses red at its gate, the next landing acquires the lease.
- [ ] A seeded lease for another assignment with no lock holder is replaced by the next landing's record.
- [ ] Two landing processes contend, the loser refuses with the lease detail, and the winner publishes.
- [ ] A `--resume` of a published landing completes while another landing holds the lease.
- [ ] `ValidateCommitment` refuses two leases on one destination.
- [ ] A tolerant transaction over a ledger with a lease, a gap, and an override writes all three back.
- [ ] The commitment record inventory names `commitment.leases.destination`, `commitment.gaps.chunk`, and `commitment.overrides.reason`.
- [ ] A directory at the landing lock path refuses with `landing lease unavailable:` before composition.
- [ ] The recovery matrix and the landing face fixtures hold `land-lease-held`.
