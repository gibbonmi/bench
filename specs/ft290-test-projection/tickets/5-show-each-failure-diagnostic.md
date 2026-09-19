# 5. Show each failure diagnostic

Blocked by: 2-count-tests-run.md
Writes: internal/testreport/
Covers: TP13, TP14, TP15, TP16, TP17, TP49

## What to build

Chunk: TP-C1.

Change the failures header to `failures[N]{package,test,line,lines}` in each mode.
The `lines` cell is the count of diagnostic lines of that test.

The default mode keeps one row for each failed test. `--full` prints one row for each diagnostic line in emitted order.
A package failure with no test name obeys the same rules over its package log.

`Outcome.FailedTests` counts distinct failed tests, so the `bench probe` verdict row does not change under `--full`.

## Acceptance

- [ ] A failed test with three lines prints one default row with `lines` 3.
- [ ] A failed test with no diagnostic prints `no diagnostic emitted` and `lines` 0.
- [ ] `--full` prints three rows for that test, and no cell holds an escaped newline.
- [ ] A package failure with two log lines prints two `--full` rows with an empty `test` cell.
- [ ] `Outcome.FailedTests` is 1 for that `--full` run.
- [ ] `bench test --package ./internal/probe` is green.
