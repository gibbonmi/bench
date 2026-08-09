# Migrate roadmap projection

Blocked by: none
Ownership fence: `internal/roadmap`
Integration surfaces: shared projection→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-projection-routes.md
Contracts: context string, original byte integer, truncated boolean, byte unit, and full-mode disposition cross `internal/roadmap/context_types.go`→shared projection, membership is every limited call, order is owner select then render, and absence is empty content, asserted by RM1
Closure: RM1/byte-cap, RM1/utf8, RM1/original-bytes, RM1/truncated, RM1/full, RM1/single-pass, RM1/route

## What to build

roadmap retains its UTF-8-safe 4096-byte default and uncapped full policy through one projection.

## Acceptance

- [ ] [RM1] (covers BP2) roadmap retains its UTF-8-safe 4096-byte default and uncapped full policy through one projection.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RM1/byte-cap | count runes instead of bytes | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| RM1/utf8 | cut inside a UTF-8 sequence | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| RM1/original-bytes | report emitted bytes as original | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| RM1/truncated | derive truncation from content | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| RM1/full | cap full mode | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| RM1/single-pass | project an already bounded value twice | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| RM1/route | bypass shared projection | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |

