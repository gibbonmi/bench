# Move the query and create verbs onto the verb runner

Blocked by: 5-migrate-cleanup-verbs.md
Writes: capture/restructure-backlog.md, internal/worktree/worktree_test.go, internal/worktree/list_actions_test.go, internal/worktree/list_selected_test.go, internal/worktree/path_identifier_test.go, internal/worktree/request_token_test.go, internal/worktree/show_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/exec_test.go, internal/worktree/exec_pwd_test.go, internal/worktree/build_test.go, internal/worktree/merge_from_sha_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR31, VR32

## What to build

Move every `create`, `list`, `path`, `show`, `exec`, `build`, `pool`, `lease-file`, `release`, and `clean` verb call in the listed files onto the verb runner. An exec call that feeds stdin passes it through the runner. A `build` call with a joins value passes that value through the runner.

Delete `runCreate` and `execAtOwnedTarget`. Move the `runCreate` callers in `merge_from_sha_test.go` onto the runner. Replace the inline 64-hex matches in `worktree_test.go` and `identifier_operand_test.go` with `mustFingerprint`. Update the `runCreate` split anchor in `capture/restructure-backlog.md`, so that the row names a declaration that still exists.

`childFailure` in `exec_test.go` calls the exec child runner, not a verb form, so it stays. The `runMerge` and `mergeFixture` uses stay for ticket 7. Keep every test name and every assertion. The over-budget files `worktree_test.go` and `exec_test.go` stay at or below their base line counts.

## Acceptance

- [ ] The VR31 command prints no line.
- [ ] The verb form command over the listed files prints no line.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
- [ ] No over-budget test file grows past its base line count.
