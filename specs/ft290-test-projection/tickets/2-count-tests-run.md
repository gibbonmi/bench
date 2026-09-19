# 2. Count the tests that ran on each packages row

Blocked by: none
Writes: internal/testreport/
Covers: TP1, TP2, TP47

## What to build

Chunk: TP-C1a.

Add the `tests_run` cell to each packages row.
The cell is the count of distinct tests and subtests of that package that emitted a run event.

The header becomes `packages[N]{package,status,elapsed_ms,tests_run}` for each form.
Update each test that pins the old header inside its current lines, and put the new tests in a new test file.

This ticket supplies the header that tickets 3, 5, and 10 extend.

## Acceptance

- [ ] A package with two run events prints `tests_run` 2.
- [ ] A `no-tests` package prints `tests_run` 0.
- [ ] A `--package` run prints no `selected_by` text.
- [ ] Each test that pinned the old header passes with the new header.
