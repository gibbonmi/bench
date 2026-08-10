# Authenticate the baseline manifest

Blocked by: none
Ownership fence: `internal/axi/compatibility`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: manifest reader, fixture schema, and pinned-subject constant→`internal/axi/compatibility`; hostile and well-formed manifest fixtures→`specs/axi-compatibility-oracle/testdata`; canonical build-input digest→`internal/freshness/freshness.go` (`Digest`, `BuildInputs`, `SealDigests`) exercised unchanged by the BM1/preimage-* rows; registry census consumer→derive-root-registry-membership.md; paired capture consumer→capture-pinned-baseline.md
Contracts: the manifest record — pinned source subject, canonical-builder seal, stable case ID, argv vector, fixture path, and the raw stdout/stderr/exit/acceptance quadruple — crosses `specs/axi-compatibility-oracle/testdata`→`internal/axi/compatibility`; its type is one TOON record per case, membership is exactly the cases the manifest lists, order is stable case ID ascending, and any absent field refuses the whole load rather than defaulting; asserted by BM1 against the real reader over real on-disk fixtures
Closure: BM1/subject, BM1/preimage-build-input-set, BM1/preimage-name-framing, BM1/preimage-length-framing, BM1/preimage-order, BM1/preimage-executable-digest, BM1/fifo-fixture, BM1/directory-fixture, BM1/dangling-symlink-fixture, BM1/symlinked-regular-fixture, BM1/unique-id, BM1/observation-stdout, BM1/observation-stderr, BM1/observation-exit, BM1/observation-acceptance

## What to build

`internal/axi/compatibility` loads a baseline manifest and refuses everything that
would let a later comparison grade itself. A load succeeds only when the record
pins source subject `974020e4af8de5ed75098c4c5934a8907952bb2b`, carries a
canonical-builder seal whose preimage is recomputed from the real build inputs
`freshness.BuildInputs` returns (`scripts/go-build.sh`, `package.json`,
`internal/releaseevidence/requirements.json` and the rest of that closure) framed
exactly as `freshness.Digest` frames them, names a fixture that `lstat` reports as
a regular non-symlink file, uses case IDs unique across the manifest, and carries
all four observations — raw stdout, raw stderr, integer exit, and the
accepted/rejected classification — for every case.

The refusal is total: one bad record refuses the load and names the offending case
ID, so a later ticket can never read a partially validated index. Every reader
entry point takes a context and the tests bound it, because the hostile fixtures
include a FIFO whose open would otherwise block forever.

The seal is the value that authorizes the whole manifest, so the mutation rows
below break its *inputs* — dropping a build input, dropping the name or length
framing, reordering the closure, swapping the recorded executable digest — rather
than only bypassing the check that reads it.

## Acceptance

