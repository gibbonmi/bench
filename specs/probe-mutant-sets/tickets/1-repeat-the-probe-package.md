# Repeat the probe package in one call

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/main.go, cmd/bench/test_command.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/command.go, internal/probe/outcome_test.go, internal/probe/package_test.go (new), internal/probe/probe.go, internal/probe/probetest/fixture.go, internal/probe/refusal_test.go, internal/probe/subject_check_fixture_test.go, internal/testreport/command.go, internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go
Covers: PM45, PM46, PM47, PM48, PM49, PM50

## What to build

A single-form `bench probe` takes `--package` more than once.
Every named package expression runs in one `go test` call for each run, in argv order.
Each run keeps one private build, so a two-package probe starts two Go children and builds twice.

Mark `--package` as repeatable in the probe grammar in `internal/probe/command.go`.
Change the usage line to `--package <expr> [--package <expr>]...`.
`selectionArgs` in `internal/probe/probe.go` passes each value of `parsed.Repeated["--package"]` in argv order.
A read of `parsed.Flags` keeps only the last value, so that read must go.

`probePackageNote` in `internal/testreport/selection_facts.go` retires.
`ProbeNotes` states `PackageExpressionNote` first and the check note second, because the two grammars now take the same repeat form.
Update the independent note copy in `internal/probe/refusal_test.go` and the probe line in `cmd/bench/help_inventory_test.go`.

Put the new rows in `internal/probe/package_test.go`.
Drive them with the canned Go runner and the build-records seam of `internal/probe/subject_check_fixture_test.go`.
The command registry files join this `Writes:` line through the binding registry, and this ticket expects no edit there.

## Acceptance

- [ ] A probe with `--package ./clamp --package ./other` passes both expressions in argv order to each `go test` start.
- [ ] That probe starts exactly two `go test` children and records one build for each run.
- [ ] That probe's `target` cell reads `./clamp ./other`.
- [ ] Every help spelling prints the usage line, `PackageExpressionNote`, and the check note, and no help text says "run one call per package".
- [ ] `ProbeNotes` starts with `PackageExpressionNote` and holds two facts.
- [ ] The help inventory's probe line names `[--package <expr>]...`.
