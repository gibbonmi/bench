# Carry the partition in the verdict and refuse its reuse

Blocked by: Record a component ancestor slot as its own class
Ownership fence: `internal/gate/verdict.go`, `internal/gate/reduced_verdict_test.go`, `internal/gate/verdict_reuse_test.go`
Assumptions: `verdict.go` today validates exactly two ready classes — full and
reduced — through `fullReadyFields`/`reducedReadyFields` and `inherits()`, and
`inspectAt` withholds `ReusableGreen` from a reduced record. Re-derive from the
tree at pickup.

## What to build

The verdict record generalizes from "reduced" to a per-component partition: the
executed set, and per skipped component the evidence that covered it — an
ancestor identity and recorded time, or the seal's source digest for build. The
existing two-class discipline holds: exact field sets, alternatives rather than
a spectrum, a fragment refused for what it is missing rather than read as
something else.

The reuse guard generalizes with it. A partial verdict is evidence about its own
tree and never a whole-tree green, so `inspectAt` withholds `ReusableGreen` from
it exactly as it does from a reduced record. That one site is what makes
`bench commit` safe without a commit-side edit: commit reads reusability through
`authorization`, which reads it here.

`Inspection` gains the partition so the consumers migrating next — status and
prep-release — can render and refuse against it. This ticket constructs partial
records in tests; nothing yet emits one.

## Acceptance

- [ ] [PC2a] a partial record round-trips its executed set and one evidence entry per skipped component, and re-reads byte-equal.
- [ ] [PC2b] a partial record missing an evidence entry for a component it lists as skipped is refused.
- [ ] [PC15] a partial green is never `ReusableGreen`, and the inspection names the partition as its reason.
- [ ] [PS23] a record mixing partial fields with reduced or full field sets is refused rather than resolved to either class.
- [ ] [PS24] the existing full and reduced classes still validate and still round-trip unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC2a | serialize the skipped set as names only, dropping the evidence | `TestPartialVerdictRoundTrips` | write a partial record, load it, compare the loaded partition field by field |
| PC2b | validate the evidence map's shape without cross-checking it against the skipped set | `TestPartialVerdictRequiresEvidencePerSkip` | write a partial record listing a skip with no evidence entry, load it, expect invalid |
| PC15 | check only `rec.Reduced` in the reuse guard | `TestPartialVerdictIsNotReusable` | write a partial green for the current tree, call `Inspect`, assert `ReusableGreen` false and the reason names the partition |
| PS23 | relax the exact-field-set comparison to a subset check | `TestMixedClassRecordRefuses` | write a record carrying both `ancestor` and the partition field, load it, expect invalid (review Sp2: the earlier swap-the-expected-set mutation cannot red this owner — a mixed set matches neither class exactly under either expectation) |
| PS24 | replace `reducedReadyFields` with the partial set | existing reduced-verdict round-trip tests | run the existing reduced-class cases unchanged |
