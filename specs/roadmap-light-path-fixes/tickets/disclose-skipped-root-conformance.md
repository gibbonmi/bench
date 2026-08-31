# Disclose skipped root conformance

Blocked by: route-prose-check-and-print-all-findings.md
Writes: internal/testreport/testreport.go, internal/testreport/testreport_test.go, internal/testreport/selection.go, internal/testreport/selection_test.go
Covers: LF27

## What to build

When TestRootConformance skips during a package run, prevent the package pass
from reading as root green. Name the graded-root named-check route.

## Acceptance

- [ ] A skipped root test is visible in the package result.
- [ ] The package result does not claim graded-root green.
- [ ] The result names bench test --check as the root route.

