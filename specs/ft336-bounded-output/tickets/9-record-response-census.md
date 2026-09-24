# 9. Record the response size of each verb in the census

Blocked by: 2-bound-every-public-response.md
Writes: internal/census/census.go, internal/census/output.go (new), internal/census/output_test.go (new), cmd/bench/census_output_test.go (new), internal/worktree/land.go, internal/worktree/exec.go, internal/worktree/land_census_output_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: BO57, BO58, BO59, BO60, BO61, BO62

## What to build

After the owner finishes, the dispatcher appends one output record to `<id>.output` in the census directory. The line holds the time, the verb head, the line count, the byte count, and `spilled` or `inline`. The record uses the assignment of the working tree, or the target assignment of `bench worktree exec`. `ExecCommand` reports the assignment that it resolved, so the dispatcher runs no second resolution. A process with neither writes no record. A failed write changes neither the response nor the exit code.

The raw-call readers keep reading only the `<id>` file. `census.Drop` removes both files. The landing prints `census output{<head>=<calls>/<bytes>,...}` on stderr beside the `census heads` line. The spec-authoring commit already names the output records in the census term of `CONTEXT.md`.

## Acceptance

- [ ] A bounded verb run in an assignment worktree appends one output line with the head, lines, bytes, and `spilled`.
- [ ] `bench worktree exec <target>` from the primary checkout records under the target assignment with the head `bench worktree exec`.
- [ ] `census.Counts` for an assignment is unchanged after 3 output records.
- [ ] The landing prints `census output{bench worktree list=2/28978}` on stderr for two canned records.
- [ ] `census.Drop` removes both the raw-call file and the output file.
- [ ] A symlink at the census directory leaves the verb's response and exit code unchanged.
