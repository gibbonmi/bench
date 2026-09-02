# Remove the section at retirement

Blocked by: add-the-handoff-document-leaf-package.md
Writes: internal/worktree/lifecycle.go, internal/worktree/land_journey_test.go, internal/worktree/worktree_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: HS16, HS17, HS18, HS19, HS20

## What to build

Verify the premise first: `executeCleanup` in internal/worktree/lifecycle.go
is the one retirement path for land, release, and clean, and it drops the
census record there. Then remove the retired assignment's section through the
leaf package beside that drop. The key is the request digest, read from
`plan.assignment` before `retireCheckout` mutates the record.
Ensure `main` in the same locked write. A missing section is a no-op.

Guard the single call site with a test in the shape of
`TestCensusDropHasOneCallSiteInThisPackage`. `internal/worktree` imports the
leaf package and never `internal/handoff`.

## Acceptance

- [ ] The landing test extended from `TestLandCommandStatesTheCensusCountAndDropsTheRecords` shows the section gone after the landing.
- [ ] `TestReleaseDropsTheCensusRecords` and `TestCleanDropsTheCensusRecords` extended show the section gone.
- [ ] The single-call-site test passes.
- [ ] After the last section is removed, the file holds `main` with no further verb run.
- [ ] Self-probe: hook the removal into `ReleaseCommand` only, and report the clean test red.
