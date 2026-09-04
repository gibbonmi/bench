# Print the removal error at retirement

Blocked by: remove-the-section-at-retirement.md
Writes: internal/worktree/lifecycle.go, internal/worktree/worktree_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: HS30

## What to build

Repair for review finding C3. Verify the premise first: `executeCleanup`
discards the error from `handoffdoc.RemoveSection`. Then report that error
through the seam the retirement already uses for an advisory failure, and
keep the retirement's verdict. Find that seam by reading how the census drop
and the other non-fatal steps in the same function report.

## Acceptance

- [ ] A release over a document `Parse` refuses completes, and its output names the removal error with the file and line.
- [ ] `TestReleaseDropsTheCensusRecords` and the single-call-site test pass.
- [ ] Self-probe: discard the error again, and report the new test red.
