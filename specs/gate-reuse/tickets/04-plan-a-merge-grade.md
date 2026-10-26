# Plan a merge grade before the merge runs

Blocked by: 03-record-a-reused-green.md
Writes: internal/landing/merge.go, internal/gate/engine.go, internal/worktree/merge.go, internal/worktree/merge_plan.go (new), internal/worktree/merge_refusal.go, internal/worktree/merge_plan_test.go (new), internal/worktree/parallel_census_test.go, internal/usage/worktree.go, cmd/bench/main.go, cmd/bench/help_inventory_test.go
Covers: GR31, GR32, GR33, GR34, GR35, GR36, GR37, GR38, GR39, GR40, GR46, GR41, GR42

## What to build

`usage.WorktreeMerge` becomes `bench worktree merge [--plan] --from <commit|target> <target>`. The merge grammar takes `--plan` as a flag with no value.

The merge help row in `cmd/bench/main.go` derives its suffix through `worktreeSuffix(usage.WorktreeMerge)`, as its sibling rows do. `cmd/bench/main.go` holds 429 lines over its limit, so the change adds no line. `TestHelpInventoryIsComplete` pins the new help line.

The landing owner's merge separates its kind decision and its composition from its authorization, so the plan and the merge read one decision. Keep the signature of `Owner.Merge`, because nine tests call it.
The plan reads the kind and the composed tree with no authorization and no ref move.

The gate exports the reuse predicate for a prospective tree. It is the `reusableEvidence` decision over the prospective subject of the tree at a root, with its cache-state refusal. The merge's whole-gate authority and the plan both read that one predicate.
A plan that read `InspectTreeContext` alone would project `reuse` beside a pending cache, and GR46 catches it.

The grade takes this order:

- `none` for a `current` kind, whatever lane the target declares.
- `lane` for a target that declares a lane.
- `reuse` for a target with no lane when the predicate holds for the composed tree.
- `whole` for every other target with no lane.

The plan prints `merge_plan{worktree=<id>,from=<commit>,kind=<kind>,previous_tip=<commit>,tree=<tree>,grade=<grade>}` at exit 0. A refusal takes the merge's face for the same cause, and the rerun from `mergeRerun` keeps `--plan` for a plan call.
The plan moves no ref, resets no checkout, writes no gate cache or evidence record, and takes no gate lock. Its predicate opens a prospective artifact bundle, which registers a private checkout and removes that registration on close.

`internal/worktree/merge.go` holds 380 of 400 lines, so the plan path goes in the new `merge_plan.go`.
The new tests use the merge sets of `merge_target_lane_test.go`, such as `goLaneMergeFixture` and `requireSelectiveLanePass`, and `wholeGateMergeFixture` of `refusal_route_follow_test.go`.
Each new test calls `t.Parallel()`. Raise `worktreeTestCount` in `parallel_census_test.go` by the count of new tests, with no new line.

## Acceptance

- [ ] The plan of a current target that declares a lane prints `grade=none` (GR31).
- [ ] The plan of a non-current target with a declared lane prints `grade=lane` (GR32).
- [ ] The plan of a target with no lane and a fresh green of the composed tree prints `grade=reuse` (GR33).
- [ ] The plan of a target with no lane and no such green prints `grade=whole` (GR34).
- [ ] The plan prints exactly one `merge_plan{...}` line with the six cells at exit 0 (GR35).
- [ ] The plan leaves the target branch tip, the checkout fingerprint, the `git worktree list` registrations, the gate cache bytes, and the evidence store names unchanged (GR36).
- [ ] After each planned grade, the merge runs that grade: a lane pass line for `lane` (GR37), the reuse line and no tally for `reuse` (GR38), one tally for `whole` (GR39), and the `current` kind for `none` (GR40).
- [ ] A no-lane target whose gate cache holds an interrupted pending record, with a fresh green of the composed tree, plans `grade=whole`, and the merge adds one tally (GR46).
- [ ] A plan of a conflicted composition refuses with the `merge-conflict` face, and its `next=` route ends with a rerun that holds `--plan` (GR41).
- [ ] `bench help` prints the merge row with `[--plan]` from the one merge grammar (GR42).
