# Remove the drift clause from the hook

Blocked by: none
Writes: internal/adopt/prepush.sh, internal/adopt/link_hook_test.go, internal/guards/guards_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: PR1, PR2, PR3, PR4, PR5

## What to build

The five registry paths after the guards test in `Writes:` are preflight
closure only. This ticket edits none of them, and the verb ticket owns
their edits. The overlap creates no blocker edge.

Verify the premise first: read internal/adopt/prepush.sh, which reads a pin
path, prints a `gate unpinned` warning, and compares the pushed `.bench`
tree with the pin. Read `TestPrePushHookAllowProtectedPushConfig`,
`runPrePushHook`, `writeHook`, and `hookTestRepo` in
internal/adopt/link_hook_test.go. Read
`TestCommandRendersRealStaleManagedPrePushHookAndRepairAction` in
internal/guards/guards_test.go, whose expected `denies` cell quotes the
hook header.

Rewrite the hook so it keeps one clause. It resolves the protected branch
from `origin/HEAD` with the baked token as the fallback. It reads
`bench.allowProtectedPush`. It loops the stdin ref lines once. When a remote
ref is the protected branch and the config is not `true`, it exits 1 with
`blocked: direct push to <branch>`.

The hook reads no pin path, prints no warning, and runs no `git rev-parse`
on the pushed oid. Keep the read loop's `|| [ -n "$line" ]` guard. Set the
header `denies` field to `direct push to the protected branch`. Keep the
`why` field free of a repo-only path. Keep the `bench:managed-pre-push`
marker line and the `__BENCH_DEFAULT_BRANCH__` token.

Add two hook tests beside the existing one. `TestPrePushHookTopicPushIsSilent`
pushes `refs/heads/topic` and expects exit 0 with an empty stderr.
`TestPrePushHookIgnoresALegacyPin` writes a `bench-gate-pin` file in the
git dir whose first line is a foreign tree hash. It then pushes
`refs/heads/topic` with the harness's non-commit oid and expects exit 0
with an empty stderr. Update the guards test's expected cell to the new
header value.

## Acceptance

- [ ] The rendered hook exits 1 on `refs/heads/main` and prints `blocked: direct push to main`.
- [ ] With `bench.allowProtectedPush` set to `true` the hook exits 0 on `refs/heads/main`, and set to `false` it exits 1.
- [ ] The hook exits 0 on `refs/heads/topic` with an empty stderr.
- [ ] With a legacy `bench-gate-pin` file that names a foreign tree, the hook exits 0 on `refs/heads/topic` for a non-commit oid with an empty stderr.
- [ ] `bench guards` renders the pre-push `denies` cell as `direct push to the protected branch`.
- [ ] Self-probe: restore the drift loop alone, and report the legacy-pin test red.
