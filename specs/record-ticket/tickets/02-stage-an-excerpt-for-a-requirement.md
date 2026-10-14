# Stage an excerpt for a requirement

Blocked by: 01-retain-the-last-bounded-response.md
Writes: internal/reviewrecord/repair.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/command_test.go, internal/reviewrecord/recordcmd/stage.go (new), internal/reviewrecord/recordcmd/stage_test.go (new), internal/reviewrecord/recordcmd/verification.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/recordtest/stage.go (new), internal/responsebound/retain.go (new), internal/responsebound/retain_test.go (new), internal/responsebound/responseboundtest/checkout.go, internal/shellcommand/shellcommand.go, tests/canary/cancel-signal-registrations/own-signal-set/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RT9, RT10, RT11, RT12, RT13, RT14, RT15, RT16, RT17, RT57, RT58, RT59

## What to build

`bench record excerpt <slug> --chunk <id> --requirement <id> [--mutant <n>]` stages the retained response for one requirement of the chunk plan at the worktree's `HEAD`.
`FormExcerpt` joins the form constants, and the record command registers the form.
The root help inventory gains the excerpt row.

The stage refuses with no retained response, with an unplanned requirement, and with bytes over the record excerpt limit.
The unplanned-requirement refusal names the planned requirements.

`shellcommand.Parse` splits the planned command, so a quoted operand folds as the shell folds it.
The stage refuses a planned command unless the parse is lexed, holds only word tokens, and starts with the word `bench`.
For a requirement with no probe, the typed verb word and its arguments must equal the planned words after `bench`.
For a requirement with a probe, the retained call must be `bench probe`, and the bytes must import through the probe reader that `importedProbe` uses.
`--mutant` is valid only for a probe requirement, and the stage keeps it for the import.

The stage lives under the assignment scope in one file per key.
The file name is the hex SHA-256 digest of the spec slug, the chunk ID, and the requirement ID, joined by NUL bytes.
No plan string forms a path component.
The stage copies the retained commit and dirty paths with the bytes.
A second stage of one requirement replaces the first.
The excerpt form writes no review record and no worktree file.

The test plan fixture holds `go test ./...` commands today, so the stage tests need a plan with `bench` commands.
`internal/reviewrecord/recordtest/fixture.go` has 5 lines of headroom, so the new fixture helper goes in a new file.
The stage tests need a pool worktree, because `AssignmentScope` answers false outside the pool.

## Acceptance

- [ ] `bench record excerpt` stages the retained `bench test` response for a planned requirement whose command matches.
- [ ] The form refuses with no retained response.
- [ ] The form refuses a retained argument list that differs from the planned command.
- [ ] The form refuses a probe requirement whose retained call is not `bench probe`.
- [ ] The form refuses a mutant-set probe answer without `--mutant`.
- [ ] The form refuses an unplanned requirement and names the planned requirements.
- [ ] A second stage of one requirement replaces the first staged bytes.
- [ ] The form refuses a retained response over the record excerpt limit.
- [ ] The form stages a run that equals a planned command with a single-quoted operand.
- [ ] The form refuses a planned command that does not start with `bench`, that holds a pipe, or that the lexer cannot read.
- [ ] A worktree release removes the retained response and the staged excerpts.
- [ ] A requirement ID of `../x` stages a file inside the assignment scope and nowhere else.
