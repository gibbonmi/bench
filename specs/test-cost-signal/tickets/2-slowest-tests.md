# Report the three slowest tests of each package

Blocked by: none
Writes: internal/testreport/testreport.go, internal/testreport/test_cost.go (new), internal/testreport/test_cost_test.go (new), internal/testreport/outcome_test.go, internal/testreport/tests_run_test.go, internal/testreport/failure_lines_test.go
Covers: TC12, TC13, TC14, TC15, TC16, TC17, TC18, TC19, TC20, TC21, TC22, TC23, TC24, TC25, TC26

## What to build

`bench test` prints a `slowest[N]{package,test,elapsed_ms}` table on every run that prints the packages table. The table names at most three top-level tests for each package, so an agent sees which tests cost the most.

`decode` keeps the `Elapsed` value of each test terminal event in a new `testResult` field, and the last terminal event of a test wins. One helper converts seconds to whole milliseconds with `math.Round`. The package cell and the test cell both use that helper.

A test is a candidate when its name contains no `/` and its last terminal action is `pass` or `fail`. A test with only `run` and `output` events is not a candidate. Rows sort by package name, then by elapsed time from high to low, then by test name. The row limit is one named constant in `internal/testreport`.

The render order becomes `packages`, `failures`, `skips`, `slowest`, and then the conditional `root_conformance`. An empty table prints `slowest[0]{package,test,elapsed_ms}:`. The full-equality expectations in `outcome_test.go`, `tests_run_test.go`, and `failure_lines_test.go` gain the new block.

This ticket supplies two values to TC-C3: the `test_cost.go` file that holds the new table helpers, and the render order with `slowest` after `skips`. Put new tests in `test_cost_test.go`, and keep `testreport.go` under 400 lines.

## Acceptance

- [ ] A top-level test with terminal `Elapsed` 1.5 renders `canned,TestA,1500` in `slowest`.
- [ ] Two packages with four top-level tests each render six rows, three for each package, in package name order.
- [ ] Rows of one package order from high to low elapsed, and an elapsed tie orders by test name.
- [ ] A subtest, a skipped test, a test with no terminal event, and the package terminal event render no row.
- [ ] A failed test renders a row, and a terminal event without `Elapsed` renders 0.
- [ ] A run with no test events prints the empty `slowest` table right after the `skips` block.
- [ ] With a skipped root conformance test, `slowest` comes before `root_conformance`.
- [ ] A test terminal with `Elapsed` 2.0 and a package terminal with `Elapsed` 0.5 render `canned,pass,500,1` in `packages`.
