# Move the merge and reauthorize verbs onto the verb runner

Blocked by: 4-migrate-query-and-create-verbs.md
Writes: internal/worktree/verb_fixture_test.go (new), internal/worktree/merge_test.go, internal/worktree/merge_caller_root_test.go, internal/worktree/merge_from_sha_test.go, internal/worktree/reset_repair_test.go, internal/worktree/worktree_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/reauthorize_test.go
Covers: VR22, VR23

## What to build

Move every `merge` and `reauthorize` verb call in the listed files onto the verb runner. Each stubbed merge call passes its joins value through the runner, so each `mergeLane` and `mergeReconcile` stub still runs. Delete `runMerge`. Change `mergeFixture` and `reauthorizeFixture` to return one named value each, declared in `verb_fixture_test.go`, and update every caller.

`mergedRecord` selects one record line from stdout and stays, because it reads a record and not a fingerprint. Keep every test name and every assertion.

## Acceptance

- [ ] The ticket's `rg` for `runMerge` prints no line.
- [ ] The tuple scan omits `mergeFixture` and `reauthorizeFixture`.
- [ ] The listed files call no merge or reauthorize verb form directly.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
- [ ] No over-budget test file grows past its base line count.
