# Compare all four public observations

Blocked by: capture-pinned-baseline.md, close-required-argv-classes.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`
Integration surfaces: provenance-checked baseline observations→capture-pinned-baseline.md; complete case index and class matrix→close-required-argv-classes.md; comparator and delta report→`internal/axi/compatibility` exported symbol `compatibility.Compare` (name fixed by the spec); in-process production seam→`cmd/bench/command_registry.go` (`Command.Run`) exercised unchanged by every PC1 row; exact-byte-class consumer→pin-default-full-and-empty-classes.md; hostile-case consumer→exercise-hostile-argv-grammar-cases.md
Contracts: the observation quadruple — raw stdout bytes, raw stderr bytes, integer exit, and the boolean accepted/rejected classification — crosses `cmd/bench/axi_compatibility_test.go`→`internal/axi/compatibility`; its type is one immutable record per case per run, membership is every case ID the class matrix produced, order is case ID then run number, and an absent observation refuses the comparison rather than comparing a zero value; asserted by PC1 against the really executed `Command.Run` and exact executables
Closure: PC1/stdout-bytes, PC1/stderr-bytes, PC1/stream-assignment, PC1/exit-code, PC1/acceptance-classification, PC1/second-run-required, PC1/fresh-state-isolation, PC1/absent-observation-refusal

## What to build

The comparator runs every case from the class matrix twice and reports a delta
when any of the four observations differs from the provenance-checked baseline
record. Ordinary cases go through the production `Command.Run` seam in package
`main`; cases the class matrix marks as process-identity, environment, or signal
cases go through the exact executables built in `capture-pinned-baseline.md`.

Two runs are not a formality: each run starts from a fresh temporary root, a fresh
`Command` value, and a fresh environment, and the comparator requires both runs to
be present and identical before it reports equality. A comparator that keeps a
cached first-run observation, or that treats an absent observation as an empty
one, would let a mutated candidate pass silently — so both are mutation rows
below.

`bench nope`, which the real dispatcher answers with `bench: unknown subcommand:
"nope"` on stderr and exit 2, is the case the stream, exit, and acceptance rows
use, because all three observations differ from a success case at once and each
must be caught by its own comparison.

## Acceptance

- [ ] [PC1] (covers CO4) the comparator reports a delta whenever raw stdout, raw stderr, the exit code, or the accepted/rejected classification differs from the baseline record, requires two runs from fresh state, and refuses a comparison whose observation is absent rather than empty.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1/stdout-bytes | compare stdout with `strings.TrimSpace` on both sides | the paired-comparator test | run `go test ./cmd/bench -run TestPairedRunComparesFourObservations/stdout -timeout 900s`; with a candidate observation whose only delta is a trailing space on the `bench status` case, it fails at the assertion that the comparator returned a stdout delta, reporting `no delta` where the raw byte slices differ; each `Command.Run` case runs under a 30s `context.WithTimeout` |
| PC1/stderr-bytes | compare stdout only and carry stderr through the report untested | the paired-comparator test | run `go test ./cmd/bench -run TestPairedRunComparesFourObservations/stderr -timeout 900s`; with the `bench nope` case whose candidate stderr reads `bench: unknown subcommand: "nope!"`, it fails at the assertion that a stderr delta was reported, printing both stderr byte slices; bounded by the 30s per-case deadline |
| PC1/stream-assignment | concatenate stdout and stderr into one buffer before comparing | the paired-comparator test | run `go test ./cmd/bench -run TestPairedRunComparesFourObservations/stream_assignment -timeout 900s`; with a candidate that emits the `bench nope` refusal on stdout instead of stderr while the concatenation is byte-identical, it fails at the assertion that both a stdout and a stderr delta were reported; bounded by the 30s per-case deadline |
| PC1/exit-code | report the exit code only when stdout already differs | the paired-comparator test | run `go test ./cmd/bench -run TestPairedRunComparesFourObservations/exit -timeout 900s`; with a candidate whose `bench nope` bytes are identical but whose exit is 1 instead of 2, it fails at the assertion that an exit delta was reported, printing want 2 got 1; bounded by the 30s per-case deadline |
| PC1/acceptance-classification | derive the accepted/rejected classification from the exit code instead of recording it independently | the paired-comparator test | run `go test ./cmd/bench -run TestPairedRunComparesFourObservations/accepted -timeout 900s`; with a candidate that newly accepts the previously rejected spelling `bench status --all=true` and still exits 0 with the same bytes as `bench status --all`, it fails at the assertion that an acceptance delta was reported; bounded by the 30s per-case deadline |
| PC1/second-run-required | return the comparison result after the first run when it shows no delta | the fresh-state test | run `go test ./cmd/bench -run TestPairedRunRequiresBothRuns -timeout 900s`; it fails at the assertion that each case record holds two run observations, reporting run count 1 for case `root-status-s`; each run is bounded by the 30s per-case deadline |
| PC1/fresh-state-isolation | reuse the first run's temporary root and environment for the second run | the fresh-state test | run `go test ./cmd/bench -run TestPairedRunRepeatsEachCaseFromFreshState -timeout 900s`; with the `bench gate` case whose first run writes a reusable verdict, it fails at the assertion that the second run's root is a distinct empty directory, reporting the shared path and the leaked verdict file; bounded by the 30s per-case deadline |
| PC1/absent-observation-refusal | treat an absent candidate observation as empty bytes | the paired-comparator test | run `go test ./cmd/bench -run TestPairedRunRefusesAnAbsentObservation -timeout 900s`; with a case whose candidate run was cut short before stderr was captured, it fails at the assertion that the comparator returned the absent-observation refusal naming the case ID and field `stderr`, rather than an empty-versus-empty equality; bounded by the 30s per-case deadline |
