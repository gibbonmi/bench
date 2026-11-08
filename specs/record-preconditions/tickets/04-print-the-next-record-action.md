# Print the next required record action with its preconditions

Blocked by: 01-collect-record-findings.md, 02-show-the-whole-plan-status.md, 03-show-one-chunk-status.md
Writes: internal/reviewrecord/recordcmd/next.go (new), internal/reviewrecord/recordcmd/next_test.go (new), internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/repair.go, internal/reviewrecord/change.go, internal/reviewrecord/coverage.go, cmd/bench/main.go, cmd/bench/axi_query_cases_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, .agents/skills/bench-craft-cli/SKILL.md, projects/benchkit.md, CHANGELOG.md
Covers: RP23, RP24, RP25, RP26, RP27, RP51, RP52, RP53, RP58

## What to build

Review chunk: RP-C2.

Add the read-only form `bench record next <slug>`. Declare it in `next.go`, add one entry to the `forms` list, and put its form-name constant in `internal/reviewrecord/repair.go`. The RR49 scan reds a bare `"next"` literal in a production `recordcmd` file, so the output table is `action[N]` and the constant lives outside `recordcmd`. The form skips the primary-checkout refusal, and it writes nothing.

The form selects one checkpoint mode at HEAD. A record entry covers a planned chunk directly or through the amendment mapping of `mappedIDs`. If a planned chunk has no covering entry, the form runs the RP-C1 collector in chunk mode on the first such chunk in plan order. Otherwise it runs completion mode. It takes the first finding in checkpoint order and never skips one. A finding with a `Repair` prints that form's command in `action[1]`, from the command helper that ticket 03 extracts in `refusal.go`. A first finding with no `Repair` prints `action[0]` and the finding as a `fails` precondition. An absent record file seeds the collector through `createRecord` (`change.go:11`), as the chunk write does, so the form prints the `record chunk` command of the first planned chunk. A record that owes nothing prints `action[0]` and `preconditions[0]`.

Each precondition row comes from a predicate that the action form's write path reads: the plan digest from `currentPlan`, the required base from the ticket 01 accessor, and the chunk entry from `ChunkEntry`. Approve the child `next` on the `record` row, and add it to `approvedAXIQueries`, the AXI envelope cases, the craft-cli approved-query table, and the AXI list in `projects/benchkit.md`. Keep the net line count of `axi_query_registry_test.go` unchanged. Add the help row to the exact inventory. Add one `CHANGELOG.md` entry under `## [Unreleased]` for RP-C2.

## Acceptance

- [ ] With a recorded first chunk and an unrecorded second chunk, the form prints the `record chunk` command for the second chunk (RP23).
- [ ] The `command` cell equals the repair command that a `bench record chunk` refusal prints for the same fixture (RP24).
- [ ] The chunk action lists the `required base` precondition with the `chunkBase` commit (RP25).
- [ ] After a plan commit, the form prints the `record amendment` command and a `fails` plan digest precondition (RP26).
- [ ] A record that owes nothing prints `action[0]` and `preconditions[0]` at exit 0 (RP27).
- [ ] A first finding with no `Repair` prints `action[0]` and that finding as `fails` (RP51).
- [ ] With no record file, the form prints the `record chunk` command of the first planned chunk (RP52).
- [ ] `bench test --check axi-query-registry` passes with the `next` child approved (RP53).
- [ ] A next read leaves the record bytes, the spec bytes, and `git status` unchanged (RP58).

## Verification

Run `bench test --package ./internal/reviewrecord/...` and `bench test --package ./cmd/bench` whole. Run `bench test --check axi-query-registry`, `bench test --check guidance-prose-budgets`, and the RR49 test through `bench test --package ./internal/conformance --run '^TestNoWriteVerbComposesARouteOutsideTheRegistry$'`. Observe each planned row red before its production edit, and record the route of each red. Self-probe: skip the `currentPlan` precondition, and report that RP26 goes red.
