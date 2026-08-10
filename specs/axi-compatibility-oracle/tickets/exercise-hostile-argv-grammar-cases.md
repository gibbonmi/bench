# Exercise the hostile argv grammar cases

Blocked by: compare-four-observations.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: four-observation comparator and candidate rebuild→compare-four-observations.md; hostile argv fixtures→`specs/axi-compatibility-oracle/testdata`; hostile-class derivation→`internal/axi/compatibility`; shared grammar owner→`internal/usage/parse.go` exercised unchanged by the HG1 flag and positional rows and by HG1/exit-0-help; root and family-home dispatchers→`cmd/bench/command_registry.go` exercised unchanged by HG1/exit-2-invalid-argv, `internal/spec/build.go` by HG1/nested-help, and `internal/publication/command.go` by HG1/exit-1-unsatisfied; the `root-unknown` and `nested-spec-build-home` census members→compare-four-observations.md; environment and process consumer→exercise-hostile-environment-and-process-cases.md
Contracts: the hostile argv case IDs `<member>-hostile-<class>` cross `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`→`specs/axi-compatibility-oracle/testdata`; their type is one baseline observation record per case ID, membership is the thirteen grammar classes this ticket derives from that inventory over every member whose class matrix marks `R`, order is stable case ID ascending, and a member marked `R` with no hostile case refuses the load; asserted by HG1 against the really rebuilt candidate executable and the real `usage.Parse` callers
Closure: HG1/root-help, HG1/wrapper-help, HG1/nested-help, HG1/unknown-flag, HG1/repeated-flag, HG1/missing-flag-value, HG1/empty-required-value, HG1/excess-positional, HG1/double-dash-terminator, HG1/no-op, HG1/exit-0-help, HG1/exit-1-unsatisfied, HG1/exit-2-invalid-argv

## What to build