- [ ] [BM1] (covers CO1) the manifest reader loads a manifest only when the pinned subject, the recomputed canonical-builder seal, every fixture's regular non-symlink file kind, case-ID uniqueness, and all four observations per case hold, and names the offending case ID otherwise.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BM1/subject | change the manifest's `source_subject` to `0000000000000000000000000000000000000000` while every other field stays valid | the manifest-reader test | run `go test ./internal/axi/compatibility -run TestManifestRejectsUnpinnedSubject -timeout 60s`; it fails at the subject-equality assertion reporting want `974020e4af8de5ed75098c4c5934a8907952bb2b`, got the zero subject; the reader's 5s `context.WithTimeout` bounds the load so a stalled read fails as a deadline |
| BM1/preimage-build-input-set | drop `scripts/go-build.sh` from the path list the seal preimage hashes, keeping the seal check itself in place | the seal-preimage test | run `go test ./internal/axi/compatibility -run TestManifestSealPreimageCoversEveryBuildInput -timeout 60s`; it fails at the set-equality assertion between the preimage path list and `freshness.BuildInputs(root)`, naming the missing `scripts/go-build.sh`; the input enumeration runs under the same 5s deadline |
| BM1/preimage-name-framing | hash each build input's contents without its slash-relative name | the seal-preimage test | run `go test ./internal/axi/compatibility -run TestManifestSealPreimageIsNameFramed -timeout 60s`; it fails at the digest-equality assertion against `freshness.Digest(root)` after two same-size inputs are swapped, reporting the two digests; the recompute runs under the 5s deadline |
| BM1/preimage-length-framing | drop the `%d:` length prefixes so name and contents concatenate unframed | the seal-preimage test | run `go test ./internal/axi/compatibility -run TestManifestSealPreimageIsLengthFramed -timeout 60s`; it fails at the digest-equality assertion against `freshness.Digest(root)` for the crafted boundary-collision pair, reporting both digests; bounded by the 5s deadline |
| BM1/preimage-order | hash the build-input closure in filesystem-walk order instead of the sorted order `freshness.BuildInputs` returns | the seal-preimage test | run `go test ./internal/axi/compatibility -run TestManifestSealPreimageIsSorted -timeout 60s`; it fails at the digest-equality assertion against `freshness.Digest(root)` on the shuffled closure, reporting both digests; bounded by the 5s deadline |
| BM1/preimage-executable-digest | record the staged build output's digest instead of the digest `freshness.SealDigests` reports for the pinned executable | the seal-authority test | run `go test ./internal/axi/compatibility -run TestManifestSealBindsThePublishedExecutableDigest -timeout 60s`; it fails at the executable-digest equality assertion, reporting the recorded digest against the sealed one; bounded by the 5s deadline |
| BM1/fifo-fixture | point one case's fixture path at a FIFO created by `syscall.Mkfifo` | the hostile-fixture test | run `go test ./internal/axi/compatibility -run TestManifestRejectsNonRegularFixture -timeout 60s`; it fails at the file-kind assertion expecting the `not a regular file` refusal naming the case ID; the reader `lstat`s before opening and the 5s deadline turns a regression that opens the FIFO into a bounded deadline failure rather than a hang |
| BM1/directory-fixture | point one case's fixture path at a directory | the hostile-fixture test | run `go test ./internal/axi/compatibility -run TestManifestRejectsDirectoryFixture -timeout 60s`; it fails at the file-kind assertion expecting the `not a regular file` refusal naming the case ID; bounded by the 5s deadline |
| BM1/dangling-symlink-fixture | point one case's fixture path at a symlink whose target does not exist | the hostile-fixture test | run `go test ./internal/axi/compatibility -run TestManifestRejectsDanglingSymlinkFixture -timeout 60s`; it fails at the symlink assertion expecting the `symlinked fixture` refusal before target resolution, naming the case ID; bounded by the 5s deadline |
| BM1/symlinked-regular-fixture | point one case's fixture path at a symlink whose target is the real regular fixture | the hostile-fixture test | run `go test ./internal/axi/compatibility -run TestManifestRejectsSymlinkedRegularFixture -timeout 60s`; it fails at the symlink assertion — a resolved-then-`stat` regression reads the same bytes and would otherwise pass — reporting the case ID and the link path; bounded by the 5s deadline |
| BM1/unique-id | give two manifest records the same case ID `status-default` with different argv | the manifest-index test | run `go test ./internal/axi/compatibility -run TestManifestRejectsDuplicateCaseID -timeout 60s`; it fails at the duplicate-ID assertion naming `status-default` and both record positions; bounded by the 5s deadline |
| BM1/observation-stdout | remove the `stdout` field from one record | the manifest-completeness test | run `go test ./internal/axi/compatibility -run TestManifestRequiresEveryObservation/stdout -timeout 60s`; it fails at the missing-field assertion naming field `stdout` and the case ID, not at a zero-value comparison; bounded by the 5s deadline |
| BM1/observation-stderr | remove the `stderr` field from one record | the manifest-completeness test | run `go test ./internal/axi/compatibility -run TestManifestRequiresEveryObservation/stderr -timeout 60s`; it fails at the missing-field assertion naming field `stderr` and the case ID — an empty stderr is a legal value, so the record must distinguish absent from empty; bounded by the 5s deadline |
| BM1/observation-exit | remove the `exit` field from one record | the manifest-completeness test | run `go test ./internal/axi/compatibility -run TestManifestRequiresEveryObservation/exit -timeout 60s`; it fails at the missing-field assertion naming field `exit` and the case ID — absent must not read as exit 0; bounded by the 5s deadline |
| BM1/observation-acceptance | remove the `accepted` field from one record | the manifest-completeness test | run `go test ./internal/axi/compatibility -run TestManifestRequiresEveryObservation/accepted -timeout 60s`; it fails at the missing-field assertion naming field `accepted` and the case ID — absent must not read as rejected; bounded by the 5s deadline |
