# Share artifacts across local processes

Blocked by: none
Ownership fence: `internal/artifactstore/`, `internal/freshness/`
Integration surfaces: opaque identity API→resolve-target-aware-artifact-identities.md; durable store API→rendezvous-artifact-consumers-across-processes.md; artifact record API→prepare-gate-artifacts-before-scheduling.md; artifact record API→migrate-gate-fixture-artifact-consumers.md; artifact record API→migrate-contract-and-preflight-artifact-consumers.md; proof backend API→register-real-build-proof-identities.md; temporary legacy coexistence→enforce-the-executable-artifact-contract.md; cross-process coordinator→serialize-gate-lineages-and-transfer-turns.md
Contracts: `ArtifactIdentity` crosses callers→`internal/artifactstore/` as an opaque exact digest, ordering is check-lock-recheck-author-publish, absence refuses without fallback, asserted by AS1 and AS3 against the real store; `ArtifactRecord` (identity, state, executable digest, immutable path, error class) crosses `internal/artifactstore/`→dependent adapters, success is durable and failure is repair-only
Closure: AS1/cross-process-singleflight, AS1/durable-reuse, AS3/sticky-failure, AS3/same-error-fanout, AS3/cancellation-no-success, AS3/atomic-success, AS3/corrupt-record-refusal, AS3/missing-artifact-refusal, AS3/symlinked-artifact-refusal, AS3/malformed-record-refusal, AS3/digest-mismatch-refusal, AS3/active-owner-repair-refusal, AS3/cross-identity-repair-refusal, AS3/isolated-test-store

## What to build

Add the repository-local content-addressed artifact store beside current builders and prove its singleflight, durable publication, failure, hostile-record, and repair semantics through real sibling processes. The legacy construction paths remain temporarily available for migration tickets, but no new consumer may join them.

## Acceptance

- [ ] [AS1] (covers SB1) repeated, concurrent, and later processes resolving one exact identity observe one backend construction and one durable artifact digest.
- [ ] [AS3] (covers SB3) one failed attempt returns the same error to every waiter without retry, cancellation publishes no success, hostile completed records refuse, exact repair cannot cross identity or an active owner, and every fault test uses an injected isolated store.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AS1/cross-process-singleflight | remove the under-lock completed-record recheck | sibling-process counting-builder test | start two bounded re-exec resolvers at one barrier, release them, and expect the backend-count-equals-one failure |
| AS1/durable-reuse | key completed records by invocation | second-process durable-reuse test | resolve once, exit the producer, resolve from a fresh process, and expect a second-build failure |
| AS3/sticky-failure | retry automatically after the injected builder error | concurrent failure waiter test | release two waiters onto one failing attempt and expect the second-attempt failure |
| AS3/same-error-fanout | synthesize a waiter-local error instead of returning the attempt result | concurrent failure waiter test | release two waiters onto one failing attempt and expect byte-identical error class and detail from both |
| AS3/cancellation-no-success | publish a completed record when the authoring context is cancelled | bounded cancellation publication test | cancel the held author, bound every waiter, and expect no completed record or executable plus the same cancellation result for every waiter |
| AS3/atomic-success | write the completed record before the executable rename | interrupted-publication test | stop construction between record and artifact publication and expect the accepted-partial failure |
| AS3/corrupt-record-refusal | accept a corrupt completed record | hostile-record table test | plant invalid record bytes and expect validation refusal before execution |
| AS3/missing-artifact-refusal | accept a record whose executable is absent | hostile-record table test | remove the recorded executable and expect validation refusal before rebuilding or execution |
| AS3/symlinked-artifact-refusal | accept a symlinked completed artifact | hostile-record table test | replace the executable with a symlink and expect validation refusal before execution |
| AS3/malformed-record-refusal | accept a completed record with an invalid field set | hostile-record table test | plant the malformed record and expect validation refusal before execution |
| AS3/digest-mismatch-refusal | accept executable bytes that do not match the record | hostile-record table test | replace the executable bytes and expect digest refusal before execution |
| AS3/active-owner-repair-refusal | omit the active-owner check from repair | exact-repair scope test | hold identity A active, request repair of A, and expect refusal before record removal |
| AS3/cross-identity-repair-refusal | broaden repair from one identity to its neighboring record | exact-repair scope test | request repair of A beside B and expect B's record and artifact to remain byte-identical |
| AS3/isolated-test-store | fall back from the injected test store to common-Git-dir discovery | store-isolation sentinel test | plant a sentinel in the live-store double, run fault cases with an injected root, and expect the sentinel-access failure |