Rejection is public language. This ticket derives one case per declared grammar
class for every member whose class matrix marks `R`, and compares all four
observations for each: the three help surfaces (a root command's own help, the
wrapper's `help|--help|-h` body, and a nested family home such as `bench spec
build` with no operation), the five malformed-argv classes `usage.Parse` keeps
distinct (unknown flag, repeated flag, missing flag value, empty required value,
excess positional), the `--` terminator, the no-op result, and the three exit
values the taxonomy assigns.

Each exit value is its own fact with its own producer, so each gets its own case
and its own row: `usage.Parse` returns the declared help body with code 0,
`runPrepare` in `internal/publication/command.go` returns 1 for an intent it
could not satisfy, and `Command.Run` returns 2 for argv the dispatcher will not
accept. One mutation cannot grade all three — a candidate that keeps help at 0
can still collapse 1 into 2.

The classes must stay distinct. Collapsing a missing flag value into an unknown
flag, or an invalid-argv 2 into an unsatisfied-intent 1, still produces a
plausible refusal — so each has its own case and its own mutation row below.

Mutations are applied to a scratch copy of the tree from which the candidate
executable is rebuilt through `scripts/go-build.sh`. The rebuild runs under a 180s
`context.WithTimeout` and every hostile case child under a 30s deadline, so a
candidate that waits on stdin instead of refusing fails as a bounded deadline
rather than hanging the gate.

## Acceptance

- [ ] [HG1] (covers CO6) every declared hostile argv class compares byte-exact on stdout, stderr, exit, and acceptance for each member whose class matrix marks it, and a candidate rebuilt with any one class relaxed or reclassified reports the delta on the case that owns it.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| HG1/root-help | make one root command's `--help` return its default success body instead of its help body in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/root_help -timeout 900s`; it fails at the raw stdout equality assertion for case `root-outline-hostile-help`, reporting the row table where the baseline carries the usage body; the rebuild is bounded at 180s and each child at 30s |
| HG1/wrapper-help | drop one line from the wrapper's help, `--help`, and `-h` heredoc in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/wrapper_help -timeout 900s`; it fails at the raw stdout equality assertion for case `wrapper-help-hostile-help`, naming the missing `bench commit -m <msg> <path>...` line; bounded by the 180s rebuild and 30s child deadlines |
| HG1/nested-help | answer `bench spec build` with no operation using exit 0 and an empty body in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/nested_help -timeout 900s`; it fails at both the stdout equality assertion for case `nested-spec-build-home-hostile-help`, which loses the nine-operation list `usage: bench spec build (missing argument: start|assign|checkpoint|integrate|review|status|promote|abandon|reclaim)` that `bench spec build` prints for a missing operation, and the exit assertion reporting 0 against the baseline 2; bounded by the 180s rebuild and 30s child deadlines |
| HG1/unknown-flag | accept an undeclared flag by ignoring it in `usage.Parse` in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/unknown_flag -timeout 900s`; it fails at the acceptance assertion for case `root-coverage-hostile-unknown-flag`, reporting accepted with exit 0 where the baseline rejects with the unknown-flag line and exit 2; bounded by the 180s rebuild and 30s child deadlines |
| HG1/repeated-flag | accept a repeated flag by keeping its last value in `usage.Parse` in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/repeated_flag -timeout 900s`; it fails at the acceptance assertion for case `nested-spec-build-assign-hostile-repeated-flag`, reporting accepted where the baseline emits the repeated-flag refusal; bounded by the 180s rebuild and 30s child deadlines |
| HG1/missing-flag-value | classify a flag given no value as an unknown flag in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/missing_flag_value -timeout 900s`; it fails at the raw stdout equality assertion for case `nested-spec-build-checkpoint-hostile-missing-value`, printing the unknown-flag line against the baseline's missing-value line at the same exit 2 — the exit alone cannot catch this; bounded by the 180s rebuild and 30s child deadlines |
| HG1/empty-required-value | accept an empty string for a `NoEmptyValue` flag in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/empty_required_value -timeout 900s`; it fails at the acceptance assertion for case `nested-spec-build-assign-hostile-empty-value`, reporting accepted for `--ticket ""` where the baseline refuses; bounded by the 180s rebuild and 30s child deadlines |
| HG1/excess-positional | ignore positional arguments past the declared `MaxArgs` in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/excess_positional -timeout 900s`; it fails at the raw stdout equality assertion for case `nested-spec-build-status-hostile-excess-positional`, losing the refusal line that names the first excess argument; bounded by the 180s rebuild and 30s child deadlines |
| HG1/double-dash-terminator | treat a `--` terminator as an ordinary positional argument in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/double_dash -timeout 900s`; it fails at the raw stdout equality assertion for case `nested-worktree-exec-hostile-double-dash`, reporting the child argv parsed one element short; bounded by the 180s rebuild and 30s child deadlines |
| HG1/no-op | report a no-op as an unsatisfied intent with exit 1 in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/no_op -timeout 900s`; it fails at the exit assertion for case `nested-worktree-release-hostile-no-op`, reporting 1 against the baseline 0 with byte-identical stdout; bounded by the 180s rebuild and 30s child deadlines |
| HG1/exit-0-help | return code 2 alongside the declared help body from the help outcome of `usage.Parse` in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/exit_zero_help -timeout 900s`; it fails at the exit assertion for case `nested-spec-build-status-hostile-help`, reporting 2 against the baseline 0 while stdout stays the byte-identical `usage: bench spec build status <slug> [--full]` line; bounded by the 180s rebuild and 30s child deadlines |
| HG1/exit-1-unsatisfied | return exit 2 from `runPrepare`'s unsatisfied-intent branch in `internal/publication/command.go` in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/exit_one_unsatisfied -timeout 900s`; it fails at the exit assertion for case `nested-release-prepare-hostile-unsatisfied`, run against a root with no approved release set, reporting 2 against the baseline 1 while the stdout `error: unsatisfied release intent — <reason>` line stays identical; bounded by the 180s rebuild and 30s child deadlines |
| HG1/exit-2-invalid-argv | return exit 1 from the unknown-subcommand arm of `Command.Run` in the candidate rebuild | the hostile-grammar test | run `go test ./cmd/bench -run TestHostileArgvObservationsAreExact/exit_two_invalid_argv -timeout 900s`; it fails at the exit assertion for case `root-unknown-hostile-refusal`, reporting 1 against the baseline 2 while the stderr line `bench: unknown subcommand: "nope"` stays identical; bounded by the 180s rebuild and 30s child deadlines |
