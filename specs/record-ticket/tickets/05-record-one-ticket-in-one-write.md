# Record one ticket in one write

Blocked by: 02-stage-an-excerpt-for-a-requirement.md, 04-extract-the-record-change-functions.md
Writes: internal/reviewrecord/ticket.go (new), internal/reviewrecord/ticket_test.go (new), internal/reviewrecord/change.go (new), internal/reviewrecord/repair.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/fold.go, internal/reviewrecord/fold_test.go, internal/reviewrecord/export_test.go, internal/reviewrecord/recordcmd/ticket.go (new), internal/reviewrecord/recordcmd/ticket_test.go (new), internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/command_test.go, internal/reviewrecord/recordcmd/refusal.go, internal/reviewrecord/recordcmd/planned_cell.go, internal/reviewrecord/recordcmd/verification.go, internal/reviewrecord/recordcmd/stage.go (new), internal/reviewrecord/recordcmd/ordering_test.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/recordtest/stage.go (new), internal/responsebound/retain.go (new), cmd/bench/record_operations_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RT18, RT19, RT20, RT21, RT22, RT23, RT24, RT25, RT26, RT27, RT28, RT29, RT30, RT31, RT32, RT33, RT34, RT35, RT36, RT38, RT55, RT56, RT60, RT61, RT62, RT63

## What to build

`bench record ticket <slug> --chunk <id> --ticket <basename> --id-prefix <prefix> --performer <session> --model <model> --effort <effort> --ref <ref> [--base <commit>] [--tip <commit>]` writes one ticket's entries.
`--tip` defaults to `HEAD`.
The owed requirements are the requirements of the chunk plan at the tip whose ticket is the named ticket.

Before the write, the form refuses in this order:

1. a ticket that the chunk plan does not list
2. a plan whose chunk requirements name no ticket
3. a planned seam cell at the tip, when the ticket is the last ticket of the chunk
4. an owed requirement with no staged excerpt, naming every missing requirement
5. a stage whose run commit is not the chunk source, naming the requirement
6. a stage whose run had a dirty path other than the record path of the spec, naming the requirement

Inside one `write` transaction, the form applies the amendment when the plan digest moved.
It then applies the freeze for the last ticket and one verification for each owed requirement.
The composed transaction lives in `internal/reviewrecord`, because `write` and `coveredAfter` are unexported there.
The transaction composes the change functions and the create function of ticket 04.
An earlier ticket's call passes the absent-record refusal as its create function.

The amendment step refuses a plan change that needs a chunk map, and the form wraps that error as a repair that names `bench record amendment`.
The base is `--base`, then the base of the existing chunk entry, then the commit that `coveredAfter` answers for the predecessor chunk's entry.
With none of them, the form refuses and asks for `--base`.
A call for an earlier ticket needs an existing chunk entry, and it never changes that entry.

Each evidence ID is the prefix, a hyphen, and the requirement ID.
The exit code comes from the stage, and a probe requirement imports its staged answer with its staged `--mutant`.
After the write succeeds, the form removes the staged excerpts that it used.

`FormTicket` joins the form constants, and the record command registers the form.
The root help inventory gains the ticket row.
The ticket form's help derives its step list from `CompletionOrder`: the steps before the first review step.

## Acceptance

- [ ] For the last ticket, the form writes the chunk entry at the tip and one verification entry for each owed requirement.
- [ ] With a moved plan digest, the form writes one amendment before the chunk entry, and with an unchanged digest it writes none.
- [ ] A refusal at the last verification and a refusal at the chunk step each leave the record bytes unchanged.
- [ ] One refusal names every owed requirement with no stage.
- [ ] An earlier ticket's call refuses with no chunk entry and names the last ticket, and with an entry it leaves the tip and base unchanged.
- [ ] A refreeze keeps the existing base, and a new entry takes the predecessor tip or its last recorded fold.
- [ ] A call for the first chunk with no base refuses and asks for `--base`.
- [ ] A plan change that removes a recorded chunk refuses and names `bench record amendment`.
- [ ] The form refuses a planned seam cell, a ticket outside the chunk, and a plan whose requirements name no ticket.
- [ ] A probe requirement records the probe values of the staged answer.
- [ ] Used stages are absent after a written call and present after a refused call.
- [ ] The form refuses a stage from another commit or from a run with a dirty Go file.
- [ ] The form accepts a stage from a run with only the record path dirty.
- [ ] The first chunk's last ticket creates an absent record with `--base`.
- [ ] `bench record` help lists the ticket and excerpt forms, and the ticket form's step list follows `CompletionOrder`.
- [ ] The probe mutation that runs the three steps as three writes turns RT21 and RT22 red.
