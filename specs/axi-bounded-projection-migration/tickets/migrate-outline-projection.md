# Migrate outline projection

Blocked by: none
Ownership fence: `internal/outline`
Integration surfaces: shared projection→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-projection-routes.md
Contracts: ordered symbol rows and tracked/scanned/skipped/total/emitted/omitted/truncated facts cross `internal/outline/outline.go`→shared projection, membership is default and full, order is owner scan then project then render, and absence is zero rows with owner metadata, asserted by OL1
Closure: OL1/default-200, OL1/full, OL1/tracked, OL1/scanned, OL1/skipped, OL1/total, OL1/emitted, OL1/omitted, OL1/truncated, OL1/route

## What to build

outline retains every row bound and owner metadata fact through the shared projection.

## Acceptance

- [ ] [OL1] (covers BP4) outline retains every row bound and owner metadata fact through the shared projection.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OL1/default-200 | change the default cap | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/full | cap full mode | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/tracked | derive tracked from rows | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/scanned | derive scanned from rows | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/skipped | omit a skipped input | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/total | use visible rows as total | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/emitted | report total as emitted | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/omitted | derive omitted before filtering | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/truncated | infer truncated from skipped | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| OL1/route | bypass shared projection | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |

