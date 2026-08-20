# Repair: prove the npm adapter's failure paths at the command surface

Blocked by: swap-workflow-and-flip-conformance.md
Writes (advisory): internal/publication/command_adapter_test.go

Source: review pickup `reviews/ft230-release-through-bench.md`, Coverage
findings 1 and 2 (both auto-fix).

## What to build

Tests only; no production change. Through `publication.Command` with
`--adapter npm`:

1. With no `npm` on PATH, `submit` exits 1 and the output carries the
   adapter's structured `npm ... failed` error as unsatisfied release intent
   (the spec's edge-inventory predicate).
2. With a stub `npm` that fails on an interior platform package, `submit`
   exits 1, the durable record marks the already-published packages and the
   failure, and a re-run against a then-successful stub resumes and completes
   (the resumability predicate the workflow's `if: always()` upload depends
   on).

## Acceptance

- [ ] Absent-binary test asserts exit 1 and the structured error.
- [ ] Mid-sequence-failure test asserts exit 1, the partial record, and a
      successful resume on re-run.
