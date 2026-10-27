# Refuse a stack member before composition

Blocked by: 01-share-the-stack-owner-rule.md, 03-admit-a-stack-publication.md, 09-attribute-a-red-stack.md
Writes: internal/refusalroute/faces_land.go, internal/worktree/land_refusal.go, internal/worktree/land_stack.go (new), internal/worktree/land_stack_rerun.go (new), internal/worktree/land_stack_preflight.go (new), internal/worktree/land_stack_refusal_test.go (new), internal/worktree/refusal_route_follow_test.go, internal/worktree/land_stack_route_test.go (new), internal/worktree/parallel_census_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_fixtures_test.go, internal/worktree/land_stack_test.go
Covers: LS70, LS71, LS72, LS73, LS74, LS75, LS76, LS77, LS88, LS90, LS91, LS92, LS93

## What to build

Add the stack proofs of the spec section "The preflight" in the new file `land_stack_preflight.go`. The stack flow of ticket 07 calls them after the source proofs of each member and before the lease. Every refusal of these proofs prints in the same preflight, and each refusal names its member index and label.

- A member whose `spec` names a staged spec refuses with `stack-member-alone`.
- A member whose assignment holds a delivery binding refuses with `stack-member-alone`. The proof calls `HoldsDeliveryBinding` of ticket 03, so a started tickets-only deliverable counts.
- A member whose tip does not contain the destination tip refuses with `land-destination-moved` and its fold route.
- A destination that moves after the grade (`landing.TipMovedError`) also refuses with the stack-raised `land-destination-moved`, as one landing does, and not with `land-handback`.
- After the member proofs, the proof calls `stackpolicy.Owners` and `stackpolicy.Overlaps` of ticket 01 over the destination tree and each member tree. It reads each tree with one `git ls-tree` entry format. Each overlap refuses with `stack-members-overlap`, which names both labels and the paths.
- The build window then decides each member with `AdmitWindow` before the lease, as one landing does. A member that the window refuses refuses with `land-build-window` and its index.

Register two land faces in `internal/refusalroute/faces_land.go`, with agent authority:

- `stack-member-alone` has one instruction step. It removes the member from the stack file and routes it to its own landing. The rerun follows.
- `stack-members-overlap` has one instruction step. It removes the later member of the pair from the stack file. The rerun follows.

The rerun composer of ticket 07 drops the override of the removed member. When one member remains, the rerun is that member's single landing.

Add the route-walk fixtures of these faces to the stack loop of ticket 09, at three members and at two members. Also add stack fixtures for a stack-raised `land-destination-moved` and `source-tip-mismatch`, each under its own cause key.

## Acceptance

- [ ] A member whose tip lacks the destination tip refuses with `land-destination-moved` and the label of that member, and the tally reads 0 (LS70).
- [ ] Two members that change one path outside `capture/` refuse with `stack-members-overlap`, which names both labels and the path, and the tally reads 0 (LS71).
- [ ] A member whose `spec` names a staged spec refuses with `stack-member-alone`, and the tally reads 0 (LS72).
- [ ] A member whose assignment holds the delivery binding of a started tickets-only deliverable refuses with `stack-member-alone`, and the tally reads 0 (LS73).
- [ ] A spec-less member whose assignment an open build binds refuses with `stack-member-alone`, and the tally reads 0 (LS88).
- [ ] With an open build and no chunk gap, a light-path member refuses with `land-build-window` and its index, and the ledger holds no lease (LS74).
- [ ] With an open build and an override for member 1 alone, the stack refuses member 2 with `land-build-window` (LS75).
- [ ] In a three-member stack, the carried-out `stack-member-alone` route ends in a stack landing that exits 0 (LS76).
- [ ] In a two-member stack, the `stack-member-alone` rerun is the single landing of the other member, and the carried-out route exits 0 (LS90).
- [ ] In a three-member stack, the carried-out `stack-members-overlap` route ends in a stack landing that exits 0 (LS77).
- [ ] In a two-member stack, the `stack-members-overlap` rerun is the single landing of the earlier member, and the carried-out route exits 0 (LS91).
- [ ] A stack-raised `land-destination-moved` refusal, carried out step by step, ends in a stack landing that exits 0 (LS92).
- [ ] A stack-raised `source-tip-mismatch` refusal, carried out step by step, ends in a stack landing that exits 0 (LS93).

The bound-member fixtures follow `TestOpenBuildsDeriveFromActiveOwners` in `internal/commitment/repository/window_test.go`. The window fixture follows `commitmenttest.OpenChunkGaps` and its closed-gap state.

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench test --package ./internal/refusalroute`
- `bench test --check root-conformance`

## Mutation proof

Make the binding proof test only the open builds. The LS73 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
