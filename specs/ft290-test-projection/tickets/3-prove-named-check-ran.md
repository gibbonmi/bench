# 3. Prove that a named check ran

Blocked by: 1-split-named-check-owner.md, 2-count-tests-run.md
Writes: internal/testreport/
Covers: TP3, TP4, TP5, TP6, TP18, TP51, TP55

## What to build

Chunk: TP-C1a.

Print the row `check[1]{name,kind,tests_run,subjects}` as the first block of each Go-backed named-check result that reached a verdict.
Add the one owner that maps a check name to its kind: `conformance`, `system`, or `prose`.

A Go-backed check prints `subjects` 0.
When `tests_run` is 0, print the `check` row and the error title `named check ran nothing`, and exit 1.

A package run with no test keeps exit 0.
Add a run event to each canned `go` builder that a named-check test uses, inside its current line.

This ticket supplies the `check` row producer and the kind owner that tickets 4, 6, and 9 read.

## Acceptance

- [ ] `--check line-routing` with one run event prints the row `line-routing,conformance,1,0` before the packages table.
- [ ] `--check system` prints the kind `system`.
- [ ] A named check with no run event exits 1 with the title `named check ran nothing`.
- [ ] `--package chosen` with no run event exits 0.
- [ ] A `--changed` run with no run event exits 0.
- [ ] A named check gets a `build-fail` event, no run event, and a nonzero child exit. It prints the compile diagnostic at exit 1, without the zero-rule title.
- [ ] The real `--check ordinary-build-census` run prints `tests_run` 1.
