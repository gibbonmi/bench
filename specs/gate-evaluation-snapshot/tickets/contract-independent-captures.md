# Contract independent capture routes

Blocked by: migrate-evaluation-identities.md
Ownership fence: `internal/gate/`
Contracts: the evaluation-owned generation contract in `internal/gate/tree_snapshot.go` crosses all identity call sites inside `internal/gate/`, asserted by CI1-CI3 through complete production-call enumeration and public equivalence controls

## What to build

Remove the superseded root-based identity capture routes and compatibility plumbing so production evaluation has one generation lifecycle, then retain structural guards and journey controls that make any later independent parser, materialization, or blob reader visible.

## Acceptance

- [ ] [CI1] No production identity resolver constructs a source adapter, parses a tree listing, or reads a snapshot blob outside the evaluation-owned generation lifecycle.
- [ ] [CI2] The single-listing parser and complete production-call enumeration reject a second parser or any direct identity-family recapture while allowing the distinct under-lock validation and post capture.
- [ ] [CI3] Stable ordinary and prospective journeys retain stdout, stderr, exits, inspection reasons, verdict fields, component/check evidence, stripped evidence, exact-tree behavior, and the accepted source-cost ceilings.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CI1 | add a direct root-based snapshot read to one identity resolver | the production capture-route structural test | add the call, run the focused structural test, expect the forbidden route report |
| CI2 | add a second ls-tree parser beside the generation source | the single-listing-parser test | add the parser signature, run the focused test, expect both files to be enumerated |
| CI3 | change one durable or operator-facing field while counts stay bounded | the ordinary and prospective journey controls | run the focused equivalence tests, expect the literal output or decoded-record mismatch |
