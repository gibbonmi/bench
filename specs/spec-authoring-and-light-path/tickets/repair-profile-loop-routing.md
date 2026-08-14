# Keep profile loop routing lean and mutation-backed

Blocked by: none
Writes: projects/benchkit.md, internal/anchors/registry_data.go, internal/canary/inventory_test.go, tests/canary/workflow-guidance-anchors/

## What to build

The Benchkit Lines rows retain cached model, effort, venue, and lean stage identity: loop 1 reviews the spec before slicing, and loop 2 reviews the ticket breakdown after the write-spec slicing step. Operational retry, reporting, advisory-verdict, and reviewer-stop policy stays single-sourced in `bench-write-spec`. Pin each lean row with an exact Require and mutation proof.

## Acceptance

- [ ] the profile no longer duplicates the command-owned loop protocol
- [ ] exact anchors and mutations independently protect the loop-1 and loop-2 routing facts
- [ ] the fixture addition and canary binding count land together
