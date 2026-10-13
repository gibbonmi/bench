# Import one mutant row into a verification entry

Blocked by: 4-grade-a-mutant-set.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/record_operations_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/probe.go, internal/probe/probetest/fixture.go, internal/probe/result.go, internal/probe/result_test.go, internal/probe/set.go (new), internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/verification.go, internal/reviewrecord/recordcmd/verification_import_test.go, internal/reviewrecord/recordcmd/verification_test.go
Covers: PM34, PM42, PM52, PM53, PM54, PM55, PM56, PM57, PM58, PM59, PM60, PM74

## What to build

`bench record verification --requirement <id> --probe-result <file> --mutant <n>` imports row n of a mutant-set answer as one verification entry.
The import binds `--exit-code` to the row's `exit_code` and the probe outcome to its `verdict`.
It binds the probe exit to the row's `focused_exit_code` and the restore to its `restored` cell.

`internal/probe` exports `ReadMutant(data []byte, n int)`, which returns a `MutantResult` with `Verdict`, `Restored`, `RowExit`, and `FocusedExit`.
`ReadResult` refuses a mutant-set answer with the exported sentinel `ErrMutantSet`.
`ReadMutant` refuses a single-form answer with the exported sentinel `ErrSingleForm`.
The two readers share today's completeness, spill-frame, strict-decode, and canonical-render rules through one code path.

The refusal lines are the exact lines of the spec's Record import section.
The probe package keeps the `--mutant` wording out of its own errors, and `recordcmd` maps each sentinel to its line.
The help text names each set field through a probe-owned field-source function, as `probeSources` does for the single form.
Build each record row's answer through the probe owner's renderer or a probe run, not a hand-typed envelope.
The command registry files join this `Writes:` line through the binding registry.

## Acceptance

- [ ] `probe.ReadResult` of a mutant-set answer returns `ErrMutantSet`.
- [ ] `--mutant 2` of a two-row answer writes one entry whose probe outcome is row 2's verdict and whose probe exit is row 2's focused exit.
- [ ] `--exit-code 1` against a row whose `exit_code` is 0 refuses with a line that names `mutants.exit_code`.
- [ ] A mutant-set answer without `--mutant`, and `--mutant` with a single-form answer, each refuse with their exact designed lines.
- [ ] `--mutant 3` of a two-row answer refuses with the line that names 2 rows.
- [ ] `--mutant 0`, `--mutant x`, and `--mutant` without `--probe-result` each print the form usage at exit 2.
- [ ] `ReadMutant` refuses a cut answer, a spill frame, a non-canonical row, and a repeated `mutants` table.
- [ ] The help inventory's record verification line names `--mutant <n>` and `mutants.exit_code`.
