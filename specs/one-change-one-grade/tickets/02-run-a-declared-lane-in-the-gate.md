# Run a declared lane in the gate package

Line: opus / medium.

Blocked by: 01-grade-named-markdown-paths.md
Writes: internal/gate/lane.go (new), internal/gate/lane_test.go (new), internal/gate/lane_record_test.go (new), internal/gate/manifest.go, internal/gate/phases.go, internal/gate/verdict.go, internal/gate/verdict_registry.go, internal/gate/verdict_registry_guard_test.go

## What to build

The gate package runs a lane and writes one lane record. A lane is an ordered
check list with the phase manifest's entry schema: name, argv, env, needs,
optional, and dir. The phase runner executes the lane in a private checkout of
the composed snapshot, under the gate timeout. The lane's Bench-owned checks take
the run binary the gate selects, so they grade with the tree's own code.

The kit root carries a built-in lane of exactly four checks: gofmt, prose, vet,
and build. The gofmt check runs `bench gate-go gofmt`, and the prose check runs
`bench gate-prose`, each through the run binary token. The vet check runs
`go vet ./...`, and the build check runs `go build ./...`. The prose entry carries
a placeholder for the named Markdown paths, and the lane run resolves it.

A phase manifest may declare a `lane` array beside its `phases` array. The loader
validates a lane entry as it validates a phase entry today, and it refuses a
malformed entry with the same three-part defect diagnostic. A manifest with no
`lane` array declares no lane. A root with no manifest that is not the kit root
declares no lane.

The lane writes one record to a lane file in the worktree's own Git dir. The
fields are `schema`, `tree`, `lane`, `outcome`, `run_binary`, and `recorded_at`,
and `outcome` is `pass` or `fail`. The record registers as the `lane record`
class with its own validator. The reader refuses reuse by class rather than by a
name suffix. The lane writes no gate cache record and no evidence record, and an
interrupted lane writes nothing.

The commit ticket consumes this package's lane run and its outcome, and it maps
that outcome to the two lane authorization kinds.

## Acceptance

- [ ] OG13 shows the kit root's lane argv for gofmt and prose names the run binary token.
- [ ] OG16 shows the gate cache holds no record for the composed tree after a lane run.
- [ ] OG17 shows the evidence store holds no record for the composed tree after a lane run.
- [ ] OG18 shows the record-class registry contains `lane record` with the fields `lane, outcome, recorded_at, run_binary, schema, tree`.
- [ ] OG19 shows `Inspect` on a lane record with outcome `pass` answers `ReusableGreen=false` and names the lane class.
- [ ] OG21 shows the kit root's lane table is exactly gofmt, prose, vet, and build with the profile's argv.
- [ ] OG22 shows a phase manifest with a `lane` array yields a lane table of exactly those checks.
- [ ] OG25 shows the lane record names the composed tree hash, the lane identity, and the outcome.
- [ ] OG32 shows the kit root's `test` phase argv stays `go test -count=1 ./...`.
