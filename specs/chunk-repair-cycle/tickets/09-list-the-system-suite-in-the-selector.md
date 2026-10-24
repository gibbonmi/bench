# List the system suite in the changed-package selector and the floor

Blocked by: 08-prepare-the-checkpoint-gate.md
Writes: internal/testreport/selection.go, internal/testreport/selection_test.go, internal/testreport/selection_cause_test.go, internal/testreport/selection_once_test.go, internal/testreport/selected_by_test.go, internal/testreport/floor.go, internal/testreport/floor_test.go, internal/testreport/command.go, internal/testreport/testreport.go, internal/testreport/named_check.go, internal/testreport/cancel_test.go, internal/testreport/lifetime_test.go, internal/gate/phases.go, internal/preflight/check_floor.go, internal/preflight/check_floor_test.go, internal/preflight/check_floor_command_test.go, internal/preflight/system_tag.go, internal/reviewrecord/recordcmd/fold.go, internal/reviewrecord/fold.go, internal/reviewrecord/fold_test.go, internal/conformance/process_lifetime_test.go, internal/systemtest/process_lifetime_callers_test.go
Covers: RC53, RC54, RC55, RC56, RC57

## What to build

Review chunk: RC-C5. This chunk lands before the delta checkpoint (RC-C6), as the FT445 precondition requires.

A changed system-suite test file, and a production change that the system suite imports, select `bench test --check system`. The floor of a ticket that writes the system suite names the same check.

The changed-package index loads `./...` with no tags (`currentPackagesWithEnvironment`, `internal/testreport/selection.go:173-201`), so it omits `internal/systemtest`, and a changed system-suite file refuses as "not in a current package" (`selection.go:235`). Add the system-suite package to the index from `gate.SystemSuite` (`internal/gate/phases.go:85-87`), which returns the tag operand and the package operand. `testreport` already imports `gate`, so this adds no import edge. Add the package only where `gate.SystemSuiteRuns` holds, so a linked repository without the suite keeps today's index. No second list of the tag or the package enters the selector (RC57). `internal/preflight/system_tag.go:21-47` already derives system tagging for the kit-pin check; if the selector needs the same fact, route both through one source. `phases.go` is at 394 of 400 lines, so an accessor there must stay within 6 lines.

A selected system-suite package maps to `bench test --check system`, not to a package operand of `go test`. Today the selected packages go to `go test` as operands (`runSelected`, `internal/testreport/command.go:213-251`), and `command.go` is at 389 of 400 lines. The `selected_by` column renders the cause (`internal/testreport/testreport.go:169-178`).

The floor reads the same index (`WritesFloor`, `internal/testreport/floor.go:24-28`), so RC55 and RC56 follow from the index. Its consumers move with it: `floorElements` and the satisfaction rules in `internal/preflight/check_floor.go:67-87,168-203` special-case only root conformance today, so a system element needs `--check system` to satisfy it and `floorCommand` to name it. Ticket 04 also writes `check_floor.go` and `check_floor_test.go`. `check_floor_test.go` is near its 400-line limit, so move an existing test to `check_floor_command_test.go` before you add rows. `foldDelta` writes `floor.Packages` into the record (`internal/reviewrecord/recordcmd/fold.go:68`, `internal/reviewrecord/fold.go:17,72`, asserted at `internal/reviewrecord/recordcmd/fold_test.go:160-169`, which compares against `WritesFloor` and needs no edit).

If the index needs a second `go list` child for the tagged package, these pins move with it: the process-lifetime census row (`internal/conformance/process_lifetime_test.go:49`), the fake `go` of `writeChangedSubjectGo` (`selection_test.go:590-601`), the fake `go` branches in `internal/testreport/cancel_test.go:161,171` and `internal/testreport/lifetime_test.go:52`, and the argv-keyed fake in `internal/systemtest/process_lifetime_callers_test.go:178-185,274`. `TestCurrentPackageGraphLoadsAllEmbedClasses` (`selection_test.go:427-434`) asserts one package on its fixture. The exact package lists at `selection_test.go:34-42,51-54,123` change only if a fixture gains a package.

The tagged system suite receives `BENCH_KIT` and the sealed run binary from `bench test --check system`; never run the tagged suite directly.

`internal/testreport/selection_test.go` is at 602 lines over its 400-line limit, and the commit lane grades growth. The planned RC53 and RC54 rows go in that file, so first move an existing self-contained test to `selection_cause_test.go` (141 lines) or `selection_once_test.go` (52 lines), so the file does not grow. Do not add a file to the crowded `internal/testreport/` directory.

## Acceptance

- [ ] RC53: a changed system-suite test file selects `bench test --check system`.
- [ ] RC54: a production change in a package that the system suite imports selects `bench test --check system`.
- [ ] RC55: the floor of a ticket that writes the system suite names `bench test --check system`.
- [ ] RC56: the floor and the changed selection agree on a system-suite path (`TestWritesFloorAgreesWithChangedSelection`).
- [ ] RC57: `bench consumers gate.SystemSuite` lists the index as a consumer, and RC53 derives its expected operands from `gate.SystemSuite`.
- [ ] A linked repository without the system suite keeps today's selection.
- [ ] The probe "drop the system-suite package from the index" bites in `bench test --package ./internal/testreport`.
