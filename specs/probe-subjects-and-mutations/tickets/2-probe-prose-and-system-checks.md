# Probe the prose check and the system suite through the check owner

Blocked by: 1-add-probe-mutation-forms.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/command.go, internal/probe/outcome_test.go, internal/probe/probe.go, internal/probe/refusal_test.go, internal/probe/replacement_test.go, internal/probe/subject_check_fixture_test.go (new), internal/probe/subject_check_test.go (new), internal/testreport/command.go, internal/testreport/runbinary_test.go, internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go
Covers: PS26, PS27, PS28, PS29, PS30, PS31, PS32, PS33, PS34, PS35, PS36, PS37, PS77

## What to build

An author probes a Markdown file under `--check prose` and a system row under `--check system --run <pattern>`.
Both runs go through `testreport.Execute`, and the probe needs no hand route and no rebuild.

Ticket 1 supplies the edit owner, the closed form matrix, and `probe.Usage`.
This ticket changes the selection grammar in `command.go` and the run sequence in `probe.go` after that ticket commits.

Remove `gradeCheckTarget` and its refusal from `probe.go`, and retire `testreport.HandProbeRoute`.
`testreport.ProbeNotes` keeps `probePackageNote`, the probe's single-`--package` expression note, and the `--check` inventory note without the not-a-target clause.
`PackageExpressionNote` stays the `bench test` note, because the probe grammar takes one `--package`.
The rebuild note leaves `ProbeNotes`, because each focused run without an inherited selection builds its own run binary.

`--run` becomes valid with `--package` or with `--check system`, and the usage line names both positions.
`--run` with any other check stays a usage refusal with exit 2.
`selectionArgs` passes `--run` to `testreport.Prepare` for both forms.

The probe parses the selection before it answers a substring miss.
The substring-miss answer reads its form, target, and run cells from the prepared request.
That request is the one source of those cells, so an unknown check wins over a substring miss.
Update the refusal order in the `Command` doc comment to match the spec.

PS33 copies a sealed executable pair from a fixture build script.
The fixture meets the four conditions in the spec's further notes, and `BENCH_RUN_BINARY` is unset, not empty.
`newFixture` sets that variable, so the PS33 test removes it for its own run.

`TestProbeRefusesProseAndSystemChecks` retires, because both checks become probe targets.
`TestProbeNotesDeriveTheirTokens` expects two notes.
The independent `probeHelpNotes` list at `refusal_test.go:291` loses the not-a-target clause, the hand route, and the rebuild note.
The literal expected probe line in `TestHelpInventoryIsComplete` (`cmd/bench/help_inventory_test.go`) changes to name both `--run` positions.

Put the fixture builder of PS33 in `internal/probe/subject_check_fixture_test.go`, so `subject_check_test.go` stays at 400 lines or fewer.
The command registry files join this `Writes:` line through the binding registry, and they change only when the help projection requires it.

Ticket 3 receives probe targets that include a prose run, which leaves the failed-name list empty.

## Acceptance

- [ ] A Markdown mutation under `--check prose` that adds a long sentence reads `bit`.
- [ ] A `--check prose` probe of a subject that `prose.IsSubjectName` rejects refuses with the not-graded line and exit 1.
- [ ] The `failed_tests` cell of that prose bite equals the count of prose findings over the mutant.
- [ ] `--check system --run <pattern>` passes the baseline and runs the mutant through the system-suite argv.
- [ ] A system probe over a root that is not the kit reads `invalid` with the cause `baseline-refused`.
- [ ] The probe request for `--check system --run ^X$` equals `testreport.Prepare` of the same selection.
- [ ] `--run` with `--check line-routing` stays the usage refusal.
- [ ] Two `Execute` calls without an inherited selection call the builder twice.
- [ ] The build script of the mutated run reads the mutant bytes when no selection is inherited.
- [ ] `testreport.ProbeNotes` holds exactly two notes, and the probe help names no copy-aside route and no not-a-target sentence.
- [ ] A substring miss with an unknown check answers the unknown-check refusal with exit 2.
- [ ] A check-form substring miss names the named-check run pattern in its run cell.
- [ ] Each new test file stays at 400 lines or fewer.
- [ ] Every remaining probe and `testreport` assertion stays green.
