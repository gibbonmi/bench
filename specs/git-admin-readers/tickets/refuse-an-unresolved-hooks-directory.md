# Refuse an unresolved hooks directory

Blocked by: add-the-named-git-admin-readers.md
Writes: internal/adopt/link_hook.go, internal/adopt/link_transaction.go, internal/adopt/link_transaction_test.go, internal/adopt/unlink.go, internal/adopt/doctor.go, internal/adopt/link_hook_test.go, internal/adopt/adopt_test.go
Covers: GR20, GR21, GR22, GR23

## What to build

Verify the premise first. Read `AdminPath` in internal/git/worktree_admin.go, which
the blocker ticket adds. Read `hooksDir` at internal/adopt/link_hook.go line 200,
which guesses the path `.git/hooks` when the Git call fails. Read `InspectPrePush`,
`noPrePushHealth`, and `PrePushAbsent` in the same file. Read `installGitHook` at
line 228. Read `removeManagedHook` at internal/adopt/unlink.go line 271, and read
the repair switch at internal/adopt/doctor.go lines 392 to 425.

Read `transactionalLink` at internal/adopt/link_transaction.go line 131. It calls
`InspectPrePush` at line 157 and stages the hook beside the record's path at line
330, so the installer never runs inside `bench link`. The doctor's repair returns 0
on an absent hook outside the kit checkout, so it never reaches the installer
either.

Read `TestPrePushHookAllowProtectedPushConfig` in internal/adopt/link_hook_test.go
as the hook harness prior art. The doctor holds a shell snippet at
internal/adopt/doctor.go line 129. That snippet keeps its own Git spelling, because
it runs without `bench` on `PATH`.

Change `hooksDir` to call `git.AdminPath` with the name `hooks` and to return the
typed failure. Remove the guessed `.git/hooks` path. In `transactionalLink`,
resolve the hooks directory through `hooksDir` before the `InspectPrePush` call,
and refuse with exit 1 before any file is staged when it fails. In the doctor's
repair, resolve it the same way before the state switch, and return 1 when it
fails. Print the root path and the action `investigate the git failure` at both
refusals. Return the failure from the hook installer and from the hook removal,
and make `bench unlink` exit 1 the same way.

Read `TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent` in
internal/adopt/link_transaction_test.go. Its root is a bare temporary directory,
so the new refusal would red its converged case. Initialize a real repository at
that root, and keep every assertion of the test unchanged.

Report the absent state with an empty path from `InspectPrePush` when the reader
fails. Add no new health state, because the reviewer decided the absent state
carries the case.

Add the new tests to the adopt suites. Put
`TestLinkRefusesUnresolvedHooksDirectory` and
`TestInspectPrePushReportsAbsentWhenHooksDirectoryIsUnresolved` in
internal/adopt/link_hook_test.go. Put `TestDoctorFixRefusesUnresolvedHooksDirectory`
and `TestUnlinkRefusesUnresolvedHooksDirectory` in internal/adopt/adopt_test.go.
Run each verb over a repository with the `fail-git-path` stub on `PATH`.

## Acceptance

- [ ] `bench link` exits 1 under the stub before it stages any file, and it prints the root path and `investigate the git failure`.
- [ ] `bench doctor --fix` exits 1 the same way in a consumer repository under the stub.
- [ ] `bench unlink` exits 1 the same way over a present managed hook.
- [ ] The hook health record holds the state `PrePushAbsent` with an empty `Path` under the stub.
- [ ] `hooksDir` returns a typed failure and never guesses a path.
- [ ] The doctor's shell snippet keeps its own Git spelling.
- [ ] The pre-existing `internal/adopt` suite passes, except the four rows this ticket changes.
- [ ] `TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent` passes over a real repository with its assertions unchanged.
- [ ] Self-probe: restore the guessed `.git/hooks` path, and report the link refusal test red.
