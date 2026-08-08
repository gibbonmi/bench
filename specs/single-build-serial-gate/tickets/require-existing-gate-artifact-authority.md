# Require existing gate artifact authority

Blocked by: share-artifacts-across-local-processes.md, resolve-target-aware-artifact-identities.md
Ownership fence: `internal/artifactstore/`, `internal/freshness/`
Integration surfaces: validated stored-record loader in `internal/artifactstore/store.go`→strict existing-record API in `internal/artifactstore/`; strict existing-record API→freshness-owned CLI/verifier authority in `internal/freshness/`; selected immutable CLI/verifier records→prepare-gate-artifacts-before-scheduling.md
Contracts: `ArtifactIdentity` and validated `ArtifactRecord` cross `internal/artifactstore/`→`internal/freshness/`, membership is one current CLI request plus one current freshness-verifier request, ordering is derive current identities then require both records then authenticate verifier and root-local CLI bytes before returning either executable, and absent, stale, malformed, or corrupt state refuses without builder, lock, record, artifact, or authority writes, asserted by GA1-GA4 against the real store and freshness owner
Closure: GA1/missing-record-refusal, GA1/valid-record-reuse, GA1/no-builder-fallback, GA1/no-store-write, GA2/corrupt-record-refusal, GA2/missing-artifact-refusal, GA2/symlinked-artifact-refusal, GA2/malformed-record-refusal, GA2/digest-mismatch-refusal, GA3/authority-atomic-publication, GA3/authority-absence-refusal, GA3/authority-field-set-refusal, GA3/cli-current-identity, GA3/verifier-current-identity, GA3/verifier-self-digest, GA3/root-cli-digest, GA4/single-record-validator

## What to build

Add the non-authoring half of the artifact store and compose it into a freshness-owned gate authority for the current CLI and freshness-verifier requests. Authority publication occurs only after both real store records validate; authority loading is strict and read-only, recomputes current identities, authenticates the running verifier and root-local CLI bytes, and returns immutable records without any construction fallback. An artifactstore-only split leaves `/tmp/gp1-precli-verifier-repro.sh` red because no pre-CLI owner consumes the strict lookup, while a freshness-only split must duplicate the store's private validator or retain authoring fallback; the debug receipt therefore requires both directories in this one repair fence. This ticket provides the independently-green seam and hostile-state coverage; `prepare-gate-artifacts-before-scheduling.md` wires its producer and gate-entry consumers.

## Acceptance

- [ ] [GA1] (covers local) strict existing-record lookup reuses a valid durable record, while an absent identity refuses without invoking a builder or writing store state.
- [ ] [GA2] (covers local) strict lookup applies the store's existing complete-record validator to every corrupt, missing-artifact, symlinked-artifact, malformed-field-set, and digest-mismatch state and refuses without repair or reconstruction.
- [ ] [GA3] (covers local) freshness publishes one exact-field-set CLI/verifier authority only after both records validate, and strict loading refuses absent, malformed, stale-identity, verifier-self-digest, or root-CLI-digest state before returning an executable or writing authority/store state.
- [ ] [GA4] (covers local) freshness consumes the artifact store's validated record interface and does not duplicate its durable record schema or parser.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GA1/missing-record-refusal | call the authoring resolver when the record is absent | strict existing-record store test | require a fresh identity with a counting builder sentinel and expect missing-record refusal with zero builder calls |
| GA1/valid-record-reuse | refuse every lookup without loading a retained success | strict existing-record store test | author one valid record through the real resolver, reopen the store, require it, and expect identical identity, digest, and immutable path |
| GA1/no-builder-fallback | inject a failing builder as the strict lookup fallback | strict existing-record store test | require a fresh identity and expect no builder marker and no sticky failure |
| GA1/no-store-write | create a lock, record, artifact, or directory during missing lookup | read-only store snapshot test | snapshot an absent injected root, require one identity, and expect byte-identical filesystem state |
| GA2/corrupt-record-refusal | accept invalid record bytes | hostile strict-record table test | plant corrupt record bytes, require the identity, and expect validation refusal with no builder marker or state change |
| GA2/missing-artifact-refusal | accept a completed record whose artifact is absent | hostile strict-record table test | remove the recorded artifact, require the identity, and expect refusal with no reconstruction |
| GA2/symlinked-artifact-refusal | accept a symlink at the immutable artifact path | hostile strict-record table test | replace the recorded artifact with a symlink, require the identity, and expect refusal before execution |
| GA2/malformed-record-refusal | ignore an unknown, missing, or contradictory record field | hostile strict-record table test | mutate the completed record field set, require the identity, and expect exact-schema refusal |
| GA2/digest-mismatch-refusal | skip comparison of immutable bytes to the recorded digest | hostile strict-record table test | replace immutable executable bytes, require the identity, and expect digest refusal |
| GA3/authority-atomic-publication | publish authority before both records validate | authority publication fault test | stop between CLI and verifier validation and expect no readable authority record |
| GA3/authority-absence-refusal | synthesize authority from ambient paths when its record is absent | strict authority load test | load from an absent authority path with executable markers planted and expect refusal with no marker or write |
| GA3/authority-field-set-refusal | tolerate an unknown or missing authority field | strict authority field-set table test | mutate each authority field class, load it, and expect exact-schema refusal before store access |
| GA3/cli-current-identity | trust the recorded CLI identity without deriving the current request | CLI request mutation test | change selected CLI source after publication, load authority, and expect stale-identity refusal |
| GA3/verifier-current-identity | trust the recorded verifier identity without deriving the current request | verifier request mutation test | change selected verifier source after publication, load authority, and expect stale-identity refusal |
| GA3/verifier-self-digest | execute a verifier whose bytes differ from its selected record | verifier self-authentication test | replace the verifier executable supplied to strict loading and expect refusal before CLI selection |
| GA3/root-cli-digest | return the selected CLI when root-local `dist/bench` differs | root CLI authentication test | replace root-local CLI bytes after authority publication and expect refusal before either CLI executes |
| GA4/single-record-validator | decode the durable store record again inside freshness | record-validator ownership audit | add a freshness-local durable record struct or JSON decoder, run the package audit, and expect the duplicate-validator diagnostic |
