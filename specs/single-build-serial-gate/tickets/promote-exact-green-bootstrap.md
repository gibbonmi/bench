# Promote an exact-green bootstrap

Blocked by: require-existing-gate-artifact-authority.md
Ownership fence: `internal/gate/`, `internal/freshness/`
Integration surfaces: exact-green gate evidence→`internal/gate/`; durable bootstrap publication→`internal/freshness/`; promoted bootstrap producer→run-promoted-gate-bootstrap.md
Contracts: `GateBootstrapRecord` crosses exact-green evidence in `internal/gate/`→durable publication in `internal/freshness/`, membership is one schema version, exact subject, host target, executable digest, and distinct bootstrap slot, ordering is exact-green validation then executable authentication then atomic executable-and-record publication, and absent, stale, red, partial, malformed, hostile, or interrupted state refuses without released-cache, artifact-store, or gate-evidence mutation, asserted by PG1-PG3
Closure: PG1/missing-green, PG1/stale-green, PG1/red-green, PG1/partial-green, PG1/executable-mismatch, PG2/released-cache-isolation, PG2/subject-binding, PG2/target-binding, PG2/digest-binding, PG2/malformed-record-refusal, PG3/symlink-destination, PG3/special-destination, PG3/interrupt-retains-prior, PG3/interrupt-cleans-stage

## What to build

Let the legacy source-verifier gate promote its already-authenticated native executable into a distinct bootstrap slot only after an exact fresh-green verdict. The producer validates the gate-owned evidence, authenticates the executable selected by that green subject, and publishes the executable plus strict record durably without changing released cache bytes or publishing externally. The dependent ticket is the first production caller and owns installed-only execution; until then the ordinary gate remains unchanged.

## Acceptance

- [ ] [PG1] (covers local) bootstrap promotion accepts only the exact fresh-green subject and its authenticated executable; absent, stale, red, partial, or executable-mismatched evidence refuses without writes.
- [ ] [PG2] (covers local) the exact-schema bootstrap record binds subject, host target, and executable digest in a namespace disjoint from the released version-target binary cache, and malformed records refuse.
- [ ] [PG3] (covers local) bootstrap executable and record publication is atomic, refuses symlink or special-file destinations, and interruption independently retains the prior valid bootstrap and removes staging residue.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PG1/missing-green | promote without gate evidence | bootstrap promotion process test | invoke promotion in a repository with no verdict and expect refusal plus a byte-identical bootstrap namespace |
| PG1/stale-green | accept evidence for an earlier tree | bootstrap promotion process test | seed a green verdict, move the tree, invoke promotion, and expect stale-subject refusal with no slot write |
| PG1/red-green | accept a red verdict | bootstrap promotion process test | seed exact red evidence, invoke promotion, and expect red-evidence refusal with no slot write |
| PG1/partial-green | accept reduced or inherited evidence | bootstrap promotion process test | seed non-fresh or partial evidence for the exact tree, invoke promotion, and expect exact-fresh refusal |
| PG1/executable-mismatch | omit comparison of executable bytes to the green subject's authenticated CLI | executable substitution test | seed exact green, supply different executable bytes, and expect digest refusal before staging |
| PG2/released-cache-isolation | publish into the ordinary version-target cache | cache namespace fixture | promote one bootstrap and expect the released-cache sentinel to remain byte-identical |
| PG2/subject-binding | omit the exact subject from the bootstrap record | strict bootstrap record test | change the selected green subject while retaining target and bytes and expect record refusal |
| PG2/target-binding | omit the host target from the bootstrap record | strict bootstrap record test | transplant a valid record under another target and expect refusal |
| PG2/digest-binding | omit executable digest from the bootstrap record | strict bootstrap record test | change executable bytes while retaining subject and target and expect record refusal |
| PG2/malformed-record-refusal | tolerate an unknown, missing, or contradictory record field | strict bootstrap record table test | mutate each field class independently and expect exact-schema refusal before publication |
| PG3/symlink-destination | replace the bootstrap destination or ancestor with a symlink | hostile publication test | promote into the hostile tree and expect refusal without touching the referent |
| PG3/special-destination | replace the bootstrap record or executable with a FIFO or other special file | bounded hostile publication test | invoke promotion under a one-second bound and expect special-file refusal without blocking |
| PG3/interrupt-retains-prior | replace the prior record or executable when interrupted | re-exec signal test | interrupt after staging and expect the prior bootstrap pair still valid |
| PG3/interrupt-cleans-stage | leave staging files when interrupted | re-exec signal test | interrupt after staging and expect no bootstrap staging residue |
