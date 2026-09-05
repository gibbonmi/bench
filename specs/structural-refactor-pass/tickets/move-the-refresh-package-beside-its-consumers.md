# Move the refresh package beside its consumers

Blocked by: none
Writes: internal/refresh/refresh.go (new), internal/refresh/refresh_test.go (new), internal/worktree/worktree.go, internal/worktree/subshell.go, internal/shift/loop.go, internal/conformance/bounds_policy_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SR36, SR37, SR38, SR39, SR40, SR41

## What to build

The line is opus / low. Move the refresh package from
`internal/worktree/refresh` to `internal/refresh` with its tests, so shift
and worktree are its peers. The three moved tests keep their names:
`TestRefreshUsesBoundedNoninteractiveFetch`,
`TestRefreshFailureAndTimeoutAreNonfatalAndDetailed`, and
`TestRefreshOfflineStartsNoGitAndNamesFlag`. They move to
`internal/refresh/refresh_test.go`.

The package gains one entry point that takes the root, a requested flag, and
the output writer. When the caller requests a refresh, the entry point runs
the refresh and renders the table. It returns the new start ref only when the
refresh succeeded, and it returns an empty ref otherwise, per decision (h).

The arg-consuming function and the shift loop both call the entry point, and
the loop maps an empty answer to `HEAD`. The shift loop states the
refreshed-start-ref rule no more. The bounds-policy registry row names
`internal/refresh/refresh.go` as the timeout consumer; that row edit is the
one mechanical rename this ticket makes, per decision (i). The refresh
timeout stays a package variable, and no test swaps it.

The exit proof for this ticket is the pre-existing suite, green with its test
logic unchanged. A mechanical rename is permitted. A needed assertion change
stops the ticket and reports.

## Acceptance

- [ ] `bench test --check bounds-policy` stays green with the row naming `internal/refresh/refresh.go`.
- [ ] `bench shift --refresh` selects the fetched remote head as the start ref and prints the refresh table.
- [ ] The entry point with the flag off writes nothing and answers an empty ref.
- [ ] The entry point with the flag on under `BENCH_OFFLINE=1` writes the offline row and answers an empty ref.
- [ ] `create --help` performs no refresh, and `create --from` with `--refresh` refuses before the refresh.
- [ ] The three moved refresh tests pass at the new path.
- [ ] `bench consumers refresh.RefreshedStartRef` lists one caller, the entry point.
- [ ] Self-probe: swap the empty answer for the remote head after a failed refresh, and report the observed red.
