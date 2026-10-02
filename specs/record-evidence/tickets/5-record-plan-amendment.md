# 5. Record a plan amendment with `bench record amendment`

Blocked by: 4-record-review-result.md
Writes: internal/reviewrecord/coverage.go, internal/reviewrecord/write.go (new), internal/reviewrecord/recordcmd/ (new), cmd/bench/help_inventory_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RE84, RE85, RE86, RE87, RE88, RE89, RE90, RE91, RE92, RE93, RE94, RE95, RE96, RE97, RE114, RE115

## What to build

Chunk: RE-C4.

Add the `amendment` form through the transaction of ticket 2. Put its record change in `internal/reviewrecord/coverage.go`, beside `mappedIDs`, so the mapping rule stays private and has one owner. `recordcmd` parses the flags and prints the output.

`from` is the record's current `plan_digest`, and `to` is the `ReadPlan` digest at the `--source` tree. The record's `plan_digest` then becomes `to`. An unchanged plan exits 1 and names `unchanged`.

The amendment keys are the recorded chunk IDs under `from`. Map each recorded chunk through the earlier amendments with `mappedIDs`, so the checkpoint and the verb never disagree. Each key maps to itself unless a `--map` names it. A key without a map must be a chunk of the new plan. Each `--map` key must be a recorded key, and each target must be a chunk of the new plan.

Print the `amendment[1]{from,to,chunks}` row of the spec. Add the amendment row to `HelpRows` and to the help golden.

## Acceptance

- [ ] `from` and `to` match the old record digest and the `ReadPlan` digest at `--source`.
- [ ] The record digest moves to `to`, and a second amendment chains from the first `to`.
- [ ] Each recorded chunk maps to itself, and `--map 1=1a,1b` writes the split.
- [ ] An unchanged plan, an unmapped missing chunk, an unrecorded map key, and an unplanned map target each exit 1 and change nothing.
- [ ] With no record the form exits 1 and names `bench record chunk`.
- [ ] After a plan commit, the amendment lets `reviewrecord.Check` accept chunk `2`.
- [ ] Each malformed `--map` exits 2, and `bench help` prints the amendment row.
