# Attribute a red stack to one member with stack-red

Blocked by: 07-land-a-green-stack.md, 08-bisect-a-red-stack.md
Writes: internal/refusalroute/faces_land.go, internal/worktree/land_refusal.go, internal/worktree/land_stack.go (new), internal/worktree/land_stack_rerun.go (new), internal/worktree/land_stack_red_test.go (new), internal/worktree/refusal_route_follow_test.go, internal/worktree/land_stack_route_test.go (new), internal/worktree/parallel_census_test.go, internal/worktree/identity_component_test.go
Covers: LS46, LS47, LS48, LS49, LS68, LS69, LS89, LS102, LS103, LS106

## What to build

Register the land face `stack-red` in `internal/refusalroute/faces_land.go`, with agent authority. Its route has four steps:

1. Repair each failure that the gate reports in the attributed member's label.
2. Commit at that label.
3. Run the review.
4. One instruction step. When the members table shows a landed member, remove each landed member from the stack file, fold the destination into each remaining member, and set the new `source_tip` of each folded member in the stack file.

The rerun comes last. The operator fills its `<repaired-tip>` slot after the instruction step, because the fold moves the tip of the attributed member a second time.

The stack flow of ticket 07 replaces its interim `land-handback` for an attributed red with `stack-red`. The face names the attributed member's label. An unattributed `inherited` stop of ticket 08 keeps `land-handback`, because no member can carry the repair.

The red journeys first record a green of the destination, so the top red reads `candidate` and the bisect runs. The inherited row uses a destination with no green. The rerun composer of ticket 07 counts `stack-red` as a face that repairs the source, so the rerun adds the slot of the attributed member. When a prefix published, the composer drops the override of each landed member and of each folded member. When one member remains, the rerun is that member's single landing.

Add the stack fixtures to the route walk. `TestLandingFacesFollowTheirRoutes` gains a loop over the stack fixtures in the new `land_stack_route_test.go`:

- Each stack fixture takes its own cause key, so `routetest.Fixtures` accepts a second fixture of one face.
- A stack rerun ends with the `stack{...}` record, so the stack walk checks that record. A single rerun keeps the `landed{` check.
- One fixture puts the stack file under a directory whose name holds a space.

Ticket 10 adds its fixtures to the same loop.

## Acceptance

- [ ] With three members and member 2 red, the stack publishes member 1 and exits 1 (LS46).
- [ ] With three members and member 2 red, the gate tally reads 3, which is 1 + ⌈log2 3⌉ (LS47).
- [ ] With three members and member 2 red, member 3 stays active and its `outcome` reads `unlanded` (LS48).
- [ ] The `next=` route of the `stack-red` refusal names the label of member 2 in its repair step (LS49).
- [ ] In a three-member stack with member 2 red, the carried-out `stack-red` route ends in a stack landing that exits 0 (LS68).
- [ ] A stack file whose path holds a space reruns from the printed `next=` route (LS69).
- [ ] In a two-member stack with member 2 red, the rerun is the single landing of member 2, and the carried-out route exits 0 (LS89).
- [ ] In a two-member stack with member 1 red, the rerun adds `--source-tip <request>=<repaired-tip>` for member 1, and the carried-out route exits 0 (LS102).
- [ ] A three-member stack called with an override for member 1, with member 2 red, publishes member 1, and its carried-out `stack-red` route exits 0 (LS103).
- [ ] A stack whose top tree reads `inherited` refuses with `land-handback` after one grade, and the gate tally reads 1 (LS106).

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench test --package ./internal/refusalroute`
- `bench test --check root-conformance`

## Mutation proof

Leave out the slot of the attributed member when nothing published. The LS102 walk must fail at `source-tip-mismatch`. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
