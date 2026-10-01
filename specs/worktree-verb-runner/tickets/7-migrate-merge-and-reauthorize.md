# Move the merge and reauthorize verbs onto the verb runner

Blocked by: 6-migrate-query-and-create-verbs.md
Writes: internal/worktree/verb_fixture_test.go, internal/worktree/merge_test.go, internal/worktree/merge_caller_root_test.go, internal/worktree/merge_from_sha_test.go, internal/worktree/reset_repair_test.go, internal/worktree/worktree_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/reauthorize_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR33, VR34, VR35

## What to build

Move every `merge` and `reauthorize` verb call in the listed files onto the verb runner. Each stubbed merge call passes its joins value through the runner, so each `mergeLane` and `mergeReconcile` stub still runs. Delete `runMerge`. Change `mergeFixture` and `reauthorizeFixture` to return one named value each, declared in `verb_fixture_test.go`, and update every caller.

`merge_test.go` is over its line budget, so it must not grow. The `mergeFixture` value carries its joins value, root, and home, and its method builds the verb call value. Each `runMerge(t, j, root, home, ...)` call therefore becomes one runner line. Each positional tuple read becomes one value read, and the deleted `runMerge` body pays the rest. `worktree_test.go` keeps the same rule.

`mergedRecord` selects one record line from stdout and stays, because it reads a record and not a fingerprint. Keep every test name and every assertion.

## Acceptance

- [ ] The VR33 command prints no line.
- [ ] The tuple scan omits `mergeFixture` and `reauthorizeFixture`.
- [ ] The verb form command over the listed files prints only the `LandCommand` line in `delegated_integration_test.go`, which ticket 9 moves.
- [ ] `merge_test.go` and `worktree_test.go` stay at or below their base line counts.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
