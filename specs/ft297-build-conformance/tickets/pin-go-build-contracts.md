# Pin the Go build contracts

Blocked by: none
Writes: internal/conformance/ordinary_build_census_test.go, internal/conformance/build_contracts_test.go (new), internal/conformance/check_bindings_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/checks.go (new), internal/conformance/registry/registry.go, internal/freshness/freshness_digest_test.go, projects/benchkit.md, reviews/ft297-build-conformance.md (new), specs/ft297-build-conformance/tickets/pin-go-build-contracts.md (new)
Covers: none

## What to build

Add two checks to the existing conformance registry for FT297.
The first check compares the Go executable owner with both shell path derivations.
The second check requires the VCS-disable flag on each direct Go build or list command.
The check includes production files and test files.

The current source has one missing flag in a freshness test diagnostic.
Protect that diagnostic without changing the test assertion.
Move the registry table and binding helpers to separate files to create line headroom.
Keep the existing registry facts single-sourced.

## Acceptance

- [x] The path check rejects a changed Go owner or either changed shell derivation.
- [x] The path check accepts matching derivations and names each inconsistent source.
- [x] The command check rejects a build or list call without the required flag.
- [x] The command check accepts the literal flag and the gate-owned constant.
- [x] The command check includes test files and named exec import aliases.
- [x] The freshness test still accepts malformed ambient Git metadata.

## Verification

Use fixture tests for matching paths, changed paths, omitted flags, and both accepted flag forms.
Run both registered checks against the assignment tree.
Run the freshness regression with malformed ambient Git metadata.
Use a swap probe on the resolver path through the registered path check.
Use a swap probe on the freshness build-input flag through the registered command check.
The coordinator uses a different mutation kind and source site for its independent probe.

## Scope

This ticket follows the approved light path at the existing conformance seam.
The user approved two isolated implementation authors and the verified local landings.
The coordinator owns roadmap, changelog, retrospective, and scorecard updates.
No sentence-boundary change belongs to this ticket.

## Retained author

Author: ft297_author
Line: gpt-6-astra / ultra
Pre-review attempt cap: 3
Post-review repair allowance: 2
Post-review repairs consumed: 1

## Accepted review repairs

Share the import resolver with the architecture census for ST1.
Use the existing bounded classifier before every source read for COV1.
Prove FIFO rejection at the Go owner, both shell paths, and the Go source sweep.
