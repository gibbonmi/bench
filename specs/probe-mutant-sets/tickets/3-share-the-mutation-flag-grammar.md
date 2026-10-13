# Share the mutation flag grammar and split the form matrix

Blocked by: 1-repeat-the-probe-package.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/combined_test.go, internal/probe/command.go, internal/probe/edit_test.go, internal/probe/outcome_test.go, internal/probe/refusal_test.go, internal/probe/replacement_test.go, internal/probe/verdict_test.go
Covers: none

## What to build

This ticket is a prefactor with no new behavior.
It prepares the seam that the mutant-set entry grammar of ticket 4 reads.
It is its own review chunk, PM-C1, so its review closes before ticket 4 starts.

Move the mutation flag declarations in `internal/probe/command.go` into one list.
The top-level grammar reads that list, and ticket 4 builds its entry grammar from the same list.
Split `gradeForms` into a mutation half and a selection half, and run both halves for the top-level argv.

Keep the rule that refuses `--expect` with `--check prose` in one named function that reads both halves.
Ticket 4 calls that function for each entry against the top-level selection.
The usage line, the flag set, every refusal line, and every exit stay byte-identical.

The existing form-matrix tests prove this ticket.
They include `TestProbeUsageNamesOneSelectionAndOneMutation`, `TestFormMatrixRefusesPartnerFaults`, `TestCombinedEditsRefuseSingleForms`, and `TestProbeSelectsTheCheckForm`.
Change a test file on this line only when a moved helper name forces the change.
The command registry files join this `Writes:` line through the binding registry, and this ticket expects no edit there.

## Acceptance

- [ ] The probe's help inventory line is unchanged, and `TestHelpInventoryIsComplete` stays green.
- [ ] Every existing form-matrix refusal keeps its line and its exit 2 through the two halves, and the form-matrix tests stay green.
- [ ] `--expect` with `--check prose` still prints the usage line at exit 2, as the prose case of `TestExpectRequiresTheNamedWitness` shows.
- [ ] Review confirms that `internal/probe/command.go` declares each mutation flag once.
