# Name the first failure and the expected witness in a v2 verdict row

Blocked by: 1-add-probe-mutation-forms.md, 2-probe-prose-and-system-checks.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/record_operations_test.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/baseline.go, internal/probe/baseline_test.go, internal/probe/combined_test.go (new), internal/probe/command.go, internal/probe/delete_test.go (new), internal/probe/edit_test.go (new), internal/probe/invalid_cause_test.go, internal/probe/omit_file_test.go, internal/probe/outcome_test.go, internal/probe/probe.go, internal/probe/probe_test.go, internal/probe/refusal_test.go, internal/probe/result.go (new), internal/probe/result_test.go (new), internal/probe/subject_check_test.go (new), internal/probe/unwrap_test.go, internal/probe/verdict.go (new), internal/probe/verdict_test.go (new), internal/reviewrecord/recordcmd/verification_import_test.go (new), internal/testreport/failure_lines_test.go, internal/testreport/full_failure_test.go, internal/testreport/named_check.go, internal/testreport/outcome.go, internal/testreport/outcome_test.go
Covers: PS38, PS39, PS40, PS41, PS42, PS43, PS54, PS55, PS56, PS57

## What to build

An author reads the first failed test in the verdict row, and `--expect <test>` gives `bit` only for that witness.
The result envelope moves to `bench-probe-result-v2`, and native verification import still reads it.

Tickets 1 and 2 supply the edit owner, the probe targets, and the request-derived selection cells.
NR-C3 lands `internal/probe/result.go`, `result_test.go`, `verification_import_test.go`, and `cmd/bench/record_operations_test.go`.
Those paths carry the `(new)` marker only because this tree does not hold them before that landing.

`testreport.Outcome` gains `FailedNames`, the failed Go test names in the order that the failures table prints them.
For a Go run, `FailedTests` and `FailedNames` derive from the one `report.failed()` walk.
A prose run keeps `FailedTests` as its finding count and leaves `FailedNames` empty.
`Outcome` then stops being comparable with `!=`, so its tests move to one field-wise helper.

Create `internal/probe/verdict.go` as the owner of `first_failure` and `--expect`.
The verdict row gains the cell `first_failure` after `failed_tests`, and the cell reads null when `FailedNames` is empty.
`--expect` matches a failed name that equals the value or starts with the value and `/`.
A failed mutated run without a match reads `invalid`, cause `expect-missed`, and exit 1.
`--expect` with the prose check is a usage refusal, and the probe reads the prose name from one `testreport` export.
The `proseCheck` constant in `probe.go` retires.

The schema token changes to `bench-probe-result-v2`, and `probe.ReadResult` reads the new field set.
A v1 envelope refuses by its schema.
`verdictCells` gains the new cell, and `gradeBaseline` builds its row with that cell.

The posture list reds these existing pins, and this ticket updates each one:

- The literal verdict-row pins at `internal/probe/invalid_cause_test.go:26`, `internal/probe/refusal_test.go:58`, and `internal/reviewrecord/recordcmd/verification_import_test.go:46`.
- The v1 token pins in `verification_import_test.go`.
- `TestExecuteClassifiesTheOutcome`, `TestFullFailedTestsCountsTests`, and the full-failure table test, which compare `Outcome` field-wise.
- `TestProbeHelpSpellings` and `TestHelpInventoryIsComplete`, which expect the `--expect` usage line.

The `probeRow` helper in `probe_test.go` derives every expected verdict row through `probeFields`.
Each caller file of that helper is on this `Writes:` line, the files of tickets 1 and 2 included.
`cmd/bench/record_operations_test.go` holds the real import journey, and it stays green over the v2 envelope.
The command registry files join this `Writes:` line through the binding registry, and they change only when the help projection requires it.

Ticket 4 receives `verdict.go` and the field-wise `Outcome` helper.

## Acceptance

- [ ] A bite with two failed tests names the first in package-then-test order in `first_failure`.
- [ ] A silent run, a build failure, a prose bite, and a baseline refusal each read `first_failure` null.
- [ ] `--expect TestClampNegative` reads `bit` over a failure of that test or of only `TestClampNegative/zero`.
- [ ] `--expect TestClampPositive` over a failure of only `TestClampNegative` reads `invalid`, `expect-missed`, and exit 1.
- [ ] `--expect` with `--check prose` prints the usage line with exit 2.
- [ ] The probe result schema cell reads `bench-probe-result-v2`.
- [ ] `probe.ReadResult` reads a v2 envelope and refuses a v1 envelope by its schema.
- [ ] `FailedNames` lists the failed tests in the failures-table order, and a prose run leaves it empty.
- [ ] The help inventory's probe line names `--expect`.
- [ ] The native verification import tests and the record journey stay green over the v2 envelope.
- [ ] Each changed file stays at 400 lines or fewer, and no test keeps a second derivation of the verdict row.
