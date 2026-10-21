# Warn first when a package uses half of the timeout

Blocked by: 1-timeout-owner.md, 2-slowest-tests.md
Writes: internal/testreport/testreport.go, internal/testreport/test_cost.go (new), internal/testreport/test_cost_test.go (new), internal/testreport/timeout_warning_test.go (new), internal/testreport/check_row_test.go, internal/testreport/outcome_test.go, internal/testreport/canned_go_test.go, internal/testreport/tests_run_test.go, internal/conformance/bounds_policy_test.go, CHANGELOG.md
Covers: TC27, TC28, TC29, TC30, TC31, TC32, TC33, TC34, TC35, TC36, TC37

## What to build

`bench test` prints a `timeout_warnings[N]{package,elapsed_ms,timeout_ms}` table when a package uses half or more of the go test timeout. The table renders only when it has a row. It is the first block of the `render` result, so it stays in the visible head of a bounded response. Under `--check`, the check row comes first, then the warning, then `packages`.

A package gets a row when its `elapsed_ms` is at least `bounds.GoTestTimeout.Milliseconds()` divided by 2. The divisor is one named constant in `internal/testreport`, so the threshold is never a second duration. The `timeout_ms` cell is `bounds.GoTestTimeout.Milliseconds()`. Rows sort by package name.

A warning changes no `Outcome` field and no exit code. A passing warned run exits 0, and a failing warned run exits 1 and still prints its failure row.

The ticket consumes `bounds.GoTestTimeout` from `1-timeout-owner.md` and the `test_cost.go` helpers and render order from `2-slowest-tests.md`. Put the warning tests in `timeout_warning_test.go`, and drive them through `commandOverEvents` with canned events. Build each threshold value from the constant: the fixture at the threshold uses half of the timeout, and the fixture below it uses one millisecond less.

The bounds policy owners map gains the entry `internal/testreport/test_cost.go` with the token `bounds.GoTestTimeout`. The `CHANGELOG.md` Unreleased section states the timeout flag and its effect on a linked repository that sets a `GOFLAGS` timeout. It also states the `slowest` table and the `timeout_warnings` table. The comment in `check_row_test.go` says that the packages table follows the check row. Change it to state that a warned run puts `timeout_warnings` between the two. Keep `testreport.go` under 400 lines, because this ticket grows it after ticket 2.

## Acceptance

- [ ] A package at the threshold renders the row `canned,<GoTestTimeout.Milliseconds()/2>,<GoTestTimeout.Milliseconds()>`.
- [ ] A package one millisecond below the threshold renders no `timeout_warnings` block.
- [ ] Two warned packages, whose events arrive in reverse name order, render in package name order.
- [ ] A warned plain run starts with the `timeout_warnings` header.
- [ ] A warned named check run renders the check row, then `timeout_warnings`, then `packages`.
- [ ] A passing warned run exits 0, and a failing warned run exits 1 with its failure row.
- [ ] A healthy run renders no `timeout_warnings` block, so the base expectations do not change.
- [ ] `bench test --check bounds-policy` turns red when `internal/testreport/test_cost.go` does not contain `bounds.GoTestTimeout`.
- [ ] The changelog Unreleased section names the flag, the linked-repository effect, and both tables.
