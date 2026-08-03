# Repair spec prose and comment consistency for the round-2 findings

Blocked by: none
Ownership fence: `specs/injected-interface-junctions/spec.md`, `internal/conformance/injected_ports_test.go`
Contracts: the spec's message count and story-2 prose cross spec→implementation; after this ticket the spec contracts five distinct failure messages including the orphan row (asserted by SC1 against the amended map) and story 2 no longer claims a real-planner unreadable path; comment edits cross nothing (none crosses beyond the spec rows named)
Assumptions: sol round-2 findings SP3/SP4 (blocking) and S2/S3 (comment accuracy/timelessness) are the authority; no behavior change anywhere; claims re-derived from the tree at pickup

## What to build

Spec: amend story 2's prose so the unreadable-metadata clause matches the
repaired row (pre-composition refusal, no owner composition); amend the
implementation decisions to contract five distinct fail-closed messages
including the orphan row; add a story-6 coverage row for the orphan diagnostic
with its observed red (the R2 unit red and the planted-orphan-row scoped red
are both recorded evidence). Comments: fix the orphan fixture comment to name
the actual orphan port, and replace the review-provenance comment with a
timeless statement of the invariant it protects.

## Acceptance

- [ ] [SC1] the spec contracts five distinct messages, story 2's prose matches its rows, the orphan diagnostic has a coverage row, and `bench coverage --check` passes.
- [ ] [SC2] the two comments read accurately and timelessly and `go test -count=1 ./internal/conformance -run TestInjectedPort` stays green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SC1 | restore the "four" message count | review grep plus the map | `rg "four distinct" specs/injected-interface-junctions/spec.md` returns no hit after the ticket |
| SC2 | restore the review-provenance comment text | review grep | `rg "sol's C2" internal/conformance/injected_ports_test.go` returns no hit after the ticket |
