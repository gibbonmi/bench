# Capture the pinned baseline independently

Blocked by: authenticate-baseline-manifest.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: authenticated manifest and pinned-subject constant→authenticate-baseline-manifest.md; capture and provenance API→`internal/axi/compatibility`; expected-observation fixtures and their provenance digest→`specs/axi-compatibility-oracle/testdata`; paired build of both subjects→`scripts/go-build.sh` exercised unchanged by the BC1/build-* and BC1/distinct-executable rows; comparator consumer→compare-four-observations.md
Contracts: the provenance record — absolute baseline executable path, its `freshness.SealDigests` executable digest, the pinned source subject, and the digest of the expected observation bytes — crosses `cmd/bench/axi_compatibility_test.go`→`specs/axi-compatibility-oracle/testdata`; its type is one record per captured case, membership is baseline-authored cases only, order is build-then-capture-then-seal, and a missing provenance record refuses before any byte comparison runs; asserted by BC1 against the two really built executables
Closure: BC1/baseline-only-authorship, BC1/distinct-executable-identity, BC1/refusal-precedes-equality, BC1/expected-bytes-immutable, BC1/digest-preimage-stdout, BC1/digest-preimage-stderr, BC1/digest-preimage-exit, BC1/digest-preimage-acceptance, BC1/digest-preimage-case-identity

## What to build

Expected observations exist only because the pinned baseline executable produced
them. `cmd/bench/axi_compatibility_test.go` builds two executables once each
through `scripts/go-build.sh` — one from a worktree checked out at
`8ae1512f95e64588487430aefa5b02c288d7de3a`, one from the candidate tree — records
each one's absolute path and `freshness.SealDigests` executable digest, and hands
only the baseline handle to the capture API. The capture API refuses any handle
whose recorded source subject is not the pinned one, and it refuses before it
compares a single byte, so a candidate that happens to agree still fails
provenance.

The expected fixtures under `specs/axi-compatibility-oracle/testdata` carry a
provenance digest over the case identity, the argv, and all four observation
fields. That digest is the value authorizing "these bytes are baseline-authored",
so the mutation rows below drop fields from its preimage as well as reroute the
capture path: a preimage missing `stderr` still passes every control-flow check
while letting a candidate rewrite stderr unnoticed.

Every build and every child run takes a `context.WithTimeout`, because a hung
child and a broken harness look identical at the gate.

## Acceptance

- [ ] [BC1] (covers CO2) expected observations load only when their provenance record names the pinned baseline executable and their digest preimage covers case identity, argv, stdout, stderr, exit, and acceptance, and the provenance refusal is reported before any byte comparison.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BC1/baseline-only-authorship | pass the candidate executable handle to the capture API that writes expected observations | the paired-capture provenance test | run `go test ./cmd/bench -run TestExpectedBytesComeOnlyFromThePinnedExecutable -timeout 600s`; it fails at the provenance assertion reporting recorded subject `<candidate HEAD>` against required `8ae1512f95e64588487430aefa5b02c288d7de3a`; both builds run under a 180s `exec.CommandContext` and each capture child under 30s |
| BC1/distinct-executable-identity | reuse the baseline executable path as the candidate path | the paired-capture provenance test | run `go test ./cmd/bench -run TestBaselineAndCandidateExecutablesAreDistinct -timeout 600s`; it fails at the executable-digest inequality assertion, printing the one `freshness.SealDigests` digest twice; bounded by the same 180s build and 30s run deadlines |
| BC1/refusal-precedes-equality | compare observation bytes first and check provenance only when they differ | the paired-capture provenance test | run `go test ./cmd/bench -run TestProvenanceRefusalPrecedesByteEquality -timeout 600s`; with candidate-authored expectations whose bytes match, it fails at the assertion that the returned error is the provenance refusal rather than `nil`; bounded by the 180s build and 30s run deadlines |
| BC1/expected-bytes-immutable | rewrite the expected fixture from the candidate observation after a comparison runs | the immutable-fixture test | run `go test ./cmd/bench -run TestExpectedFixturesAreUnchangedByACandidateRun -timeout 600s`; it fails at the assertion comparing the fixture's pre-run and post-run SHA-256, naming the rewritten fixture path; the candidate child runs under a 30s deadline |
| BC1/digest-preimage-stdout | drop the stdout bytes from the expected-fixture provenance digest preimage | the provenance-digest test | run `go test ./cmd/bench -run TestProvenanceDigestPreimageCoversEveryField/stdout -timeout 300s`; it fails at the assertion that flipping one stdout byte changes the recomputed digest, reporting the two equal digests; the recompute runs under a 30s deadline |
| BC1/digest-preimage-stderr | drop the stderr bytes from the provenance digest preimage | the provenance-digest test | run `go test ./cmd/bench -run TestProvenanceDigestPreimageCoversEveryField/stderr -timeout 300s`; it fails at the assertion that flipping one stderr byte changes the recomputed digest, reporting the two equal digests; bounded by the 30s deadline |
| BC1/digest-preimage-exit | drop the exit code from the provenance digest preimage | the provenance-digest test | run `go test ./cmd/bench -run TestProvenanceDigestPreimageCoversEveryField/exit -timeout 300s`; it fails at the assertion that changing exit 2 to 1 changes the recomputed digest, reporting the two equal digests; bounded by the 30s deadline |
| BC1/digest-preimage-acceptance | drop the accepted/rejected classification from the provenance digest preimage | the provenance-digest test | run `go test ./cmd/bench -run TestProvenanceDigestPreimageCoversEveryField/accepted -timeout 300s`; it fails at the assertion that flipping `accepted` changes the recomputed digest, reporting the two equal digests; bounded by the 30s deadline |
| BC1/digest-preimage-case-identity | drop the case ID and argv vector from the provenance digest preimage | the provenance-digest test | run `go test ./cmd/bench -run TestProvenanceDigestPreimageCoversEveryField/case_identity -timeout 300s`; it fails at the assertion that moving one record's bytes under another case ID changes the recomputed digest, reporting the two equal digests; bounded by the 30s deadline |
