# Name a run-binary or mutant build failure in its own outcome

Blocked by: 1-add-probe-mutation-forms.md, 2-probe-prose-and-system-checks.md, 3-name-the-first-failure-and-expected-witness.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/probe.go, internal/probe/subject_check_fixture_test.go, internal/probe/verdict.go (new), internal/probe/verdict_test.go (new), internal/runbinary/build_script_test.go (new), internal/runbinary/runbinary.go, internal/runbinary/runbinary_test.go, internal/testreport/canned_go_test.go (new), internal/testreport/command.go, internal/testreport/full_failure_test.go, internal/testreport/outcome.go, internal/testreport/outcome_test.go, internal/testreport/testreport.go
Covers: PS44, PS45, PS46, PS47, PS48, PS49, PS50, PS51, PS52, PS53

## What to build

A mutant that does not compile reads `build-failed`, and a `mutant` table names its first build line.
A failed run-binary build script reads `build-failed` too, and the build error text stays byte-identical.

Ticket 3 supplies `verdict.go`, the v2 envelope, and the field-wise `Outcome` helper.

`runbinary` builds a typed `BuildError` only on the exit-failure branch of `runBuildScript`.
The type carries the builder output, and its `Error()` text equals the current `<exit error>: <output>` bytes.
A cancelled build and an absent Go keep their current untyped errors.
Every reader of `Factory.Own` therefore prints the same bytes.

`testreport` maps a `BuildError` to the outcome kind `build-failed` and keeps its current refusal text.
Every other selection failure stays `refused`, and the `bench test` output and exit stay unchanged.
The `OutcomeBuildFailed` doc comment in `outcome.go` names both sources of that kind.

`testreport.Outcome` gains `BuildLine`.
For a Go test build failure, `BuildLine` is the first diagnostic line that the failures table prints for that package.
For a `BuildError`, `BuildLine` is the first builder output line that contains `.go:`, or else the first non-empty line.

`verdict.go` renders a `mutant[1]{first_build_line}` table after `probe_execution` when the baseline passed and the mutated run reads `build-failed`.
The table sits inside the counted suffix, so `trailing_bytes` includes it.
The cell passes through `sanitize.Preview` before the encoder.
A baseline whose run-binary build script fails now reads the cause `baseline-build-failed`.
`BuildLine` rides on `Outcome`, so the `render` signature stays unchanged, and the direct call at `internal/probe/outcome_test.go:266` needs no edit.

`internal/runbinary/runbinary_test.go` has 391 lines, so the PS44 and PS49 test needs headroom.
Move the build-script tests that follow `TestBuildEnvironmentCarriesTheBenchBuildCache` to `internal/runbinary/build_script_test.go` first.
`internal/testreport/outcome_test.go` passes 400 lines after tickets 3 and 4, and the growth ratchet reds it.
Move its canned-Go helpers, `cannedSets` through `cannedSelections` before `TestExecuteClassifiesTheOutcome`, to `internal/testreport/canned_go_test.go` first.
Each move keeps every assertion unchanged.

The full-failure table test and `TestExecuteClassifiesTheOutcome` gain `BuildLine` in their build-failure expectations.
The command registry files join this `Writes:` line through the binding registry, and this ticket does not change them.

## Acceptance

- [ ] A fixture `scripts/go-build.sh` that exits 1 gives `Factory.Own` an error that `errors.As` reads as `BuildError` with the output.
- [ ] `BuildError.Error()` equals the bytes of the current exit-status wrap over the same output.
- [ ] `testreport.Execute` over that failing script returns kind `build-failed` with its current refusal bytes.
- [ ] A build cancelled through its context, a build with no Go on `PATH`, and a factory `Build` with a plain error each stay `refused`.
- [ ] A `BuildError` output with a header before a `.go:` line gives `BuildLine` the `.go:` line.
- [ ] A `BuildError` output with no `.go:` line gives `BuildLine` its first non-empty line.
- [ ] A canned Go test build failure for an unused import prints the `mutant` table with that compiler line.
- [ ] A baseline whose run-binary build script fails reads the cause `baseline-build-failed`.
- [ ] The gate and lane readers of `Factory.Own` print unchanged error bytes.
- [ ] Each changed or new file stays at 400 lines or fewer, and each moved test keeps its assertions.
