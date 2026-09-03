# Repair the review round one findings

Blocked by: share-the-purity-census-and-count-process-fixtures.md, refuse-an-unheld-or-invalid-build-cache.md
Writes: specs/roadmap-light-path-fixes-3/spec.md, reviews/roadmap-light-path-fixes-3.md, internal/puritycensus, internal/worktree/lifecyclepolicy/purity_census_test.go, internal/worktree/reclaimpolicy/purity_census_test.go, internal/worktree/landingpolicy/purity_census_test.go, internal/canonicalpath/purity_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LQ2, LQ5, LQ17, LQ27, LQ28

## What to build

This ticket records the repairs the initial review accepted. The pickup
file `reviews/roadmap-light-path-fixes-3.md` names each finding.

The spec amends five rows. LQ2 and LQ5 name the package tests that close
them. LQ17 and LQ28 state the HOME-declared condition, and the cache
decision line records that exception for reviewer veto. LQ27 grades a
production function, not a file.

The `puritycensus` package unexports `Diagnose` and the `Policy` fields,
because no caller outside the package uses them. The four census wrapper
comments stop restating the forbidden imports and the ambient effects
that the helper owns. They also take the apostrophe the prose rule asks for.

## Acceptance

- [ ] `bench coverage --check` accepts the amended rows.
- [ ] The four wrapper censuses and the helper tests pass after the unexport.
- [ ] `bench gate-prose` passes on the spec and the pickup file.
- [ ] Self-probe: re-export `Diagnose`, and report that no test reds; the finding is a surface rule, not a behavior.
