# Make socket rows run at long paths

Blocked by: none
Writes: internal/testrepo/socket.go, internal/testrepo/socket_test.go, internal/worktree/clean_landed_hostile_test.go, internal/conformance/prose_budget_test.go, internal/conformance/skill_description_budgets_test.go, internal/conformance/guidance_token_sweep_test.go, internal/landing/state_test.go, internal/testreport/selection_test.go, internal/freshness/freshness_verify_test.go, internal/preflight/charge_evidence_test.go, internal/preflight/charge_file_kinds_test.go, internal/chargeevidence/store_test.go
Covers: none

## What to build

A Unix socket address has a platform length limit of about 104 to 108 bytes. Today each test that needs a socket special file binds a listener at the requested path. When the path is longer than the limit, the bind fails with "invalid argument", and the row skips through its capability route.

Four rows skip on every local run for this reason. Two are in the worktree hostile-path tests. The other two are in the prose-budget and skill-description conformance tests. A longer temporary base adds bytes to every path, so more rows can start to skip.

One helper in `internal/testrepo` makes a socket special file at a requested path. It binds the socket in a short private directory under the temporary directory, and then it renames the socket file to the requested path. The length of the requested path thus has no effect. The helper returns an error, and each site keeps its current route for a host that cannot make a socket.

Each test site that needs only a socket file at a path calls the helper. The `chargeevidence` store test changes its working directory only to keep the socket address short. It calls the helper and does not change its working directory. The `preflight` evidence test makes a short temporary directory only for its socket row, so that directory helper goes away. The `bounds` package tests keep their own bind, because `internal/testrepo` depends on `internal/bounds` and a `bounds` test cannot import it. Their socket paths are short.

## Acceptance

- [ ] The helper test makes a socket at a path longer than 108 bytes and finds a socket special file at that path.
- [ ] The two worktree socket rows and the two conformance socket rows run and pass under `go test -count=1 -v`, with no SKIP line.
- [ ] Each migrated site keeps its current route for a host that cannot make a socket.
- [ ] No test file outside the helper and the `bounds` package binds a Unix socket listener.
- [ ] The `skip-ownership` and `root-conformance` checks and each touched package pass.
