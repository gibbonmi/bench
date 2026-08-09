# Migrate sanitize projection

Blocked by: none
Ownership fence: `internal/sanitize`
Integration surfaces: shared projection→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-projection-routes.md
Contracts: preview content string, original byte integer, truncated boolean, code-point unit, and controls disposition cross `internal/sanitize`→shared projection, membership is Preview versus uncapped Controls, order is select-escape-suffix, and absence is empty input, asserted by SN1
Closure: SN1/code-points, SN1/original-bytes, SN1/controls, SN1/backslash, SN1/suffix, SN1/uncapped, SN1/route

## What to build

sanitize retains its exact preview and Controls policy through the shared projection.

## Acceptance

- [ ] [SN1] (covers BP1) sanitize retains its exact preview and Controls policy through the shared projection.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SN1/code-points | count UTF-8 bytes instead of code points | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| SN1/original-bytes | report selected bytes as original | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| SN1/controls | drop one control escape | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| SN1/backslash | drop backslash escaping | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| SN1/suffix | change the truncation suffix | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| SN1/uncapped | cap Controls at 120 | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| SN1/route | bypass shared projection | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |

