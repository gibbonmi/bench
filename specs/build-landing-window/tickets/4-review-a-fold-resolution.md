# Review a fold resolution and check the merged test packages

Blocked by: 3-record-a-fold.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/review_checkpoint_commits_test.go, internal/preflight/decision.go, internal/preflight/folded_main_test.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/preflighttest/fixture.go, internal/preflight/paths_authorized.go, internal/preflight/review.go, internal/refusalroute/faces_merge.go, internal/refusalroute/registry.go, internal/reviewrecord/check.go, internal/reviewrecord/coverage.go, internal/reviewrecord/delegated.go, internal/reviewrecord/fold_test.go (new), internal/reviewrecord/parse.go, internal/reviewrecord/record.go, internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/refusal.go, internal/reviewrecord/recordcmd/review_test.go, internal/reviewrecord/recordcmd/verification.go, internal/reviewrecord/recordcmd/verification_test.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/recordtest/fold.go (new), internal/reviewrecord/repair.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/write.go, internal/worktree/land_refusal.go, internal/worktree/merge.go, internal/worktree/merge_fold_red_test.go (new), internal/worktree/merge_refusal.go, internal/worktree/merge_route_test.go, internal/worktree/merge_target_grade_test.go, internal/worktree/parallel_census_test.go
Covers: BW39, BW40, BW42, BW43, BW44, BW45, BW93

## What to build

A fold whose resolution is not empty needs three independent reviews, and the final checks include each package of a merged test file.
A red fold routes to a hand merge, `bench commit`, and the fold record, so the repair hunks enter the resolution.

`bench record review` takes `--fold <id>` in place of `--chunk <id>`.
A fold review copies the fold's `From` and `Commit` as its base and tip.
The three axes follow the chunk rules of `checkReviews`, so a fold with two passing axes still refuses the checkpoint.
A fold with an empty resolution needs no review.
Ticket 3 supplies the fold entry, its `Resolution`, and its `TestPackages`.

`checkCompletion` (`internal/reviewrecord/check.go:71-84`) adds one derived requirement for each package in any fold's `TestPackages`.
The requirement id is `merged-test:<package>`, and its command is `bench test --package <package>`.
`bench record verification --final` accepts these ids.
The completion requires a passing result for each derived requirement at the completion source, as it requires for the planned final list.

`unauthorizedPaths` (`internal/preflight/paths_authorized.go:24`) also skips a path that a fold entry of the record at the source tip lists in `Resolution`.
A path that no fold entry lists stays red.
`internal/preflight/gather.go` is over its budget, so read the record in `internal/preflight/gather_inputs.go` and add no net line to `gather.go`.

The `merge-fold-red` face keeps reviewer authority.
Its route becomes four steps: the existing first instruction, the hand merge `<checkout-git> merge <from-commit> (repair the red in the merge commit)`, `bench commit`, and `bench record fold`.
The route and the face file no longer name FT342.
Reuse `conflictFacts` (`internal/worktree/land_refusal.go:68-74`) for the checkout Git value.
`internal/worktree/merge_route_test.go` has about 10 lines of headroom, so move `foldRedFixture` to `internal/worktree/merge_fold_red_test.go` if the new route needs room.
Raise `worktreeTestCount` in `internal/worktree/parallel_census_test.go` in place for each new top-level test.

## Acceptance

- [ ] A fold with a non-empty resolution and two passing axes refuses the checkpoint until the third axis passes.
- [ ] A fold with an empty resolution and no review passes the checkpoint.
- [ ] The complete checkpoint refuses until a passing final result exists for `merged-test:<package>`.
- [ ] `paths-authorized` passes for a resolution path that the record's fold entry lists and the fence omits.
- [ ] `paths-authorized` stays red for the same path when no fold entry lists it.
- [ ] The `merge-fold-red` route names the hand merge, `bench commit`, and `bench record fold`, and it names no `FT342`.
- [ ] The help inventory lists the `bench record fold` form, and the review form names `--fold`.
