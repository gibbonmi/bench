# Derive test wait deadlines from bounds

Blocked by: route-the-tickets-only-close-to-the-landing-verb.md
Writes: internal/bounds/bounds.go, internal/bounds/bounds_test.go, internal/runbinary/runbinary_test.go, internal/systemtest/owner_teardown_test.go, internal/systemtest/owner_test.go, internal/systemtest/otel_crash_test.go, internal/systemtest/owner_artifact_recovery_test.go, internal/worktree/classifier_shape_test.go, internal/worktree/subshell_test.go, internal/otelrecord/writer_test.go, internal/gocache/lock_test.go, internal/gate/prospective_owner_test.go, internal/status/status_producible_test.go, tests/canary/package-core-guard/bounds-duplicate-owner, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LP9, LP10

## What to build

Verify the premise first: fourteen outer wait deadlines in the named test files
carry a literal duration, and each prints a bare `t.Fatal` on expiry. Add one
renderer to `internal/bounds` beside `TestDeadline` that returns the timeout
verdict text for a wait name and a window, with no `testing` import. Migrate
every named site so its window derives from `bounds.TestDeadline` and its
expiry prints that verdict. A poll interval stays a literal. `waitForFile` in
`internal/gocache/lock_test.go` gains a verdict; today it returns silently.

The sibling ticket `sweep-literal-wait-deadlines-in-tests` reads this ticket's
migrated sites as its live-tree green set, so leave no literal deadline behind.

Run the system-tagged tests with `BENCH_KIT` and `BENCH_RUN_BINARY` set, as the
system suite requires.

## Acceptance

- [ ] Each of the fourteen named wait sites derives its window from `bounds.TestDeadline`.
- [ ] A migrated helper with a zero inner window prints the verdict text that names the wait and the window on expiry.
- [ ] `bounds_test.go` pins the verdict text for one wait name and one window.
- [ ] The dev-tier test suite for every edited package passes.
- [ ] Self-probe: restore one literal at a migrated site, and report the site name as the observed remaining literal.
