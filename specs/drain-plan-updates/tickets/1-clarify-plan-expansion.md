# Clarify in-scope plan expansion

Blocked by: none
Writes: .bench/BENCH.md, internal/anchors/registry_retained_workflow.go, internal/conformance/retained_workflow_test.go
Covers: none

## What to build

Clarify the existing plan-expansion owner for a delayed in-scope expansion. Require one pre-dispatch learning entry and one enabling plan commit.

The commit updates all affected planning and review artifacts before later-ticket dispatch.

## Acceptance

- [ ] The plan-expansion owner requires the learning entry before dispatch.
- [ ] One enabling plan commit updates the coverage map, ticket contracts, chunk table, and any required review-record amendment.
- [ ] The owner keeps the existing expansion guarantees and their anchor needles.
